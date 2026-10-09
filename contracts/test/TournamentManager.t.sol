// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {CourseToken} from "../src/CourseToken.sol";
import {CourseTypes} from "../src/CourseTypes.sol";
import {TournamentManager} from "../src/TournamentManager.sol";
import {PrizeEscrow} from "../src/PrizeEscrow.sol";
import {TournamentRegistry} from "../src/TournamentRegistry.sol";
import {ResultReview} from "../src/ResultReview.sol";

contract TournamentFixture is Test {
    CourseToken internal token;
    TournamentManager internal manager;
    PrizeEscrow internal escrow;
    address internal organizer = address(0x100);
    address internal auditor = address(0x200);
    address internal alice = address(0x301);
    address internal bob = address(0x302);
    address internal sponsor = address(0x401);
    address internal sponsor2 = address(0x402);
    address internal outsider = address(0x500);
    bytes32 internal constant HASH_A = keccak256("owner-bound-agent-alice");
    bytes32 internal constant HASH_B = keccak256("owner-bound-agent-bob");
    bytes32 internal constant RESULT = keccak256("canonical-complete-game-evidence");
    bytes32 internal constant DISPUTE = keccak256("canonical-dispute-evidence");
    bytes32 internal constant DECISION = keccak256("canonical-audit-decision");

    function setUp() public virtual {
        vm.warp(1_000_000);
        token = new CourseToken(address(this));
        manager = new TournamentManager(address(token), address(this));
        escrow = manager.prizeEscrow();
        manager.grantRole(manager.ORGANIZER_ROLE(), organizer);
        manager.grantRole(manager.AUDITOR_ROLE(), auditor);
        _giveTokens(sponsor, 10_000e6);
        _giveTokens(sponsor2, 10_000e6);
    }

    function _giveTokens(address wallet, uint256 amount) internal {
        token.mint(wallet, amount);
        vm.prank(wallet);
        token.approve(address(escrow), amount);
    }

    function _create() internal returns (uint256 id) {
        vm.prank(organizer);
        return manager.createTournament(
            "Two Agent competition",
            "evidence://rules-v1",
            auditor,
            uint64(block.timestamp + 300),
            uint64(block.timestamp + 3600),
            60
        );
    }

    function _join(uint256 id) internal {
        vm.prank(alice);
        manager.registerParticipant(id, HASH_A);
        vm.prank(bob);
        manager.registerParticipant(id, HASH_B);
    }

    function _fund(uint256 id, address wallet, uint256 amount) internal {
        vm.prank(wallet);
        manager.fundTournament(id, amount);
    }

    function _start(uint256 id) internal {
        vm.warp(manager.getTournament(id).registrationDeadline);
        vm.prank(organizer);
        manager.startTournament(id);
    }

    function _submitted() internal returns (uint256 id) {
        id = _create();
        _join(id);
        _fund(id, sponsor, 100e6);
        _fund(id, sponsor2, 50e6);
        _start(id);
        vm.prank(organizer);
        manager.submitResult(id, alice, RESULT);
    }

    function _approve(uint256 id) internal {
        vm.warp(manager.getTournament(id).reviewDeadline);
        vm.prank(auditor);
        manager.approveResult(id);
    }

    function _conserved() internal view {
        assertEq(token.balanceOf(address(escrow)), escrow.totalLiabilities());
    }
}

