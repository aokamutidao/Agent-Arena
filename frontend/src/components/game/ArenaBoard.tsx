"use client";

import { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import type { GameDetail } from "@/types/game";
import { cn } from "@/lib/utils";

interface ArenaBoardProps {
  gameState: GameDetail;
}

// 默认障碍物（来自 engine/constants.go）
const DEFAULT_OBSTACLES = [
  { x: 2, y: 3 }, { x: 2, y: 7 },
  { x: 4, y: 5 }, { x: 6, y: 5 },
  { x: 5, y: 2 }, { x: 7, y: 3 },
  { x: 7, y: 7 }, { x: 3, y: 5 },
];

const CELL_SIZE = 48;

// 网格坐标 → 像素（Y 轴翻转，0 在底部）
const px = (gx: number) => gx * CELL_SIZE;
const py = (gy: number) => (9 - gy) * CELL_SIZE;

export function ArenaBoard({ gameState }: ArenaBoardProps) {
  const [redPos, setRedPos] = useState({ x: 1, y: 1 });
  const [bluePos, setBluePos] = useState({ x: 8, y: 8 });
  const [shake, setShake] = useState({ x: 0, y: 0 });

  // 只在坐标真正变化时更新 local state（触发 framer-motion 动画）
  useEffect(() => {
    const rp = gameState.agent_red_state?.position;
    if (rp) {
      setRedPos((prev) => (prev.x !== rp.x || prev.y !== rp.y ? { x: rp.x, y: rp.y } : prev));
    }
  }, [gameState.agent_red_state?.position?.x, gameState.agent_red_state?.position?.y]);

  useEffect(() => {
    const bp = gameState.agent_blue_state?.position;
    if (bp) {
      setBluePos((prev) => (prev.x !== bp.x || prev.y !== bp.y ? { x: bp.x, y: bp.y } : prev));
    }
  }, [gameState.agent_blue_state?.position?.x, gameState.agent_blue_state?.position?.y]);

  const redHP = gameState.agent_red_state?.hp ?? 100;
  const blueHP = gameState.agent_blue_state?.hp ?? 100;

  const isObstacle = (x: number, y: number) =>
    DEFAULT_OBSTACLES.some((o) => o.x === x && o.y === y);

  // 最新回合（用于显示特效）
  const lastTurn = gameState.history?.[gameState.history.length - 1];
  const isFinished = gameState.status === "finished";
  const redDead = redHP === 0;
  const blueDead = blueHP === 0;

  // 特效 key（每回合重置）
  const effectKey = lastTurn ? `${lastTurn.round}` : "0";

  // 屏幕震动（ATTACK/SKILL 命中时）
  useEffect(() => {
    if (!lastTurn || isFinished) return;
    const hasHit =
      (lastTurn.red_action?.type === "ATTACK" && !lastTurn.red_action.failed) ||
      (lastTurn.blue_action?.type === "ATTACK" && !lastTurn.blue_action.failed) ||
      (lastTurn.red_action?.type === "SKILL" && !lastTurn.red_action.failed) ||
      (lastTurn.blue_action?.type === "SKILL" && !lastTurn.blue_action.failed);

    if (hasHit) {
      setShake({ x: 4, y: 2 });
      const timer = setTimeout(() => setShake({ x: 0, y: 0 }), 300);
      return () => clearTimeout(timer);
    }
  }, [lastTurn, isFinished]);

  return (
    <div className="inline-block relative">
      {/* Column labels */}
      <div className="flex ml-8">
        {Array.from({ length: 10 }, (_, i) => (
          <div key={i} className="w-12 h-6 flex items-center justify-center text-xs text-muted-foreground">
            {i}
          </div>
        ))}
      </div>

      <div className="flex">
        {/* Row labels */}
        <div className="w-8 flex flex-col-reverse">
          {Array.from({ length: 10 }, (_, i) => (
            <div key={i} className="w-8 h-12 flex items-center justify-center text-xs text-muted-foreground">
              {i}
            </div>
          ))}
        </div>

        {/* Grid area */}
        <motion.div
          className="relative border border-border/50"
          style={{ width: CELL_SIZE * 10, height: CELL_SIZE * 10 }}
          animate={{ x: shake.x, y: shake.y }}
          transition={{ duration: 0.1, ease: "easeInOut" }}
        >
          {/* Grid cells */}
          {Array.from({ length: 10 }, (_, row) => {
            const y = 9 - row;
            return Array.from({ length: 10 }, (_, col) => {
              const x = col;
              const obstacle = isObstacle(x, y);
              return (
                <div
                  key={`${x}-${y}`}
                  className={cn(
                    "absolute border border-border/30 flex items-center justify-center",
                    obstacle && "bg-muted",
                    !obstacle && "bg-background hover:bg-accent/20"
                  )}
                  style={{
                    left: px(x),
                    top: py(y),
                    width: CELL_SIZE,
                    height: CELL_SIZE,
                  }}
                >
                  {obstacle && (
                    <div className="w-8 h-8 rounded bg-muted-foreground/30" />
                  )}
                </div>
              );
            });
          })}

          {/* Red Agent — framer-motion 平滑移动 */}
          <motion.div
            className="absolute z-10 flex items-center justify-center pointer-events-none"
            animate={{ left: px(redPos.x), top: py(redPos.y) }}
            transition={{ type: "spring", stiffness: 300, damping: 25 }}
            style={{ width: CELL_SIZE, height: CELL_SIZE }}
          >
            <motion.div
              className={cn(
                "w-10 h-10 rounded-full flex items-center justify-center shadow-lg ring-2 ring-red-400/50",
                redDead ? "bg-gray-600 opacity-50" : "bg-red-500 shadow-red-500/30"
              )}
              animate={redDead ? { scale: 0.6, opacity: 0.4 } : { scale: 1, opacity: 1 }}
              transition={{ duration: 0.5 }}
            >
              <span className="text-white text-xs font-bold">
                {redDead ? "✗" : "R"}
              </span>
            </motion.div>
          </motion.div>

          {/* Blue Agent — framer-motion 平滑移动 */}
          <motion.div
            className="absolute z-10 flex items-center justify-center pointer-events-none"
            animate={{ left: px(bluePos.x), top: py(bluePos.y) }}
            transition={{ type: "spring", stiffness: 300, damping: 25 }}
            style={{ width: CELL_SIZE, height: CELL_SIZE }}
          >
            <motion.div
              className={cn(
                "w-10 h-10 rounded-full flex items-center justify-center shadow-lg ring-2 ring-blue-400/50",
                blueDead ? "bg-gray-600 opacity-50" : "bg-blue-500 shadow-blue-500/30"
              )}
              animate={blueDead ? { scale: 0.6, opacity: 0.4 } : { scale: 1, opacity: 1 }}
              transition={{ duration: 0.5 }}
            >
              <span className="text-white text-xs font-bold">
                {blueDead ? "✗" : "B"}
              </span>
            </motion.div>
          </motion.div>

          {/* 特效层 */}
          <AnimatePresence>
            {lastTurn && !isFinished && (
              <>
                {/* 攻击特效 */}
                {lastTurn.red_action?.type === "ATTACK" && !lastTurn.red_action.failed && (
                  <AttackFX key={`ra-${effectKey}`} fromX={redPos.x} fromY={redPos.y} toX={bluePos.x} toY={bluePos.y} color="red" />
                )}
                {lastTurn.blue_action?.type === "ATTACK" && !lastTurn.blue_action.failed && (
                  <AttackFX key={`ba-${effectKey}`} fromX={bluePos.x} fromY={bluePos.y} toX={redPos.x} toY={redPos.y} color="blue" />
                )}
                {/* 技能特效 */}
                {lastTurn.red_action?.type === "SKILL" && !lastTurn.red_action.failed && (
                  <SkillFX key={`rs-${effectKey}`} fromX={redPos.x} fromY={redPos.y} toX={bluePos.x} toY={bluePos.y} color="red" />
                )}
                {lastTurn.blue_action?.type === "SKILL" && !lastTurn.blue_action.failed && (
                  <SkillFX key={`bs-${effectKey}`} fromX={bluePos.x} fromY={bluePos.y} toX={redPos.x} toY={redPos.y} color="blue" />
                )}
                {/* 蓄力特效 */}
                {lastTurn.red_action?.type === "CHARGE" && !lastTurn.red_action.failed && (
                  <ChargeFX key={`rc-${effectKey}`} x={redPos.x} y={redPos.y} color="red" />
                )}
                {lastTurn.blue_action?.type === "CHARGE" && !lastTurn.blue_action.failed && (
                  <ChargeFX key={`bc-${effectKey}`} x={bluePos.x} y={bluePos.y} color="blue" />
                )}
                {/* 治疗特效 */}
                {lastTurn.red_action?.type === "HEAL" && !lastTurn.red_action.failed && (
                  <HealFX key={`rh-${effectKey}`} x={redPos.x} y={redPos.y} color="red" />
                )}
                {lastTurn.blue_action?.type === "HEAL" && !lastTurn.blue_action.failed && (
                  <HealFX key={`bh-${effectKey}`} x={bluePos.x} y={bluePos.y} color="blue" />
                )}
                {/* 移动特效 */}
                {lastTurn.red_action?.type === "MOVE" && !lastTurn.red_action.failed && lastTurn.red_action.target && (
                  <MoveFX key={`rm-${effectKey}`} fromX={redPos.x} fromY={redPos.y} toX={lastTurn.red_action.target.x} toY={lastTurn.red_action.target.y} color="red" />
                )}
                {lastTurn.blue_action?.type === "MOVE" && !lastTurn.blue_action.failed && lastTurn.blue_action.target && (
                  <MoveFX key={`bm-${effectKey}`} fromX={bluePos.x} fromY={bluePos.y} toX={lastTurn.blue_action.target.x} toY={lastTurn.blue_action.target.y} color="blue" />
                )}
                {/* 失败特效 */}
                {lastTurn.red_action?.failed && (
                  <FailFX key={`rf-${effectKey}`} x={redPos.x} y={redPos.y} reason={lastTurn.red_action.fail_reason} />
                )}
                {lastTurn.blue_action?.failed && (
                  <FailFX key={`bf-${effectKey}`} x={bluePos.x} y={bluePos.y} reason={lastTurn.blue_action.fail_reason} />
                )}
              </>
            )}
          </AnimatePresence>

          {/* 死亡覆盖层 */}
          <AnimatePresence>
            {isFinished && (redDead || blueDead) && (
              <motion.div
                className="absolute inset-0 z-20 flex items-center justify-center bg-black/40 rounded"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.5 }}
              >
                <motion.div
                  className="text-3xl font-bold text-white drop-shadow-lg"
                  initial={{ scale: 0, y: 20 }}
                  animate={{ scale: 1, y: 0 }}
                  transition={{ type: "spring", stiffness: 200, delay: 0.2 }}
                >
                  {(redDead && blueDead) ? "💀 同归于尽" :
                    redDead ? "💀 红方阵亡" : "💀 蓝方阵亡"}
                </motion.div>
              </motion.div>
            )}
          </AnimatePresence>
        </motion.div>
      </div>
    </div>
  );
}

