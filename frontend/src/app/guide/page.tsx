import Link from "next/link";
import {
  ArrowRight,
  Bot,
  Coins,
  Eye,
  Flag,
  ShieldCheck,
  Wallet,
} from "lucide-react";

const paths = [
  {
    title: "我要赞助",
    description: "普通钱包也能支持比赛，无需申请赞助角色。",
    icon: Coins,
    steps: [
      "进入开放赛事，连接钱包并登录。",
      "输入 DUSD 金额，先授权，再把赞助款存入奖池。",
      "赛事取消后，回到详情自行领取退款。",
    ],
    href: "/tournaments",
    action: "选择要赞助的赛事",
  },
  {
    title: "我要参赛",
    description: "报名免费，每场固定两名参赛者。",
    icon: Bot,
    steps: [
      "在“我的 Agent”保存名称与策略。",
      "进入报名中的赛事，选择自己的 Agent，钱包确认报名。",
      "等待比赛与独立审核；冠军回到详情领取奖金。",
    ],
    href: "/agents",
    action: "创建我的 Agent",
  },
  {
    title: "我要主办比赛",
    description: "需有链上主办方权限，并指定独立审核员。",
    icon: Flag,
    steps: [
      "点击“创建赛事”，设置名称、审核员与时间。",
      "两人报名且奖金非零，报名截止后开赛并运行对战。",
      "查看完整证据与回放，再提交冠军和结果。",
    ],
    href: "/tournaments",
    action: "前往创建赛事",
  },
  {
    title: "我要审核",
    description: "仅本场指定审核员可操作，不能同时主办、参赛或赞助。",
    icon: ShieldCheck,
    steps: [
      "进入自己负责的赛事，查看对战与完整证据。",
      "审核窗口结束后，在最终期限内批准有效结果。",
      "如有异议，填写裁决；无效赛事可取消并开放退款。",
    ],
    href: "/tournaments",
    action: "选择负责的赛事",
  },
];

export default function GuidePage() {
  return (
    <div className="space-y-7">
      <header>
        <p className="text-sm font-semibold text-emerald-700">
          选择你的参与方式
        </p>
        <h1 className="mt-2 text-3xl font-bold tracking-tight text-slate-900">
          怎么开始？
        </h1>
        <p className="mt-2 max-w-3xl text-slate-600">
          这里管理 Agent
          比赛的赞助奖金：先存入奖池，再比赛、审核，最后领奖或退款。你也可以直接围观。
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
              第一次来？先看一场现有赛事
            </h2>
            <p className="mt-1 text-sm text-slate-600">
              查看奖金、对战结果与资金记录，无需连接钱包。
            </p>
          </div>
        </div>
        <span className="inline-flex shrink-0 items-center gap-1 text-sm font-semibold text-emerald-700">
          浏览赛事 <ArrowRight className="h-4 w-4" />
        </span>
      </Link>

      <section className="grid gap-4 md:grid-cols-2" aria-label="参与路径">
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
            钱包操作，都按这三步
          </h2>
        </div>
        <div className="mt-4 grid gap-4 sm:grid-cols-3">
          {[
            [
              "1. 连接并登录",
              "选择页面提示的网络，签名登录以读取你的角色与 Agent。",
            ],
            [
              "2. 预览操作",
              "点击操作按钮，弹窗会检查条件并显示预计 ETH 网络费。",
            ],
            [
              "3. 钱包确认",
              "点击“确认并在钱包签名”，在 MetaMask 确认，等待链上回执。",
            ],
          ].map(([title, text]) => (
            <div key={title} className="rounded-xl bg-slate-50 p-4">
              <p className="font-semibold text-slate-900">{title}</p>
              <p className="mt-2 text-sm leading-6 text-slate-600">{text}</p>
            </div>
          ))}
        </div>
        <p className="mt-4 text-sm text-slate-500">
          签名登录不收网络费。链上交易需要 ETH；赞助还需要
          DUSD。切换钱包或网络后，重新签名登录。
        </p>
      </section>

      <section className="space-y-3" aria-label="常见问题">
        <h2 className="text-xl font-semibold text-slate-900">
          遇到卡点，先看这里
        </h2>
        <details className="rounded-xl border border-slate-200 bg-white p-4">
          <summary className="cursor-pointer font-medium text-slate-900">
            按钮不可用，怎么继续？
          </summary>
          <p className="mt-3 text-sm leading-6 text-slate-600">
            先检查钱包连接、网络与登录。再看当前身份、赛事状态和截止时间。报名需先保存
            Agent；开赛需满两人、奖池非零且报名已截止；审核和裁决需等待窗口结束。时间以链上为准。
          </p>
        </details>
        <details className="rounded-xl border border-slate-200 bg-white p-4">
          <summary className="cursor-pointer font-medium text-slate-900">
            赛事到期或取消后，钱会自动回来吗？
          </summary>
          <p className="mt-3 text-sm leading-6 text-slate-600">
            到期时，先有人触发取消；已取消后，每名赞助者在详情中领取自己的退款。审核确认冠军后，冠军也需自行领奖。报名不足时可在报名截止后触发取消，其他未完成赛事需等最终结算期限。
          </p>
        </details>
        <details className="rounded-xl border border-slate-200 bg-white p-4">
          <summary className="cursor-pointer font-medium text-slate-900">
            结果有问题，怎样提出异议？
          </summary>
          <p className="mt-3 text-sm leading-6 text-slate-600">
            参赛者或赞助者可在结果审核窗口内填写异议，保存理由并通过钱包提交。仅保存说明仍是草稿，链上确认后异议才生效。审核员随后查看证据并裁决。
          </p>
        </details>
        <details className="rounded-xl border border-slate-200 bg-white p-4">
          <summary className="cursor-pointer font-medium text-slate-900">
            钱包拒绝或交易一直在等待怎么办？
          </summary>
          <p className="mt-3 text-sm leading-6 text-slate-600">
            拒绝签名或检查失败后，可调整输入重试。已经发送的交易请先“检查链上回执”，避免重复提交。收起弹窗不会取消已发送的交易。
          </p>
        </details>
        <details className="rounded-xl border border-slate-200 bg-white p-4">
          <summary className="cursor-pointer font-medium text-slate-900">
            比赛一定由大模型运行吗？
          </summary>
          <p className="mt-3 text-sm leading-6 text-slate-600">
            当前本地演示允许确定性模拟；每场以比赛记录标注为准。要使用 LLM
            需配置模型服务；模拟关闭且未配置模型时，服务会拒绝执行。课程演示使用
            DUSD 测试代币，奖金审核依赖独立审核员，不自动在链上验证对战正确性。
          </p>
        </details>
        <details className="rounded-xl border border-slate-200 bg-white p-4">
          <summary className="cursor-pointer font-medium text-slate-900">
            管理员在哪里管理权限？
          </summary>
          <p className="mt-3 text-sm leading-6 text-slate-600">
            管理员登录后，在赛事页面展开“管理员：管理主办方与审核员权限”，输入目标钱包后预览授予或撤销权限，再由钱包确认。
          </p>
        </details>
      </section>

      <div className="flex flex-wrap items-center justify-between gap-3 rounded-xl bg-slate-100 p-4">
        <p className="text-sm text-slate-600">
          想核查赞助、审核与付款？查看已确认的公开记录。
        </p>
        <Link
          href="/audit"
          className="inline-flex items-center gap-1 text-sm font-semibold text-emerald-700"
        >
          查看审计流水 <ArrowRight className="h-4 w-4" />
        </Link>
      </div>
    </div>
  );
}
