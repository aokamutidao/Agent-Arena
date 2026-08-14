"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/lib/auth";

interface GameHistoryItem {
  game_id: number;
  red_name: string;
  blue_name: string;
  winner: "red" | "blue" | "draw";
  winner_name: string;
  final_hp_red: number;
  final_hp_blue: number;
  total_rounds: number;
  finish_tx_hash: string;
  created_at: string;
}

const ETHERSCAN_BASE = "https://sepolia.etherscan.io/tx";

export default function HistoryPage() {
  const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
  const { token } = useAuth();
  const [history, setHistory] = useState<GameHistoryItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    const load = async () => {
      try {
        const res = await fetch(`${API_URL}/api/game-history?limit=50`);
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        const data = await res.json();
        if (!cancelled) setHistory(data.history || []);
      } catch (err: any) {
        if (!cancelled) setError(err.message || `加载失败");
      } finally {
        if (!cancelled) setLoading(false);
      }
    };
    load();
    // 自动轮询：若有任意记录尚缺链上 tx，每 5 秒刷新一次；全部同步后降级为 30 秒
    const interval = setInterval(load, 5000);
    return () => {
      cancelled = true;
      clearInterval(interval);
    };
  }, []);

  // 是否所有记录都已有链上 tx（用于判断是否需要频繁刷新）
  const allSynced = history.length > 0 && history.every((h) => !!h.finish_tx_hash);

  const winnerBadge = (item: GameHistoryItem) => {
    if (item.winner === "draw") {
      return <span className="px-2 py-0.5 text-xs rounded bg-gray-500/20 text-gray-400">平局</span>;
    }
    const isRed = item.winner === "red";
    return (
      <span
        className={`px-2 py-0.5 text-xs rounded ${
          isRed ? "bg-red-500/20 text-red-400" : "bg-blue-500/20 text-blue-400"
        }`}
      >
        {isRed ? "🔴 红胜" : "🔵 蓝胜"}
      </span>
    );
  };

  return (
    <div className="container mx-auto px-4 py-8 space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">📜 对局记录</h1>
        <Button
          variant="outline"
          size="sm"
          onClick={() => window.location.reload()}
          disabled={loading}
        >
          刷新
        </Button>
      </div>

      {loading && (
        <div className="text-center py-12 text-muted-foreground">加载中...</div>
      )}
      {error && (
        <Card className="border-red-500/30 bg-red-500/5">
          <CardContent className="pt-6">
            <p className="text-red-400">❌ {error}</p>
          </CardContent>
        </Card>
      )}

      {!loading && !error && history.length === 0 && (
        <Card>
          <CardContent className="pt-6 text-center text-muted-foreground">
            暂无对局记录。完成一局对战后将在此显示。
          </CardContent>
        </Card>
      )}

      {history.length > 0 && (
        <div className="grid gap-3">
          {history.map((item) => (
            <Card key={item.game_id}>
              <CardContent className="pt-5 pb-5">
                <div className="flex items-center justify-between flex-wrap gap-3">
                  <div className="flex items-center gap-3">
                    <span className="text-xs text-muted-foreground font-mono">
                      #{item.game_id}
                    </span>
                    <div className="flex items-center gap-2 text-sm">
                      <span className={item.winner === "red" ? "font-bold text-red-400" : "text-muted-foreground"}>
                        🔴 {item.red_name}
                      </span>
                      <span className="text-muted-foreground">vs</span>
                      <span className={item.winner === "blue" ? "font-bold text-blue-400" : "text-muted-foreground"}>
                        🔵 {item.blue_name}
                      </span>
                    </div>
                    {winnerBadge(item)}
                  </div>

                  <div className="flex items-center gap-4 text-sm">
                    <span className="text-muted-foreground">
                      HP: {item.final_hp_red} / {item.final_hp_blue}
                    </span>
                    <span className="text-muted-foreground">
                      {item.total_rounds} 回合
                    </span>
                    <span className="text-xs text-muted-foreground">
                      {new Date(item.created_at).toLocaleString()}
                    </span>
                  </div>
                </div>

                <div className="mt-3 flex items-center gap-3 text-xs flex-wrap">
                  <Link
                    href={`/game/${item.game_id}`}
                    className="text-blue-400 hover:underline"
                  >
                    查看对局详情 →
                  </Link>
                  {item.finish_tx_hash ? (
                    <a
                      href={`${ETHERSCAN_BASE}/${item.finish_tx_hash}`}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="inline-flex items-center gap-1 text-green-400 hover:underline"
                      title="在 Sepolia Etherscan 查看链上交易"
                    >
                      ⛓️ 链上记录
                      <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                        <path d="M7 17L17 7M17 7H8M17 7V16" />
                      </svg>
                    </a>
                  ) : new Date(item.created_at) < new Date("2026-08-12T03:14:00") ? (
                    <span className="text-gray-500" title="合约重部署前的历史记录，未上链">
                      📝 未上链（历史遗留）
                    </span>
                  ) : (
                    <span className="text-yellow-500" title="链上交易正在打包，稍后刷新">
                      ⏳ 链上同步中
                    </span>
                  )}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
