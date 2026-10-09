// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/// @title MockUSDC - 测试用 USDC（Sepolia 部署用）
contract MockUSDC is ERC20, Ownable {
    constructor() ERC20("USD Coin", "USDC") Ownable(msg.sender) {}

    /// @notice Legacy demonstration minting is restricted to the test-asset owner.
    function mint(address to, uint256 amount) external onlyOwner {
        require(amount > 0, "zero mint");
        _mint(to, amount);
    }

    /// @notice 销毁 USDC
    function burn(address from, uint256 amount) external {
        require(from == msg.sender && amount > 0, "not token holder");
        _burn(from, amount);
    }

    /// @notice USDC 使用 6 位小数
    function decimals() public pure override returns (uint8) {
        return 6;
    }
}
