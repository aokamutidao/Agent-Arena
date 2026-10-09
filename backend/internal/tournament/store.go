package tournament

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	_ "github.com/mattn/go-sqlite3"
)

type Store struct{ db *sql.DB }

func NewStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	s := &Store{db: db}
	statements := []string{
		"PRAGMA busy_timeout=5000", "PRAGMA foreign_keys=ON", "PRAGMA journal_mode=WAL",
		"CREATE TABLE IF NOT EXISTS login_nonces (nonce TEXT PRIMARY KEY,address TEXT NOT NULL,message TEXT NOT NULL,expires_at INTEGER NOT NULL,consumed INTEGER NOT NULL DEFAULT 0)",
		"CREATE TABLE IF NOT EXISTS agent_profiles (agent_hash TEXT PRIMARY KEY,id TEXT UNIQUE NOT NULL,owner_address TEXT NOT NULL,profile TEXT NOT NULL,created_at INTEGER NOT NULL)",
		"CREATE INDEX IF NOT EXISTS agent_profiles_owner ON agent_profiles(owner_address)",
		"CREATE TABLE IF NOT EXISTS tournaments (id INTEGER PRIMARY KEY,snapshot TEXT NOT NULL)",
		"CREATE TABLE IF NOT EXISTS chain_events (id TEXT PRIMARY KEY,chain_id INTEGER NOT NULL,tournament_id INTEGER NOT NULL,tx_hash TEXT NOT NULL,log_index INTEGER NOT NULL,block_number INTEGER NOT NULL,event_name TEXT NOT NULL,event_json TEXT NOT NULL,UNIQUE(chain_id,tx_hash,log_index))",
		"CREATE INDEX IF NOT EXISTS events_tournament ON chain_events(tournament_id,block_number,log_index)",
		"CREATE TABLE IF NOT EXISTS financial_ledger (event_id TEXT PRIMARY KEY REFERENCES chain_events(id),tournament_id INTEGER NOT NULL,kind TEXT NOT NULL,actor TEXT NOT NULL,amount_raw TEXT NOT NULL)",
		"CREATE TABLE IF NOT EXISTS index_cursors (source TEXT PRIMARY KEY,next_block INTEGER NOT NULL)",
		"CREATE TABLE IF NOT EXISTS index_checkpoints (source TEXT PRIMARY KEY,height INTEGER NOT NULL,block_hash TEXT NOT NULL)",
		"CREATE TABLE IF NOT EXISTS matches (tournament_id INTEGER PRIMARY KEY,status TEXT NOT NULL,match_json TEXT NOT NULL,canonical_evidence TEXT NOT NULL DEFAULT '')",
		"CREATE TABLE IF NOT EXISTS database_identity (id INTEGER PRIMARY KEY CHECK(id=1),identity TEXT NOT NULL)",
		"CREATE TABLE IF NOT EXISTS review_notes (evidence_hash TEXT PRIMARY KEY,tournament_id INTEGER NOT NULL,note_json TEXT NOT NULL,created_at INTEGER NOT NULL)",
	}
	for _, q := range statements {
		if _, err := db.Exec(q); err != nil {
			db.Close()
			return nil, fmt.Errorf("initialize database: %w", err)
		}
	}
	return s, nil
}
func (s *Store) Close() error { return s.db.Close() }

// One database belongs to one deployment. Offline Agent profiles may be retained
// when the first real deployment is configured, before any tournament is indexed.
func (s *Store) BindDeployment(cfg Config) error {
	manager := "unconfigured"
	if cfg.ManagerAddress != "" {
		manager = address(cfg.ManagerAddress)
	}
	identity := fmt.Sprintf("%d:%s:%d", cfg.ChainID, manager, cfg.DeploymentBlock)
	var old string
	err := s.db.QueryRow("SELECT identity FROM database_identity WHERE id=1").Scan(&old)
	if err == sql.ErrNoRows {
		_, err = s.db.Exec("INSERT INTO database_identity(id,identity) VALUES(1,?)", identity)
		return err
	}
	if err != nil {
		return err
	}
	if old == identity {
		return nil
	}
	var count int
	if err := s.db.QueryRow("SELECT (SELECT COUNT(*) FROM chain_events)+(SELECT COUNT(*) FROM tournaments)+(SELECT COUNT(*) FROM matches)").Scan(&count); err != nil {
		return err
	}
	if strings.HasPrefix(old, fmt.Sprintf("%d:unconfigured:", cfg.ChainID)) && count == 0 {
		_, err = s.db.Exec("UPDATE database_identity SET identity=? WHERE id=1", identity)
		return err
	}
	return fmt.Errorf("database belongs to another blockchain deployment; use a new DATABASE_PATH")
}

