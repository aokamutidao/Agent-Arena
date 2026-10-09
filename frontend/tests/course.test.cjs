const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const Module = require("node:module");
const ts = require("typescript");
const { keccak256, toBytes, maxUint256 } = require("viem");
// Compile the actual pure domain helpers in memory, with the installed TypeScript runtime.
const file = path.resolve(__dirname, "../src/lib/course.ts");
const output = ts.transpileModule(fs.readFileSync(file, "utf8"), { compilerOptions: { module: ts.ModuleKind.CommonJS, target: ts.ScriptTarget.ES2020 } }).outputText;
const loaded = new Module(file);
loaded.filename = file;
loaded.paths = Module._nodeModulePaths(path.dirname(file));
loaded._compile(output, file);
const { parseMoney, money, registrationOpen, reviewComplete, canExpire, verifyEvidence, receiptVerdict } = loaded.exports;

test("token amounts preserve six-decimal integers beyond JavaScript's safe integer range", () => {
  assert.equal(parseMoney("9007199254740993.123456"), 9007199254740993123456n);
  assert.equal(money("9007199254740993123456"), "9007199254740993.123456");
  assert.equal(parseMoney("0.000001"), 1n);
});
test("invalid, excess-precision, exponential and overflowing amounts are rejected", () => {
  for (const input of ["", "-1", "1e6", "NaN", "1.0000001", "1,000", "0x10", "Infinity"]) assert.throws(() => parseMoney(input));
  assert.throws(() => parseMoney((maxUint256 + 1n).toString(), 0));
});
const base = { status: "open", registration_deadline: 100, settlement_deadline: 500, review_deadline: 200, participants: ["a", "b"] };
test("registration and review respect exact deadline boundaries", () => {
  assert.equal(registrationOpen(base, 99), true);
  assert.equal(registrationOpen(base, 100), false);
  assert.equal(reviewComplete(base, 199), false);
  assert.equal(reviewComplete(base, 200), true);
  assert.equal(reviewComplete(base, 500), false);
});
test("timeouts protect underfilled events and cannot cancel terminal events", () => {
  assert.equal(canExpire(base, 100), false);
  assert.equal(canExpire({ ...base, participants: ["a"] }, 100), true);
  assert.equal(canExpire({ ...base, status: "running" }, 499), false);
  assert.equal(canExpire({ ...base, status: "running" }, 500), true);
  assert.equal(canExpire({ ...base, status: "finalized" }, 500), false);
  assert.equal(canExpire({ ...base, status: "cancelled" }, 500), false);
});
test("evidence verification rejects changed actions and whitespace", () => {
  const canonical = '{"winner":"0x123","rounds":[{"action":"MOVE"}]}';
  const hash = keccak256(toBytes(canonical));
  assert.equal(verifyEvidence(canonical, hash), true);
  assert.equal(verifyEvidence(canonical.replace("MOVE", "ATTACK"), hash), false);
  assert.equal(verifyEvidence(`${canonical}\n`, hash), false);
});
test("a successful cancellation receipt does not count as a successful financial action", () => {
  assert.equal(receiptVerdict("success", "cancelled"), "replaced");
  assert.equal(receiptVerdict("success", "replaced"), "replaced");
  assert.equal(receiptVerdict("success", "repriced"), "success");
  assert.equal(receiptVerdict("reverted"), "reverted");
});
