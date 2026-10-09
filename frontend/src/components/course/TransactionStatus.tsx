"use client";
import { formatEther } from "viem";
import { useCourseConfig } from "@/hooks/useCourse";
import type { useWalletTransaction } from "@/hooks/useWalletTransaction";
import { explorer, shortAddress } from "@/lib/course";
import { Button } from "@/components/ui/button";
export function TransactionStatus({ tx }: {
    tx: ReturnType<typeof useWalletTransaction>;
}) {
    const { data: config } = useCourseConfig();
    const url = tx.hash && config ? explorer(config.chain_id, tx.hash) : null;
    return <div className="rounded-lg border p-4 space-y-3 text-sm" aria-live="polite">
    {tx.stage === "idle" && <p className="text-muted-foreground">先选择一项操作，查看 gas 估算，然后在钱包中签名。</p>}
    {tx.stage === "estimating" && <p>正在模拟交易并估算 gas…</p>}
    {tx.prepared && <div><p className="font-medium">{tx.prepared.label}</p><p>预估网络费：{formatEther(tx.prepared.fee)} ETH</p><p className="text-xs text-muted-foreground">基于当前 gas price；实际费用以回执为准，gas limit 留 20% 余量。</p><p className="text-xs mt-2" title={tx.prepared.to}>发送目标：{shortAddress(tx.prepared.to)}</p><p className="text-xs text-muted-foreground">操作参数在估算时固定；修改表单后请重新估算。</p></div>}
    {tx.stage === "ready" && <Button onClick={tx.send} disabled={!tx.canWrite}>钱包签名发送</Button>}
    {tx.stage === "signing" && <p>等待钱包批准；拒绝签名不会发送交易。</p>}
    {tx.stage === "confirming" && <p>交易已发送，等待链上确认…</p>}
    {tx.stage === "success" && <p className="text-emerald-600">交易确认成功。审计索引可能晚于回执，请等待确认深度。</p>}
    {tx.stage === "unknown" && <div><p>尚未确认最终状态。为避免重复交易，暂时锁定新操作。</p><Button variant="outline" size="sm" onClick={tx.checkReceipt}>检查回执</Button></div>}
    {tx.hash && <p className="break-all font-mono text-xs">{url ? <a href={url} target="_blank" rel="noreferrer" className="underline">{tx.hash}</a> : tx.hash}</p>}
    {tx.error && <p role="alert" className="text-red-500 break-words">{tx.error}</p>}
  </div>;
}