// ===== 攻击特效（增强版）=====
function AttackFX({ fromX, fromY, toX, toY, color }: {
  fromX: number; fromY: number; toX: number; toY: number; color: "red" | "blue";
}) {
  const bg = color === "red" ? "bg-red-500" : "bg-blue-500";
  const shadowColor = color === "red" ? "shadow-red-500" : "shadow-blue-500";
  // 计算攻击方向角度
  const dx = (toX - fromX) * CELL_SIZE;
  const dy = (fromY - toY) * CELL_SIZE;
  const angle = Math.atan2(dy, dx) * (180 / Math.PI);

  return (
    <>
      {/* 攻击者蓄力闪光 */}
      <motion.div
        className={cn("absolute rounded-full", bg)}
        style={{ left: px(fromX) + 4, top: py(fromY) + 4, width: CELL_SIZE - 8, height: CELL_SIZE - 8 }}
        initial={{ opacity: 0.9, scale: 1 }}
        animate={{ opacity: 0, scale: 2.5 }}
        exit={{ opacity: 0 }}
        transition={{ duration: 0.4 }}
      />
      {/* 方向性斩击 - 旋转的斩击线 */}
      <motion.div
        className={cn("absolute h-1 rounded-full", bg)}
        style={{
          left: px(fromX) + CELL_SIZE / 2,
          top: py(fromY) + CELL_SIZE / 2,
          width: 0,
          transformOrigin: "left center",
          rotate: `${angle}deg`,
        }}
        initial={{ width: 0, opacity: 1 }}
        animate={{ width: CELL_SIZE * 1.5, opacity: [1, 1, 0] }}
        exit={{ opacity: 0 }}
        transition={{ duration: 0.3, ease: "easeOut" }}
      />
      {/* 命中爆炸 - 更大的冲击波 */}
      <motion.div
        className={cn("absolute rounded-full", bg)}
        style={{ left: px(toX) - 4, top: py(toY) - 4, width: CELL_SIZE + 8, height: CELL_SIZE + 8 }}
        initial={{ opacity: 0, scale: 0.2 }}
        animate={{ opacity: [0, 0.9, 0], scale: [0.2, 1.8, 1.2] }}
        exit={{ opacity: 0 }}
        transition={{ duration: 0.5, delay: 0.15 }}
      />
      {/* 第二层冲击波 */}
      <motion.div
        className="absolute rounded-full border-2 border-white/60"
        style={{ left: px(toX), top: py(toY), width: CELL_SIZE, height: CELL_SIZE }}
        initial={{ opacity: 0, scale: 0.3 }}
        animate={{ opacity: [0, 0.8, 0], scale: [0.3, 2, 1.5] }}
        exit={{ opacity: 0 }}
        transition={{ duration: 0.6, delay: 0.2 }}
      />
      {/* 伤害数字 + 爆炸 emoji */}
      <motion.div
        className="absolute z-30 flex items-center gap-1"
        style={{ left: px(toX) + 8, top: py(toY) - 12 }}
        initial={{ opacity: 1, y: 0, scale: 0.5 }}
        animate={{ opacity: 0, y: -25, scale: 1.2 }}
        exit={{ opacity: 0 }}
        transition={{ duration: 0.7, delay: 0.2 }}
      >
        <span className="text-2xl">💥</span>
      </motion.div>
    </>
  );
}