func (s *Store) SaveNote(cfg Config, id uint64, author, kind, text string) (*Note, error) {
	e := NoteEvidence{1, cfg.ChainID, address(cfg.ManagerAddress), id, author, kind, text}
	canonical, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	n := &Note{NoteEvidence: e, EvidenceHash: crypto.Keccak256Hash(canonical).Hex(), CreatedAt: time.Now().Unix(), CanonicalJSON: string(canonical)}
	encoded, err := json.Marshal(n)
	if err != nil {
		return nil, err
	}
	_, err = s.db.Exec("INSERT OR IGNORE INTO review_notes(evidence_hash,tournament_id,note_json,created_at) VALUES(?,?,?,?)", n.EvidenceHash, id, string(encoded), n.CreatedAt)
	return n, err
}
func (s *Store) Notes(id uint64) ([]Note, error) {
	rows, err := s.db.Query("SELECT note_json FROM review_notes WHERE tournament_id=? ORDER BY created_at DESC LIMIT 100", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []Note{}
	for rows.Next() {
		var encoded string
		if err := rows.Scan(&encoded); err != nil {
			return nil, err
		}
		var n Note
		if err := json.Unmarshal([]byte(encoded), &n); err != nil {
			return nil, err
		}
		result = append(result, n)
	}
	return result, rows.Err()
}

func randomID() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *Store) CreateNonce(a, message, nonce string, expires int64) error {
	_, err := s.db.Exec("INSERT INTO login_nonces(nonce,address,message,expires_at) VALUES(?,?,?,?)", nonce, a, message, expires)
	return err
}
func (s *Store) Nonce(a, nonce string, now int64) (string, error) {
	var msg string
	err := s.db.QueryRow("SELECT message FROM login_nonces WHERE nonce=? AND address=? AND consumed=0 AND expires_at>?", nonce, a, now).Scan(&msg)
	if err == sql.ErrNoRows {
		return "", ErrNotFound
	}
	return msg, err
}
func (s *Store) ConsumeNonce(a, nonce string, now int64) error {
	r, err := s.db.Exec("UPDATE login_nonces SET consumed=1 WHERE nonce=? AND address=? AND consumed=0 AND expires_at>?", nonce, a, now)
	if err != nil {
		return err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return ErrConflict
	}
	return nil
}

