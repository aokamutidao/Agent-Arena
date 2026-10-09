#!/usr/bin/env bash
set -euo pipefail
ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR/contracts"
# Exact revisions already recorded by the hackathon's foundry.lock.
if [ ! -f lib/forge-std/src/Test.sol ]; then
  forge install --no-git foundry-rs/forge-std@bf647bd6046f2f7da30d0c2bf435e5c76a780c1b
fi
if [ ! -f lib/openzeppelin-contracts/contracts/access/AccessControl.sol ]; then
  forge install --no-git OpenZeppelin/openzeppelin-contracts@cab19933c33c2ad1d4c7a84864a3601dddfd16f3
fi
