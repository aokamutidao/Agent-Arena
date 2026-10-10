const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const Module = require("node:module");
const ts = require("typescript");

const file = path.resolve(__dirname, "../src/app/agents/page.tsx");
const output = ts.transpileModule(fs.readFileSync(file, "utf8"), {
  compilerOptions: {
    module: ts.ModuleKind.CommonJS,
    target: ts.ScriptTarget.ES2020,
    jsx: ts.JsxEmit.ReactJSX,
    esModuleInterop: true,
  },
}).outputText;

// Execute the real page with a minimal keyed hook lifecycle; this is not a DOM/browser test.
function createHarness() {
  let account = { address: "0xAaAa", chainId: 31337, token: "red-session" };
  let current;
  let tree;
  const instances = new Map();
  const requests = [];
  const invalidations = [];
  const profileTokens = [];
  const unmountedWrites = [];
  const Fragment = Symbol("Fragment");
  const element = (type, props, key) => ({ type, props: props || {}, key });

  function slot(initial) {
    const index = current.cursor++;
    if (!(index in current.hooks)) current.hooks[index] = initial();
    return current.hooks[index];
  }

  const react = {
    useState(initial) {
      const owner = current;
      const state = slot(() => ({ value: typeof initial === "function" ? initial() : initial }));
      return [state.value, (value) => {
        if (!owner.mounted) {
          unmountedWrites.push(owner.type.name);
          return;
        }
        state.value = typeof value === "function" ? value(state.value) : value;
      }];
    },
    useRef(initial) {
      return slot(() => ({ current: initial }));
    },
    useEffect(effect, dependencies) {
      const state = slot(() => ({ dependencies: undefined, cleanup: undefined }));
      if (!dependencies || !state.dependencies ||
          dependencies.some((value, index) => value !== state.dependencies[index])) {
        current.effects.push(() => {
          state.cleanup?.();
          state.cleanup = effect();
          state.dependencies = dependencies;
        });
      }
    },
  };

  const mocks = {
    react,
    "react/jsx-runtime": { jsx: element, jsxs: element, Fragment },
    "next/link": { __esModule: true, default: "a" },
    wagmi: { useAccount: () => account },
    "@/lib/auth": {
      useAuth: () => ({ token: account.token, user: { address: account.authAddress ?? account.address } }),
    },
    "@/lib/i18n": { useI18n: () => ({ lang: "en", l: (_zh, en) => en }) },
    "@tanstack/react-query": {
      useQueryClient: () => ({ invalidateQueries: async (options) => invalidations.push(options) }),
    },
    "@/lib/course": {
      request(url, options) {
        let resolve;
        let reject;
        const promise = new Promise((yes, no) => { resolve = yes; reject = no; });
        requests.push({ url, options, resolve, reject });
        return promise;
      },
      errorMessage: (error) => error.message,
      dateTime: (value) => String(value),
    },
    "@/hooks/useCourse": {
      useAuthenticatedResource: (_url, token) => {
        profileTokens.push(token);
        return { data: token ? { agents: [], total: 0 } : undefined, isLoading: false };
      },
    },
    "lucide-react": new Proxy({}, { get: (_target, name) => `icon-${String(name)}` }),
    "@/components/course/DeploymentNotice": { DeploymentNotice: "deployment-notice" },
    "@/components/auth/LoginButton": { LoginButton: "login-button" },
    "@/components/layout/WalletConnect": { WalletConnect: "wallet-connect" },
    "@/components/ui/button": { Button: "button" },
    "@/components/ui/input": { Input: "input" },
  };
  const loaded = new Module(file);
  loaded.filename = file;
  loaded.require = (name) => {
    assert.ok(Object.hasOwn(mocks, name), `Unexpected page dependency: ${name}`);
    return mocks[name];
  };
  loaded._compile(output, file);

  function unmount(instance) {
    instance.mounted = false;
    for (const hook of instance.hooks) hook.cleanup?.();
  }

  function renderNode(node, location, seen) {
    if (node === null || node === undefined || typeof node === "boolean") return null;
    if (Array.isArray(node)) return node.map((child, index) => renderNode(child, `${location}.${index}`, seen));
    if (typeof node !== "object") return node;
    if (node.type === Fragment) return renderNode(node.props.children, location, seen);
    if (typeof node.type === "function") {
      const identity = `${location}:${node.key ?? ""}`;
      let instance = instances.get(identity);
      if (instance && instance.type !== node.type) {
        unmount(instance);
        instances.delete(identity);
        instance = undefined;
      }
      if (!instance) {
        instance = { type: node.type, hooks: [], mounted: true, effects: [] };
        instances.set(identity, instance);
      }
      seen.add(identity);
      const previous = current;
      current = instance;
      current.cursor = 0;
      const rendered = node.type(node.props);
      current = previous;
      return renderNode(rendered, `${identity}.child`, seen);
    }
    return { ...node, props: { ...node.props, children: renderNode(node.props.children, `${location}.children`, seen) } };
  }

  function render() {
    const seen = new Set();
    tree = renderNode(element(loaded.exports.default), "root", seen);
    for (const [identity, instance] of instances) {
      if (!seen.has(identity)) {
        unmount(instance);
        instances.delete(identity);
      }
    }
    for (const instance of instances.values()) {
      for (const effect of instance.effects.splice(0)) effect();
    }
    return tree;
  }

  function findAll(predicate, node = tree) {
    if (Array.isArray(node)) return node.flatMap((child) => findAll(predicate, child));
    if (!node || typeof node !== "object") return [];
    return [...(predicate(node) ? [node] : []), ...findAll(predicate, node.props.children)];
  }

  function find(type) {
    const matches = findAll((node) => node.type === type);
    assert.equal(matches.length, 1, `Expected one ${type}`);
    return matches[0];
  }

  render();
  return {
    requests, invalidations, profileTokens, unmountedWrites, find,
    switchAccount(changes) { account = { ...account, ...changes }; render(); },
    fill(name, strategy = "defensive") {
      find("input").props.onChange({ target: { value: name } });
      find("textarea").props.onChange({ target: { value: strategy } });
      render();
    },
    submit() {
      const promise = find("form").props.onSubmit({ preventDefault() {} });
      render();
      return promise;
    },
    render,
    feedback(role) { return findAll((node) => node.props.role === role); },
    submitButton() { return findAll((node) => node.type === "button" && node.props.type === "submit")[0]; },
  };
}

