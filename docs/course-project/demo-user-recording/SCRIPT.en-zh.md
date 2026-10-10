# Agent Arena 实操演示双语讲稿

英文合成配音 · 中英文字幕 · 赛事 #13。

## 00:00 · Agent Arena

Welcome to Agent Arena: an AI tournament platform with transparent prize escrow and independent result review.

欢迎来到 Agent Arena：一个支持透明奖金托管和独立结果审核的 AI 赛事平台。

Let's follow a complete tournament, from preparation to payment.

接下来，我们将完整展示一场赛事，从筹备一直到奖金到账。

## 00:12 · A fast local demonstration

This demonstration runs locally for fast responses and a smooth, complete walkthrough.

本次演示使用本地环境，响应更快，便于流畅展示完整流程。

MetaMask signs transactions on the local blockchain, and DUSD is our demonstration token.

MetaMask 在本地区块链上签署交易，DUSD 是本次使用的演示代币。

The homepage introduces the prize pool, Agent competition, independent review, and final claim.

首页清楚呈现奖池注资、Agent 对战、独立审核和最终领奖四个环节。

## 00:32 · Clear responsibilities

Small tournament organizers need clear rules and a reliable way to distribute sponsored prizes.

小型赛事主办方需要明确的规则，以及可靠的赞助奖金分配方式。

Here, organizers set the rules, developers enter Agents, sponsors fund the pool, and an auditor reviews the result.

在这里，主办方设定规则，开发者提交 Agent，赞助方注资，审核员审查结果。

## 00:48 · Follow the funds

The dashboard separates sponsorship received, prizes paid, refunds, and funds still held in escrow.

资金看板分别展示赞助入金、已发奖金、已退赞助和仍在托管中的资金。

Now let's create and run tournament thirteen. Account switching and idle waits are shortened in this edit.

现在开始创建并运行赛事 13。剪辑中已缩短账户切换和空闲等待。

## 01:02 · 01 / Create the tournament

The organizer signs in with MetaMask and opens the tournament creation form.

主办方通过 MetaMask 签名登录，打开创建赛事表单。

We enter a title and designate the auditor who will check this competition's result.

填写赛事名称，并指定负责检查本场结果的审核员。

The form sets the registration deadline, final settlement deadline, and result review window.

表单设置报名截止时间、最终结算期限以及结果审核窗口。

For this match, registration lasts four minutes, settlement allows twenty-four hours, and review lasts sixty seconds.

本场报名窗口为 4 分钟，结算期限为 24 小时，审核窗口为 60 秒。

The transaction preview displays the action and estimated gas before MetaMask confirmation.

交易预览会在 MetaMask 确认前展示操作内容与预计 Gas。

After confirmation, tournament thirteen is created with its rules and auditor recorded on chain.

确认后，赛事 13 创建成功，规则与审核员记录在链上。

## 01:40 · 02 / Prepare the red Agent

Next, the red participant opens My Agents and selects the saved MetaMask Demo Red profile.

接着，红方参赛者打开“我的 Agent”，选择已保存的 MetaMask Demo Red。

An Agent profile contains its name, character, and instructions for deciding how to play.

Agent 定义包含名称、性格，以及决定如何行动的策略指令。

This is where a developer expresses a strategy: when to approach, attack, use a skill, or heal.

开发者可以在这里描述策略：何时接近、攻击、使用技能或治疗。

The saved definition has a content hash, which links this exact Agent to its tournament entry.

保存的定义具有内容哈希，将这个确定的 Agent 与赛事报名绑定。

## 02:08 · 03 / Register the red participant

The participant follows the tournament link, chooses the red Agent, and previews registration.

参赛者进入赛事页面，选择红方 Agent，并预览报名交易。

Registration is free. The wallet signs the entry, including the selected Agent definition hash.

报名不收取费用。钱包签署报名交易，其中包含所选 Agent 定义的哈希。

The confirmed participant card now shows the red wallet and its registered Agent.

确认后，参赛者卡片显示红方钱包及其已报名的 Agent。

## 02:29 · 04 / Register the blue participant

The blue participant signs in with a separate wallet and selects MetaMask Demo Blue.

蓝方参赛者使用另一个钱包签名登录，选择 MetaMask Demo Blue。

The tournament page gives each role a clear next action as the competition progresses.

随着赛事推进，详情页会为不同身份显示清楚的下一步操作。

The blue participant reviews and signs the second registration in MetaMask.

蓝方参赛者在 MetaMask 中查看并签署第二笔报名交易。

Both participants are registered, with their Agent definitions linked to their wallets.

双方报名完成，各自的 Agent 定义均已与钱包绑定。

## 02:52 · 05 / Sponsor 100 DUSD

The sponsor signs in and chooses to contribute one hundred DUSD to this tournament.

赞助方签名登录，选择向本场赛事赞助 100 DUSD。

There are two wallet transactions. First, approve the escrow contract to use exactly one hundred tokens.

这里有两笔钱包交易：首先授权托管合约使用精确的 100 枚代币。

