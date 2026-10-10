# 本地真实 LLM 演示：学习与操作

目标是在已有课程链上新建赛事，实际运行模型，再验证奖金结算与独立取消退款。**2026-10-10 已完成真实闭环：赛事 #5 使用 `qwen3.5-plus` 完成 14 回合、28 次成功调用，蓝方领取 100 DUSD；赛事 #6 取消后退回 25 DUSD。** 实际回执、余额、用量及局限见 [测试报告](TEST_REPORT.md)。当前可直接打开 `http://localhost:3002/tournaments/5` 查看记录，无需连接钱包。

所有命令从仓库根目录执行。先理解三条边界：Go API 读取链并保存 Agent、比赛和证据；资金交易由钱包签名；链上 hash 绑定链下内容，不能证明结果本身真实或正确。

## 1. 保留旧环境，准备新入口

| 项目 | 地址或位置 | 用途 |
| --- | --- | --- |
| 原演示 | API `18080` / 前端 `3001` | 保留已有赛事和资料 |
| 新演示 | API `http://localhost:18081` / 前端 `http://localhost:3002` | 已启动并通过验收，新配置与独立 SQLite |
| 共用本地链 | `http://127.0.0.1:18545`，chain ID `31337` | 已部署课程合约，DUSD 为 6 位小数的测试代币 |
| 新资料目录 | `.local/llm-demo/` | 数据库、准备记录、后续业务证据 |

`tooling/prepare-local-demo.py` 已用于在线备份原 SQLite 到 `.local/llm-demo/course.db`，保留已提交的 WAL 数据，并生成 `preparation.json`、`anvil-before.json` 和根目录 `.env.course`。下面是复现模板，先将路径替换为真实旧库；当前已有备份时无需重跑，脚本拒绝覆盖现有配置或目标数据库。

```bash
python3 tooling/prepare-local-demo.py --source-db "/absolute/path/to/source-course.db"
```

两个 API 共用同一条 Anvil 链。新增交易也会进入旧 API 的索引与统计；这不是链完全隔离的环境。不要重启旧 `smoke-course.mjs --serve`、重置链或自动恢复 `anvil-before.json` 来准备新赛事。

## 2. 配置预检与单次 probe

由操作者在 `.env.course` 中填写与模型服务权限匹配的 `LLM_API_KEY`、`LLM_PROVIDER`、`LLM_MODEL`、`LLM_ENDPOINT`，真实演示保持 `SIMULATION_ENABLED=false`。端点应为完整的兼容聊天接口地址；设置单次超时和每场调用预算。当前实测成功组合是模型 `qwen3.5-plus`、请求参数 `enable_thinking:false`，正式比赛也需核对实际请求参数。凭据仅留在服务端，不进入 `NEXT_PUBLIC_*`、Git、视频或证据 JSON。

配置加载后，先执行**不调用模型**的预检：

```bash
set -a
source .env.course
set +a
(cd backend && go run ./cmd/llm-check)
```

输出应明确 `request_sent=false`，只报告脱敏配置。`configured=true` 仅说明存在配置，不能替代模型成功响应。

按已授权的模型调用范围执行以下 probe；每次发送一个请求，可能产生费用。先阅读脱敏诊断结果再针对性调整配置：

```bash
(cd backend && go run ./cmd/llm-check --call)
```

probe 的 `status=succeeded` 还不等于完整比赛成功。实际诊断顺序：首次 HTTP 400 未保留错误正文；第二次确认 `qwen3-coder-plus` 不受支持，第三次确认 `qwen3-coder-next` 不受支持；第四次改用 `qwen3.5-plus`、`enable_thinking:false` 成功，HTTP 200、`WAIT`、43 tokens、1468 ms。随后另行完成了上述真实赛事验收，模型每次调用的元数据与回合一并留档。

新 API / 前端按 [运行与部署](DEPLOYMENT.md) 启动；`./start.sh` 读取 `.env.course` 的 `PORT=18081`、`FRONTEND_PORT=3002`。已运行时不重复启动。核对 `/ready` 和 `/api/config`：链 ID、Manager 和子合约必须匹配，索引须追上链，`llm_configured=true`、`execution_mode=llm`。这两个 LLM 字段仍只是配置状态，不能替代 probe 或比赛记录。

## 3. 分阶段验收新赛事

