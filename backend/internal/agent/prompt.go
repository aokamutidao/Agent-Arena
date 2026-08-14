package agent

import (
	"fmt"
	"strings"

	"agent-arena/backend/internal/engine"
)

// AgentDecisionInput Agent 决策所需的输入
type AgentDecisionInput struct {
	Self           engine.AgentState
	Enemy          engine.AgentState
	CurrentRound   uint8
	MaxRounds      uint8
	Obstacles      []engine.Position
	Weights        engine.StrategyWeights
	ActionHistory  []engine.TurnRecord // 最近回合记录（短期记忆）
}

// buildDecisionPrompt 构建决策 prompt
func buildDecisionPrompt(name, personality string, input AgentDecisionInput) string {
	dist := engine.ManhattanDistance(input.Self.Position, input.Enemy.Position)

	// 构建行动历史（最近 5 回合）
	historyText := formatActionHistory(input.ActionHistory, input.Self.ID)

	// 计算合法移动目标 + 推荐方向
	validMoves, recommendedMove := computeValidMoves(input)

	// 计算当前是否能有效攻击
	canAttack := dist <= 1
	canSkill := dist <= 4 && input.Self.SkillCooldown == 0
	canHeal := input.Self.HealCooldown == 0 && input.Self.HP < input.Self.MaxHP

	// HP 阈值：低于 40% 就值得考虑治疗
	hpRatio := float64(input.Self.HP) / float64(input.Self.MaxHP)
	hpLow := hpRatio < 0.4
	hpCritical := hpRatio < 0.2

	// 显式列出当前【合法动作】，避免 AI 误判距离
	validActions := ""
	if input.Self.IsCharging {
		// ⚡ 蓄力中：极度简化，只有 2 个选项，不给 LLM 犯错空间
		if canAttack {
			validActions = "✅ ATTACK（释放蓄力！伤害 ×2.5）"
		} else if canSkill {
			validActions = "✅ SKILL（释放蓄力！伤害 ×2.5）"
		} else {
			validActions = "❌ ATTACK（距离太远） ❌ SKILL — 只能 MOVE 靠近！下回合必须 ATTACK！"
		}
	} else {
		switch {
		case canAttack && canSkill:
			validActions = "✅ ATTACK ✅ SKILL"
		case canAttack:
			validActions = "✅ ATTACK ❌ SKILL（冷却中）"
		case canSkill:
			validActions = fmt.Sprintf("❌ ATTACK（距离 %d > 1）✅ SKILL", dist)
		default:
			validActions = fmt.Sprintf("❌ ATTACK（距离 %d）❌ SKILL", dist)
		}
		if canHeal {
			validActions += fmt.Sprintf(" ✅ HEAL（恢复 %d HP）", engine.HealAmount)
		} else if input.Self.HealCooldown > 0 {
			validActions += fmt.Sprintf(" ❌ HEAL（冷却 %d 回合）", input.Self.HealCooldown)
		}
		validActions += " ✅ MOVE ✅ CHARGE ✅ WAIT"
	}

	// 蓄力状态提示 — 放在最顶部，最强优先级
	chargeText := ""
	if input.Self.IsCharging {
		chargeText = "⚡⚡⚡ 你正在蓄力中！本回合【只能 ATTACK 或 MOVE】！绝对禁止 CHARGE！如果距离够就直接 ATTACK（×2.5伤害），距离不够就 MOVE 靠近，下回合 ATTACK！蓄力 2 回合内不释放就浪费了！⚡⚡⚡\n\n"
	}

	// 低血量提示：鼓励治疗/撤退
	hpAdvice := ""
	if hpCritical {
		hpAdvice = "\n🚨🚨 危急！HP < 20%！本回合【优先 HEAL】或【远离对手 MOVE 到障碍物后】，否则下回合可能被秒杀！"
	} else if hpLow {
		hpAdvice = "\n⚠️ 低血量（< 40%）。建议 HEAL 或边打边退；如果对手靠近且 HEAL 冷却中，考虑绕障碍物拉扯。"
	}

	// 战术提示：距离/蓄力/掩体
	tacticalAdvice := buildTacticalAdvice(input, dist, hpLow, hpCritical)

	return fmt.Sprintf(`%s你是 Agent Arena 中的战斗 AI "%s"。

## 你的性格
%s

## 当前局势
- 你的位置: (%d, %d)，HP: %d/100
- 对手位置: (%d, %d)，HP: %d/100
- 你与对手距离: %d 格（曼哈顿距离）
- 你的技能冷却: %d 回合
- 当前回合: %d/%d
- 你的状态: %s
- 对手状态: %s

## 地图障碍物（不能走上去）
%s

## 可选动作
- MOVE(x,y)  — 移动到指定坐标（必须是下方列出的合法位置之一！）
- ATTACK     — 近战攻击（距离 = 1 才有效！距离 ≥ 2 禁止！），伤害 15
- SKILL      — 远程技能（距离 ≤ 4 且冷却 = 0），伤害 12
- HEAL       — 自我治疗（恢复 %d HP，冷却 %d 回合，HP 满时禁用）
- CHARGE     — 蓄力（下回合 ATTACK/SKILL 伤害 ×2.5，蓄力后务必立刻攻击！）
- WAIT       — 原地等待

## 当前可用动作（基于距离 %d）
%s
%s
%s

## ⚠️ 重要规则（违反将导致动作失败！）
1. ATTACK 只在距离 = 1 时有效！距离 ≥ 2 时 ATTACK 100%% 失败，必须先 MOVE 靠近！
2. SKILL 只在距离 ≤ 4 且冷却 = 0 时有效。
3. MOVE(x,y) 的 x,y 必须是你当前位置 1-2 格内的坐标！不能跳到远处！
4. 如果距离 ≥ 2，绝对不要输出 ATTACK！必须 MOVE 或 SKILL！
5. 障碍物会阻挡直线视线！你与对手之间若有障碍物，ATTACK/SKILL 会失败！需绕开障碍物再攻击。
6. 蓄力后 2 回合内必须 ATTACK/SKILL 释放，但中间可以 MOVE 调整位置！蓄力中不要 CHARGE/WAIT/HEAL。

## 你可以移动到的位置（从你当前位置 %d,%d 出发，1-2 格内）
%s

%s

## 输出格式（严格按此格式！）
<动作>
思考: <你的判断>

示例:
MOVE(%s)
思考: 距离5，先向对手靠近`,
		chargeText,
		name,
		personality,
		input.Self.Position.X, input.Self.Position.Y, input.Self.HP,
		input.Enemy.Position.X, input.Enemy.Position.Y, input.Enemy.HP,
		dist,
		input.Self.SkillCooldown,
		input.CurrentRound, input.MaxRounds,
		formatStatus(input.Self.Status),
		formatStatus(input.Enemy.Status),
		formatObstacles(input.Obstacles),
		engine.HealAmount, engine.HealCooldown,
		dist,
		validActions,
		hpAdvice,
		tacticalAdvice,
		input.Self.Position.X, input.Self.Position.Y,
		validMoves,
		historyText,
		recommendedMove,
	)
}

