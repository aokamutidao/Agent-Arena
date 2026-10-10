const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const Module = require("node:module");
const ts = require("typescript");

// Compile the real pure presentation module, as the existing domain-helper tests do.
const file = path.resolve(__dirname, "../src/lib/course-presentation.ts");
const output = ts.transpileModule(fs.readFileSync(file, "utf8"), {
  compilerOptions: {
    module: ts.ModuleKind.CommonJS,
    target: ts.ScriptTarget.ES2020,
  },
}).outputText;
const loaded = new Module(file);
loaded.filename = file;
loaded.paths = Module._nodeModulePaths(path.dirname(file));
loaded._compile(output, file);
const { getTournamentPresentation, globalRoleLabel, getEventActor } =
  loaded.exports;

const base = {
  status: "open",
  registration_deadline: 100,
  settlement_deadline: 500,
  review_deadline: 200,
  review_period: 60,
  participants: ["a", "b"],
  prize_pool: "100000000",
  prize_claimed: false,
};

test("the registration boundary changes a full funded event to waiting for its organizer", () => {
  assert.equal(getTournamentPresentation(base, 99).label, "报名中");
  const deadline = getTournamentPresentation(base, 100);
  assert.equal(deadline.label, "等待开赛");
  assert.match(deadline.hint, /主办方/);
  const underfilled = getTournamentPresentation(
    { ...base, participants: ["a"] },
    100,
  );
  assert.match(underfilled.label, /可取消/);
  assert.match(underfilled.hint, /取消.*随后.*退款/);
});

test("every nonterminal event expires at the exact final deadline, with cancellation before refunds", () => {
  assert.equal(
    getTournamentPresentation({ ...base, status: "running" }, 499).label,
    "对战中",
  );
  for (const status of ["open", "running", "result_submitted", "disputed"]) {
    const expired = getTournamentPresentation({ ...base, status }, 500);
    assert.equal(expired.label, "已到期 · 待退款");
    assert.equal(expired.tone, "danger");
    assert.match(expired.hint, /先.*取消.*确认取消后.*退款/);
  }
});

test("settled and cancelled events do not become expired or pretend all money has been refunded", () => {
  const unclaimed = getTournamentPresentation(
    { ...base, status: "finalized" },
    1000,
  );
  assert.match(unclaimed.label, /待领奖/);
  const claimed = getTournamentPresentation(
    { ...base, status: "finalized", prize_claimed: true },
    1000,
  );
  assert.equal(claimed.label, "已领奖");
  const cancelled = getTournamentPresentation(
    { ...base, status: "cancelled" },
    1000,
  );
  assert.equal(cancelled.label, "已取消");
  assert.match(cancelled.hint, /剩余/);
  assert.match(cancelled.hint, /不代表所有退款已经领取/);
});

test("unknown chain time never invents a deadline or an available time-sensitive action", () => {
  for (const now of [0, -1, NaN, Infinity]) {
    for (const status of ["open", "running", "result_submitted", "disputed"]) {
      const waiting = getTournamentPresentation({ ...base, status }, now);
      assert.match(waiting.hint, /等待网络时间/);
      assert.doesNotMatch(
        waiting.label,
        /可取消|等待开赛|待退款|等待裁决|等待审核/,
      );
    }
  }
  assert.equal(
    getTournamentPresentation(
      { ...base, status: "finalized", prize_claimed: true },
      0,
    ).label,
    "已领奖",
  );
});

test("a zero or invalid prize pool cannot advertise an available start, while huge integers retain value", () => {
  for (const prize_pool of ["0", "000", "invalid"]) {
    const unfunded = getTournamentPresentation({ ...base, prize_pool }, 100);
    assert.match(unfunded.label, /无法开赛/);
    assert.match(unfunded.hint, /报名已截止/);
    assert.doesNotMatch(unfunded.hint, /继续赞助|补充赞助/);
  }
  assert.match(
    getTournamentPresentation({ ...base, prize_pool: "0" }, 99).hint,
    /开赛前需要非零赞助/,
  );
  assert.equal(
    getTournamentPresentation(
      { ...base, prize_pool: "9007199254740993123456" },
      100,
    ).label,
    "等待开赛",
  );
});

test("auditor decisions wait for the full review window", () => {
  assert.equal(
    getTournamentPresentation({ ...base, status: "result_submitted" }, 199)
      .label,
    "审核窗口中",
  );
  assert.equal(
    getTournamentPresentation({ ...base, status: "result_submitted" }, 200)
      .label,
    "等待审核",
  );
  assert.equal(
    getTournamentPresentation({ ...base, status: "disputed" }, 199).label,
    "存在异议",
  );
  assert.equal(
    getTournamentPresentation({ ...base, status: "disputed" }, 200).label,
    "等待裁决",
  );
});

