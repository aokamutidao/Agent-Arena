# S2 - 游戏规则规格 (Game Rules Specification)

## Agent Arena 格子竞技场

---

## 1. 地图设计

### 1.1 基础参数

| 参数 | 值 |
|------|-----|
| 地图尺寸 | 10 × 10 格子 |
| 坐标系统 | (0,0) 左下角 → (9,9) 右上角 |
| 障碍物数量 | 8 个（固定布局，对称设计） |
| 出生点 | Red: (1, 1)，Blue: (8, 8) |

### 1.2 地图布局

```
  y
  9 ┌───┬───┬───┬───┬───┬───┬───┬───┬───┬───┐
    │   │   │   │   │   │   │   │   │   │   │
  8 ├───┼───┼───┼───┼───┼───┼───┼───┼───┼───┤
    │   │   │   │   │   │   │   │   │ 🔵│   │
  7 ├───┼───┼───┼───┼───┼───┼───┼───┼───┼───┤
    │   │   │ ▓ │   │   │   │   │ ▓ │   │   │
  6 ├───┼───┼───┼───┼───┼───┼───┼───┼───┼───┤
    │   │   │   │   │ ▓ │   │ ▓ │   │   │   │
  5 ├───┼───┼───┼───┼───┼───┼───┼───┼───┼───┤
    │   │   │   │   │   │   │   │   │   │   │
  4 ├───┼───┼───┼───┼───┼───┼───┼───┼───┼───┤
    │   │   │   │   │   │   │   │   │   │   │
  3 ├───┼───┼───┼───┼───┼───┼───┼───┼───┼───┤
    │   │ ▓ │   │ ▓ │   │   │   │   │   │   │
  2 ├───┼───┼───┼───┼───┼───┼───┼───┼───┼───┤
    │   │   │   │   │   │ ▓ │   │   │   │   │
  1 ├───┼───┼───┼───┼───┼───┼───┼───┼───┼───┤
    │   │ 🔴│   │   │   │   │   │   │   │   │
  0 └───┴───┴───┴───┴───┴───┴───┴───┴───┴───┘
    0   1   2   3   4   5   6   7   8   9   x

    🔴 = Red 出生点 (1, 1)
    🔵 = Blue 出生点 (8, 8)
    ▓ = 障碍物（不可通过）
```

### 1.3 障碍物分布（8 个，对称）

```
Obstacles: [(2,3), (2,7), (4,5), (6,5), (5,2), (7,3), (7,7), (3,5)]
```

设计原则：
- 对称布局，双方公平
- 障碍物集中在中部，迫使绕行或远程攻击
- 出生点附近无障碍，保证开局安全

---

## 2. Agent 属性

### 2.1 基础属性

| 属性 | 值 | 说明 |
|------|-----|------|
| HP (生命值) | 100 | 归零则阵亡 |
| 攻击力 (ATK) | 15 | 近战基础伤害 |
| 移动力 (MOV) | 2 | 每回合可移动格数 |
| 攻击范围 | 1 | 近战攻击距离（相邻格子） |
| 技能范围 | 4 | 远程技能距离 |
| 技能冷却 | 3 回合 | 使用技能后需等待 |

### 2.2 状态效果

| 状态 | 效果 | 持续时间 |
|------|------|---------|
| Charging（蓄力中） | 下回合攻击力 ×2，但无法移动 | 1 回合 |
| Stunned（眩晕） | 无法行动 | 1 回合 |
| Shielded（护盾） | 减少 50% 伤害 | 2 回合 |

---

## 3. 动作系统

### 3.1 每回合可选动作

每回合 Agent 选择 **1 个动作**：

| 动作 | 代码 | 说明 | 限制 |
|------|------|------|------|
| Move | `MOVE` | 移动 1-2 格（曼哈顿距离） | 不能穿过障碍物或其他 Agent |
| Attack | `ATTACK` | 近战攻击（相邻格子） | 目标必须在攻击范围内 |
| Skill | `SKILL` | 远程技能（4 格内） | 需要冷却完毕 |
| Charge | `CHARGE` | 蓄力，下回合攻击 ×2 | 蓄力期间无法移动 |
| Wait | `WAIT` | 原地等待 | 无效果 |

