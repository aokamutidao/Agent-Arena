"use client";

import { useEffect, useId, useRef, useState } from "react";
import {
  Check,
  CircleCheck,
  Clock3,
  ExternalLink,
  LoaderCircle,
  ShieldCheck,
  TriangleAlert,
  Wallet,
  X,
} from "lucide-react";
import { formatEther } from "viem";
import { useCourseConfig } from "@/hooks/useCourse";
import type { useWalletTransaction } from "@/hooks/useWalletTransaction";
import { explorer, shortAddress } from "@/lib/course";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";

type Transaction = ReturnType<typeof useWalletTransaction>;

const headings: Record<Transaction["stage"], string> = {
  idle: "交易操作",
  estimating: "正在检查这笔操作",
  ready: "请确认本次操作",
  signing: "请在钱包中签名",
  confirming: "已发送，等待链上确认",
  success: "交易已确认成功",
  error: "本次操作未成功",
  unknown: "交易状态需要核查",
};

const descriptions: Record<Transaction["stage"], string> = {
  idle: "",
  estimating: "正在模拟交易并估算网络费用，请稍候。",
  ready: "核对操作与费用后，点击下方按钮，再到钱包中确认。",
  signing: "钱包会展示交易详情。批准后才会发送，拒绝签名不会付款。",
  confirming: "交易正在等待确认，请勿重复提交同一操作。",
  success: "链上操作已完成。审计流水可能需要稍后才能同步。",
  error: "查看下方原因，处理后重新选择操作并估算。",
  unknown: "尚未确认最终结果。先检查回执，避免重复付款。",
};

