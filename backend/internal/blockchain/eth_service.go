package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"agent-arena/backend/internal/blockchain/bindings"
	"agent-arena/backend/internal/engine"
)

// EthChainService 真实链上交互（Sepolia）
type EthChainService struct {
	client     *ethclient.Client
	arena      *bindings.AgentArenaBindings
	privateKey *ecdsa.PrivateKey
	fromAddr   common.Address
	arenaAddr  common.Address
	chainID    *big.Int

	mu       sync.RWMutex
	agentMap map[string]common.Address // agent name → address
	events   chan ContractEvent
}

// NewEthChainService 创建链上服务
func NewEthChainService(rpcURL, privateKeyHex, arenaAddrHex string) (*EthChainService, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("dial RPC: %w", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get chain ID: %w", err)
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyHex, "0x"))
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	fromAddr := crypto.PubkeyToAddress(privateKey.PublicKey)

	arenaAddr := common.HexToAddress(arenaAddrHex)
	arena, err := bindings.NewAgentArenaBindings(arenaAddr, client)
	if err != nil {
		return nil, fmt.Errorf("bind arena contract: %w", err)
	}

	return &EthChainService{
		client:     client,
		arena:      arena,
		privateKey: privateKey,
		fromAddr:   fromAddr,
		arenaAddr:  arenaAddr,
		chainID:    chainID,
		agentMap:   make(map[string]common.Address),
		events:     make(chan ContractEvent, 100),
	}, nil
}

// getTransactor 创建带 nonce/gas 的 transactor
func (s *EthChainService) getTransactor(ctx context.Context) (*bind.TransactOpts, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(s.privateKey, s.chainID)
	if err != nil {
		return nil, fmt.Errorf("create transactor: %w", err)
	}

	// 使用 latest block 而非 pending（Infura 免费层兼容性更好）
	nonce, err := s.client.NonceAt(ctx, s.fromAddr, nil)
	if err != nil {
		return nil, fmt.Errorf("get nonce: %w", err)
	}

	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("suggest gas price: %w", err)
	}

	// 提高 gas price 20% 确保交易快速确认
	gasPrice = new(big.Int).Mul(gasPrice, big.NewInt(120))
	gasPrice = new(big.Int).Div(gasPrice, big.NewInt(100))

	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = gasPrice
	auth.Context = ctx
	auth.GasLimit = 300000

	return auth, nil
}

// resolveAgent 将 agent 名称解析为地址
func (s *EthChainService) resolveAgent(name string) common.Address {
	s.mu.RLock()
	addr, ok := s.agentMap[strings.ToLower(name)]
	s.mu.RUnlock()
	if ok {
		return addr
	}
	// 默认返回 deployer 地址
	return s.fromAddr
}

// RegisterAgentAddress 注册 agent 名称到地址的映射
func (s *EthChainService) RegisterAgentAddress(name string, addr common.Address) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.agentMap[strings.ToLower(name)] = addr
}

// === Reads ===

func (s *EthChainService) GetStrategyWeights(ctx context.Context, gameID uint64) (*engine.StrategyWeights, error) {
	result, err := s.arena.GetStrategyWeights(&bind.CallOpts{Context: ctx}, big.NewInt(int64(gameID)))
	if err != nil {
		return nil, fmt.Errorf("get strategy weights: %w", err)
	}

	return &engine.StrategyWeights{
		Aggressive: uint8(result.Aggressive.Uint64()),
		Defensive:  uint8(result.Defensive.Uint64()),
		Tricky:     uint8(result.Tricky.Uint64()),
	}, nil
}

func (s *EthChainService) GetGamePool(ctx context.Context, gameID uint64) (*big.Int, *big.Int, error) {
	info, err := s.arena.GetGame(&bind.CallOpts{Context: ctx}, big.NewInt(int64(gameID)))
	if err != nil {
		return nil, nil, fmt.Errorf("get game: %w", err)
	}
	return new(big.Int).Set(info.TotalBetRed), new(big.Int).Set(info.TotalBetBlue), nil
}

func (s *EthChainService) GetOdds(ctx context.Context, gameID uint64) (*big.Int, *big.Int, error) {
	result, err := s.arena.GetOdds(&bind.CallOpts{Context: ctx}, big.NewInt(int64(gameID)))
	if err != nil {
		return nil, nil, fmt.Errorf("get odds: %w", err)
	}
	return result.OddsRed, result.OddsBlue, nil
}

