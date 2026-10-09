package tournament

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"agent-arena/backend/internal/agent"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Server struct {
	cfg     Config
	store   *Store
	chain   ChainReader
	auth    *Auth
	indexer *Indexer
	runner  *Runner
	router  *gin.Engine
	ctx     context.Context
	cancel  context.CancelFunc
	start   sync.Once
}

func NewServer(cfg Config, store *Store, chain ChainReader, llm agent.LLMClient) (*Server, error) {
	cfg.Defaults()
	if store == nil {
		return nil, fmt.Errorf("persistent store is required")
	}
	if chain != nil && !cfg.HasDeployment() {
		return nil, fmt.Errorf("chain reader requires a complete discovered deployment config")
	}
	if len(cfg.CORSOrigins) == 0 {
		cfg.CORSOrigins = []string{"http://localhost:3000"}
	}
	for _, origin := range cfg.CORSOrigins {
		if origin == "*" {
			return nil, fmt.Errorf("CORS_ORIGINS must list explicit origins")
		}
	}
	auth, err := NewAuth(cfg, store, chain)
	if err != nil {
		return nil, err
	}
	if err := store.BindDeployment(cfg); err != nil {
		return nil, err
	}
	if err := store.RecoverRuns(); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	s := &Server{cfg: cfg, store: store, chain: chain, auth: auth, indexer: NewIndexer(cfg, store, chain), runner: NewRunner(cfg, store, llm), ctx: ctx, cancel: cancel}
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	if err := router.SetTrustedProxies(cfg.TrustedProxyCIDRs); err != nil {
		cancel()
		return nil, fmt.Errorf("TRUSTED_PROXY_CIDRS must contain valid proxy IP addresses or CIDRs")
	}
	router.Use(gin.CustomRecovery(func(c *gin.Context, _ any) { apiError(c, 500, "internal_error", "The request could not be completed.") }))
	router.Use(func(c *gin.Context) { c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 65536); c.Next() })
	router.Use(cors.New(cors.Config{AllowOrigins: cfg.CORSOrigins, AllowMethods: []string{"GET", "POST", "OPTIONS"}, AllowHeaders: []string{"Origin", "Content-Type", "Authorization"}, MaxAge: time.Hour}))
	router.Use(writeLimiter())
	s.router = router
	s.routes()
	return s, nil
}
func (s *Server) Router() *gin.Engine { return s.router }
func (s *Server) Start(ctx context.Context) {
	s.start.Do(func() {
		go s.indexer.Run(s.ctx)
		go func() {
			select {
			case <-ctx.Done():
				s.cancel()
			case <-s.ctx.Done():
			}
		}()
	})
}
func (s *Server) Close()            { s.cancel(); s.runner.Wait() }
func (s *Server) Indexer() *Indexer { return s.indexer }

func apiError(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, gin.H{"error": gin.H{"code": code, "message": message}})
}
func bindJSON(c *gin.Context, v any) bool {
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(v)
	if err == nil {
		var extra any
		if trailing := decoder.Decode(&extra); trailing != io.EOF {
			err = fmt.Errorf("request must contain one JSON object")
		}
	}
	if err != nil {
		var limit *http.MaxBytesError
		if errors.As(err, &limit) {
			apiError(c, 413, "body_too_large", "Request body exceeds 64 KiB.")
		} else {
			apiError(c, 400, "invalid_request", "Provide a valid JSON request body.")
		}
		return false
	}
	return true
}
func pagination(c *gin.Context) (int, int, bool) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit < 1 || limit > 100 {
		apiError(c, 400, "invalid_pagination", "limit must be between 1 and 100.")
		return 0, 0, false
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 || offset > 1000000 {
		apiError(c, 400, "invalid_pagination", "offset must be between 0 and 1000000.")
		return 0, 0, false
	}
	return limit, offset, true
}
func tournamentID(c *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 63)
	if err != nil || id == 0 {
		apiError(c, 400, "invalid_tournament_id", "Tournament ID must be a positive integer.")
		return 0, false
	}
	return id, true
}
func (s *Server) authenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			apiError(c, 401, "authentication_required", "Sign in with your wallet.")
			return
		}
		wallet, err := s.auth.Validate(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			apiError(c, 401, "invalid_session", "Your session is invalid or expired.")
			return
		}
		c.Set("wallet", wallet)
		c.Next()
	}
}
func wallet(c *gin.Context) string { return c.GetString("wallet") }
func (s *Server) readTournament(c *gin.Context, id uint64) (*Tournament, bool) {
	if s.chain == nil {
		apiError(c, 503, "deployment_unavailable", "Configure the course contract deployment and read-only RPC.")
		return nil, false
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	t, err := s.chain.Tournament(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			apiError(c, 404, "tournament_not_found", "Tournament does not exist.")
		} else {
			apiError(c, 503, "chain_unavailable", "The configured blockchain could not be read.")
		}
		return nil, false
	}
	// Read-through latest state is returned separately. It must not overwrite the
	// confirmed database projection maintained by the indexer.
	t.Confirmed = false
	return t, true
}

