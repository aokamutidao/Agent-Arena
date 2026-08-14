package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"agent-arena/backend/internal/agent"
	"agent-arena/backend/internal/user"
)

// CustomAgentHandlers 自定义 Agent handlers
type CustomAgentHandlers struct {
	agentStore *user.AgentStore
}

// NewCustomAgentHandlers 创建 Agent handlers
func NewCustomAgentHandlers(agentStore *user.AgentStore) *CustomAgentHandlers {
	return &CustomAgentHandlers{
		agentStore: agentStore,
	}
}

// CreateAgent 创建自定义 Agent
func (h *CustomAgentHandlers) CreateAgent(c *gin.Context) {
	ownerID, _ := c.Get("user_id")
	ownerAddress, _ := c.Get("user_address")

	var req struct {
		Name         string `json:"name"`
		Personality  string `json:"personality"`
		APIEndpoint  string `json:"api_endpoint"`
		APIKey       string `json:"api_key"`
		Model        string `json:"model"`
		ChallengeFee uint64 `json:"challenge_fee"`
		CurrencyType string `json:"currency_type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// 默认值
	if req.CurrencyType == "" {
		req.CurrencyType = "ac"
	}
	if req.Model == "" {
		req.Model = "gpt-3.5-turbo"
	}

	agent, err := h.agentStore.Create(
		ownerID.(string),
		ownerAddress.(string),
		req.Name,
		req.Personality,
		req.APIEndpoint,
		req.APIKey,
		req.Model,
		req.ChallengeFee,
		req.CurrencyType,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, agent)
}

// GetAgent 获取 Agent 详情
func (h *CustomAgentHandlers) GetAgent(c *gin.Context) {
	agentID := c.Param("agentId")

	agent, err := h.agentStore.Get(agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}

	c.JSON(http.StatusOK, agent)
}

// ListAgents 列出所有上架的 Agent
func (h *CustomAgentHandlers) ListAgents(c *gin.Context) {
	agents, err := h.agentStore.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"agents": agents,
		"total":  len(agents),
	})
}

// GetMyAgents 获取当前用户的所有 Agent
func (h *CustomAgentHandlers) GetMyAgents(c *gin.Context) {
	ownerID, _ := c.Get("user_id")

	agents, err := h.agentStore.GetByOwner(ownerID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"agents": agents,
		"total":  len(agents),
	})
}

// UpdateAgent 更新 Agent
func (h *CustomAgentHandlers) UpdateAgent(c *gin.Context) {
	agentID := c.Param("agentId")
	ownerID, _ := c.Get("user_id")

	// 验证所有权
	agent, err := h.agentStore.Get(agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}
	if agent.OwnerID != ownerID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	var req struct {
		Name         string `json:"name"`
		Personality  string `json:"personality"`
		APIEndpoint  string `json:"api_endpoint"`
		APIKey       string `json:"api_key"`
		Model        string `json:"model"`
		ChallengeFee uint64 `json:"challenge_fee"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.agentStore.Update(agentID, req.Name, req.Personality, req.APIEndpoint, req.APIKey, req.Model, req.ChallengeFee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agent, _ = h.agentStore.Get(agentID)
	c.JSON(http.StatusOK, agent)
}

// SetListed 设置 Agent 上架状态
func (h *CustomAgentHandlers) SetListed(c *gin.Context) {
	agentID := c.Param("agentId")
	ownerID, _ := c.Get("user_id")

	// 验证所有权
	agent, err := h.agentStore.Get(agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}
	if agent.OwnerID != ownerID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	var req struct {
		Listed bool `json:"listed"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.agentStore.SetListed(agentID, req.Listed); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	agent, _ = h.agentStore.Get(agentID)
	c.JSON(http.StatusOK, agent)
}

// DeleteAgent 删除 Agent
func (h *CustomAgentHandlers) DeleteAgent(c *gin.Context) {
	agentID := c.Param("agentId")
	ownerID, _ := c.Get("user_id")

	if err := h.agentStore.Delete(agentID, ownerID.(string)); err != nil {
		if err.Error() == "not authorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "agent deleted"})
}

// TestAgentAPI 测试 Agent 的 API 连通性
func (h *CustomAgentHandlers) TestAgentAPI(c *gin.Context) {
	var req struct {
		APIEndpoint string `json:"api_endpoint" binding:"required"`
		APIKey      string `json:"api_key" binding:"required"`
		Model       string `json:"model"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "api_endpoint and api_key are required"})
		return
	}

	// 创建临时客户端
	client := agent.NewExternalAPIClient(req.APIEndpoint, req.APIKey, req.Model)

	// 发送测试请求
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testPrompt := "测试连接。你的位置在 (1,1)，对手在 (8,8)，距离 14 格。请回复一个动作。"
	response, err := client.Chat(ctx, testPrompt)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"response": response,
		"message":  "API 连接成功！",
	})
}
