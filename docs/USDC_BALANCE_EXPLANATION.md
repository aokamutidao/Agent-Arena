# USDC 余额说明

## 问题
用户在挑战界面看到有 USDC 余额，但钱包页面显示 0。

## 原因

### MockUSDC 机制
MockUSDC 是一个测试用的 ERC20 合约（部署在 Sepolia 测试网），它有一个公开的 `mint` 函数，任何人都可以调用它来给自己铸造测试 USDC。

- **挑战界面（BetPanel）**：通过 wagmi 直接从链上读取用户钱包地址的 USDC 余额
- **钱包页面**：通过后端 API 从链上读取用户登录地址的 USDC 余额

### 合约重部署导致余额清零
在 2026-08-12 03:14，我们重新部署了 AgentArena 和 MockUSDC 合约（因为旧合约的 owner 与后端 deployer 不匹配）。

- **旧 MockUSDC**: `0x7775B152Bba82eb8405a553E9E6A0D5bB9BD35da`
- **新 MockUSDC**: `0x9163ad7caf7bf73ff105c658a42e12eaa33f58af`

用户在旧合约上 mint 的 USDC 不会自动转移到新合约，所以新合约上用户的余额是 0。

## 解决方案

用户需要在新合约上重新 mint USDC：

1. 进入任意对局页面（如 `/game/113`）
2. 在 BetPanel 中找到"领 100 测试 USDC"按钮
3. 点击按钮，会调用新 MockUSDC 合约的 `mint` 函数
4. 铸造成功后，钱包页面会显示 100 USDC

## 技术细节

MockUSDC 合约的 mint 函数（任何人都可以调用）：
```solidity
function mint(address to, uint256 amount) public {
    _mint(to, amount);
}
```

BetPanel 中的 mint 调用：
```typescript
const { writeContract } = useWriteContract();
writeContract({
  address: USDC_ADDRESS,
  abi: USDC_ABI,
  functionName: "mint",
  args: [address, BigInt(100 * 10 ** 6)], // 100 USDC (6 decimals)
});
```
