# 运行、部署与恢复

本文件是可复现操作说明。实际执行过的步骤见 [测试报告](TEST_REPORT.md)；Docker、云服务、测试网 URL 和真实钱包演示不能仅凭配置文件视为已验证。

## 当前交付选择（2026-10-10）

用户确认优先完成本地真实 LLM 对战和新赛事的资金演示，云端前后端均可使用 Render 免费 Web Services。免费服务闲置 15 分钟后休眠，唤醒通常约一分钟；重启、休眠或重新部署会丢失 SQLite 本地文件。用户接受这一演示限制，视频与报告必须如实说明。正式录制以本地服务为主，云端部署仍以实际 URL 和外网验收为准。

本地 SQLite、链状态及对战证据需留档。链上事件可回填，但 Agent 定义、对战回放和审核理由不能仅从链上 hash 恢复。公网数据库重置后，已有报名可能找不到对应 Agent 定义；重新准备演示资料或恢复备份后再继续。

Render 部署采用两个服务，分别使用 `deploy/Dockerfile.frontend` 和 `deploy/Dockerfile.backend`，构建上下文均为仓库根目录。前端构建时配置后端公开 HTTPS `NEXT_PUBLIC_API_URL`；后端 `AUTH_DOMAIN`/`CORS_ORIGINS` 对应前端域名。`compose.yaml` 中空 API 地址依赖 Caddy 同域反代，不能原样用于两个独立 Render 域名。免费云环境不运行本机 Anvil，仍需可访问的测试网 RPC 与新课程合约。