### 3.2 动作详细规则

#### Move（移动）

```
规则：
- 可移动 1 或 2 格
- 移动距离 = 曼哈顿距离 |dx| + |dy| ≤ MOV (2)
- 不能穿过障碍物
- 不能移动到对手所在格子
- 可以沿对角线移动（如 (0,0) → (1,1) 算 2 格）

示例：
  当前位置: (3, 3)
  可移动到: (3, 4), (3, 5), (4, 3), (5, 3), (4, 4), (2, 4), ...
```

#### Attack（近战攻击）

```
规则：
- 目标必须在攻击范围内（曼哈顿距离 ≤ 1）
- 伤害 = ATK (15)
- 如果 Agent 处于 Charging 状态，伤害 = ATK × 2 (30)
- 攻击后结束回合

示例：
  我的位置: (3, 3)
  对手位置: (3, 4) → 距离 1，可以攻击
  对手位置: (4, 5) → 距离 2，超出范围
```

#### Skill（远程技能）

```
规则：
- 目标必须在技能范围内（曼哈顿距离 ≤ 4）
- 伤害 = ATK × 0.8 (12)
- 附加效果：20% 概率造成 Stunned（1 回合）
- 使用后技能进入冷却（3 回合）
- 冷却期间不可使用

示例：
  我的位置: (3, 3)
  对手位置: (6, 5) → 距离 |3| + |2| = 5，超出范围
  对手位置: (5, 6) → 距离 |2| + |3| = 5，超出范围
  对手位置: (4, 6) → 距离 |1| + |3| = 4，可以施放
```

#### Charge（蓄力）

```
规则：
- 本回合不行动，进入 Charging 状态
- 下回合攻击力翻倍（Attack 或 Skill 都生效）
- 蓄力期间如果被攻击，蓄力被打断（取消 Charging 状态）
- 连续蓄力不会叠加（只有 ×2，没有 ×3）
```

#### Wait（等待）

```
规则：
- 不做任何事，直接结束回合
- 用于等待技能冷却或观察对手
```

---

## 4. 回合流程

### 4.1 单回合时序

```
┌─────────────────────────────────────────────────────────────┐
│ 回合 N 开始                                                  │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│ 1. 状态检查 (100ms)                                         │
│    - 检查是否有 Agent HP = 0 → 结束对局                      │
│    - 检查是否达到最大回合数 → 比较 HP 判定胜负               │
│                                                             │
│ 2. 决策阶段 (2000ms)                                        │
│    - Red Agent 做决策（基于策略投票）                        │
│    - Blue Agent 做决策（基于策略投票）                       │
│    - 两个 Agent 同时决策（互不知情）                         │
│                                                             │
│ 3. 动作解析 (500ms)                                         │
│    - 先手判定：速度快的先行动（相同则随机）                   │
│    - 按顺序执行动作：                                        │
│      a. 如果有 CHARGE → 设置 Charging 状态                  │
│      b. 如果有 MOVE → 更新位置                               │
│      c. 如果有 ATTACK/SKILL → 计算伤害并扣血                │
│    - 如果蓄力中被攻击 → 打断蓄力                            │
│                                                             │
│ 4. 状态更新 (200ms)                                         │
│    - 更新冷却计时器                                          │
│    - 更新状态效果持续时间                                    │
│    - 记录本回合动作 hash                                     │
│                                                             │
│ 5. 广播 (100ms)                                             │
│    - WebSocket 推送回合状态到所有客户端                      │
│    - 前端渲染动画                                            │
│                                                             │
│ 等待 5 秒 → 进入下一回合                                     │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 4.2 同时行动与冲突解决

```
场景 1：双方都选择移动
  → 各自移动，无冲突

