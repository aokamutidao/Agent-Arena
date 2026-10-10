const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const Module = require("node:module");
const ts = require("typescript");
const { keccak256, toBytes, maxUint256, BaseError } = require("viem");
// Compile the actual pure domain helpers in memory, with the installed TypeScript runtime.
const file = path.resolve(__dirname, "../src/lib/course.ts");
const output = ts.transpileModule(fs.readFileSync(file, "utf8"), { compilerOptions: { module: ts.ModuleKind.CommonJS, target: ts.ScriptTarget.ES2020 } }).outputText;
const loaded = new Module(file);
loaded.filename = file;
loaded.paths = Module._nodeModulePaths(path.dirname(file));
loaded._compile(output, file);
const { parseMoney, money, registrationOpen, reviewComplete, canExpire, verifyEvidence, receiptVerdict, dateTime, errorMessage, ApiError } = loaded.exports;

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

test("English dates use British day-month ordering without inventing a missing timestamp", () => {
  const englishDate = new Intl.DateTimeFormat("en-GB", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false,
  });
  for (const seconds of [
    Date.UTC(2026, 9, 23, 0, 5, 9) / 1000,
    Date.UTC(2026, 2, 4, 12, 30, 45) / 1000,
  ]) {
    assert.equal(dateTime(seconds, "en"), englishDate.format(new Date(seconds * 1000)));
  }
  assert.equal(dateTime(0, "en"), "—");
  assert.equal(dateTime(0, "zh"), "—");
});

test("English amount-validation errors retain token precision and contract limits", () => {
  for (const decimals of [0, 6, 18, 36]) {
    assert.throws(() => parseMoney("-1", decimals), (error) => {
      const original = error.message;
      const translated = errorMessage(error, "en");
      assert.match(translated, /non-negative amount/i);
      assert.match(translated, new RegExp(`at most ${decimals} decimal places`));
      assert.doesNotMatch(translated, /\p{Script=Han}/u);
      assert.equal(errorMessage(error, "zh"), original);
      assert.equal(error.message, original);
      return true;
    });
  }
  assert.throws(() => parseMoney("1", 37), (error) => {
    assert.match(errorMessage(error, "en"), /Invalid token decimal configuration/);
    return true;
  });
  assert.throws(() => parseMoney((maxUint256 + 1n).toString(), 0), (error) => {
    assert.match(errorMessage(error, "en"), /exceeds the contract limit/);
    return true;
  });
  // Translating a validation failure must not change how valid financial values are parsed.
  assert.equal(parseMoney("9007199254740993.123456"), 9007199254740993123456n);
});

test("English local errors keep signing, replacement, and revert outcomes distinct", () => {
  const cases = [
    ["请先连接钱包。", /Connect your wallet first/],
    ["签名期间钱包或网络已变化，请重新登录。", /changed while signing.*Sign in again/],
    ["钱包或网络已变化，请重新估算。", /wallet or network has changed.*Estimate.*again/],
    ["目标地址没有合约代码，请检查部署配置。", /no contract.*deployment configuration/],
    ["服务未返回有效 JSON，请检查 API 地址。", /invalid JSON.*API address/],
    ["确认数不足，请稍后继续检查。", /Not enough confirmations.*Check again/],
    ["原交易已被取消或替换为其他操作，请核查新交易回执后重试。", /cancelled or replaced.*replacement receipt.*trying again/],
    ["交易已上链但执行回滚；请查看交易回执。", /included on-chain but reverted.*receipt/],
  ];
  for (const [message, meaning] of cases) {
    const error = Object.freeze(new Error(message));
    assert.match(errorMessage(error, "en"), meaning);
    assert.doesNotMatch(errorMessage(error, "en"), /\p{Script=Han}/u);
    assert.equal(errorMessage(error, "zh"), message);
    assert.equal(error.message, message);
  }
  assert.match(errorMessage(null, "en"), /Request failed.*connection/);
  assert.equal(errorMessage(new Error(""), "en"), "Request failed.");
  const pending = errorMessage("交易状态待确认：确认数不足，请稍后继续检查。", "en");
  assert.match(pending, /Transaction status unconfirmed.*Not enough confirmations/);
  assert.doesNotMatch(pending, /\p{Script=Han}/u);
});

test("English HTTP errors retain status codes and unknown backend diagnostics verbatim", () => {
  for (const status of [401, 429, 503]) {
    const error = Object.freeze(new ApiError(`服务返回 HTTP ${status}`, status));
    assert.equal(errorMessage(error, "en"), `The service returned HTTP ${status}.`);
    assert.equal(error.status, status);
    assert.equal(error.message, `服务返回 HTTP ${status}`);
  }
  const detail = "比赛证据暂不可用：tournament=17, evidence=0x123";
  const error = Object.freeze(new ApiError(detail, 503));
  assert.equal(errorMessage(error, "en"), detail);
  assert.equal(errorMessage(error, "zh"), detail);
  assert.equal(error.status, 503);
  assert.equal(error.message, detail);
  assert.equal(errorMessage("RPC unavailable: request-id=42", "en"), "RPC unavailable: request-id=42");
});

test("English provider errors use viem shortMessage without replacing diagnostic detail", () => {
  const error = new BaseError("User rejected the request.", {
    details: "request-id=wallet-42",
  });
  const original = { message: error.message, shortMessage: error.shortMessage, details: error.details };
  Object.freeze(error);
  assert.notEqual(error.message, error.shortMessage);
  assert.equal(errorMessage(error, "en"), "User rejected the request.");
  assert.equal(errorMessage(error, "zh"), "User rejected the request.");
  assert.deepEqual(
    { message: error.message, shortMessage: error.shortMessage, details: error.details },
    original,
  );
});
