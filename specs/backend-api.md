# S4 - 后端 API 规格 (Backend API Specification)

## Agent Arena Go Backend

---

## 1. 架构总览

```
┌──────────────────────────────────────────────────────────────┐
│                     Go Backend                               │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌────────────────────┐ │
│  │  HTTP API    │  │  WebSocket   │  │  Game Engine       │ │
│  │  (REST)      │  │  (实时推送)   │  │  (核心逻辑)        │ │
│  │              │  │              │  │                    │ │
│  │ /api/games   │  │ 游戏状态推送  │  │ 回合推进           │ │
│  │ /api/bets    │  │ 回合事件推送  │  │ 动作解析           │ │
│  │ /api/agents  │  │ 赔率变化推送  │  │ 胜负判定           │ │
│  └──────┬───────┘  └──────┬───────┘  └────────┬───────────┘ │
│         │                 │                    │             │
│  ┌──────▼─────────────────▼────────────────────▼───────────┐ │
│  │                   Service Layer                          │ │
│  │  GameService / BetService / AgentService / ChainService │ │
│  └──────────────────────────┬──────────────────────────────┘ │
│                             │                                │
│  ┌──────────────────────────▼──────────────────────────────┐ │
│  │                   Data Layer                             │ │
│  │         In-Memory Store (Redis optional)                 │ │
│  └──────────────────────────┬──────────────────────────────┘ │
│                             │                                │
│  ┌──────────────────────────▼──────────────────────────────┐ │
│  │               Blockchain Layer                           │ │
│  │         Contract Interactions (ethers.go)                │ │
│  └─────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────┘
```

---

## 2. HTTP REST API

### 2.1 对局相关

#### GET /api/games

获取所有对局列表。

```
Request:
  GET /api/games?status=open|playing|finished&limit=20&offset=0

Response:
{
  "games": [
    {
      "game_id": 1,
      "agent_red": {
        "id": "0xabc...",
        "name": "Berserker",
        "personality": "berserker",
        "wins": 12,
        "losses": 5,
        "win_rate": 70.6
      },
      "agent_blue": {
        "id": "0xdef...",
        "name": "Tactician",
        "personality": "tactician",
        "wins": 8,
        "losses": 9,
        "win_rate": 47.1
      },
      "status": "open",
      "total_bet_red": "50000000",     // 50 USDC (6 decimals)
      "total_bet_blue": "30000000",    // 30 USDC
      "betting_deadline": 1722600000,  // Unix timestamp
      "current_round": 0,
      "max_rounds": 30,
      "odds_red": 1.6,                 // 基于下注池计算
      "odds_blue": 2.67
    }
  ],
  "total": 15
}
```

#### GET /api/games/:gameId

获取单场对局详情。

```
Request:
  GET /api/games/1

Response:
{
  "game_id": 1,
  "agent_red": { ... },
  "agent_blue": { ... },
  "status": "playing",
  "current_round": 12,
  "max_rounds": 30,
  "total_bet_red": "120000000",
  "total_bet_blue": "80000000",
  "agent_red_state": {
    "hp": 65,
    "max_hp": 100,
    "position": {"x": 4, "y": 5},
    "status": [],
    "skill_cooldown": 0,
    "is_charging": false
  },
  "agent_blue_state": {
    "hp": 42,
    "max_hp": 100,
    "position": {"x": 7, "y": 3},
    "status": ["stunned"],
    "skill_cooldown": 2,
    "is_charging": false
  },
  "strategy_red": {
    "aggressive": 40,
    "defensive": 35,
    "tricky": 25
  },
  "strategy_blue": {
    "aggressive": 10,
    "defensive": 20,
    "tricky": 70
  },
  "history": [
    {
      "round": 1,
      "red_action": {"type": "MOVE", "target": {"x": 2, "y": 2}},
      "blue_action": {"type": "CHARGE"},
      "red_hp_after": 100,
      "blue_hp_after": 100
    },
    ...
  ]
}
```

---

### 2.2 下注相关

#### POST /api/bets/estimate

估算下注后的赔率和收益（不实际下注）。

```
Request:
{
  "game_id": 1,
  "side": "red",
  "amount": "10000000"    // 10 USDC
}

Response:
{
  "current_pool_red": "120000000",
  "current_pool_blue": "80000000",
  "new_pool_red": "130000000",
  "potential_reward": "16666666",   // 如果 Red 赢，可获得 ~16.67 USDC
  "new_odds_red": 1.54,
  "new_odds_blue": 2.83
}
```

#### POST /api/bets

提交下注（后端代合约发起链上交易）。

```
Request:
{
  "game_id": 1,
  "side": "red",
  "strategy": "aggressive",
  "amount": "10000000",
  "signature": "0x..."        // 用户签名，证明身份
}

Response:
{
  "tx_hash": "0x123...",
  "status": "pending",
  "bet_id": "bet_1_0xabc..."
}

WebSocket 推送:
{
  "type": "bet_placed",
  "game_id": 1,
  "user": "0xabc...",
  "side": "red",
  "amount": "10000000",
  "strategy": "aggressive"
}
```

