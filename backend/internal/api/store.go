package api

import (
	"fmt"
	"sync"

	"agent-arena/backend/internal/engine"
)

// GameRecord 对局记录（包含引擎状态 + 下注信息）
type GameRecord struct {
	Game               *engine.GameState
	TotalBetRed        uint64
	TotalBetBlue       uint64
	StrategyRed        engine.StrategyWeights // 百分比（用于展示）
	StrategyBlue       engine.StrategyWeights // 百分比（用于展示）
	StrategyRedVotes   [3]uint16              // 原始票数：[aggressive, defensive, tricky]
	StrategyBlueVotes  [3]uint16              // 原始票数：[aggressive, defensive, tricky]
}

// AgentRecord Agent 记录
type AgentRecord struct {
	Info    AgentInfo
	Games   []uint64 // 参与的对局 ID
}

// BetRecord 下注记录
type BetRecord struct {
	GameID   uint64 `json:"game_id"`
	Side     string `json:"side"`
	Amount   string `json:"amount"`
	Strategy string `json:"strategy"`
	Status   string `json:"status"` // pending, won, lost
	Reward   string `json:"reward"`
	Claimed  bool   `json:"claimed"`
}

// MemoryStore 内存数据存储（MVP）
type MemoryStore struct {
	games     map[uint64]*GameRecord
	agents    map[string]*AgentRecord
	bets      map[string][]*BetRecord // address -> bets
	nextID    uint64
	mu        sync.RWMutex
}

// NewMemoryStore 创建内存存储
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		games:  make(map[uint64]*GameRecord),
		agents: make(map[string]*AgentRecord),
		bets:   make(map[string][]*BetRecord),
	}
}

// NextGameID 返回下一个可用游戏 ID
func (s *MemoryStore) NextGameID() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	return s.nextID
}

// --- Game Operations ---

// CreateGame 创建对局
func (s *MemoryStore) CreateGame(game *engine.GameState) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.games[game.GameID] = &GameRecord{
		Game: game,
	}
	if game.GameID >= s.nextID {
		s.nextID = game.GameID + 1
	}
}

// GetGame 获取对局
func (s *MemoryStore) GetGame(gameID uint64) (*GameRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	game, ok := s.games[gameID]
	if !ok {
		return nil, fmt.Errorf("game %d not found", gameID)
	}
	return game, nil
}

// ListGames 列出对局
func (s *MemoryStore) ListGames(status string, limit, offset int) ([]*GameRecord, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*GameRecord
	for _, g := range s.games {
		if status != "" && string(g.Game.Status) != status {
			continue
		}
		result = append(result, g)
	}

	total := len(result)

	// offset
	if offset >= len(result) {
		return nil, total
	}
	result = result[offset:]

	// limit
	if limit > 0 && limit < len(result) {
		result = result[:limit]
	}

	return result, total
}

// UpdateBetPool 更新下注池
func (s *MemoryStore) UpdateBetPool(gameID uint64, side string, amount uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	game, ok := s.games[gameID]
	if !ok {
		return fmt.Errorf("game %d not found", gameID)
	}

	if side == "red" {
		game.TotalBetRed += amount
	} else {
		game.TotalBetBlue += amount
	}
	return nil
}

// UpdateStrategy 更新策略投票权重（对应方 +1 后重新计算百分比）
func (s *MemoryStore) UpdateStrategy(gameID uint64, side string, strategy string) (*engine.StrategyWeights, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, ok := s.games[gameID]
	if !ok {
		return nil, fmt.Errorf("game %d not found", gameID)
	}

	// 递增原始票数
	var votes *[3]uint16
	var weights *engine.StrategyWeights
	if side == "red" {
		votes = &record.StrategyRedVotes
		weights = &record.StrategyRed
	} else {
		votes = &record.StrategyBlueVotes
		weights = &record.StrategyBlue
	}

	// 对应策略 +1
	switch strategy {
	case "aggressive":
		votes[0]++
	case "defensive":
		votes[1]++
	case "tricky":
		votes[2]++
	}

	// 重新计算百分比
	total := uint16(votes[0]) + uint16(votes[1]) + uint16(votes[2])
	if total > 0 {
		weights.Aggressive = uint8(uint16(votes[0]) * 100 / total)
		weights.Defensive = uint8(uint16(votes[1]) * 100 / total)
		weights.Tricky = uint8(uint16(votes[2]) * 100 / total)
	}

	return weights, nil
}

// --- Agent Operations ---

// RegisterAgent 注册 Agent
func (s *MemoryStore) RegisterAgent(info AgentInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.agents[info.ID] = &AgentRecord{
		Info: info,
	}
}

// GetAgent 获取 Agent
func (s *MemoryStore) GetAgent(id string) (*AgentRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	agent, ok := s.agents[id]
	if !ok {
		return nil, fmt.Errorf("agent %s not found", id)
	}
	return agent, nil
}

// ListAgents 列出所有 Agent
func (s *MemoryStore) ListAgents() []AgentInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []AgentInfo
	for _, a := range s.agents {
		result = append(result, a.Info)
	}
	return result
}

// UpdateAgentStats 更新 Agent 胜负统计
func (s *MemoryStore) UpdateAgentStats(id string, won bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	agent, ok := s.agents[id]
	if !ok {
		return
	}
	if won {
		agent.Info.Wins++
	} else {
		agent.Info.Losses++
	}
	total := agent.Info.Wins + agent.Info.Losses
	if total > 0 {
		agent.Info.WinRate = float64(agent.Info.Wins) / float64(total) * 100
	}
}

// --- Bet Operations ---

// CacheBet 缓存下注记录
func (s *MemoryStore) CacheBet(address string, bet *BetRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.bets[address] = append(s.bets[address], bet)
}

// GetBetsByUser 获取用户下注记录
func (s *MemoryStore) GetBetsByUser(address string) []*BetRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.bets[address]
}

// --- Odds Calculation ---

// CalculateOdds 计算赔率
func CalculateOdds(totalRed, totalBlue uint64) (float64, float64) {
	total := float64(totalRed + totalBlue)
	if total == 0 {
		return 2.0, 2.0 // 默认赔率
	}

	// 赔率 = 总池 / 该方池（如果该方池为 0 则用最小值）
	redPool := float64(totalRed)
	bluePool := float64(totalBlue)

	oddsRed := total / max(redPool, 1)
	oddsBlue := total / max(bluePool, 1)

	return oddsRed, oddsBlue
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
