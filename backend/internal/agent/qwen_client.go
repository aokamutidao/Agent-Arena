package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// LLMClient LLM 调用接口
type LLMClient interface {
	Chat(ctx context.Context, prompt string) (string, error)
}

// QwenClient Qwen API 客户端
type QwenClient struct {
	apiKey   string
	model    string
	endpoint string
	timeout  time.Duration
}

// NewQwenClient 创建 Qwen 客户端
func NewQwenClient(apiKey string) *QwenClient {
	return &QwenClient{
		apiKey:   apiKey,
		model:    "qwen3-coder-plus",
		endpoint: "https://coding.dashscope.aliyuncs.com/v1/chat/completions",
		timeout:  10 * time.Second,
	}
}

// Chat 调用 Qwen API
func (c *QwenClient) Chat(ctx context.Context, prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	body := map[string]interface{}{
		"model": c.model,
		"messages": []map[string]string{
			{"role": "system", "content": "你是一个战斗 AI。只回复一个动作指令，不要解释，不要多余文字。格式：MOVE(x,y) 或 ATTACK 或 SKILL 或 CHARGE 或 WAIT"},
			{"role": "user", "content": prompt},
		},
		"max_tokens":  500,
		"temperature": 0.7,
	}

	reqBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errBody []byte
		errBody, _ = io.ReadAll(resp.Body)
		return "", fmt.Errorf("qwen api returned status %d: %s", resp.StatusCode, string(errBody))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("qwen api returned no choices")
	}

	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}

// ExternalAPIClient 外部 API 客户端（用于自定义 Agent）
type ExternalAPIClient struct {
	endpoint string
	apiKey   string
	model    string
	timeout  time.Duration
}

// NewExternalAPIClient 创建外部 API 客户端
func NewExternalAPIClient(endpoint, apiKey, model string) *ExternalAPIClient {
	if model == "" {
		model = "gpt-3.5-turbo"
	}
	return &ExternalAPIClient{
		endpoint: endpoint,
		apiKey:   apiKey,
		model:    model,
		timeout:  10 * time.Second,
	}
}

// Chat 调用外部 API（自动检测 OpenAI 或 Anthropic 格式）
func (c *ExternalAPIClient) Chat(ctx context.Context, prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// 根据 URL 判断使用哪种格式
	if c.isAnthropicFormat() {
		return c.chatAnthropic(ctx, prompt)
	}
	return c.chatOpenAI(ctx, prompt)
}

// isAnthropicFormat 判断是否为 Anthropic Messages API 格式
func (c *ExternalAPIClient) isAnthropicFormat() bool {
	return strings.Contains(strings.ToLower(c.endpoint), "/anthropic") ||
		strings.Contains(strings.ToLower(c.endpoint), "/v1/messages")
}

// normalizeAnthropicURL 对 Anthropic 兼容 endpoint 补全路径
// 用户常见输入：https://coding.dashscope.aliyuncs.com/apps/anthropic
// 正确 URL: .../anthropic/v1/messages
func normalizeAnthropicURL(endpoint string) string {
	e := strings.TrimRight(endpoint, "/")
	low := strings.ToLower(e)
	if strings.HasSuffix(low, "/v1/messages") {
		return e
	}
	// 以 /anthropic 结尾（或路径中含 /anthropic 但后面没有路径）→ 追加 /v1/messages
	if strings.HasSuffix(low, "/anthropic") {
		return e + "/v1/messages"
	}
	return e
}

// chatOpenAI 调用 OpenAI 兼容 API
func (c *ExternalAPIClient) chatOpenAI(ctx context.Context, prompt string) (string, error) {
	body := map[string]interface{}{
		"model": c.model,
		"messages": []map[string]string{
			{"role": "system", "content": "你是一个战斗 AI。只回复一个动作指令，不要解释，不要多余文字。格式：MOVE(x,y) 或 ATTACK 或 SKILL 或 CHARGE 或 WAIT"},
			{"role": "user", "content": prompt},
		},
		"max_tokens":  500,
		"temperature": 0.7,
	}

	reqBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errBody []byte
		errBody, _ = io.ReadAll(resp.Body)
		return "", fmt.Errorf("external api returned status %d: %s", resp.StatusCode, string(errBody))
	}

	// OpenAI 格式: choices[0].message.content
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("external api returned no choices")
	}

	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}

// chatAnthropic 调用 Anthropic Messages API（含 Qwen Anthropic-compatible proxy）
// 请求: POST /v1/messages  {model, max_tokens, system, messages:[{role,content}]}
// 响应: {content:[{type:"text", text:"..."}]}
func (c *ExternalAPIClient) chatAnthropic(ctx context.Context, prompt string) (string, error) {
	body := map[string]interface{}{
		"max_tokens": 500,
		"system":     "你是一个战斗 AI。只回复一个动作指令，不要解释，不要多余文字。格式：MOVE(x,y) 或 ATTACK 或 SKILL 或 CHARGE 或 WAIT",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	// Qwen coding plan endpoint 需要 model 参数
	if strings.Contains(c.endpoint, "dashscope") {
		body["model"] = "qwen3-coder-plus"
	}

	// 自动补全 Anthropic URL（用户常见漏写 /v1/messages）
	endpoint := normalizeAnthropicURL(c.endpoint)

	reqBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	// Anthropic 使用 x-api-key，OpenAI 使用 Authorization: Bearer
	// Qwen Anthropic proxy 同时支持两种，这里两个都发
	if c.apiKey != "" {
		req.Header.Set("x-api-key", c.apiKey)
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("anthropic-version", "2023-06-01")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("anthropic api returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Anthropic 格式: content[0].text
	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Error *struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("decode response: %w, body: %s", err, string(respBody))
	}

	if result.Error != nil {
		return "", fmt.Errorf("anthropic api error: %s - %s", result.Error.Type, result.Error.Message)
	}

	if len(result.Content) == 0 {
		return "", fmt.Errorf("anthropic api returned no content, body: %s", string(respBody))
	}

	return strings.TrimSpace(result.Content[0].Text), nil
}
