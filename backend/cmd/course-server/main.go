package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"agent-arena/backend/internal/agent"
	"agent-arena/backend/internal/tournament"
)

func number(name string, fallback uint64) (uint64, error) {
	value := os.Getenv(name)
	if value == "" {
		return fallback, nil
	}
	n, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be an unsigned integer", name)
	}
	return n, nil
}
func run() error {
	chainID, err := number("CHAIN_ID", 11155111)
	if err != nil {
		return err
	}
	block, err := number("DEPLOYMENT_BLOCK", 0)
	if err != nil {
		return err
	}
	confirmations, err := number("CONFIRMATIONS", 6)
	if err != nil {
		return err
	}
	poll, err := number("POLL_INTERVAL_SECONDS", 5)
	if err != nil {
		return err
	}
	round, err := number("ROUND_INTERVAL_MS", 250)
	if err != nil {
		return err
	}
	simulation, err := strconv.ParseBool(valueOr("SIMULATION_ENABLED", "true"))
	if err != nil {
		return fmt.Errorf("SIMULATION_ENABLED must be true or false")
	}
	cfg := tournament.Config{DatabasePath: os.Getenv("DATABASE_PATH"), JWTSecret: os.Getenv("JWT_SECRET"), AuthDomain: os.Getenv("AUTH_DOMAIN"), ChainID: chainID, ManagerAddress: os.Getenv("TOURNAMENT_MANAGER_ADDRESS"), DeploymentBlock: block, Confirmations: confirmations, SimulationEnabled: simulation, ABIPath: os.Getenv("CONTRACT_ABI_DIR"), PollInterval: time.Duration(poll) * time.Second, RoundInterval: time.Duration(round) * time.Millisecond}
	llmConfig, err := agent.CourseLLMConfig(os.Getenv)
	if err != nil {
		return err
	}
	maxCalls, err := number("LLM_MAX_CALLS", 120)
	if err != nil || maxCalls < 2 || maxCalls > 120 {
		return fmt.Errorf("LLM_MAX_CALLS must be between 2 and 120 per match")
	}
	cfg.LLMProvider, cfg.LLMModel = llmConfig.Provider, llmConfig.Model
	cfg.LLMTimeoutSeconds = int64(llmConfig.Timeout / time.Second)
	cfg.LLMMaxCalls = int(maxCalls)
	for _, origin := range strings.Split(os.Getenv("CORS_ORIGINS"), ",") {
		if trimmed := strings.TrimSpace(origin); trimmed != "" {
			cfg.CORSOrigins = append(cfg.CORSOrigins, trimmed)
		}
	}
	for _, proxy := range strings.Split(os.Getenv("TRUSTED_PROXY_CIDRS"), ",") {
		if trimmed := strings.TrimSpace(proxy); trimmed != "" {
			cfg.TrustedProxyCIDRs = append(cfg.TrustedProxyCIDRs, trimmed)
		}
	}
	cfg.Defaults()
	if len(cfg.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET is required and must contain at least 32 characters")
	}
	var chain tournament.ChainReader
	rpc := os.Getenv("RPC_URL")
	if rpc != "" || cfg.ManagerAddress != "" {
		if rpc == "" || cfg.ManagerAddress == "" {
			return fmt.Errorf("RPC_URL and TOURNAMENT_MANAGER_ADDRESS must be configured together")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		reader, err := tournament.NewEthReader(ctx, rpc, cfg)
		cancel()
		if err != nil {
			return fmt.Errorf("course deployment validation failed; verify chain, contract addresses and ABI directory")
		}
		defer reader.Close()
		cfg = reader.Config()
		chain = reader
	} else {
		log.Print("Course financial deployment is unconfigured; wallet financial actions remain disabled.")
	}
	var llm agent.LLMClient
	if llmConfig.APIKey != "" {
		llm, err = agent.NewConfiguredQwenClient(llmConfig)
		if err != nil {
			return err
		}
		log.Printf("Explicitly configured LLM mode enabled (provider %s, model %s).",
			llmConfig.Provider, llmConfig.Model)
	} else {
		log.Printf("LLM unconfigured; deterministic simulation enabled=%t.", cfg.SimulationEnabled)
	}
	store, err := tournament.NewStore(cfg.DatabasePath)
	if err != nil {
		return err
	}
	defer store.Close()
	server, err := tournament.NewServer(cfg, store, chain, llm)
	if err != nil {
		return err
	}
	defer server.Close()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	server.Start(ctx)
	port := valueOr("PORT", "8080")
	if _, err := strconv.ParseUint(port, 10, 16); err != nil {
		return fmt.Errorf("PORT must be a valid numeric port")
	}
	httpServer := &http.Server{Addr: ":" + port, Handler: server.Router(), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second}
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = httpServer.Shutdown(shutdown)
	}()
	log.Printf("Course API listening on :%s (chain %d; no transaction signer).", port, cfg.ChainID)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
func valueOr(name, fallback string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return fallback
}
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
