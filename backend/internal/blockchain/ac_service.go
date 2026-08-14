package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"agent-arena/backend/internal/blockchain/bindings"
)

// ArenaCoin 在 Sepolia 上的部署地址
const DefaultArenaCoinAddress = "0xb9a2aaf070EA64d978a41bd5E9fF89756dA4354C"

// DailyClaimAmount 每日领取数量（100 AC，18 位小数）
var DailyClaimAmount = new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18))

// ACService ArenaCoin (ERC20) 链上交互服务
type ACService struct {
	client     *ethclient.Client
	token      *bindings.ArenaCoinBindings
	privateKey *ecdsa.PrivateKey
	fromAddr   common.Address
	tokenAddr  common.Address
	chainID    *big.Int
}

// NewACService 创建 AC 服务
// rpcURL - Sepolia RPC
// deployerKeyHex - 合约部署者私钥（用于需要 deployer 签名的操作）
// ownerKeyHex - 合约 owner 私钥（用于 mint/transfer，作为 treasury）
//   如果为空则回退到 deployerKeyHex
// tokenAddrHex - ArenaCoin 合约地址（留空使用默认）
func NewACService(rpcURL, deployerKeyHex, ownerKeyHex, tokenAddrHex string) (*ACService, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("dial RPC: %w", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get chain ID: %w", err)
	}

	deployerKey, err := crypto.HexToECDSA(strings.TrimPrefix(deployerKeyHex, "0x"))
	if err != nil {
		return nil, fmt.Errorf("parse deployer private key: %w", err)
	}

	// treasury 使用 owner key（独立于 deployer）
	treasuryKey := deployerKey
	if ownerKeyHex != "" {
		treasuryKey, err = crypto.HexToECDSA(strings.TrimPrefix(ownerKeyHex, "0x"))
		if err != nil {
			return nil, fmt.Errorf("parse owner private key: %w", err)
		}
	}
	treasuryAddr := crypto.PubkeyToAddress(treasuryKey.PublicKey)

	if tokenAddrHex == "" {
		tokenAddrHex = DefaultArenaCoinAddress
	}
	tokenAddr := common.HexToAddress(tokenAddrHex)

	token, err := bindings.NewArenaCoinBindings(tokenAddr, client)
	if err != nil {
		return nil, fmt.Errorf("bind arena coin: %w", err)
	}

	return &ACService{
		client:     client,
		token:      token,
		privateKey: treasuryKey,
		fromAddr:   treasuryAddr,
		tokenAddr:  tokenAddr,
		chainID:    chainID,
	}, nil
}

// TokenAddress 合约地址
func (s *ACService) TokenAddress() common.Address {
	return s.tokenAddr
}

// TreasuryAddress 资金池地址（后端持有私钥的地址）
func (s *ACService) TreasuryAddress() common.Address {
	return s.fromAddr
}

// BalanceOf 读取链上 AC 余额（返回 AC 数量，已除以 1e18）
func (s *ACService) BalanceOf(ctx context.Context, addr common.Address) (uint64, error) {
	bal, err := s.token.BalanceOf(&bind.CallOpts{Context: ctx}, addr)
	if err != nil {
		return 0, fmt.Errorf("balanceOf: %w", err)
	}
	// 转换为 AC（除以 1e18）
	ac := new(big.Int).Div(bal, big.NewInt(1e18))
	if !ac.IsUint64() {
		return 0, fmt.Errorf("balance overflow uint64")
	}
	return ac.Uint64(), nil
}

// BalanceOfRaw 读取链上 AC 余额（wei 单位）
func (s *ACService) BalanceOfRaw(ctx context.Context, addr common.Address) (*big.Int, error) {
	return s.token.BalanceOf(&bind.CallOpts{Context: ctx}, addr)
}

// CanClaim 查询链上是否可以领取（合约内 24 小时冷却）
func (s *ACService) CanClaim(ctx context.Context, addr common.Address) (bool, error) {
	return s.token.CanClaim(&bind.CallOpts{Context: ctx}, addr)
}

// TimeUntilNextClaim 查询距离下次可领取的秒数
func (s *ACService) TimeUntilNextClaim(ctx context.Context, addr common.Address) (uint64, error) {
	t, err := s.token.TimeUntilNextClaim(&bind.CallOpts{Context: ctx}, addr)
	if err != nil {
		return 0, err
	}
	if !t.IsUint64() {
		return 0, fmt.Errorf("time overflow")
	}
	return t.Uint64(), nil
}

