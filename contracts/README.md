# Tournament prize contracts

The course application uses `TournamentManager` to coordinate `TournamentRegistry`,
`PrizeEscrow`, and `ResultReview`. The [frozen interface](../docs/course-project/CONTRACT.md)
specifies roles, timing, states, ABI tuples and 13 independently signed business operations.
Two wallets register immutable Agent profile hashes, sponsors fund a winner-takes-all
prize, and a designated independent auditor reviews the evidence. There is no viewer
betting, entry fee, platform fee, or backend transaction signer.

The designated auditor cannot organize, enter, or sponsor that tournament. Admin manages
roles but cannot withdraw escrow, replace the auditor, or change a winner. Wallets approve
**PrizeEscrow**, then call Manager's `fundTournament`. Prize/refund withdrawals are pull
payments. Timeout cancellation is permissionless, including after role revocation.

`CourseToken` is an owner-minted, 6-decimal demonstration token (DUSD), not redeemable
currency or genuine USDC. Fee-on-transfer deposits revert atomically. The auditor judges
game correctness; the chain authenticates roles and preserves hashes and decisions,
rather than replaying the AI game on-chain.

## Reproduce the checks

Run from `contracts/`, using existing `lib/` dependencies. These commands do not broadcast.

```sh
forge build --offline
python3 export_abi.py
forge fmt --check
forge test --offline -vv
forge test --offline --gas-report --gas-snapshot-emit false
forge coverage --offline --report summary
forge build --offline --sizes
```

The ABI export creates tracked JSON arrays in `abi/` and typed frontend `course-abi.ts`.
All new application events are emitted at Manager, allowing one event stream to be indexed
by chain ID + transaction hash + log index. Do not double-count ERC20 Transfer events as
additional application funding.

Tests cover business lifecycles, role/conflict enforcement, deadline boundaries, double
payments, multiple sponsors/tournaments, token callback reentrancy, payout failure rollback
and retry, incompatible-token rollback, and legacy vulnerability regressions. Two fuzz
tests each run 256 cases. A stateful invariant executes 64 sequences of 64 randomized
operations and compares escrow balance against global and per-tournament liabilities.

Measured results are in `reports/`; controlled call-gas snapshots are in `snapshots/`.
Foundry 1.5.1 / Solc 0.8.24 / optimizer 200, using `forge test --offline -vv`:

| Sponsors | Winner claim call gas | One sponsor refund call gas |
| --- | ---: | ---: |
| 1 | 66,162 | 63,976 |
| 16 | 66,162 | 63,976 |
| 64 | 66,162 | 63,976 |

These `snapshotGasLastCall` measurements use identically warmed test fixtures and exclude
transaction intrinsic gas. They demonstrate payment work independent of sponsor count,
not a public-network fee quotation. Gas-report profiling produces different measurements;
disable snapshot emission in that mode. Wallets still need live gas estimation.
Funding uses one direct sponsor-to-escrow transfer, no sponsor array, and bounded two-wallet
registration. Settlement and payments contain no sponsor iteration.

## Explicit course deployment

The script `script/DeployCourse.s.sol:DeployCourse` reads only explicitly provided
`COURSE_DEPLOYER_KEY`, `COURSE_ADMIN`, `COURSE_ORGANIZER`, `COURSE_AUDITOR`, and optional
`COURSE_TOKEN` (unset deploys CourseToken). The signer must equal COURSE_ADMIN, and the
auditor and organizer must differ. Deployment keys belong only on the deployment machine;
frontend wallets sign business actions and the backend does not require any key.

```sh
# Simulation against a running local Anvil node. Add --broadcast only for an intentional deployment.
forge script script/DeployCourse.s.sol:DeployCourse --rpc-url http://127.0.0.1:8545
```

The script prints chain ID, block, token and all four application contract addresses.
It grants organizer/auditor roles but does not seed user funds. Use only explicit local
test accounts for a local trial. Further setup is in the
[course deployment guide](../docs/course-project/DEPLOYMENT.md).

## Legacy migration

The historical `AgentArena` contracts are deprecated. Source now rejects changing stake
sides, zero/expired/unknown-game bets, duplicate Agent replacement, and unprivileged token
mint/burn. Votes preserve the original wallet. Empty winning pools and abandoned games
allow pull refunds. Legacy StrategyVoting's internal ABI changed and GameStatus now adds
Cancelled; old bindings must not be used with the repaired source. Legacy Owner-controlled
results do not meet the new independent audit design.

**Source edits cannot repair existing immutable Sepolia deployments.** Keep those historical
addresses disconnected from course writes and use a new deployment and the new Manager ABI.

## Foundry

**Foundry is a blazing fast, portable and modular toolkit for Ethereum application development written in Rust.**

Foundry consists of:

- **Forge**: Ethereum testing framework (like Truffle, Hardhat and DappTools).
- **Cast**: Swiss army knife for interacting with EVM smart contracts, sending transactions and getting chain data.
- **Anvil**: Local Ethereum node, akin to Ganache, Hardhat Network.
- **Chisel**: Fast, utilitarian, and verbose solidity REPL.

## Documentation

https://book.getfoundry.sh/

## Usage

### Build

```shell
$ forge build
```

### Test

```shell
$ forge test
```

### Format

```shell
$ forge fmt
```

### Gas Snapshots

```shell
$ forge snapshot
```

### Anvil

```shell
$ anvil
```

### Deploy

```shell
$ forge script script/DeployCourse.s.sol:DeployCourse --rpc-url <your_explicit_rpc_url>
```

### Cast

```shell
$ cast <subcommand>
```

### Help

```shell
$ forge --help
$ anvil --help
$ cast --help
```
