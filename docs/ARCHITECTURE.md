# 系统架构文档

## 架构总览

```
┌─────────────────────────────────────────────────────────────────────┐
│                        Frontend (Next.js)                           │
│                                                                     │
│  首页(对局列表) → 对局详情(观战+下注) → Agent 详情 → 个人中心       │
│                                                                     │
│  技术: Next.js 14 + Tailwind + shadcn/ui + wagmi + Framer Motion   │
└──────────────────────────────┬──────────────────────────────────────┘
                               │
                  REST API (HTTP) + WebSocket
                               │
┌──────────────────────────────▼──────────────────────────────────────┐
│                      Backend (Go + Gin)                             │
│                                                                     │
│  ┌───────────┐  ┌───────────┐  ┌───────────┐  ┌────────────────┐  │
│  │ API Layer │  │  Game     │  │  Agent    │  │  Blockchain    │  │
│  │ (REST+WS) │  │  Engine   │  │  (LLM)   │  │  Service       │  │
│  └─────┬─────┘  └─────┬─────┘  └─────┬─────┘  └───────┬────────┘  │
│        │               │              │                 │           │
│  ┌─────▼───────────────▼──────────────▼─────────────────▼────────┐ │
│  │                     Service Layer                              │ │
│  │  GameService / BetService / AgentService / ChainService       │ │
│  └─────────────────────────────┬─────────────────────────────────┘ │
│                                │                                    │
│  ┌─────────────────────────────▼─────────────────────────────────┐ │
│  │                     Data Layer (In-Memory)                     │ │
│  └───────────────────────────────────────────────────────────────┘ │
└──────────────────────────────┬──────────────────────────────────────┘
                               │
                    ethclient (go-ethereum)
                               │
┌──────────────────────────────▼──────────────────────────────────────┐
│                    Blockchain (Sepolia)                              │
│                                                                     │
│  ┌───────────────────┐  ┌─────────────────┐  ┌──────────────────┐  │
│  │ AgentArena.sol    │  │ BettingPool.sol │  │ StrategyVoting   │  │
│  │ (主合约/入口)      │→ │ (下注托管)       │  │ (策略投票)        │  │
│  └───────────────────┘  └─────────────────┘  └──────────────────┘  │
│  ┌───────────────────┐  ┌─────────────────┐                        │
│  │ GameRegistry.sol  │  │ USDC (ERC20)    │                        │
│  │ (对局+Agent记录)   │  │ (下注币种)       │                        │
│  └───────────────────┘  └─────────────────┘                        │
└─────────────────────────────────────────────────────────────────────┘
                               │
                        Qwen API (通义千问)
                               │
┌──────────────────────────────▼──────────────────────────────────────┐
│                    AI Agent Layer                                    │
│                                                                     │
│  每回合：构建 Prompt → 调用 Qwen API → 解析 Action → 验证合法性    │
│  两 Agent 并行调用，独白异步生成                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 模块依赖

```
frontend
  ├── → backend (REST API + WebSocket)
  └── → contracts (直接调用合约, via wagmi)

backend
  ├── → contracts (ethclient, 写操作)
  ├── → contracts (事件监听, 读操作)
  └── → Qwen API (Agent 决策)

contracts
  └── → USDC (ERC20 approve/transferFrom)
```

---

## 数据流

### 下注流程
```
用户钱包 → approve USDC → 调用合约 betAndVote()
                          → BettingPool: 锁定 USDC
                          → StrategyVoting: 记录投票
         → 前端轮询/WebSocket 更新下注池
```

### 对局流程
```
后端创建对局 → 合约 createGame()
等待下注 → 合约 startGame() (锁定)
游戏循环 → 每回合:
  1. 从合约读取策略权重 (只读)
  2. 两 Agent 并行调 Qwen → 得到动作
  3. 游戏引擎执行回合
  4. WebSocket 推送状态
  5. 检查胜负
对局结束 → 合约 finishGame() (提交结果 + 结算)
用户提取 → 合约 claim()
```

---

## 链上 vs 链下

| 逻辑 | 位置 | 原因 |
|------|------|------|
| 下注锁定/提取 | 链上 | 需要信任less托管 |
| 策略投票记录 | 链上 | 不可篡改 |
| 对局结果 | 链上 | 可验证 |
| 游戏状态计算 | 链下(Go) | 性能要求高 |
| AI 决策 | 链下(Qwen) | 需要 LLM |
| WebSocket 推送 | 链下(Go) | 实时性要求 |

---

## 关键接口

### 前后端接口
- REST: `GET /api/games`, `POST /api/bets`, `GET /api/agents`
- WebSocket: `ws://host/ws?game_id=X`

### 后端 ↔ 合约接口
- `betAndVote(gameId, side, amount, strategy)`
- `startGame(gameId)`
- `finishGame(gameId, redWins, actionsHash)`
- `getStrategyWeights(gameId) → (agg, def, trick)`

### 后端 ↔ AI 接口
- `QwenClient.Chat(prompt) → response`
- Prompt 构建: `buildDecisionPrompt(agent, state, weights)`
- 响应解析: `ParseAction(response) → Action`

---

## 用户系统 + 钱包（迭代 016+）

### 双币种体系
- **AC (ArenaCoin)** — 链上 ERC20 (18 decimals)，每日可领 100 AC，用于挑战押金 / 奖励
- **USDC (MockUSDC)** — 链上 ERC20 (6 decimals)，用于链上下注

### 后端服务
- `ACService` — ArenaCoin 交互（余额 / mint / 从 treasury 转账）
- `USDCService` — MockUSDC 只读交互（仅余额查询）
- `EthChainService` — AgentArena 合约交互（CreateGame / StartGame / FinishGame）

### 关键 API
- `POST /api/auth/login` — 钱包签名登录，返回 JWT
- `GET /api/auth/profile` — 用户信息 + AC 余额（链上+链下）+ USDC 余额
- `GET /api/auth/earnings` — 收益历史（挑战 + 奖励明细）
- `POST /api/auth/claim-daily` — 每日领取 100 AC
- `POST /api/auth/withdraw` — 提现 AC 到链上钱包
- `POST /api/challenges` — 创建 PVE 挑战（AC 赌注）

### 前端页面
- `/wallet` — 我的钱包（USDC/AC 余额 + 收益历史）
- `/profile` — 个人中心（领取 AC / 提现 / 管理自定义 Agent）
- Navbar 钱包胶囊 — 登录后显示 USDC + AC 余额，点击跳转钱包
