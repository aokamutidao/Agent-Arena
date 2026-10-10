# Agent Arena 演示旁白（英中对照）

英文配音、中英文双语字幕；10 章、50 组 cue。本文与 [NARRATION.en-zh.json](NARRATION.en-zh.json) 使用同一最终文案，便于复习、逐段配音或修改。

本次查看已有本地测试赛事，不执行新的钱包签名、创建或资金交易。DUSD 为测试代币，对战为确定性模拟；#2 是预设合约争议 fixture，#3 的 200 DUSD 退款是已有确认记录。画面与实际验收边界见 [演示视频说明](../DEMO_VIDEO.md)。

时间来自本轮 timed-narration.json：Samantha 英文合成语音，语速参数 145，每组后留 0.25 秒；音频总时长 596.392 秒（09:56.392）。字幕 end 不含其后停顿，章节 end 包含停顿；最终视频时长与音画/字幕验收需另行确认。

## 01 · 用途与本地演示边界 / Purpose and local demo

章节时间：00:00.000–01:00.717。

### 01-001 · 00:00.000–00:10.031

**EN:** Welcome to Agent Arena. This edited walkthrough inspects a local deployment and existing test tournaments; it does not perform new wallet signatures or payments.

**ZH:** 欢迎观看 Agent Arena。这段剪辑演示查看本地部署与已有测试赛事，本次不进行新的钱包签名或付款。

### 01-002 · 00:10.281–00:23.818

**EN:** The platform manages sponsored prizes for AI agent competitions. Organizers, developers, sponsors, and independent auditors can inspect the prize pool, results, and financial history.

**ZH:** 平台管理 AI Agent 比赛的赞助奖金，主办方、开发者、赞助者和独立审核员可以核查奖池、结果及资金历史。

### 01-003 · 00:24.068–00:35.400

**EN:** The workflow is visible here: sponsors deposit funds, two agents compete, an independent auditor reviews the result, and eligible wallets claim prizes or refunds.

**ZH:** 流程在这里清晰可见：赞助者预存资金，两个 Agent 对战，独立审核员复核结果，再由符合条件的钱包领取奖金或退款。

### 01-004 · 00:35.650–00:48.049

**EN:** This demonstration uses Anvil and DUSD test tokens, with deterministic agent simulation. DUSD has no promised monetary value, and this is not a public-cloud deployment.

**ZH:** 演示使用 Anvil、DUSD 测试代币与确定性 Agent 模拟。DUSD 不承诺现实货币价值，当前也不是公有云部署。

### 01-005 · 00:48.299–01:00.467

**EN:** The current overview shows four tournaments: three hundred twenty-five DUSD deposited, one hundred paid as prizes, two hundred twenty-five refunded, and zero remaining in escrow.

**ZH:** 当前概览显示四场赛事：累计入金 325 DUSD，已发奖金 100，已退赞助 225，托管余额为零。

## 02 · 用户、角色与指南 / Users, roles and guide

章节时间：01:00.717–02:00.876。

### 02-001 · 01:00.717–01:14.672

**EN:** The guide separates four participation paths: sponsoring, competing, organizing, and auditing. Viewing tournaments needs no wallet, while protected actions require the appropriate identity and permissions.

**ZH:** 指南分为赞助、参赛、主办和审核四条参与路径。查看赛事无需钱包，受保护操作则需要相应身份与权限。

### 02-002 · 01:14.922–01:27.577

**EN:** Administrator, organizer, and auditor are on-chain permission roles. Participant and sponsor are identities within one tournament, rather than global privileges automatically granted to everyone.

**ZH:** 管理员、主办方和审核员是链上权限角色；参赛者与赞助者是某场赛事中的身份，不是自动授予所有人的全局权限。

### 02-003 · 01:27.827–01:37.092

**EN:** An ordinary wallet can create an agent or sponsor an eligible tournament. Being labelled an ordinary user does not mean that wallet has already registered.

