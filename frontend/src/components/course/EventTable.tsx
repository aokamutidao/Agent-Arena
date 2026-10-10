"use client";

import Link from "next/link";
import { useI18n } from "@/lib/i18n";
import { getEventActor } from "@/lib/course-presentation";
import { useCourseConfig } from "@/hooks/useCourse";
import {
  type AuditEvent,
  dateTime,
  explorer,
  money,
  shortAddress,
} from "@/lib/course";

export function EventTable({ events }: { events: AuditEvent[] }) {
  const { lang, l } = useI18n();
  const EVENT_LABELS: Record<string, string> = {
    TournamentCreated: l("创建赛事", "Tournament created"),
    ParticipantRegistered: l("参赛者报名", "Participant registered"),
    ParticipantWithdrawn: l("退出报名", "Registration withdrawn"),
    TournamentFunded: l("赞助到账", "Sponsorship received"),
    TournamentStarted: l("比赛开赛", "Match started"),
    ResultSubmitted: l("提交比赛结果", "Result submitted"),
    ResultDisputed: l("提出异议", "Dispute raised"),
    ResultApproved: l("审核通过", "Result approved"),
    DisputeResolved: l("完成异议裁决", "Dispute resolved"),
    TournamentCancelled: l("取消赛事", "Tournament cancelled"),
    PrizeClaimed: l("冠军领取奖金", "Prize claimed"),
    SponsorRefunded: l("赞助退款到账", "Sponsor refunded"),
  };

  const { data: config } = useCourseConfig();
  if (!events.length) {
    return (
      <p className="rounded-xl bg-muted p-5 text-sm text-muted-foreground">
        {l(
          "暂无已确认的操作。交易确认后，记录会自动出现在这里。",
          "No confirmed activity yet. Records appear here automatically after confirmation.",
        )}
      </p>
    );
  }
  return (
    <div className="overflow-x-auto rounded-xl border">
      <table className="w-full text-left text-sm">
        <caption className="sr-only">
          {l(
            "已确认的链上资金与操作记录",
            "Confirmed on-chain funding and activity records",
          )}
        </caption>
        <thead className="bg-muted/60 text-xs text-muted-foreground">
          <tr>
            {[
              l("时间", "Time"),
              l("赛事", "Tournament"),
              l("操作", "Action"),
              l("操作钱包", "Action wallet"),
              l("金额（DUSD）", "Amount (DUSD)"),
              l("交易凭证", "Transaction"),
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
                  {dateTime(event.timestamp, lang)}
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
                    <summary className="cursor-pointer">
                      {l("技术详情", "Technical details")}
                    </summary>
                    <p className="mt-2">{event.event_name}</p>
                    <p className="mt-1">
                      {l(
                        `区块 #${event.block_number} · 日志 ${event.log_index}`,
                        `Block #${event.block_number} · Log ${event.log_index}`,
                      )}
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
                        {l("查看哈希", "View hash")}
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
