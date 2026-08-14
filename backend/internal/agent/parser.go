package agent

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"agent-arena/backend/internal/engine"
)

// (?i) 忽略大小写：LLM 常输出 "HEAl"（小写 L）/ "attack" 等变体
var actionRegex = regexp.MustCompile(`(?i)(MOVE\((\d+),(\d+)\)|ATTACK|SKILL|CHARGE|WAIT|HEAL)`)

// ParseAction 解析 LLM 返回的动作文本
func ParseAction(response string) (engine.Action, error) {
	action, _, err := ParseActionWithReasoning(response)
	return action, err
}

// ParseActionWithReasoning 解析动作 + 提取推理思路
func ParseActionWithReasoning(response string) (engine.Action, string, error) {
	cleaned := cleanResponse(response)

	// 提取 "思考:" 后面的推理文本
	reasoning := ""
	reasonRegex := regexp.MustCompile(`(?m)思考[:：]\s*(.+)`)
	if m := reasonRegex.FindStringSubmatch(cleaned); m != nil {
		reasoning = strings.TrimSpace(m[1])
		// 截断过长的推理
		if len([]rune(reasoning)) > 60 {
			runes := []rune(reasoning)
			reasoning = string(runes[:60]) + "..."
		}
	}

	// 先尝试精确匹配（整个字符串就是一个动作）
	exactRegex := regexp.MustCompile(`(?i)^(MOVE\((\d+),(\d+)\)|ATTACK|SKILL|CHARGE|WAIT|HEAL)$`)
	if matches := exactRegex.FindStringSubmatch(cleaned); matches != nil {
		// 统一大写（LLM 可能返回小写）
		matches[1] = strings.ToUpper(matches[1])
		action, err := parseActionMatches(matches)
		return action, reasoning, err
	}

	// 宽松匹配：从文本中提取第一个动作
	matches := actionRegex.FindStringSubmatch(cleaned)
	if matches == nil {
		return engine.Action{Type: engine.ActionWait}, reasoning, fmt.Errorf("invalid action: %s", cleaned)
	}
	// 统一大写
	matches[1] = strings.ToUpper(matches[1])

	action, err := parseActionMatches(matches)
	return action, reasoning, err
}

func parseActionMatches(matches []string) (engine.Action, error) {
	// MOVE(x,y) 格式
	if strings.HasPrefix(matches[1], "MOVE") {
		x, err := strconv.Atoi(matches[2])
		if err != nil {
			return engine.Action{Type: engine.ActionWait}, fmt.Errorf("invalid x coord: %w", err)
		}
		y, err := strconv.Atoi(matches[3])
		if err != nil {
			return engine.Action{Type: engine.ActionWait}, fmt.Errorf("invalid y coord: %w", err)
		}
		if x < 0 || x > 9 || y < 0 || y > 9 {
			return engine.Action{Type: engine.ActionWait}, fmt.Errorf("coordinates out of range: %d,%d", x, y)
		}
		return engine.Action{Type: engine.ActionMove, Target: engine.Position{X: uint8(x), Y: uint8(y)}}, nil
	}

	// 其他动作
	return engine.Action{Type: engine.ActionType(matches[1])}, nil
}

// cleanResponse 清理 LLM 返回的文本
func cleanResponse(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "\"'`")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	s = strings.TrimPrefix(s, "markdown")

	// 规范化中文动作名称 → 英文
	s = strings.ReplaceAll(s, "自我治疗", "HEAL")
	s = strings.ReplaceAll(s, "治疗", "HEAL")
	s = strings.ReplaceAll(s, "蓄力", "CHARGE")
	s = strings.ReplaceAll(s, "攻击", "ATTACK")
	s = strings.ReplaceAll(s, "技能", "SKILL")
	s = strings.ReplaceAll(s, "移动", "MOVE")
	s = strings.ReplaceAll(s, "等待", "WAIT")

	return strings.TrimSpace(s)
}
