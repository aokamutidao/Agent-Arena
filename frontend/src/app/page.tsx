"use client";
import Link from "next/link";
import { DeploymentNotice } from "@/components/course/DeploymentNotice";
import { useCourseResource, useCourseConfig } from "@/hooks/useCourse";
import { type Dashboard, type Tournament, errorMessage, money, STATUS_LABEL, dateTime } from "@/lib/course";
export default function DashboardPage() {
    const dashboard = useCourseResource<Dashboard>("/api/dashboard");
    const list = useCourseResource<{
        tournaments: Tournament[];
        total: number;
    }>("/api/tournaments?limit=5&offset=0");
    const { data: config } = useCourseConfig();
    const d = dashboard.data;
    const funded = BigInt(d?.escrow_in_raw || "0");
    const paid = BigInt(d?.prizes_paid_raw || "0");
    const refunded = BigInt(d?.refunds_paid_raw || "0");
    const remaining = funded - paid - refunded;
    const amounts = [["已发奖金", paid, "bg-emerald-500"], ["已退赞助", refunded, "bg-amber-500"], ["待结算资金", remaining, "bg-blue-500"]] as const;
    return <div className="space-y-6"><div className="flex flex-wrap justify-between items-end gap-4"><div><p className="text-emerald-600 text-sm">赛事资金与独立审核</p><h1 className="text-3xl font-bold mt-1">资金概览</h1><p className="text-muted-foreground mt-2">报名、赞助、结果审核与退款均通过钱包签名，在链上留下审计记录。</p></div><Link href="/tournaments" className="rounded bg-emerald-600 text-white px-4 py-2">进入赛事</Link></div><DeploymentNotice />
    {dashboard.error && <p role="alert" className="text-red-500">{errorMessage(dashboard.error)}</p>}
    {dashboard.isLoading && <p>正在读取已确认的账本…</p>}
    <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-4">{[["赛事总数", d?.tournaments_total ?? "—"], ["进行中的赛事", d?.active_tournaments ?? "—"], ["已结算赛事", d?.completed_tournaments ?? "—"], ["累计赞助 DUSD", d ? money(d.escrow_in_raw, config?.token_decimals) : "—"]].map(([title, value]) => <div key={title} className="border rounded-xl p-5"><p className="text-sm text-muted-foreground">{title}</p><p className="text-2xl font-semibold mt-2 break-all">{value}</p></div>)}</div>
    <section className="border rounded-xl p-5 space-y-4"><h2 className="text-lg font-semibold">资金去向</h2>{remaining < BigInt(0) && <p role="alert" className="text-red-500">账本总额不一致，请检查索引状态与资金流水。</p>}<div className="h-5 rounded-full overflow-hidden flex bg-muted" role="img" aria-label="已发奖金、已退赞助和待结算资金占比">{amounts.map(([label, value, color]) => <div key={label} className={color} style={{ width: `${funded > BigInt(0) ? Number((value > BigInt(0) ? value : BigInt(0)) * BigInt(10000) / funded) / 100 : 0}%` }}/>)}</div><div className="grid sm:grid-cols-3 gap-4">{amounts.map(([label, value, color]) => <div key={label}><span className={`inline-block h-2 w-2 rounded-full mr-2 ${color}`}/>{label}<p className="font-mono mt-1">{money(value.toString(), config?.token_decimals)} DUSD</p></div>)}</div><p className="text-xs text-muted-foreground">仅统计已确认应用事件；待结算余额 = 累计赞助 − 奖金支付 − 退款。索引区块 {d?.indexed_block ?? "—"} / 链头 {d?.chain_head ?? "—"}，更新于 {dateTime(d?.observed_at || 0)}。</p></section>
    <section className="border rounded-xl p-5"><h2 className="text-lg font-semibold mb-4">最近赛事</h2>{list.error && <p role="alert">{errorMessage(list.error)}</p>}{!list.data?.tournaments.length && !list.isLoading && <p className="text-muted-foreground">暂无赛事。组织者配置角色后可创建第一场赛事。</p>}<div className="space-y-3">{list.data?.tournaments.map(t => <Link href={`/tournaments/${t.id}`} key={t.id} className="flex flex-wrap items-center justify-between border rounded p-3 gap-2"><span>#{t.id} {t.title}</span><span className="text-sm text-muted-foreground">{STATUS_LABEL[t.status]} · {money(t.prize_pool, config?.token_decimals)} DUSD</span></Link>)}</div></section><p className="text-xs text-muted-foreground">DUSD 是课程演示代币，没有兑付承诺。此平台采用独立审计员审核结果，不提供观众投注。</p></div>;
}
