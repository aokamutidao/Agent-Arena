# S3 - 智能合约规格 (Contract Specification)

## Agent Arena Smart Contracts

---

## 1. 合约架构

```
┌─────────────────────────────────────────────────────────────┐
│                    AgentArena.sol                           │
│                   (主合约 / 入口)                            │
│                                                             │
│  职责：                                                     │
│  - 创建对局                                                  │
│  - 协调 Betting + Strategy + Registry                       │
│  - 提供统一的用户接口                                        │
└─────────────────────────────────────────────────────────────┘
        │                    │                    │
        ↓                    ↓                    ↓
┌───────────────┐  ┌─────────────────┐  ┌──────────────────┐
│ BettingPool   │  │ StrategyVoting  │  │ GameRegistry     │
│               │  │                 │  │                  │
│ - 接受 USDC   │  │ - 记录策略投票   │  │ - 记录对局结果   │
│ - 锁定下注    │  │ - 计算权重分布   │  │ - 记录 Agent 胜率│
│ - 结算分配    │  │ - 一次性锁定     │  │ - 批量提交 hash  │
└───────────────┘  └─────────────────┘  └──────────────────┘
        │
        ↓
┌───────────────┐
│ IERC20 (USDC) │
│ (外部合约)     │
└───────────────┘
```

---

## 2. 合约详解

### 2.1 BettingPool.sol

**职责**：管理下注资金，锁定到结算，自动分配奖金。

#### 数据结构

```solidity
struct Game {
    uint256 gameId;
    address agentRed;       // Red Agent ID (链上标识)
    address agentBlue;      // Blue Agent ID
    uint256 totalBetRed;    // Red 方总下注 (USDC, 6 decimals)
    uint256 totalBetBlue;   // Blue 方总下注
    uint256 bettingDeadline; // 下注截止时间
    GameStatus status;      // Open | Locked | Finished
    address winner;         // 获胜方 (Red Agent 或 Blue Agent)
    uint256 protocolFee;    // 协议费 (5%)
}

struct Bet {
    uint256 gameId;
    address user;
    Side side;              // Red | Blue
    uint256 amount;
    bool claimed;
}

enum GameStatus { Open, Locked, Finished }
enum Side { None, Red, Blue }
```

#### 接口定义

```solidity
// ============ 下注 ============

/// @notice 用户对某场对局下注
/// @param gameId 对局 ID
/// @param side 押哪方 (Red=true, Blue=false)
/// @param amount USDC 数量 (6 decimals)
function placeBet(
    uint256 gameId,
    bool side,          // true=Red, false=Blue
    uint256 amount
) external;

// 流程：
// 1. require(game.status == Open)
// 2. require(block.timestamp < game.bettingDeadline)
// 3. USDC.transferFrom(msg.sender, address(this), amount)
// 4. 更新 totalBetRed 或 totalBetBlue
// 5. 记录 bet 到 bets[gameId][msg.sender]
// 6. emit BetPlaced(gameId, msg.sender, side, amount)


// ============ 锁定 ============

/// @notice 锁定对局（下注截止后调用）
/// @dev 只能由 owner (后端) 调用
function lockBetting(uint256 gameId) external onlyOwner;

// 流程：
// 1. require(game.status == Open)
// 2. game.status = Locked
// 3. emit BettingLocked(gameId)


// ============ 结算 ============

/// @notice 提交对局结果并结算
/// @dev 只能由 owner (后端) 调用
/// @param gameId 对局 ID
/// @param redWins Red 是否获胜
function settle(uint256 gameId, bool redWins) external onlyOwner;

// 流程：
// 1. require(game.status == Locked)
// 2. game.winner = redWins ? Red : Blue
// 3. game.status = Finished
// 4. 计算协议费: fee = (totalBetRed + totalBetBlue) * 5 / 100
// 5. 转移协议费到 protocolTreasury
// 6. emit GameSettled(gameId, redWins, winnerPool, loserPool)


// ============ 提取 ============

/// @notice 赢家提取奖金
/// @param gameId 对局 ID
function claim(uint256 gameId) external;

// 流程：
// 1. require(game.status == Finished)
// 2. require(bet.side == game.winner)
// 3. require(!bet.claimed)
// 4. reward = bet.amount * winnerPool / (totalPool - fee)
// 5. bet.claimed = true
// 6. USDC.transfer(msg.sender, reward)
// 7. emit RewardClaimed(gameId, msg.sender, reward)


// ============ 查询 ============

/// @notice 查询用户下注信息
function getBet(uint256 gameId, address user) 
    external view returns (Side side, uint256 amount, bool claimed);

/// @notice 查询对局信息
function getGame(uint256 gameId) 
    external view returns (Game memory);

/// @notice 查询赢家可提取金额
function getReward(uint256 gameId, address user) 
    external view returns (uint256);

/// @notice 查询当前赔率
function getOdds(uint256 gameId) 
    external view returns (uint256 oddsRed, uint256 oddsBlue);
```