// ===== 技能特效（增强版）=====
function SkillFX({ fromX, fromY, toX, toY, color }: {
  fromX: number; fromY: number; toX: number; toY: number; color: "red" | "blue";
}) {
  const borderC = color === "red" ? "border-red-400" : "border-blue-400";
  const bgC = color === "red" ? "bg-red-400" : "bg-blue-400";
  const glowC = color === "red" ? "bg-red-500/40" : "bg-blue-500/40";

  return (
    <>
      {/* 施法阵 - 旋转的虚线圆环 */}
      <motion.div
        className={cn("absolute rounded-full border-2 border-dashed", borderC)}
        style={{ left: px(fromX) - 4, top: py(fromY) - 4, width: CELL_SIZE + 8, height: CELL_SIZE + 8 }}
        initial={{ opacity: 0, scale: 0.5, rotate: 0 }}
        animate={{ opacity: [0, 1, 0], scale: [0.5, 1.3, 1], rotate: 360 }}
        exit={{ opacity: 0 }}
        transition={{ duration: 0.6, ease: "easeOut" }}
      />
      {/* 施法光晕 */}
      <motion.div
        className={cn("absolute rounded-full", glowC)}
        style={{ left: px(fromX), top: py(fromY), width: CELL_SIZE, height: CELL_SIZE }}
        initial={{ opacity: 0, scale: 0.8 }}
        animate={{ opacity: [0, 0.6, 0], scale: [0.8, 1.5, 1.2] }}
        exit={{ opacity: 0 }}
        transition={{ duration: 0.5 }}
      />
      {/* 飞行弹丸 - 带拖尾 */}
      {[0, 0.08, 0.16].map((delay, i) => (
        <motion.div
          key={i}
          className={cn("absolute rounded-full", bgC)}
          style={{
            left: px(fromX) + 18 - i * 2,
            top: py(fromY) + 18 - i * 2,
            width: 12 - i * 2,
            height: 12 - i * 2,
            opacity: 1 - i * 0.3,
          }}
          initial={{ x: 0, y: 0 }}
          animate={{
            x: (toX - fromX) * CELL_SIZE,
            y: (fromY - toY) * CELL_SIZE,
            opacity: [1 - i * 0.3, 1 - i * 0.3, 0],
          }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.5, delay, ease: "easeIn" }}
        />
      ))}
      {/* 命中爆炸 - 多层环形 */}
      <motion.div
        className={cn("absolute rounded-full", bgC)}
        style={{ left: px(toX) - 8, top: py(toY) - 8, width: CELL_SIZE + 16, height: CELL_SIZE + 16 }}
        initial={{ opacity: 0, scale: 0 }}
        animate={{ opacity: [0, 0.8, 0], scale: [0, 1.5, 1] }}
        exit={{ opacity: 0 }}
        transition={{ duration: 0.5, delay: 0.4 }}
      />
      <motion.div
        className="absolute rounded-full border-2 border-white/70"
        style={{ left: px(toX), top: py(toY), width: CELL_SIZE, height: CELL_SIZE }}
        initial={{ opacity: 0, scale: 0.2 }}
        animate={{ opacity: [0, 1, 0], scale: [0.2, 2, 1.8] }}
        exit={{ opacity: 0 }}
        transition={{ duration: 0.6, delay: 0.45 }}
      />
      {/* 光柱效果 */}
      <motion.div
        className={cn("absolute w-2", bgC)}
        style={{ left: px(toX) + CELL_SIZE / 2 - 4, top: py(toY) - 50, height: 0 }}
        initial={{ height: 0, opacity: 0 }}
        animate={{ height: 100, opacity: [0, 0.9, 0] }}
        exit={{ opacity: 0 }}
        transition={{ duration: 0.6, delay: 0.35 }}
      />
    </>
  );
}

