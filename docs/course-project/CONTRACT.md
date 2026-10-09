# Course contract interface (frozen v1)

`TournamentManager` is the only application write/read entry point. It coordinates
`TournamentRegistry`, `PrizeEscrow`, and `ResultReview`. Solidity 0.8.24, OpenZeppelin
AccessControl / SafeERC20 / ReentrancyGuard. No viewer betting, entry fee, platform fee,
or backend transaction signer. Sponsors fund a winner-takes-all prize in a configured
ERC20 token; the demonstration token has 6 decimals and is not redeemable USDC.

ABI files: `contracts/abi/TournamentManager.json`, `TournamentRegistry.json`,
`PrizeEscrow.json`, `ResultReview.json`, `CourseToken.json` (plain ABI arrays).
The manager constructor is `(address token, address admin)`.

## Roles and wallet authority

- `DEFAULT_ADMIN_ROLE()` = bytes32(0). Admin may grant/revoke organizer/auditor roles.
- `ORGANIZER_ROLE()` = keccak256("ORGANIZER_ROLE"). An organizer may create a tournament.
- `AUDITOR_ROLE()` = keccak256("AUDITOR_ROLE"). An eligible auditor is designated at creation.
- `hasRole(bytes32 role,address wallet)` comes from AccessControl.
- Participants and sponsors are ordinary wallets. Participation commits a canonical owner-bound Agent profile hash.
- Only the tournament's organizer may start, submit a result, or cancel before start.
- Only its designated auditor, still holding AUDITOR_ROLE, may approve/resolve.
- The auditor is immutable and cannot be the organizer or either participant. The organizer
  also cannot participate in their own tournament. Role revocation never blocks timeout refunds.
- The designated auditor cannot sponsor that tournament, keeping financial and review authority separate.
- Admin cannot change a winner, seize escrow, reassign the auditor, or withdraw funds.
- All funds move through signed wallet calls. **Approve `prizeEscrow()`**, not the manager.

## States, timing, reads

`Status`: `Open=0`, `Active=1`, `ResultSubmitted=2`, `Disputed=3`,
`Finalized=4`, `Cancelled=5`.

Registration, withdrawal and sponsor funding are permitted while Open and
`block.timestamp < registrationDeadline`. Exactly two registered wallets are required.
Start is permitted only at/after registrationDeadline, before settlementDeadline.
Submission must leave the entire reviewPeriod before settlementDeadline.
Disputes are permitted before reviewDeadline; approval and dispute resolution are
permitted at/after reviewDeadline and strictly before settlementDeadline. Hence no
actor can shorten the review window. reviewPeriod is 60 seconds to 30 days.

Any wallet may expire a non-final tournament at/after settlementDeadline. An underfilled
Open tournament may also be expired at/after registrationDeadline. Cancellation enables
each sponsor to withdraw exactly their remaining contribution; finalization enables only
the winner to claim the entire funded prize once. Neither path loops over wallets.

Manager reads:

```solidity
function settlementToken() external view returns (address);
function tournamentRegistry() external view returns (address);
function prizeEscrow() external view returns (address);
function resultReview() external view returns (address);
function nextTournamentId() external view returns (uint256); // starts at 1
function getTournament(uint256 tournamentId) external view returns (TournamentInfo memory);
function getReview(uint256 tournamentId) external view returns (ReviewInfo memory);
function getEscrow(uint256 tournamentId) external view returns (EscrowInfo memory);
function contributionOf(uint256 tournamentId,address sponsor) external view returns (uint256);
function claimablePrize(uint256 tournamentId,address wallet) external view returns (uint256);
function refundableAmount(uint256 tournamentId,address wallet) external view returns (uint256);
```

Exact tuple field order (names also preserved in the ABI):

```solidity
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
```

getTournament / getEscrow reject unknown IDs. getReview returns an empty tuple before
submission, and rejects an unknown tournament when called through Manager. Contribution
is a remaining refundable balance (zero after refund). Finalized contributions stay as
funding records and are not refundable. `totalFunded - totalPaid` is that tournament's
outstanding escrow liability; PrizeEscrow also exposes `totalLiabilities()`.

## Wallet transactions (13 business operations)

