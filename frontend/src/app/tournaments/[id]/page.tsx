"use client";
import { useEffect, useState, type ReactNode } from "react";
import {
  ArrowDownToLine,
  ArrowLeft,
  ArrowRight,
  Check,
  CheckCircle2,
  Clock3,
  Coins,
  FileCheck2,
  Info,
  Play,
  ShieldCheck,
  Trophy,
  Users,
  Wallet,
} from "lucide-react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { useAccount, useBlock, useReadContract } from "wagmi";
import { useQueryClient } from "@tanstack/react-query";
import { encodeFunctionData, zeroAddress, type Address } from "viem";
import { tournamentManagerAbi, courseTokenAbi } from "@/lib/course-abi";
import { useAuth } from "@/lib/auth";
import {
  ApiError,
  dateTime,
  errorMessage,
  money,
  parseMoney,
  request,
  sameAddress,
  shortAddress,
  registrationOpen,
  reviewComplete as isReviewComplete,
  canExpire,
  verifyEvidence,
  type Tournament,
  type AgentProfile,
  type AuditEvent,
  type MatchResult,
  type Evidence,
  type EvidenceNote,
} from "@/lib/course";
import {
  useAuthenticatedResource,
  useCourseConfig,
  useCourseResource,
} from "@/hooks/useCourse";
import { useWalletTransaction } from "@/hooks/useWalletTransaction";
import { DeploymentNotice } from "@/components/course/DeploymentNotice";
import { TransactionStatus } from "@/components/course/TransactionStatus";
import { EventTable } from "@/components/course/EventTable";
import { MatchReplay } from "@/components/course/MatchReplay";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { WalletConnect } from "@/components/layout/WalletConnect";
import { LoginButton } from "@/components/auth/LoginButton";
import { TournamentStatusBadge } from "@/components/course/TournamentStatusBadge";
import { getTournamentPresentation } from "@/lib/course-presentation";
type IdAction =
  | "withdrawParticipant"
  | "startTournament"
  | "approveResult"
  | "expireTournament"
  | "claimPrize"
  | "claimRefund";