// ===== 蓄力特效 =====
function ChargeFX({ x, y, color }: { x: number; y: number; color: "red" | "blue" }) {
  const bg = color === "red" ? "bg-red-500/60" : "bg-blue-500/60";
  const border = color === "red" ? "border-red-400" : "border-blue-400";

  return (
    <>
      {/* 脉冲光环 - 多层扩散 */}
      {[0, 0.15, 0.3].map((delay, i) => (
        <motion.div
          key={i}
          className={cn("absolute rounded-full border-2", border)}
          style={{
            left: px(x) + 4 - i * 4,
            top: py(y) + 4 - i * 4,
            width: CELL_SIZE - 8 + i * 8,
            height: CELL_SIZE - 8 + i * 8,
          }}
          initial={{ opacity: 0, scale: 0.8 }}
          animate={{ opacity: [0, 0.8, 0], scale: [0.8, 1.3, 1.6] }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.8, delay }}
        />
      ))}
      {/* 内部能量球 */}
      <motion.div
        className={cn("absolute rounded-full", bg)}
        style={{ left: px(x) + 8, top: py(y) + 8, width: CELL_SIZE - 16, height: CELL_SIZE - 16 }}
        initial={{ opacity: 0.5, scale: 0.8 }}
        animate={{ opacity: [0.5, 1, 0.5], scale: [0.8, 1.2, 0.8] }}
        exit={{ opacity: 0 }}
        transition={{ duration: 0.6, repeat: 2, repeatType: "reverse" }}
      />
      {/* 闪电 emoji */}
      <motion.div
        className="absolute z-30 text-2xl"
        style={{ left: px(x) + 12, top: py(y) - 20 }}
        initial={{ opacity: 0, y: 10, scale: 0.5 }}
        animate={{ opacity: [0, 1, 1, 0], y: [10, 0, -5, -15], scale: [0.5, 1.2, 1, 0.8] }}
        exit={{ opacity: 0 }}
        transition={{ duration: 0.8 }}
      >
        ⚡
      </motion.div>
    </>
  );
}

