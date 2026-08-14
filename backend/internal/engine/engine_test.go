package engine

import (
	"testing"
)

func TestNewGame(t *testing.T) {
	engine := NewEngine()
	game := engine.NewGame(1, "Berserker", "Tactician", "berserker", "tactician")

	if game.Status != StatusPending {
		t.Errorf("expected status Pending, got %s", game.Status)
	}
	if game.AgentRed.HP != MaxHP {
		t.Errorf("expected Red HP %d, got %d", MaxHP, game.AgentRed.HP)
	}
	if game.AgentBlue.HP != MaxHP {
		t.Errorf("expected Blue HP %d, got %d", MaxHP, game.AgentBlue.HP)
	}
	if game.AgentRed.Position != RedSpawn {
		t.Errorf("expected Red at %v, got %v", RedSpawn, game.AgentRed.Position)
	}
	if game.AgentBlue.Position != BlueSpawn {
		t.Errorf("expected Blue at %v, got %v", BlueSpawn, game.AgentBlue.Position)
	}
}

func TestManhattanDistance(t *testing.T) {
	tests := []struct {
		a, b     Position
		expected uint8
	}{
		{Position{0, 0}, Position{1, 1}, 2},
		{Position{3, 3}, Position{4, 3}, 1},
		{Position{1, 1}, Position{8, 8}, 14},
		{Position{5, 5}, Position{5, 5}, 0},
	}

	for _, tc := range tests {
		result := ManhattanDistance(tc.a, tc.b)
		if result != tc.expected {
			t.Errorf("ManhattanDistance(%v, %v) = %d, expected %d", tc.a, tc.b, result, tc.expected)
		}
	}
}

func TestValidateAction(t *testing.T) {
	engine := NewEngine()
	self := AgentState{Position: Position{3, 3}, HP: 80, MaxHP: MaxHP}
	enemy := AgentState{Position: Position{4, 3}, HP: 60, MaxHP: MaxHP}
	obstacles := []Position{{2, 3}}

	// 合法 ATTACK（距离 1）
	err := engine.ValidateAction(Action{Type: ActionAttack}, self, enemy, obstacles)
	if err != nil {
		t.Errorf("expected ATTACK to be valid, got error: %v", err)
	}

	// 非法 ATTACK（距离太远）
	farEnemy := AgentState{Position: Position{7, 7}}
	err = engine.ValidateAction(Action{Type: ActionAttack}, self, farEnemy, obstacles)
	if err == nil {
		t.Error("expected ATTACK to be invalid at distance 8")
	}

	// 合法 MOVE
	err = engine.ValidateAction(Action{Type: ActionMove, Target: Position{4, 4}}, self, enemy, obstacles)
	if err != nil {
		t.Errorf("expected MOVE to be valid, got error: %v", err)
	}

	// 非法 MOVE（障碍物）
	err = engine.ValidateAction(Action{Type: ActionMove, Target: Position{2, 3}}, self, enemy, obstacles)
	if err == nil {
		t.Error("expected MOVE to obstacle to be invalid")
	}

	// 非法 MOVE（超出移动力）
	err = engine.ValidateAction(Action{Type: ActionMove, Target: Position{7, 7}}, self, enemy, obstacles)
	if err == nil {
		t.Error("expected MOVE beyond range to be invalid")
	}

	// CHARGE 和 WAIT 总是合法
	err = engine.ValidateAction(Action{Type: ActionCharge}, self, enemy, obstacles)
	if err != nil {
		t.Errorf("expected CHARGE to be valid, got error: %v", err)
	}

	err = engine.ValidateAction(Action{Type: ActionWait}, self, enemy, obstacles)
	if err != nil {
		t.Errorf("expected WAIT to be valid, got error: %v", err)
	}
}

func TestExecuteTurn_Attack(t *testing.T) {
	engine := NewEngine()
	game := engine.NewGame(1, "Red", "Blue", "berserker", "tactician")
	game.Status = StatusPlaying

	// 手动设置近距离
	game.AgentRed.Position = Position{3, 3}
	game.AgentBlue.Position = Position{4, 3} // 距离 1

	redAction := Action{Type: ActionAttack}
	blueAction := Action{Type: ActionAttack}

	record, err := engine.ExecuteTurn(game, redAction, blueAction)
	if err != nil {
		t.Fatalf("ExecuteTurn error: %v", err)
	}

	if record.Round != 1 {
		t.Errorf("expected round 1, got %d", record.Round)
	}

	// 红方先手，先攻击蓝方
	// Red attacks Blue: Blue HP should be reduced
	if game.AgentBlue.HP >= MaxHP {
		t.Errorf("expected Blue HP to be reduced, got %d", game.AgentBlue.HP)
	}
}

