// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/// @title ArenaCoin - Agent Arena 游戏币
contract ArenaCoin is ERC20, Ownable {
    // 水龙头：每个地址每日可领取 100 AC
    uint256 public constant DAILY_CLAIM_AMOUNT = 100 * 1e18;
    mapping(address => uint256) public lastClaimTime;

    event DailyClaimed(address indexed user, uint256 amount);

    constructor() ERC20("Arena Coin", "AC") Ownable(msg.sender) {
        // 初始铸造 100 万 AC 给部署者（用于测试）
        _mint(msg.sender, 1_000_000 * 1e18);
    }

    /// @notice 领取每日 AC（水龙头）
    function claimDaily() external {
        require(
            block.timestamp >= lastClaimTime[msg.sender] + 24 hours,
            "Can only claim once per 24 hours"
        );

        lastClaimTime[msg.sender] = block.timestamp;
        _mint(msg.sender, DAILY_CLAIM_AMOUNT);

        emit DailyClaimed(msg.sender, DAILY_CLAIM_AMOUNT);
    }

    /// @notice 查询是否可以领取
    function canClaim(address user) external view returns (bool) {
        return block.timestamp >= lastClaimTime[user] + 24 hours;
    }

    /// @notice 查询距离下次可领取的时间（秒）
    function timeUntilNextClaim(address user) external view returns (uint256) {
        uint256 nextClaim = lastClaimTime[user] + 24 hours;
        if (block.timestamp >= nextClaim) {
            return 0;
        }
        return nextClaim - block.timestamp;
    }

    /// @notice 管理员铸造（用于奖励）
    function mint(address to, uint256 amount) external onlyOwner {
        _mint(to, amount);
    }
}
