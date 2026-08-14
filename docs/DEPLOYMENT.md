# 部署流程

## 合约部署（Sepolia）

### 前置条件
- Sepolia ETH（水龙头获取）
- Sepolia USDC 合约地址: `0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238`
- Foundry 安装: `curl -L https://foundry.paradigm.xyz | bash`

### 部署步骤
```bash
cd contracts

# 1. 编译
forge build

# 2. 测试
forge test

# 3. 部署（需要环境变量）
export PRIVATE_KEY=0x...
export SEPOLIA_RPC_URL=https://...

forge script script/Deploy.s.sol:Deploy \
  --rpc-url $SEPOLIA_RPC_URL \
  --private-key $PRIVATE_KEY \
  --broadcast

# 4. 记录合约地址
echo "ARENA_ADDRESS=0x..." >> .env
```

### 验证合约
```bash
forge verify-contract <address> BettingPool \
  --chain-id 11155111 \
  --etherscan-api-key $ETHERSCAN_API_KEY
```

---

## 后端部署

### 本地开发
```bash
cd backend
export QWEN_API_KEY=...
export SEPOLIA_RPC_URL=...
export PRIVATE_KEY=...
go run cmd/server/main.go
```

### 环境变量
```bash
# .env
QWEN_API_KEY=sk-...
SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/...
PRIVATE_KEY=0x...
ARENA_CONTRACT_ADDRESS=0x...
USDC_ADDRESS=0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238
```

---

## 前端部署

### 本地开发
```bash
cd frontend
npm install
npm run dev
```

### 环境变量
```bash
# .env.local
NEXT_PUBLIC_ARENA_ADDRESS=0x...
NEXT_PUBLIC_USDC_ADDRESS=0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238
NEXT_PUBLIC_BACKEND_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080
```

---

## 部署顺序

```
1. 合约部署 → 记录合约地址
2. 后端部署 → 配置合约地址 + API Key
3. 前端部署 → 配置合约地址 + 后端 URL
```

**必须先部署合约**，因为后端和前端都依赖合约地址。
