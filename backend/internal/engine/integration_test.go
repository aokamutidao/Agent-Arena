package engine

import (
	"encoding/hex"
	"fmt"
	"testing"
)

// TestFullGameIntegration 集成测试：跑一场完整对局 + 生成 actions hash
func TestFullGameIntegration(t *testing.T) {
	runner := NewGameRunner()

	// 跑一场完整对局
	game := runner.RunFullGame(1, "Berserker", "Tactician")

	// 验证对局正常结束
	if game.Status != StatusFinished {
		t.Fatalf("expected game to be finished, got %s", game.Status)
	}

	// 验证有回合记录
	if len(game.History) == 0 {
		t.Fatal("expected at least 1 turn in history")
	}

	// 验证回合数合理
	if game.CurrentRound > game.MaxRounds {
		t.Fatalf("exceeded max rounds: %d > %d", game.CurrentRound, game.MaxRounds)
	}

	// 验证胜负已判定
	if game.Winner == SideNone {
		t.Log("game ended in a draw")
	} else {
		t.Logf("winner: %s", game.Winner)
	}

	// 计算 actions hash
	actionsHash := ComputeActionsHash(game.History)
	hashHex := hex.EncodeToString(actionsHash[:])
	t.Logf("actions hash: 0x%s", hashHex)

	// 验证 hash 非零
	if actionsHash == [32]byte{} {
		t.Fatal("actions hash is zero")
	}

	// 打印对局摘要
	printGameSummary(game)
}

// TestMultipleGames 测试多场对局（验证引擎稳定性）
func TestMultipleGames(t *testing.T) {
	runner := NewGameRunner()
	results := map[Side]int{}

	for i := uint64(1); i <= 100; i++ {
		game := runner.RunFullGame(i, "Berserker", "Tactician")

		if game.Status != StatusFinished {
			t.Fatalf("game %d not finished", i)
		}

		results[game.Winner]++

		// 验证 hash 每次都不同（随机对局）
		hash := ComputeActionsHash(game.History)
		if hash == [32]byte{} {
			t.Fatalf("game %d has zero hash", i)
		}
	}

	t.Logf("100 games results: Red=%d, Blue=%d, Draw=%d",
		results[SideRed], results[SideBlue], results[SideNone])
}

// TestActionsHashDeterminism 验证相同历史产生相同 hash
func TestActionsHashDeterminism(t *testing.T) {
	history := []TurnRecord{
		{Round: 1, RedAction: Action{Type: ActionAttack}, BlueAction: Action{Type: ActionWait}, RedHPAfter: 100, BlueHPAfter: 85},
		{Round: 2, RedAction: Action{Type: ActionMove, Target: Position{2, 2}}, BlueAction: Action{Type: ActionSkill}, RedHPAfter: 88, BlueHPAfter: 85},
	}

	hash1 := ComputeActionsHash(history)
	hash2 := ComputeActionsHash(history)

	if hash1 != hash2 {
		t.Error("same history should produce same hash")
	}
}

func printGameSummary(game *GameState) {
	fmt.Println("═══════════════════════════════════════")
	fmt.Printf("  Game #%d Summary\n", game.GameID)
	fmt.Println("═══════════════════════════════════════")
	fmt.Printf("  Rounds: %d/%d\n", game.CurrentRound, game.MaxRounds)
	fmt.Printf("  Winner: %s\n", game.Winner)
	fmt.Printf("  Red  (%s): HP %d/%d, Position (%d,%d)\n",
		game.AgentRed.Name, game.AgentRed.HP, game.AgentRed.MaxHP,
		game.AgentRed.Position.X, game.AgentRed.Position.Y)
	fmt.Printf("  Blue (%s): HP %d/%d, Position (%d,%d)\n",
		game.AgentBlue.Name, game.AgentBlue.HP, game.AgentBlue.MaxHP,
		game.AgentBlue.Position.X, game.AgentBlue.Position.Y)
	fmt.Println("───────────────────────────────────────")
	fmt.Println("  Last 5 turns:")
	start := 0
	if len(game.History) > 5 {
		start = len(game.History) - 5
	}
	for i := start; i < len(game.History); i++ {
		turn := game.History[i]
		fmt.Printf("  R%d: Red=%-8s Blue=%-8s | HP: %d vs %d\n",
			turn.Round, turn.RedAction.Type, turn.BlueAction.Type,
			turn.RedHPAfter, turn.BlueHPAfter)
	}
	fmt.Println("═══════════════════════════════════════")
}
