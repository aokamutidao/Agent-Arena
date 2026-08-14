package agent

import (
	"testing"

	"agent-arena/backend/internal/engine"
)

func TestParseAction_Attack(t *testing.T) {
	action, err := ParseAction("ATTACK")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if action.Type != engine.ActionAttack {
		t.Fatalf("expected ATTACK, got %s", action.Type)
	}
}

func TestParseAction_Skill(t *testing.T) {
	action, err := ParseAction("SKILL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if action.Type != engine.ActionSkill {
		t.Fatalf("expected SKILL, got %s", action.Type)
	}
}

func TestParseAction_Charge(t *testing.T) {
	action, err := ParseAction("CHARGE")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if action.Type != engine.ActionCharge {
		t.Fatalf("expected CHARGE, got %s", action.Type)
	}
}

func TestParseAction_Wait(t *testing.T) {
	action, err := ParseAction("WAIT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if action.Type != engine.ActionWait {
		t.Fatalf("expected WAIT, got %s", action.Type)
	}
}

func TestParseAction_Move(t *testing.T) {
	action, err := ParseAction("MOVE(5,3)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if action.Type != engine.ActionMove {
		t.Fatalf("expected MOVE, got %s", action.Type)
	}
	if action.Target.X != 5 || action.Target.Y != 3 {
		t.Fatalf("expected (5,3), got (%d,%d)", action.Target.X, action.Target.Y)
	}
}

func TestParseAction_MoveEdgeCoords(t *testing.T) {
	action, err := ParseAction("MOVE(0,0)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if action.Target.X != 0 || action.Target.Y != 0 {
		t.Fatalf("expected (0,0), got (%d,%d)", action.Target.X, action.Target.Y)
	}

	action, err = ParseAction("MOVE(9,9)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if action.Target.X != 9 || action.Target.Y != 9 {
		t.Fatalf("expected (9,9), got (%d,%d)", action.Target.X, action.Target.Y)
	}
}

func TestParseAction_MoveOutOfRange(t *testing.T) {
	_, err := ParseAction("MOVE(10,5)")
	if err == nil {
		t.Fatal("expected error for out-of-range coordinates")
	}
}

func TestParseAction_InvalidInput(t *testing.T) {
	action, err := ParseAction("invalid")
	if err == nil {
		t.Fatal("expected error for invalid input")
	}
	if action.Type != engine.ActionWait {
		t.Fatalf("invalid input should fallback to WAIT, got %s", action.Type)
	}
}

func TestParseAction_CleanMarkdown(t *testing.T) {
	action, err := ParseAction("```ATTACK```")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if action.Type != engine.ActionAttack {
		t.Fatalf("expected ATTACK after cleaning markdown, got %s", action.Type)
	}
}

