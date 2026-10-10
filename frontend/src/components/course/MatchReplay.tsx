"use client";
import { useEffect, useState } from "react";
import { useI18n } from "@/lib/i18n";
import { Pause, Play, SkipBack, SkipForward } from "lucide-react";
import { Button } from "@/components/ui/button";
import { ArenaBoard } from "@/components/game/ArenaBoard";
import { AgentCard } from "@/components/game/AgentCard";
import { TurnLog } from "@/components/game/TurnLog";
import {
  type MatchResult,
  type EngineAgentState,
  type TournamentStatus,
  shortAddress,
  matchExecutionPresentation,
  llmErrorMessage,
} from "@/lib/course";
import type { GameDetail, AgentInfo, AgentState } from "@/types/game";
export function MatchReplay({
  match,
  settlementStatus,
}: {
  match: MatchResult;
  settlementStatus: TournamentStatus;
}) {
  const { lang, l } = useI18n();
  const execution = matchExecutionPresentation(match, lang);
  const summary = match.llm_summary;
  const number = (value: number | null | undefined) =>
    value == null
      ? l("未返回", "Not reported")
      : new Intl.NumberFormat(lang === "en" ? "en-GB" : "zh-CN").format(value);
  const seconds = (milliseconds: number) =>
    `${(milliseconds / 1000).toFixed(2)} ${l("秒", "s")}`;
  const [frame, setFrame] = useState<number | null>(null);
  const [playing, setPlaying] = useState(false);
  useEffect(() => {
    if (!playing) return;
    const timer = window.setInterval(
      () =>
        setFrame((current) =>
          Math.min((current ?? 0) + 1, Math.max(0, match.rounds.length - 1)),
        ),
      900,
    );
    return () => window.clearInterval(timer);
  }, [playing, match.rounds.length]);
  useEffect(() => {
    if (frame !== null && frame >= match.rounds.length - 1) setPlaying(false);
  }, [frame, match.rounds.length]);
  const settlementNote =
    settlementStatus === "finalized"
      ? l(
          "本场已链上结算，审定冠军见资金区",
          "Settled on-chain; see the funds panel for the approved winner",
        )
      : settlementStatus === "cancelled"
        ? l(
            "赛事已取消，不发放冠军奖金",
            "Tournament cancelled; no winner payout",
          )
        : l("等待链上审核", "Awaiting on-chain review");
  const index =
    frame === null
      ? Math.max(0, match.rounds.length - 1)
      : Math.min(frame, Math.max(0, match.rounds.length - 1));
  const round = match.rounds[index];
  const agent = (address: string, name: string): AgentInfo => ({
    id: address,
    name,
    personality: shortAddress(address),
    wins: 0,
    losses: 0,
    win_rate: 0,
    description: "",
  });
  const replayState = (
    state: EngineAgentState | undefined,
    hp: number | undefined,
    position: {
      x: number;
      y: number;
    },
  ): AgentState => ({
    hp: state?.hp ?? hp ?? 100,
    max_hp: state?.max_hp ?? 100,
    position: state?.position || position,
    status: state?.status?.map((effect) => effect.type) || [],
    skill_cooldown: state?.skill_cooldown || 0,
    is_charging: state?.is_charging || false,
  });
  const game: GameDetail = {
    game_id: match.tournament_id,
    agent_red: agent(
      match.red_address,
      match.red_agent?.name || l("红方 Agent", "Red Agent"),
    ),
    agent_blue: agent(
      match.blue_address,
      match.blue_agent?.name || l("蓝方 Agent", "Blue Agent"),
    ),
    status:
      match.status === "completed" && index === match.rounds.length - 1
        ? "finished"
        : "playing",
    current_round: round?.round || 0,
    max_rounds: 30,
    total_bet_red: "0",
    total_bet_blue: "0",
    agent_red_state: replayState(
      round?.red_state,
      round?.red_hp,
      round?.red_position || { x: 1, y: 1 },
    ),
    agent_blue_state: replayState(
      round?.blue_state,
      round?.blue_hp,
      round?.blue_position || { x: 8, y: 8 },
    ),
    history: match.rounds.slice(0, index + 1).map((r) => ({
      round: r.round,
      red_action: r.red_action,
      blue_action: r.blue_action,
      red_hp_after: r.red_hp,
      blue_hp_after: r.blue_hp,
      red_reasoning: r.red_reasoning,
      blue_reasoning: r.blue_reasoning,
    })),
  };
  const matchLabel =
    (
      {
        queued: l("等待运行", "Queued"),
        running: l("对战进行中", "Match in progress"),
        completed: l("对战完成", "Match complete"),
        failed: l("执行失败", "Execution failed"),
        interrupted: l("执行中断", "Interrupted"),
      } as Record<string, string>
    )[match.status] || match.status;
  return (
    <section className="rounded-2xl border bg-card p-5 sm:p-6 space-y-5">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <h2 className="text-lg font-bold">{l("比赛回放", "Match replay")}</h2>
        <span className="rounded-full bg-muted px-3 py-1 text-xs font-medium">
          {matchLabel}
        </span>
      </div>
      <p className="text-sm text-muted-foreground">
        <span className="font-semibold text-foreground">{execution.label}</span>
        {" · "}{execution.hint}
      </p>
      {match.mode === "llm" && summary && (
        <div className="space-y-4 rounded-xl border bg-slate-50 p-4">
          <div className="flex flex-wrap items-start justify-between gap-3">
            <div className="min-w-0">
              <p className="text-xs font-medium text-muted-foreground">
                {l("记录的模型", "Recorded model")}
              </p>
              <p className="break-words font-semibold">
                {summary.model || l("模型名未记录", "Model name not recorded")}
              </p>
              <p className="text-xs text-muted-foreground">
                {summary.provider || l("提供方未记录", "Provider not recorded")}
              </p>
            </div>
            <p className="text-xs text-muted-foreground">
              {l("累计调用耗时", "Cumulative call time")}{" "}
              <span className="font-semibold text-foreground">
                {seconds(summary.latency_ms)}
              </span>
            </p>
          </div>
          <dl className="grid grid-cols-3 gap-3">
            {[
              [l("调用次数", "Attempts"), summary.attempted],
              [l("成功调用", "Successful calls"), summary.succeeded],
              [l("失败调用", "Failed calls"), summary.failed],
            ].map(([label, value]) => (
              <div key={label} className="rounded-lg border bg-white p-3">
                <dt className="text-xs text-muted-foreground">{label}</dt>
                <dd className="mt-1 text-xl font-semibold">{number(Number(value))}</dd>
              </div>
            ))}
          </dl>
          <div className="text-xs text-muted-foreground">
            <p>
              {l("Token 用量（输入 / 输出 / 总计）", "Token usage (input / output / total)")}
              {": "}{number(summary.usage.prompt_tokens)} / {number(summary.usage.completion_tokens)} / {number(summary.usage.total_tokens)}
            </p>
            <p className="mt-1">
              {l("仅展示模型服务返回的统计；未返回的用量不会记作零，失败尝试也保留记录。", "Only provider-reported usage is shown. Missing usage is not counted as zero, and failed attempts remain in the record.")}
            </p>
          </div>
          {summary.last_error_code && (
            <p className="text-xs text-amber-800">
              {l("最近记录的失败：", "Last recorded failure: ")}
              {llmErrorMessage(summary.last_error_code, lang)}
            </p>
          )}
          {!!match.llm_calls?.length && (
            <details>
              <summary className="cursor-pointer text-sm font-medium">
                {l("查看每次模型调用", "Inspect individual model calls")}
              </summary>
              <div className="mt-3 max-h-72 overflow-auto rounded-lg border bg-white">
                <table className="w-full min-w-[650px] text-left text-xs">
                  <thead className="bg-slate-100 text-muted-foreground">
                    <tr>
                      <th className="p-3">{l("调用 / 回合", "Call / round")}</th>
                      <th className="p-3">{l("参赛方 / 模型", "Side / model")}</th>
                      <th className="p-3">{l("结果 / 耗时", "Result / time")}</th>
                      <th className="p-3">{l("Token 用量", "Token usage")}</th>
                    </tr>
                  </thead>
                  <tbody>
                    {match.llm_calls.map((call) => (
                      <tr key={call.sequence} className="border-t align-top">
                        <td className="p-3">#{call.sequence} / R{call.round}</td>
                        <td className="p-3">
                          <p className={call.side === "red" ? "text-rose-700" : "text-blue-700"}>
                            {call.side === "red" ? l("红方", "Red") : l("蓝方", "Blue")}
                          </p>
                          <p className="mt-1 break-all">{call.model}</p>
                        </td>
                        <td className="p-3">
                          <p className={call.success ? "text-emerald-700" : "text-red-700"}>
                            {call.success ? l("成功", "Succeeded") : l("失败", "Failed")} · {seconds(call.latency_ms)}
                          </p>
                          {call.error_code && <p className="mt-1 max-w-xs">{llmErrorMessage(call.error_code, lang)}</p>}
                          {(call.http_status != null || call.request_id) && (
                            <details className="mt-1 text-muted-foreground">
                              <summary className="cursor-pointer">{l("请求信息", "Request details")}</summary>
                              {call.http_status != null && <p>HTTP {call.http_status}</p>}
                              {call.request_id && <p className="max-w-xs break-all">{l("请求编号：", "Request ID: ")}{call.request_id}</p>}
                            </details>
                          )}
                        </td>
                        <td className="p-3 text-muted-foreground">
                          {number(call.usage?.prompt_tokens)} / {number(call.usage?.completion_tokens)} / {number(call.usage?.total_tokens)}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </details>
          )}
        </div>
      )}
      {match.status === "completed" && (
        <p className="rounded-xl bg-muted p-4 text-sm">
          {match.winner_address
            ? l(
                `对战胜者 ${shortAddress(match.winner_address)} · ${settlementNote}`,
                `Match winner ${shortAddress(match.winner_address)} · ${settlementNote}`,
              )
            : l(
                "本场平局，没有冠军。",
                "This match ended in a draw with no winner.",
              )}
        </p>
      )}
      {match.error && (
        <p role="alert" className="text-sm text-red-700">
          {l("执行错误：", "Execution error: ")}
          {summary?.last_error_code
            ? llmErrorMessage(summary.last_error_code, lang)
            : match.error}
        </p>
      )}
      <div className="grid gap-4 sm:grid-cols-2">
        <AgentCard
          agent={game.agent_red}
          state={game.agent_red_state}
          side="red"
          showStats={false}
        />
        <AgentCard
          agent={game.agent_blue}
          state={game.agent_blue_state}
          side="blue"
          showStats={false}
        />
      </div>
      <div className="overflow-x-auto">
        <div className="flex min-w-max justify-center rounded-xl bg-slate-50 p-3">
          <ArenaBoard gameState={game} />
        </div>
      </div>
      {match.rounds.length > 0 ? (
        <div className="space-y-4 rounded-xl border p-4">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <p className="text-sm font-semibold">
              {l(`第 ${round?.round} 回合`, `Round ${round?.round}`)}{" "}
              <span className="font-normal text-muted-foreground">
                {l(
                  `/ 已记录 ${match.rounds.length} 回合`,
                  `/ ${match.rounds.length} rounds recorded`,
                )}
              </span>
            </p>
            <div className="flex gap-2">
              <Button
                size="sm"
                variant="outline"
                disabled={index === 0}
                aria-label={l("上一回合", "Previous round")}
                onClick={() => {
                  setPlaying(false);
                  setFrame(Math.max(0, index - 1));
                }}
              >
                <SkipBack size={16} />
              </Button>
              <Button
                size="sm"
                onClick={() => {
                  if (playing) setPlaying(false);
                  else {
                    setFrame(index === match.rounds.length - 1 ? 0 : index);
                    setPlaying(true);
                  }
                }}
              >
                {playing ? (
                  <Pause size={15} className="mr-1" />
                ) : (
                  <Play size={15} className="mr-1" />
                )}
                {playing ? l("暂停", "Pause") : l("播放回放", "Play replay")}
              </Button>
              <Button
                size="sm"
                variant="outline"
                disabled={index === match.rounds.length - 1}
                aria-label={l("下一回合", "Next round")}
                onClick={() => {
                  setPlaying(false);
                  setFrame(Math.min(match.rounds.length - 1, index + 1));
                }}
              >
                <SkipForward size={16} />
              </Button>
            </div>
          </div>
          <input
            aria-label={l("选择回放回合", "Choose replay round")}
            className="block w-full accent-emerald-600"
            type="range"
            min="0"
            max={match.rounds.length - 1}
            value={index}
            onChange={(e) => {
              setPlaying(false);
              setFrame(Number(e.target.value));
            }}
          />
          <button
            className="text-xs font-medium text-primary underline"
            onClick={() => {
              setPlaying(false);
              setFrame(null);
            }}
          >
            {l("跟随最新回合", "Follow latest round")}
          </button>
        </div>
      ) : (
        <p className="text-sm text-muted-foreground">
          {l(
            "暂无回合记录，执行对战后自动更新。",
            "No rounds recorded yet. This updates automatically when the match runs.",
          )}
        </p>
      )}
      <details className="border-t pt-4">
        <summary className="cursor-pointer text-sm font-medium">
          {l(
            "查看回合动作与决策记录",
            "View round actions and decision records",
          )}
        </summary>
        <div className="mt-4">
          <TurnLog turns={game.history} />
        </div>
      </details>
    </section>
  );
}
