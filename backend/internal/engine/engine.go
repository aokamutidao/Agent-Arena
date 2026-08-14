package engine

import "fmt"

// Engine 游戏引擎
type Engine struct{}

// NewEngine 创建游戏引擎
func NewEngine() *Engine {
	return &Engine{}
}

// NewGame 创建新对局
func (e *Engine) NewGame(gameID uint64, redName, blueName, redPersonality, bluePersonality string) *GameState {
	return &GameState{
		GameID:       gameID,
		Status:       StatusPending,
		CurrentRound: 0,
		MaxRounds:    MaxRounds,
		AgentRed: AgentState{
			ID:          "red",
			Name:        redName,
			Personality: redPersonality,
			HP:          MaxHP,
			MaxHP:       MaxHP,
			Position:    RedSpawn,
			Status:      []Effect{},
		},
		AgentBlue: AgentState{
			ID:          "blue",
			Name:        blueName,
			Personality: bluePersonality,
			HP:          MaxHP,
			MaxHP:       MaxHP,
			Position:    BlueSpawn,
			Status:      []Effect{},
		},
		Winner:    SideNone,
		History:   []TurnRecord{},
		Obstacles: DefaultObstacles,
	}
}

// ExecuteTurn 执行一个回合
func (e *Engine) ExecuteTurn(state *GameState, redAction, blueAction Action) (*TurnRecord, error) {
	if state.Status != StatusPlaying {
		return nil, fmt.Errorf("game not playing")
	}

	// 复制状态用于计算
	red := state.AgentRed
	blue := state.AgentBlue

	// 处理眩晕（眩晕时无法行动）
	if red.IsStunned {
		redAction = Action{Type: ActionWait}
		red.IsStunned = false
	}
	if blue.IsStunned {
		blueAction = Action{Type: ActionWait}
		blue.IsStunned = false
	}

	// 验证动作合法性（不合法标记为 Failed + 原因，保留原始意图，由 execute 函数兜底）
	if err := e.ValidateAction(redAction, red, blue, state.Obstacles); err != nil {
		redAction.Failed = true
		redAction.FailReason = err.Error()
	}
	if err := e.ValidateAction(blueAction, blue, red, state.Obstacles); err != nil {
		blueAction.Failed = true
		blueAction.FailReason = err.Error()
	}

	// 决定先手
	redFirst := e.determinePriority(red, blue, redAction, blueAction)

	// 按顺序执行动作
	if redFirst {
		e.applyAction(&red, &blue, redAction, state.Obstacles)
		e.applyAction(&blue, &red, blueAction, state.Obstacles)
	} else {
		e.applyAction(&blue, &red, blueAction, state.Obstacles)
		e.applyAction(&red, &blue, redAction, state.Obstacles)
	}

	// 更新冷却（使用技能的当回合也算一次冷却）
	if red.SkillCooldown > 0 {
		red.SkillCooldown--
	}
	if blue.SkillCooldown > 0 {
		blue.SkillCooldown--
	}
	if red.HealCooldown > 0 {
		red.HealCooldown--
	}
	if blue.HealCooldown > 0 {
		blue.HealCooldown--
	}

	// 更新状态效果持续时间
	e.tickEffects(&red)
	e.tickEffects(&blue)

	// 加时赛伤害：双方每回合各受 10 点伤害（Sudden Death）
	if state.Overtime {
		overtimeDmg := uint8(MaxHP / 10) // 10 HP
		red.HP = subtractHP(red.HP, overtimeDmg)
		blue.HP = subtractHP(blue.HP, overtimeDmg)
	}

	// 更新状态
	state.AgentRed = red
	state.AgentBlue = blue
	state.CurrentRound++

	// 记录回合
	record := &TurnRecord{
		Round:       state.CurrentRound,
		RedAction:   redAction,
		BlueAction:  blueAction,
		RedHPAfter:  red.HP,
		BlueHPAfter: blue.HP,
	}
	state.History = append(state.History, *record)

	// 检查胜负
	e.checkWinCondition(state)

	return record, nil
}

// determinePriority 决定先手
func (e *Engine) determinePriority(red, blue AgentState, redAction, blueAction Action) bool {
	// 蓄力中被攻击 → 被攻击方先手（反击）
	if red.IsCharging && (blueAction.Type == ActionAttack || blueAction.Type == ActionSkill) {
		return false // blue first
	}
	if blue.IsCharging && (redAction.Type == ActionAttack || redAction.Type == ActionSkill) {
		return true // red first
	}

	// 攻击优先于移动
	if redAction.Type == ActionAttack && blueAction.Type == ActionMove {
		return true
	}
	if blueAction.Type == ActionAttack && redAction.Type == ActionMove {
		return false
	}

	// 默认红方先手
	return true
}

