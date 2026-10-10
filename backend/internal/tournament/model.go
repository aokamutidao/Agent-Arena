// Package tournament implements the course DApp's read-only blockchain backend.
// Wallets, never this service, authorize all financial transactions.
package tournament

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"agent-arena/backend/internal/agent"
	"agent-arena/backend/internal/engine"
	"github.com/ethereum/go-ethereum/common"
)

var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("conflict")

type Config struct {
	DatabasePath      string
	JWTSecret         string
	AuthDomain        string
	ChainID           uint64
	ManagerAddress    string
	EscrowAddress     string
	RegistryAddress   string
	ReviewAddress     string
	TokenAddress      string
	TokenDecimals     uint8
	DeploymentBlock   uint64
	Confirmations     uint64
	CORSOrigins       []string
	TrustedProxyCIDRs []string
	PollInterval      time.Duration
	RoundInterval     time.Duration
	SimulationEnabled bool
	LLMProvider       string
	LLMModel          string
	LLMTimeoutSeconds int64
	LLMMaxCalls       int
	ABIPath           string
}

func (c *Config) Defaults() {
	if c.DatabasePath == "" {
		c.DatabasePath = "course.db"
	}
	if c.AuthDomain == "" {
		c.AuthDomain = "localhost:3000"
	}
	if c.ChainID == 0 {
		c.ChainID = 11155111
	}
	if c.PollInterval <= 0 {
		c.PollInterval = 5 * time.Second
	}
	if c.ABIPath == "" {
		c.ABIPath = "../contracts/abi"
	}
	if c.TokenDecimals == 0 {
		c.TokenDecimals = 6
	}
	if c.LLMMaxCalls <= 0 {
		c.LLMMaxCalls = 120
	}
}

func (c Config) HasDeployment() bool {
	for _, a := range []string{c.ManagerAddress, c.EscrowAddress, c.RegistryAddress, c.ReviewAddress, c.TokenAddress} {
		if !common.IsHexAddress(a) || common.HexToAddress(a) == (common.Address{}) {
			return false
		}
	}
	return true
}

type User struct {
	Address string   `json:"address"`
	Roles   []string `json:"roles"`
}
type ParticipantAgent struct {
	Address   string `json:"address"`
	AgentHash string `json:"agent_hash"`
}
type Tournament struct {
	ID                   uint64             `json:"id"`
	Title                string             `json:"title"`
	MetadataURI          string             `json:"metadata_uri"`
	Organizer            string             `json:"organizer"`
	Auditor              string             `json:"auditor"`
	Token                string             `json:"token"`
	EntryFee             string             `json:"entry_fee"`
	PrizePool            string             `json:"prize_pool"`
	PrizePaid            string             `json:"prize_paid"`
	PrizeClaimed         bool               `json:"prize_claimed"`
	MaxParticipants      int                `json:"max_participants"`
	RegistrationDeadline int64              `json:"registration_deadline"`
	SettlementDeadline   int64              `json:"settlement_deadline"`
	ReviewPeriod         int64              `json:"review_period"`
	ReviewDeadline       int64              `json:"review_deadline"`
	Status               string             `json:"status"`
	Participants         []string           `json:"participants"`
	ParticipantAgents    []ParticipantAgent `json:"participant_agents"`
	ProposedWinner       string             `json:"proposed_winner"`
	Winner               string             `json:"winner"`
	EvidenceHash         string             `json:"evidence_hash"`
	ChainUpdatedAt       int64              `json:"chain_updated_at"`
	ReadBlock            uint64             `json:"read_block"`
	Confirmed            bool               `json:"confirmed"`
}

type ChainEvent struct {
	ID              string          `json:"id"`
	ChainID         uint64          `json:"chain_id"`
	TournamentID    uint64          `json:"tournament_id"`
	ContractAddress string          `json:"contract_address"`
	EventName       string          `json:"event_name"`
	TxHash          string          `json:"tx_hash"`
	LogIndex        uint            `json:"log_index"`
	BlockNumber     uint64          `json:"block_number"`
	Timestamp       int64           `json:"timestamp"`
	Actor           string          `json:"actor"`
	AmountRaw       string          `json:"amount_raw"`
	Data            json.RawMessage `json:"data"`
}

// ChainReader exposes no transaction signing or sending methods.
type ChainReader interface {
	Tournament(context.Context, uint64) (*Tournament, error)
	Roles(context.Context, string) ([]string, error)
	Head(context.Context) (uint64, error)
	Events(context.Context, uint64, uint64) ([]ChainEvent, error)
}
type PinnedTournamentReader interface {
	TournamentAt(context.Context, uint64, uint64) (*Tournament, error)
}
type CheckpointReader interface {
	BlockHash(context.Context, uint64) (string, error)
}

