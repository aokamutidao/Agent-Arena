"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { api } from "@/lib/api";
import type { GameListItem } from "@/types/game";

export default function HomePage() {
  const router = useRouter();
  const [games, setGames] = useState<GameListItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);

  const refreshGames = () => {
    api
      .listGames()
      .then((data) => {
        setGames(data.games || []);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  };

  useEffect(() => {
    refreshGames();
  }, []);

  const handleCreateGame = async () => {
    setCreating(true);
    try {
      const res = await api.createGame();
      router.push(`/game/${res.game_id}`);
    } catch (err) {
      console.error("创建对局失败:", err);
      alert("创建对局失败，请检查后端是否运行");
    } finally {
      setCreating(false);
    }
  };

  const formatUSDC = (wei: string) => {
    const num = parseInt(wei) / 1e6;
    return num.toFixed(2);
  };

  const playing = games.filter((g) => g.status === "playing");
  const betting = games.filter((g) => g.status === "betting" || g.status === "pending");
  const finished = games.filter((g) => g.status === "finished");

  if (loading) {
    return (
      <div className="flex items-center justify-center py-20">
        <p className="text-muted-foreground">加载中...</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">对局列表</h1>
        <Button onClick={handleCreateGame} disabled={creating}>
          {creating ? "创建中..." : "🎮 创建新对局"}
        </Button>
      </div>

      {/* LIVE NOW */}
      <section>
        <h2 className="text-lg font-semibold mb-3 flex items-center gap-2">
          <span className="inline-block w-2 h-2 rounded-full bg-red-500 animate-pulse" />
          正在进行 ({playing.length})
        </h2>
        {playing.length === 0 ? (
          <Card>
            <CardContent className="pt-6">
              <p className="text-sm text-muted-foreground">暂无进行中的对局</p>
            </CardContent>
          </Card>
        ) : (
          <div className="grid gap-4">
            {playing.map((game) => (
              <GameCard key={game.game_id} game={game} formatUSDC={formatUSDC} live />
            ))}
          </div>
        )}
      </section>

      {/* UPCOMING */}
      <section>
        <h2 className="text-lg font-semibold mb-3">即将开始 ({betting.length})</h2>
        {betting.length === 0 ? (
          <Card>
            <CardContent className="pt-6">
              <p className="text-sm text-muted-foreground">暂无即将开始的对局</p>
            </CardContent>
          </Card>
        ) : (
          <div className="grid gap-4">
            {betting.map((game) => (
              <GameCard key={game.game_id} game={game} formatUSDC={formatUSDC} />
            ))}
          </div>
        )}
      </section>

      {/* FINISHED */}
      <section>
        <h2 className="text-lg font-semibold mb-3">已结束 ({finished.length})</h2>
        {finished.length === 0 ? (
          <Card>
            <CardContent className="pt-6">
              <p className="text-sm text-muted-foreground">暂无已结束的对局</p>
            </CardContent>
          </Card>
        ) : (
          <div className="grid gap-4">
            {finished.map((game) => (
              <GameCard key={game.game_id} game={game} formatUSDC={formatUSDC} />
            ))}
          </div>
        )}
      </section>
    </div>
  );
}

function GameCard({
  game,
  formatUSDC,
  live,
}: {
  game: GameListItem;
  formatUSDC: (wei: string) => string;
  live?: boolean;
}) {
  const redPool = formatUSDC(game.total_bet_red);
  const bluePool = formatUSDC(game.total_bet_blue);
  const totalPool = (parseFloat(redPool) + parseFloat(bluePool)).toFixed(2);

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base flex items-center justify-between">
          <span>
            🔴 {game.agent_red.name} vs {game.agent_blue.name} 🔵
          </span>
          {live && (
            <span className="inline-block w-2 h-2 rounded-full bg-red-500 animate-pulse" />
          )}
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex justify-between items-center text-sm">
          <div>
            {game.status === "playing" ? (
              <p className="text-muted-foreground">
                第 {game.current_round} / {game.max_rounds} 回合
              </p>
            ) : game.status === "finished" ? (
              <p className="text-muted-foreground">
                胜者: {game.agent_red.name === game.winner ? "🔴" : "🔵"}{" "}
                {game.winner || "none"}
              </p>
            ) : (
              <p className="text-muted-foreground">等待开始</p>
            )}
          </div>
          <div className="flex items-center gap-4">
            <span className="text-muted-foreground">池: {totalPool} USDC</span>
            <Link
              href={`/game/${game.game_id}`}
              className="text-primary hover:underline font-medium"
            >
              {game.status === "playing" ? "观战" : "查看"} →
            </Link>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