`tooling/local-course-flow.mjs` 只允许 loopback API/RPC、Anvil 和 chain ID `31337`。它使用公开的 Anvil 开发账户，不读取用户钱包密钥，不部署合约、不授予角色、不铸币；主办方、审核员和赞助者必须已有权限与测试余额。

先查看命令并做离线自检，再准备非资金资料与交易前余额：

```bash
node tooling/local-course-flow.mjs --help
node tooling/local-course-flow.mjs --self-check
node tooling/local-course-flow.mjs agents
node tooling/local-course-flow.mjs status > .local/llm-demo/preflight-balances.json
```

`agents` 会签名登录并保存两个钱包绑定的 Agent 策略，不发送链上交易；策略要求合法动作、尽快交锋，不指定赢家。`status` 完全只读，重定向保存的是执行当时余额，只有在新交易之前采集才能称为初始快照。

确认 probe 成功后，依次运行：

```bash
node tooling/local-course-flow.mjs prepare
node tooling/local-course-flow.mjs run
node tooling/local-course-flow.mjs settle
node tooling/local-course-flow.mjs cancel
```

| 阶段 | 实际动作与检查 |
| --- | --- |
| `prepare` | 创建本次全新赛事、赞助 100 DUSD、红方报名→退出→重报、蓝方报名；不调用模型。模型配置缺失时会在创建或出资之前停止。 |
| `run` | 等报名截止，主办方开赛并请求一次真实 LLM 比赛；轮询已有比赛，不重复请求正在运行或已完成的比赛。 |
| `settle` | 核验 `mode=llm`、调用统计及每回合双方成功调用记录、完整证据 hash；使用实际赢家提交结果，等完整审核窗口，指定审核员批准，再由赢家领取。 |
| `cancel` | 另建独立 Open 赛事、赞助 25 DUSD、保存取消理由 hash、主办方取消，赞助者领取原额退款。此场不运行模型。 |

报名截止、最终结算期限和审核窗口都根据**链时间**计算，默认分别为创建后 180 秒、7200 秒和提交后 60 秒。实际比赛耗时取决于模型，不能保证十分钟内完成。正常命令等待链时间推进；若链没有继续出块，等待会达到预算后停止。

本地演示需要压缩等待时，可对 `run` / `settle` 显式添加 `--allow-time-warp`。脚本记录 `evm_increaseTime` / `evm_mine` 的目标、原因和前后时间；存在其他非本次、未终结赛事时拒绝快进。快进会影响共用链，不能把压缩后的时间当成真实业务处理时长。

## 4. 中断、退款与证据

始终复用 `.local/llm-demo/business-flow.json`，不要删除检查点后重复跑资金阶段。文件保存本次赛事 ID、交易 hash、gas、前后余额、模型调用统计和时间调整。脚本只写本次新建的赛事；保存交易意图后恢复时核对原 hash，不因重跑而自动补发新交易。

模型请求或解析失败会暂停比赛，不编造动作、回合或赢家。先用 `status` 查看；确认允许额外模型调用后，才使用 `run --resume-llm`，预算仍累计历史请求。缺少调用统计的旧 `llm` 标签不能通过真实调用验收。

真实平局不会提交虚构胜者，`settle` 会等待到期后退款。失败或放弃主赛事时，可使用：

```bash
node tooling/local-course-flow.mjs cancel --main-timeout
```

它等待适用期限，再触发到期取消与赞助者退款；已 Finalized 的赛事不能走这条退款路径。只有回执和实际余额差都核验成功，才能写“已领奖 / 已退款”。检查点中不保存模型 key、JWT、完整 provider 请求提示词或原始响应正文；调用统计是后端观察记录，不是服务商签名或独立审计。

## 5. CLI 与 MetaMask 的演示区别

上述 CLI 是公开测试账户的真实本地交易签名与回执验收，**不能声称已完成 MetaMask 界面操作**。浏览器演示需另在 `3002` 连接对应测试账户和 `18545 / 31337` 网络，签名登录后按本场身份操作；交易确认弹窗中核对金额、操作和网络，再在钱包签名。切换主办方、参赛者、审核员或赞助方时，也需使用对应账户。

本地网页和 API 可匿名展示公开赛事、资金流水与已有比赛；写交易需要钱包，“我的 Agent”管理与比赛运行需要对应登录身份。Agent 定义可被公开查询，不应存入秘密。公有云页面不会让其他人的浏览器自动访问你的 `localhost`；公网部署需可访问的 API/RPC 和新测试网部署，不能直接复用本机端口或本地账户来声称 Sepolia 验收。
