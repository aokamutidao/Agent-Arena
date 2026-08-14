package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"agent-arena/backend/internal/engine"
	"agent-arena/backend/internal/user"

	"github.com/gin-gonic/gin"
)

// Handlers HTTP 请求处理器
type Handlers struct {
	store       *MemoryStore
	hub         *WSHub
	runner      *GameRunner
	historyStore *user.GameHistoryStore
}

// NewHandlers 创建处理器
func NewHandlers(store *MemoryStore, hub *WSHub) *Handlers {
	return &Handlers{store: store, hub: hub}
}

// SetHistoryStore 注入对局历史存储（重启后从 DB 兜底读取）
func (h *Handlers) SetHistoryStore(s *user.GameHistoryStore) {
	h.historyStore = s
}

// SetGameRunner 设置游戏循环管理器
func (h *Handlers) SetGameRunner(runner *GameRunner) {
	h.runner = runner
}

// ListGames GET /api/games
func (h *Handlers) ListGames(c *gin.Context) {
	status := c.Query("status")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	games, total := h.store.ListGames(status, limit, offset)

	items := make([]GameListItem, 0, len(games))
	for _, g := range games {
		oddsRed, oddsBlue := CalculateOdds(g.TotalBetRed, g.TotalBetBlue)

		item := GameListItem{
			GameID:       g.Game.GameID,
			AgentRed:     h.buildAgentInfo(g.Game.AgentRed),
			AgentBlue:    h.buildAgentInfo(g.Game.AgentBlue),
			Status:       string(g.Game.Status),
			TotalBetRed:  strconv.FormatUint(g.TotalBetRed, 10),
			TotalBetBlue: strconv.FormatUint(g.TotalBetBlue, 10),
			CurrentRound: g.Game.CurrentRound,
			MaxRounds:    g.Game.MaxRounds,
			OddsRed:      oddsRed,
			OddsBlue:     oddsBlue,
		}
		items = append(items, item)
	}

	c.JSON(http.StatusOK, GamesResponse{
		Games: items,
		Total: total,
	})
}

// GetGame GET /api/games/:gameId
func (h *Handlers) GetGame(c *gin.Context) {
	gameID, err := strconv.ParseUint(c.Param("gameId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid game id"})
		return
	}

	record, err := h.store.GetGame(gameID)
	if err != nil {
		// 内存中无记录：尝试从持久化历史兜底（重启后仍能查看已完成对局）
		if h.historyStore != nil {
			if hist, herr := h.historyStore.Get(gameID); herr == nil && hist != nil {
				resp := GameDetailResponse{
					GameID:       hist.GameID,
					Status:       string(engine.StatusFinished),
					CurrentRound: hist.TotalRounds,
					MaxRounds:    engine.MaxRounds,
					Winner:       hist.Winner,
					AgentRed:     AgentInfo{ID: "red", Name: hist.RedName, Personality: "-"},
					AgentBlue:    AgentInfo{ID: "blue", Name: hist.BlueName, Personality: "-"},
					AgentRedState: &AgentStateResponse{HP: hist.RedHP, MaxHP: engine.MaxHP, Status: []string{}},
					AgentBlueState: &AgentStateResponse{HP: hist.BlueHP, MaxHP: engine.MaxHP, Status: []string{}},
					TotalBetRed:  "0",
					TotalBetBlue: "0",
					StrategyRed:  &engine.StrategyWeights{Aggressive: 33, Defensive: 34, Tricky: 33},
					StrategyBlue: &engine.StrategyWeights{Aggressive: 33, Defensive: 34, Tricky: 33},
					History:      []engine.TurnRecord{},
					Archived:     true,
				}
				c.JSON(http.StatusOK, resp)
				return
			}
		}
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "game not found", Code: 40001})
		return
	}

	game := record.Game
	resp := GameDetailResponse{
		GameID:       game.GameID,
		AgentRed:     h.buildAgentInfo(game.AgentRed),
		AgentBlue:    h.buildAgentInfo(game.AgentBlue),
		Status:       string(game.Status),
		CurrentRound: game.CurrentRound,
		MaxRounds:    game.MaxRounds,
		TotalBetRed:  strconv.FormatUint(record.TotalBetRed, 10),
		TotalBetBlue: strconv.FormatUint(record.TotalBetBlue, 10),
		History:      game.History,
		Winner:       string(game.Winner),
		Overtime:     game.Overtime,
		StrategyRed:  &record.StrategyRed,
		StrategyBlue: &record.StrategyBlue,
	}

	// 对局进行中时返回运行时状态
	if game.Status == engine.StatusPlaying || game.Status == engine.StatusFinished {
		resp.AgentRedState = h.buildAgentState(game.AgentRed)
		resp.AgentBlueState = h.buildAgentState(game.AgentBlue)
	}

	c.JSON(http.StatusOK, resp)
}

