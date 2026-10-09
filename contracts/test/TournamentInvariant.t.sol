// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Test} from "forge-std/Test.sol";
import {StdInvariant} from "forge-std/StdInvariant.sol";
import {TournamentManager} from "../src/TournamentManager.sol";
import {CourseToken} from "../src/CourseToken.sol";
import {CourseTypes} from "../src/CourseTypes.sol";

/// @dev Bounded adversarial ordering across several isolated tournaments, using real EVM calls.
contract TournamentHandler is Test {
    TournamentManager public manager;
    address private auditor = address(0x200);
    address private alice = address(0x301);
    address private bob = address(0x302);
    address private sponsor = address(0x401);
    uint256[] private ids;

    constructor(TournamentManager m) {
        manager = m;
    }

    function initialize() external {
        for (uint256 i; i < 3; ++i) {
            uint256 id = manager.createTournament(
                "Invariant event", "", auditor, uint64(block.timestamp + 300), uint64(block.timestamp + 3600), 60
            );
            ids.push(id);
            vm.prank(alice);
            manager.registerParticipant(id, keccak256("agent-a"));
            vm.prank(bob);
            manager.registerParticipant(id, keccak256("agent-b"));
            vm.prank(sponsor);
            manager.fundTournament(id, 100e6);
        }
    }

    function fund(uint256 seed, uint96 amount) external {
        uint256 id = ids[seed % ids.length];
        CourseTypes.TournamentInfo memory t = manager.getTournament(id);
        if (t.status != CourseTypes.Status.Open || block.timestamp >= t.registrationDeadline) return;
        vm.prank(sponsor);
        manager.fundTournament(id, bound(uint256(amount), 1, 1000e6));
    }

    function advanceTime(uint32 secondsForward) external {
        vm.warp(block.timestamp + bound(secondsForward, 1, 1000));
    }

    function start(uint256 seed) external {
        uint256 id = ids[seed % ids.length];
        CourseTypes.TournamentInfo memory t = manager.getTournament(id);
        if (
            t.status == CourseTypes.Status.Open && block.timestamp >= t.registrationDeadline
                && block.timestamp + t.reviewPeriod < t.settlementDeadline
        ) manager.startTournament(id);
    }

    function submit(uint256 seed) external {
        uint256 id = ids[seed % ids.length];
        CourseTypes.TournamentInfo memory t = manager.getTournament(id);
        if (t.status == CourseTypes.Status.Active && block.timestamp + t.reviewPeriod < t.settlementDeadline) {
            manager.submitResult(id, alice, keccak256("complete evidence"));
        }
    }

    function dispute(uint256 seed) external {
        uint256 id = ids[seed % ids.length];
        CourseTypes.TournamentInfo memory t = manager.getTournament(id);
        if (t.status == CourseTypes.Status.ResultSubmitted && block.timestamp < t.reviewDeadline) {
            vm.prank(bob);
            manager.disputeResult(id, keccak256("challenge"));
        }
    }

    function review(uint256 seed, bool reject) external {
        uint256 id = ids[seed % ids.length];
        CourseTypes.TournamentInfo memory t = manager.getTournament(id);
        if (block.timestamp < t.reviewDeadline || block.timestamp >= t.settlementDeadline) return;
        if (t.status == CourseTypes.Status.ResultSubmitted) {
            vm.prank(auditor);
            manager.approveResult(id);
        } else if (t.status == CourseTypes.Status.Disputed) {
            vm.prank(auditor);
            manager.resolveDispute(id, reject ? address(0) : bob, keccak256("decision"));
        }
    }

    function cancel(uint256 seed) external {
        uint256 id = ids[seed % ids.length];
        if (manager.getTournament(id).status == CourseTypes.Status.Open) {
            manager.cancelTournament(id, keccak256("cancel reason"));
        }
    }

    function expire(uint256 seed) external {
        uint256 id = ids[seed % ids.length];
        CourseTypes.TournamentInfo memory t = manager.getTournament(id);
        if (
            t.status != CourseTypes.Status.Finalized && t.status != CourseTypes.Status.Cancelled
                && block.timestamp >= t.settlementDeadline
        ) manager.expireTournament(id);
    }

    function withdraw(uint256 seed) external {
        uint256 id = ids[seed % ids.length];
        CourseTypes.TournamentInfo memory t = manager.getTournament(id);
        if (t.status == CourseTypes.Status.Finalized && manager.claimablePrize(id, t.winner) > 0) {
            vm.prank(t.winner);
            manager.claimPrize(id);
        } else if (t.status == CourseTypes.Status.Cancelled && manager.refundableAmount(id, sponsor) > 0) {
            vm.prank(sponsor);
            manager.claimRefund(id);
        }
    }
}

contract TournamentInvariantTest is StdInvariant, Test {
    CourseToken private token;
    TournamentManager private manager;

    function setUp() public {
        vm.warp(1_000_000);
        token = new CourseToken(address(this));
        manager = new TournamentManager(address(token), address(this));
        TournamentHandler handler = new TournamentHandler(manager);
        manager.grantRole(manager.ORGANIZER_ROLE(), address(handler));
        manager.grantRole(manager.AUDITOR_ROLE(), address(0x200));
        token.mint(address(0x401), 1_000_000e6);
        address escrowAddress = address(manager.prizeEscrow());
        vm.prank(address(0x401));
        token.approve(escrowAddress, type(uint256).max);
        handler.initialize();
        bytes4[] memory selectors = new bytes4[](9);
        selectors[0] = handler.fund.selector;
        selectors[1] = handler.advanceTime.selector;
        selectors[2] = handler.start.selector;
        selectors[3] = handler.submit.selector;
        selectors[4] = handler.dispute.selector;
        selectors[5] = handler.review.selector;
        selectors[6] = handler.cancel.selector;
        selectors[7] = handler.expire.selector;
        selectors[8] = handler.withdraw.selector;
        targetSelector(FuzzSelector({addr: address(handler), selectors: selectors}));
        targetContract(address(handler));
    }

    function invariantEscrowBalanceEqualsTotalAndPerTournamentLiabilities() public view {
        uint256 total;
        for (uint256 id = 1; id < manager.nextTournamentId(); ++id) {
            CourseTypes.EscrowInfo memory pool = manager.getEscrow(id);
            assertLe(pool.totalPaid, pool.totalFunded);
            total += pool.totalFunded - pool.totalPaid;
        }
        assertEq(total, manager.prizeEscrow().totalLiabilities());
        assertEq(total, token.balanceOf(address(manager.prizeEscrow())));
    }
}
