package agent

import (
	"context"
	"fmt"
	"sync"

	"agent-arena/backend/internal/engine"
)

// Agent AI Agent
type Agent struct {
	Name        string
	Personality string
	LLM         LLMClient
}

// NewAgent 创建 Agent
func NewAgent(name, personality string, llm LLMClient) *Agent {
	return &Agent{
		Name:        name,
		Personality: personality,
		LLM:         llm,
	}
}

// DecideTurn Agent 决策（调用 LLM → 解析 → 验证）— 返回动作 + 推理
func (a *Agent) DecideTurn(input AgentDecisionInput) (engine.Action, string) {
	// 1. 构建 prompt
	prompt := buildDecisionPrompt(a.Name, a.Personality, input)

	// 2. 调用 LLM
	response, err := a.LLM.Chat(context.Background(), prompt)
	if err != nil {
		return engine.Action{Type: engine.ActionWait}, "LLM 调用失败"
	}

	// 3. 解析动作 + 推理
	action, reasoning, parseErr := ParseActionWithReasoning(response)

	if parseErr != nil {
		return engine.Action{Type: engine.ActionWait}, "解析失败: " + response
	}

	// 4. 验证合法性（不合法 → 标记 Failed，保留原始意图）
	if !validateAgentAction(action, input.Self, input.Enemy, input.Obstacles) {
		action.Failed = true
		action.FailReason = validateReason(action, input.Self, input.Enemy, input.Obstacles)
		return action, reasoning
	}

	return action, reasoning
}

// validateReason 返回动作不合法的具体原因
func validateReason(action engine.Action, self, enemy engine.AgentState, obstacles []engine.Position) string {
	switch action.Type {
	case engine.ActionMove:
		dist := engine.ManhattanDistance(self.Position, action.Target)
		if dist > 2 || dist == 0 {
			return "移动距离超限(需≤2格)"
		}
		for _, obs := range obstacles {
			if obs == action.Target {
				return "目标是障碍物"
			}
		}
		if action.Target == enemy.Position {
			return "目标是对手位置"
		}
		if action.Target.X >= 10 || action.Target.Y >= 10 {
			return "超出地图边界"
		}
	case engine.ActionAttack:
		dist := engine.ManhattanDistance(self.Position, enemy.Position)
		if dist > 1 {
			return fmt.Sprintf("攻击距离不足(距离%d,需≤1)", dist)
		}
	case engine.ActionSkill:
		dist := engine.ManhattanDistance(self.Position, enemy.Position)
		if dist > 4 {
			return fmt.Sprintf("技能距离不足(距离%d,需≤4)", dist)
		}
		if self.SkillCooldown > 0 {
			return fmt.Sprintf("技能冷却中(%d回合)", self.SkillCooldown)
		}
	case engine.ActionHeal:
		if self.HP >= self.MaxHP {
			return "HP已满"
		}
		if self.HealCooldown > 0 {
			return fmt.Sprintf("治疗冷却中(%d回合)", self.HealCooldown)
		}
	}
	return "动作不合法"
}

// GenerateMonologue 异步生成独白（不阻塞）
func (a *Agent) GenerateMonologue(input AgentDecisionInput, action engine.Action) <-chan string {
	ch := make(chan string, 1)

	go func() {
		prompt := buildMonologuePrompt(a.Name, a.Personality, input, action)
		response, err := a.LLM.Chat(context.Background(), prompt)
		if err != nil {
			ch <- "..."
			return
		}
		ch <- response
	}()

	return ch
}

// DecideResult 决策结果
type DecideResult struct {
	Action   engine.Action
	Reasoning string
}

// DecideTurnsParallel 两 Agent 并行决策
func DecideTurnsParallel(red, blue *Agent, redInput, blueInput AgentDecisionInput) (DecideResult, DecideResult) {
	var redResult, blueResult DecideResult
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		action, reasoning := red.DecideTurn(redInput)
		redResult = DecideResult{Action: action, Reasoning: reasoning}
	}()

	go func() {
		defer wg.Done()
		action, reasoning := blue.DecideTurn(blueInput)
		blueResult = DecideResult{Action: action, Reasoning: reasoning}
	}()

	wg.Wait()
	return redResult, blueResult
}

// validateAgentAction 验证 Agent 动作合法性
func validateAgentAction(action engine.Action, self, enemy engine.AgentState, obstacles []engine.Position) bool {
	switch action.Type {
	case engine.ActionMove:
		dist := engine.ManhattanDistance(self.Position, action.Target)
		if dist > 2 || dist == 0 {
			return false
		}
		for _, obs := range obstacles {
			if obs == action.Target {
				return false
			}
		}
		if action.Target == enemy.Position {
			return false
		}
		if action.Target.X >= 10 || action.Target.Y >= 10 {
			return false
		}
		return true

	case engine.ActionAttack:
		return engine.ManhattanDistance(self.Position, enemy.Position) <= 1

	case engine.ActionSkill:
		return engine.ManhattanDistance(self.Position, enemy.Position) <= 4 && self.SkillCooldown == 0

	case engine.ActionCharge:
		return true

	case engine.ActionWait:
		return true

	case engine.ActionHeal:
		return self.HP < self.MaxHP && self.HealCooldown == 0

	default:
		return false
	}
}
