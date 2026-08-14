// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/// @title 共享类型定义
library ArenaTypes {
    enum Side { None, Red, Blue }

    enum GameStatus { Open, Locked, Finished }

    enum Strategy { Aggressive, Defensive, Tricky }

    struct GameInfo {
        uint256 gameId;
        address agentRed;
        address agentBlue;
        uint256 totalBetRed;
        uint256 totalBetBlue;
        uint256 bettingDeadline;
        GameStatus status;
        Side winner;
    }

    struct StrategyVoteRecord {
        uint256 aggressive;
        uint256 defensive;
        uint256 tricky;
        bool locked;
    }
}