场景 2：双方都选择攻击，且目标在范围内
  → 先手方先攻击，如果击杀则后手方无法攻击

场景 3：双方都想移动到同一格
  → 先手方移动成功，后手方移动失败（原地等待）

场景 4：一方移动，一方攻击
  → 先攻后动（攻击先结算，再移动）
```

### 4.3 先手判定

```
规则：
1. 如果一方处于 Charging 状态且被攻击 → 被攻击方先手（反击）
2. 其他情况：比较两个 Agent 的"速度值"
3. 速度值 = 基础速度 (5) + 随机 (1-3)
4. 速度相同 → 随机决定
```

---

## 5. 伤害计算

### 5.1 基础伤害公式

```
Attack (近战):
  damage = ATK = 15
  如果攻击方 Charging: damage = ATK × 2 = 30

Skill (远程):
  damage = ATK × 0.8 = 12
  如果攻击方 Charging: damage = ATK × 2 × 0.8 = 24
  20% 概率附加 Stunned (1 回合)
```

### 5.2 伤害减免

```
如果目标有 Shielded 状态:
  final_damage = damage × 0.5 (向下取整)

如果目标 Charging 且被打断:
  final_damage = damage × 1.2 (蓄力被打断额外受 20% 伤害)
```

### 5.3 击杀判定

```
if agent.HP <= 0:
    agent.HP = 0
    agent.status = DEAD
    对局结束，对方获胜
```

---

## 6. 胜负判定

### 6.1 胜利条件

```
条件 1：一方 HP 归零
  → 另一方获胜

条件 2：达到最大回合数 (30)
  → 比较双方 HP
  → HP 高者获胜
  → HP 相同 → 平局（下注退还 50%）
```

### 6.2 回合数限制

```
max_rounds = 30

设计理由：
- 每回合 5 秒 → 30 回合 = 2.5 分钟（不含决策时间）
- 加上 AI 决策时间 (~2s/回合) → 总时长 ~3.5 分钟
- 足够产生有意义的对战，不会拖太久
```

---

## 7. 策略投票机制

### 7.1 投票规则

```
下注时，用户必须选择策略投票：
  - Aggressive (激进)
  - Defensive (稳健)
  - Tricky (诡道)

投票权重 = 下注金额 / 该方总下注金额
```

### 7.2 策略混合

```
示例：
  Red 方总下注: 100 USDC
  - 激进: 40 USDC (40%)
  - 稳健: 35 USDC (35%)
  - 诡道: 25 USDC (25%)

→ Red Agent 的决策权重:
  - 激进倾向 +40%
  - 稳健倾向 +35%
  - 诡道倾向 +25%

Agent 的 AI 决策引擎会根据这些权重调整行为概率分布。
```

### 7.3 策略对决策的影响

```
Aggressive (激进):
  - 优先选择 Attack (如果目标在范围内)
  - 倾向移动靠近对手
  - 较少选择 Wait 或 Charge

Defensive (稳健):
  - 优先保持距离
  - 倾向使用 Skill 远程攻击
  - 较多选择 Wait 等待技能冷却

Tricky (诡道):
  - 倾向 Charge 蓄力
  - 寻找绕后机会
  - 利用障碍物制造战术优势
