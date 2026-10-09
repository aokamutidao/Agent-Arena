// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

/// @notice ABI types shared by the course's tournament contracts.
library CourseTypes {
    enum Status {
        Open,
        Active,
        ResultSubmitted,
        Disputed,
        Finalized,
        Cancelled
    }

    struct TournamentInfo {
        uint256 id;
        address organizer;
        address auditor;
        uint64 registrationDeadline;
        uint64 settlementDeadline;
        uint64 reviewPeriod;
        Status status;
        address participant1;
        address participant2;
        bytes32 participant1AgentHash;
        bytes32 participant2AgentHash;
        address proposedWinner;
        address winner;
        bytes32 resultHash;
        uint64 reviewDeadline;
        string title;
        string metadataURI;
    }

    struct ReviewInfo {
        address submitter;
        address proposedWinner;
        bytes32 resultHash;
        uint64 submittedAt;
        uint64 reviewDeadline;
        address challenger;
        bytes32 disputeHash;
        bytes32 decisionHash;
    }

    struct EscrowInfo {
        uint256 totalFunded;
        uint256 totalPaid;
        address winner;
        bool cancelled;
        bool prizeClaimed;
    }
}
