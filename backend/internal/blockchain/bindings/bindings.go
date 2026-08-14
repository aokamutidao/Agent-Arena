// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// ArenaTypesGameInfo is an auto generated low-level Go binding around an user-defined struct.
type ArenaTypesGameInfo struct {
	GameId          *big.Int
	AgentRed        common.Address
	AgentBlue       common.Address
	TotalBetRed     *big.Int
	TotalBetBlue    *big.Int
	BettingDeadline *big.Int
	Status          uint8
	Winner          uint8
}

// ArenaTypesStrategyVoteRecord is an auto generated low-level Go binding around an user-defined struct.
type ArenaTypesStrategyVoteRecord struct {
	Aggressive *big.Int
	Defensive  *big.Int
	Tricky     *big.Int
	Locked     bool
}

// GameRegistryAgentInfo is an auto generated low-level Go binding around an user-defined struct.
type GameRegistryAgentInfo struct {
	Name        string
	Personality string
	Wins        *big.Int
	Losses      *big.Int
	Exists      bool
}

// GameRegistryGameRecord is an auto generated low-level Go binding around an user-defined struct.
type GameRegistryGameRecord struct {
	GameId         *big.Int
	AgentRed       common.Address
	AgentBlue      common.Address
	StartTimestamp *big.Int
	EndTimestamp   *big.Int
	RedWins        bool
	Exists         bool
}

// AgentArenaBindingsMetaData contains all meta data concerning the AgentArenaBindings contract.
var AgentArenaBindingsMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_usdc\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_protocolTreasury\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"betAndVote\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"side\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"strategy\",\"type\":\"uint8\",\"internalType\":\"enumArenaTypes.Strategy\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"bettingDuration\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"bettingPool\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractBettingPool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createGame\",\"inputs\":[{\"name\":\"agentRed\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agentBlue\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"finishGame\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"redWins\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"actionsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"gameRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractGameRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAgent\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structGameRegistry.AgentInfo\",\"components\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"personality\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"wins\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"losses\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getGame\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structArenaTypes.GameInfo\",\"components\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"agentRed\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agentBlue\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"totalBetRed\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalBetBlue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"bettingDeadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumArenaTypes.GameStatus\"},{\"name\":\"winner\",\"type\":\"uint8\",\"internalType\":\"enumArenaTypes.Side\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOdds\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"oddsRed\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"oddsBlue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReward\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStrategyWeights\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"aggressive\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"defensive\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tricky\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerAgent\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"personality\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBettingDuration\",\"inputs\":[{\"name\":\"_duration\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"startGame\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"strategyVoting\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractStrategyVoting\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"usdc\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"GameFinished\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"redWins\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GameInitialized\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"agentRed\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"agentBlue\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GameStarted\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_usdc\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_protocolTreasury\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"arena\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"bets\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"side\",\"type\":\"uint8\",\"internalType\":\"enumArenaTypes.Side\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"claimed\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"claim\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"games\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"agentRed\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agentBlue\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"totalBetRed\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalBetBlue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"bettingDeadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumArenaTypes.GameStatus\"},{\"name\":\"winner\",\"type\":\"uint8\",\"internalType\":\"enumArenaTypes.Side\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getGame\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structArenaTypes.GameInfo\",\"components\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"agentRed\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agentBlue\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"totalBetRed\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalBetBlue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"bettingDeadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumArenaTypes.GameStatus\"},{\"name\":\"winner\",\"type\":\"uint8\",\"internalType\":\"enumArenaTypes.Side\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOdds\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"oddsRed\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"oddsBlue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReward\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initGame\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"agentRed\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agentBlue\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bettingDeadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"lockBetting\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"placeBet\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"side\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"protocolFeeBps\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"protocolTreasury\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setArena\",\"inputs\":[{\"name\":\"_arena\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"settle\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"redWins\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"usdc\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"winnerPools\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"BetPlaced\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"side\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BettingLocked\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"totalBetRed\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"totalBetBlue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GameSettled\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"redWins\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"},{\"name\":\"totalPool\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"protocolFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardClaimed\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"reward\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"arena\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStrategyWeights\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"aggressiveWeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"defensiveWeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"trickyWeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUserVote\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"strategy\",\"type\":\"uint8\",\"internalType\":\"enumArenaTypes.Strategy\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getVoteRecord\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structArenaTypes.StrategyVoteRecord\",\"components\":[{\"name\":\"aggressive\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"defensive\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tricky\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"locked\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lockVotes\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setArena\",\"inputs\":[{\"name\":\"_arena\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"userVotes\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"strategy\",\"type\":\"uint8\",\"internalType\":\"enumArenaTypes.Strategy\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"vote\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"strategy\",\"type\":\"uint8\",\"internalType\":\"enumArenaTypes.Strategy\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"votes\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"aggressive\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"defensive\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tricky\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"locked\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"StrategyVoted\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"strategy\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumArenaTypes.Strategy\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VotesLocked\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"aggressive\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"defensive\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"tricky\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"actionsHashes\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"agents\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"personality\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"wins\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"losses\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"arena\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createGame\",\"inputs\":[{\"name\":\"agentRed\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agentBlue\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bettingDeadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"gameRecords\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"agentRed\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agentBlue\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"startTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"endTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"redWins\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActionsHash\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAgent\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structGameRegistry.AgentInfo\",\"components\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"personality\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"wins\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"losses\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getGameRecord\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structGameRegistry.GameRecord\",\"components\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"agentRed\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agentBlue\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"startTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"endTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"redWins\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getWinRate\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"wins\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"losses\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"winRate\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nextGameId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerAgent\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"personality\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setArena\",\"inputs\":[{\"name\":\"_arena\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitResult\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"redWins\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"actionsHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AgentRegistered\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"personality\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AgentWinUpdated\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"wins\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"losses\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GameCreated\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"agentRed\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"agentBlue\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"bettingDeadline\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GameResultSubmitted\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"redWins\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"},{\"name\":\"actionsHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// AgentArenaBindingsABI is the input ABI used to generate the binding from.
// Deprecated: Use AgentArenaBindingsMetaData.ABI instead.
var AgentArenaBindingsABI = AgentArenaBindingsMetaData.ABI

// AgentArenaBindings is an auto generated Go binding around an Ethereum contract.
type AgentArenaBindings struct {
	AgentArenaBindingsCaller     // Read-only binding to the contract
	AgentArenaBindingsTransactor // Write-only binding to the contract
	AgentArenaBindingsFilterer   // Log filterer for contract events
}

// AgentArenaBindingsCaller is an auto generated read-only Go binding around an Ethereum contract.
type AgentArenaBindingsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentArenaBindingsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AgentArenaBindingsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentArenaBindingsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AgentArenaBindingsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentArenaBindingsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AgentArenaBindingsSession struct {
	Contract     *AgentArenaBindings // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// AgentArenaBindingsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AgentArenaBindingsCallerSession struct {
	Contract *AgentArenaBindingsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// AgentArenaBindingsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AgentArenaBindingsTransactorSession struct {
	Contract     *AgentArenaBindingsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// AgentArenaBindingsRaw is an auto generated low-level Go binding around an Ethereum contract.
type AgentArenaBindingsRaw struct {
	Contract *AgentArenaBindings // Generic contract binding to access the raw methods on
}

// AgentArenaBindingsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AgentArenaBindingsCallerRaw struct {
	Contract *AgentArenaBindingsCaller // Generic read-only contract binding to access the raw methods on
}

// AgentArenaBindingsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AgentArenaBindingsTransactorRaw struct {
	Contract *AgentArenaBindingsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAgentArenaBindings creates a new instance of AgentArenaBindings, bound to a specific deployed contract.
func NewAgentArenaBindings(address common.Address, backend bind.ContractBackend) (*AgentArenaBindings, error) {
	contract, err := bindAgentArenaBindings(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AgentArenaBindings{AgentArenaBindingsCaller: AgentArenaBindingsCaller{contract: contract}, AgentArenaBindingsTransactor: AgentArenaBindingsTransactor{contract: contract}, AgentArenaBindingsFilterer: AgentArenaBindingsFilterer{contract: contract}}, nil
}

// NewAgentArenaBindingsCaller creates a new read-only instance of AgentArenaBindings, bound to a specific deployed contract.
func NewAgentArenaBindingsCaller(address common.Address, caller bind.ContractCaller) (*AgentArenaBindingsCaller, error) {
	contract, err := bindAgentArenaBindings(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AgentArenaBindingsCaller{contract: contract}, nil
}

// NewAgentArenaBindingsTransactor creates a new write-only instance of AgentArenaBindings, bound to a specific deployed contract.
func NewAgentArenaBindingsTransactor(address common.Address, transactor bind.ContractTransactor) (*AgentArenaBindingsTransactor, error) {
	contract, err := bindAgentArenaBindings(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AgentArenaBindingsTransactor{contract: contract}, nil
}

// NewAgentArenaBindingsFilterer creates a new log filterer instance of AgentArenaBindings, bound to a specific deployed contract.
func NewAgentArenaBindingsFilterer(address common.Address, filterer bind.ContractFilterer) (*AgentArenaBindingsFilterer, error) {
	contract, err := bindAgentArenaBindings(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AgentArenaBindingsFilterer{contract: contract}, nil
}

// bindAgentArenaBindings binds a generic wrapper to an already deployed contract.
func bindAgentArenaBindings(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AgentArenaBindingsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AgentArenaBindings *AgentArenaBindingsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AgentArenaBindings.Contract.AgentArenaBindingsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AgentArenaBindings *AgentArenaBindingsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.AgentArenaBindingsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AgentArenaBindings *AgentArenaBindingsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.AgentArenaBindingsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AgentArenaBindings *AgentArenaBindingsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AgentArenaBindings.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AgentArenaBindings *AgentArenaBindingsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AgentArenaBindings *AgentArenaBindingsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.contract.Transact(opts, method, params...)
}

// ActionsHashes is a free data retrieval call binding the contract method 0x0b9e17bd.
//
// Solidity: function actionsHashes(uint256 ) view returns(bytes32)
func (_AgentArenaBindings *AgentArenaBindingsCaller) ActionsHashes(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "actionsHashes", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ActionsHashes is a free data retrieval call binding the contract method 0x0b9e17bd.
//
// Solidity: function actionsHashes(uint256 ) view returns(bytes32)
func (_AgentArenaBindings *AgentArenaBindingsSession) ActionsHashes(arg0 *big.Int) ([32]byte, error) {
	return _AgentArenaBindings.Contract.ActionsHashes(&_AgentArenaBindings.CallOpts, arg0)
}

// ActionsHashes is a free data retrieval call binding the contract method 0x0b9e17bd.
//
// Solidity: function actionsHashes(uint256 ) view returns(bytes32)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) ActionsHashes(arg0 *big.Int) ([32]byte, error) {
	return _AgentArenaBindings.Contract.ActionsHashes(&_AgentArenaBindings.CallOpts, arg0)
}

// Agents is a free data retrieval call binding the contract method 0xfd66091e.
//
// Solidity: function agents(address ) view returns(string name, string personality, uint256 wins, uint256 losses, bool exists)
func (_AgentArenaBindings *AgentArenaBindingsCaller) Agents(opts *bind.CallOpts, arg0 common.Address) (struct {
	Name        string
	Personality string
	Wins        *big.Int
	Losses      *big.Int
	Exists      bool
}, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "agents", arg0)

	outstruct := new(struct {
		Name        string
		Personality string
		Wins        *big.Int
		Losses      *big.Int
		Exists      bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Name = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.Personality = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Wins = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Losses = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Exists = *abi.ConvertType(out[4], new(bool)).(*bool)

	return *outstruct, err

}

// Agents is a free data retrieval call binding the contract method 0xfd66091e.
//
// Solidity: function agents(address ) view returns(string name, string personality, uint256 wins, uint256 losses, bool exists)
func (_AgentArenaBindings *AgentArenaBindingsSession) Agents(arg0 common.Address) (struct {
	Name        string
	Personality string
	Wins        *big.Int
	Losses      *big.Int
	Exists      bool
}, error) {
	return _AgentArenaBindings.Contract.Agents(&_AgentArenaBindings.CallOpts, arg0)
}

// Agents is a free data retrieval call binding the contract method 0xfd66091e.
//
// Solidity: function agents(address ) view returns(string name, string personality, uint256 wins, uint256 losses, bool exists)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) Agents(arg0 common.Address) (struct {
	Name        string
	Personality string
	Wins        *big.Int
	Losses      *big.Int
	Exists      bool
}, error) {
	return _AgentArenaBindings.Contract.Agents(&_AgentArenaBindings.CallOpts, arg0)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _AgentArenaBindings.Contract.Allowance(&_AgentArenaBindings.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _AgentArenaBindings.Contract.Allowance(&_AgentArenaBindings.CallOpts, owner, spender)
}

// Arena is a free data retrieval call binding the contract method 0xfd3705f9.
//
// Solidity: function arena() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCaller) Arena(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "arena")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Arena is a free data retrieval call binding the contract method 0xfd3705f9.
//
// Solidity: function arena() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsSession) Arena() (common.Address, error) {
	return _AgentArenaBindings.Contract.Arena(&_AgentArenaBindings.CallOpts)
}

// Arena is a free data retrieval call binding the contract method 0xfd3705f9.
//
// Solidity: function arena() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) Arena() (common.Address, error) {
	return _AgentArenaBindings.Contract.Arena(&_AgentArenaBindings.CallOpts)
}

// Arena0 is a free data retrieval call binding the contract method 0xfd3705f9.
//
// Solidity: function arena() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCaller) Arena0(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "arena0")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Arena0 is a free data retrieval call binding the contract method 0xfd3705f9.
//
// Solidity: function arena() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsSession) Arena0() (common.Address, error) {
	return _AgentArenaBindings.Contract.Arena0(&_AgentArenaBindings.CallOpts)
}

// Arena0 is a free data retrieval call binding the contract method 0xfd3705f9.
//
// Solidity: function arena() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) Arena0() (common.Address, error) {
	return _AgentArenaBindings.Contract.Arena0(&_AgentArenaBindings.CallOpts)
}

