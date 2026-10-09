# Course frontend — Agent Arena

Next.js 14 / React / wagmi / RainbowKit / viem. Course entry points: `/` financial dashboard, `/tournaments`, `/tournaments/:id`, `/audit`, `/agents`, `/guide`.

## Run

Use Node.js 20 or newer. Install the checked-in dependencies with `npm ci`, copy `.env.example` to `.env.local`, start the course backend described in `../docs/course-project/DEPLOYMENT.md`, then run `npm run dev`. The frontend reads deployment addresses, token precision, chain ID and confirmation depth from `/api/config`. An unconfigured deployment is visibly read-only; no historical contract addresses are used.

Set `NEXT_PUBLIC_API_URL` to the public backend origin. An explicitly empty value supports a same-origin `/api` reverse proxy. `NEXT_PUBLIC_LOCAL_RPC_URL` defaults to `http://127.0.0.1:8545`; both local Anvil 31337 and Sepolia 11155111 are supported. Browser-visible RPC URLs must not contain privileged secrets.

The UI obtains a one-time server nonce for wallet login. Switching account/network clears the session. All financial actions require wallet signatures. Sponsor approval uses `PrizeEscrow` as spender; gas is estimated before signing and receipts are tracked. Course demo token DUSD uses 6 decimals and has no redemption commitment.

## Checks

```sh
npm run typecheck
npm run lint
npm test
npm run build
```

Pure domain tests compile the actual TypeScript helper module in memory using the installed TypeScript dependency; they cover exact monetary units, input precision/overflow, timeouts and evidence tampering. The Next.js build retains the ESLint gate and produces standalone output for the deployment image.

The complete Chinese user guide is `../docs/course-project/USER_GUIDE.md`. Generated ABI source is `src/lib/course-abi.ts`; regenerate it with `contracts/export_abi.py` when contracts change.
