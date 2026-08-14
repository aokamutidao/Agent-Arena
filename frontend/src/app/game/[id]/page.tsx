"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useGame } from "@/hooks/useGame";
import { useAuth } from "@/lib/auth";
import { ArenaBoard } from "@/components/game/ArenaBoard";
import { AgentCard } from "@/components/game/AgentCard";
import { TurnLog } from "@/components/game/TurnLog";
import { BetPanel } from "@/components/betting/BetPanel";
import { StrategyVote } from "@/components/betting/StrategyVote";
import { api } from "@/lib/api";

export default function GamePage() {
  const params = useParams();
  const gameId = Number(params.id);

  const { game, loading, error, wsConnected, gameFinished, balanceUpdate } = useGame(gameId);
  const { refreshUser, user } = useAuth();
  const [starting, setStarting] = useState(false);
  const [creating, setCreating] = useState(false);
  const router = useRouter();

  // 从后端读取分边策略权重（链上 getStrategyWeights 是全局的，不分边）
  // 后端的 game.strategy_red / strategy_blue 是分边的
  const backendStrategyRed = game?.strategy_red
    ? {
        aggressive: game.strategy_red.aggressive ?? 33,
        defensive: game.strategy_red.defensive ?? 33,
        tricky: game.strategy_red.tricky ?? 34,
      }
    : undefined;
  const backendStrategyBlue = game?.strategy_blue
    ? {
        aggressive: game.strategy_blue.aggressive ?? 33,
        defensive: game.strategy_blue.defensive ?? 33,
        tricky: game.strategy_blue.tricky ?? 34,
      }
    : undefined;

  // 当后端广播余额更新时，刷新当前登录用户的资料
  useEffect(() => {
    if (!balanceUpdate || !user) return;
    if (balanceUpdate.address.toLowerCase() === user.address.toLowerCase()) {
      refreshUser();
    }
  }, [balanceUpdate, user, refreshUser]);

  const handleStartGame = async () => {
    setStarting(true);
    try {
      await api.startGame(gameId);
    } catch {
      // 启动失败不影响页面
    } finally {
      setStarting(false);
    }
  };

  const handleCreateGame = async () => {
    setCreating(true);
    try {
      const res = await api.createGame();
      router.push(`/game/${res.game_id}`);
    } catch {
      // 创建失败
    } finally {
      setCreating(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center py-20">
        <p className="text-muted-foreground">加载中...</p>
      </div>
    );
  }

  if (error || !game) {
    return (
      <div className="flex items-center justify-center py-20">
        <p className="text-destructive">{error || "对局不存在"}</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">
          对局 #{game.game_id}
          {game.status === "finished" && (
            <span className="ml-3 text-sm font-normal text-muted-foreground">
              已结束
            </span>
          )}
          {game.status === "betting" && (
            <span className="ml-3 text-sm font-normal text-yellow-500">
              等待开始
            </span>
          )}
        </h1>
        <div className="flex items-center gap-3">
          {game.status === "betting" && (
            <Button
              onClick={handleStartGame}
              disabled={starting}
              size="sm"
              className="bg-green-600 hover:bg-green-700"
            >
              {starting ? "启动中..." : "⚔️ 开始对局"}
            </Button>
          )}
          <div className="flex items-center gap-2 text-sm">
            <span
              className={`inline-block w-2 h-2 rounded-full ${
                wsConnected ? "bg-green-500" : "bg-gray-400"
              }`}
            />
            <span className="text-muted-foreground">
              {wsConnected ? "实时连接" : "未连接"}
            </span>
          </div>
        </div>
      </div>

      {/* Game Finished Banner — 来自 WS 实时事件或 REST 初始数据 */}
      {(gameFinished || game.status === "finished") && (
        <Card className="border-green-500/30 bg-green-500/5">
          <CardContent className="pt-6 text-center space-y-3">
            <p className="text-lg font-semibold">
              🏆 {game.winner === "red" ? "红方" : game.winner === "blue" ? "蓝方" : ""}获胜！
            </p>
            <p className="text-muted-foreground">
              胜者: {gameFinished?.winner_name || game.winner || "—"} |
              共 {gameFinished?.total_rounds || game.current_round} 回合
            </p>
            <p className="text-sm text-muted-foreground mt-1">
              最终 HP — 红方: {game.agent_red_state?.hp ?? 0} / 蓝方: {game.agent_blue_state?.hp ?? 0}
            </p>
            <Button
              onClick={handleCreateGame}
              disabled={creating}
              className="mt-2"
            >
              {creating ? "创建中..." : "🆕 新建对局"}
            </Button>
          </CardContent>
        </Card>
      )}

      {/* Overtime Banner */}
      {game?.overtime && !gameFinished && (
        <Card className="border-orange-500/30 bg-orange-500/5">
          <CardContent className="pt-6 text-center">
            <p className="text-lg font-semibold text-orange-500">⚡ 加时赛！</p>
            <p className="text-sm text-muted-foreground">
              双方 HP 相同，进入 Sudden Death — 每回合双方各受 10 点伤害
            </p>
          </CardContent>
        </Card>
      )}

      {/* Archived Banner — 重启后从 DB 恢复的已完成对局，无回合明细 */}
      {game?.archived && (
        <Card className="border-yellow-500/30 bg-yellow-500/5">
          <CardContent className="pt-6 text-center">
            <p className="text-base font-semibold text-yellow-500">
              🗄️ 已归档（后端重启）
            </p>
            <p className="text-sm text-muted-foreground mt-1">
              最终结果来自数据库，回合明细在重启前未持久化。
              完整链上记录请前往 <a href="/history" className="text-blue-400 hover:underline">对局记录</a> 查看。
            </p>
          </CardContent>
        </Card>
      )}

      {/* Agent Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <AgentCard
          agent={game.agent_red}
          state={game.agent_red_state}
          side="red"
        />
        <AgentCard
          agent={game.agent_blue}
          state={game.agent_blue_state}
          side="blue"
        />
      </div>

      {/* Arena Board */}
      <Card>
        <CardContent className="pt-6">
          <div className="flex justify-center overflow-auto">
            <ArenaBoard gameState={game} />
          </div>
          <div className="text-center mt-4 text-sm text-muted-foreground">
            第 {game.current_round} / {game.max_rounds} 回合
          </div>
        </CardContent>
      </Card>

      {/* Bet Panel */}
      <BetPanel
        gameId={game.game_id}
        totalBetRed={game.total_bet_red}
        totalBetBlue={game.total_bet_blue}
        status={game.status}
      />

      {/* Strategy Vote — 从后端读取分边策略权重 */}
      <StrategyVote
        gameId={game.game_id}
        strategyRed={backendStrategyRed || game.strategy_red}
        strategyBlue={backendStrategyBlue || game.strategy_blue}
      />

      {/* Turn Log */}
      <Card>
        <CardContent className="pt-6">
          <h3 className="text-lg font-semibold mb-3">📜 回合日志</h3>
          <TurnLog turns={game.history} />
        </CardContent>
      </Card>
    </div>
  );
}
