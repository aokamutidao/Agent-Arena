"use client";

import { useEffect, useRef, useState, useCallback } from "react";

interface WSMessage {
  type: string;
  data: unknown;
}

interface UseWebSocketReturn {
  connected: boolean;
  gameState: GameStateSnapshot | null;
  turnHistory: TurnUpdate[];
  bettingUpdate: BettingSnapshot | null;
  gameFinished: GameFinishedEvent | null;
  strategyUpdate: StrategyUpdateEvent | null;
  overtimeStarted: OvertimeEvent | null;
  balanceUpdate: BalanceUpdateEvent | null;
}

interface GameStateSnapshot {
  game_id: number;
  status: string;
  current_round: number;
  agent_red_hp: number;
  agent_blue_hp: number;
  agent_red_pos_x: number;
  agent_red_pos_y: number;
  agent_blue_pos_x: number;
  agent_blue_pos_y: number;
}

interface TurnUpdate {
  round: number;
  red_action: { type: string; target?: { x: number; y: number }; failed?: boolean; fail_reason?: string };
  blue_action: { type: string; target?: { x: number; y: number }; failed?: boolean; fail_reason?: string };
  red_hp: number;
  blue_hp: number;
  red_pos_x: number;
  red_pos_y: number;
  blue_pos_x: number;
  blue_pos_y: number;
  red_reasoning?: string;
  blue_reasoning?: string;
}

interface BettingSnapshot {
  total_bet_red: string;
  total_bet_blue: string;
  odds_red: number;
  odds_blue: number;
}

interface GameFinishedEvent {
  game_id: number;
  winner: string;
  winner_name: string;
  total_rounds: number;
  final_hp_red: number;
  final_hp_blue: number;
}

interface StrategyUpdateEvent {
  side: string;
  aggressive: number;
  defensive: number;
  tricky: number;
}

interface OvertimeEvent {
  game_id: number;
  extra_rounds: number;
  overtime_dmg: number;
}

interface BalanceUpdateEvent {
  address: string;
  ac_balance: number;
  reason: string;
}

export function useWebSocket(gameId: number): UseWebSocketReturn {
  const [connected, setConnected] = useState(false);
  const [gameState, setGameState] = useState<GameStateSnapshot | null>(null);
  const [turnHistory, setTurnHistory] = useState<TurnUpdate[]>([]);
  const [bettingUpdate, setBettingUpdate] = useState<BettingSnapshot | null>(null);
  const [gameFinished, setGameFinished] = useState<GameFinishedEvent | null>(null);
  const [strategyUpdate, setStrategyUpdate] = useState<StrategyUpdateEvent | null>(null);
  const [overtimeStarted, setOvertimeStarted] = useState<OvertimeEvent | null>(null);
  const [balanceUpdate, setBalanceUpdate] = useState<BalanceUpdateEvent | null>(null);
  const wsRef = useRef<WebSocket | null>(null);

  const WS_BASE = process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8080";

  const handleMessage = useCallback((event: MessageEvent) => {
    try {
      const msg: WSMessage = JSON.parse(event.data);

      switch (msg.type) {
        case "game_state":
          setGameState(msg.data as GameStateSnapshot);
          break;
        case "turn_update":
          setTurnHistory((prev) => [msg.data as TurnUpdate, ...prev].slice(0, 10));
          break;
        case "betting_update":
          setBettingUpdate(msg.data as BettingSnapshot);
          break;
        case "game_started":
          // Could update status
          break;
        case "game_finished":
          setGameFinished(msg.data as GameFinishedEvent);
          break;
        case "strategy_update":
          setStrategyUpdate(msg.data as StrategyUpdateEvent);
          break;
        case "overtime_started":
          setOvertimeStarted(msg.data as OvertimeEvent);
          break;
        case "user_balance_update":
          setBalanceUpdate(msg.data as BalanceUpdateEvent);
          break;
      }
    } catch {
      // Ignore parse errors
    }
  }, []);

  useEffect(() => {
    if (!gameId) return;

    const ws = new WebSocket(`${WS_BASE}/ws?game_id=${gameId}`);
    wsRef.current = ws;

    ws.onopen = () => setConnected(true);
    ws.onclose = () => setConnected(false);
    ws.onerror = () => setConnected(false);
    ws.onmessage = handleMessage;

    return () => {
      ws.close();
      wsRef.current = null;
    };
  }, [gameId, WS_BASE, handleMessage]);

  return { connected, gameState, turnHistory, bettingUpdate, gameFinished, strategyUpdate, overtimeStarted, balanceUpdate };
}
