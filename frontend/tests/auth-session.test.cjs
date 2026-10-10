const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const ts = require("typescript");

const RED = "0x" + "a".repeat(40);
const BLUE = "0x" + "b".repeat(40);
const session = (address = RED, token = "test-session", chain_id = 31337) => ({
  token, user: { address, roles: ["participant"] }, chain_id,
});
const connected = (address = RED, chainId = 31337) => ({ address, chainId, status: "connected" });
function compile(relative) {
  return ts.transpileModule(fs.readFileSync(path.resolve(__dirname, relative), "utf8"), {
    compilerOptions: { module: ts.ModuleKind.CommonJS, target: ts.ScriptTarget.ES2020, jsx: ts.JsxEmit.ReactJSX, esModuleInterop: true },
  }).outputText;
}
const authSource = compile("../src/lib/auth.tsx");
const providerSource = compile("../src/components/wallet-providers.tsx");
class ApiError extends Error {
  constructor(status) { super("Test API failure"); this.status = status; }
}

// Run the actual auth component with controlled requests and storage; no browser or real JWT.
function authHarness(initialAccount, initialSession) {
  let account = initialAccount;
  let cursor = 0;
  let dirty = false;
  let value;
  const hooks = [];
  const effects = [];
  const timers = new Map();
  const requests = [];
  let persisted = initialSession === undefined ? null : JSON.stringify(initialSession);
  let removed = 0;
  let nextTimer = 0;
  const storage = {
    getItem: () => persisted,
    setItem: (_key, text) => { persisted = text; },
    removeItem: () => { persisted = null; removed += 1; },
  };
  function slot(init) {
    const index = cursor++;
    if (!(index in hooks)) hooks[index] = init();
    return hooks[index];
  }
  const changed = (old, deps) => !old || deps.some((item, i) => item !== old[i]);
  const react = {
    createContext: () => ({ Provider: "AuthProvider" }),
    useContext: () => value,
    useState(initial) {
      const state = slot(() => ({ value: initial }));
      return [state.value, (next) => {
        const nextValue = typeof next === "function" ? next(state.value) : next;
        if (!Object.is(nextValue, state.value)) { state.value = nextValue; dirty = true; }
      }];
    },
    useRef: (initial) => slot(() => ({ current: initial })),
    useCallback(callback, deps) {
      const state = slot(() => ({ deps: null }));
      if (changed(state.deps, deps)) { state.callback = callback; state.deps = deps; }
      return state.callback;
    },
    useEffect(effect, deps) {
      const state = slot(() => ({ deps: null }));
      if (changed(state.deps, deps)) {
        state.deps = deps;
        effects.push(() => { state.cleanup?.(); state.cleanup = effect(); });
      }
    },
  };
  const exports = {};
  const mocks = {
    react,
    "react/jsx-runtime": { jsx: (type, props) => ({ type, props }) },
    wagmi: { useAccount: () => account },
    "@/lib/course": {
      ApiError,
      request(url, options) {
        let resolve;
        let reject;
        const promise = new Promise((yes, no) => { resolve = yes; reject = no; });
        requests.push({ url, options, resolve, reject });
        return promise;
      },
    },
  };
  new Function("require", "exports", "localStorage", "setInterval", "clearInterval", authSource)(
    (name) => { assert.ok(Object.hasOwn(mocks, name), name); return mocks[name]; },
    exports, storage,
    (callback) => { timers.set(++nextTimer, callback); return nextTimer; },
    (id) => timers.delete(id),
  );
  function render(runEffects = true) {
    cursor = 0;
    dirty = false;
    value = exports.AuthProvider({ children: null }).props.value;
    if (runEffects) flush();
    return value;
  }
  function flush() {
    let passes = 0;
    while (effects.length || dirty) {
      assert.ok(++passes < 20, "auth effects must settle");
      for (const effect of effects.splice(0)) effect();
      if (dirty) render(false);
    }
  }
  render();
  return {
    requests,
    get value() { return value; },
    get persisted() { return persisted && JSON.parse(persisted); },
    get removed() { return removed; },
    account(next, runEffects = true) { account = next; return render(runEffects); },
    flush,
    async settle() { await Promise.resolve(); await Promise.resolve(); flush(); },
  };
}

