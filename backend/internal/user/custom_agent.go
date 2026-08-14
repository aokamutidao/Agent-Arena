package user

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
)

// CustomAgent 用户自定义 Agent
type CustomAgent struct {
	ID            string    `json:"id"`
	OwnerID       string    `json:"owner_id"`
	OwnerAddress  string    `json:"owner_address"`
	Name          string    `json:"name"`
	Personality   string    `json:"personality"`   // 性格描述（≤500 字）
	APIEndpoint   string    `json:"api_endpoint"`  // 必填：外部 AI API endpoint
	APIKey        string    `json:"api_key"`       // 必填：API key（加密存储）
	Model         string    `json:"model"`         // 必填：模型名称（如 gpt-4, gpt-3.5-turbo）
	ChallengeFee  uint64    `json:"challenge_fee"` // 挑战费用（AC）
	CurrencyType  string    `json:"currency_type"` // "ac" 或 "usdc"
	IsListed      bool      `json:"is_listed"`     // 是否上架到市场
	Wins          uint32    `json:"wins"`
	Losses        uint32    `json:"losses"`
	CreatedAt     time.Time `json:"created_at"`
}

// AgentStore 自定义 Agent 存储（SQLite）
type AgentStore struct {
	db *sql.DB
	mu sync.RWMutex
}

// NewAgentStore 创建 Agent 存储
func NewAgentStore(dbPath string) (*AgentStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// 创建表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS custom_agents (
			id TEXT PRIMARY KEY,
			owner_id TEXT NOT NULL,
			owner_address TEXT NOT NULL,
			name TEXT NOT NULL,
			personality TEXT NOT NULL,
			api_endpoint TEXT NOT NULL,
			api_key TEXT NOT NULL,
			model TEXT NOT NULL DEFAULT 'gpt-3.5-turbo',
			challenge_fee INTEGER NOT NULL,
			currency_type TEXT NOT NULL,
			is_listed INTEGER NOT NULL DEFAULT 0,
			wins INTEGER NOT NULL DEFAULT 0,
			losses INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_agents_owner ON custom_agents(owner_id);
	`)
	if err != nil {
		return nil, fmt.Errorf("create table: %w", err)
	}

	// 迁移：为旧表添加 model 列（如果不存在）
	_, _ = db.Exec(`ALTER TABLE custom_agents ADD COLUMN model TEXT NOT NULL DEFAULT 'gpt-3.5-turbo'`)

	return &AgentStore{db: db}, nil
}

// Close 关闭数据库
func (s *AgentStore) Close() error {
	return s.db.Close()
}

// Create 创建自定义 Agent
func (s *AgentStore) Create(ownerID, ownerAddress, name, personality, apiEndpoint, apiKey, model string, challengeFee uint64, currencyType string) (*CustomAgent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 验证
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if len(name) > 50 {
		return nil, fmt.Errorf("name too long (max 50 chars)")
	}
	if personality == "" {
		return nil, fmt.Errorf("personality is required")
	}
	if len(personality) > 500 {
		return nil, fmt.Errorf("personality too long (max 500 chars)")
	}
	if currencyType != "ac" && currencyType != "usdc" {
		return nil, fmt.Errorf("currency_type must be 'ac' or 'usdc'")
	}
	// 自定义 Agent 必须提供外部 AI API
	if apiEndpoint == "" {
		return nil, fmt.Errorf("api_endpoint is required for custom agents")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("api_key is required for custom agents")
	}
	if model == "" {
		model = "gpt-3.5-turbo" // 默认模型
	}

	id := fmt.Sprintf("agent_%d", time.Now().UnixNano())
	now := time.Now()

	_, err := s.db.Exec(
		`INSERT INTO custom_agents
		(id, owner_id, owner_address, name, personality, api_endpoint, api_key, model, challenge_fee, currency_type, is_listed, wins, losses, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, 0, 0, ?)`,
		id, ownerID, ownerAddress, name, personality, apiEndpoint, apiKey, model, challengeFee, currencyType, now.Format(time.RFC3339),
	)
	if err != nil {
		return nil, err
	}

	return &CustomAgent{
		ID:           id,
		OwnerID:      ownerID,
		OwnerAddress: ownerAddress,
		Name:         name,
		Personality:  personality,
		APIEndpoint:  apiEndpoint,
		APIKey:       apiKey,
		Model:        model,
		ChallengeFee: challengeFee,
		CurrencyType: currencyType,
		IsListed:     false,
		Wins:         0,
		Losses:       0,
		CreatedAt:    now,
	}, nil
}

// Get 获取 Agent
func (s *AgentStore) Get(id string) (*CustomAgent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	row := s.db.QueryRow(
		`SELECT id, owner_id, owner_address, name, personality, api_endpoint, api_key, model, challenge_fee, currency_type, is_listed, wins, losses, created_at
		FROM custom_agents WHERE id = ?`,
		id,
	)

	var a CustomAgent
	var createdAt string
	var isListed int
	err := row.Scan(&a.ID, &a.OwnerID, &a.OwnerAddress, &a.Name, &a.Personality, &a.APIEndpoint, &a.APIKey, &a.Model, &a.ChallengeFee, &a.CurrencyType, &isListed, &a.Wins, &a.Losses, &createdAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("agent %s not found", id)
	}
	if err != nil {
		return nil, err
	}

	a.IsListed = isListed > 0
	a.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return &a, nil
}

// GetByOwner 获取用户的所有 Agent
func (s *AgentStore) GetByOwner(ownerID string) ([]*CustomAgent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query(
		`SELECT id, owner_id, owner_address, name, personality, api_endpoint, api_key, model, challenge_fee, currency_type, is_listed, wins, losses, created_at
		FROM custom_agents WHERE owner_id = ? ORDER BY created_at DESC`,
		ownerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []*CustomAgent
	for rows.Next() {
		var a CustomAgent
		var createdAt string
		var isListed int
		if err := rows.Scan(&a.ID, &a.OwnerID, &a.OwnerAddress, &a.Name, &a.Personality, &a.APIEndpoint, &a.APIKey, &a.Model, &a.ChallengeFee, &a.CurrencyType, &isListed, &a.Wins, &a.Losses, &createdAt); err != nil {
			return nil, err
		}
		a.IsListed = isListed > 0
		a.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		agents = append(agents, &a)
	}

	return agents, nil
}

// List 列出所有上架的 Agent
func (s *AgentStore) List() ([]*CustomAgent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query(
		`SELECT id, owner_id, owner_address, name, personality, api_endpoint, api_key, model, challenge_fee, currency_type, is_listed, wins, losses, created_at
		FROM custom_agents WHERE is_listed = 1 ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []*CustomAgent
	for rows.Next() {
		var a CustomAgent
		var createdAt string
		var isListed int
		if err := rows.Scan(&a.ID, &a.OwnerID, &a.OwnerAddress, &a.Name, &a.Personality, &a.APIEndpoint, &a.APIKey, &a.Model, &a.ChallengeFee, &a.CurrencyType, &isListed, &a.Wins, &a.Losses, &createdAt); err != nil {
			return nil, err
		}
		a.IsListed = isListed > 0
		a.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		agents = append(agents, &a)
	}

	return agents, nil
}

