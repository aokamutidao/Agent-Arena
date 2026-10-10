package agent

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"agent-arena/backend/internal/engine"
)

func providerFixture(t *testing.T, handler http.HandlerFunc, timeout time.Duration) *QwenClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	client, err := NewConfiguredQwenClient(QwenConfig{
		APIKey: "test-credential-do-not-log", Provider: "fixture", Model: "fixture-model",
		Endpoint: server.URL, Timeout: timeout,
	})
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestConfiguredProviderSuccessUsageAndHEALProtocol(t *testing.T) {
	client := providerFixture(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.Header.Get("Authorization") != "Bearer test-credential-do-not-log" {
			t.Error("configured method/authentication missing")
		}
		var request struct {
			Model          string  `json:"model"`
			MaxTokens      int     `json:"max_tokens"`
			Temperature    float64 `json:"temperature"`
			EnableThinking *bool   `json:"enable_thinking"`
			Messages       []struct {
				Content string `json:"content"`
			} `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Error(err)
		}
		if request.Model != "qwen3.5-plus" || request.MaxTokens != 500 || request.Temperature != 0.7 ||
			request.EnableThinking == nil || *request.EnableThinking || len(request.Messages) != 2 ||
			!strings.Contains(request.Messages[0].Content, "HEAL") ||
			!strings.Contains(request.Messages[0].Content, "NOTE:") {
			t.Error("model/decision protocol mismatch")
		}
		w.Header().Set("X-Request-Id", "request-123")
		_, _ = w.Write([]byte(`{"model":"fixture-model-v2","choices":[{"message":{"content":"HEAL\nNOTE: Restore health."}}],"usage":{"prompt_tokens":32,"completion_tokens":6,"total_tokens":38}}`))
	}, time.Second)
	client.model = "qwen3.5-plus"
	result := NewAgent("Guard", "Defend", client).DecideTurnStrict(context.Background(),
		AgentDecisionInput{Self: engine.AgentState{HP: 50, MaxHP: 100}, Enemy: engine.AgentState{}})
	if result.Err != nil || result.Action.Type != engine.ActionHeal || result.Note != "Restore health." {
		t.Fatalf("valid short HEAL decision rejected: %+v", result)
	}
	metadata := result.Metadata
	if !metadata.Success || metadata.Provider != "fixture" || metadata.Model != "fixture-model-v2" ||
		metadata.RequestID != "request-123" || metadata.HTTPStatus != 200 || metadata.Usage == nil ||
		metadata.Usage.TotalTokens == nil || *metadata.Usage.TotalTokens != 38 {
		t.Fatalf("missing provider facts: %+v", metadata)
	}
}

func TestProviderErrorsDiscardBodiesAndUnsafeIdentifiers(t *testing.T) {
	for _, test := range []struct {
		status int
		code   string
	}{{401, "authentication_error"}, {429, "rate_limited"}, {500, "provider_error"}} {
		t.Run(test.code, func(t *testing.T) {
			client := providerFixture(t, func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("X-Request-Id", "test-credential-do-not-log")
				w.WriteHeader(test.status)
				_, _ = w.Write([]byte("private body test-credential-do-not-log"))
			}, time.Second)
			response, err := client.ChatWithMetadata(context.Background(), "test")
			encoded, _ := json.Marshal(response.Metadata)
			if err == nil || ErrorCode(err) != test.code || response.Metadata.Success ||
				response.Metadata.RequestID != "" || strings.Contains(err.Error(), "private body") ||
				strings.Contains(string(encoded), "test-credential-do-not-log") {
				t.Fatalf("unsafe or misclassified provider failure: %+v %v", response.Metadata, err)
			}
		})
	}
}

func TestProviderCompletionIDFallbackAndNoRedirect(t *testing.T) {
	client := providerFixture(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"id":"completion-456","choices":[{"message":{"content":"WAIT"}}]}`))
	}, time.Second)
	result, err := client.ChatWithMetadata(context.Background(), "test")
	if err != nil || result.Metadata.RequestID != "completion-456" || result.Metadata.Usage != nil {
		t.Fatal("completion id fallback or unknown usage mishandled")
	}
	var forwarded atomic.Bool
	target := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		forwarded.Store(true)
	}))
	defer target.Close()
	redirect := providerFixture(t, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusFound)
	}, time.Second)
	_, err = redirect.ChatWithMetadata(context.Background(), "test")
	if err == nil || forwarded.Load() {
		t.Fatal("provider redirect forwarded an authenticated request")
	}
	misconfigured := NewQwenClient("sk-synthetic-secret")
	misconfigured.model, misconfigured.provider = "sk-synthetic-secret", "sk-synthetic-secret"
	info := misconfigured.LLMInfo()
	if info.Model != "" || info.Provider != "" {
		t.Fatal("safe configuration exposed a credential-shaped identifier")
	}
}