来源：[Render 免费服务说明](https://render.com/docs/free)（查阅 2026-10-10）。后文持久化 Compose 模板保留，供后续长期部署学习使用。

## 1. 环境与第一次构建

- Go 1.25.4、Node.js 24.11.1 / npm、Foundry 1.5.1、Solidity 0.8.24。
- SQLite 使用 `go-sqlite3`，需要 CGO 和系统 C 编译器。
- 首次安装/构建预留数 GB 磁盘；依赖已存在时无需重复安装。
- 所有以下命令默认从仓库根目录执行，括号内切换目录不会改变下一步的位置。

```bash
bash tooling/install-contract-deps.sh
(cd contracts && forge build)
(cd frontend && npm ci)
(cd backend && mkdir -p tmp && go build -o tmp/course-server ./cmd/course-server)
```

合约依赖安装脚本固定 forge-std 和 OpenZeppelin 的提交，不使用漂移的 latest。ABI 来自本次合约：修改合约后执行 `(cd contracts && forge build && python3 export_abi.py)`，同时检查生成的前端 ABI。

## 2. 一条命令验证本地闭环

```bash
node tooling/smoke-course.mjs
```

脚本自动启动新的 loopback Anvil，部署测试代币和课程合约，分配不同钱包角色，再运行 HTTP 与钱包交易联调。它不连接已有 RPC，也不读取项目 `.env` 中的部署私钥。

验证范围：单次 nonce 登录、防重放、越权运行拒绝、真实代币赞助/报名/开赛、模拟对战证据哈希、独立审核/奖金或平局退款、独立争议退款、事件资金对账和后端重启恢复。争议场景中的 fixture 结果会明确标识，不冒充实际 AI 比赛。

成功结果保存在 `.artifacts/course-smoke.json`，包括每笔交易 gasUsed。临时链、数据库和测试签名身份只用于本地验证，脚本退出后清理自己创建的进程和临时目录。

## 3. 保留演示并打开前端

```bash
node tooling/smoke-course.mjs --serve
```

完成验证后保留 RPC `http://127.0.0.1:8545` 与 API `http://localhost:8080`，输出各钱包地址，并创建一场可查看的开放赛事。在另一个终端运行：

```bash
(cd frontend && NEXT_PUBLIC_API_URL=http://localhost:8080 npm run dev)
```

打开 `http://localhost:3000`。本地链 chain ID 为 **31337**，代币 DUSD 为 **6 位小数**。`--serve` 使用 Anvil 公开测试助记词及账号索引：

| 索引 | 本地演示身份 |
| --- | --- |
| 0 | 管理员、测试币发行者 |
| 1 | 主办方 |
| 2 | 指定审核员 |
| 3 / 4 | 两位参赛者 |
| 5 | 赞助方 |

标准 Anvil 测试助记词为 `test test test test test test test test test test test junk`。只导入独立测试钱包配置，绝不向这些公开身份存入真实资产。MetaMask 添加本地网络后，可用测试账户分别演示操作；账号地址以脚本实际输出为准。

脚本会推进本地链时间以验证截止与审核窗口；保留演示时每秒主动挖一个本地区块，使报名/审核截止时间能继续推进。需要继续手工演示时，以链上的时间和状态为准。`Ctrl+C` 停止脚本创建的 API/Anvil，结束本次临时演示；它不会终止其他服务。固定端口已占用时脚本拒绝启动。可为自己的新演示选择其他端口：

```bash
COURSE_SMOKE_RPC_PORT=18545 COURSE_SMOKE_API_PORT=18080 COURSE_SMOKE_UI_PORT=3001 node tooling/smoke-course.mjs --serve
(cd frontend && NEXT_PUBLIC_API_URL=http://localhost:18080 NEXT_PUBLIC_LOCAL_RPC_URL=http://127.0.0.1:18545 npm run dev -- --port 3001)
```

使用非默认 RPC 时，MetaMask 中对应网络的 RPC 地址也须一致。

## 4. 日常持久化开发

上面的临时演示适合自检；需要保存自己的 Agent 和赛事时，使用独立持久 Anvil 数据与 SQLite。

1. 在自己的本地节点部署 `DeployCourse.s.sol`，记下新 Manager 地址及**真实部署回执**所在区块。`forge script` 模拟阶段打印的区块不能代替实际回执。
2. 复制 `.env.course.example` 为 `.env.course`，用 `openssl rand -hex 32` 生成 JWT secret，填写新部署配置。
3. 运行 `./start.sh`，默认 API 8080 / 前端 3000。脚本仅加载根目录 `.env.course`，不加载旧 `backend/.env`。

关键变量：

| 变量 | 含义 |
| --- | --- |
| `JWT_SECRET` | 至少 32 字符的随机会话签名密钥；不能使用示例或旧硬编码值 |
| `DATABASE_PATH` | 从 backend 工作目录解析；默认 `course.db`；不得指向旧 `arena.db` |
| `AUTH_DOMAIN` | 实际前端 host，可含端口；本地为 `localhost:3000` |
| `CORS_ORIGINS` | 精确允许的来源，逗号分隔；不能写 `*` |
| `TRUSTED_PROXY_CIDRS` | 可选、逗号分隔的受控反向代理 IP/CIDR；默认空，不信任伪造转发头 |
| `CHAIN_ID` | 本地 31337 / Sepolia 11155111 |
| `RPC_URL` | 只读 RPC；不需要账户私钥 |
| `TOURNAMENT_MANAGER_ADDRESS` | 新课程 Manager；子合约/代币由后端读取并验证 |
| `DEPLOYMENT_BLOCK` | Manager 实际部署回执区块，用于全量回填 |
| `CONFIRMATIONS` | 本地可为 0；测试网示例 6，需按风险选择 |
| `CONTRACT_ABI_DIR` | 默认 `../contracts/abi` |
| `SIMULATION_ENABLED` | 无模型配置时是否允许明确标注的确定性模拟 |
| `LLM_API_KEY` | 服务端模型凭据；真实模型演示设置 `SIMULATION_ENABLED=false`；不自动读取旧 `backend/.env` |
| `LLM_PROVIDER` / `LLM_MODEL` | 调用记录中的服务商标签 / 实际请求模型名称 |
| `LLM_ENDPOINT` | 服务端配置的完整兼容接口地址，包含 `/chat/completions`；不通过 API 返回或由参赛者设置 |
| `LLM_TIMEOUT_SECONDS` / `LLM_MAX_CALLS` | 单次请求超时 / 每场调用预算；后续重试仍累计预算 |
| `NEXT_PUBLIC_API_URL` | 浏览器访问的 API 基地址；同域反代时可为空字符串 |

`.env.course.example` 预填了本地 RPC，但 Manager 留空。必须填写 Manager，或同时清空 RPC 和 Manager 后进入未部署查看模式；只填写其中一个会拒绝启动。未部署时 `/health` 可用，`/ready` 返回 503，资金操作不可用。

手工编译启动不自动加载任何 `.env`，需要明确导出配置。不要将 `PRIVATE_KEY`、`OWNER_PRIVATE_KEY` 或 `COURSE_DEPLOYER_KEY` 注入课程 API 容器。JWT secret 与模型 key 也不应进入 `NEXT_PUBLIC_*`。

## 5. 测试网合约部署

脚本 `contracts/script/DeployCourse.s.sol` 使用明确的部署钱包；没有 `--broadcast` 时只有模拟。部署环境与后端环境分开：

- `COURSE_DEPLOYER_KEY`：测试网部署钱包，通过本地秘密管理提供。
- `COURSE_ADMIN`：必须为该部署钱包地址。
- `COURSE_ORGANIZER` / `COURSE_AUDITOR`：不同钱包地址。
- `COURSE_TOKEN`：可省略/零地址，脚本会新建 owner-only mint 的 DUSD；使用其他代币前需单独评估兼容性。

在配置这些变量并核对目标测试网之后执行：

```bash
(cd contracts && forge script script/DeployCourse.s.sol:DeployCourse --rpc-url "$RPC_URL")
# 查看模拟结果，确认部署身份、chain ID、费用和角色后，由部署操作者执行：
(cd contracts && forge script script/DeployCourse.s.sol:DeployCourse --rpc-url "$RPC_URL" --broadcast)
```

保存部署交易、chain ID、所有合约地址、实际部署区块、编译器与优化设置。测试币 owner 可给演示赞助者 mint DUSD；该测试资产不能写成真实 USDC。赞助前 ERC20 `approve` 的 spender 必须是 **PrizeEscrow**，不是 Manager。

旧 Sepolia 地址属于黑客松版。此次源码修复不会升级那些不可变合约，不能直接复用。

## 6. 公有云模板

仓库提供 `compose.yaml`、`deploy/Dockerfile.backend`、`deploy/Dockerfile.frontend`、`deploy/Caddyfile`。设计为单实例 API + SQLite 持久卷 + Next.js standalone + Caddy HTTPS。主机需要 Docker Compose、域名 DNS 指向，以及入站 80/443。

在服务器另建忽略于 Git 的 `.env.course`，填写 `SITE_DOMAIN`、JWT、链配置、实际部署区块；云环境 chain ID 与 RPC 必须匹配。通过下面命令预检并启动：

```bash
docker compose --env-file .env.course config --quiet
docker compose --env-file .env.course up -d --build
```

前端构建时 API URL 为空，浏览器使用同域 `/api`，Caddy 分流到 Go 服务。容器仅代理暴露端口，SQLite 位于命名卷 `course_data`。云端应将 `TRUSTED_PROXY_CIDRS` 明确配置为实际 Caddy 容器 IP 或独立受控 Docker 网络 CIDR（通过 `docker network inspect` 核对），不要设成信任全网。否则后端不会使用转发的客户端 IP，代理后的用户会共享同一个每分钟写请求限流桶。直接本地访问保持为空。不要对 API 简单扩容多个副本：当前任务调度与 SQLite 是单实例设计。

验收时检查：公开 HTTPS URL、`/health`、`/ready`、钱包登录消息中的正确域名、账户/网络切换、所有角色流程、容器重启后数据保持、拒签/断网提示、手机宽度及 RPC 确认延迟。保存实际截图和检查时间，而不只保存 Dockerfile。

当前 CI 配置含合约、后端 race/vet、前端和 Anvil 集成四个 job；只有推送后真实成功的 Actions 记录才能称为 CI 通过。

## 7. 备份与故障恢复

链上日志可以重建，但 Agent 定义、对战记录和审核理由不能从 hash 还原。定期备份整份 SQLite，并验证恢复。

本地可使用 SQLite 在线备份接口：

```bash
sqlite3 backend/course.db ".backup 'backend/course.backup.db'"
```

云部署在有 SQLite 工具的运维环境执行等价备份，或停止 API 后一致地备份整个数据卷；不要在运行中只复制主 `.db` 文件而遗漏 WAL。恢复到新路径后，以**同一部署身份**启动服务，检查 `/ready`、证据 hash、资金累计和事件条数。更换 JWT secret 会使旧登录失效。

- **RPC 暂时失败**：`/ready` 不通过，保留数据库，修复 RPC 后继续索引。
- **链 checkpoint 改变**：暂停操作，检查是否重组或换错节点。先备份；停止 API；在同一部署的备份副本上清空 `financial_ledger`、`chain_events`、`tournaments`、`index_cursors`、`index_checkpoints` 后重新回填。保留 `agent_profiles`、`matches`、`review_notes`、`database_identity`。重新核对链下比赛是否仍对应当前链上报名/结果，不能把旧证据自动当成新事实。
- **换链或重新部署**：新建数据库路径，不绕过 deployment identity 检查。
- **对战中断**：主办方重新请求运行，恢复最后保存的回合；已完成比赛不可重跑覆盖。
- **审核员失联/撤权**：等结算期限后，任何钱包触发 expire，赞助者自行 claimRefund。

课程版尚未提供多区域备份、自动故障转移、KMS 或完整监控告警；报告需如实界定。
