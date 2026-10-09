// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Test} from "forge-std/Test.sol";
import {AgentArena} from "../src/AgentArena.sol";
import {MockUSDC} from "../src/MockUSDC.sol";
import {BettingPool} from "../src/BettingPool.sol";
import {ArenaTypes} from "../src/ArenaTypes.sol";

/// @notice Regression tests for previously exploitable legacy source. Existing deployments are immutable.
contract LegacySecurityTest is Test {
    MockUSDC token;
    AgentArena arena;
    BettingPool pool;
    address alice = address(0x301);
    address bob = address(0x302);

    function setUp() public {
        token = new MockUSDC();
        arena = new AgentArena(address(token), address(0xFEE));
        pool = arena.bettingPool();
        arena.registerAgent(address(1), "Agent Red", "aggressive");
        arena.registerAgent(address(2), "Agent Blue", "defensive");
        token.mint(alice, 1000e6);
        token.mint(bob, 1000e6);
        vm.prank(alice);
        token.approve(address(arena), 1000e6);
        vm.prank(bob);
        token.approve(address(arena), 1000e6);
    }

    function _bet(address wallet, uint256 id, bool red, uint256 amount) internal {
        vm.prank(wallet);
        arena.betAndVote(id, red, amount, ArenaTypes.Strategy.Aggressive);
    }

    function testSideSwitchCannotDrainOtherGames() public {
        uint256 id = arena.createGame(address(1), address(2));
        uint256 other = arena.createGame(address(1), address(2));
        _bet(bob, other, true, 900e6);
        _bet(alice, id, false, 100e6);
        vm.prank(alice);
        vm.expectRevert("cannot switch side");
        arena.betAndVote(id, true, 1e6, ArenaTypes.Strategy.Aggressive);
        vm.prank(alice);
        vm.expectRevert("zero bet");
        arena.betAndVote(id, true, 0, ArenaTypes.Strategy.Aggressive);
        assertEq(token.balanceOf(address(pool)), 1000e6);
        assertEq(pool.getGame(id).totalBetRed, 0);
    }

    function testBettingDeadlineEnforcedAndEarlyLockRejected() public {
        uint256 id = arena.createGame(address(1), address(2));
        vm.expectRevert("betting still open");
        arena.startGame(id);
        vm.warp(arena.getGame(id).bettingDeadline);
        vm.prank(alice);
        vm.expectRevert("betting closed");
        arena.betAndVote(id, true, 1e6, ArenaTypes.Strategy.Aggressive);
        arena.startGame(id);
    }

    function testUncreatedGameCannotAcceptFunds() public {
        vm.prank(alice);
        vm.expectRevert("game not found");
        arena.betAndVote(9999, true, 10e6, ArenaTypes.Strategy.Aggressive);
        assertEq(token.balanceOf(alice), 1000e6);
        assertEq(token.balanceOf(address(pool)), 0);
    }

    function testVotesBelongToUserAndAggregate() public {
        uint256 id = arena.createGame(address(1), address(2));
        _bet(alice, id, true, 10e6);
        _bet(alice, id, true, 2e6);
        (, uint256 aliceAmount) = arena.strategyVoting().getUserVote(id, alice);
        (, uint256 arenaAmount) = arena.strategyVoting().getUserVote(id, address(arena));
        assertEq(aliceAmount, 12e6);
        assertEq(arenaAmount, 0);
    }

    function testAnyWalletCannotBurnEscrowOrMintTestAssets() public {
        uint256 id = arena.createGame(address(1), address(2));
        _bet(alice, id, true, 10e6);
        vm.prank(bob);
        vm.expectRevert("not token holder");
        token.burn(address(pool), 10e6);
        vm.prank(bob);
        vm.expectRevert();
        token.mint(bob, 10e6);
        assertEq(token.balanceOf(address(pool)), 10e6);
    }

    function testNoWinnerPoolRefundsInsteadOfDivisionByZero() public {
        uint256 id = arena.createGame(address(1), address(2));
        _bet(alice, id, false, 10e6);
        vm.warp(arena.getGame(id).bettingDeadline);
        arena.startGame(id);
        arena.finishGame(id, true, keccak256("legacy-evidence"));
        assertEq(uint256(arena.getGame(id).status), uint256(ArenaTypes.GameStatus.Cancelled));
        assertEq(pool.getReward(id, alice), 0);
        vm.prank(alice);
        pool.claimRefund(id);
        assertEq(token.balanceOf(alice), 1000e6);
        vm.prank(alice);
        vm.expectRevert("no refund");
        pool.claimRefund(id);
    }

    function testAbandonedGameHasPermissionlessTimeoutRefund() public {
        uint256 id = arena.createGame(address(1), address(2));
        _bet(alice, id, true, 10e6);
        vm.warp(arena.getGame(id).bettingDeadline + 1 days);
        vm.prank(bob);
        pool.cancelExpired(id);
        vm.prank(alice);
        pool.claimRefund(id);
        assertEq(token.balanceOf(alice), 1000e6);
    }

    function testDuplicateAgentCannotEraseStatistics() public {
        vm.expectRevert("invalid agent");
        arena.registerAgent(address(1), "Replacement", "aggressive");
        vm.expectRevert("invalid agents");
        arena.createGame(address(1), address(1));
    }
}
