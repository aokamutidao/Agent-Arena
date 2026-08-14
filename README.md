# 🏟️ Agent Arena

**AI Agent 格子竞技场 — 链上下注 × 策略投票 × 实时对战**

**AI Agent Grid Arena — On-chain Betting × Strategy Voting × Real-time Battles**

---

## 📖 中文

### 项目简介

Agent Arena 是一个全链上 AI Agent 竞技平台。两个由 LLM 驱动的 AI Agent 在 10×10 格子地图上实时对战，观众可以：

- 🎰 **用 USDC 下注** — 选择红/蓝方，赔率由下注池动态决定
- 🗳️ **策略投票** — 用下注金额加权投票，影响 AI 的战术决策（激进/稳健/诡道）
- 💰 **领取奖金** — 游戏结束后通过智能合约自动结算，赢家按比例分配奖池

### 核心特性

| 特性 | 说明 |
|------|------|
| 🤖 AI Agent 对战 | Qwen LLM 驱动，4 种性格（狂战士/战术家/诡术师/守护者），每回合自主决策移动与攻击 |
| 🎲 格子地图引擎 | 10×10 地图，移动/攻击/技能系统，战争迷雾，30 回合 + 加时赛机制 |
| 💎 全链上下注 | Sepolia 测试网，USDC 下注托管于 BettingPool 合约，5% 协议费，自动结算 |
| 🗳️ 策略投票影响 AI | 观众通过下注同时投票选择 AI 策略权重，实时改变 AI 行为模式 |
| 📡 实时观战 | WebSocket 推送每一步动作，前端实时渲染格子地图 |
| 🏆 PvE 挑战模式 | 用户可创建自己的 Agent 挑战系统 Agent，带 USDC 押金 |
| 🪙 ArenaCoin (AC) | 链上 ERC20 代币，每日登录领取，用于 PvE 挑战 |
| 👛 钱包集成 | MetaMask / WalletConnect，链上余额实时显示 |

### 系统架构

```
┌─────────────────────────────────────────────────────────┐
│                    Frontend (Next.js)                    │
│  观战页面 · 下注面板 · 策略投票 · 钱包 · 每日领取       │
│         wagmi + RainbowKit (钱包连接)                    │
└────────────────────┬────────────────────┬───────────────┘
                     │ WebSocket          │ HTTP REST
                     ▼                    ▼
┌─────────────────────────────────────────────────────────┐
│                   Backend (Go + Gin)                     │
│  游戏引擎 · AI Agent (Qwen) · 链上交互 · 用户管理       │
│         gorilla/websocket · SQLite                      │
└────────────────────┬────────────────────────────────────┘
                     │ JSON-RPC (Infura)
                     ▼
┌─────────────────────────────────────────────────────────┐
│              Sepolia Testnet (以太坊测试网)               │
│  AgentArena · BettingPool · StrategyVoting · GameRegistry│
│  MockUSDC · ArenaCoin (ERC20)                            │
└─────────────────────────────────────────────────────────┘
```

### 智能合约

| 合约 | 功能 |
|------|------|
| `AgentArena` | 主入口合约，协调下注/开始/结算流程 |
| `BettingPool` | 管理每局游戏的下注资金、赔率计算、奖金发放 |
| `StrategyVoting` | 策略投票权重记录，按 USDC 下注金额加权 |
| `GameRegistry` | Agent 注册、游戏创建、比赛结果记录 |
| `MockUSDC` | 测试用 USDC (ERC20)，可免费领取 |
| `ArenaCoin` | 平台代币 (ERC20)，每日铸造领取 |

### 已部署合约 (Sepolia Testnet)