// Arena1 is a free data retrieval call binding the contract method 0xfd3705f9.
//
// Solidity: function arena() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCaller) Arena1(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "arena1")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Arena1 is a free data retrieval call binding the contract method 0xfd3705f9.
//
// Solidity: function arena() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsSession) Arena1() (common.Address, error) {
	return _AgentArenaBindings.Contract.Arena1(&_AgentArenaBindings.CallOpts)
}

// Arena1 is a free data retrieval call binding the contract method 0xfd3705f9.
//
// Solidity: function arena() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) Arena1() (common.Address, error) {
	return _AgentArenaBindings.Contract.Arena1(&_AgentArenaBindings.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _AgentArenaBindings.Contract.BalanceOf(&_AgentArenaBindings.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _AgentArenaBindings.Contract.BalanceOf(&_AgentArenaBindings.CallOpts, account)
}

// Bets is a free data retrieval call binding the contract method 0xf644b3bb.
//
// Solidity: function bets(uint256 , address ) view returns(uint8 side, uint256 amount, bool claimed)
func (_AgentArenaBindings *AgentArenaBindingsCaller) Bets(opts *bind.CallOpts, arg0 *big.Int, arg1 common.Address) (struct {
	Side    uint8
	Amount  *big.Int
	Claimed bool
}, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "bets", arg0, arg1)

	outstruct := new(struct {
		Side    uint8
		Amount  *big.Int
		Claimed bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Side = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.Amount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Claimed = *abi.ConvertType(out[2], new(bool)).(*bool)

	return *outstruct, err

}

// Bets is a free data retrieval call binding the contract method 0xf644b3bb.
//
// Solidity: function bets(uint256 , address ) view returns(uint8 side, uint256 amount, bool claimed)
func (_AgentArenaBindings *AgentArenaBindingsSession) Bets(arg0 *big.Int, arg1 common.Address) (struct {
	Side    uint8
	Amount  *big.Int
	Claimed bool
}, error) {
	return _AgentArenaBindings.Contract.Bets(&_AgentArenaBindings.CallOpts, arg0, arg1)
}

// Bets is a free data retrieval call binding the contract method 0xf644b3bb.
//
// Solidity: function bets(uint256 , address ) view returns(uint8 side, uint256 amount, bool claimed)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) Bets(arg0 *big.Int, arg1 common.Address) (struct {
	Side    uint8
	Amount  *big.Int
	Claimed bool
}, error) {
	return _AgentArenaBindings.Contract.Bets(&_AgentArenaBindings.CallOpts, arg0, arg1)
}

// BettingDuration is a free data retrieval call binding the contract method 0xcd7636c4.
//
// Solidity: function bettingDuration() view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsCaller) BettingDuration(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "bettingDuration")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BettingDuration is a free data retrieval call binding the contract method 0xcd7636c4.
//
// Solidity: function bettingDuration() view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsSession) BettingDuration() (*big.Int, error) {
	return _AgentArenaBindings.Contract.BettingDuration(&_AgentArenaBindings.CallOpts)
}

// BettingDuration is a free data retrieval call binding the contract method 0xcd7636c4.
//
// Solidity: function bettingDuration() view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) BettingDuration() (*big.Int, error) {
	return _AgentArenaBindings.Contract.BettingDuration(&_AgentArenaBindings.CallOpts)
}

// BettingPool is a free data retrieval call binding the contract method 0xa55e7756.
//
// Solidity: function bettingPool() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCaller) BettingPool(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "bettingPool")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BettingPool is a free data retrieval call binding the contract method 0xa55e7756.
//
// Solidity: function bettingPool() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsSession) BettingPool() (common.Address, error) {
	return _AgentArenaBindings.Contract.BettingPool(&_AgentArenaBindings.CallOpts)
}

// BettingPool is a free data retrieval call binding the contract method 0xa55e7756.
//
// Solidity: function bettingPool() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) BettingPool() (common.Address, error) {
	return _AgentArenaBindings.Contract.BettingPool(&_AgentArenaBindings.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_AgentArenaBindings *AgentArenaBindingsCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_AgentArenaBindings *AgentArenaBindingsSession) Decimals() (uint8, error) {
	return _AgentArenaBindings.Contract.Decimals(&_AgentArenaBindings.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) Decimals() (uint8, error) {
	return _AgentArenaBindings.Contract.Decimals(&_AgentArenaBindings.CallOpts)
}

// GameRecords is a free data retrieval call binding the contract method 0x5882532f.
//
// Solidity: function gameRecords(uint256 ) view returns(uint256 gameId, address agentRed, address agentBlue, uint256 startTimestamp, uint256 endTimestamp, bool redWins, bool exists)
func (_AgentArenaBindings *AgentArenaBindingsCaller) GameRecords(opts *bind.CallOpts, arg0 *big.Int) (struct {
	GameId         *big.Int
	AgentRed       common.Address
	AgentBlue      common.Address
	StartTimestamp *big.Int
	EndTimestamp   *big.Int
	RedWins        bool
	Exists         bool
}, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "gameRecords", arg0)

	outstruct := new(struct {
		GameId         *big.Int
		AgentRed       common.Address
		AgentBlue      common.Address
		StartTimestamp *big.Int
		EndTimestamp   *big.Int
		RedWins        bool
		Exists         bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.GameId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.AgentRed = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.AgentBlue = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.StartTimestamp = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.EndTimestamp = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.RedWins = *abi.ConvertType(out[5], new(bool)).(*bool)
	outstruct.Exists = *abi.ConvertType(out[6], new(bool)).(*bool)

	return *outstruct, err

}

// GameRecords is a free data retrieval call binding the contract method 0x5882532f.
//
// Solidity: function gameRecords(uint256 ) view returns(uint256 gameId, address agentRed, address agentBlue, uint256 startTimestamp, uint256 endTimestamp, bool redWins, bool exists)
func (_AgentArenaBindings *AgentArenaBindingsSession) GameRecords(arg0 *big.Int) (struct {
	GameId         *big.Int
	AgentRed       common.Address
	AgentBlue      common.Address
	StartTimestamp *big.Int
	EndTimestamp   *big.Int
	RedWins        bool
	Exists         bool
}, error) {
	return _AgentArenaBindings.Contract.GameRecords(&_AgentArenaBindings.CallOpts, arg0)
}

// GameRecords is a free data retrieval call binding the contract method 0x5882532f.
//
// Solidity: function gameRecords(uint256 ) view returns(uint256 gameId, address agentRed, address agentBlue, uint256 startTimestamp, uint256 endTimestamp, bool redWins, bool exists)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) GameRecords(arg0 *big.Int) (struct {
	GameId         *big.Int
	AgentRed       common.Address
	AgentBlue      common.Address
	StartTimestamp *big.Int
	EndTimestamp   *big.Int
	RedWins        bool
	Exists         bool
}, error) {
	return _AgentArenaBindings.Contract.GameRecords(&_AgentArenaBindings.CallOpts, arg0)
}

// GameRegistry is a free data retrieval call binding the contract method 0xda090755.
//
// Solidity: function gameRegistry() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCaller) GameRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "gameRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GameRegistry is a free data retrieval call binding the contract method 0xda090755.
//
// Solidity: function gameRegistry() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsSession) GameRegistry() (common.Address, error) {
	return _AgentArenaBindings.Contract.GameRegistry(&_AgentArenaBindings.CallOpts)
}

// GameRegistry is a free data retrieval call binding the contract method 0xda090755.
//
// Solidity: function gameRegistry() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) GameRegistry() (common.Address, error) {
	return _AgentArenaBindings.Contract.GameRegistry(&_AgentArenaBindings.CallOpts)
}

// Games is a free data retrieval call binding the contract method 0x117a5b90.
//
// Solidity: function games(uint256 ) view returns(uint256 gameId, address agentRed, address agentBlue, uint256 totalBetRed, uint256 totalBetBlue, uint256 bettingDeadline, uint8 status, uint8 winner)
func (_AgentArenaBindings *AgentArenaBindingsCaller) Games(opts *bind.CallOpts, arg0 *big.Int) (struct {
	GameId          *big.Int
	AgentRed        common.Address
	AgentBlue       common.Address
	TotalBetRed     *big.Int
	TotalBetBlue    *big.Int
	BettingDeadline *big.Int
	Status          uint8
	Winner          uint8
}, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "games", arg0)

	outstruct := new(struct {
		GameId          *big.Int
		AgentRed        common.Address
		AgentBlue       common.Address
		TotalBetRed     *big.Int
		TotalBetBlue    *big.Int
		BettingDeadline *big.Int
		Status          uint8
		Winner          uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.GameId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.AgentRed = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.AgentBlue = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.TotalBetRed = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.TotalBetBlue = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.BettingDeadline = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.Status = *abi.ConvertType(out[6], new(uint8)).(*uint8)
	outstruct.Winner = *abi.ConvertType(out[7], new(uint8)).(*uint8)

	return *outstruct, err

}

// Games is a free data retrieval call binding the contract method 0x117a5b90.
//
// Solidity: function games(uint256 ) view returns(uint256 gameId, address agentRed, address agentBlue, uint256 totalBetRed, uint256 totalBetBlue, uint256 bettingDeadline, uint8 status, uint8 winner)
func (_AgentArenaBindings *AgentArenaBindingsSession) Games(arg0 *big.Int) (struct {
	GameId          *big.Int
	AgentRed        common.Address
	AgentBlue       common.Address
	TotalBetRed     *big.Int
	TotalBetBlue    *big.Int
	BettingDeadline *big.Int
	Status          uint8
	Winner          uint8
}, error) {
	return _AgentArenaBindings.Contract.Games(&_AgentArenaBindings.CallOpts, arg0)
}

// Games is a free data retrieval call binding the contract method 0x117a5b90.
//
// Solidity: function games(uint256 ) view returns(uint256 gameId, address agentRed, address agentBlue, uint256 totalBetRed, uint256 totalBetBlue, uint256 bettingDeadline, uint8 status, uint8 winner)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) Games(arg0 *big.Int) (struct {
	GameId          *big.Int
	AgentRed        common.Address
	AgentBlue       common.Address
	TotalBetRed     *big.Int
	TotalBetBlue    *big.Int
	BettingDeadline *big.Int
	Status          uint8
	Winner          uint8
}, error) {
	return _AgentArenaBindings.Contract.Games(&_AgentArenaBindings.CallOpts, arg0)
}

// GetActionsHash is a free data retrieval call binding the contract method 0x100e4a59.
//
// Solidity: function getActionsHash(uint256 gameId) view returns(bytes32)
func (_AgentArenaBindings *AgentArenaBindingsCaller) GetActionsHash(opts *bind.CallOpts, gameId *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "getActionsHash", gameId)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetActionsHash is a free data retrieval call binding the contract method 0x100e4a59.
//
// Solidity: function getActionsHash(uint256 gameId) view returns(bytes32)
func (_AgentArenaBindings *AgentArenaBindingsSession) GetActionsHash(gameId *big.Int) ([32]byte, error) {
	return _AgentArenaBindings.Contract.GetActionsHash(&_AgentArenaBindings.CallOpts, gameId)
}

// GetActionsHash is a free data retrieval call binding the contract method 0x100e4a59.
//
// Solidity: function getActionsHash(uint256 gameId) view returns(bytes32)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) GetActionsHash(gameId *big.Int) ([32]byte, error) {
	return _AgentArenaBindings.Contract.GetActionsHash(&_AgentArenaBindings.CallOpts, gameId)
}

// GetAgent is a free data retrieval call binding the contract method 0xfb3551ff.
//
// Solidity: function getAgent(address agentId) view returns((string,string,uint256,uint256,bool))
func (_AgentArenaBindings *AgentArenaBindingsCaller) GetAgent(opts *bind.CallOpts, agentId common.Address) (GameRegistryAgentInfo, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "getAgent", agentId)

	if err != nil {
		return *new(GameRegistryAgentInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(GameRegistryAgentInfo)).(*GameRegistryAgentInfo)

	return out0, err

}

// GetAgent is a free data retrieval call binding the contract method 0xfb3551ff.
//
// Solidity: function getAgent(address agentId) view returns((string,string,uint256,uint256,bool))
func (_AgentArenaBindings *AgentArenaBindingsSession) GetAgent(agentId common.Address) (GameRegistryAgentInfo, error) {
	return _AgentArenaBindings.Contract.GetAgent(&_AgentArenaBindings.CallOpts, agentId)
}

// GetAgent is a free data retrieval call binding the contract method 0xfb3551ff.
//
// Solidity: function getAgent(address agentId) view returns((string,string,uint256,uint256,bool))
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) GetAgent(agentId common.Address) (GameRegistryAgentInfo, error) {
	return _AgentArenaBindings.Contract.GetAgent(&_AgentArenaBindings.CallOpts, agentId)
}

// GetAgent0 is a free data retrieval call binding the contract method 0xfb3551ff.
//
// Solidity: function getAgent(address agentId) view returns((string,string,uint256,uint256,bool))
func (_AgentArenaBindings *AgentArenaBindingsCaller) GetAgent0(opts *bind.CallOpts, agentId common.Address) (GameRegistryAgentInfo, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "getAgent0", agentId)

	if err != nil {
		return *new(GameRegistryAgentInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(GameRegistryAgentInfo)).(*GameRegistryAgentInfo)

	return out0, err

}

// GetAgent0 is a free data retrieval call binding the contract method 0xfb3551ff.
//
// Solidity: function getAgent(address agentId) view returns((string,string,uint256,uint256,bool))
func (_AgentArenaBindings *AgentArenaBindingsSession) GetAgent0(agentId common.Address) (GameRegistryAgentInfo, error) {
	return _AgentArenaBindings.Contract.GetAgent0(&_AgentArenaBindings.CallOpts, agentId)
}

// GetAgent0 is a free data retrieval call binding the contract method 0xfb3551ff.
//
// Solidity: function getAgent(address agentId) view returns((string,string,uint256,uint256,bool))
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) GetAgent0(agentId common.Address) (GameRegistryAgentInfo, error) {
	return _AgentArenaBindings.Contract.GetAgent0(&_AgentArenaBindings.CallOpts, agentId)
}

// GetGame is a free data retrieval call binding the contract method 0xa2f77bcc.
//
// Solidity: function getGame(uint256 gameId) view returns((uint256,address,address,uint256,uint256,uint256,uint8,uint8))
func (_AgentArenaBindings *AgentArenaBindingsCaller) GetGame(opts *bind.CallOpts, gameId *big.Int) (ArenaTypesGameInfo, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "getGame", gameId)

	if err != nil {
		return *new(ArenaTypesGameInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(ArenaTypesGameInfo)).(*ArenaTypesGameInfo)

	return out0, err

}

