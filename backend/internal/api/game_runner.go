package api

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"time"

	"agent-arena/backend/internal/agent"
	"agent-arena/backend/internal/blockchain"
	"agent-arena/backend/internal/engine"
	"agent-arena/backend/internal/user"
)

const (
	TurnInterval = 3 * time.Second // 每回合间隔（API 响应约 1 秒，留 2 秒缓冲）
)

// GameRunner 游戏循环管理器
type GameRunner struct {
	engine *engine.Engine
	store  *MemoryStore
	hub    *WSHub
	qwen   agent.LLMClient
	chain  blockchain.ChainService // 可选：链上交互

	// 游戏结束回调：(gameID, winner) → 由 main.go 注入奖励/挑战结算逻辑
	onGameFinish []func(gameID uint64, winner string)
	// 链上 FinishGame 交易哈希回调：由 main.go 注入以更新 game_history.finish_tx_hash
	onFinishTx []func(gameID uint64, txHash string)
	// 对局历史存储（可选，nil 时跳过持久化）
	historyStore *user.GameHistoryStore
}

// SetHistoryStore 注入对局历史存储
func (r *GameRunner) SetHistoryStore(s *user.GameHistoryStore) {
	r.historyStore = s
}

// NewGameRunner 创建游戏循环管理器
func NewGameRunner(eng *engine.Engine, store *MemoryStore, hub *WSHub, qwen agent.LLMClient) *GameRunner {
	return &GameRunner{
		engine:       eng,
		store:        store,
		hub:          hub,
		qwen:         qwen,
		onGameFinish: make([]func(gameID uint64, winner string), 0),
		onFinishTx:   make([]func(gameID uint64, txHash string), 0),
	}
}

// OnGameFinish 注册游戏结束回调
func (r *GameRunner) OnGameFinish(callback func(gameID uint64, winner string)) {
	r.onGameFinish = append(r.onGameFinish, callback)
}

// OnFinishTx 注册链上 FinishGame 交易回调（txHash 已确认）
func (r *GameRunner) OnFinishTx(callback func(gameID uint64, txHash string)) {
	r.onFinishTx = append(r.onFinishTx, callback)
}

// SetChainService 设置链上服务（可选）
func (r *GameRunner) SetChainService(chain blockchain.ChainService) {
	r.chain = chain
}

// StartGame 启动一个对局的游戏循环
func (r *GameRunner) StartGame(gameID uint64) {
	go r.run(gameID)
}

// CreateNewGame 创建新对局并存储
func (r *GameRunner) CreateNewGame(gameID uint64, redName, blueName, redPersonality, bluePersonality string) *engine.GameState {
	game := r.engine.NewGame(gameID, redName, blueName, redPersonality, bluePersonality)
	game.Status = engine.StatusBetting
	r.store.CreateGame(game)
	return game
}

// CreateNewGameWithAgent 创建新对局（带自定义 Agent API 配置）
func (r *GameRunner) CreateNewGameWithAgent(
	gameID uint64,
	redName, blueName,
	redPersonality, bluePersonality,
	redAPIEndpoint, redAPIKey, redModel string,
) *engine.GameState {
	game := r.engine.NewGame(gameID, redName, blueName, redPersonality, bluePersonality)
	game.Status = engine.StatusBetting

	// 存储红方（挑战者）的 API 配置
	if redAPIEndpoint != "" {
		game.AgentRedAPIEndpoint = redAPIEndpoint
		game.AgentRedAPIKey = redAPIKey
		game.AgentRedModel = redModel
	}

	r.store.CreateGame(game)
	return game
}

