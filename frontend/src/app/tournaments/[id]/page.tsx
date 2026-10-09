"use client";
import { useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { useAccount, useBlock, useReadContract } from "wagmi";
import { useQueryClient } from "@tanstack/react-query";
import { encodeFunctionData, zeroAddress, type Address } from "viem";
import { tournamentManagerAbi, courseTokenAbi } from "@/lib/course-abi";
import { useAuth } from "@/lib/auth";
import { ApiError, dateTime, errorMessage, money, parseMoney, request, sameAddress, shortAddress, STATUS_LABEL, registrationOpen, reviewComplete as isReviewComplete, canExpire, verifyEvidence, type Tournament, type AgentProfile, type AuditEvent, type MatchResult, type Evidence, type EvidenceNote } from "@/lib/course";
import { useAuthenticatedResource, useCourseConfig, useCourseResource } from "@/hooks/useCourse";
import { useWalletTransaction } from "@/hooks/useWalletTransaction";
import { DeploymentNotice } from "@/components/course/DeploymentNotice";
import { TransactionStatus } from "@/components/course/TransactionStatus";
import { EventTable } from "@/components/course/EventTable";
import { MatchReplay } from "@/components/course/MatchReplay";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
type IdAction = "withdrawParticipant" | "startTournament" | "approveResult" | "expireTournament" | "claimPrize" | "claimRefund";
export default function TournamentPage() {
    const params = useParams();
    const id = Number(params.id);
    const validId = Number.isSafeInteger(id) && id > 0;
    const { address } = useAccount();
    const { token, user } = useAuth();
    const { data: config } = useCourseConfig();
    const tournament = useCourseResource<Tournament>(`/api/tournaments/${id}`, validId);
    const t = tournament.data;
    const events = useCourseResource<{
        events: AuditEvent[];
        total: number;
    }>(`/api/tournaments/${id}/events?limit=50&offset=0`, !!t);
    const profiles = useAuthenticatedResource<{
        agents: AgentProfile[];
        total: number;
    }>("/api/agents/my", token);
    const match = useCourseResource<MatchResult>(`/api/tournaments/${id}/match`, !!t && t.status !== "open");
    const evidence = useCourseResource<Evidence>(`/api/tournaments/${id}/evidence`, match.data?.status === "completed");
    const notes = useCourseResource<{
        notes: EvidenceNote[];
        total: number;
    }>(`/api/tournaments/${id}/notes`, !!t);
    const queryClient = useQueryClient();
    const tx = useWalletTransaction();
    const [amount, setAmount] = useState("");
    const [agentHash, setAgentHash] = useState("");
    const [reason, setReason] = useState("");
    const [decisionWinner, setDecisionWinner] = useState("");
    const [busy, setBusy] = useState(false);
    const [error, setError] = useState("");
    const [viewedEvidence, setViewedEvidence] = useState("");
    const block = useBlock({ chainId: config?.chain_id, watch: true, query: { enabled: !!config?.deployment_configured, refetchInterval: 5000 } });
    const now = block.data ? Number(block.data.timestamp) : 0;
    const contractArgs = address && validId ? [BigInt(id), address] as const : undefined;
    const readEnabled = !!config?.deployment_configured && !!address && validId;
    const contribution = useReadContract({ address: config?.manager_address, abi: tournamentManagerAbi, functionName: "contributionOf", args: contractArgs, chainId: config?.chain_id, query: { enabled: readEnabled, refetchInterval: 10000 } });
    const refund = useReadContract({ address: config?.manager_address, abi: tournamentManagerAbi, functionName: "refundableAmount", args: contractArgs, chainId: config?.chain_id, query: { enabled: readEnabled, refetchInterval: 10000 } });
    const prize = useReadContract({ address: config?.manager_address, abi: tournamentManagerAbi, functionName: "claimablePrize", args: contractArgs, chainId: config?.chain_id, query: { enabled: readEnabled, refetchInterval: 10000 } });
    const allowance = useReadContract({ address: config?.token_address, abi: courseTokenAbi, functionName: "allowance", args: address && config ? [address, config.escrow_address] : undefined, chainId: config?.chain_id, query: { enabled: readEnabled, refetchInterval: 10000 } });
    const balance = useReadContract({ address: config?.token_address, abi: courseTokenAbi, functionName: "balanceOf", args: address ? [address] : undefined, chainId: config?.chain_id, query: { enabled: readEnabled, refetchInterval: 10000 } });
    const organizer = !!user?.roles.includes("organizer") && sameAddress(address, t?.organizer);
    const auditor = !!user?.roles.includes("auditor") && sameAddress(address, t?.auditor);
    const participant = t?.participants.some(a => sameAddress(a, address)) || false;
    const open = !!t && now > 0 && registrationOpen(t, now);
    const reviewComplete = !!t && now > 0 && isReviewComplete(t, now);
    const expired = !!t && now > 0 && canExpire(t, now);
    const verified = !!evidence.data && verifyEvidence(evidence.data.canonical_json, evidence.data.evidence_hash) && evidence.data.evidence_hash === match.data?.evidence_hash;
    const evidenceRead = !!evidence.data && viewedEvidence === evidence.data.evidence_hash;
    const submittedEvidenceMatches = !!evidence.data && t?.evidence_hash === evidence.data.evidence_hash;
    const disabled = !tx.canWrite || tx.locked || busy;
    async function idAction(name: IdAction, label: string) {
        if (!config)
            return;
        await tx.prepare(`${label} · 赛事 #${id}`, config.manager_address, encodeFunctionData({ abi: tournamentManagerAbi, functionName: name, args: [BigInt(id)] }));
    }
    async function fund(approve: boolean) {
        if (!config)
            return;
        setError("");
        try {
            const raw = parseMoney(amount, config.token_decimals);
            if (raw <= BigInt(0))
                throw new Error("赞助金额应大于零。");
            if (sameAddress(address, t?.auditor))
                throw new Error("指定审计员不能赞助自己审核的赛事。");
            if (!approve && (allowance.data === undefined || allowance.data < raw))
                throw new Error("请先授权赞助金额，等待授权确认，再提交赞助。");
            if (!approve && (balance.data === undefined || balance.data < raw))
                throw new Error("DUSD 余额不足。");
            await tx.prepare(`${approve ? "授权赞助金额" : "赞助奖金"} ${money(raw.toString(), config.token_decimals)} DUSD · 赛事 #${id}`, approve ? config.token_address : config.manager_address, approve ? encodeFunctionData({ abi: courseTokenAbi, functionName: "approve", args: [config.escrow_address, raw] }) : encodeFunctionData({ abi: tournamentManagerAbi, functionName: "fundTournament", args: [BigInt(id), raw] }));
        }
        catch (err) {
            setError(errorMessage(err));
        }
    }
    async function register() {
        if (!config)
            return;
        setError("");
        const selected = profiles.data?.agents.find(a => a.agent_hash === agentHash);
        if (!selected || !sameAddress(selected.owner_address, address)) {
            setError("请选择自己钱包创建的 Agent 定义。");
            return;
        }
        await tx.prepare(`报名赛事 #${id} · ${selected.name} · ${shortAddress(selected.agent_hash)}`, config.manager_address, encodeFunctionData({ abi: tournamentManagerAbi, functionName: "registerParticipant", args: [BigInt(id), selected.agent_hash] }));
    }
    async function run() {
        if (!token)
            return;
        setBusy(true);
        setError("");
        try {
            await request<MatchResult>(`/api/tournaments/${id}/run`, { method: "POST", headers: { Authorization: `Bearer ${token}` } });
            await queryClient.invalidateQueries({ queryKey: ["course", `/api/tournaments/${id}/match`] });
        }
        catch (err) {
            setError(errorMessage(err));
        }
        finally {
            setBusy(false);
        }
    }
    async function submit() {
        if (!config || !verified || !evidenceRead || !match.data?.winner_address || !evidence.data)
            return;
        await tx.prepare(`提交赛事 #${id} 结果 · 冠军 ${shortAddress(match.data.winner_address)} · 证据 ${shortAddress(evidence.data.evidence_hash)}`, config.manager_address, encodeFunctionData({ abi: tournamentManagerAbi, functionName: "submitResult", args: [BigInt(id), match.data.winner_address, evidence.data.evidence_hash] }));
    }
    async function noteAction(kind: "dispute" | "decision" | "cancellation") {
        if (!config || !token)
            return;
        setBusy(true);
        setError("");
        try {
            if (!reason.trim() || reason.length > 4000)
                throw new Error("说明需要 1–4000 字符。");
            if (kind === "decision" && !decisionWinner)
                throw new Error("请先选择裁决结果。");
            const note = await request<Evidence>(`/api/tournaments/${id}/notes`, { method: "POST", headers: { Authorization: `Bearer ${token}` }, body: JSON.stringify({ kind, text: reason.trim() }) });
            if (!verifyEvidence(note.canonical_json, note.evidence_hash))
                throw new Error("说明内容与哈希不一致，无法提交。");
            await queryClient.invalidateQueries({ queryKey: ["course", `/api/tournaments/${id}/notes`] });
            const data = kind === "decision" ? encodeFunctionData({ abi: tournamentManagerAbi, functionName: "resolveDispute", args: [BigInt(id), decisionWinner as Address, note.evidence_hash] }) : encodeFunctionData({ abi: tournamentManagerAbi, functionName: kind === "dispute" ? "disputeResult" : "cancelTournament", args: [BigInt(id), note.evidence_hash] });
            await tx.prepare(`${kind === "dispute" ? "提交异议" : kind === "decision" ? `裁决异议：${decisionWinner === zeroAddress ? "取消并退款" : shortAddress(decisionWinner)}` : "取消赛事"} · #${id} · 说明 ${shortAddress(note.evidence_hash)}`, config.manager_address, data);
        }
        catch (err) {
            setError(errorMessage(err));
        }
        finally {
            setBusy(false);
        }
    }
    function downloadEvidence() {
        if (!evidence.data)
            return;
        const url = URL.createObjectURL(new Blob([evidence.data.canonical_json], { type: "application/json" }));
        const a = document.createElement("a");
        a.href = url;
        a.download = `tournament-${id}-evidence.json`;
        a.click();
        URL.revokeObjectURL(url);
        setViewedEvidence(evidence.data?.evidence_hash || "");
    }
    if (!validId)
        return <p role="alert">赛事编号无效。</p>;
    return <div className="space-y-6"><Link href="/tournaments" className="text-sm underline">← 返回赛事</Link><DeploymentNotice />{tournament.isLoading && <p>读取赛事中…</p>}{tournament.error && <p role="alert" className="text-red-500">{errorMessage(tournament.error)}</p>}{t && <>
    <div><div className="flex flex-wrap justify-between gap-3"><h1 className="text-3xl font-bold">#{t.id} {t.title}</h1><span className="border rounded-full px-3 py-1 text-sm self-start">{STATUS_LABEL[t.status]}</span></div><p className="text-muted-foreground mt-3 text-sm">组织者 {shortAddress(t.organizer)} · 独立审计员 {shortAddress(t.auditor)}</p>{t.metadata_uri && <p className="text-sm mt-2 break-all">规则：{/^https?:\/\//.test(t.metadata_uri) ? <a href={t.metadata_uri} target="_blank" rel="noreferrer" className="underline">{t.metadata_uri}</a> : t.metadata_uri}</p>}</div>
    <p className="text-xs text-muted-foreground">链时间：{now ? dateTime(now) : "读取中…"}（截止校验以链上区块时间为准）。{block.error && ` ${errorMessage(block.error)}`}</p><section className="grid sm:grid-cols-2 lg:grid-cols-4 gap-4">{[["奖金总额", `${money(t.prize_pool, config?.token_decimals)} DUSD`], ["报名人数", `${t.participants.length} / 2`], ["报名截止", dateTime(t.registration_deadline)], ["最终结算期限", dateTime(t.settlement_deadline)]].map(([label, value]) => <div key={label} className="border rounded-xl p-4"><p className="text-sm text-muted-foreground">{label}</p><p className="font-semibold mt-2 break-words">{value}</p></div>)}</section>
    <section className="border rounded-xl p-5 space-y-3"><h2 className="font-semibold">参赛与资金</h2><div className="space-y-2">{t.participant_agents.map((p, i) => <p key={p.address} className="text-sm break-all">{i === 0 ? "红方" : "蓝方"}：<span className="font-mono">{p.address}</span><span className="block text-xs text-muted-foreground">Agent 哈希：{p.agent_hash}</span></p>)}{!t.participants.length && <p className="text-muted-foreground">尚无参赛者。</p>}</div>{address && <p className="text-sm">钱包余额 {balance.data === undefined ? "—" : money(balance.data.toString(), config?.token_decimals)} DUSD · 我的赞助 {contribution.data === undefined ? "—" : money(contribution.data.toString(), config?.token_decimals)} DUSD</p>}
      {open && <><div className="flex flex-wrap gap-3 items-end"><label className="text-sm space-y-1">赞助金额（DUSD）<Input inputMode="decimal" value={amount} onChange={e => setAmount(e.target.value)} placeholder="例如 100.25"/></label><Button variant="outline" disabled={disabled || sameAddress(address, t.auditor)} onClick={() => fund(true)}>1. 估算授权</Button><Button disabled={disabled || sameAddress(address, t.auditor)} onClick={() => fund(false)}>2. 估算赞助</Button></div><p className="text-xs text-muted-foreground break-all">{sameAddress(address, t.auditor) && "指定审计员不能赞助自己审核的赛事。"}授权 spender 是奖金托管合约 {config?.escrow_address}。授权确认后再赞助；各步骤都由当前钱包签名。</p>{!participant && !sameAddress(address, t.organizer) && !sameAddress(address, t.auditor) && t.participants.length < 2 && <div className="flex flex-wrap gap-3 items-center"><select aria-label="选择报名 Agent" className="border rounded p-2 bg-background max-w-full" value={agentHash} onChange={e => setAgentHash(e.target.value)}><option value="">选择自己的 Agent 定义</option>{token && profiles.data?.agents.map(a => <option key={a.id} value={a.agent_hash}>{a.name} · {shortAddress(a.agent_hash)}</option>)}</select><Button disabled={disabled || !agentHash || !token} onClick={register}>估算报名</Button><Link href="/agents" className="underline text-sm">创建 Agent</Link></div>}{participant && <Button variant="outline" disabled={disabled} onClick={() => idAction("withdrawParticipant", "退出报名")}>估算退出报名</Button>}</>}
      {prize.data !== undefined && prize.data > BigInt(0) && <Button disabled={disabled} onClick={() => idAction("claimPrize", "领取冠军奖金")}>领取 {money(prize.data.toString(), config?.token_decimals)} DUSD 奖金（估算）</Button>}{refund.data !== undefined && refund.data > BigInt(0) && <Button disabled={disabled} onClick={() => idAction("claimRefund", "领取赞助退款")}>领取 {money(refund.data.toString(), config?.token_decimals)} DUSD 退款（估算）</Button>}
      {t.status === "finalized" && <p className="text-emerald-600 break-all">冠军：{t.winner} · {t.prize_claimed ? "奖金已领取" : "等待冠军钱包领取奖金"}</p>}{t.status === "cancelled" && <p>赛事已取消。每位赞助者可领取自己尚未退回的赞助。</p>}
    </section>
    {organizer && <section className="border rounded-xl p-5 space-y-3"><h2 className="font-semibold">组织者操作</h2>{t.status === "open" && <><Button disabled={disabled || t.participants.length !== 2 || now < t.registration_deadline} onClick={() => idAction("startTournament", "开始赛事")}>估算开赛</Button><p className="text-xs text-muted-foreground">需要两位参赛者、非零奖金，且报名已截止。</p></>}{t.status === "running" && <><Button disabled={!token || busy || match.data?.status === "running" || match.data?.status === "completed"} onClick={run}>{busy ? "请求中…" : match.data?.status === "interrupted" ? "恢复对战执行" : "运行 Agent 对战"}</Button><p className="text-xs text-muted-foreground">后端只生成对战证据，结果提交仍需组织者钱包签名。</p>{match.data?.status === "completed" && <Button disabled={disabled || !verified || !evidenceRead || !match.data.winner_address} onClick={submit}>估算提交结果</Button>}{match.data?.outcome === "draw" && <p className="text-amber-600">本场平局，没有冠军。请等待到期取消，赞助者随后退款。</p>}</>}</section>}
    {match.data && <MatchReplay match={match.data} settlementStatus={t.status}/>}{match.error && !(match.error instanceof ApiError && match.error.status === 404) && <p role="alert">对战数据：{errorMessage(match.error)}</p>}
    {evidence.data && <section className="border rounded-xl p-5 space-y-3"><h2 className="font-semibold">完整对战证据</h2><p className="font-mono text-xs break-all">{evidence.data.evidence_hash}</p><p className={verified ? "text-emerald-600" : "text-red-500"}>{verified ? "下载内容哈希与对战证据哈希一致" : "证据哈希不一致，已禁止提交"}</p>{t.evidence_hash && <p className={submittedEvidenceMatches ? "text-emerald-600" : "text-red-500"}>{submittedEvidenceMatches ? "链上提交哈希也一致" : "链上哈希与当前证据不一致"}</p>}<Button variant="outline" onClick={downloadEvidence}>下载完整证据 JSON</Button><details onToggle={e => { if (e.currentTarget.open)
            setViewedEvidence(evidence.data?.evidence_hash || ""); }}><summary className="cursor-pointer underline">展开证据内容</summary><pre className="mt-3 p-3 bg-muted rounded max-h-80 overflow-auto text-xs whitespace-pre-wrap break-all">{evidence.data.canonical_json}</pre></details><p className="text-xs text-muted-foreground">哈希一致证明内容承诺一致；审计员仍需独立审核比赛与规则。</p></section>}
    {evidence.error && <p role="alert">证据加载失败：{errorMessage(evidence.error)}</p>}
    {(t.status === "result_submitted" || t.status === "disputed") && <section className="border rounded-xl p-5 space-y-3"><h2 className="font-semibold">结果审核</h2><p className="break-all">提名冠军 {t.proposed_winner || "—"} · 审核截止 {dateTime(t.review_deadline)}</p>{auditor && t.status === "result_submitted" && <Button disabled={disabled || !reviewComplete || !verified || !submittedEvidenceMatches || !evidenceRead} onClick={() => idAction("approveResult", "审核通过并结算")}>估算审核通过</Button>}<p className="text-xs text-muted-foreground">审计员需先查看证据，审核窗口结束后才可批准；窗口期间参赛者或赞助者可提出异议。</p></section>}
    {((t.status === "result_submitted" && now < t.review_deadline && (participant || (contribution.data || BigInt(0)) > BigInt(0))) || (auditor && t.status === "disputed") || (organizer && t.status === "open")) && <section className="border rounded-xl p-5 space-y-3"><h2 className="font-semibold">异议、裁决或取消说明</h2><textarea className="w-full min-h-28 border rounded p-3 bg-background" aria-label="证据说明" maxLength={4000} value={reason} onChange={e => setReason(e.target.value)} placeholder="写明依据、发现的问题和相关证据。说明先持久化为草稿，再由钱包提交哈希。"/>{t.status === "result_submitted" && <Button disabled={disabled || !token} onClick={() => noteAction("dispute")}>保存说明并估算异议</Button>}{organizer && t.status === "open" && <Button variant="outline" disabled={disabled || !token} onClick={() => noteAction("cancellation")}>保存理由并估算取消</Button>}{auditor && t.status === "disputed" && <><select aria-label="异议裁决结果" className="border rounded p-2 bg-background block" value={decisionWinner} onChange={e => setDecisionWinner(e.target.value)}><option value="">选择裁决结果</option>{t.participants.map(p => <option key={p} value={p}>冠军：{shortAddress(p)}</option>)}<option value={zeroAddress}>否决结果，取消赛事并退款</option></select><Button disabled={disabled || !token || !reviewComplete || !evidenceRead || !verified || !submittedEvidenceMatches} onClick={() => noteAction("decision")}>保存裁决并估算提交</Button></>}</section>}
    {expired && <section className="border rounded-xl p-5 space-y-3"><h2 className="font-semibold">到期保护</h2><p className="text-sm">赛事超过期限，或报名截止后人数不足。任意钱包可以触发取消，赞助者自行领取退款。</p><Button variant="outline" disabled={disabled} onClick={() => idAction("expireTournament", "到期取消并启用退款")}>估算到期取消</Button></section>}
    {error && <p role="alert" className="text-red-500">{error}</p>}<TransactionStatus tx={tx}/>
    <section className="border rounded-xl p-5 space-y-3"><h2 className="font-semibold">审核说明草稿</h2><p className="text-xs text-muted-foreground">以下说明保存在后端。只有对应哈希出现在下方已确认合约事件时，才能视为已提交异议、取消或裁决。</p>{notes.data?.notes.map(note => <details key={note.evidence_hash} className="border rounded p-3"><summary className="cursor-pointer text-sm">{note.kind} · {shortAddress(note.author_address)} · {dateTime(note.created_at)}</summary><p className="whitespace-pre-wrap text-sm mt-2">{note.text}</p><p className="font-mono text-xs break-all mt-2">{note.evidence_hash}</p></details>)}{notes.error && <p role="alert">{errorMessage(notes.error)}</p>}</section>
    <section className="border rounded-xl p-5"><h2 className="font-semibold mb-3">赛事事件记录</h2>{events.error && <p role="alert">{errorMessage(events.error)}</p>}<EventTable events={events.data?.events || []}/><Link href="/audit" className="underline text-sm block mt-3">查看全部分页流水</Link></section>
  </>}</div>;
}
