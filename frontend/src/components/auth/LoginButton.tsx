"use client";

import { useState } from "react";
import { LoaderCircle, LogOut, UserRound } from "lucide-react";
import { useAccount, useSignMessage } from "wagmi";
import { useAuth } from "@/lib/auth";
import { errorMessage, request } from "@/lib/course";
import { globalRoleLabel } from "@/lib/course-presentation";
import { cn } from "@/lib/utils";
import { useI18n } from "@/lib/i18n";
import { useCourseConfig } from "@/hooks/useCourse";
import { Button } from "@/components/ui/button";

export function LoginButton({ className }: { className?: string }) {
  const { lang, l } = useI18n();
  const { address, chainId } = useAccount();
  const { signMessageAsync } = useSignMessage();
  const { login, logout, user } = useAuth();
  const { data: config } = useCourseConfig();
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<unknown>(null);
  async function handleLogin() {
    if (!address || !config || chainId !== config.chain_id || pending) return;
    setPending(true);
    setError(null);
    try {
      const challenge = await request<{
        nonce: string;
        message: string;
        expires_at: number;
      }>("/api/auth/nonce", {
        method: "POST",
        body: JSON.stringify({ address: address.toLowerCase() }),
      });
      const signature = await signMessageAsync({ message: challenge.message });
      await login(address, challenge.nonce, signature);
    } catch (err) {
      setError(err);
    } finally {
      setPending(false);
    }
  }
  if (!address) return null;
  return (
    <div className={cn("relative", className)}>
      {user ? (
        <div className="flex items-center gap-1.5">
          <span className="inline-flex items-center gap-1.5 rounded-full bg-emerald-50 px-2.5 py-1.5 text-xs font-medium text-emerald-700">
            <UserRound className="h-3.5 w-3.5" aria-hidden="true" />
            {globalRoleLabel(user.roles, lang)}
          </span>
          <Button
            type="button"
            size="sm"
            variant="ghost"
            onClick={() => {
              logout();
              setError(null);
            }}
            className="px-2 text-slate-500"
            title={l("退出签名登录", "Sign out")}
            aria-label={l("退出登录", "Sign out")}
          >
            <LogOut className="h-4 w-4" aria-hidden="true" />
          </Button>
        </div>
      ) : (
        <Button
          type="button"
          size="sm"
          variant="outline"
          disabled={pending || !config || chainId !== config.chain_id}
          onClick={handleLogin}
          title={
            !config
              ? l(
                  "等待服务连接后登录",
                  "Wait for the service to connect before signing in",
                )
              : chainId !== config.chain_id
                ? l(
                    "请先切换至课程网络",
                    "Switch to the tournament network first",
                  )
                : l(
                    "签名只用于验证身份，不发送资金交易",
                    "This signature verifies your identity and does not send a transaction",
                  )
          }
        >
          {pending && (
            <LoaderCircle
              className="mr-2 h-4 w-4 animate-spin"
              aria-hidden="true"
            />
          )}
          {pending
            ? l("等待签名…", "Waiting for signature…")
            : l("签名登录", "Sign in")}
        </Button>
      )}
      {error !== null && (
        <p
          role="alert"
          className="absolute right-0 top-full z-40 mt-2 w-72 max-w-[calc(100vw-2rem)] rounded-lg border border-red-200 bg-red-50 p-3 text-xs leading-5 text-red-700 shadow-sm"
        >
          {l("登录失败：", "Sign-in failed: ")}
          {errorMessage(error, lang)}
        </p>
      )}
    </div>
  );
}
