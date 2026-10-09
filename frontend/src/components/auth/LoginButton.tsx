"use client";
import { useState } from "react";
import { useAccount, useSignMessage } from "wagmi";
import { useAuth } from "@/lib/auth";
import { errorMessage, request } from "@/lib/course";
import { useCourseConfig } from "@/hooks/useCourse";
import { Button } from "@/components/ui/button";
export function LoginButton() {
    const { address, chainId } = useAccount();
    const { signMessageAsync } = useSignMessage();
    const { login, logout, user } = useAuth();
    const { data: config } = useCourseConfig();
    const [pending, setPending] = useState(false);
    const [error, setError] = useState("");
    async function handleLogin() {
        if (!address || !config || chainId !== config.chain_id)
            return;
        setPending(true);
        setError("");
        try {
            const challenge = await request<{
                nonce: string;
                message: string;
                expires_at: number;
            }>("/api/auth/nonce", { method: "POST", body: JSON.stringify({ address: address.toLowerCase() }) });
            const signature = await signMessageAsync({ message: challenge.message });
            await login(address, challenge.nonce, signature);
        }
        catch (err) {
            setError(errorMessage(err));
        }
        finally {
            setPending(false);
        }
    }
    if (!address)
        return null;
    return <div className="max-w-xs"><Button size="sm" variant="outline" disabled={pending || (!user && (!config || chainId !== config.chain_id))} onClick={user ? logout : handleLogin}>{user ? "退出登录" : pending ? "等待签名…" : "签名登录"}</Button>{user && <p className="text-xs text-muted-foreground mt-1">{user.roles.join(" / ")}</p>}{error && <p role="alert" className="text-xs text-red-500 mt-1 break-words">{error}</p>}</div>;
}
