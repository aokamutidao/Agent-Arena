package tournament

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"agent-arena/backend/internal/agent"
	"agent-arena/backend/internal/engine"
	"github.com/ethereum/go-ethereum/crypto"
)

type Runner struct {
	store *Store
	cfg   Config
	llm   agent.LLMClient
	mu    sync.Mutex
	jobs  map[uint64]bool
	wg    sync.WaitGroup
}

func NewRunner(cfg Config, s *Store, llm agent.LLMClient) *Runner {
	cfg.Defaults()
	return &Runner{store: s, cfg: cfg, llm: llm, jobs: map[uint64]bool{}}
}
func (r *Runner) Wait() { r.wg.Wait() }

func (r *Runner) Start(ctx context.Context, t *Tournament) (*Match, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.jobs[t.ID] {
		m, _, err := r.store.Match(t.ID)
		return m, err
	}
	if len(r.jobs) >= 4 {
		return nil, fmt.Errorf("four matches are already running; retry later")
	}
	if t.Status != "running" || len(t.ParticipantAgents) != 2 || len(t.Participants) != 2 {
		return nil, fmt.Errorf("tournament must be running with exactly two registered participants")
	}
	profiles := make([]*AgentProfile, 2)
	for n, p := range t.ParticipantAgents {
		if address(p.Address) != address(t.Participants[n]) {
			return nil, fmt.Errorf("participant slot mismatch")
		}
		a, err := r.store.AgentByHash(p.AgentHash)
		if err != nil {
			return nil, fmt.Errorf("participant %d Agent hash is not available on this service", n+1)
		}
		canonical, err := json.Marshal(a.Definition())
		if err != nil {
			return nil, err
		}
		if a.OwnerAddress != address(p.Address) || crypto.Keccak256Hash(canonical).Hex() != p.AgentHash {
			return nil, fmt.Errorf("participant %d Agent profile ownership/hash mismatch", n+1)
		}
		profiles[n] = a
	}
	m, _, err := r.store.Match(t.ID)
	if err != nil && err != ErrNotFound {
		return nil, err
	}
	if err == nil {
		if m.Status == "completed" {
			return nil, fmt.Errorf("completed match evidence is immutable")
		}
		if m.Status != "interrupted" {
			return nil, fmt.Errorf("match already exists; inspect its status before retrying")
		}
		if m.RedAgent != profiles[0].Definition() || m.BlueAgent != profiles[1].Definition() {
			return nil, fmt.Errorf("registered Agents changed since persisted run")
		}
		if (m.Mode == "llm" && r.llm == nil) || (m.Mode == "deterministic_simulation" && !r.cfg.SimulationEnabled) {
			return nil, fmt.Errorf("original match execution mode is unavailable")
		}
		if m.Mode == "llm" && len(m.LLMCalls)+2 > r.cfg.LLMMaxCalls {
			return nil, fmt.Errorf("LLM call budget exhausted; no additional model request is allowed")
		}
	} else {
		mode := "llm"
		if r.llm == nil {
			if !r.cfg.SimulationEnabled {
				return nil, fmt.Errorf("LLM is unconfigured and simulation is disabled")
			}
			mode = "deterministic_simulation"
		}
		state := engine.NewEngine().NewGame(t.ID, profiles[0].Name, profiles[1].Name, profiles[0].Personality, profiles[1].Personality)
		state.Status = engine.StatusPlaying
		m = &Match{TournamentID: t.ID, Status: "queued", Mode: mode, RedAddress: t.Participants[0], BlueAddress: t.Participants[1], Rounds: []Round{}, State: *state, StartedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix(), RedAgent: profiles[0].Definition(), BlueAgent: profiles[1].Definition()}
	}
	m.Status = "queued"
	m.RetryAllowed = false
	m.Error = ""
	m.UpdatedAt = time.Now().Unix()
	if err := r.store.SaveMatch(m, ""); err != nil {
		return nil, err
	}
	response, _, err := r.store.Match(t.ID)
	if err != nil {
		return nil, err
	}
	r.jobs[t.ID] = true
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		defer func() { r.mu.Lock(); delete(r.jobs, t.ID); r.mu.Unlock() }()
		r.run(ctx, m)
	}()
	return response, nil
}

