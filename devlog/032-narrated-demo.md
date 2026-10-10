# 032：英文配音与双语字幕讲解视频

日期：2026-10-10。起点版本：`a42593c`；分支：`codex/course-prize-escrow`。

## 需求与范围

用户希望得到约十分钟的课程演示视频，由助手操作并讲解，采用英文配音、中英文双语字幕和英文界面。本轮围绕当前本地课程平台，制作十章经过编辑的 walkthrough，说明用途、角色、奖金托管、Agent 比赛、证据、独立审核、领奖、取消退款、架构和验证边界。

镜头查看已有本地赛事和链上事件，不新增 MetaMask 登录签名、赛事创建、报名或转账。角色权限、资金流程、合约、后端 API 和前端业务功能代码未变；视频制作不能计作新增的业务验收或公网部署。

## 已完成的制作与学习材料

- 已完成十章画面采集（实际平台界面及说明卡），截图序列以 5 fps 采集，按章节编辑；镜头不足旁白时长时保留末帧。这不是连续十分钟的新交易录屏，输出 30 fps 也不代表操作以 30 fps 采集。
- 已准备 10 章、50 组一一对应的英文旁白和中文字幕，并保留可读脚本与 JSON 来源。
- 已生成 macOS Samantha 英文合成配音，语速参数 145，每组旁白后保留 0.25 秒间隔。实际音频时间轴为 **596.392 秒，约 9 分 56 秒**，来自 `deliverables/demo-video/timed-narration.json`。
- 已留存架构与验证说明卡，以及本地配音、字幕、时间轴和视频合成工具。说明卡是项目解释素材，不能当作交易页面或运行中的测试输出。

| 学习关注点 | 入口（相对于仓库根目录） |
| --- | --- |
| 观看入口、章节时间、制作命令与限制 | [docs/course-project/DEMO_VIDEO.md](../docs/course-project/DEMO_VIDEO.md) |
| 逐段英文与中文对照 | [docs/course-project/demo/SCRIPT.en-zh.md](../docs/course-project/demo/SCRIPT.en-zh.md) |
| 50 组最终旁白的数据来源 | [docs/course-project/demo/NARRATION.en-zh.json](../docs/course-project/demo/NARRATION.en-zh.json) |
| 架构与验证解释素材 | [docs/course-project/demo/EXPLANATION_CARDS.html](../docs/course-project/demo/EXPLANATION_CARDS.html) |
| 可复用的本地合成流程 | [tooling/build-demo-video.py](../tooling/build-demo-video.py) |

主查看目标是 `deliverables/demo-video/demo-captioned.mp4`；同目录保留无烧录字幕的视频、音频和中英字幕文件。媒体、截图与缓存保持忽略，Git 留存文案、说明卡及生成工具。干净克隆需要另取媒体或重新采集真实界面，不能仅凭旁白恢复原画面。

## 演示证据与边界

画面对应本地 Anvil 与 DUSD 测试代币。采集快照有 4 场赛事，累计赞助 325 DUSD、已付奖金 100、已退 225、托管余额 0。#1 展示已经完成的 12 回合确定性模拟及奖金记录；#3 展示已有的 200 DUSD 退款；#2 的争议、裁决和 25 DUSD 退款是预置合约测试案例，不是真实客户投诉。

比赛记录明确标注确定性模拟，不能把回放称为真实 LLM 评价。钱包入口和创建权限只用于介绍操作入口；实际资金步骤通过已有事件、交易记录、余额和证据哈希说明。视频没有执行新的钱包签名，也没有重新运行合约、后端或前端业务测试；历史验证结果继续以[课程测试报告](../docs/course-project/TEST_REPORT.md)为准。

本轮未完成公网云部署、新课程 Sepolia 验证、真实用户评价或外部安全审计。哈希一致只能支持数据核验，不能证明比赛公平或链下数据真实。报告、Slides 和真实利益相关者验证仍需按[提交计划](../docs/course-project/SUBMISSION_PLAN.md)继续准备。

## 验证结果

- 最终字幕版与无字幕版均已生成。字幕版 596.400 秒、1920 × 1080、H.264 / AAC，30,305,416 bytes；英文音轨 596.392 秒，音画时长差 8 毫秒。
- 最终字幕视频经 FFmpeg 全文件解码，退出码 0、无 error 输出。查看封面及 14 个时间点的画面，字幕未覆盖平台操作区，中文字体正常；100 秒、448 秒画面按原尺寸复查。
- 50 组源文案、配音缓存文本、时间轴逐条一致；三份 SRT 均为 50 条、无重叠且在音频范围内。音量检测 mean −15.9 dB、peak −1.5 dB。
- 审阅修正了“报名不足可取消”的期限表达，明确必须报名截止后不足两人；只重新生成对应配音及受影响章节。合成器提供显式 --reuse-chapters，并验证源清单、修改时间、视频参数与时长后才复用。
- 修复中文孤立标点和混排英文单词换行；生成脚本编译检查与 git diff --check 通过。

原始信息见[媒体验证记录](../docs/course-project/evidence/demo-video-verification-2026-10-10.json)、[封面](../docs/course-project/evidence/demo-video-cover-2026-10-10.png)；视频说明保留复用命令和章节目录。抽帧与程序检查不等于真实观众逐秒听审。

本轮不修改角色或业务功能，不重跑业务测试。文档入口、迭代计划及日志白名单同步更新；此前迭代历史和测试网、公有云正式部署待办保留。录制使用的临时标签页和说明卡服务器已关闭，平台服务与数据保持原状。
