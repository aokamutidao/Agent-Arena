"use client";

import Link from "next/link";
import { useI18n } from "@/lib/i18n";
import { getEventActor } from "@/lib/course-presentation";
import { useState } from "react";
import {
  ArrowRight,
  Download,
  FileCheck2,
  Search,
  ShieldCheck,
} from "lucide-react";
import { useCourseResource } from "@/hooks/useCourse";
import { type Tournament, type AuditEvent, errorMessage } from "@/lib/course";
import { EventTable } from "@/components/course/EventTable";
import { DeploymentNotice } from "@/components/course/DeploymentNotice";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

const EVENT_SEARCH_LABELS: Record<string, string> = {
  TournamentCreated: "创建赛事 Create tournament",
  ParticipantRegistered: "参赛者报名 Participant registration",
  ParticipantWithdrawn: "退出报名 Withdraw registration",
  TournamentFunded: "赞助到账入金 Sponsorship deposit funded",
  TournamentStarted: "比赛开赛 Tournament started",
  ResultSubmitted: "提交比赛结果 Submit match result",
  ResultDisputed: "提出异议 Raise dispute",
  ResultApproved: "审核通过 Result approved",
  DisputeResolved: "完成异议裁决 Dispute resolved",
  TournamentCancelled: "取消赛事 Cancel tournament",
  PrizeClaimed: "冠军领取奖金 Winner claimed prize",
  SponsorRefunded: "赞助退款到账 Sponsor refunded",
};

