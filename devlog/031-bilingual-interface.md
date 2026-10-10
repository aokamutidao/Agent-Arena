# 031：课程平台中英文切换

日期：2026-10-10。起点版本：`eef2889`；分支：`codex/course-prize-escrow`。

## 需求与范围

用户需要英文界面用于课程展示，并能随时切回中文。本轮在已有界面接入语言选择，覆盖资金概览、赛事列表与详情、Agent、审计、指南、钱包连接、交易状态和比赛回放。

顶部的 English / 中文 按钮直接切换，手机上也可见。长操作按钮允许换行，避免英文在窄屏溢出。语言偏好保存在本机浏览器的 `lang` 键，首次默认中文；存储不可用时仍能在本次页面内切换。无新增 npm 依赖。

## 怎么实现，学习时看哪里

| 关注点 | 代码入口（相对于 `frontend/src/`） |
| --- | --- |
| 全局语言与浏览器偏好 | `lib/i18n.tsx`：`lang` / `setLang` / `l(中文, English)` |
| 明显且可访问的切换入口 | `components/layout/LanguageSwitcher.tsx` 与 `Navbar.tsx` |
| 页脚、跳转链接、页面布局 | `components/layout/AppShell.tsx`；根 `app/layout.tsx` 继续保留服务端字体与 metadata |
| 钱包弹窗语言 | `components/providers.tsx` 将 I18nProvider 放在动态钱包树外，`wallet-providers.tsx` 给 RainbowKit 传 locale |
| 状态语义与角色 | `lib/course-presentation.ts` 接收语言参数，默认中文，保持原状态判断 |
| 日期和已知错误提示 | `lib/course.ts` 的 `dateTime` / `errorMessage`；未知 API / 钱包诊断保持原文 |
| 页面文案与动态值 | 各课程 `app/**/page.tsx` 在渲染时调用 `l`，数量和操作说明使用完整句子的双语模板 |

语言变化是普通 React context 更新，没有按语言设置 Provider 的 `key`。钱包、登录和交易状态不因语言主动重建。表单 state 保留原值，临时验证错误保存原始错误或双语信息后再渲染，避免错误提示永久锁定为产生时的语言。

增加文案时，把一对翻译放在显示位置：`l("报名参赛", "Register")`；金额或数量使用两套完整模板。不要把数据库字段、合约数据、用户输入或待签名数据送入自动翻译。已有 `t(key)` 字典保留供历史组件兼容，课程页面使用新的双语入口。

## 有意保留原文的内容

- 赛事名称、Agent 名称/策略、用户取消或裁决说明、规则链接。
- canonical JSON、证据哈希、事件原始代码、交易编号、钱包地址及服务原始诊断。
- 已估算交易的操作标签是当时的快照；切换语言不修改固定参数，后续重新估算使用当前语言。
- RainbowKit 弹窗随平台切换；MetaMask 扩展自身语言由钱包设置决定。

本轮没有修改合约、Go API、数据库或既有演示赛事。创建英文演示数据时，由用户填写英文名称和说明，不能篡改已经绑定哈希的证据来做翻译。

## 验证与版本对比

实际执行结果和截图见[课程测试报告](../docs/course-project/TEST_REPORT.md)。原始数据保留不变；仅重建和重启前端展示服务。

```bash
git diff eef2889..HEAD -- frontend/src
git diff eef2889..HEAD -- docs/course-project devlog/031-bilingual-interface.md
(cd frontend && npm test && npm run typecheck && npm run lint)
```

需要对照前一轮体验改造时，先阅读 [030 日志](030-frontend-usability.md)，再看本轮差异；两轮功能独立留档，避免把国际化误认为又做了一次业务重构。
