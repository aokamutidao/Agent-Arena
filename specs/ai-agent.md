# S5 - AI Agent 规格 (AI Agent Specification)

## Agent Arena — LLM-Powered Decision Engine

---

## 1. 架构总览

```
┌──────────────────────────────────────────────────────────────┐
│                   AI Agent 架构                              │
│                                                              │
│  每回合：                                                    │
│                                                              │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐              │
│  │ 构建      │    │ 调用      │    │ 解析      │              │
│  │ Prompt   │ →  │ Qwen API │ →  │ Action   │              │
│  └──────────┘    └──────────┘    └──────────┘              │
│       │                                  │                  │
│       │ 输入:                             │ 输出:             │
│       │ - 游戏状态                        │ - MOVE(x,y)      │
│       │ - 自身状态                        │ - ATTACK         │
│       │ - 对手状态                        │ - SKILL          │
│       │ - 性格描述（自然语言）            │ - CHARGE         │
│       │ - 策略投票权重                    │ - WAIT           │
│       │ - 可选动作列表                    │                  │
│       │ - 约束规则                                          │
│       └──────────────────────────────────┘                  │
│                                                              │
│  两个 Agent 并行调用 Qwen API，不串行                         │
└──────────────────────────────────────────────────────────────┘
```

**核心设计**：Agent 的"大脑"就是 LLM 本身。不需要手写规则引擎，性格通过自然语言 prompt 定义，策略投票通过 prompt 注入。

---

## 2. 四种预设 Agent

### 2.1 Berserker（狂战士）

```yaml
Name: Berserker
Personality: |
  你是一个狂战士 AI，信奉"最好的防守就是进攻"。
  你追求最短距离接近对手，造成最大伤害。
  你会使用蓄力攻击来打出爆发伤害。
  你几乎不会选择 WAIT，也不会远距离移动。
  你的目标是尽快击败对手，哪怕自己受伤也无所谓。
Detail: "激进攻型，喜欢近战冲锋，蓄力爆发"
Avatar: "/avatars/berserker.png"
```

### 2.2 Tactician（战术家）

```yaml
Name: Tactician
Personality: |
  你是一个战术家 AI，信奉"控制距离就是控制战局"。
  你善于利用技能远程消耗对手，同时保持安全距离。
  你会观察对手的模式，找到最佳进攻时机。
  你不会轻易冒险，但一旦出手就要有效。
  你很少主动靠近对手，更倾向于等待对手犯错。
Detail: "稳扎稳打，远程消耗，等待时机"
Avatar: "/avatars/tactician.png"
```

### 2.3 Trickster（诡术师）

```yaml
Name: Trickster
Personality: |
  你是一个诡术师 AI，信奉"欺骗是最好的武器"。
  你善于利用蓄力攻击打出双倍伤害。
  你会假装后退，然后突然蓄力反击。
  你喜欢利用障碍物绕后，打出出其不意的攻击。
  你经常使用 CHARGE，让对手猜不透你的下一步。
Detail: "善于欺骗，蓄力偷袭，利用地形"
Avatar: "/avatars/trickster.png"
```

### 2.4 Defender（守护者）

```yaml
Name: Defender
Personality: |
  你是一个守护者 AI，信奉"耐心是胜利的关键"。
  你善于防守，等待对手犯错。
  你会用技能远程骚扰，但绝不轻易靠近。
  你的目标是让对手先犯错，然后抓住机会反击。
  你几乎从不选择 CHARGE，因为那太冒险了。
Detail: "铁壁防守，等待反击，绝不冒险"
Avatar: "/avatars/defender.png"
```

---

## 3. Prompt 设计

### 3.1 决策 Prompt（每回合调用）

```go
func buildDecisionPrompt(a *Agent, state GameState, weights StrategyWeights) string {
    return fmt.Sprintf(`你是 Agent Arena 中的战斗 AI "%s"。

## 你的性格
%s

## 观众策略投票（影响你的打法倾向）
- 激进倾向: %d%% — 偏好 ATTACK、主动靠近
- 稳健倾向: %d%% — 偏好 SKILL 远程、保持距离
- 诡道倾向: %d%% — 偏好 CHARGE 蓄力、绕后偷袭
→ 请根据投票倾向调整你的打法。

## 当前局势
- 你的位置: (%d, %d)，HP: %d/100
- 对手位置: (%d, %d)，HP: %d/100
- 你与对手距离: %d 格（曼哈顿距离）
- 你的技能冷却: %d 回合后可用
- 当前回合: %d/30
- 你的状态: %s
- 对手的状态: %s

