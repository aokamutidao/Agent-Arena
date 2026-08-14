"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/lib/auth";

interface SystemAgent {
  id: string;
  name: string;
  personality: string;
  difficulty: string;
}

interface CustomAgent {
  id: string;
  name: string;
  personality: string;
}

const SYSTEM_AGENTS: SystemAgent[] = [
  {
    id: "berserker",
    name: "🔥 Berserker",
    personality: "激进战士，最短路径攻击，喜欢蓄力爆发",
    difficulty: "中等",
  },
  {
    id: "tactician",
    name: "🎯 Tactician",
    personality: "战术大师，控制距离，远程技能，耐心等待",
    difficulty: "困难",
  },
  {
    id: "trickster",
    name: "🎭 Trickster",
    personality: "诡计多端，假撤退真蓄力，利用障碍物",
    difficulty: "困难",
  },
  {
    id: "defender",
    name: "🛡️ Defender",
    personality: "铁壁防守，等待对手过度延伸，反击专家",
    difficulty: "简单",
  },
];

export default function PVEPage() {
  const router = useRouter();
  const { user, token, isAuthenticated } = useAuth();
  const [challenging, setChallenging] = useState<string | null>(null);
  const [message, setMessage] = useState("");
  const [myAgents, setMyAgents] = useState<CustomAgent[]>([]);
  const [selectedAgentId, setSelectedAgentId] = useState<string>("");

  useEffect(() => {
    if (!token) return;

    // 获取用户的 Agent 列表
    const fetchAgents = async () => {
      try {
        const res = await fetch("http://localhost:8080/api/auth/agents/my", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        if (res.ok) {
          const data = await res.json();
          const agents = data.agents || [];
          setMyAgents(agents);
          if (agents.length > 0) {
            setSelectedAgentId(agents[0].id); // 默认选择第一个
          }
        }
      } catch (err) {
        console.error("Failed to fetch agents:", err);
      }
    };

    fetchAgents();
  }, [token]);

  const handleChallenge = async (opponentId: string) => {
    if (!token) {
      setMessage("请先登录");
      return;
    }

    if (!selectedAgentId) {
      setMessage("请先创建一个自定义 Agent");
      return;
    }

    setChallenging(opponentId);
    setMessage("");

    try {
      const res = await fetch("http://localhost:8080/api/auth/challenges", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          challenger_agent_id: selectedAgentId,
          opponent_id: opponentId,
          opponent_type: "system",
          stake: 10,
          currency_type: "ac",
        }),
      });

      if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || "挑战失败");
      }

      const challenge = await res.json();
      setMessage(`✅ 挑战已创建！即将跳转到对局页面...`);
      // 跳转到对局页面（处于 betting 状态，等待挑战者手动开始）
      setTimeout(() => {
        window.location.href = `/game/${challenge.game_id}`;
      }, 800);
    } catch (err: any) {
      setMessage(`❌ ${err.message}`);
    } finally {
      setChallenging(null);
    }
  };

  if (!isAuthenticated || !user) {
    return (
      <Card className="max-w-md mx-auto mt-20">
        <CardHeader>
          <CardTitle>请先登录</CardTitle>
          <CardDescription>登录后才能挑战系统 Agent</CardDescription>
        </CardHeader>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">🎮 PVE 挑战</h1>
        <div className="text-right">
          <div className="text-sm text-muted-foreground">你的 AC 余额</div>
          <div className="text-2xl font-bold text-primary">
            {(user.ac_on_chain_balance ?? user.ac_balance ?? 0).toLocaleString()} AC
          </div>
        </div>
      </div>

      {/* 选择你的 Agent */}
      <Card>
        <CardHeader>
          <CardTitle>选择你的 Agent</CardTitle>
          <CardDescription>选择你要派出的自定义 Agent 进行挑战</CardDescription>
        </CardHeader>
        <CardContent>
          {myAgents.length === 0 ? (
            <div className="space-y-3">
              <p className="text-sm text-muted-foreground">
                你还没有创建自定义 Agent。
              </p>
              <Button onClick={() => router.push("/agent/create")}>
                + 创建第一个 Agent
              </Button>
            </div>
          ) : (
            <div className="space-y-3">
              <div className="grid gap-2">
                {myAgents.map((agent) => (
                  <label
                    key={agent.id}
                    className={`flex items-center gap-3 p-3 border rounded-lg cursor-pointer transition-colors ${
                      selectedAgentId === agent.id
                        ? "border-primary bg-primary/5"
                        : "hover:bg-muted"
                    }`}
                  >
                    <input
                      type="radio"
                      name="agent"
                      value={agent.id}
                      checked={selectedAgentId === agent.id}
                      onChange={(e) => setSelectedAgentId(e.target.value)}
                      className="w-4 h-4"
                    />
                    <div className="flex-1">
                      <div className="font-medium">{agent.name}</div>
                      <div className="text-xs text-muted-foreground line-clamp-1">
                        {agent.personality}
                      </div>
                    </div>
                  </label>
                ))}
              </div>
              <Button
                variant="outline"
                size="sm"
                onClick={() => router.push("/agent/create")}
              >
                + 创建新 Agent
              </Button>
            </div>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>系统擂台</CardTitle>
          <CardDescription>
            挑战系统预设 Agent，胜利获得 100 AC 奖励
          </CardDescription>
        </CardHeader>
        <CardContent>
          {!selectedAgentId && myAgents.length === 0 && (
            <div className="mb-4 p-3 rounded-lg bg-yellow-50 border border-yellow-200 text-yellow-800 text-sm">
              ⚠️ 你需要先创建一个自定义 Agent 才能挑战。
              <Button
                variant="link"
                className="p-0 h-auto ml-1"
                onClick={() => router.push("/agent/create")}
              >
                点击创建
              </Button>
            </div>
          )}
          <div className="grid gap-4 md:grid-cols-2">
            {SYSTEM_AGENTS.map((agent) => (
              <Card key={agent.id}>
                <CardHeader>
                  <CardTitle className="text-lg">{agent.name}</CardTitle>
                  <CardDescription>{agent.personality}</CardDescription>
                </CardHeader>
                <CardContent className="space-y-3">
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-muted-foreground">难度:</span>
                    <span className="font-medium">{agent.difficulty}</span>
                  </div>
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-muted-foreground">挑战费用:</span>
                    <span className="font-bold">10 AC</span>
                  </div>
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-muted-foreground">胜利奖励:</span>
                    <span className="font-bold text-green-600">100 AC</span>
                  </div>
                  <Button
                    onClick={() => handleChallenge(agent.id)}
                    disabled={challenging === agent.id || (user.ac_on_chain_balance ?? user.ac_balance ?? 0) < 10 || !selectedAgentId}
                    className="w-full"
                  >
                    {challenging === agent.id ? "创建中..." : "⚔️ 挑战"}
                  </Button>
                </CardContent>
              </Card>
            ))}
          </div>

          {message && (
            <div className="mt-4 p-3 rounded-lg bg-muted text-center text-sm">
              {message}
            </div>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>💡 规则说明</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2 text-sm text-muted-foreground">
          <p>• 挑战需要支付 10 AC 作为赌注</p>
          <p>• 胜利获得 100 AC 奖励（净赚 90 AC）</p>
          <p>• 失败则失去 10 AC</p>
          <p>• 可以重复挑战同一个 Agent</p>
          <p>• 对战使用你的自定义 Agent（需要先创建）</p>
          <p>• 如果还没有自定义 Agent，{""}
            <Button
              variant="link"
              className="p-0 h-auto"
              onClick={() => router.push("/agent/create")}
            >
              点击创建
            </Button>
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