type AgentDefinition struct {
	SchemaVersion int    `json:"schema_version"`
	OwnerAddress  string `json:"owner_address"`
	Name          string `json:"name"`
	Personality   string `json:"personality"`
}
type AgentProfile struct {
	ID           string `json:"id"`
	AgentHash    string `json:"agent_hash"`
	Name         string `json:"name"`
	Personality  string `json:"personality"`
	OwnerAddress string `json:"owner_address"`
	CreatedAt    int64  `json:"created_at"`
}

func (a AgentProfile) Definition() AgentDefinition {
	return AgentDefinition{1, a.OwnerAddress, a.Name, a.Personality}
}

type Round struct {
	Round         uint8             `json:"round"`
	RedAction     engine.Action     `json:"red_action"`
	BlueAction    engine.Action     `json:"blue_action"`
	RedHP         uint8             `json:"red_hp"`
	BlueHP        uint8             `json:"blue_hp"`
	RedPosition   engine.Position   `json:"red_position"`
	BluePosition  engine.Position   `json:"blue_position"`
	RedReasoning  string            `json:"red_reasoning"`
	BlueReasoning string            `json:"blue_reasoning"`
	RedState      engine.AgentState `json:"red_state"`
	BlueState     engine.AgentState `json:"blue_state"`
	RedLLMCallID  int               `json:"red_llm_call_id,omitempty"`
	BlueLLMCallID int               `json:"blue_llm_call_id,omitempty"`
}

// DecisionCall is the persisted provider attempt journal, including retries.
type DecisionCall struct {
	Sequence int    `json:"sequence"`
	Round    uint8  `json:"round"`
	Side     string `json:"side"`
	agent.CallMetadata
}

// UsageSummary keeps unknown aggregate token usage null rather than zero.
type UsageSummary struct {
	PromptTokens     *uint64 `json:"prompt_tokens"`
	CompletionTokens *uint64 `json:"completion_tokens"`
	TotalTokens      *uint64 `json:"total_tokens"`
}

type LLMSummary struct {
	Provider      string       `json:"provider"`
	Model         string       `json:"model"`
	Attempted     int          `json:"attempted"`
	Succeeded     int          `json:"succeeded"`
	Failed        int          `json:"failed"`
	LatencyMS     int64        `json:"latency_ms"`
	Usage         UsageSummary `json:"usage"`
	LastErrorCode string       `json:"last_error_code"`
}

type NoteEvidence struct {
	SchemaVersion  int    `json:"schema_version"`
	ChainID        uint64 `json:"chain_id"`
	ManagerAddress string `json:"manager_address"`
	TournamentID   uint64 `json:"tournament_id"`
	AuthorAddress  string `json:"author_address"`
	Kind           string `json:"kind"`
	Text           string `json:"text"`
}
type Note struct {
	NoteEvidence
	EvidenceHash  string `json:"evidence_hash"`
	CreatedAt     int64  `json:"created_at"`
	CanonicalJSON string `json:"canonical_json"`
}
type Match struct {
	TournamentID  uint64           `json:"tournament_id"`
	Status        string           `json:"status"`
	Mode          string           `json:"mode"`
	RedAddress    string           `json:"red_address"`
	BlueAddress   string           `json:"blue_address"`
	CurrentRound  uint8            `json:"current_round"`
	WinnerAddress string           `json:"winner_address"`
	Outcome       string           `json:"outcome"`
	EvidenceHash  string           `json:"evidence_hash"`
	Rounds        []Round          `json:"rounds"`
	State         engine.GameState `json:"state"`
	StartedAt     int64            `json:"started_at"`
	UpdatedAt     int64            `json:"updated_at"`
	Error         string           `json:"error"`
	RedAgent      AgentDefinition  `json:"red_agent"`
	BlueAgent     AgentDefinition  `json:"blue_agent"`
	LLMCalls      []DecisionCall   `json:"llm_calls,omitempty"`
	LLMSummary    *LLMSummary      `json:"llm_summary,omitempty"`
	RetryAllowed  bool             `json:"retry_allowed"`
}
type Evidence struct {
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
	LLMCalls       []DecisionCall   `json:"llm_calls,omitempty"`
}

func address(a string) string { return strings.ToLower(common.HexToAddress(a).Hex()) }
func nonzeroAddress(a common.Address) string {
	if a == (common.Address{}) {
		return ""
	}
	return strings.ToLower(a.Hex())
}