```solidity
createTournament(string title,string metadataURI,address auditor,
  uint64 registrationDeadline,uint64 settlementDeadline,uint64 reviewPeriod) returns (uint256)
registerParticipant(uint256 tournamentId,bytes32 agentHash)
withdrawParticipant(uint256 tournamentId)
fundTournament(uint256 tournamentId,uint256 amount)
startTournament(uint256 tournamentId)
submitResult(uint256 tournamentId,address proposedWinner,bytes32 resultHash)
disputeResult(uint256 tournamentId,bytes32 disputeHash)
approveResult(uint256 tournamentId)
resolveDispute(uint256 tournamentId,address winner,bytes32 decisionHash)
cancelTournament(uint256 tournamentId,bytes32 reasonHash)
expireTournament(uint256 tournamentId)
claimPrize(uint256 tournamentId)
claimRefund(uint256 tournamentId)
```

Registration binds participant1/participant2 to participant1AgentHash/participant2AgentHash in the same slot. A wallet may leave and register a different profile only before registration closes. All hashes and amounts must be nonzero; `resolveDispute`
accepts `winner=address(0)` to cancel/refund, otherwise winner must be an entrant.
Participants or existing sponsors may open a dispute during the review window.
The chain binds each opaque Agent hash to the signing wallet; it cannot parse the profile.
The backend must independently verify the canonical profile's owner and hash before a run,
and retain that profile snapshot in the result evidence for the auditor.
`resultHash` must commit to the complete canonical report/log produced by the backend,
not just final HP. The contract authenticates roles and preserves evidence/decisions;
an independent auditor still judges correctness (it is not an on-chain game verifier).

## Events (all emitted at the Manager address)

```solidity
event TournamentCreated(uint256 indexed tournamentId,address indexed organizer,address indexed auditor,
  string title,string metadataURI,uint64 registrationDeadline,uint64 settlementDeadline,uint64 reviewPeriod);
event ParticipantRegistered(uint256 indexed tournamentId,address indexed participant,bytes32 agentHash,uint8 slot);
event ParticipantWithdrawn(uint256 indexed tournamentId,address indexed participant);
event TournamentFunded(uint256 indexed tournamentId,address indexed sponsor,uint256 amount,uint256 totalFunded);
event TournamentStarted(uint256 indexed tournamentId);
event ResultSubmitted(uint256 indexed tournamentId,address indexed submitter,address indexed proposedWinner,
  bytes32 resultHash,uint64 reviewDeadline);
event ResultDisputed(uint256 indexed tournamentId,address indexed challenger,bytes32 disputeHash);
event ResultApproved(uint256 indexed tournamentId,address indexed auditor,address indexed winner,bytes32 resultHash);
event DisputeResolved(uint256 indexed tournamentId,address indexed auditor,address indexed winner,bytes32 decisionHash);
event TournamentCancelled(uint256 indexed tournamentId,address indexed caller,bytes32 reasonHash,bool expired);
event PrizeClaimed(uint256 indexed tournamentId,address indexed winner,uint256 amount);
event SponsorRefunded(uint256 indexed tournamentId,address indexed sponsor,uint256 amount);
```

The AccessControl RoleGranted / RoleRevoked events are also emitted at Manager.
Time-expiry uses `keccak256("DEADLINE_EXPIRED")` as reasonHash. Dispute cancellation
emits DisputeResolved(winner=0) and TournamentCancelled(reasonHash=decisionHash,expired=false).
Index by chain ID + tx hash + log index; do not double-count ERC20 Transfer as another
application deposit. All financial amounts are uint256 raw units and should be strings
at REST/JSON boundaries. Read token `decimals()`; never use JS floating-point amounts.

## Deployment and limits

`CourseToken` is a separately deployed demonstration ERC20 (`USD Demonstration`, `DUSD`,
6 decimals), owner-only `mint(address,uint256)`. This course prototype uses fixed two-wallet
events, a trusted but independent designated auditor, and no real-money betting. Re-deploy
after ABI changes; old immutable Sepolia deployments cannot be repaired by editing source.
The deploy script requires an explicit admin, organizer and auditor and prints all addresses.
No public-network transaction is broadcast by this task. Public deployment remains a separate authorized step.
