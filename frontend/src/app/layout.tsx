import type { Metadata } from "next";
import Link from "next/link";
import localFont from "next/font/local";
import "./globals.css";
import { Providers } from "@/components/providers";
import { Navbar } from "@/components/layout/Navbar";

const geistSans = localFont({
  src: "./fonts/GeistVF.woff",
  variable: "--font-geist-sans",
  weight: "100 900",
});
const geistMono = localFont({
  src: "./fonts/GeistMonoVF.woff",
  variable: "--font-geist-mono",
  weight: "100 900",
});

export const metadata: Metadata = {
  title: "Agent Arena — 赛事奖金托管",
  description: "赛事奖金提前托管，比赛结果独立审核，领奖与退款全程可查。",
};

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="zh-CN">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <Providers>
          <div className="flex min-h-screen flex-col">
            <a
              href="#main-content"
              className="fixed left-4 top-4 z-50 -translate-y-24 rounded-lg bg-emerald-700 px-4 py-2 text-sm font-medium text-white focus:translate-y-0"
            >
              跳转到主要内容
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
                <span>Agent Arena · 奖金托管与赛事结算</span>
                <Link
                  href="/guide"
                  className="rounded text-slate-600 hover:text-emerald-700"
                >
                  了解操作流程
                </Link>
              </div>
            </footer>
          </div>
        </Providers>
      </body>
    </html>
  );
}
