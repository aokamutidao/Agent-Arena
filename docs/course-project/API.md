# Course API contract

The course server is `backend/cmd/course-server`. It does not load `.env`, hold a transaction-signing key, or send blockchain transactions. All financial writes are signed by the user's wallet against the addresses returned by `/api/config`. There is no mock financial deployment. An unconfigured chain returns an empty dashboard/list, and tournament operations return a structured 503.

All addresses are lowercase Ethereum hex addresses; raw token amounts are **decimal strings**. IDs and block numbers are integers. All time fields are **Unix seconds**, including `registration_deadline`, `created_at`, `expires_at`, and event `timestamp`. Token decimals come from config. States and roles are lowercase strings. HTTP errors have `{ "error": { "code": "...", "message": "..." } }`.

## Configuration and authentication

- `GET /api/config`: `{chain_id, deployment_configured, manager_address, escrow_address, registry_address, review_address, token_address, token_decimals, deployment_block, confirmations, simulation_enabled, auth_domain}`. Configure only the read-only RPC and manager address: startup discovers the three immutable child addresses, settlement token and actual token decimals from the manager, checks the chain ID and deployed manager code, and returns the discovered addresses. The readiness endpoint additionally checks indexing health.
- `POST /api/auth/nonce` body `{address}`: `{nonce,message,expires_at}`. Sign the returned message unchanged with `personal_sign` / `signMessage`; it binds the domain, address, chain, manager deployment resource, random nonce, and five-minute expiry. The server validates that a persisted challenge still belongs to the current domain/chain/deployment before accepting its signature.
- `POST /api/auth/login` body `{address,nonce,signature}`: `{token,user:{address,roles:["participant",...]}}`. A nonce can be used once. Additional roles (`admin`, `organizer`, `auditor`) come from the configured manager's `hasRole`, never the request. Tokens expire after one hour and validate domain audience, chain ID and manager address on every authenticated request.
- `GET /api/auth/profile`, `Authorization: Bearer <token>`: `{address,roles:string[]}`. Roles are checked again against the chain.

## Real model execution and diagnostics (iteration 033)

`GET /api/config` additionally returns `llm_configured`, `execution_mode` (`llm`, `deterministic_simulation`, or `unavailable`), `llm_provider`, `llm_model`, `llm_timeout_seconds`, and `llm_max_calls`. Configured means a server-side client exists; it is not a successful provider connectivity test. The API never returns its endpoint or API key.

For real model matches, `Match` additionally includes:

- `llm_calls`: a persisted attempt journal. Each entry has `sequence`, `round`, `side`, `provider`, `model`, `success`, `latency_ms`, and optional `http_status`, `request_id`, `usage`, `error_code`. A provider response is only a successful decision after action parsing succeeds. Missing token usage is unknown, not zero.
- `llm_summary`: aggregate `provider`, `model`, `attempted`, `succeeded`, `failed`, cumulative `latency_ms`, nullable `usage.{prompt_tokens,completion_tokens,total_tokens}`, and `last_error_code`. The cumulative time sums both players' requests and is not match wall-clock duration.
- `retry_allowed`: the server's current eligibility for resuming an interrupted run; the frontend also requires the correct organizer, on-chain running status, and a valid settlement deadline.
- Each committed `Round` references its two decisions using `red_llm_call_id` and `blue_llm_call_id`.

Both requests are reserved in the journal before sending. Failed HTTP calls or malformed action output interrupt the match without executing a replacement WAIT or submitting a winner. Previously committed rounds remain available; the organizer can explicitly resume eligible interruptions. Resuming consumes the same persisted attempt budget (default and maximum 120); it does not reset the count. No automatic provider retry occurs. An interrupted in-flight attempt whose response is lost remains a consumed attempt.

The model protocol is one action (`MOVE(x,y)`, `ATTACK`, `SKILL`, `HEAL`, `CHARGE`, or `WAIT`) with an optional short `NOTE:` describing the action's purpose. It does not request internal chain-of-thought. Qwen 3 model requests explicitly set `enable_thinking=false`; other compatible model families do not receive this vendor-specific field. Invalid game moves are recorded as failed actions according to the game rules; they are distinct from HTTP or output-parsing failures.

New model evidence that includes the attempt journal uses schema version 2. Existing saved evidence remains unchanged, and deterministic simulations without a journal retain schema version 1. A hash binds the recorded evidence; provider/model labels and request IDs are operational records, not cryptographic proof supplied by the model vendor.