// EstimateBet POST /api/bets/estimate
func (h *Handlers) EstimateBet(c *gin.Context) {
	var req EstimateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	record, err := h.store.GetGame(req.GameID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "game not found", Code: 40001})
		return
	}

	amount, err := strconv.ParseUint(req.Amount, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid amount", Code: 40003})
		return
	}

	newRed := record.TotalBetRed
	newBlue := record.TotalBetBlue
	if req.Side == "red" {
		newRed += amount
	} else {
		newBlue += amount
	}

	oddsRed, oddsBlue := CalculateOdds(newRed, newBlue)

	// 计算潜在收益
	var potentialReward uint64
	totalPool := newRed + newBlue
	if req.Side == "red" && newRed > 0 {
		// 扣除 5% 协议费后的收益
		winnerShare := uint64(float64(totalPool) * 0.95)
		potentialReward = winnerShare * amount / newRed
	} else if req.Side == "blue" && newBlue > 0 {
		winnerShare := uint64(float64(totalPool) * 0.95)
		potentialReward = winnerShare * amount / newBlue
	}

	c.JSON(http.StatusOK, EstimateResponse{
		CurrentPoolRed:  strconv.FormatUint(record.TotalBetRed, 10),
		CurrentPoolBlue: strconv.FormatUint(record.TotalBetBlue, 10),
		NewPoolRed:      strconv.FormatUint(newRed, 10),
		PotentialReward: strconv.FormatUint(potentialReward, 10),
		NewOddsRed:      oddsRed,
		NewOddsBlue:     oddsBlue,
	})
}

// ListAgents GET /api/agents
func (h *Handlers) ListAgents(c *gin.Context) {
	agents := h.store.ListAgents()
	c.JSON(http.StatusOK, AgentsResponse{Agents: agents})
}

// GetAgent GET /api/agents/:agentId
func (h *Handlers) GetAgent(c *gin.Context) {
	agentID := c.Param("agentId")

	record, err := h.store.GetAgent(agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "agent not found"})
		return
	}

	c.JSON(http.StatusOK, record.Info)
}

// GetUserBets GET /api/users/:address/bets
func (h *Handlers) GetUserBets(c *gin.Context) {
	address := c.Param("address")
	bets := h.store.GetBetsByUser(address)

	if bets == nil {
		bets = []*BetRecord{}
	}

	c.JSON(http.StatusOK, gin.H{
		"address": address,
		"bets":    bets,
	})
}

// --- Helper Functions ---

func (h *Handlers) buildAgentInfo(agent engine.AgentState) AgentInfo {
	record, err := h.store.GetAgent(agent.ID)
	if err != nil {
		return AgentInfo{
			ID:          agent.ID,
			Name:        agent.Name,
			Personality: agent.Personality,
		}
	}
	return record.Info
}

func (h *Handlers) buildAgentState(agent engine.AgentState) *AgentStateResponse {
	var statusList []string
	for _, eff := range agent.Status {
		statusList = append(statusList, string(eff.Type))
	}
	if statusList == nil {
		statusList = []string{}
	}

	return &AgentStateResponse{
		HP:            agent.HP,
		MaxHP:         agent.MaxHP,
		Position:      agent.Position,
		Status:        statusList,
		SkillCooldown: agent.SkillCooldown,
		IsCharging:    agent.IsCharging,
	}
}

// VoteStrategy POST /api/strategy/vote
func (h *Handlers) VoteStrategy(c *gin.Context) {
	var req VoteStrategyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	record, err := h.store.GetGame(req.GameID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "game not found"})
		return
	}

	// 游戏结束前都允许投票（下注在链上随时可以发生，策略同步不应被阻断）
	if record.Game.Status == engine.StatusFinished {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "game is finished, voting locked"})
		return
	}

	weights, err := h.store.UpdateStrategy(req.GameID, req.Side, req.Strategy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// 广播策略更新
	if h.hub != nil {
		h.hub.Broadcast(req.GameID, "strategy_update", StrategyUpdateData{
			Side:       req.Side,
			Aggressive: weights.Aggressive,
			Defensive:  weights.Defensive,
			Tricky:     weights.Tricky,
		})
	}

	c.JSON(http.StatusOK, VoteStrategyResponse{
		Side:       req.Side,
		Aggressive: weights.Aggressive,
		Defensive:  weights.Defensive,
		Tricky:     weights.Tricky,
	})
}