test("starting must leave strictly more than the full review period before settlement", () => {
  assert.equal(getTournamentPresentation(base, 439).label, "等待开赛");
  for (const now of [440, 499]) {
    const tooLate = getTournamentPresentation(base, now);
    assert.match(tooLate.label, /审核时间不足.*无法开赛/);
    assert.doesNotMatch(tooLate.label, /已到期|待退款/);
    assert.match(tooLate.hint, /主办方可取消.*最终期限到达后.*任意钱包取消/);
  }
  assert.equal(getTournamentPresentation(base, 500).label, "已到期 · 待退款");
});

test("the default global role is an ordinary user, not a registered tournament participant", () => {
  assert.equal(globalRoleLabel([]), "普通用户");
  assert.equal(globalRoleLabel(["participant"]), "普通用户");
  assert.equal(globalRoleLabel(["unknown", "participant"]), "普通用户");
  assert.equal(
    globalRoleLabel(["auditor", "participant", "admin", "organizer", "admin"]),
    "管理员 · 主办方 · 审核员",
  );
});

test("English presentation preserves severity and on-chain data across all business states", () => {
  const scenarios = [
    ["registration with a place available", { participants: ["a"] }, 99],
    ["full registration", {}, 99],
    ["empty pool before registration closes", { prize_pool: "0" }, 99],
    ["underfilled registration", { participants: ["a"] }, 100],
    ["unexpected participant count", { participants: ["a", "b", "c"] }, 100],
    ["ready to start", {}, 100],
    ["last available start time", {}, 439],
    ["insufficient review time", {}, 440],
    ["match in progress", { status: "running" }, 150],
    ["open review window", { status: "result_submitted" }, 199],
    ["approval pending", { status: "result_submitted" }, 200],
    ["dispute during review", { status: "disputed" }, 199],
    ["resolution pending", { status: "disputed" }, 200],
    ["settled without a claim", { status: "finalized" }, 1000],
    ["claimed prize", { status: "finalized", prize_claimed: true }, 0],
    ["cancelled with refunds to claim", { status: "cancelled" }, 1000],
    ["large exact prize pool", { prize_pool: "9007199254740993123456" }, 100],
  ];
  for (const prize_pool of ["0", "000", "invalid"]) {
    scenarios.push([`unfunded pool ${prize_pool}`, { prize_pool }, 100]);
  }
  for (const status of ["open", "running", "result_submitted", "disputed"]) {
    scenarios.push([`${status} at final deadline`, { status }, 500]);
    for (const now of [0, -1, NaN, Infinity]) {
      scenarios.push([`${status} with unknown chain time ${now}`, { status }, now]);
    }
  }

  for (const [name, overrides, now] of scenarios) {
    const tournament = {
      ...base,
      ...overrides,
      participants: [...(overrides.participants || base.participants)],
    };
    const before = JSON.parse(JSON.stringify(tournament));
    Object.freeze(tournament.participants);
    Object.freeze(tournament);
    const zh = getTournamentPresentation(tournament, now, "zh");
    const en = getTournamentPresentation(tournament, now, "en");
    assert.doesNotMatch(`${en.label} ${en.hint}`, /\p{Script=Han}/u, name);
    assert.match(en.label, /[A-Za-z]/, name);
    assert.match(en.hint, /[A-Za-z]/, name);
    assert.equal(en.tone, zh.tone, `${name}: language must not change severity`);
    assert.deepEqual(tournament, before, `${name}: presentation must not change financial state`);
  }
});

test("English distinguishes pending claims from completed transfers and requires cancellation before refunds", () => {
  const settled = getTournamentPresentation({ ...base, status: "finalized" }, 1000, "en");
  assert.match(settled.label, /Ready to claim/);
  assert.match(settled.hint, /winner.*own wallet.*claim/i);
  assert.doesNotMatch(settled.label, /claimed/i);
  const claimed = getTournamentPresentation(
    { ...base, status: "finalized", prize_claimed: true },
    1000,
    "en",
  );
  assert.match(claimed.label, /Prize claimed/);

  const cancelled = getTournamentPresentation({ ...base, status: "cancelled" }, 1000, "en");
  assert.match(cancelled.hint, /can claim.*remaining refundable contribution/i);
  assert.match(cancelled.hint, /does not mean all refunds have been claimed/i);
  assert.doesNotMatch(cancelled.label, /refunded|paid|claimed/i);
  for (const status of ["open", "running", "result_submitted", "disputed"]) {
    const expired = getTournamentPresentation({ ...base, status }, 500, "en");
    assert.match(expired.label, /Refund pending/);
    assert.match(expired.hint, /first trigger cancellation.*after confirmation.*claim.*refund/i);
  }
});

