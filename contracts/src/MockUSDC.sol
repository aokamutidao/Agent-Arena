// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/// @title MockUSDC - 测试用 USDC（Sepolia 部署用）
contract MockUSDC is ERC20, Ownable {
    constructor() ERC20("USD Coin", "USDC") Ownable(msg.sender) {}

    /// @notice 铸造 USDC（任何人可调用，测试用）
    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }

    /// @notice 销毁 USDC
    function burn(address from, uint256 amount) external {
        _burn(from, amount);
    }

    /// @notice USDC 使用 6 位小数
    function decimals() public pure override returns (uint8) {
        return 6;
    }
}