test("Agent feedback and drafts reset for wallet, chain, session, and sign-out changes", async () => {
  const page = createHarness();
  page.fill("Red agent");
  const saving = page.submit();
  page.requests[0].resolve({});
  await saving;
  page.render();
  assert.equal(page.feedback("status").length, 1);

  page.switchAccount({ address: "0xBbBb", authAddress: "0xAaAa" });
  assert.equal(page.submitButton(), undefined, "a new wallet cannot submit with the old wallet's auth session");
  assert.equal(page.profileTokens.at(-1), null, "the old session cannot read the new wallet's profile list");
  assert.equal(page.feedback("status").length, 0);
  page.switchAccount({ token: "blue-session", authAddress: "0xBbBb" });
  assert.equal(page.feedback("status").length, 0);
  assert.equal(page.find("input").props.value, "");
  assert.equal(page.find("textarea").props.value, "");
  page.fill("Blue draft");
  page.switchAccount({ address: "0xBBBB" });
  assert.equal(page.find("input").props.value, "Blue draft", "address casing is the same account");

  for (const change of [{ chainId: 1 }, { token: "renewed-blue-session" }]) {
    page.fill("Private draft");
    page.switchAccount(change);
    assert.equal(page.find("input").props.value, "");
    assert.equal(page.find("textarea").props.value, "");
    assert.equal(page.feedback("status").length, 0);
  }
  page.fill("Draft before sign-out");
  page.switchAccount({ token: null });
  assert.equal(page.submitButton(), undefined);
  page.switchAccount({ token: "another-blue-session" });
  assert.equal(page.find("input").props.value, "");
});

for (const outcome of ["success", "failure"]) {
  test(`a previous wallet's delayed ${outcome} cannot alter the new wallet's active save`, async () => {
    const page = createHarness();
    page.fill(" Red agent ", " aggressive ");
    const oldSave = page.submit();
    page.switchAccount({ address: "0xBbBb", token: "blue-session" });
    page.fill("Blue agent", "Blue strategy");
    const newSave = page.submit();
    assert.equal(page.requests.length, 2);

    if (outcome === "success") page.requests[0].resolve({});
    else page.requests[0].reject(new Error("Red wallet save failed"));
    await oldSave;
    page.render();
    assert.equal(page.find("input").props.value, "Blue agent");
    assert.equal(page.find("textarea").props.value, "Blue strategy");
    assert.equal(page.submitButton().props.children, "Saving…");
    assert.equal(page.submitButton().props.disabled, true);
    assert.equal(page.feedback("status").length, 0);
    assert.equal(page.feedback("alert").length, 0);
    assert.deepEqual(page.invalidations, []);
    assert.deepEqual(page.unmountedWrites, [], "old request must stop before updating an unmounted workspace");

    await page.submit();
    assert.equal(page.requests.length, 2, "the current pending save still blocks duplicate submission");
    assert.equal(page.requests[0].options.headers.Authorization, "Bearer red-session");
    assert.deepEqual(JSON.parse(page.requests[0].options.body), { name: "Red agent", personality: "aggressive" });
    assert.equal(page.requests[1].options.headers.Authorization, "Bearer blue-session");
    page.requests[1].resolve({});
    await newSave;
    page.render();
    assert.equal(page.feedback("status").length, 1);
    assert.equal(page.find("input").props.value, "");
    assert.equal(page.submitButton().props.children, "Save agent and prepare to compete");
    assert.deepEqual(page.invalidations, [{ queryKey: ["course", "/api/agents/my", "blue-session"], exact: true }]);
  });
}

test("a current save failure preserves its draft and becomes retryable", async () => {
  const page = createHarness();
  page.fill("Retry agent", "Retry strategy");
  const saving = page.submit();
  page.requests[0].reject(new Error("Save unavailable"));
  await saving;
  page.render();
  assert.equal(page.feedback("alert")[0].props.children, "Save unavailable");
  assert.equal(page.find("input").props.value, "Retry agent");
  assert.equal(page.find("textarea").props.value, "Retry strategy");
  assert.equal(page.submitButton().props.disabled, false);
  assert.deepEqual(page.invalidations, []);
  page.switchAccount({ address: "0xBbBb", token: "blue-session" });
  assert.equal(page.feedback("alert").length, 0, "switching wallets clears the previous wallet's error");
  assert.equal(page.find("input").props.value, "");
  assert.equal(page.find("textarea").props.value, "");
});