## 地图障碍物坐标
%s

## 可选动作（只能选一个）
- MOVE(x,y)  — 移动 1-2 格到目标坐标（曼哈顿距离 ≤ 2）
- ATTACK     — 近战攻击（需距离 ≤ 1），伤害 15，蓄力时 30
- SKILL      — 远程技能（需距离 ≤ 4 且冷却完毕），伤害 12，20%% 概率眩晕
- CHARGE     — 蓄力（下回合攻击伤害 ×2，但本回合不能移动）
- WAIT       — 原地等待

## 规则约束
- MOVE 目标不能是障碍物坐标，不能超出地图 (0-9, 0-9)
- ATTACK 只在距离 ≤ 1 时有效
- SKILL 只在距离 ≤ 4 且冷却 = 0 时有效
- 如果你正在蓄力（Charging），被攻击会打断蓄力

## 输出要求
只回复一个动作，严格按格式:
MOVE(x,y) 或 ATTACK 或 SKILL 或 CHARGE 或 WAIT
不要解释，不要多余文字，不要 markdown。`,

        a.Name,
        a.Personality,
        weights.Aggressive, weights.Defensive, weights.Tricky,
        state.Self.Position.X, state.Self.Position.Y, state.Self.HP,
        state.Enemy.Position.X, state.Enemy.Position.Y, state.Enemy.HP,
        manhattanDistance(state.Self.Position, state.Enemy.Position),
        state.Self.SkillCooldown,
        state.CurrentRound,
        formatStatus(state.Self.Status),
        formatStatus(state.Enemy.Status),
        formatObstacles(state.Obstacles),
    )
}
```

### 3.2 独白 Prompt（每回合异步调用，不阻塞）

```go
func buildMonologuePrompt(a *Agent, state GameState, action Action) string {
    return fmt.Sprintf(`你是 Agent Arena 中的战斗 AI "%s"。

刚才的回合你选择了 %s。
当前你的 HP: %d，对手 HP: %d。

用一句话表达你此刻的想法（20 字以内，符合你的性格）。
只回复这句话，不要其他内容。

你的性格: %s`,
        a.Name,
        formatAction(action),
        state.Self.HP, state.Enemy.HP,
        a.Personality,
    )
}
```

---

## 4. Qwen API 客户端

### 4.1 接口定义

```go
type LLMClient interface {
    Chat(ctx context.Context, prompt string) (string, error)
}
```

### 4.2 Qwen 实现

```go
type QwenClient struct {
    apiKey   string
    model    string        // "qwen-turbo" (快) 或 "qwen-plus" (强)
    endpoint string
    timeout  time.Duration
}

func NewQwenClient(apiKey string) *QwenClient {
    return &QwenClient{
        apiKey:   apiKey,
        model:    "qwen-turbo",   // MVP 用 turbo（快+便宜）
        endpoint: "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
        timeout:  5 * time.Second,
    }
}

func (c *QwenClient) Chat(ctx context.Context, prompt string) (string, error) {
    ctx, cancel := context.WithTimeout(ctx, c.timeout)
    defer cancel()

    body := map[string]interface{}{
        "model": c.model,
        "messages": []map[string]string{
            {"role": "system", "content": "你是一个战斗 AI，只回复动作指令。"},
            {"role": "user", "content": prompt},
        },
        "max_tokens": 50,        // 决策只需要很短的回复
        "temperature": 0.8,      // 保持一定随机性
    }

    // HTTP POST 请求
    reqBody, _ := json.Marshal(body)
    req, _ := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+c.apiKey)

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // 解析响应
    var result struct {
        Choices []struct {
            Message struct {
                Content string `json:"content"`
            } `json:"message"`
        } `json:"choices"`
    }
    json.NewDecoder(resp.Body).Decode(&result)

    return strings.TrimSpace(result.Choices[0].Message.Content), nil
}
```

---

## 5. 动作解析器

