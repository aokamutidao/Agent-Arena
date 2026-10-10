# 用户实操录屏版演示

本版使用用户在 2026-10-10 录制的四段原始视频，展示赛事 **#13 / agent fight test 4** 的真实浏览器操作。编排为平台介绍 → 创建、报名、赞助 → LLM 对战与回放 → 提交、审核、领奖 → 审计与历史退款 → 架构讲解。

旁白为英文合成语音，字幕中英双语。开篇说明本地环境响应快、便于完整演示；DUSD 是本地演示代币。主要内容展示功能和实现，不设项目不足章节。

## 已交付成片

- [完整演示 MP4：英文配音＋中英字幕](../../../deliverables/demo-video-user/Agent-Arena-Demo-EN-ZH.mp4)，10:00，1920×1080，30 fps，约 35.9 MB。
- [无烧录字幕版](../../../deliverables/demo-video-user/Agent-Arena-Demo-clean.mp4)，便于后续重剪。
- [双语 SRT](../../../deliverables/demo-video-user/subtitles-bilingual.srt)、[独立配音](../../../deliverables/demo-video-user/narration.m4a)、[章节时间轴](../../../deliverables/demo-video-user/chaptertimestamps.md)。
- [媒体验证记录](VALIDATION.json)：完整解码通过，90 条字幕无时间重叠，赞助到账、审核与领奖关键帧已复核。

以上媒体文件保存在本机并排除于 Git；讲稿和制作工具纳入版本管理。

## 素材与事实

| 用户原片 | 内容 | 源片时长 |
| --- | --- | --- |
| `1赛事准备.mov` | 主办方创建、指定审核员、红蓝报名、100 DUSD 赞助 | 165.833 秒 |
| `2 开始对战.mov` | 开赛、真实模型对战、证据、回放、提交结果 | 123.433 秒 |
| `3审计检查.mov` | 审核员登录、查看证据、审核确认 | 67.850 秒 |
| `4 领取奖励.mov` | 蓝方登录、领取 100 DUSD、到账展示 | 36.000 秒 |

本场界面显示 `qwen3.5-plus`，15 回合，30 次成功模型调用，0 次失败调用，蓝方胜出并剩余 10 HP。领奖画面显示蓝方余额 300→400 DUSD，本场托管 100→0 DUSD。对应界面在补录时只读复核；本轮视频制作没有再发起比赛、钱包交易或模型调用。

开篇首页、资金概览和末尾审计流水由助手补录真实 `localhost:3002` 页面。历史退款镜头查看赛事 #6 的 25 DUSD 入金、取消和退款记录，不在新视频中执行退款。该历史记录来自此前 CLI 测试，与 #13 的用户钱包实操区分。

## 剪辑规则

- 原始 MOV 文件留在桌面，保持不变；输出在 `deliverables/demo-video-user/`，中间文件在 `.local/user-demo-edit/`。
- 裁去顶部无关 Chrome 标签、macOS 菜单栏和底部 Dock，保留平台与交易确认侧栏。
- 删除大部分账户选择列表、重复查看和空闲等待；保留创建、报名、授权、赞助、开赛、提交、审核、领奖的操作及成功画面。
- 各素材内部保持时间顺序。操作画面按原速播放，旁白需要阅读时间时延长末帧；字幕不将等待剪辑描述为链上耗时。
- 移除原录屏音轨，使用本地 macOS `Samantha` 英文合成配音。完整原始素材未上传外部配音或剪辑服务。
- 英文/中文字幕在独立底部区域，避免覆盖按钮、棋盘、钱包金额和结果。
- 架构卡是解释图，按当前课程版源码绘制；实际使用 Next.js、Go/Gin、SQLite、Solidity、LLM，不引入尚未实现的 MySQL、Redis 或 RAG 叙述。

## 可复用文件

- [配音与剪辑清单](NARRATION.en-zh.json)：章节时长、逐句中英文本、源片选段与关键画面对齐点。
- [双语阅读稿](SCRIPT.en-zh.md)：方便复习项目、准备演示和答辩。
- [架构旁白来源](ARCHITECTURE_DRAFT.json)：核对过的架构描述与源码依据。
- `tooling/render-demo-explainer-cards.py`：生成七张解释图。
- `tooling/build-demo-video.py`：生成并缓存逐句英文配音。
- `tooling/build-user-recorded-demo.py`：组合实际视频、补录、架构图、配音与字幕。

## 复现命令

先将四个原始 MOV 文件放在本机桌面，保留上述名称，并保留 `.local/user-demo-edit/intro/` 中的实际补录帧。随后在仓库根目录运行：

```sh
python3 tooling/render-demo-explainer-cards.py
python3 tooling/build-demo-video.py --prepare-audio --narration docs/course-project/demo-user-recording/NARRATION.en-zh.json --output .local/user-demo-edit/voice --voice Samantha --rate 150 --silence 0.2
python3 tooling/build-user-recorded-demo.py --all
```

解释图生成需要 Pillow；配音依赖 macOS `say`，剪辑依赖 FFmpeg/FFprobe。当前电脑的 Pillow Python 路径由工作区运行时提供；若系统 Python 没有 Pillow，使用已配置运行时执行第一行。`--plan` 可单独查看时间轴，`--audio` 可单独输出配音和字幕。各步骤需等待上一步完成，再执行下一步。

视频和本地缓存不进入 Git；版本管理保存脚本、旁白、剪辑清单与制作记录。最终媒体信息和验证结果见 [035 制作记录](../../../devlog/035-user-recorded-demo.md)。