contract TournamentManagerTest is TournamentFixture {
    event TournamentFunded(uint256 indexed tournamentId, address indexed sponsor, uint256 amount, uint256 totalFunded);
    event SponsorRefunded(uint256 indexed tournamentId, address indexed sponsor, uint256 amount);

    function testFullMultiSponsorPrizeLifecycle() public {
        uint256 id = _submitted();
        CourseTypes.TournamentInfo memory t = manager.getTournament(id);
        assertEq(t.participant1, alice);
        assertEq(t.participant2, bob);
        assertEq(t.participant1AgentHash, HASH_A);
        assertEq(t.participant2AgentHash, HASH_B);
        assertEq(manager.getReview(id).resultHash, RESULT);
        assertEq(manager.claimablePrize(id, alice), 0);
        _approve(id);
        assertEq(manager.claimablePrize(id, alice), 150e6);
        assertEq(manager.claimablePrize(id, bob), 0);
        vm.prank(bob);
        vm.expectRevert("not winner");
        manager.claimPrize(id);
        vm.prank(alice);
        manager.claimPrize(id);
        assertEq(token.balanceOf(alice), 150e6);
        assertEq(manager.getEscrow(id).totalPaid, 150e6);
        assertEq(manager.claimablePrize(id, alice), 0);
        vm.prank(alice);
        vm.expectRevert("prize already claimed");
        manager.claimPrize(id);
        vm.prank(sponsor);
        vm.expectRevert("not cancelled");
        manager.claimRefund(id);
        _conserved();
    }

    function testDisputeCorrectsWinnerAfterFullWindow() public {
        uint256 id = _submitted();
        vm.prank(bob);
        manager.disputeResult(id, DISPUTE);
        vm.prank(auditor);
        vm.expectRevert("review still open");
        manager.resolveDispute(id, bob, DECISION);
        vm.prank(auditor);
        vm.expectRevert("not submitted");
        manager.approveResult(id);
        vm.warp(manager.getTournament(id).reviewDeadline);
        vm.prank(auditor);
        manager.resolveDispute(id, bob, DECISION);
        assertEq(manager.getTournament(id).winner, bob);
        assertEq(manager.getReview(id).challenger, bob);
        assertEq(manager.getReview(id).disputeHash, DISPUTE);
        assertEq(manager.getReview(id).decisionHash, DECISION);
        vm.prank(bob);
        manager.claimPrize(id);
        assertEq(token.balanceOf(bob), 150e6);
        _conserved();
    }

    function testAuditRejectsDisputedResultAndRefundsExactly() public {
        uint256 id = _submitted();
        vm.prank(sponsor2);
        manager.disputeResult(id, DISPUTE);
        vm.warp(manager.getTournament(id).reviewDeadline);
        vm.prank(auditor);
        manager.resolveDispute(id, address(0), DECISION);
        assertEq(uint256(manager.getTournament(id).status), uint256(CourseTypes.Status.Cancelled));
        assertEq(manager.refundableAmount(id, sponsor), 100e6);
        vm.expectEmit(true, true, false, true, address(manager));
        emit SponsorRefunded(id, sponsor, 100e6);
        vm.prank(sponsor);
        manager.claimRefund(id);
        vm.prank(sponsor2);
        manager.claimRefund(id);
        assertEq(token.balanceOf(sponsor), 10_000e6);
        assertEq(token.balanceOf(sponsor2), 10_000e6);
        assertEq(manager.contributionOf(id, sponsor), 0);
        assertEq(manager.getEscrow(id).totalPaid, 150e6);
        vm.prank(sponsor);
        vm.expectRevert("no refund");
        manager.claimRefund(id);
        vm.prank(outsider);
        vm.expectRevert("no refund");
        manager.claimRefund(id);
        _conserved();
    }

    function testWithdrawReplacesWalletAndHashInSameSlot() public {
        uint256 id = _create();
        _join(id);
        vm.prank(alice);
        manager.withdrawParticipant(id);
        assertEq(manager.getTournament(id).participant1, address(0));
        assertEq(manager.getTournament(id).participant1AgentHash, bytes32(0));
        vm.prank(outsider);
        manager.registerParticipant(id, keccak256("outsider-profile"));
        assertEq(manager.getTournament(id).participant1, outsider);
        assertEq(manager.getTournament(id).participant1AgentHash, keccak256("outsider-profile"));
        assertEq(manager.getTournament(id).participant2AgentHash, HASH_B);
        vm.warp(manager.getTournament(id).registrationDeadline);
        vm.prank(outsider);
        vm.expectRevert("registration closed");
        manager.withdrawParticipant(id);
    }

    function testFundingEmittedAtManagerAndSumsContributions() public {
        uint256 id = _create();
        vm.expectEmit(true, true, false, true, address(manager));
        emit TournamentFunded(id, sponsor, 10e6, 10e6);
        _fund(id, sponsor, 10e6);
        _fund(id, sponsor, 2e6);
        _fund(id, sponsor2, 7e6);
        assertEq(manager.getEscrow(id).totalFunded, 19e6);
        assertEq(manager.contributionOf(id, sponsor), 12e6);
        assertEq(manager.contributionOf(id, sponsor2), 7e6);
        _conserved();
    }

    function testRegistrationFundingAndStartBoundary() public {
        uint256 id = _create();
        _join(id);
        _fund(id, sponsor, 1e6);
        vm.prank(organizer);
        vm.expectRevert("registration still open");
        manager.startTournament(id);
        vm.warp(manager.getTournament(id).registrationDeadline);
        vm.prank(outsider);
        vm.expectRevert("registration closed");
        manager.registerParticipant(id, RESULT);
        vm.prank(sponsor);
        vm.expectRevert("registration closed");
        manager.fundTournament(id, 1e6);
        vm.prank(organizer);
        manager.startTournament(id);
        vm.prank(organizer);
        vm.expectRevert("not open");
        manager.startTournament(id);
    }

    function testReviewCannotBeShortenedAndDisputeDeadlineExclusive() public {
        uint256 id = _submitted();
        uint64 end = manager.getTournament(id).reviewDeadline;
        vm.warp(end - 1);
        vm.prank(auditor);
        vm.expectRevert("review still open");
        manager.approveResult(id);
        vm.warp(end);
        vm.prank(bob);
        vm.expectRevert("review window closed");
        manager.disputeResult(id, DISPUTE);
        vm.prank(auditor);
        manager.approveResult(id);
    }

    function testLastReviewSecondStillAllowsDispute() public {
        uint256 id = _submitted();
        vm.warp(manager.getTournament(id).reviewDeadline - 1);
        vm.prank(bob);
        manager.disputeResult(id, DISPUTE);
        vm.warp(manager.getTournament(id).reviewDeadline);
        vm.prank(auditor);
        manager.resolveDispute(id, alice, DECISION);
    }

    function testLateSubmissionCannotCutReviewShort() public {
        uint256 id = _create();
        _join(id);
        _fund(id, sponsor, 1e6);
        _start(id);
        vm.warp(manager.getTournament(id).settlementDeadline - 60);
        vm.prank(organizer);
        vm.expectRevert("insufficient review time");
        manager.submitResult(id, alice, RESULT);
    }

    function testExpiryRefundsActiveTournamentDespiteRevokedRoles() public {
        uint256 id = _create();
        _join(id);
        _fund(id, sponsor, 100e6);
        _start(id);
        manager.revokeRole(manager.ORGANIZER_ROLE(), organizer);
        manager.revokeRole(manager.AUDITOR_ROLE(), auditor);
        vm.warp(manager.getTournament(id).settlementDeadline);
        vm.prank(outsider);
        manager.expireTournament(id);
        vm.prank(sponsor);
        manager.claimRefund(id);
        assertEq(token.balanceOf(sponsor), 10_000e6);
        _conserved();
    }

    function testExpiryRefundsSubmittedTournamentAndBlocksApproval() public {
        uint256 id = _submitted();
        vm.warp(manager.getTournament(id).settlementDeadline);
        vm.prank(auditor);
        vm.expectRevert("settlement expired");
        manager.approveResult(id);
        vm.prank(outsider);
        manager.expireTournament(id);
        vm.prank(sponsor);
        manager.claimRefund(id);
        vm.prank(sponsor2);
        manager.claimRefund(id);
        _conserved();
    }

    function testExpiryRefundsUnresolvedDispute() public {
        uint256 id = _submitted();
        vm.prank(alice);
        manager.disputeResult(id, DISPUTE);
        vm.warp(manager.getTournament(id).settlementDeadline);
        vm.prank(auditor);
        vm.expectRevert("settlement expired");
        manager.resolveDispute(id, bob, DECISION);
        vm.prank(outsider);
        manager.expireTournament(id);
        vm.prank(sponsor);
        manager.claimRefund(id);
        vm.prank(sponsor2);
        manager.claimRefund(id);
        _conserved();
    }

    function testUnderfilledEventExpiresAtRegistrationDeadline() public {
        uint256 id = _create();
        vm.prank(alice);
        manager.registerParticipant(id, HASH_A);
        _fund(id, sponsor, 100e6);
        vm.prank(outsider);
        vm.expectRevert("not expired");
        manager.expireTournament(id);
        vm.warp(manager.getTournament(id).registrationDeadline);
        vm.prank(outsider);
        manager.expireTournament(id);
        vm.prank(sponsor);
        manager.claimRefund(id);
        _conserved();
    }

    function testFullOpenEventCannotExpireUntilSettlementDeadline() public {
        uint256 id = _create();
        _join(id);
        _fund(id, sponsor, 1e6);
        vm.warp(manager.getTournament(id).registrationDeadline);
        vm.prank(outsider);
        vm.expectRevert("not expired");
        manager.expireTournament(id);
        vm.warp(manager.getTournament(id).settlementDeadline);
        vm.prank(outsider);
        manager.expireTournament(id);
    }

    function testOrganizerCancelsOpenEventWithPullRefunds() public {
        uint256 id = _create();
        _join(id);
        _fund(id, sponsor, 20e6);
        vm.prank(organizer);
        manager.cancelTournament(id, DECISION);
        vm.prank(sponsor);
        manager.claimRefund(id);
        vm.prank(organizer);
        vm.expectRevert("already started");
        manager.cancelTournament(id, DECISION);
        vm.expectRevert("terminal tournament");
        manager.expireTournament(id);
        _conserved();
    }

    function testOrganizerCannotCancelAfterStartOrChangeResult() public {
        uint256 id = _submitted();
        vm.prank(organizer);
        vm.expectRevert("already started");
        manager.cancelTournament(id, DECISION);
        vm.prank(organizer);
        vm.expectRevert("not active");
        manager.submitResult(id, bob, DISPUTE);
        _approve(id);
        vm.prank(outsider);
        vm.expectRevert("terminal tournament");
        manager.expireTournament(id);
        vm.prank(auditor);
        vm.expectRevert("not submitted");
        manager.approveResult(id);
    }

    function testRolesCannotReplaceSpecificOrganizerOrAuditor() public {
        address otherOrganizer = address(0x101);
        address otherAuditor = address(0x201);
        manager.grantRole(manager.ORGANIZER_ROLE(), otherOrganizer);
        manager.grantRole(manager.AUDITOR_ROLE(), otherAuditor);
        uint256 id = _submitted();
        vm.prank(otherOrganizer);
        vm.expectRevert("not tournament organizer");
        manager.submitResult(id, alice, RESULT);
        vm.prank(otherAuditor);
        vm.expectRevert("not designated auditor");
        manager.approveResult(id);
        vm.expectRevert("not designated auditor");
        manager.approveResult(id);
        bytes32 auditorRole = manager.AUDITOR_ROLE();
        vm.prank(outsider);
        vm.expectRevert();
        manager.grantRole(auditorRole, outsider);
    }

    function testRoleRevocationBlocksReviewButTimeoutRefundsRemain() public {
        uint256 id = _submitted();
        manager.revokeRole(manager.AUDITOR_ROLE(), auditor);
        vm.warp(manager.getTournament(id).reviewDeadline);
        vm.prank(auditor);
        vm.expectRevert();
        manager.approveResult(id);
        vm.warp(manager.getTournament(id).settlementDeadline);
        vm.prank(sponsor);
        manager.expireTournament(id);
        vm.prank(sponsor);
        manager.claimRefund(id);
    }

    function testOnlyParticipantsAndSponsorsMayDispute() public {
        uint256 id = _submitted();
        vm.prank(outsider);
        vm.expectRevert("not stakeholder");
        manager.disputeResult(id, DISPUTE);
        vm.prank(auditor);
        vm.expectRevert("not stakeholder");
        manager.disputeResult(id, DISPUTE);
        vm.prank(sponsor);
        manager.disputeResult(id, DISPUTE);
        vm.prank(bob);
        vm.expectRevert("not submitted");
        manager.disputeResult(id, DISPUTE);
    }

    function testAuditorAndOrganizerCannotCompeteOrAuditorSponsor() public {
        uint256 id = _create();
        vm.prank(auditor);
        vm.expectRevert("conflicted participant");
        manager.registerParticipant(id, HASH_A);
        vm.prank(organizer);
        vm.expectRevert("conflicted participant");
        manager.registerParticipant(id, HASH_A);
        vm.prank(auditor);
        vm.expectRevert("conflicted sponsor");
        manager.fundTournament(id, 1e6);
    }

    function testZeroHashDuplicateAndOverCapacityRejected() public {
        uint256 id = _create();
        vm.prank(alice);
        vm.expectRevert("zero agent hash");
        manager.registerParticipant(id, bytes32(0));
        _join(id);
        vm.prank(alice);
        vm.expectRevert("already registered");
        manager.registerParticipant(id, HASH_B);
        vm.prank(outsider);
        vm.expectRevert("tournament full");
        manager.registerParticipant(id, RESULT);
        vm.prank(outsider);
        vm.expectRevert("not participant");
        manager.withdrawParticipant(id);
        vm.prank(sponsor);
        vm.expectRevert("zero funding");
        manager.fundTournament(id, 0);
    }

    function testInvalidWinnerAndEvidenceHashesRejected() public {
        uint256 id = _create();
        _join(id);
        _fund(id, sponsor, 1e6);
        _start(id);
        vm.prank(organizer);
        vm.expectRevert("winner not participant");
        manager.submitResult(id, outsider, RESULT);
        vm.prank(organizer);
        vm.expectRevert("zero result hash");
        manager.submitResult(id, alice, bytes32(0));
        vm.prank(organizer);
        manager.submitResult(id, alice, RESULT);
        vm.prank(bob);
        vm.expectRevert("zero dispute hash");
        manager.disputeResult(id, bytes32(0));
        vm.prank(bob);
        manager.disputeResult(id, DISPUTE);
        vm.warp(manager.getTournament(id).reviewDeadline);
        vm.prank(auditor);
        vm.expectRevert("zero decision hash");
        manager.resolveDispute(id, alice, bytes32(0));
        vm.prank(auditor);
        vm.expectRevert("winner not participant");
        manager.resolveDispute(id, outsider, DECISION);
    }

    function testUnfundedOrUnderfilledTournamentCannotStart() public {
        uint256 id = _create();
        _join(id);
        vm.warp(manager.getTournament(id).registrationDeadline);
        vm.prank(organizer);
        vm.expectRevert("unfunded tournament");
        manager.startTournament(id);
        uint256 second = _create();
        _fund(second, sponsor, 1e6);
        vm.warp(manager.getTournament(second).registrationDeadline);
        vm.prank(organizer);
        vm.expectRevert("need two participants");
        manager.startTournament(second);
    }

    function testApprovalToWrongSpenderRevertsAndLeavesAccountingUntouched() public {
        uint256 id = _create();
        token.mint(outsider, 1e6);
        vm.prank(outsider);
        token.approve(address(manager), 1e6);
        vm.prank(outsider);
        vm.expectRevert();
        manager.fundTournament(id, 1e6);
        assertEq(manager.getEscrow(id).totalFunded, 0);
        assertEq(manager.contributionOf(id, outsider), 0);
        _conserved();
    }

    function testUnknownIdsRejectedForEveryFinancialPath() public {
        vm.expectRevert("unknown tournament");
        manager.getTournament(0);
        vm.expectRevert("unknown tournament");
        manager.getEscrow(1);
        vm.expectRevert("unknown tournament");
        manager.getReview(999);
        vm.expectRevert("unknown tournament");
        manager.contributionOf(999, sponsor);
        vm.expectRevert("unknown tournament");
        manager.claimablePrize(999, alice);
        vm.expectRevert("unknown tournament");
        manager.refundableAmount(999, sponsor);
        vm.prank(sponsor);
        vm.expectRevert("unknown tournament");
        manager.fundTournament(999, 1e6);
        vm.expectRevert("unknown tournament");
        manager.expireTournament(999);
        vm.expectRevert("unknown tournament");
        manager.claimPrize(999);
        vm.expectRevert("unknown tournament");
        manager.claimRefund(999);
    }

    function testDirectChildMutationsRejected() public {
        uint256 id = _create();
        TournamentRegistry registry = manager.tournamentRegistry();
        ResultReview review = manager.resultReview();
        vm.expectRevert("not manager");
        escrow.setWinner(id, outsider);
        vm.expectRevert("not manager");
        escrow.fund(id, sponsor, 1e6);
        vm.expectRevert("not manager");
        escrow.refund(id, sponsor);
        vm.expectRevert("not manager");
        registry.cancel(id);
        vm.expectRevert("not manager");
        review.recordDecision(id, DECISION);
    }

    function testCreateValidatesRolesTimeAndMetadata() public {
        vm.prank(outsider);
        vm.expectRevert();
        manager.createTournament("event", "", auditor, 1001000, 1003600, 60);
        vm.startPrank(organizer);
        vm.expectRevert("invalid auditor");
        manager.createTournament("event", "", organizer, 1001000, 1003600, 60);
        vm.expectRevert("invalid auditor");
        manager.createTournament("event", "", outsider, 1001000, 1003600, 60);
        vm.expectRevert("invalid title");
        manager.createTournament("", "", auditor, 1001000, 1003600, 60);
        vm.expectRevert("invalid review period");
        manager.createTournament("event", "", auditor, 1001000, 1003600, 1);
        vm.expectRevert("invalid registration deadline");
        manager.createTournament("event", "", auditor, 1000000, 1003600, 60);
        vm.expectRevert("invalid settlement deadline");
        manager.createTournament("event", "", auditor, 1001000, 1001060, 60);
        vm.stopPrank();
    }

    function testDemonstrationTokenCannotBeMintedByOrdinaryWallet() public {
        vm.prank(outsider);
        vm.expectRevert();
        token.mint(outsider, 1e6);
        vm.expectRevert("zero mint");
        token.mint(outsider, 0);
        assertEq(token.decimals(), 6);
    }

    function testTwoTournamentsNeverSubsidizeEachOther() public {
        uint256 first = _create();
        _join(first);
        _fund(first, sponsor, 10e6);
        uint256 second = _create();
        _join(second);
        _fund(second, sponsor2, 900e6);
        _start(first);
        vm.prank(organizer);
        manager.submitResult(first, alice, RESULT);
        _approve(first);
        vm.prank(alice);
        manager.claimPrize(first);
        assertEq(token.balanceOf(alice), 10e6);
        assertEq(token.balanceOf(address(escrow)), 900e6);
        assertEq(manager.getEscrow(second).totalFunded, 900e6);
        vm.prank(alice);
        vm.expectRevert("prize already claimed");
        manager.claimPrize(first);
        vm.prank(organizer);
        manager.cancelTournament(second, DECISION);
        vm.prank(sponsor2);
        manager.claimRefund(second);
        _conserved();
    }

    function testFuzzMultipleSponsorRefundConservation(uint96 a, uint96 b) public {
        uint256 amountA = bound(uint256(a), 1, 1000e6);
        uint256 amountB = bound(uint256(b), 1, 1000e6);
        uint256 id = _create();
        _fund(id, sponsor, amountA);
        _fund(id, sponsor2, amountB);
        assertEq(manager.getEscrow(id).totalFunded, amountA + amountB);
        _conserved();
        vm.prank(organizer);
        manager.cancelTournament(id, DECISION);
        vm.prank(sponsor2);
        manager.claimRefund(id);
        assertEq(escrow.totalLiabilities(), amountA);
        _conserved();
        vm.prank(sponsor);
        manager.claimRefund(id);
        assertEq(token.balanceOf(sponsor), 10_000e6);
        assertEq(token.balanceOf(sponsor2), 10_000e6);
        _conserved();
    }

    function testFuzzWinnerReceivesExactlyBothContributions(uint96 a, uint96 b) public {
        uint256 amountA = bound(uint256(a), 1, 1000e6);
        uint256 amountB = bound(uint256(b), 1, 1000e6);
        uint256 id = _create();
        _join(id);
        _fund(id, sponsor, amountA);
        _fund(id, sponsor2, amountB);
        _start(id);
        vm.prank(organizer);
        manager.submitResult(id, alice, RESULT);
        _approve(id);
        vm.prank(alice);
        manager.claimPrize(id);
        assertEq(token.balanceOf(alice), amountA + amountB);
        assertEq(manager.getEscrow(id).totalPaid, manager.getEscrow(id).totalFunded);
        _conserved();
    }
}

