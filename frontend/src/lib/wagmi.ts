"use client";

import { getDefaultConfig } from "@rainbow-me/rainbowkit";
import { sepolia } from "wagmi/chains";

export const config = getDefaultConfig({
  appName: "Agent Arena",
  projectId: "f9798ce4e27e38aa273066336b8632f2",
  chains: [sepolia],
  ssr: true,
});
