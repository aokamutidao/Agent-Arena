# Agent Arena - Claude 工作指引

## 项目概述

AI Agent 格子竞技场：两个 AI Agent 在 10x10 格子上对战，观众用 USDC 下注并通过策略投票影响 AI 决策。区块链（Sepolia）托管下注和结算，Go 后端驱动游戏引擎和 LLM Agent，Next.js 前端提供观战和下注界面。

**黑客松 Agent 赛道参赛项目。**

---

## 工作原则

1. **小步迭代**：每次只做 1-2 个任务，做完验证后再继续
2. **日志追踪**：每次迭代必须写 devlog，记录完成项和待办项
3. **Spec 驱动**：实现前读 spec，实现后对照 spec 验收
4. **不确定就问**：不要自己猜需求或技术方案

---

## 核心文件路径

| 文件 | 用途 | 何时读 |
|------|------|--------|
| `CLAUDE.md` | 本文件，工作入口 | **每次开始工作前** |
| `docs/ITERATION_PLAN.md` | 当前迭代计划和任务拆解 | 每次开始前确认当前任务 |
| `docs/ARCHITECTURE.md` | 系统架构 + 模块依赖 | 涉及跨模块改动时 |
| `docs/CODE_STANDARDS.md` | 编码规范（Go/Solidity/TS） | 写代码时参照 |
| `docs/TESTING.md` | 测试标准 + 覆盖率要求 | 写测试时参照 |
| `docs/GIT_WORKFLOW.md` | Git 提交规范 | 提交代码时参照 |
| `docs/DEPLOYMENT.md` | 部署流程 | 部署时参照 |
| `specs/` | 产品和技术规格（已锁定） | 实现前确认需求 |
| `devlog/` | 开发日志（按迭代编号） | 开始前读最新一篇 |

---

## 每次开始工作前的检查清单

```
□ 1. 读 CLAUDE.md（本文件）
□ 2. 读 docs/ITERATION_PLAN.md → 确认当前任务
□ 3. 读最新的 devlog → 了解上次进度和遗留问题
□ 4. 读对应的 spec → 确认实现要求
□ 5. 读 docs/CODE_STANDARDS.md → 确认编码规范
```

---

## 每次完成工作后的检查清单

```
□ 1. 写 devlog/XXX-<title>.md（记录完成项 + 待办项 + 遇到的问题）
□ 2. 更新 docs/ITERATION_PLAN.md（勾选已完成任务）
□ 3. 如果有架构变动 → 更新 docs/ARCHITECTURE.md
□ 4. 如果测试通过 → 在 devlog 中记录测试结果
```

---

## 项目目录结构

```
agent-arena/
├── CLAUDE.md              # 本文件
├── README.md              # 项目介绍
├── specs/                 # Spec 文档（已锁定，不改）
├── docs/                  # 开发标准文件
├── devlog/                # 开发日志
├── contracts/             # Solidity 合约 (Foundry)
│   ├── src/
│   ├── test/
│   ├── script/
│   └── foundry.toml
├── backend/               # Go 后端 (Gin)
│   ├── cmd/server/
│   ├── internal/
│   │   ├── engine/        # 游戏引擎
│   │   ├── agent/         # AI Agent (LLM)
│   │   ├── blockchain/    # 合约交互
│   │   ├── betting/       # 下注逻辑
│   │   └── api/           # HTTP/WS 服务
│   ├── go.mod
│   └── go.sum
├── frontend/              # Next.js 前端
│   ├── src/
│   └── package.json
└── scripts/               # 辅助脚本（部署、测试等）
```

---

## 技术栈速查

| 层 | 技术 | 版本 |
|----|------|------|
| 链 | Sepolia 测试网 | - |
| 下注币 | USDC (ERC20) | - |
| 合约 | Solidity + Foundry | 0.8.24 |
| 后端 | Go + Gin + gorilla/websocket | 1.21+ |
| LLM | Qwen API (qwen-turbo) | - |
| 前端 | Next.js 14 + Tailwind + shadcn/ui | - |
| 钱包 | wagmi + RainbowKit | - |

---

## 迭代节奏

```
每个迭代：
  1. 确认任务（读 ITERATION_PLAN.md）
  2. 读相关 spec
  3. 实现代码
  4. 测试验证
  5. 写 devlog
  6. 更新 ITERATION_PLAN.md

不要跳步。不要一次做多个迭代。
```
