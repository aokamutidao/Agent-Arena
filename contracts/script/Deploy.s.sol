// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "forge-std/Script.sol";
import "../src/AgentArena.sol";
import "../src/MockUSDC.sol";

/// @title Deploy - 部署 Agent Arena 到 Sepolia
contract Deploy is Script {
    function run() external {
        uint256 deployerKey = vm.envUint("PRIVATE_KEY");
        address deployer = vm.addr(deployerKey);

        console.log("Deployer:", deployer);

        vm.startBroadcast(deployerKey);

        // 1. 部署 MockUSDC
        MockUSDC usdc = new MockUSDC();
        console.log("MockUSDC:", address(usdc));

        // 给 deployer 铸造 10000 USDC（测试用）
        usdc.mint(deployer, 10_000 * 1e6);
        console.log("Minted 10000 USDC to deployer");

        // 2. 部署 AgentArena（内部部署 BettingPool + StrategyVoting + GameRegistry）
        AgentArena arena = new AgentArena(address(usdc), deployer);
        console.log("AgentArena:", address(arena));
        console.log("BettingPool:", address(arena.bettingPool()));
        console.log("StrategyVoting:", address(arena.strategyVoting()));
        console.log("GameRegistry:", address(arena.gameRegistry()));

        // 3. AgentArena is the legacy transferFrom spender.
        usdc.approve(address(arena), 10_000 * 1e6);
        console.log("Approved AgentArena for 10000 test USDC");

        vm.stopBroadcast();

        console.log("");
        console.log("=== Deployment Complete ===");
        console.log("MockUSDC:       ", address(usdc));
        console.log("AgentArena:     ", address(arena));
        console.log("BettingPool:    ", address(arena.bettingPool()));
        console.log("StrategyVoting: ", address(arena.strategyVoting()));
        console.log("GameRegistry:   ", address(arena.gameRegistry()));
    }
}
