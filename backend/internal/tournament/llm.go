package tournament

import (
	"agent-arena/backend/internal/agent"
)

func summarizeCalls(calls []DecisionCall) *LLMSummary {
	if len(calls) == 0 {
		return nil
	}
	summary := &LLMSummary{Attempted: len(calls)}
	var prompt, completion, total uint64
	promptKnown, completionKnown, totalKnown := true, true, true
	for _, call := range calls {
		summary.Provider, summary.Model = call.Provider, call.Model
		summary.LatencyMS += call.LatencyMS
		if call.Success {
			summary.Succeeded++
		} else if call.ErrorCode != "in_progress" {
			summary.Failed++
			summary.LastErrorCode = call.ErrorCode
		}
		if call.Usage == nil || call.Usage.PromptTokens == nil {
			promptKnown = false
		} else {
			prompt += *call.Usage.PromptTokens
		}
		if call.Usage == nil || call.Usage.CompletionTokens == nil {
			completionKnown = false
		} else {
			completion += *call.Usage.CompletionTokens
		}
		if call.Usage == nil || call.Usage.TotalTokens == nil {
			totalKnown = false
		} else {
			total += *call.Usage.TotalTokens
		}
	}
	if promptKnown {
		summary.Usage.PromptTokens = &prompt
	}
	if completionKnown {
		summary.Usage.CompletionTokens = &completion
	}
	if totalKnown {
		summary.Usage.TotalTokens = &total
	}
	return summary
}

func clientInfo(client agent.LLMClient, cfg Config) agent.LLMInfo {
	if described, ok := client.(agent.LLMDescriber); ok {
		return described.LLMInfo()
	}
	return agent.LLMInfo{
		Provider: cfg.LLMProvider, Model: cfg.LLMModel,
		TimeoutSeconds: cfg.LLMTimeoutSeconds,
	}
}

func (r *Runner) PublicMatch(m *Match) *Match {
	r.mu.Lock()
	busy := r.jobs[m.TournamentID]
	r.mu.Unlock()
	m.RetryAllowed = m.Status == "interrupted" && !busy
	if m.Mode == "llm" {
		m.RetryAllowed = m.RetryAllowed && r.llm != nil && len(m.LLMCalls)+2 <= r.cfg.LLMMaxCalls
	} else {
		m.RetryAllowed = m.RetryAllowed && r.cfg.SimulationEnabled
	}
	return m
}
