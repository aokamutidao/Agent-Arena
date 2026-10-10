"use client";

import Link from "next/link";
import type { ReactNode } from "react";
import { Navbar } from "./Navbar";
import { useI18n } from "@/lib/i18n";

export function AppShell({ children }: { children: ReactNode }) {
  const { l } = useI18n();
  return (
    <div className="flex min-h-screen flex-col">
      <a
        href="#main-content"
        className="fixed left-4 top-4 z-50 -translate-y-24 rounded-lg bg-emerald-700 px-4 py-2 text-sm font-medium text-white focus:translate-y-0"
      >
        {l("跳转到主要内容", "Skip to main content")}
      </a>
      <Navbar />
      <main
        id="main-content"
        tabIndex={-1}
        className="mx-auto w-full max-w-7xl flex-1 px-4 py-6 outline-none sm:px-6 sm:py-8 lg:px-8"
      >
        {children}
      </main>
      <footer className="border-t border-slate-200 bg-white">
        <div className="mx-auto flex max-w-7xl flex-wrap items-center justify-between gap-2 px-4 py-4 text-xs text-slate-500 sm:px-6 lg:px-8">
          <span>
            {l(
              "Agent Arena · 奖金托管与赛事结算",
              "Agent Arena · Prize escrow & tournament settlement",
            )}
          </span>
          <Link
            href="/guide"
            className="rounded text-slate-600 hover:text-emerald-700"
          >
            {l("了解操作流程", "How it works")}
          </Link>
        </div>
      </footer>
    </div>
  );
}
