package api

import (
	"context"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"

	"agent-arena/backend/internal/blockchain"
	"agent-arena/backend/internal/user"
)

// AuthHandlers 认证相关 handler
type AuthHandlers struct {
	userStore      *user.Store
	acService      *blockchain.ACService   // ArenaCoin 链上服务（可为 nil）
	usdcService    *blockchain.USDCService // USDC 只读服务（可为 nil）
	challengeStore *user.ChallengeStore    // 挑战记录（用于收益页，可为 nil）
}

// NewAuthHandlers 创建认证 handlers
func NewAuthHandlers(userStore *user.Store, acService *blockchain.ACService, usdcService *blockchain.USDCService) *AuthHandlers {
	return &AuthHandlers{
		userStore:   userStore,
		acService:   acService,
		usdcService: usdcService,
	}
}

// SetACService 设置链上 AC 服务
func (h *AuthHandlers) SetACService(ac *blockchain.ACService) {
	h.acService = ac
}

// SetUSDCService 设置链上 USDC 服务
func (h *AuthHandlers) SetUSDCService(u *blockchain.USDCService) {
	h.usdcService = u
}

// SetChallengeStore 设置挑战记录存储（用于收益页）
func (h *AuthHandlers) SetChallengeStore(cs *user.ChallengeStore) {
	h.challengeStore = cs
}

// Login 登录（钱包签名验证）
func (h *AuthHandlers) Login(c *gin.Context) {
	var req user.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// 验证签名
	valid, err := user.VerifySignature(req.Message, req.Signature, req.Address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "signature verification failed: " + err.Error()})
		return
	}
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
		return
	}

	// 获取或创建用户
	userRecord, err := h.userStore.GetOrCreate(req.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	// 生成 JWT
	token, err := user.GenerateJWT(userRecord.ID, userRecord.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, user.LoginResponse{
		Token: token,
		User:  userRecord,
	})
}

// GetProfile 获取当前用户信息（需要认证）
func (h *AuthHandlers) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userRecord, err := h.userStore.GetByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// AC 完全链上：使用链上余额作为主余额
	var acBalance uint64
	var onChainBalance *uint64
	var tokenAddress string
	var treasuryAddress string
	if h.acService != nil {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()
		userAddr := common.HexToAddress(userRecord.Address)
		if bal, err := h.acService.BalanceOf(ctx, userAddr); err == nil {
			acBalance = bal
			onChainBalance = &bal
		} else {
			log.Printf("[GetProfile] on-chain balance read failed for %s: %v", userRecord.Address, err)
			// 降级：使用 DB 余额
			acBalance = userRecord.ACBalance
		}
		tokenAddress = h.acService.TokenAddress().Hex()
		treasuryAddress = h.acService.TreasuryAddress().Hex()
	} else {
		// 未配置链上 AC，使用 DB 余额
		acBalance = userRecord.ACBalance
	}

	// 如果启用了 USDC，读取链上 USDC 余额
	var usdcBalance *uint64
	var usdcBalanceRaw string
	var usdcTokenAddress string
	if h.usdcService != nil {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()
		userAddr := common.HexToAddress(userRecord.Address)
		if rawBal, err := h.usdcService.BalanceOfRaw(ctx, userAddr); err == nil {
			usdcBalanceRaw = rawBal.String()
			// USDC 6 位小数，转为人类可读的整数
			oneM := new(big.Int).SetUint64(1_000_000)
			intVal := new(big.Int).Div(rawBal, oneM).Uint64()
			usdcBalance = &intVal
		} else {
			log.Printf("[GetProfile] USDC balance read failed for %s: %v", userRecord.Address, err)
		}
		usdcTokenAddress = h.usdcService.TokenAddress().Hex()
	}

	resp := map[string]interface{}{
		"id":          userRecord.ID,
		"address":     userRecord.Address,
		"username":    userRecord.Username,
		"ac_balance":  acBalance, // 链上余额（或 DB 降级）
		"last_claim_at": userRecord.LastClaimAt,
		"created_at":  userRecord.CreatedAt,
		"ac_on_chain": onChainBalance != nil,
	}
	if onChainBalance != nil {
		resp["ac_on_chain_balance"] = *onChainBalance
	}
	if tokenAddress != "" {
		resp["ac_token_address"] = tokenAddress
	}
	if treasuryAddress != "" {
		resp["ac_treasury_address"] = treasuryAddress
	}
	if usdcBalance != nil {
		resp["usdc_balance"] = *usdcBalance
		resp["usdc_on_chain"] = true
	}
	if usdcBalanceRaw != "" {
		resp["usdc_balance_raw"] = usdcBalanceRaw
	}
	if usdcTokenAddress != "" {
		resp["usdc_token_address"] = usdcTokenAddress
	}

	c.JSON(http.StatusOK, resp)
}

