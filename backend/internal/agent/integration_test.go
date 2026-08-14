package agent

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"agent-arena/backend/internal/engine"
)

// MockLLMClient 模拟 LLM 客户端
type MockLLMClient struct {
	mu        sync.Mutex
	responses []string
	index     int
	callCount int
}

// NewMockLLM 创建模拟 LLM
func NewMockLLM(responses []string) *MockLLMClient {
	return &MockLLMClient{responses: responses}
}

// Chat 返回预设的响应（循环）
func (m *MockLLMClient) Chat(ctx context.Context, prompt string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callCount++
	resp := m.responses[m.index%len(m.responses)]
	m.index++
	return resp, nil
}

// FailLLMClient 总是失败的 LLM 客户端
type FailLLMClient struct{}

func (f *FailLLMClient) Chat(ctx context.Context, prompt string) (string, error) {
	return "", fmt.Errorf("mock api error")
}

func TestDecideTurn_Attack(t *testing.T) {
	mock := NewMockLLM([]string{"ATTACK"})
	agent := NewAgent("Berserker", Berserker.Description, mock)

	// 对手在距离 1
	input := AgentDecisionInput{
		Self:    engine.AgentState{Position: engine.Position{X: 3, Y: 3}},
		Enemy:   engine.AgentState{Position: engine.Position{X: 4, Y: 3}},
		Weights: engine.StrategyWeights{Aggressive: 100},
	}

	action, _ := agent.DecideTurn(input)
	if action.Type != engine.ActionAttack {
		t.Fatalf("expected ATTACK, got %s", action.Type)
	}
}

func TestDecideTurn_Move(t *testing.T) {
	mock := NewMockLLM([]string{"MOVE(5,3)"})
	agent := NewAgent("Tactician", Tactician.Description, mock)

	input := AgentDecisionInput{
		Self:    engine.AgentState{Position: engine.Position{X: 3, Y: 3}},
		Enemy:   engine.AgentState{Position: engine.Position{X: 8, Y: 8}},
		Weights: engine.StrategyWeights{Defensive: 100},
	}

	action, _ := agent.DecideTurn(input)
	if action.Type != engine.ActionMove {
		t.Fatalf("expected MOVE, got %s", action.Type)
	}
	if action.Target.X != 5 || action.Target.Y != 3 {
		t.Fatalf("expected (5,3), got (%d,%d)", action.Target.X, action.Target.Y)
	}
}

func TestDecideTurn_InvalidFallbackWait(t *testing.T) {
	mock := NewMockLLM([]string{"invalid response"})
	agent := NewAgent("Berserker", Berserker.Description, mock)

	input := AgentDecisionInput{
		Self:    engine.AgentState{Position: engine.Position{X: 3, Y: 3}},
		Enemy:   engine.AgentState{Position: engine.Position{X: 4, Y: 3}},
		Weights: engine.StrategyWeights{},
	}

	action, _ := agent.DecideTurn(input)
	if action.Type != engine.ActionWait {
		t.Fatalf("invalid LLM response should fallback to WAIT, got %s", action.Type)
	}
}

func TestDecideTurn_APIErrorFallbackWait(t *testing.T) {
	failLLM := &FailLLMClient{}
	agent := NewAgent("Berserker", Berserker.Description, failLLM)

	input := AgentDecisionInput{
		Self:    engine.AgentState{Position: engine.Position{X: 3, Y: 3}},
		Enemy:   engine.AgentState{Position: engine.Position{X: 4, Y: 3}},
		Weights: engine.StrategyWeights{},
	}

	action, _ := agent.DecideTurn(input)
	if action.Type != engine.ActionWait {
		t.Fatalf("API error should fallback to WAIT, got %s", action.Type)
	}
}

func TestDecideTurn_IllegalMoveMarkedFailed(t *testing.T) {
	// LLM 返回了合法的格式，但动作不合法（攻击距离太远）
	mock := NewMockLLM([]string{"ATTACK"})
	agent := NewAgent("Berserker", Berserker.Description, mock)

	// 对手在距离 5，ATTACK 不合法
	input := AgentDecisionInput{
		Self:    engine.AgentState{Position: engine.Position{X: 0, Y: 0}},
		Enemy:   engine.AgentState{Position: engine.Position{X: 5, Y: 0}},
		Weights: engine.StrategyWeights{},
	}

	action, _ := agent.DecideTurn(input)
	if action.Type != engine.ActionAttack {
		t.Fatalf("expected ATTACK (marked Failed), got %s", action.Type)
	}
	if !action.Failed {
		t.Fatal("expected Failed=true")
	}
	if action.FailReason == "" {
		t.Fatal("expected non-empty FailReason")
	}
}

