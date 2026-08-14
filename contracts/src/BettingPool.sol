// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "./ArenaTypes.sol";

/// @title BettingPool - 管理下注资金、锁定、结算
contract BettingPool is ReentrancyGuard {
    using SafeERC20 for IERC20;

    address public owner;
    address public usdc;
    address public arena; // AgentArena 主合约
    uint256 public protocolFeeBps = 500; // 5%
    address public protocolTreasury;

    // gameId => GameInfo
    mapping(uint256 => ArenaTypes.GameInfo) public games;

    // gameId => user => bet info
    struct BetInfo {
        ArenaTypes.Side side;
        uint256 amount;
        bool claimed;
    }
    mapping(uint256 => mapping(address => BetInfo)) public bets;

    // gameId => total bet per side (for winner payout calc)
    mapping(uint256 => uint256) public winnerPools;

    event BetPlaced(uint256 indexed gameId, address indexed user, bool side, uint256 amount);
    event BettingLocked(uint256 indexed gameId, uint256 totalBetRed, uint256 totalBetBlue);
    event GameSettled(uint256 indexed gameId, bool redWins, uint256 totalPool, uint256 protocolFee);
    event RewardClaimed(uint256 indexed gameId, address indexed user, uint256 reward);

    modifier onlyOwner() {
        require(msg.sender == owner, "not owner");
        _;
    }

    modifier onlyArena() {
        require(msg.sender == arena, "not arena");
        _;
    }

    constructor(address _usdc, address _protocolTreasury) {
        owner = msg.sender;
        usdc = _usdc;
        protocolTreasury = _protocolTreasury;
    }

    function setArena(address _arena) external onlyOwner {
        arena = _arena;
    }

    /// @notice 初始化一场对局的下注
    function initGame(
        uint256 gameId,
        address agentRed,
        address agentBlue,
        uint256 bettingDeadline
    ) external onlyArena {
        games[gameId] = ArenaTypes.GameInfo({
            gameId: gameId,
            agentRed: agentRed,
            agentBlue: agentBlue,
            totalBetRed: 0,
            totalBetBlue: 0,
            bettingDeadline: bettingDeadline,
            status: ArenaTypes.GameStatus.Open,
            winner: ArenaTypes.Side.None
        });
    }

    /// @notice 下注（USDC 已由 Arena 转入）
    function placeBet(uint256 gameId, address user, bool side, uint256 amount) external onlyArena {
        ArenaTypes.GameInfo storage game = games[gameId];
        require(game.status == ArenaTypes.GameStatus.Open, "betting not open");

        // Update totals
        if (side) {
            game.totalBetRed += amount;
        } else {
            game.totalBetBlue += amount;
        }

        // Record bet (aggregate per user per game)
        BetInfo storage bet = bets[gameId][user];
        bet.side = side ? ArenaTypes.Side.Red : ArenaTypes.Side.Blue;
        bet.amount += amount;

        emit BetPlaced(gameId, user, side, amount);
    }

    /// @notice 锁定下注
    function lockBetting(uint256 gameId) external onlyArena {
        ArenaTypes.GameInfo storage game = games[gameId];
        require(game.status == ArenaTypes.GameStatus.Open, "not open");
        game.status = ArenaTypes.GameStatus.Locked;
        emit BettingLocked(gameId, game.totalBetRed, game.totalBetBlue);
    }

    /// @notice 结算
    function settle(uint256 gameId, bool redWins) external onlyArena nonReentrant {
        ArenaTypes.GameInfo storage game = games[gameId];
        require(game.status == ArenaTypes.GameStatus.Locked, "not locked");

        game.winner = redWins ? ArenaTypes.Side.Red : ArenaTypes.Side.Blue;
        game.status = ArenaTypes.GameStatus.Finished;

        uint256 totalPool = game.totalBetRed + game.totalBetBlue;
        uint256 fee = totalPool * protocolFeeBps / 10000;
        uint256 winnerPool = redWins ? game.totalBetRed : game.totalBetBlue;

        winnerPools[gameId] = winnerPool;

        // Transfer fee to treasury
        if (fee > 0 && protocolTreasury != address(0)) {
            IERC20(usdc).safeTransfer(protocolTreasury, fee);
        }

        emit GameSettled(gameId, redWins, totalPool, fee);
    }

    /// @notice 赢家提取奖金
    function claim(uint256 gameId) external nonReentrant {
        ArenaTypes.GameInfo storage game = games[gameId];
        require(game.status == ArenaTypes.GameStatus.Finished, "not finished");

        BetInfo storage bet = bets[gameId][msg.sender];
        require(bet.amount > 0, "no bet");
        require(bet.side == game.winner, "not winner");
        require(!bet.claimed, "already claimed");

        uint256 totalPool = game.totalBetRed + game.totalBetBlue;
        uint256 fee = totalPool * protocolFeeBps / 10000;
        uint256 distributable = totalPool - fee;
        uint256 reward = bet.amount * distributable / winnerPools[gameId];

        bet.claimed = true;
        IERC20(usdc).safeTransfer(msg.sender, reward);

        emit RewardClaimed(gameId, msg.sender, reward);
    }

    /// @notice 查询可提取金额
    function getReward(uint256 gameId, address user) external view returns (uint256) {
        ArenaTypes.GameInfo storage game = games[gameId];
        if (game.status != ArenaTypes.GameStatus.Finished) return 0;

        BetInfo storage bet = bets[gameId][user];
        if (bet.side != game.winner || bet.claimed) return 0;

        uint256 totalPool = game.totalBetRed + game.totalBetBlue;
        uint256 fee = totalPool * protocolFeeBps / 10000;
        uint256 distributable = totalPool - fee;
        return bet.amount * distributable / winnerPools[gameId];
    }

    /// @notice 查询赔率 (基于下注池)
    function getOdds(uint256 gameId) external view returns (uint256 oddsRed, uint256 oddsBlue) {
        ArenaTypes.GameInfo storage game = games[gameId];
        uint256 total = game.totalBetRed + game.totalBetBlue;
        if (total == 0) return (1 ether, 1 ether);
        // odds = total / side (scaled by 1e18)
        oddsRed = game.totalBetRed > 0 ? (total * 1e18 / game.totalBetRed) : type(uint256).max;
        oddsBlue = game.totalBetBlue > 0 ? (total * 1e18 / game.totalBetBlue) : type(uint256).max;
    }

    function getGame(uint256 gameId) external view returns (ArenaTypes.GameInfo memory) {
        return games[gameId];
    }
}
