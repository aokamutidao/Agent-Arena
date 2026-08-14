package user

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
)

// GameHistory 已完成对局的记录（持久化）
type GameHistory struct {
	GameID         uint64    `json:"game_id"`
	RedName        string    `json:"red_name"`
	BlueName       string    `json:"blue_name"`
	Winner         string    `json:"winner"`          // "red" / "blue" / "draw"
	WinnerName     string    `json:"winner_name"`
	RedHP          uint8     `json:"final_hp_red"`
	BlueHP         uint8     `json:"final_hp_blue"`
	TotalRounds    uint8     `json:"total_rounds"`
	FinishTxHash   string    `json:"finish_tx_hash,omitempty"` // 链上 FinishGame 交易
	CreatedAt      time.Time `json:"created_at"`
}

// GameHistoryStore 对局历史存储（SQLite）
type GameHistoryStore struct {
	db *sql.DB
	mu sync.RWMutex
}

// NewGameHistoryStore 创建对局历史存储
func NewGameHistoryStore(dbPath string) (*GameHistoryStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS game_history (
			game_id       INTEGER PRIMARY KEY,
			red_name      TEXT NOT NULL,
			blue_name     TEXT NOT NULL,
			winner        TEXT NOT NULL,
			winner_name   TEXT NOT NULL,
			final_hp_red  INTEGER NOT NULL,
			final_hp_blue INTEGER NOT NULL,
			total_rounds  INTEGER NOT NULL,
			finish_tx_hash TEXT,
			created_at    DATETIME NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_game_history_created ON game_history(created_at DESC);
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("create table: %w", err)
	}

	return &GameHistoryStore{db: db}, nil
}

// Record 记录一局完成的对局（可重复调用；finish_tx_hash 通过 UpdateFinishTx 更新）
func (s *GameHistoryStore) Record(g GameHistory) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`
		INSERT OR REPLACE INTO game_history
		(game_id, red_name, blue_name, winner, winner_name, final_hp_red, final_hp_blue, total_rounds, finish_tx_hash, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		g.GameID, g.RedName, g.BlueName, g.Winner, g.WinnerName,
		g.RedHP, g.BlueHP, g.TotalRounds, g.FinishTxHash, g.CreatedAt,
	)
	return err
}

// UpdateFinishTx 更新链上 FinishGame 交易哈希
func (s *GameHistoryStore) UpdateFinishTx(gameID uint64, txHash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`UPDATE game_history SET finish_tx_hash = ? WHERE game_id = ?`,
		txHash, gameID)
	return err
}

// List 返回最近的对局历史（按创建时间倒序）
func (s *GameHistoryStore) List(limit, offset int) ([]GameHistory, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := s.db.Query(`
		SELECT game_id, red_name, blue_name, winner, winner_name,
		       final_hp_red, final_hp_blue, total_rounds, finish_tx_hash, created_at
		FROM game_history
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []GameHistory
	for rows.Next() {
		var g GameHistory
		var txHash sql.NullString
		if err := rows.Scan(&g.GameID, &g.RedName, &g.BlueName, &g.Winner, &g.WinnerName,
			&g.RedHP, &g.BlueHP, &g.TotalRounds, &txHash, &g.CreatedAt); err != nil {
			return nil, err
		}
		if txHash.Valid {
			g.FinishTxHash = txHash.String
		}
		out = append(out, g)
	}
	if out == nil {
		out = []GameHistory{}
	}
	return out, rows.Err()
}

// Get 获取单条记录
func (s *GameHistoryStore) Get(gameID uint64) (*GameHistory, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	row := s.db.QueryRow(`
		SELECT game_id, red_name, blue_name, winner, winner_name,
		       final_hp_red, final_hp_blue, total_rounds, finish_tx_hash, created_at
		FROM game_history WHERE game_id = ?`, gameID)

	var g GameHistory
	var txHash sql.NullString
	if err := row.Scan(&g.GameID, &g.RedName, &g.BlueName, &g.Winner, &g.WinnerName,
		&g.RedHP, &g.BlueHP, &g.TotalRounds, &txHash, &g.CreatedAt); err != nil {
		return nil, err
	}
	if txHash.Valid {
		g.FinishTxHash = txHash.String
	}
	return &g, nil
}

// Close 关闭 DB
func (s *GameHistoryStore) Close() error {
	return s.db.Close()
}
