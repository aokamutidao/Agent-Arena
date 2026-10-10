"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useState } from "react";
import { Menu, X, ShieldCheck } from "lucide-react";
import { WalletConnect } from "./WalletConnect";
import { LoginButton } from "@/components/auth/LoginButton";
import { useCourseConfig } from "@/hooks/useCourse";
import { cn } from "@/lib/utils";

const links = [
  ["/", "资金概览"],
  ["/tournaments", "奖金赛事"],
  ["/audit", "审计流水"],
  ["/agents", "我的 Agent"],
  ["/guide", "使用指南"],
] as const;

export function Navbar() {
  const [open, setOpen] = useState(false);
  const pathname = usePathname();
  const { data: config, error } = useCourseConfig();
  const networkLabel = error
    ? "服务未连接"
    : !config
      ? "连接服务中"
      : !config.deployment_configured
        ? "尚未部署"
        : config.chain_id === 31337
          ? "本地演示"
          : config.chain_id === 11155111
            ? "Sepolia 测试网"
            : `网络 ${config.chain_id}`;
  const isActive = (href: string) =>
    href === "/" ? pathname === "/" : pathname.startsWith(href);

  return (
    <header className="sticky top-0 z-30 border-b border-slate-200/80 bg-white/95 backdrop-blur">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex flex-wrap items-center justify-between gap-x-4 gap-y-3 py-3">
          <div className="flex min-w-0 items-center gap-3">
            <Link
              href="/"
              className="flex shrink-0 items-center gap-2.5 rounded-lg"
              aria-label="Agent Arena 赛事奖金托管首页"
            >
              <span className="flex h-10 w-10 items-center justify-center rounded-xl bg-emerald-600 text-white">
                <ShieldCheck className="h-6 w-6" aria-hidden="true" />
              </span>
              <span>
                <span className="block text-base font-bold tracking-tight text-slate-900">
                  Agent Arena
                </span>
                <span className="block text-xs text-slate-500">
                  赛事奖金托管
                </span>
              </span>
            </Link>
            <span
              className={cn(
                "hidden items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-medium sm:inline-flex lg:hidden xl:inline-flex",
                error
                  ? "bg-red-50 text-red-700"
                  : config && !config.deployment_configured
                    ? "bg-amber-50 text-amber-800"
                    : "bg-slate-100 text-slate-600",
              )}
            >
              <span
                className={cn(
                  "h-1.5 w-1.5 rounded-full",
                  error
                    ? "bg-red-500"
                    : config && !config.deployment_configured
                      ? "bg-amber-500"
                      : "bg-slate-400",
                )}
                aria-hidden="true"
              />
              {networkLabel}
            </span>
          </div>
          <nav
            aria-label="主导航"
            className="hidden items-center gap-1 lg:flex"
          >
            {links.map(([href, label]) => (
              <Link
                key={href}
                href={href}
                aria-current={isActive(href) ? "page" : undefined}
                className={cn(
                  "rounded-lg px-3 py-2 text-sm font-medium transition-colors",
                  isActive(href)
                    ? "bg-emerald-50 text-emerald-700"
                    : "text-slate-600 hover:bg-slate-50 hover:text-slate-900",
                )}
              >
                {label}
              </Link>
            ))}
          </nav>
          <button
            type="button"
            aria-label={open ? "收起导航" : "打开导航"}
            aria-expanded={open}
            aria-controls="mobile-navigation"
            onClick={() => setOpen(!open)}
            className="rounded-lg border border-slate-200 p-2 text-slate-600 hover:bg-slate-50 lg:hidden"
          >
            {open ? (
              <X size={20} aria-hidden="true" />
            ) : (
              <Menu size={20} aria-hidden="true" />
            )}
          </button>
          <div className="flex w-full flex-wrap items-center justify-end gap-2 border-t border-slate-100 pt-3 lg:w-auto lg:border-0 lg:pt-0">
            <LoginButton />
            <WalletConnect />
          </div>
        </div>
        {open && (
          <nav
            id="mobile-navigation"
            aria-label="移动端导航"
            className="grid grid-cols-2 gap-1 border-t border-slate-100 py-3 lg:hidden"
          >
            {links.map(([href, label]) => (
              <Link
                onClick={() => setOpen(false)}
                key={href}
                href={href}
                aria-current={isActive(href) ? "page" : undefined}
                className={cn(
                  "rounded-lg px-3 py-2.5 text-sm font-medium",
                  isActive(href)
                    ? "bg-emerald-50 text-emerald-700"
                    : "text-slate-600 hover:bg-slate-50",
                )}
              >
                {label}
              </Link>
            ))}
          </nav>
        )}
      </div>
    </header>
  );
}
