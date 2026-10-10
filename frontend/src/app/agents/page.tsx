"use client";

import Link from "next/link";
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
  const { token } = useAuth();
  const profiles = useAuthenticatedResource<{
    agents: AgentProfile[];
    total: number;
  }>("/api/agents/my", token);
  const queryClient = useQueryClient();
  const [name, setName] = useState("");
  const [personality, setPersonality] = useState("");
  const [pending, setPending] = useState(false);
  const [error, setError] = useState("");
  const [savedName, setSavedName] = useState("");

  async function create(e: React.FormEvent) {
    e.preventDefault();
    if (!token || !name.trim() || !personality.trim()) return;
    setPending(true);
    setError("");
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
      setError(errorMessage(err));
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
            参赛准备
          </p>
          <h1 className="text-3xl font-bold tracking-tight text-slate-900">
            我的 Agent
          </h1>
          <p className="mt-3 text-base text-slate-600">
            给你的 Agent 起名字、选策略，保存后就能去赛事报名。
          </p>
        </div>
        <Link
          href="/tournaments"
          className="inline-flex items-center gap-2 font-semibold text-emerald-700 hover:text-emerald-800"
        >
          选择参赛赛事 <ArrowRight size={18} />
        </Link>
      </div>

      <div className="grid gap-3 sm:grid-cols-3">
        {[
          ["01", "创建 Agent", "名字与策略由你决定"],
          ["02", "报名赛事", "选择自己保存的 Agent"],
          ["03", "比赛与领奖", "审核后由冠军领取奖金"],
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
            先连接钱包，再签名登录
          </h2>
          <p className="mt-3 max-w-xl text-slate-600">
            钱包用来确认 Agent 属于谁。登录签名不转账，也不支付网络费。
          </p>
          <div className="mt-6 flex flex-wrap items-start gap-3">
            <WalletConnect />
            <LoginButton />
          </div>
          <Link
            href="/tournaments"
            className="mt-6 inline-flex items-center gap-2 text-sm font-semibold text-emerald-700"
          >
            先去看比赛，无需钱包 <ArrowRight size={16} />
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
                  <strong>{savedName}</strong> 已保存，接下来去赛事报名。
                </span>
              </p>
              <Link
                href="/tournaments"
                className="inline-flex items-center gap-2 font-semibold"
              >
                去报名 <ArrowRight size={16} />
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
                  创建参赛 Agent
                </h2>
              </div>
              <label className="block text-sm font-semibold text-slate-700">
                Agent 名称
                <Input
                  required
                  maxLength={50}
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="mt-2 h-11 bg-white"
                  placeholder="例如：我的进攻 Agent"
                />
                <span className="mt-1.5 block text-sm font-normal text-slate-500">
                  最多 50 字，比赛里会显示这个名字。
                </span>
              </label>
              <fieldset className="mt-6">
                <legend className="text-sm font-semibold text-slate-700">
                  先选一个策略起点
                </legend>
                <div className="mt-3 grid gap-3 sm:grid-cols-2">
                  <button
                    type="button"
                    aria-pressed={personality === "aggressive"}
                    onClick={() =>
                      chooseStrategy("aggressive", "我的进攻 Agent")
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
                      进攻型
                    </span>
                    <span className="mt-1 block text-sm text-slate-500">
                      主动接近对手，优先进攻
                    </span>
                  </button>
                  <button
                    type="button"
                    aria-pressed={personality === "defensive"}
                    onClick={() =>
                      chooseStrategy("defensive", "我的防守 Agent")
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
                      防守型
                    </span>
                    <span className="mt-1 block text-sm text-slate-500">
                      关注距离与恢复，稳健应对
                    </span>
                  </button>
                </div>
              </fieldset>
              <label className="mt-6 block text-sm font-semibold text-slate-700">
                策略与性格
                <textarea
                  className="mt-2 min-h-32 w-full rounded-lg border border-slate-300 bg-white p-3 font-normal focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
                  required
                  maxLength={500}
                  value={personality}
                  onChange={(e) => setPersonality(e.target.value)}
                  placeholder="填写策略，或先选择上方进攻型 / 防守型。"
                />
                <span className="mt-1 block text-sm font-normal text-slate-500">
                  最多 500 字。模拟演示使用 aggressive / defensive
                  等关键词；比赛页会标明是否实际使用 LLM。
                </span>
              </label>
              {error && (
                <p
                  role="alert"
                  className="mt-4 rounded-lg bg-red-50 p-3 text-sm text-red-700"
                >
                  {error}
                </p>
              )}
              <Button
                type="submit"
                size="lg"
                className="mt-6 w-full bg-emerald-600 text-white hover:bg-emerald-700"
                disabled={pending || !name.trim() || !personality.trim()}
              >
                {pending ? "正在保存…" : "保存 Agent，准备参赛"}
              </Button>
              <p className="mt-3 text-center text-sm text-slate-500">
                保存不需要链上交易；报名时才在钱包确认。
              </p>
              <details className="mt-6 border-t border-slate-100 pt-4 text-sm text-slate-600">
                <summary className="cursor-pointer font-medium">
                  为什么保存后不能修改策略？
                </summary>
                <p className="mt-3 leading-6">
                  报名会记录这份定义，便于赛后核对。调整策略请创建新的
                  Agent；报名截止前可以退出后重新选择。
                </p>
              </details>
            </form>

            <section className="space-y-4">
              <div className="flex items-center justify-between gap-3">
                <h2 className="text-xl font-bold text-slate-900">
                  已保存的 Agent
                </h2>
                <span className="rounded-full bg-slate-200/60 px-3 py-1 text-sm text-slate-600">
                  {profiles.data?.total ?? "—"} 个
                </span>
              </div>
              <p className="text-sm text-slate-500">
                这里只显示属于当前钱包的定义。去赛事详情选择它，即可报名。
              </p>
              {profiles.error && (
                <p
                  role="alert"
                  className="rounded-xl border border-red-200 bg-white p-4 text-red-700"
                >
                  {errorMessage(profiles.error)}
                </p>
              )}
              {profiles.isLoading && (
                <div
                  className="rounded-xl border border-slate-200 bg-white p-6 text-slate-500"
                  role="status"
                >
                  正在读取你的 Agent…
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
                        {dateTime(agent.created_at)}
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
                      带它去参赛 <ArrowRight size={16} />
                    </Link>
                    <details className="max-w-full text-sm text-slate-500">
                      <summary className="cursor-pointer">查看定义指纹</summary>
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
                    第一个 Agent，就从这里开始
                  </h3>
                  <p className="mt-2 text-sm text-slate-500">
                    在创建表单里选好策略、填个名字，然后保存。
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
