package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"

	"agent-arena/backend/internal/agent"
	"agent-arena/backend/internal/api"
	"agent-arena/backend/internal/blockchain"
	"agent-arena/backend/internal/engine"
)

func main() {
	// 加载 .env
	_ = godotenv.Load()

	server, err := api.NewServer()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}
	store := server.Store()

	// Seed: 注册 4 个 Agent
	seedAgents(store)

	// Seed: 创建测试对局
	eng := engine.NewEngine()
	seedGames(server, store, eng)

	// 创建 GameRunner
	qwenAPIKey := os.Getenv("QWEN_API_KEY")
	if qwenAPIKey == "" {
		log.Println("WARNING: QWEN_API_KEY not set, LLM calls will fail")
	}
	qwen := agent.NewQwenClient(qwenAPIKey)

	runner := api.NewGameRunner(eng, store, server.WSHub(), qwen)
	server.SetGameRunner(runner)

	// 注入对局历史存储：每局结束自动入库
	if gh := server.GameHistory(); gh != nil {
		runner.SetHistoryStore(gh)
		server.Handlers().SetHistoryStore(gh)
	}

	// 链上 FinishGame txHash 回调：更新 game_history
	runner.OnFinishTx(func(gameID uint64, txHash string) {
		if gh := server.GameHistory(); gh != nil {
			if err := gh.UpdateFinishTx(gameID, txHash); err != nil {
				log.Printf("[Main] update finish tx for game %d FAILED: %v", gameID, err)
			} else {
				log.Printf("[Main] game %d finish tx saved: %s", gameID, txHash)
			}
		}
	})

	// 注册游戏结束回调：更新挑战状态并发放奖励（链上 + DB）
	runner.OnGameFinish(func(gameID uint64, winner string) {
		log.Printf("[Main] Game %d finished, winner: %s", gameID, winner)

		// 查找关联的挑战
		challenge, err := server.ChallengeStore().GetByGameID(gameID)
		if err != nil {
			log.Printf("[Main] No challenge found for game %d: %v", gameID, err)
			return
		}

		// AC 发放辅助函数：优先链上，DB 同步更新
		rewardAC := func(toAddr string, amount uint64, reason string) {
			if amount == 0 {
				return
			}
			// 链上：mint 到用户钱包
			if ac := server.ACService(); ac != nil {
				ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
				defer cancel()
				userAddr := common.HexToAddress(toAddr)
				txHash, err := ac.MintToUser(ctx, userAddr, amount)
				if err != nil {
					log.Printf("[Main] on-chain mint FAILED for %s (%s): %v", toAddr, reason, err)
				} else {
					log.Printf("[Main] on-chain mint OK: %s -> %d AC (%s), tx=%s",
						toAddr, amount, reason, txHash)
				}
			}
			// DB：ledger 同步
			if err := server.UserStore().AddACBalance(toAddr, amount); err != nil {
				log.Printf("[Main] DB add AC FAILED for %s: %v", toAddr, err)
			}
		}

		// 判定胜负：红方是挑战者，蓝方是对手
		var challengeWinner string
		var reward uint64
		if winner == "red" {
			challengeWinner = "challenger"
			// 退还赌注 + 奖励
			reward = challenge.Stake + 100
			rewardAC(challenge.ChallengerAddr, reward, "challenge win (stake+bonus)")
			log.Printf("[Main] Challenge %s: challenger %s wins, awarded %d AC",
				challenge.ID, challenge.ChallengerAddr, reward)
		} else if winner == "blue" {
			challengeWinner = "opponent"
			reward = 0
			// 对手赢得赌注
			if challenge.OpponentType == "user" {
				// 查找对手 Agent 的 owner
				agent, err := server.AgentStore().Get(challenge.OpponentID)
				if err == nil && agent.OwnerAddress != "" {
					reward = challenge.Stake
					rewardAC(agent.OwnerAddress, reward, "opponent win (stake)")
					log.Printf("[Main] Challenge %s: opponent owner %s wins %d AC",
						challenge.ID, agent.OwnerAddress, reward)
				} else {
					log.Printf("[Main] Challenge %s: opponent win, but cannot find owner: %v",
						challenge.ID, err)
				}
			} else {
				log.Printf("[Main] Challenge %s: system opponent wins, challenger loses %d AC",
					challenge.ID, challenge.Stake)
			}
		} else {
			challengeWinner = "draw"
			reward = challenge.Stake
			rewardAC(challenge.ChallengerAddr, reward, "draw refund")
			log.Printf("[Main] Challenge %s: draw, refunded %d AC to %s",
				challenge.ID, reward, challenge.ChallengerAddr)
		}

		// 更新挑战状态
		server.ChallengeStore().UpdateStatus(challenge.ID, "finished", challengeWinner, reward)

		// 更新 Agent 统计
		if challenge.OpponentType == "user" {
			opponentWon := (winner == "blue")
			server.AgentStore().UpdateStats(challenge.OpponentID, opponentWon)
		}

		// 通知前端刷新余额（广播到该 gameID 的 WS 房间）
		notifyBalanceChange := func(addr string) {
			user, err := server.UserStore().GetByAddress(addr)
			if err != nil {
				return
			}
			server.WSHub().Broadcast(gameID, "user_balance_update", map[string]interface{}{
				"address":    strings.ToLower(addr),
				"ac_balance": user.ACBalance,
				"reason":     challengeWinner,
			})
		}
		notifyBalanceChange(challenge.ChallengerAddr)
		if challenge.OpponentType == "user" {
			if ag, err := server.AgentStore().Get(challenge.OpponentID); err == nil && ag.OwnerAddress != "" {
				notifyBalanceChange(ag.OwnerAddress)
			}
		}
	})

	// 链上交互: 优先使用 EthChainService，否则回退到 Mock
	chainService := initChainService(server, runner)
	_ = chainService // 用于后续 API handler 接入

	fmt.Println("Agent Arena server starting on :8080")
	fmt.Println("Use POST /api/games/:gameId/start to start a game")
	log.Fatal(server.Run(":8080"))
}