func (r *GameRunner) run(gameID uint64) {
	record, err := r.store.GetGame(gameID)
	if err != nil {
		log.Printf("[GameRunner] game %d not found: %v", gameID, err)
		return
	}

	state := record.Game

	// 状态转换: betting → playing
	state.Status = engine.StatusPlaying
	redWins := state.Winner == engine.SideRed
	log.Printf("[GameRunner] Game %d started: %s vs %s",
		gameID, state.AgentRed.Name, state.AgentBlue.Name)

	// 链上: StartGame — 同步执行（必须等链上下注锁定后才能开始游戏循环）
	// 如果不锁定就开赛，后续 settle() 会因 status != Locked 而 revert
	if r.chain != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
		defer cancel()

		// 带重试的 StartGame（最多 3 次，应对网络波动/nonce 冲突）
		var txHash string
		var startErr error
		for attempt := 1; attempt <= 3; attempt++ {
			txHash, startErr = r.chain.StartGame(ctx, gameID)
			if startErr == nil {
				log.Printf("[GameRunner] chain.StartGame tx: %s, gameID: %d (attempt %d)", txHash, gameID, attempt)
				break
			}
			log.Printf("[GameRunner] chain.StartGame attempt %d failed: %v", attempt, startErr)
			if attempt < 3 {
				time.Sleep(5 * time.Second)
			}
		}
		if startErr != nil {
			log.Printf("[GameRunner] FATAL: chain.StartGame failed after 3 attempts: %v", startErr)
			log.Printf("[GameRunner] Game %d will run without on-chain lock — settlement will likely fail", gameID)
			// 继续游戏（本地状态仍然 playing），但 warn 用户链上结算可能失败
			// 广播错误信息给前端
			r.hub.Broadcast(gameID, "chain_error", map[string]interface{}{
				"game_id": gameID,
				"step":    "start_game",
				"error":   startErr.Error(),
				"message": "链上锁定下注失败，游戏将正常进行但结算可能受影响",
			})
		}
	}

	// 创建 Agent - 红方（自定义 Agent）必须使用外部 API
	if state.AgentRedAPIEndpoint == "" {
		log.Printf("[GameRunner] ERROR: Red agent (custom) has no API endpoint configured")
		// 回退到默认 Qwen（理论上不应该发生）
		state.AgentRedAPIEndpoint = ""
	}

	var redLLM agent.LLMClient
	if state.AgentRedAPIEndpoint != "" {
		log.Printf("[GameRunner] Red agent using external API: %s", state.AgentRedAPIEndpoint)
		redLLM = agent.NewExternalAPIClient(state.AgentRedAPIEndpoint, state.AgentRedAPIKey, state.AgentRedModel)
	} else {
		log.Printf("[GameRunner] WARNING: Red agent falling back to default Qwen (should not happen)")
		redLLM = r.qwen
	}

	redAgent := agent.NewAgent(state.AgentRed.Name, state.AgentRed.Personality, redLLM)
	blueAgent := agent.NewAgent(state.AgentBlue.Name, state.AgentBlue.Personality, r.qwen)

	// 加时赛广播追踪
	overtimeAnnounced := state.Overtime // 如果已经处于加时赛，不重复广播

	// 广播 game_started
	r.hub.Broadcast(gameID, "game_started", GameStartedData{
		GameID:        gameID,
		Status:        string(state.Status),
		BettingLocked: true,
	})

	// 游戏循环
	for state.Status == engine.StatusPlaying {
		time.Sleep(TurnInterval)

		// 构建决策输入
		redInput := agent.AgentDecisionInput{
			Self:          state.AgentRed,
			Enemy:         state.AgentBlue,
			CurrentRound:  state.CurrentRound,
			MaxRounds:     state.MaxRounds,
			Obstacles:     state.Obstacles,
			Weights:       record.StrategyRed,
			ActionHistory: state.History,
		}
		blueInput := agent.AgentDecisionInput{
			Self:          state.AgentBlue,
			Enemy:         state.AgentRed,
			CurrentRound:  state.CurrentRound,
			MaxRounds:     state.MaxRounds,
			Obstacles:     state.Obstacles,
			Weights:       record.StrategyBlue,
			ActionHistory: state.History,
		}

		// 并行决策
		log.Printf("[GameRunner] Round %d: agents deciding...", state.CurrentRound+1)
		redResult, blueResult := agent.DecideTurnsParallel(redAgent, blueAgent, redInput, blueInput)

		// 执行回合
		turnRecord, err := r.engine.ExecuteTurn(state, redResult.Action, blueResult.Action)
		if err != nil {
			log.Printf("[GameRunner] ExecuteTurn error: %v", err)
			break
		}

		// 附加推理思路到回合记录
		turnRecord.RedReasoning = redResult.Reasoning
		turnRecord.BlueReasoning = blueResult.Reasoning

		// 检测加时赛开始
		if state.Overtime && !overtimeAnnounced {
			overtimeAnnounced = true
			log.Printf("[GameRunner] Game %d entered OVERTIME! Both agents take 10 HP damage per round", gameID)
			r.hub.Broadcast(gameID, "overtime_started", OvertimeStartedData{
				GameID:      gameID,
				ExtraRounds: 10,
				OvertimeDmg: 10,
			})
		}

		log.Printf("[GameRunner] Round %d: Red=%s Blue=%s | HP: %d/%d",
			turnRecord.Round,
			redResult.Action.Type, blueResult.Action.Type,
			state.AgentRed.HP, state.AgentBlue.HP)

		// 广播 turn_update
		r.hub.Broadcast(gameID, "turn_update", toTurnUpdateData(turnRecord, state))

		// 更新 redWins（每回合重新计算）
		redWins = state.Winner == engine.SideRed

		// 检查对局结束
		if state.Status == engine.StatusFinished {
			winnerName := ""
			if state.Winner == engine.SideRed {
				winnerName = state.AgentRed.Name
			} else if state.Winner == engine.SideBlue {
				winnerName = state.AgentBlue.Name
			}

			log.Printf("[GameRunner] Game %d finished! Winner: %s (Rounds: %d)",
				gameID, winnerName, state.CurrentRound)

			// 更新 Agent 胜负统计
			if state.Winner == engine.SideRed {
				r.store.UpdateAgentStats(state.AgentRed.ID, true)
				r.store.UpdateAgentStats(state.AgentBlue.ID, false)
			} else if state.Winner == engine.SideBlue {
				r.store.UpdateAgentStats(state.AgentBlue.ID, true)
				r.store.UpdateAgentStats(state.AgentRed.ID, false)
			}

			// 持久化对局历史（DB）
			if r.historyStore != nil {
				winnerSide := string(state.Winner)
				h := user.GameHistory{
					GameID:      gameID,
					RedName:     state.AgentRed.Name,
					BlueName:    state.AgentBlue.Name,
					Winner:      winnerSide,
					WinnerName:  winnerName,
					RedHP:       state.AgentRed.HP,
					BlueHP:      state.AgentBlue.HP,
					TotalRounds: state.CurrentRound,
					CreatedAt:   time.Now(),
				}
				if err := r.historyStore.Record(h); err != nil {
					log.Printf("[GameRunner] record game history FAILED: %v", err)
				} else {
					log.Printf("[GameRunner] game %d recorded in history", gameID)
				}
			}

			// 链上: FinishGame（提交结果 + 结算）— 异步执行，带重试
			if r.chain != nil {
				actionsHash := computeActionsHash(gameID, state.CurrentRound, state.AgentRed.HP, state.AgentBlue.HP)
				gameWinner := state.Winner
				go func() {
					ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
					defer cancel()

					var txHash string
					var finishErr error
					for attempt := 1; attempt <= 3; attempt++ {
						log.Printf("[GameRunner] chain.FinishGame: gameID=%d, winner=%v (attempt %d)", gameID, gameWinner, attempt)
						txHash, finishErr = r.chain.FinishGame(ctx, gameID, redWins, actionsHash)
						if finishErr == nil {
							log.Printf("[GameRunner] chain.FinishGame tx: %s (attempt %d)", txHash, attempt)
							for _, cb := range r.onFinishTx {
								cb(gameID, txHash)
							}
							return
						}
						log.Printf("[GameRunner] chain.FinishGame attempt %d failed: %v", attempt, finishErr)
						if attempt < 3 {
							time.Sleep(5 * time.Second)
						}
					}
					log.Printf("[GameRunner] FATAL: chain.FinishGame failed after 3 attempts: %v", finishErr)
					r.hub.Broadcast(gameID, "chain_error", map[string]interface{}{
						"game_id": gameID,
						"step":    "finish_game",
						"error":   finishErr.Error(),
						"message": "链上结算失败，奖金可能无法领取。请检查 deployer 钱包余额和 owner 权限",
					})
				}()
			}

			r.hub.Broadcast(gameID, "game_finished", GameFinishedData{
				GameID:      gameID,
				Winner:      string(state.Winner),
				WinnerName:  winnerName,
				TotalRounds: state.CurrentRound,
				FinalHPRed:  state.AgentRed.HP,
				FinalHPBlue: state.AgentBlue.HP,
			})

			// 触发游戏结束回调
			winnerSide := string(state.Winner)
			for _, cb := range r.onGameFinish {
				go cb(gameID, winnerSide)
			}
		}
	}
}