// buildTacticalAdvice 基于局势给出战术提示（伏击/拉扯/掩体/追击）
func buildTacticalAdvice(input AgentDecisionInput, dist uint8, hpLow, hpCritical bool) string {
	var tips []string

	// 1. 危急状态：立即治疗
	if hpCritical {
		if input.Self.HealCooldown == 0 {
			tips = append(tips, "🚑 【立即 HEAL】！下回合可能被秒杀，先回血 25 HP 再说！")
		} else {
			tips = append(tips, fmt.Sprintf("🏃 【远离对手】！HEAL 冷却 %d 回合，绕障碍物拉扯，等冷却好立刻治疗！", input.Self.HealCooldown))
		}
	} else if hpLow && input.Self.HealCooldown == 0 && dist > 1 {
		tips = append(tips, "💊 血量偏低。对手还离得远，本回合 HEAL 比冒险靠近更安全。")
	}

	// 2. 伏击：对手远 + 路径上有障碍物 → 蓄力埋伏
	if dist >= 3 && !input.Self.IsCharging && input.Self.HealCooldown == 0 && !hpLow {
		// 检查是否靠近障碍物（曼哈顿距离 ≤ 2 内有障碍物）
		nearObs := false
		for _, o := range input.Obstacles {
			if engine.ManhattanDistance(input.Self.Position, o) <= 2 {
				nearObs = true
				break
			}
		}
		if nearObs {
			tips = append(tips, "🗡️ 伏击机会！对手还远 + 你在掩体旁，本回合 CHARGE，下回合他靠近时一击必杀（×2.5 伤害）！")
		}
	}

	// 3. 风筝：对手近 + 自己有远程技能 → 边打边退
	if dist == 1 && canSkillNow(input) && !hpCritical {
		// 检查身后是否有退路（离自己 1-2 格且离对手更远的合法位置）
		hasRetreat := false
		for dx := -2; dx <= 2; dx++ {
			for dy := -2; dy <= 2; dy++ {
				nx := int(input.Self.Position.X) - dx
				ny := int(input.Self.Position.Y) - dy
				if nx < 0 || nx > 9 || ny < 0 || ny > 9 {
					continue
				}
				newPos := engine.Position{X: uint8(nx), Y: uint8(ny)}
				if engine.ManhattanDistance(newPos, input.Enemy.Position) > 1 {
					hasRetreat = true
					break
				}
			}
			if hasRetreat {
				break
			}
		}
		if hasRetreat {
			tips = append(tips, "🏹 风筝打法：先 SKILL 远程输出，下回合 MOVE 拉开距离，让近战对手摸不到你！")
		}
	}

	// 4. 掩体使用：对手远 + 自己血量健康 + 中间有障碍物
	if dist >= 4 && !hpLow {
		tips = append(tips, "🛡️ 对手还远，MOVE 靠近时优先选【障碍物后的位置】——让他必须绕路，你保持射程优势。")
	}

	// 5. 蓄力被打断风险提示
	if input.Self.IsCharging {
		if dist == 1 {
			tips = append(tips, "⚡ 蓄力中 + 对手贴脸：本回合必须 ATTACK（×2.5 = 37 伤害），可以直接带走对手大片血！")
		} else {
			tips = append(tips, "⚡ 蓄力中 + 对手较远：本回合建议 MOVE 靠近，下回合再攻击释放蓄力伤害。")
		}
	}

	if len(tips) == 0 {
		return "（无特别提示，根据性格自由发挥）"
	}
	return "• " + tips[0] + func() string {
		if len(tips) > 1 {
			return "\n• " + tips[1]
		}
		return ""
	}()
}

