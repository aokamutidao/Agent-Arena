"use client";
import Link from "next/link";
import { useState } from "react";
import { Menu, X, ShieldCheck } from "lucide-react";
import { WalletConnect } from "./WalletConnect";
import { LoginButton } from "@/components/auth/LoginButton";
const links = [["/", "资金概览"], ["/tournaments", "赛事"], ["/audit", "审计流水"], ["/agents", "我的 Agent"], ["/guide", "使用说明"]];
export function Navbar() {
    const [open, setOpen] = useState(false);
    return <header className="border-b bg-background"><div className="container mx-auto px-4 py-3"><div className="flex flex-wrap items-center justify-between gap-3"><Link href="/" className="flex items-center gap-2 font-bold"><ShieldCheck className="h-6 w-6 text-emerald-500"/>Agent Arena</Link><button aria-label={open ? "收起导航" : "打开导航"} aria-expanded={open} onClick={() => setOpen(!open)} className="md:hidden border rounded p-2">{open ? <X size={18}/> : <Menu size={18}/>}</button><nav className="hidden md:flex gap-5 text-sm">{links.map(([href, label]) => <Link key={href} href={href} className="hover:text-emerald-500">{label}</Link>)}</nav><div className="flex flex-wrap items-center gap-3"><LoginButton /><WalletConnect /></div></div>{open && <nav className="md:hidden flex flex-col gap-3 pt-4">{links.map(([href, label]) => <Link onClick={() => setOpen(false)} key={href} href={href}>{label}</Link>)}</nav>}</div></header>;
}
