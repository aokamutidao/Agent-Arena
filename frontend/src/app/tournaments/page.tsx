"use client";
import Link from "next/link";
import { useState } from "react";
import { encodeFunctionData, isAddress, keccak256, toBytes, type Address } from "viem";
import { usePublicClient } from "wagmi";
import { useAuth } from "@/lib/auth";
import { tournamentManagerAbi } from "@/lib/course-abi";
import { errorMessage, money, STATUS_LABEL, dateTime, type Tournament } from "@/lib/course";
import { useCourseConfig, useCourseResource } from "@/hooks/useCourse";
import { useWalletTransaction } from "@/hooks/useWalletTransaction";
import { DeploymentNotice } from "@/components/course/DeploymentNotice";
import { TransactionStatus } from "@/components/course/TransactionStatus";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
export default function TournamentsPage() {
    const { user } = useAuth();
    const { data: config } = useCourseConfig();
    const publicClient = usePublicClient({ chainId: config?.chain_id });
    const [offset, setOffset] = useState(0);
    const list = useCourseResource<{
        tournaments: Tournament[];
        total: number;
    }>(`/api/tournaments?limit=20&offset=${offset}`);
    const tx = useWalletTransaction();
    const [error, setError] = useState("");
    const [title, setTitle] = useState("");
    const [uri, setURI] = useState("");
    const [auditor, setAuditor] = useState("");
    const [registrationMinutes, setRegistrationMinutes] = useState("5");
    const [settlementHours, setSettlementHours] = useState("24");
    const [reviewSeconds, setReviewSeconds] = useState("60");
    const [roleWallet, setRoleWallet] = useState("");
    const [role, setRole] = useState("ORGANIZER_ROLE");
    const organizer = user?.roles.includes("organizer");
    const admin = user?.roles.includes("admin");
    async function create() {
        setError("");
        try {
            if (!config || !publicClient || !isAddress(auditor))
                throw new Error("请输入有效审计员钱包。");
            if (!title.trim() || toBytes(title.trim()).length > 120)
                throw new Error("赛事名称需要 1–120 个 UTF-8 字节（汉字通常占 3 字节）。");
            if (toBytes(uri.trim()).length > 512)
                throw new Error("规则链接最多 512 个 UTF-8 字节。");
            if (!/^\d+$/.test(registrationMinutes) || !/^\d+$/.test(settlementHours) || !/^\d+$/.test(reviewSeconds))
                throw new Error("时间参数必须为正整数。");
            if (![registrationMinutes, settlementHours, reviewSeconds].every(value => Number.isSafeInteger(Number(value))))
                throw new Error("时间数值超出安全范围。");
            const currentBlock = await publicClient.getBlock();
            const registration = currentBlock.timestamp + BigInt(registrationMinutes) * BigInt(60);
            const settlement = currentBlock.timestamp + BigInt(settlementHours) * BigInt(3600);
            const review = BigInt(reviewSeconds);
            if (Number(registrationMinutes) < 1 || review < BigInt(60) || review > BigInt(2592000) || settlement <= registration + review)
                throw new Error("结算期限必须晚于报名截止与审核窗口；审核时间为 60 秒至 30 天。");
            await tx.prepare(`创建 ${title.trim()} · 审计员 ${auditor} · 报名截止 ${dateTime(Number(registration))}`, config.manager_address, encodeFunctionData({ abi: tournamentManagerAbi, functionName: "createTournament", args: [title.trim(), uri.trim(), auditor as Address, registration, settlement, review] }));
        }
        catch (err) {
            setError(errorMessage(err));
        }
    }
    async function changeRole(revoke: boolean) {
        setError("");
        try {
            if (!config || !isAddress(roleWallet))
                throw new Error("请输入有效钱包地址。");
            await tx.prepare(`${revoke ? "撤销角色" : "授予角色"} ${role} · ${roleWallet}`, config.manager_address, encodeFunctionData({ abi: tournamentManagerAbi, functionName: revoke ? "revokeRole" : "grantRole", args: [keccak256(toBytes(role)), roleWallet as Address] }));
        }
        catch (err) {
            setError(errorMessage(err));
        }
    }
    return <div className="space-y-6"><div><h1 className="text-3xl font-bold">赛事</h1><p className="mt-2 text-muted-foreground">两名参赛者、公开赞助奖金、独立审计员审核。报名费为零。</p></div><DeploymentNotice />
    {organizer && <section className="border rounded-xl p-5 space-y-4"><h2 className="font-semibold text-lg">创建赛事</h2><div className="grid md:grid-cols-2 gap-4"><label className="space-y-1 text-sm">赛事名称<Input value={title} onChange={e => setTitle(e.target.value)} placeholder="Agent 开发者挑战赛" maxLength={120}/></label><label className="space-y-1 text-sm">指定审计员钱包<Input value={auditor} onChange={e => setAuditor(e.target.value)} placeholder="0x…"/></label><label className="space-y-1 text-sm">规则链接（可选）<Input value={uri} onChange={e => setURI(e.target.value)} placeholder="https://… 或 ipfs://…"/></label><label className="space-y-1 text-sm">距报名截止（分钟）<Input type="number" min="1" value={registrationMinutes} onChange={e => setRegistrationMinutes(e.target.value)}/></label><label className="space-y-1 text-sm">距最终结算期限（小时）<Input type="number" min="1" value={settlementHours} onChange={e => setSettlementHours(e.target.value)}/></label><label className="space-y-1 text-sm">结果审核窗口（秒）<Input type="number" min="60" value={reviewSeconds} onChange={e => setReviewSeconds(e.target.value)}/></label></div><p className="text-sm text-muted-foreground">审计员需已拥有 auditor 角色，且不能是组织者或参赛者。</p><Button onClick={create} disabled={!tx.canWrite || tx.locked}>模拟并估算创建费用</Button></section>}
    {!organizer && <p className="text-sm text-muted-foreground">创建赛事需要 organizer 角色。请先签名登录；角色由链上权限决定。</p>}
    {admin && <section className="border rounded-xl p-5 space-y-3"><h2 className="font-semibold">管理员角色管理</h2><div className="flex flex-wrap gap-3"><select aria-label="角色" className="border rounded bg-background p-2" value={role} onChange={e => setRole(e.target.value)}><option value="ORGANIZER_ROLE">组织者</option><option value="AUDITOR_ROLE">审计员</option></select><Input className="max-w-md" aria-label="授权钱包地址" value={roleWallet} onChange={e => setRoleWallet(e.target.value)} placeholder="目标钱包 0x…"/><Button disabled={!tx.canWrite || tx.locked} onClick={() => changeRole(false)}>估算授权</Button><Button variant="outline" disabled={!tx.canWrite || tx.locked} onClick={() => changeRole(true)}>估算撤权</Button></div></section>}
    {error && <p role="alert" className="text-red-500">{error}</p>}<TransactionStatus tx={tx}/>
    <section className="space-y-3">{list.isLoading && <p>正在读取赛事…</p>}{list.error && <p role="alert">{errorMessage(list.error)}</p>}{list.data && !list.data.tournaments.length && <p className="text-muted-foreground">暂无赛事记录。</p>}{list.data?.tournaments.map(t => <Link className="block border rounded-xl p-5 hover:border-emerald-500" href={`/tournaments/${t.id}`} key={t.id}><div className="flex justify-between flex-wrap gap-3"><h2 className="font-semibold">#{t.id} {t.title}</h2><span className="text-sm">{STATUS_LABEL[t.status]}</span></div><div className="flex flex-wrap gap-6 text-sm text-muted-foreground mt-3"><span>奖金 {money(t.prize_pool, config?.token_decimals)} DUSD</span><span>参赛 {t.participants.length} / 2</span><span>报名截止 {dateTime(t.registration_deadline)}</span></div></Link>)}</section><div className="flex gap-3"><Button variant="outline" disabled={offset === 0} onClick={() => setOffset(Math.max(0, offset - 20))}>上一页</Button><Button variant="outline" disabled={offset + 20 >= (list.data?.total || 0)} onClick={() => setOffset(offset + 20)}>下一页</Button></div>
  </div>;
}