**注意**：实际下注需要用户先 approve USDC，然后前端直接调用合约。此 API 仅用于记录后端状态。

---

### 2.3 Agent 相关

#### GET /api/agents

获取所有 Agent 列表。

```
Response:
{
  "agents": [
    {
      "id": "0xabc...",
      "name": "Berserker",
      "personality": "berserker",
      "wins": 12,
      "losses": 5,
      "win_rate": 70.6,
      "description": "激进攻型，喜欢近战冲锋",
      "avatar_url": "/avatars/berserker.png"
    },
    ...
  ]
}
```

#### GET /api/agents/:agentId

获取单个 Agent 详情。

```
Response:
{
  "id": "0xabc...",
  "name": "Berserker",
  "personality": "berserker",
  "wins": 12,
  "losses": 5,
  "win_rate": 70.6,
  "description": "激进攻型，喜欢近战冲锋",
  "avatar_url": "/avatars/berserker.png",
  "recent_games": [
    {
      "game_id": 1,
      "opponent": "Tactician",
      "result": "win",
      "rounds": 22
    }
  ]
}
```

---

### 2.4 用户相关

#### GET /api/users/:address/bets

获取用户的下注历史。

```
Response:
{
  "address": "0xabc...",
  "bets": [
    {
      "game_id": 1,
      "side": "red",
      "amount": "10000000",
      "strategy": "aggressive",
      "status": "won",
      "reward": "16666666",
      "claimed": false
    }
  ]
}
```

#### GET /api/users/:address/claims

获取用户可提取的奖金列表。

```
Response:
{
  "address": "0xabc...",
  "claimable": [
    {
      "game_id": 1,
      "reward": "16666666",
      "status": "unclaimed"
    }
  ],
  "total_claimable": "16666666"
}
```

---

## 3. WebSocket API

### 3.1 连接

```
URL: ws://localhost:8080/ws?game_id=1

前端连接后，会收到当前游戏状态的快照，然后持续收到实时更新。
```

### 3.2 消息类型

#### 服务端 → 客户端

##### game_state（完整状态快照）

```json
{
  "type": "game_state",
  "data": {
    "game_id": 1,
    "status": "playing",
    "current_round": 12,
    "agent_red": { "hp": 65, "position": {"x": 4, "y": 5}, ... },
    "agent_blue": { "hp": 42, "position": {"x": 7, "y": 3}, ... }
  }
}
```

##### turn_update（回合更新）

```json
{
  "type": "turn_update",
  "data": {
    "round": 13,
    "red_action": {
      "type": "ATTACK",
      "target": {"x": 5, "y": 5},
      "damage": 15,
      "hit": true
    },
    "blue_action": {
      "type": "MOVE",
      "target": {"x": 6, "y": 2}
    },
    "red_hp": 65,
    "blue_hp": 27,
    "effects": [
      {"target": "blue", "type": "damage", "value": 15}
    ]
  }
}
```

##### betting_update（下注池变化）

```json
{
  "type": "betting_update",
  "data": {
    "total_bet_red": "130000000",
    "total_bet_blue": "80000000",
    "odds_red": 1.54,
    "odds_blue": 2.83
  }
}
```

##### game_started（对局开始）

```json
{
  "type": "game_started",
  "data": {
    "game_id": 1,
    "status": "playing",
    "betting_locked": true,
    "strategy_red": {"aggressive": 40, "defensive": 35, "tricky": 25},
    "strategy_blue": {"aggressive": 10, "defensive": 20, "tricky": 70}
  }
}
```

##### game_finished（对局结束）

```json
{
  "type": "game_finished",
  "data": {
    "game_id": 1,
    "winner": "red",
    "winner_name": "Berserker",
    "total_rounds": 22,
    "final_hp_red": 45,
    "final_hp_blue": 0,
    "total_pool": "210000000",
    "protocol_fee": "10500000",
    "winner_pool": "130000000"
  }
}
```

##### action_hash（动作 hash 提交）

```json
{
  "type": "action_hash",
  "data": {
    "game_id": 1,
    "tx_hash": "0x456...",
    "actions_hash": "0x789...",
    "status": "confirmed"
  }
}
```

#### 客户端 → 服务端

##### subscribe（订阅对局）

```json
{
  "type": "subscribe",
  "game_id": 1
}
```

##### unsubscribe（取消订阅）

```json
{
  "type": "unsubscribe",
  "game_id": 1
}
```

---

## 4. 核心服务层

### 4.1 GameService

