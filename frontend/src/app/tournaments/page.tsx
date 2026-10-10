"use client";

import Link from "next/link";
import { useState } from "react";
import {
  ArrowRight,
  CalendarClock,
  Coins,
  Plus,
  ShieldCheck,
  Swords,
  Users,
  Wallet,
  X,
} from "lucide-react";
import {
  encodeFunctionData,
  isAddress,
  keccak256,
  toBytes,
  type Address,
} from "viem";
import { useBlock, usePublicClient } from "wagmi";
import { useAuth } from "@/lib/auth";
import { tournamentManagerAbi } from "@/lib/course-abi";
import {
  canExpire,
  errorMessage,
  money,
  dateTime,
  type Tournament,
} from "@/lib/course";
import {
  getTournamentPresentation,
  globalRoleLabel,
} from "@/lib/course-presentation";
import { useCourseConfig, useCourseResource } from "@/hooks/useCourse";
import { useWalletTransaction } from "@/hooks/useWalletTransaction";
import { DeploymentNotice } from "@/components/course/DeploymentNotice";
import { TournamentStatusBadge } from "@/components/course/TournamentStatusBadge";
import { TransactionStatus } from "@/components/course/TransactionStatus";
import { LoginButton } from "@/components/auth/LoginButton";
import { WalletConnect } from "@/components/layout/WalletConnect";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

function detailLabel(t: Tournament, now: number) {
  if (t.status === "cancelled" || (now > 0 && canExpire(t, now)))
    return "查看退款";
  if (t.status === "finalized") return "查看结果与领奖";
  if (t.status === "running") return "查看对战";
  if (t.status === "result_submitted" || t.status === "disputed")
    return "查看结果审核";
  return "进入赛事";
}

