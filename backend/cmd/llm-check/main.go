// llm-check prints safe configuration and optionally makes exactly one request.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"net/url"
	"os"
	"time"

	"agent-arena/backend/internal/agent"
)

func main() {
	call := flag.Bool("call", false, "Explicitly send one potentially billable provider request")
	flag.Parse()
	cfg, err := agent.CourseLLMConfig(os.Getenv)
	if err != nil {
		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
			"status": "invalid_configuration", "error": err.Error(),
		})
		os.Exit(1)
	}
	endpoint, _ := url.Parse(cfg.Endpoint)
	result := map[string]any{
		"configured": cfg.APIKey != "", "provider": cfg.Provider, "model": cfg.Model,
		"timeout_seconds": int64(cfg.Timeout / time.Second), "endpoint_host": endpoint.Hostname(),
		"request_sent": false,
	}
	if !*call {
		result["status"] = "configuration_only"
		_ = json.NewEncoder(os.Stdout).Encode(result)
		return
	}
	client, err := agent.NewConfiguredQwenClient(cfg)
	if err != nil || cfg.APIKey == "" {
		result["status"] = "not_configured"
		_ = json.NewEncoder(os.Stdout).Encode(result)
		os.Exit(1)
	}
	result["request_sent"] = true
	response, err := client.ChatWithMetadata(context.Background(),
		"Connectivity test: return exactly WAIT on the first line and NOTE: Ready on the second line.")
	if err == nil {
		action, note, parseErr := agent.ParseStrictDecision(response.Content)
		if parseErr != nil {
			err = parseErr
			response.Metadata.Success = false
			response.Metadata.ErrorCode = "parse_error"
		} else {
			result["action"], result["note"] = action, note
		}
	}
	result["call"] = response.Metadata
	if err != nil {
		result["status"] = "failed"
		result["error_code"] = agent.ErrorCode(err)
	} else {
		result["status"] = "succeeded"
	}
	_ = json.NewEncoder(os.Stdout).Encode(result)
	if err != nil {
		os.Exit(1)
	}
}
