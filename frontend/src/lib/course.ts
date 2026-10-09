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
    red_agent?: {
        name: string;
        personality: string;
    };
    blue_agent?: {
        name: string;
        personality: string;
    };
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
export function errorMessage(error: unknown): string {
    if (error instanceof Error) {
        const brief = "shortMessage" in error && typeof error.shortMessage === "string" ? error.shortMessage : error.message;
        return brief || "请求失败";
    }
    return typeof error === "string" ? error : "请求失败，请检查网络并重试。";
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
export function dateTime(seconds: number): string { return seconds ? new Date(seconds * 1000).toLocaleString("zh-CN") : "—"; }
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
