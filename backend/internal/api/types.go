package api

import (
	"agent-arena/backend/internal/engine"
)

// --- API Response Types ---

// AgentInfo API 返回的 Agent 信息
type AgentInfo struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Personality string  `json:"personality"`
	Wins        int     `json:"wins"`
	Losses      int     `json:"losses"`
	WinRate     float64 `json:"win_rate"`
	Description string  `json:"description"`
}

// GameListItem 对局列表项
type GameListItem struct {
	GameID        uint64    `json:"game_id"`
	AgentRed      AgentInfo `json:"agent_red"`
	AgentBlue     AgentInfo `json:"agent_blue"`
	Status        string    `json:"status"`
	TotalBetRed   string    `json:"total_bet_red"`
	TotalBetBlue  string    `json:"total_bet_blue"`
	CurrentRound  uint8     `json:"current_round"`
	MaxRounds     uint8     `json:"max_rounds"`
	OddsRed       float64   `json:"odds_red"`
	OddsBlue      float64   `json:"odds_blue"`
}

// GamesResponse GET /api/games 响应
type GamesResponse struct {
	Games []GameListItem `json:"games"`
	Total int            `json:"total"`
}

// AgentStateResponse Agent 运行时状态
type AgentStateResponse struct {
	HP            uint8             `json:"hp"`
	MaxHP         uint8             `json:"max_hp"`
	Position      engine.Position   `json:"position"`
	Status        []string          `json:"status"`
	SkillCooldown uint8             `json:"skill_cooldown"`
	IsCharging    bool              `json:"is_charging"`
}

// GameDetailResponse GET /api/games/:gameId 响应
type GameDetailResponse struct {
	GameID        uint64             `json:"game_id"`
	AgentRed      AgentInfo          `json:"agent_red"`
	AgentBlue     AgentInfo          `json:"agent_blue"`
	Status        string             `json:"status"`
	CurrentRound  uint8              `json:"current_round"`
	MaxRounds     uint8              `json:"max_rounds"`
	TotalBetRed   string             `json:"total_bet_red"`
	TotalBetBlue  string             `json:"total_bet_blue"`
	AgentRedState *AgentStateResponse `json:"agent_red_state,omitempty"`
	AgentBlueState *AgentStateResponse `json:"agent_blue_state,omitempty"`
	StrategyRed   *engine.StrategyWeights `json:"strategy_red,omitempty"`
	StrategyBlue  *engine.StrategyWeights `json:"strategy_blue,omitempty"`
	Overtime      bool               `json:"overtime"`
	History       []engine.TurnRecord `json:"history"`
	Winner        string             `json:"winner,omitempty"`
	Archived      bool               `json:"archived,omitempty"` // 重启后从 DB 兜底恢复的已完成对局（无回合明细）
}

// EstimateRequest POST /api/bets/estimate 请求
type EstimateRequest struct {
	GameID uint64 `json:"game_id" binding:"required"`
	Side   string `json:"side" binding:"required,oneof=red blue"`
	Amount string `json:"amount" binding:"required"`
}

// EstimateResponse POST /api/bets/estimate 响应
type EstimateResponse struct {
	CurrentPoolRed  string  `json:"current_pool_red"`
	CurrentPoolBlue string  `json:"current_pool_blue"`
	NewPoolRed      string  `json:"new_pool_red"`
	PotentialReward string  `json:"potential_reward"`
	NewOddsRed      float64 `json:"new_odds_red"`
	NewOddsBlue     float64 `json:"new_odds_blue"`
}

// AgentsResponse GET /api/agents 响应
type AgentsResponse struct {
	Agents []AgentInfo `json:"agents"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code,omitempty"`
}

// VoteStrategyRequest POST /api/strategy/vote 请求
type VoteStrategyRequest struct {
	GameID   uint64 `json:"game_id" binding:"required"`
	Side     string `json:"side" binding:"required,oneof=red blue"`
	Strategy string `json:"strategy" binding:"required,oneof=aggressive defensive tricky"`
	User     string `json:"user"`
}

// VoteStrategyResponse POST /api/strategy/vote 响应
type VoteStrategyResponse struct {
	Side       string `json:"side"`
	Aggressive uint8  `json:"aggressive"`
	Defensive  uint8  `json:"defensive"`
	Tricky     uint8  `json:"tricky"`
}