**ZH:** 普通钱包可以创建 Agent，也可按规则赞助赛事。标为普通用户，并不表示这个钱包已经报名。

### 02-004 · 01:37.342–01:50.298

**EN:** The designated auditor must remain separate from that tournament's organizer, participants, and sponsors. This reduces conflicting wallet roles, although different wallets do not prove independent people.

**ZH:** 指定审核员必须与本场主办方、参赛者和赞助者分离。这降低钱包角色冲突，但不同钱包并不能证明现实身份独立。

### 02-005 · 01:50.548–02:00.626

**EN:** These restrictions are enforced by smart contracts, not by hiding buttons. The guide also explains when to wait, when to review, and who claims funds.

**ZH:** 限制由智能合约执行，而不只是隐藏按钮。指南还说明什么时候等待、什么时候审核，以及由谁领取资金。

## 03 · 钱包入口与创建权限 / Wallet connection and creation permissions

章节时间：02:00.876–02:54.116。

### 03-001 · 02:00.876–02:10.698

**EN:** The wallet connector opens the connection interface. I close it without selecting a wallet or signing, so this recording makes no claim of a fresh transaction.

**ZH:** 钱包入口可以打开连接界面。我关闭它，不选择钱包或签名，因此本次录像不声称发送了新的交易。

### 03-002 · 02:10.948–02:22.674

**EN:** Wallet login uses a server-generated, expiring nonce bound to the address, domain, and chain. Signing in identifies the user without transferring tokens or paying gas.

**ZH:** 钱包登录使用服务端生成的限时 nonce，绑定地址、域名和链。登录签名用于确认身份，不转移代币，也不支付 Gas。

### 03-003 · 02:22.924–02:32.212

**EN:** Switching the wallet account or network clears the login session. Users must sign in again and check the displayed address before performing another role's actions.

**ZH:** 切换钱包账户或网络会清除登录会话。执行另一种身份的操作前，用户必须重新登录并核对页面地址。

### 03-004 · 02:32.462–02:42.586

**EN:** On the tournaments page, Create tournament reveals the entry point. Because this browser is not signed in as an organizer, it displays the permission requirement.

**ZH:** 赛事页的 Create tournament 展开创建入口。当前浏览器未以主办方身份登录，因此显示相应权限要求。

### 03-005 · 02:42.836–02:53.866

**EN:** An authorized organizer would set the name, auditor, and deadlines, then preview the transaction fee. The contract still checks authorization and timing when the wallet sends.

**ZH:** 有权限的主办方可设置名称、审核员与期限，再预览交易费用。钱包发送时，合约仍会检查授权与时间条件。

## 04 · 报名及 Agent 定义 / Registration and committed agent definitions

章节时间：02:54.116–03:52.348。

### 04-001 · 02:54.116–03:06.144

**EN:** This completed tournament records two participants and their registered agent definitions. Registration was free; the prize came from sponsorship, rather than admission fees or spectator betting.

**ZH:** 这场已完成赛事记录了两名参赛者及报名 Agent 定义。报名免费，奖金来自赞助，而非入场费或观众下注。

### 04-002 · 03:06.394–03:17.028

**EN:** An agent profile contains a name, strategy, owner, and content hash. It is saved off-chain, while registration commits that hash to the tournament on-chain.

**ZH:** Agent 定义包含名称、策略、所有者与内容哈希。定义保存在链下，报名则将该哈希承诺到链上赛事。

### 04-003 · 03:17.278–03:27.936

**EN:** The definitions are immutable. Changing a strategy requires creating another profile; the match runner checks ownership and includes the registered definition in its evidence.

**ZH:** 定义不可变。调整策略需要另建定义；比赛运行器校验所有权，并把已报名定义纳入证据。

### 04-004 · 03:28.186–03:40.075

**EN:** The rules panel shows registration and final settlement deadlines. Review timing is recorded with submitted results. Eligibility follows blockchain time, not the computer clock.

