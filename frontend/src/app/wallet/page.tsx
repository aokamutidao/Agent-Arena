"use client";

import { useEffect, useState, useRef } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Coins, Gem, Wallet, ArrowUpRight, Trophy, Swords, RefreshCw, TrendingUp } from "lucide-react";
import { useAuth } from "@/lib/auth";
import { LoginPage } from "@/components/auth/LoginButton";
import Link from "next/link";
import { useAccount, useReadContract, useReadContracts, useWriteContract, useWaitForTransactionReceipt } from "wagmi";
import { formatUnits, parseUnits } from "viem";
import { POOL_ADDRESS, POOL_ABI, USDC_ADDRESS, USDC_ABI } from "@/lib/contracts";
import { useI18n } from "@/lib/i18n";

interface EarningEntry {
  challenge_id: string;
  game_id: number;
  opponent_id: string;
  opponent_type: string;
  stake: number;
  currency: string; // "ac" | "usdc"
  winner: string;   // "challenger" | "opponent" | "draw" | ""
  reward: number;
  status: string;
  created_at: string;
  finished_at: string;
}

interface EarningsResponse {
  earnings: EarningEntry[];
  total_reward_ac: number;
  total_reward_usdc: number;
}

// 系统 Agent 名称映射
const SYSTEM_AGENT_NAMES: Record<string, string> = {
  berserker: "🔥 狂战士 Berserker",
  tactician: "🧠 战术家 Tactician",
  trickster: "🎭 诡术师 Trickster",
  defender: "🛡️ 守护者 Defender",
};

const opponentLabel = (e: EarningEntry) => {
  if (e.opponent_type === "system" && SYSTEM_AGENT_NAMES[e.opponent_id]) {
    return SYSTEM_AGENT_NAMES[e.opponent_id];
  }
  if (e.opponent_type === "user") {
    return `👤 ${e.opponent_id.slice(0, 8)}...`;
  }
  return e.opponent_id;
};

