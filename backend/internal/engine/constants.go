package engine

// 常量定义（来自 game-rules.md）
const (
	MapWidth  = 10
	MapHeight = 10
	MaxHP     = 100
	ATK       = 15
	MOV       = 2     // 移动力
	AtkRange  = 1     // 近战范围
	SkillRange    = 4 // 技能范围
	SkillCooldown = 3 // 技能冷却
	MaxRounds     = 30

	SkillDamage    = 12  // ATK * 0.8
	ChargeMulti    = 2.5 // 蓄力伤害倍率
	StunChance     = 0.2 // 眩晕概率
	StunDuration   = 1   // 眩晕持续回合
	ShieldReduce   = 0.5 // 护盾减伤
	ChargeBreakDmg = 1.2 // 蓄力被打断额外伤害

	// 治疗
	HealAmount   = 25 // ≈ 1.5 × ATK（ATK=15）
	HealCooldown = 3  // 回合
	HealCap      = MaxHP // 不能超过最大 HP
)

// 默认障碍物布局（对称，来自 game-rules.md）
var DefaultObstacles = []Position{
	{2, 3}, {2, 7},
	{4, 5}, {6, 5},
	{5, 2}, {7, 3},
	{7, 7}, {3, 5},
}

// 出生点
var (
	RedSpawn  = Position{1, 1}
	BlueSpawn = Position{8, 8}
)
