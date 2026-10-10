"use client";

import Link from "next/link";
import { useI18n } from "@/lib/i18n";
import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import {
  ArrowRight,
  Bot,
  CheckCircle2,
  Fingerprint,
  Shield,
  Swords,
  Wallet,
} from "lucide-react";
import { useAuth } from "@/lib/auth";
import {
  errorMessage,
  request,
  type AgentProfile,
  dateTime,
} from "@/lib/course";
import { useAuthenticatedResource } from "@/hooks/useCourse";
import { DeploymentNotice } from "@/components/course/DeploymentNotice";
import { LoginButton } from "@/components/auth/LoginButton";
import { WalletConnect } from "@/components/layout/WalletConnect";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

export default function AgentsPage() {
  const { lang, l } = useI18n();
  const { token } = useAuth();
  const profiles = useAuthenticatedResource<{
    agents: AgentProfile[];
    total: number;
  }>("/api/agents/my", token);
  const queryClient = useQueryClient();
  const [name, setName] = useState("");
  const [personality, setPersonality] = useState("");
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<unknown>(null);
  const [savedName, setSavedName] = useState("");

  async function create(e: React.FormEvent) {
    e.preventDefault();
    if (!token || !name.trim() || !personality.trim()) return;
    setPending(true);
    setError(null);
    setSavedName("");
    try {
      await request<AgentProfile>("/api/agents", {
        method: "POST",
        headers: { Authorization: "Bearer " + token },
        body: JSON.stringify({
          name: name.trim(),
          personality: personality.trim(),
        }),
      });
      setSavedName(name.trim());
      setName("");
      setPersonality("");
      await queryClient.invalidateQueries({
        queryKey: ["course", "/api/agents/my"],
      });
    } catch (err) {
      setError(err);
    } finally {
      setPending(false);
    }
  }

  function chooseStrategy(strategy: string, suggestedName: string) {
    setPersonality(strategy);
    if (!name.trim()) setName(suggestedName);
    setSavedName("");
  }

  return (
    <div className="space-y-7">
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <p className="mb-2 text-sm font-semibold text-emerald-700">
            {l("参赛准备", "Prepare to compete")}
          </p>
          <h1 className="text-3xl font-bold tracking-tight text-slate-900">
            {l("我的 Agent", "My agents")}
          </h1>
          <p className="mt-3 text-base text-slate-600">
            {l(
              "给你的 Agent 起名字、选策略，保存后就能去赛事报名。",
              "Name your agent, choose a strategy, then save it to register for a tournament.",
            )}
          </p>
        </div>
        <Link
          href="/tournaments"
          className="inline-flex items-center gap-2 font-semibold text-emerald-700 hover:text-emerald-800"
        >
          {l("选择参赛赛事", "Choose a tournament")} <ArrowRight size={18} />
        </Link>
      </div>

      <div className="grid gap-3 sm:grid-cols-3">
        {[
          [
            "01",
            l("创建 Agent", "Create an agent"),
            l("名字与策略由你决定", "Choose the name and strategy"),
          ],
          [
            "02",
            l("报名赛事", "Register for a tournament"),
            l("选择自己保存的 Agent", "Select an agent you saved"),
          ],
          [
            "03",
            l("比赛与领奖", "Compete and claim"),
            l(
              "审核后由冠军领取奖金",
              "The winner claims the prize after review",
            ),
          ],
        ].map(([number, title, hint]) => (
          <div
            key={number}
            className="flex items-center gap-4 rounded-xl border border-slate-200 bg-white p-4"
          >
            <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-emerald-50 text-sm font-bold text-emerald-700">
              {number}
            </span>
            <div>
              <p className="font-semibold text-slate-900">{title}</p>
              <p className="mt-1 text-sm text-slate-500">{hint}</p>
            </div>
          </div>
        ))}
      </div>
      <DeploymentNotice />

      {!token ? (
        <section className="rounded-2xl border border-emerald-200 bg-white p-6 sm:p-9">
          <Wallet className="mb-4 h-9 w-9 text-emerald-600" />
          <h2 className="text-xl font-bold text-slate-900">
            {l("先连接钱包，再签名登录", "Connect your wallet, then sign in")}
          </h2>
          <p className="mt-3 max-w-xl text-slate-600">
            {l(
              "钱包用来确认 Agent 属于谁。登录签名不转账，也不支付网络费。",
              "Your wallet proves ownership of the agent. Signing in does not transfer funds or cost a network fee.",
            )}
          </p>
          <div className="mt-6 flex flex-wrap items-start gap-3">
            <WalletConnect />
            <LoginButton />
          </div>
          <Link
            href="/tournaments"
            className="mt-6 inline-flex items-center gap-2 text-sm font-semibold text-emerald-700"
          >
            {l("先去看比赛，无需钱包", "Browse tournaments without a wallet")}{" "}
            <ArrowRight size={16} />
          </Link>
        </section>
      ) : (
        <>
          {savedName && (
            <div
              role="status"
              className="flex flex-wrap items-center justify-between gap-3 rounded-xl border border-emerald-200 bg-emerald-50 p-4 text-emerald-900"
            >
              <p className="flex min-w-0 items-center gap-2">
                <CheckCircle2 className="h-5 w-5 shrink-0" />
                <span className="break-words">
                  <strong>{savedName}</strong>{" "}
                  {l(
                    "已保存，接下来去赛事报名。",
                    "has been saved. You can now register for a tournament.",
                  )}
                </span>
              </p>
              <Link
                href="/tournaments"
                className="inline-flex items-center gap-2 font-semibold"
              >
                {l("去报名", "Register")} <ArrowRight size={16} />
              </Link>
            </div>
          )}
          <div className="grid items-start gap-6 lg:grid-cols-2">
            <form
              className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm"
              onSubmit={create}
            >
              <div className="mb-6 flex items-center gap-3">
                <Bot className="h-6 w-6 text-emerald-600" />
                <h2 className="text-xl font-bold text-slate-900">
                  {l("创建参赛 Agent", "Create an agent to compete")}
                </h2>
              </div>
              <label className="block text-sm font-semibold text-slate-700">
                {l("Agent 名称", "Agent name")}
                <Input
                  required
                  maxLength={50}
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="mt-2 h-11 bg-white"
                  placeholder={l(
                    "例如：我的进攻 Agent",
                    "e.g. My Aggressive Agent",
                  )}
                />
                <span className="mt-1.5 block text-sm font-normal text-slate-500">
                  {l(
                    "最多 50 字，比赛里会显示这个名字。",
                    "Up to 50 characters. This name will appear in the match.",
                  )}
                </span>
              </label>
              <fieldset className="mt-6">
                <legend className="text-sm font-semibold text-slate-700">
                  {l("先选一个策略起点", "Choose a starting strategy")}
                </legend>
                <div className="mt-3 grid gap-3 sm:grid-cols-2">
                  <button
                    type="button"
                    aria-pressed={personality === "aggressive"}
                    onClick={() =>
                      chooseStrategy(
                        "aggressive",
                        l("我的进攻 Agent", "My Aggressive Agent"),
                      )
                    }
                    className={
                      "rounded-xl border p-4 text-left transition-colors " +
                      (personality === "aggressive"
                        ? "border-emerald-500 bg-emerald-50"
                        : "border-slate-200 hover:border-emerald-300")
                    }
                  >
                    <Swords className="mb-2 h-5 w-5 text-emerald-700" />
                    <span className="block font-semibold text-slate-900">
                      {l("进攻型", "Aggressive")}
                    </span>
                    <span className="mt-1 block text-sm text-slate-500">
                      {l(
                        "主动接近对手，优先进攻",
                        "Approach the opponent and prioritize attacks",
                      )}
                    </span>
                  </button>
                  <button
                    type="button"
                    aria-pressed={personality === "defensive"}
                    onClick={() =>
                      chooseStrategy(
                        "defensive",
                        l("我的防守 Agent", "My Defensive Agent"),
                      )
                    }
                    className={
                      "rounded-xl border p-4 text-left transition-colors " +
                      (personality === "defensive"
                        ? "border-emerald-500 bg-emerald-50"
                        : "border-slate-200 hover:border-emerald-300")
                    }
                  >
                    <Shield className="mb-2 h-5 w-5 text-emerald-700" />
                    <span className="block font-semibold text-slate-900">
                      {l("防守型", "Defensive")}
                    </span>
                    <span className="mt-1 block text-sm text-slate-500">
                      {l(
                        "关注距离与恢复，稳健应对",
                        "Manage distance and recovery for a steady response",
                      )}
                    </span>
                  </button>
                </div>
              </fieldset>
              <label className="mt-6 block text-sm font-semibold text-slate-700">
                {l("策略与性格", "Strategy and personality")}
                <textarea
                  className="mt-2 min-h-32 w-full rounded-lg border border-slate-300 bg-white p-3 font-normal focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
                  required
                  maxLength={500}
                  value={personality}
                  onChange={(e) => setPersonality(e.target.value)}
                  placeholder={l(
                    "填写策略，或先选择上方进攻型 / 防守型。",
                    "Write a strategy or start with Aggressive / Defensive above.",
                  )}
                />
                <span className="mt-1 block text-sm font-normal text-slate-500">
                  {l(
                    "最多 500 字。模拟演示使用 aggressive / defensive 等关键词；比赛页会标明是否实际使用 LLM。",
                    "Up to 500 characters. The simulation uses keywords such as aggressive / defensive; the match page states whether an LLM was used.",
                  )}
                </span>
              </label>
              {error !== null && (
                <p
                  role="alert"
                  className="mt-4 rounded-lg bg-red-50 p-3 text-sm text-red-700"
                >
                  {errorMessage(error, lang)}
                </p>
              )}
              <Button
                type="submit"
                size="lg"
                className="mt-6 w-full bg-emerald-600 text-white hover:bg-emerald-700"
                disabled={pending || !name.trim() || !personality.trim()}
              >
                {pending
                  ? l("正在保存…", "Saving…")
                  : l(
                      "保存 Agent，准备参赛",
                      "Save agent and prepare to compete",
                    )}
              </Button>
              <p className="mt-3 text-center text-sm text-slate-500">
                {l(
                  "保存不需要链上交易；报名时才在钱包确认。",
                  "Saving does not require an on-chain transaction. You confirm registration in your wallet later.",
                )}
              </p>
              <details className="mt-6 border-t border-slate-100 pt-4 text-sm text-slate-600">
                <summary className="cursor-pointer font-medium">
                  {l(
                    "为什么保存后不能修改策略？",
                    "Why can’t I edit a saved strategy?",
                  )}
                </summary>
                <p className="mt-3 leading-6">
                  {l(
                    "报名会记录这份定义，便于赛后核对。调整策略请创建新的 Agent；报名截止前可以退出后重新选择。",
                    "Registration records this definition for later verification. Create a new agent to change its strategy; before registration closes, you can withdraw and register again.",
                  )}
                </p>
              </details>
            </form>

            <section className="space-y-4">
              <div className="flex items-center justify-between gap-3">
                <h2 className="text-xl font-bold text-slate-900">
                  {l("已保存的 Agent", "Saved agents")}
                </h2>
                <span className="rounded-full bg-slate-200/60 px-3 py-1 text-sm text-slate-600">
                  {l(
                    (profiles.data?.total ?? "—") + " 个",
                    (profiles.data?.total ?? "—") + " agents",
                  )}
                </span>
              </div>
              <p className="text-sm text-slate-500">
                {l(
                  "这里只显示属于当前钱包的定义。去赛事详情选择它，即可报名。",
                  "Only definitions owned by your current wallet appear here. Select one on a tournament page to register.",
                )}
              </p>
              {profiles.error && (
                <p
                  role="alert"
                  className="rounded-xl border border-red-200 bg-white p-4 text-red-700"
                >
                  {errorMessage(profiles.error, lang)}
                </p>
              )}
              {profiles.isLoading && (
                <div
                  className="rounded-xl border border-slate-200 bg-white p-6 text-slate-500"
                  role="status"
                >
                  {l("正在读取你的 Agent…", "Loading your agents…")}
                </div>
              )}
              {profiles.data?.agents.map((agent) => (
                <article
                  className="rounded-2xl border border-slate-200 bg-white p-5"
                  key={agent.id}
                >
                  <div className="flex items-start gap-3">
                    <div className="rounded-xl bg-emerald-50 p-3 text-emerald-700">
                      <Bot size={23} />
                    </div>
                    <div className="min-w-0">
                      <h3 className="break-words text-lg font-semibold text-slate-900">
                        {agent.name}
                      </h3>
                      <p className="mt-1 text-sm text-slate-500">
                        {dateTime(agent.created_at, lang)}
                      </p>
                    </div>
                  </div>
                  <p className="mt-4 whitespace-pre-wrap break-words text-sm leading-6 text-slate-600">
                    {agent.personality}
                  </p>
                  <div className="mt-5 flex flex-wrap items-center justify-between gap-3">
                    <Link
                      href="/tournaments"
                      className="inline-flex items-center gap-2 text-sm font-semibold text-emerald-700"
                    >
                      {l("带它去参赛", "Compete with this agent")}{" "}
                      <ArrowRight size={16} />
                    </Link>
                    <details className="max-w-full text-sm text-slate-500">
                      <summary className="cursor-pointer">
                        {l("查看定义指纹", "View definition fingerprint")}
                      </summary>
                      <p className="mt-2 flex items-start gap-2 break-all font-mono text-xs">
                        <Fingerprint className="h-4 w-4 shrink-0" />
                        {agent.agent_hash}
                      </p>
                    </details>
                  </div>
                </article>
              ))}
              {profiles.data?.agents.length === 0 && (
                <div className="rounded-2xl border border-dashed border-slate-300 bg-white p-8 text-center">
                  <Bot className="mx-auto mb-3 h-8 w-8 text-slate-400" />
                  <h3 className="font-semibold text-slate-900">
                    {l(
                      "第一个 Agent，就从这里开始",
                      "Start with your first agent",
                    )}
                  </h3>
                  <p className="mt-2 text-sm text-slate-500">
                    {l(
                      "在创建表单里选好策略、填个名字，然后保存。",
                      "Choose a strategy in the form, enter a name, then save.",
                    )}
                  </p>
                </div>
              )}
            </section>
          </div>
        </>
      )}
    </div>
  );
}
