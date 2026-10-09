// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import "./CourseTypes.sol";

/// @notice Append-only submission, dispute and decision evidence, controlled by Manager.
contract ResultReview {
    address public immutable manager;
    mapping(uint256 => CourseTypes.ReviewInfo) private reviews;

    modifier onlyManager() {
        require(msg.sender == manager, "not manager");
        _;
    }

    constructor(address _manager) {
        require(_manager != address(0), "zero manager");
        manager = _manager;
    }

    function getReview(uint256 tournamentId) external view returns (CourseTypes.ReviewInfo memory) {
        return reviews[tournamentId];
    }

    function submit(
        uint256 tournamentId,
        address submitter,
        address proposedWinner,
        bytes32 resultHash,
        uint64 reviewDeadline
    ) external onlyManager {
        require(reviews[tournamentId].submitter == address(0), "result already recorded");
        reviews[tournamentId] = CourseTypes.ReviewInfo({
            submitter: submitter,
            proposedWinner: proposedWinner,
            resultHash: resultHash,
            submittedAt: uint64(block.timestamp),
            reviewDeadline: reviewDeadline,
            challenger: address(0),
            disputeHash: bytes32(0),
            decisionHash: bytes32(0)
        });
    }

    function dispute(uint256 tournamentId, address challenger, bytes32 disputeHash) external onlyManager {
        CourseTypes.ReviewInfo storage review = reviews[tournamentId];
        require(review.submitter != address(0) && review.challenger == address(0), "invalid dispute");
        review.challenger = challenger;
        review.disputeHash = disputeHash;
    }

    function recordDecision(uint256 tournamentId, bytes32 decisionHash) external onlyManager {
        CourseTypes.ReviewInfo storage review = reviews[tournamentId];
        require(review.submitter != address(0) && review.decisionHash == bytes32(0), "invalid decision");
        review.decisionHash = decisionHash;
    }
}
