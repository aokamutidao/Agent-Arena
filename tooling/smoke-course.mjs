// Real local-EVM + HTTP integration. Never connects to an existing or public RPC.
// `--serve` keeps the verified local demo alive; Ctrl+C only stops this script's children.
import assert from 'node:assert/strict';
import { spawn } from 'node:child_process';
import { createRequire } from 'node:module';
import { mkdtempSync, readFileSync, rmSync, openSync, closeSync, existsSync, mkdirSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { randomBytes } from 'node:crypto';
import net from 'node:net';

const root = resolve(fileURLToPath(new URL('..', import.meta.url)));
const require = createRequire(join(root, 'frontend/package.json'));
const { createPublicClient, createWalletClient, http, keccak256, stringToHex, parseUnits, zeroAddress } = require('viem');
const { mnemonicToAccount } = require('viem/accounts');
const keep = process.argv.includes('--serve');
const work = mkdtempSync(join(tmpdir(), 'agent-arena-smoke-'));
const children = [];
const logDescriptors = [];
const gas = [];
let rpcURL, apiURL, backend, backendOptions, miningTimer;
const uiPort = Number(process.env.COURSE_SMOKE_UI_PORT || 3000);
const uiOrigin = `http://localhost:${uiPort}`;
const pause = (ms) => new Promise((r) => setTimeout(r, ms));
const accounts = Array.from({ length: 7 }, (_, addressIndex) => mnemonicToAccount(
  'test test test test test test test test test test test junk', { addressIndex },
)); // Public Anvil mnemonic, for the new loopback-only test chain.
const [admin, organizer, auditor, red, blue, sponsor] = accounts;
const artifact = (name) => JSON.parse(readFileSync(join(root, `contracts/out/${name}.sol/${name}.json`), 'utf8'));
const managerABI = JSON.parse(readFileSync(join(root, 'contracts/abi/TournamentManager.json'), 'utf8'));
const tokenABI = JSON.parse(readFileSync(join(root, 'contracts/abi/CourseToken.json'), 'utf8'));

async function freePort(preferred = 0) {
  const server = net.createServer();
  await new Promise((yes, no) => server.once('error', no).listen(preferred, '127.0.0.1', yes));
  const port = server.address().port;
  await new Promise((r) => server.close(r));
  return port;
}
function child(command, args, options = {}) {
  const fd = openSync(join(work, `${children.length}.log`), 'a'); logDescriptors.push(fd);
  const process = spawn(command, args, { cwd: root, stdio: ['ignore', fd, fd], ...options });
  process.on('error', (e) => { process.spawnError = e; });
  children.push(process); return process;
}
async function stop(process) {
  if (!process || process.exitCode !== null || process.signalCode !== null) return;
  process.kill('SIGTERM');
  for (let i = 0; i < 40 && process.exitCode === null && process.signalCode === null; i++) await pause(50);
  if (process.exitCode === null && process.signalCode === null) process.kill('SIGKILL');
}
async function until(check, label, attempts = 100) {
  let last;
  for (let i = 0; i < attempts; i++) {
    try { const value = await check(); if (value) return value; } catch (e) { last = e; }
    await pause(100);
  }
  throw new Error(`Timed out: ${label}${last ? ` (${last.message})` : ''}`);
}
async function api(path, { token, body, expected = 200, method } = {}) {
  const response = await fetch(`${apiURL}${path}`, {
    method: method || (body === undefined ? 'GET' : 'POST'),
    headers: { 'Content-Type': 'application/json', ...(token ? { Authorization: `Bearer ${token}` } : {}) },
    ...(body === undefined ? {} : { body: JSON.stringify(body) }), signal: AbortSignal.timeout(10000),
  });
  const value = await response.json();
  assert.equal(response.status, expected, `${path}: ${JSON.stringify(value)}`);
  return value;
}
async function login(account) {
  const challenge = await api('/api/auth/nonce', { body: { address: account.address } });
  const body = { address: account.address, nonce: challenge.nonce, signature: await account.signMessage({ message: challenge.message }) };
  const session = await api('/api/auth/login', { body });
  await api('/api/auth/login', { body, expected: 401 });
  return session;
}

try {
  const rpcPort = await freePort(keep ? Number(process.env.COURSE_SMOKE_RPC_PORT || 8545) : 0);
  const apiPort = await freePort(keep ? Number(process.env.COURSE_SMOKE_API_PORT || 8080) : 0);
  rpcURL = `http://127.0.0.1:${rpcPort}`; apiURL = `http://127.0.0.1:${apiPort}`;
  const anvil = child('anvil', ['--host', '127.0.0.1', '--port', String(rpcPort), '--chain-id', '31337', '--silent']);
  const chain = { id: 31337, name: 'Course Anvil', nativeCurrency: { name: 'Ether', symbol: 'ETH', decimals: 18 }, rpcUrls: { default: { http: [rpcURL] } } };
  const publicClient = createPublicClient({ chain, transport: http(rpcURL, { retryCount: 0 }) });
  await until(async () => {
    if (anvil.spawnError) throw anvil.spawnError;
    assert.equal(await publicClient.getChainId(), 31337); return true;
  }, 'new local Anvil chain');
  // Assert this is our fresh chain; never use a pre-existing server on --serve's port.
  assert.equal(await publicClient.getBlockNumber({ cacheTime: 0 }), 0n, 'refusing an existing chain');
  assert.equal(anvil.exitCode, null, 'Anvil failed to bind');
  const wallet = (a) => createWalletClient({ account: a, chain, transport: http(rpcURL) });
  async function receipt(hash, label) {
    const r = await publicClient.waitForTransactionReceipt({ hash });
    assert.equal(r.status, 'success', label); gas.push({ operation: label, gas_used: r.gasUsed.toString(), tx_hash: hash }); return r;
  }
  async function deploy(name, args) {
    const a = artifact(name);
    return (await receipt(await wallet(admin).deployContract({ abi: a.abi, bytecode: a.bytecode.object, args }), `deploy ${name}`)).contractAddress;
  }
  const tokenAddress = await deploy('CourseToken', [admin.address]);
  const manager = await deploy('TournamentManager', [tokenAddress, admin.address]);
  const deploymentBlock = Number(await publicClient.getBlockNumber({ cacheTime: 0 }));
  const read = (name, args = []) => publicClient.readContract({ address: manager, abi: managerABI, functionName: name, args });
  const write = async (actor, name, args = []) => receipt(await wallet(actor).writeContract({ address: manager, abi: managerABI, functionName: name, args }), name);
  const tokenWrite = async (actor, name, args) => receipt(await wallet(actor).writeContract({ address: tokenAddress, abi: tokenABI, functionName: name, args }), name);
  await write(admin, 'grantRole', [keccak256(stringToHex('ORGANIZER_ROLE')), organizer.address]);
  await write(admin, 'grantRole', [keccak256(stringToHex('AUDITOR_ROLE')), auditor.address]);
  const escrow = await read('prizeEscrow');
  const balance = (address) => publicClient.readContract({ address: tokenAddress, abi: tokenABI, functionName: 'balanceOf', args: [address] });
  await tokenWrite(admin, 'mint', [sponsor.address, parseUnits('1000', 6)]);
  await tokenWrite(sponsor, 'approve', [escrow, parseUnits('1000', 6)]);

  const binary = join(root, 'backend/tmp/course-server');
  assert.ok(existsSync(binary), 'Build backend/tmp/course-server before running this integration test');
  backendOptions = { cwd: join(root, 'backend'), env: {
    ...process.env, JWT_SECRET: randomBytes(32).toString('hex'), DATABASE_PATH: join(work, 'course.db'),
    AUTH_DOMAIN: `localhost:${uiPort}`, CORS_ORIGINS: uiOrigin, PORT: String(apiPort),
    CHAIN_ID: '31337', RPC_URL: rpcURL, TOURNAMENT_MANAGER_ADDRESS: manager,
    DEPLOYMENT_BLOCK: String(deploymentBlock), CONFIRMATIONS: '0', CONTRACT_ABI_DIR: join(root, 'contracts/abi'),
    SIMULATION_ENABLED: 'true', QWEN_API_KEY: '', ROUND_INTERVAL_MS: '1', POLL_INTERVAL_SECONDS: '1',
    PRIVATE_KEY: '', OWNER_PRIVATE_KEY: '', COURSE_DEPLOYER_KEY: '',
  } };
  backend = child(binary, [], backendOptions);
  await until(() => api('/health'), 'backend');
  const config = await api('/api/config'); assert.equal(config.deployment_configured, true);
  assert.equal(config.manager_address.toLowerCase(), manager.toLowerCase());
  const orgSession = await login(organizer), redSession = await login(red), blueSession = await login(blue), auditSession = await login(auditor);
  assert.ok(orgSession.user.roles.includes('organizer')); assert.ok(auditSession.user.roles.includes('auditor'));
  assert.ok(!redSession.user.roles.includes('organizer'));
  const redAgent = await api('/api/agents', { token: redSession.token, body: { name: 'Aggressive Scholar', personality: 'aggressive' }, expected: 201 });
  const blueAgent = await api('/api/agents', { token: blueSession.token, body: { name: 'Patient Guardian', personality: 'defensive' }, expected: 201 });
  const blockTime = async () => Number((await publicClient.getBlock()).timestamp);
  const advanceTo = async (timestamp) => {
    await publicClient.request({ method: 'evm_setNextBlockTimestamp', params: [timestamp] });
    await publicClient.request({ method: 'evm_mine', params: [] });
  };
  async function create(title) {
    const now = await blockTime(), id = await read('nextTournamentId');
    await write(organizer, 'createTournament', [title, '', auditor.address, now + 120, now + 7200, 60]); return id;
  }
  const flowStarted = performance.now();
  const tournamentID = await create('Course integration: audited prize');
  await write(sponsor, 'fundTournament', [tournamentID, parseUnits('100', 6)]);
  await write(red, 'registerParticipant', [tournamentID, redAgent.agent_hash]);
  await write(blue, 'registerParticipant', [tournamentID, blueAgent.agent_hash]);
  const registered = await read('getTournament', [tournamentID]);
  await advanceTo(Number(registered.registrationDeadline));
  await write(organizer, 'startTournament', [tournamentID]);
  await api(`/api/tournaments/${tournamentID}/run`, { token: redSession.token, method: 'POST', body: {}, expected: 403 });
  await api(`/api/tournaments/${tournamentID}/run`, { token: orgSession.token, method: 'POST', body: {}, expected: 202 });
  const match = await until(async () => {
    const m = await api(`/api/tournaments/${tournamentID}/match`);
    assert.notEqual(m.status, 'failed', m.error); return m.status === 'completed' ? m : false;
  }, 'completed deterministic match', 600);
  const evidence = await api(`/api/tournaments/${tournamentID}/evidence`);
  assert.equal(keccak256(stringToHex(evidence.canonical_json)), evidence.evidence_hash);
  assert.ok(evidence.evidence.rounds.length > 0);
  assert.equal(evidence.evidence.red_agent.owner_address.toLowerCase(), red.address.toLowerCase());
  await api(`/api/tournaments/${tournamentID}/run`, { token: orgSession.token, body: {}, expected: 409 });
  if (match.winner_address) {
    await write(organizer, 'submitResult', [tournamentID, match.winner_address, evidence.evidence_hash]);
    const pending = await read('getTournament', [tournamentID]);
    await advanceTo(Number(pending.reviewDeadline));
    await write(auditor, 'approveResult', [tournamentID]);
    const winner = match.winner_address.toLowerCase() === red.address.toLowerCase() ? red : blue;
    const winnerBefore = await balance(winner.address);
    await write(winner, 'claimPrize', [tournamentID]);
    assert.equal((await balance(winner.address)) - winnerBefore, parseUnits('100', 6));
    assert.equal(await read('claimablePrize', [tournamentID, winner.address]), 0n);
  } else {
    // Draws remain real draws; the immutable rules provide a timeout refund path.
    await advanceTo(Number(registered.settlementDeadline));
    await write(sponsor, 'expireTournament', [tournamentID]);
    const sponsorBefore = await balance(sponsor.address);
    await write(sponsor, 'claimRefund', [tournamentID]);
    assert.equal((await balance(sponsor.address)) - sponsorBefore, parseUnits('100', 6));
  }

  // Independent scenario verifies an auditor's dispute cancellation and exact refund.
  const disputedID = await create('Course integration: disputed result');
  await write(sponsor, 'fundTournament', [disputedID, parseUnits('25', 6)]);
  await write(red, 'registerParticipant', [disputedID, redAgent.agent_hash]);
  await write(blue, 'registerParticipant', [disputedID, blueAgent.agent_hash]);
  await advanceTo(Number((await read('getTournament', [disputedID])).registrationDeadline));
  await write(organizer, 'startTournament', [disputedID]);
  await write(organizer, 'submitResult', [disputedID, red.address, keccak256(stringToHex('fixture evidence: contract-only dispute scenario'))]);
  const dispute = await api(`/api/tournaments/${disputedID}/notes`, { token: blueSession.token, body: { kind: 'dispute', text: 'Contract-only test fixture: contest this result and request sponsor refunds.' }, expected: 201 });
  assert.equal(keccak256(stringToHex(dispute.canonical_json)), dispute.evidence_hash);
  await write(blue, 'disputeResult', [disputedID, dispute.evidence_hash]);
  await advanceTo(Number((await read('getTournament', [disputedID])).reviewDeadline));
  const decision = await api(`/api/tournaments/${disputedID}/notes`, { token: auditSession.token, body: { kind: 'decision', text: 'Contract-only test fixture: result rejected; refund the sponsors.' }, expected: 201 });
  assert.equal(keccak256(stringToHex(decision.canonical_json)), decision.evidence_hash);
  await write(auditor, 'resolveDispute', [disputedID, zeroAddress, decision.evidence_hash]);
  assert.equal(await read('refundableAmount', [disputedID, sponsor.address]), parseUnits('25', 6));
  const sponsorBeforeRefund = await balance(sponsor.address);
  await write(sponsor, 'claimRefund', [disputedID]);
  assert.equal((await balance(sponsor.address)) - sponsorBeforeRefund, parseUnits('25', 6));
  assert.equal(await balance(escrow), 0n);
  assert.equal(await publicClient.readContract({ address: escrow, abi: artifact('PrizeEscrow').abi, functionName: 'totalLiabilities' }), 0n);
  assert.equal(await read('refundableAmount', [disputedID, sponsor.address]), 0n);
  await until(async () => (await api('/ready')).ready, 'indexer readiness');
  const targetBlock = Number(await publicClient.getBlockNumber({ cacheTime: 0 }));
  const before = await until(async () => { const d = await api('/api/dashboard'); return d.indexed_block >= targetBlock ? d : false; }, 'confirmed ledger');
  assert.equal(before.escrow_in_raw, '125000000');
  assert.equal(BigInt(before.prizes_paid_raw) + BigInt(before.refunds_paid_raw), 125000000n);
  const flowMilliseconds = performance.now() - flowStarted;
  const latency = {};
  for (const path of ['/api/dashboard', `/api/tournaments/${tournamentID}/evidence`]) {
    const samples = [];
    for (let i = 0; i < 30; i++) { const started = performance.now(); await api(path); samples.push(performance.now() - started); }
    samples.sort((a, b) => a - b);
    latency[path] = { samples: samples.length, concurrency: 1, p50_ms: Number(samples[14].toFixed(3)), p95_ms: Number(samples[28].toFixed(3)) };
  }
  await stop(backend); backend = child(binary, [], backendOptions);
  await until(() => api('/health'), 'backend restart');
  await until(async () => (await api('/ready')).ready, 'restarted indexer readiness');
  const after = await api('/api/dashboard');
  for (const key of ['prizes_paid_raw', 'refunds_paid_raw', 'ledger_entries']) assert.equal(after[key], before[key]);
  assert.equal(after.escrow_in_raw, before.escrow_in_raw); assert.equal(after.event_count, before.event_count);
  assert.equal((await api(`/api/tournaments/${tournamentID}/evidence`)).evidence_hash, evidence.evidence_hash);
  const notes = await api(`/api/tournaments/${disputedID}/notes`);
  assert.ok(notes.notes.some((n) => n.evidence_hash === dispute.evidence_hash));
  assert.ok(notes.notes.some((n) => n.evidence_hash === decision.evidence_hash));
  const events = await api(`/api/tournaments/${disputedID}/events`);
  assert.ok(JSON.stringify(events).toLowerCase().includes(dispute.evidence_hash.toLowerCase().slice(2)));
  assert.ok(JSON.stringify(events).toLowerCase().includes(decision.evidence_hash.toLowerCase().slice(2)));
  mkdirSync(join(root, '.artifacts'), { recursive: true });
  const result = { checked_at: new Date().toISOString(), chain_id: 31337, manager, mode: 'deterministic_simulation', outcome: match.outcome, nonce_replay_rejected: true, unauthorized_run_rejected: true, evidence_verified: true, actual_erc20_balances_verified: true, escrow_zero_liability_verified: true, dispute_refund_verified: true, review_notes_verified: true, completed_rerun_rejected: true, restart_verified: true, flow_ms: Number(flowMilliseconds.toFixed(1)), local_http_latency: latency, measurement_limits: 'Loopback, 30 sequential samples per endpoint, deterministic simulation, Anvil time-warp; not a production load test or real-world payout duration.', dashboard: before, transactions: gas };
  writeFileSync(join(root, '.artifacts/course-smoke.json'), JSON.stringify(result, null, 2));
  console.log(JSON.stringify({ status: 'passed', transactions: gas.length, outcome: match.outcome, evidence: evidence.evidence_hash, artifact: '.artifacts/course-smoke.json' }, null, 2));
  if (keep) {
    const demoID = await create('学习演示：Agent 奖金托管');
    await write(sponsor, 'fundTournament', [demoID, parseUnits('200', 6)]);
    await write(red, 'registerParticipant', [demoID, redAgent.agent_hash]);
    await write(blue, 'registerParticipant', [demoID, blueAgent.agent_hash]);
    miningTimer = setInterval(() => { void publicClient.request({ method: 'evm_mine', params: [] }).catch(() => {}); }, 1000);
    console.log(JSON.stringify({ process_id: process.pid, ui_origin: uiOrigin, rpc: rpcURL, api: apiURL, demo_tournament: String(demoID), wallets: { admin: admin.address, organizer: organizer.address, auditor: auditor.address, participant1: red.address, participant2: blue.address, sponsor: sponsor.address }, note: 'Local public Anvil accounts only. Start frontend separately. Ctrl+C stops this isolated stack.' }, null, 2));
    await new Promise((r) => { process.once('SIGINT', r); process.once('SIGTERM', r); });
  }
} catch (error) {
  console.error(error);
  for (let i = 0; i < children.length; i++) {
    try { console.error(readFileSync(join(work, `${i}.log`), 'utf8').slice(-6000)); } catch { /* missing spawn log */ }
  }
  process.exitCode = 1;
} finally {
  if (miningTimer) clearInterval(miningTimer);
  for (const process of children.reverse()) await stop(process);
  for (const fd of logDescriptors) closeSync(fd);
  rmSync(work, { recursive: true, force: true });
}