// applyAction 执行动作
func (e *Engine) applyAction(self, target *AgentState, action Action, obstacles []Position) {
	switch action.Type {
	case ActionMove:
		e.executeMove(self, action.Target, target.Position, obstacles)
	case ActionAttack:
		e.executeAttack(self, target, obstacles)
	case ActionSkill:
		e.executeSkill(self, target, obstacles)
	case ActionCharge:
		self.IsCharging = true
		// Remaining=3：CHARGE 当回合末 -1 → 2，下一回合（可 MOVE 调整位置）末 -1 → 1，
		// 再下一回合仍可释放。如果 Remaining=2 的话，MOVE 一回合后就过期了。
		self.Status = append(self.Status, Effect{Type: EffectCharging, Remaining: 3})
	case ActionWait:
		// 不做任何事
	case ActionHeal:
		e.executeHeal(self)
	}
}

// executeHeal 执行治疗：回复 HealAmount HP（不超过 MaxHP）
func (e *Engine) executeHeal(self *AgentState) {
	if self.HealCooldown > 0 || self.HP >= self.MaxHP {
		return
	}
	healed := uint8(HealAmount)
	if self.HP+healed > self.MaxHP || self.HP+healed < self.HP {
		// 避免溢出，直接回满
		self.HP = self.MaxHP
	} else {
		self.HP += healed
	}
	// +1 因为当回合末会 -1
	self.HealCooldown = HealCooldown + 1
}

// executeMove 执行移动
func (e *Engine) executeMove(self *AgentState, target Position, enemyPos Position, obstacles []Position) {
	// 安全检查：目标不能越界
	if target.X >= MapWidth || target.Y >= MapHeight {
		return
	}

	// 检查距离
	dist := ManhattanDistance(self.Position, target)
	if dist > MOV || dist == 0 {
		return
	}

	// 检查是否是障碍物
	for _, obs := range obstacles {
		if obs == target {
			return
		}
	}

	// 检查是否是对手位置
	if target == enemyPos {
		return
	}

	self.Position = target
	// 注意：MOVE 不清除蓄力！AI 可以先 CHARGE，再 MOVE 调整位置，再 ATTACK 释放。
}

// executeAttack 执行近战攻击
func (e *Engine) executeAttack(self, target *AgentState, obstacles []Position) {
	dist := ManhattanDistance(self.Position, target.Position)
	if dist > AtkRange {
		return
	}
	// 障碍物阻挡攻击
	if !hasLineOfSight(self.Position, target.Position, obstacles) {
		return
	}

	damage := ATK
	if self.IsCharging {
		atkF := float64(ATK)
		multi := float64(ChargeMulti)
		damage = int(atkF * multi)
	}

	// 蓄力被打断
	if target.IsCharging {
		damage = int(float64(damage) * ChargeBreakDmg)
		target.IsCharging = false
		e.removeEffect(target, EffectCharging)
	}

	// 护盾减伤
	if e.hasEffect(target, EffectShielded) {
		damage = int(float64(damage) * ShieldReduce)
	}

	target.HP = subtractHP(target.HP, uint8(damage))
	self.IsCharging = false
	e.removeEffect(self, EffectCharging)
}

// executeSkill 执行远程技能
func (e *Engine) executeSkill(self, target *AgentState, obstacles []Position) {
	dist := ManhattanDistance(self.Position, target.Position)
	if dist > SkillRange || self.SkillCooldown > 0 {
		return
	}
	// 障碍物阻挡技能
	if !hasLineOfSight(self.Position, target.Position, obstacles) {
		return
	}

	damage := SkillDamage
	if self.IsCharging {
		atkF := float64(ATK)
		multi := float64(ChargeMulti)
		damage = int(atkF * multi * 0.8)
	}

	// 蓄力被打断
	if target.IsCharging {
		damage = int(float64(damage) * ChargeBreakDmg)
		target.IsCharging = false
		e.removeEffect(target, EffectCharging)
	}

	// 护盾减伤
	if e.hasEffect(target, EffectShielded) {
		damage = int(float64(damage) * ShieldReduce)
	}

	target.HP = subtractHP(target.HP, uint8(damage))
	self.SkillCooldown = SkillCooldown + 1 // +1 因为当回合末会 -1
	self.IsCharging = false
	e.removeEffect(self, EffectCharging)
}

// tickEffects 更新状态效果
func (e *Engine) tickEffects(agent *AgentState) {
	var remaining []Effect
	for _, eff := range agent.Status {
		eff.Remaining--
		if eff.Remaining > 0 {
			remaining = append(remaining, eff)
		} else if eff.Type == EffectCharging {
			// 蓄力效果过期，清除 IsCharging 标志
			agent.IsCharging = false
		}
	}
	agent.Status = remaining
}