#### 事件定义

```solidity
event BetPlaced(
    uint256 indexed gameId,
    address indexed user,
    bool side,           // true=Red, false=Blue
    uint256 amount
);

event BettingLocked(
    uint256 indexed gameId,
    uint256 totalBetRed,
    uint256 totalBetBlue
);

event GameSettled(
    uint256 indexed gameId,
    bool redWins,
    uint256 totalPool,
    uint256 protocolFee
);

event RewardClaimed(
    uint256 indexed gameId,
    address indexed user,
    uint256 reward
);
```

---

### 2.2 StrategyVoting.sol

**职责**：记录策略投票，计算权重分布，一次性锁定。

#### 数据结构

```solidity
struct StrategyVoteRecord {
    uint256 aggressive;  // 激进投票总额 (USDC)
    uint256 defensive;   // 稳健投票总额
    uint256 tricky;      // 诡道投票总额
    bool locked;         // 是否已锁定
}

enum Strategy { Aggressive, Defensive, Tricky }
```

#### 接口定义

```solidity
// ============ 投票 ============

/// @notice 用户投票选择策略（与下注同时调用）
/// @param gameId 对局 ID
/// @param strategy 选择的策略
/// @param amount 关联的下注金额（用于计算权重）
function vote(
    uint256 gameId,
    Strategy strategy,
    uint256 amount
) external;

// 流程：
// 1. require(!voteRecord[gameId].locked)
// 2. require(block.timestamp < bettingDeadline)
// 3. 更新对应策略的总额
// 4. emit StrategyVoted(gameId, msg.sender, strategy, amount)


// ============ 锁定 ============

/// @notice 锁定投票（开局时调用）
function lockVotes(uint256 gameId) external onlyOwner;

// 流程：
// 1. voteRecord[gameId].locked = true
// 2. emit VotesLocked(gameId, aggressive, defensive, tricky)


// ============ 查询 ============

/// @notice 查询策略分布（给 AI Agent 用）
function getStrategyWeights(uint256 gameId) 
    external view returns (
        uint256 aggressiveWeight,  // 百分比 0-100
        uint256 defensiveWeight,
        uint256 trickyWeight
    );

// 计算：
// total = aggressive + defensive + tricky
// if total == 0: return (33, 33, 34)  // 默认均分
// return (aggressive*100/total, defensive*100/total, tricky*100/total)

/// @notice 查询用户的投票
function getUserVote(uint256 gameId, address user) 
    external view returns (Strategy strategy, uint256 amount);
```

#### 事件定义

```solidity
event StrategyVoted(
    uint256 indexed gameId,
    address indexed user,
    Strategy strategy,
    uint256 amount
);

event VotesLocked(
    uint256 indexed gameId,
    uint256 aggressive,
    uint256 defensive,
    uint256 tricky
);
```

---

### 2.3 GameRegistry.sol

**职责**：记录对局元数据、Agent 信息、动作 hash（防篡改）、胜率统计。

#### 数据结构

```solidity
struct GameRecord {
    uint256 gameId;
    address agentRed;
    address agentBlue;
    uint256 startTimestamp;
    uint256 endTimestamp;
    bool redWins;
    bool exists;
}

struct AgentInfo {
    string name;
    string personality;     // "Berserker", "Tactician", "Trickster", "Defender"
    uint256 wins;
    uint256 losses;
    bool exists;
}
```

