# 课程版测试与评价记录

测试日期：2026-10-09，环境为本地 macOS / Go 1.25.4 / Node.js 24.11.1 / Foundry 1.5.1 / Solidity 0.8.24（optimizer 200）。本报告区分已执行证据与待验收项目。它不是外部安全审计，也不等于生产部署认证。

## 本轮验证概览

| 层 | 实际结果 | 证据 / 命令 |
| --- | --- | --- |
| 合约 | 49 项通过，含正常业务、越权/期限/资金边界、旧缺陷回归 | `contracts/reports/foundry-gas-report.txt` |
| 模糊与状态测试 | 两项资金 fuzz 各 256 runs；stateful invariant 64 × 64 = 4,096 次调用 | 同上，`contracts/test/` |
| 合约覆盖率 | Manager 行 98.59%、分支 92.50%；Escrow/Review/Token 行 100%；Registry 行 96.72% | `contracts/reports/foundry-coverage.txt` |
| 后端 | 全部包 test / race / vet、最终 course-server build 通过 | [race 原始输出](evidence/backend-race-2026-10-09.txt)，含课程与历史模块 |
| 前端 | TypeScript、lint、6 项业务单元测试、生产 build 已通过；UI 检查见下方后续记录 | `frontend/tests/course.test.cjs` |
| 本地全链路 | 23 笔成功钱包交易，实际 ERC20 付款/退款、零负债、事件索引、证据与重启通过 | [原始 JSON](evidence/local-integration-2026-10-09.json) |
| 格式/脚本 | Foundry fmt、`git diff --check`、Shell/Node 脚本语法检查通过 | 仓库命令可复现 |

旧合约和部署脚本也计入全仓覆盖，总行覆盖 **80.54%**；不能将局部合约 100% 写成全仓 100%。覆盖率衡量执行过的代码，不证明不存在漏洞。

## 合约验证场景

- 正常创建、赞助、报名/退出、开始、提交、审核、领奖。
- 主办方取消、报名不足超时、结算超时、独立争议裁决与退款。
- 不存在赛事、零金额/哈希、重复报名/领奖/退款、过期操作、错误状态和权限拒绝。
- 指定审核员不能兼任该赛事主办方、参赛者或赞助者。
- 完整异议窗口必须保留；撤销角色不能阻塞 permissionless timeout 退出路径。
- 恶意 ERC20 重入、出款失败状态回滚、手续费代币实际到账不符拒绝。
- 随机资金操作下累计已付不超过实际入金，合约负债与托管余额相容。
- 历史换边下注多领、无效对局入金、截止不执行、投票身份、任意烧毁等源码回归。

历史源码修复不改变旧 Sepolia 部署。新的课程测试只证明此次本地源码对应部署的行为。

## Gas：可复现的规模对照与实际交易

受控实验固定同一编译配置，以 1、16、64 个赞助者各资助同一场赛事，分别测量一次提现。快照见 `contracts/snapshots/PrizeScale.json` 与 `RefundScale.json`。

| 操作（受控合约调用） | 1 位赞助者 | 16 位 | 64 位 |
| --- | ---: | ---: | ---: |
| claimPrize | 66,162 | 66,162 | 66,162 |
| claimRefund | 63,976 | 63,976 | 63,976 |

结论仅限这个实现与样本：领取不遍历赞助者，赞助者增加没有增加单次领取开销。SafeERC20、重入锁、独立账本与事件仍有固定成本。

独立 Anvil 钱包交易回执测得：

| 业务交易 | gasUsed |
| --- | ---: |
| createTournament（第一场） | 230,380 |
| fundTournament（第一场） | 192,196 |
| startTournament | 88,526 |
| submitResult（第一场） | 238,117 |
| approveResult | 152,699 |
| claimPrize | 136,106 |
| disputeResult | 124,506 |
| resolveDispute（取消） | 135,170 |
| claimRefund | 114,020 |

受控调用 Gas 与独立交易回执的 gasUsed 不是同一种测量：交易固有费用、冷热存储访问及测试调用上下文不同，不能混用或声称减少了差值。前端在每笔操作前根据当前钱包和链状态估算 Gas；实际 ETH 成本还依赖网络费率。

部署大小：Manager runtime 13,614 bytes、initcode 26,742 bytes；Registry 5,005、Escrow 4,515、Review 1,705、Token 2,259 bytes。见 `contracts/reports/foundry-sizes.txt`。

## 本地全链路的实际记录

命令：`node tooling/smoke-course.mjs`（本次保留服务还使用 `--serve` 和独立端口）。脚本启动全新 Anvil 与独立临时数据库，通过 viem 使用公开本地测试账户签名。后端不持有交易私钥。

1. 部署 DUSD 和 Manager，分配组织者/审核员角色，给赞助者测试币并授权 Escrow。
2. 钱包 nonce 登录；重放被拒绝；普通参赛者不能调用主办方运行接口。
3. 保存属于不同钱包的不可变 Agent 定义，链上赞助 100 DUSD、报名并绑定 hash。
4. 运行 aggressive 与 defensive 基线；12 回合后蓝方获胜，模式明确为 `deterministic_simulation`。
5. 核对完整证据 Keccak-256，完成比赛不可再次运行覆盖；主办方提交，窗口后审核员确认，蓝方领取 100 DUSD。
6. 第二场注入 25 DUSD；使用明确标为 fixture 的合约争议场景。保存异议/决定原文，核对 hash，与链上事件关联；审核员裁决取消，赞助者收到 25 DUSD。
7. 断言收款钱包实际 ERC20 增量正确，Escrow 最终余额和 totalLiabilities 都为零。
8. 确认流水为入金 125 DUSD、奖金 100 DUSD、退款 25 DUSD，共 4 条财务记录；游标和头部均为区块 27，事件 21 条。
9. 停止并重启 API，使用同一数据库重新索引；资金、事件数、财务记录数、比赛证据、异议/裁决原文保持不变。

