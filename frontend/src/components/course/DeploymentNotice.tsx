"use client";

import Link from "next/link";
import {
  CheckCircle2,
  ChevronDown,
  LoaderCircle,
  RefreshCw,
  TriangleAlert,
} from "lucide-react";
import { useAccount, useSwitchChain } from "wagmi";
import { useCourseConfig, useCourseResource } from "@/hooks/useCourse";
import { errorMessage } from "@/lib/course";
import { cn } from "@/lib/utils";
import { useI18n } from "@/lib/i18n";
import { Button } from "@/components/ui/button";

interface Readiness {
  ready: boolean;
  indexed_block: number;
  chain_head: number;
}

export function DeploymentNotice() {
  const { lang, l } = useI18n();
  const { data, error, isLoading, isFetching, refetch } = useCourseConfig();
  const readiness = useCourseResource<Readiness>(
    "/ready",
    !!data?.deployment_configured,
  );
  const { address, chainId } = useAccount();
  const { switchChain, isPending, error: switchError } = useSwitchChain();
  const wrongNetwork = !!address && !!data && chainId !== data.chain_id;
  const networkLabel =
    data?.chain_id === 31337
      ? l("本地演示网络", "Local demo network")
      : data?.chain_id === 11155111
        ? l("Sepolia 测试网络", "Sepolia testnet")
        : l(
            `网络 ${data?.chain_id || "—"}`,
            `Network ${data?.chain_id || "—"}`,
          );
  const retry = () => {
    void refetch();
    if (data?.deployment_configured) void readiness.refetch();
  };

  if (isLoading)
    return (
      <div
        role="status"
        className="flex items-center gap-2 text-sm text-slate-500"
      >
        <LoaderCircle className="h-4 w-4 animate-spin" aria-hidden="true" />
        {l("正在连接赛事服务…", "Connecting to the tournament service…")}
      </div>
    );
  if (error || !data)
    return (
      <div
        role="alert"
        className="flex flex-wrap items-center justify-between gap-3 rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-800"
      >
        <div className="flex min-w-0 items-start gap-2">
          <TriangleAlert
            className="mt-0.5 h-4 w-4 shrink-0"
            aria-hidden="true"
          />
          <div>
            <p className="font-semibold">
              {l("赛事服务连接失败", "Tournament service unavailable")}
            </p>
            <p className="mt-1 break-words text-xs leading-5">
              {l(
                "检查服务是否启动，恢复连接后再继续操作。",
                "Check that the service is running, then reconnect before continuing.",
              )}
              {error ? ` ${errorMessage(error, lang)}` : ""}
            </p>
          </div>
        </div>
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={retry}
          disabled={isFetching}
        >
          <RefreshCw
            className={cn("mr-1.5 h-3.5 w-3.5", isFetching && "animate-spin")}
            aria-hidden="true"
          />
          {l("重新连接", "Reconnect")}
        </Button>
      </div>
    );

  const ready = !!readiness.data?.ready && !readiness.error;
  const syncing = data.deployment_configured && readiness.isLoading;
  return (
    <div className="space-y-2">
      <details className="group rounded-xl border border-slate-200 bg-white text-sm">
        <summary className="flex cursor-pointer list-none flex-wrap items-center gap-x-3 gap-y-2 px-4 py-3 [&::-webkit-details-marker]:hidden">
          <span
            className={cn(
              "inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-medium",
              !data.deployment_configured
                ? "bg-amber-50 text-amber-800"
                : ready
                  ? "bg-emerald-50 text-emerald-700"
                  : "bg-slate-100 text-slate-600",
            )}
          >
            {!data.deployment_configured ? (
              <TriangleAlert className="h-3.5 w-3.5" aria-hidden="true" />
            ) : ready ? (
              <CheckCircle2 className="h-3.5 w-3.5" aria-hidden="true" />
            ) : (
              <LoaderCircle
                className={cn("h-3.5 w-3.5", syncing && "animate-spin")}
                aria-hidden="true"
              />
            )}
            {!data.deployment_configured
              ? l("查看模式 · 尚未部署", "View only · Not deployed")
              : ready
                ? l("链数据已同步", "Chain data synced")
                : syncing
                  ? l("正在同步链数据", "Syncing chain data")
                  : l("链数据暂未就绪", "Chain data not ready")}
          </span>
          <span className="text-xs text-slate-500">{networkLabel}</span>
          <span className="ml-auto inline-flex items-center gap-1 text-xs text-slate-500">
            {l("网络详情", "Network details")}
            <ChevronDown
              className="h-4 w-4 transition-transform group-open:rotate-180"
              aria-hidden="true"
            />
          </span>
        </summary>
        <div className="space-y-3 border-t border-slate-100 px-4 py-4">
          {!data.deployment_configured && (
            <p className="text-amber-800">
              {l(
                "课程合约尚未配置，暂时无法发送资金交易。请按",
                "Tournament contracts are not configured, so fund transactions are unavailable. Follow the ",
              )}
              <Link
                href="/guide"
                className="ml-1 font-medium underline underline-offset-2"
              >
                {l("使用指南", "getting started guide")}
              </Link>
              {l("完成部署配置。", " to complete deployment setup.")}
            </p>
          )}
          <dl className="grid gap-2 text-xs sm:grid-cols-2">
            <div>
              <dt className="text-slate-500">
                {l("链 ID / 确认数", "Chain ID / Confirmations")}
              </dt>
              <dd className="mt-1 font-medium text-slate-700">
                {data.chain_id} / {data.confirmations}
              </dd>
            </div>
            {readiness.data && (
              <div>
                <dt className="text-slate-500">
                  {l("已索引 / 链头区块", "Indexed block / Chain head")}
                </dt>
                <dd className="mt-1 font-medium text-slate-700">
                  {readiness.data.indexed_block} / {readiness.data.chain_head}
                </dd>
              </div>
            )}
            {(
              [
                [l("赛事合约", "Tournament contract"), data.manager_address],
                [l("奖金托管", "Prize escrow"), data.escrow_address],
                [l("赛事注册", "Tournament registry"), data.registry_address],
                [l("结果审核", "Result review"), data.review_address],
                [l("结算代币", "Settlement token"), data.token_address],
              ] as const
            ).map(([label, value]) => (
              <div key={label} className="min-w-0">
                <dt className="text-slate-500">{label}</dt>
                <dd className="mt-1 break-all font-mono text-slate-700">
                  {value || l("未配置", "Not configured")}
                </dd>
              </div>
            ))}
          </dl>
          {data.simulation_enabled && (
            <p className="rounded-lg bg-amber-50 px-3 py-2 text-xs leading-5 text-amber-800">
              {l(
                "运行器允许确定性模拟。每场对战的实际模式以比赛记录和证据为准。",
                "Deterministic simulation is enabled. Check each match record and its evidence for the mode actually used.",
              )}
            </p>
          )}
        </div>
      </details>
      {data.deployment_configured && readiness.error && (
        <div
          role="alert"
          className="flex flex-wrap items-center justify-between gap-2 rounded-lg border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-800"
        >
          <span>
            {l(
              "链数据暂未就绪，请检查节点连接，资金流水可能尚未更新。",
              "Chain data is not ready. Check the node connection; financial records may not be up to date.",
            )}
          </span>
          <Button
            type="button"
            size="sm"
            variant="ghost"
            onClick={retry}
            disabled={readiness.isFetching}
            className="h-8 text-amber-800"
          >
            {l("重新检查", "Check again")}
          </Button>
        </div>
      )}
      {wrongNetwork && (
        <div
          role="alert"
          className="flex flex-wrap items-center justify-between gap-3 rounded-xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-900"
        >
          <span className="flex items-center gap-2">
            <TriangleAlert className="h-4 w-4 shrink-0" aria-hidden="true" />
            {l(
              `钱包网络不匹配，请切换到${networkLabel}后再操作。`,
              `Wallet network mismatch. Switch to ${networkLabel} before continuing.`,
            )}
          </span>
          <Button
            type="button"
            size="sm"
            disabled={isPending || ![31337, 11155111].includes(data.chain_id)}
            onClick={() =>
              switchChain({ chainId: data.chain_id as 31337 | 11155111 })
            }
          >
            {isPending
              ? l("切换中…", "Switching…")
              : l("切换到正确网络", "Switch to the correct network")}
          </Button>
        </div>
      )}
      {switchError && (
        <p
          role="alert"
          className="rounded-lg bg-red-50 px-3 py-2 text-xs text-red-700"
        >
          {l(
            `网络切换未完成：${errorMessage(switchError, lang)}。也可在钱包中手动选择${networkLabel}。`,
            `Network switch incomplete: ${errorMessage(switchError, lang)}. You can also select ${networkLabel} manually in your wallet.`,
          )}
        </p>
      )}
    </div>
  );
}
