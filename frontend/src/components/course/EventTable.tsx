"use client";
import { useCourseConfig } from "@/hooks/useCourse";
import { type AuditEvent, dateTime, explorer, money, shortAddress } from "@/lib/course";
export function EventTable({ events }: {
    events: AuditEvent[];
}) {
    const { data: config } = useCourseConfig();
    if (!events.length)
        return <p className="text-sm text-muted-foreground py-4">暂无已确认的合约事件。</p>;
    return <div className="overflow-x-auto"><table className="w-full text-sm text-left"><thead><tr className="border-b">{["时间 / 区块", "赛事", "事件", "钱包", "DUSD 金额", "交易"].map(t => <th key={t} className="p-3 whitespace-nowrap">{t}</th>)}</tr></thead><tbody>{events.map(event => {
            const url = explorer(event.chain_id, event.tx_hash);
            return <tr key={event.id} className="border-b"><td className="p-3 whitespace-nowrap">{dateTime(event.timestamp)}<p className="text-xs text-muted-foreground">#{event.block_number} · log {event.log_index}</p></td><td className="p-3">#{event.tournament_id || "—"}</td><td className="p-3">{event.event_name}<details className="mt-1"><summary className="cursor-pointer text-xs underline">事件详情</summary><pre className="text-xs max-w-sm overflow-auto whitespace-pre-wrap break-all mt-2">{JSON.stringify(event.data, null, 2)}</pre></details></td><td className="p-3 font-mono" title={event.actor}>{shortAddress(event.actor)}</td><td className="p-3 font-mono">{event.amount_raw ? money(event.amount_raw, config?.token_decimals) : "—"}</td><td className="p-3 font-mono">{url ? <a className="underline" target="_blank" rel="noreferrer" href={url}>{shortAddress(event.tx_hash)}</a> : <span title={event.tx_hash}>{shortAddress(event.tx_hash)}</span>}</td></tr>;
        })}</tbody></table></div>;
}
