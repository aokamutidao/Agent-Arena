package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"agent-arena/backend/internal/engine"
)

// MockChainService 模拟链上交互（用于开发和测试）
type MockChainService struct {
	mu            sync.RWMutex
	games         map[uint64]*mockGame
	strategies    map[uint64]*engine.StrategyWeights
	nextGameID    uint64
	events        chan ContractEvent
}

type mockGame struct {
	Status       string // open, locked, finished
	TotalBetRed  *big.Int
	TotalBetBlue *big.Int
	RedWins      bool
	ActionsHash  [32]byte
}

// NewMockChainService 创建模拟链服务
func NewMockChainService() *MockChainService {
	return &MockChainService{
		games:      make(map[uint64]*mockGame),
		strategies: make(map[uint64]*engine.StrategyWeights),
		nextGameID: 1,
		events:     make(chan ContractEvent, 100),
	}
}

// GetStrategyWeights 获取策略权重
func (m *MockChainService) GetStrategyWeights(ctx context.Context, gameID uint64) (*engine.StrategyWeights, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	w, ok := m.strategies[gameID]
	if !ok {
		// 默认均分
		return &engine.StrategyWeights{Aggressive: 33, Defensive: 33, Tricky: 34}, nil
	}
	return w, nil
}

// GetGamePool 获取下注池
func (m *MockChainService) GetGamePool(ctx context.Context, gameID uint64) (*big.Int, *big.Int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	game, ok := m.games[gameID]
	if !ok {
		return big.NewInt(0), big.NewInt(0), nil
	}
	return new(big.Int).Set(game.TotalBetRed), new(big.Int).Set(game.TotalBetBlue), nil
}

// GetOdds 获取赔率
func (m *MockChainService) GetOdds(ctx context.Context, gameID uint64) (*big.Int, *big.Int, error) {
	totalRed, totalBlue, err := m.GetGamePool(ctx, gameID)
	if err != nil {
		return nil, nil, err
	}

	total := new(big.Int).Add(totalRed, totalBlue)
	if total.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(2), big.NewInt(2), nil
	}

	// 简化赔率计算
	oddsRed := new(big.Int).Div(total, maxBig(totalRed, big.NewInt(1)))
	oddsBlue := new(big.Int).Div(total, maxBig(totalBlue, big.NewInt(1)))

	return oddsRed, oddsBlue, nil
}

// GetGameStatus 获取对局状态
func (m *MockChainService) GetGameStatus(ctx context.Context, gameID uint64) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	game, ok := m.games[gameID]
	if !ok {
		return "", fmt.Errorf("game %d not found", gameID)
	}
	return game.Status, nil
}

// CreateGame 创建对局
func (m *MockChainService) CreateGame(ctx context.Context, agentRed, agentBlue string, bettingDuration uint64) (string, uint64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	gameID := m.nextGameID
	m.nextGameID++

	m.games[gameID] = &mockGame{
		Status:       "open",
		TotalBetRed:  big.NewInt(0),
		TotalBetBlue: big.NewInt(0),
	}

	txHash := fmt.Sprintf("0xmock_create_%d", gameID)

	// 发送事件
	m.emitEvent(ContractEvent{
		Type:   EventGameCreated,
		GameID: gameID,
		TxHash: txHash,
		Data: GameCreatedData{
			AgentRed:  agentRed,
			AgentBlue: agentBlue,
		},
	})

	return txHash, gameID, nil
}

// StartGame 开始对局
func (m *MockChainService) StartGame(ctx context.Context, gameID uint64) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	game, ok := m.games[gameID]
	if !ok {
		return "", fmt.Errorf("game %d not found", gameID)
	}
	if game.Status != "open" {
		return "", fmt.Errorf("game %d not in open state", gameID)
	}

	game.Status = "locked"
	txHash := fmt.Sprintf("0xmock_start_%d", gameID)

	m.emitEvent(ContractEvent{
		Type:   EventBettingLocked,
		GameID: gameID,
		TxHash: txHash,
		Data: BettingLockedData{
			TotalBetRed:  new(big.Int).Set(game.TotalBetRed),
			TotalBetBlue: new(big.Int).Set(game.TotalBetBlue),
		},
	})

	return txHash, nil
}

// FinishGame 结束对局
func (m *MockChainService) FinishGame(ctx context.Context, gameID uint64, redWins bool, actionsHash [32]byte) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	game, ok := m.games[gameID]
	if !ok {
		return "", fmt.Errorf("game %d not found", gameID)
	}
	if game.Status != "locked" {
		return "", fmt.Errorf("game %d not in locked state", gameID)
	}

	game.Status = "finished"
	game.RedWins = redWins
	game.ActionsHash = actionsHash

	totalPool := new(big.Int).Add(game.TotalBetRed, game.TotalBetBlue)
	fee := new(big.Int).Div(totalPool, big.NewInt(20)) // 5%

	txHash := fmt.Sprintf("0xmock_finish_%d", gameID)

	m.emitEvent(ContractEvent{
		Type:   EventGameSettled,
		GameID: gameID,
		TxHash: txHash,
		Data: GameSettledData{
			RedWins:     redWins,
			TotalPool:   totalPool,
			ProtocolFee: fee,
		},
	})

	return txHash, nil
}

// ListenEvents 监听事件
func (m *MockChainService) ListenEvents(ctx context.Context) (<-chan ContractEvent, error) {
	return m.events, nil
}

// --- Mock 辅助方法 ---

// SetStrategyWeights 设置策略权重（测试用）
func (m *MockChainService) SetStrategyWeights(gameID uint64, weights *engine.StrategyWeights) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.strategies[gameID] = weights
}

// AddBet 添加下注（测试用）
func (m *MockChainService) AddBet(gameID uint64, side string, amount *big.Int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	game, ok := m.games[gameID]
	if !ok {
		return fmt.Errorf("game %d not found", gameID)
	}
	if game.Status != "open" {
		return fmt.Errorf("game %d not open", gameID)
	}

	if side == "red" {
		game.TotalBetRed.Add(game.TotalBetRed, amount)
	} else {
		game.TotalBetBlue.Add(game.TotalBetBlue, amount)
	}

	m.emitEvent(ContractEvent{
		Type:   EventBetPlaced,
		GameID: gameID,
		TxHash: fmt.Sprintf("0xmock_bet_%d", gameID),
		Data: BetPlacedData{
			User:   "0xmock_user",
			Side:   side == "red",
			Amount: amount,
		},
	})

	return nil
}

func (m *MockChainService) emitEvent(event ContractEvent) {
	select {
	case m.events <- event:
	default:
		// 事件通道满，丢弃
	}
}

func maxBig(a, b *big.Int) *big.Int {
	if a.Cmp(b) > 0 {
		return a
	}
	return b
}
