"use client";

import type { AgentInfo, AgentState } from "@/types/game";
import { cn } from "@/lib/utils";

interface AgentCardProps {
  agent: AgentInfo;
  state?: AgentState;
  side: "red" | "blue";
}

export function AgentCard({ agent, state, side }: AgentCardProps) {
  const hp = state?.hp ?? 100;
  const maxHp = state?.max_hp ?? 100;
  const hpPercent = Math.round((hp / maxHp) * 100);

  const sideColor = side === "red" ? "text-red-500" : "text-blue-500";
  const sideBg = side === "red" ? "bg-red-500" : "bg-blue-500";
  const sideBorder = side === "red" ? "border-red-500/30" : "border-blue-500/30";

  return (
    <div className={cn("rounded-lg border p-4", sideBorder)}>
      <div className="flex items-center gap-3 mb-3">
        <div className={cn("w-10 h-10 rounded-full flex items-center justify-center text-white font-bold", sideBg)}>
          {agent.name[0]}
        </div>
        <div>
          <h3 className={cn("font-semibold", sideColor)}>{agent.name}</h3>
          <p className="text-xs text-muted-foreground">{agent.personality}</p>
        </div>
      </div>

      {/* HP Bar */}
      <div className="space-y-1">
        <div className="flex justify-between text-sm">
          <span className="text-muted-foreground">HP</span>
          <span className="font-mono">
            {hp}/{maxHp}
          </span>
        </div>
        <div className="h-3 bg-secondary rounded-full overflow-hidden">
          <div
            className={cn("h-full rounded-full transition-all duration-500", sideBg)}
            style={{ width: `${hpPercent}%` }}
          />
        </div>
      </div>

      {/* Position */}
      {state?.position && (
        <p className="text-xs text-muted-foreground mt-2">
          位置: ({state.position.x}, {state.position.y})
        </p>
      )}

      {/* Status Effects */}
      {state?.status && state.status.length > 0 && (
        <div className="flex gap-1 mt-2">
          {state.status.map((eff) => (
            <span
              key={eff}
              className="text-xs px-2 py-0.5 rounded bg-secondary text-secondary-foreground"
            >
              {eff}
            </span>
          ))}
        </div>
      )}

      {/* Stats */}
      <div className="text-xs text-muted-foreground mt-2">
        胜率: {agent.win_rate.toFixed(1)}% ({agent.wins}W {agent.losses}L)
      </div>
    </div>
  );
}
