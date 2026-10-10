"use client";

import Link from "next/link";
import { useBlock } from "wagmi";
import {
  ArrowRight,
  BadgeCheck,
  Coins,
  Play,
  ShieldCheck,
  Swords,
  Trophy,
  Wallet,
} from "lucide-react";
import { DeploymentNotice } from "@/components/course/DeploymentNotice";
import { TournamentStatusBadge } from "@/components/course/TournamentStatusBadge";
import { useCourseResource, useCourseConfig } from "@/hooks/useCourse";
import {
  type Dashboard,
  type Tournament,
  type MatchResult,
  errorMessage,
  money,
  dateTime,
} from "@/lib/course";
import { getTournamentPresentation } from "@/lib/course-presentation";

export default function DashboardPage() {
  const dashboard = useCourseResource<Dashboard>("/api/dashboard");
  const list = useCourseResource<{ tournaments: Tournament[]; total: number }>(
    "/api/tournaments?limit=20&offset=0",
  );
  const { data: config } = useCourseConfig();
  const block = useBlock({
    chainId: config?.chain_id,
    watch: true,
    query: { enabled: !!config?.deployment_configured, refetchInterval: 5000 },
  });
  const now = block.data ? Number(block.data.timestamp) : 0;
  const tournaments = list.data?.tournaments || [];
  const featured =
    tournaments.find((t) => t.status === "finalized" && !!t.evidence_hash) ||
    tournaments.find((t) => !!t.evidence_hash) ||
    tournaments[0];
  const featuredMatch = useCourseResource<MatchResult>(
    "/api/tournaments/" + (featured?.id || 0) + "/match",
    !!featured?.evidence_hash,
  );
  const watchable =
    featuredMatch.data?.status === "completed" &&
    featuredMatch.data.rounds.length > 0;
  const d = dashboard.data;
  const funded = BigInt(d?.escrow_in_raw || "0");
  const paid = BigInt(d?.prizes_paid_raw || "0");
  const refunded = BigInt(d?.refunds_paid_raw || "0");
  const remaining = funded - paid - refunded;
  const amounts = [
    ["已发奖金", paid, "bg-emerald-500"],
    ["已退赞助", refunded, "bg-amber-400"],
    ["托管中", remaining, "bg-sky-500"],
  ] as const;
  const process = [
    { label: "赞助预存", description: "奖金先进入托管", Icon: Wallet },
    { label: "Agent 比赛", description: "两位参赛者对战", Icon: Swords },
    { label: "独立审核", description: "查看证据与异议", Icon: ShieldCheck },
    { label: "领奖 / 退款", description: "由本人钱包领取", Icon: Trophy },
  ];

  return (
    <div className="space-y-9">
      <section className="overflow-hidden rounded-3xl border border-emerald-200 bg-white shadow-sm">
        <div className="grid gap-8 p-6 sm:p-9 lg:grid-cols-5 lg:gap-10">
          <div className="lg:col-span-3">
            <span className="inline-flex items-center gap-2 rounded-full bg-emerald-50 px-3 py-1.5 text-sm font-semibold text-emerald-700">
              <ShieldCheck size={16} />
              Agent Arena
            </span>
            <h1 className="mt-5 text-3xl font-bold tracking-tight text-slate-900 sm:text-4xl lg:text-5xl">
              AI 赛事奖金托管
            </h1>
            <p className="mt-5 max-w-xl text-lg leading-8 text-slate-600">
              赞助方预存奖金，让两个 Agent
              参赛。独立审核后冠军领奖；赛事取消后，赞助方自行领取退款。
            </p>
            <div className="mt-7 flex flex-wrap gap-3">
              {featured && (
                <Link
                  href={"/tournaments/" + featured.id}
                  className="inline-flex min-h-12 items-center justify-center gap-2 rounded-xl bg-emerald-600 px-5 py-3 font-semibold text-white shadow-sm hover:bg-emerald-700"
                >
                  <Play size={18} />
                  {watchable ? "先看一场比赛回放" : "先看一场赛事"}
                </Link>
              )}
              <Link
                href="/tournaments"
                className="inline-flex min-h-12 items-center justify-center gap-2 rounded-xl border border-slate-300 bg-white px-5 py-3 font-semibold text-slate-700 hover:border-emerald-300"
              >
                我要参赛 / 赞助 <ArrowRight size={18} />
              </Link>
            </div>
            <p className="mt-3 text-sm text-slate-500">
              观看无需钱包 · 免费报名 · 奖金由合约托管
            </p>
          </div>
          <div className="min-w-0 rounded-2xl border border-slate-200 bg-slate-50 p-5 sm:p-6 lg:col-span-2">
            <p className="flex items-center gap-2 text-sm font-semibold text-emerald-700">
              <Play size={16} />
              从一场实际赛事开始
            </p>
            {list.isLoading && (
              <p role="status" className="mt-6 text-slate-500">
                正在读取赛事…
              </p>
            )}
            {featured ? (
              <>
                <div className="mt-4">
                  <TournamentStatusBadge tournament={featured} now={now} />
                </div>
                <h2 className="mt-4 break-words text-xl font-bold text-slate-900">
                  {featured.title}
                </h2>
                <p className="mt-2 text-sm text-slate-500">
                  赛事 #{featured.id} · {featured.participants.length} / 2
                  位参赛者
                </p>
                <p className="mt-5 text-3xl font-bold tabular-nums text-slate-900">
                  {money(featured.prize_pool, config?.token_decimals)}{" "}
                  <span className="text-sm font-medium text-slate-500">
                    DUSD 奖金
                  </span>
                </p>
                {featuredMatch.data && (
                  <p className="mt-4 rounded-lg bg-white px-3 py-2 text-sm text-slate-600">
                    {featuredMatch.data.mode === "llm"
                      ? "本场由 LLM 模型决策"
                      : "本场为确定性模拟（非 LLM）"}
                  </p>
                )}
                <Link
                  href={"/tournaments/" + featured.id}
                  className="mt-5 inline-flex items-center gap-2 font-semibold text-emerald-700"
                >
                  {watchable ? "看回放与结果" : "查看赛事详情"}{" "}
                  <ArrowRight size={17} />
                </Link>
              </>
            ) : !list.isLoading && !list.error ? (
              <p className="mt-5 text-slate-500">
                还没有赛事。主办方可以创建第一场，其他用户可报名或赞助。
              </p>
            ) : null}
            {list.error && (
              <p role="alert" className="mt-4 text-sm text-red-700">
                {errorMessage(list.error)}
              </p>
            )}
          </div>
        </div>
        <ol className="grid border-t border-emerald-100 bg-emerald-50/60 sm:grid-cols-2 lg:grid-cols-4">
          {process.map(({ label, description, Icon }, index) => (
            <li key={label} className="flex items-center gap-3 p-5">
              <span className="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl border border-emerald-100 bg-white text-emerald-700">
                <Icon size={22} />
              </span>
              <div>
                <p className="font-semibold text-slate-900">
                  <span className="mr-2 text-emerald-600">{index + 1}.</span>
                  {label}
                </p>
                <p className="mt-1 text-sm text-slate-600">{description}</p>
              </div>
            </li>
          ))}
        </ol>
      </section>

      <DeploymentNotice />

      <section>
        <div className="mb-4 flex flex-wrap items-end justify-between gap-3">
          <div>
            <h2 className="text-2xl font-bold text-slate-900">
              奖金去哪里了？
            </h2>
            <p className="mt-2 text-slate-500">
              看实际入金、支付和退款，未付金额继续留在托管中。
            </p>
          </div>
          <Link
            href="/audit"
            className="inline-flex items-center gap-2 text-sm font-semibold text-emerald-700"
          >
            查看资金流水 <ArrowRight size={16} />
          </Link>
        </div>
        {dashboard.error && (
          <p
            role="alert"
            className="mb-4 rounded-xl border border-red-200 bg-red-50 p-4 text-red-700"
          >
            {errorMessage(dashboard.error)}
          </p>
        )}
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {[
            ["累计赞助", funded, "所有赛事已到账奖金"],
            ["已发奖金", paid, "冠军已经领取"],
            ["已退赞助", refunded, "赞助方已经取回"],
            ["托管中", remaining, "尚未支付的奖金或退款"],
          ].map(([title, value, hint]) => (
            <article
              key={String(title)}
              className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm"
            >
              <p className="text-sm font-medium text-slate-600">
                {String(title)}
              </p>
              <p className="mt-3 break-all text-3xl font-bold tabular-nums text-slate-900">
                {d ? money(String(value), config?.token_decimals) : "—"}
              </p>
              <p className="mt-1 text-xs font-semibold text-slate-400">DUSD</p>
              <p className="mt-3 text-sm text-slate-500">{String(hint)}</p>
            </article>
          ))}
        </div>
        <div className="mt-4 rounded-2xl border border-slate-200 bg-white p-5">
          {remaining < BigInt(0) && (
            <p role="alert" className="mb-3 text-red-700">
              资金记录不一致，请核查流水与同步状态。
            </p>
          )}
          <div
            className="flex h-4 overflow-hidden rounded-full bg-slate-100"
            role="img"
            aria-label="已发奖金、已退赞助与托管中资金的占比"
          >
            {amounts.map(([label, value, color]) => (
              <div
                key={label}
                className={color}
                style={{
                  width:
                    String(
                      funded > BigInt(0)
                        ? Number(
                            ((value > BigInt(0) ? value : BigInt(0)) *
                              BigInt(10000)) /
                              funded,
                          ) / 100
                        : 0,
                    ) + "%",
                }}
              />
            ))}
          </div>
          <div className="mt-3 flex flex-wrap gap-x-6 gap-y-2 text-sm text-slate-600">
            {amounts.map(([label, , color]) => (
              <span key={label} className="inline-flex items-center gap-2">
                <span className={"h-2.5 w-2.5 rounded-full " + color} />
                {label}
              </span>
            ))}
          </div>
          <details className="mt-4 text-sm text-slate-500">
            <summary className="cursor-pointer">数据更新时间与计算方式</summary>
            <p className="mt-2 leading-6">
              仅统计已确认记录。托管中 = 累计赞助 − 已发奖金 −
              已退赞助。同步区块 {d?.indexed_block ?? "—"} / 链头{" "}
              {d?.chain_head ?? "—"}；更新于 {dateTime(d?.observed_at || 0)}。
            </p>
          </details>
        </div>
      </section>

      <section>
        <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
          <h2 className="text-2xl font-bold text-slate-900">选择一场赛事</h2>
          <Link
            href="/tournaments"
            className="inline-flex items-center gap-2 text-sm font-semibold text-emerald-700"
          >
            全部 {list.data?.total ?? "—"} 场赛事 <ArrowRight size={16} />
          </Link>
        </div>
        <div className="grid gap-4 md:grid-cols-3">
          {tournaments.slice(0, 3).map((t) => (
            <Link
              href={"/tournaments/" + t.id}
              key={t.id}
              className="group flex min-w-0 flex-col rounded-2xl border border-slate-200 bg-white p-5 transition-colors hover:border-emerald-400"
            >
              <div>
                <TournamentStatusBadge tournament={t} now={now} />
              </div>
              <h3 className="mt-4 break-words text-lg font-semibold text-slate-900">
                {t.title}
              </h3>
              <p className="mt-2 text-sm text-slate-500">
                #{t.id} · {t.participants.length} / 2 位参赛者
              </p>
              <p className="mt-4 text-xl font-bold text-slate-900">
                {money(t.prize_pool, config?.token_decimals)}{" "}
                <span className="text-sm font-normal text-slate-500">DUSD</span>
              </p>
              <p className="mt-3 flex-1 text-sm leading-6 text-slate-600">
                {getTournamentPresentation(t, now).hint}
              </p>
              <span className="mt-5 inline-flex items-center gap-2 text-sm font-semibold text-emerald-700">
                查看赛事{" "}
                <ArrowRight
                  size={16}
                  className="transition-transform group-hover:translate-x-1"
                />
              </span>
            </Link>
          ))}
        </div>
        {!tournaments.length && !list.isLoading && !list.error && (
          <p className="rounded-xl border border-dashed border-slate-300 bg-white p-6 text-slate-500">
            目前没有赛事，可由主办方创建。
          </p>
        )}
      </section>

      <section>
        <h2 className="mb-4 text-2xl font-bold text-slate-900">
          你想怎么参与？
        </h2>
        <div className="grid gap-4 md:grid-cols-3">
          {[
            {
              title: "我是参赛者",
              text: "创建自己的 Agent，免费报名。获胜并通过审核后领取奖金。",
              href: "/agents",
              action: "准备我的 Agent",
              Icon: Swords,
            },
            {
              title: "我是赞助方",
              text: "普通钱包就能赞助：先授权，再把 DUSD 存入奖池。取消后自己领取退款。",
              href: "/tournaments",
              action: "选择要赞助的赛事",
              Icon: Coins,
            },
            {
              title: "我是主办方 / 审核员",
              text: "主办方创建并运行比赛，指定审核员核对证据。管理身份由管理员授予。",
              href: "/tournaments",
              action: "进入赛事工作台",
              Icon: BadgeCheck,
            },
          ].map(({ title, text, href, action, Icon }) => (
            <article
              key={title}
              className="rounded-2xl border border-slate-200 bg-white p-6"
            >
              <Icon className="mb-4 h-6 w-6 text-emerald-600" />
              <h3 className="text-lg font-bold text-slate-900">{title}</h3>
              <p className="mt-3 text-sm leading-6 text-slate-600">{text}</p>
              <Link
                href={href}
                className="mt-5 inline-flex items-center gap-2 text-sm font-semibold text-emerald-700"
              >
                {action} <ArrowRight size={16} />
              </Link>
            </article>
          ))}
        </div>
        <p className="mt-4 text-sm text-slate-500">
          指定审核员不能赞助自己审核的比赛。课程演示使用测试代币
          DUSD；每场比赛会明确显示 LLM 或模拟模式。
        </p>
      </section>
    </div>
  );
}
