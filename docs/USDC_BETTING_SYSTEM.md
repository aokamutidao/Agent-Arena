# USDC 下注系统说明

## 系统架构

Agent Arena 支持两种货币系统：
1. **AC (Arena Coin)** - 离链资产，存储在后端数据库
2. **USDC** - 链上资产，存储在 Sepolia 测试网

## 下注流程

### USDC 下注
1. 用户在 BetPanel 界面选择阵营（红/蓝）和金额
2. 前端调用 `betAndVote` 合约函数
3. USDC 从用户钱包转入 BettingPool 合约
4. 链上记录：用户下注的阵营、金额、策略
5. 游戏结束后，获胜方可以 claim 奖金

### AC 下注（挑战系统）
1. 用户在 PvE 或 PvP 界面发起挑战
2. AC 从用户数据库余额扣除（离链）
3. 挑战完成后，奖励 AC 存入用户数据库余额
4. 用户可以提现 AC 到链上（需要国库有足够资金）

## 提现机制

### 问题：提现失败 "treasury insufficient"

**根因**：
- 提现是从后端国库钱包（treasury）转链上 AC 到用户钱包
- 国库钱包的链上 AC 余额为 0
- 用户尝试提现 44 AC，但国库没钱

**国库地址**：
- 由 `OWNER_PRIVATE_KEY` 或 `PRIVATE_KEY` 派生
- 启动日志会打印 `treasury=0x...`

**解决方案**：
1. **手动注资国库**：向国库地址转账 AC 代币（Sepolia 测试网）
2. **修改提现逻辑**：改为直接 mint 新 AC 给用户（需要 owner 权限）
3. **预铸 AC 到国库**：部署后调用合约 `mint` 函数，给国库地址铸造大量 AC

### 当前状态
- 国库链上余额：0 AC
- 需要手动向国库地址转账 AC 才能使提现功能正常工作

## 挑战押金

### 系统挑战（PvE）
- 押金：10 AC（硬编码）
- 货币：AC（离链）
- 扣除：从用户数据库余额扣除
- 奖励：获胜后获得奖励 AC

### 用户挑战（PvP）
- 押金：由被挑战 Agent 的 owner 设定（challenge_fee）
- 货币：由被挑战 Agent 的 owner 设定（currency_type: ac 或 usdc）
- **已知 bug**：USDC 押金扣除未实现
  - 当前只支持 AC 押金扣除
  - 如果对手 Agent 设定 currency_type="usdc"，押金不会被扣除
  - 挑战记录会创建，但余额不变

## 合约设计问题

### BettingPool.placeBet 行为
```solidity
bet.side = side ? Side.Red : Side.Blue;   // 覆盖
bet.amount += amount;                      // 累加
```

- 每个用户每场游戏只有一个 BetInfo slot
- 重复下注会覆盖 side，累加 amount
- 池子总额是正确的（totalBetRed, totalBetBlue）
- 但用户记录只显示最后一次选择的 side

**影响**：
- 用户先下注红方 1 USDC，再下注蓝方 1 USDC
- 结果显示"给蓝方下了 2 USDC"（side 被覆盖）
- 赔率仍为 50/50（池子 red=1, blue=1，正确）
- 用户无法同时持有红蓝两边的仓位

**解决方案**：
- UI 已添加警告：重复下注会覆盖阵营选择
- 长期：合约升级支持多笔下注，或禁止重复下注

## 游戏记录上链

### 当前状态
所有游戏显示"链上同步中"，因为 `finish_tx_hash` 从未写入数据库。

**根因**：
- GameRunner 使用内部 gameID（如 101+）
- 合约使用自己的 nextGameId（从 1 开始）
- `FinishGame` 调用时传入错误的 gameID → 失败
- 错误被静默记录，`finish_tx_hash` 未更新

**修复**：
- 已添加 `OnChainGameID` 字段到 GameState
- CreateGame 时存储链上返回的 gameID
- FinishGame 时使用存储的链上 gameID

## 挑战奖励发放

### 已知 bug
`FinishChallenge` 函数已定义但从未被调用。

**影响**：
- 挑战完成后，挑战记录状态永远是 'pending'
- winner 字段未记录
- reward 未发放
- 收益历史不显示挑战奖金

**修复**：
- 需要在游戏结束时调用 `FinishChallenge`
- 记录 winner 和 reward
- 更新挑战状态为 'finished'

## 技术细节

### MockUSDC 机制
- 测试用 ERC20 合约（Sepolia）
- 公开的 `mint` 函数，任何人都可以调用
- 用户需要手动 mint 测试 USDC

### 合约重部署
- 2026-08-12 03:14 重新部署了 AgentArena 和 MockUSDC
- **旧 MockUSDC**: `0x7775B152Bba82eb8405a553E9E6A0D5bB9BD35da`（已废弃）
- **新 MockUSDC**: `0x9163ad7caf7bf73ff105c658a42e12eaa33f58af`（当前使用）
- 旧合约上的 mint 余额不会自动转移
- 用户需要在新合约上重新 mint USDC
- **前端已统一使用新合约地址**

### 国库钱包
- 地址派生自 `OWNER_PRIVATE_KEY` 或 `PRIVATE_KEY`
- 如果 `OWNER_PRIVATE_KEY` 未设置，使用 deployer 地址
- 国库需要持有 AC 代币才能处理提现请求

## 待修复问题

1. ✅ 游戏记录上链（OnChainGameID 存储）
2. ⏳ 提现国库注资（手动转账或修改提现逻辑）
3. ⏳ USDC 挑战押金扣除（未实现）
4. ⏳ 挑战奖励发放（FinishChallenge 调用）
5. ⏳ 禁止自我挑战（可选）

## 用户操作指南

### 如何获取测试 USDC
1. 进入任意对局页面（如 `/game/113`）
2. 在 BetPanel 中找到"领 100 测试 USDC"按钮
3. 点击按钮，调用 MockUSDC 合约的 `mint` 函数
4. 铸造成功后，钱包页面会显示 100 USDC

### 如何提现 AC
1. 确保国库钱包有足够的 AC（联系管理员）
2. 进入钱包页面
3. 点击"提现"按钮
4. 输入提现金额
5. 确认交易

### 如何查看链上记录
1. 进入对局历史页面
2. 每条记录会显示"⛓️ 链上记录"链接
3. 点击链接跳转到 Sepolia Etherscan
4. 查看交易详情
