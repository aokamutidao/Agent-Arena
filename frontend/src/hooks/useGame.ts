"use client";

import { useEffect, useState, useRef, useCallback } from "react";
import { api } from "@/lib/api";
import { useWebSocket } from "./useWebSocket";
import type { GameDetail, AgentState, TurnRecord } from "@/types/game";

interface UseGameReturn {
  game: GameDetail | null;
  loading: boolean;
  error: string | null;
  wsConnected: boolean;
  turnHistory: Array<{
    round: number;
    red_action: { type: string; target?: { x: number; y: number } };
    blue_action: { type: string; target?: { x: number; y: number } };
    red_hp: number;
    blue_hp: number;
  }>;
  bettingUpdate: {
    total_bet_red: string;
    total_bet_blue: string;
    odds_red: number;
    odds_blue: number;
  } | null;
  gameFinished: {
    game_id: number;
    winner: string;
    winner_name: string;
    total_rounds: number;
    final_hp_red: number;
    final_hp_blue: number;
  } | null;
  overtimeStarted: {
    game_id: number;
    extra_rounds: number;
    overtime_dmg: number;
  } | null;
  balanceUpdate: {
    address: string;
    ac_balance: number;
    reason: string;
  } | null;
}

export function useGame(gameId: number): UseGameReturn {
  const [game, setGame] = useState<GameDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const {
    connected: wsConnected,
    gameState,
    turnHistory,
    bettingUpdate,
    gameFinished,
    strategyUpdate,
    overtimeStarted,
    balanceUpdate,
  } = useWebSocket(gameId);

  // 初始加载 REST 数据
  useEffect(() => {
    if (!gameId) return;

    setLoading(true);
    setError(null);

    api
      .getGame(gameId)
      .then((data) => {
        setGame(data);
        setLoading(false);
      })
      .catch((err) => {
        setError(err.message || "Failed to load game");
        setLoading(false);
      });
  }, [gameId]);

  // REST 轮询兜底：每 3 秒拉一次最新状态（WS 可能断连/丢消息）
  useEffect(() => {
    if (!gameId || !game) return;
    if (game.status === "finished") return;

    const interval = setInterval(() => {
      api.getGame(gameId).then((data) => {
        setGame((prev) => {
          if (!prev) return prev;
          // 只在数据有变化时更新
          if (
            data.current_round !== prev.current_round ||
            data.status !== prev.status ||
            data.history.length !== prev.history.length
          ) {
            return data;
          }
          return prev;
        });
      }).catch(() => {
        // 轮询失败不处理
      });
    }, 3000);

    return () => clearInterval(interval);
  }, [gameId, game?.status]);

  // WS game_state 快照（初始连接时发送一次）
  useEffect(() => {
    if (!game || !gameState) return;

    setGame((prev) => {
      if (!prev) return prev;
      return {
        ...prev,
        current_round: gameState.current_round,
        status: gameState.status,
        agent_red_state: {
          ...prev.agent_red_state,
          hp: gameState.agent_red_hp,
          position: {
            x: gameState.agent_red_pos_x,
            y: gameState.agent_red_pos_y,
          },
        } as AgentState,
        agent_blue_state: {
          ...prev.agent_blue_state,
          hp: gameState.agent_blue_hp,
          position: {
            x: gameState.agent_blue_pos_x,
            y: gameState.agent_blue_pos_y,
          },
        } as AgentState,
      };
    });
  }, [gameState]);

  // WS turn_update：实时更新位置/HP/回合 + 合并历史
  useEffect(() => {
    if (!game || turnHistory.length === 0) return;

    setGame((prev) => {
      if (!prev) return prev;

      let changed = false;
      const next = { ...prev };

      // 取最新的 turn_update（turnHistory 是倒序，[0] 最新）
      const latest = turnHistory[0];
      if (latest) {
        // 实时更新 agent 位置和 HP
        if (
          prev.agent_red_state?.hp !== latest.red_hp ||
          prev.agent_red_state?.position.x !== latest.red_pos_x ||
          prev.agent_red_state?.position.y !== latest.red_pos_y ||
          prev.agent_blue_state?.hp !== latest.blue_hp ||
          prev.agent_blue_state?.position.x !== latest.blue_pos_x ||
          prev.agent_blue_state?.position.y !== latest.blue_pos_y ||
          prev.current_round !== latest.round
        ) {
          next.agent_red_state = {
            ...prev.agent_red_state,
            hp: latest.red_hp,
            position: { x: latest.red_pos_x, y: latest.red_pos_y },
          } as AgentState;
          next.agent_blue_state = {
            ...prev.agent_blue_state,
            hp: latest.blue_hp,
            position: { x: latest.blue_pos_x, y: latest.blue_pos_y },
          } as AgentState;
          next.current_round = latest.round;
          changed = true;
        }
      }

      // 合并新的回合到历史（从 prev.history 推导已有回合，避免重复）
      const existingRounds = new Set(prev.history.map((t) => t.round));
      const newTurns = turnHistory
        .filter((t) => !existingRounds.has(t.round))
        .map(
          (t) =>
            ({
              round: t.round,
              red_action: t.red_action,
              blue_action: t.blue_action,
              red_hp_after: t.red_hp,
              blue_hp_after: t.blue_hp,
              red_reasoning: t.red_reasoning,
              blue_reasoning: t.blue_reasoning,
            }) as TurnRecord
        );

      if (newTurns.length > 0) {
        next.history = [...prev.history, ...newTurns].sort(
          (a, b) => a.round - b.round
        );
        changed = true;
      }

      return changed ? next : prev;
    });
  }, [turnHistory]);

  // 下注更新
  useEffect(() => {
    if (!game || !bettingUpdate) return;
    setGame((prev) => {
      if (!prev) return prev;
      return {
        ...prev,
        total_bet_red: bettingUpdate.total_bet_red,
        total_bet_blue: bettingUpdate.total_bet_blue,
      };
    });
  }, [bettingUpdate]);

  // 对局结束
  useEffect(() => {
    if (!game || !gameFinished) return;
    setGame((prev) => {
      if (!prev) return prev;
      return {
        ...prev,
        status: "finished",
        winner: gameFinished.winner,
        current_round: gameFinished.total_rounds,
        agent_red_state: {
          ...prev.agent_red_state,
          hp: gameFinished.final_hp_red,
        } as AgentState,
        agent_blue_state: {
          ...prev.agent_blue_state,
          hp: gameFinished.final_hp_blue,
        } as AgentState,
      };
    });
  }, [gameFinished]);

  // 策略权重更新
  useEffect(() => {
    if (!game || !strategyUpdate) return;
    setGame((prev) => {
      if (!prev) return prev;
      const weights = {
        aggressive: strategyUpdate.aggressive,
        defensive: strategyUpdate.defensive,
        tricky: strategyUpdate.tricky,
      };
      if (strategyUpdate.side === "red") {
        return { ...prev, strategy_red: weights };
      } else {
        return { ...prev, strategy_blue: weights };
      }
    });
  }, [strategyUpdate]);

  // 加时赛开始
  useEffect(() => {
    if (!game || !overtimeStarted) return;
    setGame((prev) => {
      if (!prev) return prev;
      return { ...prev, overtime: true };
    });
  }, [overtimeStarted]);

  return {
    game,
    loading,
    error,
    wsConnected,
    turnHistory,
    bettingUpdate,
    gameFinished,
    overtimeStarted,
    balanceUpdate,
  };
}
