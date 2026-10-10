import { formatUnits, parseUnits, keccak256, toBytes, maxUint256, type Address, type Hex } from "viem";
export const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";
export const TOKEN_DECIMALS = 6;
export type Role = "admin" | "organizer" | "auditor" | "participant";
export type TournamentStatus = "open" | "running" | "result_submitted" | "disputed" | "finalized" | "cancelled";
export interface CourseConfig {
    chain_id: number;
    deployment_configured: boolean;
    manager_address: Address;
    escrow_address: Address;
    registry_address: Address;
    review_address: Address;
    token_address: Address;
    token_decimals: number;
    deployment_block: number;
    confirmations: number;
    simulation_enabled: boolean;
    auth_domain: string;
    llm_configured?: boolean;
    execution_mode?: "llm" | "deterministic_simulation" | "unavailable";
    llm_provider?: string;
    llm_model?: string;
    llm_timeout_seconds?: number;
}
export interface Tournament {
    id: number;
    title: string;
    organizer: Address;
    auditor: Address;
    token: Address;
    entry_fee: string;
    prize_pool: string;
    prize_paid: string;
    prize_claimed: boolean;
    max_participants: number;
    registration_deadline: number;
    settlement_deadline: number;
    review_period: number;
    review_deadline: number;
    status: TournamentStatus;
    participants: Address[];
    participant_agents: {
        address: Address;
        agent_hash: Hex;
    }[];
    winner: Address | "";
    proposed_winner: Address | "";
    evidence_hash: Hex | "";
    metadata_uri: string;
    chain_updated_at: number;
}
export interface AuditEvent {
    id: string;
    chain_id: number;
    tournament_id: number;
    contract_address: Address;
    event_name: string;
    actor: string;
    amount_raw: string;
    tx_hash: Hex;
    block_number: number;
    log_index: number;
    timestamp: number;
    data: Record<string, unknown>;
}
export interface AgentProfile {
    id: string;
    agent_hash: Hex;
    name: string;
    personality: string;
    owner_address: Address;
    created_at: number;
}
export interface Dashboard {
    tournaments_total: number;
    registered_participants: number;
    active_tournaments: number;
    completed_tournaments: number;
    escrow_in_raw: string;
    prizes_paid_raw: string;
    refunds_paid_raw: string;
    ledger_entries: number;
    event_count: number;
    indexed_block: number;
    chain_head: number;
    observed_at: number;
}
export interface MatchRound {
    round: number;
    red_hp: number;
    blue_hp: number;
    red_position: {
        x: number;
        y: number;
    };
    blue_position: {
        x: number;
        y: number;
    };
    red_action: {
        type: string;
        target?: {
            x: number;
            y: number;
        };
        failed?: boolean;
        fail_reason?: string;
    };
    blue_action: {
        type: string;
        target?: {
            x: number;
            y: number;
        };
        failed?: boolean;
        fail_reason?: string;
    };
    red_reasoning?: string;
    blue_reasoning?: string;
    red_state?: EngineAgentState;
    blue_state?: EngineAgentState;
    red_llm_call_id?: number;
    blue_llm_call_id?: number;
}
export interface EngineAgentState {
    hp: number;
    max_hp: number;
    position: {
        x: number;
        y: number;
    };
    status: {
        type: string;
        remaining: number;
    }[];
    skill_cooldown: number;
    is_charging: boolean;
}
export interface MatchResult {
    tournament_id: number;
    status: string;
    mode: string;
    red_address: Address;
    blue_address: Address;
    current_round: number;
    winner_address: Address | "";
    outcome: string;
    evidence_hash: Hex | "";
    rounds: MatchRound[];
    state: Record<string, unknown>;
    error: string;
    retry_allowed?: boolean;
    llm_summary?: {
        provider: string;
        model: string;
        attempted: number;
        succeeded: number;
        failed: number;
        latency_ms: number;
        usage: {
            prompt_tokens: number | null;
            completion_tokens: number | null;
            total_tokens: number | null;
        };
        last_error_code: string;
    };
    llm_calls?: {
        sequence: number;
        round: number;
        side: "red" | "blue";
        provider: string;
        model: string;
        success: boolean;
        latency_ms: number;
        http_status?: number;
        request_id?: string;
        usage?: {
            prompt_tokens?: number;
            completion_tokens?: number;
            total_tokens?: number;
        };
        error_code?: string;
    }[];
    red_agent?: {
        name: string;
        personality: string;
    };
    blue_agent?: {
        name: string;
        personality: string;
    };
}
/** A configured LLM mode alone is not evidence of a successful model call. */
export function matchExecutionPresentation(
    match: Pick<MatchResult, "mode" | "status" | "llm_summary">,
    lang: "zh" | "en" = "zh",
): { label: string; hint: string } {
    const l = (zh: string, en: string) => lang === "en" ? en : zh;
    if (match.mode === "deterministic_simulation") return {
        label: l("确定性模拟", "Deterministic simulation"),
        hint: l("本场使用确定性模拟，没有调用 LLM；可用于体验赛事和奖金流程。", "This match uses a deterministic simulation without an LLM, for trying the tournament and prize flow."),
    };
    if (match.mode !== "llm") return {
        label: l("执行模式未确认", "Execution mode unverified"),
        hint: l("本场未提供可识别的执行模式，请核查记录。", "This record does not identify a supported execution mode. Inspect the record before making model claims."),
    };
    const summary = match.llm_summary;
    if (!summary) return {
        label: l("历史 LLM · 调用未验证", "Historical LLM · calls unverified"),
        hint: l("历史记录没有模型调用统计，不能仅凭 LLM 模式确认真实调用成功。", "This historical record has no model-call statistics. LLM mode alone does not verify a successful real call."),
    };
    if (summary.succeeded <= 0) return {
        label: l("LLM · 尚无成功调用", "LLM · no successful calls yet"),
        hint: summary.failed > 0
            ? l("尚无成功模型调用，已记录失败；请查看错误原因。这不是模拟对战成功。", "No successful model call is recorded. Inspect the failures below; this does not count as a successful simulated or LLM match.")
            : l("已选择 LLM 模式，正在等待成功调用记录；模型配置不代表调用已成功。", "LLM mode is selected; waiting for a successful call record. Configuration does not prove that a call succeeded."),
    };
    if (match.status === "completed") return {
        label: summary.failed > 0
            ? l("LLM 对战完成 · 含重试", "LLM match complete · includes retries")
            : l("LLM 对战完成", "LLM match complete"),
        hint: summary.failed > 0
            ? l("对战已完成，包含此前失败及重试记录；成功回合来自已记录的模型决策。", "The match completed with earlier failures and retries retained. Successful rounds use the recorded model decisions.")
            : l("对战已完成，模型成功调用已记录，未记录失败调用。审核与付款状态见赛事进度和资金区。", "The match completed with successful model calls and no recorded call failures. See the tournament progress and funds panel for review and payment status."),
    };
    return {
        label: l("LLM · 已有成功调用", "LLM · successful calls recorded"),
        hint: l("已收到成功模型决策，对战是否完成以执行状态为准。", "Successful model decisions are recorded. Check the execution status to see whether the match has completed."),
    };
}