**ZH:** 规则区展示报名及最终结算期限，已提交结果还记录审核时间。操作资格依据链时间判断，而非电脑时钟。

### 04-005 · 03:40.325–03:52.098

**EN:** Starting requires closed registration, two participants, a funded pool, and enough review time. Approval and refund actions also depend on the recorded state and deadlines.

**ZH:** 开赛要求报名已截止、满两人、奖池有资金且审核时间充足。批准与退款操作也取决于记录状态及期限。

## 05 · 已有比赛回放 / Recorded match replay

章节时间：03:52.348–04:47.097。

### 05-001 · 03:52.348–04:02.472

**EN:** The replay displays this tournament's recorded match, rather than a newly launched game. The available history contains twelve turns, with actions and health changes.

**ZH:** 回放展示该赛事已记录的比赛，而不是刚运行的新对局。历史包含十二个回合，以及动作与生命值变化。

### 05-002 · 04:02.722–04:13.403

**EN:** I can play, pause, step between turns, or drag the timeline. These controls inspect persisted data and do not submit transactions or alter the result.

**ZH:** 可以播放、暂停、切换回合或拖动时间轴。这些控件只查看持久化数据，不发送交易，也不修改结果。

### 05-003 · 04:13.653–04:25.681

**EN:** The mode is explicitly identified as deterministic simulation, not a real language-model match. Aggressive and defensive strategies provide reproducible baselines for this local scenario.

**ZH:** 模式明确标为确定性模拟，而非真实大模型对局。进攻型与防守型策略为本地场景提供可复现的基线。

### 05-004 · 04:25.931–04:37.262

**EN:** A configured model service can support LLM decisions, but that capability is distinct from this recorded demonstration. If neither mode is available, execution is rejected.

**ZH:** 配置模型服务后可以支持 LLM 决策，但这项能力与本次记录演示不同。两种模式均不可用时，执行会被拒绝。

### 05-005 · 04:37.512–04:46.847

**EN:** The backend runs the game and persists its turns. It does not control the prize wallet, and a completed match cannot simply be rerun to overwrite evidence.

**ZH:** 后端运行游戏并保存回合，但不控制奖金钱包。已经完成的比赛也不能通过直接重跑覆盖证据。

## 06 · 证据与哈希核验 / Evidence and hash verification

章节时间：04:47.097–05:41.219。

### 06-001 · 04:47.097–04:59.636

**EN:** The evidence panel exposes the canonical match record. It includes the tournament context, registered agents, participants, actions, and result needed for later inspection and review.

**ZH:** 证据区展示规范化比赛记录，包含赛事上下文、报名 Agent、参赛者、动作与结果，供后续核查及审核。

### 06-002 · 04:59.886–05:08.663

**EN:** The browser recomputes a Keccak hash of the canonical evidence. It compares that value with the returned record and the result commitment stored on-chain.

**ZH:** 浏览器重新计算规范证据的 Keccak 哈希，并将其与返回记录及链上保存的结果承诺进行比较。

### 06-003 · 05:08.913–05:19.896

**EN:** Downloading the JSON preserves the underlying record for independent checking. The hash binds the submitted result to specific content, without storing the entire match on-chain.

**ZH:** 下载 JSON 可保留原始记录供独立核查。哈希把提交结果绑定到特定内容，无需将整场比赛永久写入链上。

### 06-004 · 05:20.146–05:29.109

**EN:** A matching hash proves consistency with the committed record. It does not prove that the original information was truthful or that every decision was fair.

**ZH:** 哈希一致证明内容与已承诺记录一致，但不证明原始信息真实，也不证明每个决策公平。

### 06-005 · 05:29.359–05:40.969

**EN:** Independent auditing remains necessary. Original evidence also requires backups: an on-chain hash can detect a changed file, but it cannot reconstruct a file that was lost.

**ZH:** 因此仍需要独立审核。原始证据也必须备份：链上哈希能发现文件变化，却不能恢复已丢失的文件。

## 07 · 审核及领奖记录 / Approval and prize payment records

