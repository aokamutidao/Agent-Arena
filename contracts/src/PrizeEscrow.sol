// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "./CourseTypes.sol";

/// @notice Independent per-tournament ERC20 liabilities, with pull prize/refund payments.
contract PrizeEscrow is ReentrancyGuard {
    using SafeERC20 for IERC20;

    address public immutable manager;
    IERC20 public immutable token;
    uint256 public totalLiabilities;
    mapping(uint256 => bool) private initialized;
    mapping(uint256 => CourseTypes.EscrowInfo) private pools;
    mapping(uint256 => mapping(address => uint256)) private contributions;

    modifier onlyManager() {
        require(msg.sender == manager, "not manager");
        _;
    }

    modifier existingPool(uint256 tournamentId) {
        require(initialized[tournamentId], "unknown tournament");
        _;
    }

    constructor(address _token, address _manager) {
        require(_token.code.length > 0 && _manager != address(0), "invalid escrow config");
        token = IERC20(_token);
        manager = _manager;
    }

    function initialize(uint256 tournamentId) external onlyManager {
        require(tournamentId > 0 && !initialized[tournamentId], "invalid pool");
        initialized[tournamentId] = true;
    }

    function getEscrow(uint256 tournamentId)
        external
        view
        existingPool(tournamentId)
        returns (CourseTypes.EscrowInfo memory)
    {
        return pools[tournamentId];
    }

    function contributionOf(uint256 tournamentId, address sponsor)
        external
        view
        existingPool(tournamentId)
        returns (uint256)
    {
        return contributions[tournamentId][sponsor];
    }

    /// @notice Transfer directly from the approving sponsor into escrow, exactly once.
    function fund(uint256 tournamentId, address sponsor, uint256 amount)
        external
        onlyManager
        existingPool(tournamentId)
        nonReentrant
    {
        require(amount > 0 && sponsor != address(0), "invalid funding");
        CourseTypes.EscrowInfo storage pool = pools[tournamentId];
        require(!pool.cancelled && pool.winner == address(0), "pool closed");
        contributions[tournamentId][sponsor] += amount;
        pool.totalFunded += amount;
        totalLiabilities += amount;
        uint256 beforeBalance = token.balanceOf(address(this));
        token.safeTransferFrom(sponsor, address(this), amount);
        require(token.balanceOf(address(this)) == beforeBalance + amount, "unsupported token transfer");
    }

    function setWinner(uint256 tournamentId, address winner) external onlyManager existingPool(tournamentId) {
        CourseTypes.EscrowInfo storage pool = pools[tournamentId];
        require(winner != address(0) && !pool.cancelled && pool.winner == address(0), "pool closed");
        pool.winner = winner;
    }

    function cancel(uint256 tournamentId) external onlyManager existingPool(tournamentId) {
        CourseTypes.EscrowInfo storage pool = pools[tournamentId];
        require(!pool.cancelled && pool.winner == address(0), "pool closed");
        pool.cancelled = true;
    }

    function claimablePrize(uint256 tournamentId, address wallet)
        external
        view
        existingPool(tournamentId)
        returns (uint256)
    {
        CourseTypes.EscrowInfo storage pool = pools[tournamentId];
        return pool.winner != address(0) && pool.winner == wallet && !pool.prizeClaimed
            ? pool.totalFunded - pool.totalPaid
            : 0;
    }

    function refundableAmount(uint256 tournamentId, address wallet)
        external
        view
        existingPool(tournamentId)
        returns (uint256)
    {
        return pools[tournamentId].cancelled ? contributions[tournamentId][wallet] : 0;
    }

    function payPrize(uint256 tournamentId, address wallet)
        external
        onlyManager
        existingPool(tournamentId)
        nonReentrant
        returns (uint256 amount)
    {
        CourseTypes.EscrowInfo storage pool = pools[tournamentId];
        require(pool.winner != address(0) && pool.winner == wallet, "not winner");
        require(!pool.prizeClaimed, "prize already claimed");
        amount = pool.totalFunded - pool.totalPaid;
        require(amount > 0, "no prize");
        pool.prizeClaimed = true;
        pool.totalPaid += amount;
        totalLiabilities -= amount;
        token.safeTransfer(wallet, amount);
    }

    function refund(uint256 tournamentId, address wallet)
        external
        onlyManager
        existingPool(tournamentId)
        nonReentrant
        returns (uint256 amount)
    {
        require(pools[tournamentId].cancelled, "not cancelled");
        amount = contributions[tournamentId][wallet];
        require(amount > 0, "no refund");
        contributions[tournamentId][wallet] = 0;
        pools[tournamentId].totalPaid += amount;
        totalLiabilities -= amount;
        token.safeTransfer(wallet, amount);
    }
}