/** Resume only the interrupted state that the backend explicitly marks retryable. */
export function canResumeMatch(match: Pick<MatchResult, "status" | "retry_allowed"> | undefined): boolean {
    return match?.status === "interrupted" && match.retry_allowed === true;
}

export function llmErrorMessage(code: string, lang: "zh" | "en" = "zh"): string {
    const labels: Record<string, [string, string]> = {
        not_configured: ["模型未配置，请联系维护者。", "The model is not configured. Contact the maintainer."],
        provider_error: ["模型服务请求失败，请检查服务后重试。", "The model provider request failed. Check the service before retrying."],
        authentication_error: ["模型凭据被拒绝，请维护者检查配置。", "The model credentials were rejected. Ask the maintainer to check the configuration."],
        rate_limited: ["模型服务限流，请稍后重试。", "The model service rate-limited the request. Try again later."],
        timeout: ["模型请求超时。", "The model request timed out."],
        cancelled: ["模型请求已中断。", "The model request was interrupted."],
        invalid_response: ["模型返回无效响应，该调用失败。", "The model returned an invalid response and this call failed."],
        response_too_large: ["模型响应超过大小限制，该调用失败。", "The model response exceeded the size limit and this call failed."],
        parse_error: ["模型回复无法解析为有效动作。", "The model response could not be parsed as a valid action."],
    };
    return labels[code]?.[lang === "en" ? 1 : 0] || code;
}
export interface Evidence {
    evidence_hash: Hex;
    canonical_json: string;
    evidence: Record<string, unknown>;
}
export interface EvidenceNote {
    author_address: Address;
    kind: "dispute" | "decision" | "cancellation";
    text: string;
    evidence_hash: Hex;
    created_at: number;
    canonical_json: string;
}
export class ApiError extends Error {
    constructor(message: string, public status: number) { super(message); this.name = "ApiError"; }
}
const ENGLISH_ERRORS: Record<string, string> = {
  请求失败: "Request failed.",
  "请求失败，请检查网络并重试。":
    "Request failed. Check your connection and try again.",
  "服务未返回有效 JSON，请检查 API 地址。":
    "The service returned invalid JSON. Check the API address.",
  "代币精度配置无效。": "Invalid token decimal configuration.",
  "金额超出合约范围。": "The amount exceeds the contract limit.",
  "连接钱包并切换至已部署网络后，才能发送交易。":
    "Connect your wallet and switch to the deployed network before sending a transaction.",
  "目标地址没有合约代码，请检查部署配置。":
    "There is no contract at the target address. Check the deployment configuration.",
  "钱包或网络已变化，请重新估算。":
    "Your wallet or network has changed. Estimate the transaction again.",
  "原交易已被取消或替换为其他操作，请核查新交易回执后重试。":
    "The original transaction was cancelled or replaced. Check the replacement receipt before trying again.",
  "交易已上链但执行回滚；请查看交易回执。":
    "The transaction was included on-chain but reverted. Check its receipt.",
  "确认数不足，请稍后继续检查。":
    "Not enough confirmations yet. Check again shortly.",
  "原交易已被取消或替换为其他操作，请核查新交易回执。":
    "The original transaction was cancelled or replaced. Check the replacement receipt.",
  "交易回执显示执行回滚。": "The transaction receipt indicates a revert.",
  "请先连接钱包。": "Connect your wallet first.",
  "签名期间钱包或网络已变化，请重新登录。":
    "Your wallet or network changed while signing. Sign in again.",
};