func TestProviderTimeoutCancellationAndInvalidResponse(t *testing.T) {
	client := providerFixture(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		select {
		case <-r.Context().Done():
		case <-time.After(100 * time.Millisecond):
		}
	}, 30*time.Millisecond)
	response, err := client.ChatWithMetadata(context.Background(), "test")
	if err == nil || ErrorCode(err) != "timeout" || response.Metadata.LatencyMS < 20 {
		t.Fatalf("request timeout was not recorded: %+v %v", response.Metadata, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = client.ChatWithMetadata(ctx, "test")
	if err == nil || ErrorCode(err) != "cancelled" {
		t.Fatalf("cancellation not propagated: %v", err)
	}
	invalid := providerFixture(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"choices":[]}`))
	}, time.Second)
	response, err = invalid.ChatWithMetadata(context.Background(), "test")
	if err == nil || ErrorCode(err) != "invalid_response" || response.Metadata.Success {
		t.Fatal("empty completion was accepted")
	}
}

func TestStrictParserRejectsAmbiguityAndInternalReasoning(t *testing.T) {
	for _, reply := range []string{
		"WAIT or ATTACK", "analysis: ATTACK", "MOVE(99,99)",
		"ATTACK\n<think>internal reasoning</think>", "ATTACK\nNOTE: short\nmore",
	} {
		if _, _, err := ParseStrictDecision(reply); err == nil {
			t.Fatalf("ambiguous protocol accepted: %q", reply)
		}
	}
	for _, reply := range []string{"WAIT", "HEAL\nNOTE: Recover.", "MOVE(2,3)\n说明：接近对手。"} {
		if _, _, err := ParseStrictDecision(reply); err != nil {
			t.Fatalf("valid final decision rejected: %q", reply)
		}
	}
	client := providerFixture(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"analysis: ATTACK"}}]}`))
	}, time.Second)
	decision := NewAgent("Agent", "Profile", client).DecideTurnStrict(context.Background(), AgentDecisionInput{})
	if decision.Err == nil || decision.Metadata.Success || decision.Metadata.ErrorCode != "parse_error" {
		t.Fatal("strict decision silently used WAIT for a parser failure")
	}
}

func TestCourseConfigIsExplicitAndGeneralAPIDefault(t *testing.T) {
	values := map[string]string{}
	getenv := func(key string) string { return values[key] }
	cfg, err := CourseLLMConfig(getenv)
	if err != nil || cfg.APIKey != "" || cfg.Model != "qwen-plus" ||
		cfg.Endpoint != "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions" {
		t.Fatalf("unsafe course defaults: provider=%s model=%s err=%v", cfg.Provider, cfg.Model, err)
	}
	values["QWEN_API_KEY"] = "legacy-fixture"
	values["LLM_API_KEY"] = "explicit-fixture"
	values["LLM_ENDPOINT"] = "https://coding.dashscope.aliyuncs.com/v1/chat/completions"
	values["LLM_MODEL"] = "qwen3-coder-plus"
	cfg, err = CourseLLMConfig(getenv)
	if err != nil || cfg.APIKey != "explicit-fixture" || cfg.Model != "qwen3-coder-plus" {
		t.Fatal("explicit operator configuration/alias precedence failed")
	}
	values["LLM_ENDPOINT"] = "https://user:secret@example.com/api?key=secret"
	if _, err := CourseLLMConfig(getenv); err == nil {
		t.Fatal("credential-bearing endpoint URL accepted")
	}
	legacy := NewQwenClient("synthetic")
	if legacy.model != "qwen3-coder-plus" || legacy.timeout != 10*time.Second {
		t.Fatal("legacy defaults unexpectedly changed")
	}
}
