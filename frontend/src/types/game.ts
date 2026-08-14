// 类型定义 - 映射自后端 API types.go

export interface Position {
  x: number;
  y: number;
}

export interface AgentInfo {
  id: string;
  name: string;
  personality: string;
  wins: number;
  losses: number;
  win_rate: number;
  description: string;
}

export interface AgentState {
  hp: number;
  max_hp: number;
  position: Position;
  status: string[];
  skill_cooldown: number;
  is_charging: boolean;
}

export interface StrategyWeights {
  aggressive: number;
  defensive: number;
  tricky: number;
}

export interface TurnRecord {
  round: number;
  red_action: { type: string; target?: { x: number; y: number }; failed?: boolean; fail_reason?: string };
  blue_action: { type: string; target?: { x: number; y: number }; failed?: boolean; fail_reason?: string };
  red_hp_after: number;
  blue_hp_after: number;
  red_reasoning?: string;
  blue_reasoning?: string;
}

export interface GameListItem {
  game_id: number;
  agent_red: AgentInfo;
  agent_blue: AgentInfo;
  status: string;
  total_bet_red: string;
  total_bet_blue: string;
  current_round: number;
  max_rounds: number;
  odds_red: number;
  odds_blue: number;
  winner?: string;
}

export interface GamesResponse {
  games: GameListItem[];
  total: number;
}

export interface AgentsResponse {
  agents: AgentInfo[];
}

export interface GameDetail {
  game_id: number;
  agent_red: AgentInfo;
  agent_blue: AgentInfo;
  status: string;
  current_round: number;
  max_rounds: number;
  total_bet_red: string;
  total_bet_blue: string;
  agent_red_state?: AgentState;
  agent_blue_state?: AgentState;
  strategy_red?: StrategyWeights;
  strategy_blue?: StrategyWeights;
  overtime?: boolean;
  history: TurnRecord[];
  winner?: string;
  archived?: boolean; // 重启后从 DB 兜底恢复的已完成对局（无回合明细）
}

export interface EstimateRequest {
  game_id: number;
  side: "red" | "blue";
  amount: string;
}

export interface EstimateResponse {
  current_pool_red: string;
  current_pool_blue: string;
  new_pool_red: string;
  potential_reward: string;
  new_odds_red: number;
  new_odds_blue: number;
}

export interface VoteRequest {
  game_id: number;
  side: "red" | "blue";
  strategy: "aggressive" | "defensive" | "tricky";
  user?: string;
}

export interface VoteResponse {
  side: string;
  aggressive: number;
  defensive: number;
  tricky: number;
}