章节时间：05:41.219–06:40.380。

### 07-001 · 05:41.219–05:53.015

**EN:** The audit history separates submission, approval, and payment. The organizer proposed a winner and evidence hash; approval followed the review period under the designated auditor's authority.

**ZH:** 审计历史区分提交、批准与付款。主办方提出冠军和证据哈希，指定审核员在审核窗口后行使权限批准。

### 07-002 · 05:53.265–06:03.319

**EN:** Participants or actual sponsors could dispute a submitted result during its review window. The organizer could not bypass that window and finalize the prize alone.

**ZH:** 参赛者或实际赞助者可以在结果审核窗口内提出异议。主办方不能绕过窗口，独自完成奖金结算。

### 07-003 · 06:03.569–06:14.459

**EN:** Approval determined the winner, but payment required a separate claim. This tournament's PrizeClaimed record shows that the winner received one hundred DUSD in the test scenario.

**ZH:** 批准确定冠军，但付款仍需要单独领取。本场 PrizeClaimed 记录显示，冠军在测试场景中领取了 100 DUSD。

### 07-004 · 06:14.709–06:26.250

**EN:** This pull-payment design lets each eligible wallet claim its own funds. Duplicate prize claims are rejected, and neither the organizer nor administrator can arbitrarily drain escrow.

**ZH:** 领取式付款让每个有资格的钱包自行领取资金。重复领奖会被拒绝，主办方与管理员也不能任意提空托管余额。

### 07-005 · 06:26.500–06:40.130

**EN:** These are previously confirmed events, not transactions sent during this recording. Event timestamps, responsible wallet addresses, and transaction identifiers make the approval and payment distinguishable.

**ZH:** 这些是此前已确认的事件，不是录像时发送的交易。事件时间、操作钱包及交易标识，使批准与付款可以分别核查。

## 08 · 取消及赞助退款 / Cancellation and sponsor refunds

章节时间：06:40.380–07:42.745。

### 08-001 · 06:40.380–06:51.131

**EN:** Tournament three illustrates the refund path. Its earlier two-hundred-DUSD sponsorship has already been refunded, and the confirmed event can be inspected without connecting a wallet.

**ZH:** 第三场赛事展示退款路径。此前的 200 DUSD 赞助已经退回，已确认事件无需连接钱包即可查看。

### 08-002 · 06:51.381–07:01.598

**EN:** Expiry alone does not transfer money. Someone first triggers the permitted cancellation, then the original sponsor claims their own refund using the wallet that funded it.

**ZH:** 仅到期不会转移资金。需要先有人按规则触发取消，再由原赞助者用出资钱包领取自己的退款。

### 08-003 · 07:01.848–07:16.871

**EN:** Underfilled registration after its deadline, or an expired final settlement deadline, permits timeout cancellation without an organizer's cooperation. Sponsors are accounted for separately, and duplicate refunds are rejected.

**ZH:** 报名截止后不足两人，或最终结算期限已到，可不依赖主办方配合而触发超时取消。赞助者独立记账，重复退款会被拒绝。

### 08-004 · 07:17.121–07:29.404

**EN:** Tournament two is a prepared contract dispute fixture, not a real customer complaint. It demonstrates recorded dispute and decision hashes, cancellation, and a twenty-five-DUSD refund.

**ZH:** 第二场是预设合约争议测试场景，不是真实客户投诉。它展示异议与裁决哈希、取消以及 25 DUSD 退款。

### 08-005 · 07:29.654–07:42.495

**EN:** Together, the two refunds explain the overview's two-hundred-twenty-five-DUSD total. Saved explanations are only drafts until linked to confirmed on-chain actions; the interface makes that difference explicit.

**ZH:** 两笔退款合计解释概览中的 225 DUSD。保存说明在关联已确认链上操作前只是草稿，界面明确展示这一区别。

## 09 · 架构与测试证据 / Architecture and verification

章节时间：07:42.745–08:50.195。

### 09-001 · 07:42.745–07:56.631

