"use client";
import { useAccount, useSwitchChain } from "wagmi";
import { useCourseConfig } from "@/hooks/useCourse";
import { errorMessage, shortAddress } from "@/lib/course";
import { Button } from "@/components/ui/button";
export function DeploymentNotice() {
    const { data, error, isLoading } = useCourseConfig();
    const { address, chainId } = useAccount();
    const { switchChain, isPending, error: switchError } = useSwitchChain();
    if (isLoading)
        return <p className="text-sm text-muted-foreground">正在读取部署配置…</p>;
    if (error || !data)
        return <p role="alert" className="rounded-lg border border-red-500/40 p-4 text-sm">后端连接失败：{errorMessage(error)}。链上操作已禁用。</p>;
    return <div className="rounded-lg border p-4 text-sm space-y-2 bg-muted/30">{!data.deployment_configured ? <p className="text-amber-600">课程合约尚未配置部署，当前仅可查看。配置并部署后才可签名发送交易。</p> : <p>网络：{data.chain_id === 31337 ? "本地 Anvil" : "Sepolia"} · 合约：<span className="font-mono">{shortAddress(data.manager_address)}</span> · 确认数：{data.confirmations}</p>}{data.simulation_enabled && <p className="text-amber-600">对战执行器允许确定性模拟；每场证据会标明实际执行模式。</p>}{address && chainId !== data.chain_id && <div className="flex flex-wrap items-center gap-3"><span className="text-amber-600">钱包网络与部署不一致。</span><Button size="sm" disabled={isPending} onClick={() => switchChain({ chainId: data.chain_id as 31337 | 11155111 })}>{isPending ? "切换中…" : "切换网络"}</Button></div>}{switchError && <p role="alert">{errorMessage(switchError)}</p>}</div>;
}