// computeActionsHash 计算对局动作哈希（用于链上验证）
func computeActionsHash(gameID uint64, rounds uint8, hpRed, hpBlue uint8) [32]byte {
	buf := make([]byte, 12)
	binary.BigEndian.PutUint64(buf[0:8], gameID)
	buf[8] = rounds
	buf[9] = hpRed
	buf[10] = hpBlue
	buf[11] = 0 // padding
	return sha256.Sum256(buf)
}

// toTurnUpdateData 转换引擎回合记录为 WS 消息
func toTurnUpdateData(record *engine.TurnRecord, state *engine.GameState) TurnUpdateData {
	return TurnUpdateData{
		Round:         record.Round,
		RedAction:     actionToBrief(record.RedAction),
		BlueAction:    actionToBrief(record.BlueAction),
		RedHP:         state.AgentRed.HP,
		BlueHP:        state.AgentBlue.HP,
		RedPosX:       state.AgentRed.Position.X,
		RedPosY:       state.AgentRed.Position.Y,
		BluePosX:      state.AgentBlue.Position.X,
		BluePosY:      state.AgentBlue.Position.Y,
		RedReasoning:  record.RedReasoning,
		BlueReasoning: record.BlueReasoning,
	}
}

// actionToBrief 转换引擎动作到 WS 摘要
func actionToBrief(action engine.Action) ActionBrief {
	brief := ActionBrief{
		Type:       string(action.Type),
		Failed:     action.Failed,
		FailReason: action.FailReason,
	}
	if action.Type == engine.ActionMove {
		brief.Target = &Pos{
			X: action.Target.X,
			Y: action.Target.Y,
		}
	}
	return brief
}