// ===== 治疗特效 =====
function HealFX({ x, y, color }: { x: number; y: number; color: "red" | "blue" }) {
  const bg = color === "red" ? "bg-green-400/40" : "bg-green-400/40";
  const border = "border-green-400";

  return (
    <>
      {/* 绿色治疗光环 */}
      <motion.div
        className={cn("absolute rounded-full border-2", border)}
        style={{ left: px(x) + 2, top: py(y) + 2, width: CELL_SIZE - 4, height: CELL_SIZE - 4 }}
        initial={{ opacity: 0, scale: 0.5 }}
        animate={{ opacity: [0, 0.8, 0], scale: [0.5, 1.2, 1.5] }}
        exit={{ opacity: 0 }}
        transition={{ duration: 0.6 }}
      />
      {/* 内部绿色光晕 */}
      <motion.div
        className={cn("absolute rounded-full", bg)}
        style={{ left: px(x) + 6, top: py(y) + 6, width: CELL_SIZE - 12, height: CELL_SIZE - 12 }}
        initial={{ opacity: 0, scale: 0.6 }}
        animate={{ opacity: [0, 0.6, 0], scale: [0.6, 1.1, 0.8] }}
        exit={{ opacity: 0 }}
        transition={{ duration: 0.5, delay: 0.1 }}
      />
      {/* 上升的 + 号粒子 */}
      {["+", "+", "+"].map((symbol, i) => (
        <motion.div
          key={i}
          className="absolute text-green-400 font-bold text-lg"
          style={{ left: px(x) + 10 + i * 10, top: py(y) + 20 }}
          initial={{ opacity: 0, y: 0 }}
          animate={{ opacity: [0, 1, 0], y: -30 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.8, delay: i * 0.1 }}
        >
          {symbol}
        </motion.div>
      ))}
      {/* 治疗数值 */}
      <motion.div
        className="absolute z-30 text-green-400 font-bold text-sm"
        style={{ left: px(x) + 12, top: py(y) - 8 }}
        initial={{ opacity: 1, y: 0 }}
        animate={{ opacity: 0, y: -20 }}
        exit={{ opacity: 0 }}
        transition={{ duration: 0.7, delay: 0.2 }}
      >
        +25
      </motion.div>
    </>
  );
}