func (r *Runner) run(ctx context.Context, m *Match) {
	eng := engine.NewEngine()
	m.Status = "running"
	m.UpdatedAt = time.Now().Unix()
	if err := r.store.SaveMatch(m, ""); err != nil {
		return
	}
	var red, blue *agent.Agent
	if m.Mode == "llm" {
		red = agent.NewAgent(m.RedAgent.Name, m.RedAgent.Personality, r.llm)
		blue = agent.NewAgent(m.BlueAgent.Name, m.BlueAgent.Personality, r.llm)
	}
	fail := func(status, message string) {
		m.Status = status
		m.Error = message
		m.UpdatedAt = time.Now().Unix()
		_ = r.store.SaveMatch(m, "")
	}
	for m.State.Status == engine.StatusPlaying {
		if m.State.CurrentRound >= 60 {
			fail("failed", "engine exceeded the bounded round limit")
			return
		}
		if ctx.Err() != nil {
			fail("interrupted", "Server stopped; organizer can resume from the persisted round.")
			return
		}
		if r.cfg.RoundInterval > 0 {
			timer := time.NewTimer(r.cfg.RoundInterval)
			select {
			case <-ctx.Done():
				timer.Stop()
				fail("interrupted", "Server stopped; resume is available.")
				return
			case <-timer.C:
			}
		}
		var redAction, blueAction engine.Action
		var redCallID, blueCallID int
		redStyle, blueStyle := simulationStyle(m.RedAgent), simulationStyle(m.BlueAgent)
		redReason, blueReason := "Deterministic baseline: "+redStyle+"; no LLM request.", "Deterministic baseline: "+blueStyle+"; no LLM request."
		if m.Mode == "deterministic_simulation" {
			redAction = simulationAction(eng, m.State.AgentRed, m.State.AgentBlue, m.State.Obstacles, redStyle)
			blueAction = simulationAction(eng, m.State.AgentBlue, m.State.AgentRed, m.State.Obstacles, blueStyle)
		} else {
			if len(m.LLMCalls)+2 > r.cfg.LLMMaxCalls {
				fail("failed", "LLM call budget exhausted; no result was generated.")
				return
			}
			info := clientInfo(r.llm, r.cfg)
			redCallID, blueCallID = len(m.LLMCalls)+1, len(m.LLMCalls)+2
			for n, side := range []string{"red", "blue"} {
				m.LLMCalls = append(m.LLMCalls, DecisionCall{
					Sequence: redCallID + n, Round: m.State.CurrentRound + 1, Side: side,
					CallMetadata: agent.CallMetadata{
						Provider: info.Provider, Model: info.Model, ErrorCode: "in_progress",
					},
				})
			}
			m.LLMSummary = summarizeCalls(m.LLMCalls)
			// Reserve the budget before the requests: a crash cannot reset spent attempts.
			if err := r.store.SaveMatch(m, ""); err != nil {
				fail("failed", "Persisting LLM attempts failed; no requests were sent.")
				return
			}
			weights := engine.StrategyWeights{Aggressive: 34, Defensive: 33, Tricky: 33}
			a, b := agent.DecideTurnsStrictParallel(ctx, red, blue, agent.AgentDecisionInput{Self: m.State.AgentRed, Enemy: m.State.AgentBlue, CurrentRound: m.State.CurrentRound, MaxRounds: m.State.MaxRounds, Obstacles: m.State.Obstacles, Weights: weights, ActionHistory: m.State.History}, agent.AgentDecisionInput{Self: m.State.AgentBlue, Enemy: m.State.AgentRed, CurrentRound: m.State.CurrentRound, MaxRounds: m.State.MaxRounds, Obstacles: m.State.Obstacles, Weights: weights, ActionHistory: m.State.History})
			m.LLMCalls[redCallID-1].CallMetadata = a.Metadata
			m.LLMCalls[blueCallID-1].CallMetadata = b.Metadata
			m.LLMSummary = summarizeCalls(m.LLMCalls)
			if a.Err != nil || b.Err != nil {
				code := a.Metadata.ErrorCode
				if a.Err == nil {
					code = b.Metadata.ErrorCode
				}
				fail("interrupted", "LLM decision interrupted ("+code+"); no round or result was committed. Organizer may retry.")
				return
			}
			redAction, blueAction = a.Action, b.Action
			redReason, blueReason = a.Note, b.Note
		}
		if ctx.Err() != nil {
			fail("interrupted", "Server stopped before the next round was committed.")
			return
		}
		turn, err := eng.ExecuteTurn(&m.State, redAction, blueAction)
		if err != nil {
			fail("failed", "engine could not execute round")
			return
		}
		// The inherited engine checks red KO first. A simultaneous KO must remain
		// a draw in financial evidence instead of assigning a fabricated winner.
		if m.State.AgentRed.HP == 0 && m.State.AgentBlue.HP == 0 {
			m.State.Winner = engine.SideNone
			m.State.Status = engine.StatusFinished
		}
		turn.RedReasoning, turn.BlueReasoning = redReason, blueReason
		m.State.History[len(m.State.History)-1] = *turn
		m.Rounds = append(m.Rounds, Round{Round: turn.Round, RedAction: turn.RedAction, BlueAction: turn.BlueAction, RedHP: m.State.AgentRed.HP, BlueHP: m.State.AgentBlue.HP, RedPosition: m.State.AgentRed.Position, BluePosition: m.State.AgentBlue.Position, RedReasoning: redReason, BlueReasoning: blueReason, RedState: m.State.AgentRed, BlueState: m.State.AgentBlue, RedLLMCallID: redCallID, BlueLLMCallID: blueCallID})
		m.CurrentRound = m.State.CurrentRound
		m.UpdatedAt = time.Now().Unix()
		if err := r.store.SaveMatch(m, ""); err != nil {
			fail("failed", "persisting a round failed; no result was submitted")
			return
		}
	}
	m.Outcome = string(m.State.Winner)
	switch m.State.Winner {
	case engine.SideRed:
		m.WinnerAddress = m.RedAddress
	case engine.SideBlue:
		m.WinnerAddress = m.BlueAddress
	default:
		m.Outcome = "draw"
		m.WinnerAddress = ""
	}
	evidence, canonical, hash, err := BuildEvidence(r.cfg, m)
	_ = evidence
	if err != nil {
		fail("failed", "building complete evidence failed")
		return
	}
	m.Status = "completed"
	m.EvidenceHash = hash
	m.UpdatedAt = time.Now().Unix()
	if err := r.store.SaveMatch(m, canonical); err != nil {
		fail("failed", "persisting evidence failed; no result was submitted")
	}
}

