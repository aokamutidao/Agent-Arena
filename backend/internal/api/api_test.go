package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"agent-arena/backend/internal/engine"
)

func setupTestServer() *Server {
	server := NewServer()

	// 注册测试 Agent
	server.store.RegisterAgent(AgentInfo{
		ID:          "berserker",
		Name:        "Berserker",
		Personality: "berserker",
		Wins:        10,
		Losses:      5,
		WinRate:     66.7,
		Description: "激进攻型，喜欢近战冲锋",
	})
	server.store.RegisterAgent(AgentInfo{
		ID:          "tactician",
		Name:        "Tactician",
		Personality: "tactician",
		Wins:        8,
		Losses:      7,
		WinRate:     53.3,
		Description: "稳扎稳打，远程消耗",
	})

	// 创建测试对局
	eng := engine.NewEngine()
	game := eng.NewGame(1, "Berserker", "Tactician", "berserker", "tactician")
	game.Status = engine.StatusBetting
	server.store.CreateGame(game)

	// 创建第二场对局（进行中）
	game2 := eng.NewGame(2, "Berserker", "Tactician", "berserker", "tactician")
	game2.Status = engine.StatusPlaying
	game2.CurrentRound = 5
	game2.AgentRed.HP = 85
	game2.AgentBlue.HP = 70
	server.store.CreateGame(game2)
	server.store.UpdateBetPool(2, "red", 50000000)
	server.store.UpdateBetPool(2, "blue", 30000000)

	return server
}

func TestHealthCheck(t *testing.T) {
	server := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestListGames(t *testing.T) {
	server := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/games", nil)
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp GamesResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Total != 2 {
		t.Fatalf("expected 2 games, got %d", resp.Total)
	}
	if len(resp.Games) != 2 {
		t.Fatalf("expected 2 game items, got %d", len(resp.Games))
	}
}

func TestListGames_FilterByStatus(t *testing.T) {
	server := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/games?status=betting", nil)
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp GamesResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Total != 1 {
		t.Fatalf("expected 1 betting game, got %d", resp.Total)
	}
	if resp.Games[0].Status != "betting" {
		t.Fatalf("expected betting status, got %s", resp.Games[0].Status)
	}
}

func TestGetGame_Found(t *testing.T) {
	server := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/games/2", nil)
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp GameDetailResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.GameID != 2 {
		t.Fatalf("expected game_id 2, got %d", resp.GameID)
	}
	if resp.CurrentRound != 5 {
		t.Fatalf("expected round 5, got %d", resp.CurrentRound)
	}
	if resp.AgentRedState == nil {
		t.Fatal("expected agent_red_state to be present")
	}
	if resp.AgentRedState.HP != 85 {
		t.Fatalf("expected red HP 85, got %d", resp.AgentRedState.HP)
	}
	if resp.TotalBetRed != "50000000" {
		t.Fatalf("expected total_bet_red 50000000, got %s", resp.TotalBetRed)
	}
}

func TestGetGame_NotFound(t *testing.T) {
	server := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/games/999", nil)
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestGetGame_InvalidID(t *testing.T) {
	server := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/games/abc", nil)
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestEstimateBet(t *testing.T) {
	server := setupTestServer()

	body := EstimateRequest{
		GameID: 2,
		Side:   "red",
		Amount: "10000000",
	}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/bets/estimate", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp EstimateResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.CurrentPoolRed != "50000000" {
		t.Fatalf("expected current pool red 50000000, got %s", resp.CurrentPoolRed)
	}
	if resp.NewPoolRed != "60000000" {
		t.Fatalf("expected new pool red 60000000, got %s", resp.NewPoolRed)
	}
	if resp.NewOddsRed <= 0 || resp.NewOddsBlue <= 0 {
		t.Fatalf("odds should be positive: red=%f blue=%f", resp.NewOddsRed, resp.NewOddsBlue)
	}
}

func TestEstimateBet_InvalidBody(t *testing.T) {
	server := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/bets/estimate", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestEstimateBet_GameNotFound(t *testing.T) {
	server := setupTestServer()

	body := EstimateRequest{GameID: 999, Side: "red", Amount: "10000000"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/bets/estimate", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestListAgents(t *testing.T) {
	server := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/agents", nil)
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp AgentsResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	if len(resp.Agents) != 2 {
		t.Fatalf("expected 2 agents, got %d", len(resp.Agents))
	}
}

func TestGetAgent_Found(t *testing.T) {
	server := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/agents/berserker", nil)
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp AgentInfo
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Name != "Berserker" {
		t.Fatalf("expected Berserker, got %s", resp.Name)
	}
	if resp.Wins != 10 {
		t.Fatalf("expected 10 wins, got %d", resp.Wins)
	}
}

func TestGetAgent_NotFound(t *testing.T) {
	server := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/agents/unknown", nil)
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestGetUserBets_Empty(t *testing.T) {
	server := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/users/0xabc/bets", nil)
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	bets := resp["bets"].([]interface{})
	if len(bets) != 0 {
		t.Fatalf("expected 0 bets, got %d", len(bets))
	}
}

func TestGetUserBets_WithBets(t *testing.T) {
	server := setupTestServer()

	server.store.CacheBet("0xabc", &BetRecord{
		GameID:   2,
		Side:     "red",
		Amount:   "10000000",
		Strategy: "aggressive",
		Status:   "pending",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/users/0xabc/bets", nil)
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	bets := resp["bets"].([]interface{})
	if len(bets) != 1 {
		t.Fatalf("expected 1 bet, got %d", len(bets))
	}
}

func TestCalculateOdds(t *testing.T) {
	// 无下注 → 默认赔率
	oddsRed, oddsBlue := CalculateOdds(0, 0)
	if oddsRed != 2.0 || oddsBlue != 2.0 {
		t.Fatalf("expected 2.0/2.0, got %f/%f", oddsRed, oddsBlue)
	}

	// 均等下注 → 赔率接近 2.0
	oddsRed, oddsBlue = CalculateOdds(100, 100)
	if oddsRed != 2.0 || oddsBlue != 2.0 {
		t.Fatalf("expected 2.0/2.0, got %f/%f", oddsRed, oddsBlue)
	}

	// 红方多 → 红方赔率低
	oddsRed, oddsBlue = CalculateOdds(150, 50)
	if oddsRed >= oddsBlue {
		t.Fatalf("red odds should be lower: red=%f blue=%f", oddsRed, oddsBlue)
	}
}

func TestMemoryStore_UpdateAgentStats(t *testing.T) {
	store := NewMemoryStore()
	store.RegisterAgent(AgentInfo{ID: "test", Name: "Test"})

	store.UpdateAgentStats("test", true)
	store.UpdateAgentStats("test", true)
	store.UpdateAgentStats("test", false)

	agent, _ := store.GetAgent("test")
	if agent.Info.Wins != 2 {
		t.Fatalf("expected 2 wins, got %d", agent.Info.Wins)
	}
	if agent.Info.Losses != 1 {
		t.Fatalf("expected 1 loss, got %d", agent.Info.Losses)
	}
	if agent.Info.WinRate < 66.0 || agent.Info.WinRate > 67.0 {
		t.Fatalf("expected ~66.7%% winrate, got %f", agent.Info.WinRate)
	}
}
