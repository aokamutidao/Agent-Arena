"use client";

import { ConnectButton } from "@rainbow-me/rainbowkit";
import { ChevronDown, Wallet, TriangleAlert } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useCourseConfig } from "@/hooks/useCourse";
import { shortAddress } from "@/lib/course";
import { cn } from "@/lib/utils";
import { useI18n } from "@/lib/i18n";

export function WalletConnect({ className }: { className?: string }) {
  const { lang, l } = useI18n();
  const { data: config } = useCourseConfig();
  return (
    <ConnectButton.Custom>
      {({
        account,
        chain,
        mounted,
        authenticationStatus,
        openAccountModal,
        openChainModal,
        openConnectModal,
      }) => {
        const ready = mounted && authenticationStatus !== "loading";
        const connected =
          ready &&
          account &&
          chain &&
          (!authenticationStatus || authenticationStatus === "authenticated");
        const wrongNetwork =
          !!chain &&
          (chain.unsupported || (!!config && chain.id !== config.chain_id));
        const networkLabel =
          chain?.id === 31337
            ? l("本地链", "Local chain")
            : chain?.id === 11155111
              ? "Sepolia"
              : l("选择网络", "Select network");
        return (
          <div
            className={cn("flex flex-wrap items-center gap-2", className)}
            aria-hidden={!ready || undefined}
            style={
              !ready
                ? { opacity: 0, pointerEvents: "none", userSelect: "none" }
                : undefined
            }
          >
            {!connected ? (
              <Button type="button" size="sm" onClick={openConnectModal}>
                <Wallet className="mr-2 h-4 w-4" aria-hidden="true" />
                {l("连接钱包", "Connect wallet")}
              </Button>
            ) : (
              <>
                <Button
                  type="button"
                  size="sm"
                  variant="ghost"
                  onClick={openChainModal}
                  className={
                    wrongNetwork
                      ? "bg-amber-50 text-amber-800 hover:bg-amber-100 hover:text-amber-900"
                      : "text-slate-600"
                  }
                  aria-label={
                    wrongNetwork
                      ? l(
                          "钱包网络不匹配，选择正确网络",
                          "Wallet network mismatch; select the correct network",
                        )
                      : l("更换钱包网络", "Change wallet network")
                  }
                >
                  {wrongNetwork && (
                    <TriangleAlert
                      className="mr-1.5 h-4 w-4"
                      aria-hidden="true"
                    />
                  )}
                  {wrongNetwork
                    ? l("切换网络", "Switch network")
                    : networkLabel}
                  <ChevronDown
                    className="ml-1 h-3.5 w-3.5"
                    aria-hidden="true"
                  />
                </Button>
                <Button
                  type="button"
                  size="sm"
                  variant="outline"
                  onClick={openAccountModal}
                  title={l(
                    "查看账户、切换账户或断开钱包",
                    "View or switch accounts, or disconnect your wallet",
                  )}
                  aria-label={
                    lang === "zh"
                      ? `切换账户，当前 ${account.displayName}`
                      : `Switch account; current account ${account.displayName}`
                  }
                  className="bg-white"
                >
                  <Wallet
                    className="mr-2 h-4 w-4 text-slate-500"
                    aria-hidden="true"
                  />
                  <span className="font-mono text-xs">
                    {account.ensName || shortAddress(account.address)}
                  </span>
                  <ChevronDown
                    className="ml-2 h-3.5 w-3.5 text-slate-500"
                    aria-hidden="true"
                  />
                </Button>
              </>
            )}
          </div>
        );
      }}
    </ConnectButton.Custom>
  );
}