// initChainService 初始化链上服务
func initChainService(server *api.Server, runner *api.GameRunner) blockchain.ChainService {
	rpcURL := os.Getenv("RPC_URL")
	privateKey := os.Getenv("PRIVATE_KEY")
	arenaAddr := os.Getenv("ARENA_ADDRESS")

	if rpcURL != "" && privateKey != "" && arenaAddr != "" {
		log.Printf("Initializing EthChainService (Sepolia)...")
		log.Printf("  RPC: %s", rpcURL)
		log.Printf("  Arena: %s", arenaAddr)

		ethService, err := blockchain.NewEthChainService(rpcURL, privateKey, arenaAddr)
		if err != nil {
			log.Printf("WARNING: EthChainService init failed: %v, falling back to Mock", err)
			mock := blockchain.NewMockChainService()
			server.SetChainService(mock)
			runner.SetChainService(mock)
			return mock
		}

		log.Printf("  Deployer: %s", ethService.TransactorAddress().Hex())
		log.Printf("EthChainService initialized successfully")

		server.SetChainService(ethService)
		runner.SetChainService(ethService)

		return ethService
	}

	log.Println("Using MockChainService (no RPC_URL/PRIVATE_KEY/ARENA_ADDRESS)")
	mock := blockchain.NewMockChainService()
	server.SetChainService(mock)
	runner.SetChainService(mock)
	return mock
}

func seedAgents(store *api.MemoryStore) {
	agents := []api.AgentInfo{
		{
			ID:          "berserker",
			Name:        "Berserker",
			Personality: "狂战士",
			Wins:        0,
			Losses:      0,
			WinRate:     0,
			Description: "最好的防守就是进攻。",
		},
		{
			ID:          "tactician",
			Name:        "Tactician",
			Personality: "战术家",
			Wins:        0,
			Losses:      0,
			WinRate:     0,
			Description: "耐心是胜利的关键。",
		},
		{
			ID:          "trickster",
			Name:        "Trickster",
			Personality: "诡术师",
			Wins:        0,
			Losses:      0,
			WinRate:     0,
			Description: "出其不意，攻其不备。",
		},
		{
			ID:          "defender",
			Name:        "Defender",
			Personality: "守护者",
			Wins:        0,
			Losses:      0,
			WinRate:     0,
			Description: "坚如磐石，不动如山。",
		},
	}

	for _, a := range agents {
		store.RegisterAgent(a)
	}
	fmt.Printf("Seeded %d agents\n", len(agents))
}

func seedGames(server *api.Server, store *api.MemoryStore, eng *engine.Engine) {
	// 跳过 SQLite 历史中已存在的 game ID，避免重启后分配出重复 ID
	// （否则新对局会覆盖旧记录，FinishGame 链上交易也会对不上）
	var maxUsed uint64
	if gh := server.GameHistory(); gh != nil {
		items, err := gh.List(1, 0) // 按 created_at DESC 取 1 条
		if err == nil && len(items) > 0 {
			// 不能只信第一条 — 列表按时间倒序，但 game_id 未必最大
			// 用一条较大 limit 扫描更稳妥
			all, err := gh.List(1000, 0)
			if err == nil {
				for _, it := range all {
					if it.GameID > maxUsed {
						maxUsed = it.GameID
					}
				}
			}
		}
	}
	// 也跳过链上 Arena 合约已有的 ID（1-5）+ 安全缓冲 100
	skip := uint64(100)
	if maxUsed+1 > skip {
		skip = maxUsed + 1
	}
	for i := uint64(1); i <= skip; i++ {
		store.NextGameID()
	}

	fmt.Printf("Seeded 0 games (skipped game_id 1-%d from DB, next new game = %d)\n",
		skip, skip+1)
}