#### 接口定义

```solidity
// ============ 对局管理 ============

/// @notice 创建新对局
/// @param agentRed Red 方 Agent ID
/// @param agentBlue Blue 方 Agent ID
/// @param bettingDuration 下注持续时间（秒）
function createGame(
    address agentRed,
    address agentBlue,
    uint256 bettingDuration
) external onlyOwner returns (uint256 gameId);

// 流程：
// 1. gameId = nextGameId++
// 2. 记录 gameRecord[gameId]
// 3. emit GameCreated(gameId, agentRed, agentBlue, bettingDeadline)


/// @notice 提交对局结果
function submitResult(
    uint256 gameId,
    bool redWins,
    bytes32 actionsHash    // 所有回合动作的 hash
) external onlyOwner;

// 流程：
// 1. require(gameRecord[gameId].exists)
// 2. gameRecord[gameId].redWins = redWins
// 3. gameRecord[gameId].endTimestamp = block.timestamp
// 4. actionsHashes[gameId] = actionsHash
// 5. 更新 Agent 胜率
// 6. emit GameResultSubmitted(gameId, redWins, actionsHash)


// ============ Agent 管理 ============

/// @notice 注册 Agent
function registerAgent(
    address agentId,
    string calldata name,
    string calldata personality
) external onlyOwner;

// 流程：
// 1. agents[agentId] = AgentInfo(name, personality, 0, 0, true)
// 2. emit AgentRegistered(agentId, name, personality)


// ============ 查询 ============

/// @notice 查询 Agent 信息
function getAgent(address agentId) 
    external view returns (AgentInfo memory);

/// @notice 查询 Agent 胜率
function getWinRate(address agentId) 
    external view returns (uint256 wins, uint256 losses, uint256 winRate);

// winRate = wins * 100 / (wins + losses)  // 百分比

/// @notice 查询对局动作 hash（用于验证）
function getActionsHash(uint256 gameId) 
    external view returns (bytes32);

/// @notice 查询对局记录
function getGameRecord(uint256 gameId) 
    external view returns (GameRecord memory);
```

#### 事件定义

```solidity
event GameCreated(
    uint256 indexed gameId,
    address agentRed,
    address agentBlue,
    uint256 bettingDeadline
);

event GameResultSubmitted(
    uint256 indexed gameId,
    bool redWins,
    bytes32 actionsHash
);

event AgentRegistered(
    address indexed agentId,
    string name,
    string personality
);

event AgentWinUpdated(
    address indexed agentId,
    uint256 wins,
    uint256 losses
);
```

---

### 2.4 AgentArena.sol（主合约）

**职责**：统一入口，协调三个子合约。

```solidity
contract AgentArena {
    BettingPool public bettingPool;
    StrategyVoting public strategyVoting;
    GameRegistry public gameRegistry;

    /// @notice 用户下注 + 投票（一次性调用）
    function betAndVote(
        uint256 gameId,
        bool side,              // true=Red, false=Blue
        uint256 amount,
        StrategyVoting.Strategy strategy
    ) external {
        // 1. 先下注
        bettingPool.placeBet(gameId, side, amount);
        // 2. 再投票
        strategyVoting.vote(gameId, strategy, amount);
    }

    /// @notice 后端调用：开局（锁定下注 + 锁定投票）
    function startGame(uint256 gameId) external onlyOwner {
        bettingPool.lockBetting(gameId);
        strategyVoting.lockVotes(gameId);
    }

    /// @notice 后端调用：结算（提交结果 + 结算奖金）
    function finishGame(
        uint256 gameId,
        bool redWins,
        bytes32 actionsHash
    ) external onlyOwner {
        gameRegistry.submitResult(gameId, redWins, actionsHash);
        bettingPool.settle(gameId, redWins);
    }
}
```

---

## 3. USDC 交互流程

### 3.1 用户下注流程

```
用户钱包                    AgentArena                  USDC
   │                           │                          │
   │── 1. approve(amount) ────→│                          │
   │                           │                          │
   │── 2. betAndVote() ───────→│                          │
   │                           │── 3. transferFrom() ────→│
   │                           │   (锁定 USDC 到合约)      │
   │                           │                          │
   │                           │── 4. 记录 bet + vote ───→│
   │                           │                          │
   │←── 5. emit events ───────│                          │
```

