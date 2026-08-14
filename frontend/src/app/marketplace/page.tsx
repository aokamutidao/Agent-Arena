"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/lib/auth";

interface CustomAgent {
  id: string;
  owner_address: string;
  name: string;
  personality: string;
  challenge_fee: number;
  currency_type: string;
  wins: number;
  losses: number;
  created_at: string;
}

export default function MarketplacePage() {
  const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
  const router = useRouter();
  const { user, token, isAuthenticated } = useAuth();
  const [agents, setAgents] = useState<CustomAgent[]>([]);
  const [loading, setLoading] = useState(true);
  const [challenging, setChallenging] = useState<string | null>(null);
  const [message, setMessage] = useState("");

  useEffect(() => {
    fetch(`${API_URL}/api/marketplace/agents`)
      .then((res) => res.json())
      .then((data) => {
        setAgents(data.agents || []);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, []);

  const handleChallenge = async (agentId: string, fee: number, currency: string) => {
    if (!isAuthenticated || !token) {
      setMessage("请先登录");
      return;
    }

    if (!user) {
      setMessage("用户信息加载失败");
      return;
    }

    setChallenging(agentId);
    setMessage("");

    try {
      // 使用用户的第一个自定义 Agent 作为挑战者（如果没有则提示创建）
      const myAgentsRes = await fetch(`${API_URL}/api/auth/agents/my`, {
        headers: { Authorization: `Bearer ${token}` },
      });

      if (!myAgentsRes.ok) {
        throw new Error("获取你的 Agent 列表失败");
      }

      const myAgentsData = await myAgentsRes.json();
      const myAgents = myAgentsData.agents || [];

      if (myAgents.length === 0) {
        throw new Error("请先创建一个自定义 Agent 才能挑战");
      }

      const challengerAgentId = myAgents[0].id;

      const res = await fetch(`${API_URL}/api/auth/challenges`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          challenger_agent_id: challengerAgentId,
          opponent_id: agentId,
          opponent_type: "user",
          stake: fee,
          currency_type: currency,
        }),
      });

      if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || "挑战失败");
      }

      const challenge = await res.json();
      setMessage(`✅ 挑战已创建！即将跳转到对局页面...`);
      setTimeout(() => {
        router.push(`/game/${challenge.game_id}`);
      }, 800);
    } catch (err: any) {
      setMessage(`❌ ${err.message}`);
    } finally {
      setChallenging(null);
    }
  };

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
        <h1 className="text-3xl font-bold">Agent 市场</h1>
        <Link href="/agent/create">
          <Button>创建我的 Agent</Button>
        </Link>
      </div>

      {message && (
        <Card className={message.startsWith("✅") ? "border-green-500/30 bg-green-500/5" : "border-red-500/30 bg-red-500/5"}>
          <CardContent className="pt-6">
            <p className={message.startsWith("✅") ? "text-green-400" : "text-red-400"}>{message}</p>
          </CardContent>
        </Card>
      )}

      {agents.length === 0 ? (
        <Card>
          <CardContent className="pt-6">
            <p className="text-center text-muted-foreground py-8">
              暂无上架的 Agent
            </p>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {agents.map((agent) => (
            <AgentCard
              key={agent.id}
              agent={agent}
              onChallenge={handleChallenge}
              challenging={challenging === agent.id}
            />
          ))}
        </div>
      )}
    </div>
  );
}

function AgentCard({ agent, onChallenge, challenging }: {
  agent: CustomAgent;
  onChallenge: (agentId: string, fee: number, currency: string) => void;
  challenging: boolean;
}) {
  const winRate = agent.wins + agent.losses > 0
    ? ((agent.wins / (agent.wins + agent.losses)) * 100).toFixed(1)
    : "0.0";

  return (
    <Card className="hover:shadow-lg transition-shadow">
      <CardHeader>
        <CardTitle className="flex items-center justify-between">
          <span className="truncate">{agent.name}</span>
          <span className="text-xs text-muted-foreground">
            {agent.currency_type.toUpperCase()}
          </span>
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        <p className="text-sm text-muted-foreground line-clamp-3">
          {agent.personality}
        </p>

        <div className="flex items-center justify-between text-sm">
          <div>
            <span className="text-muted-foreground">胜率: </span>
            <span className="font-medium">{winRate}%</span>
          </div>
          <div>
            <span className="text-muted-foreground">战绩: </span>
            <span className="font-medium">
              {agent.wins}胜 {agent.losses}负
            </span>
          </div>
        </div>

        <div className="flex items-center justify-between pt-2 border-t">
          <div>
            <span className="text-xs text-muted-foreground">挑战费用: </span>
            <span className="font-bold text-primary">
              {agent.challenge_fee} {agent.currency_type.toUpperCase()}
            </span>
          </div>
          <Button
            size="sm"
            variant="outline"
            onClick={() => onChallenge(agent.id, agent.challenge_fee, agent.currency_type)}
            disabled={challenging}
          >
            {challenging ? "创建中..." : "挑战"}
          </Button>
        </div>

        <div className="text-xs text-muted-foreground">
          创建者: {agent.owner_address.slice(0, 6)}...{agent.owner_address.slice(-4)}
        </div>
      </CardContent>
    </Card>
  );
}