func BuildEvidence(cfg Config, m *Match) (Evidence, string, string, error) {
	initial := engine.NewEngine().NewGame(m.TournamentID, m.RedAgent.Name, m.BlueAgent.Name, m.RedAgent.Personality, m.BlueAgent.Personality)
	initial.Status = engine.StatusPlaying
	e := Evidence{SchemaVersion: 1, EngineVersion: "agent-arena-grid-v1", ChainID: cfg.ChainID, ManagerAddress: address(cfg.ManagerAddress), TournamentID: m.TournamentID, Mode: m.Mode, RedAgent: m.RedAgent, BlueAgent: m.BlueAgent, InitialState: *initial, Rounds: m.Rounds, FinalState: m.State, Outcome: m.Outcome, WinnerAddress: m.WinnerAddress}
	if len(m.LLMCalls) > 0 {
		e.SchemaVersion = 2
		e.LLMCalls = m.LLMCalls
	}
	b, err := json.Marshal(e)
	if err != nil {
		return e, "", "", err
	}
	return e, string(b), crypto.Keccak256Hash(b).Hex(), nil
}

// simulationAction chooses legal moves deterministically, including a shortest
// grid path around obstacles. It is a demonstration baseline, not a fake LLM.
func simulationStyle(a AgentDefinition) string {
	p := strings.ToLower(a.Name + " " + a.Personality)
	for _, word := range []string{"defensive", "defend", "防守", "防御", "稳健"} {
		if strings.Contains(p, word) {
			return "defensive"
		}
	}
	for _, word := range []string{"aggressive", "attack", "offensive", "进攻", "狂战"} {
		if strings.Contains(p, word) {
			return "aggressive"
		}
	}
	b, _ := json.Marshal(a)
	h := crypto.Keccak256(b)
	return []string{"aggressive", "defensive", "tactical"}[int(h[0])%3]
}
func simulationAction(eng *engine.Engine, self, enemy engine.AgentState, obstacles []engine.Position, style string) engine.Action {
	dist := engine.ManhattanDistance(self.Position, enemy.Position)
	valid := func(kind engine.ActionType) bool {
		return eng.ValidateAction(engine.Action{Type: kind}, self, enemy, obstacles) == nil
	}
	if style == "defensive" && self.HP <= 70 && valid(engine.ActionHeal) {
		return engine.Action{Type: engine.ActionHeal}
	}
	if style == "aggressive" && !self.IsCharging && dist <= 2 && valid(engine.ActionCharge) {
		return engine.Action{Type: engine.ActionCharge}
	}
	if dist <= engine.SkillRange && valid(engine.ActionSkill) {
		return engine.Action{Type: engine.ActionSkill}
	}
	if dist <= engine.AtkRange && valid(engine.ActionAttack) {
		return engine.Action{Type: engine.ActionAttack}
	}
	if self.HP <= 20 && valid(engine.ActionHeal) {
		return engine.Action{Type: engine.ActionHeal}
	}
	blocked := map[engine.Position]bool{enemy.Position: true}
	for _, p := range obstacles {
		blocked[p] = true
	}
	queue := []engine.Position{self.Position}
	previous := map[engine.Position]engine.Position{}
	visited := map[engine.Position]bool{self.Position: true}
	var goal engine.Position
	found := false
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		if engine.ManhattanDistance(p, enemy.Position) <= 1 {
			goal = p
			found = true
			break
		}
		for _, d := range [][2]int{{0, -1}, {-1, 0}, {0, 1}, {1, 0}} {
			x, y := int(p.X)+d[0], int(p.Y)+d[1]
			if x < 0 || x >= 10 || y < 0 || y >= 10 {
				continue
			}
			next := engine.Position{X: uint8(x), Y: uint8(y)}
			if visited[next] || blocked[next] {
				continue
			}
			visited[next] = true
			previous[next] = p
			queue = append(queue, next)
		}
	}
	if found {
		path := []engine.Position{}
		for goal != self.Position {
			path = append([]engine.Position{goal}, path...)
			goal = previous[goal]
		}
		for n := len(path) - 1; n >= 0; n-- {
			action := engine.Action{Type: engine.ActionMove, Target: path[n]}
			if eng.ValidateAction(action, self, enemy, obstacles) == nil {
				return action
			}
		}
	}
	return engine.Action{Type: engine.ActionWait}
}