### 3.2 结算流程

```
Go 后端                     AgentArena                  USDC
   │                           │                          │
   │── 1. finishGame() ───────→│                          │
   │                           │── 2. 计算协议费 (5%) ───→│
   │                           │── 3. transfer(fee) ─────→│ (protocolTreasury)
   │                           │── 4. 记录 winner ───────→│
   │                           │                          │
   │←── 5. emit events ───────│                          │
   │                           │                          │

赢家用户                    AgentArena                  USDC
   │                           │                          │
   │── 6. claim() ───────────→│                          │
   │                           │── 7. transfer(reward) ──→│
   │←── 8. USDC 到账 ─────────│                          │
```

---

## 4. 状态机

### 4.1 对局状态（链上视角）

```
  createGame()
      │
      ▼
  ┌─────────┐
  │  Open   │ ← 下注 + 投票开放
  └────┬────┘
       │ startGame()
       ▼
  ┌─────────┐
  │ Locked  │ ← 下注截止，投票锁定，对战进行中（链下）
  └────┬────┘
       │ finishGame()
       ▼
  ┌─────────┐
  │Finished │ ← 可提取奖金
  └─────────┘
```

### 4.2 约束规则

```
- Open → Locked: 只有 owner 可调用
- Locked → Finished: 只有 owner 可调用
- Finished 是终态，不可回退
- claim() 只能在 Finished 状态调用
- settle() 只能在 Locked 状态调用
```

---

## 5. 权限控制

| 函数 | 调用者 | 说明 |
|------|--------|------|
| `placeBet()` | 任何用户 | 公开 |
| `vote()` | 任何用户 | 公开 |
| `claim()` | 任何用户 | 赢家公开 |
| `createGame()` | onlyOwner | Go 后端 |
| `startGame()` | onlyOwner | Go 后端 |
| `finishGame()` | onlyOwner | Go 后端 |
| `registerAgent()` | onlyOwner | Go 后端 |
| `lockBetting()` | onlyOwner | Go 后端 |
| `lockVotes()` | onlyOwner | Go 后端 |
| `settle()` | onlyOwner | Go 后端 |

**Owner**：Go 后端的运营钱包地址（部署时设定）。

---

## 6. Gas 优化

| 优化点 | 方案 |
|--------|------|
| USDC 精度 | 使用 6 decimals（USDC 标准） |
| 数组存储 | 不存储完整 bet 列表，只存总额 + mapping |
| 批量查询 | 提供 `getGame()` 一次性返回所有信息 |
| 事件索引 | gameId 做 indexed，方便前端过滤 |
| 动作 hash | 只存一个 bytes32，不存每回合 |

**预估 Gas**：
- `placeBet()`: ~80k Gas
- `vote()`: ~50k Gas
- `betAndVote()`: ~120k Gas (合并调用)
- `claim()`: ~60k Gas
- `finishGame()`: ~100k Gas

---

## 7. 安全考虑

| 风险 | 应对 |
|------|------|
| 重入攻击 | 使用 OpenZeppelin ReentrancyGuard |
| USDC 精度问题 | 全程使用 uint256，6 decimals |
| Owner 作恶 | MVP 阶段 owner 是可信后端；未来可改多签 |
| 结算不公 | 结算逻辑简单透明，链上可验证 |
| USDC 合约兼容性 | 使用 SafeERC20 处理非标准返回值 |

---

## 8. 部署配置

```toml
# foundry.toml
[profile.default]
src = "src"
out = "out"
libs = ["lib"]
solc_version = "0.8.24"
optimizer = true
optimizer_runs = 200

[profile.sepolia]
rpc_url = "${SEPOLIA_RPC_URL}"
chain_id = 11155111
```

```
环境变量：
  SEPOLIA_RPC_URL=...
  PRIVATE_KEY=...         # 部署者私钥
  USDC_ADDRESS=...        # Sepolia USDC 合约地址
  PROTOCOL_TREASURY=...   # 协议费收款地址
```

---

**文档版本**：v1.0
**最后更新**：2026-08-02
**状态**：待评审