```go
func ParseAction(response string) (Action, error) {
    // 清理响应（去掉可能的 markdown、引号等）
    response = cleanResponse(response)

    // 正则匹配
    // MOVE(x,y) | ATTACK | SKILL | CHARGE | WAIT
    re := regexp.MustCompile(`^(MOVE\((\d+),(\d+)\)|ATTACK|SKILL|CHARGE|WAIT)$`)
    matches := re.FindStringSubmatch(response)

    if matches == nil {
        // 解析失败 → 默认 WAIT（安全回退）
        return Action{Type: ActionWait}, fmt.Errorf("invalid action: %s", response)
    }

    action := Action{Type: ActionType(matches[1])}

    // 解析 MOVE 的目标坐标
    if matches[1] == "MOVE" && matches[2] != "" {
        x, _ := strconv.Atoi(matches[2])
        y, _ := strconv.Atoi(matches[3])
        // 验证坐标范围
        if x < 0 || x > 9 || y < 0 || y > 9 {
            return Action{Type: ActionWait}, fmt.Errorf("invalid coordinates: %d,%d", x, y)
        }
        action.Target = Position{X: uint8(x), Y: uint8(y)}
    }

    return action, nil
}

func cleanResponse(s string) string {
    s = strings.TrimSpace(s)
    s = strings.Trim(s, "\"'`")
    s = strings.TrimPrefix(s, "```")
    s = strings.TrimSuffix(s, "```")
    return strings.TrimSpace(s)
}
```

### 5.1 合法性验证

```go
// 验证 LLM 返回的动作是否合法
func ValidateAction(action Action, self AgentState, enemy AgentState, obstacles []Position) bool {
    switch action.Type {
    case ActionMove:
        // 检查距离 ≤ 2
        dist := manhattanDistance(self.Position, action.Target)
        if dist > 2 || dist == 0 {
            return false
        }
        // 检查不是障碍物
        if isObstacle(action.Target, obstacles) {
            return false
        }
        // 检查不是对手位置
        if action.Target == enemy.Position {
            return false
        }
        return true

    case ActionAttack:
        return manhattanDistance(self.Position, enemy.Position) <= 1

    case ActionSkill:
        return manhattanDistance(self.Position, enemy.Position) <= 4 && self.SkillCooldown == 0

    case ActionCharge:
        return true  // 总是合法

    case ActionWait:
        return true  // 总是合法

    default:
        return false
    }
}

// 如果不合法，回退到 WAIT
func SafeAction(action Action, err error, self AgentState, enemy AgentState, obstacles []Position) Action {
    if err != nil || !ValidateAction(action, self, enemy, obstacles) {
        return Action{Type: ActionWait}
    }
    return action
}
```

---

## 6. 完整决策流程

```go
type Agent struct {
    Name         string
    Personality  string
    LLM          LLMClient
}

func (a *Agent) DecideTurn(state GameState, weights StrategyWeights) Action {
    // 1. 构建 prompt
    prompt := buildDecisionPrompt(a, state, weights)

    // 2. 调用 Qwen API
    response, err := a.LLM.Chat(context.Background(), prompt)
    if err != nil {
        // API 失败 → 默认 WAIT
        return Action{Type: ActionWait}
    }

    // 3. 解析动作
    action, parseErr := ParseAction(response)

    // 4. 验证合法性（LLM 可能返回不合法的动作）
    return SafeAction(action, parseErr, state.Self, state.Enemy, state.Obstacles)
}
```

### 6.1 两 Agent 并行决策

```go
// GameService 中每回合的决策逻辑
func (s *GameService) decideActions(game *Game) (Action, Action) {
    var redAction, blueAction Action
    var wg sync.WaitGroup

    wg.Add(2)

    // Red Agent 决策（并行）
    go func() {
        defer wg.Done()
        redAction = game.AgentRed.DecideTurn(game.State, game.StrategyRed)
    }()

    // Blue Agent 决策（并行）
    go func() {
        defer wg.Done()
        blueAction = game.AgentBlue.DecideTurn(game.State, game.StrategyBlue)
    }()

    wg.Wait()
    return redAction, blueAction
}
```

---

## 7. 独白生成（异步，不阻塞）

```go
func (a *Agent) GenerateMonologue(state GameState, action Action) <-chan string {
    ch := make(chan string, 1)

    go func() {
        prompt := buildMonologuePrompt(a, state, action)
        response, err := a.LLM.Chat(context.Background(), prompt)
        if err != nil {
            ch <- "..."
            return
        }
        ch <- response
    }()

    return ch
}

