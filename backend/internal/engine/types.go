package engine

// Position 格子坐标
type Position struct {
	X uint8 `json:"x"`
	Y uint8 `json:"y"`
}

// ActionType 动作类型
type ActionType string

const (
	ActionMove   ActionType = "MOVE"
	ActionAttack ActionType = "ATTACK"
	ActionSkill  ActionType = "SKILL"
	ActionCharge ActionType = "CHARGE"
	ActionWait   ActionType = "WAIT"
	ActionHeal   ActionType = "HEAL" // 治疗（恢复 HealAmount HP，冷却 HealCooldown 回合）
)

// Action 玩家动作
type Action struct {
	Type       ActionType `json:"type"`
	Target     Position   `json:"target,omitempty"`
	Failed     bool       `json:"failed,omitempty"`      // 动作尝试但不合法
	FailReason string     `json:"fail_reason,omitempty"` // 失败原因
}

// EffectType 状态效果
type EffectType string

const (
	EffectCharging EffectType = "charging"
	EffectStunned  EffectType = "stunned"
	EffectShielded EffectType = "shielded"
)

// Effect 状态效果
type Effect struct {
	Type      EffectType `json:"type"`
	Remaining uint8      `json:"remaining"` // 剩余回合数
}

// AgentState Agent 状态
type AgentState struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Personality   string   `json:"personality"`
	HP            uint8    `json:"hp"`
	MaxHP         uint8    `json:"max_hp"`
	Position      Position `json:"position"`
	Status        []Effect `json:"status"`
	SkillCooldown uint8    `json:"skill_cooldown"`
	HealCooldown  uint8    `json:"heal_cooldown"` // 治疗冷却（每回合 -1）
	IsCharging    bool     `json:"is_charging"`
	IsStunned     bool     `json:"is_stunned"`
}

// Side 阵营
type Side string

const (
	SideNone  Side = "none"
	SideRed   Side = "red"
	SideBlue  Side = "blue"
)

// GameStatus 对局状态
type GameStatus string

const (
	StatusPending  GameStatus = "pending"
	StatusBetting  GameStatus = "betting"
	StatusPlaying  GameStatus = "playing"
	StatusFinished GameStatus = "finished"
)

// TurnRecord 回合记录
type TurnRecord struct {
	Round        uint8  `json:"round"`
	RedAction    Action `json:"red_action"`
	BlueAction   Action `json:"blue_action"`
	RedHPAfter   uint8  `json:"red_hp_after"`
	BlueHPAfter  uint8  `json:"blue_hp_after"`
	RedReasoning string `json:"red_reasoning,omitempty"`
	BlueReasoning string `json:"blue_reasoning,omitempty"`
}

// GameState 游戏状态
type GameState struct {
	GameID       uint64      `json:"game_id"`
	Status       GameStatus  `json:"status"`
	CurrentRound uint8       `json:"current_round"`
	MaxRounds    uint8       `json:"max_rounds"`
	AgentRed     AgentState  `json:"agent_red"`
	AgentBlue    AgentState  `json:"agent_blue"`
	Winner       Side        `json:"winner"`
	History      []TurnRecord `json:"history"`
	Obstacles    []Position  `json:"obstacles"`
	Overtime     bool        `json:"overtime"`      // 加时赛标志

	// 自定义 Agent API 配置（可选）
	AgentRedAPIEndpoint string `json:"agent_red_api_endpoint,omitempty"`
	AgentRedAPIKey      string `json:"-"` // 不序列化到 JSON（敏感信息）
	AgentRedModel       string `json:"agent_red_model,omitempty"`

	// 链上游戏 ID（用于 FinishGame 调用）
	OnChainGameID uint64 `json:"-"` // 不序列化到 JSON
}

// StrategyWeights 策略投票权重（百分比 0-100）
type StrategyWeights struct {
	Aggressive uint8 `json:"aggressive"`
	Defensive  uint8 `json:"defensive"`
	Tricky     uint8 `json:"tricky"`
}