// GetEarnings 返回当前用户的收益历史（挑战 + 奖励明细）
func (h *AuthHandlers) GetEarnings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userRecord, err := h.userStore.GetByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if h.challengeStore == nil {
		c.JSON(http.StatusOK, gin.H{"earnings": []interface{}{}, "total_reward_ac": 0})
		return
	}

	challenges, err := h.challengeStore.GetByAddress(userRecord.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed: " + err.Error()})
		return
	}

	// 组装收益条目
	var totalRewardAC uint64
	var totalRewardUSDC uint64
	earnings := make([]map[string]interface{}, 0, len(challenges))
	for _, ch := range challenges {
		entry := map[string]interface{}{
			"challenge_id":  ch.ID,
			"game_id":       ch.GameID,
			"opponent_id":   ch.OpponentID,
			"opponent_type": ch.OpponentType,
			"stake":         ch.Stake,
			"currency":      ch.CurrencyType,
			"winner":        ch.Winner,
			"reward":        ch.Reward,
			"status":        ch.Status,
			"created_at":    ch.CreatedAt,
			"finished_at":   ch.FinishedAt,
		}
		earnings = append(earnings, entry)

		if ch.Status == "finished" {
			switch ch.CurrencyType {
			case "usdc":
				totalRewardUSDC += ch.Reward
			default:
				totalRewardAC += ch.Reward
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"earnings":        earnings,
		"total_reward_ac": totalRewardAC,
		"total_reward_usdc": totalRewardUSDC,
	})
}

// UpdateProfile 更新用户信息（需要认证）
func (h *AuthHandlers) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.userStore.UpdateUsername(userID.(string), req.Username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	userRecord, _ := h.userStore.GetByID(userID.(string))
	c.JSON(http.StatusOK, userRecord)
}

// ClaimDailyAC 领取每日 AC（需要认证，24小时冷却）
// 如果配置了链上 AC：调用合约 owner.mint(user, 100 AC) 真正铸造到用户地址
// 如果未配置：只更新数据库余额（兼容旧模式）
func (h *AuthHandlers) ClaimDailyAC(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userRecord, err := h.userStore.GetByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// 检查是否可以领取（24小时冷却）
	canClaim, err := h.userStore.CanClaimDaily(userRecord.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check claim status"})
		return
	}

	if !canClaim {
		var hoursLeft float64
		if userRecord.LastClaimAt != nil {
			hoursLeft = 24 - time.Since(*userRecord.LastClaimAt).Hours()
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "每日只能领取一次",
			"hours_left": hoursLeft,
			"last_claim": userRecord.LastClaimAt,
		})
		return
	}

	const dailyAmount = uint64(100)
	var txHash string
	var onChain bool

	// 优先链上铸造
	if h.acService != nil {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
		defer cancel()

		userAddr := common.HexToAddress(userRecord.Address)
		hash, err := h.acService.MintToUser(ctx, userAddr, dailyAmount)
		if err != nil {
			log.Printf("[ClaimDailyAC] on-chain mint failed for %s: %v", userRecord.Address, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "on-chain mint failed: " + err.Error()})
			return
		}
		txHash = hash
		onChain = true
		log.Printf("[ClaimDailyAC] on-chain mint success: %s -> %s AC, tx=%s",
			userRecord.Address, dailyAmount, txHash)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "on-chain AC not configured"})
		return
	}

	// 只更新领取时间（不再更新 DB 余额，AC 完全链上）
	now := time.Now()
	if err := h.userStore.UpdateLastClaimAt(userRecord.Address, now); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update claim time"})
		return
	}

	// 读取链上最新余额
	var chainBalance uint64
	if h.acService != nil {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()
		userAddr := common.HexToAddress(userRecord.Address)
		if bal, err := h.acService.BalanceOf(ctx, userAddr); err == nil {
			chainBalance = bal
		}
	}

	userRecord, _ = h.userStore.GetByID(userID.(string))
	resp := gin.H{
		"message":     "claimed 100 AC",
		"ac_balance":  chainBalance, // 返回链上余额
		"next_claim":  now.Add(24 * time.Hour),
		"on_chain":    onChain,
	}
	if txHash != "" {
		resp["tx_hash"] = txHash
	}
	c.JSON(http.StatusOK, resp)
}

// WithdrawAC 已废弃：AC 现在是完全链上的 ERC20，用户可直接在钱包中操作
// 保留此端点但返回提示信息
func (h *AuthHandlers) WithdrawAC(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "AC is fully on-chain now. Your AC tokens are already in your wallet - no withdrawal needed. You can transfer/use them directly from your wallet (e.g. MetaMask).",
	})
}

// JWTMiddleware JWT 认证中间件
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// 解析 Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 验证 JWT
		claims, err := user.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// 将用户信息存入 context
		c.Set("user_id", claims.UserID)
		c.Set("user_address", claims.Address)

		c.Next()
	}
}
