package blockchain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"agent-arena/backend/internal/blockchain/bindings"
)

// DefaultUSDCAddress MockUSDC 在 Sepolia 上的部署地址（6 位小数）
const DefaultUSDCAddress = "0x9163ad7caf7bf73ff105c658a42e12eaa33f58af"

// USDCService MockUSDC (ERC20, 6 decimals) 只读服务
type USDCService struct {
	client    *ethclient.Client
	token     *bindings.MockUSDCBindings
	tokenAddr common.Address
}

// NewUSDCService 创建 USDC 只读服务
func NewUSDCService(rpcURL, tokenAddrHex string) (*USDCService, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("dial RPC: %w", err)
	}

	if tokenAddrHex == "" {
		tokenAddrHex = DefaultUSDCAddress
	}
	tokenAddr := common.HexToAddress(tokenAddrHex)

	token, err := bindings.NewMockUSDCBindings(tokenAddr, client)
	if err != nil {
		return nil, fmt.Errorf("bind USDC: %w", err)
	}

	return &USDCService{
		client:    client,
		token:     token,
		tokenAddr: tokenAddr,
	}, nil
}

// TokenAddress USDC 合约地址
func (s *USDCService) TokenAddress() common.Address {
	return s.tokenAddr
}

// BalanceOf 读取链上 USDC 余额（返回 USDC 数量，已除以 1e6）
func (s *USDCService) BalanceOf(ctx context.Context, addr common.Address) (uint64, error) {
	bal, err := s.token.BalanceOf(&bind.CallOpts{Context: ctx}, addr)
	if err != nil {
		return 0, fmt.Errorf("balanceOf: %w", err)
	}
	// USDC 是 6 位小数
	usdc := new(big.Int).Div(bal, big.NewInt(1e6))
	if !usdc.IsUint64() {
		return 0, fmt.Errorf("balance overflow uint64")
	}
	return usdc.Uint64(), nil
}

// BalanceOfRaw 读取链上 USDC 余额（最小单位）
func (s *USDCService) BalanceOfRaw(ctx context.Context, addr common.Address) (*big.Int, error) {
	return s.token.BalanceOf(&bind.CallOpts{Context: ctx}, addr)
}