// ClaimDaily 调用合约 claimDaily()，为用户铸造 100 AC
// 要求后端地址为合约 owner（只有 owner 能 mint）
// 注意：此方法会直接修改链上状态
func (s *ACService) ClaimDaily(ctx context.Context, userAddr common.Address) (string, error) {
	// 检查是否可以领取
	canClaim, err := s.token.CanClaim(&bind.CallOpts{Context: ctx}, userAddr)
	if err != nil {
		return "", fmt.Errorf("check canClaim: %w", err)
	}
	if !canClaim {
		return "", fmt.Errorf("already claimed within 24 hours")
	}

	auth, err := s.getTransactor(ctx)
	if err != nil {
		return "", err
	}

	// 调用 claimDaily() —— 合约以 msg.sender 为受益人铸造
	// 注意：这里 msg.sender 是后端地址，不是 userAddr
	// 所以我们用 transfer 方式替代（见 TransferFromTreasury）
	_ = auth
	return "", fmt.Errorf("claimDaily mints to msg.sender (treasury), use TransferFromTreasury instead")
}

// TransferFromTreasury 从资金池向用户转账指定数量的 AC
// 这是实际使用的发放方式：后端作为资金池，向用户转账
func (s *ACService) TransferFromTreasury(ctx context.Context, to common.Address, amount *big.Int) (string, error) {
	// 检查资金池余额
	bal, err := s.token.BalanceOf(&bind.CallOpts{Context: ctx}, s.fromAddr)
	if err != nil {
		return "", fmt.Errorf("check treasury balance: %w", err)
	}
	if bal.Cmp(amount) < 0 {
		return "", fmt.Errorf("treasury insufficient: have %s, need %s", bal.String(), amount.String())
	}

	auth, err := s.getTransactor(ctx)
	if err != nil {
		return "", err
	}

	tx, err := s.token.Transfer(auth, to, amount)
	if err != nil {
		return "", fmt.Errorf("transfer: %w", err)
	}

	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return "", fmt.Errorf("wait mined: %w", err)
	}

	return receipt.TxHash.Hex(), nil
}

// TransferACFromTreasury 转账指定数量的 AC（以 AC 为单位，自动乘以 1e18）
func (s *ACService) TransferACFromTreasury(ctx context.Context, to common.Address, acAmount uint64) (string, error) {
	amount := new(big.Int).Mul(big.NewInt(int64(acAmount)), big.NewInt(1e18))
	return s.TransferFromTreasury(ctx, to, amount)
}

// MintToUser 通过 owner 权限向用户铸造 AC（备用）
func (s *ACService) MintToUser(ctx context.Context, to common.Address, acAmount uint64) (string, error) {
	amount := new(big.Int).Mul(big.NewInt(int64(acAmount)), big.NewInt(1e18))

	auth, err := s.getTransactor(ctx)
	if err != nil {
		return "", err
	}

	tx, err := s.token.Mint(auth, to, amount)
	if err != nil {
		return "", fmt.Errorf("mint: %w", err)
	}

	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return "", fmt.Errorf("wait mined: %w", err)
	}

	return receipt.TxHash.Hex(), nil
}

// TreasuryBalance 查询资金池余额（AC 单位）
func (s *ACService) TreasuryBalance(ctx context.Context) (uint64, error) {
	return s.BalanceOf(ctx, s.fromAddr)
}

// getTransactor 创建带 nonce/gas 的 transactor
func (s *ACService) getTransactor(ctx context.Context) (*bind.TransactOpts, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(s.privateKey, s.chainID)
	if err != nil {
		return nil, fmt.Errorf("create transactor: %w", err)
	}

	nonce, err := s.client.NonceAt(ctx, s.fromAddr, nil)
	if err != nil {
		return nil, fmt.Errorf("get nonce: %w", err)
	}

	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("suggest gas price: %w", err)
	}
	// 提高 20% 确保快速确认
	gasPrice = new(big.Int).Mul(gasPrice, big.NewInt(120))
	gasPrice = new(big.Int).Div(gasPrice, big.NewInt(100))

	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = gasPrice
	auth.Context = ctx
	auth.GasLimit = 120000 // ERC20 转账足够

	return auth, nil
}

// ACServiceOrNil 返回 AC 服务（如果未配置则返回 nil）
// 调用方需要检查是否为 nil，nil 时降级到 DB-only 模式
type ACServiceOrNil struct {
	Service *ACService
}

// IsOnChain 是否启用了链上 AC
func (n *ACServiceOrNil) IsOnChain() bool {
	return n.Service != nil
}

// WaitForReceipt 带超时等待交易
func WaitForReceipt(ctx context.Context, client *ethclient.Client, txHash string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	hash := common.HexToHash(txHash)
	for {
		_, isPending, err := client.TransactionByHash(ctx, hash)
		if err != nil {
			return err
		}
		if !isPending {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}
}
