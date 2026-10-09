# Agent Arena 工作指引

## 当前目标

当前为课程版：**AI 赛事奖金托管、审核、结算与退款 DApp**。用户已确认从黑客松下注原型重构，并要求保留学习文档；先完成课程，后续再做 RAG / Agent 面试扩展。

课程入口是 `backend/cmd/course-server`，领域代码在 `backend/internal/tournament`，前端默认展示课程界面。旧 `cmd/server` 仅可显式开启只读历史查看；不要恢复旧公开资金写接口，也不要复用旧合约地址。

## 开始前阅读

1. `docs/course-project/README.md` — 当前文档入口。
2. `docs/ITERATION_PLAN.md` 与最新 `devlog/` — 当前进度与待办。
3. `docs/course-project/REQUIREMENTS.md` — 冻结范围与验收。
4. 对应 `ARCHITECTURE.md`、`CONTRACT.md`、`API.md` — 模块边界和接口。
5. `docs/CODE_STANDARDS.md`、`docs/GIT_WORKFLOW.md` — 编码与提交规范。

旧 `specs/` 保留黑客松历史，不修改来伪装为新需求。当前用户要求与课程版需求优先。

## 实现原则

- 先明确需求和接口，再实现、测试、更新文档与开发日志。
- 资金写操作只能由用户钱包签名；后端 ChainReader 保持只读。
- 钱包登录使用一次性 nonce；链上权限由合约最终校验。
- 金额使用整数单位，不能使用浮点数保存或计算资金。
- 原始证据、Agent 定义和审核理由保存链下，哈希绑定链上；明确真实性信任边界。
- 模拟和真实 LLM 模式明确标注，测试网不能描述成真实资产。
- 个人贡献、用户访谈、性能、部署与测试数据只记录实际证据。
- 每轮完成后更新 devlog、迭代进度、相关接口文档和测试报告。提交仅包含本轮文件，保留用户已有未跟踪文件。

## 验证入口

见 `docs/course-project/TEST_REPORT.md` 和 `.github/workflows/course-ci.yml`。优先跑受影响模块，再跑完整本地链验证。Go 的 SQLite 依赖 CGO；受磁盘限制时串行构建并复用缓存。

## 技术栈

Solidity 0.8.24 / Foundry 1.5.1；Go 1.25.4 / Gin / SQLite；Next.js 14 / TypeScript / wagmi / RainbowKit；Node.js 24.11.1。云部署模板使用 Caddy HTTPS、独立 API 容器与持久 SQLite 卷。
