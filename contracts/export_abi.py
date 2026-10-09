#!/usr/bin/env python3
"""Export compiled course ABIs without Node dependencies or private configuration."""
import json
from pathlib import Path

ROOT = Path(__file__).resolve().parent
NAMES = ("TournamentManager", "TournamentRegistry", "PrizeEscrow", "ResultReview", "CourseToken")


def main():
    target = ROOT / "abi"
    target.mkdir(exist_ok=True)
    abis = {}
    for name in NAMES:
        artifact = json.loads((ROOT / "out" / (name + ".sol") / (name + ".json")).read_text())
        abis[name] = artifact["abi"]
        (target / (name + ".json")).write_text(json.dumps(artifact["abi"], indent=2) + "\n")
    frontend = ROOT.parent / "frontend" / "src" / "lib" / "course-abi.ts"
    frontend.write_text(
        "// Generated from contracts/abi by contracts/export_abi.py. Do not edit by hand.\n"
        + "export const tournamentManagerAbi = " + json.dumps(abis["TournamentManager"], indent=2) + " as const;\n\n"
        + "export const courseTokenAbi = " + json.dumps(abis["CourseToken"], indent=2) + " as const;\n"
    )
    print("Exported course ABIs and frontend course-abi.ts")


if __name__ == "__main__":
    main()