export function TransactionStatus({ tx }: { tx: Transaction }) {
  const { data: config } = useCourseConfig();
  const dialog = useRef<HTMLDialogElement>(null);
  const primaryButton = useRef<HTMLButtonElement>(null);
  const previousStage = useRef<Transaction["stage"]>("idle");
  const sending = useRef(false);
  const [expanded, setExpanded] = useState(false);
  const [dismissed, setDismissed] = useState(false);
  const [checking, setChecking] = useState(false);
  const titleId = useId();
  const descriptionId = useId();
  const url = tx.hash && config ? explorer(config.chain_id, tx.hash) : null;
  const waiting = ["estimating", "signing", "confirming", "unknown"].includes(
    tx.stage,
  );
  const step =
    tx.stage === "estimating" || tx.stage === "ready"
      ? 0
      : tx.stage === "signing" || (tx.stage === "error" && !tx.hash)
        ? 1
        : 2;

  useEffect(() => {
    // Each new prepare opens the review; later receipt updates never undo a user's close.
    if (
      (tx.stage === "estimating" && previousStage.current !== "estimating") ||
      (previousStage.current === "idle" && tx.stage !== "idle")
    ) {
      setDismissed(false);
      setExpanded(true);
    }
    if (tx.stage === "idle") setExpanded(false);
    previousStage.current = tx.stage;
  }, [tx.stage]);

  useEffect(() => {
    const element = dialog.current;
    if (!element) return;
    if (expanded && tx.stage !== "idle" && !element.open) element.showModal();
    if ((!expanded || tx.stage === "idle") && element.open) element.close();
  }, [expanded, tx.stage]);

  useEffect(() => {
    if (expanded && tx.stage === "ready") primaryButton.current?.focus();
  }, [expanded, tx.stage]);

  function close() {
    setExpanded(false);
    if (tx.stage === "success" || tx.stage === "error") setDismissed(true);
    // Closing presentation must never reset the transaction hook or cancel receipt polling.
  }

  async function sendOnce() {
    if (sending.current || tx.stage !== "ready" || !tx.canWrite) return;
    sending.current = true;
    try {
      await tx.send();
    } finally {
      sending.current = false;
    }
  }

  async function checkReceipt() {
    if (checking) return;
    setChecking(true);
    try {
      await tx.checkReceipt();
    } finally {
      setChecking(false);
    }
  }

  if (tx.stage === "idle") return null;

  return (
    <>
      {!expanded && !dismissed && (
        <div
          className={cn(
            "fixed bottom-4 left-4 right-4 z-40 rounded-xl border bg-white p-4 shadow-lg sm:left-auto sm:w-96",
            tx.stage === "unknown" || tx.stage === "error"
              ? "border-amber-300"
              : "border-emerald-200",
          )}
          role="status"
        >
          <div className="flex items-start gap-3">
            {waiting ? (
              <Clock3
                className="mt-0.5 h-5 w-5 shrink-0 text-amber-600"
                aria-hidden="true"
              />
            ) : (
              <ShieldCheck
                className="mt-0.5 h-5 w-5 shrink-0 text-emerald-600"
                aria-hidden="true"
              />
            )}
            <div className="min-w-0 flex-1">
              <p className="text-sm font-semibold text-slate-900">
                {headings[tx.stage]}
              </p>
              <p className="mt-1 truncate text-xs text-slate-500">
                {tx.prepared?.label || "当前交易"}
              </p>
            </div>
          </div>
          <Button
            type="button"
            size="sm"
            className="mt-3 w-full"
            onClick={() => setExpanded(true)}
          >
            {tx.stage === "ready"
              ? "继续确认交易"
              : tx.stage === "unknown"
                ? "查看并检查回执"
                : "查看交易进度"}
          </Button>
        </div>
      )}

      <dialog
        ref={dialog}
        aria-labelledby={titleId}
        aria-describedby={descriptionId}
        onCancel={close}
        onClose={close}
        className="course-transaction-dialog m-auto max-h-[calc(100dvh-2rem)] w-[calc(100%-2rem)] max-w-lg overflow-y-auto rounded-2xl border border-slate-200 bg-white p-0 text-slate-900 shadow-2xl"
      >
        <div className="flex items-start justify-between gap-3 border-b border-slate-100 px-5 py-5 sm:px-6">
          <div className="flex items-center gap-3">
            <span
              className={cn(
                "flex h-11 w-11 shrink-0 items-center justify-center rounded-xl",
                tx.stage === "error"
                  ? "bg-red-50 text-red-600"
                  : tx.stage === "unknown"
                    ? "bg-amber-50 text-amber-700"
                    : "bg-emerald-50 text-emerald-600",
              )}
            >
              {["estimating", "signing", "confirming"].includes(tx.stage) ? (
                <LoaderCircle
                  className="h-6 w-6 animate-spin"
                  aria-hidden="true"
                />
              ) : tx.stage === "success" ? (
                <CircleCheck className="h-6 w-6" aria-hidden="true" />
              ) : tx.stage === "error" || tx.stage === "unknown" ? (
                <TriangleAlert className="h-6 w-6" aria-hidden="true" />
              ) : (
                <ShieldCheck className="h-6 w-6" aria-hidden="true" />
              )}
            </span>
            <div>
              <p className="text-xs font-medium text-slate-500">钱包交易</p>
              <h2
                id={titleId}
                className="mt-1 text-lg font-semibold tracking-tight"
              >
                {headings[tx.stage]}
              </h2>
            </div>
          </div>
          <button
            type="button"
            onClick={close}
            aria-label="关闭交易窗口"
            className="rounded-lg p-1.5 text-slate-400 hover:bg-slate-100 hover:text-slate-700"
          >
            <X className="h-5 w-5" aria-hidden="true" />
          </button>
        </div>

        <div className="space-y-5 px-5 py-5 sm:px-6">
          <ol aria-label="交易步骤" className="flex gap-2">
            {["检查与估算", "钱包签名", "链上确认"].map((label, index) => (
              <li
                key={label}
                aria-current={index === step ? "step" : undefined}
                className="flex flex-1 flex-col items-center gap-2 text-center"
              >
                <span
                  className={cn(
                    "flex h-7 w-7 items-center justify-center rounded-full text-xs font-semibold",
                    index < step || tx.stage === "success"
                      ? "bg-emerald-600 text-white"
                      : index === step
                        ? "bg-emerald-50 text-emerald-700 ring-1 ring-emerald-200"
                        : "bg-slate-100 text-slate-400",
                  )}
                >
                  {index < step || tx.stage === "success" ? (
                    <Check className="h-4 w-4" aria-hidden="true" />
                  ) : (
                    index + 1
                  )}
                </span>
                <span
                  className={cn(
                    "text-xs",
                    index === step
                      ? "font-semibold text-slate-800"
                      : "text-slate-500",
                  )}
                >
                  {label}
                </span>
              </li>
            ))}
          </ol>

          <p
            id={descriptionId}
            className="text-sm leading-6 text-slate-600"
            role="status"
          >
            {descriptions[tx.stage]}
          </p>

          {tx.prepared && (
            <div className="rounded-xl border border-slate-200 bg-slate-50 p-4">
              <p className="text-xs text-slate-500">本次操作</p>
              <p className="mt-1 break-words text-base font-semibold text-slate-900">
                {tx.prepared.label}
              </p>
              <div className="mt-4 border-t border-slate-200 pt-3">
                <p className="text-xs text-slate-500">预估网络费</p>
                <p className="mt-1 break-all font-mono text-sm font-medium text-slate-800">
                  {formatEther(tx.prepared.fee)} ETH
                </p>
                <p className="mt-1 text-xs leading-5 text-slate-500">
                  这是网络费用，不包含操作中的代币金额。实际费用以交易回执为准。
                </p>
              </div>
            </div>
          )}

          {tx.error && tx.stage !== "success" && (
            <p
              role="alert"
              className="break-words rounded-xl border border-red-200 bg-red-50 p-3 text-sm leading-6 text-red-800"
            >
              {tx.error}
            </p>
          )}
          {tx.stage === "ready" && !tx.canWrite && (
            <p
              role="alert"
              className="rounded-lg bg-amber-50 p-3 text-sm text-amber-800"
            >
              钱包或网络状态已变化，请恢复连接并重新估算。
            </p>
          )}

          {tx.hash && (
            <div className="flex flex-wrap items-center justify-between gap-2 text-xs text-slate-500">
              <span>
                交易编号：
                <span className="font-mono" title={tx.hash}>
                  {shortAddress(tx.hash)}
                </span>
              </span>
              {url && (
                <a
                  href={url}
                  target="_blank"
                  rel="noreferrer"
                  className="inline-flex items-center gap-1 rounded font-medium text-emerald-700 hover:underline"
                >
                  查看区块浏览器
                  <ExternalLink className="h-3.5 w-3.5" aria-hidden="true" />
                </a>
              )}
            </div>
          )}

          {(tx.prepared || tx.hash) && (
            <details className="rounded-lg border border-slate-200 px-3 py-2 text-xs">
              <summary className="cursor-pointer font-medium text-slate-600">
                交易详情
              </summary>
              <dl className="mt-3 space-y-2 text-slate-500">
                {tx.prepared && (
                  <>
                    <div>
                      <dt>签名账户</dt>
                      <dd className="mt-1 break-all font-mono text-slate-700">
                        {tx.prepared.account}
                      </dd>
                    </div>
                    <div>
                      <dt>目标合约</dt>
                      <dd className="mt-1 break-all font-mono text-slate-700">
                        {tx.prepared.to}
                      </dd>
                    </div>
                    <div>
                      <dt>链 ID / Gas 上限</dt>
                      <dd className="mt-1 font-mono text-slate-700">
                        {tx.prepared.chainId} / {tx.prepared.gas.toString()}
                      </dd>
                    </div>
                    <p className="leading-5">
                      操作参数已在估算时固定。修改表单后，应重新估算再签名；Gas
                      上限预留 20% 余量。
                    </p>
                  </>
                )}
                {tx.hash && (
                  <div>
                    <dt>完整交易编号</dt>
                    <dd className="mt-1 break-all font-mono text-slate-700">
                      {tx.hash}
                    </dd>
                  </div>
                )}
              </dl>
            </details>
          )}
        </div>

        <div className="space-y-3 border-t border-slate-100 bg-slate-50/70 px-5 py-4 sm:px-6">
          {tx.stage === "ready" && (
            <Button
              ref={primaryButton}
              type="button"
              size="lg"
              onClick={() => void sendOnce()}
              disabled={!tx.canWrite}
              className="w-full"
            >
              <Wallet className="mr-1 h-4 w-4" aria-hidden="true" />
              确认并在钱包签名
            </Button>
          )}
          {tx.stage === "estimating" && (
            <Button type="button" size="lg" disabled className="w-full">
              正在检查与估算…
            </Button>
          )}
          {tx.stage === "signing" && (
            <Button type="button" size="lg" disabled className="w-full">
              等待钱包批准…
            </Button>
          )}
          {(tx.stage === "confirming" || tx.stage === "unknown") && (
            <Button
              type="button"
              size="lg"
              variant={tx.stage === "unknown" ? "default" : "outline"}
              onClick={() => void checkReceipt()}
              disabled={checking}
              className="w-full"
            >
              {checking && (
                <LoaderCircle
                  className="mr-1 h-4 w-4 animate-spin"
                  aria-hidden="true"
                />
              )}
              {checking ? "正在检查回执…" : "检查链上回执"}
            </Button>
          )}
          {tx.stage === "success" && (
            <Button type="button" size="lg" onClick={close} className="w-full">
              完成
            </Button>
          )}
          {tx.stage === "error" && (
            <Button
              type="button"
              size="lg"
              variant="outline"
              onClick={close}
              className="w-full"
            >
              返回并调整操作
            </Button>
          )}
          <div className="flex items-center justify-between gap-3 text-xs text-slate-500">
            <span>
              {tx.stage === "confirming" || tx.stage === "unknown"
                ? "关闭窗口不会取消已发送的交易。"
                : tx.stage === "signing"
                  ? "关闭此窗口后，钱包仍可能在等待签名。"
                  : tx.stage === "estimating" || tx.stage === "ready"
                    ? "确认之前不会发送资金交易。"
                    : "交易记录可在赛事和审计流水中查看。"}
            </span>
            {[
              "estimating",
              "ready",
              "signing",
              "confirming",
              "unknown",
            ].includes(tx.stage) && (
              <button
                type="button"
                onClick={close}
                className="shrink-0 rounded font-medium text-slate-600 underline underline-offset-2 hover:text-slate-900"
              >
                收起窗口
              </button>
            )}
          </div>
        </div>
      </dialog>
    </>
  );
}
