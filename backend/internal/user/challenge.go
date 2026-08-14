package user

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
)

// Challenge PVE 挑战记录
type Challenge struct {
	ID            string    `json:"id"`
	ChallengerID  string    `json:"challenger_id"`
	ChallengerAddr string   `json:"challenger_address"`
	OpponentID    string    `json:"opponent_id"`      // 对手 Agent ID（系统或用户）
	OpponentType  string    `json:"opponent_type"`    // "system" or "user"
	GameID        uint64    `json:"game_id"`          // 关联的游戏 ID
	Stake         uint64    `json:"stake"`            // 赌注（AC）
	CurrencyType  string    `json:"currency_type"`    // "ac" or "usdc"
	Winner        string    `json:"winner"`           // "challenger" or "opponent" or ""
	Reward        uint64    `json:"reward"`           // 奖励（AC）
	Status        string    `json:"status"`           // "pending", "playing", "finished"
	CreatedAt     time.Time `json:"created_at"`
	FinishedAt    time.Time `json:"finished_at"`
}

// ChallengeStore 挑战记录存储（SQLite）
type ChallengeStore struct {
	db *sql.DB
	mu sync.RWMutex
}

// NewChallengeStore 创建挑战存储
func NewChallengeStore(dbPath string) (*ChallengeStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// 创建表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS challenges (
			id TEXT PRIMARY KEY,
			challenger_id TEXT NOT NULL,
			challenger_address TEXT NOT NULL,
			opponent_id TEXT NOT NULL,
			opponent_type TEXT NOT NULL,
			game_id INTEGER NOT NULL,
			stake INTEGER NOT NULL,
			currency_type TEXT NOT NULL,
			winner TEXT,
			reward INTEGER,
			status TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			finished_at DATETIME
		);
		CREATE INDEX IF NOT EXISTS idx_challenges_game ON challenges(game_id);
		CREATE INDEX IF NOT EXISTS idx_challenges_user ON challenges(challenger_id);
	`)
	if err != nil {
		return nil, fmt.Errorf("create table: %w", err)
	}

	return &ChallengeStore{db: db}, nil
}

// Close 关闭数据库
func (s *ChallengeStore) Close() error {
	return s.db.Close()
}

// Create 创建挑战
func (s *ChallengeStore) Create(challengerID, challengerAddr, opponentID, opponentType string, gameID uint64, stake uint64, currencyType string) (*Challenge, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if opponentType != "system" && opponentType != "user" {
		return nil, fmt.Errorf("invalid opponent_type")
	}
	if currencyType != "ac" && currencyType != "usdc" {
		return nil, fmt.Errorf("invalid currency_type")
	}

	id := fmt.Sprintf("challenge_%d", time.Now().UnixNano())
	now := time.Now()

	_, err := s.db.Exec(
		`INSERT INTO challenges
		(id, challenger_id, challenger_address, opponent_id, opponent_type, game_id, stake, currency_type, winner, reward, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, '', 0, 'pending', ?)`,
		id, challengerID, challengerAddr, opponentID, opponentType, gameID, stake, currencyType, now.Format(time.RFC3339),
	)
	if err != nil {
		return nil, err
	}

	return &Challenge{
		ID:             id,
		ChallengerID:   challengerID,
		ChallengerAddr: challengerAddr,
		OpponentID:     opponentID,
		OpponentType:   opponentType,
		GameID:         gameID,
		Stake:          stake,
		CurrencyType:   currencyType,
		Winner:         "",
		Reward:         0,
		Status:         "pending",
		CreatedAt:      now,
	}, nil
}

// Get 获取挑战
func (s *ChallengeStore) Get(id string) (*Challenge, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	row := s.db.QueryRow(
		`SELECT id, challenger_id, challenger_address, opponent_id, opponent_type, game_id, stake, currency_type, winner, reward, status, created_at, finished_at
		FROM challenges WHERE id = ?`,
		id,
	)

	var c Challenge
	var createdAt, finishedAt sql.NullString
	err := row.Scan(&c.ID, &c.ChallengerID, &c.ChallengerAddr, &c.OpponentID, &c.OpponentType, &c.GameID, &c.Stake, &c.CurrencyType, &c.Winner, &c.Reward, &c.Status, &createdAt, &finishedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("challenge %s not found", id)
	}
	if err != nil {
		return nil, err
	}

	if createdAt.Valid {
		c.CreatedAt, _ = time.Parse(time.RFC3339, createdAt.String)
	}
	if finishedAt.Valid {
		c.FinishedAt, _ = time.Parse(time.RFC3339, finishedAt.String)
	}

	return &c, nil
}

// GetByGameID 根据游戏 ID 获取挑战
func (s *ChallengeStore) GetByGameID(gameID uint64) (*Challenge, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	row := s.db.QueryRow(
		`SELECT id, challenger_id, challenger_address, opponent_id, opponent_type, game_id, stake, currency_type, winner, reward, status, created_at, finished_at
		FROM challenges WHERE game_id = ?`,
		gameID,
	)

	var c Challenge
	var createdAt, finishedAt sql.NullString
	err := row.Scan(&c.ID, &c.ChallengerID, &c.ChallengerAddr, &c.OpponentID, &c.OpponentType, &c.GameID, &c.Stake, &c.CurrencyType, &c.Winner, &c.Reward, &c.Status, &createdAt, &finishedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("challenge for game %d not found", gameID)
	}
	if err != nil {
		return nil, err
	}

	if createdAt.Valid {
		c.CreatedAt, _ = time.Parse(time.RFC3339, createdAt.String)
	}
	if finishedAt.Valid {
		c.FinishedAt, _ = time.Parse(time.RFC3339, finishedAt.String)
	}

	return &c, nil
}

// GetByUser 获取用户的所有挑战
func (s *ChallengeStore) GetByUser(userID string) ([]*Challenge, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query(
		`SELECT id, challenger_id, challenger_address, opponent_id, opponent_type, game_id, stake, currency_type, winner, reward, status, created_at, finished_at
		FROM challenges WHERE challenger_id = ? ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var challenges []*Challenge
	for rows.Next() {
		var c Challenge
		var createdAt, finishedAt sql.NullString
		if err := rows.Scan(&c.ID, &c.ChallengerID, &c.ChallengerAddr, &c.OpponentID, &c.OpponentType, &c.GameID, &c.Stake, &c.CurrencyType, &c.Winner, &c.Reward, &c.Status, &createdAt, &finishedAt); err != nil {
			return nil, err
		}
		if createdAt.Valid {
			c.CreatedAt, _ = time.Parse(time.RFC3339, createdAt.String)
		}
		if finishedAt.Valid {
			c.FinishedAt, _ = time.Parse(time.RFC3339, finishedAt.String)
		}
		challenges = append(challenges, &c)
	}

	return challenges, nil
}