Server configuration uses `LLM_API_KEY`, `LLM_PROVIDER`, `LLM_MODEL`, `LLM_ENDPOINT`, `LLM_TIMEOUT_SECONDS`, and `LLM_MAX_CALLS`. Explicit `QWEN_API_KEY`, `QWEN_MODEL`, `QWEN_ENDPOINT`, and `QWEN_TIMEOUT_SECONDS` remain compatibility aliases. Default course settings use the general DashScope endpoint and `qwen-plus`. Provider credentials are never loaded automatically from the historical `backend/.env`.

`backend/cmd/llm-check` performs a redacted configuration check by default. A real network request requires the explicit `-call` flag. Keep actual provider checks and their request counts separate from local mocked tests.

## Immutable Agent profiles

- `GET /api/agents`: `{agents:AgentProfile[],total}`; returns the latest 100 profiles at most, with `total` equal to the returned count.
- `GET /api/agents/my`, authenticated: same envelope, only the current wallet's profiles.
- `POST /api/agents`, authenticated, body `{name,personality}`: HTTP 201 `AgentProfile`.
- `AgentProfile`: `{id,agent_hash,name,personality,owner_address,created_at}`.

Profiles are immutable; create another profile to change a definition. No endpoint accepts an API key or an arbitrary external URL. `agent_hash = keccak256(UTF8(canonical_profile_JSON))`, with the schema described by the evidence endpoint. The wallet passes this hash to the on-chain registration. A match requires both participant hashes to resolve to profiles owned by those participants.

## Tournaments

- `GET /api/tournaments?limit=20&offset=0`: `{tournaments:Tournament[],total}`; maximum limit 100.
- `GET /api/tournaments/:id`: `Tournament`; reads through to the configured chain rather than accepting client financial data.
- `Tournament`: `{id,title,metadata_uri,organizer,auditor,token,entry_fee,prize_pool,prize_paid,prize_claimed,max_participants,registration_deadline,settlement_deadline,review_period,review_deadline,status,participants:string[],participant_agents:[{address,agent_hash}],proposed_winner,winner,evidence_hash,chain_updated_at,read_block,confirmed}`. Raw token amounts are decimal strings. Zero winner/hash are represented as empty strings. Status names are `open`, `running`, `result_submitted`, `disputed`, `finalized`, `cancelled`. The list is a confirmed projection. A direct detail read is pinned consistently to one latest block and marked `confirmed:false`; it never overwrites the confirmed dashboard projection. This course model fixes `entry_fee:"0"` and `max_participants:2`.
- `GET /api/tournaments/:id/events?limit=50&offset=0`: `{events:ChainEvent[],total}`.
- `ChainEvent`: `{id,chain_id,tournament_id,contract_address,event_name,tx_hash,log_index,block_number,timestamp,actor,amount_raw,data}`. `amount_raw` is a decimal string or empty for nonfinancial events; `data` contains decoded event fields. Only logs past the configured confirmation depth are indexed. Event IDs include chain ID, transaction hash and log index; replay is idempotent.
- `POST /api/tournaments/:id/run`, authenticated, no body required: HTTP 202 `Match`. Only that tournament's organizer can request a run; the chain must be `running`, with exactly two registered participant Agent hashes. Completed matches cannot be rerun. Interrupted jobs can resume from their last persisted round. A run creates evidence only: it **never** submits a result or transfers tokens.
- `GET /api/tournaments/:id/match`: `Match`, or 404 if never started.
- `Match`: `{tournament_id,status,mode,red_address,blue_address,current_round,winner_address,outcome,evidence_hash,rounds:Round[],state,started_at,updated_at,error,red_agent,blue_agent}`. `red_agent` and `blue_agent` are the immutable `{schema_version,owner_address,name,personality}` definitions bound to the registered hashes. Job status is `queued`, `running`, `interrupted`, `completed`, or `failed`; mode is `deterministic_simulation` or `llm`. `state` reuses the existing game engine JSON (HP, positions, obstacles, cooldowns and action history). Red/blue map to the chain's participant order. `outcome` is `red`, `blue`, or `draw`; a draw has an empty winner address and requires cancellation/refunds rather than inventing a winner.
- `Round`: `{round,red_action,blue_action,red_hp,blue_hp,red_position:{x,y},blue_position:{x,y},red_reasoning,blue_reasoning,red_state,blue_state}`. Actions use the original engine's `{type,target?,failed?,fail_reason?}` representation. Full Agent states include cooldowns, charge flags and effects. Deterministic simulation chooses an aggressive, defensive or tactical baseline from explicit profile keywords (then the profile definition hash as a fallback); every reasoning field labels this mode. A simultaneous KO is a draw.
- `GET /api/tournaments/:id/evidence`: `{evidence_hash,canonical_json,evidence}` for a completed job. `canonical_json` is the exact compact UTF-8 JSON whose Keccak-256 equals `evidence_hash`. `evidence` includes schema/engine versions, chain and manager identity, immutable Agent snapshots, initial state, every action and resulting round state, final state, outcome and winner. Download/inspect this before organizer wallet submission; the auditor reviews it independently.
- Optional `GET /ws?tournament_id=<id>`: read-only WebSocket snapshots `{type:"match_update",data:Match}`. HTTP `/match` polling remains a supported fallback.

