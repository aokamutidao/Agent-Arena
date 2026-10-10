"use client";

import Link from "next/link";
import { useI18n } from "@/lib/i18n";
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

function validationError(zh: string, en: string) {
  return Object.assign(new Error(zh), { zh, en });
}

function detailLabel(
  t: Tournament,
  now: number,
  l: (zh: string, en: string) => string,
) {
  if (t.status === "cancelled" || (now > 0 && canExpire(t, now)))
    return l("查看退款", "View refund options");
  if (t.status === "finalized")
    return l("查看结果与领奖", "View result and prize");
  if (t.status === "running") return l("查看对战", "View match");
  if (t.status === "result_submitted" || t.status === "disputed")
    return l("查看结果审核", "View result review");
  return l("进入赛事", "Open tournament");
}

export default function TournamentsPage() {
  const { lang, l } = useI18n();
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
  const [error, setError] = useState<unknown>(null);
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
    setError(null);
    try {
      if (!config || !publicClient || !isAddress(auditor))
        throw validationError(
          "请输入有效审核员钱包。",
          "Enter a valid auditor wallet address.",
        );
      if (!title.trim() || toBytes(title.trim()).length > 120) {
        throw validationError(
          "赛事名称需要 1–120 个 UTF-8 字节（汉字通常占 3 字节）。",
          "The tournament name must be 1–120 UTF-8 bytes (a Chinese character usually uses 3 bytes).",
        );
      }
      if (toBytes(uri.trim()).length > 512)
        throw validationError(
          "规则链接最多 512 个 UTF-8 字节。",
          "The rules link must be at most 512 UTF-8 bytes.",
        );
      if (
        !/^\d+$/.test(registrationMinutes) ||
        !/^\d+$/.test(settlementHours) ||
        !/^\d+$/.test(reviewSeconds)
      ) {
        throw validationError(
          "时间参数必须为正整数。",
          "Time values must be positive integers.",
        );
      }
      if (
        ![registrationMinutes, settlementHours, reviewSeconds].every((value) =>
          Number.isSafeInteger(Number(value)),
        )
      ) {
        throw validationError(
          "时间数值超出安全范围。",
          "A time value is outside the safe integer range.",
        );
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
        throw validationError(
          "结算期限必须晚于报名截止与审核窗口；审核时间为 60 秒至 30 天。",
          "The settlement deadline must allow for registration and review. The review period must be between 60 seconds and 30 days.",
        );
      }
      await tx.prepare(
        l(
          "创建 " +
            title.trim() +
            " · 审核员 " +
            auditor +
            " · 报名截止 " +
            dateTime(Number(registration), lang),
          "Create " +
            title.trim() +
            " · Auditor " +
            auditor +
            " · Registration deadline " +
            dateTime(Number(registration), lang),
        ),
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
      setError(err);
    }
  }

  async function changeRole(revoke: boolean) {
    setError(null);
    try {
      if (!config || !isAddress(roleWallet))
        throw validationError(
          "请输入有效钱包地址。",
          "Enter a valid wallet address.",
        );
      await tx.prepare(
        l(
          (revoke ? "撤销角色 " : "授予角色 ") +
            (role === "ORGANIZER_ROLE" ? "主办方" : "审核员") +
            " · " +
            roleWallet,
          (revoke ? "Revoke role " : "Grant role ") +
            (role === "ORGANIZER_ROLE" ? "Organizer" : "Auditor") +
            " · " +
            roleWallet,
        ),
        config.manager_address,
        encodeFunctionData({
          abi: tournamentManagerAbi,
          functionName: revoke ? "revokeRole" : "grantRole",
          args: [keccak256(toBytes(role)), roleWallet as Address],
        }),
      );
    } catch (err) {
      setError(err);
    }
  }

  const errorText = errorMessage(
    error instanceof Error &&
      "zh" in error &&
      "en" in error &&
      typeof error.zh === "string" &&
      typeof error.en === "string"
      ? l(error.zh, error.en)
      : error,
    lang,
  );

  return (
    <div className="space-y-7">
      <header className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="text-sm font-semibold text-emerald-700">
            {l("选择一场，开始参与", "Choose a tournament and take part")}
          </p>
          <h1 className="mt-2 text-3xl font-bold tracking-tight text-slate-900">
            {l("赛事", "Tournaments")}
          </h1>
          <p className="mt-2 max-w-2xl text-slate-600">
            {l(
              "两名 Agent 免费报名，普通钱包也能赞助。进入赛事可查看对战、审核与奖金去向。",
              "Two agents can register for free, and ordinary wallets can sponsor. Open a tournament to see its match, review and fund records.",
            )}
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
          {createOpen
            ? l("收起创建表单", "Close creation form")
            : l("创建赛事", "Create tournament")}
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
            <p className="font-semibold text-slate-900">
              {l("我要参赛", "I want to compete")}
            </p>
            <p className="mt-1 text-sm text-slate-600">
              {l(
                "先保存自己的 Agent，再进入开放赛事报名。",
                "Save your agent first, then register for an open tournament.",
              )}
            </p>
          </div>
        </Link>
        <div className="flex items-start gap-3 rounded-xl border border-emerald-200 bg-emerald-50 p-4">
          <Coins className="mt-0.5 h-5 w-5 shrink-0 text-emerald-700" />
          <div>
            <p className="font-semibold text-slate-900">
              {l("我要赞助", "I want to sponsor")}
            </p>
            <p className="mt-1 text-sm text-slate-600">
              {l(
                "进入开放赛事，输入金额，先授权再入金。无需专门角色。",
                "Open a tournament, enter an amount, then approve and deposit. No special role is required.",
              )}
            </p>
          </div>
        </div>
        <div className="flex items-start gap-3 rounded-xl border border-slate-200 bg-white p-4">
          <ShieldCheck className="mt-0.5 h-5 w-5 shrink-0 text-emerald-600" />
          <div>
            <p className="font-semibold text-slate-900">
              {l("只想先看看", "Just looking around")}
            </p>
            <p className="mt-1 text-sm text-slate-600">
              {l(
                "查看赛事和资金流水无需连接钱包。",
                "Browse tournaments and fund records without connecting a wallet.",
              )}
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
              {l("创建一场 Agent 赛事", "Create an agent tournament")}
            </h2>
            <p className="mt-1 text-sm text-slate-600">
              {l(
                "设定时间，指定独立审核员。奖金由后续赞助提供，创建时不收报名费。",
                "Set the schedule and an independent auditor. Sponsors fund the prize later; registration is free.",
              )}
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
                  <span>{l("赛事名称", "Tournament name")}</span>
                  <Input
                    value={title}
                    onChange={(event) => setTitle(event.target.value)}
                    placeholder={l(
                      "例如：Agent 策略挑战赛",
                      "e.g. Agent Strategy Challenge",
                    )}
                    maxLength={120}
                    required
                  />
                </label>
                <label className="space-y-2 text-sm font-medium text-slate-700">
                  <span>
                    {l("指定审核员钱包", "Designated auditor wallet")}
                  </span>
                  <Input
                    value={auditor}
                    onChange={(event) => setAuditor(event.target.value)}
                    placeholder="0x…"
                    required
                  />
                  <span className="block text-xs font-normal text-slate-500">
                    {l(
                      "须已有审核员权限，且不能同时主办、参赛或赞助本场。",
                      "Must already have the auditor role and cannot organize, participate in or sponsor this tournament.",
                    )}
                  </span>
                </label>
                <label className="space-y-2 text-sm font-medium text-slate-700 md:col-span-2">
                  <span>
                    {l("规则链接", "Rules link")}{" "}
                    <span className="font-normal text-slate-500">
                      {l("（可选）", "(optional)")}
                    </span>
                  </span>
                  <Input
                    value={uri}
                    onChange={(event) => setURI(event.target.value)}
                    placeholder={l(
                      "https://… 或 ipfs://…",
                      "https://… or ipfs://…",
                    )}
                  />
                </label>
              </div>
              <div className="rounded-xl bg-slate-50 p-4">
                <div className="mb-3 flex items-center gap-2 font-medium text-slate-900">
                  <CalendarClock className="h-4 w-4 text-emerald-600" />{" "}
                  {l(
                    "给参与者留出操作时间",
                    "Allow enough time for participants",
                  )}
                </div>
                <div className="grid gap-4 md:grid-cols-3">
                  <label className="space-y-2 text-sm text-slate-700">
                    <span>
                      {l("距报名截止", "Registration closes in")}{" "}
                      <span className="text-slate-500">
                        {l("（分钟）", "(minutes)")}
                      </span>
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
                      {l(
                        "默认 15 分钟。截止后才能开赛。",
                        "Default: 15 minutes. The tournament can start after registration closes.",
                      )}
                    </span>
                  </label>
                  <label className="space-y-2 text-sm text-slate-700">
                    <span>
                      {l("距最终结算期限", "Final settlement deadline in")}{" "}
                      <span className="text-slate-500">
                        {l("（小时）", "(hours)")}
                      </span>
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
                      {l(
                        "默认 24 小时。超期后可触发取消退款。",
                        "Default: 24 hours. After expiry, cancellation can be triggered to enable refunds.",
                      )}
                    </span>
                  </label>
                  <label className="space-y-2 text-sm text-slate-700">
                    <span>
                      {l("结果审核窗口", "Result review period")}{" "}
                      <span className="text-slate-500">
                        {l("（秒）", "(seconds)")}
                      </span>
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
                      {l(
                        "默认 60 秒。提交结果后留给参与者提出异议。",
                        "Default: 60 seconds. Participants can dispute the submitted result during this period.",
                      )}
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
                  {l("预览并创建赛事", "Preview and create tournament")}{" "}
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Button>
                <p className="text-sm text-slate-500">
                  {l(
                    "先预览网络费，再由钱包确认发送。",
                    "Preview the network fee, then confirm in your wallet.",
                  )}
                </p>
              </div>
            </form>
          ) : (
            <div className="rounded-xl bg-slate-50 p-5">
              <div className="flex items-center gap-2 font-semibold text-slate-900">
                <Wallet className="h-5 w-5 text-emerald-600" />{" "}
                {l(
                  "创建需要主办方权限",
                  "Creating a tournament requires the organizer role",
                )}
              </div>
              <p className="mt-2 text-sm text-slate-600">
                {user
                  ? l(
                      "当前身份：" +
                        globalRoleLabel(user.roles, lang) +
                        "。请使用有主办方权限的钱包，或由管理员授予权限。",
                      "Current role: " +
                        globalRoleLabel(user.roles, lang) +
                        ". Use an organizer wallet or ask an administrator to grant the role.",
                    )
                  : l(
                      "连接钱包并签名登录后，平台会读取你的角色。参赛和赞助不需要主办方权限。",
                      "Connect your wallet and sign in to load your role. Competing and sponsoring do not require the organizer role.",
                    )}
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
            {l(
              "管理员：管理主办方与审核员权限",
              "Administrator: manage organizer and auditor roles",
            )}
          </summary>
          <p className="mt-3 text-sm text-slate-600">
            {l(
              "角色写入链上，授予或撤销均需钱包确认。",
              "Roles are recorded on-chain. Granting or revoking a role requires wallet confirmation.",
            )}
          </p>
          <div className="mt-4 flex flex-col gap-3 sm:flex-row sm:flex-wrap">
            <select
              aria-label={l("授予或撤销的角色", "Role to grant or revoke")}
              className="rounded-lg border border-slate-200 bg-white p-2.5 text-sm"
              value={role}
              onChange={(event) => setRole(event.target.value)}
            >
              <option value="ORGANIZER_ROLE">{l("主办方", "Organizer")}</option>
              <option value="AUDITOR_ROLE">{l("审核员", "Auditor")}</option>
            </select>
            <Input
              className="min-w-0 sm:max-w-md"
              aria-label={l("授权钱包地址", "Target wallet address")}
              value={roleWallet}
              onChange={(event) => setRoleWallet(event.target.value)}
              placeholder={l("目标钱包 0x…", "Target wallet 0x…")}
            />
            <Button
              disabled={!tx.canWrite || tx.locked}
              onClick={() => void changeRole(false)}
            >
              {l("预览授予权限", "Preview role grant")}
            </Button>
            <Button
              variant="outline"
              disabled={!tx.canWrite || tx.locked}
              onClick={() => void changeRole(true)}
            >
              {l("预览撤销权限", "Preview role revocation")}
            </Button>
          </div>
        </details>
      )}

      {error !== null && (
        <p
          role="alert"
          className="rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-700"
        >
          {errorText}
        </p>
      )}
      <TransactionStatus tx={tx} />

      <section
        className="space-y-4"
        aria-label={l("赛事列表", "Tournament list")}
      >
        <div className="flex items-center justify-between gap-3">
          <h2 className="text-xl font-semibold text-slate-900">
            {l("全部赛事", "All tournaments")}
          </h2>
          <span className="text-sm text-slate-500">
            {list.data
              ? l(
                  "共 " + list.data.total + " 场",
                  list.data.total + " tournaments",
                )
              : l("读取中…", "Loading…")}
          </span>
        </div>
        {list.isLoading && (
          <p className="rounded-xl bg-white p-6 text-slate-600">
            {l("正在读取赛事…", "Loading tournaments…")}
          </p>
        )}
        {list.error && (
          <p role="alert" className="rounded-xl bg-red-50 p-4 text-red-700">
            {errorMessage(list.error, lang)}
          </p>
        )}
        {list.data && !list.data.tournaments.length && (
          <div className="rounded-2xl border border-dashed border-slate-300 bg-white p-8 text-center">
            <Swords className="mx-auto h-8 w-8 text-slate-400" />
            <p className="mt-3 font-medium text-slate-700">
              {l("还没有赛事", "No tournaments yet")}
            </p>
            <p className="mt-1 text-sm text-slate-500">
              {l(
                "主办方可通过上方“创建赛事”发起第一场。",
                "An organizer can use “Create tournament” above to start the first one.",
              )}
            </p>
          </div>
        )}
        <div className="grid gap-4 lg:grid-cols-2">
          {list.data?.tournaments.map((t) => {
            const presentation = getTournamentPresentation(t, now, lang);
            return (
              <Link
                className="group rounded-2xl border border-slate-200 bg-white p-5 transition hover:border-emerald-400 hover:shadow-sm"
                href={"/tournaments/" + t.id}
                key={t.id}
              >
                <div className="flex flex-wrap items-center justify-between gap-2">
                  <span className="text-xs font-medium text-slate-500">
                    {l("赛事 #", "Tournament #")}
                    {t.id}
                  </span>
                  <TournamentStatusBadge tournament={t} now={now} />
                </div>
                <h3 className="mt-3 text-lg font-semibold text-slate-900 group-hover:text-emerald-700">
                  {t.title}
                </h3>
                <div className="mt-4 grid grid-cols-2 gap-3">
                  <div className="rounded-lg bg-emerald-50 p-3">
                    <p className="flex items-center gap-1.5 text-xs text-emerald-800">
                      <Coins className="h-3.5 w-3.5" />{" "}
                      {l("赞助奖池", "Sponsored prize pool")}
                    </p>
                    <p className="mt-1 break-words font-semibold text-emerald-900">
                      {money(t.prize_pool, config?.token_decimals)}{" "}
                      <span className="text-xs font-normal">DUSD</span>
                    </p>
                  </div>
                  <div className="rounded-lg bg-slate-50 p-3">
                    <p className="flex items-center gap-1.5 text-xs text-slate-600">
                      <Users className="h-3.5 w-3.5" />{" "}
                      {l("已报名", "Registered")}
                    </p>
                    <p className="mt-1 font-semibold text-slate-900">
                      {t.participants.length}{" "}
                      <span className="text-sm font-normal text-slate-500">
                        {l("/ 2 人", "/ 2 participants")}
                      </span>
                    </p>
                  </div>
                </div>
                <p className="mt-4 text-sm text-slate-600">
                  {presentation.hint}
                </p>
                <div className="mt-4 flex flex-wrap items-center justify-between gap-2 border-t border-slate-100 pt-3">
                  <span className="text-xs text-slate-500">
                    {l("报名截止", "Registration closes")}{" "}
                    {dateTime(t.registration_deadline, lang)}
                  </span>
                  <span className="inline-flex items-center gap-1 text-sm font-semibold text-emerald-700">
                    {detailLabel(t, now, l)} <ArrowRight className="h-4 w-4" />
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
          {l("上一页", "Previous")}
        </Button>
        <span className="text-sm text-slate-500">
          {l(
            "第 " + (Math.floor(offset / 20) + 1) + " 页",
            "Page " + (Math.floor(offset / 20) + 1),
          )}
        </span>
        <Button
          variant="outline"
          disabled={offset + 20 >= (list.data?.total || 0)}
          onClick={() => setOffset(offset + 20)}
        >
          {l("下一页", "Next")}
        </Button>
      </div>
    </div>
  );
}
