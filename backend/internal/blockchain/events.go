package blockchain

import "math/big"

// EventType 事件类型
type EventType string

const (
	EventBetPlaced    EventType = "BetPlaced"
	EventBettingLocked EventType = "BettingLocked"
	EventGameSettled  EventType = "GameSettled"
	EventVotesLocked  EventType = "VotesLocked"
	EventRewardClaimed EventType = "RewardClaimed"
	EventGameCreated  EventType = "GameCreated"
	EventGameResultSubmitted EventType = "GameResultSubmitted"
)

// ContractEvent 合约事件
type ContractEvent struct {
	Type    EventType
	GameID  uint64
	Data    interface{}
	Block   uint64
	TxHash  string
}

// BetPlacedData BetPlaced 事件数据
type BetPlacedData struct {
	User   string
	Side   bool   // true=Red, false=Blue
	Amount *big.Int
}

// BettingLockedData BettingLocked 事件数据
type BettingLockedData struct {
	TotalBetRed  *big.Int
	TotalBetBlue *big.Int
}

// GameSettledData GameSettled 事件数据
type GameSettledData struct {
	RedWins     bool
	TotalPool   *big.Int
	ProtocolFee *big.Int
}

// VotesLockedData VotesLocked 事件数据
type VotesLockedData struct {
	Aggressive uint64
	Defensive  uint64
	Tricky     uint64
}

// RewardClaimedData RewardClaimed 事件数据
type RewardClaimedData struct {
	User   string
	Reward *big.Int
}

// GameCreatedData GameCreated 事件数据
type GameCreatedData struct {
	AgentRed        string
	AgentBlue       string
	BettingDeadline uint64
}

// GameResultSubmittedData GameResultSubmitted 事件数据
type GameResultSubmittedData struct {
	RedWins     bool
	ActionsHash [32]byte
}