func (s *Store) CreateAgent(a, name, personality string) (*AgentProfile, error) {
	id, err := randomID()
	if err != nil {
		return nil, err
	}
	p := AgentProfile{ID: id, OwnerAddress: a, Name: name, Personality: personality, CreatedAt: time.Now().Unix()}
	canonical, err := json.Marshal(p.Definition())
	if err != nil {
		return nil, err
	}
	p.AgentHash = crypto.Keccak256Hash(canonical).Hex()
	encoded, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	_, err = s.db.Exec("INSERT OR IGNORE INTO agent_profiles(agent_hash,id,owner_address,profile,created_at) VALUES(?,?,?,?,?)", p.AgentHash, p.ID, a, string(encoded), p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return s.AgentByHash(p.AgentHash)
}
func (s *Store) AgentByHash(h string) (*AgentProfile, error) {
	var encoded string
	err := s.db.QueryRow("SELECT profile FROM agent_profiles WHERE agent_hash=?", h).Scan(&encoded)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	var p AgentProfile
	err = json.Unmarshal([]byte(encoded), &p)
	return &p, err
}
func (s *Store) Agents(owner string) ([]AgentProfile, error) {
	q := "SELECT profile FROM agent_profiles"
	args := []any{}
	if owner != "" {
		q += " WHERE owner_address=?"
		args = append(args, owner)
	}
	q += " ORDER BY created_at DESC LIMIT 100"
	rows, err := s.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	p := []AgentProfile{}
	for rows.Next() {
		var encoded string
		if err := rows.Scan(&encoded); err != nil {
			return nil, err
		}
		var a AgentProfile
		if err := json.Unmarshal([]byte(encoded), &a); err != nil {
			return nil, err
		}
		p = append(p, a)
	}
	return p, rows.Err()
}

func (s *Store) SaveTournament(t *Tournament) error {
	b, err := json.Marshal(t)
	if err != nil {
		return err
	}
	_, err = s.db.Exec("INSERT INTO tournaments(id,snapshot) VALUES(?,?) ON CONFLICT(id) DO UPDATE SET snapshot=excluded.snapshot", t.ID, string(b))
	return err
}
func (s *Store) Tournaments(limit, offset int) ([]Tournament, int, error) {
	var total int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM tournaments").Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := s.db.Query("SELECT snapshot FROM tournaments ORDER BY id DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	tournaments := []Tournament{}
	for rows.Next() {
		var encoded string
		if err := rows.Scan(&encoded); err != nil {
			return nil, 0, err
		}
		var t Tournament
		if err := json.Unmarshal([]byte(encoded), &t); err != nil {
			return nil, 0, err
		}
		tournaments = append(tournaments, t)
	}
	return tournaments, total, rows.Err()
}
func (s *Store) Events(id uint64, limit, offset int) ([]ChainEvent, int, error) {
	var total int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM chain_events WHERE tournament_id=?", id).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := s.db.Query("SELECT event_json FROM chain_events WHERE tournament_id=? ORDER BY block_number DESC,log_index DESC LIMIT ? OFFSET ?", id, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	events := []ChainEvent{}
	for rows.Next() {
		var encoded string
		if err := rows.Scan(&encoded); err != nil {
			return nil, 0, err
		}
		var e ChainEvent
		if err := json.Unmarshal([]byte(encoded), &e); err != nil {
			return nil, 0, err
		}
		events = append(events, e)
	}
	return events, total, rows.Err()
}
func (s *Store) Cursor(source string, initial uint64) (uint64, error) {
	var n uint64
	err := s.db.QueryRow("SELECT next_block FROM index_cursors WHERE source=?", source).Scan(&n)
	if err == sql.ErrNoRows {
		return initial, nil
	}
	return n, err
}

// ApplyBatch commits event identities, financial rows, snapshots and the cursor
// atomically. Retries cannot double-count funds, even across process restarts.
func (s *Store) ApplyBatch(source string, next uint64, events []ChainEvent, snapshots []*Tournament) error {
	return s.CommitBatch(source, next, events, snapshots, "")
}
func (s *Store) Checkpoint(source string) (uint64, string, error) {
	var height uint64
	var hash string
	err := s.db.QueryRow("SELECT height,block_hash FROM index_checkpoints WHERE source=?", source).Scan(&height, &hash)
	if err == sql.ErrNoRows {
		return 0, "", nil
	}
	return height, hash, err
}
func (s *Store) CommitBatch(source string, next uint64, events []ChainEvent, snapshots []*Tournament, checkpoint string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, e := range events {
		b, err := json.Marshal(e)
		if err != nil {
			return err
		}
		r, err := tx.Exec("INSERT OR IGNORE INTO chain_events(id,chain_id,tournament_id,tx_hash,log_index,block_number,event_name,event_json) VALUES(?,?,?,?,?,?,?,?)", e.ID, e.ChainID, e.TournamentID, e.TxHash, e.LogIndex, e.BlockNumber, e.EventName, string(b))
		if err != nil {
			return err
		}
		n, err := r.RowsAffected()
		if err != nil {
			return err
		}
		if n == 0 {
			continue
		}
		kind := ""
		switch e.EventName {
		case "TournamentFunded":
			kind = "deposit"
		case "PrizeClaimed":
			kind = "prize"
		case "SponsorRefunded":
			kind = "refund"
		}
		if kind != "" {
			amount, ok := new(big.Int).SetString(e.AmountRaw, 10)
			if !ok || amount.Sign() < 0 {
				return fmt.Errorf("invalid confirmed event amount")
			}
			_, err = tx.Exec("INSERT INTO financial_ledger(event_id,tournament_id,kind,actor,amount_raw) VALUES(?,?,?,?,?)", e.ID, e.TournamentID, kind, e.Actor, e.AmountRaw)
			if err != nil {
				return err
			}
		}
	}
	for _, t := range snapshots {
		b, err := json.Marshal(t)
		if err != nil {
			return err
		}
		if _, err = tx.Exec("INSERT INTO tournaments(id,snapshot) VALUES(?,?) ON CONFLICT(id) DO UPDATE SET snapshot=excluded.snapshot", t.ID, string(b)); err != nil {
			return err
		}
	}
	if _, err = tx.Exec("INSERT INTO index_cursors(source,next_block) VALUES(?,?) ON CONFLICT(source) DO UPDATE SET next_block=excluded.next_block", source, next); err != nil {
		return err
	}
	if checkpoint != "" {
		if _, err = tx.Exec("INSERT INTO index_checkpoints(source,height,block_hash) VALUES(?,?,?) ON CONFLICT(source) DO UPDATE SET height=excluded.height,block_hash=excluded.block_hash", source, next-1, checkpoint); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) LedgerTotals() (map[string]string, int, int, error) {
	totals := map[string]*big.Int{"deposit": new(big.Int), "prize": new(big.Int), "refund": new(big.Int)}
	rows, err := s.db.Query("SELECT kind,amount_raw FROM financial_ledger")
	if err != nil {
		return nil, 0, 0, err
	}
	n := 0
	for rows.Next() {
		var kind, raw string
		if err := rows.Scan(&kind, &raw); err != nil {
			rows.Close()
			return nil, 0, 0, err
		}
		v, ok := new(big.Int).SetString(raw, 10)
		if !ok || totals[kind] == nil {
			rows.Close()
			return nil, 0, 0, fmt.Errorf("invalid ledger data")
		}
		totals[kind].Add(totals[kind], v)
		n++
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return nil, 0, 0, err
	}
	var count int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM chain_events").Scan(&count); err != nil {
		return nil, 0, 0, err
	}
	result := map[string]string{}
	for k, v := range totals {
		result[k] = v.String()
	}
	return result, n, count, nil
}

func (s *Store) SaveMatch(m *Match, canonical string) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = s.db.Exec("INSERT INTO matches(tournament_id,status,match_json,canonical_evidence) VALUES(?,?,?,?) ON CONFLICT(tournament_id) DO UPDATE SET status=excluded.status,match_json=excluded.match_json,canonical_evidence=CASE WHEN excluded.canonical_evidence='' THEN matches.canonical_evidence ELSE excluded.canonical_evidence END", m.TournamentID, m.Status, string(b), canonical)
	return err
}
func (s *Store) Match(id uint64) (*Match, string, error) {
	var encoded, canonical string
	err := s.db.QueryRow("SELECT match_json,canonical_evidence FROM matches WHERE tournament_id=?", id).Scan(&encoded, &canonical)
	if err == sql.ErrNoRows {
		return nil, "", ErrNotFound
	}
	if err != nil {
		return nil, "", err
	}
	var m Match
	err = json.Unmarshal([]byte(encoded), &m)
	return &m, canonical, err
}
func (s *Store) RecoverRuns() error {
	rows, err := s.db.Query("SELECT match_json FROM matches WHERE status IN ('queued','running')")
	if err != nil {
		return err
	}
	var recoverable []Match
	for rows.Next() {
		var encoded string
		if err := rows.Scan(&encoded); err != nil {
			rows.Close()
			return err
		}
		var m Match
		if err := json.Unmarshal([]byte(encoded), &m); err != nil {
			rows.Close()
			return err
		}
		recoverable = append(recoverable, m)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return err
	}
	for i := range recoverable {
		m := &recoverable[i]
		m.Status = "interrupted"
		m.Error = "Server restarted; organizer can resume the persisted round after chain validation."
		m.UpdatedAt = time.Now().Unix()
		if err := s.SaveMatch(m, ""); err != nil {
			return err
		}
	}
	return nil
}
