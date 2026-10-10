"use client";

import Link from "next/link";
import { getEventActor } from "@/lib/course-presentation";
import { useCourseConfig } from "@/hooks/useCourse";
import {
  type AuditEvent,
  dateTime,
  explorer,
  money,
  shortAddress,
} from "@/lib/course";

const EVENT_LABELS: Record<string, string> = {
  TournamentCreated: "创建赛事",
  ParticipantRegistered: "参赛者报名",
  ParticipantWithdrawn: "退出报名",
  TournamentFunded: "赞助到账",
  TournamentStarted: "比赛开赛",
  ResultSubmitted: "提交比赛结果",
  ResultDisputed: "提出异议",
  ResultApproved: "审核通过",
  DisputeResolved: "完成异议裁决",
  TournamentCancelled: "取消赛事",
  PrizeClaimed: "冠军领取奖金",
  SponsorRefunded: "赞助退款到账",
};

export function EventTable({ events }: { events: AuditEvent[] }) {
  const { data: config } = useCourseConfig();
  if (!events.length) {
    return (
      <p className="rounded-xl bg-muted p-5 text-sm text-muted-foreground">
        暂无已确认的操作。交易确认后，记录会自动出现在这里。
      </p>
    );
  }
  return (
    <div className="overflow-x-auto rounded-xl border">
      <table className="w-full text-left text-sm">
        <caption className="sr-only">已确认的链上资金与操作记录</caption>
        <thead className="bg-muted/60 text-xs text-muted-foreground">
          <tr>
            {[
              "时间",
              "赛事",
              "操作",
              "操作钱包",
              "金额（DUSD）",
              "交易凭证",
            ].map((label) => (
              <th
                scope="col"
                key={label}
                className="whitespace-nowrap p-3 font-medium"
              >
                {label}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {events.map((event) => {
            const url = explorer(event.chain_id, event.tx_hash);
            const actor = getEventActor(event);
            const payment = ["PrizeClaimed", "SponsorRefunded"].includes(
              event.event_name,
            );
            return (
              <tr
                key={event.id}
                className="border-t align-top hover:bg-muted/30"
              >
                <td className="whitespace-nowrap p-3 text-xs text-muted-foreground">
                  {dateTime(event.timestamp)}
                </td>
                <td className="p-3">
                  {event.tournament_id ? (
                    <Link
                      href={`/tournaments/${event.tournament_id}`}
                      className="font-medium text-primary hover:underline"
                    >
                      #{event.tournament_id}
                    </Link>
                  ) : (
                    "—"
                  )}
                </td>
                <td className="min-w-36 p-3">
                  <span className="font-medium">
                    {EVENT_LABELS[event.event_name] || event.event_name}
                  </span>
                  <details className="mt-2 text-xs text-muted-foreground">
                    <summary className="cursor-pointer">技术详情</summary>
                    <p className="mt-2">{event.event_name}</p>
                    <p className="mt-1">
                      区块 #{event.block_number} · 日志 {event.log_index}
                    </p>
                    <pre className="mt-2 max-w-sm overflow-auto whitespace-pre-wrap break-all">
                      {JSON.stringify(event.data, null, 2)}
                    </pre>
                  </details>
                </td>
                <td className="p-3 font-mono text-xs" title={actor}>
                  {shortAddress(actor)}
                </td>
                <td
                  className={`whitespace-nowrap p-3 font-semibold tabular-nums ${payment ? "text-emerald-700" : ""}`}
                >
                  {event.amount_raw
                    ? money(event.amount_raw, config?.token_decimals)
                    : "—"}
                </td>
                <td className="p-3 font-mono text-xs">
                  {url ? (
                    <a
                      className="text-primary underline"
                      target="_blank"
                      rel="noreferrer"
                      href={url}
                    >
                      {shortAddress(event.tx_hash)} ↗
                    </a>
                  ) : (
                    <details>
                      <summary className="cursor-pointer whitespace-nowrap">
                        查看哈希
                      </summary>
                      <p className="mt-2 max-w-48 break-all">{event.tx_hash}</p>
                    </details>
                  )}
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}
