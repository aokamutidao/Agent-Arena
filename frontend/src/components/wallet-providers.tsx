"use client";

import React from "react";
import { WagmiProvider } from "wagmi";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { RainbowKitProvider } from "@rainbow-me/rainbowkit";
import { config } from "@/lib/wagmi";
import { AuthProvider } from "@/lib/auth";
import { useI18n } from "@/lib/i18n";

import "@rainbow-me/rainbowkit/styles.css";

const queryClient = new QueryClient();

interface ProvidersProps {
  children: React.ReactNode;
}

export default function WalletProviders({ children }: ProvidersProps) {
  const { lang } = useI18n();
  return (
    <WagmiProvider
      config={config}
      initialState={{ ...config.state, status: "reconnecting" }}
    >
      <QueryClientProvider client={queryClient}>
        <RainbowKitProvider locale={lang === "zh" ? "zh-CN" : "en-US"}>
          <AuthProvider>{children}</AuthProvider>
        </RainbowKitProvider>
      </QueryClientProvider>
    </WagmiProvider>
  );
}