Second, confirm the sponsorship deposit, which transfers the tokens into the tournament's escrow.

然后确认赞助入金，将代币转入本场赛事的托管奖池。

The confirmed prize pool now holds one hundred DUSD.

确认后，奖池中现已托管 100 DUSD。

Funding is now in place before the competition begins.

比赛开始前，奖金已实际落实。

## 03:26 · 06 / Start the tournament

Once registration closes, the organizer returns to a ready-to-start tournament.

报名窗口结束后，主办方回到已准备开赛的赛事。

Both participant slots are filled and the sponsored prize is held in escrow.

两个参赛名额已满，赞助奖金保存在托管合约中。

The organizer confirms Start match in MetaMask, moving the tournament into its active state.

主办方在 MetaMask 中确认开赛，将赛事推进到进行中状态。

## 03:46 · 07 / Launch the real AI match

Run Agent match now starts the backend battle engine using the two registered Agent definitions.

点击运行 Agent 对战，后端引擎开始使用双方已报名的 Agent 定义执行比赛。

This match uses the Qwen three point five Plus language model for decisions.

本场比赛使用 qwen3.5-plus 大模型生成决策。

## 04:00 · 08 / Watch the battle unfold

Each round, the Agents receive the current arena state and choose their next actions.

每一回合，Agent 都会接收当前竞技场状态，并选择下一步动作。

The game engine checks those actions against movement, attack, healing, and obstacle rules.

游戏引擎依据移动、攻击、治疗和障碍物规则检查这些动作。

As the battle develops, the board shows positions, while the health bars show the effect of each exchange.

随着战斗推进，棋盘显示双方位置，生命条显示每轮交锋的影响。

The Agents close the distance and trade attacks. Their decisions are recorded alongside the resulting game state.

双方逐步缩短距离并交互攻击。每次决策都会与执行后的游戏状态一起保存。

At low health, healing becomes part of the strategy. The outcome follows the executed actions and game rules.

生命值较低时，治疗也成为策略的一部分。结果由实际执行的动作和游戏规则决定。

The live display keeps updating until the engine reaches the final result.

实时画面持续更新，直到引擎判定最终结果。

## 04:40 · 09 / Inspect the result

After fifteen rounds, the blue Agent wins with ten health points remaining.

经过 15 回合，蓝方以 10 点剩余生命值获胜。

The action history shows each round, and the evidence view brings together the Agents, actions, and final result.

动作历史展示每一回合，证据视图汇总双方 Agent、动作与最终结果。

The complete evidence can also be downloaded as JSON.

完整证据还可以下载为 JSON 文件。

## 05:00 · 10 / Replay the match

The replay lets us revisit the battle after it has finished.

比赛结束后，可以通过回放重新查看对局。

Play the recorded sequence or move between rounds to inspect positions, health, actions, and decision records.

播放已记录的过程，或逐回合检查位置、生命值、动作与决策记录。

This gives participants and the auditor a concrete record to examine before the result is approved.

这样，参赛者和审核员在结果确认之前，就有具体的对局记录可以检查。

## 05:22 · 11 / Submit the result on chain

The organizer reviews the evidence, then submits the blue winner and the evidence hash.

主办方查看证据后，提交蓝方赢家地址和证据哈希。

The wallet preview makes the proposed winner and commitment visible before signing.

钱包签名前，交易预览清楚展示拟定赢家和提交的证据承诺。

After confirmation, the tournament enters its result review window.

交易确认后，赛事进入结果审核窗口。

The prize remains in escrow while the designated auditor checks the result.

指定审核员检查结果期间，奖金继续保存在托管合约中。

## 05:47 · 12 / Independent audit

We now sign in with the auditor's wallet, which has a separate responsibility from the organizer.

现在使用审核员钱包签名登录，其职责与主办方分开。

The auditor inspects the submitted winner, the match record, and the evidence associated with this tournament.

审核员检查已提交的赢家、对局记录，以及本场对应的证据。

The evidence hash connects this stored content to the commitment submitted on chain.

证据哈希将保存的内容与链上提交的承诺连接起来。

The review step gives a dedicated role control over approving the settlement result.

审核步骤由专门的角色负责，确认用于奖金结算的比赛结果。

## 06:13 · 13 / Approve the settlement

After checking the evidence and reaching the review deadline, the auditor approves the result through MetaMask.

检查证据并到达审核截止时间后，审核员通过 MetaMask 批准结果。

The tournament is finalized, and the blue winner can now claim the sponsored prize.

赛事结果确认完成，蓝方冠军现在可以领取赞助奖金。

## 06:31 · 14 / Claim the prize

The blue participant reconnects and signs in as the approved winner.

蓝方参赛者重新连接并签名登录，此时身份为已确认的冠军。

The page displays a claimable prize of one hundred DUSD.

页面显示可以领取的奖金为 100 DUSD。

The winner previews Claim prize and confirms the transfer in MetaMask.

冠军预览领取奖金操作，并在 MetaMask 中确认转账。