test("wallet bootstrap gates saved login before asynchronous hydration, then clears a genuine failed reconnect", async () => {
  const { createConfig, createStorage, hydrate, getAccount } = require("@wagmi/core");
  const { http } = require("viem");
  const { foundry } = require("viem/chains");
  const config = createConfig({
    chains: [foundry], ssr: true,
    transports: { [foundry.id]: http("http://localhost:9") },
    storage: createStorage({ storage: { getItem: () => null, setItem() {}, removeItem() {} } }),
  });
  assert.equal(getAccount(config).status, "disconnected");
  const exports = {};
  const mocks = {
    react: {},
    "react/jsx-runtime": { jsx: (type, props) => ({ type, props }) },
    wagmi: { WagmiProvider: "WagmiProvider" },
    "@tanstack/react-query": { QueryClient: class {}, QueryClientProvider: "QueryClientProvider" },
    "@rainbow-me/rainbowkit": { RainbowKitProvider: "RainbowKitProvider" },
    "@/lib/wagmi": { config },
    "@/lib/auth": { AuthProvider: "AuthProvider" },
    "@/lib/i18n": { useI18n: () => ({ lang: "en" }) },
    "@rainbow-me/rainbowkit/styles.css": {},
  };
  new Function("require", "exports", providerSource)((name) => mocks[name], exports);
  const element = exports.default({ children: null });
  const start = hydrate(config, { initialState: element.props.initialState, reconnectOnMount: true });
  assert.equal(getAccount(config).status, "reconnecting");
  const auth = authHarness(getAccount(config), session());
  assert.equal(auth.value.token, null);
  assert.equal(auth.persisted.token, "test-session");
  assert.equal(auth.removed, 0);
  const settled = new Promise((resolve) => {
    const unsubscribe = config.subscribe((state) => state.status, (status) => {
      if (status === "connected" || status === "disconnected") {
        unsubscribe();
        resolve();
      }
    });
  });
  await start.onMount(); // No connectors exist, so no network or wallet is accessed.
  await settled;
  assert.equal(getAccount(config).status, "disconnected");
  auth.account(getAccount(config));
  assert.equal(auth.persisted, null);
});

test("a matching session survives startup and restores only after wallet reconnects", () => {
  const auth = authHarness({ status: "reconnecting" }, session());
  assert.equal(auth.value.token, null);
  assert.equal(auth.requests.length, 0);
  assert.equal(auth.removed, 0);
  auth.account(connected());
  assert.equal(auth.value.token, "test-session");
  assert.equal(auth.requests[0].url, "/api/auth/profile");
  assert.equal(auth.requests[0].options.headers.Authorization, "Bearer test-session");
});

test("real disconnect, account and chain changes hide old tokens before cleanup", () => {
  for (const next of [{ status: "disconnected" }, connected(BLUE), connected(RED, 1)]) {
    const auth = authHarness(connected(), session());
    const count = auth.requests.length;
    auth.account(next, false);
    assert.equal(auth.value.token, null);
    assert.equal(auth.value.user, null);
    assert.equal(auth.value.isAuthenticated, false);
    assert.equal(auth.requests.length, count);
    auth.flush();
    assert.equal(auth.persisted, null);
    auth.account(connected());
    assert.equal(auth.value.token, null, "returning does not restore a cleared session");
  }
});

test("startup rejects sessions from another wallet or chain and malformed data", () => {
  for (const saved of [session(BLUE), session(RED, "test-session", 1), { token: "bad" }]) {
    const auth = authHarness(connected(), saved);
    assert.equal(auth.value.token, null);
    assert.equal(auth.persisted, null);
    assert.equal(auth.requests.length, 0);
  }
});

for (const outcome of ["success", "unauthorized"]) {
  test(`a stale profile ${outcome} cannot overwrite or clear a newly signed-in session`, async () => {
    const auth = authHarness(connected(), session());
    const oldProfile = auth.requests[0];
    const login = auth.value.login(RED, "test-nonce", "test-signature");
    auth.requests[1].resolve({ token: "new-session", user: { address: RED, roles: ["organizer"] } });
    await login;
    auth.flush();
    if (outcome === "success") oldProfile.resolve({ address: RED, roles: ["participant"] });
    else oldProfile.reject(new ApiError(401));
    await auth.settle();
    assert.equal(auth.value.token, "new-session");
    assert.deepEqual(auth.value.user.roles, ["organizer"]);
    assert.equal(auth.persisted.token, "new-session");
  });
}

test("a current unauthorized profile removes expired authentication", async () => {
  const auth = authHarness(connected(), session());
  auth.requests[0].reject(new ApiError(401));
  await auth.settle();
  assert.equal(auth.value.token, null);
  assert.equal(auth.persisted, null);
});

test("a late login after switching account, chain or signing out cannot restore the old session", async () => {
  for (const change of ["account", "chain", "logout"]) {
    const auth = authHarness(connected());
    const login = auth.value.login(RED, "test-nonce", "test-signature");
    const rejected = assert.rejects(login, /重新登录/);
    if (change === "account") { auth.account(connected(BLUE)); auth.account(connected()); }
    else if (change === "chain") { auth.account(connected(RED, 1)); auth.account(connected()); }
    else { auth.value.logout(); auth.flush(); }
    auth.requests[0].resolve({ token: "late-session", user: { address: RED, roles: ["participant"] } });
    await rejected;
    auth.flush();
    assert.equal(auth.value.token, null);
    assert.equal(auth.persisted, null);
  }
});
