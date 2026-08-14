// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "./ArenaTypes.sol";
import "./BettingPool.sol";
import "./StrategyVoting.sol";
import "./GameRegistry.sol";

/// @title AgentArena - 主合约，统一入口
contract AgentArena {
    using SafeERC20 for IERC20;

    address public owner;
    address public usdc;

    BettingPool public bettingPool;
    StrategyVoting public strategyVoting;
    GameRegistry public gameRegistry;

    uint256 public bettingDuration = 120; // 2 minutes

    event GameInitialized(uint256 indexed gameId, address agentRed, address agentBlue);
    event GameStarted(uint256 indexed gameId);
    event GameFinished(uint256 indexed gameId, bool redWins);

    modifier onlyOwner() {
        require(msg.sender == owner, "not owner");
        _;
    }

    constructor(
        address _usdc,
        address _protocolTreasury
    ) {
        owner = msg.sender;
        usdc = _usdc;

        // Deploy sub-contracts
        bettingPool = new BettingPool(_usdc, _protocolTreasury);
        strategyVoting = new StrategyVoting();
        gameRegistry = new GameRegistry();

        // Set arena address in sub-contracts
        bettingPool.setArena(address(this));
        strategyVoting.setArena(address(this));
        gameRegistry.setArena(address(this));
    }

    /// @notice 注册 Agent
    function registerAgent(address agentId, string calldata name, string calldata personality) external onlyOwner {
        gameRegistry.registerAgent(agentId, name, personality);
    }

    /// @notice 创建新对局
    function createGame(address agentRed, address agentBlue) external onlyOwner returns (uint256 gameId) {
        uint256 deadline = block.timestamp + bettingDuration;

        // Create in registry
        gameId = gameRegistry.createGame(agentRed, agentBlue, deadline);

        // Init betting pool
        bettingPool.initGame(gameId, agentRed, agentBlue, deadline);

        emit GameInitialized(gameId, agentRed, agentBlue);
    }

    /// @notice 用户下注 + 投票（一步完成）
    function betAndVote(
        uint256 gameId,
        bool side,
        uint256 amount,
        ArenaTypes.Strategy strategy
    ) external {
        // Transfer USDC from user to this contract, then to BettingPool
        IERC20(usdc).safeTransferFrom(msg.sender, address(this), amount);
        IERC20(usdc).safeTransfer(address(bettingPool), amount);

        // Record bet and vote
        bettingPool.placeBet(gameId, msg.sender, side, amount);
        strategyVoting.vote(gameId, strategy, amount);
    }

    /// @notice 开局（锁定下注 + 锁定投票）
    function startGame(uint256 gameId) external onlyOwner {
        bettingPool.lockBetting(gameId);
        strategyVoting.lockVotes(gameId);
        emit GameStarted(gameId);
    }

    /// @notice 结算（提交结果 + 结算奖金）
    function finishGame(uint256 gameId, bool redWins, bytes32 actionsHash) external onlyOwner {
        gameRegistry.submitResult(gameId, redWins, actionsHash);
        bettingPool.settle(gameId, redWins);
        emit GameFinished(gameId, redWins);
    }

    /// @notice 设置下注时长
    function setBettingDuration(uint256 _duration) external onlyOwner {
        bettingDuration = _duration;
    }

    // ============ View functions ============

    function getStrategyWeights(uint256 gameId)
        external
        view
        returns (uint256 aggressive, uint256 defensive, uint256 tricky)
    {
        return strategyVoting.getStrategyWeights(gameId);
    }

    function getGame(uint256 gameId) external view returns (ArenaTypes.GameInfo memory) {
        return bettingPool.getGame(gameId);
    }

    function getOdds(uint256 gameId) external view returns (uint256 oddsRed, uint256 oddsBlue) {
        return bettingPool.getOdds(gameId);
    }

    function getReward(uint256 gameId, address user) external view returns (uint256) {
        return bettingPool.getReward(gameId, user);
    }

    function getAgent(address agentId) external view returns (GameRegistry.AgentInfo memory) {
        return gameRegistry.getAgent(agentId);
    }
}
