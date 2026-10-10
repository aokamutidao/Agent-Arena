#!/usr/bin/env python3
"""Copy an existing local course demo safely; never reset its chain or database."""
import argparse
import json
import secrets
import shlex
import sqlite3
from datetime import datetime, timezone
from pathlib import Path
from urllib.parse import urlparse
from urllib.request import Request, urlopen

ROOT = Path(__file__).resolve().parents[1]


def local_url(value):
    url = urlparse(value)
    if url.scheme != "http" or url.hostname not in {"localhost", "127.0.0.1", "::1"}:
        raise ValueError("Only loopback HTTP endpoints are allowed")
    if url.username or url.password or url.query or url.fragment:
        raise ValueError("Endpoint must not contain credentials, query or fragment")
    return value.rstrip("/")


def request(url, data=None):
    payload = None if data is None else json.dumps(data).encode()
    req = Request(url, data=payload, headers={"Content-Type": "application/json"})
    with urlopen(req, timeout=20) as response:
        result = json.load(response)
    if isinstance(result, dict) and result.get("error"):
        raise ValueError("Local endpoint returned an error")
    return result


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--source-db", type=Path, required=True)
    parser.add_argument("--source-api", default="http://127.0.0.1:18080")
    parser.add_argument("--rpc", default="http://127.0.0.1:18545")
    args = parser.parse_args()
    source_api, rpc = local_url(args.source_api), local_url(args.rpc)
    target = ROOT / ".local" / "llm-demo"
    database = target / "course.db"
    env = ROOT / ".env.course"
    if env.exists() or database.exists():
        raise SystemExit("Existing .env.course or demo database found; kept unchanged.")
    cfg = request(source_api + "/api/config")
    chain = request(rpc, {"jsonrpc": "2.0", "id": 1, "method": "eth_chainId", "params": []})
    if cfg.get("chain_id") != 31337 or int(chain["result"], 16) != 31337:
        raise SystemExit("Refusing to prepare any network other than local chain 31337")
    if not cfg.get("deployment_configured"):
        raise SystemExit("Source API has no configured course deployment")
    manager = cfg["manager_address"].lower()
    code = request(rpc, {"jsonrpc": "2.0", "id": 2, "method": "eth_getCode", "params": [manager, "latest"]})
    if code["result"] in {"0x", "0x0"}:
        raise SystemExit("Manager is not deployed on the selected local RPC")
    source = args.source_db.resolve()
    with sqlite3.connect(source.as_uri() + "?mode=ro", uri=True) as src:
        identity = src.execute("SELECT identity FROM database_identity WHERE id=1").fetchone()
        expected = f"31337:{manager}:{cfg['deployment_block']}"
        if not identity or identity[0] != expected:
            raise SystemExit("Source database does not match the selected API deployment")
        if src.execute("SELECT COUNT(*) FROM matches WHERE status IN ('queued','running')").fetchone()[0]:
            raise SystemExit("Wait for source matches to finish before making a demo copy")
        target.mkdir(parents=True, exist_ok=True)
        # SQLite's online backup includes committed WAL changes; copying only the
        # .db file could lose Agent definitions or evidence still present in WAL.
        with sqlite3.connect(database) as dst:
            src.backup(dst)
    database.chmod(0o600)
    snapshot = request(rpc, {"jsonrpc": "2.0", "id": 3, "method": "anvil_dumpState", "params": []})
    (target / "anvil-before.json").write_text(json.dumps(snapshot["result"]) + "\n")
    values = {
        "JWT_SECRET": secrets.token_hex(32), "DATABASE_PATH": str(database),
        "PORT": "18081", "FRONTEND_PORT": "3002", "AUTH_DOMAIN": "localhost:3002",
        "CORS_ORIGINS": "http://localhost:3002", "CHAIN_ID": "31337", "RPC_URL": rpc,
        "TOURNAMENT_MANAGER_ADDRESS": manager, "DEPLOYMENT_BLOCK": str(cfg["deployment_block"]),
        "CONFIRMATIONS": "0", "CONTRACT_ABI_DIR": str(ROOT / "contracts" / "abi"),
        "POLL_INTERVAL_SECONDS": "1", "ROUND_INTERVAL_MS": "300", "SIMULATION_ENABLED": "false",
        "LLM_PROVIDER": "qwen", "LLM_MODEL": "qwen-plus",
        "LLM_ENDPOINT": "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions",
        "LLM_TIMEOUT_SECONDS": "30", "LLM_MAX_CALLS": "120", "LLM_API_KEY": "",
        "NEXT_PUBLIC_API_URL": "http://localhost:18081", "NEXT_PUBLIC_LOCAL_RPC_URL": rpc,
    }
    with env.open("x") as output:
        env.chmod(0o600)
        output.write("# Local demo only. Fill model configuration; never commit this file.\n")
        output.write("\n".join(f"{key}={shlex.quote(value)}" for key, value in values.items()) + "\n")
    manifest = {
        "prepared_at": datetime.now(timezone.utc).isoformat(), "chain_id": 31337,
        "manager": manager, "source_api": source_api, "source_database": str(source),
        "database": str(database), "rpc": rpc, "api": "http://localhost:18081",
        "frontend": "http://localhost:3002", "database_backup_method": "sqlite_online_backup",
        "note": "Same local chain, separate API/database. New transactions will also appear in the original API.",
    }
    (target / "preparation.json").write_text(json.dumps(manifest, indent=2) + "\n")
    print(json.dumps(manifest, indent=2))


if __name__ == "__main__":
    main()
