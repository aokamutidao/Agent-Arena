# 测试标准

## 测试层级

```
1. 单元测试     → 每个函数/方法
2. 集成测试     → 模块间交互
3. 合约测试     → Foundry test
4. 端到端测试   → 完整流程（Demo 前）
```

## Go 测试

### 规则
- 每个 `.go` 文件有对应的 `_test.go`
- 测试覆盖率目标: > 70%
- 使用标准 `testing` 包，不用第三方测试框架

### 命名
```go
func TestFunctionName(t *testing.T)          // 单元测试
func TestFunctionName_Scenario(t *testing.T) // 场景测试
func BenchmarkFunctionName(b *testing.B)     // 性能测试
```

### 必须测试的部分
- 游戏引擎: 动作解析、伤害计算、胜负判定
- AI Agent: prompt 构建、响应解析、动作验证
- API: 请求/响应格式

## Solidity 测试

### 规则
- 使用 Foundry test (`forge test`)
- 每个合约有对应的测试文件
- 覆盖: 正常路径 + 异常路径 + 边界情况

### 必须测试的部分
- 下注: 正常下注、超额下注、截止后下注
- 策略投票: 投票、锁定后投票
- 结算: 正常结算、重复结算
- 提取: 赢家提取、输家提取、重复提取

## 手动测试

### 对局测试
```
□ 创建对局 → 下注 → 开始 → 对战 → 结算 → 提取
□ Agent 行为符合性格设定
□ 策略投票影响 Agent 决策
□ WebSocket 实时推送正常
```

## 运行测试

```bash
# Go
cd backend && go test ./...

# Solidity
cd contracts && forge test -v

# 覆盖率
cd backend && go test -cover ./...
cd contracts && forge coverage
```
