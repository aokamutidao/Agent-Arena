"use client";

import Link from "next/link";
import { WalletConnect } from "./WalletConnect";
import { LoginButton } from "@/components/auth/LoginButton";
import { LanguageSwitcher } from "./LanguageSwitcher";
import { Swords, Coins, Gem, Wallet } from "lucide-react";
import { useAuth } from "@/lib/auth";
import { useI18n } from "@/lib/i18n";
import { formatUnits } from "viem";

export function Navbar() {
  const { user, isAuthenticated } = useAuth();
  const { t } = useI18n();

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
              {t("nav.agents")}
            </Link>
            <Link
              href="/marketplace"
              className="text-muted-foreground hover:text-foreground transition-colors"
            >
              {t("nav.marketplace")}
            </Link>
            <Link
              href="/pve"
              className="text-muted-foreground hover:text-foreground transition-colors"
            >
              {t("nav.pve")}
            </Link>
            <Link
              href="/history"
              className="text-muted-foreground hover:text-foreground transition-colors"
            >
              {t("nav.history")}
            </Link>
            <Link
              href="/wallet"
              className="text-muted-foreground hover:text-foreground transition-colors"
            >
              {t("nav.wallet")}
            </Link>
            <Link
              href="/profile"
              className="text-muted-foreground hover:text-foreground transition-colors"
            >
              {t("nav.profile")}
            </Link>
          </div>
        </div>

        {/* Balance + Login + Wallet Connect + Language Switcher */}
        <div className="flex items-center gap-3">
          {isAuthenticated && user && (
            <Link
              href="/wallet"
              className="flex items-center gap-2 rounded-lg border bg-muted/50 px-3 py-1.5 text-sm hover:bg-muted transition-colors"
              title={t("nav.wallet")}
            >
              <Wallet className="h-4 w-4 text-muted-foreground" />
              <span className="flex items-center gap-1" title="USDC">
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
              <span className="flex items-center gap-1" title="AC">
                <Gem className="h-3.5 w-3.5 text-purple-500" />
                <span className="font-mono">
                  {(user.ac_on_chain_balance ?? user.ac_balance ?? 0).toLocaleString()}
                </span>
              </span>
            </Link>
          )}
          <LoginButton />
          <WalletConnect />
          <LanguageSwitcher />
        </div>
      </div>
    </nav>
  );
}