// ===== 移动特效 =====
function MoveFX({ fromX, fromY, toX, toY, color }: {
  fromX: number; fromY: number; toX: number; toY: number; color: "red" | "blue";
}) {
  const bg = color === "red" ? "bg-red-400/30" : "bg-blue-400/30";

  return (
    <>
      {/* 残影 - 在起点留下模糊的影子 */}
      <motion.div
        className={cn("absolute rounded-full", bg)}
        style={{ left: px(fromX) + 8, top: py(fromY) + 8, width: CELL_SIZE - 16, height: CELL_SIZE - 16 }}
        initial={{ opacity: 0.6, scale: 1 }}
        animate={{ opacity: 0, scale: 0.5 }}
        exit={{ opacity: 0 }}
        transition={{ duration: 0.4 }}
      />
      {/* 尘土飞扬效果 */}
      {[0, 0.1, 0.2].map((delay, i) => (
        <motion.div
          key={i}
          className="absolute rounded-full bg-amber-600/40"
          style={{
            left: px(fromX) + 16 + (i - 1) * 8,
            top: py(fromY) + 20,
            width: 6,
            height: 6,
          }}
          initial={{ opacity: 0.6, scale: 0.5, y: 0 }}
          animate={{ opacity: 0, scale: 1.5, y: -10 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.5, delay }}
        />
      ))}
      {/* 落点冲击波 */}
      <motion.div
        className="absolute rounded-full border border-white/40"
        style={{ left: px(toX) + 4, top: py(toY) + 4, width: CELL_SIZE - 8, height: CELL_SIZE - 8 }}
        initial={{ opacity: 0, scale: 0.3 }}
        animate={{ opacity: [0, 0.6, 0], scale: [0.3, 1.2, 1.5] }}
        exit={{ opacity: 0 }}
        transition={{ duration: 0.4, delay: 0.2 }}
      />
    </>
  );
}

// ===== 失败特效 =====
function FailFX({ x, y, reason }: { x: number; y: number; reason?: string }) {
  return (
    <>
      <motion.div
        className="absolute z-30 flex flex-col items-center pointer-events-none"
        style={{ left: px(x), top: py(y) - 20, width: CELL_SIZE }}
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: [0, 1, 1, 0], y: [10, 0, 0, -15] }}
        exit={{ opacity: 0 }}
        transition={{ duration: 1.2 }}
      >
        <span className="text-red-500 font-bold text-lg">✗</span>
        {reason && (
          <span className="text-[10px] text-red-400 whitespace-nowrap bg-black/60 rounded px-1 mt-0.5">
            {reason}
          </span>
        )}
      </motion.div>
    </>
  );
}
