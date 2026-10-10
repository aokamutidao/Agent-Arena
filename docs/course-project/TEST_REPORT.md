# 课程版测试与评价记录

测试日期：2026-10-09，环境为本地 macOS / Go 1.25.4 / Node.js 24.11.1 / Foundry 1.5.1 / Solidity 0.8.24（optimizer 200）。本报告区分已执行证据与待验收项目。它不是外部安全审计，也不等于生产部署认证。

## 2026-10-10：MetaMask 浏览器验收（034，进行中）

已保存[验收前快照](evidence/metamask-before-2026-10-10.json)，确认 API `18081` ready、真实模型配置及六角色余额。Chrome `localhost:3002/tournaments` 实际显示六场赛事；点击“连接钱包”并选择已安装的 MetaMask 后，页面等待扩展确认。

原生 Chrome 窗口读取超时，暂未读取到钱包扩展确认页面；已请用户解锁（如需要）并展开 MetaMask 标签页。本轮截至此记录未完成钱包连接、登录签名或新交易，未发起新模型比赛。完整角色与资金 UI 流程、回执及新视频仍待接续，详见 [034 记录](../../devlog/034-metamask-browser-acceptance.md)。

## 2026-10-10：真实模型与本地资金闭环（033）

本轮新建 #5 主赛事与 #6 取消退款赛事，继续使用已有 Anvil 31337 课程合约。新 API `18081` 使用原 SQLite 的在线备份，前端入口为 `http://localhost:3002/tournaments/5`；原链未重置。实际业务签名来自公开 Anvil 开发账户的 CLI 脚本，**没有将脚本验收记为 MetaMask 弹窗操作或新视频录制**。

| 验收 | 实际结果 |
| --- | --- |
| 模型连通性 | 4 次实际诊断：前 3 次 HTTP 400；脱敏错误确认两个 Coder 模型名不受当前接口支持。第 4 次 `qwen3.5-plus` / `enable_thinking=false` 返回 HTTP 200、合法 WAIT、43 tokens、1,468 ms |
| 真实对战 #5 | 14 回合，28 次实际模型调用，28 成功、0 失败、无重试；蓝方自然获胜，红方 HP 0 / 蓝方 HP 25；没有确定性模拟或预设赢家 |
| 模型用量与耗时 | 供应商返回输入 28,480 / 输出 483 / 总计 28,963 tokens；比赛墙钟 30 秒（时间戳精度 1 秒）。累计双边请求耗时 41,484 ms，因并行不等于墙钟时长 |
| 证据 | schema 2；逐回合关联两次成功调用。canonical JSON Keccak-256 与提交上链的结果哈希一致 |
| 奖金 #5 | 创建、赞助、报名/退出/重报、开赛、提交、审核、领取；蓝方 DUSD 余额从 100 增至 200，实际到账 100 |
| 退款 #6 | 独立新赛事赞助 25 DUSD 后取消，赞助者余额从 775 增至 800，实际退回 25；没有调用第二场模型对战 |
| 钱包交易 | 本轮 14 笔成功回执，覆盖 10 种业务交易类型，总 Gas 1,924,746；最终托管代币余额与总负债均为 0 |
| 时间窗口 | 为缩短本地报名等待，明确记录一次 164 秒的 Anvil 时间快进；结果提交后完整等待 60 秒审核窗口，没有快进审核窗口 |
| 后端回归 | 全量 test / vet 通过；agent 与 tournament 的 race 通过，无 DATA RACE。新增 10 个顶层测试覆盖失败中断、取消、严格动作、重试/重启预算、脱敏和旧 schema 1 字节/哈希兼容 |
| 前端 | 30/30 单元测试、类型检查、lint 和生产构建通过；真实页面自动更新调用与回合，证据核验通过，已领奖与退款状态可查 |
| 手机与双语 | 390 × 844 下页面宽度 390，无整页横向溢出；英文调用统计可读，切回中文后选中的第 13 回合保持不变 |
| 留档 | 操作前后 SQLite 在线备份与 Anvil 状态快照留在 `.local/llm-demo/`；备份恢复演练未执行 |

证据：[连通性](evidence/llm-connectivity-2026-10-10.json)、[本地业务与余额](evidence/local-llm-integration-2026-10-10.json)、[完整比赛证据](evidence/local-llm-match-5-2026-10-10.json)、[后端测试](evidence/backend-llm-tests-2026-10-10.txt)、[race](evidence/backend-llm-race-2026-10-10.txt)、[前端测试](evidence/frontend-llm-tests-2026-10-10.txt)、[构建](evidence/frontend-llm-build-2026-10-10.txt)、[手机画面](evidence/llm-mobile-2026-10-10.jpg)。`go vet` 退出 0 且没有诊断输出；race 保留 macOS 链接器既有警告，前端保留 WalletConnect 可选 `pino-pretty` 警告。

