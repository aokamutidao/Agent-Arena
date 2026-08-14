# V2 迭代计划 — 用户系统 + 自定义 Agent + PVE 模式

## 版本概述

基于 V1 基线版本（当前）进行大版本迭代，新增：
1. **用户系统**：登录、个人中心、收益查看、认领
2. **自定义 Agent**：用户创建 Agent、定义提示词、提供 API 接口
3. **PVE 模式**：系统擂台、用户擂台、挑战机制
4. **双币经济**：Arena Coins（游戏币）+ USDC（真实代币）

**基线版本**：V1（迭代 014 完成）— 可回退

---

## 核心设计决策（已确认）

### 1. 货币体系：双币制

| 货币 | 获取方式 | 用途 | 区块链 |
|------|---------|------|--------|
| **Arena Coins (AC)** | 免费领取、PVE 胜利奖励 | PVE 挑战、休闲对战、测试 | **链上（ERC20 on Sepolia）** |
| **USDC** | 钱包充值 | 正式 PvP、锦标赛、高风险下注 | 链上（Sepolia） |

**✅ 已确认：AC 上链（ERC20）**
- 发行 ArenaCoin 代币到 Sepolia
- 用户钱包持有，可转账、交易
- 需要 gas，但完全去中心化

### 2. 用户认证：钱包签名

- 用户用 MetaMask 签名登录（不需要密码）
- 签名消息：`"Login to Agent Arena: {timestamp}"`
- 后端验证签名 → 生成 JWT token
- 用户信息：address, username, AC balance, agent list, bet history

**✅ 已确认：钱包签名认证**

### 3. Agent 系统：用户创建 + 平台预设

| Agent 类型 | 来源 | 提示词 | API | 货币 |
|-----------|------|--------|-----|------|
| **系统 Agent** | 平台预设 | 固定（4 种性格） | 后端 Qwen | AC 奖励 |
| **用户 Agent** | 用户创建 | 用户自定义（≤500 字） | 可选：用户提供的 API endpoint | AC 或 USDC |

- 用户创建 Agent 时选择：
  - 名称、头像、性格描述（≤500 字）
  - 是否提供外部 API（可选）
  - 挑战货币类型（AC 或 USDC）
- 如果用户提供 API，后端调用该 API 而非 Qwen
- Agent 可以上架"擂台"供他人挑战

### 4. 对战模式

**✅ 已确认：MVP 先做 PVE**

#### PVE 模式（AC）— Sprint 6 优先实现
- **系统擂台**：用户挑战平台预设 Agent（Berserker/Tactician/Trickster/Defender）
  - 胜利奖励：100 AC
  - 失败无惩罚
  - 可以重复挑战
- **用户擂台**：用户挑战其他用户的 Agent
  - 擂主设定赌注（AC）
  - 挑战者支付相同赌注
  - 胜利者获得池子（扣除 5% 手续费）

#### PvP 模式（USDC）— Sprint 7 后续实现
- **匹配对战**：用户用自己的 Agent 匹配其他用户
  - 双方各出 USDC
  - 胜利者获得池子（扣除 5% 手续费）
- **锦标赛**：多人淘汰赛
  - 报名费 USDC
  - 冠军获得奖池

#### 观战下注
- 所有对局（PVE/PvP）都开放观战和下注
- 下注货币与对局货币一致（AC 局用 AC 下注，USDC 局用 USDC 下注）
- 策略投票机制保持不变

---

## 技术架构变更

### 后端新增模块

```
backend/internal/
├── user/              # 用户系统
│   ├── user.go        # User struct, CRUD
│   ├── auth.go        # 钱包签名验证, JWT 生成
│   └── balance.go     # AC 余额管理
├── agent/             # Agent 系统（扩展现有）
│   ├── custom.go      # 用户自定义 Agent
│   ├── api_client.go  # 外部 API 调用
│   └── marketplace.go # Agent 上架/下架
├── match/             # 对战匹配
│   ├── pve.go         # PVE 挑战逻辑
│   ├── pvp.go         # PvP 匹配逻辑
│   └── tournament.go  # 锦标赛逻辑
└── currency/          # 货币系统
    ├── arena_coin.go  # AC 管理（链下）
    └── usdc.go        # USDC 管理（链上）
```

### 数据库变更

**新增表**：
- `users`: address, username, ac_balance, created_at
- `custom_agents`: owner, name, personality, api_endpoint, challenge_fee, currency_type
- `matches`: match_type (pve/pvp/tournament), challenger, defender, stake, currency, winner
- `transactions`: user_id, amount, currency, type (reward/bet/claim/transfer), match_id

### 合约变更

**新增合约**：
- `ArenaCoin.sol`: ERC20 游戏币（可选上链，或纯链下）
- `TournamentPool.sol`: 锦标赛奖池管理
- 现有 `BettingPool.sol` 扩展：支持多币种（AC 或 USDC）

### 前端新增页面

```
/                   → 首页：对局列表（PVE/PvP/锦标赛）
/login              → 登录页（钱包签名）
/profile            → 个人中心
  /profile/agents   → 我的 Agent 列表
  /profile/bets     → 我的下注历史
  /profile/rewards  → 我的收益
/arena              → 擂台列表
  /arena/pve        → PVE 挑战（系统 Agent）
  /arena/user       → 用户擂台
  /arena/tournament → 锦标赛
/agent/create       → 创建 Agent
/agent/[id]         → Agent 详情
/game/[id]          → 对局详情（现有）
```

