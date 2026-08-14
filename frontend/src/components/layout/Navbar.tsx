"use client";

import Link from "next/link";
import { WalletConnect } from "./WalletConnect";
import { LoginButton } from "@/components/auth/LoginButton";
import { Swords, Coins, Gem, Wallet } from "lucide-react";
import { useAuth } from "@/lib/auth";
import { formatUnits } from "viem";

export function Navbar() {
  const { user, isAuthenticated } = useAuth();

  return (
    <nav className="border-b bg-background">
      <div className="container mx-auto flex h-16 items-center justify-between px-4">
        {/* Logo + Nav Links */}
        <div className="flex items-center gap-6">
          <Link href="/" className="flex items-center gap-2 font-bold text-lg">
            <Swords className="h-6 w-6" />
            <span>Agent Arena</span>
          </Link>
          <div className="flex items-center gap-4 text-sm">
            <Link
              href="/"
              className="text-muted-foreground hover:text-foreground transition-colors"
            >
              对局列表
            </Link>
            <Link
              href="/agents"
              className="text-muted-foreground hover:text-foreground transition-colors"
            >
              Agent
            </Link>
            <Link
              href="/marketplace"
              className="text-muted-foreground hover:text-foreground transition-colors"
            >
              市场
            </Link>
            <Link
              href="/pve"
              className="text-muted-foreground hover:text-foreground transition-colors"
            >
              PVE
            </Link>
            <Link
              href="/history"
              className="text-muted-foreground hover:text-foreground transition-colors"
            >
              对局记录
            </Link>
            <Link
              href="/wallet"
              className="text-muted-foreground hover:text-foreground transition-colors"
            >
              我的钱包
            </Link>
            <Link
              href="/profile"
              className="text-muted-foreground hover:text-foreground transition-colors"
            >
              个人中心
            </Link>
          </div>
        </div>

        {/* Balance + Login + Wallet Connect */}
        <div className="flex items-center gap-3">
          {isAuthenticated && user && (
            <Link
              href="/wallet"
              className="flex items-center gap-2 rounded-lg border bg-muted/50 px-3 py-1.5 text-sm hover:bg-muted transition-colors"
              title="我的钱包"
            >
              <Wallet className="h-4 w-4 text-muted-foreground" />
              <span className="flex items-center gap-1" title="USDC（链上）">
                <Coins className="h-3.5 w-3.5 text-yellow-500" />
                <span className="font-mono">
                  {user.usdc_balance_raw
                    ? formatUnits(BigInt(user.usdc_balance_raw), 6)
                    : user.usdc_balance != null
                    ? user.usdc_balance.toLocaleString()
                    : "—"}
                </span>
              </span>
              <span className="text-muted-foreground">|</span>
              <span className="flex items-center gap-1" title="AC（链上余额）">
                <Gem className="h-3.5 w-3.5 text-purple-500" />
                <span className="font-mono">
                  {(user.ac_on_chain_balance ?? user.ac_balance ?? 0).toLocaleString()}
                </span>
              </span>
            </Link>
          )}
          <LoginButton />
          <WalletConnect />
        </div>
      </div>
    </nav>
  );
}
