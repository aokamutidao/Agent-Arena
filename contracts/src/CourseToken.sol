// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/// @notice Owner-minted test asset. DUSD is a demonstration token, not redeemable currency.
contract CourseToken is ERC20, Ownable {
    constructor(address initialOwner) ERC20("USD Demonstration", "DUSD") Ownable(initialOwner) {}

    function decimals() public pure override returns (uint8) {
        return 6;
    }

    function mint(address recipient, uint256 amount) external onlyOwner {
        require(amount > 0, "zero mint");
        _mint(recipient, amount);
    }
}