func TestParseAction_CleanQuotes(t *testing.T) {
	action, err := ParseAction(`"SKILL"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if action.Type != engine.ActionSkill {
		t.Fatalf("expected SKILL after cleaning quotes, got %s", action.Type)
	}
}

func TestParseAction_EmptyInput(t *testing.T) {
	action, err := ParseAction("")
	if err == nil {
		t.Fatal("expected error for empty input")
	}
	if action.Type != engine.ActionWait {
		t.Fatalf("empty input should fallback to WAIT, got %s", action.Type)
	}
}

func TestCleanResponse(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"ATTACK", "ATTACK"},
		{"  ATTACK  ", "ATTACK"},
		{"\"ATTACK\"", "ATTACK"},
		{"```ATTACK```", "ATTACK"},
		{"`MOVE(3,4)`", "MOVE(3,4)"},
	}

	for _, tt := range tests {
		result := cleanResponse(tt.input)
		if result != tt.expected {
			t.Errorf("cleanResponse(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestFormatStatus_Empty(t *testing.T) {
	result := formatStatus(nil)
	if result != "无" {
		t.Fatalf("expected 无, got %s", result)
	}
}

func TestFormatStatus_WithEffects(t *testing.T) {
	effects := []engine.Effect{
		{Type: engine.EffectCharging, Remaining: 1},
	}
	result := formatStatus(effects)
	if result != "charging(1回合)" {
		t.Fatalf("expected charging(1回合), got %s", result)
	}
}

func TestFormatObstacles(t *testing.T) {
	obstacles := []engine.Position{
		{X: 3, Y: 3},
		{X: 6, Y: 6},
	}
	result := formatObstacles(obstacles)
	if result != "(3,3) (6,6)" {
		t.Fatalf("expected (3,3) (6,6), got %s", result)
	}
}

func TestFormatAction(t *testing.T) {
	tests := []struct {
		action   engine.Action
		expected string
	}{
		{engine.Action{Type: engine.ActionAttack}, "ATTACK"},
		{engine.Action{Type: engine.ActionSkill}, "SKILL"},
		{engine.Action{Type: engine.ActionMove, Target: engine.Position{X: 5, Y: 3}}, "MOVE(5,3)"},
	}

	for _, tt := range tests {
		result := formatAction(tt.action)
		if result != tt.expected {
			t.Errorf("formatAction(%v) = %q, want %q", tt.action, result, tt.expected)
		}
	}
}

func TestBuildDecisionPrompt_NotEmpty(t *testing.T) {
	input := AgentDecisionInput{
		Self: engine.AgentState{
			Name:          "Berserker",
			HP:            100,
			MaxHP:         100,
			Position:      engine.Position{X: 0, Y: 0},
			SkillCooldown: 0,
		},
		Enemy: engine.AgentState{
			Name:          "Tactician",
			HP:            85,
			MaxHP:         100,
			Position:      engine.Position{X: 3, Y: 3},
			SkillCooldown: 0,
		},
		CurrentRound: 1,
		MaxRounds:    30,
		Obstacles:    engine.DefaultObstacles,
		Weights:      engine.StrategyWeights{Aggressive: 60, Defensive: 20, Tricky: 20},
	}

	prompt := buildDecisionPrompt("Berserker", "你是一个狂战士", input)
	if prompt == "" {
		t.Fatal("prompt should not be empty")
	}
	if len(prompt) < 200 {
		t.Fatalf("prompt too short: %d chars", len(prompt))
	}
}

func TestBuildMonologuePrompt_NotEmpty(t *testing.T) {
	input := AgentDecisionInput{
		Self:  engine.AgentState{Name: "Berserker", HP: 90},
		Enemy: engine.AgentState{Name: "Tactician", HP: 85},
	}

	prompt := buildMonologuePrompt("Berserker", "狂战士", input, engine.Action{Type: engine.ActionAttack})
	if prompt == "" {
		t.Fatal("monologue prompt should not be empty")
	}
}

// TestParseAction_HealVariants 测试 HEAL 大小写变体（LLM 常输出 "HEAl" 小写 L）
func TestParseAction_HealVariants(t *testing.T) {
	cases := []string{"HEAL", "heal", "Heal", "HEAl", "heal "}
	for _, c := range cases {
		action, err := ParseAction(c)
		if err != nil {
			t.Fatalf("ParseAction(%q) error: %v", c, err)
		}
		if action.Type != engine.ActionHeal {
			t.Fatalf("ParseAction(%q) = %s, want HEAL", c, action.Type)
		}
	}
}

// TestParseAction_CaseInsensitive 测试所有动作的大小写容错
func TestParseAction_CaseInsensitive(t *testing.T) {
	cases := map[string]engine.ActionType{
		"attack": engine.ActionAttack,
		"Attack": engine.ActionAttack,
		"SKILL":  engine.ActionSkill,
		"skill":  engine.ActionSkill,
		"charge": engine.ActionCharge,
		"wait":   engine.ActionWait,
		"move(3,4)": engine.ActionMove,
	}
	for input, want := range cases {
		action, err := ParseAction(input)
		if err != nil {
			t.Fatalf("ParseAction(%q) error: %v", input, err)
		}
		if action.Type != want {
			t.Fatalf("ParseAction(%q) = %s, want %s", input, action.Type, want)
		}
	}
}