/// @dev A token callback attempts a second winner withdrawal during the first transfer.
contract CallbackToken is ERC20 {
    TournamentManager public manager;
    uint256 public tournamentId;
    address public escrow;
    bool public attack;
    bool public reentrySucceeded;
    bool public rejectPayout;
    constructor() ERC20("Callback test", "CBT") {}

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }

    function configure(TournamentManager m, uint256 id) external {
        manager = m;
        tournamentId = id;
        escrow = address(m.prizeEscrow());
    }

    function setAttack(bool value) external {
        attack = value;
    }

    function setRejectPayout(bool value) external {
        rejectPayout = value;
    }

    function _update(address from, address to, uint256 amount) internal override {
        if (from == escrow && rejectPayout) revert("token payout rejected");
        super._update(from, to, amount);
        if (attack && from == escrow) {
            attack = false;
            (reentrySucceeded,) = address(manager).call(abi.encodeCall(TournamentManager.claimPrize, (tournamentId)));
        }
    }
}

contract FeeToken is ERC20 {
    constructor() ERC20("Fee test", "FEE") {}

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }

    function _update(address from, address to, uint256 amount) internal override {
        if (from != address(0) && to != address(0)) {
            super._update(from, to, amount * 9 / 10);
            super._update(from, address(0), amount - amount * 9 / 10);
        } else {
            super._update(from, to, amount);
        }
    }
}