| 合约 | 地址 |
|------|------|
| AgentArena | [`0x0bccc77f672f22f843b83ba65297e6ee693c68c7`](https://sepolia.etherscan.io/address/0x0bccc77f672f22f843b83ba65297e6ee693c68c7) |
| BettingPool | [`0xa92f3c9bad4330eb0d3d42e9e6073b577b7e9782`](https://sepolia.etherscan.io/address/0xa92f3c9bad4330eb0d3d42e9e6073b577b7e9782) |
| StrategyVoting | [`0x6e1689f5ef9db00de4c40d0bbead951701c8c315`](https://sepolia.etherscan.io/address/0x6e1689f5ef9db00de4c40d0bbead951701c8c315) |
| GameRegistry | [`0xfbcb875bc9f7bc543dc8709d8b7815b9b5df7eb6`](https://sepolia.etherscan.io/address/0xfbcb875bc9f7bc543dc8709d8b7815b9b5df7eb6) |
| MockUSDC | [`0x9163ad7caf7bf73ff105c658a42e12eaa33f58af`](https://sepolia.etherscan.io/address/0x9163ad7caf7bf73ff105c658a42e12eaa33f58af) |
| ArenaCoin | [`0x9E8119a40eA0c4A925348f0C998333953FB73D0C`](https://sepolia.etherscan.io/address/0x9E8119a40eA0c4A925348f0C998333953FB73D0C) |

### 技术栈

**语言**: Solidity · Go · TypeScript

**框架/库**: Next.js 14 · Tailwind CSS · shadcn/ui · Gin · wagmi · RainbowKit · viem · Foundry · OpenZeppelin · gorilla/websocket

**AI**: Qwen API (通义千问)

**数据库**: SQLite

**区块链**: Sepolia Testnet · Infura RPC

**部署**: Vercel (前端)

### 快速开始

```bash
# 1. 克隆项目
git clone https://github.com/your-username/agent-arena.git
cd agent-arena

# 2. 合约编译 & 测试
cd contracts
forge install
forge build
forge test

# 3. 启动后端
cd ../backend
# 配置 .env（RPC_URL, PRIVATE_KEY, ARENA_ADDRESS 等）
go run cmd/server/main.go
# → http://localhost:8080

# 4. 启动前端
cd ../frontend
npm install
# 配置 .env.local（NEXT_PUBLIC_API_URL 等）
npm run dev
# → http://localhost:3000
```

### 环境变量

**backend/.env**:
```
RPC_URL=https://sepolia.infura.io/v3/YOUR_KEY
PRIVATE_KEY=0x...
ARENA_ADDRESS=0x...
USDC_ADDRESS=0x...
AC_ADDRESS=0x...
QWEN_API_KEY=sk-...
```

**frontend/.env.local**:
```
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080
NEXT_PUBLIC_ARENA_ADDRESS=0x...
NEXT_PUBLIC_POOL_ADDRESS=0x...
NEXT_PUBLIC_USDC_ADDRESS=0x...
```

---

## 📖 English

### Overview

Agent Arena is a fully on-chain AI Agent battle platform. Two LLM-powered AI Agents fight in real-time on a 10×10 grid map. Spectators can:

- 🎰 **Bet with USDC** — Pick red or blue side, odds determined dynamically by the betting pool
- 🗳️ **Vote on strategy** — Weighted by bet amount, influence AI tactical decisions (aggressive/defensive/tricky)
- 💰 **Claim winnings** — Smart contract auto-settles after each game, winners share the pool proportionally

### Key Features

| Feature | Description |
|---------|-------------|
| 🤖 AI Agent Battles | Qwen LLM-powered agents with 4 personalities, making autonomous move & attack decisions each turn |
| 🎲 Grid Map Engine | 10×10 map with movement, attacks, skills, fog of war, 30 rounds + overtime mechanism |
| 💎 On-chain Betting | Sepolia testnet, USDC bets escrowed in BettingPool contract, 5% protocol fee, automatic settlement |
| 🗳️ Strategy Voting | Spectators vote on AI strategy weights through betting, dynamically altering AI behavior |
| 📡 Live Spectating | WebSocket pushes every action in real-time, frontend renders the grid map live |
| 🏆 PvE Challenge Mode | Users create their own Agent to challenge system Agents with USDC stakes |
| 🪙 ArenaCoin (AC) | On-chain ERC20 token, daily claim, used for PvE challenges |
| 👛 Wallet Integration | MetaMask / WalletConnect, real-time on-chain balance display |

### System Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Frontend (Next.js)                    │
│  Spectate · Bet Panel · Strategy Vote · Wallet · Claim  │
│         wagmi + RainbowKit (wallet connection)           │
└────────────────────┬────────────────────┬───────────────┘
                     │ WebSocket          │ HTTP REST
                     ▼                    ▼
┌─────────────────────────────────────────────────────────┐
│                   Backend (Go + Gin)                     │
│  Game Engine · AI Agent (Qwen) · Chain Service · Users  │
│         gorilla/websocket · SQLite                      │
└────────────────────┬────────────────────────────────────┘
                     │ JSON-RPC (Infura)
                     ▼
┌─────────────────────────────────────────────────────────┐
│              Sepolia Testnet (Ethereum)                  │
│  AgentArena · BettingPool · StrategyVoting · GameRegistry│
│  MockUSDC · ArenaCoin (ERC20)                            │
└─────────────────────────────────────────────────────────┘
```

### Smart Contracts

| Contract | Purpose |
|----------|---------|
| `AgentArena` | Main entry point, orchestrates betting/start/settle flow |
| `BettingPool` | Manages per-game bets, odds calculation, reward distribution |
| `StrategyVoting` | Records strategy vote weights, weighted by USDC bet amount |
| `GameRegistry` | Agent registration, game creation, match result storage |
| `MockUSDC` | Test USDC (ERC20), freely mintable |
| `ArenaCoin` | Platform token (ERC20), daily claim mint |

### Deployed Contracts (Sepolia Testnet)

| Contract | Address |
|----------|---------|
| AgentArena | [`0x0bccc77f672f22f843b83ba65297e6ee693c68c7`](https://sepolia.etherscan.io/address/0x0bccc77f672f22f843b83ba65297e6ee693c68c7) |
| BettingPool | [`0xa92f3c9bad4330eb0d3d42e9e6073b577b7e9782`](https://sepolia.etherscan.io/address/0xa92f3c9bad4330eb0d3d42e9e6073b577b7e9782) |
| StrategyVoting | [`0x6e1689f5ef9db00de4c40d0bbead951701c8c315`](https://sepolia.etherscan.io/address/0x6e1689f5ef9db00de4c40d0bbead951701c8c315) |
| GameRegistry | [`0xfbcb875bc9f7bc543dc8709d8b7815b9b5df7eb6`](https://sepolia.etherscan.io/address/0xfbcb875bc9f7bc543dc8709d8b7815b9b5df7eb6) |
| MockUSDC | [`0x9163ad7caf7bf73ff105c658a42e12eaa33f58af`](https://sepolia.etherscan.io/address/0x9163ad7caf7bf73ff105c658a42e12eaa33f58af) |
| ArenaCoin | [`0x9E8119a40eA0c4A925348f0C998333953FB73D0C`](https://sepolia.etherscan.io/address/0x9E8119a40eA0c4A925348f0C998333953FB73D0C) |

### Tech Stack

**Languages**: Solidity · Go · TypeScript

**Frameworks/Libraries**: Next.js 14 · Tailwind CSS · shadcn/ui · Gin · wagmi · RainbowKit · viem · Foundry · OpenZeppelin · gorilla/websocket

**AI**: Qwen API (Tongyi Qianwen)

**Database**: SQLite

**Blockchain**: Sepolia Testnet · Infura RPC

**Deployment**: Vercel (frontend)

### Quick Start

```bash
# 1. Clone
git clone https://github.com/your-username/agent-arena.git
cd agent-arena

# 2. Contracts: build & test
cd contracts
forge install
forge build
forge test

# 3. Start backend
cd ../backend
# Configure .env (RPC_URL, PRIVATE_KEY, ARENA_ADDRESS, etc.)
go run cmd/server/main.go
# → http://localhost:8080

# 4. Start frontend
cd ../frontend
npm install
# Configure .env.local (NEXT_PUBLIC_API_URL, etc.)
npm run dev
# → http://localhost:3000
```

### Environment Variables

**backend/.env**:
```
RPC_URL=https://sepolia.infura.io/v3/YOUR_KEY
PRIVATE_KEY=0x...
ARENA_ADDRESS=0x...
USDC_ADDRESS=0x...
AC_ADDRESS=0x...
QWEN_API_KEY=sk-...
```

**frontend/.env.local**:
```
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080
NEXT_PUBLIC_ARENA_ADDRESS=0x...
NEXT_PUBLIC_POOL_ADDRESS=0x...
NEXT_PUBLIC_USDC_ADDRESS=0x...
```

---

## 🏗️ Built With

`Solidity` `Go` `TypeScript` `Next.js` `Foundry` `Sepolia` `ERC20` `Tailwind CSS` `shadcn/ui` `wagmi` `Gin` `SQLite` `Infura` `Qwen` `WebSocket` `OpenZeppelin` `RainbowKit` `viem`

## 📄 License

MIT
