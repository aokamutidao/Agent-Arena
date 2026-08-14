# 编码规范

## Go 代码规范

### 基本规则
- Go 版本: 1.21+
- 使用 `go fmt` 格式化所有代码
- 使用 `go vet` 检查常见错误
- 行长度不超过 100 字符

### 命名
```go
// 包名: 小写，单词
package engine
package agent

// 导出: 首字母大写
type GameState struct { ... }
func NewGame() *Game { ... }

// 非导出: 首字母小写
type agentState struct { ... }
func calculateScore() int { ... }

// 接口: 以 er 结尾或用行为命名
type LLMClient interface { ... }
type Store interface { ... }
```

### 错误处理
```go
// 总是检查错误
result, err := doSomething()
if err != nil {
    return fmt.Errorf("doSomething failed: %w", err)
}

// 不要吞掉错误
// ❌ result, _ := doSomething()
// ✅ result, err := doSomething(); if err != nil { ... }
```

### 项目结构
```
backend/
├── cmd/server/main.go     # 入口
├── internal/              # 内部包（不导出）
│   ├── engine/            # 游戏引擎
│   ├── agent/             # AI Agent
│   ├── blockchain/        # 合约交互
│   ├── api/               # HTTP/WS 服务
│   └── config/            # 配置
└── go.mod
```

---

## Solidity 代码规范

### 基本规则
- Solidity 版本: 0.8.24
- 使用 Foundry 框架
- 遵循 Solidity Style Guide

### 命名
```solidity
// 合约: PascalCase
contract BettingPool { ... }

// 函数: camelCase
function placeBet() external { ... }

// 状态变量: camelCase
uint256 public totalBetRed;

// 常量: UPPER_SNAKE_CASE
uint256 constant PROTOCOL_FEE_BPS = 500; // 5%

// 事件: PascalCase
event BetPlaced(uint256 indexed gameId, address indexed user);

// 枚举: PascalCase
enum GameStatus { Open, Locked, Finished }
```

### 安全
- 使用 OpenZeppelin 库（ReentrancyGuard, SafeERC20）
- 所有外部输入做校验
- 使用 `require` + 有意义的错误信息
- 注意整数溢出（Solidity 0.8+ 自动检查）

### Gas 优化
- 不存储可以计算的数据
- 使用 `calldata` 而非 `memory`（只读参数）
- 合并状态变量读写
- 事件用 `indexed` 方便前端过滤

---

## TypeScript / 前端代码规范

### 基本规则
- TypeScript strict mode
- 使用函数式组件 + Hooks
- 组件文件 PascalCase，工具文件 camelCase

### 命名
```typescript
// 组件: PascalCase
function BetPanel() { ... }

// Hooks: use 前缀
function useWebSocket() { ... }

// 类型: PascalCase
interface GameState { ... }
type Side = 'red' | 'blue'

// 常量: UPPER_SNAKE_CASE
const MAX_ROUNDS = 30
```

### 组件结构
```typescript
// 1. 类型定义
interface Props { gameId: number }

// 2. 组件
export function BetPanel({ gameId }: Props) {
  // 3. Hooks
  const { data } = useGame(gameId)
  
  // 4. 计算
  const odds = useMemo(() => calculateOdds(data), [data])
  
  // 5. 渲染
  return <div>...</div>
}
```

---

## 通用规则

### 注释
- 公共 API 必须有注释
- 不写显而易见的注释
- 用 TODO/FIXME 标记待办

### 日志
- 使用结构化日志（zap 或 slog）
- 关键操作必须记录日志
- 日志级别: Debug / Info / Warn / Error

### 配置
- 环境变量 > 配置文件 > 默认值
- 敏感信息（私钥、API Key）只用环境变量
- 提供 `.env.example` 文件