func (s *Server) routes() {
	s.router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	s.router.GET("/ready", func(c *gin.Context) {
		indexed, head, err, _, ready := s.indexer.State()
		status := 200
		if !ready {
			status = 503
		}
		c.JSON(status, gin.H{"ready": ready, "chain_configured": s.chain != nil, "indexed_block": indexed, "chain_head": head, "last_error": err})
	})
	api := s.router.Group("/api")
	api.GET("/config", func(c *gin.Context) {
		c.JSON(200, gin.H{"chain_id": s.cfg.ChainID, "deployment_configured": s.chain != nil && s.cfg.HasDeployment(), "manager_address": s.cfg.ManagerAddress, "escrow_address": s.cfg.EscrowAddress, "registry_address": s.cfg.RegistryAddress, "review_address": s.cfg.ReviewAddress, "token_address": s.cfg.TokenAddress, "token_decimals": s.cfg.TokenDecimals, "deployment_block": s.cfg.DeploymentBlock, "confirmations": s.cfg.Confirmations, "simulation_enabled": s.cfg.SimulationEnabled, "auth_domain": s.cfg.AuthDomain})
	})
	api.POST("/auth/nonce", func(c *gin.Context) {
		var req struct {
			Address string `json:"address"`
		}
		if !bindJSON(c, &req) {
			return
		}
		n, err := s.auth.IssueNonce(req.Address)
		if err != nil {
			apiError(c, 400, "invalid_wallet", "Provide a nonzero Ethereum wallet address.")
			return
		}
		c.JSON(200, n)
	})
	api.POST("/auth/login", func(c *gin.Context) {
		var req LoginRequest
		if !bindJSON(c, &req) {
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()
		result, err := s.auth.Login(ctx, req)
		if err != nil {
			apiError(c, 401, "login_failed", "Signature invalid, challenge expired/used, or chain role verification unavailable.")
			return
		}
		c.JSON(200, result)
	})
	api.GET("/auth/profile", s.authenticated(), func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()
		roles, err := s.auth.Roles(ctx, wallet(c))
		if err != nil {
			apiError(c, 503, "chain_unavailable", "Role verification is temporarily unavailable.")
			return
		}
		c.JSON(200, User{wallet(c), roles})
	})
	api.GET("/agents", func(c *gin.Context) { s.listAgents(c, "") })
	api.GET("/agents/my", s.authenticated(), func(c *gin.Context) { s.listAgents(c, wallet(c)) })
	api.POST("/agents", s.authenticated(), func(c *gin.Context) {
		var req struct {
			Name        string `json:"name"`
			Personality string `json:"personality"`
		}
		if !bindJSON(c, &req) {
			return
		}
		req.Name = strings.TrimSpace(req.Name)
		req.Personality = strings.TrimSpace(req.Personality)
		if utf8.RuneCountInString(req.Name) < 1 || utf8.RuneCountInString(req.Name) > 50 || utf8.RuneCountInString(req.Personality) < 1 || utf8.RuneCountInString(req.Personality) > 500 {
			apiError(c, 400, "invalid_agent", "Name must contain 1–50 characters and personality 1–500 characters.")
			return
		}
		a, err := s.store.CreateAgent(wallet(c), req.Name, req.Personality)
		if err != nil {
			apiError(c, 500, "storage_error", "Agent profile could not be persisted.")
			return
		}
		c.JSON(201, a)
	})
	api.GET("/tournaments", func(c *gin.Context) {
		limit, offset, ok := pagination(c)
		if !ok {
			return
		}
		items, total, err := s.store.Tournaments(limit, offset)
		if err != nil {
			apiError(c, 500, "storage_error", "Tournament projection could not be read.")
			return
		}
		c.JSON(200, gin.H{"tournaments": items, "total": total})
	})
	api.GET("/tournaments/:id", func(c *gin.Context) {
		id, ok := tournamentID(c)
		if !ok {
			return
		}
		t, ok := s.readTournament(c, id)
		if ok {
			c.JSON(200, t)
		}
	})
	api.GET("/tournaments/:id/events", func(c *gin.Context) {
		id, ok := tournamentID(c)
		if !ok {
			return
		}
		limit, offset, ok := pagination(c)
		if !ok {
			return
		}
		items, total, err := s.store.Events(id, limit, offset)
		if err != nil {
			apiError(c, 500, "storage_error", "Confirmed event log could not be read.")
			return
		}
		c.JSON(200, gin.H{"events": items, "total": total})
	})
	api.POST("/tournaments/:id/run", s.authenticated(), func(c *gin.Context) {
		id, ok := tournamentID(c)
		if !ok {
			return
		}
		t, ok := s.readTournament(c, id)
		if !ok {
			return
		}
		if address(t.Organizer) != wallet(c) {
			apiError(c, 403, "organizer_required", "Only this tournament's organizer may run its match.")
			return
		}
		m, err := s.runner.Start(s.ctx, t)
		if err != nil {
			apiError(c, 409, "run_unavailable", err.Error())
			return
		}
		c.JSON(202, m)
	})
	api.GET("/tournaments/:id/match", func(c *gin.Context) {
		id, ok := tournamentID(c)
		if !ok {
			return
		}
		m, _, err := s.store.Match(id)
		if err != nil {
			apiError(c, 404, "match_not_found", "No persisted match exists for this tournament.")
			return
		}
		c.JSON(200, m)
	})
	api.GET("/tournaments/:id/evidence", func(c *gin.Context) {
		id, ok := tournamentID(c)
		if !ok {
			return
		}
		m, canonical, err := s.store.Match(id)
		if err != nil || m.Status != "completed" || canonical == "" {
			apiError(c, 404, "evidence_unavailable", "Completed match evidence is not available.")
			return
		}
		var e Evidence
		if err := json.Unmarshal([]byte(canonical), &e); err != nil {
			apiError(c, 500, "storage_error", "Stored evidence could not be decoded.")
			return
		}
		c.JSON(200, gin.H{"evidence_hash": m.EvidenceHash, "canonical_json": canonical, "evidence": e})
	})
	api.GET("/tournaments/:id/notes", func(c *gin.Context) {
		id, ok := tournamentID(c)
		if !ok {
			return
		}
		notes, err := s.store.Notes(id)
		if err != nil {
			apiError(c, 500, "storage_error", "Review notes could not be read.")
			return
		}
		c.JSON(200, gin.H{"notes": notes, "total": len(notes)})
	})
	api.POST("/tournaments/:id/notes", s.authenticated(), func(c *gin.Context) {
		id, ok := tournamentID(c)
		if !ok {
			return
		}
		var req struct {
			Kind string `json:"kind"`
			Text string `json:"text"`
		}
		if !bindJSON(c, &req) {
			return
		}
		if req.Kind != "dispute" && req.Kind != "decision" && req.Kind != "cancellation" {
			apiError(c, 400, "invalid_note", "kind must be dispute, decision, or cancellation.")
			return
		}
		req.Text = strings.TrimSpace(req.Text)
		if utf8.RuneCountInString(req.Text) < 1 || utf8.RuneCountInString(req.Text) > 4000 {
			apiError(c, 400, "invalid_note", "Note text must contain 1–4000 characters.")
			return
		}
		t, ok := s.readTournament(c, id)
		if !ok {
			return
		}
		author := wallet(c)
		if req.Kind == "decision" && address(t.Auditor) != author {
			apiError(c, 403, "auditor_required", "Only the designated auditor may write a decision note.")
			return
		}
		if req.Kind == "cancellation" && address(t.Organizer) != author && address(t.Auditor) != author {
			apiError(c, 403, "review_authority_required", "Only the organizer or designated auditor may write this note.")
			return
		}
		n, err := s.store.SaveNote(s.cfg, id, author, req.Kind, req.Text)
		if err != nil {
			apiError(c, 500, "storage_error", "Review note could not be persisted.")
			return
		}
		c.JSON(201, gin.H{"evidence_hash": n.EvidenceHash, "canonical_json": n.CanonicalJSON, "evidence": n.NoteEvidence})
	})
	api.GET("/dashboard", s.dashboard)
	s.router.GET("/ws", s.ws)
	s.router.NoRoute(func(c *gin.Context) { apiError(c, 404, "route_not_found", "API route does not exist.") })
}
func (s *Server) listAgents(c *gin.Context, owner string) {
	items, err := s.store.Agents(owner)
	if err != nil {
		apiError(c, 500, "storage_error", "Agent profiles could not be read.")
		return
	}
	c.JSON(200, gin.H{"agents": items, "total": len(items)})
}
func (s *Server) dashboard(c *gin.Context) {
	totals, ledger, eventCount, err := s.store.LedgerTotals()
	if err != nil {
		apiError(c, 500, "storage_error", "Confirmed financial ledger could not be read.")
		return
	}
	rows, err := s.store.db.Query("SELECT snapshot FROM tournaments")
	if err != nil {
		apiError(c, 500, "storage_error", "Tournament projection could not be read.")
		return
	}
	defer rows.Close()
	count, participants, active, completed := 0, 0, 0, 0
	for rows.Next() {
		var encoded string
		if err := rows.Scan(&encoded); err != nil {
			apiError(c, 500, "storage_error", "Tournament projection is unavailable.")
			return
		}
		var t Tournament
		if err := json.Unmarshal([]byte(encoded), &t); err != nil {
			apiError(c, 500, "storage_error", "Tournament projection is unavailable.")
			return
		}
		count++
		participants += len(t.Participants)
		if t.Status == "running" || t.Status == "result_submitted" || t.Status == "disputed" {
			active++
		}
		if t.Status == "finalized" {
			completed++
		}
	}
	if rows.Err() != nil {
		apiError(c, 500, "storage_error", "Tournament projection is unavailable.")
		return
	}
	indexed, head, _, observed, _ := s.indexer.State()
	c.JSON(200, gin.H{"tournaments_total": count, "registered_participants": participants, "active_tournaments": active, "completed_tournaments": completed, "escrow_in_raw": totals["deposit"], "prizes_paid_raw": totals["prize"], "refunds_paid_raw": totals["refund"], "ledger_entries": ledger, "event_count": eventCount, "indexed_block": indexed, "chain_head": head, "observed_at": observed})
}
func (s *Server) ws(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("tournament_id"), 10, 63)
	if err != nil || id == 0 {
		apiError(c, 400, "invalid_tournament_id", "Provide a positive tournament_id.")
		return
	}
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true
		}
		for _, allowed := range s.cfg.CORSOrigins {
			if origin == allowed {
				return true
			}
		}
		return false
	}}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	done := make(chan struct{})
	go func() {
		defer close(done)
		conn.SetReadLimit(1024)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		m, _, err := s.store.Match(id)
		if err == nil {
			_ = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err := conn.WriteJSON(gin.H{"type": "match_update", "data": m}); err != nil {
				return
			}
		}
		select {
		case <-s.ctx.Done():
			return
		case <-done:
			return
		case <-ticker.C:
		}
	}
}

// Bound anonymous authentication and note/profile writes. Forwarded addresses
// are accepted only from explicitly configured proxies. Expired buckets are
// periodically removed.
func writeLimiter() gin.HandlerFunc {
	type bucket struct {
		since time.Time
		count int
	}
	buckets := map[string]bucket{}
	var mu sync.Mutex
	return func(c *gin.Context) {
		if c.Request.Method != "POST" {
			c.Next()
			return
		}
		now := time.Now()
		mu.Lock()
		key := c.ClientIP()
		b := buckets[key]
		if now.Sub(b.since) >= time.Minute {
			b = bucket{since: now}
		}
		b.count++
		buckets[key] = b
		allowed := b.count <= 30
		if len(buckets) > 10000 {
			for k, v := range buckets {
				if now.Sub(v.since) > time.Minute {
					delete(buckets, k)
				}
			}
		}
		mu.Unlock()
		if !allowed {
			apiError(c, 429, "rate_limited", "Too many write requests; retry in one minute.")
			return
		}
		c.Next()
	}
}
