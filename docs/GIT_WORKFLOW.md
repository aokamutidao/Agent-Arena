# Git 工作流

## 分支策略

```
main           ← 稳定版本（Demo 用）
  └── dev      ← 开发分支
       └── feat/001-init       ← 功能分支
       └── feat/002-contracts  ← 功能分支
       └── fix/xxx             ← 修复分支
```

## 提交规范

```
<type>(<scope>): <description>

type:
  feat     - 新功能
  fix      - 修复 bug
  docs     - 文档
  test     - 测试
  refactor - 重构
  chore    - 杂项

scope:
  contract - 合约
  engine   - 游戏引擎
  agent    - AI Agent
  api      - 后端 API
  frontend - 前端
  devlog   - 开发日志

示例:
  feat(contract): implement BettingPool
  fix(engine): fix collision detection
  docs(devlog): add iteration 001 log
```

## 提交频率

- 每个迭代完成后至少提交一次
- 大功能可以拆分为多个提交
- devlog 和代码一起提交

## 不提交的内容

- `.env` 文件（包含私钥等敏感信息）
- `node_modules/`
- `out/` (Foundry 编译输出)
- `.DS_Store`