Review text stays available for independent inspection:

- `POST /api/tournaments/:id/notes`, authenticated, `{kind,text}`: HTTP 201 `{evidence_hash,canonical_json,evidence:{schema_version:1,chain_id,manager_address,tournament_id,author_address,kind,text}}`. Kind is `dispute`, `decision` or `cancellation`; text is 1–4000 characters. Decision notes require the designated auditor; cancellation notes require the organizer or auditor. Anyone logged in may save a dispute draft, while the contract validates standing when their wallet submits it.
- `GET /api/tournaments/:id/notes`: `{notes:[{schema_version,chain_id,manager_address,tournament_id,author_address,kind,text,evidence_hash,created_at,canonical_json}],total}`; returns the latest 100 notes at most, with `total` equal to the returned count. A saved note is a draft until its hash appears in a confirmed contract event; saving it never triggers a blockchain transaction.

## Dashboard and service health

- `GET /api/dashboard`: `{tournaments_total,registered_participants,active_tournaments,completed_tournaments,escrow_in_raw,prizes_paid_raw,refunds_paid_raw,ledger_entries,event_count,indexed_block,chain_head,observed_at}`. Token totals are derived only from confirmed escrow logs, not frontend callbacks, and cover the single configured token. `escrow_in_raw` is cumulative sponsorship received; subtract `prizes_paid_raw` and `refunds_paid_raw` to calculate outstanding indexed funds. Tournament counters use only confirmed indexed snapshots; `indexed_block` and `chain_head` make lag visible. No-chain mode contains zeros.
- `GET /health`: `{status:"ok"}` liveness.
- `GET /ready`: `{ready,chain_configured,indexed_block,chain_head,last_error}`; HTTP 503 when financial deployment/RPC/indexing is unavailable.

The application requires `JWT_SECRET` of at least 32 characters, a SQLite file on persistent storage, explicit CORS origins, and read-only RPC configuration. `SIMULATION_ENABLED=true` enables a clearly labelled deterministic game demonstration when no LLM is explicitly configured. The simulation still requires a real on-chain tournament and wallet-signed registrations.

The executable reads explicit environment configuration and does not automatically load `.env`: `DATABASE_PATH`, `JWT_SECRET`, `AUTH_DOMAIN`, `CHAIN_ID`, `RPC_URL`, `TOURNAMENT_MANAGER_ADDRESS`, `DEPLOYMENT_BLOCK`, `CONFIRMATIONS`, `CORS_ORIGINS`, `PORT`, `CONTRACT_ABI_DIR`, `SIMULATION_ENABLED`, and optional `QWEN_API_KEY`, `ROUND_INTERVAL_MS`, `POLL_INTERVAL_SECONDS`, `TRUSTED_PROXY_CIDRS`. Keep `RPC_URL` credentials and the LLM key out of logs. CORS is an explicit comma-separated list. Localhost login messages use an HTTP URI; production domains use HTTPS.

POST requests are limited to 30 per minute per client IP. Forwarded IP headers are ignored by default. Behind a reverse proxy, set `TRUSTED_PROXY_CIDRS` to a comma-separated list of its actual controlled IP addresses/CIDRs; malformed entries fail startup. Without that configuration, clients behind the same proxy share its limit. Never trust arbitrary Internet peers or a broad `0.0.0.0/0`/`::/0` range.

Canonical encoding is defined as the compact `encoding/json` output of the fixed, versioned structures returned in `canonical_json` (fixed field order, UTF-8, no whitespace, standard JSON escaping). Hash the exact returned string with Keccak-256; do not reformat it before verifying. Agent definitions use fields in this order: `{schema_version:1,owner_address,name,personality}`. Match evidence intentionally excludes wall-clock start/finish timestamps, so an identical deterministic replay hashes identically.

One SQLite database is bound to one chain ID, manager address and deployment block. Switching deployments requires a new database path. An offline-only profile database can be bound to its first real deployment before any tournament/event/match data exists. Indexing persists a block-hash checkpoint; a changed confirmed checkpoint pauses indexing and fails readiness until the operator rebuilds the projection from its deployment block. This detects reorganizations and does not claim mathematically irreversible finality.