contract TournamentAdversarialTest is Test {
    function testReentrantPrizeWithdrawalBlocked() public {
        CallbackToken token = new CallbackToken();
        TournamentManager manager = new TournamentManager(address(token), address(this));
        address auditor = address(0x200);
        address bob = address(0x302);
        manager.grantRole(manager.ORGANIZER_ROLE(), address(this));
        manager.grantRole(manager.AUDITOR_ROLE(), auditor);
        uint256 id = manager.createTournament(
            "Callback test", "", auditor, uint64(block.timestamp + 300), uint64(block.timestamp + 3600), 60
        );
        vm.prank(address(token));
        manager.registerParticipant(id, keccak256("token-agent"));
        vm.prank(bob);
        manager.registerParticipant(id, keccak256("bob-agent"));
        token.mint(address(this), 100e6);
        token.approve(address(manager.prizeEscrow()), 100e6);
        manager.fundTournament(id, 100e6);
        vm.warp(manager.getTournament(id).registrationDeadline);
        manager.startTournament(id);
        manager.submitResult(id, address(token), keccak256("result"));
        vm.warp(manager.getTournament(id).reviewDeadline);
        vm.prank(auditor);
        manager.approveResult(id);
        token.configure(manager, id);
        token.setAttack(true);
        vm.prank(address(token));
        manager.claimPrize(id);
        assertFalse(token.reentrySucceeded());
        assertEq(token.balanceOf(address(token)), 100e6);
        assertEq(manager.prizeEscrow().totalLiabilities(), 0);
    }

    function testPayoutFailureRollsBackClaimAndCanBeRetried() public {
        CallbackToken token = new CallbackToken();
        TournamentManager manager = new TournamentManager(address(token), address(this));
        address alice = address(0x301);
        address bob = address(0x302);
        address auditor = address(0x200);
        manager.grantRole(manager.ORGANIZER_ROLE(), address(this));
        manager.grantRole(manager.AUDITOR_ROLE(), auditor);
        uint256 id = manager.createTournament(
            "Payout failure", "", auditor, uint64(block.timestamp + 300), uint64(block.timestamp + 3600), 60
        );
        vm.prank(alice);
        manager.registerParticipant(id, keccak256("a"));
        vm.prank(bob);
        manager.registerParticipant(id, keccak256("b"));
        token.mint(address(this), 100e6);
        token.approve(address(manager.prizeEscrow()), 100e6);
        manager.fundTournament(id, 100e6);
        vm.warp(manager.getTournament(id).registrationDeadline);
        manager.startTournament(id);
        manager.submitResult(id, alice, keccak256("result"));
        vm.warp(manager.getTournament(id).reviewDeadline);
        vm.prank(auditor);
        manager.approveResult(id);
        token.configure(manager, id);
        token.setRejectPayout(true);
        vm.prank(alice);
        vm.expectRevert("token payout rejected");
        manager.claimPrize(id);
        assertFalse(manager.getEscrow(id).prizeClaimed);
        assertEq(manager.prizeEscrow().totalLiabilities(), 100e6);
        token.setRejectPayout(false);
        vm.prank(alice);
        manager.claimPrize(id);
        assertEq(token.balanceOf(alice), 100e6);
    }

    function testFeeOnTransferCannotCreateUnbackedLiabilities() public {
        FeeToken token = new FeeToken();
        TournamentManager manager = new TournamentManager(address(token), address(this));
        manager.grantRole(manager.ORGANIZER_ROLE(), address(this));
        address auditor = address(0x200);
        manager.grantRole(manager.AUDITOR_ROLE(), auditor);
        uint256 id = manager.createTournament(
            "Token compatibility", "", auditor, uint64(block.timestamp + 300), uint64(block.timestamp + 3600), 60
        );
        token.mint(address(this), 100e6);
        token.approve(address(manager.prizeEscrow()), 100e6);
        vm.expectRevert("unsupported token transfer");
        manager.fundTournament(id, 100e6);
        assertEq(token.balanceOf(address(this)), 100e6);
        assertEq(manager.getEscrow(id).totalFunded, 0);
        assertEq(manager.prizeEscrow().totalLiabilities(), 0);
    }
}

