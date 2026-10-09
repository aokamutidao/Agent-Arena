package tournament

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

func testConfig() Config {
	return Config{JWTSecret: strings.Repeat("s", 40), AuthDomain: "localhost:3000", ChainID: 31337, ManagerAddress: "0x1111111111111111111111111111111111111111", EscrowAddress: "0x2222222222222222222222222222222222222222", RegistryAddress: "0x3333333333333333333333333333333333333333", ReviewAddress: "0x4444444444444444444444444444444444444444", TokenAddress: "0x5555555555555555555555555555555555555555", TokenDecimals: 6, DeploymentBlock: 10, Confirmations: 3, SimulationEnabled: true, PollInterval: time.Second}
}
func testStore(t *testing.T) *Store {
	t.Helper()
	s, err := NewStore(filepath.Join(t.TempDir(), "course.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}
func testKey(t *testing.T) (*ecdsa.PrivateKey, string) {
	t.Helper()
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	return key, address(crypto.PubkeyToAddress(key.PublicKey).Hex())
}
func signature(t *testing.T, key *ecdsa.PrivateKey, message string) string {
	t.Helper()
	hash := crypto.Keccak256([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)))
	sig, err := crypto.Sign(hash, key)
	if err != nil {
		t.Fatal(err)
	}
	sig[64] += 27
	return "0x" + fmt.Sprintf("%x", sig)
}

func TestWalletNonceBindingReplayExpiryAndWrongWallet(t *testing.T) {
	s := testStore(t)
	cfg := testConfig()
	auth, err := NewAuth(cfg, s, nil)
	if err != nil {
		t.Fatal(err)
	}
	clock := time.Date(2026, 10, 9, 0, 0, 0, 0, time.UTC)
	auth.now = func() time.Time { return clock }
	key, wallet := testKey(t)
	n, err := auth.IssueNonce(wallet)
	if err != nil {
		t.Fatal(err)
	}
	for _, part := range []string{"localhost:3000", "URI: http://localhost:3000", "Chain ID: 31337", n.Nonce, "Expiration Time:"} {
		if !strings.Contains(n.Message, part) {
			t.Fatalf("challenge missing %s", part)
		}
	}
	wrong, _ := testKey(t)
	if _, err := auth.Login(context.Background(), LoginRequest{wallet, n.Nonce, signature(t, wrong, n.Message)}); err == nil {
		t.Fatal("wrong wallet accepted")
	}
	for _, altered := range []Config{func() Config { c := cfg; c.AuthDomain = "other.example"; return c }(), func() Config { c := cfg; c.ChainID = 11155111; return c }(), func() Config { c := cfg; c.ManagerAddress = "0x9999999999999999999999999999999999999999"; return c }()} {
		other, err := NewAuth(altered, s, nil)
		if err != nil {
			t.Fatal(err)
		}
		other.now = auth.now
		if _, err := other.Login(context.Background(), LoginRequest{wallet, n.Nonce, signature(t, key, n.Message)}); err == nil {
			t.Fatal("persisted nonce accepted by a different domain/chain/deployment")
		}
	}
	login, err := auth.Login(context.Background(), LoginRequest{wallet, n.Nonce, signature(t, key, n.Message)})
	if err != nil {
		t.Fatal(err)
	}
	if len(login.User.Roles) != 1 || login.User.Roles[0] != "participant" {
		t.Fatal("unconfigured RPC must not invent privileged roles")
	}
	if got, err := auth.Validate(login.Token); err != nil || got != wallet {
		t.Fatal("valid session rejected")
	}
	otherConfig := cfg
	otherConfig.ChainID = 11155111
	otherAuth, err := NewAuth(otherConfig, s, nil)
	if err != nil {
		t.Fatal(err)
	}
	otherAuth.now = auth.now
	if _, err := otherAuth.Validate(login.Token); err == nil {
		t.Fatal("session accepted by a different chain")
	}
	if _, err := auth.Login(context.Background(), LoginRequest{wallet, n.Nonce, signature(t, key, n.Message)}); err == nil {
		t.Fatal("replayed challenge accepted")
	}
	expired, err := auth.IssueNonce(wallet)
	if err != nil {
		t.Fatal(err)
	}
	clock = clock.Add(5 * time.Minute)
	if _, err := auth.Login(context.Background(), LoginRequest{wallet, expired.Nonce, signature(t, key, expired.Message)}); err == nil {
		t.Fatal("expired challenge accepted")
	}
	clock = clock.Add(time.Hour)
	if _, err := auth.Validate(login.Token); err == nil {
		t.Fatal("expired JWT accepted")
	}
	if _, err := NewAuth(Config{JWTSecret: "known-default", AuthDomain: "localhost"}, s, nil); err == nil {
		t.Fatal("short secret accepted")
	}
}

func TestNonceConsumedExactlyOnceUnderConcurrency(t *testing.T) {
	s := testStore(t)
	auth, err := NewAuth(testConfig(), s, nil)
	if err != nil {
		t.Fatal(err)
	}
	key, wallet := testKey(t)
	n, err := auth.IssueNonce(wallet)
	if err != nil {
		t.Fatal(err)
	}
	req := LoginRequest{wallet, n.Nonce, signature(t, key, n.Message)}
	var wg sync.WaitGroup
	var mu sync.Mutex
	success := 0
	for k := 0; k < 8; k++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := auth.Login(context.Background(), req); err == nil {
				mu.Lock()
				success++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if success != 1 {
		t.Fatalf("nonce accepted %d times", success)
	}
}

type fakeChain struct {
	tournament   Tournament
	roleWallet   string
	head         uint64
	events       []ChainEvent
	checkpoint   string
	pinnedBlocks []uint64
}

func (f *fakeChain) Tournament(context.Context, uint64) (*Tournament, error) {
	t := f.tournament
	t.Participants = append([]string{}, t.Participants...)
	t.ParticipantAgents = append([]ParticipantAgent{}, t.ParticipantAgents...)
	return &t, nil
}
func (f *fakeChain) Roles(_ context.Context, wallet string) ([]string, error) {
	if wallet == f.roleWallet {
		return []string{"organizer"}, nil
	}
	return []string{}, nil
}
func (f *fakeChain) Head(context.Context) (uint64, error) { return f.head, nil }
func (f *fakeChain) Events(_ context.Context, from, to uint64) ([]ChainEvent, error) {
	out := []ChainEvent{}
	for _, e := range f.events {
		if e.BlockNumber >= from && e.BlockNumber <= to {
			out = append(out, e)
		}
	}
	return out, nil
}
func (f *fakeChain) TournamentAt(_ context.Context, _ uint64, block uint64) (*Tournament, error) {
	f.pinnedBlocks = append(f.pinnedBlocks, block)
	t := f.tournament
	t.Status = "open"
	return &t, nil
}
func (f *fakeChain) BlockHash(context.Context, uint64) (string, error) { return f.checkpoint, nil }

func TestConfirmedEventsAreIdempotentAcrossRestartAndPinned(t *testing.T) {
	path := filepath.Join(t.TempDir(), "course.db")
	s, err := NewStore(path)
	if err != nil {
		t.Fatal(err)
	}
	cfg := testConfig()
	if err := s.BindDeployment(cfg); err != nil {
		t.Fatal(err)
	}
	e := ChainEvent{ID: "31337:tx:0", ChainID: 31337, TournamentID: 1, ContractAddress: cfg.ManagerAddress, EventName: "TournamentFunded", TxHash: "tx", LogIndex: 0, BlockNumber: 10, Actor: "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", AmountRaw: "1000000000000000000000000000000", Data: json.RawMessage(`{}`)}
	chain := &fakeChain{tournament: Tournament{ID: 1, Status: "finalized", EntryFee: "0", PrizePool: e.AmountRaw, PrizePaid: "0", Participants: []string{}, ParticipantAgents: []ParticipantAgent{}}, head: 15, events: []ChainEvent{e}, checkpoint: "block-A"}
	i := NewIndexer(cfg, s, chain)
	if err := i.SyncOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(chain.pinnedBlocks) != 1 || chain.pinnedBlocks[0] != 12 {
		t.Fatal("snapshot was not pinned to confirmed batch height")
	}
	items, _, err := s.Tournaments(20, 0)
	if err != nil || len(items) != 1 || items[0].Status != "open" || !items[0].Confirmed {
		t.Fatal("latest/unconfirmed state contaminated projection")
	}
	if err := s.ApplyBatch(i.source, 13, []ChainEvent{e}, nil); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	s, err = NewStore(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	i = NewIndexer(cfg, s, chain)
	if err := i.SyncOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	totals, n, count, err := s.LedgerTotals()
	if err != nil || n != 1 || count != 1 || totals["deposit"] != e.AmountRaw {
		t.Fatalf("duplicate or imprecise financial ledger: %v", err)
	}
	chain.checkpoint = "block-B"
	if err := i.SyncOnce(context.Background()); err == nil {
		t.Fatal("changed confirmed checkpoint accepted")
	}
	_, _, _, _, ready := i.State()
	if ready {
		t.Fatal("reorg must fail readiness closed")
	}
	other := cfg
	other.ManagerAddress = "0x9999999999999999999999999999999999999999"
	if err := s.BindDeployment(other); err == nil {
		t.Fatal("different deployment reused financial DB")
	}
}
func requestJSON(t *testing.T, s *Server, method, path string, body any, token string) *httptest.ResponseRecorder {
	t.Helper()
	var b bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&b).Encode(body); err != nil {
			t.Fatal(err)
		}
	}
	req := httptest.NewRequest(method, path, &b)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	out := httptest.NewRecorder()
	s.Router().ServeHTTP(out, req)
	return out
}
func loginToken(t *testing.T, s *Server, key *ecdsa.PrivateKey, wallet string) string {
	t.Helper()
	n, err := s.auth.IssueNonce(wallet)
	if err != nil {
		t.Fatal(err)
	}
	result, err := s.auth.Login(context.Background(), LoginRequest{wallet, n.Nonce, signature(t, key, n.Message)})
	if err != nil {
		t.Fatal(err)
	}
	return result.Token
}

func TestRunAuthorizationAndImmutableAgentOwnership(t *testing.T) {
	store := testStore(t)
	cfg := testConfig()
	organizerKey, organizer := testKey(t)
	participantKey, participant := testKey(t)
	_, blue := testKey(t)
	redAgent, err := store.CreateAgent(participant, "Attack Agent", "aggressive")
	if err != nil {
		t.Fatal(err)
	}
	blueAgent, err := store.CreateAgent(blue, "Guard Agent", "defensive")
	if err != nil {
		t.Fatal(err)
	}
	chain := &fakeChain{tournament: Tournament{ID: 1, Organizer: organizer, Status: "open", Participants: []string{participant, blue}, ParticipantAgents: []ParticipantAgent{{participant, redAgent.AgentHash}, {blue, blueAgent.AgentHash}}}, roleWallet: organizer}
	s, err := NewServer(cfg, store, chain, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if got := requestJSON(t, s, "POST", "/api/tournaments/1/run", nil, ""); got.Code != 401 {
		t.Fatal("anonymous organizer operation accepted")
	}
	participantToken := loginToken(t, s, participantKey, participant)
	if got := requestJSON(t, s, "POST", "/api/tournaments/1/run", nil, participantToken); got.Code != 403 {
		t.Fatalf("other wallet accepted: %d", got.Code)
	}
	token := loginToken(t, s, organizerKey, organizer)
	if got := requestJSON(t, s, "POST", "/api/tournaments/1/run", nil, token); got.Code != 409 {
		t.Fatal("not-running tournament accepted")
	}
	chain.tournament.Status = "running"
	chain.tournament.ParticipantAgents[1].AgentHash = redAgent.AgentHash
	if got := requestJSON(t, s, "POST", "/api/tournaments/1/run", nil, token); got.Code != 409 {
		t.Fatal("another participant's profile accepted")
	}
	chain.tournament.ParticipantAgents[1].AgentHash = blueAgent.AgentHash
	if got := requestJSON(t, s, "POST", "/api/tournaments/1/run", nil, token); got.Code != 202 {
		t.Fatalf("valid run rejected: %d %s", got.Code, got.Body.String())
	}
	s.runner.Wait()
	match, canonical, err := store.Match(1)
	if err != nil || match.Status != "completed" {
		t.Fatalf("simulation did not finish: %v %+v", err, match)
	}
	if match.Mode != "deterministic_simulation" || len(match.Rounds) == 0 || len(match.State.History) != len(match.Rounds) {
		t.Fatal("simulation evidence is incomplete")
	}
	if crypto.Keccak256Hash([]byte(canonical)).Hex() != match.EvidenceHash {
		t.Fatal("evidence hash does not cover exact canonical document")
	}
	var evidence Evidence
	if err := json.Unmarshal([]byte(canonical), &evidence); err != nil {
		t.Fatal(err)
	}
	if evidence.RedAgent.OwnerAddress != participant || evidence.BlueAgent.OwnerAddress != blue {
		t.Fatal("Agent ownership missing from evidence")
	}
	original := match.EvidenceHash
	match.Rounds[0].RedAction.Type = "WAIT"
	_, _, mutated, err := BuildEvidence(cfg, match)
	if err != nil || mutated == original {
		t.Fatal("action mutation did not change full evidence hash")
	}
	if got := requestJSON(t, s, "POST", "/api/tournaments/1/run", nil, token); got.Code != 409 {
		t.Fatal("immutable completed evidence rerun accepted")
	}
	t.Logf("deterministic aggressive/defensive outcome=%s winner=%s rounds=%d", match.Outcome, match.WinnerAddress, match.CurrentRound)
}

func TestInterruptedJobResumesPersistedState(t *testing.T) {
	s := testStore(t)
	cfg := testConfig()
	_, red := testKey(t)
	_, blue := testKey(t)
	a, err := s.CreateAgent(red, "Attack", "aggressive")
	if err != nil {
		t.Fatal(err)
	}
	b, err := s.CreateAgent(blue, "Guard", "defensive")
	if err != nil {
		t.Fatal(err)
	}
	tournament := &Tournament{ID: 9, Status: "running", Participants: []string{red, blue}, ParticipantAgents: []ParticipantAgent{{red, a.AgentHash}, {blue, b.AgentHash}}}
	cfg.RoundInterval = 10 * time.Millisecond
	r := NewRunner(cfg, s, nil)
	ctx, cancel := context.WithCancel(context.Background())
	if _, err := r.Start(ctx, tournament); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		m, _, err := s.Match(9)
		if err == nil && m.CurrentRound >= 1 {
			break
		}
		time.Sleep(time.Millisecond)
	}
	cancel()
	r.Wait()
	before, _, err := s.Match(9)
	if err != nil || before.Status != "interrupted" || before.CurrentRound < 1 {
		t.Fatal("persisted resumable round missing")
	}
	if err := s.RecoverRuns(); err != nil {
		t.Fatal(err)
	}
	cfg.RoundInterval = 0
	r = NewRunner(cfg, s, nil)
	if _, err := r.Start(context.Background(), tournament); err != nil {
		t.Fatal(err)
	}
	r.Wait()
	after, _, err := s.Match(9)
	if err != nil || after.Status != "completed" || after.CurrentRound <= before.CurrentRound {
		t.Fatal("interrupted round did not resume")
	}
	firstBefore, _ := json.Marshal(before.Rounds[0])
	firstAfter, _ := json.Marshal(after.Rounds[0])
	if !bytes.Equal(firstBefore, firstAfter) {
		t.Fatal("resume changed already committed round")
	}
}

func TestNoDeploymentCannotExposeFinancialWritesAndErrorsAreStructured(t *testing.T) {
	s := testStore(t)
	cfg := testConfig()
	cfg.ManagerAddress = ""
	cfg.EscrowAddress = ""
	cfg.RegistryAddress = ""
	cfg.ReviewAddress = ""
	cfg.TokenAddress = ""
	server, err := NewServer(cfg, s, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	for _, path := range []string{"/api/tournaments/1", "/ready"} {
		got := requestJSON(t, server, http.MethodGet, path, nil, "")
		if got.Code != 503 {
			t.Fatalf("unconfigured financial deployment appeared ready: %s %d", path, got.Code)
		}
	}
	if got := requestJSON(t, server, "GET", "/api/tournaments?offset=-1", nil, ""); got.Code != 400 || !strings.Contains(got.Body.String(), "invalid_pagination") {
		t.Fatal("negative offset is not handled safely")
	}
	if got := requestJSON(t, server, "POST", "/api/bets/place", map[string]any{"amount": "999999"}, ""); got.Code != 404 {
		t.Fatal("old financial mutation route mounted")
	}
	key, wallet := testKey(t)
	token := loginToken(t, server, key, wallet)
	got := requestJSON(t, server, "POST", "/api/agents", map[string]any{"name": "A", "personality": "aggressive", "api_key": "must-never-be-saved", "api_endpoint": "http://127.0.0.1/private"}, token)
	if got.Code != 400 || strings.Contains(got.Body.String(), "must-never-be-saved") {
		t.Fatal("Agent endpoint accepted secret/endpoint fields")
	}
	large := requestJSON(t, server, "POST", "/api/agents", map[string]any{"name": strings.Repeat("A", 70000), "personality": "aggressive"}, token)
	if large.Code != 413 {
		t.Fatalf("oversized request not bounded: %d", large.Code)
	}
}

func TestOfflineProfilesBindFirstDeploymentAndRollbackInvalidFinancialBatch(t *testing.T) {
	s := testStore(t)
	offline := testConfig()
	offline.ManagerAddress = ""
	offline.DeploymentBlock = 0
	if err := s.BindDeployment(offline); err != nil {
		t.Fatal(err)
	}
	_, owner := testKey(t)
	p, err := s.CreateAgent(owner, "A", "aggressive")
	if err != nil {
		t.Fatal(err)
	}
	configured := testConfig()
	configured.DeploymentBlock = 42
	if err := s.BindDeployment(configured); err != nil {
		t.Fatal("profile-only DB could not bind its first deployment", err)
	}
	if got, err := s.AgentByHash(p.AgentHash); err != nil || got.OwnerAddress != owner {
		t.Fatal("first deployment lost offline Agent profiles")
	}
	bad := ChainEvent{ID: "bad", ChainID: 31337, TournamentID: 1, TxHash: "bad", EventName: "TournamentFunded", BlockNumber: 42, AmountRaw: "untrusted-invalid", Data: json.RawMessage(`{}`)}
	if err := s.ApplyBatch("test", 43, []ChainEvent{bad}, nil); err == nil {
		t.Fatal("invalid amount reached financial ledger")
	}
	_, ledger, events, err := s.LedgerTotals()
	if err != nil || ledger != 0 || events != 0 {
		t.Fatal("failed batch was only partially rolled back")
	}
	next, err := s.Cursor("test", 42)
	if err != nil || next != 42 {
		t.Fatal("failed batch advanced its cursor")
	}
}

func TestReviewNoteAuthorizationAndCanonicalHash(t *testing.T) {
	s := testStore(t)
	cfg := testConfig()
	auditorKey, auditor := testKey(t)
	userKey, user := testKey(t)
	chain := &fakeChain{tournament: Tournament{ID: 1, Organizer: cfg.ManagerAddress, Auditor: auditor, Status: "result_submitted"}}
	server, err := NewServer(cfg, s, chain, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()
	token := loginToken(t, server, userKey, user)
	if got := requestJSON(t, server, "POST", "/api/tournaments/1/notes", map[string]any{"kind": "decision", "text": "approve"}, token); got.Code != 403 {
		t.Fatal("ordinary wallet created an auditor decision")
	}
	auditorToken := loginToken(t, server, auditorKey, auditor)
	got := requestJSON(t, server, "POST", "/api/tournaments/1/notes", map[string]any{"kind": "decision", "text": "Round logs reproduce the result."}, auditorToken)
	if got.Code != 201 {
		t.Fatalf("auditor note rejected: %d", got.Code)
	}
	var result struct {
		EvidenceHash  string `json:"evidence_hash"`
		CanonicalJSON string `json:"canonical_json"`
	}
	if err := json.Unmarshal(got.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if crypto.Keccak256Hash([]byte(result.CanonicalJSON)).Hex() != result.EvidenceHash {
		t.Fatal("review note hash does not cover the stored note")
	}
	notes, err := s.Notes(1)
	if err != nil || len(notes) != 1 || notes[0].AuthorAddress != auditor {
		t.Fatal("review note not persisted with its author")
	}
}

func TestWriteLimiterTrustsOnlyExplicitProxyPeers(t *testing.T) {
	for _, tc := range []struct {
		name       string
		proxies    []string
		remoteAddr string
		lastStatus int
	}{
		{"default rejects spoofed forwarding", nil, "10.20.0.2:12345", 429},
		{"explicit proxy separates client IPs", []string{"10.20.0.0/24"}, "10.20.0.2:12345", 200},
		{"untrusted peer still rejects spoofing", []string{"10.20.0.0/24"}, "10.21.0.2:12345", 429},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg := testConfig()
			cfg.TrustedProxyCIDRs = tc.proxies
			server, err := NewServer(cfg, testStore(t), nil, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer server.Close()
			for i := 1; i <= 31; i++ {
				req := httptest.NewRequest(http.MethodPost, "/api/auth/nonce", strings.NewReader(`{"address":"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`))
				req.RemoteAddr = tc.remoteAddr
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-Forwarded-For", fmt.Sprintf("192.0.2.%d", i))
				req.Header.Set("X-Real-IP", fmt.Sprintf("192.0.2.%d", i))
				out := httptest.NewRecorder()
				server.Router().ServeHTTP(out, req)
				expected := 200
				if i == 31 {
					expected = tc.lastStatus
				}
				if out.Code != expected {
					t.Fatalf("request %d returned %d, expected %d: %s", i, out.Code, expected, out.Body.String())
				}
			}
		})
	}
	cfg := testConfig()
	cfg.TrustedProxyCIDRs = []string{"invalid-proxy/999"}
	if server, err := NewServer(cfg, testStore(t), nil, nil); err == nil {
		server.Close()
		t.Fatal("invalid trusted proxy configuration accepted")
	}
}