// canSkillNow 辅助：判断本回合是否能用技能
func canSkillNow(input AgentDecisionInput) bool {
	dist := engine.ManhattanDistance(input.Self.Position, input.Enemy.Position)
	return dist <= engine.SkillRange && input.Self.SkillCooldown == 0
}

// computeValidMoves 计算合法移动目标 + 推荐方向
func computeValidMoves(input AgentDecisionInput) (string, string) {
	self := input.Self
	enemy := input.Enemy
	obstacles := input.Obstacles

	// 构建障碍物集合
	obsSet := make(map[engine.Position]bool)
	for _, o := range obstacles {
		obsSet[o] = true
	}

	var moves []string
	var towardEnemy []string // 靠近对手的方向

	for dx := -2; dx <= 2; dx++ {
		for dy := -2; dy <= 2; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			absD := abs(dx) + abs(dy)
			if absD > 2 || absD == 0 {
				continue
			}
			nx := int(self.Position.X) + dx
			ny := int(self.Position.Y) + dy
			if nx < 0 || nx > 9 || ny < 0 || ny > 9 {
				continue
			}
			pos := engine.Position{X: uint8(nx), Y: uint8(ny)}
			if obsSet[pos] {
				continue
			}
			if pos == enemy.Position {
				continue
			}
			moveStr := fmt.Sprintf("(%d,%d)", nx, ny)
			moves = append(moves, moveStr)

			// 是否靠近对手
			newDist := engine.ManhattanDistance(pos, enemy.Position)
			curDist := engine.ManhattanDistance(self.Position, enemy.Position)
			if newDist < curDist {
				towardEnemy = append(towardEnemy, moveStr)
			}
		}
	}

	movesText := strings.Join(moves, " ")
	if len(moves) == 0 {
		movesText = "无（被包围或贴边）"
	}

	// 推荐移动：选最靠近对手的合法位置
	recommended := ""
	if len(towardEnemy) > 0 {
		// 选第一个靠近对手的位置
		recommended = strings.Trim(towardEnemy[0], "()")
	} else if len(moves) > 0 {
		recommended = strings.Trim(moves[0], "()")
	} else {
		recommended = "0,0"
	}

	return movesText, recommended
}