export default function WalletPage() {
  const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
  const { user, token, isAuthenticated, refreshUser } = useAuth();
  const { address, isConnected } = useAccount();
  const { t } = useI18n();
  const [earnings, setEarnings] = useState<EarningsResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [claimMessage, setClaimMessage] = useState("");
  const hasFetched = useRef(false);

  // ===== 铸造测试 USDC =====
  const { data: usdcBalance, refetch: refetchUsdcBalance } = useReadContract({
    address: USDC_ADDRESS,
    abi: USDC_ABI,
    functionName: "balanceOf",
    args: address ? [address] : undefined,
    query: { enabled: !!address, refetchInterval: 10000 },
  });

  const { writeContract: writeMint, data: mintTxHash, isPending: mintPending } = useWriteContract();
  const { isLoading: mintConfirming, isSuccess: mintSuccess } = useWaitForTransactionReceipt({ hash: mintTxHash });

  const handleMintUSDC = () => {
    if (!address) return;
    writeMint({
      address: USDC_ADDRESS,
      abi: USDC_ABI,
      functionName: "mint",
      args: [address, parseUnits("100", 6)],
    });
  };

  useEffect(() => {
    if (mintSuccess) {
      refetchUsdcBalance();
      refreshUser();
      setClaimMessage("✅ 已铸造 100 测试 USDC");
      setTimeout(() => setClaimMessage(""), 4000);
    }
  }, [mintSuccess, refetchUsdcBalance, refreshUser]);

  // ===== 可{t("wallet.claim")}的 USDC 赌注 =====
  // 架构：后端返回用户的下注 gameID 列表 → 前端查链上 bets() + games() 手动计算奖励
  // （因为 deployed BettingPool 没有 getReward 函数）
  const [userBetGameIds, setUserBetGameIds] = useState<number[]>([]);
  const [queryEnabled, setQueryEnabled] = useState(false);
  const [querying, setQuerying] = useState(false);

  // 根据后端返回的 gameID 列表构造链上查询
  const gameIdsBigInt = userBetGameIds.map((id) => BigInt(id));

  const betContracts = gameIdsBigInt.map((gid) => ({
    address: POOL_ADDRESS,
    abi: POOL_ABI,
    functionName: "bets",
    args: [gid, address as `0x${string}`],
  }));

  const gameContracts = gameIdsBigInt.map((gid) => ({
    address: POOL_ADDRESS,
    abi: POOL_ABI,
    functionName: "games",
    args: [gid],
  }));

  const { data: betsData, isLoading: betsLoading } = useReadContracts({
    contracts: address && queryEnabled && userBetGameIds.length > 0 ? betContracts : [],
    query: { enabled: !!address && queryEnabled && userBetGameIds.length > 0, gcTime: 60_000, staleTime: 60_000 },
  });

  const { data: gamesData, isLoading: gamesLoading } = useReadContracts({
    contracts: address && queryEnabled && userBetGameIds.length > 0 ? gameContracts : [],
    query: { enabled: !!address && queryEnabled && userBetGameIds.length > 0, gcTime: 60_000, staleTime: 60_000 },
  });

  // 查询完成时清除 loading
  useEffect(() => {
    if (queryEnabled && !betsLoading && !gamesLoading) {
      setQuerying(false);
    }
  }, [queryEnabled, betsLoading, gamesLoading]);

  // 手动触发查询：获取所有已结束游戏的 ID，然后逐个查链上奖励
  const handleQueryClaimable = async () => {
    if (!address || !token) return;
    setQuerying(true);
    setQueryEnabled(false); // 先禁用链上查询

    try {
      // 1. 从后端获取所有已结束游戏的历史记录（持久化在 DB 中）
      const res = await fetch(`${API_URL}/api/game-history?limit=100`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (res.ok) {
        const data = await res.json();
        const history = data.history || [];
        // 提取所有已结束游戏的 gameID
        const ids = history.map((g: { game_id: number }) => g.game_id) as number[];
        setUserBetGameIds(ids);
      }
    } catch (err) {
      console.error("fetch game history failed:", err);
    }

    // 2. 启用链上查询（此时 gameIds 已更新）
    setQueryEnabled(true);
  };

  // 解析链上数据：找出可{t("wallet.claim")}的赌注
  // 手动计算奖励（因为 deployed BettingPool 没有 getReward 函数）
  // 公式：reward = bet.amount * (totalPool - fee) / winnerPool
  // 其中 fee = totalPool * 5% = totalPool * 500 / 10000
  type ClaimableBet = { gameId: number; reward: bigint };
  const claimableBets: ClaimableBet[] = [];
  let totalClaimable = BigInt(0);

  if (betsData && gamesData && address && userBetGameIds.length > 0) {
    for (let i = 0; i < userBetGameIds.length; i++) {
      const betResult = betsData[i] as any;
      const gameResult = gamesData[i] as any;
      if (!betResult?.result || !gameResult?.result) continue;

      // Parse bet info: (side: uint8, amount: uint256, claimed: bool)
      const [betSide, betAmount, betClaimed] = betResult.result as [number, bigint, boolean];

      // Skip if no bet or already claimed
      if (betAmount === BigInt(0) || betClaimed) continue;

      // Parse game info: (gameId, agentRed, agentBlue, totalBetRed, totalBetBlue, bettingDeadline, status, winner)
      const [, , , totalBetRed, totalBetBlue, , status, winner] = gameResult.result as [
        bigint, `0x${string}`, `0x${string}`, bigint, bigint, bigint, number, number
      ];

      // status: 0=Open, 1=Locked, 2=Finished
      // winner: 0=None, 1=Red, 2=Blue
      // betSide: 0=None, 1=Red, 2=Blue

      // Game must be finished
      if (status !== 2) continue;

      // User must be on the winning side
      if (betSide !== winner || winner === 0) continue;

      // Calculate reward
      const totalPool = totalBetRed + totalBetBlue;
      const fee = (totalPool * BigInt(500)) / BigInt(10000); // 5% protocol fee
      const distributable = totalPool - fee;
      const winnerPool = winner === 1 ? totalBetRed : totalBetBlue;

      // Avoid division by zero
      if (winnerPool === BigInt(0)) continue;

      const reward = (betAmount * distributable) / winnerPool;

      if (reward > BigInt(0)) {
        claimableBets.push({ gameId: userBetGameIds[i], reward });
        totalClaimable += reward;
      }
    }
  }

  // Claim 写入
  const {
    writeContract: writeClaim,
    data: claimTxHash,
    isPending: isClaimPending,
    error: claimWriteError,
  } = useWriteContract();

  const { isLoading: isClaimConfirming, isSuccess: claimConfirmed } =
    useWaitForTransactionReceipt({ hash: claimTxHash });

  const [claimGameId, setClaimGameId] = useState<number | null>(null);

  const handleClaim = (gameId: number) => {
    setClaimGameId(gameId);
    setClaimMessage("");
    writeClaim({
      address: POOL_ADDRESS,
      abi: POOL_ABI,
      functionName: "claim",
      args: [BigInt(gameId)],
    });
  };

  // 监听 claim 交易结果
  useEffect(() => {
    if (claimConfirmed) {
      setClaimMessage(`✅ Game #${claimGameId} {t("wallet.claim")}成功！`);
      setClaimGameId(null);
      // 延迟清除消息
      setTimeout(() => setClaimMessage(""), 4000);
    }
  }, [claimConfirmed, claimGameId]);

  useEffect(() => {
    if (claimWriteError) {
      setClaimMessage(`❌ {t("wallet.claim")}失败: ${claimWriteError.message}`);
      setTimeout(() => setClaimMessage(""), 5000);
    }
  }, [claimWriteError]);

  useEffect(() => {
    if (!token) return;
    const load = async () => {
      // 只在首次加载时显示 loading，后续刷新保持旧数据
      if (!hasFetched.current) setLoading(true);
      try {
        const [profileRes, earnRes] = await Promise.all([
          fetch(`${API_URL}/api/auth/profile`, {
            headers: { Authorization: `Bearer ${token}` },
          }),
          fetch(`${API_URL}/api/auth/earnings`, {
            headers: { Authorization: `Bearer ${token}` },
          }),
        ]);
        if (profileRes.ok) {
          // 触发 auth context 更新，使 navbar 余额刷新
          await refreshUser();
        }
        if (earnRes.ok) {
          const data = await earnRes.json();
          setEarnings(data);
        }
      } catch (err) {
        console.error("wallet load failed:", err);
      } finally {
        hasFetched.current = true;
        setLoading(false);
        setRefreshing(false);
      }
    };
    load();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [token]); // 只依赖 token，避免 refreshUser 导致的无限循环

  // 手动刷新
  const handleRefresh = () => {
    setRefreshing(true);
    if (!token) return;
    const load = async () => {
      try {
        const [profileRes, earnRes] = await Promise.all([
          fetch(`${API_URL}/api/auth/profile`, {
            headers: { Authorization: `Bearer ${token}` },
          }),
          fetch(`${API_URL}/api/auth/earnings`, {
            headers: { Authorization: `Bearer ${token}` },
          }),
        ]);
        if (profileRes.ok) {
          await refreshUser();
        }
        if (earnRes.ok) {
          const data = await earnRes.json();
          setEarnings(data);
        }
      } catch (err) {
        console.error("wallet refresh failed:", err);
      } finally {
        setRefreshing(false);
      }
    };
    load();
  };

  if (!isAuthenticated || !user) {
    return <LoginPage />;
  }

  const profit = (e: EarningEntry) => {
    if (e.status !== "finished") return 0;
    if (e.winner === "challenger") return e.reward - e.stake;
    if (e.winner === "draw") return 0; // 退回本金
    return -e.stake; // 输了
  };

  return (
    <div className="container mx-auto px-4 py-8 space-y-6">
      <div>
        <h1 className="text-3xl font-bold flex items-center gap-2">
          <Wallet className="h-8 w-8" />
          {t("wallet.title")}
        </h1>
        <p className="text-muted-foreground mt-1">
          {t("profile.address")}<span className="font-mono">{user.address}</span>
        </p>
      </div>

      {/* 余额卡片 */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card>
          <CardHeader className="pb-2">
            <CardDescription className="flex items-center gap-2">
              <Coins className="h-4 w-4 text-yellow-500" />
              USDC（链上）
            </CardDescription>
            <CardTitle className="text-3xl font-mono">
              {usdcBalance !== undefined
                ? formatUnits(usdcBalance as bigint, 6)
                : user.usdc_balance_raw
                ? formatUnits(BigInt(user.usdc_balance_raw), 6)
                : user.usdc_balance != null
                ? user.usdc_balance.toLocaleString()
                : "—"}
            </CardTitle>
          </CardHeader>
          <CardContent className="text-xs text-muted-foreground space-y-2">
            <div>{t("wallet.forBetting")}</div>
            <Button
              size="sm"
              variant="outline"
              onClick={handleMintUSDC}
              disabled={!isConnected || mintPending || mintConfirming}
              className="w-full"
            >
              {mintPending || mintConfirming ? t("wallet.minting") : t("wallet.mintTest")}
            </Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardDescription className="flex items-center gap-2">
              <Gem className="h-4 w-4 text-purple-500" />
              AC（链上余额）
            </CardDescription>
            <CardTitle className="text-3xl font-mono">
              {(user.ac_on_chain_balance ?? user.ac_balance ?? 0).toLocaleString()}
            </CardTitle>
          </CardHeader>
          <CardContent className="text-xs text-muted-foreground">
            {t("wallet.erc20Desc")}
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardDescription className="flex items-center gap-2">
              <Trophy className="h-4 w-4 text-green-500" />
              {t("wallet.totalEarnings")}
            </CardDescription>
            <CardTitle className="text-3xl font-mono">
              {earnings
                ? `+${earnings.total_reward_ac.toLocaleString()} AC`
                : "—"}
            </CardTitle>
          </CardHeader>
          <CardContent className="text-xs text-muted-foreground">
            {earnings
              ? `+${earnings.total_reward_usdc.toLocaleString()} USDC`
              : t("wallet.challengeReward")}
          </CardContent>
        </Card>
      </div>

      {/* 可{t("wallet.claim")}的 USDC 赌注 */}
      <Card>
        <CardHeader className="pb-2">
          <div className="flex items-center justify-between">
            <div>
              <CardDescription className="flex items-center gap-2">
                <TrendingUp className="h-4 w-4 text-blue-500" />
                {t("wallet.claimableBets")}
              </CardDescription>
              <CardTitle className="text-2xl font-mono mt-1">
                {queryEnabled ? `${formatUnits(totalClaimable, 6)} USDC` : "— USDC"}
              </CardTitle>
            </div>
            <Button
              variant="outline"
              size="sm"
              onClick={handleQueryClaimable}
              disabled={!address || querying}
            >
              {querying ? t("wallet.querying") : t("wallet.queryClaimable")}
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          {claimMessage && (
            <div className="mb-3 p-2 rounded bg-blue-50 dark:bg-blue-950 text-blue-700 dark:text-blue-300 text-sm">
              {claimMessage}
            </div>
          )}
          {!address ? (
            <p className="text-muted-foreground text-sm">{t("wallet.connectToView")}</p>
          ) : !queryEnabled ? (
            <p className="text-muted-foreground text-sm">
              点击「{t("wallet.queryClaimable")}」按钮，查询链上可{t("wallet.claim")}的赢利。
            </p>
          ) : querying ? (
            <p className="text-muted-foreground text-sm">{t("wallet.queryingChain")}</p>
          ) : claimableBets.length === 0 ? (
            <p className="text-muted-foreground text-sm">
              {t("wallet.noClaimable")}。参与下注并获胜后可在此{t("wallet.claim")}赢利。
            </p>
          ) : (
            <div className="space-y-2">
              {claimableBets.map((bet) => (
                <div
                  key={bet.gameId}
                  className="flex items-center justify-between p-2 rounded border"
                >
                  <div>
                    <div className="font-medium text-sm">
                      Game #{bet.gameId}
                    </div>
                    <div className="text-xs text-muted-foreground">
                      {t("wallet.canClaim")}: <span className="font-mono text-green-600">{formatUnits(bet.reward, 6)} USDC</span>
                    </div>
                  </div>
                  <Button
                    size="sm"
                    onClick={() => handleClaim(bet.gameId)}
                    disabled={isClaimPending || claimGameId === bet.gameId}
                  >
                    {claimGameId === bet.gameId && (isClaimPending || isClaimConfirming)
                      ? t("wallet.claiming")
                      : t("wallet.claim")}
                  </Button>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* 操作按钮 */}
      <div className="flex gap-3">
        <Button variant="outline" onClick={handleRefresh} disabled={refreshing}>
          <RefreshCw className={`h-4 w-4 mr-2 ${refreshing ? 'animate-spin' : ''}`} />
          {refreshing ? '{t("wallet.refreshing")}' : '{t("wallet.refresh")}'}
        </Button>
        <Link href="/profile">
          <Button variant="outline">
            <ArrowUpRight className="h-4 w-4 mr-2" />
            {t("wallet.claim")}每日 AC
          </Button>
        </Link>
        <Link href="/pve">
          <Button>
            <Swords className="h-4 w-4 mr-2" />
            {t("wallet.startChallenge")}
          </Button>
        </Link>
      </div>

      {/* {t("wallet.history")} */}
      <Card>
        <CardHeader>
          <CardTitle>{t("wallet.history")}</CardTitle>
          <CardDescription>{t("wallet.recentRecords")}</CardDescription>
        </CardHeader>
        <CardContent>
          {loading ? (
            <p className="text-muted-foreground">{t("common.loading")}</p>
          ) : !earnings || earnings.earnings.length === 0 ? (
            <p className="text-muted-foreground">{t("wallet.noRecords")}。完成第一次挑战吧！</p>
          ) : (
            <div className="space-y-2">
              {earnings.earnings.map((e) => {
                const p = profit(e);
                const isWin = p > 0;
                const isLoss = p < 0;
                const currencySymbol = e.currency === "usdc" ? "USDC" : "AC";
                return (
                  <div
                    key={e.challenge_id}
                    className="flex items-center justify-between p-3 rounded-lg border hover:bg-muted/50 transition-colors"
                  >
                    <div className="flex items-center gap-3">
                      <div
                        className={`w-8 h-8 rounded-full flex items-center justify-center text-white text-sm font-bold ${
                          isWin ? "bg-green-500" : isLoss ? "bg-red-500" : "bg-gray-500"
                        }`}
                      >
                        {isWin ? "W" : isLoss ? "L" : "D"}
                      </div>
                      <div>
                        <div className="font-medium">
                          vs {opponentLabel(e)}
                          <span className="text-muted-foreground text-xs ml-2">
                            ({e.opponent_type === 'system' ? '系统' : '用户'})
                          </span>
                        </div>
                        <div className="text-xs text-muted-foreground">
                          Game #{e.game_id} · 押金 {e.stake} {currencySymbol}
                          {e.finished_at
                            ? ` · ${new Date(e.finished_at).toLocaleDateString()}`
                            : ""}
                        </div>
                      </div>
                    </div>
                    <div className="text-right">
                      <div
                        className={`font-mono font-bold ${
                          isWin ? "text-green-600" : isLoss ? "text-red-600" : "text-muted-foreground"
                        }`}
                      >
                        {p > 0 ? "+" : ""}
                        {p} {currencySymbol}
                      </div>
                      {e.status !== "finished" && (
                        <div className="text-xs text-muted-foreground">
                          {e.status === "playing" ? t("wallet.playing") : e.status}
                        </div>
                      )}
                      <Link
                        href={`/game/${e.game_id}`}
                        className="text-xs text-blue-500 hover:underline"
                      >
                        {t("wallet.viewGame")}
                      </Link>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
