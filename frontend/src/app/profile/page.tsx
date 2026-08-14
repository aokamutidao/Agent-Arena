"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/lib/auth";
import { LoginPage } from "@/components/auth/LoginButton";
import { useI18n } from "@/lib/i18n";

interface CustomAgent {
  id: string;
  name: string;
  personality: string;
  api_endpoint: string;
  challenge_fee: number;
  currency_type: string;
  is_listed: boolean;
  wins: number;
  losses: number;
  created_at: string;
}

interface ProfileUser {
  id: string;
  address: string;
  username: string;
  ac_balance: number;
  ac_on_chain: boolean;
  ac_on_chain_balance?: number;
  ac_token_address?: string;
  ac_treasury_address?: string;
  created_at: string;
}

export default function ProfilePage() {
  const { t } = useI18n();
  const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
  const { user, token, isAuthenticated, updateBalance } = useAuth();
  const router = useRouter();
  const [claiming, setClaiming] = useState(false);
  const [message, setMessage] = useState("");
  const [messageType, setMessageType] = useState<"success" | "error" | "">("");
  const [agents, setAgents] = useState<CustomAgent[]>([]);
  const [loadingAgents, setLoadingAgents] = useState(true);
  const [refreshing, setRefreshing] = useState(false);

  // AC 完全链上：使用链上余额
  const profileUser = user as unknown as ProfileUser;
  const acBalance = profileUser?.ac_on_chain_balance ?? user?.ac_balance ?? 0;

  useEffect(() => {
    if (!token) return;

    // 获取用户的 Agent 列表
    const fetchAgents = async () => {
      try {
        const res = await fetch(`${API_URL}/api/auth/agents/my`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        if (res.ok) {
          const data = await res.json();
          setAgents(data.agents || []);
        }
      } catch (err) {
        console.error("Failed to fetch agents:", err);
      } finally {
        setLoadingAgents(false);
      }
    };

    fetchAgents();
  }, [token]);

  // 刷新用户信息（读取最新余额，包括链上）
  const refreshProfile = async () => {
    if (!token) return;
    setRefreshing(true);
    try {
      const res = await fetch(`${API_URL}/api/auth/profile`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (res.ok) {
        const data = await res.json();
        updateBalance(data.ac_balance);
      }
    } catch (err) {
      console.error("Failed to refresh profile:", err);
    } finally {
      setRefreshing(false);
    }
  };

  const handleClaimDaily = async () => {
    if (!token) return;

    setClaiming(true);
    setMessage("");
    setMessageType("");

    try {
      const res = await fetch(`${API_URL}/api/auth/claim-daily`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      const data = await res.json();
      if (!res.ok) {
        throw new Error(data.error || "领取失败");
      }

      // 领取后刷新 profile 以获取最新链上余额
      await refreshProfile();
      let msg = `✅ 成功领取 100 AC（已铸造到你的钱包）`;
      if (data.tx_hash) {
        msg += `\nTx: ${data.tx_hash}`;
      }
      if (data.next_claim) {
        const dateStr = new Date(data.next_claim).toLocaleString("zh-CN");
        msg += `\n下次可领取: ${dateStr}`;
      }
      setMessage(msg);
      setMessageType("success");
    } catch (err: any) {
      setMessage(`❌ ${err.message}`);
      setMessageType("error");
    } finally {
      setClaiming(false);
    }
  };

  if (!isAuthenticated || !user) {
    return <LoginPage />;
  }

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold">{t("profile.title")}</h1>

      <div className="grid gap-6 md:grid-cols-2">
        {/* 用户信息 */}
        <Card>
          <CardHeader>
            <CardTitle>👤 用户信息</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div>
              <div className="text-sm text-muted-foreground">用户名</div>
              <div className="font-medium">{user.username}</div>
            </div>
            <div>
              <div className="text-sm text-muted-foreground">钱包{t("profile.address")}</div>
              <div className="font-mono text-sm">{user.address}</div>
            </div>
            <div>
              <div className="text-sm text-muted-foreground">用户 ID</div>
              <div className="font-mono text-sm">{user.id}</div>
            </div>
            <div>
              <div className="text-sm text-muted-foreground">注册时间</div>
              <div className="text-sm">
                {new Date(user.created_at).toLocaleString("zh-CN")}
              </div>
            </div>
          </CardContent>
        </Card>

        {/* AC 余额 */}
        <Card>
          <CardHeader>
            <CardTitle>💰 Arena Coins</CardTitle>
            <CardDescription>
              游戏币，用于 PVE 挑战和下注
              {profileUser?.ac_on_chain && (
                <span className="ml-2 text-green-600">● 链上 ERC20</span>
              )}
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="text-4xl font-bold text-primary">
              {acBalance.toLocaleString()} AC
            </div>
            <p className="text-xs text-muted-foreground">
              💡 AC 是完全链上的 ERC20 代币，余额从 Sepolia 链上实时读取
            </p>

            {/* 链上信息 */}
            {profileUser?.ac_on_chain && (
              <div className="rounded-lg border border-green-200 bg-green-50 p-3 text-xs space-y-1">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Token</span>
                  <a
                    href={`https://sepolia.etherscan.io/token/${profileUser.ac_token_address}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="font-mono text-green-700 hover:underline"
                  >
                    {profileUser.ac_token_address?.slice(0, 8)}...{profileUser.ac_token_address?.slice(-6)}
                  </a>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-muted-foreground">刷新余额</span>
                  <Button
                    size="sm"
                    variant="ghost"
                    onClick={refreshProfile}
                    disabled={refreshing}
                    className="ml-2 h-6 px-2"
                  >
                    {refreshing ? "..." : "↻ 刷新"}
                  </Button>
                </div>
              </div>
            )}

            <div className="flex gap-2">
              <Button
                onClick={handleClaimDaily}
                disabled={claiming}
                className="flex-1"
              >
                {claiming ? "领取中（链上铸造）..." : "🎁 领取每日 100 AC"}
              </Button>
            </div>

            <p className="text-xs text-muted-foreground">
              ✅ AC 是完全链上的 ERC20 代币，领取后直接到你的钱包。无需提现，可在 MetaMask 等钱包中直接查看和使用。
            </p>

            {message && (
              <pre className={`text-xs whitespace-pre-wrap rounded p-2 ${
                messageType === "error"
                  ? "bg-red-50 text-red-800"
                  : "bg-green-50 text-green-800"
              }`}>{message}</pre>
            )}
          </CardContent>
        </Card>
      </div>

      {/* {t("profile.myAgents")} */}
      <Card>
        <CardHeader>
          <CardTitle>🤖 {t("profile.myAgents")}</CardTitle>
          <CardDescription>创建和管理你的自定义 Agent</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <Button onClick={() => router.push("/agent/create")}>
            + 创建新 Agent
          </Button>

          {loadingAgents ? (
            <p className="text-sm text-muted-foreground">{t("common.loading")}...</p>
          ) : agents.length === 0 ? (
            <p className="text-sm text-muted-foreground">
              还没有{t("agent.create")}，点击上方按钮创建你的第一个 Agent！
            </p>
          ) : (
            <div className="space-y-3">
              {agents.map((agent) => (
                <div
                  key={agent.id}
                  className="border rounded-lg p-4 space-y-2"
                >
                  <div className="flex justify-between items-start">
                    <div>
                      <h3 className="font-semibold text-lg">{agent.name}</h3>
                      <p className="text-sm text-muted-foreground line-clamp-2">
                        {agent.personality}
                      </p>
                    </div>
                    <div className="text-right">
                      <div className="text-sm font-medium">
                        {agent.wins}{t("agent.wins")} {agent.losses}{t("agent.losses")}
                      </div>
                      <div className="text-xs text-muted-foreground">
                        {agent.is_listed ? "✅ {t("agent.listed")}" : "❌ {t("agent.notListed")}"}
                      </div>
                    </div>
                  </div>
                  <div className="flex gap-2 text-xs text-muted-foreground">
                    <span>挑战费: {agent.challenge_fee} {agent.currency_type.toUpperCase()}</span>
                    <span>•</span>
                    <span>API: {agent.api_endpoint ? "自定义" : "默认"}</span>
                    <span>•</span>
                    <span>创建于: {new Date(agent.created_at).toLocaleDateString("zh-CN")}</span>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => router.push(`/agent/edit/${agent.id}`)}
                    >
                      {t("common.edit")}
                    </Button>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={async () => {
                        if (!token) return;
                        try {
                          const res = await fetch(
                            `${API_URL}/api/auth/agents/${agent.id}/listed`,
                            {
                              method: "PUT",
                              headers: {
                                "Content-Type": "application/json",
                                Authorization: `Bearer ${token}`,
                              },
                              body: JSON.stringify({ listed: !agent.is_listed }),
                            }
                          );
                          if (res.ok) {
                            // 刷新列表
                            setAgents(agents.map(a =>
                              a.id === agent.id ? { ...a, is_listed: !a.is_listed } : a
                            ));
                          }
                        } catch (err) {
                          console.error("Failed to toggle listing:", err);
                        }
                      }}
                    >
                      {agent.is_listed ? "下架" : "上架"}
                    </Button>
                    <Button
                      size="sm"
                      variant="destructive"
                      onClick={async () => {
                        if (!token) return;
                        if (!confirm("确定要删除这个 Agent 吗？")) return;
                        try {
                          const res = await fetch(
                            `${API_URL}/api/auth/agents/${agent.id}`,
                            {
                              method: "DELETE",
                              headers: {
                                Authorization: `Bearer ${token}`,
                              },
                            }
                          );
                          if (res.ok) {
                            setAgents(agents.filter(a => a.id !== agent.id));
                          }
                        } catch (err) {
                          console.error("Failed to delete agent:", err);
                        }
                      }}
                    >
                      删除
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* 下注历史 */}
      <Card>
        <CardHeader>
          <CardTitle>📊 下注历史</CardTitle>
          <CardDescription>查看你的下注记录和收益</CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">
            功能开发中...
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
