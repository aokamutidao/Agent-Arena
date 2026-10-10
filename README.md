# Agent Arena — AI 赛事奖金托管 DApp

面向小型赛事主办方、赞助方、Agent 开发者和审核员的课程项目。赞助资金先进入合约，主办方提交对战证据，独立审核员确认或裁决，获胜者领取奖金；赛事取消或超时后，赞助者按原出资领取退款。

本版本由黑客松 AI 对战原型重构，保留 Go 游戏引擎与 Agent 决策能力。当前规则为两名参赛者、免费报名、赞助奖池、赢家全额领取、无平台抽成。

## 已实现的模块

| 层 | 实现 |
| --- | --- |
| 合约 | TournamentManager 协调登记、托管、审核三个合约；13 类业务写交易；角色、时间窗口、独立审核、超时退款 |
| 后端 | Go / Gin REST API、钱包 nonce 登录、链上角色检查、SQLite、确认事件索引、持久化对战与证据 |
| 前端 | Next.js / wagmi / RainbowKit；MetaMask、资金 Dashboard、赛事操作、审计流水、Agent 定义、Gas 估算 |
| Agent | 不可变定义绑定参赛钱包；显式配置的 Qwen 模型或清楚标注的确定性模拟；完整回合记录 |
| 工程 | Foundry 正常/对抗/fuzz/invariant 测试、后端与前端测试、本地 Anvil 集成脚本、CI 和 Docker 部署配置 |

每笔资金操作由钱包签名；课程后端没有交易签名私钥。测试代币为 6 位小数的 **DUSD**，没有现实兑付承诺。结果哈希支持记录核对，比赛真实性仍依赖运行器与独立审核。

## 先从这里阅读

完整文档入口：[课程文档导航](docs/course-project/README.md)。

- [需求与作业映射](docs/course-project/REQUIREMENTS.md)
- [架构与设计取舍](docs/course-project/ARCHITECTURE.md)
- [合约与 13 类交易](docs/course-project/CONTRACT.md)、[后端 API](docs/course-project/API.md)
- [运行与部署](docs/course-project/DEPLOYMENT.md)、[用户手册](docs/course-project/USER_GUIDE.md)、[本地真实模型演示](docs/course-project/LOCAL_LLM_DEMO.md)
- [测试报告](docs/course-project/TEST_REPORT.md)、[代码学习路线](docs/course-project/LEARNING_GUIDE.md)
- [报告与演示准备](docs/course-project/SUBMISSION_PLAN.md)、[本轮开发日志](devlog/028-course-refactor.md)

## 本地验证

需要 Go 1.25.4、Node.js 24.11.1 / npm、Foundry 1.5.1（含 Anvil）、C 编译器（SQLite 使用 CGO）。首次安装/构建请预留数 GB 空间。

```bash
bash tooling/install-contract-deps.sh
(cd contracts && forge build && forge test)
(cd frontend && npm ci)
(cd backend && mkdir -p tmp && go build -o tmp/course-server ./cmd/course-server)
node tooling/smoke-course.mjs
```

脚本创建独立本地链和临时数据库，验证钱包登录、赞助、报名、对战、审核、领奖/退款、账本和重启；结束后清理自己创建的进程与临时数据。结果输出到 `.artifacts/course-smoke.json`。

需要保留本地演示服务时使用 `node tooling/smoke-course.mjs --serve`，另开终端运行 `(cd frontend && npm run dev)`。具体钱包、端口和持久部署方法见运行文档。已有 8545/8080 服务时脚本会拒绝覆盖。

## 状态与边界

执行结果以[测试报告](docs/course-project/TEST_REPORT.md)为准。仓库提供云部署配置，但课程版的公网 URL、Sepolia 部署和真实用户验证仍需实际执行并留证；不应据此宣称已经达到生产级或完成所有课程交付。

033 迭代补齐模型调用记录、失败中断、持续计数的调用预算与分阶段本地验收工具。模型配置不等于连通成功；比赛页面分别展示实际成功/失败次数、模型、耗时与返回的用量。当前课程交付选择本地演示优先、Render 免费服务作为公网展示，冷启动与 SQLite 重置限制需写入报告和视频。

旧下注/PVE 写接口已经停用。历史资料见 [黑客松 README](docs/archive/HACKATHON_README.md) 与 `specs/`；旧 Sepolia 地址不适用于课程版，源代码修复不会改变已经部署的旧合约。

RAG、工具调用、记忆与面试拓展安排在课程版验收之后。