// GetGame is a free data retrieval call binding the contract method 0xa2f77bcc.
//
// Solidity: function getGame(uint256 gameId) view returns((uint256,address,address,uint256,uint256,uint256,uint8,uint8))
func (_AgentArenaBindings *AgentArenaBindingsSession) GetGame(gameId *big.Int) (ArenaTypesGameInfo, error) {
	return _AgentArenaBindings.Contract.GetGame(&_AgentArenaBindings.CallOpts, gameId)
}

// GetGame is a free data retrieval call binding the contract method 0xa2f77bcc.
//
// Solidity: function getGame(uint256 gameId) view returns((uint256,address,address,uint256,uint256,uint256,uint8,uint8))
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) GetGame(gameId *big.Int) (ArenaTypesGameInfo, error) {
	return _AgentArenaBindings.Contract.GetGame(&_AgentArenaBindings.CallOpts, gameId)
}

// GetGame0 is a free data retrieval call binding the contract method 0xa2f77bcc.
//
// Solidity: function getGame(uint256 gameId) view returns((uint256,address,address,uint256,uint256,uint256,uint8,uint8))
func (_AgentArenaBindings *AgentArenaBindingsCaller) GetGame0(opts *bind.CallOpts, gameId *big.Int) (ArenaTypesGameInfo, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "getGame0", gameId)

	if err != nil {
		return *new(ArenaTypesGameInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(ArenaTypesGameInfo)).(*ArenaTypesGameInfo)

	return out0, err

}

// GetGame0 is a free data retrieval call binding the contract method 0xa2f77bcc.
//
// Solidity: function getGame(uint256 gameId) view returns((uint256,address,address,uint256,uint256,uint256,uint8,uint8))
func (_AgentArenaBindings *AgentArenaBindingsSession) GetGame0(gameId *big.Int) (ArenaTypesGameInfo, error) {
	return _AgentArenaBindings.Contract.GetGame0(&_AgentArenaBindings.CallOpts, gameId)
}

// GetGame0 is a free data retrieval call binding the contract method 0xa2f77bcc.
//
// Solidity: function getGame(uint256 gameId) view returns((uint256,address,address,uint256,uint256,uint256,uint8,uint8))
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) GetGame0(gameId *big.Int) (ArenaTypesGameInfo, error) {
	return _AgentArenaBindings.Contract.GetGame0(&_AgentArenaBindings.CallOpts, gameId)
}

// GetGameRecord is a free data retrieval call binding the contract method 0xe233d0a5.
//
// Solidity: function getGameRecord(uint256 gameId) view returns((uint256,address,address,uint256,uint256,bool,bool))
func (_AgentArenaBindings *AgentArenaBindingsCaller) GetGameRecord(opts *bind.CallOpts, gameId *big.Int) (GameRegistryGameRecord, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "getGameRecord", gameId)

	if err != nil {
		return *new(GameRegistryGameRecord), err
	}

	out0 := *abi.ConvertType(out[0], new(GameRegistryGameRecord)).(*GameRegistryGameRecord)

	return out0, err

}

// GetGameRecord is a free data retrieval call binding the contract method 0xe233d0a5.
//
// Solidity: function getGameRecord(uint256 gameId) view returns((uint256,address,address,uint256,uint256,bool,bool))
func (_AgentArenaBindings *AgentArenaBindingsSession) GetGameRecord(gameId *big.Int) (GameRegistryGameRecord, error) {
	return _AgentArenaBindings.Contract.GetGameRecord(&_AgentArenaBindings.CallOpts, gameId)
}

// GetGameRecord is a free data retrieval call binding the contract method 0xe233d0a5.
//
// Solidity: function getGameRecord(uint256 gameId) view returns((uint256,address,address,uint256,uint256,bool,bool))
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) GetGameRecord(gameId *big.Int) (GameRegistryGameRecord, error) {
	return _AgentArenaBindings.Contract.GetGameRecord(&_AgentArenaBindings.CallOpts, gameId)
}

// GetOdds is a free data retrieval call binding the contract method 0x126c4166.
//
// Solidity: function getOdds(uint256 gameId) view returns(uint256 oddsRed, uint256 oddsBlue)
func (_AgentArenaBindings *AgentArenaBindingsCaller) GetOdds(opts *bind.CallOpts, gameId *big.Int) (struct {
	OddsRed  *big.Int
	OddsBlue *big.Int
}, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "getOdds", gameId)

	outstruct := new(struct {
		OddsRed  *big.Int
		OddsBlue *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.OddsRed = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.OddsBlue = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetOdds is a free data retrieval call binding the contract method 0x126c4166.
//
// Solidity: function getOdds(uint256 gameId) view returns(uint256 oddsRed, uint256 oddsBlue)
func (_AgentArenaBindings *AgentArenaBindingsSession) GetOdds(gameId *big.Int) (struct {
	OddsRed  *big.Int
	OddsBlue *big.Int
}, error) {
	return _AgentArenaBindings.Contract.GetOdds(&_AgentArenaBindings.CallOpts, gameId)
}

// GetOdds is a free data retrieval call binding the contract method 0x126c4166.
//
// Solidity: function getOdds(uint256 gameId) view returns(uint256 oddsRed, uint256 oddsBlue)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) GetOdds(gameId *big.Int) (struct {
	OddsRed  *big.Int
	OddsBlue *big.Int
}, error) {
	return _AgentArenaBindings.Contract.GetOdds(&_AgentArenaBindings.CallOpts, gameId)
}

// GetOdds0 is a free data retrieval call binding the contract method 0x126c4166.
//
// Solidity: function getOdds(uint256 gameId) view returns(uint256 oddsRed, uint256 oddsBlue)
func (_AgentArenaBindings *AgentArenaBindingsCaller) GetOdds0(opts *bind.CallOpts, gameId *big.Int) (struct {
	OddsRed  *big.Int
	OddsBlue *big.Int
}, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "getOdds0", gameId)

	outstruct := new(struct {
		OddsRed  *big.Int
		OddsBlue *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.OddsRed = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.OddsBlue = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetOdds0 is a free data retrieval call binding the contract method 0x126c4166.
//
// Solidity: function getOdds(uint256 gameId) view returns(uint256 oddsRed, uint256 oddsBlue)
func (_AgentArenaBindings *AgentArenaBindingsSession) GetOdds0(gameId *big.Int) (struct {
	OddsRed  *big.Int
	OddsBlue *big.Int
}, error) {
	return _AgentArenaBindings.Contract.GetOdds0(&_AgentArenaBindings.CallOpts, gameId)
}

// GetOdds0 is a free data retrieval call binding the contract method 0x126c4166.
//
// Solidity: function getOdds(uint256 gameId) view returns(uint256 oddsRed, uint256 oddsBlue)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) GetOdds0(gameId *big.Int) (struct {
	OddsRed  *big.Int
	OddsBlue *big.Int
}, error) {
	return _AgentArenaBindings.Contract.GetOdds0(&_AgentArenaBindings.CallOpts, gameId)
}

// GetReward is a free data retrieval call binding the contract method 0x008f33d7.
//
// Solidity: function getReward(uint256 gameId, address user) view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsCaller) GetReward(opts *bind.CallOpts, gameId *big.Int, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "getReward", gameId, user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetReward is a free data retrieval call binding the contract method 0x008f33d7.
//
// Solidity: function getReward(uint256 gameId, address user) view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsSession) GetReward(gameId *big.Int, user common.Address) (*big.Int, error) {
	return _AgentArenaBindings.Contract.GetReward(&_AgentArenaBindings.CallOpts, gameId, user)
}

// GetReward is a free data retrieval call binding the contract method 0x008f33d7.
//
// Solidity: function getReward(uint256 gameId, address user) view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) GetReward(gameId *big.Int, user common.Address) (*big.Int, error) {
	return _AgentArenaBindings.Contract.GetReward(&_AgentArenaBindings.CallOpts, gameId, user)
}

// GetReward0 is a free data retrieval call binding the contract method 0x008f33d7.
//
// Solidity: function getReward(uint256 gameId, address user) view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsCaller) GetReward0(opts *bind.CallOpts, gameId *big.Int, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "getReward0", gameId, user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetReward0 is a free data retrieval call binding the contract method 0x008f33d7.
//
// Solidity: function getReward(uint256 gameId, address user) view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsSession) GetReward0(gameId *big.Int, user common.Address) (*big.Int, error) {
	return _AgentArenaBindings.Contract.GetReward0(&_AgentArenaBindings.CallOpts, gameId, user)
}

// GetReward0 is a free data retrieval call binding the contract method 0x008f33d7.
//
// Solidity: function getReward(uint256 gameId, address user) view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) GetReward0(gameId *big.Int, user common.Address) (*big.Int, error) {
	return _AgentArenaBindings.Contract.GetReward0(&_AgentArenaBindings.CallOpts, gameId, user)
}

// GetStrategyWeights is a free data retrieval call binding the contract method 0x8d511d9d.
//
// Solidity: function getStrategyWeights(uint256 gameId) view returns(uint256 aggressive, uint256 defensive, uint256 tricky)
func (_AgentArenaBindings *AgentArenaBindingsCaller) GetStrategyWeights(opts *bind.CallOpts, gameId *big.Int) (struct {
	Aggressive *big.Int
	Defensive  *big.Int
	Tricky     *big.Int
}, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "getStrategyWeights", gameId)

	outstruct := new(struct {
		Aggressive *big.Int
		Defensive  *big.Int
		Tricky     *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Aggressive = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Defensive = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Tricky = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetStrategyWeights is a free data retrieval call binding the contract method 0x8d511d9d.
//
// Solidity: function getStrategyWeights(uint256 gameId) view returns(uint256 aggressive, uint256 defensive, uint256 tricky)
func (_AgentArenaBindings *AgentArenaBindingsSession) GetStrategyWeights(gameId *big.Int) (struct {
	Aggressive *big.Int
	Defensive  *big.Int
	Tricky     *big.Int
}, error) {
	return _AgentArenaBindings.Contract.GetStrategyWeights(&_AgentArenaBindings.CallOpts, gameId)
}

// GetStrategyWeights is a free data retrieval call binding the contract method 0x8d511d9d.
//
// Solidity: function getStrategyWeights(uint256 gameId) view returns(uint256 aggressive, uint256 defensive, uint256 tricky)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) GetStrategyWeights(gameId *big.Int) (struct {
	Aggressive *big.Int
	Defensive  *big.Int
	Tricky     *big.Int
}, error) {
	return _AgentArenaBindings.Contract.GetStrategyWeights(&_AgentArenaBindings.CallOpts, gameId)
}

// GetStrategyWeights0 is a free data retrieval call binding the contract method 0x8d511d9d.
//
// Solidity: function getStrategyWeights(uint256 gameId) view returns(uint256 aggressiveWeight, uint256 defensiveWeight, uint256 trickyWeight)
func (_AgentArenaBindings *AgentArenaBindingsCaller) GetStrategyWeights0(opts *bind.CallOpts, gameId *big.Int) (struct {
	AggressiveWeight *big.Int
	DefensiveWeight  *big.Int
	TrickyWeight     *big.Int
}, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "getStrategyWeights0", gameId)

	outstruct := new(struct {
		AggressiveWeight *big.Int
		DefensiveWeight  *big.Int
		TrickyWeight     *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.AggressiveWeight = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.DefensiveWeight = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.TrickyWeight = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetStrategyWeights0 is a free data retrieval call binding the contract method 0x8d511d9d.
//
// Solidity: function getStrategyWeights(uint256 gameId) view returns(uint256 aggressiveWeight, uint256 defensiveWeight, uint256 trickyWeight)
func (_AgentArenaBindings *AgentArenaBindingsSession) GetStrategyWeights0(gameId *big.Int) (struct {
	AggressiveWeight *big.Int
	DefensiveWeight  *big.Int
	TrickyWeight     *big.Int
}, error) {
	return _AgentArenaBindings.Contract.GetStrategyWeights0(&_AgentArenaBindings.CallOpts, gameId)
}

// GetStrategyWeights0 is a free data retrieval call binding the contract method 0x8d511d9d.
//
// Solidity: function getStrategyWeights(uint256 gameId) view returns(uint256 aggressiveWeight, uint256 defensiveWeight, uint256 trickyWeight)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) GetStrategyWeights0(gameId *big.Int) (struct {
	AggressiveWeight *big.Int
	DefensiveWeight  *big.Int
	TrickyWeight     *big.Int
}, error) {
	return _AgentArenaBindings.Contract.GetStrategyWeights0(&_AgentArenaBindings.CallOpts, gameId)
}

// GetUserVote is a free data retrieval call binding the contract method 0x03c7881a.
//
// Solidity: function getUserVote(uint256 gameId, address user) view returns(uint8 strategy, uint256 amount)
func (_AgentArenaBindings *AgentArenaBindingsCaller) GetUserVote(opts *bind.CallOpts, gameId *big.Int, user common.Address) (struct {
	Strategy uint8
	Amount   *big.Int
}, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "getUserVote", gameId, user)

	outstruct := new(struct {
		Strategy uint8
		Amount   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Strategy = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.Amount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetUserVote is a free data retrieval call binding the contract method 0x03c7881a.
//
// Solidity: function getUserVote(uint256 gameId, address user) view returns(uint8 strategy, uint256 amount)
func (_AgentArenaBindings *AgentArenaBindingsSession) GetUserVote(gameId *big.Int, user common.Address) (struct {
	Strategy uint8
	Amount   *big.Int
}, error) {
	return _AgentArenaBindings.Contract.GetUserVote(&_AgentArenaBindings.CallOpts, gameId, user)
}

// GetUserVote is a free data retrieval call binding the contract method 0x03c7881a.
//
// Solidity: function getUserVote(uint256 gameId, address user) view returns(uint8 strategy, uint256 amount)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) GetUserVote(gameId *big.Int, user common.Address) (struct {
	Strategy uint8
	Amount   *big.Int
}, error) {
	return _AgentArenaBindings.Contract.GetUserVote(&_AgentArenaBindings.CallOpts, gameId, user)
}

