package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"agent-arena/backend/internal/engine"
)

func TestLegacyWritesAreRetired(t *testing.T) {
	server := setupTestServer(t)
	for _, path := range []string{
		"/api/auth/login", "/api/games/create", "/api/games/1/start",
		"/api/bets/place", "/api/strategy/vote", "/api/auth/challenges",
		"/api/auth/agents/test-api", "/api/auth/claim-daily",
	} {
		t.Run(path, func(t *testing.T) {
			w := httptest.NewRecorder()
			server.Router().ServeHTTP(w, httptest.NewRequest(http.MethodPost, path, nil))
			if w.Code != http.StatusGone {
				t.Fatalf("retired endpoint returned %d, want 410", w.Code)
			}
		})
	}
}

func TestEvidenceHashIncludesActionsWithSameOutcome(t *testing.T) {
	state := &engine.GameState{
		GameID: 1, CurrentRound: 1,
		History: []engine.TurnRecord{{Round: 1, RedAction: engine.Action{Type: engine.ActionWait}}},
	}
	before := computeActionsHash(state)
	// The old implementation hashed only ID, round count and HP: this alteration was invisible.
	state.History[0].RedAction.Type = engine.ActionCharge
	if before == computeActionsHash(state) {
		t.Fatal("altering actions without changing final HP must change the evidence hash")
	}
}