func (s *EthChainService) GetGameStatus(ctx context.Context, gameID uint64) (string, error) {
	info, err := s.arena.GetGame(&bind.CallOpts{Context: ctx}, big.NewInt(int64(gameID)))
	if err != nil {
		return "", fmt.Errorf("get game: %w", err)
	}

	// status enum: 0=Open, 1=Locked, 2=Finished
	switch info.Status {
	case 0:
		return "open", nil
	case 1:
		return "locked", nil
	case 2:
		return "finished", nil
	default:
		return fmt.Sprintf("unknown(%d)", info.Status), nil
	}
}

// === Writes ===

func (s *EthChainService) CreateGame(ctx context.Context, agentRed, agentBlue string, bettingDuration uint64) (string, uint64, error) {
	auth, err := s.getTransactor(ctx)
	if err != nil {
		return "", 0, err
	}

	addrRed := s.resolveAgent(agentRed)
	addrBlue := s.resolveAgent(agentBlue)

	tx, err := s.arena.CreateGame(auth, addrRed, addrBlue)
	if err != nil {
		return "", 0, fmt.Errorf("create game tx: %w", err)
	}

	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return tx.Hash().Hex(), 0, fmt.Errorf("wait mined: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return tx.Hash().Hex(), 0, fmt.Errorf("create game tx failed: status %d", receipt.Status)
	}

	// 从事件日志中提取 gameID
	gameID, err := s.extractGameIDFromLogs(receipt)
	if err != nil {
		return tx.Hash().Hex(), 0, fmt.Errorf("extract game ID: %w", err)
	}

	log.Printf("CreateGame tx: %s, gameID: %d", tx.Hash().Hex(), gameID)
	return tx.Hash().Hex(), gameID, nil
}

func (s *EthChainService) StartGame(ctx context.Context, gameID uint64) (string, error) {
	auth, err := s.getTransactor(ctx)
	if err != nil {
		return "", err
	}

	tx, err := s.arena.StartGame(auth, big.NewInt(int64(gameID)))
	if err != nil {
		return "", fmt.Errorf("start game tx: %w", err)
	}

	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return tx.Hash().Hex(), fmt.Errorf("wait mined: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return tx.Hash().Hex(), fmt.Errorf("start game tx failed: status %d", receipt.Status)
	}

	log.Printf("StartGame tx: %s, gameID: %d", tx.Hash().Hex(), gameID)
	return tx.Hash().Hex(), nil
}

func (s *EthChainService) FinishGame(ctx context.Context, gameID uint64, redWins bool, actionsHash [32]byte) (string, error) {
	auth, err := s.getTransactor(ctx)
	if err != nil {
		return "", err
	}

	tx, err := s.arena.FinishGame(auth, big.NewInt(int64(gameID)), redWins, actionsHash)
	if err != nil {
		return "", fmt.Errorf("finish game tx: %w", err)
	}

	receipt, err := bind.WaitMined(ctx, s.client, tx)
	if err != nil {
		return tx.Hash().Hex(), fmt.Errorf("wait mined: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return tx.Hash().Hex(), fmt.Errorf("finish game tx failed: status %d", receipt.Status)
	}

	log.Printf("FinishGame tx: %s, gameID: %d, redWins: %v", tx.Hash().Hex(), gameID, redWins)
	return tx.Hash().Hex(), nil
}

// === Events ===

func (s *EthChainService) ListenEvents(ctx context.Context) (<-chan ContractEvent, error) {
	go s.pollEvents(ctx)
	return s.events, nil
}

// pollEvents 轮询合约事件（每 15 秒查一次新区块）
func (s *EthChainService) pollEvents(ctx context.Context) {
	arenaABI, err := abi.JSON(strings.NewReader(bindings.AgentArenaBindingsMetaData.ABI))
	if err != nil {
		log.Printf("ERROR: parse arena ABI for events: %v", err)
		return
	}

	lastBlock, err := s.client.BlockNumber(ctx)
	if err != nil {
		log.Printf("ERROR: get current block: %v", err)
		return
	}

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			currentBlock, err := s.client.BlockNumber(ctx)
			if err != nil {
				log.Printf("ERROR: poll block number: %v", err)
				continue
			}
			if currentBlock <= lastBlock {
				continue
			}

			fromBlock := lastBlock + 1
			toBlock := currentBlock

			query := ethereum.FilterQuery{
				FromBlock: big.NewInt(int64(fromBlock)),
				ToBlock:   big.NewInt(int64(toBlock)),
				Addresses: []common.Address{s.arenaAddr},
			}

			logs, err := s.client.FilterLogs(ctx, query)
			if err != nil {
				log.Printf("ERROR: filter logs: %v", err)
				continue
			}

			for _, vLog := range logs {
				event, ok := s.parseEvent(arenaABI, vLog)
				if ok {
					select {
					case s.events <- event:
					default:
						log.Printf("WARN: event channel full, dropping event")
					}
				}
			}

			lastBlock = currentBlock
		}
	}
}