test("English action guidance respects exact registration and review boundaries", () => {
  assert.match(getTournamentPresentation(base, 99, "en").label, /Registration open/);
  assert.match(getTournamentPresentation(base, 100, "en").label, /Ready to start/);
  assert.match(getTournamentPresentation(base, 439, "en").label, /Ready to start/);
  assert.match(getTournamentPresentation(base, 440, "en").label, /Cannot start/);
  for (const status of ["result_submitted", "disputed"]) {
    const reviewing = getTournamentPresentation({ ...base, status }, 199, "en");
    assert.match(reviewing.hint, /after.*window closes/i);
    const reviewed = getTournamentPresentation({ ...base, status }, 200, "en");
    assert.match(reviewed.label, /Awaiting (approval|resolution)/);
    assert.match(reviewed.hint, /window is closed/i);
    const unknown = getTournamentPresentation({ ...base, status }, 0, "en");
    assert.match(unknown.hint, /Waiting for network time/);
    assert.doesNotMatch(unknown.label, /Awaiting|Ready|Can cancel|Refund pending/);
  }
  for (const prize_pool of ["0", "000", "invalid"]) {
    const unfunded = getTournamentPresentation({ ...base, prize_pool }, 100, "en");
    assert.match(unfunded.label, /Cannot start/);
    assert.doesNotMatch(unfunded.hint, /continue sponsoring|fund the prize pool/i);
  }
});

test("English global roles do not imply entry in a tournament and preserve the role list", () => {
  for (const roles of [[], ["participant"], ["unknown", "participant"]]) {
    const before = [...roles];
    Object.freeze(roles);
    assert.equal(globalRoleLabel(roles, "en"), "Standard user");
    assert.doesNotMatch(globalRoleLabel(roles, "en"), /participant/i);
    assert.deepEqual(roles, before);
  }
  const roles = Object.freeze(["auditor", "participant", "admin", "organizer", "admin"]);
  assert.equal(globalRoleLabel(roles, "en"), "Admin · Organizer · Auditor");
  assert.deepEqual(roles, ["auditor", "participant", "admin", "organizer", "admin"]);
});

test("approval and dispute resolution identify the auditor rather than the winning beneficiary", () => {
  const auditor = "0x1111111111111111111111111111111111111111";
  const winner = "0x2222222222222222222222222222222222222222";
  for (const event_name of ["ResultApproved", "DisputeResolved"]) {
    assert.equal(
      getEventActor({ event_name, actor: winner, data: { auditor, winner } }),
      auditor,
    );
  }
  assert.equal(
    getEventActor({
      event_name: "DisputeResolved",
      actor: "0x0000000000000000000000000000000000000000",
      data: { auditor, winner: "0x0000000000000000000000000000000000000000" },
    }),
    auditor,
  );
});

test("prize claims identify the receiving winner and keep financial data unchanged", () => {
  const winner = "0x2222222222222222222222222222222222222222";
  const event = {
    event_name: "PrizeClaimed",
    actor: "0x1111111111111111111111111111111111111111",
    amount_raw: "9007199254740993123456",
    data: { winner, amount: "9007199254740993123456" },
  };
  const before = JSON.stringify(event);
  assert.equal(getEventActor(event), winner);
  assert.equal(JSON.stringify(event), before);
});

test("missing, malformed, and unrecorded operator fields never guess a wallet", () => {
  const wallet = "0x1111111111111111111111111111111111111111";
  const missing = {
    event_name: "ResultApproved",
    actor: wallet,
    data: { winner: wallet },
  };
  assert.equal(getEventActor(missing), "");
  for (const auditor of [
    null,
    123,
    "not-a-wallet",
    "0x0000000000000000000000000000000000000000",
  ]) {
    assert.equal(getEventActor({ ...missing, data: { auditor } }), "");
  }
  for (const event_name of [
    "TournamentStarted",
    "RoleAdminChanged",
    "UnknownEvent",
  ]) {
    assert.equal(
      getEventActor({
        event_name,
        actor: wallet,
        data: { organizer: wallet, sender: wallet },
      }),
      "",
    );
  }
});

test("sponsorship, participant actions, and role grants distinguish operator fields from other wallets", () => {
  const wallet = "0x1111111111111111111111111111111111111111";
  const other = "0x2222222222222222222222222222222222222222";
  for (const event_name of ["TournamentFunded", "SponsorRefunded"]) {
    assert.equal(
      getEventActor({
        event_name,
        actor: other,
        data: { sponsor: wallet, winner: other },
      }),
      wallet,
    );
  }
  for (const event_name of ["ParticipantRegistered", "ParticipantWithdrawn"]) {
    assert.equal(
      getEventActor({
        event_name,
        actor: other,
        data: { participant: wallet },
      }),
      wallet,
    );
  }
  for (const event_name of ["RoleGranted", "RoleRevoked"]) {
    assert.equal(
      getEventActor({
        event_name,
        actor: other,
        data: { account: other, sender: wallet },
      }),
      wallet,
    );
  }
});
