package tournament

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"agent-arena/backend/internal/agent"
	"agent-arena/backend/internal/engine"
	"github.com/ethereum/go-ethereum/crypto"
)

func llmTournament(t *testing.T, store *Store) *Tournament {
	t.Helper()
	_, red := testKey(t)
	_, blue := testKey(t)
	a, err := store.CreateAgent(red, "Red", "Attack")
	if err != nil {
		t.Fatal(err)
	}
	b, err := store.CreateAgent(blue, "Blue", "Defend")
	if err != nil {
		t.Fatal(err)
	}
	return &Tournament{ID: 44, Status: "running", Participants: []string{red, blue},
		ParticipantAgents: []ParticipantAgent{{red, a.AgentHash}, {blue, b.AgentHash}}}
}

func llmFixture(t *testing.T, handler http.HandlerFunc) *agent.QwenClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	client, err := agent.NewConfiguredQwenClient(agent.QwenConfig{
		APIKey: "synthetic-test-secret", Provider: "fixture", Model: "fixture-model",
		Endpoint: server.URL, Timeout: time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestLLMFailureCannotCompleteOrFabricateEvidence(t *testing.T) {
	for _, test := range []struct {
		name, body, code string
		status           int
		partialSuccess   bool
	}{
		{"provider", "body contains synthetic-test-secret", "authentication_error", 401, false},
		{"parser", `{"choices":[{"message":{"content":"analysis: ATTACK"}}]}`, "parse_error", 200, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			store := testStore(t)
			tournament := llmTournament(t, store)
			var requests atomic.Int64
			client := llmFixture(t, func(w http.ResponseWriter, _ *http.Request) {
				if requests.Add(1) == 1 && test.partialSuccess {
					_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"WAIT"}}]}`))
					return
				}
				w.WriteHeader(test.status)
				_, _ = w.Write([]byte(test.body))
			})
			runner := NewRunner(testConfig(), store, client)
			if _, err := runner.Start(context.Background(), tournament); err != nil {
				t.Fatal(err)
			}
			runner.Wait()
			match, canonical, err := store.Match(tournament.ID)
			if err != nil || match.Status != "interrupted" || match.Mode != "llm" ||
				match.CurrentRound != 0 || len(match.Rounds) != 0 || canonical != "" ||
				match.WinnerAddress != "" || match.EvidenceHash != "" || match.Outcome != "" {
				t.Fatalf("LLM failure created a false round/result: %+v %v", match, err)
			}
			failed, succeeded := 2, 0
			if test.partialSuccess {
				failed, succeeded = 1, 1
			}
			if len(match.LLMCalls) != 2 || match.LLMSummary.Failed != failed ||
				match.LLMSummary.Succeeded != succeeded ||
				match.LLMSummary.LastErrorCode != test.code || !runner.PublicMatch(match).RetryAllowed {
				t.Fatalf("failure was not recorded as a resumable attempt: %+v", match.LLMSummary)
			}
			encoded, _ := json.Marshal(match)
			if strings.Contains(string(encoded), "synthetic-test-secret") {
				t.Fatal("provider body/credentials leaked into public match")
			}
		})
	}
}

func TestLLMRetryRetainsAttemptsAndCompletedRoundProofs(t *testing.T) {
	store := testStore(t)
	tournament := llmTournament(t, store)
	var healthy atomic.Bool
	client := llmFixture(t, func(w http.ResponseWriter, _ *http.Request) {
		if !healthy.Load() {
			w.WriteHeader(429)
			return
		}
		w.Header().Set("X-Request-Id", "local-request")
		_, _ = w.Write([]byte(`{"model":"fixture-model-v2","choices":[{"message":{"content":"WAIT\nNOTE: Hold position."}}],"usage":{"prompt_tokens":20,"completion_tokens":8,"total_tokens":28}}`))
	})
	runner := NewRunner(testConfig(), store, client)
	if _, err := runner.Start(context.Background(), tournament); err != nil {
		t.Fatal(err)
	}
	runner.Wait()
	healthy.Store(true)
	if _, err := runner.Start(context.Background(), tournament); err != nil {
		t.Fatal(err)
	}
	runner.Wait()
	match, canonical, err := store.Match(tournament.ID)
	// Equal HP after the initial 30 rounds causes the engine's ten-round
	// overtime; two failed calls precede the 80 successful decision calls.
	const expectedRounds = 40
	const successfulCalls = expectedRounds * 2
	const attemptedCalls = successfulCalls + 2
	if err != nil || match.Status != "completed" || match.Mode != "llm" ||
		match.Outcome != "draw" || match.WinnerAddress != "" || len(match.Rounds) != expectedRounds ||
		match.LLMSummary.Attempted != attemptedCalls || match.LLMSummary.Succeeded != successfulCalls ||
		match.LLMSummary.Failed != 2 {
		t.Fatalf("retry execution history mismatch: status=%s outcome=%s rounds=%d summary=%+v err=%v",
			match.Status, match.Outcome, len(match.Rounds), match.LLMSummary, err)
	}
	for _, round := range match.Rounds {
		for n, sequence := range []int{round.RedLLMCallID, round.BlueLLMCallID} {
			if sequence <= 0 || sequence > len(match.LLMCalls) {
				t.Fatal("executed round has no call proof")
			}
			call := match.LLMCalls[sequence-1]
			if !call.Success || call.Round != round.Round || call.Sequence != sequence ||
				call.Side != []string{"red", "blue"}[n] || call.RequestID != "local-request" {
				t.Fatal("executed round refers to failed/wrong call")
			}
		}
	}
	var evidence Evidence
	if json.Unmarshal([]byte(canonical), &evidence) != nil || evidence.SchemaVersion != 2 ||
		len(evidence.LLMCalls) != attemptedCalls || crypto.Keccak256Hash([]byte(canonical)).Hex() != match.EvidenceHash {
		t.Fatal("provider attempts are not bound into v2 evidence")
	}
	if _, err := runner.Start(context.Background(), tournament); err == nil {
		t.Fatal("completed LLM evidence could be replaced")
	}
}

func TestLLMBudgetSurvivesRestartAndDoesNotSendExtraCalls(t *testing.T) {
	path := filepath.Join(t.TempDir(), "calls.db")
	store, err := NewStore(path)
	if err != nil {
		t.Fatal(err)
	}
	tournament := llmTournament(t, store)
	var requests atomic.Int64
	client := llmFixture(t, func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		w.WriteHeader(401)
	})
	cfg := testConfig()
	cfg.LLMMaxCalls = 2
	runner := NewRunner(cfg, store, client)
	if _, err := runner.Start(context.Background(), tournament); err != nil {
		t.Fatal(err)
	}
	runner.Wait()
	if requests.Load() != 2 {
		t.Fatal("one failed round should contain exactly two requests, without automatic retry")
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	store, err = NewStore(path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	runner = NewRunner(cfg, store, client)
	if _, err := runner.Start(context.Background(), tournament); err == nil || requests.Load() != 2 {
		t.Fatal("restarting reset the persisted model-call budget")
	}
	match, _, _ := store.Match(tournament.ID)
	if runner.PublicMatch(match).RetryAllowed {
		t.Fatal("exhausted call budget advertised a usable retry")
	}
	match.Status = "running"
	match.LLMCalls[0].ErrorCode = "in_progress"
	if err := store.SaveMatch(match, ""); err != nil {
		t.Fatal(err)
	}
	if err := store.RecoverRuns(); err != nil {
		t.Fatal(err)
	}
	match, _, _ = store.Match(tournament.ID)
	if match.Status != "interrupted" || match.LLMCalls[0].ErrorCode != "cancelled" ||
		len(match.LLMCalls) != 2 || match.LLMSummary.Failed != 2 {
		t.Fatal("unknown crashed request was lost or incorrectly called successful")
	}
}

func TestLegacyEvidenceCanonicalBytesRemainUnchanged(t *testing.T) {
	cfg := testConfig()
	m := &Match{TournamentID: 7, Mode: "deterministic_simulation",
		RedAgent:  AgentDefinition{1, "red", "Red", "Attack"},
		BlueAgent: AgentDefinition{1, "blue", "Blue", "Defend"}, Rounds: []Round{}}
	evidence, canonical, hash, err := BuildEvidence(cfg, m)
	if err != nil || evidence.SchemaVersion != 1 || strings.Contains(canonical, "llm_calls") {
		t.Fatal("legacy evidence acquired new hash fields")
	}
	// The original v1 field order and types are a compatibility contract.
	legacy := struct {
		SchemaVersion  int              `json:"schema_version"`
		EngineVersion  string           `json:"engine_version"`
		ChainID        uint64           `json:"chain_id"`
		ManagerAddress string           `json:"manager_address"`
		TournamentID   uint64           `json:"tournament_id"`
		Mode           string           `json:"mode"`
		RedAgent       AgentDefinition  `json:"red_agent"`
		BlueAgent      AgentDefinition  `json:"blue_agent"`
		InitialState   engine.GameState `json:"initial_state"`
		Rounds         []Round          `json:"rounds"`
		FinalState     engine.GameState `json:"final_state"`
		Outcome        string           `json:"outcome"`
		WinnerAddress  string           `json:"winner_address"`
	}{evidence.SchemaVersion, evidence.EngineVersion, evidence.ChainID, evidence.ManagerAddress,
		evidence.TournamentID, evidence.Mode, evidence.RedAgent, evidence.BlueAgent,
		evidence.InitialState, evidence.Rounds, evidence.FinalState, evidence.Outcome, evidence.WinnerAddress}
	old, err := json.Marshal(legacy)
	if err != nil || canonical != string(old) || hash != crypto.Keccak256Hash(old).Hex() {
		t.Fatal("legacy canonical/hash changed")
	}
}