// checkWinCondition 检查胜负
func (e *Engine) checkWinCondition(state *GameState) {
	if state.AgentRed.HP == 0 {
		state.Winner = SideBlue
		state.Status = StatusFinished
	} else if state.AgentBlue.HP == 0 {
		state.Winner = SideRed
		state.Status = StatusFinished
	} else if state.CurrentRound >= state.MaxRounds {
		// 回合用完
		if state.AgentRed.HP > state.AgentBlue.HP {
			state.Winner = SideRed
			state.Status = StatusFinished
		} else if state.AgentBlue.HP > state.AgentRed.HP {
			state.Winner = SideBlue
			state.Status = StatusFinished
		} else if !state.Overtime {
			// HP 相等 → 进入加时赛（Sudden Death），不结束
			state.Overtime = true
			state.MaxRounds += 10 // 再加 10 回合
		}
		// 加时赛中仍相等则继续（不结束，等下一回合再检查）
	}
}

// ValidateAction 验证动作合法性
func (e *Engine) ValidateAction(action Action, self AgentState, enemy AgentState, obstacles []Position) error {
	// 蓄力中只能 ATTACK/SKILL/MOVE（MOVE 用于调整位置后释放）。
	// CHARGE/WAIT/HEAL 在蓄力中是浪费。
	if self.IsCharging && action.Type != ActionAttack && action.Type != ActionSkill && action.Type != ActionMove {
		return fmt.Errorf("charging: must ATTACK, SKILL, or MOVE to reposition")
	}
	switch action.Type {
	case ActionMove:
		dist := ManhattanDistance(self.Position, action.Target)
		if dist > MOV || dist == 0 {
			return fmt.Errorf("move distance invalid")
		}
		for _, obs := range obstacles {
			if obs == action.Target {
				return fmt.Errorf("target is obstacle")
			}
		}
		if action.Target == enemy.Position {
			return fmt.Errorf("target is enemy position")
		}
		if action.Target.X >= MapWidth || action.Target.Y >= MapHeight {
			return fmt.Errorf("out of bounds")
		}
	case ActionAttack:
		dist := ManhattanDistance(self.Position, enemy.Position)
		if dist > AtkRange {
			return fmt.Errorf("enemy out of attack range")
		}
		if !hasLineOfSight(self.Position, enemy.Position, obstacles) {
			return fmt.Errorf("obstacle blocks line of sight")
		}
	case ActionSkill:
		dist := ManhattanDistance(self.Position, enemy.Position)
		if dist > SkillRange {
			return fmt.Errorf("enemy out of skill range")
		}
		if self.SkillCooldown > 0 {
			return fmt.Errorf("skill on cooldown")
		}
		if !hasLineOfSight(self.Position, enemy.Position, obstacles) {
			return fmt.Errorf("obstacle blocks line of sight")
		}
	case ActionCharge:
		// 总是合法
	case ActionWait:
		// 总是合法
	case ActionHeal:
		if self.HP >= self.MaxHP {
			return fmt.Errorf("hp already full")
		}
		if self.HealCooldown > 0 {
			return fmt.Errorf("heal on cooldown (%d turns)", self.HealCooldown)
		}
	default:
		return fmt.Errorf("unknown action type: %s", action.Type)
	}
	return nil
}

// 辅助函数

// hasLineOfSight 检查 from 到 to 的直线路径是否被障碍物阻挡。
// 使用 Bresenham 线段算法遍历路径上的每一格（不含端点）。
// 只要路径经过任一障碍物即返回 false。
func hasLineOfSight(from, to Position, obstacles []Position) bool {
	if len(obstacles) == 0 {
		return true
	}
	obsSet := make(map[Position]bool, len(obstacles))
	for _, o := range obstacles {
		obsSet[o] = true
	}

	x0 := int(from.X)
	y0 := int(from.Y)
	x1 := int(to.X)
	y1 := int(to.Y)

	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx := 1
	if x0 >= x1 {
		sx = -1
	}
	sy := 1
	if y0 >= y1 {
		sy = -1
	}
	err := dx - dy

	for {
		// 到达终点时停止（终点是对手格子，不视为阻挡）
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
		// 跳过起点（自身格子）和终点（对手格子）
		if (x0 == int(from.X) && y0 == int(from.Y)) || (x0 == x1 && y0 == y1) {
			continue
		}
		if obsSet[Position{X: uint8(x0), Y: uint8(y0)}] {
			return false
		}
	}
	return true
}

func ManhattanDistance(a, b Position) uint8 {
	dx := abs(int(a.X) - int(b.X))
	dy := abs(int(a.Y) - int(b.Y))
	return uint8(dx + dy)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (e *Engine) hasEffect(agent *AgentState, effectType EffectType) bool {
	for _, eff := range agent.Status {
		if eff.Type == effectType {
			return true
		}
	}
	return false
}

func (e *Engine) removeEffect(agent *AgentState, effectType EffectType) {
	var remaining []Effect
	for _, eff := range agent.Status {
		if eff.Type != effectType {
			remaining = append(remaining, eff)
		}
	}
	agent.Status = remaining
}

func subtractHP(current uint8, damage uint8) uint8 {
	if damage >= current {
		return 0
	}
	return current - damage
}
