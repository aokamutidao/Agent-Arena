// 合约地址（Sepolia）— 实际在用的部署
export const USDC_ADDRESS = (process.env.NEXT_PUBLIC_USDC_ADDRESS ||
  "0x9163ad7caf7bf73ff105c658a42e12eaa33f58af") as `0x${string}`;

export const ARENA_ADDRESS = (process.env.NEXT_PUBLIC_ARENA_ADDRESS ||
  "0x0bccc77f672f22f843b83ba65297e6ee693c68c7") as `0x${string}`;

// BettingPool 是 AgentArena 内部部署的独立合约
export const POOL_ADDRESS = (process.env.NEXT_PUBLIC_POOL_ADDRESS ||
  "0xa92f3c9bad4330eb0d3d42e9e6073b577b7e9782") as `0x${string}`;

// MockUSDC ABI — 只需要 mint/approve/balanceOf/allowance
export const USDC_ABI = [
  {
    name: "mint",
    type: "function",
    stateMutability: "nonpayable",
    inputs: [
      { name: "to", type: "address" },
      { name: "amount", type: "uint256" },
    ],
    outputs: [],
  },
  {
    name: "approve",
    type: "function",
    stateMutability: "nonpayable",
    inputs: [
      { name: "spender", type: "address" },
      { name: "amount", type: "uint256" },
    ],
    outputs: [{ name: "", type: "bool" }],
  },
  {
    name: "balanceOf",
    type: "function",
    stateMutability: "view",
    inputs: [{ name: "account", type: "address" }],
    outputs: [{ name: "", type: "uint256" }],
  },
  {
    name: "allowance",
    type: "function",
    stateMutability: "view",
    inputs: [
      { name: "owner", type: "address" },
      { name: "spender", type: "address" },
    ],
    outputs: [{ name: "", type: "uint256" }],
  },
  {
    name: "decimals",
    type: "function",
    stateMutability: "view",
    inputs: [],
    outputs: [{ name: "", type: "uint8" }],
  },
] as const;

// AgentArena ABI — betAndVote + 读取函数
export const ARENA_ABI = [
  {
    name: "betAndVote",
    type: "function",
    stateMutability: "nonpayable",
    inputs: [
      { name: "gameId", type: "uint256" },
      { name: "side", type: "bool" },
      { name: "amount", type: "uint256" },
      { name: "strategy", type: "uint8" },
    ],
    outputs: [],
  },
  {
    name: "getOdds",
    type: "function",
    stateMutability: "view",
    inputs: [{ name: "gameId", type: "uint256" }],
    outputs: [
      { name: "oddsRed", type: "uint256" },
      { name: "oddsBlue", type: "uint256" },
    ],
  },
  {
    name: "getStrategyWeights",
    type: "function",
    stateMutability: "view",
    inputs: [{ name: "gameId", type: "uint256" }],
    outputs: [
      { name: "aggressive", type: "uint256" },
      { name: "defensive", type: "uint256" },
      { name: "tricky", type: "uint256" },
    ],
  },
  {
    name: "getGame",
    type: "function",
    stateMutability: "view",
    inputs: [{ name: "gameId", type: "uint256" }],
    outputs: [
      { name: "gameId", type: "uint256" },
      { name: "agentRed", type: "address" },
      { name: "agentBlue", type: "address" },
      { name: "totalBetRed", type: "uint256" },
      { name: "totalBetBlue", type: "uint256" },
      { name: "bettingDeadline", type: "uint256" },
      { name: "status", type: "uint8" },
      { name: "winner", type: "uint8" },
    ],
  },
] as const;

// 策略枚举映射（匹配 ArenaTypes.Strategy）
export const STRATEGY_ENUM: Record<string, number> = {
  aggressive: 0,
  defensive: 1,
  tricky: 2,
};

// BettingPool ABI — 查询下注 + 领奖
export const POOL_ABI = [
  {
    name: "getReward",
    type: "function",
    stateMutability: "view",
    inputs: [
      { name: "gameId", type: "uint256" },
      { name: "user", type: "address" },
    ],
    outputs: [{ name: "", type: "uint256" }],
  },
  {
    name: "bets",
    type: "function",
    stateMutability: "view",
    inputs: [
      { name: "gameId", type: "uint256" },
      { name: "user", type: "address" },
    ],
    outputs: [
      { name: "side", type: "uint8" },
      { name: "amount", type: "uint256" },
      { name: "claimed", type: "bool" },
    ],
  },
  {
    name: "games",
    type: "function",
    stateMutability: "view",
    inputs: [{ name: "gameId", type: "uint256" }],
    outputs: [
      { name: "gameId", type: "uint256" },
      { name: "agentRed", type: "address" },
      { name: "agentBlue", type: "address" },
      { name: "totalBetRed", type: "uint256" },
      { name: "totalBetBlue", type: "uint256" },
      { name: "bettingDeadline", type: "uint256" },
      { name: "status", type: "uint8" },
      { name: "winner", type: "uint8" },
    ],
  },
  {
    name: "claim",
    type: "function",
    stateMutability: "nonpayable",
    inputs: [{ name: "gameId", type: "uint256" }],
    outputs: [],
  },
] as const;