func TestExecuteTurn_ChargeAndAttack(t *testing.T) {
	engine := NewEngine()
	game := engine.NewGame(1, "Red", "Blue", "berserker", "tactician")
	game.Status = StatusPlaying

	game.AgentRed.Position = Position{3, 3}
	game.AgentBlue.Position = Position{4, 3}

	// Round 1: Red charges
	record, _ := engine.ExecuteTurn(game, Action{Type: ActionCharge}, Action{Type: ActionWait})
	if !game.AgentRed.IsCharging {
		t.Error("expected Red to be charging after CHARGE action")
	}
	if record.RedAction.Type != ActionCharge {
		t.Errorf("expected red action CHARGE, got %s", record.RedAction.Type)
	}

	// Round 2: Red attacks with charge bonus
	blueHPBefore := game.AgentBlue.HP
	engine.ExecuteTurn(game, Action{Type: ActionAttack}, Action{Type: ActionWait})

	damage := blueHPBefore - game.AgentBlue.HP
	atkF := float64(ATK)
	expected := uint8(int(atkF * float64(ChargeMulti)))
	if damage != expected {
		t.Errorf("expected charged damage %d, got %d", expected, damage)
	}
}

// TestExecuteTurn_ChargeMoveAttack 蓄力 → MOVE 调整位置 → ATTACK 释放
// 验证蓄力在 MOVE 后不消失，且可以在第 3 回合释放
func TestExecuteTurn_ChargeMoveAttack(t *testing.T) {
	engine := NewEngine()
	game := engine.NewGame(1, "Red", "Blue", "berserker", "tactician")
	game.Status = StatusPlaying
	game.Obstacles = []Position{} // 清空障碍物

	game.AgentRed.Position = Position{1, 1}
	game.AgentBlue.Position = Position{5, 1} // 距离 4，ATTACK 打不到

	// Round 1: Red 蓄力
	engine.ExecuteTurn(game, Action{Type: ActionCharge}, Action{Type: ActionWait})
	if !game.AgentRed.IsCharging {
		t.Fatal("Round 1: Red should be charging after CHARGE")
	}

	// Round 2: Red MOVE 靠近（距离 4 → 2），蓄力应保留
	engine.ExecuteTurn(game, Action{Type: ActionMove, Target: Position{3, 1}}, Action{Type: ActionWait})
	if !game.AgentRed.IsCharging {
		t.Fatal("Round 2: Red should STILL be charging after MOVE (charge persists through move)")
	}

	// Round 3: Red ATTACK（距离 2 → 需要先靠近到 1）
	// 当前 Red 在 (3,1)，Blue 在 (5,1)，距离 2，先 MOVE 到 (4,1) 再 ATTACK
	// 但一回合只能一个动作，所以 Round 3 MOVE，Round 4 ATTACK
	engine.ExecuteTurn(game, Action{Type: ActionMove, Target: Position{4, 1}}, Action{Type: ActionWait})
	if game.AgentRed.IsCharging {
		// 还剩 1 回合剩余，应该还是 charging
		// EffectCharging Remaining: 起始 3，R1 末=2，R2 末=1，R3 末=0 → 清除
		// 所以 R3 末 IsCharging 应该变 false（因为效果过期）
	}

	// 重新设计：CHARGE 在 R1，MOVE 在 R2，R3 时 Remaining=1，R3 末 Remaining=0 过期
	// 所以 R3 仍可 ATTACK 释放蓄力。验证这个场景：
	game2 := engine.NewGame(2, "Red", "Blue", "berserker", "tactician")
	game2.Status = StatusPlaying
	game2.Obstacles = []Position{}

	game2.AgentRed.Position = Position{2, 1}
	game2.AgentBlue.Position = Position{4, 1} // 距离 2

	// R1: CHARGE
	engine.ExecuteTurn(game2, Action{Type: ActionCharge}, Action{Type: ActionWait})
	if !game2.AgentRed.IsCharging {
		t.Fatal("R1: should be charging")
	}

	// R2: MOVE 靠近 (2,1) → (3,1)，距离变为 1
	engine.ExecuteTurn(game2, Action{Type: ActionMove, Target: Position{3, 1}}, Action{Type: ActionWait})
	if !game2.AgentRed.IsCharging {
		t.Fatal("R2: charge should persist through MOVE")
	}

	// R3: ATTACK（距离 1），应该享受蓄力加成
	blueHPBefore := game2.AgentBlue.HP
	engine.ExecuteTurn(game2, Action{Type: ActionAttack}, Action{Type: ActionWait})
	damage := blueHPBefore - game2.AgentBlue.HP

	atkF := float64(ATK)
	expected := uint8(int(atkF * float64(ChargeMulti)))
	if damage != expected {
		t.Errorf("R3 charged ATTACK should deal %d, got %d", expected, damage)
	}
	if game2.AgentRed.IsCharging {
		t.Error("R3: charge should be consumed after ATTACK")
	}
}

