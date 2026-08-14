"use client";

import type { StrategyWeights } from "@/types/game";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { cn } from "@/lib/utils";

interface StrategyVoteProps {
  gameId: number;
  strategyRed?: StrategyWeights;
  strategyBlue?: StrategyWeights;
}

const STRATEGIES = [
  { key: "aggressive" as const, label: "🔥 激进", color: "bg-red-500" },
  { key: "defensive" as const, label: "🛡️ 稳健", color: "bg-green-500" },
  { key: "tricky" as const, label: "🎲 诡道", color: "bg-purple-500" },
];

function StrategyPanel({
  side,
  weights,
}: {
  side: "red" | "blue";
  weights: StrategyWeights;
}) {
  const sideColor = side === "red" ? "border-red-500/30" : "border-blue-500/30";
  const sideTitle = side === "red" ? "🔴 Red 策略偏好" : "🔵 Blue 策略偏好";

  return (
    <Card className={cn("border", sideColor)}>
      <CardHeader className="pb-2">
        <CardTitle className="text-base">{sideTitle}</CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        {STRATEGIES.map((s) => {
          const weight = weights[s.key];
          return (
            <div key={s.key} className="space-y-1">
              <div className="flex justify-between text-sm">
                <span>{s.label}</span>
                <span className="text-muted-foreground font-mono">{weight}%</span>
              </div>
              <div className="h-2 bg-secondary rounded-full overflow-hidden">
                <div
                  className={cn("h-full rounded-full transition-all duration-500", s.color)}
                  style={{ width: `${weight}%` }}
                />
              </div>
            </div>
          );
        })}

        <p className="text-xs text-muted-foreground text-center pt-2">
          下注时选择策略自动影响本方权重
        </p>
      </CardContent>
    </Card>
  );
}

export function StrategyVote({ strategyRed, strategyBlue }: StrategyVoteProps) {
  // 使用后端分边数据（下注时 strategy 参数影响对应方权重）
  const red = strategyRed || { aggressive: 33, defensive: 33, tricky: 34 };
  const blue = strategyBlue || { aggressive: 33, defensive: 33, tricky: 34 };

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      <StrategyPanel side="red" weights={red} />
      <StrategyPanel side="blue" weights={blue} />
    </div>
  );
}
