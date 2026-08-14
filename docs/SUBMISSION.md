# Agent Arena
## AI Agent 格子竞技场 — 链上下注 × 策略投票 × 实时对战

---

## 📖 项目概述

Agent Arena 是一个全链上 AI Agent 竞技平台。两个由 LLM 驱动的 AI Agent 在 10×10 格子地图上实时对战，观众可以：

- 🎰 **用 USDC 下注** — 选择红/蓝方，赔率由下注池动态决定
- 🗳️ **策略投票** — 用下注金额加权投票，影响 AI 的战术决策（激进/稳健/诡道）
- 💰 **领取奖金** — 游戏结束后通过智能合约自动结算，赢家按比例分配奖池

---

## 🏗️ 系统架构

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

---

## 🎯 核心特性

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

---

## 🔗 智能合约 (Sepolia Testnet)

| 合约 | 地址 |
|------|------|
| AgentArena | `0x0bccc77f672f22f843b83ba65297e6ee693c68c7` |
| BettingPool | `0xa92f3c9bad4330eb0d3d42e9e6073b577b7e9782` |
| StrategyVoting | `0x6e1689f5ef9db00de4c40d0bbead951701c8c315` |
| GameRegistry | `0xfbcb875bc9f7bc543dc8709d8b7815b9b5df7eb6` |
| MockUSDC | `0x9163ad7caf7bf73ff105c658a42e12eaa33f58af` |
| ArenaCoin | `0x9E8119a40eA0c4A925348f0C998333953FB73D0C` |

---

## 🛠️ 技术栈

**语言**: Solidity · Go · TypeScript

**框架/库**: Next.js 14 · Tailwind CSS · shadcn/ui · Gin · wagmi · RainbowKit · viem · Foundry · OpenZeppelin · gorilla/websocket

**AI**: Qwen API (通义千问)

**数据库**: SQLite

**区块链**: Sepolia Testnet · Infura RPC

**部署**: Vercel (前端) · Railway (后端)

---

## 📝 项目链接

- **源代码**: https://github.com/aokamutidao/Agent-Arena
- **在线演示**: https://agent-arena-gold.vercel.app

---

## 📸 界面截图

（请在此处插入截图，可以从 http://localhost:3000 截取以下页面：）

1. **首页** - 显示当前进行的游戏列表
2. **游戏观战页面** - 10×10 格子地图，红蓝双方 Agent 实时对战
3. **下注面板** - 选择红/蓝方，显示动态赔率
4. **策略投票** - 红蓝分边投票面板，显示当前策略权重
5. **钱包页面** - 连接钱包，显示 USDC 余额和可领取奖金
6. **历史战绩** - 过往游戏记录和结果

---

## 🎬 功能演示

### 游戏流程
1. 系统创建游戏，两个 AI Agent 准备对战
2. 观众在 30 秒下注窗口内用 USDC 下注
3. 下注同时投票选择 AI 策略（激进/稳健/诡道）
4. 游戏开始，AI 根据策略权重自主决策
5. WebSocket 实时推送每一步动作
6. 前端实时渲染格子地图和战斗动画
7. 游戏结束，赢家通过智能合约领取奖金

### AI Agent 决策
- 每回合根据策略权重选择行动（移动/攻击/技能）
- 考虑血量、位置、敌方状态、战争迷雾
- 加时赛机制确保必有胜负

---

## 🚀 快速开始

```bash
# 1. 克隆项目
git clone https://github.com/aokamutidao/Agent-Arena.git
cd Agent-Arena

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

---

## 📄 License

MIT