/** Translate known UI failures only; retain unknown service/provider diagnostics verbatim. */
export function errorMessage(error: unknown, lang: "zh" | "en" = "zh"): string {
  const brief =
    error instanceof Error
      ? ("shortMessage" in error && typeof error.shortMessage === "string"
          ? error.shortMessage
          : error.message) || "请求失败"
      : typeof error === "string"
        ? error
        : "请求失败，请检查网络并重试。";
  if (lang === "zh") return brief;
  if (ENGLISH_ERRORS[brief]) return ENGLISH_ERRORS[brief];
  const precision = brief.match(/^请输入非负金额，最多 (\d+) 位小数。$/);
  if (precision)
    return `Enter a non-negative amount with at most ${precision[1]} decimal places.`;
  const http = brief.match(/^服务返回 HTTP (\d+)$/);
  if (http) return `The service returned HTTP ${http[1]}.`;
  if (brief.startsWith("交易状态待确认："))
    return `Transaction status unconfirmed: ${errorMessage(brief.slice("交易状态待确认：".length), lang)}`;
  return brief;
}
export async function request<T>(path: string, init?: RequestInit): Promise<T> {
    const res = await fetch(`${API_URL}${path}`, { ...init, signal: init?.signal ?? AbortSignal.timeout(20000), headers: { "Content-Type": "application/json", ...init?.headers } });
    const body: unknown = await res.json().catch(() => null);
    if (!res.ok) {
        const detail = body && typeof body === "object" && "error" in body ? body.error : null;
        const message = detail && typeof detail === "object" && "message" in detail ? detail.message : detail;
        throw new ApiError(typeof message === "string" ? message : `服务返回 HTTP ${res.status}`, res.status);
    }
    if (body === null) throw new ApiError("服务未返回有效 JSON，请检查 API 地址。", res.status);
    return body as T;
}
export function money(raw: string | undefined, decimals = TOKEN_DECIMALS): string { return formatUnits(BigInt(raw || "0"), decimals); }
export function parseMoney(value: string, decimals = TOKEN_DECIMALS): bigint {
    if (!Number.isInteger(decimals) || decimals < 0 || decimals > 36)
        throw new Error("代币精度配置无效。");
    const pattern = decimals > 0 ? `^\\d+(\\.\\d{1,${decimals}})?$` : "^\\d+$";
    if (!new RegExp(pattern).test(value.trim()))
        throw new Error(`请输入非负金额，最多 ${decimals} 位小数。`);
    const amount = parseUnits(value.trim(), decimals);
    if (amount > maxUint256)
        throw new Error("金额超出合约范围。");
    return amount;
}
export function shortAddress(value: string | undefined): string { return value ? `${value.slice(0, 6)}…${value.slice(-4)}` : "—"; }
export function sameAddress(a?: string, b?: string): boolean { return !!a && !!b && a.toLowerCase() === b.toLowerCase(); }
export const STATUS_LABEL: Record<TournamentStatus, string> = { open: "开放报名", running: "对战中", result_submitted: "结果审核", disputed: "存在异议", finalized: "已结算", cancelled: "已取消" };
export function explorer(chainId: number, tx: string): string | null { return chainId === 11155111 ? `https://sepolia.etherscan.io/tx/${tx}` : null; }
export function dateTime(seconds: number, lang: "zh" | "en" = "zh"): string { return seconds ? new Date(seconds * 1000).toLocaleString(lang === "zh" ? "zh-CN" : "en-GB") : "—"; }
export function registrationOpen(t: Tournament, now: number): boolean { return t.status === "open" && now < t.registration_deadline; }
export function reviewComplete(t: Tournament, now: number): boolean { return now >= t.review_deadline && now < t.settlement_deadline; }
export function canExpire(t: Tournament, now: number): boolean {
    return !["finalized", "cancelled"].includes(t.status) && (now >= t.settlement_deadline || (t.status === "open" && t.participants.length < 2 && now >= t.registration_deadline));
}
export function verifyEvidence(canonicalJSON: string, hash: string): boolean { return keccak256(toBytes(canonicalJSON)) === hash; }
export function receiptVerdict(status: "success" | "reverted", replacement = ""): "success" | "reverted" | "replaced" {
    if (replacement === "cancelled" || replacement === "replaced") return "replaced";
    return status;
}
