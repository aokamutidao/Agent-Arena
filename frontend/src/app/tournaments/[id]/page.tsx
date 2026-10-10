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
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { encodeFunctionData, zeroAddress, type Address } from "viem";
import { tournamentManagerAbi, courseTokenAbi } from "@/lib/course-abi";
import { useAuth } from "@/lib/auth";
import { useI18n } from "@/lib/i18n";
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
  canResumeMatch,
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
function validationError(zh: string, en: string) {
  return Object.assign(new Error(zh), { zh, en });
}

type IdAction =
  | "withdrawParticipant"
  | "startTournament"
  | "approveResult"
  | "expireTournament"
  | "claimPrize"
  | "claimRefund";
export default function TournamentPage() {
  const { lang, l } = useI18n();
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
  const match = useQuery({
    queryKey: ["course", `/api/tournaments/${id}/match`],
    queryFn: () => request<MatchResult>(`/api/tournaments/${id}/match`),
    enabled: !!t && t.status !== "open",
    refetchInterval: (query) =>
      ["queued", "running"].includes(query.state.data?.status || "")
        ? 2000
        : 10000,
    retry: 1,
  });
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
  const [error, setError] = useState<unknown>(null);
  const [viewedEvidence, setViewedEvidence] = useState("");
  const [showEvidence, setShowEvidence] = useState(false);
  const [panel, setPanel] = useState<"join" | "fund">("join");
  useEffect(() => {
    setViewedEvidence("");
    setShowEvidence(false);
    setAgentHash("");
    setReason("");
    setDecisionWinner("");
    setError(null);
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
      l(`${label} · 赛事 #${id}`, `${label} · Tournament #${id}`),
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
    setError(null);
    try {
      const raw = parseMoney(amount, config.token_decimals);
      if (raw <= BigInt(0))
        throw validationError(
          "赞助金额应大于零。",
          "The sponsorship amount must be greater than zero.",
        );
      if (sameAddress(address, t?.auditor))
        throw validationError(
          "指定审计员不能赞助自己审核的赛事。",
          "The designated auditor cannot sponsor a tournament they review.",
        );
      if (!approve && (allowance.data === undefined || allowance.data < raw))
        throw validationError(
          "请先授权赞助金额，等待授权确认，再提交赞助。",
          "Approve the sponsorship amount and wait for confirmation before funding the prize pool.",
        );
      if (!approve && (balance.data === undefined || balance.data < raw))
        throw validationError("DUSD 余额不足。", "Insufficient DUSD balance.");
      await tx.prepare(
        l(
          `${approve ? "授权赞助金额" : "赞助奖金"} ${money(raw.toString(), config.token_decimals)} DUSD · 赛事 #${id}`,
          `${approve ? "Approve sponsorship amount" : "Sponsor prize"} ${money(raw.toString(), config.token_decimals)} DUSD · Tournament #${id}`,
        ),
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
      setError(err);
    }
  }
  async function register() {
    if (!config) return;
    setError(null);
    const selected = profiles.data?.agents.find(
      (a) => a.agent_hash === agentHash,
    );
    if (!selected || !sameAddress(selected.owner_address, address)) {
      setError(
        validationError(
          "请选择自己钱包创建的 Agent 定义。",
          "Choose an Agent definition created by your wallet.",
        ),
      );
      return;
    }
    await tx.prepare(
      l(
        `报名赛事 #${id} · ${selected.name} · ${shortAddress(selected.agent_hash)}`,
        `Register for tournament #${id} · ${selected.name} · ${shortAddress(selected.agent_hash)}`,
      ),
      config.manager_address,
      encodeFunctionData({
        abi: tournamentManagerAbi,
        functionName: "registerParticipant",
        args: [BigInt(id), selected.agent_hash],
      }),
    );
  }
  async function run() {
    if (
      !token ||
      !organizer ||
      !t ||
      t.status !== "running" ||
      now <= 0 ||
      canExpire(t, now) ||
      config?.execution_mode === "unavailable" ||
      (match.data && !canResumeMatch(match.data))
    )
      return;
    setBusy(true);
    setError(null);
    try {
      const result = await request<MatchResult>(`/api/tournaments/${id}/run`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
      });
      queryClient.setQueryData(
        ["course", `/api/tournaments/${id}/match`],
        result,
      );
      await queryClient.invalidateQueries({
        queryKey: ["course", `/api/tournaments/${id}/match`],
      });
    } catch (err) {
      setError(err);
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
      l(
        `提交赛事 #${id} 结果 · 冠军 ${shortAddress(match.data.winner_address)} · 证据 ${shortAddress(evidence.data.evidence_hash)}`,
        `Submit tournament #${id} result · Winner ${shortAddress(match.data.winner_address)} · Evidence ${shortAddress(evidence.data.evidence_hash)}`,
      ),
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
    setError(null);
    try {
      if (!reason.trim() || reason.length > 4000)
        throw validationError(
          "说明需要 1–4000 字符。",
          "The explanation must contain 1–4000 characters.",
        );
      if (kind === "decision" && !decisionWinner)
        throw validationError(
          "请先选择裁决结果。",
          "Choose a dispute decision first.",
        );
      const note = await request<Evidence>(`/api/tournaments/${id}/notes`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
        body: JSON.stringify({ kind, text: reason.trim() }),
      });
      if (!verifyEvidence(note.canonical_json, note.evidence_hash))
        throw validationError(
          "说明内容与哈希不一致，无法提交。",
          "The explanation does not match its hash and cannot be submitted.",
        );
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
        l(
          `${kind === "dispute" ? "提交异议" : kind === "decision" ? `裁决异议：${decisionWinner === zeroAddress ? "取消并退款" : shortAddress(decisionWinner)}` : "取消赛事"} · #${id} · 说明 ${shortAddress(note.evidence_hash)}`,
          `${kind === "dispute" ? "Submit dispute" : kind === "decision" ? `Resolve dispute: ${decisionWinner === zeroAddress ? "Cancel and enable refunds" : shortAddress(decisionWinner)}` : "Cancel tournament"} · #${id} · Explanation ${shortAddress(note.evidence_hash)}`,
        ),
        config.manager_address,
        data,
      );
    } catch (err) {
      setError(err);
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
  const errorText = errorMessage(
    error instanceof Error &&
      "zh" in error &&
      "en" in error &&
      typeof error.zh === "string" &&
      typeof error.en === "string"
      ? l(error.zh, error.en)
      : error,
    lang,
  );
  if (!validId)
    return (
      <div role="alert" className="rounded-2xl border bg-card p-8">
        {l("赛事编号无效。", "Invalid tournament ID.")}
        <Link href="/tournaments" className="ml-3 text-primary underline">
          {l("返回赛事大厅", "Back to tournaments")}
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
          {l("赛事大厅", "Tournaments")}
        </Link>
        <DeploymentNotice />
        <div
          role={tournament.error ? "alert" : "status"}
          className="rounded-2xl border bg-card p-8"
        >
          {tournament.error
            ? errorMessage(tournament.error, lang)
            : l(
                "正在加载赛事与奖金信息…",
                "Loading the tournament and prize details…",
              )}
          {tournament.error && (
            <Button
              className="ml-4"
              variant="outline"
              onClick={() => void tournament.refetch()}
            >
              {l("重新加载", "Reload")}
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
    ownerWallet && l("本场主办方", "Tournament organizer"),
    auditWallet && l("本场审核员", "Tournament auditor"),
    sponsor && l("赞助方", "Sponsor"),
    participant &&
      (winner
        ? l("冠军 · 参赛者", "Winner · Participant")
        : l("参赛者", "Participant")),
  ].filter(Boolean) as string[];
  const presentation = getTournamentPresentation(t, now, lang);
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
    if (requested <= BigInt(0))
      amountError = l(
        "请输入大于零的赞助金额。",
        "Enter a sponsorship amount greater than zero.",
      );
  } catch (err) {
    amountError = errorMessage(err, lang);
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
      {l("查看比赛证据", "View match evidence")}
    </Button>
  );
  let nextTitle = presentation.label;
  let nextDescription = presentation.hint;
  let nextContent: ReactNode = null;
  if (!address) {
    nextTitle = l(
      "连接钱包，查看你的参赛与资金状态",
      "Connect your wallet to see your participation and funds",
    );
    nextDescription = expired
      ? l(
          "这场赛事已到期。赞助者连接原出资钱包后，可先取消赛事，再领取退款。",
          "This tournament has expired. Sponsors can connect the wallet they funded with, cancel the tournament, then claim their refunds.",
        )
      : l(
          "回放和资金记录可以直接浏览；连接钱包后，才能报名、赞助或处理赛事。",
          "Replays and financial records are available to everyone. Connect a wallet to register, sponsor, or manage the tournament.",
        );
    nextContent = <WalletConnect />;
  } else if (financialError) {
    nextTitle = l(
      "暂时无法读取你的资金状态",
      "Your financial status is temporarily unavailable",
    );
    nextDescription = l(
      "请检查钱包网络，重新读取后再操作。",
      "Check your wallet network and reload your financial status before taking action.",
    );
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
        {l("重新读取资金", "Reload financial status")}
      </Button>
    );
  } else if (!walletReady || (now === 0 && notFinal)) {
    nextTitle = l("正在核对你的身份与资金", "Checking your identity and funds");
    nextDescription = l(
      "读取链上记录后，会在这里显示你当前能做的操作。",
      "Your available actions will appear here after the on-chain records load.",
    );
  } else if (expired) {
    nextTitle = sponsor
      ? l(
          "赛事已到期，取回你的赞助",
          "Tournament expired: reclaim your sponsorship",
        )
      : l(
          "赛事已到期，可以取消并开启退款",
          "Tournament expired: cancel to enable refunds",
        );
    nextDescription = sponsor
      ? l(
          `你的 ${amountLabel(held)} DUSD 仍在托管。先确认取消，之后再将退款领回这个钱包。`,
          `Your ${amountLabel(held)} DUSD is still in escrow. Confirm cancellation first, then claim the refund to this wallet.`,
        )
      : l(
          "任何钱包都可以触发到期取消；取消确认后，每位赞助者自行领取退款。",
          "Any wallet can cancel an expired tournament. Once cancellation is confirmed, each sponsor claims their own refund.",
        );
    nextContent = (
      <div className="space-y-4">
        <div className="flex flex-wrap items-center gap-3 text-sm">
          <span className="font-semibold">
            {l("① 确认取消赛事", "① Confirm cancellation")}
          </span>
          <ArrowRight size={16} />
          <span className="text-muted-foreground">
            {l("② 赞助者领取退款", "② Sponsors claim refunds")}
          </span>
        </div>
        <Button
          size="lg"
          disabled={disabled}
          onClick={() =>
            idAction(
              "expireTournament",
              l(
                "取消到期赛事，开启退款",
                "Cancel expired tournament and enable refunds",
              ),
            )
          }
        >
          <ArrowDownToLine size={18} className="mr-2" />
          {l("取消赛事，开启退款", "Cancel tournament and enable refunds")}
        </Button>
      </div>
    );
  } else if ((refund.data || BigInt(0)) > BigInt(0)) {
    nextTitle = l(
      `${amountLabel(refund.data)} DUSD 退款待领取`,
      `${amountLabel(refund.data)} DUSD refund ready to claim`,
    );
    nextDescription = l(
      "赛事已经取消。这笔退款属于当前钱包，确认后将直接退回。",
      "The tournament is cancelled. This refund belongs to your connected wallet and will be returned directly after confirmation.",
    );
    nextContent = (
      <Button
        size="lg"
        disabled={disabled}
        onClick={() =>
          idAction(
            "claimRefund",
            l(
              `领取 ${amountLabel(refund.data)} DUSD 退款`,
              `Claim ${amountLabel(refund.data)} DUSD refund`,
            ),
          )
        }
      >
        <ArrowDownToLine size={18} className="mr-2" />
        {l(
          `领取 ${amountLabel(refund.data)} DUSD 退款`,
          `Claim ${amountLabel(refund.data)} DUSD refund`,
        )}
      </Button>
    );
  } else if ((prize.data || BigInt(0)) > BigInt(0)) {
    nextTitle = l(
      `你赢得了 ${amountLabel(prize.data)} DUSD`,
      `You won ${amountLabel(prize.data)} DUSD`,
    );
    nextDescription = l(
      "审核已通过。领取后，奖金将直接进入当前钱包。",
      "The result has been approved. Claim the prize to receive it directly in your connected wallet.",
    );
    nextContent = (
      <Button
        size="lg"
        disabled={disabled}
        onClick={() =>
          idAction(
            "claimPrize",
            l(
              `领取 ${amountLabel(prize.data)} DUSD 奖金`,
              `Claim ${amountLabel(prize.data)} DUSD prize`,
            ),
          )
        }
      >
        <Trophy size={18} className="mr-2" />
        {l(
          `领取 ${amountLabel(prize.data)} DUSD 奖金`,
          `Claim ${amountLabel(prize.data)} DUSD prize`,
        )}
      </Button>
    );
  } else if (t.status === "cancelled") {
    nextTitle = l("赛事已取消", "Tournament cancelled");
    nextDescription = l(
      "当前钱包没有可领取的退款。你可以在下方核对出资与退款记录。",
      "Your connected wallet has no refund to claim. Check the sponsorship and refund records below.",
    );
  } else if (t.status === "finalized") {
    nextTitle =
      winner && t.prize_claimed
        ? l(
            "奖金已领取，恭喜获胜",
            "Prize claimed. Congratulations on winning!",
          )
        : t.prize_claimed
          ? l("比赛已完成，奖金已发放", "Tournament completed and prize paid")
          : l(
              "审核已完成，等待冠军领取奖金",
              "Review complete: awaiting the winner’s prize claim",
            );
    nextDescription = l(
      "查看比赛回放与资金记录，了解本场赛事的完整过程。",
      "Explore the replay and financial records to follow the full tournament process.",
    );
  } else if ((ownerWallet || auditWallet) && !token) {
    nextTitle = l(
      `${ownerWallet ? "主办方" : "审核员"}，请先签名登录`,
      `${ownerWallet ? "Organizer" : "Auditor"}: sign in with your wallet first`,
    );
    nextDescription = l(
      "钱包已连接。登录验证身份后，就能处理这场赛事；登录签名不扣款。",
      "Your wallet is connected. Sign in to verify your identity and manage this tournament. The sign-in signature does not charge you.",
    );
    nextContent = <LoginButton />;
  } else if ((ownerWallet && !organizer) || (auditWallet && !auditor)) {
    nextTitle = l(
      "当前账户缺少管理权限",
      "Your account lacks the required management role",
    );
    nextDescription = l(
      "该钱包是赛事指定账户，但当前登录权限未通过校验。请重新登录或联系管理员；到期退款仍可按规则处理。",
      "This wallet is designated for the tournament, but your signed-in role could not be verified. Sign in again or contact an administrator. Expiry refunds remain available under the rules.",
    );
  } else if (organizer && t.status === "open") {
    const ready =
      now > 0 &&
      now >= t.registration_deadline &&
      t.participants.length === 2 &&
      pool > BigInt(0) &&
      canSubmitInTime;
    nextTitle = !canSubmitInTime
      ? l(
          "剩余时间不足以开赛并审核",
          "Not enough time remains to start and review the match",
        )
      : ready
        ? l(
            "报名已截止，可以开始比赛",
            "Registration closed: you can start the match",
          )
        : l(
            "准备好奖池和参赛者，等待开赛",
            "Prepare the prize pool and participants before starting",
          );
    nextDescription = !canSubmitInTime
      ? l(
          "你可以在下方取消赛事，让赞助者领取退款。最终期限到达后，任意钱包也可触发取消。",
          "You can cancel the tournament below so sponsors can claim refunds. After the final deadline, any wallet can also trigger cancellation.",
        )
      : l(
          `报名截止：${dateTime(t.registration_deadline, lang)}。两位参赛者和非零奖池齐备后，由你确认开赛。`,
          `Registration deadline: ${dateTime(t.registration_deadline, lang)}. Once both participants and a nonzero prize pool are ready, confirm the tournament start.`,
        );
    nextContent = (
      <div className="space-y-4">
        <div className="flex flex-wrap gap-4 text-sm">
          {[
            [pool > BigInt(0), l("奖池已落实", "Prize pool funded")],
            [
              t.participants.length === 2,
              l("两位参赛者已报名", "Both participants registered"),
            ],
            [
              now >= t.registration_deadline && now > 0,
              l("报名已截止", "Registration closed"),
            ],
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
          onClick={() =>
            idAction("startTournament", l("开始赛事", "Start tournament"))
          }
        >
          <Play size={18} className="mr-2" />
          {l("开始比赛", "Start match")}
        </Button>
      </div>
    );
  } else if (organizer && t.status === "running") {
    if (match.data?.status === "completed") {
      nextTitle = match.data.winner_address
        ? l(
            "对战已完成，请查看证据并提交结果",
            "Match complete: review the evidence and submit the result",
          )
        : l("本场平局，等待到期退款", "Match drawn: awaiting expiry refunds");
      nextDescription = !match.data.winner_address
        ? l(
            "当前规则不为平局指定冠军。最终期限到达后，可取消赛事并退回赞助。",
            "The rules do not assign a winner to a draw. Once the final deadline is reached, the tournament can be cancelled and sponsors can claim refunds.",
          )
        : !canSubmitInTime
          ? l(
              "剩余时间已不足一个完整审核窗口。最终期限到达后，请按规则取消退款。",
              "There is no longer enough time for the full review window. Once the final deadline is reached, cancel the tournament to enable refunds.",
            )
          : l(
              "先查看下方回放和完整证据，再将结果交给独立审核员。",
              "Review the replay and complete evidence below, then submit the result to the independent auditor.",
            );
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
            {l("提交比赛结果", "Submit match result")}
            <ArrowRight size={17} className="ml-2" />
          </Button>
        </div>
      );
    } else {
      const running = ["queued", "running"].includes(match.data?.status || "");
      const failed = match.data?.status === "failed";
      const interrupted = match.data?.status === "interrupted";
      const retryable = canResumeMatch(match.data);
      const unavailable = config?.execution_mode === "unavailable";
      nextTitle = running
        ? l("Agent 正在对战", "Agents are playing the match")
        : failed
          ? l("对战执行失败", "Match execution failed")
          : interrupted
            ? retryable
              ? l(
                  "对战中断，可以继续执行",
                  "Match interrupted: execution can resume",
                )
              : l("对战中断，暂不能恢复", "Match interrupted: resume unavailable")
            : unavailable
              ? l("比赛服务尚不可用", "Match execution is unavailable")
              : l(
                  "比赛已开赛，现在运行 Agent 对战",
                  "Tournament started: run the Agent match",
                );
      nextDescription = running
        ? l(
            "回合和模型调用记录会自动更新。下方展示成功、失败和用量；配置了模型不代表已经调用成功。",
            "Round and model-call records update automatically below, including successes, failures, and usage. Configuring a model does not prove a successful call.",
          )
        : retryable
          ? l(
              "后台允许从已保存的进度恢复；此前失败调用仍会保留，执行模式不变。先核对下方错误，模型凭据或服务有问题时请维护者修复。",
              "The backend allows resuming the saved progress in the same mode. Earlier failed calls remain recorded. Inspect the error below and ask the maintainer to fix model credentials or service issues first.",
            )
          : failed || interrupted
            ? l(
              "后台尚未允许恢复，请查看错误记录并联系维护者；最终期限到达后可按规则退款。",
              "The backend has not allowed resuming. Inspect the error record and contact the maintainer; refunds become available under the rules after the final deadline.",
            )
            : unavailable
              ? l(
                  "未配置可用模型，且没有启用模拟。请维护者配置模型服务或显式启用演示模拟；页面不会把不可用模型当作成功对战。",
                  "No model is configured and simulation is disabled. Ask the maintainer to configure a model or explicitly enable demo simulation; an unavailable model is not treated as a successful match.",
                )
              : l(
                  "系统执行对战并保存回合记录。模型模式与成功调用将在下方明确展示；完成后仍需提交结果供审核。",
                  "The system runs the match and saves each round. Its mode and successful calls are shown below; you must still submit the completed result for review.",
                );
      nextContent = (
        <div className="space-y-3">
          {config?.execution_mode === "llm" && !match.data && (
            <p className="text-sm">
              {l("配置的模型：", "Configured model: ")}
              <span className="font-semibold">{config.llm_model || l("名称待确认", "Name unconfirmed")}</span>
              {config.llm_provider && ` · ${config.llm_provider}`}
              {config.llm_timeout_seconds != null && l(` · 单次超时 ${config.llm_timeout_seconds} 秒`, ` · ${config.llm_timeout_seconds}s call timeout`)}
            </p>
          )}
          {config?.execution_mode === "deterministic_simulation" && !match.data && (
            <p className="text-sm text-amber-800">
              {l("本次将运行确定性模拟，不调用 LLM。", "This run will use deterministic simulation, without calling an LLM.")}
            </p>
          )}
          {match.data?.mode === "llm" && match.data.llm_summary && (
            <p className="text-sm">
              <span className="font-semibold">{match.data.llm_summary.model}</span>
              {l(` · 成功 ${match.data.llm_summary.succeeded} 次 / 失败 ${match.data.llm_summary.failed} 次`, ` · ${match.data.llm_summary.succeeded} successful / ${match.data.llm_summary.failed} failed calls`)}
            </p>
          )}
          <Button
            size="lg"
            disabled={!token || busy || running || unavailable || now <= 0 || (!!match.data && !retryable)}
            onClick={run}
          >
            <Play size={18} className="mr-2" />
            {busy
              ? l("正在请求执行…", "Requesting execution…")
              : running
                ? l("对战进行中…", "Match in progress…")
                : retryable
                  ? l("恢复对战执行", "Resume match execution")
                  : l("运行 Agent 对战", "Run Agent match")}
          </Button>
        </div>
      );
    }
  } else if (auditor && reviewState) {
    nextTitle =
      t.status === "disputed"
        ? l(
            "这场比赛有异议，等待你的裁决",
            "This match is disputed and awaits your decision",
          )
        : l("比赛结果等待你审核", "The match result awaits your review");
    nextDescription = l(
      `先查看证据${t.status === "disputed" ? "和异议说明" : ""}。审核窗口于 ${dateTime(t.review_deadline, lang)} 结束，之后可确认${t.status === "disputed" ? "裁决" : "结果"}。`,
      `Review the evidence${t.status === "disputed" ? " and dispute explanations" : ""} first. The review window ends at ${dateTime(t.review_deadline, lang)}; you can then confirm ${t.status === "disputed" ? "your decision" : "the result"}.`,
    );
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
            onClick={() =>
              idAction(
                "approveResult",
                l("审核通过并结算", "Approve result and finalize"),
              )
            }
          >
            {l("审核通过，确认冠军", "Approve result and confirm winner")}
          </Button>
        ) : (
          <a
            href="#review-action"
            className="inline-flex items-center rounded-lg bg-primary px-5 py-3 text-sm font-semibold text-primary-foreground"
          >
            {l("查看异议并裁决", "Review dispute and decide")}
            <ArrowRight size={16} className="ml-2" />
          </a>
        )}
      </div>
    );
  } else if (open) {
    nextTitle = participant
      ? l("你已报名，等待比赛开始", "You are registered: awaiting the match")
      : sponsor
        ? l("你的赞助已进入奖池", "Your sponsorship is in the prize pool")
        : l(
            "报名参赛，或为这场比赛赞助奖金",
            "Register to compete or sponsor this tournament",
          );
    nextDescription = participant
      ? l(
          "报名截止后由主办方开赛。你的 Agent 已绑定本次报名。",
          "The organizer starts the tournament after registration closes. Your Agent is bound to this registration.",
        )
      : sponsor
        ? l(
            "你可以继续赞助，并在比赛后查看结果；赛事取消时可按出资退款。",
            "You can add more sponsorship and view the result after the match. If the tournament is cancelled, refunds follow each sponsor’s contribution.",
          )
        : l(
            "参赛需要一个自己的 Agent；赞助金额会先进入奖金托管合约。",
            "You need your own Agent to compete. Sponsorship is deposited into the prize escrow contract.",
          );
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
          ? l("查看我的报名", "View my registration")
          : sponsor
            ? l("查看我的赞助", "View my sponsorship")
            : l("选择参与方式", "Choose how to participate")}
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
        {l("返回赛事大厅", "Back to tournaments")}
      </Link>
      <DeploymentNotice />
      <header className="flex flex-wrap items-start justify-between gap-5">
        <div className="min-w-0">
          <div className="mb-3 flex items-center gap-3">
            <span className="text-xs font-semibold tracking-widest text-muted-foreground">
              {l("赛事 #", "Tournament #")}
              {t.id}
            </span>
            <TournamentStatusBadge tournament={t} now={now} />
          </div>
          <h1 className="break-words text-2xl font-bold tracking-tight sm:text-3xl">
            {t.title}
          </h1>
          <p className="mt-2 text-sm text-muted-foreground">
            {l(
              "两位 Agent 参赛 · 赞助奖金预存 · 独立审核后领奖",
              "Two competing Agents · Sponsored prize held in escrow · Claim after independent review",
            )}
          </p>
        </div>
        <div className="rounded-2xl border bg-card px-6 py-4">
          <p className="text-sm text-muted-foreground">
            {l("本场奖金总额", "Total tournament prize")}
          </p>
          <p className="mt-1 text-3xl font-bold tabular-nums">
            {money(t.prize_pool, config?.token_decimals)}{" "}
            <span className="text-sm font-medium text-muted-foreground">
              DUSD
            </span>
          </p>
        </div>
      </header>

      <section
        aria-label={l("赛事进度", "Tournament progress")}
        className="rounded-2xl border bg-card p-4 sm:p-5"
      >
        {t.status === "cancelled" || expired ? (
          <div className="flex flex-wrap items-center gap-3 text-sm">
            <span className="rounded-full bg-amber-100 p-2 text-amber-800">
              <ArrowDownToLine size={18} />
            </span>
            <strong>
              {t.status === "cancelled"
                ? l("赛事已取消", "Tournament cancelled")
                : l("赛事已到期", "Tournament expired")}
            </strong>
            <ArrowRight size={16} className="text-muted-foreground" />
            <span>
              {t.status === "cancelled"
                ? l(
                    "赞助者可领取尚未退回的资金",
                    "Sponsors can claim funds not yet refunded",
                  )
                : l(
                    "确认取消后，赞助者可领取退款",
                    "Sponsors can claim refunds after cancellation is confirmed",
                  )}
            </span>
          </div>
        ) : (
          <ol className="grid grid-cols-2 gap-3 sm:grid-cols-5">
            {[
              l("报名与赞助", "Register & sponsor"),
              l("Agent 对战", "Agent match"),
              l("结果审核", "Result review"),
              l("冠军领奖", "Winner claims prize"),
              l("完成", "Complete"),
            ].map((label, i) => (
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
            ))}
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
                {l("你的下一步", "Your next step")}
              </span>
              {address && (
                <span className="rounded-full border bg-white px-3 py-1 text-xs font-semibold">
                  {identityLabels.length
                    ? identityLabels.join(" / ")
                    : walletReady
                      ? l("观赛者", "Spectator")
                      : l("核对身份中", "Checking identity")}
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
                {l(
                  "钱包网络未就绪，请先处理上方的网络提示。",
                  "Your wallet network is not ready. Follow the network notice above first.",
                )}
              </p>
            )}
            {!!error && (
              <p
                role="alert"
                className="mt-4 rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-800"
              >
                {errorText}
              </p>
            )}
          </section>

          {open && (
            <section
              id="participate"
              className="scroll-mt-28 rounded-2xl border bg-card p-6"
            >
              <div className="mb-5 flex flex-wrap items-center justify-between gap-3">
                <h2 className="text-lg font-bold">
                  {l("参与这场赛事", "Participate in this tournament")}
                </h2>
                <span className="text-sm text-muted-foreground">
                  {t.participants.length} {l("/ 2 人已报名", "/ 2 registered")}
                </span>
              </div>
              {!address ? (
                <div className="space-y-4">
                  <p className="text-sm text-muted-foreground">
                    {l(
                      "连接钱包后，可以选择自己的 Agent 报名，或赞助奖金。",
                      "Connect a wallet to register your Agent or sponsor the prize.",
                    )}
                  </p>
                  <WalletConnect />
                </div>
              ) : (
                <>
                  <div
                    className="mb-5 flex gap-2 rounded-xl bg-muted p-1"
                    role="group"
                    aria-label={l("参与方式", "Participation options")}
                  >
                    <Button
                      className="flex-1"
                      variant={panel === "join" ? "default" : "ghost"}
                      aria-pressed={panel === "join"}
                      onClick={() => setPanel("join")}
                    >
                      <Users size={16} className="mr-2" />
                      {l("我要参赛", "Compete")}
                    </Button>
                    <Button
                      className="flex-1"
                      variant={panel === "fund" ? "default" : "ghost"}
                      aria-pressed={panel === "fund"}
                      onClick={() => setPanel("fund")}
                    >
                      <Coins size={16} className="mr-2" />
                      {l("赞助奖金", "Sponsor prize")}
                    </Button>
                  </div>
                  {panel === "join" ? (
                    participant ? (
                      <div className="space-y-4">
                        <p className="flex items-center gap-2 font-medium text-emerald-700">
                          <CheckCircle2 size={19} />
                          {l("你已成功报名", "You are registered")}
                        </p>
                        <p className="text-sm text-muted-foreground">
                          {l(
                            "报名截止前可以退出；开赛后由已报名的 Agent 参加对战。",
                            "You can withdraw before registration closes. Once the tournament starts, your registered Agent will compete.",
                          )}
                        </p>
                        <Button
                          variant="outline"
                          disabled={disabled}
                          onClick={() =>
                            idAction(
                              "withdrawParticipant",
                              l("退出报名", "Withdraw registration"),
                            )
                          }
                        >
                          {l("退出本场报名", "Withdraw from this tournament")}
                        </Button>
                      </div>
                    ) : !canJoin ? (
                      <p className="rounded-xl bg-muted p-4 text-sm">
                        {ownerWallet || auditWallet
                          ? l(
                              "本场主办方和审核员保持独立，不能参赛。",
                              "The organizer and auditor must remain independent and cannot compete in this tournament.",
                            )
                          : l(
                              "本场两个参赛名额已满，你仍可以观赛或赞助。",
                              "Both participant slots are filled. You can still watch or sponsor the tournament.",
                            )}
                      </p>
                    ) : !token ? (
                      <div className="space-y-3">
                        <p className="text-sm text-muted-foreground">
                          {l(
                            "签名登录后，即可选择自己的 Agent 报名。报名不收费用。",
                            "Sign in with your wallet signature to register your own Agent. Registration is free.",
                          )}
                        </p>
                        <LoginButton />
                      </div>
                    ) : (
                      <div className="space-y-4">
                        {profiles.isLoading ? (
                          <p role="status">
                            {l("正在读取你的 Agent…", "Loading your Agents…")}
                          </p>
                        ) : profiles.error ? (
                          <p role="alert" className="text-sm text-red-700">
                            {errorMessage(profiles.error, lang)}
                          </p>
                        ) : !profiles.data?.agents.length ? (
                          <div className="rounded-xl border border-dashed p-5">
                            <p className="mb-3 text-sm">
                              {l(
                                "你还没有参赛 Agent，先为它取名并设定策略。",
                                "You do not have an Agent yet. Create one with a name and strategy first.",
                              )}
                            </p>
                            <Link
                              href="/agents"
                              className="inline-flex items-center font-semibold text-primary"
                            >
                              {l("创建我的 Agent", "Create my Agent")}
                              <ArrowRight size={16} className="ml-2" />
                            </Link>
                          </div>
                        ) : (
                          <>
                            <label className="block text-sm font-medium">
                              {l(
                                "选择报名 Agent",
                                "Choose an Agent to register",
                              )}
                              <select
                                aria-label={l(
                                  "选择报名 Agent",
                                  "Choose an Agent to register",
                                )}
                                className="mt-2 min-h-11 w-full rounded-lg border bg-background p-3"
                                value={agentHash}
                                onChange={(e) => setAgentHash(e.target.value)}
                              >
                                <option value="">
                                  {l("请选择你的 Agent", "Choose your Agent")}
                                </option>
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
                              {l(
                                "确认报名，查看费用",
                                "Register and review fees",
                              )}
                            </Button>
                            <p className="text-xs text-muted-foreground">
                              {l(
                                "报名费 0 DUSD，仅需本地网络费。",
                                "Registration costs 0 DUSD. Only local network fees apply.",
                              )}
                              <Link href="/agents" className="ml-2 underline">
                                {l("管理我的 Agent", "Manage my Agents")}
                              </Link>
                            </p>
                          </>
                        )}
                      </div>
                    )
                  ) : auditWallet ? (
                    <p className="rounded-xl bg-muted p-4 text-sm">
                      {l(
                        "你是本场审核员，需要保持独立，不能赞助这场赛事。",
                        "You are this tournament’s auditor and must remain independent. You cannot sponsor this tournament.",
                      )}
                    </p>
                  ) : (
                    <div className="space-y-4">
                      <label className="block text-sm font-medium">
                        {l("赞助金额（DUSD）", "Sponsorship amount (DUSD)")}
                        <Input
                          className="mt-2 h-12 text-lg"
                          inputMode="decimal"
                          value={amount}
                          aria-invalid={!!amountError}
                          aria-describedby={
                            amountError ? "sponsorship-amount-error" : undefined
                          }
                          onChange={(e) => setAmount(e.target.value)}
                          placeholder={l("例如 100", "For example, 100")}
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
                        {l("钱包可用", "Wallet balance")}{" "}
                        {amountLabel(balance.data)}{" "}
                        {l(
                          "DUSD · 本场仍托管",
                          "DUSD · Escrowed for this tournament",
                        )}{" "}
                        {amountLabel(held)} DUSD
                      </p>
                      <div className="rounded-xl bg-muted p-4">
                        <p className="text-sm font-medium">
                          {enoughAllowance
                            ? l(
                                "✓ 金额已授权，可以确认赞助",
                                "✓ Amount approved: ready to sponsor",
                              )
                            : l(
                                "先授权金额，再确认赞助",
                                "Approve the amount, then confirm sponsorship",
                              )}
                        </p>
                        <p className="mt-1 text-xs leading-5 text-muted-foreground">
                          {l(
                            "授权不扣款。确认赞助后，资金进入奖池；赛事取消时按原出资退款。",
                            "Approval does not transfer funds. Confirming sponsorship deposits funds in the prize pool. Cancellation refunds each sponsor’s contribution.",
                          )}
                        </p>
                      </div>
                      <div className="flex flex-wrap gap-3">
                        <Button
                          variant={enoughAllowance ? "outline" : "default"}
                          disabled={disabled || requested <= BigInt(0)}
                          onClick={() => fund(true)}
                        >
                          {l("1. 授权", "1. Approve")}{" "}
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
                          {l("2. 确认赞助", "2. Confirm sponsorship")}
                        </Button>
                      </div>
                      {requested > BigInt(0) &&
                        balance.data !== undefined &&
                        balance.data < requested && (
                          <p className="text-sm text-red-700">
                            {l(
                              "DUSD 余额不足，请换到有测试币的赞助钱包。",
                              "Insufficient DUSD balance. Switch to a sponsor wallet with test tokens.",
                            )}
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
              <h2 className="text-lg font-bold">
                {l("参赛阵容", "Participants")}
              </h2>
              <span className="text-sm text-muted-foreground">
                {t.participants.length} {l("/ 2 席", "/ 2 slots")}
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
                      {i === 0 ? l("红方", "Red team") : l("蓝方", "Blue team")}
                    </p>
                    <p className="break-words font-semibold">
                      {p
                        ? info?.name || l("已报名 Agent", "Registered Agent")
                        : l("等待参赛者", "Awaiting participant")}
                    </p>
                    {p ? (
                      <>
                        <p
                          className="mt-2 font-mono text-xs text-muted-foreground"
                          title={p.address}
                        >
                          {shortAddress(p.address)}
                          {sameAddress(p.address, address) &&
                            l(" · 你", " · You")}
                        </p>
                        <details className="mt-3 text-xs text-muted-foreground">
                          <summary className="cursor-pointer">
                            {l("报名凭证", "Registration proof")}
                          </summary>
                          <p className="mt-2 break-all font-mono">
                            {p.agent_hash}
                          </p>
                        </details>
                      </>
                    ) : (
                      <p className="mt-2 text-sm text-muted-foreground">
                        {l(
                          "创建自己的 Agent 后即可报名",
                          "Create your own Agent to register",
                        )}
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
                  ? l(
                      "这场赛事尚未生成对战回放",
                      "No match replay has been generated for this tournament",
                    )
                  : l(
                      "对战开始后，在这里观看回放",
                      "Watch the replay here once the match starts",
                    )}
              </h2>
              <p className="mt-2 text-sm text-muted-foreground">
                {l(
                  "地图、血量和每回合动作会随对战记录显示。",
                  "The map, health, and actions for each round appear with the match records.",
                )}
              </p>
            </section>
          )}
          {match.error &&
            !(
              match.error instanceof ApiError && match.error.status === 404
            ) && (
              <p role="alert" className="text-sm text-red-700">
                {l("对战数据：", "Match data: ")}
                {errorMessage(match.error, lang)}
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
                  {l("比赛证据", "Match evidence")}
                </h2>
                <span
                  className={`rounded-full px-3 py-1 text-xs font-semibold ${verified && (!t.evidence_hash || submittedEvidenceMatches) ? "bg-emerald-50 text-emerald-700" : "bg-red-50 text-red-700"}`}
                >
                  {verified
                    ? t.evidence_hash
                      ? submittedEvidenceMatches
                        ? l(
                            "内容与链上记录一致",
                            "Content matches the on-chain record",
                          )
                        : l("链上记录不一致", "On-chain record does not match")
                      : l("内容校验通过", "Content verification passed")
                    : l("内容校验失败", "Content verification failed")}
                </span>
              </div>
              <p className="mt-3 text-sm leading-6 text-muted-foreground">
                {l(
                  "包含双方 Agent、每回合动作和结果。主办方提交与审核员批准前，都需要先查看。",
                  "Includes both Agents, every round’s actions, and the result. The organizer and auditor must review it before submitting or approving the result.",
                )}
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
                  {showEvidence
                    ? l("收起完整证据", "Hide full evidence")
                    : l("展开完整证据", "Show full evidence")}
                </Button>
                <Button variant="outline" onClick={downloadEvidence}>
                  <ArrowDownToLine size={16} className="mr-2" />
                  {l("下载证据 JSON", "Download evidence JSON")}
                </Button>
                {evidenceRead && (
                  <span className="flex items-center gap-1 text-sm text-emerald-700">
                    <Check size={15} />
                    {l("已查看", "Viewed")}
                  </span>
                )}
              </div>
              {showEvidence && (
                <pre className="mt-4 max-h-80 overflow-auto whitespace-pre-wrap break-all rounded-xl bg-muted p-4 text-xs">
                  {evidence.data.canonical_json}
                </pre>
              )}
              <details className="mt-4 text-xs text-muted-foreground">
                <summary className="cursor-pointer">
                  {l("证据指纹与核验说明", "Evidence hash and verification")}
                </summary>
                <p className="mt-2 break-all font-mono">
                  {evidence.data.evidence_hash}
                </p>
                <p className="mt-2 leading-5">
                  {l(
                    "哈希核对用于发现内容变更。审核员仍需判断比赛结果和规则是否正确。",
                    "Hash verification detects content changes. The auditor must still judge whether the result and rules are correct.",
                  )}
                </p>
              </details>
            </section>
          )}
          {evidence.error && (
            <p role="alert" className="text-sm text-red-700">
              {l("证据加载失败：", "Evidence could not be loaded: ")}
              {errorMessage(evidence.error, lang)}
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
                    ? l("审核异议并裁决", "Review and resolve the dispute")
                    : canDispute
                      ? l(
                          "对结果有异议？提交复核申请",
                          "Disagree with the result? Submit a dispute",
                        )
                      : l(
                          "需要取消这场赛事？",
                          "Need to cancel this tournament?",
                        )}
                </summary>
                <div className="mt-4 space-y-4">
                  <p className="text-sm leading-6 text-muted-foreground">
                    {canDispute
                      ? l(
                          `在 ${dateTime(t.review_deadline, lang)} 前说明问题并提交，审核员会处理你的异议。`,
                          `Explain and submit your dispute before ${dateTime(t.review_deadline, lang)}. The auditor will review it.`,
                        )
                      : t.status === "disputed"
                        ? l(
                            "结合回放、证据和下方说明进行判断。审核窗口结束后才可提交裁决。",
                            "Consider the replay, evidence, and explanations below. You can submit a decision only after the review window ends.",
                          )
                        : l(
                            "开赛前可以取消。取消确认后，赞助者可自行领取退款。",
                            "You can cancel before the tournament starts. After cancellation is confirmed, sponsors can claim their own refunds.",
                          )}
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
                          {l("提交者", "Submitted by")}{" "}
                          {shortAddress(n.author_address)} ·{" "}
                          {events.data?.events.some((e) =>
                            JSON.stringify(e.data)
                              .toLowerCase()
                              .includes(n.evidence_hash.toLowerCase()),
                          )
                            ? l("已关联链上事件", "Linked to an on-chain event")
                            : l(
                                "说明草稿，需核对链上事件",
                                "Draft explanation: check the on-chain event",
                              )}
                        </p>
                      </blockquote>
                    ))}
                  {!token ? (
                    <LoginButton />
                  ) : (
                    <>
                      {auditor && t.status === "disputed" && (
                        <label className="block text-sm font-medium">
                          {l("裁决结果", "Dispute decision")}
                          <select
                            aria-label={l(
                              "异议裁决结果",
                              "Dispute resolution outcome",
                            )}
                            className="mt-2 block min-h-11 w-full rounded-lg border bg-background p-3"
                            value={decisionWinner}
                            onChange={(e) => setDecisionWinner(e.target.value)}
                          >
                            <option value="">
                              {l("选择裁决结果", "Choose a decision")}
                            </option>
                            {t.participants.map((p, i) => (
                              <option key={p} value={p}>
                                {l("冠军：", "Winner: ")}
                                {i === 0
                                  ? l("红方", "Red team")
                                  : l("蓝方", "Blue team")}{" "}
                                {shortAddress(p)}
                              </option>
                            ))}
                            <option value={zeroAddress}>
                              {l(
                                "否决结果，取消赛事并退款",
                                "Reject result: cancel tournament and enable refunds",
                              )}
                            </option>
                          </select>
                        </label>
                      )}
                      <label className="block text-sm font-medium">
                        {canDispute
                          ? l("异议依据", "Grounds for dispute")
                          : t.status === "disputed"
                            ? l("裁决理由", "Reason for decision")
                            : l("取消原因", "Reason for cancellation")}
                        <textarea
                          className="mt-2 min-h-28 w-full rounded-lg border bg-background p-3 font-normal"
                          aria-label={l("证据说明", "Evidence explanation")}
                          maxLength={4000}
                          value={reason}
                          onChange={(e) => setReason(e.target.value)}
                          placeholder={l(
                            "简要说明原因和相关证据",
                            "Briefly explain the reasons and supporting evidence",
                          )}
                        />
                      </label>
                      {canDispute && (
                        <Button
                          disabled={disabled || !reason.trim()}
                          onClick={() => noteAction("dispute")}
                        >
                          {l(
                            "提交异议，查看费用",
                            "Submit dispute and review fees",
                          )}
                        </Button>
                      )}
                      {organizer && t.status === "open" && (
                        <Button
                          variant="outline"
                          disabled={disabled || !reason.trim()}
                          onClick={() => noteAction("cancellation")}
                        >
                          {l(
                            "取消赛事，查看费用",
                            "Cancel tournament and review fees",
                          )}
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
                            {l(
                              "提交裁决，查看费用",
                              "Submit decision and review fees",
                            )}
                          </Button>
                          {decisionWinner === zeroAddress && (
                            <p className="rounded-lg bg-amber-50 p-3 text-sm text-amber-800">
                              {l(
                                "此裁决会取消赛事并允许赞助者退款，不会发放奖金。即使原比赛证据缺失或不一致，也可在说明原因后否决结果。",
                                "This decision cancels the tournament and enables sponsor refunds. No prize will be paid. You may reject a result with missing or inconsistent evidence after explaining your reasons.",
                              )}
                            </p>
                          )}
                          {(!reviewComplete ||
                            (decisionWinner !== zeroAddress &&
                              (!evidenceRead ||
                                !verified ||
                                !submittedEvidenceMatches))) && (
                            <p className="text-sm text-muted-foreground">
                              {!reviewComplete
                                ? l(
                                    "请等待审核窗口结束。",
                                    "Wait for the review window to end.",
                                  )
                                : !evidenceRead
                                  ? l(
                                      "请先查看比赛证据。",
                                      "View the match evidence first.",
                                    )
                                  : l(
                                      "原比赛证据未通过核验，不能批准发放奖金。",
                                      "The original match evidence failed verification. Prize payment cannot be approved.",
                                    )}
                            </p>
                          )}
                        </>
                      )}
                      <p className="text-xs text-muted-foreground">
                        {l(
                          "说明会先保存为草稿；钱包交易确认后，操作才正式生效。",
                          "The explanation is saved as a draft first. The action takes effect only after the wallet transaction is confirmed.",
                        )}
                      </p>
                    </>
                  )}
                </div>
              </details>
            </section>
          )}

          <section className="rounded-2xl border bg-card p-6">
            <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
              <h2 className="text-lg font-bold">
                {l("资金与操作记录", "Financial and activity records")}
              </h2>
              <Link
                href={`/audit`}
                className="text-sm font-medium text-primary"
              >
                {l("查看全部流水", "View all records")}
                <ArrowRight size={14} className="ml-1 inline" />
              </Link>
            </div>
            {events.error && (
              <p role="alert" className="text-sm text-red-700">
                {errorMessage(events.error, lang)}
              </p>
            )}
            <EventTable events={events.data?.events || []} />
            {notes.data && notes.data.notes.length > 0 && (
              <details className="mt-5 border-t pt-4">
                <summary className="cursor-pointer text-sm font-medium">
                  {l(
                    `查看全部审核说明与草稿（${notes.data.notes.length}）`,
                    `View all review notes and drafts (${notes.data.notes.length})`,
                  )}
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
                            dispute: l("异议说明", "Dispute explanation"),
                            decision: l("裁决说明", "Decision explanation"),
                            cancellation: l(
                              "取消说明",
                              "Cancellation explanation",
                            ),
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
                          ? l(
                              "已关联本页链上事件",
                              "Linked to an on-chain event on this page",
                            )
                          : l(
                              "草稿或事件不在当前页，请核对完整流水",
                              "Draft or event not on this page: check the complete records",
                            )}
                      </p>
                      <details className="mt-2 text-xs">
                        <summary className="cursor-pointer">
                          {l("内容哈希", "Content hash")}
                        </summary>
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
                {l("审核说明：", "Review notes: ")}
                {errorMessage(notes.error, lang)}
              </p>
            )}
          </section>
        </div>

        <aside className="min-w-0 space-y-5 lg:sticky lg:top-28">
          <section className="rounded-2xl border bg-card p-5">
            <h2 className="mb-4 flex items-center gap-2 font-bold">
              <Wallet size={18} />
              {l("我的赛事账户", "My tournament account")}
            </h2>
            {address ? (
              <>
                <div className="mb-4 flex flex-wrap gap-2">
                  {(identityLabels.length
                    ? identityLabels
                    : [
                        walletReady
                          ? l("观赛者", "Spectator")
                          : l("核对身份中", "Checking identity"),
                      ]
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
                        ? l(
                            "我的赞助 · 已用于奖金",
                            "My sponsorship · Allocated to the prize",
                          )
                        : t.status === "cancelled"
                          ? l("我的待退赞助", "My sponsorship awaiting refund")
                          : l(
                              "我的赞助 · 仍在托管",
                              "My sponsorship · Held in escrow",
                            )}
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
                        ? l("可领取奖金", "Claimable prize")
                        : expired && (held || BigInt(0)) > BigInt(0)
                          ? l("取消后可退回", "Refundable after cancellation")
                          : l("可领取退款", "Claimable refund")}
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
                    <dt className="text-muted-foreground">
                      {l("钱包可用余额", "Available wallet balance")}
                    </dt>
                    <dd className="font-semibold">
                      {amountLabel(balance.data)} DUSD
                    </dd>
                  </div>
                </dl>
                {financialError && (
                  <p role="alert" className="mt-3 text-sm text-red-700">
                    {l(
                      "资金读取失败，请重试。",
                      "Financial data could not be loaded. Try again.",
                    )}
                  </p>
                )}
              </>
            ) : (
              <div className="space-y-4">
                <p className="text-sm leading-6 text-muted-foreground">
                  {l(
                    "连接钱包，即可查看你的出资、奖金和退款。",
                    "Connect your wallet to view your sponsorship, prize, and refunds.",
                  )}
                </p>
                <WalletConnect />
              </div>
            )}
          </section>
          <section className="rounded-2xl border bg-card p-5">
            <h2 className="mb-4 flex items-center gap-2 font-bold">
              <ShieldCheck size={18} />
              {l("奖金与规则", "Prize and rules")}
            </h2>
            <dl className="space-y-4 text-sm">
              <div className="flex justify-between gap-3">
                <dt className="text-muted-foreground">
                  {l("仍在托管", "Held in escrow")}
                </dt>
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
                <dt className="text-muted-foreground">
                  {l("报名费", "Registration fee")}
                </dt>
                <dd className="font-semibold">{l("免费", "Free")}</dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-muted-foreground">
                  {l("奖金分配", "Prize distribution")}
                </dt>
                <dd>
                  {l("冠军领取全部奖金", "Winner claims the entire prize")}
                </dd>
              </div>
              <div className="border-t pt-4">
                <dt className="mb-1 flex items-center gap-2 text-muted-foreground">
                  <Clock3 size={14} />
                  {l("报名截止", "Registration deadline")}
                </dt>
                <dd>{dateTime(t.registration_deadline, lang)}</dd>
              </div>
              {reviewState && (
                <div>
                  <dt className="mb-1 text-muted-foreground">
                    {l("审核窗口结束", "Review window ends")}
                  </dt>
                  <dd>{dateTime(t.review_deadline, lang)}</dd>
                </div>
              )}
              <div>
                <dt className="mb-1 text-muted-foreground">
                  {l("最迟完成结算", "Final settlement deadline")}
                </dt>
                <dd>{dateTime(t.settlement_deadline, lang)}</dd>
              </div>
              {t.winner && (
                <div className="border-t pt-4">
                  <dt className="mb-1 text-muted-foreground">
                    {l("审定冠军", "Approved winner")}
                  </dt>
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
                    {l("查看赛事规则 ↗", "View tournament rules ↗")}
                  </a>
                ) : (
                  t.metadata_uri
                )}
              </p>
            )}
            <details className="mt-5 border-t pt-4 text-xs text-muted-foreground">
              <summary className="cursor-pointer">
                {l(
                  "查看组织者、审核员与链时间",
                  "View organizer, auditor, and chain time",
                )}
              </summary>
              <p className="mt-3 break-all">
                {l("主办方：", "Organizer: ")}
                {t.organizer}
              </p>
              <p className="mt-2 break-all">
                {l("独立审核员：", "Independent auditor: ")}
                {t.auditor}
              </p>
              <p className="mt-2">
                {l("当前链时间：", "Current chain time: ")}
                {now ? dateTime(now, lang) : l("读取中…", "Loading…")}
              </p>
              {block.error && (
                <p role="alert" className="mt-2 text-red-700">
                  {errorMessage(block.error, lang)}
                </p>
              )}
              <p className="mt-2">
                {l(
                  "操作期限以链时间为准。",
                  "Action deadlines follow the chain time.",
                )}
              </p>
            </details>
          </section>
          <Link
            href="/guide"
            className="flex items-start gap-3 rounded-xl border border-dashed p-4 text-sm"
          >
            <Info size={18} className="mt-0.5 shrink-0 text-primary" />
            <span>
              <strong>{l("第一次体验？", "First time here?")}</strong>
              <span className="mt-1 block text-muted-foreground">
                {l(
                  "了解钱包、角色和完整操作流程。",
                  "Learn about wallets, roles, and the full workflow.",
                )}
              </span>
            </span>
          </Link>
        </aside>
      </div>
      <TransactionStatus tx={tx} />
    </div>
  );
}
