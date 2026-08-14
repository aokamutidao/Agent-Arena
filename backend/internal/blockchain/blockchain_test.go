package blockchain

import (
	"context"
	"math/big"
	"testing"
	"time"

	"agent-arena/backend/internal/engine"
)

func TestMockChainService_CreateGame(t *testing.T) {
	mock := NewMockChainService()
	ctx := context.Background()

	txHash, gameID, err := mock.CreateGame(ctx, "berserker", "tactician", 120)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gameID != 1 {
		t.Fatalf("expected gameID 1, got %d", gameID)
	}
	if txHash == "" {
		t.Fatal("tx hash should not be empty")
	}

	// 验证状态
	status, err := mock.GetGameStatus(ctx, gameID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != "open" {
		t.Fatalf("expected open, got %s", status)
	}
}

func TestMockChainService_StartGame(t *testing.T) {
	mock := NewMockChainService()
	ctx := context.Background()

	_, gameID, _ := mock.CreateGame(ctx, "berserker", "tactician", 120)

	txHash, err := mock.StartGame(ctx, gameID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if txHash == "" {
		t.Fatal("tx hash should not be empty")
	}

	status, _ := mock.GetGameStatus(ctx, gameID)
	if status != "locked" {
		t.Fatalf("expected locked, got %s", status)
	}
}

func TestMockChainService_StartGame_NotOpen(t *testing.T) {
	mock := NewMockChainService()
	ctx := context.Background()

	_, gameID, _ := mock.CreateGame(ctx, "berserker", "tactician", 120)
	mock.StartGame(ctx, gameID)

	_, err := mock.StartGame(ctx, gameID)
	if err == nil {
		t.Fatal("expected error for starting a locked game")
	}
}

func TestMockChainService_FinishGame(t *testing.T) {
	mock := NewMockChainService()
	ctx := context.Background()

	_, gameID, _ := mock.CreateGame(ctx, "berserker", "tactician", 120)
	mock.StartGame(ctx, gameID)

	hash := [32]byte{1, 2, 3}
	txHash, err := mock.FinishGame(ctx, gameID, true, hash)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if txHash == "" {
		t.Fatal("tx hash should not be empty")
	}

	status, _ := mock.GetGameStatus(ctx, gameID)
	if status != "finished" {
		t.Fatalf("expected finished, got %s", status)
	}
}

func TestMockChainService_GetGamePool(t *testing.T) {
	mock := NewMockChainService()
	ctx := context.Background()

	_, gameID, _ := mock.CreateGame(ctx, "berserker", "tactician", 120)

	mock.AddBet(gameID, "red", big.NewInt(50000000))
	mock.AddBet(gameID, "blue", big.NewInt(30000000))

	totalRed, totalBlue, err := mock.GetGamePool(ctx, gameID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if totalRed.Cmp(big.NewInt(50000000)) != 0 {
		t.Fatalf("expected 50000000, got %s", totalRed)
	}
	if totalBlue.Cmp(big.NewInt(30000000)) != 0 {
		t.Fatalf("expected 30000000, got %s", totalBlue)
	}
}

func TestMockChainService_GetStrategyWeights_Default(t *testing.T) {
	mock := NewMockChainService()
	ctx := context.Background()

	weights, err := mock.GetStrategyWeights(ctx, 999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if weights.Aggressive != 33 || weights.Defensive != 33 || weights.Tricky != 34 {
		t.Fatalf("expected default weights 33/33/34, got %d/%d/%d",
			weights.Aggressive, weights.Defensive, weights.Tricky)
	}
}

func TestMockChainService_GetStrategyWeights_Custom(t *testing.T) {
	mock := NewMockChainService()
	ctx := context.Background()

	custom := &engine.StrategyWeights{Aggressive: 60, Defensive: 20, Tricky: 20}
	mock.SetStrategyWeights(1, custom)

	weights, err := mock.GetStrategyWeights(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if weights.Aggressive != 60 {
		t.Fatalf("expected 60, got %d", weights.Aggressive)
	}
}

func TestMockChainService_GetOdds(t *testing.T) {
	mock := NewMockChainService()
	ctx := context.Background()

	_, gameID, _ := mock.CreateGame(ctx, "berserker", "tactician", 120)

	// 无下注 → 默认赔率
	oddsRed, oddsBlue, _ := mock.GetOdds(ctx, gameID)
	if oddsRed.Cmp(big.NewInt(2)) != 0 || oddsBlue.Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("expected 2/2, got %s/%s", oddsRed, oddsBlue)
	}

	// 有下注
	mock.AddBet(gameID, "red", big.NewInt(150))
	mock.AddBet(gameID, "blue", big.NewInt(50))

	oddsRed, oddsBlue, _ = mock.GetOdds(ctx, gameID)
	// total=200, redOdds=200/150=1, blueOdds=200/50=4
	if oddsRed.Cmp(big.NewInt(0)) <= 0 || oddsBlue.Cmp(big.NewInt(0)) <= 0 {
		t.Fatalf("odds should be positive: red=%s blue=%s", oddsRed, oddsBlue)
	}
}

func TestMockChainService_ListenEvents(t *testing.T) {
	mock := NewMockChainService()
	ctx := context.Background()

	ch, err := mock.ListenEvents(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 创建对局应产生事件
	mock.CreateGame(ctx, "berserker", "tactician", 120)

	select {
	case event := <-ch:
		if event.Type != EventGameCreated {
			t.Fatalf("expected GameCreated, got %s", event.Type)
		}
		if event.GameID != 1 {
			t.Fatalf("expected gameID 1, got %d", event.GameID)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for event")
	}
}

func TestMockChainService_EventFlow(t *testing.T) {
	mock := NewMockChainService()
	ctx := context.Background()

	ch, _ := mock.ListenEvents(ctx)

	// 完整流程
	_, gameID, _ := mock.CreateGame(ctx, "berserker", "tactician", 120)
	mock.AddBet(gameID, "red", big.NewInt(100))
	mock.StartGame(ctx, gameID)
	mock.FinishGame(ctx, gameID, true, [32]byte{})

	// 收集事件
	var events []ContractEvent
	timeout := time.After(2 * time.Second)
	for {
		select {
		case event := <-ch:
			events = append(events, event)
		case <-timeout:
			goto done
		}
	}
done:

	if len(events) < 4 {
		t.Fatalf("expected at least 4 events, got %d", len(events))
	}

	// 验证事件顺序
	expectedTypes := []EventType{EventGameCreated, EventBetPlaced, EventBettingLocked, EventGameSettled}
	for i, expected := range expectedTypes {
		if i >= len(events) {
			break
		}
		if events[i].Type != expected {
			t.Fatalf("event %d: expected %s, got %s", i, expected, events[i].Type)
		}
	}
}

func TestMockChainService_NonExistentGame(t *testing.T) {
	mock := NewMockChainService()
	ctx := context.Background()

	_, err := mock.GetGameStatus(ctx, 999)
	if err == nil {
		t.Fatal("expected error for non-existent game")
	}

	_, err = mock.StartGame(ctx, 999)
	if err == nil {
		t.Fatal("expected error for non-existent game")
	}

	_, err = mock.FinishGame(ctx, 999, true, [32]byte{})
	if err == nil {
		t.Fatal("expected error for non-existent game")
	}
}