```

### 7.4 投票规则

```
- 对局开始前 2 分钟：开放下注 + 策略投票
- 策略投票为一次性操作，开局时确定，整场不变
- 对局进行中可追加下注，但不能修改策略投票
- 策略投票不影响收益分配，纯粹影响 AI 决策行为（参与感）
- 押对方向（Red/Blue）即获得收益，与策略选择无关
```

---

## 8. 数据结构定义

### 8.1 游戏状态 (GameState)

```go
type GameState struct {
    GameID       uint64       `json:"game_id"`
    Status       GameStatus   `json:"status"`
    CurrentRound uint8        `json:"current_round"`
    MaxRounds    uint8        `json:"max_rounds"`
    AgentRed     AgentState   `json:"agent_red"`
    AgentBlue    AgentState   `json:"agent_blue"`
    BettingRed   uint256      `json:"betting_red"`
    BettingBlue  uint256      `json:"betting_blue"`
    StrategyRed  StrategyVote `json:"strategy_red"`
    StrategyBlue StrategyVote `json:"strategy_blue"`
    History      []TurnRecord `json:"history"`
}
```

### 8.2 Agent 状态 (AgentState)

```go
type AgentState struct {
    ID           uint64     `json:"id"`
    Name         string     `json:"name"`
    Personality  string     `json:"personality"`
    HP           uint8      `json:"hp"`
    MaxHP        uint8      `json:"max_hp"`
    Position     Position   `json:"position"`
    Status       []Effect   `json:"status"`
    SkillCooldown uint8     `json:"skill_cooldown"`
    IsCharging   bool       `json:"is_charging"`
    Wins         uint64     `json:"wins"`
    Losses       uint64     `json:"losses"`
}
```

### 8.3 回合记录 (TurnRecord)

```go
type TurnRecord struct {
    Round      uint8      `json:"round"`
    RedAction  Action     `json:"red_action"`
    BlueAction Action     `json:"blue_action"`
    RedHPAfter uint8      `json:"red_hp_after"`
    BlueHPAfter uint8     `json:"blue_hp_after"`
    ActionHash [32]byte   `json:"action_hash"` // 用于链上验证
}
```

### 8.4 动作 (Action)

```go
type Action struct {
    Type   ActionType `json:"type"` // MOVE, ATTACK, SKILL, CHARGE, WAIT
    Target Position   `json:"target,omitempty"` // 移动目标或攻击目标
}

type ActionType string
const (
    ActionMove   ActionType = "MOVE"
    ActionAttack ActionType = "ATTACK"
    ActionSkill  ActionType = "SKILL"
    ActionCharge ActionType = "CHARGE"
    ActionWait   ActionType = "WAIT"
)
```

### 8.5 策略投票 (StrategyVote)

```go
type StrategyVote struct {
    Aggressive uint256 `json:"aggressive"` // 激进投票总额
    Defensive  uint256 `json:"defensive"`  // 稳健投票总额
    Tricky     uint256 `json:"tricky"`     // 诡道投票总额
}

// 计算权重比例
func (sv StrategyVote) Weights() (float64, float64, float64) {
    total := sv.Aggressive + sv.Defensive + sv.Tricky
    if total == 0 {
        return 0.33, 0.33, 0.34 // 默认均分
    }
    return float64(sv.Aggressive) / float64(total),
           float64(sv.Defensive) / float64(total),
           float64(sv.Tricky) / float64(total)
}
```

---

## 9. 随机性控制

### 9.1 需要随机的地方

| 场景 | 随机内容 | 范围 |
|------|---------|------|
| 先手判定 | 速度相同时 | 1-2 |
| Skill Stunned 效果 | 是否触发 | 20% 概率 |
| AI 决策 | 同等权重动作选择 | 按权重随机 |

### 9.2 随机种子

```
使用区块 hash 作为随机种子（链上可验证）
seed = blockhash(block.number - 1)

→ 保证随机性不可预测，但事后可验证
```

---

## 10. 边界情况处理

| 情况 | 处理方式 |
|------|---------|
| 双方同时 HP 归零 | 先触发的一方输 |
| 30 回合 HP 相同 | 平局，退还 50% 下注 |
| Agent 被障碍物包围无法移动 | 只能选择 Attack/Skill/Wait |
| Skill 冷却中 + 对手在远程范围外 | 只能 Move 或 Wait |
| AI API 调用失败 | 使用默认策略（Wait） |
| WebSocket 断开 | 客户端自动重连，从合约读取状态 |

---

**文档版本**：v1.0
**最后更新**：2026-08-02
**状态**：待评审