// GetVoteRecord is a free data retrieval call binding the contract method 0x2488d909.
//
// Solidity: function getVoteRecord(uint256 gameId) view returns((uint256,uint256,uint256,bool))
func (_AgentArenaBindings *AgentArenaBindingsCaller) GetVoteRecord(opts *bind.CallOpts, gameId *big.Int) (ArenaTypesStrategyVoteRecord, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "getVoteRecord", gameId)

	if err != nil {
		return *new(ArenaTypesStrategyVoteRecord), err
	}

	out0 := *abi.ConvertType(out[0], new(ArenaTypesStrategyVoteRecord)).(*ArenaTypesStrategyVoteRecord)

	return out0, err

}

// GetVoteRecord is a free data retrieval call binding the contract method 0x2488d909.
//
// Solidity: function getVoteRecord(uint256 gameId) view returns((uint256,uint256,uint256,bool))
func (_AgentArenaBindings *AgentArenaBindingsSession) GetVoteRecord(gameId *big.Int) (ArenaTypesStrategyVoteRecord, error) {
	return _AgentArenaBindings.Contract.GetVoteRecord(&_AgentArenaBindings.CallOpts, gameId)
}

// GetVoteRecord is a free data retrieval call binding the contract method 0x2488d909.
//
// Solidity: function getVoteRecord(uint256 gameId) view returns((uint256,uint256,uint256,bool))
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) GetVoteRecord(gameId *big.Int) (ArenaTypesStrategyVoteRecord, error) {
	return _AgentArenaBindings.Contract.GetVoteRecord(&_AgentArenaBindings.CallOpts, gameId)
}

// GetWinRate is a free data retrieval call binding the contract method 0x914a1e76.
//
// Solidity: function getWinRate(address agentId) view returns(uint256 wins, uint256 losses, uint256 winRate)
func (_AgentArenaBindings *AgentArenaBindingsCaller) GetWinRate(opts *bind.CallOpts, agentId common.Address) (struct {
	Wins    *big.Int
	Losses  *big.Int
	WinRate *big.Int
}, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "getWinRate", agentId)

	outstruct := new(struct {
		Wins    *big.Int
		Losses  *big.Int
		WinRate *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Wins = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Losses = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.WinRate = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetWinRate is a free data retrieval call binding the contract method 0x914a1e76.
//
// Solidity: function getWinRate(address agentId) view returns(uint256 wins, uint256 losses, uint256 winRate)
func (_AgentArenaBindings *AgentArenaBindingsSession) GetWinRate(agentId common.Address) (struct {
	Wins    *big.Int
	Losses  *big.Int
	WinRate *big.Int
}, error) {
	return _AgentArenaBindings.Contract.GetWinRate(&_AgentArenaBindings.CallOpts, agentId)
}

// GetWinRate is a free data retrieval call binding the contract method 0x914a1e76.
//
// Solidity: function getWinRate(address agentId) view returns(uint256 wins, uint256 losses, uint256 winRate)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) GetWinRate(agentId common.Address) (struct {
	Wins    *big.Int
	Losses  *big.Int
	WinRate *big.Int
}, error) {
	return _AgentArenaBindings.Contract.GetWinRate(&_AgentArenaBindings.CallOpts, agentId)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AgentArenaBindings *AgentArenaBindingsCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AgentArenaBindings *AgentArenaBindingsSession) Name() (string, error) {
	return _AgentArenaBindings.Contract.Name(&_AgentArenaBindings.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) Name() (string, error) {
	return _AgentArenaBindings.Contract.Name(&_AgentArenaBindings.CallOpts)
}

// NextGameId is a free data retrieval call binding the contract method 0xb135bbb0.
//
// Solidity: function nextGameId() view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsCaller) NextGameId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "nextGameId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextGameId is a free data retrieval call binding the contract method 0xb135bbb0.
//
// Solidity: function nextGameId() view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsSession) NextGameId() (*big.Int, error) {
	return _AgentArenaBindings.Contract.NextGameId(&_AgentArenaBindings.CallOpts)
}

// NextGameId is a free data retrieval call binding the contract method 0xb135bbb0.
//
// Solidity: function nextGameId() view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) NextGameId() (*big.Int, error) {
	return _AgentArenaBindings.Contract.NextGameId(&_AgentArenaBindings.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsSession) Owner() (common.Address, error) {
	return _AgentArenaBindings.Contract.Owner(&_AgentArenaBindings.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) Owner() (common.Address, error) {
	return _AgentArenaBindings.Contract.Owner(&_AgentArenaBindings.CallOpts)
}

// Owner0 is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCaller) Owner0(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "owner0")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner0 is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsSession) Owner0() (common.Address, error) {
	return _AgentArenaBindings.Contract.Owner0(&_AgentArenaBindings.CallOpts)
}

// Owner0 is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) Owner0() (common.Address, error) {
	return _AgentArenaBindings.Contract.Owner0(&_AgentArenaBindings.CallOpts)
}

// Owner1 is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCaller) Owner1(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "owner1")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner1 is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsSession) Owner1() (common.Address, error) {
	return _AgentArenaBindings.Contract.Owner1(&_AgentArenaBindings.CallOpts)
}

// Owner1 is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) Owner1() (common.Address, error) {
	return _AgentArenaBindings.Contract.Owner1(&_AgentArenaBindings.CallOpts)
}

// Owner2 is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCaller) Owner2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "owner2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner2 is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsSession) Owner2() (common.Address, error) {
	return _AgentArenaBindings.Contract.Owner2(&_AgentArenaBindings.CallOpts)
}

// Owner2 is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) Owner2() (common.Address, error) {
	return _AgentArenaBindings.Contract.Owner2(&_AgentArenaBindings.CallOpts)
}

// Owner3 is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCaller) Owner3(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "owner3")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner3 is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsSession) Owner3() (common.Address, error) {
	return _AgentArenaBindings.Contract.Owner3(&_AgentArenaBindings.CallOpts)
}

// Owner3 is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) Owner3() (common.Address, error) {
	return _AgentArenaBindings.Contract.Owner3(&_AgentArenaBindings.CallOpts)
}

// ProtocolFeeBps is a free data retrieval call binding the contract method 0x35659fb8.
//
// Solidity: function protocolFeeBps() view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsCaller) ProtocolFeeBps(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "protocolFeeBps")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ProtocolFeeBps is a free data retrieval call binding the contract method 0x35659fb8.
//
// Solidity: function protocolFeeBps() view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsSession) ProtocolFeeBps() (*big.Int, error) {
	return _AgentArenaBindings.Contract.ProtocolFeeBps(&_AgentArenaBindings.CallOpts)
}

// ProtocolFeeBps is a free data retrieval call binding the contract method 0x35659fb8.
//
// Solidity: function protocolFeeBps() view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) ProtocolFeeBps() (*big.Int, error) {
	return _AgentArenaBindings.Contract.ProtocolFeeBps(&_AgentArenaBindings.CallOpts)
}

// ProtocolTreasury is a free data retrieval call binding the contract method 0x803db96d.
//
// Solidity: function protocolTreasury() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCaller) ProtocolTreasury(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "protocolTreasury")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ProtocolTreasury is a free data retrieval call binding the contract method 0x803db96d.
//
// Solidity: function protocolTreasury() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsSession) ProtocolTreasury() (common.Address, error) {
	return _AgentArenaBindings.Contract.ProtocolTreasury(&_AgentArenaBindings.CallOpts)
}

// ProtocolTreasury is a free data retrieval call binding the contract method 0x803db96d.
//
// Solidity: function protocolTreasury() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) ProtocolTreasury() (common.Address, error) {
	return _AgentArenaBindings.Contract.ProtocolTreasury(&_AgentArenaBindings.CallOpts)
}

// StrategyVoting is a free data retrieval call binding the contract method 0xe9c45ade.
//
// Solidity: function strategyVoting() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCaller) StrategyVoting(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "strategyVoting")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StrategyVoting is a free data retrieval call binding the contract method 0xe9c45ade.
//
// Solidity: function strategyVoting() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsSession) StrategyVoting() (common.Address, error) {
	return _AgentArenaBindings.Contract.StrategyVoting(&_AgentArenaBindings.CallOpts)
}

// StrategyVoting is a free data retrieval call binding the contract method 0xe9c45ade.
//
// Solidity: function strategyVoting() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) StrategyVoting() (common.Address, error) {
	return _AgentArenaBindings.Contract.StrategyVoting(&_AgentArenaBindings.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_AgentArenaBindings *AgentArenaBindingsCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_AgentArenaBindings *AgentArenaBindingsSession) Symbol() (string, error) {
	return _AgentArenaBindings.Contract.Symbol(&_AgentArenaBindings.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) Symbol() (string, error) {
	return _AgentArenaBindings.Contract.Symbol(&_AgentArenaBindings.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsSession) TotalSupply() (*big.Int, error) {
	return _AgentArenaBindings.Contract.TotalSupply(&_AgentArenaBindings.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) TotalSupply() (*big.Int, error) {
	return _AgentArenaBindings.Contract.TotalSupply(&_AgentArenaBindings.CallOpts)
}

// Usdc is a free data retrieval call binding the contract method 0x3e413bee.
//
// Solidity: function usdc() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCaller) Usdc(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "usdc")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Usdc is a free data retrieval call binding the contract method 0x3e413bee.
//
// Solidity: function usdc() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsSession) Usdc() (common.Address, error) {
	return _AgentArenaBindings.Contract.Usdc(&_AgentArenaBindings.CallOpts)
}

// Usdc is a free data retrieval call binding the contract method 0x3e413bee.
//
// Solidity: function usdc() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) Usdc() (common.Address, error) {
	return _AgentArenaBindings.Contract.Usdc(&_AgentArenaBindings.CallOpts)
}

// Usdc0 is a free data retrieval call binding the contract method 0x3e413bee.
//
// Solidity: function usdc() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCaller) Usdc0(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "usdc0")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Usdc0 is a free data retrieval call binding the contract method 0x3e413bee.
//
// Solidity: function usdc() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsSession) Usdc0() (common.Address, error) {
	return _AgentArenaBindings.Contract.Usdc0(&_AgentArenaBindings.CallOpts)
}

// Usdc0 is a free data retrieval call binding the contract method 0x3e413bee.
//
// Solidity: function usdc() view returns(address)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) Usdc0() (common.Address, error) {
	return _AgentArenaBindings.Contract.Usdc0(&_AgentArenaBindings.CallOpts)
}

// UserVotes is a free data retrieval call binding the contract method 0xfe5b3e3b.
//
// Solidity: function userVotes(uint256 , address ) view returns(uint8 strategy, uint256 amount)
func (_AgentArenaBindings *AgentArenaBindingsCaller) UserVotes(opts *bind.CallOpts, arg0 *big.Int, arg1 common.Address) (struct {
	Strategy uint8
	Amount   *big.Int
}, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "userVotes", arg0, arg1)

	outstruct := new(struct {
		Strategy uint8
		Amount   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Strategy = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.Amount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// UserVotes is a free data retrieval call binding the contract method 0xfe5b3e3b.
//
// Solidity: function userVotes(uint256 , address ) view returns(uint8 strategy, uint256 amount)
func (_AgentArenaBindings *AgentArenaBindingsSession) UserVotes(arg0 *big.Int, arg1 common.Address) (struct {
	Strategy uint8
	Amount   *big.Int
}, error) {
	return _AgentArenaBindings.Contract.UserVotes(&_AgentArenaBindings.CallOpts, arg0, arg1)
}

// UserVotes is a free data retrieval call binding the contract method 0xfe5b3e3b.
//
// Solidity: function userVotes(uint256 , address ) view returns(uint8 strategy, uint256 amount)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) UserVotes(arg0 *big.Int, arg1 common.Address) (struct {
	Strategy uint8
	Amount   *big.Int
}, error) {
	return _AgentArenaBindings.Contract.UserVotes(&_AgentArenaBindings.CallOpts, arg0, arg1)
}

// Votes is a free data retrieval call binding the contract method 0x5df81330.
//
// Solidity: function votes(uint256 ) view returns(uint256 aggressive, uint256 defensive, uint256 tricky, bool locked)
func (_AgentArenaBindings *AgentArenaBindingsCaller) Votes(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Aggressive *big.Int
	Defensive  *big.Int
	Tricky     *big.Int
	Locked     bool
}, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "votes", arg0)

	outstruct := new(struct {
		Aggressive *big.Int
		Defensive  *big.Int
		Tricky     *big.Int
		Locked     bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Aggressive = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Defensive = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Tricky = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Locked = *abi.ConvertType(out[3], new(bool)).(*bool)

	return *outstruct, err

}

// Votes is a free data retrieval call binding the contract method 0x5df81330.
//
// Solidity: function votes(uint256 ) view returns(uint256 aggressive, uint256 defensive, uint256 tricky, bool locked)
func (_AgentArenaBindings *AgentArenaBindingsSession) Votes(arg0 *big.Int) (struct {
	Aggressive *big.Int
	Defensive  *big.Int
	Tricky     *big.Int
	Locked     bool
}, error) {
	return _AgentArenaBindings.Contract.Votes(&_AgentArenaBindings.CallOpts, arg0)
}

// Votes is a free data retrieval call binding the contract method 0x5df81330.
//
// Solidity: function votes(uint256 ) view returns(uint256 aggressive, uint256 defensive, uint256 tricky, bool locked)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) Votes(arg0 *big.Int) (struct {
	Aggressive *big.Int
	Defensive  *big.Int
	Tricky     *big.Int
	Locked     bool
}, error) {
	return _AgentArenaBindings.Contract.Votes(&_AgentArenaBindings.CallOpts, arg0)
}

// WinnerPools is a free data retrieval call binding the contract method 0x80b2a011.
//
// Solidity: function winnerPools(uint256 ) view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsCaller) WinnerPools(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AgentArenaBindings.contract.Call(opts, &out, "winnerPools", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WinnerPools is a free data retrieval call binding the contract method 0x80b2a011.
//
// Solidity: function winnerPools(uint256 ) view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsSession) WinnerPools(arg0 *big.Int) (*big.Int, error) {
	return _AgentArenaBindings.Contract.WinnerPools(&_AgentArenaBindings.CallOpts, arg0)
}