func TestDecideTurn_MoveToObstacleMarkedFailed(t *testing.T) {
	mock := NewMockLLM([]string{"MOVE(3,3)"})
	agent := NewAgent("Berserker", Berserker.Description, mock)

	input := AgentDecisionInput{
		Self:      engine.AgentState{Position: engine.Position{X: 1, Y: 3}},
		Enemy:     engine.AgentState{Position: engine.Position{X: 8, Y: 8}},
		Obstacles: []engine.Position{{X: 3, Y: 3}},
		Weights:   engine.StrategyWeights{},
	}

	action, _ := agent.DecideTurn(input)
	if action.Type != engine.ActionMove {
		t.Fatalf("expected MOVE (marked Failed), got %s", action.Type)
	}
	if !action.Failed {
		t.Fatal("expected Failed=true")
	}
	if action.FailReason == "" {
		t.Fatal("expected non-empty FailReason")
	}
}

func TestValidateAgentAction_Move(t *testing.T) {
	self := engine.AgentState{Position: engine.Position{X: 3, Y: 3}}
	enemy := engine.AgentState{Position: engine.Position{X: 8, Y: 8}}

	// 合法移动
	if !validateAgentAction(engine.Action{Type: engine.ActionMove, Target: engine.Position{X: 4, Y: 3}}, self, enemy, nil) {
		t.Error("valid move should pass")
	}

	// 超出移动范围
	if validateAgentAction(engine.Action{Type: engine.ActionMove, Target: engine.Position{X: 6, Y: 3}}, self, enemy, nil) {
		t.Error("move distance > 2 should fail")
	}

	// 移动到障碍物
	if validateAgentAction(engine.Action{Type: engine.ActionMove, Target: engine.Position{X: 4, Y: 4}}, self, enemy, []engine.Position{{X: 4, Y: 4}}) {
		t.Error("move to obstacle should fail")
	}

	// 移动到敌人位置
	if validateAgentAction(engine.Action{Type: engine.ActionMove, Target: engine.Position{X: 8, Y: 8}}, self, enemy, nil) {
		t.Error("move to enemy should fail")
	}
}

func TestValidateAgentAction_Attack(t *testing.T) {
	self := engine.AgentState{Position: engine.Position{X: 3, Y: 3}}

	// 距离 1 → 合法
	nearEnemy := engine.AgentState{Position: engine.Position{X: 4, Y: 3}}
	if !validateAgentAction(engine.Action{Type: engine.ActionAttack}, self, nearEnemy, nil) {
		t.Error("attack at distance 1 should pass")
	}

	// 距离 3 → 不合法
	farEnemy := engine.AgentState{Position: engine.Position{X: 6, Y: 3}}
	if validateAgentAction(engine.Action{Type: engine.ActionAttack}, self, farEnemy, nil) {
		t.Error("attack at distance 3 should fail")
	}
}

func TestValidateAgentAction_Skill(t *testing.T) {
	self := engine.AgentState{Position: engine.Position{X: 3, Y: 3}, SkillCooldown: 0}
	enemy := engine.AgentState{Position: engine.Position{X: 6, Y: 3}}

	// 距离 3 + 冷却 0 → 合法
	if !validateAgentAction(engine.Action{Type: engine.ActionSkill}, self, enemy, nil) {
		t.Error("skill at distance 3 with no cooldown should pass")
	}

	// 距离 3 + 冷却中 → 不合法
	self.SkillCooldown = 2
	if validateAgentAction(engine.Action{Type: engine.ActionSkill}, self, enemy, nil) {
		t.Error("skill on cooldown should fail")
	}
}

func TestDecideTurnsParallel(t *testing.T) {
	redMock := NewMockLLM([]string{"ATTACK"})
	blueMock := NewMockLLM([]string{"MOVE(6,3)"})

	red := NewAgent("Berserker", Berserker.Description, redMock)
	blue := NewAgent("Tactician", Tactician.Description, blueMock)

	redInput := AgentDecisionInput{
		Self:    engine.AgentState{Position: engine.Position{X: 3, Y: 3}},
		Enemy:   engine.AgentState{Position: engine.Position{X: 4, Y: 3}},
		Weights: engine.StrategyWeights{Aggressive: 100},
	}
	blueInput := AgentDecisionInput{
		Self:    engine.AgentState{Position: engine.Position{X: 4, Y: 3}},
		Enemy:   engine.AgentState{Position: engine.Position{X: 3, Y: 3}},
		Weights: engine.StrategyWeights{Defensive: 100},
	}

	redResult, blueResult := DecideTurnsParallel(red, blue, redInput, blueInput)

	if redResult.Action.Type != engine.ActionAttack {
		t.Fatalf("expected red ATTACK, got %s", redResult.Action.Type)
	}
	if blueResult.Action.Type != engine.ActionMove {
		t.Fatalf("expected blue MOVE, got %s", blueResult.Action.Type)
	}

	// 验证两边都被调用了
	if redMock.callCount != 1 || blueMock.callCount != 1 {
		t.Fatalf("expected 1 call each, got red=%d blue=%d", redMock.callCount, blueMock.callCount)
	}
}