// TestExecuteTurn_ChargeExpiry 蓄力 3 回合不攻击则过期
func TestExecuteTurn_ChargeExpiry(t *testing.T) {
	engine := NewEngine()
	game := engine.NewGame(1, "Red", "Blue", "berserker", "tactician")
	game.Status = StatusPlaying
	game.Obstacles = []Position{}

	game.AgentRed.Position = Position{1, 1}
	game.AgentBlue.Position = Position{9, 9} // 远离

	// R1: CHARGE
	engine.ExecuteTurn(game, Action{Type: ActionCharge}, Action{Type: ActionWait})
	if !game.AgentRed.IsCharging {
		t.Fatal("R1: should be charging")
	}

	// R2: WAIT（蓄力仍在）
	engine.ExecuteTurn(game, Action{Type: ActionWait}, Action{Type: ActionWait})
	if !game.AgentRed.IsCharging {
		t.Fatal("R2: should still be charging (2 turns to release)")
	}

	// R3: WAIT（蓄力应过期）
	engine.ExecuteTurn(game, Action{Type: ActionWait}, Action{Type: ActionWait})
	if game.AgentRed.IsCharging {
		t.Error("R3: charge should expire after 2 turns of not attacking")
	}
}

func TestWinCondition_KO(t *testing.T) {
	engine := NewEngine()
	game := engine.NewGame(1, "Red", "Blue", "berserker", "tactician")
	game.Status = StatusPlaying

	game.AgentRed.Position = Position{3, 3}
	game.AgentBlue.Position = Position{4, 3}
	game.AgentBlue.HP = 10 // 很低血量

	engine.ExecuteTurn(game, Action{Type: ActionAttack}, Action{Type: ActionWait})

	if game.Status != StatusFinished {
		t.Errorf("expected game to be finished, got %s", game.Status)
	}
	if game.Winner != SideRed {
		t.Errorf("expected Red to win, got %s", game.Winner)
	}
}

func TestWinCondition_MaxRounds(t *testing.T) {
	engine := NewEngine()
	game := engine.NewGame(1, "Red", "Blue", "berserker", "tactician")
	game.Status = StatusPlaying

	// 远离对方，等待到最大回合数
	game.AgentRed.Position = Position{0, 0}
	game.AgentBlue.Position = Position{9, 9}

	for game.Status != StatusFinished {
		engine.ExecuteTurn(game, Action{Type: ActionWait}, Action{Type: ActionWait})
	}

	// HP 相同 → 进入加时赛（MaxRounds += 10），不再平局
	if !game.Overtime {
		t.Error("expected overtime to be triggered")
	}
	// 加时赛每回合双方各受 10 伤害，100 HP / 10 = 10 回合后双方 HP = 0
	// checkWinCondition 先检查 Red HP == 0 → Blue 胜
	if game.Winner != SideBlue {
		t.Errorf("expected Blue to win in overtime, got %s", game.Winner)
	}
	// 总回合数 = 30 (regular) + 10 (overtime) = 40
	if game.CurrentRound != 40 {
		t.Errorf("expected 40 rounds (overtime), got %d", game.CurrentRound)
	}
}

func TestSkillCooldown(t *testing.T) {
	engine := NewEngine()
	game := engine.NewGame(1, "Red", "Blue", "tactician", "berserker")
	game.Status = StatusPlaying

	game.AgentRed.Position = Position{3, 3}
	game.AgentBlue.Position = Position{5, 3} // 距离 2，在技能范围内

	// 使用技能
	engine.ExecuteTurn(game, Action{Type: ActionSkill}, Action{Type: ActionWait})
	if game.AgentRed.SkillCooldown != SkillCooldown {
		t.Errorf("expected cooldown %d, got %d", SkillCooldown, game.AgentRed.SkillCooldown)
	}

	// 冷却中不能使用技能
	err := engine.ValidateAction(Action{Type: ActionSkill}, game.AgentRed, game.AgentBlue, game.Obstacles)
	if err == nil {
		t.Error("expected skill to be on cooldown")
	}

	// 等待冷却
	for i := 0; i < SkillCooldown; i++ {
		engine.ExecuteTurn(game, Action{Type: ActionWait}, Action{Type: ActionWait})
	}

	if game.AgentRed.SkillCooldown != 0 {
		t.Errorf("expected cooldown 0 after %d rounds, got %d", SkillCooldown, game.AgentRed.SkillCooldown)
	}
}
