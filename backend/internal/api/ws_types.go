package api

import "encoding/json"

// WebSocket 消息类型

// WSMessage WebSocket 消息基类
type WSMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// --- 服务端 → 客户端消息 ---

// GameStateData 游戏状态快照
type GameStateData struct {
	GameID          uint64 `json:"game_id"`
	Status          string `json:"status"`
	CurrentRound    uint8  `json:"current_round"`
	AgentRedHP      uint8  `json:"agent_red_hp"`
	AgentBlueHP     uint8  `json:"agent_blue_hp"`
	AgentRedPosX    uint8  `json:"agent_red_pos_x"`
	AgentRedPosY    uint8  `json:"agent_red_pos_y"`
	AgentBluePosX   uint8  `json:"agent_blue_pos_x"`
	AgentBluePosY   uint8  `json:"agent_blue_pos_y"`
}

// TurnUpdateData 回合更新
type TurnUpdateData struct {
	Round         uint8        `json:"round"`
	RedAction     ActionBrief  `json:"red_action"`
	BlueAction    ActionBrief  `json:"blue_action"`
	RedHP         uint8        `json:"red_hp"`
	BlueHP        uint8        `json:"blue_hp"`
	RedPosX       uint8        `json:"red_pos_x"`
	RedPosY       uint8        `json:"red_pos_y"`
	BluePosX      uint8        `json:"blue_pos_x"`
	BluePosY      uint8        `json:"blue_pos_y"`
	RedReasoning  string       `json:"red_reasoning,omitempty"`
	BlueReasoning string       `json:"blue_reasoning,omitempty"`
}

// ActionBrief 动作摘要
type ActionBrief struct {
	Type       string `json:"type"`
	Target     *Pos   `json:"target,omitempty"`
	Failed     bool   `json:"failed,omitempty"`
	FailReason string `json:"fail_reason,omitempty"`
}

// Pos 坐标
type Pos struct {
	X uint8 `json:"x"`
	Y uint8 `json:"y"`
}

// BettingUpdateData 下注池变化
type BettingUpdateData struct {
	TotalBetRed  string  `json:"total_bet_red"`
	TotalBetBlue string  `json:"total_bet_blue"`
	OddsRed      float64 `json:"odds_red"`
	OddsBlue     float64 `json:"odds_blue"`
}

// GameStartedData 对局开始
type GameStartedData struct {
	GameID       uint64 `json:"game_id"`
	Status       string `json:"status"`
	BettingLocked bool  `json:"betting_locked"`
}

// GameFinishedData 对局结束
type GameFinishedData struct {
	GameID     uint64 `json:"game_id"`
	Winner     string `json:"winner"`
	WinnerName string `json:"winner_name"`
	TotalRounds uint8 `json:"total_rounds"`
	FinalHPRed  uint8 `json:"final_hp_red"`
	FinalHPBlue uint8 `json:"final_hp_blue"`
}

// StrategyUpdateData 策略权重更新
type StrategyUpdateData struct {
	Side          string `json:"side"`
	Aggressive    uint8  `json:"aggressive"`
	Defensive     uint8  `json:"defensive"`
	Tricky        uint8  `json:"tricky"`
}

// OvertimeStartedData 加时赛开始
type OvertimeStartedData struct {
	GameID       uint64 `json:"game_id"`
	ExtraRounds  uint8  `json:"extra_rounds"`
	OvertimeDmg  uint8  `json:"overtime_dmg"`
}

// --- 客户端 → 服务端消息 ---

// WSSubscribe 订阅对局
type WSSubscribe struct {
	Type   string `json:"type"`
	GameID uint64 `json:"game_id"`
}