原始报告中的合约地址仅对应这条临时本地链，不是公网地址。`--serve` 会额外创建第三场 200 DUSD 的开放演示赛事，因此浏览器 Dashboard 在演示期间显示累计 325 DUSD 是预期差异，不能与前两场测试快照混淆。

## API 小样本响应测量

主业务完成并索引后，在同一机器用 Node fetch 顺序请求，每个端点 30 次、并发 1，计算排序后的 p50/p95：

| 端点 | p50 | p95 |
| --- | ---: | ---: |
| `/api/dashboard` | 0.497 ms | 1.494 ms |
| `/api/tournaments/1/evidence` | 1.181 ms | 1.664 ms |

这是本机 loopback 的小样本观测，没有隔离其他主机负载，不能作为云端 SLA 或高并发结论。脚本业务段约 1,048.8 ms 使用 Anvil 时间快进和模拟决策，不代表现实奖金处理周期。正式影响评价仍需真实用户任务测试、同规则的人工对账基线与较充分样本。

## 复现命令

```bash
(cd contracts && forge fmt --check && forge test --gas-report --gas-snapshot-emit false)
(cd contracts && forge coverage --report summary)
(cd backend && go test ./... && go vet ./...)
(cd backend && go test -race ./...)
(cd frontend && npm run lint && npm run typecheck && npm test && npm run build)
(cd backend && mkdir -p tmp && go build -o tmp/course-server ./cmd/course-server)
node tooling/smoke-course.mjs
```

首次后端编译与前端构建应预留空间。本次沙箱禁止 WebSocket 测试绑定本地端口，授权本地监听后相同 API 测试通过；这不是跳过测试。重复构建时复用了专用 Go 缓存，完成后已清理本轮缓存以释放空间。最后执行了全量 `go test -race -p 1 -count=1 ./...`，所有包通过、无数据竞争；macOS CGO 链接器输出 LC_DYSYMTAB warning，但测试退出码为 0，原始日志保留。

## 实际浏览器与响应式检查

使用真实 Chrome，访问本地生产构建 `http://localhost:3001`，连接本轮独立 API 18080 / Anvil 18545（未连接已有的 3000/8545 服务）：

- Dashboard 准确显示三场赛事、累计赞助 325、已付奖金 100、退款 25、待结算资金 200 DUSD。
- 已完成赛事能显示实际 12 回合数据；回放滑块可切换至第一回合；下载/展示证据 hash 与链上提交一致。
- 审计页能读取真实 SponsorRefunded 事件，Refund 筛选有效。
- 390 px 手机宽度首页与详情页无整页横向溢出，记录 `scrollWidth = clientWidth = 390`；审计表格在局部滚动，移动导航可到 Agent 登录提示。
- 最终复验已修复已结算回放的状态文案、手机长地址换行、名称 UTF-8 长度与合约 120-byte 限制一致性，并通过最终生产构建。
- 只读浏览器检查没有执行钱包签名，不能替代真实 MetaMask 多角色手工验收。

截图： [桌面资金页](evidence/dashboard-desktop.jpg)、[手机资金页](evidence/dashboard-mobile.jpg)、[退款审计](evidence/audit-refund-desktop.jpg)、[比赛回放](evidence/tournament-replay-desktop.jpg)、[手机赛事详情](evidence/tournament-detail-mobile.jpg)。

生产构建首次出现 WalletConnect 在 SSR 阶段访问 IndexedDB 的错误日志，已通过延迟至浏览器加载钱包 Provider 修复并重新构建通过。构建仍提示其日志库的可选 `pino-pretty` 依赖缺失，未导致类型、lint、页面生成或此次浏览器验证失败。

## 待验收事项

2026-10-10 体验准备只读核查：前端 3001、API 18080、Anvil 18545 进程运行，API `/ready` 正常，主办方/审核员/管理员链上角色有效；六个公开测试账户均有本地 ETH，赞助方余额 700 DUSD。#3 已超过最终结算期限，需到期取消及赞助者退款，不能作为新对战起点。[原始快照](evidence/experience-preflight-2026-10-10.json)、[操作指南](EXPERIENCE_GUIDE.md)。本轮没有重跑测试套件或通过 MetaMask 发送交易，下列待验收状态不变。

- 公有云 Docker 镜像构建、HTTPS、公开 URL、持久卷与云端恢复；当前机器没有 Docker CLI。
- 新课程合约的 Sepolia 广播/验证、真实网络费用及确认深度行为。
- MetaMask 实际扩展中的完整多角色流程、用户拒签、网络切换和移动端钱包体验。
- 真实 Qwen 服务调用、模型错误率、成本与对局公平性；本次闭环使用确定性模拟。
- 多用户并发/长时间运行/资源耗尽/自动灾备测试；当前为单实例 SQLite 方案。
- 真实利益相关者反馈与作业报告、Slides、视频、个人贡献的最终版本。

CI 工作流已经配置，但本次没有推送 GitHub，不能记录远程 CI 已通过。
