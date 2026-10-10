package agent

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"agent-arena/backend/internal/engine"
)

// TokenUsage is provider-reported usage. Missing counts are unknown, not zero.
type TokenUsage struct {
	PromptTokens     *uint64 `json:"prompt_tokens,omitempty"`
	CompletionTokens *uint64 `json:"completion_tokens,omitempty"`
	TotalTokens      *uint64 `json:"total_tokens,omitempty"`
}

// CallMetadata contains only safe operational facts, never credentials or CoT.
type CallMetadata struct {
	Provider   string      `json:"provider"`
	Model      string      `json:"model"`
	Success    bool        `json:"success"`
	LatencyMS  int64       `json:"latency_ms"`
	HTTPStatus int         `json:"http_status,omitempty"`
	RequestID  string      `json:"request_id,omitempty"`
	Usage      *TokenUsage `json:"usage,omitempty"`
	ErrorCode  string      `json:"error_code,omitempty"`
}

// LLMResponse keeps the final response separate from persistable call metadata.
type LLMResponse struct {
	Content  string
	Metadata CallMetadata
}

// DetailedLLMClient extends the legacy Chat interface without breaking callers.
type DetailedLLMClient interface {
	LLMClient
	ChatWithMetadata(context.Context, string) (LLMResponse, error)
}

// LLMInfo is safe public configuration; endpoint and credentials are excluded.
type LLMInfo struct {
	Provider       string
	Model          string
	TimeoutSeconds int64
}

type LLMDescriber interface {
	LLMInfo() LLMInfo
}

// QwenConfig is supplied by the server operator, never by a user HTTP request.
type QwenConfig struct {
	APIKey   string
	Provider string
	Model    string
	Endpoint string
	Timeout  time.Duration
}

// CourseLLMConfig reads explicitly exported settings; it never loads .env files.
// LLM_* settings take precedence over backwards-compatible QWEN_* aliases.
func CourseLLMConfig(getenv func(string) string) (QwenConfig, error) {
	value := func(primary, alias, fallback string) string {
		if v := strings.TrimSpace(getenv(primary)); v != "" {
			return v
		}
		if alias != "" {
			if v := strings.TrimSpace(getenv(alias)); v != "" {
				return v
			}
		}
		return fallback
	}
	seconds, err := strconv.ParseUint(
		value("LLM_TIMEOUT_SECONDS", "QWEN_TIMEOUT_SECONDS", "30"), 10, 16,
	)
	if err != nil || seconds < 1 || seconds > 120 {
		return QwenConfig{}, fmt.Errorf("LLM_TIMEOUT_SECONDS must be between 1 and 120")
	}
	cfg := QwenConfig{
		APIKey:   value("LLM_API_KEY", "QWEN_API_KEY", ""),
		Provider: value("LLM_PROVIDER", "", "qwen"),
		Model:    value("LLM_MODEL", "QWEN_MODEL", "qwen-plus"),
		Endpoint: value("LLM_ENDPOINT", "QWEN_ENDPOINT",
			"https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"),
		Timeout: time.Duration(seconds) * time.Second,
	}
	return cfg, cfg.Validate()
}

var identifierPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:/-]{0,127}$`)

// Validate rejects credentials in URLs and remote plaintext HTTP endpoints.
func (c QwenConfig) Validate() error {
	if !identifierPattern.MatchString(c.Provider) || !identifierPattern.MatchString(c.Model) {
		return fmt.Errorf("LLM_PROVIDER and LLM_MODEL must be valid short identifiers")
	}
	for _, value := range []string{c.Provider, c.Model} {
		if strings.HasPrefix(strings.ToLower(value), "sk-") ||
			(c.APIKey != "" && strings.Contains(value, c.APIKey)) {
			return fmt.Errorf("LLM_PROVIDER and LLM_MODEL must not contain credentials")
		}
	}
	if c.APIKey != "" && strings.Contains(c.Endpoint, c.APIKey) {
		return fmt.Errorf("LLM_ENDPOINT must not contain credentials")
	}
	u, err := url.Parse(c.Endpoint)
	if err != nil || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return fmt.Errorf("LLM_ENDPOINT must be an absolute URL without credentials, query or fragment")
	}
	loopback := u.Hostname() == "localhost" || u.Hostname() == "127.0.0.1" || u.Hostname() == "::1"
	if u.Scheme != "https" && !(u.Scheme == "http" && loopback) {
		return fmt.Errorf("LLM_ENDPOINT requires HTTPS (HTTP is allowed only for local diagnostics)")
	}
	if c.Timeout <= 0 || c.Timeout > 120*time.Second {
		return fmt.Errorf("LLM timeout must be positive and no greater than 120 seconds")
	}
	return nil
}

// CallError exposes a fixed category instead of provider bodies or secret URLs.
type CallError struct{ Code string }

func (e *CallError) Error() string { return "LLM request failed: " + e.Code }

// ErrorCode deliberately discards arbitrary error messages from custom clients.
func ErrorCode(err error) string {
	var call *CallError
	if errors.As(err, &call) {
		return call.Code
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return "timeout"
	}
	if errors.Is(err, context.Canceled) {
		return "cancelled"
	}
	return "provider_error"
}

// DecisionResult retains the action intent and a short explanation, not CoT.
type DecisionResult struct {
	Action   engine.Action
	Note     string
	Metadata CallMetadata
	Err      error
}

var strictAction = regexp.MustCompile(`(?i)^(MOVE\((\d+),(\d+)\)|ATTACK|SKILL|CHARGE|WAIT|HEAL)$`)

// ParseStrictDecision requires one action and an optional single brief NOTE.
func ParseStrictDecision(content string) (engine.Action, string, error) {
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) < 1 || len(lines) > 2 || !strictAction.MatchString(strings.TrimSpace(lines[0])) {
		return engine.Action{}, "", &CallError{Code: "parse_error"}
	}
	action, err := ParseAction(strings.TrimSpace(lines[0]))
	if err != nil {
		return engine.Action{}, "", &CallError{Code: "parse_error"}
	}
	note := ""
	if len(lines) == 2 {
		line := strings.TrimSpace(lines[1])
		for _, prefix := range []string{"NOTE:", "Note:", "说明:", "说明："} {
			if strings.HasPrefix(line, prefix) {
				note = strings.TrimSpace(strings.TrimPrefix(line, prefix))
				break
			}
		}
		if note == "" || len([]rune(note)) > 80 {
			return engine.Action{}, "", &CallError{Code: "parse_error"}
		}
	}
	return action, note, nil
}

// DecideTurnStrict never replaces a failed model request/parse with WAIT.
func (a *Agent) DecideTurnStrict(ctx context.Context, input AgentDecisionInput) DecisionResult {
	start := time.Now()
	metadata := CallMetadata{Provider: "custom", Model: "unknown"}
	if described, ok := a.LLM.(LLMDescriber); ok {
		info := described.LLMInfo()
		metadata.Provider, metadata.Model = info.Provider, info.Model
	}
	if a.LLM == nil {
		metadata.ErrorCode = "not_configured"
		return DecisionResult{Metadata: metadata, Err: &CallError{Code: metadata.ErrorCode}}
	}
	prompt := buildDecisionPrompt(a.Name, a.Personality, input)
	var content string
	var err error
	if detailed, ok := a.LLM.(DetailedLLMClient); ok {
		var response LLMResponse
		response, err = detailed.ChatWithMetadata(ctx, prompt)
		content, metadata = response.Content, response.Metadata
	} else {
		content, err = a.LLM.Chat(ctx, prompt)
		metadata.LatencyMS = time.Since(start).Milliseconds()
	}
	if err != nil {
		metadata.Success, metadata.ErrorCode = false, ErrorCode(err)
		return DecisionResult{Metadata: metadata, Err: &CallError{Code: metadata.ErrorCode}}
	}
	action, note, err := ParseStrictDecision(content)
	if err != nil {
		metadata.Success, metadata.ErrorCode = false, "parse_error"
		return DecisionResult{Metadata: metadata, Err: &CallError{Code: "parse_error"}}
	}
	metadata.Success, metadata.ErrorCode = true, ""
	if !validateAgentAction(action, input.Self, input.Enemy, input.Obstacles) {
		action.Failed = true
		action.FailReason = validateReason(action, input.Self, input.Enemy, input.Obstacles)
	}
	return DecisionResult{Action: action, Note: note, Metadata: metadata}
}

// DecideTurnsStrictParallel uses the caller's cancellation for both requests.
func DecideTurnsStrictParallel(ctx context.Context, red, blue *Agent,
	redInput, blueInput AgentDecisionInput) (DecisionResult, DecisionResult) {
	var a, b DecisionResult
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); a = red.DecideTurnStrict(ctx, redInput) }()
	go func() { defer wg.Done(); b = blue.DecideTurnStrict(ctx, blueInput) }()
	wg.Wait()
	return a, b
}