export default function TournamentPage() {
  const params = useParams();
  const id = Number(params.id);
  const validId = Number.isSafeInteger(id) && id > 0;
  const { address } = useAccount();
  const { token, user } = useAuth();
  const { data: config } = useCourseConfig();
  const tournament = useCourseResource<Tournament>(
    `/api/tournaments/${id}`,
    validId,
  );
  const t = tournament.data;
  const events = useCourseResource<{
    events: AuditEvent[];
    total: number;
  }>(`/api/tournaments/${id}/events?limit=50&offset=0`, !!t);
  const profiles = useAuthenticatedResource<{
    agents: AgentProfile[];
    total: number;
  }>("/api/agents/my", token);
  const match = useCourseResource<MatchResult>(
    `/api/tournaments/${id}/match`,
    !!t && t.status !== "open",
  );
  const evidence = useCourseResource<Evidence>(
    `/api/tournaments/${id}/evidence`,
    match.data?.status === "completed",
  );
  const notes = useCourseResource<{
    notes: EvidenceNote[];
    total: number;
  }>(`/api/tournaments/${id}/notes`, !!t);
  const queryClient = useQueryClient();
  const tx = useWalletTransaction();
  const [amount, setAmount] = useState("100");
  const [agentHash, setAgentHash] = useState("");
  const [reason, setReason] = useState("");
  const [decisionWinner, setDecisionWinner] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const [viewedEvidence, setViewedEvidence] = useState("");
  const [showEvidence, setShowEvidence] = useState(false);
  const [panel, setPanel] = useState<"join" | "fund">("join");
  useEffect(() => {
    setViewedEvidence("");
    setShowEvidence(false);
    setAgentHash("");
    setReason("");
    setDecisionWinner("");
    setError("");
  }, [address, id]);
  const block = useBlock({
    chainId: config?.chain_id,
    watch: true,
    query: { enabled: !!config?.deployment_configured, refetchInterval: 5000 },
  });
  const now = block.data ? Number(block.data.timestamp) : 0;
  const contractArgs =
    address && validId ? ([BigInt(id), address] as const) : undefined;
  const readEnabled = !!config?.deployment_configured && !!address && validId;
  const contribution = useReadContract({
    address: config?.manager_address,
    abi: tournamentManagerAbi,
    functionName: "contributionOf",
    args: contractArgs,
    chainId: config?.chain_id,
    query: { enabled: readEnabled, refetchInterval: 10000 },
  });
  const refund = useReadContract({
    address: config?.manager_address,
    abi: tournamentManagerAbi,
    functionName: "refundableAmount",
    args: contractArgs,
    chainId: config?.chain_id,
    query: { enabled: readEnabled, refetchInterval: 10000 },
  });
  const prize = useReadContract({
    address: config?.manager_address,
    abi: tournamentManagerAbi,
    functionName: "claimablePrize",
    args: contractArgs,
    chainId: config?.chain_id,
    query: { enabled: readEnabled, refetchInterval: 10000 },
  });
  const allowance = useReadContract({
    address: config?.token_address,
    abi: courseTokenAbi,
    functionName: "allowance",
    args: address && config ? [address, config.escrow_address] : undefined,
    chainId: config?.chain_id,
    query: { enabled: readEnabled, refetchInterval: 10000 },
  });
  const balance = useReadContract({
    address: config?.token_address,
    abi: courseTokenAbi,
    functionName: "balanceOf",
    args: address ? [address] : undefined,
    chainId: config?.chain_id,
    query: { enabled: readEnabled, refetchInterval: 10000 },
  });
  const organizer =
    !!user?.roles.includes("organizer") && sameAddress(address, t?.organizer);
  const auditor =
    !!user?.roles.includes("auditor") && sameAddress(address, t?.auditor);
  const participant =
    t?.participants.some((a) => sameAddress(a, address)) || false;
  const open = !!t && now > 0 && registrationOpen(t, now);
  const reviewComplete = !!t && now > 0 && isReviewComplete(t, now);
  const expired = !!t && now > 0 && canExpire(t, now);
  const verified =
    !!evidence.data &&
    verifyEvidence(evidence.data.canonical_json, evidence.data.evidence_hash) &&
    evidence.data.evidence_hash === match.data?.evidence_hash;
  const evidenceRead =
    !!evidence.data && viewedEvidence === evidence.data.evidence_hash;
  const submittedEvidenceMatches =
    !!evidence.data && t?.evidence_hash === evidence.data.evidence_hash;
  const disabled = !tx.canWrite || tx.locked || busy;
  async function idAction(name: IdAction, label: string) {
    if (!config) return;
    await tx.prepare(
      `${label} · 赛事 #${id}`,
      config.manager_address,
      encodeFunctionData({
        abi: tournamentManagerAbi,
        functionName: name,
        args: [BigInt(id)],
      }),
    );
  }
  async function fund(approve: boolean) {
    if (!config) return;
    setError("");
    try {
      const raw = parseMoney(amount, config.token_decimals);
      if (raw <= BigInt(0)) throw new Error("赞助金额应大于零。");
      if (sameAddress(address, t?.auditor))
        throw new Error("指定审计员不能赞助自己审核的赛事。");
      if (!approve && (allowance.data === undefined || allowance.data < raw))
        throw new Error("请先授权赞助金额，等待授权确认，再提交赞助。");
      if (!approve && (balance.data === undefined || balance.data < raw))
        throw new Error("DUSD 余额不足。");
      await tx.prepare(
        `${approve ? "授权赞助金额" : "赞助奖金"} ${money(raw.toString(), config.token_decimals)} DUSD · 赛事 #${id}`,
        approve ? config.token_address : config.manager_address,
        approve
          ? encodeFunctionData({
              abi: courseTokenAbi,
              functionName: "approve",
              args: [config.escrow_address, raw],
            })
          : encodeFunctionData({
              abi: tournamentManagerAbi,
              functionName: "fundTournament",
              args: [BigInt(id), raw],
            }),
      );
    } catch (err) {
      setError(errorMessage(err));
    }
  }
  async function register() {
    if (!config) return;
    setError("");
    const selected = profiles.data?.agents.find(
      (a) => a.agent_hash === agentHash,
    );
    if (!selected || !sameAddress(selected.owner_address, address)) {
      setError("请选择自己钱包创建的 Agent 定义。");
      return;
    }
    await tx.prepare(
      `报名赛事 #${id} · ${selected.name} · ${shortAddress(selected.agent_hash)}`,
      config.manager_address,
      encodeFunctionData({
        abi: tournamentManagerAbi,
        functionName: "registerParticipant",
        args: [BigInt(id), selected.agent_hash],
      }),
    );
  }
  async function run() {
    if (!token) return;
    setBusy(true);
    setError("");
    try {
      await request<MatchResult>(`/api/tournaments/${id}/run`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
      });
      await queryClient.invalidateQueries({
        queryKey: ["course", `/api/tournaments/${id}/match`],
      });
    } catch (err) {
      setError(errorMessage(err));
    } finally {
      setBusy(false);
    }
  }
  async function submit() {
    if (
      !config ||
      !verified ||
      !evidenceRead ||
      !match.data?.winner_address ||
      !evidence.data
    )
      return;
    await tx.prepare(
      `提交赛事 #${id} 结果 · 冠军 ${shortAddress(match.data.winner_address)} · 证据 ${shortAddress(evidence.data.evidence_hash)}`,
      config.manager_address,
      encodeFunctionData({
        abi: tournamentManagerAbi,
        functionName: "submitResult",
        args: [
          BigInt(id),
          match.data.winner_address,
          evidence.data.evidence_hash,
        ],
      }),
    );
  }
  async function noteAction(kind: "dispute" | "decision" | "cancellation") {
    if (!config || !token) return;
    setBusy(true);
    setError("");
    try {
      if (!reason.trim() || reason.length > 4000)
        throw new Error("说明需要 1–4000 字符。");
      if (kind === "decision" && !decisionWinner)
        throw new Error("请先选择裁决结果。");
      const note = await request<Evidence>(`/api/tournaments/${id}/notes`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
        body: JSON.stringify({ kind, text: reason.trim() }),
      });
      if (!verifyEvidence(note.canonical_json, note.evidence_hash))
        throw new Error("说明内容与哈希不一致，无法提交。");
      await queryClient.invalidateQueries({
        queryKey: ["course", `/api/tournaments/${id}/notes`],
      });
      const data =
        kind === "decision"
          ? encodeFunctionData({
              abi: tournamentManagerAbi,
              functionName: "resolveDispute",
              args: [BigInt(id), decisionWinner as Address, note.evidence_hash],
            })
          : encodeFunctionData({
              abi: tournamentManagerAbi,
              functionName:
                kind === "dispute" ? "disputeResult" : "cancelTournament",
              args: [BigInt(id), note.evidence_hash],
            });
      await tx.prepare(
        `${kind === "dispute" ? "提交异议" : kind === "decision" ? `裁决异议：${decisionWinner === zeroAddress ? "取消并退款" : shortAddress(decisionWinner)}` : "取消赛事"} · #${id} · 说明 ${shortAddress(note.evidence_hash)}`,
        config.manager_address,
        data,
      );
    } catch (err) {
      setError(errorMessage(err));
    } finally {
      setBusy(false);
    }
  }
  function downloadEvidence() {
    if (!evidence.data) return;
    const url = URL.createObjectURL(
      new Blob([evidence.data.canonical_json], { type: "application/json" }),
    );
    const a = document.createElement("a");
    a.href = url;
    a.download = `tournament-${id}-evidence.json`;
    a.click();
    URL.revokeObjectURL(url);
    setViewedEvidence(evidence.data?.evidence_hash || "");
  }
  function inspectEvidence() {
    setShowEvidence(true);
    if (evidence.data) setViewedEvidence(evidence.data.evidence_hash);
    requestAnimationFrame(() =>
      document
        .getElementById("match-evidence")
        ?.scrollIntoView({ behavior: "smooth", block: "start" }),
    );
  }
  if (!validId)
    return (
      <div role="alert" className="rounded-2xl border bg-card p-8">
        赛事编号无效。
        <Link href="/tournaments" className="ml-3 text-primary underline">
          返回赛事大厅
        </Link>
      </div>
    );
  if (!t)
    return (
      <div className="space-y-5">
        <Link
          href="/tournaments"
          className="inline-flex items-center gap-2 text-sm text-muted-foreground"
        >
          <ArrowLeft size={16} />
          赛事大厅
        </Link>
        <DeploymentNotice />
        <div
          role={tournament.error ? "alert" : "status"}
          className="rounded-2xl border bg-card p-8"
        >
          {tournament.error
            ? errorMessage(tournament.error)
            : "正在加载赛事与奖金信息…"}
          {tournament.error && (
            <Button
              className="ml-4"
              variant="outline"
              onClick={() => void tournament.refetch()}
            >
              重新加载
            </Button>
          )}
        </div>
      </div>
    );

  const ownerWallet = sameAddress(address, t.organizer);
  const auditWallet = sameAddress(address, t.auditor);
  const held = contribution.data;
  const hasSponsorshipEvent = events.data?.events.some(
    (e) =>
      ["TournamentFunded", "SponsorRefunded"].includes(e.event_name) &&
      sameAddress(e.actor, address),
  );
  const sponsor =
    (held !== undefined && held > BigInt(0)) || !!hasSponsorshipEvent;
  const winner = sameAddress(address, t.winner);
  const financialError =
    contribution.error || refund.error || prize.error || balance.error;
  const walletReady =
    held !== undefined &&
    refund.data !== undefined &&
    prize.data !== undefined &&
    balance.data !== undefined;
  const identityLabels = [
    ownerWallet && "本场主办方",
    auditWallet && "本场审核员",
    sponsor && "赞助方",
    participant && (winner ? "冠军 · 参赛者" : "参赛者"),
  ].filter(Boolean) as string[];
  const presentation = getTournamentPresentation(t, now);
  const pool = BigInt(t.prize_pool);
  const poolRemaining = pool - BigInt(t.prize_paid);
  const notFinal = !["finalized", "cancelled"].includes(t.status);
  const reviewState = ["result_submitted", "disputed"].includes(t.status);
  const canSubmitInTime =
    now > 0 && now + t.review_period < t.settlement_deadline;
  const canDispute =
    t.status === "result_submitted" &&
    now > 0 &&
    now < t.review_deadline &&
    (participant || (held || BigInt(0)) > BigInt(0));
  const canJoin =
    open &&
    !participant &&
    !ownerWallet &&
    !auditWallet &&
    t.participants.length < 2;
  let requested = BigInt(0);
  let amountError = "";
  try {
    requested = parseMoney(amount, config?.token_decimals);
    if (requested <= BigInt(0)) amountError = "请输入大于零的赞助金额。";
  } catch (err) {
    amountError = errorMessage(err);
  }
  const enoughAllowance =
    requested > BigInt(0) &&
    allowance.data !== undefined &&
    allowance.data >= requested;
  const step =
    t.status === "finalized"
      ? t.prize_claimed
        ? 5
        : 4
      : reviewState
        ? 3
        : t.status === "running"
          ? 2
          : 1;
  const amountLabel = (value: bigint | undefined) =>
    value === undefined ? "—" : money(value.toString(), config?.token_decimals);
  const verifyButton = (
    <Button variant="outline" onClick={inspectEvidence}>
      <FileCheck2 size={16} className="mr-2" />
      查看比赛证据
    </Button>
  );
  let nextTitle = presentation.label;
  let nextDescription = presentation.hint;
  let nextContent: ReactNode = null;
  if (!address) {
    nextTitle = "连接钱包，查看你的参赛与资金状态";
    nextDescription = expired
      ? "这场赛事已到期。赞助者连接原出资钱包后，可先取消赛事，再领取退款。"
      : "回放和资金记录可以直接浏览；连接钱包后，才能报名、赞助或处理赛事。";
    nextContent = <WalletConnect />;
  } else if (financialError) {
    nextTitle = "暂时无法读取你的资金状态";
    nextDescription = "请检查钱包网络，重新读取后再操作。";
    nextContent = (
      <Button
        variant="outline"
        onClick={() => {
          void contribution.refetch();
          void refund.refetch();
          void prize.refetch();
          void balance.refetch();
        }}
      >
        重新读取资金
      </Button>
    );
  } else if (!walletReady || (now === 0 && notFinal)) {
    nextTitle = "正在核对你的身份与资金";
    nextDescription = "读取链上记录后，会在这里显示你当前能做的操作。";
  } else if (expired) {
    nextTitle = sponsor
      ? "赛事已到期，取回你的赞助"
      : "赛事已到期，可以取消并开启退款";
    nextDescription = sponsor
      ? `你的 ${amountLabel(held)} DUSD 仍在托管。先确认取消，之后再将退款领回这个钱包。`
      : "任何钱包都可以触发到期取消；取消确认后，每位赞助者自行领取退款。";
    nextContent = (
      <div className="space-y-4">
        <div className="flex flex-wrap items-center gap-3 text-sm">
          <span className="font-semibold">① 确认取消赛事</span>
          <ArrowRight size={16} />
          <span className="text-muted-foreground">② 赞助者领取退款</span>
        </div>
        <Button
          size="lg"
          disabled={disabled}
          onClick={() => idAction("expireTournament", "取消到期赛事，开启退款")}
        >
          <ArrowDownToLine size={18} className="mr-2" />
          取消赛事，开启退款
        </Button>
      </div>
    );
  } else if ((refund.data || BigInt(0)) > BigInt(0)) {
    nextTitle = `${amountLabel(refund.data)} DUSD 退款待领取`;
    nextDescription = "赛事已经取消。这笔退款属于当前钱包，确认后将直接退回。";
    nextContent = (
      <Button
        size="lg"
        disabled={disabled}
        onClick={() =>
          idAction("claimRefund", `领取 ${amountLabel(refund.data)} DUSD 退款`)
        }
      >
        <ArrowDownToLine size={18} className="mr-2" />
        领取 {amountLabel(refund.data)} DUSD 退款
      </Button>
    );
  } else if ((prize.data || BigInt(0)) > BigInt(0)) {
    nextTitle = `你赢得了 ${amountLabel(prize.data)} DUSD`;
    nextDescription = "审核已通过。领取后，奖金将直接进入当前钱包。";
    nextContent = (
      <Button
        size="lg"
        disabled={disabled}
        onClick={() =>
          idAction("claimPrize", `领取 ${amountLabel(prize.data)} DUSD 奖金`)
        }
      >
        <Trophy size={18} className="mr-2" />
        领取 {amountLabel(prize.data)} DUSD 奖金
      </Button>
    );
  } else if (t.status === "cancelled") {
    nextTitle = "赛事已取消";
    nextDescription =
      "当前钱包没有可领取的退款。你可以在下方核对出资与退款记录。";
  } else if (t.status === "finalized") {
    nextTitle =
      winner && t.prize_claimed
        ? "奖金已领取，恭喜获胜"
        : t.prize_claimed
          ? "比赛已完成，奖金已发放"
          : "审核已完成，等待冠军领取奖金";
    nextDescription = "查看比赛回放与资金记录，了解本场赛事的完整过程。";
  } else if ((ownerWallet || auditWallet) && !token) {
    nextTitle = `${ownerWallet ? "主办方" : "审核员"}，请先签名登录`;
    nextDescription =
      "钱包已连接。登录验证身份后，就能处理这场赛事；登录签名不扣款。";
    nextContent = <LoginButton />;
  } else if ((ownerWallet && !organizer) || (auditWallet && !auditor)) {
    nextTitle = "当前账户缺少管理权限";
    nextDescription =
      "该钱包是赛事指定账户，但当前登录权限未通过校验。请重新登录或联系管理员；到期退款仍可按规则处理。";
  } else if (organizer && t.status === "open") {
    const ready =
      now > 0 &&
      now >= t.registration_deadline &&
      t.participants.length === 2 &&
      pool > BigInt(0) &&
      canSubmitInTime;
    nextTitle = !canSubmitInTime
      ? "剩余时间不足以开赛并审核"
      : ready
        ? "报名已截止，可以开始比赛"
        : "准备好奖池和参赛者，等待开赛";
    nextDescription = !canSubmitInTime
      ? "你可以在下方取消赛事，让赞助者领取退款。最终期限到达后，任意钱包也可触发取消。"
      : `报名截止：${dateTime(t.registration_deadline)}。两位参赛者和非零奖池齐备后，由你确认开赛。`;
    nextContent = (
      <div className="space-y-4">
        <div className="flex flex-wrap gap-4 text-sm">
          {[
            [pool > BigInt(0), "奖池已落实"],
            [t.participants.length === 2, "两位参赛者已报名"],
            [now >= t.registration_deadline && now > 0, "报名已截止"],
          ].map(([ok, label]) => (
            <span
              key={String(label)}
              className={ok ? "text-emerald-700" : "text-muted-foreground"}
            >
              {ok ? "✓" : "○"} {label}
            </span>
          ))}
        </div>
        <Button
          size="lg"
          disabled={disabled || !ready}
          onClick={() => idAction("startTournament", "开始赛事")}
        >
          <Play size={18} className="mr-2" />
          开始比赛
        </Button>
      </div>
    );
  } else if (organizer && t.status === "running") {
    if (match.data?.status === "completed") {
      nextTitle = match.data.winner_address
        ? "对战已完成，请查看证据并提交结果"
        : "本场平局，等待到期退款";
      nextDescription = !match.data.winner_address
        ? "当前规则不为平局指定冠军。最终期限到达后，可取消赛事并退回赞助。"
        : !canSubmitInTime
          ? "剩余时间已不足一个完整审核窗口。最终期限到达后，请按规则取消退款。"
          : "先查看下方回放和完整证据，再将结果交给独立审核员。";
      nextContent = match.data.winner_address && (
        <div className="flex flex-wrap gap-3">
          {verifyButton}
          <Button
            size="lg"
            disabled={
              disabled || !verified || !evidenceRead || !canSubmitInTime
            }
            onClick={submit}
          >
            提交比赛结果
            <ArrowRight size={17} className="ml-2" />
          </Button>
        </div>
      );
    } else {
      const running = ["queued", "running"].includes(match.data?.status || "");
      const failed = match.data?.status === "failed";
      nextTitle = running
        ? "Agent 正在对战"
        : failed
          ? "对战执行失败"
          : match.data?.status === "interrupted"
            ? "对战中断，可以继续执行"
            : "比赛已开赛，现在运行 Agent 对战";
      nextDescription = running
        ? "回合记录会自动更新。你可以在下方观察对战进度。"
        : failed
          ? "请查看下方错误记录并联系维护者；最终期限到达后可走退款流程。"
          : "系统执行对战并保存回合记录。完成后，你还需要提交结果供审核。";
      nextContent = (
        <Button
          size="lg"
          disabled={!token || busy || running || failed}
          onClick={run}
        >
          <Play size={18} className="mr-2" />
          {running
            ? "对战进行中…"
            : match.data?.status === "interrupted"
              ? "恢复对战执行"
              : "运行 Agent 对战"}
        </Button>
      );
    }
  } else if (auditor && reviewState) {
    nextTitle =
      t.status === "disputed"
        ? "这场比赛有异议，等待你的裁决"
        : "比赛结果等待你审核";
    nextDescription = `先查看证据${t.status === "disputed" ? "和异议说明" : ""}。审核窗口于 ${dateTime(t.review_deadline)} 结束，之后可确认${t.status === "disputed" ? "裁决" : "结果"}。`;
    nextContent = (
      <div className="flex flex-wrap gap-3">
        {verifyButton}
        {t.status === "result_submitted" ? (
          <Button
            size="lg"
            disabled={
              disabled ||
              !reviewComplete ||
              !verified ||
              !submittedEvidenceMatches ||
              !evidenceRead
            }
            onClick={() => idAction("approveResult", "审核通过并结算")}
          >
            审核通过，确认冠军
          </Button>
        ) : (
          <a
            href="#review-action"
            className="inline-flex items-center rounded-lg bg-primary px-5 py-3 text-sm font-semibold text-primary-foreground"
          >
            查看异议并裁决
            <ArrowRight size={16} className="ml-2" />
          </a>
        )}
      </div>
    );
  } else if (open) {
    nextTitle = participant
      ? "你已报名，等待比赛开始"
      : sponsor
        ? "你的赞助已进入奖池"
        : "报名参赛，或为这场比赛赞助奖金";
    nextDescription = participant
      ? `报名截止后由主办方开赛。你的 Agent 已绑定本次报名。`
      : sponsor
        ? "你可以继续赞助，并在比赛后查看结果；赛事取消时可按出资退款。"
        : "参赛需要一个自己的 Agent；赞助金额会先进入奖金托管合约。";
    nextContent = (
      <Button
        variant="outline"
        onClick={() => {
          setPanel(!participant && sponsor ? "fund" : "join");
          document
            .getElementById("participate")
            ?.scrollIntoView({ behavior: "smooth", block: "start" });
        }}
      >
        {participant
          ? "查看我的报名"
          : sponsor
            ? "查看我的赞助"
            : "选择参与方式"}
        <ArrowRight size={16} className="ml-2" />
      </Button>
    );
  }

  return (
    <div className="space-y-6 pb-8">
      <Link
        href="/tournaments"
        className="inline-flex items-center gap-2 text-sm font-medium text-muted-foreground hover:text-primary"
      >
        <ArrowLeft size={16} />
        返回赛事大厅
      </Link>
      <DeploymentNotice />
      <header className="flex flex-wrap items-start justify-between gap-5">
        <div className="min-w-0">
          <div className="mb-3 flex items-center gap-3">
            <span className="text-xs font-semibold tracking-widest text-muted-foreground">
              赛事 #{t.id}
            </span>
            <TournamentStatusBadge tournament={t} now={now} />
          </div>
          <h1 className="break-words text-2xl font-bold tracking-tight sm:text-3xl">
            {t.title}
          </h1>
          <p className="mt-2 text-sm text-muted-foreground">
            两位 Agent 参赛 · 赞助奖金预存 · 独立审核后领奖
          </p>
        </div>
        <div className="rounded-2xl border bg-card px-6 py-4">
          <p className="text-sm text-muted-foreground">本场奖金总额</p>
          <p className="mt-1 text-3xl font-bold tabular-nums">
            {money(t.prize_pool, config?.token_decimals)}{" "}
            <span className="text-sm font-medium text-muted-foreground">
              DUSD
            </span>
          </p>
        </div>
      </header>

      <section
        aria-label="赛事进度"
        className="rounded-2xl border bg-card p-4 sm:p-5"
      >
        {t.status === "cancelled" || expired ? (
          <div className="flex flex-wrap items-center gap-3 text-sm">
            <span className="rounded-full bg-amber-100 p-2 text-amber-800">
              <ArrowDownToLine size={18} />
            </span>
            <strong>
              {t.status === "cancelled" ? "赛事已取消" : "赛事已到期"}
            </strong>
            <ArrowRight size={16} className="text-muted-foreground" />
            <span>
              {t.status === "cancelled"
                ? "赞助者可领取尚未退回的资金"
                : "确认取消后，赞助者可领取退款"}
            </span>
          </div>
        ) : (
          <ol className="grid grid-cols-2 gap-3 sm:grid-cols-5">
            {["报名与赞助", "Agent 对战", "结果审核", "冠军领奖", "完成"].map(
              (label, i) => (
                <li
                  key={label}
                  aria-current={step === i + 1 ? "step" : undefined}
                  className={`flex items-center gap-2 rounded-lg px-3 py-2 text-sm ${step === i + 1 ? "bg-emerald-50 font-semibold text-emerald-800" : step > i + 1 ? "text-emerald-700" : "text-muted-foreground"}`}
                >
                  <span
                    className={`flex h-6 w-6 shrink-0 items-center justify-center rounded-full text-xs ${step >= i + 1 ? "bg-emerald-600 text-white" : "bg-muted"}`}
                  >
                    {step > i + 1 ? <Check size={14} /> : i + 1}
                  </span>
                  {label}
                </li>
              ),
            )}
          </ol>
        )}
      </section>

      <div className="grid items-start gap-6 lg:grid-cols-[minmax(0,1fr)_320px]">
        <div className="min-w-0 space-y-6">
          <section
            aria-labelledby="next-action-title"
            className={`rounded-2xl border p-6 sm:p-7 ${expired || (refund.data || BigInt(0)) > BigInt(0) ? "border-amber-200 bg-amber-50/70" : "border-emerald-200 bg-emerald-50/70"}`}
          >
            <div className="mb-4 flex flex-wrap items-center gap-2">
              <span className="text-xs font-bold tracking-widest text-muted-foreground">
                你的下一步
              </span>
              {address && (
                <span className="rounded-full border bg-white px-3 py-1 text-xs font-semibold">
                  {identityLabels.length
                    ? identityLabels.join(" / ")
                    : walletReady
                      ? "观赛者"
                      : "核对身份中"}
                </span>
              )}
            </div>
            <h2
              id="next-action-title"
              className="text-xl font-bold sm:text-2xl"
            >
              {nextTitle}
            </h2>
            <p className="mt-3 max-w-2xl text-sm leading-7 text-slate-600">
              {nextDescription}
            </p>
            {nextContent && <div className="mt-5">{nextContent}</div>}
            {address && !tx.canWrite && (
              <p className="mt-4 text-sm font-medium text-amber-800">
                钱包网络未就绪，请先处理上方的网络提示。
              </p>
            )}
            {error && (
              <p
                role="alert"
                className="mt-4 rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-800"
              >
                {error}
              </p>
            )}
          </section>

          {open && (
            <section
              id="participate"
              className="scroll-mt-28 rounded-2xl border bg-card p-6"
            >
              <div className="mb-5 flex flex-wrap items-center justify-between gap-3">
                <h2 className="text-lg font-bold">参与这场赛事</h2>
                <span className="text-sm text-muted-foreground">
                  {t.participants.length} / 2 人已报名
                </span>
              </div>
              {!address ? (
                <div className="space-y-4">
                  <p className="text-sm text-muted-foreground">
                    连接钱包后，可以选择自己的 Agent 报名，或赞助奖金。
                  </p>
                  <WalletConnect />
                </div>
              ) : (
                <>
                  <div
                    className="mb-5 flex gap-2 rounded-xl bg-muted p-1"
                    role="group"
                    aria-label="参与方式"
                  >
                    <Button
                      className="flex-1"
                      variant={panel === "join" ? "default" : "ghost"}
                      aria-pressed={panel === "join"}
                      onClick={() => setPanel("join")}
                    >
                      <Users size={16} className="mr-2" />
                      我要参赛
                    </Button>
                    <Button
                      className="flex-1"
                      variant={panel === "fund" ? "default" : "ghost"}
                      aria-pressed={panel === "fund"}
                      onClick={() => setPanel("fund")}
                    >
                      <Coins size={16} className="mr-2" />
                      赞助奖金
                    </Button>
                  </div>
                  {panel === "join" ? (
                    participant ? (
                      <div className="space-y-4">
                        <p className="flex items-center gap-2 font-medium text-emerald-700">
                          <CheckCircle2 size={19} />
                          你已成功报名
                        </p>
                        <p className="text-sm text-muted-foreground">
                          报名截止前可以退出；开赛后由已报名的 Agent 参加对战。
                        </p>
                        <Button
                          variant="outline"
                          disabled={disabled}
                          onClick={() =>
                            idAction("withdrawParticipant", "退出报名")
                          }
                        >
                          退出本场报名
                        </Button>
                      </div>
                    ) : !canJoin ? (
                      <p className="rounded-xl bg-muted p-4 text-sm">
                        {ownerWallet || auditWallet
                          ? "本场主办方和审核员保持独立，不能参赛。"
                          : "本场两个参赛名额已满，你仍可以观赛或赞助。"}
                      </p>
                    ) : !token ? (
                      <div className="space-y-3">
                        <p className="text-sm text-muted-foreground">
                          签名登录后，即可选择自己的 Agent 报名。报名不收费用。
                        </p>
                        <LoginButton />
                      </div>
                    ) : (
                      <div className="space-y-4">
                        {profiles.isLoading ? (
                          <p role="status">正在读取你的 Agent…</p>
                        ) : profiles.error ? (
                          <p role="alert" className="text-sm text-red-700">
                            {errorMessage(profiles.error)}
                          </p>
                        ) : !profiles.data?.agents.length ? (
                          <div className="rounded-xl border border-dashed p-5">
                            <p className="mb-3 text-sm">
                              你还没有参赛 Agent，先为它取名并设定策略。
                            </p>
                            <Link
                              href="/agents"
                              className="inline-flex items-center font-semibold text-primary"
                            >
                              创建我的 Agent
                              <ArrowRight size={16} className="ml-2" />
                            </Link>
                          </div>
                        ) : (
                          <>
                            <label className="block text-sm font-medium">
                              选择报名 Agent
                              <select
                                aria-label="选择报名 Agent"
                                className="mt-2 min-h-11 w-full rounded-lg border bg-background p-3"
                                value={agentHash}
                                onChange={(e) => setAgentHash(e.target.value)}
                              >
                                <option value="">请选择你的 Agent</option>
                                {profiles.data.agents.map((a) => (
                                  <option key={a.id} value={a.agent_hash}>
                                    {a.name}
                                  </option>
                                ))}
                              </select>
                            </label>
                            <Button
                              size="lg"
                              disabled={disabled || !agentHash}
                              onClick={register}
                            >
                              确认报名，查看费用
                            </Button>
                            <p className="text-xs text-muted-foreground">
                              报名费 0 DUSD，仅需本地网络费。
                              <Link href="/agents" className="ml-2 underline">
                                管理我的 Agent
                              </Link>
                            </p>
                          </>
                        )}
                      </div>
                    )
                  ) : auditWallet ? (
                    <p className="rounded-xl bg-muted p-4 text-sm">
                      你是本场审核员，需要保持独立，不能赞助这场赛事。
                    </p>
                  ) : (
                    <div className="space-y-4">
                      <label className="block text-sm font-medium">
                        赞助金额（DUSD）
                        <Input
                          className="mt-2 h-12 text-lg"
                          inputMode="decimal"
                          value={amount}
                          aria-invalid={!!amountError}
                          aria-describedby={
                            amountError ? "sponsorship-amount-error" : undefined
                          }
                          onChange={(e) => setAmount(e.target.value)}
                          placeholder="例如 100"
                        />
                      </label>
                      {amountError && (
                        <p
                          id="sponsorship-amount-error"
                          className="text-sm text-red-700"
                        >
                          {amountError}
                        </p>
                      )}
                      <div className="flex gap-2">
                        {[25, 100, 200].map((value) => (
                          <Button
                            key={value}
                            size="sm"
                            variant="outline"
                            onClick={() => setAmount(String(value))}
                          >
                            {value} DUSD
                          </Button>
                        ))}
                      </div>
                      <p className="text-sm text-muted-foreground">
                        钱包可用 {amountLabel(balance.data)} DUSD · 本场仍托管{" "}
                        {amountLabel(held)} DUSD
                      </p>
                      <div className="rounded-xl bg-muted p-4">
                        <p className="text-sm font-medium">
                          {enoughAllowance
                            ? "✓ 金额已授权，可以确认赞助"
                            : "先授权金额，再确认赞助"}
                        </p>
                        <p className="mt-1 text-xs leading-5 text-muted-foreground">
                          授权不扣款。确认赞助后，资金进入奖池；赛事取消时按原出资退款。
                        </p>
                      </div>
                      <div className="flex flex-wrap gap-3">
                        <Button
                          variant={enoughAllowance ? "outline" : "default"}
                          disabled={disabled || requested <= BigInt(0)}
                          onClick={() => fund(true)}
                        >
                          1. 授权{" "}
                          {requested > BigInt(0)
                            ? money(
                                requested.toString(),
                                config?.token_decimals,
                              )
                            : ""}{" "}
                          DUSD
                        </Button>
                        <Button
                          disabled={
                            disabled ||
                            !enoughAllowance ||
                            balance.data === undefined ||
                            balance.data < requested
                          }
                          onClick={() => fund(false)}
                        >
                          2. 确认赞助
                        </Button>
                      </div>
                      {requested > BigInt(0) &&
                        balance.data !== undefined &&
                        balance.data < requested && (
                          <p className="text-sm text-red-700">
                            DUSD 余额不足，请换到有测试币的赞助钱包。
                          </p>
                        )}
                    </div>
                  )}
                </>
              )}
            </section>
          )}

          <section className="rounded-2xl border bg-card p-6">
            <div className="mb-5 flex items-center justify-between">
              <h2 className="text-lg font-bold">参赛阵容</h2>
              <span className="text-sm text-muted-foreground">
                {t.participants.length} / 2 席
              </span>
            </div>
            <div className="grid gap-4 sm:grid-cols-2">
              {[0, 1].map((i) => {
                const p = t.participant_agents[i];
                const info =
                  i === 0 ? match.data?.red_agent : match.data?.blue_agent;
                return (
                  <div
                    key={i}
                    className={`min-w-0 rounded-xl border p-4 ${i === 0 ? "border-rose-100 bg-rose-50/40" : "border-blue-100 bg-blue-50/40"}`}
                  >
                    <p
                      className={`mb-2 text-xs font-bold ${i === 0 ? "text-rose-700" : "text-blue-700"}`}
                    >
                      {i === 0 ? "红方" : "蓝方"}
                    </p>
                    <p className="break-words font-semibold">
                      {p ? info?.name || "已报名 Agent" : "等待参赛者"}
                    </p>
                    {p ? (
                      <>
                        <p
                          className="mt-2 font-mono text-xs text-muted-foreground"
                          title={p.address}
                        >
                          {shortAddress(p.address)}
                          {sameAddress(p.address, address) && " · 你"}
                        </p>
                        <details className="mt-3 text-xs text-muted-foreground">
                          <summary className="cursor-pointer">报名凭证</summary>
                          <p className="mt-2 break-all font-mono">
                            {p.agent_hash}
                          </p>
                        </details>
                      </>
                    ) : (
                      <p className="mt-2 text-sm text-muted-foreground">
                        创建自己的 Agent 后即可报名
                      </p>
                    )}
                  </div>
                );
              })}
            </div>
          </section>

          {match.data ? (
            <MatchReplay match={match.data} settlementStatus={t.status} />
          ) : (
            <section className="rounded-2xl border border-dashed bg-card p-6 text-center">
              <Play size={26} className="mx-auto mb-3 text-muted-foreground" />
              <h2 className="font-semibold">
                {t.status === "cancelled" || expired
                  ? "这场赛事尚未生成对战回放"
                  : "对战开始后，在这里观看回放"}
              </h2>
              <p className="mt-2 text-sm text-muted-foreground">
                地图、血量和每回合动作会随对战记录显示。
              </p>
            </section>
          )}
          {match.error &&
            !(
              match.error instanceof ApiError && match.error.status === 404
            ) && (
              <p role="alert" className="text-sm text-red-700">
                对战数据：{errorMessage(match.error)}
              </p>
            )}

          {evidence.data && (
            <section
              id="match-evidence"
              className="scroll-mt-28 rounded-2xl border bg-card p-6"
            >
              <div className="flex flex-wrap items-center justify-between gap-3">
                <h2 className="flex items-center gap-2 text-lg font-bold">
                  <FileCheck2 size={20} />
                  比赛证据
                </h2>
                <span
                  className={`rounded-full px-3 py-1 text-xs font-semibold ${verified && (!t.evidence_hash || submittedEvidenceMatches) ? "bg-emerald-50 text-emerald-700" : "bg-red-50 text-red-700"}`}
                >
                  {verified
                    ? t.evidence_hash
                      ? submittedEvidenceMatches
                        ? "内容与链上记录一致"
                        : "链上记录不一致"
                      : "内容校验通过"
                    : "内容校验失败"}
                </span>
              </div>
              <p className="mt-3 text-sm leading-6 text-muted-foreground">
                包含双方
                Agent、每回合动作和结果。主办方提交与审核员批准前，都需要先查看。
              </p>
              <div className="mt-4 flex flex-wrap gap-3">
                <Button
                  variant="outline"
                  onClick={() => {
                    setShowEvidence(!showEvidence);
                    if (!showEvidence)
                      setViewedEvidence(evidence.data!.evidence_hash);
                  }}
                >
                  {showEvidence ? "收起完整证据" : "展开完整证据"}
                </Button>
                <Button variant="outline" onClick={downloadEvidence}>
                  <ArrowDownToLine size={16} className="mr-2" />
                  下载证据 JSON
                </Button>
                {evidenceRead && (
                  <span className="flex items-center gap-1 text-sm text-emerald-700">
                    <Check size={15} />
                    已查看
                  </span>
                )}
              </div>
              {showEvidence && (
                <pre className="mt-4 max-h-80 overflow-auto whitespace-pre-wrap break-all rounded-xl bg-muted p-4 text-xs">
                  {evidence.data.canonical_json}
                </pre>
              )}
              <details className="mt-4 text-xs text-muted-foreground">
                <summary className="cursor-pointer">证据指纹与核验说明</summary>
                <p className="mt-2 break-all font-mono">
                  {evidence.data.evidence_hash}
                </p>
                <p className="mt-2 leading-5">
                  哈希核对用于发现内容变更。审核员仍需判断比赛结果和规则是否正确。
                </p>
              </details>
            </section>
          )}
          {evidence.error && (
            <p role="alert" className="text-sm text-red-700">
              证据加载失败：{errorMessage(evidence.error)}
            </p>
          )}

          {(canDispute ||
            (auditor && t.status === "disputed" && !expired) ||
            (organizer && t.status === "open" && !expired)) && (
            <section
              id="review-action"
              className="scroll-mt-28 rounded-2xl border bg-card p-6"
            >
              <details
                open={auditor && t.status === "disputed" ? true : undefined}
              >
                <summary className="cursor-pointer font-semibold">
                  {auditor && t.status === "disputed"
                    ? "审核异议并裁决"
                    : canDispute
                      ? "对结果有异议？提交复核申请"
                      : "需要取消这场赛事？"}
                </summary>
                <div className="mt-4 space-y-4">
                  <p className="text-sm leading-6 text-muted-foreground">
                    {canDispute
                      ? `在 ${dateTime(t.review_deadline)} 前说明问题并提交，审核员会处理你的异议。`
                      : t.status === "disputed"
                        ? "结合回放、证据和下方说明进行判断。审核窗口结束后才可提交裁决。"
                        : "开赛前可以取消。取消确认后，赞助者可自行领取退款。"}
                  </p>
                  {notes.data?.notes
                    .filter((n) => n.kind === "dispute")
                    .map((n) => (
                      <blockquote
                        key={n.evidence_hash}
                        className="rounded-lg border-l-4 border-amber-300 bg-amber-50 p-3 text-sm"
                      >
                        <p className="whitespace-pre-wrap break-words">
                          {n.text}
                        </p>
                        <p className="mt-2 text-xs text-muted-foreground">
                          提交者 {shortAddress(n.author_address)} ·{" "}
                          {events.data?.events.some((e) =>
                            JSON.stringify(e.data)
                              .toLowerCase()
                              .includes(n.evidence_hash.toLowerCase()),
                          )
                            ? "已关联链上事件"
                            : "说明草稿，需核对链上事件"}
                        </p>
                      </blockquote>
                    ))}
                  {!token ? (
                    <LoginButton />
                  ) : (
                    <>
                      {auditor && t.status === "disputed" && (
                        <label className="block text-sm font-medium">
                          裁决结果
                          <select
                            aria-label="异议裁决结果"
                            className="mt-2 block min-h-11 w-full rounded-lg border bg-background p-3"
                            value={decisionWinner}
                            onChange={(e) => setDecisionWinner(e.target.value)}
                          >
                            <option value="">选择裁决结果</option>
                            {t.participants.map((p, i) => (
                              <option key={p} value={p}>
                                冠军：{i === 0 ? "红方" : "蓝方"}{" "}
                                {shortAddress(p)}
                              </option>
                            ))}
                            <option value={zeroAddress}>
                              否决结果，取消赛事并退款
                            </option>
                          </select>
                        </label>
                      )}
                      <label className="block text-sm font-medium">
                        {canDispute
                          ? "异议依据"
                          : t.status === "disputed"
                            ? "裁决理由"
                            : "取消原因"}
                        <textarea
                          className="mt-2 min-h-28 w-full rounded-lg border bg-background p-3 font-normal"
                          aria-label="证据说明"
                          maxLength={4000}
                          value={reason}
                          onChange={(e) => setReason(e.target.value)}
                          placeholder="简要说明原因和相关证据"
                        />
                      </label>
                      {canDispute && (
                        <Button
                          disabled={disabled || !reason.trim()}
                          onClick={() => noteAction("dispute")}
                        >
                          提交异议，查看费用
                        </Button>
                      )}
                      {organizer && t.status === "open" && (
                        <Button
                          variant="outline"
                          disabled={disabled || !reason.trim()}
                          onClick={() => noteAction("cancellation")}
                        >
                          取消赛事，查看费用
                        </Button>
                      )}
                      {auditor && t.status === "disputed" && (
                        <>
                          <Button
                            disabled={
                              disabled ||
                              !reviewComplete ||
                              (decisionWinner !== zeroAddress &&
                                (!evidenceRead ||
                                  !verified ||
                                  !submittedEvidenceMatches)) ||
                              !decisionWinner ||
                              !reason.trim()
                            }
                            onClick={() => noteAction("decision")}
                          >
                            提交裁决，查看费用
                          </Button>
                          {decisionWinner === zeroAddress && (
                            <p className="rounded-lg bg-amber-50 p-3 text-sm text-amber-800">
                              此裁决会取消赛事并允许赞助者退款，不会发放奖金。即使原比赛证据缺失或不一致，也可在说明原因后否决结果。
                            </p>
                          )}
                          {(!reviewComplete ||
                            (decisionWinner !== zeroAddress &&
                              (!evidenceRead ||
                                !verified ||
                                !submittedEvidenceMatches))) && (
                            <p className="text-sm text-muted-foreground">
                              {!reviewComplete
                                ? "请等待审核窗口结束。"
                                : !evidenceRead
                                  ? "请先查看比赛证据。"
                                  : "原比赛证据未通过核验，不能批准发放奖金。"}
                            </p>
                          )}
                        </>
                      )}
                      <p className="text-xs text-muted-foreground">
                        说明会先保存为草稿；钱包交易确认后，操作才正式生效。
                      </p>
                    </>
                  )}
                </div>
              </details>
            </section>
          )}

          <section className="rounded-2xl border bg-card p-6">
            <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
              <h2 className="text-lg font-bold">资金与操作记录</h2>
              <Link
                href={`/audit`}
                className="text-sm font-medium text-primary"
              >
                查看全部流水
                <ArrowRight size={14} className="ml-1 inline" />
              </Link>
            </div>
            {events.error && (
              <p role="alert" className="text-sm text-red-700">
                {errorMessage(events.error)}
              </p>
            )}
            <EventTable events={events.data?.events || []} />
            {notes.data && notes.data.notes.length > 0 && (
              <details className="mt-5 border-t pt-4">
                <summary className="cursor-pointer text-sm font-medium">
                  查看全部审核说明与草稿（{notes.data.notes.length}）
                </summary>
                <div className="mt-3 space-y-3">
                  {notes.data.notes.map((note) => (
                    <article
                      key={note.evidence_hash}
                      className="rounded-xl bg-muted p-4 text-sm"
                    >
                      <p className="mb-2 text-xs text-muted-foreground">
                        {
                          {
                            dispute: "异议说明",
                            decision: "裁决说明",
                            cancellation: "取消说明",
                          }[note.kind]
                        }{" "}
                        · {shortAddress(note.author_address)}
                      </p>
                      <p className="whitespace-pre-wrap break-words">
                        {note.text}
                      </p>
                      <p className="mt-2 text-xs text-muted-foreground">
                        {events.data?.events.some((e) =>
                          JSON.stringify(e.data)
                            .toLowerCase()
                            .includes(note.evidence_hash.toLowerCase()),
                        )
                          ? "已关联本页链上事件"
                          : "草稿或事件不在当前页，请核对完整流水"}
                      </p>
                      <details className="mt-2 text-xs">
                        <summary className="cursor-pointer">内容哈希</summary>
                        <p className="break-all font-mono">
                          {note.evidence_hash}
                        </p>
                      </details>
                    </article>
                  ))}
                </div>
              </details>
            )}
            {notes.error && (
              <p role="alert" className="mt-3 text-sm text-red-700">
                审核说明：{errorMessage(notes.error)}
              </p>
            )}
          </section>
        </div>

        <aside className="min-w-0 space-y-5 lg:sticky lg:top-28">
          <section className="rounded-2xl border bg-card p-5">
            <h2 className="mb-4 flex items-center gap-2 font-bold">
              <Wallet size={18} />
              我的赛事账户
            </h2>
            {address ? (
              <>
                <div className="mb-4 flex flex-wrap gap-2">
                  {(identityLabels.length
                    ? identityLabels
                    : [walletReady ? "观赛者" : "核对身份中"]
                  ).map((label) => (
                    <span
                      key={label}
                      className="rounded-full bg-emerald-50 px-3 py-1 text-sm font-semibold text-emerald-800"
                    >
                      {label}
                    </span>
                  ))}
                </div>
                <p className="mb-5 break-all font-mono text-xs text-muted-foreground">
                  {address}
                </p>
                <dl className="space-y-4">
                  <div>
                    <dt className="text-sm text-muted-foreground">
                      {t.status === "finalized"
                        ? "我的赞助 · 已用于奖金"
                        : t.status === "cancelled"
                          ? "我的待退赞助"
                          : "我的赞助 · 仍在托管"}
                    </dt>
                    <dd className="mt-1 text-2xl font-bold tabular-nums">
                      {amountLabel(held)}{" "}
                      <span className="text-xs font-normal text-muted-foreground">
                        DUSD
                      </span>
                    </dd>
                  </div>
                  <div className="border-t pt-4">
                    <dt className="text-sm text-muted-foreground">
                      {(prize.data || BigInt(0)) > BigInt(0)
                        ? "可领取奖金"
                        : expired && (held || BigInt(0)) > BigInt(0)
                          ? "取消后可退回"
                          : "可领取退款"}
                    </dt>
                    <dd
                      className={`mt-1 text-2xl font-bold tabular-nums ${(refund.data || prize.data || BigInt(0)) > BigInt(0) || (expired && (held || BigInt(0)) > BigInt(0)) ? "text-emerald-700" : ""}`}
                    >
                      {amountLabel(
                        (prize.data || BigInt(0)) > BigInt(0)
                          ? prize.data
                          : expired && (held || BigInt(0)) > BigInt(0)
                            ? held
                            : refund.data,
                      )}{" "}
                      <span className="text-xs font-normal text-muted-foreground">
                        DUSD
                      </span>
                    </dd>
                  </div>
                  <div className="flex justify-between border-t pt-4 text-sm">
                    <dt className="text-muted-foreground">钱包可用余额</dt>
                    <dd className="font-semibold">
                      {amountLabel(balance.data)} DUSD
                    </dd>
                  </div>
                </dl>
                {financialError && (
                  <p role="alert" className="mt-3 text-sm text-red-700">
                    资金读取失败，请重试。
                  </p>
                )}
              </>
            ) : (
              <div className="space-y-4">
                <p className="text-sm leading-6 text-muted-foreground">
                  连接钱包，即可查看你的出资、奖金和退款。
                </p>
                <WalletConnect />
              </div>
            )}
          </section>
          <section className="rounded-2xl border bg-card p-5">
            <h2 className="mb-4 flex items-center gap-2 font-bold">
              <ShieldCheck size={18} />
              奖金与规则
            </h2>
            <dl className="space-y-4 text-sm">
              <div className="flex justify-between gap-3">
                <dt className="text-muted-foreground">仍在托管</dt>
                <dd className="font-semibold">
                  {money(
                    (poolRemaining > BigInt(0)
                      ? poolRemaining
                      : BigInt(0)
                    ).toString(),
                    config?.token_decimals,
                  )}{" "}
                  DUSD
                </dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-muted-foreground">报名费</dt>
                <dd className="font-semibold">免费</dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-muted-foreground">奖金分配</dt>
                <dd>冠军领取全部奖金</dd>
              </div>
              <div className="border-t pt-4">
                <dt className="mb-1 flex items-center gap-2 text-muted-foreground">
                  <Clock3 size={14} />
                  报名截止
                </dt>
                <dd>{dateTime(t.registration_deadline)}</dd>
              </div>
              {reviewState && (
                <div>
                  <dt className="mb-1 text-muted-foreground">审核窗口结束</dt>
                  <dd>{dateTime(t.review_deadline)}</dd>
                </div>
              )}
              <div>
                <dt className="mb-1 text-muted-foreground">最迟完成结算</dt>
                <dd>{dateTime(t.settlement_deadline)}</dd>
              </div>
              {t.winner && (
                <div className="border-t pt-4">
                  <dt className="mb-1 text-muted-foreground">审定冠军</dt>
                  <dd className="break-all font-mono text-xs">{t.winner}</dd>
                </div>
              )}
            </dl>
            {t.metadata_uri && (
              <p className="mt-4 break-all text-sm">
                {/^https?:\/\//.test(t.metadata_uri) ? (
                  <a
                    href={t.metadata_uri}
                    target="_blank"
                    rel="noreferrer"
                    className="font-medium text-primary underline"
                  >
                    查看赛事规则 ↗
                  </a>
                ) : (
                  t.metadata_uri
                )}
              </p>
            )}
            <details className="mt-5 border-t pt-4 text-xs text-muted-foreground">
              <summary className="cursor-pointer">
                查看组织者、审核员与链时间
              </summary>
              <p className="mt-3 break-all">主办方：{t.organizer}</p>
              <p className="mt-2 break-all">独立审核员：{t.auditor}</p>
              <p className="mt-2">
                当前链时间：{now ? dateTime(now) : "读取中…"}
              </p>
              {block.error && (
                <p role="alert" className="mt-2 text-red-700">
                  {errorMessage(block.error)}
                </p>
              )}
              <p className="mt-2">操作期限以链时间为准。</p>
            </details>
          </section>
          <Link
            href="/guide"
            className="flex items-start gap-3 rounded-xl border border-dashed p-4 text-sm"
          >
            <Info size={18} className="mt-0.5 shrink-0 text-primary" />
            <span>
              <strong>第一次体验？</strong>
              <span className="mt-1 block text-muted-foreground">
                了解钱包、角色和完整操作流程。
              </span>
            </span>
          </Link>
        </aside>
      </div>
      <TransactionStatus tx={tx} />
    </div>
  );
}