export default function TournamentsPage() {
  const { user } = useAuth();
  const { data: config } = useCourseConfig();
  const publicClient = usePublicClient({ chainId: config?.chain_id });
  const block = useBlock({
    chainId: config?.chain_id,
    watch: true,
    query: { enabled: !!config?.deployment_configured, refetchInterval: 5000 },
  });
  const now = block.data ? Number(block.data.timestamp) : 0;
  const [offset, setOffset] = useState(0);
  const list = useCourseResource<{ tournaments: Tournament[]; total: number }>(
    "/api/tournaments?limit=20&offset=" + offset,
  );
  const tx = useWalletTransaction();
  const [error, setError] = useState("");
  const [createOpen, setCreateOpen] = useState(false);
  const [title, setTitle] = useState("");
  const [uri, setURI] = useState("");
  const [auditor, setAuditor] = useState("");
  const [registrationMinutes, setRegistrationMinutes] = useState("15");
  const [settlementHours, setSettlementHours] = useState("24");
  const [reviewSeconds, setReviewSeconds] = useState("60");
  const [roleWallet, setRoleWallet] = useState("");
  const [role, setRole] = useState("ORGANIZER_ROLE");
  const organizer = !!user?.roles.includes("organizer");
  const admin = !!user?.roles.includes("admin");

  async function create() {
    setError("");
    try {
      if (!config || !publicClient || !isAddress(auditor))
        throw new Error("请输入有效审核员钱包。");
      if (!title.trim() || toBytes(title.trim()).length > 120) {
        throw new Error(
          "赛事名称需要 1–120 个 UTF-8 字节（汉字通常占 3 字节）。",
        );
      }
      if (toBytes(uri.trim()).length > 512)
        throw new Error("规则链接最多 512 个 UTF-8 字节。");
      if (
        !/^\d+$/.test(registrationMinutes) ||
        !/^\d+$/.test(settlementHours) ||
        !/^\d+$/.test(reviewSeconds)
      ) {
        throw new Error("时间参数必须为正整数。");
      }
      if (
        ![registrationMinutes, settlementHours, reviewSeconds].every((value) =>
          Number.isSafeInteger(Number(value)),
        )
      ) {
        throw new Error("时间数值超出安全范围。");
      }
      const currentBlock = await publicClient.getBlock();
      const registration =
        currentBlock.timestamp + BigInt(registrationMinutes) * BigInt(60);
      const settlement =
        currentBlock.timestamp + BigInt(settlementHours) * BigInt(3600);
      const review = BigInt(reviewSeconds);
      if (
        Number(registrationMinutes) < 1 ||
        review < BigInt(60) ||
        review > BigInt(2592000) ||
        settlement <= registration + review
      ) {
        throw new Error(
          "结算期限必须晚于报名截止与审核窗口；审核时间为 60 秒至 30 天。",
        );
      }
      await tx.prepare(
        "创建 " +
          title.trim() +
          " · 审核员 " +
          auditor +
          " · 报名截止 " +
          dateTime(Number(registration)),
        config.manager_address,
        encodeFunctionData({
          abi: tournamentManagerAbi,
          functionName: "createTournament",
          args: [
            title.trim(),
            uri.trim(),
            auditor as Address,
            registration,
            settlement,
            review,
          ],
        }),
      );
    } catch (err) {
      setError(errorMessage(err));
    }
  }

  async function changeRole(revoke: boolean) {
    setError("");
    try {
      if (!config || !isAddress(roleWallet))
        throw new Error("请输入有效钱包地址。");
      await tx.prepare(
        (revoke ? "撤销角色" : "授予角色") + " " + role + " · " + roleWallet,
        config.manager_address,
        encodeFunctionData({
          abi: tournamentManagerAbi,
          functionName: revoke ? "revokeRole" : "grantRole",
          args: [keccak256(toBytes(role)), roleWallet as Address],
        }),
      );
    } catch (err) {
      setError(errorMessage(err));
    }
  }

  return (
    <div className="space-y-7">
      <header className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="text-sm font-semibold text-emerald-700">
            选择一场，开始参与
          </p>
          <h1 className="mt-2 text-3xl font-bold tracking-tight text-slate-900">
            赛事
          </h1>
          <p className="mt-2 max-w-2xl text-slate-600">
            两名 Agent
            免费报名，普通钱包也能赞助。进入赛事可查看对战、审核与奖金去向。
          </p>
        </div>
        <Button
          className="shrink-0 bg-emerald-600 hover:bg-emerald-700"
          aria-expanded={createOpen}
          aria-controls="create-tournament"
          onClick={() => setCreateOpen(!createOpen)}
        >
          {createOpen ? (
            <X className="mr-2 h-4 w-4" />
          ) : (
            <Plus className="mr-2 h-4 w-4" />
          )}
          {createOpen ? "收起创建表单" : "创建赛事"}
        </Button>
      </header>

      <DeploymentNotice />

      <div className="grid gap-3 md:grid-cols-3">
        <Link
          href="/agents"
          className="flex items-start gap-3 rounded-xl border border-slate-200 bg-white p-4 hover:border-emerald-300"
        >
          <Swords className="mt-0.5 h-5 w-5 shrink-0 text-emerald-600" />
          <div>
            <p className="font-semibold text-slate-900">我要参赛</p>
            <p className="mt-1 text-sm text-slate-600">
              先保存自己的 Agent，再进入开放赛事报名。
            </p>
          </div>
        </Link>
        <div className="flex items-start gap-3 rounded-xl border border-emerald-200 bg-emerald-50 p-4">
          <Coins className="mt-0.5 h-5 w-5 shrink-0 text-emerald-700" />
          <div>
            <p className="font-semibold text-slate-900">我要赞助</p>
            <p className="mt-1 text-sm text-slate-600">
              进入开放赛事，输入金额，先授权再入金。无需专门角色。
            </p>
          </div>
        </div>
        <div className="flex items-start gap-3 rounded-xl border border-slate-200 bg-white p-4">
          <ShieldCheck className="mt-0.5 h-5 w-5 shrink-0 text-emerald-600" />
          <div>
            <p className="font-semibold text-slate-900">只想先看看</p>
            <p className="mt-1 text-sm text-slate-600">
              查看赛事和资金流水无需连接钱包。
            </p>
          </div>
        </div>
      </div>

      {createOpen && (
        <section
          id="create-tournament"
          className="rounded-2xl border border-emerald-200 bg-white p-5 sm:p-6"
        >
          <div className="mb-5">
            <h2 className="text-xl font-semibold text-slate-900">
              创建一场 Agent 赛事
            </h2>
            <p className="mt-1 text-sm text-slate-600">
              设定时间，指定独立审核员。奖金由后续赞助提供，创建时不收报名费。
            </p>
          </div>
          {organizer ? (
            <form
              className="space-y-5"
              onSubmit={(event) => {
                event.preventDefault();
                void create();
              }}
            >
              <div className="grid gap-4 md:grid-cols-2">
                <label className="space-y-2 text-sm font-medium text-slate-700">
                  <span>赛事名称</span>
                  <Input
                    value={title}
                    onChange={(event) => setTitle(event.target.value)}
                    placeholder="例如：Agent 策略挑战赛"
                    maxLength={120}
                    required
                  />
                </label>
                <label className="space-y-2 text-sm font-medium text-slate-700">
                  <span>指定审核员钱包</span>
                  <Input
                    value={auditor}
                    onChange={(event) => setAuditor(event.target.value)}
                    placeholder="0x…"
                    required
                  />
                  <span className="block text-xs font-normal text-slate-500">
                    须已有审核员权限，且不能同时主办、参赛或赞助本场。
                  </span>
                </label>
                <label className="space-y-2 text-sm font-medium text-slate-700 md:col-span-2">
                  <span>
                    规则链接{" "}
                    <span className="font-normal text-slate-500">（可选）</span>
                  </span>
                  <Input
                    value={uri}
                    onChange={(event) => setURI(event.target.value)}
                    placeholder="https://… 或 ipfs://…"
                  />
                </label>
              </div>
              <div className="rounded-xl bg-slate-50 p-4">
                <div className="mb-3 flex items-center gap-2 font-medium text-slate-900">
                  <CalendarClock className="h-4 w-4 text-emerald-600" />{" "}
                  给参与者留出操作时间
                </div>
                <div className="grid gap-4 md:grid-cols-3">
                  <label className="space-y-2 text-sm text-slate-700">
                    <span>
                      距报名截止{" "}
                      <span className="text-slate-500">（分钟）</span>
                    </span>
                    <Input
                      type="number"
                      min="1"
                      step="1"
                      value={registrationMinutes}
                      onChange={(event) =>
                        setRegistrationMinutes(event.target.value)
                      }
                      required
                    />
                    <span className="block text-xs text-slate-500">
                      默认 15 分钟。截止后才能开赛。
                    </span>
                  </label>
                  <label className="space-y-2 text-sm text-slate-700">
                    <span>
                      距最终结算期限{" "}
                      <span className="text-slate-500">（小时）</span>
                    </span>
                    <Input
                      type="number"
                      min="1"
                      step="1"
                      value={settlementHours}
                      onChange={(event) =>
                        setSettlementHours(event.target.value)
                      }
                      required
                    />
                    <span className="block text-xs text-slate-500">
                      默认 24 小时。超期后可触发取消退款。
                    </span>
                  </label>
                  <label className="space-y-2 text-sm text-slate-700">
                    <span>
                      结果审核窗口{" "}
                      <span className="text-slate-500">（秒）</span>
                    </span>
                    <Input
                      type="number"
                      min="60"
                      step="1"
                      value={reviewSeconds}
                      onChange={(event) => setReviewSeconds(event.target.value)}
                      required
                    />
                    <span className="block text-xs text-slate-500">
                      默认 60 秒。提交结果后留给参与者提出异议。
                    </span>
                  </label>
                </div>
              </div>
              <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
                <Button
                  type="submit"
                  disabled={!tx.canWrite || tx.locked}
                  className="bg-emerald-600 hover:bg-emerald-700"
                >
                  预览并创建赛事 <ArrowRight className="ml-2 h-4 w-4" />
                </Button>
                <p className="text-sm text-slate-500">
                  先预览网络费，再由钱包确认发送。
                </p>
              </div>
            </form>
          ) : (
            <div className="rounded-xl bg-slate-50 p-5">
              <div className="flex items-center gap-2 font-semibold text-slate-900">
                <Wallet className="h-5 w-5 text-emerald-600" />{" "}
                创建需要主办方权限
              </div>
              <p className="mt-2 text-sm text-slate-600">
                {user
                  ? "当前身份：" +
                    globalRoleLabel(user.roles) +
                    "。请使用有主办方权限的钱包，或由管理员授予权限。"
                  : "连接钱包并签名登录后，平台会读取你的角色。参赛和赞助不需要主办方权限。"}
              </p>
              {!user && (
                <div className="mt-4 flex flex-wrap gap-3">
                  <WalletConnect />
                  <LoginButton />
                </div>
              )}
            </div>
          )}
        </section>
      )}

      {admin && (
        <details className="rounded-xl border border-slate-200 bg-white p-5">
          <summary className="cursor-pointer font-semibold text-slate-900">
            管理员：管理主办方与审核员权限
          </summary>
          <p className="mt-3 text-sm text-slate-600">
            角色写入链上，授予或撤销均需钱包确认。
          </p>
          <div className="mt-4 flex flex-col gap-3 sm:flex-row sm:flex-wrap">
            <select
              aria-label="授予或撤销的角色"
              className="rounded-lg border border-slate-200 bg-white p-2.5 text-sm"
              value={role}
              onChange={(event) => setRole(event.target.value)}
            >
              <option value="ORGANIZER_ROLE">主办方</option>
              <option value="AUDITOR_ROLE">审核员</option>
            </select>
            <Input
              className="min-w-0 sm:max-w-md"
              aria-label="授权钱包地址"
              value={roleWallet}
              onChange={(event) => setRoleWallet(event.target.value)}
              placeholder="目标钱包 0x…"
            />
            <Button
              disabled={!tx.canWrite || tx.locked}
              onClick={() => void changeRole(false)}
            >
              预览授予权限
            </Button>
            <Button
              variant="outline"
              disabled={!tx.canWrite || tx.locked}
              onClick={() => void changeRole(true)}
            >
              预览撤销权限
            </Button>
          </div>
        </details>
      )}

      {error && (
        <p
          role="alert"
          className="rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-700"
        >
          {error}
        </p>
      )}
      <TransactionStatus tx={tx} />

      <section className="space-y-4" aria-label="赛事列表">
        <div className="flex items-center justify-between gap-3">
          <h2 className="text-xl font-semibold text-slate-900">全部赛事</h2>
          <span className="text-sm text-slate-500">
            {list.data ? "共 " + list.data.total + " 场" : "读取中…"}
          </span>
        </div>
        {list.isLoading && (
          <p className="rounded-xl bg-white p-6 text-slate-600">
            正在读取赛事…
          </p>
        )}
        {list.error && (
          <p role="alert" className="rounded-xl bg-red-50 p-4 text-red-700">
            {errorMessage(list.error)}
          </p>
        )}
        {list.data && !list.data.tournaments.length && (
          <div className="rounded-2xl border border-dashed border-slate-300 bg-white p-8 text-center">
            <Swords className="mx-auto h-8 w-8 text-slate-400" />
            <p className="mt-3 font-medium text-slate-700">还没有赛事</p>
            <p className="mt-1 text-sm text-slate-500">
              主办方可通过上方“创建赛事”发起第一场。
            </p>
          </div>
        )}
        <div className="grid gap-4 lg:grid-cols-2">
          {list.data?.tournaments.map((t) => {
            const presentation = getTournamentPresentation(t, now);
            return (
              <Link
                className="group rounded-2xl border border-slate-200 bg-white p-5 transition hover:border-emerald-400 hover:shadow-sm"
                href={"/tournaments/" + t.id}
                key={t.id}
              >
                <div className="flex flex-wrap items-center justify-between gap-2">
                  <span className="text-xs font-medium text-slate-500">
                    赛事 #{t.id}
                  </span>
                  <TournamentStatusBadge tournament={t} now={now} />
                </div>
                <h3 className="mt-3 text-lg font-semibold text-slate-900 group-hover:text-emerald-700">
                  {t.title}
                </h3>
                <div className="mt-4 grid grid-cols-2 gap-3">
                  <div className="rounded-lg bg-emerald-50 p-3">
                    <p className="flex items-center gap-1.5 text-xs text-emerald-800">
                      <Coins className="h-3.5 w-3.5" /> 赞助奖池
                    </p>
                    <p className="mt-1 break-words font-semibold text-emerald-900">
                      {money(t.prize_pool, config?.token_decimals)}{" "}
                      <span className="text-xs font-normal">DUSD</span>
                    </p>
                  </div>
                  <div className="rounded-lg bg-slate-50 p-3">
                    <p className="flex items-center gap-1.5 text-xs text-slate-600">
                      <Users className="h-3.5 w-3.5" /> 已报名
                    </p>
                    <p className="mt-1 font-semibold text-slate-900">
                      {t.participants.length}{" "}
                      <span className="text-sm font-normal text-slate-500">
                        / 2 人
                      </span>
                    </p>
                  </div>
                </div>
                <p className="mt-4 text-sm text-slate-600">
                  {presentation.hint}
                </p>
                <div className="mt-4 flex flex-wrap items-center justify-between gap-2 border-t border-slate-100 pt-3">
                  <span className="text-xs text-slate-500">
                    报名截止 {dateTime(t.registration_deadline)}
                  </span>
                  <span className="inline-flex items-center gap-1 text-sm font-semibold text-emerald-700">
                    {detailLabel(t, now)} <ArrowRight className="h-4 w-4" />
                  </span>
                </div>
              </Link>
            );
          })}
        </div>
      </section>
      <div className="flex items-center justify-center gap-3">
        <Button
          variant="outline"
          disabled={offset === 0}
          onClick={() => setOffset(Math.max(0, offset - 20))}
        >
          上一页
        </Button>
        <span className="text-sm text-slate-500">
          第 {Math.floor(offset / 20) + 1} 页
        </span>
        <Button
          variant="outline"
          disabled={offset + 20 >= (list.data?.total || 0)}
          onClick={() => setOffset(offset + 20)}
        >
          下一页
        </Button>
      </div>
    </div>
  );
}