// abs 绝对值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// formatActionHistory 格式化最近行动历史（短期记忆）
func formatActionHistory(history []engine.TurnRecord, selfID string) string {
	if len(history) == 0 {
		return "## 最近行动记录\n（第一回合，暂无记录）"
	}

	// 取最近 5 回合
	start := len(history) - 5
	if start < 0 {
		start = 0
	}
	recent := history[start:]

	var lines []string
	lines = append(lines, "## 最近行动记录（你=我）")
	for _, turn := range recent {
		var myAction, enemyAction engine.Action
		if selfID == "red" {
			myAction = turn.RedAction
			enemyAction = turn.BlueAction
		} else {
			myAction = turn.BlueAction
			enemyAction = turn.RedAction
		}

		myText := formatActionWithResult(myAction)
		enemyText := formatActionWithResult(enemyAction)
		lines = append(lines, fmt.Sprintf("- R%d: 我=%s | 对手=%s | HP: %d/%d",
			turn.Round, myText, enemyText, turn.RedHPAfter, turn.BlueHPAfter))
	}
	return strings.Join(lines, "\n")
}

// formatActionWithResult 格式化动作（含失败标记）
func formatActionWithResult(action engine.Action) string {
	text := formatAction(action)
	if action.Failed {
		text += "(失败:" + action.FailReason + ")"
	}
	return text
}

// buildMonologuePrompt 构建独白 prompt
func buildMonologuePrompt(name, personality string, input AgentDecisionInput, action engine.Action) string {
	return fmt.Sprintf(`你是 Agent Arena 中的战斗 AI "%s"。

刚才的回合你选择了 %s。
当前你的 HP: %d，对手 HP: %d。

用一句话表达你此刻的想法（20 字以内，符合你的性格）。
只回复这句话，不要其他内容。

你的性格: %s`,
		name,
		formatAction(action),
		input.Self.HP, input.Enemy.HP,
		personality,
	)
}

// formatStatus 格式化状态效果
func formatStatus(effects []engine.Effect) string {
	if len(effects) == 0 {
		return "无"
	}
	var parts []string
	for _, eff := range effects {
		parts = append(parts, fmt.Sprintf("%s(%d回合)", eff.Type, eff.Remaining))
	}
	return strings.Join(parts, ", ")
}

// formatObstacles 格式化障碍物坐标
func formatObstacles(obstacles []engine.Position) string {
	if len(obstacles) == 0 {
		return "无"
	}
	var parts []string
	for _, obs := range obstacles {
		parts = append(parts, fmt.Sprintf("(%d,%d)", obs.X, obs.Y))
	}
	return strings.Join(parts, " ")
}

// formatAction 格式化动作文本
func formatAction(action engine.Action) string {
	if action.Type == engine.ActionMove {
		return fmt.Sprintf("MOVE(%d,%d)", action.Target.X, action.Target.Y)
	}
	return string(action.Type)
}
