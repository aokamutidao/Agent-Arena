"use client";

import dynamic from "next/dynamic";
import type { ReactNode } from "react";
import { I18nProvider, useI18n } from "@/lib/i18n";

function LoadingWallet() {
  const { l } = useI18n();
  return (
    <div
      className="min-h-screen flex items-center justify-center px-4 text-sm text-muted-foreground"
      role="status"
    >
      {l(
        "正在初始化钱包与赛事资金界面…",
        "Loading wallet and tournament dashboard…",
      )}
    </div>
  );
}

// WalletConnect creates browser storage while constructing its connectors.
// Load the wallet tree in the browser so SSR never touches IndexedDB.
const WalletProviders = dynamic(() => import("./wallet-providers"), {
  ssr: false,
  loading: LoadingWallet,
});

export function Providers({ children }: { children: ReactNode }) {
  return (
    <I18nProvider>
      <WalletProviders>{children}</WalletProviders>
    </I18nProvider>
  );
}