// GetByAddress 按钱包地址查询挑战记录（用于钱包/收益页）
func (s *ChallengeStore) GetByAddress(address string) ([]*Challenge, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query(
		`SELECT id, challenger_id, challenger_address, opponent_id, opponent_type, game_id, stake, currency_type, winner, reward, status, created_at, finished_at
		FROM challenges WHERE LOWER(challenger_address) = LOWER(?) ORDER BY created_at DESC LIMIT 50`,
		address,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var challenges []*Challenge
	for rows.Next() {
		var c Challenge
		var createdAt, finishedAt sql.NullString
		if err := rows.Scan(&c.ID, &c.ChallengerID, &c.ChallengerAddr, &c.OpponentID, &c.OpponentType, &c.GameID, &c.Stake, &c.CurrencyType, &c.Winner, &c.Reward, &c.Status, &createdAt, &finishedAt); err != nil {
			return nil, err
		}
		if createdAt.Valid {
			c.CreatedAt, _ = time.Parse(time.RFC3339, createdAt.String)
		}
		if finishedAt.Valid {
			c.FinishedAt, _ = time.Parse(time.RFC3339, finishedAt.String)
		}
		challenges = append(challenges, &c)
	}

	return challenges, nil
}

// UpdateStatus 更新挑战状态
func (s *ChallengeStore) UpdateStatus(id, status, winner string, reward uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var finishedAt sql.NullString
	if status == "finished" {
		finishedAt = sql.NullString{String: time.Now().Format(time.RFC3339), Valid: true}
	}

	result, err := s.db.Exec(
		"UPDATE challenges SET status = ?, winner = ?, reward = ?, finished_at = ? WHERE id = ?",
		status, winner, reward, finishedAt, id,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("challenge %s not found", id)
	}
	return nil
}

// ListActive 列出所有进行中的挑战
func (s *ChallengeStore) ListActive() ([]*Challenge, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query(
		`SELECT id, challenger_id, challenger_address, opponent_id, opponent_type, game_id, stake, currency_type, winner, reward, status, created_at, finished_at
		FROM challenges WHERE status IN ('pending', 'playing') ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var challenges []*Challenge
	for rows.Next() {
		var c Challenge
		var createdAt, finishedAt sql.NullString
		if err := rows.Scan(&c.ID, &c.ChallengerID, &c.ChallengerAddr, &c.OpponentID, &c.OpponentType, &c.GameID, &c.Stake, &c.CurrencyType, &c.Winner, &c.Reward, &c.Status, &createdAt, &finishedAt); err != nil {
			return nil, err
		}
		if createdAt.Valid {
			c.CreatedAt, _ = time.Parse(time.RFC3339, createdAt.String)
		}
		if finishedAt.Valid {
			c.FinishedAt, _ = time.Parse(time.RFC3339, finishedAt.String)
		}
		challenges = append(challenges, &c)
	}

	return challenges, nil
}