// WinnerPools is a free data retrieval call binding the contract method 0x80b2a011.
//
// Solidity: function winnerPools(uint256 ) view returns(uint256)
func (_AgentArenaBindings *AgentArenaBindingsCallerSession) WinnerPools(arg0 *big.Int) (*big.Int, error) {
	return _AgentArenaBindings.Contract.WinnerPools(&_AgentArenaBindings.CallOpts, arg0)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_AgentArenaBindings *AgentArenaBindingsTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_AgentArenaBindings *AgentArenaBindingsSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.Approve(&_AgentArenaBindings.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.Approve(&_AgentArenaBindings.TransactOpts, spender, value)
}

// BetAndVote is a paid mutator transaction binding the contract method 0x74651cbf.
//
// Solidity: function betAndVote(uint256 gameId, bool side, uint256 amount, uint8 strategy) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactor) BetAndVote(opts *bind.TransactOpts, gameId *big.Int, side bool, amount *big.Int, strategy uint8) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "betAndVote", gameId, side, amount, strategy)
}

// BetAndVote is a paid mutator transaction binding the contract method 0x74651cbf.
//
// Solidity: function betAndVote(uint256 gameId, bool side, uint256 amount, uint8 strategy) returns()
func (_AgentArenaBindings *AgentArenaBindingsSession) BetAndVote(gameId *big.Int, side bool, amount *big.Int, strategy uint8) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.BetAndVote(&_AgentArenaBindings.TransactOpts, gameId, side, amount, strategy)
}

// BetAndVote is a paid mutator transaction binding the contract method 0x74651cbf.
//
// Solidity: function betAndVote(uint256 gameId, bool side, uint256 amount, uint8 strategy) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) BetAndVote(gameId *big.Int, side bool, amount *big.Int, strategy uint8) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.BetAndVote(&_AgentArenaBindings.TransactOpts, gameId, side, amount, strategy)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address from, uint256 amount) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactor) Burn(opts *bind.TransactOpts, from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "burn", from, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address from, uint256 amount) returns()
func (_AgentArenaBindings *AgentArenaBindingsSession) Burn(from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.Burn(&_AgentArenaBindings.TransactOpts, from, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address from, uint256 amount) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) Burn(from common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.Burn(&_AgentArenaBindings.TransactOpts, from, amount)
}

// Claim is a paid mutator transaction binding the contract method 0x379607f5.
//
// Solidity: function claim(uint256 gameId) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactor) Claim(opts *bind.TransactOpts, gameId *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "claim", gameId)
}

// Claim is a paid mutator transaction binding the contract method 0x379607f5.
//
// Solidity: function claim(uint256 gameId) returns()
func (_AgentArenaBindings *AgentArenaBindingsSession) Claim(gameId *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.Claim(&_AgentArenaBindings.TransactOpts, gameId)
}

// Claim is a paid mutator transaction binding the contract method 0x379607f5.
//
// Solidity: function claim(uint256 gameId) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) Claim(gameId *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.Claim(&_AgentArenaBindings.TransactOpts, gameId)
}

// CreateGame is a paid mutator transaction binding the contract method 0xa6f979ff.
//
// Solidity: function createGame(address agentRed, address agentBlue) returns(uint256 gameId)
func (_AgentArenaBindings *AgentArenaBindingsTransactor) CreateGame(opts *bind.TransactOpts, agentRed common.Address, agentBlue common.Address) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "createGame", agentRed, agentBlue)
}

// CreateGame is a paid mutator transaction binding the contract method 0xa6f979ff.
//
// Solidity: function createGame(address agentRed, address agentBlue) returns(uint256 gameId)
func (_AgentArenaBindings *AgentArenaBindingsSession) CreateGame(agentRed common.Address, agentBlue common.Address) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.CreateGame(&_AgentArenaBindings.TransactOpts, agentRed, agentBlue)
}

// CreateGame is a paid mutator transaction binding the contract method 0xa6f979ff.
//
// Solidity: function createGame(address agentRed, address agentBlue) returns(uint256 gameId)
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) CreateGame(agentRed common.Address, agentBlue common.Address) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.CreateGame(&_AgentArenaBindings.TransactOpts, agentRed, agentBlue)
}

// CreateGame0 is a paid mutator transaction binding the contract method 0xb718d70c.
//
// Solidity: function createGame(address agentRed, address agentBlue, uint256 bettingDeadline) returns(uint256 gameId)
func (_AgentArenaBindings *AgentArenaBindingsTransactor) CreateGame0(opts *bind.TransactOpts, agentRed common.Address, agentBlue common.Address, bettingDeadline *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "createGame0", agentRed, agentBlue, bettingDeadline)
}

// CreateGame0 is a paid mutator transaction binding the contract method 0xb718d70c.
//
// Solidity: function createGame(address agentRed, address agentBlue, uint256 bettingDeadline) returns(uint256 gameId)
func (_AgentArenaBindings *AgentArenaBindingsSession) CreateGame0(agentRed common.Address, agentBlue common.Address, bettingDeadline *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.CreateGame0(&_AgentArenaBindings.TransactOpts, agentRed, agentBlue, bettingDeadline)
}

// CreateGame0 is a paid mutator transaction binding the contract method 0xb718d70c.
//
// Solidity: function createGame(address agentRed, address agentBlue, uint256 bettingDeadline) returns(uint256 gameId)
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) CreateGame0(agentRed common.Address, agentBlue common.Address, bettingDeadline *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.CreateGame0(&_AgentArenaBindings.TransactOpts, agentRed, agentBlue, bettingDeadline)
}

// FinishGame is a paid mutator transaction binding the contract method 0x18df11c9.
//
// Solidity: function finishGame(uint256 gameId, bool redWins, bytes32 actionsHash) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactor) FinishGame(opts *bind.TransactOpts, gameId *big.Int, redWins bool, actionsHash [32]byte) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "finishGame", gameId, redWins, actionsHash)
}

// FinishGame is a paid mutator transaction binding the contract method 0x18df11c9.
//
// Solidity: function finishGame(uint256 gameId, bool redWins, bytes32 actionsHash) returns()
func (_AgentArenaBindings *AgentArenaBindingsSession) FinishGame(gameId *big.Int, redWins bool, actionsHash [32]byte) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.FinishGame(&_AgentArenaBindings.TransactOpts, gameId, redWins, actionsHash)
}

// FinishGame is a paid mutator transaction binding the contract method 0x18df11c9.
//
// Solidity: function finishGame(uint256 gameId, bool redWins, bytes32 actionsHash) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) FinishGame(gameId *big.Int, redWins bool, actionsHash [32]byte) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.FinishGame(&_AgentArenaBindings.TransactOpts, gameId, redWins, actionsHash)
}

// InitGame is a paid mutator transaction binding the contract method 0x9aeb7753.
//
// Solidity: function initGame(uint256 gameId, address agentRed, address agentBlue, uint256 bettingDeadline) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactor) InitGame(opts *bind.TransactOpts, gameId *big.Int, agentRed common.Address, agentBlue common.Address, bettingDeadline *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "initGame", gameId, agentRed, agentBlue, bettingDeadline)
}

// InitGame is a paid mutator transaction binding the contract method 0x9aeb7753.
//
// Solidity: function initGame(uint256 gameId, address agentRed, address agentBlue, uint256 bettingDeadline) returns()
func (_AgentArenaBindings *AgentArenaBindingsSession) InitGame(gameId *big.Int, agentRed common.Address, agentBlue common.Address, bettingDeadline *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.InitGame(&_AgentArenaBindings.TransactOpts, gameId, agentRed, agentBlue, bettingDeadline)
}

// InitGame is a paid mutator transaction binding the contract method 0x9aeb7753.
//
// Solidity: function initGame(uint256 gameId, address agentRed, address agentBlue, uint256 bettingDeadline) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) InitGame(gameId *big.Int, agentRed common.Address, agentBlue common.Address, bettingDeadline *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.InitGame(&_AgentArenaBindings.TransactOpts, gameId, agentRed, agentBlue, bettingDeadline)
}

// LockBetting is a paid mutator transaction binding the contract method 0xca28cf19.
//
// Solidity: function lockBetting(uint256 gameId) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactor) LockBetting(opts *bind.TransactOpts, gameId *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "lockBetting", gameId)
}

// LockBetting is a paid mutator transaction binding the contract method 0xca28cf19.
//
// Solidity: function lockBetting(uint256 gameId) returns()
func (_AgentArenaBindings *AgentArenaBindingsSession) LockBetting(gameId *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.LockBetting(&_AgentArenaBindings.TransactOpts, gameId)
}

// LockBetting is a paid mutator transaction binding the contract method 0xca28cf19.
//
// Solidity: function lockBetting(uint256 gameId) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) LockBetting(gameId *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.LockBetting(&_AgentArenaBindings.TransactOpts, gameId)
}

// LockVotes is a paid mutator transaction binding the contract method 0x81cc4a8a.
//
// Solidity: function lockVotes(uint256 gameId) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactor) LockVotes(opts *bind.TransactOpts, gameId *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "lockVotes", gameId)
}

// LockVotes is a paid mutator transaction binding the contract method 0x81cc4a8a.
//
// Solidity: function lockVotes(uint256 gameId) returns()
func (_AgentArenaBindings *AgentArenaBindingsSession) LockVotes(gameId *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.LockVotes(&_AgentArenaBindings.TransactOpts, gameId)
}

// LockVotes is a paid mutator transaction binding the contract method 0x81cc4a8a.
//
// Solidity: function lockVotes(uint256 gameId) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) LockVotes(gameId *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.LockVotes(&_AgentArenaBindings.TransactOpts, gameId)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactor) Mint(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "mint", to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_AgentArenaBindings *AgentArenaBindingsSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.Mint(&_AgentArenaBindings.TransactOpts, to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.Mint(&_AgentArenaBindings.TransactOpts, to, amount)
}

// PlaceBet is a paid mutator transaction binding the contract method 0x7d0cf2e0.
//
// Solidity: function placeBet(uint256 gameId, address user, bool side, uint256 amount) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactor) PlaceBet(opts *bind.TransactOpts, gameId *big.Int, user common.Address, side bool, amount *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "placeBet", gameId, user, side, amount)
}

// PlaceBet is a paid mutator transaction binding the contract method 0x7d0cf2e0.
//
// Solidity: function placeBet(uint256 gameId, address user, bool side, uint256 amount) returns()
func (_AgentArenaBindings *AgentArenaBindingsSession) PlaceBet(gameId *big.Int, user common.Address, side bool, amount *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.PlaceBet(&_AgentArenaBindings.TransactOpts, gameId, user, side, amount)
}

// PlaceBet is a paid mutator transaction binding the contract method 0x7d0cf2e0.
//
// Solidity: function placeBet(uint256 gameId, address user, bool side, uint256 amount) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) PlaceBet(gameId *big.Int, user common.Address, side bool, amount *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.PlaceBet(&_AgentArenaBindings.TransactOpts, gameId, user, side, amount)
}

// RegisterAgent is a paid mutator transaction binding the contract method 0x7216fe48.
//
// Solidity: function registerAgent(address agentId, string name, string personality) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactor) RegisterAgent(opts *bind.TransactOpts, agentId common.Address, name string, personality string) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "registerAgent", agentId, name, personality)
}

// RegisterAgent is a paid mutator transaction binding the contract method 0x7216fe48.
//
// Solidity: function registerAgent(address agentId, string name, string personality) returns()
func (_AgentArenaBindings *AgentArenaBindingsSession) RegisterAgent(agentId common.Address, name string, personality string) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.RegisterAgent(&_AgentArenaBindings.TransactOpts, agentId, name, personality)
}

// RegisterAgent is a paid mutator transaction binding the contract method 0x7216fe48.
//
// Solidity: function registerAgent(address agentId, string name, string personality) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) RegisterAgent(agentId common.Address, name string, personality string) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.RegisterAgent(&_AgentArenaBindings.TransactOpts, agentId, name, personality)
}

// RegisterAgent0 is a paid mutator transaction binding the contract method 0x7216fe48.
//
// Solidity: function registerAgent(address agentId, string name, string personality) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactor) RegisterAgent0(opts *bind.TransactOpts, agentId common.Address, name string, personality string) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "registerAgent0", agentId, name, personality)
}

// RegisterAgent0 is a paid mutator transaction binding the contract method 0x7216fe48.
//
// Solidity: function registerAgent(address agentId, string name, string personality) returns()
func (_AgentArenaBindings *AgentArenaBindingsSession) RegisterAgent0(agentId common.Address, name string, personality string) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.RegisterAgent0(&_AgentArenaBindings.TransactOpts, agentId, name, personality)
}

// RegisterAgent0 is a paid mutator transaction binding the contract method 0x7216fe48.
//
// Solidity: function registerAgent(address agentId, string name, string personality) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) RegisterAgent0(agentId common.Address, name string, personality string) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.RegisterAgent0(&_AgentArenaBindings.TransactOpts, agentId, name, personality)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AgentArenaBindings *AgentArenaBindingsSession) RenounceOwnership() (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.RenounceOwnership(&_AgentArenaBindings.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.RenounceOwnership(&_AgentArenaBindings.TransactOpts)
}

// SetArena is a paid mutator transaction binding the contract method 0x1bd5ff28.
//
// Solidity: function setArena(address _arena) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactor) SetArena(opts *bind.TransactOpts, _arena common.Address) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "setArena", _arena)
}

// SetArena is a paid mutator transaction binding the contract method 0x1bd5ff28.
//
// Solidity: function setArena(address _arena) returns()
func (_AgentArenaBindings *AgentArenaBindingsSession) SetArena(_arena common.Address) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.SetArena(&_AgentArenaBindings.TransactOpts, _arena)
}

// SetArena is a paid mutator transaction binding the contract method 0x1bd5ff28.
//
// Solidity: function setArena(address _arena) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) SetArena(_arena common.Address) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.SetArena(&_AgentArenaBindings.TransactOpts, _arena)
}

// SetArena0 is a paid mutator transaction binding the contract method 0x1bd5ff28.
//
// Solidity: function setArena(address _arena) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactor) SetArena0(opts *bind.TransactOpts, _arena common.Address) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "setArena0", _arena)
}

