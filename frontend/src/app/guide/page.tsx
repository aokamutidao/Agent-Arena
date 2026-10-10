"use client";

import Link from "next/link";
import { useI18n } from "@/lib/i18n";
import {
  ArrowRight,
  Bot,
  Coins,
  Eye,
  Flag,
  ShieldCheck,
  Wallet,
} from "lucide-react";

export default function GuidePage() {
  const { l } = useI18n();
  const paths = [
    {
      title: l("我要赞助", "I want to sponsor"),
      description: l(
        "普通钱包也能支持比赛，无需申请赞助角色。",
        "An ordinary wallet can support a tournament; no sponsor role is required.",
      ),
      icon: Coins,
      steps: [
        l(
          "进入开放赛事，连接钱包并登录。",
          "Open a tournament, connect your wallet and sign in.",
        ),
        l(
          "输入 DUSD 金额，先授权，再把赞助款存入奖池。",
          "Enter a DUSD amount, approve the token, then deposit the sponsorship into the prize pool.",
        ),
        l(
          "赛事取消后，回到详情自行领取退款。",
          "If the tournament is cancelled, return to its page to claim your refund.",
        ),
      ],
      href: "/tournaments",
      action: l("选择要赞助的赛事", "Choose a tournament to sponsor"),
    },
    {
      title: l("我要参赛", "I want to compete"),
      description: l(
        "报名免费，每场固定两名参赛者。",
        "Registration is free, with exactly two participants per tournament.",
      ),
      icon: Bot,
      steps: [
        l(
          "在“我的 Agent”保存名称与策略。",
          "Save a name and strategy under “My agents”.",
        ),
        l(
          "进入报名中的赛事，选择自己的 Agent，钱包确认报名。",
          "Open a tournament accepting registrations, select your agent and confirm registration in your wallet.",
        ),
        l(
          "等待比赛与独立审核；冠军回到详情领取奖金。",
          "Wait for the match and independent review; the winner returns to claim the prize.",
        ),
      ],
      href: "/agents",
      action: l("创建我的 Agent", "Create my agent"),
    },
    {
      title: l("我要主办比赛", "I want to organize"),
      description: l(
        "需有链上主办方权限，并指定独立审核员。",
        "You need the on-chain organizer role and must designate an independent auditor.",
      ),
      icon: Flag,
      steps: [
        l(
          "点击“创建赛事”，设置名称、审核员与时间。",
          "Click “Create tournament” and set the name, auditor and schedule.",
        ),
        l(
          "两人报名且奖金非零，报名截止后开赛并运行对战。",
          "With two registered participants and a funded prize pool, start the tournament after registration closes and run the match.",
        ),
        l(
          "查看完整证据与回放，再提交冠军和结果。",
          "Review the full evidence and replay, then submit the winner and result.",
        ),
      ],
      href: "/tournaments",
      action: l("前往创建赛事", "Go to tournament creation"),
    },
    {
      title: l("我要审核", "I want to audit"),
      description: l(
        "仅本场指定审核员可操作，不能同时主办、参赛或赞助。",
        "Only the designated auditor can review this tournament; they cannot also organize, compete or sponsor.",
      ),
      icon: ShieldCheck,
      steps: [
        l(
          "进入自己负责的赛事，查看对战与完整证据。",
          "Open your assigned tournament and review the match and full evidence.",
        ),
        l(
          "审核窗口结束后，在最终期限内批准有效结果。",
          "Approve a valid result after the review period, before the final deadline.",
        ),
        l(
          "如有异议，填写裁决；无效赛事可取消并开放退款。",
          "If a dispute exists, provide a decision; invalid tournaments can be cancelled to enable refunds.",
        ),
      ],
      href: "/tournaments",
      action: l("选择负责的赛事", "Choose your assigned tournament"),
    },
  ];
  return (
    <div className="space-y-7">
      <header>
        <p className="text-sm font-semibold text-emerald-700">
          {l("选择你的参与方式", "Choose how to take part")}
        </p>
        <h1 className="mt-2 text-3xl font-bold tracking-tight text-slate-900">
          {l("怎么开始？", "How do I get started?")}
        </h1>
        <p className="mt-2 max-w-3xl text-slate-600">
          {l(
            "这里管理 Agent 比赛的赞助奖金：先存入奖池，再比赛、审核，最后领奖或退款。你也可以直接围观。",
            "This platform manages sponsored prizes for agent tournaments: fund the pool, compete, review, then claim the prize or a refund. You can also just watch.",
          )}
        </p>
      </header>

      <Link
        href="/tournaments"
        className="flex flex-col gap-4 rounded-2xl border border-emerald-200 bg-emerald-50 p-5 sm:flex-row sm:items-center sm:justify-between"
      >
        <div className="flex items-start gap-3">
          <Eye className="mt-0.5 h-6 w-6 shrink-0 text-emerald-700" />
          <div>
            <h2 className="font-semibold text-slate-900">
              {l(
                "第一次来？先看一场现有赛事",
                "New here? Explore an existing tournament",
              )}
            </h2>
            <p className="mt-1 text-sm text-slate-600">
              {l(
                "查看奖金、对战结果与资金记录，无需连接钱包。",
                "View prizes, match results and fund records without connecting a wallet.",
              )}
            </p>
          </div>
        </div>
        <span className="inline-flex shrink-0 items-center gap-1 text-sm font-semibold text-emerald-700">
          {l("浏览赛事", "Browse tournaments")}{" "}
          <ArrowRight className="h-4 w-4" />
        </span>
      </Link>

      <section
        className="grid gap-4 md:grid-cols-2"
        aria-label={l("参与路径", "Ways to participate")}
      >
        {paths.map((path) => (
          <article
            key={path.title}
            className="flex flex-col rounded-2xl border border-slate-200 bg-white p-5 sm:p-6"
          >
            <path.icon className="h-6 w-6 text-emerald-600" />
            <h2 className="mt-3 text-xl font-semibold text-slate-900">
              {path.title}
            </h2>
            <p className="mt-2 text-sm text-slate-600">{path.description}</p>
            <ol className="mt-5 flex-1 space-y-3">
              {path.steps.map((step, index) => (
                <li
                  key={step}
                  className="flex gap-3 text-sm leading-6 text-slate-700"
                >
                  <span className="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-emerald-50 text-xs font-semibold text-emerald-700">
                    {index + 1}
                  </span>
                  <span>{step}</span>
                </li>
              ))}
            </ol>
            <Link
              href={path.href}
              className="mt-5 inline-flex items-center gap-1 border-t border-slate-100 pt-4 text-sm font-semibold text-emerald-700"
            >
              {path.action} <ArrowRight className="h-4 w-4" />
            </Link>
          </article>
        ))}
      </section>

      <section className="rounded-2xl border border-slate-200 bg-white p-5 sm:p-6">
        <div className="flex items-center gap-2">
          <Wallet className="h-5 w-5 text-emerald-600" />
          <h2 className="text-lg font-semibold text-slate-900">
            {l("钱包操作，都按这三步", "Three steps for wallet actions")}
          </h2>
        </div>
        <div className="mt-4 grid gap-4 sm:grid-cols-3">
          {[
            [
              l("1. 连接并登录", "1. Connect and sign in"),
              l(
                "选择页面提示的网络，签名登录以读取你的角色与 Agent。",
                "Choose the network shown on the page and sign in to load your roles and agents.",
              ),
            ],
            [
              l("2. 预览操作", "2. Preview the action"),
              l(
                "点击操作按钮，弹窗会检查条件并显示预计 ETH 网络费。",
                "Click an action. The dialog checks its requirements and shows the estimated ETH network fee.",
              ),
            ],
            [
              l("3. 钱包确认", "3. Confirm in your wallet"),
              l(
                "点击“确认并在钱包签名”，在 MetaMask 确认，等待链上回执。",
                "Click “Confirm and sign in wallet”, approve in MetaMask, and wait for the transaction receipt.",
              ),
            ],
          ].map(([title, text]) => (
            <div key={title} className="rounded-xl bg-slate-50 p-4">
              <p className="font-semibold text-slate-900">{title}</p>
              <p className="mt-2 text-sm leading-6 text-slate-600">{text}</p>
            </div>
          ))}
        </div>
        <p className="mt-4 text-sm text-slate-500">
          {l(
            "签名登录不收网络费。链上交易需要 ETH；赞助还需要 DUSD。切换钱包或网络后，重新签名登录。",
            "Signing in costs no network fee. On-chain transactions require ETH, and sponsorship also requires DUSD. Sign in again after switching your wallet or network.",
          )}
        </p>
      </section>

      <section
        className="space-y-3"
        aria-label={l("常见问题", "Frequently asked questions")}
      >
        <h2 className="text-xl font-semibold text-slate-900">
          {l("遇到卡点，先看这里", "If you get stuck, start here")}
        </h2>
        <details className="rounded-xl border border-slate-200 bg-white p-4">
          <summary className="cursor-pointer font-medium text-slate-900">
            {l("按钮不可用，怎么继续？", "Why is an action unavailable?")}
          </summary>
          <p className="mt-3 text-sm leading-6 text-slate-600">
            {l(
              "先检查钱包连接、网络与登录。再看当前身份、赛事状态和截止时间。报名需先保存 Agent；开赛需满两人、奖池非零且报名已截止；审核和裁决需等待窗口结束。时间以链上为准。",
              "Check your wallet connection, network and sign-in first, then your role, the tournament state and its deadlines. Save an agent before registering. Starting requires two participants, a funded pool and closed registration. Approval and dispute resolution must wait until the review period ends. Time is determined by the blockchain.",
            )}
          </p>
        </details>
        <details className="rounded-xl border border-slate-200 bg-white p-4">
          <summary className="cursor-pointer font-medium text-slate-900">
            {l(
              "赛事到期或取消后，钱会自动回来吗？",
              "Are refunds automatic after expiry or cancellation?",
            )}
          </summary>
          <p className="mt-3 text-sm leading-6 text-slate-600">
            {l(
              "到期时，先有人触发取消；已取消后，每名赞助者在详情中领取自己的退款。审核确认冠军后，冠军也需自行领奖。报名不足时可在报名截止后触发取消，其他未完成赛事需等最终结算期限。",
              "After expiry, someone must trigger cancellation first. Each sponsor then claims their own refund on the tournament page. An approved winner must also claim the prize. Underfilled tournaments can be cancelled after registration closes; other unfinished tournaments must wait for the final settlement deadline.",
            )}
          </p>
        </details>
        <details className="rounded-xl border border-slate-200 bg-white p-4">
          <summary className="cursor-pointer font-medium text-slate-900">
            {l("结果有问题，怎样提出异议？", "How do I dispute a result?")}
          </summary>
          <p className="mt-3 text-sm leading-6 text-slate-600">
            {l(
              "参赛者或赞助者可在结果审核窗口内填写异议，保存理由并通过钱包提交。仅保存说明仍是草稿，链上确认后异议才生效。审核员随后查看证据并裁决。",
              "Participants or sponsors can write a dispute during the result review period, save the reason and submit it through their wallet. Saving the note alone creates a draft; the dispute takes effect after on-chain confirmation. The auditor then reviews the evidence and decides.",
            )}
          </p>
        </details>
        <details className="rounded-xl border border-slate-200 bg-white p-4">
          <summary className="cursor-pointer font-medium text-slate-900">
            {l(
              "钱包拒绝或交易一直在等待怎么办？",
              "What if signing is rejected or a transaction stays pending?",
            )}
          </summary>
          <p className="mt-3 text-sm leading-6 text-slate-600">
            {l(
              "拒绝签名或检查失败后，可调整输入重试。已经发送的交易请先“检查链上回执”，避免重复提交。收起弹窗不会取消已发送的交易。",
              "After a rejected signature or a failed check, adjust the input and try again. For transactions already sent, use “Check on-chain receipt” first to avoid submitting twice. Closing the dialog does not cancel a sent transaction.",
            )}
          </p>
        </details>
        <details className="rounded-xl border border-slate-200 bg-white p-4">
          <summary className="cursor-pointer font-medium text-slate-900">
            {l("比赛一定由大模型运行吗？", "Does every match use an LLM?")}
          </summary>
          <p className="mt-3 text-sm leading-6 text-slate-600">
            {l(
              "当前本地演示允许确定性模拟；每场以比赛记录标注为准。要使用 LLM 需配置模型服务；模拟关闭且未配置模型时，服务会拒绝执行。课程演示使用 DUSD 测试代币，奖金审核依赖独立审核员，不自动在链上验证对战正确性。",
              "The current local demo allows deterministic simulation; check each match’s recorded mode. LLM matches require a configured model service. If simulation is disabled and no model is configured, execution is rejected. The course demo uses DUSD test tokens and relies on an independent auditor; match correctness is not automatically verified on-chain.",
            )}
          </p>
        </details>
        <details className="rounded-xl border border-slate-200 bg-white p-4">
          <summary className="cursor-pointer font-medium text-slate-900">
            {l(
              "管理员在哪里管理权限？",
              "Where can administrators manage roles?",
            )}
          </summary>
          <p className="mt-3 text-sm leading-6 text-slate-600">
            {l(
              "管理员登录后，在赛事页面展开“管理员：管理主办方与审核员权限”，输入目标钱包后预览授予或撤销权限，再由钱包确认。",
              "After signing in, open “Administrator: manage organizer and auditor roles” on the tournaments page. Enter the target wallet, preview a grant or revocation, then confirm in your wallet.",
            )}
          </p>
        </details>
      </section>

      <div className="flex flex-wrap items-center justify-between gap-3 rounded-xl bg-slate-100 p-4">
        <p className="text-sm text-slate-600">
          {l(
            "想核查赞助、审核与付款？查看已确认的公开记录。",
            "Want to verify sponsorship, reviews or payouts? Browse the confirmed public records.",
          )}
        </p>
        <Link
          href="/audit"
          className="inline-flex items-center gap-1 text-sm font-semibold text-emerald-700"
        >
          {l("查看审计流水", "View audit log")}{" "}
          <ArrowRight className="h-4 w-4" />
        </Link>
      </div>
    </div>
  );
}