// 在 GameService 中使用
func (s *GameService) generateMonologues(game *Game, redAction, blueAction Action) {
    redCh := game.AgentRed.GenerateMonologue(game.State, redAction)
    blueCh := game.AgentBlue.GenerateMonologue(game.State, blueAction)

    // 不等待，通过 WebSocket 异步推送
    go func() {
        if text := <-redCh; text != "" {
            s.wsHub.Broadcast(game.ID, MonologueEvent{
                Side: "red", Agent: game.AgentRed.Name, Text: text,
            })
        }
    }()
    go func() {
        if text := <-blueCh; text != "" {
            s.wsHub.Broadcast(game.ID, MonologueEvent{
                Side: "blue", Agent: game.AgentBlue.Name, Text: text,
            })
        }
    }()
}
```

---

## 8. 性能与成本

### 8.1 延迟控制

```
决策调用:  qwen-turbo ~300-500ms
独白调用:  qwen-turbo ~300-500ms (异步，不阻塞)

每回合时序:
  t=0.0s   并行调用 Red + Blue 决策
  t=0.5s   两边返回，执行动作
  t=0.5s   异步启动独白生成
  t=1.0s   WebSocket 推送回合状态
  t=1.0s   独白可能返回，推送独白
  t=5.0s   等待结束，下一回合开始

→ 每回合间隔 5 秒，AI 决策只占 0.5 秒，完全够用
```

### 8.2 成本控制

```
Qwen Turbo 定价:
  输入: 0.002 元 / 1000 tokens
  输出: 0.006 元 / 1000 tokens

每回合:
  决策 prompt: ~350 tokens 输入 + ~10 tokens 输出
  独白 prompt: ~150 tokens 输入 + ~20 tokens 输出
  每 Agent 每回合: ~0.0015 元
  每回合（2 Agent）: ~0.003 元

每场对局（30 回合）: ~0.09 元
每天测试 100 场: ~9 元

→ 成本可忽略
```

### 8.3 容错

```
LLM 返回不合法动作:
  → ValidateAction 检测到 → 回退到 WAIT

LLM API 超时/失败:
  → 5 秒超时 → 回退到 WAIT

LLM 返回无法解析的内容:
  → ParseAction 失败 → 回退到 WAIT

所有容错都是静默回退，不会中断对局。
```

---

## 9. 代码结构

```
backend/internal/agent/
├── agent.go           # Agent 结构体 + DecideTurn 主逻辑 (~50 行)
├── prompt.go          # buildDecisionPrompt + buildMonologuePrompt (~80 行)
├── parser.go          # ParseAction + ValidateAction + SafeAction (~60 行)
├── qwen_client.go     # QwenClient 实现 (~60 行)
├── personalities.go   # 4 种性格定义 (~30 行)
└── agent_test.go      # 单元测试 (~80 行)

总计: ~360 行 Go 代码
```

---

## 10. 测试

### 10.1 单元测试

```go
// 测试动作解析
func TestParseAction(t *testing.T) {
    tests := []struct {
        input    string
        expected Action
    }{
        {"ATTACK", Action{Type: ActionAttack}},
        {"MOVE(5,3)", Action{Type: ActionMove, Target: Position{5, 3}}},
        {"SKILL", Action{Type: ActionSkill}},
        {"CHARGE", Action{Type: ActionCharge}},
        {"WAIT", Action{Type: ActionWait}},
        {"```ATTACK```", Action{Type: ActionAttack}},   // 清理 markdown
        {"invalid", Action{Type: ActionWait}},           // 无效 → WAIT
    }
    // ...
}

// 测试动作验证
func TestValidateAction(t *testing.T) {
    self := AgentState{Position: Position{3, 3}, SkillCooldown: 0}
    enemy := AgentState{Position: Position{4, 3}}

    // ATTACK 距离 1 → 合法
    assert.True(t, ValidateAction(Action{Type: ActionAttack}, self, enemy, nil))

    // ATTACK 距离 3 → 不合法
    farEnemy := AgentState{Position: Position{6, 3}}
    assert.False(t, ValidateAction(Action{Type: ActionAttack}, self, farEnemy, nil))
}
```

### 10.2 集成测试

```go
// 测试完整对局（用 mock LLM）
func TestFullGameWithMockLLM(t *testing.T) {
    mockLLM := &MockLLMClient{
        responses: []string{"ATTACK", "MOVE(5,5)", "SKILL", "CHARGE", "WAIT"},
    }

    red := NewAgent("Berserker", berserkerPersonality, mockLLM)
    blue := NewAgent("Tactician", tacticianPersonality, mockLLM)

    game := NewGame(red, blue)
    game.Run()

    assert.Equal(t, GameFinished, game.Status)
    assert.NotEqual(t, SideNone, game.Winner)
}
```

---

**文档版本**：v2.0 (LLM-Powered)
**最后更新**：2026-08-02
**状态**：待评审
