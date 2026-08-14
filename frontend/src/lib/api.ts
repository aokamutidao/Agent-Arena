import type {
  GamesResponse,
  GameDetail,
  AgentsResponse,
  EstimateRequest,
  EstimateResponse,
  VoteRequest,
  VoteResponse,
} from "@/types/game";

const API_BASE =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

async function fetchJSON<T>(url: string): Promise<T> {
  const res = await fetch(url);
  if (!res.ok) {
    throw new Error(`API error: ${res.status}`);
  }
  return res.json();
}

async function postJSON<T>(url: string, body: unknown): Promise<T> {
  const res = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    throw new Error(`API error: ${res.status}`);
  }
  return res.json();
}

export const api = {
  // 对局列表
  listGames(status?: string): Promise<GamesResponse> {
    const params = status ? `?status=${status}` : "";
    return fetchJSON(`${API_BASE}/api/games${params}`);
  },

  // 对局详情
  getGame(gameId: number): Promise<GameDetail> {
    return fetchJSON(`${API_BASE}/api/games/${gameId}`);
  },

  // Agent 列表
  listAgents(): Promise<AgentsResponse> {
    return fetchJSON(`${API_BASE}/api/agents`);
  },

  // 下注预估
  estimateBet(req: EstimateRequest): Promise<EstimateResponse> {
    return postJSON(`${API_BASE}/api/bets/estimate`, req);
  },

  // 策略投票
  voteStrategy(req: VoteRequest): Promise<VoteResponse> {
    return postJSON(`${API_BASE}/api/strategy/vote`, req);
  },

  // 手动开始对局
  startGame(gameId: number): Promise<{ game_id: number; status: string }> {
    return postJSON(`${API_BASE}/api/games/${gameId}/start`, {});
  },

  // 下注（内存模式）
  placeBet(req: {
    game_id: number;
    side: "red" | "blue";
    amount: string;
    strategy?: string;
    user?: string;
  }): Promise<{
    game_id: number;
    side: string;
    amount: string;
    total_bet_red: string;
    total_bet_blue: string;
    odds_red: number;
    odds_blue: number;
  }> {
    return postJSON(`${API_BASE}/api/bets/place`, req);
  },

  // 创建新对局
  createGame(): Promise<{
    game_id: number;
    agent_red: string;
    agent_blue: string;
    status: string;
  }> {
    return postJSON(`${API_BASE}/api/games/create`, {});
  },
};