全部署 Dashboard 现在有 6 场赛事，累计入金 450 DUSD、奖金 200、退款 250；其中包含旧演示，不能将它们全部归到本轮。本轮新增入金为 125、奖金 100、退款 25 DUSD。

本轮共 4 次诊断 + 28 次比赛模型请求。成功调用和本地测试资金闭环不代表套餐允许一般应用后端使用、不代表公网部署、真实资产付款、供应商签名证明或外部独立审计。云端继续按已确认的 Render 免费方案准备；课程公网 URL、测试网部署与新版英文配音实操视频尚未在本轮完成。

复现与代码阅读：[本地真实模型指南](LOCAL_LLM_DEMO.md)、[033 开发记录](../../devlog/033-local-llm-demo.md)。

## 2026-10-10：英文讲解与双语字幕视频（032）

已生成 [9 分 56 秒演示视频](DEMO_VIDEO.md)，含实际本地界面操作、已有交易记录回顾及明确标记的架构说明卡。本轮不改变业务代码，不重新运行合约/API/前端测试，不新增钱包签名或资金交易。

| 检查 | 实际结果 |
| --- | --- |
| 成片格式 | H.264，1920 × 1080，30 fps；AAC 48 kHz 单声道；596.400 秒，30,305,416 bytes |
| 解码完整性 | FFmpeg 对最终字幕版完整解码，退出码 0，无 error 输出 |
| 配音与字幕 | 10 章、50 组英文配音/双语文案逐条一致；中英、纯英文、纯中文 SRT 各 50 条，无重叠、无越过音轨结束时间 |
| 时间与音量 | 音轨 596.392 秒，音视频时长差 8 ms；音量检测 mean −15.9 dB / peak −1.5 dB |
| 画面抽查 | 查看封面和 14 个时点的画面，双语字幕位于独立底部区域，中文字体正常；对长字幕与退款场景按原尺寸复查 |
| 生成器 | Python 编译检查、章节复用条件检查、中文换行检查和 git diff --check 通过 |

来源：[媒体验证 JSON（含 SHA-256）](evidence/demo-video-verification-2026-10-10.json)、[成片封面](evidence/demo-video-cover-2026-10-10.png)、[逐句学习稿](demo/SCRIPT.en-zh.md)。视频由约 5 fps 的真实浏览器截图按实际时间合成，章节尾部保留画面以配合讲解；输出 30 fps 不等于原始采集帧率。抽帧、解码与音量检查没有替代真实观众逐秒听审。

讲解使用已有 4 场本地赛事、325 DUSD 入金、100 奖金、225 退款和 0 托管余额；#2 是预置合约争议测试案例。视频明确本地网络、确定性模拟和未验证项目，不将其记作新的 MetaMask 实操、真实 LLM 评价、公网部署或利益相关者验证。

## 2026-10-10：中英文界面（031）

新增全局 English / 中文 切换、浏览器偏好保存和完整课程界面翻译；没有增加 npm 依赖、修改合约/API 或重置运行数据。`npm test` **27/27 通过**（本轮新增 9 项），`npm run typecheck`、`npm run lint` 和最终生产构建通过，`git diff --check` 通过。

原始输出：[单元测试](evidence/frontend-i18n-tests-2026-10-10.txt)、[生产构建](evidence/frontend-i18n-build-2026-10-10.txt)。构建仍有既有 WalletConnect 可选 `pino-pretty` warning，不影响本次构建。

构建连接 API `18080` / RPC `18545`，只替换 `3001` 前端；原 API PID 6178 和 Anvil PID 6175 保留。重启前只读快照已经是 **4 场赛事、325 DUSD 入金、100 DUSD 奖金、225 DUSD 退款、0 DUSD 托管余额**，#3/#4 已取消；不能把这个快照与 030 的早期三场数据混为一谈。本轮未发送链上交易。

| 验证范围 | 实际结果 |
| --- | --- |
| 语言覆盖 | 浏览器实际查看首页、赛事列表、#1 详情与回放、审计、未登录 Agent 页、指南；英文平台文案显示正常，原始中文赛事名保留 |
| 偏好保存 | 点击 English 后跳转赛事页、重新加载页面，仍显示英文；文档 lang/title/description 同步英文 |
| 钱包弹窗 | 在英文模式打开 RainbowKit，实际显示 “Connect a Wallet” 和英文说明；随后关闭，没有选择钱包或签名 |
| 筛选输入保留 | 审计选 #3、输入 `refund`，切回中文后赛事和输入不变，仍匹配一条 200 DUSD 退款记录；中英文搜索词均作为别名索引 |
| 回放状态保留 | 手机详情页从第 12 回合退到第 11 回合，切回中文后滑块 value 仍为 10，没有重置 |
| 手机 390 × 844 | 英文详情与中英文首页的 document scrollWidth / clientWidth 均为 390；切换入口可见，移动导航可跳转首页；长按钮允许换行 |
| 键盘与诊断 | Enter 可触发语言切换；本轮测试标签页未收集到 console error；测试后恢复正常视口 |
| 语义与整数安全 | 40 个英文赛事场景覆盖期限/零奖池/未知链时间/终态，中文与英文 tone 一致且输入不变；默认普通用户不误标为参赛者；测试保留退款顺序、大整数和原始钱包诊断 |

