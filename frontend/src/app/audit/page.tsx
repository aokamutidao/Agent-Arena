"use client";
import { useState } from "react";
import { useCourseResource } from "@/hooks/useCourse";
import { type Tournament, type AuditEvent, errorMessage } from "@/lib/course";
import { EventTable } from "@/components/course/EventTable";
import { DeploymentNotice } from "@/components/course/DeploymentNotice";
import { Button } from "@/components/ui/button";
export default function AuditPage() {
    const list = useCourseResource<{
        tournaments: Tournament[];
        total: number;
    }>("/api/tournaments?limit=100&offset=0");
    const [selected, setSelected] = useState("");
    const [offset, setOffset] = useState(0);
    const [filter, setFilter] = useState("");
    const id = selected || list.data?.tournaments[0]?.id.toString() || "";
    const events = useCourseResource<{
        events: AuditEvent[];
        total: number;
    }>(`/api/tournaments/${id}/events?limit=50&offset=${offset}`, !!id);
    const visible = events.data?.events.filter(event => `${event.event_name} ${event.actor} ${event.tx_hash}`.toLowerCase().includes(filter.toLowerCase())) || [];
    function download() { const url = URL.createObjectURL(new Blob([JSON.stringify(events.data?.events || [], null, 2)], { type: "application/json" })); const a = document.createElement("a"); a.href = url; a.download = `tournament-${id}-events-${offset}.json`; a.click(); URL.revokeObjectURL(url); }
    return <div className="space-y-6"><h1 className="text-3xl font-bold">审计流水</h1><DeploymentNotice /><p className="text-muted-foreground">公开查看已确认的赞助、结果审核、领奖与退款事件。事件按链、交易哈希和日志序号唯一索引。</p><div className="flex flex-wrap gap-3"><label className="text-sm">赛事<select className="ml-2 border rounded p-2 bg-background" value={id} onChange={e => { setSelected(e.target.value); setOffset(0); }}>{list.data?.tournaments.map(t => <option key={t.id} value={t.id}>#{t.id} {t.title}</option>)}</select></label><input className="border rounded p-2 bg-background" aria-label="筛选当前页事件" placeholder="事件 / 地址 / 交易哈希" value={filter} onChange={e => setFilter(e.target.value)}/><Button variant="outline" disabled={!events.data} onClick={download}>导出当前页 JSON</Button></div>{list.error && <p role="alert">{errorMessage(list.error)}</p>}{events.error && <p role="alert">{errorMessage(events.error)}</p>}{events.isLoading && id && <p>加载事件中…</p>}<section className="border rounded-xl"><EventTable events={visible}/></section><p className="text-xs text-muted-foreground">赛事选择器显示最新 100 场；筛选仅作用于当前 50 条事件。总事件数 {events.data?.total ?? 0}。</p><div className="flex gap-3"><Button variant="outline" disabled={offset === 0} onClick={() => setOffset(Math.max(0, offset - 50))}>上一页</Button><Button variant="outline" disabled={offset + 50 >= (events.data?.total || 0)} onClick={() => setOffset(offset + 50)}>下一页</Button></div></div>;
}
