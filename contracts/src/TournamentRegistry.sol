// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import "./CourseTypes.sol";

/// @notice Immutable manager-controlled tournament and participant records.
contract TournamentRegistry {
    address public immutable manager;
    uint256 public nextTournamentId = 1;
    mapping(uint256 => CourseTypes.TournamentInfo) private tournaments;

    modifier onlyManager() {
        require(msg.sender == manager, "not manager");
        _;
    }

    constructor(address _manager) {
        require(_manager != address(0), "zero manager");
        manager = _manager;
    }

    /// @notice Read an existing tournament. Zero/default records are never valid.
    function getTournament(uint256 tournamentId) external view returns (CourseTypes.TournamentInfo memory) {
        require(tournaments[tournamentId].organizer != address(0), "unknown tournament");
        return tournaments[tournamentId];
    }

    /// @dev Called atomically with creation of the isolated escrow pool.
    function createTournament(
        address organizer,
        address auditor,
        string calldata title,
        string calldata metadataURI,
        uint64 registrationDeadline,
        uint64 settlementDeadline,
        uint64 reviewPeriod
    ) external onlyManager returns (uint256 tournamentId) {
        tournamentId = nextTournamentId++;
        CourseTypes.TournamentInfo storage t = tournaments[tournamentId];
        t.id = tournamentId;
        t.organizer = organizer;
        t.auditor = auditor;
        t.registrationDeadline = registrationDeadline;
        t.settlementDeadline = settlementDeadline;
        t.reviewPeriod = reviewPeriod;
        t.title = title;
        t.metadataURI = metadataURI;
    }

    /// @notice Bind a wallet to its immutable canonical Agent profile hash for this event.
    function registerParticipant(uint256 tournamentId, address participant, bytes32 agentHash)
        external
        onlyManager
        returns (uint8 slot)
    {
        CourseTypes.TournamentInfo storage t = tournaments[tournamentId];
        require(t.status == CourseTypes.Status.Open, "not open");
        require(t.participant1 != participant && t.participant2 != participant, "already registered");
        if (t.participant1 == address(0)) {
            t.participant1 = participant;
            t.participant1AgentHash = agentHash;
            return 1;
        }
        require(t.participant2 == address(0), "tournament full");
        t.participant2 = participant;
        t.participant2AgentHash = agentHash;
        return 2;
    }

    /// @notice A wallet can leave before registration closes; its hash is removed with it.
    function withdrawParticipant(uint256 tournamentId, address participant) external onlyManager {
        CourseTypes.TournamentInfo storage t = tournaments[tournamentId];
        if (t.participant1 == participant) {
            t.participant1 = address(0);
            t.participant1AgentHash = bytes32(0);
        } else {
            require(t.participant2 == participant, "not participant");
            t.participant2 = address(0);
            t.participant2AgentHash = bytes32(0);
        }
    }

    function startTournament(uint256 tournamentId) external onlyManager {
        require(tournaments[tournamentId].status == CourseTypes.Status.Open, "not open");
        tournaments[tournamentId].status = CourseTypes.Status.Active;
    }

    function submitResult(uint256 tournamentId, address proposedWinner, bytes32 resultHash, uint64 reviewDeadline)
        external
        onlyManager
    {
        CourseTypes.TournamentInfo storage t = tournaments[tournamentId];
        require(t.status == CourseTypes.Status.Active, "not active");
        t.proposedWinner = proposedWinner;
        t.resultHash = resultHash;
        t.reviewDeadline = reviewDeadline;
        t.status = CourseTypes.Status.ResultSubmitted;
    }

    function disputeResult(uint256 tournamentId) external onlyManager {
        require(tournaments[tournamentId].status == CourseTypes.Status.ResultSubmitted, "not submitted");
        tournaments[tournamentId].status = CourseTypes.Status.Disputed;
    }

    function finalize(uint256 tournamentId, address winner) external onlyManager {
        CourseTypes.TournamentInfo storage t = tournaments[tournamentId];
        require(
            t.status == CourseTypes.Status.ResultSubmitted || t.status == CourseTypes.Status.Disputed, "not reviewable"
        );
        t.winner = winner;
        t.status = CourseTypes.Status.Finalized;
    }

    function cancel(uint256 tournamentId) external onlyManager {
        CourseTypes.TournamentInfo storage t = tournaments[tournamentId];
        require(
            t.status != CourseTypes.Status.Finalized && t.status != CourseTypes.Status.Cancelled, "terminal tournament"
        );
        t.status = CourseTypes.Status.Cancelled;
    }
}