---

## 迭代计划

### Sprint 5: 用户系统 + 基础 Agent 管理（2-3 周）

| # | 迭代 | 任务 | 状态 | 依赖 |
|---|------|------|------|------|
| 016 | 用户认证 | 钱包签名登录 + JWT + 用户表 | ⏳ 待开始 | 无 |
| 017 | 个人中心 | 用户信息页 + AC 余额 + 下注历史 | ⏳ 待开始 | 016 |
| 018 | 自定义 Agent | Agent CRUD + 提示词管理 + API 接口注册 | ⏳ 待开始 | 016 |
| 019 | Agent 市场 | Agent 列表页 + 上架/下架 | ⏳ 待开始 | 018 |

### Sprint 6: PVE 模式（2 周）

| # | 迭代 | 任务 | 状态 | 依赖 |
|---|------|------|------|------|
| 020 | AC 系统 | Arena Coin 管理（链下）+ 水龙头 | ⏳ 待开始 | 016 |
| 021 | 系统擂台 | PVE 挑战系统 Agent + 奖励机制 | ⏳ 待开始 | 018, 020 |
| 022 | 用户擂台 | 用户 Agent 上架 + 挑战逻辑 | ⏳ 待开始 | 019, 020 |

### Sprint 7: PvP + 锦标赛（2-3 周）

| # | 迭代 | 任务 | 状态 | 依赖 |
|---|------|------|------|------|
| 023 | PvP 匹配 | 匹配队列 + USDC 对战 | ⏳ 待开始 | 016, 020 |
| 024 | 锦标赛 | 多人淘汰赛 + 报名 + 奖池分配 | ⏳ 待开始 | 023 |
| 025 | 双币下注 | BettingPool 支持 AC/USDC 双币种 | ⏳ 待开始 | 020 |

### Sprint 8: 排行榜 + 社交（1-2 周）

| # | 迭代 | 任务 | 状态 | 依赖 |
|---|------|------|------|------|
| 026 | 排行榜 | 玩家排行（按收益）+ Agent 排行（按胜率） | ⏳ 待开始 | 025 |
| 027 | 社交功能 | Agent 评论 + 关注 + 分享 | ⏳ 待开始 | 019 |

---

## 关键问题与决策点

### 1. Arena Coin 是否上链？

**选项 A：纯链下（推荐）**
- AC 存在后端数据库，不发行 ERC20
- 优点：简单、快速、无 gas 费
- 缺点：不够 Web3

**选项 B：上链（ERC20）**
- 发行 ArenaCoin 代币，用户钱包持有
- 优点：完全去中心化、可交易
- 缺点：复杂、需要 gas、用户门槛高

**建议**：MVP 用链下，后续版本再考虑上链。

### 2. 外部 API 调用如何处理？

用户提供的 Agent API 可能：
- 响应慢（>5s）
- 返回错误
- 甚至宕机

**方案**：
- 设置 5s 超时
- 失败则 fallback 到 WAIT
- 记录 API 可用性，影响 Agent 信誉分

### 3. 防作弊机制

**问题**：用户可以创建"必胜"Agent 然后自己挑战自己刷 AC

**方案**：
- 每日 AC 奖励上限（如 1000 AC/天）
- 同 IP/同钱包挑战自己的 Agent 不奖励
- 高收益触发人工审核

---

## 验收标准

### Sprint 5 验收
- [ ] 用户可以用 MetaMask 签名登录
- [ ] 登录后可以看到 AC 余额和下注历史
- [ ] 用户可以创建自定义 Agent（设置名称、提示词）
- [ ] 用户可以选择是否提供外部 API endpoint
- [ ] Agent 可以上架到市场供他人查看

### Sprint 6 验收
- [ ] 用户可以免费领取 AC（每日上限）
- [ ] 用户可以挑战系统 Agent（PVE）
- [ ] 挑战胜利获得 100 AC 奖励
- [ ] 用户可以将自己的 Agent 上架为擂台
- [ ] 其他用户可以挑战用户擂台（消耗 AC）
- [ ] 挑战胜利获得赌注池（扣除 5% 手续费）

### Sprint 7 验收
- [ ] 用户可以用 USDC 进行 PvP 匹配
- [ ] 匹配成功后自动开始对战
- [ ] 胜利者获得 USDC 奖池
- [ ] 锦标赛可以报名、淘汰、分配奖池
- [ ] 下注支持 AC 和 USDC 两种货币

---

## 风险与缓解

| 风险 | 影响 | 缓解措施 |
|------|------|---------|
| 外部 API 不稳定 | 对战中断 | 超时 fallback 到 WAIT + 信誉分系统 |
| AC 刷币 | 经济系统崩溃 | 每日上限 + 防作弊检测 |
| 匹配不到人 | 用户体验差 | PVE 保底 + 机器人填充 |
| 合约漏洞 | 资金损失 | 代码审计 + 小额测试 |

---

## 下一步行动

1. ~~**确认货币方案**~~：✅ AC 上链（ERC20）
2. ~~**确认 MVP 范围**~~：✅ 先 PVE
3. ~~**确认认证方式**~~：✅ 钱包签名
4. **建立 V1 基线**：git tag v1.0-baseline
5. **部署 ArenaCoin 合约**：ERC20 代币到 Sepolia
6. **开始 Sprint 5**：迭代 016 用户认证

---

**文档版本**：v1.1（已确认）
**创建日期**：2026-08-10
**状态**：已确认，开始实施
