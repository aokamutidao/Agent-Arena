"use client";
import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { useAuth } from "@/lib/auth";
import { errorMessage, request, type AgentProfile, dateTime } from "@/lib/course";
import { useAuthenticatedResource } from "@/hooks/useCourse";
import { DeploymentNotice } from "@/components/course/DeploymentNotice";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
export default function AgentsPage() {
    const { token } = useAuth();
    const profiles = useAuthenticatedResource<{
        agents: AgentProfile[];
        total: number;
    }>("/api/agents/my", token);
    const queryClient = useQueryClient();
    const [name, setName] = useState("");
    const [personality, setPersonality] = useState("");
    const [pending, setPending] = useState(false);
    const [error, setError] = useState("");
    async function create(e: React.FormEvent) {
        e.preventDefault();
        if (!token)
            return;
        setPending(true);
        setError("");
        try {
            await request<AgentProfile>("/api/agents", { method: "POST", headers: { Authorization: `Bearer ${token}` }, body: JSON.stringify({ name: name.trim(), personality: personality.trim() }) });
            setName("");
            setPersonality("");
            await queryClient.invalidateQueries({ queryKey: ["course", "/api/agents/my"] });
        }
        catch (err) {
            setError(errorMessage(err));
        }
        finally {
            setPending(false);
        }
    }
    return <div className="space-y-6"><h1 className="text-3xl font-bold">我的 Agent</h1><DeploymentNotice /><p className="text-muted-foreground">创建不可变的参赛定义，在赛事报名时选择它。名称、策略与所有者形成内容哈希，赛后证据保留完整定义。</p>{!token ? <p className="border rounded-xl p-5">请连接钱包并在顶部签名登录，然后创建你的 Agent。</p> : <><form className="border rounded-xl p-5 space-y-4" onSubmit={create}><h2 className="text-lg font-semibold">创建参赛定义</h2><label className="block text-sm space-y-1">Agent 名称（最多 50 字）<Input required maxLength={50} value={name} onChange={e => setName(e.target.value)}/></label><label className="block text-sm space-y-1">策略与性格（最多 500 字）<textarea className="w-full min-h-32 border rounded p-3 bg-background" required maxLength={500} value={personality} onChange={e => setPersonality(e.target.value)} placeholder="例如：优先保持安全距离，在敌方蓄力时移动，低血量时恢复。"/></label><p className="text-xs text-muted-foreground">定义保存后不可修改；调整策略请创建新定义。模型连接由平台统一配置。</p><Button type="submit" disabled={pending}>{pending ? "保存中…" : "保存 Agent 定义"}</Button></form>{error && <p role="alert" className="text-red-500">{error}</p>}{profiles.error && <p role="alert">{errorMessage(profiles.error)}</p>}{profiles.isLoading && <p>正在加载…</p>}<div className="grid md:grid-cols-2 gap-4">{profiles.data?.agents.map(agent => <article className="border rounded-xl p-5 space-y-3" key={agent.id}><h2 className="font-semibold">{agent.name}</h2><p className="text-sm whitespace-pre-wrap">{agent.personality}</p><p className="font-mono text-xs break-all">{agent.agent_hash}</p><p className="text-xs text-muted-foreground">创建于 {dateTime(agent.created_at)}</p></article>)}</div>{profiles.data?.agents.length === 0 && <p className="text-muted-foreground">还没有 Agent 定义。</p>}</>}</div>;
}