**EN:** Next.js and TypeScript provide the interface, while wagmi integrates wallet transactions. The Go REST backend and SQLite persist identities, agent definitions, match evidence, and indexed events.

**ZH:** Next.js 与 TypeScript 提供界面，wagmi 集成钱包交易。Go REST 后端和 SQLite 持久化身份、Agent 定义、比赛证据及索引事件。

### 09-002 · 07:56.881–08:09.652

**EN:** TournamentManager coordinates TournamentRegistry, PrizeEscrow, and ResultReview alongside the test token. The backend only reads blockchain state and logs; it holds no funds-signing private key.

**ZH:** TournamentManager 协调 TournamentRegistry、PrizeEscrow、ResultReview 与测试代币。后端只读链状态和日志，不持有资金签名私钥。

### 09-003 · 08:09.902–08:21.604

**EN:** Amounts use integer token units. Confirmed events have unique identities and are indexed transactionally, allowing projections to restart without counting the same financial event twice.

**ZH:** 金额采用整数最小单位。已确认事件具有唯一标识并事务化索引，使投影重启后不会重复累计同一笔财务事件。

### 09-004 · 08:21.854–08:35.345

**EN:** Recorded verification includes forty-nine passing contract tests, twenty-seven frontend tests, backend race checks, and twenty-three successful signed transactions in the separate local integration script.

**ZH:** 已有验证记录包括 49 项合约测试、27 项前端测试、后端 race 检查，以及独立本地集成脚本中的 23 笔成功签名交易。

### 09-005 · 08:35.595–08:49.945

**EN:** Security scenarios cover unauthorized calls, replay, reentrancy, deadlines, duplicate claims, and fund conservation. Controlled gas samples also show withdrawals avoid looping through all sponsors.

**ZH:** 安全场景覆盖越权、重放、重入、期限、重复领取与资金守恒。受控 Gas 样本也表明，领取操作不遍历全部赞助者。

## 10 · 双语体验、限制与结束 / Usability, limits and closing

章节时间：08:50.195–09:56.392。

### 10-001 · 08:50.195–09:02.084

**EN:** The responsive interface supports English and Chinese. Language changes preserve user content and page state; mobile navigation and replay controls remain available on a narrow screen.

**ZH:** 响应式界面支持中英文切换。切换保留用户内容与页面状态；窄屏仍可使用移动导航和回放控件。

### 10-002 · 09:02.334–09:16.800

**EN:** The local workflow and tests demonstrate functioning software, not production certification. Public hosting, Sepolia deployment, full MetaMask acceptance, and real model evaluation still need separate evidence.

**ZH:** 本地流程与测试证明软件能够运行，不等于生产认证。公网运行、Sepolia 部署、完整 MetaMask 验收及真实模型评价仍需独立证据。

### 10-003 · 09:17.050–09:28.730

**EN:** Real stakeholder interviews and measured business improvements are still pending. Local API timings describe this machine only, and cannot be presented as cloud service-level guarantees.

**ZH:** 真实利益相关者访谈及业务改善测量仍待完成。本机 API 耗时只描述这台机器的观测，不能作为云服务水平保证。

### 10-004 · 09:28.980–09:40.659

**EN:** Blockchain is useful when parties lack a shared trusted operator. Where that trust already exists, conventional escrow may be simpler; privacy and auditor collusion remain concerns.

**ZH:** 当多方缺乏共同可信运营者时，区块链具有价值。已有共同信任时，传统托管可能更简单；隐私与审核员串谋仍是问题。

### 10-005 · 09:40.909–09:56.142

**EN:** Agent Arena demonstrates transparent funding, reviewable evidence, controlled prize claims, and recoverable sponsor funds. This walkthrough showed existing records honestly, with explicit local-demo boundaries and limitations.

**ZH:** Agent Arena 展示透明入金、可复核证据、受控领奖及赞助资金退出机制。本次演示如实查看已有记录，并明确本地演示边界与限制。