// SetArena0 is a paid mutator transaction binding the contract method 0x1bd5ff28.
//
// Solidity: function setArena(address _arena) returns()
func (_AgentArenaBindings *AgentArenaBindingsSession) SetArena0(_arena common.Address) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.SetArena0(&_AgentArenaBindings.TransactOpts, _arena)
}

// SetArena0 is a paid mutator transaction binding the contract method 0x1bd5ff28.
//
// Solidity: function setArena(address _arena) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) SetArena0(_arena common.Address) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.SetArena0(&_AgentArenaBindings.TransactOpts, _arena)
}

// SetArena1 is a paid mutator transaction binding the contract method 0x1bd5ff28.
//
// Solidity: function setArena(address _arena) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactor) SetArena1(opts *bind.TransactOpts, _arena common.Address) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "setArena1", _arena)
}

// SetArena1 is a paid mutator transaction binding the contract method 0x1bd5ff28.
//
// Solidity: function setArena(address _arena) returns()
func (_AgentArenaBindings *AgentArenaBindingsSession) SetArena1(_arena common.Address) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.SetArena1(&_AgentArenaBindings.TransactOpts, _arena)
}

// SetArena1 is a paid mutator transaction binding the contract method 0x1bd5ff28.
//
// Solidity: function setArena(address _arena) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) SetArena1(_arena common.Address) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.SetArena1(&_AgentArenaBindings.TransactOpts, _arena)
}

// SetBettingDuration is a paid mutator transaction binding the contract method 0xfd16a2e7.
//
// Solidity: function setBettingDuration(uint256 _duration) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactor) SetBettingDuration(opts *bind.TransactOpts, _duration *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "setBettingDuration", _duration)
}

// SetBettingDuration is a paid mutator transaction binding the contract method 0xfd16a2e7.
//
// Solidity: function setBettingDuration(uint256 _duration) returns()
func (_AgentArenaBindings *AgentArenaBindingsSession) SetBettingDuration(_duration *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.SetBettingDuration(&_AgentArenaBindings.TransactOpts, _duration)
}

// SetBettingDuration is a paid mutator transaction binding the contract method 0xfd16a2e7.
//
// Solidity: function setBettingDuration(uint256 _duration) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) SetBettingDuration(_duration *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.SetBettingDuration(&_AgentArenaBindings.TransactOpts, _duration)
}

// Settle is a paid mutator transaction binding the contract method 0xfe417a47.
//
// Solidity: function settle(uint256 gameId, bool redWins) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactor) Settle(opts *bind.TransactOpts, gameId *big.Int, redWins bool) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "settle", gameId, redWins)
}

// Settle is a paid mutator transaction binding the contract method 0xfe417a47.
//
// Solidity: function settle(uint256 gameId, bool redWins) returns()
func (_AgentArenaBindings *AgentArenaBindingsSession) Settle(gameId *big.Int, redWins bool) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.Settle(&_AgentArenaBindings.TransactOpts, gameId, redWins)
}

// Settle is a paid mutator transaction binding the contract method 0xfe417a47.
//
// Solidity: function settle(uint256 gameId, bool redWins) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) Settle(gameId *big.Int, redWins bool) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.Settle(&_AgentArenaBindings.TransactOpts, gameId, redWins)
}

// StartGame is a paid mutator transaction binding the contract method 0xe5ed1d59.
//
// Solidity: function startGame(uint256 gameId) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactor) StartGame(opts *bind.TransactOpts, gameId *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "startGame", gameId)
}

// StartGame is a paid mutator transaction binding the contract method 0xe5ed1d59.
//
// Solidity: function startGame(uint256 gameId) returns()
func (_AgentArenaBindings *AgentArenaBindingsSession) StartGame(gameId *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.StartGame(&_AgentArenaBindings.TransactOpts, gameId)
}

// StartGame is a paid mutator transaction binding the contract method 0xe5ed1d59.
//
// Solidity: function startGame(uint256 gameId) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) StartGame(gameId *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.StartGame(&_AgentArenaBindings.TransactOpts, gameId)
}

// SubmitResult is a paid mutator transaction binding the contract method 0xe47a8059.
//
// Solidity: function submitResult(uint256 gameId, bool redWins, bytes32 actionsHash) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactor) SubmitResult(opts *bind.TransactOpts, gameId *big.Int, redWins bool, actionsHash [32]byte) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "submitResult", gameId, redWins, actionsHash)
}

// SubmitResult is a paid mutator transaction binding the contract method 0xe47a8059.
//
// Solidity: function submitResult(uint256 gameId, bool redWins, bytes32 actionsHash) returns()
func (_AgentArenaBindings *AgentArenaBindingsSession) SubmitResult(gameId *big.Int, redWins bool, actionsHash [32]byte) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.SubmitResult(&_AgentArenaBindings.TransactOpts, gameId, redWins, actionsHash)
}

// SubmitResult is a paid mutator transaction binding the contract method 0xe47a8059.
//
// Solidity: function submitResult(uint256 gameId, bool redWins, bytes32 actionsHash) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) SubmitResult(gameId *big.Int, redWins bool, actionsHash [32]byte) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.SubmitResult(&_AgentArenaBindings.TransactOpts, gameId, redWins, actionsHash)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_AgentArenaBindings *AgentArenaBindingsTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_AgentArenaBindings *AgentArenaBindingsSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.Transfer(&_AgentArenaBindings.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.Transfer(&_AgentArenaBindings.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_AgentArenaBindings *AgentArenaBindingsTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_AgentArenaBindings *AgentArenaBindingsSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.TransferFrom(&_AgentArenaBindings.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.TransferFrom(&_AgentArenaBindings.TransactOpts, from, to, value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AgentArenaBindings *AgentArenaBindingsSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.TransferOwnership(&_AgentArenaBindings.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.TransferOwnership(&_AgentArenaBindings.TransactOpts, newOwner)
}

// Vote is a paid mutator transaction binding the contract method 0xe520251a.
//
// Solidity: function vote(uint256 gameId, uint8 strategy, uint256 amount) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactor) Vote(opts *bind.TransactOpts, gameId *big.Int, strategy uint8, amount *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.contract.Transact(opts, "vote", gameId, strategy, amount)
}

// Vote is a paid mutator transaction binding the contract method 0xe520251a.
//
// Solidity: function vote(uint256 gameId, uint8 strategy, uint256 amount) returns()
func (_AgentArenaBindings *AgentArenaBindingsSession) Vote(gameId *big.Int, strategy uint8, amount *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.Vote(&_AgentArenaBindings.TransactOpts, gameId, strategy, amount)
}

// Vote is a paid mutator transaction binding the contract method 0xe520251a.
//
// Solidity: function vote(uint256 gameId, uint8 strategy, uint256 amount) returns()
func (_AgentArenaBindings *AgentArenaBindingsTransactorSession) Vote(gameId *big.Int, strategy uint8, amount *big.Int) (*types.Transaction, error) {
	return _AgentArenaBindings.Contract.Vote(&_AgentArenaBindings.TransactOpts, gameId, strategy, amount)
}

