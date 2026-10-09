"use client";
import { useEffect, useRef, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { useAccount, usePublicClient, useWalletClient } from "wagmi";
import { type Address, type Hex } from "viem";
import { errorMessage, receiptVerdict } from "@/lib/course";
import { useCourseConfig } from "./useCourse";
type Stage = "idle" | "estimating" | "ready" | "signing" | "confirming" | "success" | "error" | "unknown";
interface Prepared {
    label: string;
    to: Address;
    data: Hex;
    account: Address;
    chainId: number;
    gas: bigint;
    fee: bigint;
}
export function useWalletTransaction() {
    const { address, chainId } = useAccount();
    const { data: config } = useCourseConfig();
    const publicClient = usePublicClient({ chainId: config?.chain_id });
    const { data: walletClient } = useWalletClient();
    const queryClient = useQueryClient();
    const [stage, setStage] = useState<Stage>("idle");
    const [prepared, setPrepared] = useState<Prepared | null>(null);
    const [hash, setHash] = useState<Hex | null>(null);
    const [error, setError] = useState("");
    const currentIdentity = `${address || ""}:${chainId || ""}`;
    const identity = useRef(currentIdentity);
    const replacement = useRef("");
    identity.current = currentIdentity;
    const canWrite = !!address && !!config?.deployment_configured && chainId === config.chain_id;
    const locked = ["estimating", "signing", "confirming", "unknown"].includes(stage);
    useEffect(() => {
        setPrepared(null);
        setStage((previous) => ["signing", "confirming", "unknown"].includes(previous) ? previous : "idle");
    }, [address, chainId]);
    async function prepare(label: string, to: Address, data: Hex) {
        if (locked)
            return;
        setError("");
        setHash(null);
        setPrepared(null);
        if (!canWrite || !address || !config || !publicClient) {
            setError("连接钱包并切换至已部署网络后，才能发送交易。");
            setStage("error");
            return;
        }
        const requestIdentity = currentIdentity;
        setStage("estimating");
        try {
            if (!(await publicClient.getBytecode({ address: to })))
                throw new Error("目标地址没有合约代码，请检查部署配置。");
            const gas = await publicClient.estimateGas({ account: address, to, data });
            const price = await publicClient.getGasPrice();
            if (identity.current !== requestIdentity)
                throw new Error("钱包或网络已变化，请重新估算。");
            setPrepared({ label, to, data, account: address, chainId: config.chain_id, gas: gas * BigInt(120) / BigInt(100), fee: gas * price });
            setStage("ready");
        }
        catch (err) {
            setError(errorMessage(err));
            setStage("error");
        }
    }
    async function send() {
        if (!prepared || !canWrite || !walletClient || !publicClient || !config)
            return;
        if (prepared.account !== address || prepared.chainId !== chainId) {
            setError("钱包或网络已变化，请重新估算。");
            setStage("error");
            return;
        }
        setStage("signing");
        setError("");
        const requestIdentity = currentIdentity;
        let sent: Hex | null = null;
        let receiptKnown = false;
        replacement.current = "";
        try {
            // Re-estimate immediately before signing; the displayed fee remains an estimate.
            const gas = await publicClient.estimateGas({ account: prepared.account, to: prepared.to, data: prepared.data });
            if (identity.current !== requestIdentity || await walletClient.getChainId() !== prepared.chainId)
                throw new Error("钱包或网络已变化，请重新估算。");
            sent = await walletClient.sendTransaction({ account: prepared.account, chain: walletClient.chain, to: prepared.to, data: prepared.data, gas: gas * BigInt(120) / BigInt(100) });
            setHash(sent);
            setStage("confirming");
            const receipt = await publicClient.waitForTransactionReceipt({ hash: sent, confirmations: Math.max(1, config.confirmations), timeout: 120000, onReplaced: ({ transaction, reason }) => { setHash(transaction.hash); replacement.current = reason; } });
            receiptKnown = true;
            if (receiptVerdict(receipt.status, replacement.current) === "replaced")
                throw new Error("原交易已被取消或替换为其他操作，请核查新交易回执后重试。");
            if (receipt.status !== "success")
                throw new Error("交易已上链但执行回滚；请查看交易回执。");
            setStage("success");
            await queryClient.invalidateQueries({ queryKey: ["course"] });
            await queryClient.invalidateQueries({ queryKey: ["readContract"] });
        }
        catch (err) {
            setError(errorMessage(err));
            setStage(sent && !receiptKnown ? "unknown" : "error");
        }
    }
    async function checkReceipt() {
        if (!hash || !publicClient || !config)
            return;
        try {
            const receipt = await publicClient.getTransactionReceipt({ hash });
            const head = await publicClient.getBlockNumber();
            if (head - receipt.blockNumber + BigInt(1) < BigInt(Math.max(1, config.confirmations)))
                throw new Error("确认数不足，请稍后继续检查。");
            if (receiptVerdict(receipt.status, replacement.current) === "replaced") {
                setStage("error");
                setError("原交易已被取消或替换为其他操作，请核查新交易回执。");
                return;
            }
            setStage(receipt.status === "success" ? "success" : "error");
            setError(receipt.status === "success" ? "" : "交易回执显示执行回滚。");
            await queryClient.invalidateQueries({ queryKey: ["course"] });
            await queryClient.invalidateQueries({ queryKey: ["readContract"] });
        }
        catch (err) {
            setError(`交易状态待确认：${errorMessage(err)}`);
        }
    }
    return { prepare, send, checkReceipt, prepared, stage, hash, error, canWrite, locked };
}
