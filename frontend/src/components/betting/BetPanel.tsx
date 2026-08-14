"use client";

import { useState, useEffect, useRef } from "react";
import { useAccount, useReadContract, useWriteContract, useWaitForTransactionReceipt } from "wagmi";
import { parseUnits, formatUnits } from "viem";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { USDC_ADDRESS, ARENA_ADDRESS, POOL_ADDRESS, USDC_ABI, ARENA_ABI, POOL_ABI, STRATEGY_ENUM } from "@/lib/contracts";
import { api } from "@/lib/api";

interface BetPanelProps {
  gameId: number;
  totalBetRed: string;
  totalBetBlue: string;
  status: string;
}

export function BetPanel({ gameId, totalBetRed: propRed, totalBetBlue: propBlue, status }: BetPanelProps) {
  const { address, isConnected } = useAccount();
  const [amount, setAmount] = useState("");
  const [side, setSide] = useState<"red" | "blue">("red");
  const [strategy, setStrategy] = useState<"aggressive" | "defensive" | "tricky">("aggressive");
  const [message, setMessage] = useState("");
  const [step, setStep] = useState<"idle" | "approving" | "betting" | "done">("idle");

  const pendingBet = useRef<{ amountWei: bigint } | null>(null);
  const isLocked = status === "playing" || status === "finished";

  // === 链上读取 ===

  // 下注池 (from chain)
  const { data: gameData, refetch: refetchGame } = useReadContract({
    address: ARENA_ADDRESS,
    abi: ARENA_ABI,
    functionName: "getGame",
    args: [BigInt(gameId)],
    query: { refetchInterval: 5000 },
  });

  // 赔率 (from chain)
  const { data: oddsData, refetch: refetchOdds } = useReadContract({
    address: ARENA_ADDRESS,
    abi: ARENA_ABI,
    functionName: "getOdds",
    args: [BigInt(gameId)],
    query: { refetchInterval: 5000 },
  });

  // 解析链上数据
  const chainRed = gameData ? formatUnits((gameData as readonly bigint[])[3] as bigint, 6) : null;
  const chainBlue = gameData ? formatUnits((gameData as readonly bigint[])[4] as bigint, 6) : null;

  // 优先用链上数据，回退到 props
  const redAmount = chainRed ?? (parseInt(propRed) / 1e6).toFixed(2);
  const blueAmount = chainBlue ?? (parseInt(propBlue) / 1e6).toFixed(2);
  const totalPool = parseFloat(redAmount) + parseFloat(blueAmount);
  const redPct = totalPool > 0 ? Math.round((parseFloat(redAmount) / totalPool) * 100) : 50;
  const bluePct = 100 - redPct;

  // 赔率（链上 scaled by 1e18）
  const oddsRed = oddsData ? Number((oddsData as readonly bigint[])[0]) / 1e18 : 0;
  const oddsBlue = oddsData ? Number((oddsData as readonly bigint[])[1]) / 1e18 : 0;

  // USDC 余额
  const { data: balance, refetch: refetchBalance } = useReadContract({
    address: USDC_ADDRESS,
    abi: USDC_ABI,
    functionName: "balanceOf",
    args: address ? [address] : undefined,
    query: { enabled: !!address },
  });

  // 授权额度
  const { data: allowance, refetch: refetchAllowance } = useReadContract({
    address: USDC_ADDRESS,
    abi: USDC_ABI,
    functionName: "allowance",
    args: address ? [address, ARENA_ADDRESS] : undefined,
    query: { enabled: !!address },
  });

  // 游戏信息 (from chain) - 用于计算奖励
  const { data: poolGameData, refetch: refetchPoolGame } = useReadContract({
    address: POOL_ADDRESS,
    abi: POOL_ABI,
    functionName: "games",
    args: [BigInt(gameId)],
    query: { refetchInterval: 5000 },
  });

  // 我的下注记录 (from chain)
  const { data: myBet, refetch: refetchMyBet } = useReadContract({
    address: POOL_ADDRESS,
    abi: POOL_ABI,
    functionName: "bets",
    args: address ? [BigInt(gameId), address] : undefined,
    query: { enabled: !!address, refetchInterval: 5000 },
  });

  // 解析我的下注
  const hasBet = myBet ? (myBet as readonly [number, bigint, boolean])[1] > BigInt(0) : false;
  const myBetSide = myBet ? (myBet as readonly [number, bigint, boolean])[0] : 0; // 1=Red, 2=Blue
  const myBetAmount = myBet ? formatUnits((myBet as readonly [number, bigint, boolean])[1], 6) : "0";
  const myBetClaimed = myBet ? (myBet as readonly [number, bigint, boolean])[2] : false;

  // 手动计算可领取的奖励（因为 deployed BettingPool 没有 getReward 函数）
  // 公式：reward = bet.amount * (totalPool - fee) / winnerPool
  let rewardAmount = "0";
  if (poolGameData && myBet && hasBet && status === "finished" && !myBetClaimed) {
    // games() returns: (gameId, agentRed, agentBlue, totalBetRed, totalBetBlue, bettingDeadline, status, winner)
    const [, , , totalBetRed, totalBetBlue, , , winner] = poolGameData as readonly [
      bigint, `0x${string}`, `0x${string}`, bigint, bigint, bigint, number, number
    ];
    // winner: 0=None, 1=Red, 2=Blue; betSide: 0=None, 1=Red, 2=Blue
    if (myBetSide === winner && winner !== 0) {
      const betAmount = (myBet as readonly [number, bigint, boolean])[1];
      const totalPool = totalBetRed + totalBetBlue;
      const fee = (totalPool * BigInt(500)) / BigInt(10000); // 5%
      const distributable = totalPool - fee;
      const winnerPool = winner === 1 ? totalBetRed : totalBetBlue;
      if (winnerPool > BigInt(0)) {
        const reward = (betAmount * distributable) / winnerPool;
        rewardAmount = formatUnits(reward, 6);
      }
    }
  }

  // === 写入 ===
  const { writeContract: writeMint, data: mintTxHash, isPending: mintPending } = useWriteContract();
  const { isLoading: mintConfirming, isSuccess: mintSuccess } = useWaitForTransactionReceipt({ hash: mintTxHash });

  const { writeContract: writeApprove, data: approveTxHash, isPending: approvePending } = useWriteContract();
  const { isLoading: approveConfirming, isSuccess: approveSuccess } = useWaitForTransactionReceipt({ hash: approveTxHash });

  const { writeContract: writeBet, data: betTxHash, isPending: betPending } = useWriteContract();
  const { isLoading: betConfirming, isSuccess: betSuccess } = useWaitForTransactionReceipt({ hash: betTxHash });

  // 领奖交易
  const { writeContract: writeClaim, data: claimTxHash, isPending: claimPending } = useWriteContract();
  const { isLoading: claimConfirming, isSuccess: claimSuccess } = useWaitForTransactionReceipt({ hash: claimTxHash });

  // === Effects ===
  useEffect(() => {
    if (mintSuccess) {
      refetchBalance();
      setMessage("✅ 已铸造 100 测试 USDC");
    }
  }, [mintSuccess, refetchBalance]);

  useEffect(() => {
    if (approveSuccess && pendingBet.current) {
      const { amountWei } = pendingBet.current;
      pendingBet.current = null;
      refetchAllowance();
      setMessage("✅ 授权成功，正在下注...");
      setStep("betting");
      writeBet({
        address: ARENA_ADDRESS,
        abi: ARENA_ABI,
        functionName: "betAndVote",
        args: [BigInt(gameId), side === "red", amountWei, STRATEGY_ENUM[strategy]],
      });
    }
  }, [approveSuccess, gameId, side, strategy, writeBet, refetchAllowance]);

  useEffect(() => {
    if (betSuccess) {
      refetchBalance();
      refetchAllowance();
      refetchGame();
      refetchOdds();
      refetchMyBet();
      refetchPoolGame();
      setMessage("✅ 链上下注成功！");
      setAmount("");
      setStep("done");

      // 同步策略投票到后端（后端内存需要知道策略权重，否则前端展示全0）
      api.voteStrategy({
        game_id: gameId,
        side,
        strategy,
      }).catch((err) => {
        console.warn("[BetPanel] voteStrategy sync failed:", err);
        // 不阻断用户体验，策略同步失败不影响下注成功
      });
    }
  }, [betSuccess, refetchBalance, refetchAllowance, refetchGame, refetchOdds, refetchMyBet, refetchPoolGame, gameId, side, strategy]);

  useEffect(() => {
    if (claimSuccess) {
      refetchBalance();
      refetchMyBet();
      refetchPoolGame();
      setMessage("✅ 奖金已领取！");
    }
  }, [claimSuccess, refetchBalance, refetchMyBet, refetchPoolGame]);

  const usdcBalance = balance ? formatUnits(balance as bigint, 6) : "0";
  const usdcAllowance = allowance ? (allowance as bigint) : BigInt(0);

  // 预估收益
  const estimateReward = (betAmount: number, betSide: "red" | "blue") => {
    if (totalPool === 0) return 0;
    const newPool = totalPool + betAmount;
    const newSideTotal = (betSide === "red" ? parseFloat(redAmount) : parseFloat(blueAmount)) + betAmount;
    if (newSideTotal === 0) return 0;
    const afterFee = newPool * 0.95;
    return (betAmount / newSideTotal) * afterFee;
  };

  const handleMint = () => {
    if (!address) return;
    setMessage("");
    writeMint({
      address: USDC_ADDRESS,
      abi: USDC_ABI,
      functionName: "mint",
      args: [address, parseUnits("100", 6)],
    });
  };

  const handleBet = () => {
    if (!amount || !address || parseFloat(amount) <= 0) return;
    setMessage("");
    const amountWei = parseUnits(amount, 6);

    if (balance && amountWei > (balance as bigint)) {
      setMessage("❌ USDC 余额不足，请先点击 [领 100 测试 USDC]");
      return;
    }

    if (amountWei > usdcAllowance) {
      pendingBet.current = { amountWei };
      setStep("approving");
      writeApprove({
        address: USDC_ADDRESS,
        abi: USDC_ABI,
        functionName: "approve",
        args: [ARENA_ADDRESS, amountWei],
      });
      setMessage("⏳ 第 1 步：授权 Arena 合约...请在钱包确认");
      return;
    }

    setStep("betting");
    writeBet({
      address: ARENA_ADDRESS,
      abi: ARENA_ABI,
      functionName: "betAndVote",
      args: [BigInt(gameId), side === "red", amountWei, STRATEGY_ENUM[strategy]],
    });
    setMessage("⏳ 正在下注...请在钱包确认交易");
  };

  const handleClaim = () => {
    setMessage("");
    writeClaim({
      address: POOL_ADDRESS,
      abi: POOL_ABI,
      functionName: "claim",
      args: [BigInt(gameId)],
    });
  };

  const isBusy = mintPending || mintConfirming || approvePending || approveConfirming || betPending || betConfirming || claimPending || claimConfirming;
  const potentialReward = amount ? estimateReward(parseFloat(amount), side) : 0;

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-lg">📊 链上下注</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* 池占比条 */}
        <div>
          <div className="flex justify-between text-sm mb-1">
            <span className="text-red-500 font-medium">🔴 {redAmount} USDC ({redPct}%)</span>
            <span className="text-blue-500 font-medium">🔵 {blueAmount} USDC ({bluePct}%)</span>
          </div>
          <div className="h-3 bg-secondary rounded-full overflow-hidden flex">
            <div className="h-full bg-red-500 transition-all duration-500" style={{ width: `${redPct}%` }} />
            <div className="h-full bg-blue-500 transition-all duration-500" style={{ width: `${bluePct}%` }} />
          </div>
          <div className="text-center text-xs text-muted-foreground mt-1">
            总池: {totalPool.toFixed(2)} USDC · Sepolia 链上
          </div>
        </div>

        {/* 赔率展示 */}
        <div className="grid grid-cols-2 gap-3">
          <div className="rounded-lg bg-red-500/5 border border-red-500/20 p-3 text-center">
            <div className="text-xs text-muted-foreground mb-1">🔴 红方赔率</div>
            <div className="text-xl font-bold text-red-500">
              {oddsRed > 0 ? `${oddsRed.toFixed(2)}x` : "—"}
            </div>
          </div>
          <div className="rounded-lg bg-blue-500/5 border border-blue-500/20 p-3 text-center">
            <div className="text-xs text-muted-foreground mb-1">🔵 蓝方赔率</div>
            <div className="text-xl font-bold text-blue-500">
              {oddsBlue > 0 ? `${oddsBlue.toFixed(2)}x` : "—"}
            </div>
          </div>
        </div>

        {/* 预估收益 */}
        {amount && parseFloat(amount) > 0 && potentialReward > 0 && (
          <div className="rounded bg-primary/5 border border-primary/20 px-3 py-2 text-center">
            <span className="text-xs text-muted-foreground">
              下注 {amount} USDC 到{side === "red" ? "🔴红方" : "🔵蓝方"}，预估收益:
            </span>
            <span className="text-sm font-bold text-primary ml-1">
              {potentialReward.toFixed(2)} USDC
            </span>
            <span className="text-xs text-muted-foreground ml-1">
              (含 5% 协议费)
            </span>
          </div>
        )}

        {/* 我的下注记录 */}
        {hasBet && isConnected && (
          <div className="rounded-lg bg-primary/5 border border-primary/20 p-3 space-y-2">
            <div className="text-sm font-medium flex items-center justify-between">
              <span>🎯 我的下注</span>
              <span className="text-xs text-muted-foreground">链上</span>
            </div>
            <div className="text-sm flex items-center gap-2">
              <span>{myBetSide === 1 ? "🔴 红方" : "🔵 蓝方"}</span>
              <span className="text-muted-foreground">·</span>
              <span className="font-mono">{myBetAmount} USDC</span>
            </div>

            {status === "finished" && (
              <>
                {parseFloat(rewardAmount) > 0 && !myBetClaimed ? (
                  <div className="space-y-2 pt-2 border-t">
                    <div className="text-sm text-green-600 font-medium">
                      🎉 可领取 {rewardAmount} USDC
                    </div>
                    <Button
                      size="sm"
                      className="w-full"
                      onClick={handleClaim}
                      disabled={isBusy}
                    >
                      {claimPending || claimConfirming ? "领取中...请确认交易" : "💰 领取奖金"}
                    </Button>
                  </div>
                ) : myBetClaimed ? (
                  <div className="text-sm text-muted-foreground pt-2 border-t">
                    ✅ 奖金已领取
                  </div>
                ) : (
                  <div className="text-sm text-muted-foreground pt-2 border-t">
                    未获胜，下次加油！
                  </div>
                )}
              </>
            )}
          </div>
        )}

        {/* 下注区 */}
        {!isLocked && (
          <div className="space-y-3 pt-2 border-t">
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">钱包状态</span>
              {isConnected ? (
                <div className="flex items-center gap-3">
                  <span className="text-sm">💰 {parseFloat(usdcBalance).toFixed(2)} USDC</span>
                  <span className="text-xs text-green-500">● 已连接</span>
                </div>
              ) : (
                <span className="text-sm text-yellow-500">请先点右上角连接钱包</span>
              )}
            </div>

            {isConnected && (
              <>
                {parseFloat(usdcBalance) === 0 && step === "idle" && (
                  <div className="flex items-center justify-between rounded bg-yellow-500/10 border border-yellow-500/20 px-3 py-2">
                    <span className="text-xs text-yellow-600">💡 你还没有测试 USDC，点击领取</span>
                    <Button variant="outline" size="sm" onClick={handleMint} disabled={isBusy}>
                      {mintPending || mintConfirming ? "铸造中..." : "🚰 领 100 USDC"}
                    </Button>
                  </div>
                )}

                <div>
                  <label className="text-sm text-muted-foreground mb-1 block">下注金额 (USDC)</label>
                  <Input
                    type="number"
                    placeholder="输入金额"
                    value={amount}
                    onChange={(e) => setAmount(e.target.value)}
                    min="0"
                    disabled={isBusy}
                  />
                </div>

                <div>
                  <label className="text-sm text-muted-foreground mb-2 block">选择阵营</label>
                  {hasBet && (
                    <div className="mb-2 p-2 rounded bg-yellow-500/10 border border-yellow-500/30 text-xs text-yellow-600">
                      ⚠️ 注意：你已经下注过。再次下注会覆盖之前的阵营选择，金额会累加。
                    </div>
                  )}
                  <div className="flex gap-3">
                    <button
                      onClick={() => setSide("red")}
                      disabled={isBusy}
                      className={`flex-1 py-2 px-3 rounded border text-sm font-medium transition-colors ${
                        side === "red" ? "border-red-500 bg-red-500/10 text-red-500" : "border-border hover:bg-accent"
                      }`}
                    >
                      🔴 Red
                    </button>
                    <button
                      onClick={() => setSide("blue")}
                      disabled={isBusy}
                      className={`flex-1 py-2 px-3 rounded border text-sm font-medium transition-colors ${
                        side === "blue" ? "border-blue-500 bg-blue-500/10 text-blue-500" : "border-border hover:bg-accent"
                      }`}
                    >
                      🔵 Blue
                    </button>
                  </div>
                </div>

                <div>
                  <label className="text-sm text-muted-foreground mb-2 block">策略偏好</label>
                  <div className="flex gap-2">
                    {(["aggressive", "defensive", "tricky"] as const).map((s) => (
                      <button
                        key={s}
                        onClick={() => setStrategy(s)}
                        disabled={isBusy}
                        className={`flex-1 py-1.5 px-2 rounded border text-xs font-medium transition-colors ${
                          strategy === s ? "border-primary bg-primary/10 text-primary" : "border-border hover:bg-accent"
                        }`}
                      >
                        {s === "aggressive" && "🔥 激进"}
                        {s === "defensive" && "🛡️ 稳健"}
                        {s === "tricky" && "🎲 诡道"}
                      </button>
                    ))}
                  </div>
                </div>

                <Button
                  className="w-full"
                  disabled={!amount || parseFloat(amount) <= 0 || isBusy}
                  onClick={handleBet}
                >
                  {step === "approving" && (approvePending || approveConfirming)
                    ? "授权中...请确认交易"
                    : step === "betting" && (betPending || betConfirming)
                    ? "下注中...请确认交易"
                    : "确认下注（链上）"}
                </Button>

                {message && <p className="text-xs text-center">{message}</p>}
              </>
            )}

            {!isConnected && (
              <p className="text-xs text-muted-foreground text-center py-2">
                连接钱包后在 Sepolia 链上下注 · 可用 MockUSDC 免费领取测试币
              </p>
            )}
          </div>
        )}

        {isLocked && (
          <div className="text-center text-sm text-muted-foreground py-2">
            {status === "finished" ? "对局已结束" : "下注已锁定"}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