contract TournamentGasTest is TournamentFixture {
    function _fundMany(uint256 id, uint256 sponsors) internal {
        for (uint256 i; i < sponsors; ++i) {
            address wallet = address(uint160(0x1000 + i));
            _giveTokens(wallet, 1e6);
            _fund(id, wallet, 1e6);
        }
    }

    function _claimGas(uint256 sponsors, string memory label) internal {
        uint256 id = _create();
        _join(id);
        _fundMany(id, sponsors);
        _start(id);
        vm.prank(organizer);
        manager.submitResult(id, alice, RESULT);
        _approve(id);
        vm.prank(alice);
        manager.claimPrize(id);
        uint256 measured = vm.snapshotGasLastCall("PrizeScale", label);
        emit log_named_uint(label, measured);
        assertEq(token.balanceOf(alice), sponsors * 1e6);
        _conserved();
    }

    function _refundGas(uint256 sponsors, string memory label) internal {
        uint256 id = _create();
        _fundMany(id, sponsors);
        vm.prank(organizer);
        manager.cancelTournament(id, DECISION);
        vm.prank(address(0x1000));
        manager.claimRefund(id);
        uint256 measured = vm.snapshotGasLastCall("RefundScale", label);
        emit log_named_uint(label, measured);
        assertEq(token.balanceOf(address(0x1000)), 1e6);
        _conserved();
    }

    function testGasClaimOneSponsor() public {
        _claimGas(1, "claim_1_sponsor");
    }

    function testGasClaimSixteenSponsors() public {
        _claimGas(16, "claim_16_sponsors");
    }

    function testGasClaimSixtyFourSponsors() public {
        _claimGas(64, "claim_64_sponsors");
    }

    function testGasRefundOneSponsor() public {
        _refundGas(1, "refund_1_sponsor");
    }

    function testGasRefundSixteenSponsors() public {
        _refundGas(16, "refund_16_sponsors");
    }

    function testGasRefundSixtyFourSponsors() public {
        _refundGas(64, "refund_64_sponsors");
    }
}
