"use client";

import { getDefaultConfig } from "@rainbow-me/rainbowkit";
import { sepolia } from "wagmi/chains";
import { defineChain, http } from "viem";

export const localChain = defineChain({
  id: 31337,
  name: "Local Anvil",
  nativeCurrency: { name: "Ether", symbol: "ETH", decimals: 18 },
  rpcUrls: { default: { http: [process.env.NEXT_PUBLIC_LOCAL_RPC_URL || "http://127.0.0.1:8545"] } },
});

export const config = getDefaultConfig({
  appName: "Agent Arena",
  projectId: "f9798ce4e27e38aa273066336b8632f2",
  chains: [localChain, sepolia],
  transports: {
    [localChain.id]: http(process.env.NEXT_PUBLIC_LOCAL_RPC_URL || "http://127.0.0.1:8545"),
    [sepolia.id]: http(process.env.NEXT_PUBLIC_SEPOLIA_RPC_URL || undefined),
  },
  ssr: true,
});
