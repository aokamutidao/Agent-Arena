package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

const decisionProtocol = `You are a battle agent. Return one action on the first line:
MOVE(x,y), ATTACK, SKILL, HEAL, CHARGE or WAIT.
Optionally add a second line: NOTE: <a brief description of the action's purpose,
at most 80 characters>. Do not disclose chain-of-thought or internal reasoning.
Do not include markdown, analysis, other commands, or extra lines.`

// NewConfiguredQwenClient uses only explicitly supplied server configuration.
func NewConfiguredQwenClient(cfg QwenConfig) (*QwenClient, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &QwenClient{
		apiKey: cfg.APIKey, provider: cfg.Provider, model: cfg.Model,
		endpoint: cfg.Endpoint, timeout: cfg.Timeout,
	}, nil
}

// LLMInfo never exposes the API key or endpoint.
func (c *QwenClient) LLMInfo() LLMInfo {
	return LLMInfo{
		Provider: c.safeIdentifier(c.provider), Model: c.safeIdentifier(c.model),
		TimeoutSeconds: int64(c.timeout / time.Second),
	}
}

// Do not forward credentials to redirect targets, including HTTPS downgrades.
var providerHTTPClient = &http.Client{
	CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func (c *QwenClient) safeIdentifier(value string) string {
	if !identifierPattern.MatchString(value) || strings.HasPrefix(strings.ToLower(value), "sk-") {
		return ""
	}
	if c.apiKey != "" && strings.Contains(value, c.apiKey) {
		return ""
	}
	return value
}

// ChatWithMetadata reads bounded responses and discards provider error bodies.
func (c *QwenClient) ChatWithMetadata(ctx context.Context, prompt string) (result LLMResponse, err error) {
	started := time.Now()
	result.Metadata = CallMetadata{Provider: c.safeIdentifier(c.provider), Model: c.safeIdentifier(c.model)}
	defer func() {
		result.Metadata.LatencyMS = time.Since(started).Milliseconds()
		result.Metadata.Success = err == nil
		if err != nil {
			result.Metadata.ErrorCode = ErrorCode(err)
		}
	}()
	if c.apiKey == "" {
		return result, &CallError{Code: "not_configured"}
	}
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	body := map[string]any{
		"model": c.model, "max_tokens": 500, "temperature": 0.7,
		"messages": []map[string]string{
			{"role": "system", "content": decisionProtocol},
			{"role": "user", "content": prompt},
		},
	}
	// Qwen 3 thinking models must return only the final action protocol. Keep
	// this vendor extension out of requests to other compatible model families.
	if strings.HasPrefix(strings.ToLower(c.model), "qwen3") {
		body["enable_thinking"] = false
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		return result, &CallError{Code: "invalid_response"}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(encoded))
	if err != nil {
		return result, &CallError{Code: "provider_error"}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	response, err := providerHTTPClient.Do(req)
	if err != nil {
		return result, &CallError{Code: ErrorCode(err)}
	}
	defer response.Body.Close()
	result.Metadata.HTTPStatus = response.StatusCode
	for _, header := range []string{"X-Request-Id", "X-Dashscope-Request-Id"} {
		if identifier := c.safeIdentifier(response.Header.Get(header)); identifier != "" {
			result.Metadata.RequestID = identifier
			break
		}
	}
	if response.StatusCode != http.StatusOK {
		code := "provider_error"
		switch response.StatusCode {
		case http.StatusUnauthorized, http.StatusForbidden:
			code = "authentication_error"
		case http.StatusTooManyRequests:
			code = "rate_limited"
		}
		return result, &CallError{Code: code}
	}
	const bodyLimit = 1 << 20
	data, err := io.ReadAll(io.LimitReader(response.Body, bodyLimit+1))
	if err != nil {
		return result, &CallError{Code: ErrorCode(err)}
	}
	if len(data) > bodyLimit {
		return result, &CallError{Code: "response_too_large"}
	}
	var decoded struct {
		ID        string      `json:"id"`
		Model     string      `json:"model"`
		RequestID string      `json:"request_id"`
		Usage     *TokenUsage `json:"usage"`
		Choices   []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if json.Unmarshal(data, &decoded) != nil || len(decoded.Choices) == 0 {
		return result, &CallError{Code: "invalid_response"}
	}
	if model := c.safeIdentifier(decoded.Model); model != "" {
		result.Metadata.Model = model
	}
	if result.Metadata.RequestID == "" {
		result.Metadata.RequestID = c.safeIdentifier(decoded.RequestID)
	}
	if result.Metadata.RequestID == "" {
		result.Metadata.RequestID = c.safeIdentifier(decoded.ID)
	}
	result.Metadata.Usage = decoded.Usage
	result.Content = strings.TrimSpace(decoded.Choices[0].Message.Content)
	if result.Content == "" || strings.Contains(result.Content, c.apiKey) {
		result.Content = ""
		return result, &CallError{Code: "invalid_response"}
	}
	return result, nil
}