```go
type GameService struct {
    store       Store
    engine      *GameEngine
    chain       *ChainService
    wsHub       *WebSocketHub
}

// 创建新对局
func (s *GameService) CreateGame(agentRed, agentBlue string) (*Game, error)

// 开始对局（锁定下注，启动游戏循环）
func (s *GameService) StartGame(gameID uint64) error

// 游戏主循环（每 5 秒一回合）
func (s *GameService) gameLoop(game *Game) {
    for game.CurrentRound < game.MaxRounds {
        // 1. 获取策略权重
        weights := s.chain.GetStrategyWeights(game.ID)

        // 2. AI 决策
        redAction := s.agentDecide(game.AgentRed, game, weights.Red)
        blueAction := s.agentDecide(game.AgentBlue, game, weights.Blue)

        // 3. 执行回合
        turnResult := s.engine.ExecuteTurn(game, redAction, blueAction)

        // 4. 广播
        s.wsHub.Broadcast(game.ID, turnResult)

        // 5. 检查胜负
        if s.engine.CheckWin(game) {
            break
        }

        // 6. 等待 5 秒
        time.Sleep(5 * time.Second)
    }

    // 对局结束
    s.finishGame(game)
}

// 结束对局
func (s *GameService) finishGame(game *Game) {
    // 1. 计算动作 hash
    actionsHash := computeActionsHash(game.History)

    // 2. 提交到链上
    s.chain.SubmitResult(game.ID, game.Winner, actionsHash)

    // 3. 结算
    s.chain.Settle(game.ID, game.Winner == "red")

    // 4. 广播
    s.wsHub.Broadcast(game.ID, GameFinishedEvent{...})
}
```

### 4.2 ChainService

```go
type ChainService struct {
    client     *ethclient.Client
    contracts  *ContractBindings
    privateKey *ecdsa.PrivateKey
}

// 下注（后端辅助，实际由用户前端直接调用合约）
func (c *ChainService) PlaceBet(gameID uint64, side bool, amount *big.Int) (string, error)

// 获取策略权重（从合约读取）
func (c *ChainService) GetStrategyWeights(gameID uint64) (*StrategyWeights, error)

// 提交对局结果
func (c *ChainService) SubmitResult(gameID uint64, actionsHash [32]byte) (string, error)

// 结算
func (c *ChainService) Settle(gameID uint64, redWins bool) (string, error)

// 监听合约事件
func (c *ChainService) ListenEvents() <-chan ContractEvent
```

### 4.3 WebSocketHub

```go
type WebSocketHub struct {
    rooms map[uint64]map[*Client]bool  // gameID -> clients
    mu    sync.RWMutex
}

// 客户端连接
func (h *WebSocketHub) HandleConnect(conn *websocket.Conn, gameID uint64)

// 广播消息到房间
func (h *WebSocketHub) Broadcast(gameID uint64, msg interface{})

// 客户端断开
func (h *WebSocketHub) HandleDisconnect(client *Client)
```

---

## 5. 数据存储

### 5.1 内存存储（MVP）

```go
type Store interface {
    // Games
    CreateGame(game *Game) error
    GetGame(gameID uint64) (*Game, error)
    ListGames(filter GameFilter) ([]*Game, error)
    UpdateGame(game *Game) error

    // Agents
    GetAgent(agentID string) (*Agent, error)
    ListAgents() ([]*Agent, error)
    UpdateAgent(agent *Agent) error

    // Bets (本地缓存，链上为准)
    CacheBet(bet *Bet) error
    GetBetsByUser(address string) ([]*Bet, error)
}

// 实现：MemoryStore (MVP)
type MemoryStore struct {
    games   map[uint64]*Game
    agents  map[string]*Agent
    bets    map[string][]*Bet  // address -> bets
    mu      sync.RWMutex
}
```

### 5.2 数据同步

```
链上事件 → ChainService.ListenEvents() → 更新本地 Store → 推送 WebSocket

关键事件监听：
  - BetPlaced: 更新本地 bet 缓存
  - BettingLocked: 更新游戏状态
  - GameSettled: 更新赢家信息
```

---

## 6. 配置

```yaml
# config.yaml
server:
  http_port: 8080
  ws_port: 8080

blockchain:
  rpc_url: "${SEPOLIA_RPC_URL}"
  chain_id: 11155111
  private_key: "${PRIVATE_KEY}"
  contracts:
    arena: "${ARENA_CONTRACT_ADDRESS}"
    betting_pool: "${BETTING_POOL_ADDRESS}"
    strategy_voting: "${STRATEGY_VOTING_ADDRESS}"
    game_registry: "${GAME_REGISTRY_ADDRESS}"
    usdc: "${USDC_ADDRESS}"

game:
  turn_duration: 5s
  max_rounds: 30
  betting_duration: 2m

ai:
  qwen_api_key: "${QWEN_API_KEY}"
  qwen_model: "qwen-plus"
  timeout: 10s
```

---

## 7. 错误码

| 错误码 | 说明 |
|--------|------|
| 40001 | 对局不存在 |
| 40002 | 下注已截止 |
| 40003 | 金额无效 |
| 40004 | 对局已结束，不可下注 |
| 40005 | 策略投票已锁定 |
| 50001 | 链上交易失败 |
| 50002 | AI 决策超时 |
| 50003 | WebSocket 连接失败 |

---

**文档版本**：v1.0
**最后更新**：2026-08-02
**状态**：待评审
