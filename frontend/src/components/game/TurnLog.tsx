"use client";

import type { TurnRecord } from "@/types/game";

interface TurnLogProps {
  turns: TurnRecord[];
}

function formatAction(action: {
  type: string;
  target?: { x: number; y: number };
  failed?: boolean;
  fail_reason?: string;
} | string): { text: string; failed: boolean; failReason: string } {
  if (typeof action === "string") return { text: action, failed: false, failReason: "" };
  let text: string;
  if (action.target) {
    text = `${action.type}(${action.target.x},${action.target.y})`;
  } else {
    text = action.type;
  }
  return { text, failed: !!action.failed, failReason: action.fail_reason || "" };
}

function actionIcon(type: string): string {
  switch (type) {
    case "ATTACK": return "⚔️";
    case "SKILL":  return "🎯";
    case "MOVE":   return "🏃";
    case "CHARGE": return "💪";
    case "WAIT":   return "⏳";
    case "HEAL":   return "💊";
    default:       return "❓";
  }
}

export function TurnLog({ turns }: TurnLogProps) {
  if (!turns || turns.length === 0) {
    return (
      <div className="text-sm text-muted-foreground text-center py-4">
        暂无回合记录
      </div>
    );
  }

  return (
    <div className="space-y-1.5 max-h-96 overflow-y-auto">
      {turns.map((turn) => {
        const redType = typeof turn.red_action === "string" ? turn.red_action : turn.red_action?.type ?? "WAIT";
        const blueType = typeof turn.blue_action === "string" ? turn.blue_action : turn.blue_action?.type ?? "WAIT";
        const redFmt = formatAction(turn.red_action);
        const blueFmt = formatAction(turn.blue_action);
        return (
          <div key={turn.round} className="px-2 py-1.5 rounded hover:bg-accent/50 space-y-0.5">
            <div className="flex items-center gap-2 text-sm flex-wrap">
              <span className="text-muted-foreground font-mono w-10 shrink-0">
                R{turn.round}
              </span>
              {/* Red action */}
              <span className={redFmt.failed ? "text-red-500/50 line-through" : "text-red-500 font-medium"}>
                {actionIcon(redType)} {redFmt.text}
                {redFmt.failed && <span className="no-underline ml-1 text-red-400 text-xs font-semibold not-italic">✗</span>}
              </span>
              <span className="text-muted-foreground">vs</span>
              {/* Blue action */}
              <span className={blueFmt.failed ? "text-blue-500/50 line-through" : "text-blue-500 font-medium"}>
                {actionIcon(blueType)} {blueFmt.text}
                {blueFmt.failed && <span className="no-underline ml-1 text-blue-400 text-xs font-semibold not-italic">✗</span>}
              </span>
              <span className="text-muted-foreground ml-auto font-mono text-xs">
                {turn.red_hp_after}/{turn.blue_hp_after} HP
              </span>
            </div>
            {/* 失败原因 */}
            {(redFmt.failed || blueFmt.failed) && (
              <div className="flex gap-2 pl-10 text-xs">
                {redFmt.failed && redFmt.failReason && (
                  <span className="text-red-400/60">🔴 {redFmt.failReason}</span>
                )}
                {blueFmt.failed && blueFmt.failReason && (
                  <span className="text-blue-400/60">🔵 {blueFmt.failReason}</span>
                )}
              </div>
            )}
            {/* Agent 思考/意图 */}
            {(turn.red_reasoning || turn.blue_reasoning) && (
              <div className="flex gap-2 pl-10 text-xs text-muted-foreground">
                {turn.red_reasoning && (
                  <span className="text-red-400/70">🔴 {turn.red_reasoning}</span>
                )}
                {turn.blue_reasoning && (
                  <span className="text-blue-400/70">🔵 {turn.blue_reasoning}</span>
                )}
              </div>
            )}
          </div>
        );
      })}
    </div>
  );
}
