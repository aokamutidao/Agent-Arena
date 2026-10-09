"use client";

import dynamic from "next/dynamic";
import type { ReactNode } from "react";

// WalletConnect creates browser storage while constructing its connectors.
// Load the wallet tree in the browser so SSR never touches IndexedDB.
const WalletProviders = dynamic(() => import("./wallet-providers"), {
  ssr: false,
  loading: () => (
    <div className="min-h-screen flex items-center justify-center px-4 text-sm text-muted-foreground" role="status">
      正在初始化钱包与赛事资金界面…
    </div>
  ),
});

export function Providers({ children }: { children: ReactNode }) {
  return <WalletProviders>{children}</WalletProviders>;
}