// StartGame POST /api/games/:gameId/start — 手动开始对局
func (h *Handlers) StartGame(c *gin.Context) {
	gameID, err := strconv.ParseUint(c.Param("gameId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid game id"})
		return
	}

	record, err := h.store.GetGame(gameID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "game not found"})
		return
	}

	if record.Game.Status != engine.StatusBetting {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "game is not in betting phase"})
		return
	}

	if h.runner == nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "game runner not initialized"})
		return
	}

	// 启动游戏循环
	h.runner.StartGame(gameID)

	c.JSON(http.StatusOK, gin.H{
		"game_id": gameID,
		"status":  "starting",
	})
}

// PlaceBet POST /api/bets/place — 下注（内存模式，无需钱包）
func (h *Handlers) PlaceBet(c *gin.Context) {
	var req struct {
		GameID   uint64 `json:"game_id" binding:"required"`
		Side     string `json:"side" binding:"required,oneof=red blue"`
		Amount   string `json:"amount" binding:"required"`
		Strategy string `json:"strategy"` // 可选策略偏好
		User     string `json:"user"`     // 可选用户标识
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	record, err := h.store.GetGame(req.GameID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "game not found"})
		return
	}

	if record.Game.Status != engine.StatusBetting {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "betting is closed"})
		return
	}

	amount, err := strconv.ParseUint(req.Amount, 10, 64)
	if err != nil || amount == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid amount"})
		return
	}

	// 更新下注池
	if err := h.store.UpdateBetPool(req.GameID, req.Side, amount); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// 记录用户下注
	user := req.User
	if user == "" {
		user = fmt.Sprintf("anon_%d", record.Game.GameID)
	}
	h.store.CacheBet(user, &BetRecord{
		GameID:   req.GameID,
		Side:     req.Side,
		Amount:   req.Amount,
		Strategy: req.Strategy,
		Status:   "pending",
	})

	// 如果有策略偏好，同时投票
	if req.Strategy != "" {
		h.store.UpdateStrategy(req.GameID, req.Side, req.Strategy)
	}

	// 广播下注更新
	record, _ = h.store.GetGame(req.GameID)
	oddsRed, oddsBlue := CalculateOdds(record.TotalBetRed, record.TotalBetBlue)
	if h.hub != nil {
		h.hub.Broadcast(req.GameID, "betting_update", BettingUpdateData{
			TotalBetRed:  strconv.FormatUint(record.TotalBetRed, 10),
			TotalBetBlue: strconv.FormatUint(record.TotalBetBlue, 10),
			OddsRed:      oddsRed,
			OddsBlue:     oddsBlue,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"game_id":       req.GameID,
		"side":          req.Side,
		"amount":        req.Amount,
		"total_bet_red": strconv.FormatUint(record.TotalBetRed, 10),
		"total_bet_blue": strconv.FormatUint(record.TotalBetBlue, 10),
		"odds_red":      oddsRed,
		"odds_blue":     oddsBlue,
	})
}

// CreateGame POST /api/games/create — 创建新对局
func (h *Handlers) CreateGame(c *gin.Context) {
	// 从已注册的 Agent 中随机选 2 个
	agents := h.store.ListAgents()
	if len(agents) < 2 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "need at least 2 agents"})
		return
	}

	redAgent := agents[0]
	blueAgent := agents[1]
	// 简单 shuffle
	if len(agents) > 2 {
		blueAgent = agents[len(agents)/2]
	}

	// 先同步在链上创建对局，获取链上 game ID
	var gameID uint64
	if h.runner != nil && h.runner.chain != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		txHash, onChainID, err := h.runner.chain.CreateGame(ctx, redAgent.Name, blueAgent.Name, 120)
		if err != nil {
			log.Printf("[CreateGame] chain.CreateGame FAILED: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "链上创建对局失败，请检查部署者钱包余额"})
			return
		}
		log.Printf("[CreateGame] chain.CreateGame tx: %s, onChainID: %d", txHash, onChainID)
		gameID = onChainID
	} else {
		gameID = h.store.NextGameID()
	}

	game := h.runner.CreateNewGame(gameID, redAgent.Name, blueAgent.Name, redAgent.Personality, blueAgent.Personality)

	c.JSON(http.StatusOK, gin.H{
		"game_id":    game.GameID,
		"agent_red":  redAgent.Name,
		"agent_blue": blueAgent.Name,
		"status":     string(game.Status),
	})
}
