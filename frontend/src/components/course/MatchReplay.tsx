"use client";
import { useState } from "react";
import { ArenaBoard } from "@/components/game/ArenaBoard";
import { AgentCard } from "@/components/game/AgentCard";
import { TurnLog } from "@/components/game/TurnLog";
import { type MatchResult, type EngineAgentState, type TournamentStatus, shortAddress } from "@/lib/course";
import type { GameDetail, AgentInfo, AgentState } from "@/types/game";
export function MatchReplay({ match, settlementStatus }: {
    match: MatchResult;
    settlementStatus: TournamentStatus;
}) {
    const [frame, setFrame] = useState<number | null>(null);
    const settlementNote = settlementStatus === "finalized" ? "本场已链上结算，审定冠军见资金区" : settlementStatus === "cancelled" ? "赛事已取消，不发放冠军奖金" : "等待链上审核";
    const index = frame === null ? Math.max(0, match.rounds.length - 1) : Math.min(frame, Math.max(0, match.rounds.length - 1));
    const round = match.rounds[index];
    const agent = (address: string, name: string): AgentInfo => ({ id: address, name, personality: shortAddress(address), wins: 0, losses: 0, win_rate: 0, description: "" });
    const replayState = (state: EngineAgentState | undefined, hp: number | undefined, position: {
        x: number;
        y: number;
    }): AgentState => ({ hp: state?.hp ?? hp ?? 100, max_hp: state?.max_hp ?? 100, position: state?.position || position, status: state?.status?.map(effect => effect.type) || [], skill_cooldown: state?.skill_cooldown || 0, is_charging: state?.is_charging || false });
    const game: GameDetail = { game_id: match.tournament_id, agent_red: agent(match.red_address, match.red_agent?.name || "红方 Agent"), agent_blue: agent(match.blue_address, match.blue_agent?.name || "蓝方 Agent"), status: match.status === "completed" && index === match.rounds.length - 1 ? "finished" : "playing", current_round: round?.round || 0, max_rounds: 30, total_bet_red: "0", total_bet_blue: "0", agent_red_state: replayState(round?.red_state, round?.red_hp, round?.red_position || { x: 1, y: 1 }), agent_blue_state: replayState(round?.blue_state, round?.blue_hp, round?.blue_position || { x: 8, y: 8 }), history: match.rounds.slice(0, index + 1).map(r => ({ round: r.round, red_action: r.red_action, blue_action: r.blue_action, red_hp_after: r.red_hp, blue_hp_after: r.blue_hp, red_reasoning: r.red_reasoning, blue_reasoning: r.blue_reasoning })) };
    return <section className="border rounded-xl p-5 space-y-4"><div className="flex flex-wrap justify-between gap-3"><h2 className="text-lg font-semibold">对战与回合回放</h2><p className="text-sm">{match.mode === "llm" ? "LLM 模型决策" : "确定性模拟（非 LLM）"} · {match.status}</p></div>{match.status === "completed" && <p className="text-sm">对战结果：{match.winner_address ? `拟冠军 ${shortAddress(match.winner_address)}（${settlementNote}）` : "平局，没有冠军"}</p>}{match.error && <p role="alert" className="text-red-500">{match.error}</p>}<div className="grid md:grid-cols-2 gap-4"><AgentCard agent={game.agent_red} state={game.agent_red_state} side="red" showStats={false}/><AgentCard agent={game.agent_blue} state={game.agent_blue_state} side="blue" showStats={false}/></div><div className="overflow-x-auto flex justify-start md:justify-center"><ArenaBoard gameState={game}/></div>{match.rounds.length > 0 && <label className="block text-sm">回放回合：{round?.round}<input aria-label="回放回合" className="block w-full mt-2" type="range" min="0" max={match.rounds.length - 1} value={index} onChange={e => setFrame(Number(e.target.value))}/><button className="underline text-xs mt-2" onClick={() => setFrame(null)}>跟随最新回合</button></label>}<TurnLog turns={game.history}/></section>;
}
