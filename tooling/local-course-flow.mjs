#!/usr/bin/env node
// Existing-chain acceptance flow. Only the public Anvil development accounts
// sign; no deployment, chain reset, real-network RPC or environment keys.
// Run --help first. --self-check is offline and sends nothing.
import assert from 'node:assert/strict';
import { randomUUID } from 'node:crypto';
import { createRequire } from 'node:module';
import { existsSync, mkdirSync, openSync, closeSync, readFileSync, writeFileSync, renameSync, unlinkSync } from 'node:fs';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const root = resolve(fileURLToPath(new URL('..', import.meta.url)));
const require = createRequire(join(root, 'frontend/package.json'));
const { createPublicClient, createWalletClient, http, keccak256, stringToHex, encodeFunctionData, decodeEventLog, parseUnits } = require('viem');
const { mnemonicToAccount } = require('viem/accounts');
const managerABI = JSON.parse(readFileSync(join(root, 'contracts/abi/TournamentManager.json'), 'utf8'));
const tokenABI = JSON.parse(readFileSync(join(root, 'contracts/abi/CourseToken.json'), 'utf8'));
const EXPECTED_MANAGER = '0xe7f1725e7734ce288f8367e1bb143e90bb3f0512';
const CHAIN_ID = 31337;
const STATUS = ['Open', 'Active', 'ResultSubmitted', 'Disputed', 'Finalized', 'Cancelled'];
const pause = (ms) => new Promise((r) => setTimeout(r, ms));
const lower = (value) => String(value).toLowerCase();
const nowISO = () => new Date().toISOString();
const json = (value) => JSON.stringify(value, (_, v) => typeof v === 'bigint' ? v.toString() : v, 2);
const integer = (value, name, min, max) => {
  const n = Number(value);
  assert.ok(Number.isSafeInteger(n) && n >= min && n <= max, `${name} must be an integer between ${min} and ${max}`);
  return n;
};
function loopbackURL(value, name) {
  const u = new URL(value);
  assert.ok(['http:', 'https:'].includes(u.protocol), `${name} must use HTTP(S)`);
  assert.ok(['127.0.0.1', 'localhost', '[::1]', '::1'].includes(u.hostname), `${name} must be loopback; public RPC/API is forbidden`);
  assert.ok(!u.username && !u.password && !u.search && !u.hash && (u.pathname === '/' || u.pathname === ''), `${name} cannot contain credentials, a path, query or fragment`);
  return u.origin;
}
function options(argv) {
  const result = { phase: undefined, rpc: 'http://127.0.0.1:18545', api: 'http://127.0.0.1:18081', state: join(root, '.local/llm-demo/business-flow.json'), manager: EXPECTED_MANAGER, wait: 900, explicit: new Set() };
  const values = new Map([['--rpc', 'rpc'], ['--api', 'api'], ['--state', 'state'], ['--expected-manager', 'manager'], ['--wait-seconds', 'wait'], ['--prize', 'prize'], ['--refund-prize', 'refundPrize'], ['--registration-seconds', 'registration'], ['--settlement-seconds', 'settlement'], ['--review-seconds', 'review']]);
  for (let i = 0; i < argv.length; i++) {
    const arg = argv[i];
    if (values.has(arg)) {
      assert.ok(argv[i + 1] && !argv[i + 1].startsWith('--'), `${arg} requires a value`);
      const key = values.get(arg); result[key] = argv[++i]; result.explicit.add(key);
    } else if (arg === '--allow-time-warp') result.timeWarp = true;
    else if (arg === '--resume-llm') result.resume = true;
    else if (arg === '--main-timeout') result.mainTimeout = true;
    else if (arg === '--help' || arg === '-h') result.help = true;
    else if (arg === '--self-check') result.selfCheck = true;
    else {
      assert.ok(!arg.startsWith('-') && !result.phase, 'Unknown or duplicate argument (--help lists accepted options)');
      result.phase = arg;
    }
  }
  result.rpc = loopbackURL(result.rpc, 'RPC'); result.api = loopbackURL(result.api, 'API');
  result.state = resolve(result.state);
  assert.match(result.manager, /^0x[0-9a-fA-F]{40}$/, 'Invalid expected manager');
  result.manager = lower(result.manager);
  result.wait = integer(result.wait, 'wait seconds', 1, 7200);
  return result;
}
function help() {
  console.log(`Existing local course chain / real-LLM acceptance (no services are started)

  node tooling/local-course-flow.mjs <agents|prepare|run|settle|cancel|status> [options]
  node tooling/local-course-flow.mjs --self-check

agents   Wallet-login and save owner-bound Agent profiles; no chain transaction.
prepare  Require configured real LLM first, then NEW tournament, sponsorship,
         red registration -> withdrawal -> re-registration, blue registration.
run      Wait for registration closure, start, request ONE LLM run, await evidence.
settle   Verify real call journal/evidence; submit, wait full review, approve, claim.
         An actual draw waits for expiry and refunds; never invents a winner.
cancel   NEW independent Open tournament -> sponsor -> cancel -> exact refund.
         --main-timeout instead expires/refunds this script's main tournament.
status   Read-only snapshot (no wallet login, file write, transaction or model call).

Defaults: --rpc http://127.0.0.1:18545 --api http://127.0.0.1:18081
          --state .local/llm-demo/business-flow.json
          --expected-manager ${EXPECTED_MANAGER}
          --prize 100 --refund-prize 25 --registration-seconds 180
          --settlement-seconds 7200 --review-seconds 60 --wait-seconds 900

--allow-time-warp  Explicit local evm_increaseTime/evm_mine, recorded in JSON.
                  Refused while other unowned tournaments are nonterminal.
--resume-llm       Explicit approval for one interrupted/uncertain run request.
                  This can incur additional model charges; never automatic.
Parameters bind to the state file. Reuse it to resume, not a fresh file.
No keys/JWT/provider request prompts/raw response bodies are written. The mnemonic is ONLY
Anvil's public development seed. Existing tournaments are never written.
This is a scripted wallet acceptance flow, not proof of a MetaMask UI signing.
LLM configured does not prove provider availability; calls may still interrupt.`);
}
function safeError(error) {
  // Provider/RPC errors can contain serialized requests or credentials. Persist
  // only our local assertions and fixed categories, never arbitrary bodies.
  if (error instanceof LocalError) return error.message;
  if (error?.code === 'ERR_ASSERTION') return String(error.message).split('\n')[0].slice(0, 300);
  if (error?.name === 'AbortError' || error?.name === 'TimeoutError') return 'Request timed out; inspect status before retrying.';
  return `Operation stopped (${String(error?.name || 'unknown').replace(/[^A-Za-z0-9_]/g, '').slice(0, 80)}); inspect the saved checkpoint and service logs privately.`;
}
class LocalError extends Error { constructor(message) { super(message); this.name = 'LocalFlowError'; } }
function isNotFound(error) {
  for (let e = error; e; e = e.cause) {
    if (['TransactionNotFoundError', 'TransactionReceiptNotFoundError'].includes(e.name)) return true;
  }
  return false;
}
function usage(value) {
  const count = (n) => Number.isSafeInteger(n) && n >= 0 ? n : null;
  return { prompt_tokens: count(value?.prompt_tokens), completion_tokens: count(value?.completion_tokens), total_tokens: count(value?.total_tokens) };
}
function realLLMProof(match) {
  assert.equal(match.status, 'completed', 'Match is not complete; no result can be committed');
  assert.equal(match.mode, 'llm', 'Simulation cannot be presented as real LLM execution');
  const summary = match.llm_summary;
  assert.ok(summary && Number.isSafeInteger(summary.succeeded) && summary.succeeded > 0, 'Missing successful real-LLM summary; legacy mode labels are insufficient');
  assert.ok(Array.isArray(match.llm_calls) && Array.isArray(match.rounds) && match.rounds.length > 0, 'Missing per-call journal or completed rounds');
  const byID = new Map();
  for (const call of match.llm_calls) {
    assert.ok(call.sequence !== undefined && !byID.has(String(call.sequence)), 'Invalid or duplicate call sequence');
    byID.set(String(call.sequence), call);
  }
  for (const round of match.rounds) {
    for (const side of ['red', 'blue']) {
      const id = round[`${side}_llm_call_id`];
      const call = byID.get(String(id));
      assert.ok(id !== undefined && call && call.success === true && call.side === side && Number(call.round) === Number(round.round), 'Executed round is not backed by two successful real LLM calls');
      assert.equal(call.provider, summary.provider, 'LLM provider changed within the match');
      assert.equal(call.model, summary.model, 'LLM model changed within the match');
    }
  }
  const succeeded = match.llm_calls.filter((c) => c.success === true).length;
  assert.equal(summary.attempted, match.llm_calls.length, 'Call count differs from summary');
  assert.equal(summary.succeeded, succeeded, 'Successful count differs from journal');
  assert.equal(summary.failed, match.llm_calls.length - succeeded, 'Failed count differs from journal');
  return {
    verified: true, status: match.status, mode: match.mode, outcome: match.outcome,
    winner_address: match.winner_address || '', rounds: match.rounds.length,
    provider: summary.provider, model: summary.model, attempted: summary.attempted,
    succeeded: summary.succeeded, failed: summary.failed, latency_ms: summary.latency_ms,
    usage: usage(summary.usage), last_error_code: summary.last_error_code || '',
    calls: match.llm_calls.map((c) => ({ sequence: c.sequence, round: c.round, side: c.side, provider: c.provider, model: c.model, success: c.success, latency_ms: c.latency_ms, http_status: c.http_status ?? null, usage: usage(c.usage), error_code: c.error_code || '' })),
    trust_boundary: 'Call metadata records provider execution as observed by the backend; it is not a provider signature, independent truth of the game, or an external audit.',
  };
}
async function selfCheck() {
  for (const url of ['http://127.0.0.1:18545', 'http://localhost:18081', 'http://[::1]:18545']) assert.ok(loopbackURL(url, 'test'));
  for (const url of ['https://example.com', 'http://127.0.0.1.example.com', 'http://127.0.0.1@evil.test', 'http://127.0.0.1:18545/?key=secret', 'file:///tmp/rpc']) assert.throws(() => loopbackURL(url, 'test'));
  const m = { status: 'completed', mode: 'llm', outcome: 'draw', winner_address: '', rounds: [{ round: 1, red_llm_call_id: 2, blue_llm_call_id: 3 }], llm_summary: { provider: 'test', model: 'test', attempted: 3, succeeded: 2, failed: 1 }, llm_calls: [{ sequence: 1, round: 1, side: 'red', provider: 'test', model: 'test', success: false }, { sequence: 2, round: 1, side: 'red', provider: 'test', model: 'test', success: true }, { sequence: 3, round: 1, side: 'blue', provider: 'test', model: 'test', success: true }] };
  assert.equal(realLLMProof(m).failed, 1, 'Recovered genuine calls retain earlier failures');
  for (const bad of [{ ...m, mode: 'deterministic_simulation' }, { ...m, llm_summary: undefined }, { ...m, status: 'interrupted' }, { ...m, rounds: [{ round: 1, red_llm_call_id: 1, blue_llm_call_id: 3 }] }]) assert.throws(() => realLLMProof(bad));
  // Offline determinism is necessary to rebroadcast the same hash after a crash.
  const a = mnemonicToAccount('test test test test test test test test test test test junk', { addressIndex: 5 });
  const request = { chainId: CHAIN_ID, nonce: 0, to: EXPECTED_MANAGER, data: '0x', value: 0n, gas: 21000n, maxFeePerGas: 1000000000n, maxPriorityFeePerGas: 100000000n, type: 'eip1559' };
  const one = await a.signTransaction(request), two = await a.signTransaction(request);
  assert.equal(keccak256(one), keccak256(two));
  console.log('Offline checks passed: loopback rejection, real-call proof, draw preservation, deterministic transaction recovery. No network requests, transactions or model calls.');
}