export default function AuditPage() {
  const { lang, l } = useI18n();
  const list = useCourseResource<{ tournaments: Tournament[]; total: number }>(
    "/api/tournaments?limit=100&offset=0",
  );
  const [selected, setSelected] = useState("");
  const [offset, setOffset] = useState(0);
  const [filter, setFilter] = useState("");
  const id = selected || list.data?.tournaments[0]?.id.toString() || "";
  const events = useCourseResource<{ events: AuditEvent[]; total: number }>(
    "/api/tournaments/" + id + "/events?limit=50&offset=" + offset,
    !!id,
  );
  const visible =
    events.data?.events.filter((event) =>
      [
        EVENT_SEARCH_LABELS[event.event_name],
        event.event_name,
        getEventActor(event),
        event.tx_hash,
      ]
        .join(" ")
        .toLowerCase()
        .includes(filter.toLowerCase()),
    ) || [];

  function download() {
    const url = URL.createObjectURL(
      new Blob([JSON.stringify(events.data?.events || [], null, 2)], {
        type: "application/json",
      }),
    );
    const anchor = document.createElement("a");
    anchor.href = url;
    anchor.download = "tournament-" + id + "-events-" + offset + ".json";
    anchor.click();
    URL.revokeObjectURL(url);
  }

  return (
    <div className="space-y-6">
      <header className="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="text-sm font-semibold text-emerald-700">
            {l("每一笔，都有记录", "Every action has a record")}
          </p>
          <h1 className="mt-2 text-3xl font-bold tracking-tight text-slate-900">
            {l("审计流水", "Audit log")}
          </h1>
          <p className="mt-2 max-w-2xl text-slate-600">
            {l(
              "查看谁赞助了赛事、结果如何审核、奖金与退款是否到账。所有人都能查看，无需钱包。",
              "See who sponsored a tournament, how results were reviewed, and whether prizes and refunds were claimed. Anyone can view the records without a wallet.",
            )}
          </p>
        </div>
        <Link
          href="/tournaments"
          className="inline-flex shrink-0 items-center gap-1 text-sm font-semibold text-emerald-700"
        >
          {l("返回赛事", "Back to tournaments")}{" "}
          <ArrowRight className="h-4 w-4" />
        </Link>
      </header>

      <DeploymentNotice />

      <div className="grid gap-3 sm:grid-cols-2">
        <div className="flex gap-3 rounded-xl border border-slate-200 bg-white p-4">
          <ShieldCheck className="h-5 w-5 shrink-0 text-emerald-600" />
          <div>
            <p className="font-medium text-slate-900">
              {l("只显示已确认的操作", "Confirmed actions only")}
            </p>
            <p className="mt-1 text-sm text-slate-600">
              {l(
                "钱包尚未发送或正在确认的交易，需要等待后才会出现在这里。",
                "Transactions that have not been sent or are still confirming will appear here only after confirmation.",
              )}
            </p>
          </div>
        </div>
        <div className="flex gap-3 rounded-xl border border-slate-200 bg-white p-4">
          <FileCheck2 className="h-5 w-5 shrink-0 text-emerald-600" />
          <div>
            <p className="font-medium text-slate-900">
              {l("记录可核查、可导出", "Verifiable, exportable records")}
            </p>
            <p className="mt-1 text-sm text-slate-600">
              {l(
                "展开交易详情查看凭证，也可以下载当前页原始记录。",
                "Expand transaction details to inspect the evidence, or download this page’s original records.",
              )}
            </p>
          </div>
        </div>
      </div>

      <section className="rounded-2xl border border-slate-200 bg-white p-4 sm:p-5">
        <div className="grid gap-4 lg:grid-cols-[1fr_1fr_auto] lg:items-end">
          <label className="min-w-0 space-y-2 text-sm font-medium text-slate-700">
            <span>{l("选择赛事", "Select a tournament")}</span>
            <select
              className="block w-full min-w-0 rounded-lg border border-slate-200 bg-white p-2.5 text-sm"
              value={id}
              disabled={!list.data?.tournaments.length}
              onChange={(event) => {
                setSelected(event.target.value);
                setOffset(0);
              }}
            >
              {!list.data?.tournaments.length && (
                <option value="">{l("暂无赛事", "No tournaments")}</option>
              )}
              {list.data?.tournaments.map((t) => (
                <option key={t.id} value={t.id}>
                  #{t.id} {t.title}
                </option>
              ))}
            </select>
          </label>
          <label className="min-w-0 space-y-2 text-sm font-medium text-slate-700">
            <span className="flex items-center gap-1.5">
              <Search className="h-4 w-4" />{" "}
              {l("搜索当前页", "Search this page")}
            </span>
            <Input
              aria-label={l("筛选当前页事件", "Filter events on this page")}
              placeholder={l(
                "例如：赞助、退款、钱包地址",
                "e.g. sponsorship, refund, wallet address",
              )}
              value={filter}
              onChange={(event) => setFilter(event.target.value)}
            />
          </label>
          <Button variant="outline" disabled={!events.data} onClick={download}>
            <Download className="mr-2 h-4 w-4" />{" "}
            {l("导出当前页 JSON", "Export this page as JSON")}
          </Button>
        </div>
        <p className="mt-3 text-xs text-slate-500">
          {l(
            "显示最新 100 场赛事，每页 50 条记录。搜索仅筛选本页，导出保留本页全部原始记录。",
            "Shows the latest 100 tournaments, with 50 records per page. Search filters this page only; export includes all original records on the page.",
          )}
        </p>
      </section>

      {list.error && (
        <p role="alert" className="rounded-xl bg-red-50 p-4 text-red-700">
          {errorMessage(list.error, lang)}
        </p>
      )}
      {events.error && (
        <p role="alert" className="rounded-xl bg-red-50 p-4 text-red-700">
          {errorMessage(events.error, lang)}
        </p>
      )}
      {list.isLoading && (
        <p className="text-slate-600">
          {l("正在读取赛事…", "Loading tournaments…")}
        </p>
      )}
      {events.isLoading && id && (
        <p className="text-slate-600">
          {l("正在读取操作记录…", "Loading action records…")}
        </p>
      )}

      <section
        className="space-y-3"
        aria-label={l("已确认的赛事记录", "Confirmed tournament records")}
      >
        <div className="flex flex-wrap items-center justify-between gap-2">
          <h2 className="text-lg font-semibold text-slate-900">
            {id
              ? l(
                  "赛事 #" + id + " 的操作记录",
                  "Action records for tournament #" + id,
                )
              : l("操作记录", "Action records")}
          </h2>
          <span className="text-sm text-slate-500">
            {events.data
              ? l(
                  "共 " + events.data.total + " 条",
                  events.data.total + " records",
                )
              : "—"}
            {filter &&
              l(
                " · 本页匹配 " + visible.length + " 条",
                " · " + visible.length + " matches on this page",
              )}
          </span>
        </div>
        {filter && !visible.length && events.data ? (
          <div className="rounded-xl border border-dashed border-slate-300 bg-white p-6 text-sm text-slate-600">
            {l(
              "本页没有匹配记录。换一个关键词，或清空搜索后翻页查看。",
              "No matching records on this page. Try another keyword, or clear the search and browse other pages.",
            )}
          </div>
        ) : (
          <EventTable events={visible} />
        )}
      </section>

      <div className="flex items-center justify-center gap-3">
        <Button
          variant="outline"
          disabled={offset === 0}
          onClick={() => setOffset(Math.max(0, offset - 50))}
        >
          {l("上一页", "Previous")}
        </Button>
        <span className="text-sm text-slate-500">
          {l(
            "第 " + (Math.floor(offset / 50) + 1) + " 页",
            "Page " + (Math.floor(offset / 50) + 1),
          )}
        </span>
        <Button
          variant="outline"
          disabled={offset + 50 >= (events.data?.total || 0)}
          onClick={() => setOffset(offset + 50)}
        >
          {l("下一页", "Next")}
        </Button>
      </div>
    </div>
  );
}
