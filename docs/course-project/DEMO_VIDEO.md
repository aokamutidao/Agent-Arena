# 作业演示视频

主查看文件为 [demo-captioned.mp4](../../deliverables/demo-video/demo-captioned.mp4)：英文配音与烧录的中英文双语字幕。建议用本地播放器全屏查看；课程系统提交前核对其格式、大小和上传规则。

成片为 **596.400 秒（9 分 56 秒）、1920 × 1080、H.264 + AAC**，字幕版约 **30.3 MB**。英文配音使用 macOS Samantha 合成声音，语速参数 145。完整解码通过；50 组字幕与配音时间轴逐条对应，音视频时长差 8 毫秒。已抽查 14 个关键时点的画面与双语字幕排版。详见[媒体验证记录](evidence/demo-video-verification-2026-10-10.json)和[封面](evidence/demo-video-cover-2026-10-10.png)。

## 文件与学习资料

| 文件 | 用途 |
| --- | --- |
| [demo-captioned.mp4](../../deliverables/demo-video/demo-captioned.mp4) | 主视频，英文配音、中英字幕直接显示在画面上 |
| [demo-clean.mp4](../../deliverables/demo-video/demo-clean.mp4) | 同一画面和配音，不烧录字幕，便于重新编辑 |
| [audio.m4a](../../deliverables/demo-video/audio.m4a) / [narration.wav](../../deliverables/demo-video/audio/narration.wav) | 完整英文合成配音 |
| [双语 SRT](../../deliverables/demo-video/subtitles-bilingual.srt)、[英文 SRT](../../deliverables/demo-video/subtitles-en.srt)、[中文 SRT](../../deliverables/demo-video/subtitles-zh.srt)、[ASS](../../deliverables/demo-video/subtitles.ass) | 外挂字幕或重排字幕；主视频已烧录字幕，不需再次叠加 |
| [双语脚本](demo/SCRIPT.en-zh.md) / [旁白 JSON](demo/NARRATION.en-zh.json) | 最终 10 章、50 组对应文案，可追溯和复用 |
| [说明卡 HTML](demo/EXPLANATION_CARDS.html) | 架构及验证说明，分别以 #architecture / #validation 查看 |
| [生成脚本](../../tooling/build-demo-video.py) | 本地配音、时间轴、字幕及视频合成 |

媒体、音频缓存与界面截图不作为 Git 源码提交；学习文案和说明卡保留在 docs。干净克隆需要另取媒体或重新采集界面，不能仅凭旁白 JSON 恢复原截图。

## 章节时间

以下为配音时间轴，单位为分:秒.毫秒。最终视频可能存在编码级时长差异。

| 章节 | 时间 | 内容 |
| --- | --- | --- |
| 01 | 00:00.000–01:00.717 | 用途与本地演示边界 |
| 02 | 01:00.717–02:00.876 | 用户、角色与指南 |
| 03 | 02:00.876–02:54.116 | 钱包入口与创建权限 |
| 04 | 02:54.116–03:52.348 | 报名及 Agent 定义 |
| 05 | 03:52.348–04:47.097 | 已有比赛回放 |
| 06 | 04:47.097–05:41.219 | 证据与哈希核验 |
| 07 | 05:41.219–06:40.380 | 审核及领奖记录 |
| 08 | 06:40.380–07:42.745 | 取消及赞助退款 |
| 09 | 07:42.745–08:50.195 | 架构与测试证据 |
| 10 | 08:50.195–09:56.392 | 双语体验、限制与结束 |

## 画面与内容边界

- 画面来自当前本地平台的实际界面截图序列（采集 5 fps），按章节剪辑。不足旁白时长的章节保留末帧；这是经过编辑的 walkthrough，不是连续十分钟的新增交易录屏。生成器默认输出 30 fps，不会把 5 fps 的采集变成真实 30 fps 操作记录。
- 平台使用英文界面；用户已有赛事名、策略、说明和原始证据保留原文。配音是英文合成声音，不是用户本人录音；字幕是对应的中英双语文本。
- 本次只查看已有本地 Anvil 记录，不新增 MetaMask 登录签名、报名、创建或转账。钱包连接弹窗只展示后关闭；创建入口展示权限提示。已有的本地集成测试由 viem 使用公开测试账户签名，并非 MetaMask 扩展现场操作。
- 展示快照为 4 场赛事、累计赞助 325 DUSD、已付奖金 100、退款 225、托管余额 0。#1 为已完成的 12 回合确定性模拟比赛；#3 已确认退回 200 DUSD；#2 的争议/裁决/25 DUSD 退款是 prepared contract dispute fixture，不是真实客户投诉。
- 第 9 章架构与验证卡片属于项目说明，不是平台交易页面。49 项合约测试及后端 race/vet 等是历史验证记录；27 项前端测试来自中英文改造。视频没有重新执行这些测试，也不构成外部安全审计。
- 本地 DUSD 不承诺现实资产价值。本轮没有完成公网云部署、新课程 Sepolia 验证、真实 LLM 评价、完整 MetaMask 人工验收或真实利益相关者验证；哈希一致不能证明比赛真实、公平。课程报告应继续使用 [测试报告](TEST_REPORT.md) 与 [需求证据边界](REQUIREMENTS.md) 的实际记录。

## 生成与复用

已查看并执行生成器的 --help；它只合成输入素材，不操作浏览器、不自动采集画面。需要 Python 3 标准库、macOS say、ffmpeg 与 ffprobe，字幕烧录需 ffmpeg 支持 subtitles/libass。下列命令从仓库根目录运行，构建会更新指定输出目录中的同名媒体文件。

先将每章真实截图放入 .local/demo-video/captures/01 至 10。各章建议提供 manifest.json，记录真实相对时间与图像路径，例如：

~~~json
{
  "duration": 12.8,
  "frames": [
    { "file": "frames001.png", "t": 0 },
    { "file": "frames002.png", "t": 0.2 }
  ]
}
~~~

帧路径相对该章目录；没有 manifest 的编号图像按 --capture-fps 推定时间。当前采集素材属于本机本轮工作文件，复用时应先核对内容、身份、金额及章节完整性，不用旧素材冒充新部署。

~~~bash
python3 tooling/build-demo-video.py --help

python3 tooling/build-demo-video.py --prepare-audio \
  --narration docs/course-project/demo/NARRATION.en-zh.json \
  --output deliverables/demo-video \
  --voice Samantha --rate 145 --silence 0.25

python3 tooling/build-demo-video.py --render \
  --captures .local/demo-video/captures \
  --output deliverables/demo-video \
  --capture-fps 5 --fps 30
~~~

也可合并 --prepare-audio 与 --render。只修改字幕或部分旁白后，可为 --render 添加 --reuse-chapters：生成器核对截图清单、修改时间、分辨率、帧率与目标时长后复用未变化章节；输入不一致时仍重新编码，最终音轨与字幕始终重新合成。--ffmpeg、--ffprobe、--say 可指定现有可执行文件路径。语音工具若生成空文件，需排查本机音频文件访问权限；不要把空音频当作成功。修改英文文案、语速或声音后应重新准备音频，再根据新的 timed-narration.json 检查镜头与字幕时间。

当前说明卡是独立 1280 × 720 HTML，可在浏览器打开 EXPLANATION_CARDS.html#architecture 或 #validation。它用于解释模块与验证边界，不会调用钱包或资金接口。完整解码、字幕与时间轴检查、关键帧抽查已完成；这些媒体验证不替代真实钱包多角色操作、云部署或用户评价。
