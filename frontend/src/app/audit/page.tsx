"use client";

import Link from "next/link";
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
  TournamentCreated: "创建赛事",
  ParticipantRegistered: "参赛者报名",
  ParticipantWithdrawn: "退出报名",
  TournamentFunded: "赞助到账入金",
  TournamentStarted: "比赛开赛",
  ResultSubmitted: "提交比赛结果",
  ResultDisputed: "提出异议",
  ResultApproved: "审核通过",
  DisputeResolved: "完成异议裁决",
  TournamentCancelled: "取消赛事",
  PrizeClaimed: "冠军领取奖金",
  SponsorRefunded: "赞助退款到账",
};

export default function AuditPage() {
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
            每一笔，都有记录
          </p>
          <h1 className="mt-2 text-3xl font-bold tracking-tight text-slate-900">
            审计流水
          </h1>
          <p className="mt-2 max-w-2xl text-slate-600">
            查看谁赞助了赛事、结果如何审核、奖金与退款是否到账。所有人都能查看，无需钱包。
          </p>
        </div>
        <Link
          href="/tournaments"
          className="inline-flex shrink-0 items-center gap-1 text-sm font-semibold text-emerald-700"
        >
          返回赛事 <ArrowRight className="h-4 w-4" />
        </Link>
      </header>

      <DeploymentNotice />

      <div className="grid gap-3 sm:grid-cols-2">
        <div className="flex gap-3 rounded-xl border border-slate-200 bg-white p-4">
          <ShieldCheck className="h-5 w-5 shrink-0 text-emerald-600" />
          <div>
            <p className="font-medium text-slate-900">只显示已确认的操作</p>
            <p className="mt-1 text-sm text-slate-600">
              钱包尚未发送或正在确认的交易，需要等待后才会出现在这里。
            </p>
          </div>
        </div>
        <div className="flex gap-3 rounded-xl border border-slate-200 bg-white p-4">
          <FileCheck2 className="h-5 w-5 shrink-0 text-emerald-600" />
          <div>
            <p className="font-medium text-slate-900">记录可核查、可导出</p>
            <p className="mt-1 text-sm text-slate-600">
              展开交易详情查看凭证，也可以下载当前页原始记录。
            </p>
          </div>
        </div>
      </div>

      <section className="rounded-2xl border border-slate-200 bg-white p-4 sm:p-5">
        <div className="grid gap-4 lg:grid-cols-[1fr_1fr_auto] lg:items-end">
          <label className="min-w-0 space-y-2 text-sm font-medium text-slate-700">
            <span>选择赛事</span>
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
                <option value="">暂无赛事</option>
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
              <Search className="h-4 w-4" /> 搜索当前页
            </span>
            <Input
              aria-label="筛选当前页事件"
              placeholder="例如：赞助、退款、钱包地址"
              value={filter}
              onChange={(event) => setFilter(event.target.value)}
            />
          </label>
          <Button variant="outline" disabled={!events.data} onClick={download}>
            <Download className="mr-2 h-4 w-4" /> 导出当前页 JSON
          </Button>
        </div>
        <p className="mt-3 text-xs text-slate-500">
          显示最新 100 场赛事，每页 50
          条记录。搜索仅筛选本页，导出保留本页全部原始记录。
        </p>
      </section>

      {list.error && (
        <p role="alert" className="rounded-xl bg-red-50 p-4 text-red-700">
          {errorMessage(list.error)}
        </p>
      )}
      {events.error && (
        <p role="alert" className="rounded-xl bg-red-50 p-4 text-red-700">
          {errorMessage(events.error)}
        </p>
      )}
      {list.isLoading && <p className="text-slate-600">正在读取赛事…</p>}
      {events.isLoading && id && (
        <p className="text-slate-600">正在读取操作记录…</p>
      )}

      <section className="space-y-3" aria-label="已确认的赛事记录">
        <div className="flex flex-wrap items-center justify-between gap-2">
          <h2 className="text-lg font-semibold text-slate-900">
            {id ? "赛事 #" + id + " 的操作记录" : "操作记录"}
          </h2>
          <span className="text-sm text-slate-500">
            {events.data ? "共 " + events.data.total + " 条" : "—"}
            {filter && " · 本页匹配 " + visible.length + " 条"}
          </span>
        </div>
        {filter && !visible.length && events.data ? (
          <div className="rounded-xl border border-dashed border-slate-300 bg-white p-6 text-sm text-slate-600">
            本页没有匹配记录。换一个关键词，或清空搜索后翻页查看。
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
          上一页
        </Button>
        <span className="text-sm text-slate-500">
          第 {Math.floor(offset / 50) + 1} 页
        </span>
        <Button
          variant="outline"
          disabled={offset + 50 >= (events.data?.total || 0)}
          onClick={() => setOffset(offset + 50)}
        >
          下一页
        </Button>
      </div>
    </div>
  );
}