截图：[英文桌面首页](evidence/frontend-english-home-2026-10-10.jpg)、[英文手机详情](evidence/frontend-english-mobile-2026-10-10.jpg)。浏览器为 Codex 内置浏览器，桌面实际 1280 × 720；移动使用临时 390 × 844 视口。

**本轮未执行：** 真实钱包登录/签名、已登录各角色完整英文流程、在待签名交易期间切换语言、禁用 localStorage 的真实浏览器场景。Provider 与交易代码已审查为不按语言重建、签名参数不受翻译影响，但代码审查不能代替这部分人工验收。MetaMask 扩展自身语言不由平台控制。已准备交易的说明保留估算时语言；未知 API/钱包诊断、用户内容和原始证据不自动翻译。

资料：[031 开发日志](../../devlog/031-bilingual-interface.md)、[用户手册](USER_GUIDE.md#中英文切换)。

## 2026-10-10：前端可用性改进（030）

本轮只调整前端及文档，继续使用原本的 3 场演示赛事，不重置数据。最终 `npm test` **18/18 通过**（原 6 项 + 12 项展示状态/事件身份测试），`npm run typecheck`、`npm run lint` 均通过，生产构建成功。[构建原始输出](evidence/frontend-ux-build-2026-10-10.txt)仍有 WalletConnect 依赖对可选 `pino-pretty` 的既有 warning，没有新增依赖来掩盖它。

构建参数：`NEXT_PUBLIC_API_URL=http://localhost:18080 NEXT_PUBLIC_LOCAL_RPC_URL=http://127.0.0.1:18545 npm run build`。新 standalone 服务在 `127.0.0.1:3001` 运行；原 Anvil `18545`、API `18080` 保留，`/ready` 返回 `ready: true`。

| 验证范围 | 实际检查结果 |
| --- | --- |
| 首页桌面 1280 × 900 | 清晰显示平台用途、真实 #1 比赛入口、模拟标识和 325 / 100 / 25 / 200 DUSD，见[截图](evidence/frontend-home-2026-10-10.jpg) |
| 赛事列表 | #3 从“开放报名”改为“已到期 · 待退款”，#1 为“已领奖”，#2 为“已取消” |
| 回放 | 点击“播放回放”从最后回合返回首回合，控件变“暂停”，实际推进至第 12 回合；“跟随最新回合”可用 |
| 手机 390 × 844 | 主导航展开可用，列表创建表单显示所需主办方权限；详情页 document 宽度与 viewport 均为 390，没有整页横向溢出；棋盘和流水保留局部横向滚动，见[截图](evidence/frontend-mobile-2026-10-10.jpg) |
| 匿名 Agent 入口 / 最终流水 | “先连接钱包，再签名登录”及连接按钮可见；最后构建中 #1“审核通过”行实际显示 `0x3c44…93bc`，与指定审核员一致 |
| Chrome 已连接普通钱包 | 实际显示“观赛者”、0 个人赞助和可见的到期取消入口，没有把普通钱包标成参赛者 |
| 精确边界 / 大整数 / 事件身份 | 新测试覆盖最终期限、审核窗口、未知链时间、零奖池、终态优先、大额金额、审核员与冠军地址区分、缺失操作钱包不猜测 |

**未计作通过的项目：** 本次没有执行 MetaMask 资金签名；Chrome 原生自动检查交易弹窗时工具超时，因此弹窗实际开启、关闭恢复、拒签、替换交易和全部角色完整流程仍待人工复测。不能以静态代码检查代替这部分验收。到期取消和退款不是本次已执行操作，原 #3 的 200 DUSD 仍按用户操作前状态保存。合约/后端本轮未改动，没有重复运行旧全链路脚本。

资料：[030 开发日志](../../devlog/030-frontend-usability.md)、[更新后的体验指南](EXPERIENCE_GUIDE.md)。

## 2026-10-09：原课程版验证概览

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
- 真实利益相关者反馈与作业报告、Slides、个人贡献的最终版本。视频讲解版已于 032 产出；新增钱包签名的现场演示仍未录制。

CI 工作流已经配置，但本次没有推送 GitHub，不能记录远程 CI 已通过。
