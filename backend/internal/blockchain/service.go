package blockchain

import (
	"context"
	"math/big"

	"agent-arena/backend/internal/engine"
)

// ChainService 链上交互接口
type ChainService interface {
	// === 读取 ===

	// GetStrategyWeights 获取策略投票权重（从合约读取）
	GetStrategyWeights(ctx context.Context, gameID uint64) (*engine.StrategyWeights, error)

	// GetGamePool 获取下注池总额
	GetGamePool(ctx context.Context, gameID uint64) (totalRed, totalBlue *big.Int, err error)

	// GetOdds 获取当前赔率
	GetOdds(ctx context.Context, gameID uint64) (oddsRed, oddsBlue *big.Int, err error)

	// GetGameStatus 获取对局链上状态
	GetGameStatus(ctx context.Context, gameID uint64) (string, error)

	// === 写入 ===

	// CreateGame 创建对局（链上注册）
	CreateGame(ctx context.Context, agentRed, agentBlue string, bettingDuration uint64) (txHash string, gameID uint64, err error)

	// StartGame 开始对局（锁定下注 + 锁定投票）
	StartGame(ctx context.Context, gameID uint64) (txHash string, err error)

	// FinishGame 结束对局（提交结果 + 结算）
	FinishGame(ctx context.Context, gameID uint64, redWins bool, actionsHash [32]byte) (txHash string, err error)

	// === 事件 ===

	// ListenEvents 监听合约事件
	ListenEvents(ctx context.Context) (<-chan ContractEvent, error)
}