// Update 更新 Agent
func (s *AgentStore) Update(id, name, personality, apiEndpoint, apiKey, model string, challengeFee uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	result, err := s.db.Exec(
		"UPDATE custom_agents SET name = ?, personality = ?, api_endpoint = ?, api_key = ?, model = ?, challenge_fee = ? WHERE id = ?",
		name, personality, apiEndpoint, apiKey, model, challengeFee, id,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("agent %s not found", id)
	}
	return nil
}

// SetListed 设置上架状态
func (s *AgentStore) SetListed(id string, listed bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	isListed := 0
	if listed {
		isListed = 1
	}
	result, err := s.db.Exec("UPDATE custom_agents SET is_listed = ? WHERE id = ?", isListed, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("agent %s not found", id)
	}
	return nil
}

// UpdateStats 更新胜负统计
func (s *AgentStore) UpdateStats(id string, won bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var query string
	if won {
		query = "UPDATE custom_agents SET wins = wins + 1 WHERE id = ?"
	} else {
		query = "UPDATE custom_agents SET losses = losses + 1 WHERE id = ?"
	}
	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("agent %s not found", id)
	}
	return nil
}

// Delete 删除 Agent
func (s *AgentStore) Delete(id, ownerID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	result, err := s.db.Exec("DELETE FROM custom_agents WHERE id = ? AND owner_id = ?", id, ownerID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("agent %s not found or not authorized", id)
	}
	return nil
}
