package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"agent-arena/backend/internal/user"
)

// ChallengeHandlers PVE 挑战 handlers
type ChallengeHandlers struct {
	challengeStore *user.ChallengeStore
	userStore      *user.Store
	agentStore     *user.AgentStore
	gameRunner     *GameRunner
}

// NewChallengeHandlers 创建挑战 handlers
func NewChallengeHandlers(challengeStore *user.ChallengeStore, userStore *user.Store, agentStore *user.AgentStore) *ChallengeHandlers {
	return &ChallengeHandlers{
		challengeStore: challengeStore,
		userStore:      userStore,
		agentStore:     agentStore,
	}
}

// SetGameRunner 设置 GameRunner（延迟注入，避免循环依赖）
func (h *ChallengeHandlers) SetGameRunner(runner *GameRunner) {
	h.gameRunner = runner
}

// CreateChallenge 创建挑战（PVE）
func (h *ChallengeHandlers) CreateChallenge(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userAddr, _ := c.Get("user_address")

	var req struct {
		ChallengerAgentID string `json:"challenger_agent_id"` // 用户的自定义 Agent ID
		OpponentID        string `json:"opponent_id"`
		OpponentType      string `json:"opponent_type"` // "system" or "user"
		Stake             uint64 `json:"stake"`
		CurrencyType      string `json:"currency_type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// 验证用户的 Agent
	if req.ChallengerAgentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "challenger_agent_id is required"})
		return
	}

	challengerAgent, err := h.agentStore.Get(req.ChallengerAgentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "challenger agent not found"})
		return
	}

	// 验证 Agent 所有权
	if challengerAgent.OwnerID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your agent"})
		return
	}

	// 确定对手信息
	var opponentName, opponentPersonality string
	if req.OpponentType == "user" {
		agent, err := h.agentStore.Get(req.OpponentID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "opponent agent not found"})
			return
		}
		if !agent.IsListed {
			c.JSON(http.StatusBadRequest, gin.H{"error": "opponent agent not listed"})
			return
		}
		opponentName = agent.Name
		opponentPersonality = agent.Personality
		req.Stake = agent.ChallengeFee
		req.CurrencyType = agent.CurrencyType
	} else if req.OpponentType == "system" {
		// 系统 Agent 预设
		systemAgents := map[string]struct{ name, personality string }{
			"berserker": {"Berserker", "Best defense is offense. Shortest path to enemy, max damage, rarely waits."},
			"tactician": {"Tactician", "Control distance = control the game. Kiting, ranged skills, patient."},
			"trickster": {"Trickster", "Deception is the best weapon. Fakes retreats, surprise charges, uses obstacles."},
			"defender":  {"Defender", "Patience is the key to victory. Sits back, lets opponent overextend, never charges."},
		}
		sa, ok := systemAgents[req.OpponentID]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unknown system agent"})
			return
		}
		opponentName = sa.name
		opponentPersonality = sa.personality
		req.Stake = 10
		req.CurrencyType = "ac"
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid opponent_type"})
		return
	}

	// 检查用户余额（AC）
	if req.CurrencyType == "ac" {
		userRecord, err := h.userStore.GetByID(userID.(string))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if userRecord.ACBalance < req.Stake {
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient AC balance"})
			return
		}
		// 扣除赌注
		h.userStore.AddACBalance(userAddr.(string), ^(req.Stake - 1)) // 减去
	}

	// 创建游戏
	var gameID uint64
	if h.gameRunner != nil {
		// 先在链上创建对局（如果有链上服务）
		if h.gameRunner.chain != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			txHash, onChainID, err := h.gameRunner.chain.CreateGame(ctx, challengerAgent.Name, opponentName, 120)
			if err != nil {
				log.Printf("[CreateChallenge] chain.CreateGame FAILED: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "链上创建对局失败，请检查部署者钱包余额"})
				return
			}
			log.Printf("[CreateChallenge] chain.CreateGame tx: %s, onChainID: %d", txHash, onChainID)
			gameID = onChainID
		} else {
			// 没有链上服务，使用内存 ID
			gameID = h.gameRunner.store.NextGameID()
		}

		challenge, err := h.challengeStore.Create(
			userID.(string),
			userAddr.(string),
			req.OpponentID,
			req.OpponentType,
			gameID,
			req.Stake,
			req.CurrencyType,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 用户的自定义 Agent 作为红方，对手作为蓝方
		game := h.gameRunner.CreateNewGameWithAgent(
			gameID,
			challengerAgent.Name,
			opponentName,
			challengerAgent.Personality,
			opponentPersonality,
			challengerAgent.APIEndpoint,
			challengerAgent.APIKey,
			challengerAgent.Model,
		)

		// 保持 betting 状态：挑战者需要在前端手动点击"开始对局"
		// 观战者可以在 betting 窗口内下注
		_ = game

		c.JSON(http.StatusCreated, gin.H{
			"message":      "challenge created, waiting for challenger to start",
			"challenge_id": challenge.ID,
			"game_id":      gameID,
			"status":       "betting",
			"challenger":   challengerAgent.Name,
			"opponent":     opponentName,
			"stake":        req.Stake,
			"currency":     req.CurrencyType,
		})
		return
	}

	// 如果没有 GameRunner，只创建挑战记录（测试用）
	challenge, err := h.challengeStore.Create(
		userID.(string),
		userAddr.(string),
		req.OpponentID,
		req.OpponentType,
		gameID,
		req.Stake,
		req.CurrencyType,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, challenge)
}

// GetChallenge 获取挑战详情
func (h *ChallengeHandlers) GetChallenge(c *gin.Context) {
	challengeID := c.Param("challengeId")

	challenge, err := h.challengeStore.Get(challengeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "challenge not found"})
		return
	}

	c.JSON(http.StatusOK, challenge)
}

// GetMyChallenges 获取我的挑战历史
func (h *ChallengeHandlers) GetMyChallenges(c *gin.Context) {
	userID, _ := c.Get("user_id")

	challenges, err := h.challengeStore.GetByUser(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"challenges": challenges,
		"total":      len(challenges),
	})
}

// ListActiveChallenges 列出所有进行中的挑战
func (h *ChallengeHandlers) ListActiveChallenges(c *gin.Context) {
	challenges, err := h.challengeStore.ListActive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"challenges": challenges,
		"total":      len(challenges),
	})
}

// FinishChallenge 完成挑战（内部调用）
func (h *ChallengeHandlers) FinishChallenge(challengeID, winner string, reward uint64) error {
	return h.challengeStore.UpdateStatus(challengeID, "finished", winner, reward)
}