// parseEvent 解析单条日志为 ContractEvent
func (s *EthChainService) parseEvent(contractABI abi.ABI, vLog types.Log) (ContractEvent, bool) {
	if len(vLog.Topics) == 0 {
		return ContractEvent{}, false
	}

	eventSig := vLog.Topics[0]

	switch {
	case eventSig == crypto.Keccak256Hash([]byte("GameInitialized(uint256,address,address)")):
		var event struct {
			AgentRed  common.Address
			AgentBlue common.Address
		}
		if err := contractABI.UnpackIntoInterface(&event, "GameInitialized", vLog.Data); err != nil {
			return ContractEvent{}, false
		}
		gameID := uint64(0)
		if len(vLog.Topics) > 1 {
			gameID = new(big.Int).SetBytes(vLog.Topics[1].Bytes()).Uint64()
		}
		return ContractEvent{
			Type:   EventGameCreated,
			GameID: gameID,
			Block:  vLog.BlockNumber,
			TxHash: vLog.TxHash.Hex(),
			Data: GameCreatedData{
				AgentRed:  event.AgentRed.Hex(),
				AgentBlue: event.AgentBlue.Hex(),
			},
		}, true

	case eventSig == crypto.Keccak256Hash([]byte("BetPlaced(uint256,address,bool,uint256)")):
		var event struct {
			User   common.Address
			Side   bool
			Amount *big.Int
		}
		if err := contractABI.UnpackIntoInterface(&event, "BetPlaced", vLog.Data); err != nil {
			return ContractEvent{}, false
		}
		gameID := uint64(0)
		if len(vLog.Topics) > 1 {
			gameID = new(big.Int).SetBytes(vLog.Topics[1].Bytes()).Uint64()
		}
		return ContractEvent{
			Type:   EventBetPlaced,
			GameID: gameID,
			Block:  vLog.BlockNumber,
			TxHash: vLog.TxHash.Hex(),
			Data: BetPlacedData{
				User:   event.User.Hex(),
				Side:   event.Side,
				Amount: event.Amount,
			},
		}, true

	case eventSig == crypto.Keccak256Hash([]byte("BettingLocked(uint256,uint256,uint256)")):
		var event struct {
			TotalBetRed  *big.Int
			TotalBetBlue *big.Int
		}
		if err := contractABI.UnpackIntoInterface(&event, "BettingLocked", vLog.Data); err != nil {
			return ContractEvent{}, false
		}
		gameID := uint64(0)
		if len(vLog.Topics) > 1 {
			gameID = new(big.Int).SetBytes(vLog.Topics[1].Bytes()).Uint64()
		}
		return ContractEvent{
			Type:   EventBettingLocked,
			GameID: gameID,
			Block:  vLog.BlockNumber,
			TxHash: vLog.TxHash.Hex(),
			Data: BettingLockedData{
				TotalBetRed:  event.TotalBetRed,
				TotalBetBlue: event.TotalBetBlue,
			},
		}, true

	case eventSig == crypto.Keccak256Hash([]byte("GameStarted(uint256)")):
		gameID := uint64(0)
		if len(vLog.Topics) > 1 {
			gameID = new(big.Int).SetBytes(vLog.Topics[1].Bytes()).Uint64()
		}
		return ContractEvent{
			Type:   EventBettingLocked,
			GameID: gameID,
			Block:  vLog.BlockNumber,
			TxHash: vLog.TxHash.Hex(),
		}, true

	case eventSig == crypto.Keccak256Hash([]byte("GameFinished(uint256,bool)")):
		var event struct {
			RedWins bool
		}
		if err := contractABI.UnpackIntoInterface(&event, "GameFinished", vLog.Data); err != nil {
			return ContractEvent{}, false
		}
		gameID := uint64(0)
		if len(vLog.Topics) > 1 {
			gameID = new(big.Int).SetBytes(vLog.Topics[1].Bytes()).Uint64()
		}
		return ContractEvent{
			Type:   EventGameSettled,
			GameID: gameID,
			Block:  vLog.BlockNumber,
			TxHash: vLog.TxHash.Hex(),
			Data: GameSettledData{
				RedWins: event.RedWins,
			},
		}, true

	case eventSig == crypto.Keccak256Hash([]byte("GameSettled(uint256,bool,uint256,uint256)")):
		var event struct {
			RedWins     bool
			TotalPool   *big.Int
			ProtocolFee *big.Int
		}
		if err := contractABI.UnpackIntoInterface(&event, "GameSettled", vLog.Data); err != nil {
			return ContractEvent{}, false
		}
		gameID := uint64(0)
		if len(vLog.Topics) > 1 {
			gameID = new(big.Int).SetBytes(vLog.Topics[1].Bytes()).Uint64()
		}
		return ContractEvent{
			Type:   EventGameSettled,
			GameID: gameID,
			Block:  vLog.BlockNumber,
			TxHash: vLog.TxHash.Hex(),
			Data: GameSettledData{
				RedWins:     event.RedWins,
				TotalPool:   event.TotalPool,
				ProtocolFee: event.ProtocolFee,
			},
		}, true

	case eventSig == crypto.Keccak256Hash([]byte("VotesLocked(uint256,uint256,uint256,uint256)")):
		var event struct {
			Aggressive *big.Int
			Defensive  *big.Int
			Tricky     *big.Int
		}
		if err := contractABI.UnpackIntoInterface(&event, "VotesLocked", vLog.Data); err != nil {
			return ContractEvent{}, false
		}
		gameID := uint64(0)
		if len(vLog.Topics) > 1 {
			gameID = new(big.Int).SetBytes(vLog.Topics[1].Bytes()).Uint64()
		}
		return ContractEvent{
			Type:   EventVotesLocked,
			GameID: gameID,
			Block:  vLog.BlockNumber,
			TxHash: vLog.TxHash.Hex(),
			Data: VotesLockedData{
				Aggressive: event.Aggressive.Uint64(),
				Defensive:  event.Defensive.Uint64(),
				Tricky:     event.Tricky.Uint64(),
			},
		}, true

	case eventSig == crypto.Keccak256Hash([]byte("RewardClaimed(uint256,address,uint256)")):
		var event struct {
			User   common.Address
			Reward *big.Int
		}
		if err := contractABI.UnpackIntoInterface(&event, "RewardClaimed", vLog.Data); err != nil {
			return ContractEvent{}, false
		}
		gameID := uint64(0)
		if len(vLog.Topics) > 1 {
			gameID = new(big.Int).SetBytes(vLog.Topics[1].Bytes()).Uint64()
		}
		return ContractEvent{
			Type:   EventRewardClaimed,
			GameID: gameID,
			Block:  vLog.BlockNumber,
			TxHash: vLog.TxHash.Hex(),
			Data: RewardClaimedData{
				User:   event.User.Hex(),
				Reward: event.Reward,
			},
		}, true

	case eventSig == crypto.Keccak256Hash([]byte("GameResultSubmitted(uint256,bool,bytes32)")):
		var event struct {
			RedWins     bool
			ActionsHash [32]byte
		}
		if err := contractABI.UnpackIntoInterface(&event, "GameResultSubmitted", vLog.Data); err != nil {
			return ContractEvent{}, false
		}
		gameID := uint64(0)
		if len(vLog.Topics) > 1 {
			gameID = new(big.Int).SetBytes(vLog.Topics[1].Bytes()).Uint64()
		}
		return ContractEvent{
			Type:   EventGameResultSubmitted,
			GameID: gameID,
			Block:  vLog.BlockNumber,
			TxHash: vLog.TxHash.Hex(),
			Data: GameResultSubmittedData{
				RedWins:     event.RedWins,
				ActionsHash: event.ActionsHash,
			},
		}, true
	}

	return ContractEvent{}, false
}

// extractGameIDFromLogs 从交易回执中提取 gameID
func (s *EthChainService) extractGameIDFromLogs(receipt *types.Receipt) (uint64, error) {
	gameInitializedSig := crypto.Keccak256Hash([]byte("GameInitialized(uint256,address,address)"))

	for _, vLog := range receipt.Logs {
		if len(vLog.Topics) > 0 && vLog.Topics[0] == gameInitializedSig {
			if len(vLog.Topics) > 1 {
				gameID := new(big.Int).SetBytes(vLog.Topics[1].Bytes())
				return gameID.Uint64(), nil
			}
		}
	}

	return 0, fmt.Errorf("GameInitialized event not found in tx receipt")
}

// Client 返回 ethclient
func (s *EthChainService) Client() *ethclient.Client {
	return s.client
}

// ArenaAddress 返回 Arena 合约地址
func (s *EthChainService) ArenaAddress() common.Address {
	return s.arenaAddr
}

// TransactorAddress 返回交易发送者地址
func (s *EthChainService) TransactorAddress() common.Address {
	return s.fromAddr
}