let state, statePath, lockPath, lockFD, publicClient, cfg, accounts, deployment, opts;
const sessions = new Map();
function save() {
  assert.ok(opts.phase !== 'status', 'status is read-only');
  state.updated_at = nowISO();
  const temporary = `${statePath}.tmp-${process.pid}`;
  writeFileSync(temporary, `${json(state)}\n`, { mode: 0o600 }); renameSync(temporary, statePath);
}
function log(message) { console.log(`[local-course-flow] ${message}`); }
async function api(path, { token, body, expected = [200], method, allow404 = false } = {}) {
  const response = await fetch(`${opts.api}${path}`, { method: method || (body === undefined ? 'GET' : 'POST'), redirect: 'error', headers: { 'Content-Type': 'application/json', ...(token ? { Authorization: `Bearer ${token}` } : {}) }, ...(body === undefined ? {} : { body: JSON.stringify(body) }), signal: AbortSignal.timeout(15000) });
  if (allow404 && response.status === 404) return null;
  if (!expected.includes(response.status)) throw new LocalError(`${method || (body === undefined ? 'GET' : 'POST')} ${path}: HTTP ${response.status}; response body omitted.`);
  return response.json();
}
const read = (name, args = []) => publicClient.readContract({ address: opts.manager, abi: managerABI, functionName: name, args });
const tokenRead = (name, args = []) => publicClient.readContract({ address: deployment.token, abi: tokenABI, functionName: name, args });
const blockTime = async () => Number((await publicClient.getBlock({ blockTag: 'latest' })).timestamp);
async function login(actor) {
  if (sessions.has(actor)) return sessions.get(actor);
  const account = accounts[actor];
  const challenge = await api('/api/auth/nonce', { body: { address: account.address } });
  assert.ok(typeof challenge.message === 'string' && challenge.message.startsWith(`${cfg.auth_domain} wants you to sign in with your Ethereum account:\n${lower(account.address)}\n`) && challenge.message.includes(`\nChain ID: ${CHAIN_ID}\n`) && challenge.message.endsWith(`\nResources:\n- urn:agent-arena:${CHAIN_ID}:${opts.manager}`), 'Unexpected wallet login challenge domain/chain/deployment/address');
  const response = await api('/api/auth/login', { body: { address: account.address, nonce: challenge.nonce, signature: await account.signMessage({ message: challenge.message }) } });
  assert.equal(lower(response.user.address), lower(account.address), 'Login identity mismatch');
  sessions.set(actor, response.token); return response.token;
}
function requireLLM() {
  assert.ok(cfg.llm_configured === true && cfg.execution_mode === 'llm' && typeof cfg.llm_provider === 'string' && cfg.llm_provider.length && typeof cfg.llm_model === 'string' && cfg.llm_model.length, 'Real LLM is unavailable: no tournament/funding/run created. Configure the separate demo API first; agents/status are safe without it.');
}
function argsEqual(a, b) { return json(a) === json(b); }
async function receiptOrNull(hash) {
  try { return await publicClient.getTransactionReceipt({ hash }); }
  catch (e) { if (isNotFound(e)) return null; throw e; }
}
function unsigned(request) {
  const out = {};
  for (const key of ['chainId', 'nonce', 'to', 'data', 'value', 'gas', 'gasPrice', 'maxFeePerGas', 'maxPriorityFeePerGas', 'type']) if (request[key] !== undefined) out[key] = typeof request[key] === 'bigint' ? request[key].toString() : request[key];
  return out;
}
function restoreUnsigned(request) {
  const out = { ...request };
  for (const key of ['value', 'gas', 'gasPrice', 'maxFeePerGas', 'maxPriorityFeePerGas']) if (out[key] !== undefined) out[key] = BigInt(out[key]);
  return out;
}
async function transact(key, actor, name, args = [], isToken = false) {
  const account = accounts[actor], to = isToken ? deployment.token : opts.manager, abi = isToken ? tokenABI : managerABI;
  const data = encodeFunctionData({ abi, functionName: name, args });
  let item = state.transactions.find((t) => t.key === key);
  if (item) {
    assert.equal(item.actor, actor, 'Checkpoint signer mismatch'); assert.equal(lower(item.to), lower(to), 'Checkpoint target mismatch'); assert.equal(item.request.data, data, 'Checkpoint calldata mismatch');
    assert.notEqual(item.status, 'reverted', 'A prior transaction reverted; do not silently create another');
  } else {
    assert.equal(await publicClient.getChainId(), CHAIN_ID, 'Chain changed before signing');
    const wallet = createWalletClient({ account, chain: publicClient.chain, transport: http(opts.rpc, { retryCount: 0, fetchOptions: { redirect: 'error' } }) });
    const request = unsigned(await wallet.prepareTransactionRequest({ account, to, data, value: 0n }));
    const signed = await account.signTransaction(restoreUnsigned(request));
    item = { key, actor, address: lower(account.address), operation: name, to: lower(to), request, tx_hash: keccak256(signed), status: 'prepared', prepared_at: nowISO() };
    // The hash and exact unsigned request survive crashes; signer objects and
    // serialized signed payloads never go into evidence.
    state.transactions.push(item); save();
  }
  let receipt = await receiptOrNull(item.tx_hash);
  if (!receipt) {
    let pending = false;
    try { await publicClient.getTransaction({ hash: item.tx_hash }); pending = true; }
    catch (e) { if (!isNotFound(e)) throw e; }
    if (!pending) {
      const nonce = await publicClient.getTransactionCount({ address: account.address, blockTag: 'pending' });
      assert.equal(nonce, item.request.nonce, 'Saved transaction is absent but nonce changed; manual reconciliation required, no extra transaction sent');
      assert.equal(await publicClient.getChainId(), CHAIN_ID, 'Chain changed before broadcasting');
      const serializedTransaction = await account.signTransaction(restoreUnsigned(item.request));
      assert.equal(keccak256(serializedTransaction), item.tx_hash, 'Recovery signature hash mismatch');
      item.broadcast_attempts = (item.broadcast_attempts || 0) + 1; item.status = 'broadcasting'; save();
      const hash = await publicClient.sendRawTransaction({ serializedTransaction });
      assert.equal(hash, item.tx_hash, 'Broadcast transaction hash mismatch');
    }
    log(`${name}: checking receipt ${item.tx_hash}`);
    receipt = await publicClient.waitForTransactionReceipt({ hash: item.tx_hash, timeout: 120000, pollingInterval: 1000 });
  }
  assert.equal(lower(receipt.from), lower(account.address), 'Receipt signer mismatch'); assert.equal(lower(receipt.to), lower(to), 'Receipt target mismatch');
  item.status = receipt.status === 'success' ? 'confirmed' : 'reverted'; item.block_number = receipt.blockNumber.toString(); item.block_hash = receipt.blockHash; item.gas_used = receipt.gasUsed.toString(); item.effective_gas_price = receipt.effectiveGasPrice.toString(); item.gas_cost_wei = (receipt.gasUsed * receipt.effectiveGasPrice).toString(); item.confirmed_at ||= nowISO(); save();
  assert.equal(receipt.status, 'success', `${name} reverted; no automatic replacement`);
  return receipt;
}
async function observe(label) {
  const block = await publicClient.getBlock();
  const balances = {};
  for (const actor of ['organizer', 'auditor', 'red', 'blue', 'sponsor']) balances[actor] = { address: lower(accounts[actor].address), dusd_raw: (await tokenRead('balanceOf', [accounts[actor].address])).toString(), eth_wei: (await publicClient.getBalance({ address: accounts[actor].address })).toString() };
  balances.escrow = { address: deployment.escrow, dusd_raw: (await tokenRead('balanceOf', [deployment.escrow])).toString() };
  const snapshot = { label, recorded_at: nowISO(), block_number: block.number.toString(), chain_timestamp: Number(block.timestamp), balances };
  if (opts.phase !== 'status') { state.snapshots.push(snapshot); save(); }
  return snapshot;
}
async function owned(which) {
  const item = state.tournaments[which]; assert.ok(item?.id, `No ${which} tournament in this checkpoint; run prepare/cancel first`);
  assert.ok(BigInt(item.id) >= BigInt(state.initial_next_tournament_id), 'Refusing historical tournament id');
  const t = await read('getTournament', [BigInt(item.id)]);
  assert.equal(t.title, item.title, 'Tournament title does not match this run'); assert.equal(lower(t.organizer), lower(accounts.organizer.address), 'Tournament organizer changed'); assert.equal(lower(t.auditor), lower(accounts.auditor.address), 'Tournament auditor mismatch');
  assert.equal(Number(t.registrationDeadline), item.registration_deadline, 'Registration deadline mismatch'); assert.equal(Number(t.settlementDeadline), item.settlement_deadline, 'Settlement deadline mismatch');
  return { ...item, chain: t, id: BigInt(item.id) };
}
async function create(which, amount) {
  let item = state.tournaments[which];
  if (!item) {
    requireLLM(); // A missing key must never strand a new funded tournament.
    const chainTime = await blockTime();
    item = { title: `Local LLM ${which} ${state.run_id.slice(0, 12)}`, registration_deadline: chainTime + state.parameters.registration_seconds, settlement_deadline: chainTime + state.parameters.settlement_seconds, review_seconds: state.parameters.review_seconds, sponsor_amount_raw: amount.toString() };
    state.tournaments[which] = item; save();
  }
  if (!item.id) {
    const receipt = await transact(`${which}.create`, 'organizer', 'createTournament', [item.title, '', accounts.auditor.address, item.registration_deadline, item.settlement_deadline, item.review_seconds]);
    const event = receipt.logs.filter((l) => lower(l.address) === opts.manager).map((l) => { try { return decodeEventLog({ abi: managerABI, data: l.data, topics: l.topics }); } catch { return null; } }).find((e) => e?.eventName === 'TournamentCreated' && e.args.title === item.title);
    assert.ok(event, 'Creation receipt lacks expected TournamentCreated event');
    item.id = event.args.tournamentId.toString(); save();
  }
  return owned(which);
}
async function profiles() {
  for (const [actor, name, personality] of [
    ['red', 'Local LLM Red', 'Aggressive but legal: close distance with a listed legal MOVE, then ATTACK or SKILL when legal. Seek a quick fair encounter; avoid endless WAIT. Return exactly one allowed action and optionally one brief NOTE. Never choose an illegal action or assume a fixed winner.'],
    ['blue', 'Local LLM Blue', 'Tactical but decisive: move toward the opponent using a listed legal MOVE, use ATTACK or SKILL when legal, and HEAL only when useful. Seek a quick fair encounter; avoid endless retreat or WAIT. Return exactly one allowed action and optionally one brief NOTE. Never choose an illegal action or assume a fixed winner.'],
  ]) {
    const definition = { schema_version: 1, owner_address: lower(accounts[actor].address), name, personality };
    const hash = keccak256(stringToHex(JSON.stringify(definition)));
    if (!state.agents[actor]) {
      const agent = await api('/api/agents', { token: await login(actor), body: { name, personality }, expected: [201] });
      assert.equal(lower(agent.agent_hash), hash, 'Agent canonical hash mismatch'); assert.equal(lower(agent.owner_address), definition.owner_address, 'Agent owner mismatch');
      state.agents[actor] = { ...definition, agent_hash: hash }; save();
    } else assert.ok(argsEqual(state.agents[actor], { ...definition, agent_hash: hash }), 'Saved Agent changed');
  }
}
async function ensureFunded(which) {
  const t = await owned(which), amount = BigInt(t.sponsor_amount_raw);
  const key = `${which}.fund`;
  if (!state.transactions.some((v) => v.key === key)) {
    assert.equal(await read('contributionOf', [t.id, accounts.sponsor.address]), 0n, 'Unexpected prior funding; refuse duplicate sponsorship');
    assert.equal(Number(t.chain.status), 0, 'New funding requires Open tournament');
    assert.ok(await blockTime() < t.registration_deadline, 'Registration expired before sponsorship; no funds sent');
    assert.ok(await tokenRead('balanceOf', [accounts.sponsor.address]) >= amount, 'Sponsor has insufficient existing DUSD; this script never mints');
    const allowance = await tokenRead('allowance', [accounts.sponsor.address, deployment.escrow]);
    if (allowance < amount) await transact(`${which}.approve-token`, 'sponsor', 'approve', [deployment.escrow, amount], true);
  }
  await transact(key, 'sponsor', 'fundTournament', [t.id, amount]);
  // contributionOf becomes zero after a refund. The immutable funded total is
  // the recovery check when a crash happened after a successful payout.
  assert.equal((await read('getEscrow', [t.id])).totalFunded, amount, 'Sponsorship total mismatch');
}
async function waitChain(target, reason) {
  assert.ok(!state.time_adjustments.some((v) => v.status === 'uncertain_not_repeated'), 'Prior time adjustment is uncertain; manual clock/checkpoint reconciliation required before any further increase');
  const unfinished = state.time_adjustments.find((v) => ['intent', 'increase_acknowledged'].includes(v.status));
  if (unfinished) {
    assert.ok(opts.timeWarp, 'An explicit time adjustment was interrupted; rerun with --allow-time-warp to reconcile the pending mine');
    assert.equal(await publicClient.getChainId(), CHAIN_ID, 'Wrong chain for time adjustment recovery');
    // Mine once without another increase: if the RPC response was lost after
    // increaseTime, applying the delta again would corrupt the intended clock.
    await publicClient.request({ method: 'evm_mine', params: [] });
    unfinished.after_timestamp = await blockTime();
    unfinished.status = unfinished.after_timestamp >= unfinished.requested_target ? 'recovered_after_interruption' : 'uncertain_not_repeated'; save();
    assert.ok(unfinished.after_timestamp >= unfinished.requested_target, 'Interrupted time adjustment is uncertain; no additional increase sent. Reconcile the local chain clock manually');
  }
  let time = await blockTime(); if (time >= target) return;
  if (opts.timeWarp) {
    const ownedIDs = new Set(Object.values(state.tournaments).filter((t) => t.id).map((t) => String(t.id)));
    const next = await read('nextTournamentId');
    assert.ok(next <= 1000n, 'Too many tournaments to safely inspect time-warp effects');
    for (let id = 1n; id < next; id++) {
      if (ownedIDs.has(String(id))) continue;
      const other = await read('getTournament', [id]);
      assert.ok(Number(other.status) >= 4, `Time warp refused: unrelated tournament #${id} is nonterminal; wait naturally`);
    }
    const warp = { reason, before_timestamp: time, requested_target: target, increase_seconds: target - time, recorded_at: nowISO(), status: 'intent', affects_shared_chain_clock: true };
    state.time_adjustments.push(warp); save();
    assert.equal(await publicClient.getChainId(), CHAIN_ID, 'Wrong chain for time warp');
    await publicClient.request({ method: 'evm_increaseTime', params: [target - time] });
    warp.status = 'increase_acknowledged'; save();
    await publicClient.request({ method: 'evm_mine', params: [] });
    warp.after_timestamp = await blockTime(); warp.status = 'applied'; save();
    assert.ok(warp.after_timestamp >= target, 'Time warp did not reach target');
    log(`Explicit local time adjustment: ${reason} (+${warp.increase_seconds}s)`); return;
  }
  const started = Date.now(); let lastLog = 0;
  while ((time = await blockTime()) < target) {
    if (Date.now() - started > opts.wait * 1000) throw new LocalError(`Still waiting for ${reason}; target chain timestamp ${target}. Rerun the same phase later, or explicitly allow local time warp.`);
    if (Date.now() - lastLog > 15000) { log(`${reason}: ${target - time}s of chain time remaining`); lastLog = Date.now(); }
    await pause(2000);
  }
}
async function indexed(id) {
  const started = Date.now(), target = Number(await publicClient.getBlockNumber({ cacheTime: 0 }));
  while (Date.now() - started < 30000) {
    const ready = await api('/ready', { expected: [200, 503] });
    if (ready.ready && ready.indexed_block >= target) return api(`/api/tournaments/${id}`);
    await pause(1000);
  }
  throw new LocalError('API indexer has not caught up; no model request made. Retry after the independent demo API is ready.');
}
async function prepare() {
  requireLLM(); await profiles();
  if (!state.checkpoints.prepare_started) { await observe('before_main_prepare'); state.checkpoints.prepare_started = nowISO(); save(); }
  const t = await create('main', parseUnits(state.parameters.prize, deployment.decimals));
  await ensureFunded('main');
  for (const [key, actor, name, args] of [['main.register-red-first', 'red', 'registerParticipant', [t.id, state.agents.red.agent_hash]], ['main.withdraw-red', 'red', 'withdrawParticipant', [t.id]], ['main.register-red-again', 'red', 'registerParticipant', [t.id, state.agents.red.agent_hash]], ['main.register-blue', 'blue', 'registerParticipant', [t.id, state.agents.blue.agent_hash]]]) await transact(key, actor, name, args);
  const after = await owned('main');
  assert.equal(lower(after.chain.participant1), lower(accounts.red.address)); assert.equal(lower(after.chain.participant2), lower(accounts.blue.address));
  assert.equal(after.chain.participant1AgentHash, state.agents.red.agent_hash); assert.equal(after.chain.participant2AgentHash, state.agents.blue.agent_hash);
  if (!state.checkpoints.prepared) { await observe('after_main_prepare'); state.checkpoints.prepared = nowISO(); save(); }
  log(`New main tournament #${t.id} prepared; no model request yet. Registration closes at chain timestamp ${t.registration_deadline}.`);
}
async function completedMatch() {
  const t = await owned('main'); const m = await api(`/api/tournaments/${t.id}/match`, { allow404: true });
  assert.ok(m, 'No match exists; run the run phase first');
  const proof = realLLMProof(m);
  assert.equal(Number(m.tournament_id), Number(t.id), 'Match id mismatch'); assert.equal(lower(m.red_address), lower(t.chain.participant1), 'Red wallet mismatch'); assert.equal(lower(m.blue_address), lower(t.chain.participant2), 'Blue wallet mismatch');
  assert.equal(keccak256(stringToHex(JSON.stringify(m.red_agent))), t.chain.participant1AgentHash, 'Red match Agent differs from registered hash'); assert.equal(keccak256(stringToHex(JSON.stringify(m.blue_agent))), t.chain.participant2AgentHash, 'Blue match Agent differs from registered hash');
  const evidence = await api(`/api/tournaments/${t.id}/evidence`);
  assert.equal(keccak256(stringToHex(evidence.canonical_json)), evidence.evidence_hash, 'Canonical evidence hash mismatch'); assert.equal(evidence.evidence_hash, m.evidence_hash, 'Match/evidence hash mismatch');
  const e = JSON.parse(evidence.canonical_json);
  assert.equal(e.chain_id, CHAIN_ID); assert.equal(lower(e.manager_address), opts.manager); assert.equal(Number(e.tournament_id), Number(t.id)); assert.equal(e.mode, 'llm'); assert.equal(e.outcome, m.outcome); assert.equal(lower(e.winner_address || ''), lower(m.winner_address || ''));
  assert.ok(argsEqual(e.rounds, m.rounds) && argsEqual(e.red_agent, m.red_agent) && argsEqual(e.blue_agent, m.blue_agent), 'Evidence snapshot differs from completed match');
  if (m.winner_address) assert.ok([lower(t.chain.participant1), lower(t.chain.participant2)].includes(lower(m.winner_address)) && m.outcome !== 'draw', 'Winner is not a real participant/outcome');
  else assert.equal(m.outcome, 'draw', 'No winner requires an explicit draw');
  state.llm = proof; state.llm.evidence_hash = evidence.evidence_hash; state.llm.evidence_canonical_hash_verified = true; state.checkpoints.completed_match_verified = nowISO(); save();
  return { t, m, evidence };
}
async function run() {
  let t = await owned('main'); assert.ok(state.checkpoints.prepared, 'prepare is not complete');
  let match = await api(`/api/tournaments/${t.id}/match`, { allow404: true });
  if (match?.status === 'completed') { await completedMatch(); log('Existing completed real match verified; no repeat run/model call.'); return; }
  requireLLM();
  if (Number(t.chain.status) === 0) {
    await waitChain(t.registration_deadline, 'registration closure');
    await transact('main.start', 'organizer', 'startTournament', [t.id]); t = await owned('main');
  }
  assert.equal(Number(t.chain.status), 1, 'Main tournament is not Active; no LLM request permitted');
  assert.ok(await blockTime() + t.review_seconds < t.settlement_deadline, 'Too little settlement time; no model charge attempted');
  await indexed(t.id);
  match = await api(`/api/tournaments/${t.id}/match`, { allow404: true });
  if (match?.status === 'failed') throw new LocalError('Match failed; no automatic new model run. Inspect metadata; use cancel --main-timeout for the honest timeout refund.');
  const needsRequest = !match || match.status === 'interrupted';
  if (needsRequest) {
    if ((match || state.run_requests.length) && !opts.resume) throw new LocalError('Run was interrupted or a prior request is uncertain. No automatic model spend; inspect status and explicitly use --resume-llm to retry once.');
    const intent = { requested_at: nowISO(), reason: match ? 'explicit_interrupted_resume' : state.run_requests.length ? 'explicit_uncertain_retry' : 'first_run', previous_status: match?.status || 'absent', configured_provider: cfg.llm_provider, configured_model: cfg.llm_model, status: 'intent' };
    state.run_requests.push(intent); save();
    const response = await api(`/api/tournaments/${t.id}/run`, { token: await login('organizer'), method: 'POST', body: {}, expected: [202, 409] });
    intent.status = response.status || 'accepted_or_conflict'; intent.responded_at = nowISO(); save();
  }
  const started = Date.now(); let lastLog = 0;
  while (Date.now() - started < opts.wait * 1000) {
    match = await api(`/api/tournaments/${t.id}/match`);
    assert.equal(match.mode, 'llm', 'Backend used simulation; refuse to commit any result');
    if (match.status === 'completed') { await completedMatch(); log(`Real LLM completed (${state.llm.rounds} rounds, ${state.llm.succeeded} successful calls, outcome ${state.llm.outcome}).`); return; }
    if (['interrupted', 'failed'].includes(match.status)) {
      state.llm = { verified: false, status: match.status, mode: match.mode, provider: match.llm_summary?.provider || '', model: match.llm_summary?.model || '', attempted: match.llm_summary?.attempted ?? null, succeeded: match.llm_summary?.succeeded ?? null, failed: match.llm_summary?.failed ?? null, usage: usage(match.llm_summary?.usage), last_error_code: match.llm_summary?.last_error_code || '', rounds: match.rounds?.length ?? 0 }; save();
      throw new LocalError(`Real model run ${match.status}; no fabricated round/winner, no automatic retry. Inspect status; --resume-llm explicitly permits another interrupted attempt.`);
    }
    if (Date.now() - lastLog > 15000) { log(`LLM ${match.status}: round ${match.current_round}, successful calls ${match.llm_summary?.succeeded ?? 'unknown'}`); lastLog = Date.now(); }
    await pause(2000);
  }
  throw new LocalError('LLM wait budget reached; the server may still be running. Rerun run/status with the same state; a running/completed match will not be requested again.');
}
async function moneyOut(which, prize = false) {
  const t = await owned(which), amount = BigInt(t.sponsor_amount_raw), actor = prize ? (lower(t.chain.winner) === lower(accounts.red.address) ? 'red' : 'blue') : 'sponsor';
  if (prize) assert.ok([lower(accounts.red.address), lower(accounts.blue.address)].includes(lower(t.chain.winner)), 'Unrecognized prize recipient');
  assert.equal((await read('getEscrow', [t.id])).totalFunded, amount, 'Unexpected additional sponsorship; reconcile before scripted payout');
  const key = `${which}.${prize ? 'claim' : 'refund'}`;
  state.payouts[key] ||= { actor, amount_raw: amount.toString(), before_raw: (await tokenRead('balanceOf', [accounts[actor].address])).toString() }; save();
  const payout = state.payouts[key];
  const receipt = await transact(key, actor, prize ? 'claimPrize' : 'claimRefund', [t.id]);
  const transfer = receipt.logs.filter((l) => lower(l.address) === deployment.token).map((l) => { try { return decodeEventLog({ abi: tokenABI, data: l.data, topics: l.topics }); } catch { return null; } }).find((e) => e?.eventName === 'Transfer' && lower(e.args.from) === deployment.escrow && lower(e.args.to) === lower(accounts[actor].address) && e.args.value === amount);
  assert.ok(transfer, 'Receipt lacks exact escrow-to-recipient token Transfer');
  if (!payout.verified) {
    payout.after_raw = (await tokenRead('balanceOf', [accounts[actor].address])).toString();
    assert.equal(BigInt(payout.after_raw) - BigInt(payout.before_raw), amount, 'Actual recipient balance change differs from expected payout');
    payout.verified = true; payout.tx_hash = receipt.transactionHash; payout.balance_delta_raw = amount.toString(); save();
  }
  assert.equal(await read(prize ? 'claimablePrize' : 'refundableAmount', [t.id, accounts[actor].address]), 0n, 'Payout still claimable after withdrawal');
  const escrow = await read('getEscrow', [t.id]); assert.equal(escrow.totalFunded, amount); assert.equal(escrow.totalPaid, amount, 'Per-tournament funds not fully paid');
  state.checkpoints[`${which}_paid`] ||= nowISO(); save();
}
async function timeoutMain() {
  let t = await owned('main');
  assert.notEqual(Number(t.chain.status), 4, 'Finalized tournament cannot be refunded; use settle');
  if (Number(t.chain.status) !== 5) {
    const underfilled = Number(t.chain.status) === 0 && [t.chain.participant1, t.chain.participant2].some((a) => lower(a) === '0x0000000000000000000000000000000000000000');
    await waitChain(underfilled ? t.registration_deadline : t.settlement_deadline, underfilled ? 'underfilled registration deadline' : 'actual final settlement deadline');
    await transact('main.expire', 'sponsor', 'expireTournament', [t.id]);
  }
  await moneyOut('main'); state.checkpoints.main_timeout_refund = nowISO(); save();
  log('Main tournament honestly expired and sponsor refunded; no winner was invented.');
}
async function settle() {
  const { m, evidence } = await completedMatch(); let t = await owned('main');
  if (!m.winner_address) { log('Actual draw: no result submission or prize assignment. Waiting for timeout refund.'); await timeoutMain(); return; }
  if (Number(t.chain.status) === 0) throw new LocalError('Tournament was never started');
  if (Number(t.chain.status) === 5) throw new LocalError('Main tournament was cancelled; no result/prize will be committed');
  if (Number(t.chain.status) === 1) {
    await transact('main.submit', 'organizer', 'submitResult', [t.id, m.winner_address, evidence.evidence_hash]); t = await owned('main');
  }
  assert.equal(t.chain.resultHash, evidence.evidence_hash, 'On-chain result hash differs from actual evidence'); assert.equal(lower(t.chain.proposedWinner), lower(m.winner_address), 'On-chain proposed winner mismatch');
  if (Number(t.chain.status) === 3) throw new LocalError('A stakeholder disputed this result; scripted approval is forbidden. Independent review is required.');
  if (Number(t.chain.status) === 2) { await waitChain(Number(t.chain.reviewDeadline), 'full stakeholder review window'); await transact('main.approve', 'auditor', 'approveResult', [t.id]); t = await owned('main'); }
  assert.equal(Number(t.chain.status), 4, 'Main tournament was not Finalized'); assert.equal(lower(t.chain.winner), lower(m.winner_address), 'Audited winner differs from real match');
  await moneyOut('main', true); await observe('after_main_settlement');
  log(`Actual winner ${m.winner_address} claimed the sponsored prize; local balance/receipt verified.`);
}
async function cancel() {
  if (opts.mainTimeout) { await timeoutMain(); return; }
  requireLLM();
  const t = await create('refund', parseUnits(state.parameters.refund_prize, deployment.decimals));
  if (!state.checkpoints.refund_paid) {
    await ensureFunded('refund');
    if (!state.notes.refund_cancellation) {
      await indexed(t.id);
      const note = await api(`/api/tournaments/${t.id}/notes`, { token: await login('organizer'), body: { kind: 'cancellation', text: 'Separate local acceptance scenario: organizer cancels before starting; sponsor receives the full original contribution.' }, expected: [201] });
      assert.equal(keccak256(stringToHex(note.canonical_json)), note.evidence_hash, 'Cancellation note hash mismatch');
      state.notes.refund_cancellation = { evidence_hash: note.evidence_hash, canonical_hash_verified: true }; save();
    }
    await transact('refund.cancel', 'organizer', 'cancelTournament', [t.id, state.notes.refund_cancellation.evidence_hash]);
    assert.equal(Number((await owned('refund')).chain.status), 5);
    await moneyOut('refund'); await observe('after_independent_cancel_refund');
  } else { await owned('refund'); log('Independent refund already complete; no new tournament or payment.'); }
}
async function status() {
  const tournaments = {};
  for (const which of Object.keys(state?.tournaments || {})) {
    if (!state.tournaments[which].id) { tournaments[which] = { status: 'creation_checkpoint_pending' }; continue; }
    const t = await owned(which), escrow = await read('getEscrow', [t.id]);
    const match = which === 'main' ? await api(`/api/tournaments/${t.id}/match`, { allow404: true }) : null;
    tournaments[which] = { id: String(t.id), title: t.title, status: STATUS[Number(t.chain.status)], registration_deadline: t.registration_deadline, settlement_deadline: t.settlement_deadline, review_deadline: Number(t.chain.reviewDeadline), winner: t.chain.winner, total_funded_raw: escrow.totalFunded.toString(), total_paid_raw: escrow.totalPaid.toString(), prize_claimed: escrow.prizeClaimed, sponsor_refundable_raw: (await read('refundableAmount', [t.id, accounts.sponsor.address])).toString(), match: match ? { status: match.status, mode: match.mode, rounds: match.rounds?.length || 0, outcome: match.outcome, summary: match.llm_summary ? { provider: match.llm_summary.provider, model: match.llm_summary.model, attempted: match.llm_summary.attempted, succeeded: match.llm_summary.succeeded, failed: match.llm_summary.failed, usage: usage(match.llm_summary.usage), last_error_code: match.llm_summary.last_error_code || '' } : null } : null };
  }
  console.log(json({ checked_at: nowISO(), rpc: opts.rpc, api: opts.api, chain_id: CHAIN_ID, chain_timestamp: await blockTime(), manager: opts.manager, execution_mode: cfg.execution_mode || 'unknown_old_api', llm_configured: cfg.llm_configured === true, provider: cfg.llm_provider || '', model: cfg.llm_model || '', state_file: statePath, run_id: state?.run_id || null, transactions: state?.transactions?.length || 0, tournaments, balances: (await observe('status_readonly')).balances }));
}
async function init() {
  statePath = opts.state;
  if (existsSync(statePath)) state = JSON.parse(readFileSync(statePath, 'utf8'));
  else assert.ok(['agents', 'prepare', 'cancel', 'status'].includes(opts.phase), 'No state checkpoint: run agents/prepare first');
  if (opts.phase !== 'status') {
    mkdirSync(dirname(statePath), { recursive: true }); lockPath = `${statePath}.lock`;
    try { lockFD = openSync(lockPath, 'wx', 0o600); writeFileSync(lockFD, `${process.pid}\n`); }
    catch { throw new LocalError(`Checkpoint lock exists: ${lockPath}. Stop any prior script and verify it is stopped before removing its stale lock.`); }
  }
  const chain = { id: CHAIN_ID, name: 'Local Course Anvil', nativeCurrency: { name: 'Ether', symbol: 'ETH', decimals: 18 }, rpcUrls: { default: { http: [opts.rpc] } } };
  publicClient = createPublicClient({ chain, transport: http(opts.rpc, { retryCount: 0, timeout: 10000, fetchOptions: { redirect: 'error' } }) });
  assert.equal(await publicClient.getChainId(), CHAIN_ID, 'Refusing non-31337 chain');
  const version = await publicClient.request({ method: 'web3_clientVersion' }); assert.match(version, /anvil/i, 'Refusing non-Anvil node');
  cfg = await api('/api/config'); assert.equal(cfg.chain_id, CHAIN_ID, 'API chain mismatch'); assert.equal(cfg.deployment_configured, true, 'API deployment unavailable'); assert.equal(lower(cfg.manager_address), opts.manager, 'Unexpected manager; explicit local deployment confirmation required');
  assert.ok(await publicClient.getBytecode({ address: opts.manager }), 'Manager bytecode missing');
  deployment = { token: lower(await read('settlementToken')), escrow: lower(await read('prizeEscrow')), registry: lower(await read('tournamentRegistry')), review: lower(await read('resultReview')) };
  for (const [key, field] of [['token', 'token_address'], ['escrow', 'escrow_address'], ['registry', 'registry_address'], ['review', 'review_address']]) assert.equal(deployment[key], lower(cfg[field]), `API ${key} differs from manager`);
  deployment.decimals = Number(await tokenRead('decimals')); assert.equal(deployment.decimals, cfg.token_decimals, 'Token precision mismatch');
  // This fixed public seed is not read from any user's wallet/config/environment.
  const labels = ['admin', 'organizer', 'auditor', 'red', 'blue', 'sponsor'];
  accounts = Object.fromEntries(labels.map((name, addressIndex) => [name, mnemonicToAccount('test test test test test test test test test test test junk', { addressIndex })]));
  if (opts.phase !== 'status') {
    const organizerRole = await read('ORGANIZER_ROLE'), auditorRole = await read('AUDITOR_ROLE');
    assert.equal(await read('hasRole', [organizerRole, accounts.organizer.address]), true, 'Public development organizer lacks role; script never grants roles');
    assert.equal(await read('hasRole', [auditorRole, accounts.auditor.address]), true, 'Public development auditor lacks role; script never grants roles');
  }
  if (state) {
    assert.equal(state.schema_version, 1, 'Unsupported checkpoint format'); assert.equal(state.rpc, opts.rpc, 'Checkpoint RPC mismatch'); assert.equal(state.api, opts.api, 'Checkpoint API mismatch'); assert.equal(state.manager, opts.manager, 'Checkpoint manager mismatch'); assert.equal(state.chain_id, CHAIN_ID);
    const anchor = await publicClient.getBlock({ blockNumber: BigInt(state.chain_anchor.block_number) }); assert.equal(anchor.hash, state.chain_anchor.block_hash, 'Chain reset/fork detected: checkpoint cannot be reused');
    assert.ok(argsEqual(state.deployment, deployment), 'Checkpoint deployment mismatch');
    for (const [key, field] of [['prize', 'prize'], ['refundPrize', 'refund_prize'], ['registration', 'registration_seconds'], ['settlement', 'settlement_seconds'], ['review', 'review_seconds']]) if (opts.explicit.has(key)) assert.equal(String(opts[key]), String(state.parameters[field]), `Cannot alter ${key} on an existing run`);
  } else if (opts.phase !== 'status') {
    const parameters = { prize: String(opts.prize ?? '100'), refund_prize: String(opts.refundPrize ?? '25'), registration_seconds: integer(opts.registration ?? 180, 'registration seconds', 60, 86400), settlement_seconds: integer(opts.settlement ?? 7200, 'settlement seconds', 180, 604800), review_seconds: integer(opts.review ?? 60, 'review seconds', 60, 2592000) };
    for (const amount of [parameters.prize, parameters.refund_prize]) { assert.match(amount, /^\d+(\.\d+)?$/, 'Prize must be a positive decimal without exponents'); assert.ok((amount.split('.')[1]?.length || 0) <= deployment.decimals, 'Prize exceeds token precision'); assert.ok(parseUnits(amount, deployment.decimals) > 0n, 'Zero prize forbidden'); }
    assert.ok(parameters.settlement_seconds > parameters.registration_seconds + parameters.review_seconds + 60, 'Settlement needs sufficient independent review/execution time');
    const anchor = await publicClient.getBlock();
    state = { schema_version: 1, run_id: randomUUID(), created_at: nowISO(), rpc: opts.rpc, api: opts.api, chain_id: CHAIN_ID, manager: opts.manager, deployment, chain_anchor: { block_number: anchor.number.toString(), block_hash: anchor.hash, timestamp: Number(anchor.timestamp) }, initial_next_tournament_id: (await read('nextTournamentId')).toString(), wallets: Object.fromEntries(labels.map((name) => [name, lower(accounts[name].address)])), parameters, agents: {}, tournaments: {}, transactions: [], run_requests: [], time_adjustments: [], snapshots: [], payouts: {}, notes: {}, checkpoints: {}, phases: [], scope: 'Project-local public Anvil accounts; newly created run-owned tournaments only.', limitations: ['Local test funds only; no Sepolia or cloud validation.', 'Scripted signing, not MetaMask browser signing evidence.', 'Provider-call journal and hash integrity are not independent truth/audit.', 'Time warp, if opted in, affects the shared local chain clock.'] };
    save();
  }
}
async function main() {
  opts = options(process.argv.slice(2));
  if (opts.help) { help(); return; }
  if (opts.selfCheck) { await selfCheck(); return; }
  assert.ok(['agents', 'prepare', 'run', 'settle', 'cancel', 'status'].includes(opts.phase), 'Choose agents, prepare, run, settle, cancel or status (--help for commands)');
  assert.ok(!opts.mainTimeout || opts.phase === 'cancel', '--main-timeout only applies to cancel');
  assert.ok(!opts.resume || opts.phase === 'run', '--resume-llm only applies to run');
  await init();
  if (opts.phase === 'status') { await status(); return; }
  const entry = { phase: opts.phase, started_at: nowISO(), status: 'started' }; state.phases.push(entry); save();
  try {
    await ({ agents: profiles, prepare, run, settle, cancel })[opts.phase]();
    entry.status = 'completed'; entry.completed_at = nowISO();
    const confirmed = state.transactions.filter((t) => t.status === 'confirmed');
    state.validation_summary = { confirmed_transactions: confirmed.length, unique_business_transaction_types: [...new Set(confirmed.filter((t) => t.to === opts.manager).map((t) => t.operation))], actual_prize_balance_verified: state.payouts['main.claim']?.verified === true, actual_refund_balance_verified: Object.entries(state.payouts).some(([key, p]) => key.endsWith('.refund') && p.verified), real_llm_evidence_verified: state.llm?.verified === true, time_warp_used: state.time_adjustments.length > 0, model_run_requests: state.run_requests.length };
    save();
    log(`${opts.phase} complete. Evidence/checkpoint: ${statePath}`);
  } catch (e) { entry.status = 'stopped'; entry.stopped_at = nowISO(); entry.reason = safeError(e); save(); throw e; }
}
main().catch((e) => { console.error(`[local-course-flow] ${safeError(e)}`); process.exitCode = 1; }).finally(() => {
  if (lockFD !== undefined) { closeSync(lockFD); try { unlinkSync(lockPath); } catch { /* Never remove anyone else's lock. */ } }
});