The escrow contract pays the eligible wallet and records that this tournament's prize has been claimed.

托管合约向符合条件的钱包付款，并记录本场奖金已经领取。

This completes the final financial step of the tournament.

赛事的最后一个资金环节由此完成。

## 06:59 · 15 / Verify the payout

The winner's displayed balance has increased from three hundred to four hundred DUSD.

冠军页面显示的余额从 300 DUSD 增加到 400 DUSD。

This tournament's escrow is now zero, and its status reads Prize claimed.

本场赛事的托管余额归零，状态显示为奖金已领取。

## 07:14 · 16 / One traceable financial history

The audit page ties the workflow together: creation, registrations, sponsorship, match start, result submission, approval, and prize claim.

审计页串联起创建、报名、赞助、开赛、结果提交、审核与领奖的完整流程。

Each confirmed event links the action to its transaction record.

每个已确认事件都将业务操作与交易记录对应起来。

## 07:30 · 17 / Cancellation and refunds

The platform also supports cancellation and sponsor refunds. Here we are viewing a historical refund example.

平台也支持取消赛事和退还赞助。这里展示的是一条历史退款记录。

Tournament six records a twenty-five DUSD deposit, cancellation, and refund to the sponsor.

赛事 6 记录了 25 DUSD 的入金、赛事取消以及退回赞助方的过程。

This section only displays the saved history.

这一段只展示已保存的历史记录。

## 07:48 · 18 / Full-stack architecture

The workflow connects four layers: a Next.js frontend, a Go backend, SQLite storage, and Solidity smart contracts.

完整流程连接四个层次：Next.js 前端、Go 后端、SQLite 存储和 Solidity 智能合约。

The interface uses wagmi and RainbowKit to connect MetaMask and guide each transaction.

界面通过 wagmi 和 RainbowKit 连接 MetaMask，引导用户完成每笔交易。

Wallets sign blockchain operations, while REST APIs provide the profiles, match execution, and dashboard data.

钱包负责签署链上操作，REST API 提供 Agent 资料、对战执行和看板数据。

## 08:12 · 19 / Interacting smart contracts

Tournament Manager coordinates three dedicated smart contract modules.

Tournament Manager 协调三个专门的智能合约模块。

Tournament Registry stores registration, deadlines, and lifecycle state.

Tournament Registry 保存报名信息、期限与赛事状态。

Prize Escrow holds sponsored tokens and handles prize claims and refunds for each tournament.

Prize Escrow 按赛事托管赞助代币，并处理领奖和退款。

Result Review manages submitted evidence, review windows, and dispute decisions. DUSD provides the ERC-twenty token used for settlement.

Result Review 管理提交证据、审核窗口与争议裁决；DUSD 提供用于结算的 ERC-20 代币。

## 08:38 · 20 / Backend and persistent records

The Go service uses Gin for its REST endpoints and SQLite for persistent application records.

Go 服务使用 Gin 提供 REST 接口，使用 SQLite 持久化应用记录。

Agent definitions, round histories, match evidence, and review notes are stored for later inspection.

Agent 定义、回合历史、比赛证据和审核记录都会保存，便于后续查看。

A blockchain reader indexes confirmed events into the database, making financial history available to the dashboard.

链上读取器将已确认事件索引到数据库，为资金看板提供可查询的历史记录。

## 09:00 · 21 / From Agent decisions to evidence

The Agent chooses an action from the current game state, and the engine applies deterministic game rules.

Agent 根据当前游戏状态选择动作，引擎依据确定的游戏规则执行。

Saved rounds power both the replay and the complete evidence record.

保存的回合数据同时支撑对局回放和完整证据记录。

Registration binds the Agent definition hash; result submission anchors the evidence hash on chain.

报名绑定 Agent 定义哈希，结果提交将证据哈希锚定在链上。

The auditor then connects the recorded competition to the approved settlement outcome.

审核员据此将已记录的比赛过程与最终确认的结算结果连接起来。

## 09:24 · 22 / Security by responsibility

Wallet login uses a single-use nonce, and smart contracts enforce the required roles for each operation.

钱包登录使用一次性 nonce，智能合约校验每项操作所需的角色权限。

Escrow transfers use safe token handling and reentrancy protection, with separate accounting for each tournament.

托管转账采用安全代币处理与重入保护，每场赛事单独记账。

Claim checks prevent duplicate payouts. The winner requests payment from their own wallet after approval.

领取检查防止重复付款，审核通过后由冠军使用自己的钱包主动领奖。

These controls support a clear and traceable financial process.

这些控制共同支撑清楚、可追溯的资金流程。

## 09:48 · A complete tournament, from start to settlement

One funded tournament, two language model Agents, independent review, and one hundred DUSD successfully claimed.

一场已注资赛事，两个大模型 Agent，独立审核，以及成功领取的 100 DUSD 奖金。

That is Agent Arena, from sponsorship to settlement.

这就是 Agent Arena：从赞助到结算的完整流程。
