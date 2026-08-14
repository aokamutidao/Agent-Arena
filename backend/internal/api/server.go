package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"agent-arena/backend/internal/blockchain"
	"agent-arena/backend/internal/user"
)

// Server HTTP 服务器
type Server struct {
	router          *gin.Engine
	handlers        *Handlers
	authHandlers    *AuthHandlers
	agentHandlers   *CustomAgentHandlers
	challengeHandlers *ChallengeHandlers
	store           *MemoryStore
	userStore       *user.Store
	agentStore      *user.AgentStore
	challengeStore  *user.ChallengeStore
	gameHistory     *user.GameHistoryStore
	wsHub           *WSHub
	gameRunner      *GameRunner
	chain           blockchain.ChainService
	acService       *blockchain.ACService // ArenaCoin 链上服务
	usdcService     *blockchain.USDCService // USDC 只读服务（余额查询）
}

// NewServer 创建 HTTP 服务器
func NewServer() (*Server, error) {
	store := NewMemoryStore()

	// 初始化数据库（所有 store 共享同一个 SQLite 文件）
	dbPath := "arena.db"
	userStore, err := user.NewStore(dbPath)
	if err != nil {
		return nil, fmt.Errorf("init user store: %w", err)
	}

	agentStore, err := user.NewAgentStore(dbPath)
	if err != nil {
		return nil, fmt.Errorf("init agent store: %w", err)
	}

	challengeStore, err := user.NewChallengeStore(dbPath)
	if err != nil {
		return nil, fmt.Errorf("init challenge store: %w", err)
	}

	gameHistory, err := user.NewGameHistoryStore(dbPath)
	if err != nil {
		return nil, fmt.Errorf("init game history store: %w", err)
	}

	// 初始化 ArenaCoin 链上服务（可选，失败时降级为 DB-only）
	var acService *blockchain.ACService
	rpcURL := os.Getenv("RPC_URL")
	deployerKey := os.Getenv("PRIVATE_KEY")     // 合约 deployer（签名链上管理操作）
	ownerKey := os.Getenv("OWNER_PRIVATE_KEY")  // 合约 owner/treasury（mint/转账 AC）
	acAddress := os.Getenv("AC_ADDRESS")        // 可选，默认使用 Sepolia 已部署地址
	if rpcURL != "" && deployerKey != "" {
		acSvc, err := blockchain.NewACService(rpcURL, deployerKey, ownerKey, acAddress)
		if err != nil {
			log.Printf("WARNING: ACService init failed: %v, falling back to DB-only", err)
		} else {
			acService = acSvc
			log.Printf("ACService initialized (token=%s, treasury=%s)",
				acSvc.TokenAddress().Hex(), acSvc.TreasuryAddress().Hex())
			if ownerKey != "" {
				log.Printf("  Treasury (owner) ≠ Deployer — separate accounts")
			} else {
				log.Printf("  Treasury = Deployer (OWNER_PRIVATE_KEY not set)")
			}
		}
	} else {
		log.Println("ACService disabled (no RPC_URL/PRIVATE_KEY), using DB-only AC")
	}

	// 初始化 USDC 只读服务（可选，用于显示用户 USDC 余额）
	var usdcService *blockchain.USDCService
	usdcAddress := os.Getenv("USDC_ADDRESS")
	if rpcURL != "" {
		usdcSvc, err := blockchain.NewUSDCService(rpcURL, usdcAddress)
		if err != nil {
			log.Printf("WARNING: USDCService init failed: %v", err)
		} else {
			usdcService = usdcSvc
			log.Printf("USDCService initialized (token=%s)", usdcSvc.TokenAddress().Hex())
		}
	}

	wsHub := NewWSHub()
	handlers := NewHandlers(store, wsHub)
	authHandlers := NewAuthHandlers(userStore, acService, usdcService)
	authHandlers.SetChallengeStore(challengeStore)
	agentHandlers := NewCustomAgentHandlers(agentStore)
	challengeHandlers := NewChallengeHandlers(challengeStore, userStore, agentStore)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	server := &Server{
		router:            router,
		handlers:          handlers,
		authHandlers:      authHandlers,
		agentHandlers:     agentHandlers,
		challengeHandlers: challengeHandlers,
		store:             store,
		userStore:         userStore,
		agentStore:        agentStore,
		challengeStore:    challengeStore,
		gameHistory:       gameHistory,
		wsHub:             wsHub,
		acService:         acService,
		usdcService:       usdcService,
	}

	// 启动 WebSocket Hub
	go wsHub.Run()

	server.registerRoutes()
	return server, nil
}

