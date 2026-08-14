// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "./ArenaTypes.sol";

/// @title StrategyVoting - 记录策略投票，一次性锁定
contract StrategyVoting {
    address public owner;
    address public arena;

    // gameId => StrategyVoteRecord
    mapping(uint256 => ArenaTypes.StrategyVoteRecord) public votes;

    // gameId => user => (strategy, amount)
    struct UserVote {
        ArenaTypes.Strategy strategy;
        uint256 amount;
    }
    mapping(uint256 => mapping(address => UserVote)) public userVotes;

    event StrategyVoted(uint256 indexed gameId, address indexed user, ArenaTypes.Strategy strategy, uint256 amount);
    event VotesLocked(uint256 indexed gameId, uint256 aggressive, uint256 defensive, uint256 tricky);

    modifier onlyOwner() {
        require(msg.sender == owner, "not owner");
        _;
    }

    modifier onlyArena() {
        require(msg.sender == arena, "not arena");
        _;
    }

    constructor() {
        owner = msg.sender;
    }

    function setArena(address _arena) external onlyOwner {
        arena = _arena;
    }

    /// @notice 投票（下注时调用）
    function vote(uint256 gameId, ArenaTypes.Strategy strategy, uint256 amount) external onlyArena {
        ArenaTypes.StrategyVoteRecord storage record = votes[gameId];
        require(!record.locked, "votes locked");

        // Update totals
        if (strategy == ArenaTypes.Strategy.Aggressive) {
            record.aggressive += amount;
        } else if (strategy == ArenaTypes.Strategy.Defensive) {
            record.defensive += amount;
        } else {
            record.tricky += amount;
        }

        // Record user vote
        userVotes[gameId][msg.sender] = UserVote({strategy: strategy, amount: amount});

        emit StrategyVoted(gameId, msg.sender, strategy, amount);
    }

    /// @notice 锁定投票（开局时）
    function lockVotes(uint256 gameId) external onlyArena {
        ArenaTypes.StrategyVoteRecord storage record = votes[gameId];
        require(!record.locked, "already locked");
        record.locked = true;
        emit VotesLocked(gameId, record.aggressive, record.defensive, record.tricky);
    }

    /// @notice 查询策略权重（百分比 0-100）
    function getStrategyWeights(uint256 gameId)
        external
        view
        returns (uint256 aggressiveWeight, uint256 defensiveWeight, uint256 trickyWeight)
    {
        ArenaTypes.StrategyVoteRecord storage record = votes[gameId];
        uint256 total = record.aggressive + record.defensive + record.tricky;

        if (total == 0) {
            return (33, 33, 34); // 默认均分
        }

        aggressiveWeight = record.aggressive * 100 / total;
        defensiveWeight = record.defensive * 100 / total;
        trickyWeight = 100 - aggressiveWeight - defensiveWeight; // 确保总和 = 100
    }

    function getUserVote(uint256 gameId, address user)
        external
        view
        returns (ArenaTypes.Strategy strategy, uint256 amount)
    {
        UserVote storage uv = userVotes[gameId][user];
        return (uv.strategy, uv.amount);
    }

    function getVoteRecord(uint256 gameId) external view returns (ArenaTypes.StrategyVoteRecord memory) {
        return votes[gameId];
    }
}
