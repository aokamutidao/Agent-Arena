"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { api } from "@/lib/api";
import { useI18n } from "@/lib/i18n";
import type { GameListItem } from "@/types/game";

export default function HomePage() {
  const router = useRouter();
  const { t } = useI18n();
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
      console.error("Create game failed:", err);
      alert(t("common.error"));
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
        <p className="text-muted-foreground">{t("common.loading")}</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">{t("nav.home")}</h1>
        <Button onClick={handleCreateGame} disabled={creating}>
          {creating ? t("common.loading") : "🎮 " + t("home.startGame")}
        </Button>
      </div>

      {/* LIVE NOW */}
      <section>
        <h2 className="text-lg font-semibold mb-3 flex items-center gap-2">
          <span className="inline-block w-2 h-2 rounded-full bg-red-500 animate-pulse" />
          {t("game.inProgress")} ({playing.length})
        </h2>
        {playing.length === 0 ? (
          <Card>
            <CardContent className="pt-6">
              <p className="text-sm text-muted-foreground">{t("home.noGames")}</p>
            </CardContent>
          </Card>
        ) : (
          <div className="grid gap-4">
            {playing.map((game) => (
              <GameCard key={game.game_id} game={game} formatUSDC={formatUSDC} live t={t} />
            ))}
          </div>
        )}
      </section>

      {/* UPCOMING */}
      <section>
        <h2 className="text-lg font-semibold mb-3">{t("game.waiting")} ({betting.length})</h2>
        {betting.length === 0 ? (
          <Card>
            <CardContent className="pt-6">
              <p className="text-sm text-muted-foreground">{t("home.noGames")}</p>
            </CardContent>
          </Card>
        ) : (
          <div className="grid gap-4">
            {betting.map((game) => (
              <GameCard key={game.game_id} game={game} formatUSDC={formatUSDC} t={t} />
            ))}
          </div>
        )}
      </section>

      {/* FINISHED */}
      <section>
        <h2 className="text-lg font-semibold mb-3">{t("game.finished")} ({finished.length})</h2>
        {finished.length === 0 ? (
          <Card>
            <CardContent className="pt-6">
              <p className="text-sm text-muted-foreground">{t("home.noGames")}</p>
            </CardContent>
          </Card>
        ) : (
          <div className="grid gap-4">
            {finished.map((game) => (
              <GameCard key={game.game_id} game={game} formatUSDC={formatUSDC} t={t} />
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
  t,
}: {
  game: GameListItem;
  formatUSDC: (wei: string) => string;
  live?: boolean;
  t: (key: string) => string;
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
                {t("game.round")} {game.current_round} / {game.max_rounds}
              </p>
            ) : game.status === "finished" ? (
              <p className="text-muted-foreground">
                {t("game.winner")}: {game.agent_red.name === game.winner ? "🔴" : "🔵"}{" "}
                {game.winner || "none"}
              </p>
            ) : (
              <p className="text-muted-foreground">{t("game.waiting")}</p>
            )}
          </div>
          <div className="flex items-center gap-4">
            <span className="text-muted-foreground">Pool: {totalPool} USDC</span>
            <Link
              href={`/game/${game.game_id}`}
              className="text-primary hover:underline font-medium"
            >
              {game.status === "playing" ? t("game.title") : t("common.confirm")} →
            </Link>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