// registerRoutes 注册路由
func (s *Server) registerRoutes() {
	api := s.router.Group("/api")
	{
		// 认证（公开）
		api.POST("/auth/login", s.authHandlers.Login)

		// 认证（需要 JWT）
		auth := api.Group("/auth")
		auth.Use(JWTMiddleware())
		{
			auth.GET("/profile", s.authHandlers.GetProfile)
			auth.PUT("/profile", s.authHandlers.UpdateProfile)
			auth.POST("/claim-daily", s.authHandlers.ClaimDailyAC)
			auth.POST("/withdraw", s.authHandlers.WithdrawAC)
			auth.GET("/earnings", s.authHandlers.GetEarnings)

			// 自定义 Agent（需要认证）
			auth.POST("/agents", s.agentHandlers.CreateAgent)
			auth.GET("/agents/my", s.agentHandlers.GetMyAgents)
			auth.PUT("/agents/:agentId", s.agentHandlers.UpdateAgent)
			auth.DELETE("/agents/:agentId", s.agentHandlers.DeleteAgent)
			auth.PUT("/agents/:agentId/listed", s.agentHandlers.SetListed)
			auth.POST("/agents/test-api", s.agentHandlers.TestAgentAPI)

			// PVE 挑战（需要认证）
			auth.POST("/challenges", s.challengeHandlers.CreateChallenge)
			auth.GET("/challenges/my", s.challengeHandlers.GetMyChallenges)
		}

		// 对局
		api.GET("/games", s.handlers.ListGames)
		api.GET("/games/:gameId", s.handlers.GetGame)
		api.POST("/games/:gameId/start", s.handlers.StartGame)
		api.POST("/games/create", s.handlers.CreateGame)

		// 对局历史（持久化）
		api.GET("/game-history", func(c *gin.Context) {
			limitStr := c.DefaultQuery("limit", "20")
			offsetStr := c.DefaultQuery("offset", "0")
			var limit, offset int
			fmt.Sscanf(limitStr, "%d", &limit)
			fmt.Sscanf(offsetStr, "%d", &offset)
			items, err := s.gameHistory.List(limit, offset)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"history": items,
				"total":   len(items),
			})
		})
		api.GET("/game-history/:gameId", func(c *gin.Context) {
			var gameID uint64
			fmt.Sscanf(c.Param("gameId"), "%d", &gameID)
			item, err := s.gameHistory.Get(gameID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
				return
			}
			c.JSON(http.StatusOK, item)
		})

		// 下注
		api.POST("/bets/estimate", s.handlers.EstimateBet)
		api.POST("/bets/place", s.handlers.PlaceBet)

		// 策略投票
		api.POST("/strategy/vote", s.handlers.VoteStrategy)

		// Agent（公开）
		api.GET("/agents", s.handlers.ListAgents)
		api.GET("/agents/:agentId", s.handlers.GetAgent)

		// 自定义 Agent 市场（公开）
		api.GET("/marketplace/agents", s.agentHandlers.ListAgents)
		api.GET("/marketplace/agents/:agentId", s.agentHandlers.GetAgent)

		// 挑战（公开）
		api.GET("/challenges/active", s.challengeHandlers.ListActiveChallenges)
		api.GET("/challenges/:challengeId", s.challengeHandlers.GetChallenge)

		// 用户
		api.GET("/users/:address/bets", s.handlers.GetUserBets)
	}

	// WebSocket
	s.router.GET("/ws", s.WSHandler)

	// 健康检查
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}

// Store 获取数据存储（供外部使用）
func (s *Server) Store() *MemoryStore {
	return s.store
}

// UserStore 获取用户存储
func (s *Server) UserStore() *user.Store {
	return s.userStore
}

// AgentStore 获取自定义 Agent 存储
func (s *Server) AgentStore() *user.AgentStore {
	return s.agentStore
}

// ChallengeStore 获取挑战存储
func (s *Server) ChallengeStore() *user.ChallengeStore {
	return s.challengeStore
}

// GameHistory 获取对局历史存储
func (s *Server) GameHistory() *user.GameHistoryStore {
	return s.gameHistory
}

// WSHub 获取 WebSocket Hub（供外部使用）
func (s *Server) WSHub() *WSHub {
	return s.wsHub
}

// SetGameRunner 设置游戏循环管理器
func (s *Server) SetGameRunner(runner *GameRunner) {
	s.gameRunner = runner
	s.handlers.SetGameRunner(runner)
	s.challengeHandlers.SetGameRunner(runner)
}

// Handlers 获取 HTTP 处理器（用于注入外部依赖）
func (s *Server) Handlers() *Handlers {
	return s.handlers
}

// GameRunner 获取游戏循环管理器
func (s *Server) GameRunner() *GameRunner {
	return s.gameRunner
}

// SetChainService 设置链上服务
func (s *Server) SetChainService(chain blockchain.ChainService) {
	s.chain = chain
}

// ChainService 获取链上服务
func (s *Server) ChainService() blockchain.ChainService {
	return s.chain
}

// SetACService 设置 ArenaCoin 链上服务
func (s *Server) SetACService(ac *blockchain.ACService) {
	s.acService = ac
	if s.authHandlers != nil {
		s.authHandlers.SetACService(ac)
	}
}

// ACService 获取 ArenaCoin 链上服务（可能为 nil）
func (s *Server) ACService() *blockchain.ACService {
	return s.acService
}

// Run 启动服务器
func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

// Router 获取 gin router（用于测试）
func (s *Server) Router() *gin.Engine {
	return s.router
}
