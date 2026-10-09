// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "./CourseTypes.sol";
import "./TournamentRegistry.sol";
import "./PrizeEscrow.sol";
import "./ResultReview.sol";

/// @title Audited tournament prize escrow (course prototype)
/// @notice Wallets sign every business transaction. Admin cannot spend or redirect escrow.
contract TournamentManager is AccessControl, ReentrancyGuard {
    bytes32 public constant ORGANIZER_ROLE = keccak256("ORGANIZER_ROLE");
    bytes32 public constant AUDITOR_ROLE = keccak256("AUDITOR_ROLE");
    address public immutable settlementToken;
    TournamentRegistry public immutable tournamentRegistry;
    PrizeEscrow public immutable prizeEscrow;
    ResultReview public immutable resultReview;

    event TournamentCreated(
        uint256 indexed tournamentId,
        address indexed organizer,
        address indexed auditor,
        string title,
        string metadataURI,
        uint64 registrationDeadline,
        uint64 settlementDeadline,
        uint64 reviewPeriod
    );
    event ParticipantRegistered(
        uint256 indexed tournamentId, address indexed participant, bytes32 agentHash, uint8 slot
    );
    event ParticipantWithdrawn(uint256 indexed tournamentId, address indexed participant);
    event TournamentFunded(uint256 indexed tournamentId, address indexed sponsor, uint256 amount, uint256 totalFunded);
    event TournamentStarted(uint256 indexed tournamentId);
    event ResultSubmitted(
        uint256 indexed tournamentId,
        address indexed submitter,
        address indexed proposedWinner,
        bytes32 resultHash,
        uint64 reviewDeadline
    );
    event ResultDisputed(uint256 indexed tournamentId, address indexed challenger, bytes32 disputeHash);
    event ResultApproved(
        uint256 indexed tournamentId, address indexed auditor, address indexed winner, bytes32 resultHash
    );
    event DisputeResolved(
        uint256 indexed tournamentId, address indexed auditor, address indexed winner, bytes32 decisionHash
    );
    event TournamentCancelled(uint256 indexed tournamentId, address indexed caller, bytes32 reasonHash, bool expired);
    event PrizeClaimed(uint256 indexed tournamentId, address indexed winner, uint256 amount);
    event SponsorRefunded(uint256 indexed tournamentId, address indexed sponsor, uint256 amount);

    constructor(address token, address admin) {
        require(token.code.length > 0 && admin != address(0), "invalid manager config");
        settlementToken = token;
        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        tournamentRegistry = new TournamentRegistry(address(this));
        prizeEscrow = new PrizeEscrow(token, address(this));
        resultReview = new ResultReview(address(this));
    }

    function nextTournamentId() external view returns (uint256) {
        return tournamentRegistry.nextTournamentId();
    }

    function getTournament(uint256 tournamentId) public view returns (CourseTypes.TournamentInfo memory) {
        return tournamentRegistry.getTournament(tournamentId);
    }

    function getReview(uint256 tournamentId) external view returns (CourseTypes.ReviewInfo memory) {
        getTournament(tournamentId);
        return resultReview.getReview(tournamentId);
    }

    function getEscrow(uint256 tournamentId) external view returns (CourseTypes.EscrowInfo memory) {
        return prizeEscrow.getEscrow(tournamentId);
    }

    function contributionOf(uint256 tournamentId, address sponsor) public view returns (uint256) {
        return prizeEscrow.contributionOf(tournamentId, sponsor);
    }

    function claimablePrize(uint256 tournamentId, address wallet) external view returns (uint256) {
        return prizeEscrow.claimablePrize(tournamentId, wallet);
    }

    function refundableAmount(uint256 tournamentId, address wallet) external view returns (uint256) {
        return prizeEscrow.refundableAmount(tournamentId, wallet);
    }

    /// @notice Create a fixed two-Agent competition with an immutable independent auditor.
    function createTournament(
        string calldata title,
        string calldata metadataURI,
        address auditor,
        uint64 registrationDeadline,
        uint64 settlementDeadline,
        uint64 reviewPeriod
    ) external onlyRole(ORGANIZER_ROLE) nonReentrant returns (uint256 tournamentId) {
        require(bytes(title).length > 0 && bytes(title).length <= 120, "invalid title");
        require(bytes(metadataURI).length <= 512, "metadata too long");
        require(auditor != address(0) && auditor != msg.sender && hasRole(AUDITOR_ROLE, auditor), "invalid auditor");
        require(reviewPeriod >= 60 && reviewPeriod <= 30 days, "invalid review period");
        require(registrationDeadline > block.timestamp, "invalid registration deadline");
        require(
            uint256(settlementDeadline) > uint256(registrationDeadline) + reviewPeriod, "invalid settlement deadline"
        );
        tournamentId = tournamentRegistry.createTournament(
            msg.sender, auditor, title, metadataURI, registrationDeadline, settlementDeadline, reviewPeriod
        );
        prizeEscrow.initialize(tournamentId);
        emit TournamentCreated(
            tournamentId,
            msg.sender,
            auditor,
            title,
            metadataURI,
            registrationDeadline,
            settlementDeadline,
            reviewPeriod
        );
    }

    /// @notice Register caller's immutable, owner-bound canonical Agent profile hash.
    function registerParticipant(uint256 tournamentId, bytes32 agentHash) external nonReentrant {
        CourseTypes.TournamentInfo memory t = getTournament(tournamentId);
        _requireRegistrationOpen(t);
        require(agentHash != bytes32(0), "zero agent hash");
        require(msg.sender != t.auditor && msg.sender != t.organizer, "conflicted participant");
        uint8 slot = tournamentRegistry.registerParticipant(tournamentId, msg.sender, agentHash);
        emit ParticipantRegistered(tournamentId, msg.sender, agentHash, slot);
    }

    function withdrawParticipant(uint256 tournamentId) external nonReentrant {
        CourseTypes.TournamentInfo memory t = getTournament(tournamentId);
        _requireRegistrationOpen(t);
        tournamentRegistry.withdrawParticipant(tournamentId, msg.sender);
        emit ParticipantWithdrawn(tournamentId, msg.sender);
    }

    /// @notice Sponsor deposits use PrizeEscrow as the ERC20 approval spender.
    function fundTournament(uint256 tournamentId, uint256 amount) external nonReentrant {
        CourseTypes.TournamentInfo memory t = getTournament(tournamentId);
        _requireRegistrationOpen(t);
        require(amount > 0, "zero funding");
        require(msg.sender != t.auditor, "conflicted sponsor");
        prizeEscrow.fund(tournamentId, msg.sender, amount);
        emit TournamentFunded(tournamentId, msg.sender, amount, prizeEscrow.getEscrow(tournamentId).totalFunded);
    }

    function startTournament(uint256 tournamentId) external nonReentrant {
        CourseTypes.TournamentInfo memory t = getTournament(tournamentId);
        _requireOrganizer(t);
        require(t.status == CourseTypes.Status.Open, "not open");
        require(block.timestamp >= t.registrationDeadline, "registration still open");
        require(block.timestamp + t.reviewPeriod < t.settlementDeadline, "insufficient review time");
        require(t.participant1 != address(0) && t.participant2 != address(0), "need two participants");
        require(prizeEscrow.getEscrow(tournamentId).totalFunded > 0, "unfunded tournament");
        tournamentRegistry.startTournament(tournamentId);
        emit TournamentStarted(tournamentId);
    }

    /// @notice Commit a complete evidence report and reserve the entire dispute window.
    function submitResult(uint256 tournamentId, address proposedWinner, bytes32 resultHash) external nonReentrant {
        CourseTypes.TournamentInfo memory t = getTournament(tournamentId);
        _requireOrganizer(t);
        require(t.status == CourseTypes.Status.Active, "not active");
        require(_isParticipant(t, proposedWinner), "winner not participant");
        require(resultHash != bytes32(0), "zero result hash");
        require(block.timestamp + t.reviewPeriod < t.settlementDeadline, "insufficient review time");
        uint64 reviewDeadline = uint64(block.timestamp + t.reviewPeriod);
        tournamentRegistry.submitResult(tournamentId, proposedWinner, resultHash, reviewDeadline);
        resultReview.submit(tournamentId, msg.sender, proposedWinner, resultHash, reviewDeadline);
        emit ResultSubmitted(tournamentId, msg.sender, proposedWinner, resultHash, reviewDeadline);
    }

    function disputeResult(uint256 tournamentId, bytes32 disputeHash) external nonReentrant {
        CourseTypes.TournamentInfo memory t = getTournament(tournamentId);
        require(t.status == CourseTypes.Status.ResultSubmitted, "not submitted");
        require(block.timestamp < t.reviewDeadline, "review window closed");
        require(disputeHash != bytes32(0), "zero dispute hash");
        require(_isParticipant(t, msg.sender) || contributionOf(tournamentId, msg.sender) > 0, "not stakeholder");
        tournamentRegistry.disputeResult(tournamentId);
        resultReview.dispute(tournamentId, msg.sender, disputeHash);
        emit ResultDisputed(tournamentId, msg.sender, disputeHash);
    }

    function approveResult(uint256 tournamentId) external nonReentrant {
        CourseTypes.TournamentInfo memory t = getTournament(tournamentId);
        _requireAuditor(t);
        require(t.status == CourseTypes.Status.ResultSubmitted, "not submitted");
        _requireReviewComplete(t);
        tournamentRegistry.finalize(tournamentId, t.proposedWinner);
        resultReview.recordDecision(tournamentId, t.resultHash);
        prizeEscrow.setWinner(tournamentId, t.proposedWinner);
        emit ResultApproved(tournamentId, msg.sender, t.proposedWinner, t.resultHash);
    }

    /// @notice Auditor may correct the winner or reject the entire result for sponsor refunds.
    function resolveDispute(uint256 tournamentId, address winner, bytes32 decisionHash) external nonReentrant {
        CourseTypes.TournamentInfo memory t = getTournament(tournamentId);
        _requireAuditor(t);
        require(t.status == CourseTypes.Status.Disputed, "not disputed");
        _requireReviewComplete(t);
        require(decisionHash != bytes32(0), "zero decision hash");
        require(winner == address(0) || _isParticipant(t, winner), "winner not participant");
        resultReview.recordDecision(tournamentId, decisionHash);
        if (winner == address(0)) {
            _cancel(tournamentId, decisionHash, false);
        } else {
            tournamentRegistry.finalize(tournamentId, winner);
            prizeEscrow.setWinner(tournamentId, winner);
        }
        emit DisputeResolved(tournamentId, msg.sender, winner, decisionHash);
    }

    function cancelTournament(uint256 tournamentId, bytes32 reasonHash) external nonReentrant {
        CourseTypes.TournamentInfo memory t = getTournament(tournamentId);
        _requireOrganizer(t);
        require(t.status == CourseTypes.Status.Open, "already started");
        require(reasonHash != bytes32(0), "zero reason hash");
        _cancel(tournamentId, reasonHash, false);
    }

    /// @notice Permissionless timeout refunds remain possible even if every role is revoked.
    function expireTournament(uint256 tournamentId) external nonReentrant {
        CourseTypes.TournamentInfo memory t = getTournament(tournamentId);
        require(
            t.status != CourseTypes.Status.Finalized && t.status != CourseTypes.Status.Cancelled, "terminal tournament"
        );
        bool underfilled = t.status == CourseTypes.Status.Open
            && (t.participant1 == address(0) || t.participant2 == address(0))
            && block.timestamp >= t.registrationDeadline;
        require(block.timestamp >= t.settlementDeadline || underfilled, "not expired");
        _cancel(tournamentId, keccak256("DEADLINE_EXPIRED"), true);
    }

    function claimPrize(uint256 tournamentId) external nonReentrant {
        require(getTournament(tournamentId).status == CourseTypes.Status.Finalized, "not finalized");
        uint256 amount = prizeEscrow.payPrize(tournamentId, msg.sender);
        emit PrizeClaimed(tournamentId, msg.sender, amount);
    }

    function claimRefund(uint256 tournamentId) external nonReentrant {
        require(getTournament(tournamentId).status == CourseTypes.Status.Cancelled, "not cancelled");
        uint256 amount = prizeEscrow.refund(tournamentId, msg.sender);
        emit SponsorRefunded(tournamentId, msg.sender, amount);
    }

    function _requireRegistrationOpen(CourseTypes.TournamentInfo memory t) private view {
        require(t.status == CourseTypes.Status.Open && block.timestamp < t.registrationDeadline, "registration closed");
    }

    function _requireOrganizer(CourseTypes.TournamentInfo memory t) private view {
        require(msg.sender == t.organizer, "not tournament organizer");
        _checkRole(ORGANIZER_ROLE, msg.sender);
    }

    function _requireAuditor(CourseTypes.TournamentInfo memory t) private view {
        require(msg.sender == t.auditor, "not designated auditor");
        _checkRole(AUDITOR_ROLE, msg.sender);
    }

    function _requireReviewComplete(CourseTypes.TournamentInfo memory t) private view {
        require(block.timestamp >= t.reviewDeadline, "review still open");
        require(block.timestamp < t.settlementDeadline, "settlement expired");
    }

    function _isParticipant(CourseTypes.TournamentInfo memory t, address wallet) private pure returns (bool) {
        return wallet != address(0) && (wallet == t.participant1 || wallet == t.participant2);
    }

    function _cancel(uint256 tournamentId, bytes32 reasonHash, bool expired) private {
        tournamentRegistry.cancel(tournamentId);
        prizeEscrow.cancel(tournamentId);
        emit TournamentCancelled(tournamentId, msg.sender, reasonHash, expired);
    }
}