// AgentArenaBindingsAgentRegisteredIterator is returned from FilterAgentRegistered and is used to iterate over the raw logs and unpacked data for AgentRegistered events raised by the AgentArenaBindings contract.
type AgentArenaBindingsAgentRegisteredIterator struct {
	Event *AgentArenaBindingsAgentRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentArenaBindingsAgentRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentArenaBindingsAgentRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentArenaBindingsAgentRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentArenaBindingsAgentRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentArenaBindingsAgentRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentArenaBindingsAgentRegistered represents a AgentRegistered event raised by the AgentArenaBindings contract.
type AgentArenaBindingsAgentRegistered struct {
	AgentId     common.Address
	Name        string
	Personality string
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAgentRegistered is a free log retrieval operation binding the contract event 0x37b3aaf4f807403887412a5994aa1832eface4a8c15501a3311bf21d3279ea3f.
//
// Solidity: event AgentRegistered(address indexed agentId, string name, string personality)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) FilterAgentRegistered(opts *bind.FilterOpts, agentId []common.Address) (*AgentArenaBindingsAgentRegisteredIterator, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.FilterLogs(opts, "AgentRegistered", agentIdRule)
	if err != nil {
		return nil, err
	}
	return &AgentArenaBindingsAgentRegisteredIterator{contract: _AgentArenaBindings.contract, event: "AgentRegistered", logs: logs, sub: sub}, nil
}

// WatchAgentRegistered is a free log subscription operation binding the contract event 0x37b3aaf4f807403887412a5994aa1832eface4a8c15501a3311bf21d3279ea3f.
//
// Solidity: event AgentRegistered(address indexed agentId, string name, string personality)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) WatchAgentRegistered(opts *bind.WatchOpts, sink chan<- *AgentArenaBindingsAgentRegistered, agentId []common.Address) (event.Subscription, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.WatchLogs(opts, "AgentRegistered", agentIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentArenaBindingsAgentRegistered)
				if err := _AgentArenaBindings.contract.UnpackLog(event, "AgentRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAgentRegistered is a log parse operation binding the contract event 0x37b3aaf4f807403887412a5994aa1832eface4a8c15501a3311bf21d3279ea3f.
//
// Solidity: event AgentRegistered(address indexed agentId, string name, string personality)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) ParseAgentRegistered(log types.Log) (*AgentArenaBindingsAgentRegistered, error) {
	event := new(AgentArenaBindingsAgentRegistered)
	if err := _AgentArenaBindings.contract.UnpackLog(event, "AgentRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentArenaBindingsAgentWinUpdatedIterator is returned from FilterAgentWinUpdated and is used to iterate over the raw logs and unpacked data for AgentWinUpdated events raised by the AgentArenaBindings contract.
type AgentArenaBindingsAgentWinUpdatedIterator struct {
	Event *AgentArenaBindingsAgentWinUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentArenaBindingsAgentWinUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentArenaBindingsAgentWinUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentArenaBindingsAgentWinUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentArenaBindingsAgentWinUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentArenaBindingsAgentWinUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentArenaBindingsAgentWinUpdated represents a AgentWinUpdated event raised by the AgentArenaBindings contract.
type AgentArenaBindingsAgentWinUpdated struct {
	AgentId common.Address
	Wins    *big.Int
	Losses  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentWinUpdated is a free log retrieval operation binding the contract event 0x49aac4873b6a95a9d44b8a9a8b3568b4a7483801f596fc975f46951bd2d9328f.
//
// Solidity: event AgentWinUpdated(address indexed agentId, uint256 wins, uint256 losses)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) FilterAgentWinUpdated(opts *bind.FilterOpts, agentId []common.Address) (*AgentArenaBindingsAgentWinUpdatedIterator, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.FilterLogs(opts, "AgentWinUpdated", agentIdRule)
	if err != nil {
		return nil, err
	}
	return &AgentArenaBindingsAgentWinUpdatedIterator{contract: _AgentArenaBindings.contract, event: "AgentWinUpdated", logs: logs, sub: sub}, nil
}

// WatchAgentWinUpdated is a free log subscription operation binding the contract event 0x49aac4873b6a95a9d44b8a9a8b3568b4a7483801f596fc975f46951bd2d9328f.
//
// Solidity: event AgentWinUpdated(address indexed agentId, uint256 wins, uint256 losses)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) WatchAgentWinUpdated(opts *bind.WatchOpts, sink chan<- *AgentArenaBindingsAgentWinUpdated, agentId []common.Address) (event.Subscription, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.WatchLogs(opts, "AgentWinUpdated", agentIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentArenaBindingsAgentWinUpdated)
				if err := _AgentArenaBindings.contract.UnpackLog(event, "AgentWinUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAgentWinUpdated is a log parse operation binding the contract event 0x49aac4873b6a95a9d44b8a9a8b3568b4a7483801f596fc975f46951bd2d9328f.
//
// Solidity: event AgentWinUpdated(address indexed agentId, uint256 wins, uint256 losses)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) ParseAgentWinUpdated(log types.Log) (*AgentArenaBindingsAgentWinUpdated, error) {
	event := new(AgentArenaBindingsAgentWinUpdated)
	if err := _AgentArenaBindings.contract.UnpackLog(event, "AgentWinUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentArenaBindingsApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the AgentArenaBindings contract.
type AgentArenaBindingsApprovalIterator struct {
	Event *AgentArenaBindingsApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentArenaBindingsApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentArenaBindingsApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentArenaBindingsApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentArenaBindingsApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentArenaBindingsApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentArenaBindingsApproval represents a Approval event raised by the AgentArenaBindings contract.
type AgentArenaBindingsApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*AgentArenaBindingsApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &AgentArenaBindingsApprovalIterator{contract: _AgentArenaBindings.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *AgentArenaBindingsApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentArenaBindingsApproval)
				if err := _AgentArenaBindings.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) ParseApproval(log types.Log) (*AgentArenaBindingsApproval, error) {
	event := new(AgentArenaBindingsApproval)
	if err := _AgentArenaBindings.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentArenaBindingsBetPlacedIterator is returned from FilterBetPlaced and is used to iterate over the raw logs and unpacked data for BetPlaced events raised by the AgentArenaBindings contract.
type AgentArenaBindingsBetPlacedIterator struct {
	Event *AgentArenaBindingsBetPlaced // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentArenaBindingsBetPlacedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentArenaBindingsBetPlaced)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentArenaBindingsBetPlaced)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentArenaBindingsBetPlacedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentArenaBindingsBetPlacedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentArenaBindingsBetPlaced represents a BetPlaced event raised by the AgentArenaBindings contract.
type AgentArenaBindingsBetPlaced struct {
	GameId *big.Int
	User   common.Address
	Side   bool
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBetPlaced is a free log retrieval operation binding the contract event 0x4af71b021e799c62c158bd54636ca8da2fa26115a21a2dc6efe486ec104fd15f.
//
// Solidity: event BetPlaced(uint256 indexed gameId, address indexed user, bool side, uint256 amount)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) FilterBetPlaced(opts *bind.FilterOpts, gameId []*big.Int, user []common.Address) (*AgentArenaBindingsBetPlacedIterator, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.FilterLogs(opts, "BetPlaced", gameIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return &AgentArenaBindingsBetPlacedIterator{contract: _AgentArenaBindings.contract, event: "BetPlaced", logs: logs, sub: sub}, nil
}

// WatchBetPlaced is a free log subscription operation binding the contract event 0x4af71b021e799c62c158bd54636ca8da2fa26115a21a2dc6efe486ec104fd15f.
//
// Solidity: event BetPlaced(uint256 indexed gameId, address indexed user, bool side, uint256 amount)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) WatchBetPlaced(opts *bind.WatchOpts, sink chan<- *AgentArenaBindingsBetPlaced, gameId []*big.Int, user []common.Address) (event.Subscription, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.WatchLogs(opts, "BetPlaced", gameIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentArenaBindingsBetPlaced)
				if err := _AgentArenaBindings.contract.UnpackLog(event, "BetPlaced", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBetPlaced is a log parse operation binding the contract event 0x4af71b021e799c62c158bd54636ca8da2fa26115a21a2dc6efe486ec104fd15f.
//
// Solidity: event BetPlaced(uint256 indexed gameId, address indexed user, bool side, uint256 amount)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) ParseBetPlaced(log types.Log) (*AgentArenaBindingsBetPlaced, error) {
	event := new(AgentArenaBindingsBetPlaced)
	if err := _AgentArenaBindings.contract.UnpackLog(event, "BetPlaced", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentArenaBindingsBettingLockedIterator is returned from FilterBettingLocked and is used to iterate over the raw logs and unpacked data for BettingLocked events raised by the AgentArenaBindings contract.
type AgentArenaBindingsBettingLockedIterator struct {
	Event *AgentArenaBindingsBettingLocked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentArenaBindingsBettingLockedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentArenaBindingsBettingLocked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentArenaBindingsBettingLocked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentArenaBindingsBettingLockedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentArenaBindingsBettingLockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentArenaBindingsBettingLocked represents a BettingLocked event raised by the AgentArenaBindings contract.
type AgentArenaBindingsBettingLocked struct {
	GameId       *big.Int
	TotalBetRed  *big.Int
	TotalBetBlue *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterBettingLocked is a free log retrieval operation binding the contract event 0xedaa0e3892e0f762797f49c5028c0a64f5587b1d1f970d98402062d5b6dd4d9d.
//
// Solidity: event BettingLocked(uint256 indexed gameId, uint256 totalBetRed, uint256 totalBetBlue)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) FilterBettingLocked(opts *bind.FilterOpts, gameId []*big.Int) (*AgentArenaBindingsBettingLockedIterator, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.FilterLogs(opts, "BettingLocked", gameIdRule)
	if err != nil {
		return nil, err
	}
	return &AgentArenaBindingsBettingLockedIterator{contract: _AgentArenaBindings.contract, event: "BettingLocked", logs: logs, sub: sub}, nil
}

// WatchBettingLocked is a free log subscription operation binding the contract event 0xedaa0e3892e0f762797f49c5028c0a64f5587b1d1f970d98402062d5b6dd4d9d.
//
// Solidity: event BettingLocked(uint256 indexed gameId, uint256 totalBetRed, uint256 totalBetBlue)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) WatchBettingLocked(opts *bind.WatchOpts, sink chan<- *AgentArenaBindingsBettingLocked, gameId []*big.Int) (event.Subscription, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.WatchLogs(opts, "BettingLocked", gameIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentArenaBindingsBettingLocked)
				if err := _AgentArenaBindings.contract.UnpackLog(event, "BettingLocked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBettingLocked is a log parse operation binding the contract event 0xedaa0e3892e0f762797f49c5028c0a64f5587b1d1f970d98402062d5b6dd4d9d.
//
// Solidity: event BettingLocked(uint256 indexed gameId, uint256 totalBetRed, uint256 totalBetBlue)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) ParseBettingLocked(log types.Log) (*AgentArenaBindingsBettingLocked, error) {
	event := new(AgentArenaBindingsBettingLocked)
	if err := _AgentArenaBindings.contract.UnpackLog(event, "BettingLocked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentArenaBindingsGameCreatedIterator is returned from FilterGameCreated and is used to iterate over the raw logs and unpacked data for GameCreated events raised by the AgentArenaBindings contract.
type AgentArenaBindingsGameCreatedIterator struct {
	Event *AgentArenaBindingsGameCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentArenaBindingsGameCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentArenaBindingsGameCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentArenaBindingsGameCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentArenaBindingsGameCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentArenaBindingsGameCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentArenaBindingsGameCreated represents a GameCreated event raised by the AgentArenaBindings contract.
type AgentArenaBindingsGameCreated struct {
	GameId          *big.Int
	AgentRed        common.Address
	AgentBlue       common.Address
	BettingDeadline *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterGameCreated is a free log retrieval operation binding the contract event 0x6200407c0ea392b8107b21a9be480acd41fda186d04bed28cc7da2d4b53d56e2.
//
// Solidity: event GameCreated(uint256 indexed gameId, address agentRed, address agentBlue, uint256 bettingDeadline)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) FilterGameCreated(opts *bind.FilterOpts, gameId []*big.Int) (*AgentArenaBindingsGameCreatedIterator, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.FilterLogs(opts, "GameCreated", gameIdRule)
	if err != nil {
		return nil, err
	}
	return &AgentArenaBindingsGameCreatedIterator{contract: _AgentArenaBindings.contract, event: "GameCreated", logs: logs, sub: sub}, nil
}

// WatchGameCreated is a free log subscription operation binding the contract event 0x6200407c0ea392b8107b21a9be480acd41fda186d04bed28cc7da2d4b53d56e2.
//
// Solidity: event GameCreated(uint256 indexed gameId, address agentRed, address agentBlue, uint256 bettingDeadline)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) WatchGameCreated(opts *bind.WatchOpts, sink chan<- *AgentArenaBindingsGameCreated, gameId []*big.Int) (event.Subscription, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.WatchLogs(opts, "GameCreated", gameIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentArenaBindingsGameCreated)
				if err := _AgentArenaBindings.contract.UnpackLog(event, "GameCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseGameCreated is a log parse operation binding the contract event 0x6200407c0ea392b8107b21a9be480acd41fda186d04bed28cc7da2d4b53d56e2.
//
// Solidity: event GameCreated(uint256 indexed gameId, address agentRed, address agentBlue, uint256 bettingDeadline)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) ParseGameCreated(log types.Log) (*AgentArenaBindingsGameCreated, error) {
	event := new(AgentArenaBindingsGameCreated)
	if err := _AgentArenaBindings.contract.UnpackLog(event, "GameCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentArenaBindingsGameFinishedIterator is returned from FilterGameFinished and is used to iterate over the raw logs and unpacked data for GameFinished events raised by the AgentArenaBindings contract.
type AgentArenaBindingsGameFinishedIterator struct {
	Event *AgentArenaBindingsGameFinished // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentArenaBindingsGameFinishedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentArenaBindingsGameFinished)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentArenaBindingsGameFinished)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentArenaBindingsGameFinishedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentArenaBindingsGameFinishedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentArenaBindingsGameFinished represents a GameFinished event raised by the AgentArenaBindings contract.
type AgentArenaBindingsGameFinished struct {
	GameId  *big.Int
	RedWins bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterGameFinished is a free log retrieval operation binding the contract event 0xb1134e1759d5f0c30f41cf6236fce19daad8b3c8ffd208c3c29388a64aa19729.
//
// Solidity: event GameFinished(uint256 indexed gameId, bool redWins)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) FilterGameFinished(opts *bind.FilterOpts, gameId []*big.Int) (*AgentArenaBindingsGameFinishedIterator, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.FilterLogs(opts, "GameFinished", gameIdRule)
	if err != nil {
		return nil, err
	}
	return &AgentArenaBindingsGameFinishedIterator{contract: _AgentArenaBindings.contract, event: "GameFinished", logs: logs, sub: sub}, nil
}

// WatchGameFinished is a free log subscription operation binding the contract event 0xb1134e1759d5f0c30f41cf6236fce19daad8b3c8ffd208c3c29388a64aa19729.
//
// Solidity: event GameFinished(uint256 indexed gameId, bool redWins)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) WatchGameFinished(opts *bind.WatchOpts, sink chan<- *AgentArenaBindingsGameFinished, gameId []*big.Int) (event.Subscription, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.WatchLogs(opts, "GameFinished", gameIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentArenaBindingsGameFinished)
				if err := _AgentArenaBindings.contract.UnpackLog(event, "GameFinished", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseGameFinished is a log parse operation binding the contract event 0xb1134e1759d5f0c30f41cf6236fce19daad8b3c8ffd208c3c29388a64aa19729.
//
// Solidity: event GameFinished(uint256 indexed gameId, bool redWins)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) ParseGameFinished(log types.Log) (*AgentArenaBindingsGameFinished, error) {
	event := new(AgentArenaBindingsGameFinished)
	if err := _AgentArenaBindings.contract.UnpackLog(event, "GameFinished", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentArenaBindingsGameInitializedIterator is returned from FilterGameInitialized and is used to iterate over the raw logs and unpacked data for GameInitialized events raised by the AgentArenaBindings contract.
type AgentArenaBindingsGameInitializedIterator struct {
	Event *AgentArenaBindingsGameInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentArenaBindingsGameInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentArenaBindingsGameInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentArenaBindingsGameInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentArenaBindingsGameInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentArenaBindingsGameInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentArenaBindingsGameInitialized represents a GameInitialized event raised by the AgentArenaBindings contract.
type AgentArenaBindingsGameInitialized struct {
	GameId    *big.Int
	AgentRed  common.Address
	AgentBlue common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterGameInitialized is a free log retrieval operation binding the contract event 0x62300909f39859a4dacb5c7c8698930d598ec7cda8b5d526433878604a30845c.
//
// Solidity: event GameInitialized(uint256 indexed gameId, address agentRed, address agentBlue)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) FilterGameInitialized(opts *bind.FilterOpts, gameId []*big.Int) (*AgentArenaBindingsGameInitializedIterator, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.FilterLogs(opts, "GameInitialized", gameIdRule)
	if err != nil {
		return nil, err
	}
	return &AgentArenaBindingsGameInitializedIterator{contract: _AgentArenaBindings.contract, event: "GameInitialized", logs: logs, sub: sub}, nil
}

// WatchGameInitialized is a free log subscription operation binding the contract event 0x62300909f39859a4dacb5c7c8698930d598ec7cda8b5d526433878604a30845c.
//
// Solidity: event GameInitialized(uint256 indexed gameId, address agentRed, address agentBlue)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) WatchGameInitialized(opts *bind.WatchOpts, sink chan<- *AgentArenaBindingsGameInitialized, gameId []*big.Int) (event.Subscription, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.WatchLogs(opts, "GameInitialized", gameIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentArenaBindingsGameInitialized)
				if err := _AgentArenaBindings.contract.UnpackLog(event, "GameInitialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseGameInitialized is a log parse operation binding the contract event 0x62300909f39859a4dacb5c7c8698930d598ec7cda8b5d526433878604a30845c.
//
// Solidity: event GameInitialized(uint256 indexed gameId, address agentRed, address agentBlue)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) ParseGameInitialized(log types.Log) (*AgentArenaBindingsGameInitialized, error) {
	event := new(AgentArenaBindingsGameInitialized)
	if err := _AgentArenaBindings.contract.UnpackLog(event, "GameInitialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentArenaBindingsGameResultSubmittedIterator is returned from FilterGameResultSubmitted and is used to iterate over the raw logs and unpacked data for GameResultSubmitted events raised by the AgentArenaBindings contract.
type AgentArenaBindingsGameResultSubmittedIterator struct {
	Event *AgentArenaBindingsGameResultSubmitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentArenaBindingsGameResultSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentArenaBindingsGameResultSubmitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentArenaBindingsGameResultSubmitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentArenaBindingsGameResultSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentArenaBindingsGameResultSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentArenaBindingsGameResultSubmitted represents a GameResultSubmitted event raised by the AgentArenaBindings contract.
type AgentArenaBindingsGameResultSubmitted struct {
	GameId      *big.Int
	RedWins     bool
	ActionsHash [32]byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterGameResultSubmitted is a free log retrieval operation binding the contract event 0xe366c87d746dacfdd06a572c2761fbe019e661720af2bc28520b2187ab07be0b.
//
// Solidity: event GameResultSubmitted(uint256 indexed gameId, bool redWins, bytes32 actionsHash)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) FilterGameResultSubmitted(opts *bind.FilterOpts, gameId []*big.Int) (*AgentArenaBindingsGameResultSubmittedIterator, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.FilterLogs(opts, "GameResultSubmitted", gameIdRule)
	if err != nil {
		return nil, err
	}
	return &AgentArenaBindingsGameResultSubmittedIterator{contract: _AgentArenaBindings.contract, event: "GameResultSubmitted", logs: logs, sub: sub}, nil
}

// WatchGameResultSubmitted is a free log subscription operation binding the contract event 0xe366c87d746dacfdd06a572c2761fbe019e661720af2bc28520b2187ab07be0b.
//
// Solidity: event GameResultSubmitted(uint256 indexed gameId, bool redWins, bytes32 actionsHash)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) WatchGameResultSubmitted(opts *bind.WatchOpts, sink chan<- *AgentArenaBindingsGameResultSubmitted, gameId []*big.Int) (event.Subscription, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.WatchLogs(opts, "GameResultSubmitted", gameIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentArenaBindingsGameResultSubmitted)
				if err := _AgentArenaBindings.contract.UnpackLog(event, "GameResultSubmitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseGameResultSubmitted is a log parse operation binding the contract event 0xe366c87d746dacfdd06a572c2761fbe019e661720af2bc28520b2187ab07be0b.
//
// Solidity: event GameResultSubmitted(uint256 indexed gameId, bool redWins, bytes32 actionsHash)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) ParseGameResultSubmitted(log types.Log) (*AgentArenaBindingsGameResultSubmitted, error) {
	event := new(AgentArenaBindingsGameResultSubmitted)
	if err := _AgentArenaBindings.contract.UnpackLog(event, "GameResultSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentArenaBindingsGameSettledIterator is returned from FilterGameSettled and is used to iterate over the raw logs and unpacked data for GameSettled events raised by the AgentArenaBindings contract.
type AgentArenaBindingsGameSettledIterator struct {
	Event *AgentArenaBindingsGameSettled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentArenaBindingsGameSettledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentArenaBindingsGameSettled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentArenaBindingsGameSettled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentArenaBindingsGameSettledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentArenaBindingsGameSettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentArenaBindingsGameSettled represents a GameSettled event raised by the AgentArenaBindings contract.
type AgentArenaBindingsGameSettled struct {
	GameId      *big.Int
	RedWins     bool
	TotalPool   *big.Int
	ProtocolFee *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterGameSettled is a free log retrieval operation binding the contract event 0xd6f7947fbe1708ac19baf453d7b00bc2b588e34ba696213de9275eaaa6d70065.
//
// Solidity: event GameSettled(uint256 indexed gameId, bool redWins, uint256 totalPool, uint256 protocolFee)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) FilterGameSettled(opts *bind.FilterOpts, gameId []*big.Int) (*AgentArenaBindingsGameSettledIterator, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.FilterLogs(opts, "GameSettled", gameIdRule)
	if err != nil {
		return nil, err
	}
	return &AgentArenaBindingsGameSettledIterator{contract: _AgentArenaBindings.contract, event: "GameSettled", logs: logs, sub: sub}, nil
}

// WatchGameSettled is a free log subscription operation binding the contract event 0xd6f7947fbe1708ac19baf453d7b00bc2b588e34ba696213de9275eaaa6d70065.
//
// Solidity: event GameSettled(uint256 indexed gameId, bool redWins, uint256 totalPool, uint256 protocolFee)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) WatchGameSettled(opts *bind.WatchOpts, sink chan<- *AgentArenaBindingsGameSettled, gameId []*big.Int) (event.Subscription, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.WatchLogs(opts, "GameSettled", gameIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentArenaBindingsGameSettled)
				if err := _AgentArenaBindings.contract.UnpackLog(event, "GameSettled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseGameSettled is a log parse operation binding the contract event 0xd6f7947fbe1708ac19baf453d7b00bc2b588e34ba696213de9275eaaa6d70065.
//
// Solidity: event GameSettled(uint256 indexed gameId, bool redWins, uint256 totalPool, uint256 protocolFee)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) ParseGameSettled(log types.Log) (*AgentArenaBindingsGameSettled, error) {
	event := new(AgentArenaBindingsGameSettled)
	if err := _AgentArenaBindings.contract.UnpackLog(event, "GameSettled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentArenaBindingsGameStartedIterator is returned from FilterGameStarted and is used to iterate over the raw logs and unpacked data for GameStarted events raised by the AgentArenaBindings contract.
type AgentArenaBindingsGameStartedIterator struct {
	Event *AgentArenaBindingsGameStarted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentArenaBindingsGameStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentArenaBindingsGameStarted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentArenaBindingsGameStarted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentArenaBindingsGameStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentArenaBindingsGameStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentArenaBindingsGameStarted represents a GameStarted event raised by the AgentArenaBindings contract.
type AgentArenaBindingsGameStarted struct {
	GameId *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterGameStarted is a free log retrieval operation binding the contract event 0x50ad08f58a27f2851d7e3a1b3a6a46b290f2ce677e99642d30ff639721e77790.
//
// Solidity: event GameStarted(uint256 indexed gameId)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) FilterGameStarted(opts *bind.FilterOpts, gameId []*big.Int) (*AgentArenaBindingsGameStartedIterator, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.FilterLogs(opts, "GameStarted", gameIdRule)
	if err != nil {
		return nil, err
	}
	return &AgentArenaBindingsGameStartedIterator{contract: _AgentArenaBindings.contract, event: "GameStarted", logs: logs, sub: sub}, nil
}

// WatchGameStarted is a free log subscription operation binding the contract event 0x50ad08f58a27f2851d7e3a1b3a6a46b290f2ce677e99642d30ff639721e77790.
//
// Solidity: event GameStarted(uint256 indexed gameId)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) WatchGameStarted(opts *bind.WatchOpts, sink chan<- *AgentArenaBindingsGameStarted, gameId []*big.Int) (event.Subscription, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.WatchLogs(opts, "GameStarted", gameIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentArenaBindingsGameStarted)
				if err := _AgentArenaBindings.contract.UnpackLog(event, "GameStarted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseGameStarted is a log parse operation binding the contract event 0x50ad08f58a27f2851d7e3a1b3a6a46b290f2ce677e99642d30ff639721e77790.
//
// Solidity: event GameStarted(uint256 indexed gameId)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) ParseGameStarted(log types.Log) (*AgentArenaBindingsGameStarted, error) {
	event := new(AgentArenaBindingsGameStarted)
	if err := _AgentArenaBindings.contract.UnpackLog(event, "GameStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentArenaBindingsOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the AgentArenaBindings contract.
type AgentArenaBindingsOwnershipTransferredIterator struct {
	Event *AgentArenaBindingsOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentArenaBindingsOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentArenaBindingsOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentArenaBindingsOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentArenaBindingsOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentArenaBindingsOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentArenaBindingsOwnershipTransferred represents a OwnershipTransferred event raised by the AgentArenaBindings contract.
type AgentArenaBindingsOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*AgentArenaBindingsOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &AgentArenaBindingsOwnershipTransferredIterator{contract: _AgentArenaBindings.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AgentArenaBindingsOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentArenaBindingsOwnershipTransferred)
				if err := _AgentArenaBindings.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) ParseOwnershipTransferred(log types.Log) (*AgentArenaBindingsOwnershipTransferred, error) {
	event := new(AgentArenaBindingsOwnershipTransferred)
	if err := _AgentArenaBindings.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentArenaBindingsRewardClaimedIterator is returned from FilterRewardClaimed and is used to iterate over the raw logs and unpacked data for RewardClaimed events raised by the AgentArenaBindings contract.
type AgentArenaBindingsRewardClaimedIterator struct {
	Event *AgentArenaBindingsRewardClaimed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentArenaBindingsRewardClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentArenaBindingsRewardClaimed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentArenaBindingsRewardClaimed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentArenaBindingsRewardClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentArenaBindingsRewardClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentArenaBindingsRewardClaimed represents a RewardClaimed event raised by the AgentArenaBindings contract.
type AgentArenaBindingsRewardClaimed struct {
	GameId *big.Int
	User   common.Address
	Reward *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRewardClaimed is a free log retrieval operation binding the contract event 0x24b5efa61dd1cfc659205a97fb8ed868f3cb8c81922bab2b96423e5de1de2cb7.
//
// Solidity: event RewardClaimed(uint256 indexed gameId, address indexed user, uint256 reward)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) FilterRewardClaimed(opts *bind.FilterOpts, gameId []*big.Int, user []common.Address) (*AgentArenaBindingsRewardClaimedIterator, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.FilterLogs(opts, "RewardClaimed", gameIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return &AgentArenaBindingsRewardClaimedIterator{contract: _AgentArenaBindings.contract, event: "RewardClaimed", logs: logs, sub: sub}, nil
}

// WatchRewardClaimed is a free log subscription operation binding the contract event 0x24b5efa61dd1cfc659205a97fb8ed868f3cb8c81922bab2b96423e5de1de2cb7.
//
// Solidity: event RewardClaimed(uint256 indexed gameId, address indexed user, uint256 reward)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) WatchRewardClaimed(opts *bind.WatchOpts, sink chan<- *AgentArenaBindingsRewardClaimed, gameId []*big.Int, user []common.Address) (event.Subscription, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.WatchLogs(opts, "RewardClaimed", gameIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentArenaBindingsRewardClaimed)
				if err := _AgentArenaBindings.contract.UnpackLog(event, "RewardClaimed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRewardClaimed is a log parse operation binding the contract event 0x24b5efa61dd1cfc659205a97fb8ed868f3cb8c81922bab2b96423e5de1de2cb7.
//
// Solidity: event RewardClaimed(uint256 indexed gameId, address indexed user, uint256 reward)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) ParseRewardClaimed(log types.Log) (*AgentArenaBindingsRewardClaimed, error) {
	event := new(AgentArenaBindingsRewardClaimed)
	if err := _AgentArenaBindings.contract.UnpackLog(event, "RewardClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentArenaBindingsStrategyVotedIterator is returned from FilterStrategyVoted and is used to iterate over the raw logs and unpacked data for StrategyVoted events raised by the AgentArenaBindings contract.
type AgentArenaBindingsStrategyVotedIterator struct {
	Event *AgentArenaBindingsStrategyVoted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentArenaBindingsStrategyVotedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentArenaBindingsStrategyVoted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentArenaBindingsStrategyVoted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentArenaBindingsStrategyVotedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentArenaBindingsStrategyVotedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentArenaBindingsStrategyVoted represents a StrategyVoted event raised by the AgentArenaBindings contract.
type AgentArenaBindingsStrategyVoted struct {
	GameId   *big.Int
	User     common.Address
	Strategy uint8
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterStrategyVoted is a free log retrieval operation binding the contract event 0x74e9434d1a2e6e628776e1b8d3b385445687b844e81c3265e7a12bfffedba110.
//
// Solidity: event StrategyVoted(uint256 indexed gameId, address indexed user, uint8 strategy, uint256 amount)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) FilterStrategyVoted(opts *bind.FilterOpts, gameId []*big.Int, user []common.Address) (*AgentArenaBindingsStrategyVotedIterator, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.FilterLogs(opts, "StrategyVoted", gameIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return &AgentArenaBindingsStrategyVotedIterator{contract: _AgentArenaBindings.contract, event: "StrategyVoted", logs: logs, sub: sub}, nil
}

// WatchStrategyVoted is a free log subscription operation binding the contract event 0x74e9434d1a2e6e628776e1b8d3b385445687b844e81c3265e7a12bfffedba110.
//
// Solidity: event StrategyVoted(uint256 indexed gameId, address indexed user, uint8 strategy, uint256 amount)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) WatchStrategyVoted(opts *bind.WatchOpts, sink chan<- *AgentArenaBindingsStrategyVoted, gameId []*big.Int, user []common.Address) (event.Subscription, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.WatchLogs(opts, "StrategyVoted", gameIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentArenaBindingsStrategyVoted)
				if err := _AgentArenaBindings.contract.UnpackLog(event, "StrategyVoted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseStrategyVoted is a log parse operation binding the contract event 0x74e9434d1a2e6e628776e1b8d3b385445687b844e81c3265e7a12bfffedba110.
//
// Solidity: event StrategyVoted(uint256 indexed gameId, address indexed user, uint8 strategy, uint256 amount)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) ParseStrategyVoted(log types.Log) (*AgentArenaBindingsStrategyVoted, error) {
	event := new(AgentArenaBindingsStrategyVoted)
	if err := _AgentArenaBindings.contract.UnpackLog(event, "StrategyVoted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentArenaBindingsTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the AgentArenaBindings contract.
type AgentArenaBindingsTransferIterator struct {
	Event *AgentArenaBindingsTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentArenaBindingsTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentArenaBindingsTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentArenaBindingsTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentArenaBindingsTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentArenaBindingsTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentArenaBindingsTransfer represents a Transfer event raised by the AgentArenaBindings contract.
type AgentArenaBindingsTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AgentArenaBindingsTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AgentArenaBindingsTransferIterator{contract: _AgentArenaBindings.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *AgentArenaBindingsTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentArenaBindingsTransfer)
				if err := _AgentArenaBindings.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) ParseTransfer(log types.Log) (*AgentArenaBindingsTransfer, error) {
	event := new(AgentArenaBindingsTransfer)
	if err := _AgentArenaBindings.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentArenaBindingsVotesLockedIterator is returned from FilterVotesLocked and is used to iterate over the raw logs and unpacked data for VotesLocked events raised by the AgentArenaBindings contract.
type AgentArenaBindingsVotesLockedIterator struct {
	Event *AgentArenaBindingsVotesLocked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentArenaBindingsVotesLockedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentArenaBindingsVotesLocked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentArenaBindingsVotesLocked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentArenaBindingsVotesLockedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentArenaBindingsVotesLockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentArenaBindingsVotesLocked represents a VotesLocked event raised by the AgentArenaBindings contract.
type AgentArenaBindingsVotesLocked struct {
	GameId     *big.Int
	Aggressive *big.Int
	Defensive  *big.Int
	Tricky     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVotesLocked is a free log retrieval operation binding the contract event 0xbb0862b7edaa06a3c636517b28f9d5f1d00cfd56d8d041ecdf06867a0788a503.
//
// Solidity: event VotesLocked(uint256 indexed gameId, uint256 aggressive, uint256 defensive, uint256 tricky)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) FilterVotesLocked(opts *bind.FilterOpts, gameId []*big.Int) (*AgentArenaBindingsVotesLockedIterator, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.FilterLogs(opts, "VotesLocked", gameIdRule)
	if err != nil {
		return nil, err
	}
	return &AgentArenaBindingsVotesLockedIterator{contract: _AgentArenaBindings.contract, event: "VotesLocked", logs: logs, sub: sub}, nil
}

// WatchVotesLocked is a free log subscription operation binding the contract event 0xbb0862b7edaa06a3c636517b28f9d5f1d00cfd56d8d041ecdf06867a0788a503.
//
// Solidity: event VotesLocked(uint256 indexed gameId, uint256 aggressive, uint256 defensive, uint256 tricky)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) WatchVotesLocked(opts *bind.WatchOpts, sink chan<- *AgentArenaBindingsVotesLocked, gameId []*big.Int) (event.Subscription, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _AgentArenaBindings.contract.WatchLogs(opts, "VotesLocked", gameIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentArenaBindingsVotesLocked)
				if err := _AgentArenaBindings.contract.UnpackLog(event, "VotesLocked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseVotesLocked is a log parse operation binding the contract event 0xbb0862b7edaa06a3c636517b28f9d5f1d00cfd56d8d041ecdf06867a0788a503.
//
// Solidity: event VotesLocked(uint256 indexed gameId, uint256 aggressive, uint256 defensive, uint256 tricky)
func (_AgentArenaBindings *AgentArenaBindingsFilterer) ParseVotesLocked(log types.Log) (*AgentArenaBindingsVotesLocked, error) {
	event := new(AgentArenaBindingsVotesLocked)
	if err := _AgentArenaBindings.contract.UnpackLog(event, "VotesLocked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