func TestGenerateMonologue(t *testing.T) {
	mock := NewMockLLM([]string{"来吧，战个痛快！"})
	agent := NewAgent("Berserker", Berserker.Description, mock)

	input := AgentDecisionInput{
		Self:  engine.AgentState{Name: "Berserker", HP: 90},
		Enemy: engine.AgentState{Name: "Tactician", HP: 85},
	}

	ch := agent.GenerateMonologue(input, engine.Action{Type: engine.ActionAttack})
	result := <-ch

	if result != "来吧，战个痛快！" {
		t.Fatalf("unexpected monologue: %s", result)
	}
}

func TestGenerateMonologue_Error(t *testing.T) {
	failLLM := &FailLLMClient{}
	agent := NewAgent("Berserker", Berserker.Description, failLLM)

	input := AgentDecisionInput{
		Self:  engine.AgentState{HP: 90},
		Enemy: engine.AgentState{HP: 85},
	}

	ch := agent.GenerateMonologue(input, engine.Action{Type: engine.ActionAttack})
	result := <-ch

	if result != "..." {
		t.Fatalf("error monologue should be ..., got %s", result)
	}
}

func TestGetPersonality(t *testing.T) {
	p := GetPersonality("Berserker")
	if p == nil {
		t.Fatal("Berserker not found")
	}
	if p.Name != "Berserker" {
		t.Fatalf("expected Berserker, got %s", p.Name)
	}
	if p.Detail == "" {
		t.Fatal("Detail should not be empty")
	}

	// 不存在的性格
	p = GetPersonality("Unknown")
	if p != nil {
		t.Fatal("Unknown should return nil")
	}
}

func TestAllPersonalities(t *testing.T) {
	if len(AllPersonalities) != 4 {
		t.Fatalf("expected 4 personalities, got %d", len(AllPersonalities))
	}

	names := map[string]bool{}
	for _, p := range AllPersonalities {
		if p.Description == "" {
			t.Fatalf("%s has empty Description", p.Name)
		}
		if p.Detail == "" {
			t.Fatalf("%s has empty Detail", p.Name)
		}
		names[p.Name] = true
	}

	expected := []string{"Berserker", "Tactician", "Trickster", "Defender"}
	for _, name := range expected {
		if !names[name] {
			t.Fatalf("missing personality: %s", name)
		}
	}
}

// TestFullGameWithMockLLM 用 MockLLM 跑一场完整对局
func TestFullGameWithMockLLM(t *testing.T) {
	redEngine := engine.NewEngine()
	redMock := NewMockLLM([]string{"ATTACK", "MOVE(1,0)", "CHARGE", "WAIT", "SKILL"})
	blueMock := NewMockLLM([]string{"MOVE(8,9)", "ATTACK", "SKILL", "WAIT", "CHARGE"})

	red := NewAgent("Berserker", Berserker.Description, redMock)
	blue := NewAgent("Tactician", Tactician.Description, blueMock)

	game := redEngine.NewGame(1, red.Name, blue.Name, "berserker", "tactician")
	game.Status = engine.StatusPlaying

	weights := engine.StrategyWeights{Aggressive: 50, Defensive: 30, Tricky: 20}

	for game.Status == engine.StatusPlaying {
		redInput := AgentDecisionInput{
			Self:         game.AgentRed,
			Enemy:        game.AgentBlue,
			CurrentRound: game.CurrentRound,
			MaxRounds:    game.MaxRounds,
			Obstacles:    game.Obstacles,
			Weights:      weights,
		}
		blueInput := AgentDecisionInput{
			Self:         game.AgentBlue,
			Enemy:        game.AgentRed,
			CurrentRound: game.CurrentRound,
			MaxRounds:    game.MaxRounds,
			Obstacles:    game.Obstacles,
			Weights:      weights,
		}

		redResult, blueResult := DecideTurnsParallel(red, blue, redInput, blueInput)
		redEngine.ExecuteTurn(game, redResult.Action, blueResult.Action)
	}

	if game.Status != engine.StatusFinished {
		t.Fatalf("expected finished, got %s", game.Status)
	}
	if game.CurrentRound > game.MaxRounds {
		t.Fatalf("exceeded max rounds: %d > %d", game.CurrentRound, game.MaxRounds)
	}

	t.Logf("mock game finished: rounds=%d, winner=%s, redHP=%d, blueHP=%d",
		game.CurrentRound, game.Winner, game.AgentRed.HP, game.AgentBlue.HP)
	t.Logf("LLM calls: red=%d, blue=%d", redMock.callCount, blueMock.callCount)
}
