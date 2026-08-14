package user

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// User 用户信息
type User struct {
	ID          string     `json:"id"`
	Address     string     `json:"address"` // 钱包地址（小写）
	Username    string     `json:"username"`
	ACBalance   uint64     `json:"ac_balance"`    // Arena Coins 余额
	LastClaimAt *time.Time `json:"last_claim_at"` // 上次领取时间
	CreatedAt   time.Time  `json:"created_at"`
}

// Store 用户存储（SQLite 实现）
type Store struct {
	db *sql.DB
	mu sync.RWMutex
}

// NewStore 创建用户存储
func NewStore(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// 创建表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			address TEXT UNIQUE NOT NULL,
			username TEXT NOT NULL,
			ac_balance INTEGER NOT NULL DEFAULT 0,
			last_claim_at DATETIME,
			created_at DATETIME NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_users_address ON users(address);
	`)
	if err != nil {
		return nil, fmt.Errorf("create table: %w", err)
	}

	return &Store{db: db}, nil
}

// Close 关闭数据库
func (s *Store) Close() error {
	return s.db.Close()
}

// GetByAddress 根据地址获取用户
func (s *Store) GetByAddress(address string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	addr := normalizeAddress(address)
	row := s.db.QueryRow(
		"SELECT id, address, username, ac_balance, last_claim_at, created_at FROM users WHERE address = ?",
		addr,
	)

	var user User
	var createdAt string
	var lastClaimAt sql.NullString
	err := row.Scan(&user.ID, &user.Address, &user.Username, &user.ACBalance, &lastClaimAt, &createdAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user %s not found", address)
	}
	if err != nil {
		return nil, err
	}

	user.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	if lastClaimAt.Valid {
		t, _ := time.Parse(time.RFC3339, lastClaimAt.String)
		user.LastClaimAt = &t
	}
	return &user, nil
}

// GetByID 根据 ID 获取用户
func (s *Store) GetByID(id string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	row := s.db.QueryRow(
		"SELECT id, address, username, ac_balance, last_claim_at, created_at FROM users WHERE id = ?",
		id,
	)

	var user User
	var createdAt string
	var lastClaimAt sql.NullString
	err := row.Scan(&user.ID, &user.Address, &user.Username, &user.ACBalance, &lastClaimAt, &createdAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user %s not found", id)
	}
	if err != nil {
		return nil, err
	}

	user.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	if lastClaimAt.Valid {
		t, _ := time.Parse(time.RFC3339, lastClaimAt.String)
		user.LastClaimAt = &t
	}
	return &user, nil
}

// Create 创建用户
func (s *Store) Create(address string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	addr := normalizeAddress(address)

	// 检查是否已存在
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE address = ?", addr).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("user %s already exists", address)
	}

	// 生成用户 ID（简单用地址前 8 位）
	id := fmt.Sprintf("u_%s", addr[2:10])
	now := time.Now()

	_, err = s.db.Exec(
		"INSERT INTO users (id, address, username, ac_balance, created_at) VALUES (?, ?, ?, 0, ?)",
		id, addr, fmt.Sprintf("Player_%s", addr[2:6]), now.Format(time.RFC3339),
	)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        id,
		Address:   addr,
		Username:  fmt.Sprintf("Player_%s", addr[2:6]),
		ACBalance: 0,
		CreatedAt: now,
	}, nil
}

// GetOrCreate 获取或创建用户
func (s *Store) GetOrCreate(address string) (*User, error) {
	user, err := s.GetByAddress(address)
	if err == nil {
		return user, nil
	}
	return s.Create(address)
}

// UpdateUsername 更新用户名
func (s *Store) UpdateUsername(address, username string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	addr := normalizeAddress(address)
	result, err := s.db.Exec("UPDATE users SET username = ? WHERE address = ?", username, addr)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user %s not found", address)
	}
	return nil
}

// UpdateACBalance 更新 AC 余额
func (s *Store) UpdateACBalance(address string, balance uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	addr := normalizeAddress(address)
	result, err := s.db.Exec("UPDATE users SET ac_balance = ? WHERE address = ?", balance, addr)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user %s not found", address)
	}
	return nil
}

// AddACBalance 增加 AC 余额
func (s *Store) AddACBalance(address string, amount uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	addr := normalizeAddress(address)
	result, err := s.db.Exec("UPDATE users SET ac_balance = ac_balance + ? WHERE address = ?", amount, addr)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user %s not found", address)
	}
	return nil
}

// UpdateLastClaimAt 更新上次领取时间
func (s *Store) UpdateLastClaimAt(address string, claimTime time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	addr := normalizeAddress(address)
	result, err := s.db.Exec("UPDATE users SET last_claim_at = ? WHERE address = ?", claimTime.Format(time.RFC3339), addr)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user %s not found", address)
	}
	return nil
}

// CanClaimDaily 检查是否可以领取每日 AC（24小时冷却）
func (s *Store) CanClaimDaily(address string) (bool, error) {
	user, err := s.GetByAddress(address)
	if err != nil {
		return false, err
	}

	if user.LastClaimAt == nil {
		return true, nil // 从未领取过
	}

	// 检查是否超过 24 小时
	hoursSinceLastClaim := time.Since(*user.LastClaimAt).Hours()
	return hoursSinceLastClaim >= 24, nil
}

// normalizeAddress 标准化地址（小写 + 0x 前缀）
func normalizeAddress(address string) string {
	addr := strings.ToLower(address)
	if len(addr) >= 2 && addr[:2] == "0x" {
		return addr
	}
	return "0x" + addr
}
