# worker

Cloudflare Worker API adapter for `comical-tool`.

```bash
npm install
npm run cf-typegen
npm run build
npm run dry-run
```

`npm run build` runs `npm run build:wasm` first. That command compiles
`worker/wasm/policy` to `worker/src/gen/policy.wasm` and copies Go's
`wasm_exec.js`; generated files are ignored by Git.

Bindings:

- D1: metadata
- R2: images and files
- KV: clipboard and volatile data

The Worker does not run Hertz. It reuses slug, TTL, expiry, and visit-limit
policy from `server/pkg/policy` through the generated Go WASM module, while D1,
R2, KV, and routing remain Worker-specific.
