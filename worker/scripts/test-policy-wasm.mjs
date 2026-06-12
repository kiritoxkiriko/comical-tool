import { readFile } from "node:fs/promises";
import { resolve } from "node:path";

await import("../src/gen/wasm_exec.js");

const go = new globalThis.Go();
const bytes = await readFile(resolve(import.meta.dirname, "../src/gen/policy.wasm"));
const { instance } = await WebAssembly.instantiate(bytes, go.importObject);
void go.run(instance);
await Promise.resolve();

const policy = globalThis.comicalPolicy;
if (!policy?.validateSlug("abc_123-test")) throw new Error("validateSlug failed");
if (policy.validateSlug("ab")) throw new Error("validateSlug accepted short slug");
if (policy.expiryUnix("bad", 0) !== -1) throw new Error("invalid ttl did not fail");
if (policy.expiryUnix("", 0) !== 0) throw new Error("empty ttl should not expire");
if (!policy.expiredUnix(1)) throw new Error("expiredUnix failed");
if (!policy.visitLimitExceeded(3, 3)) throw new Error("visitLimitExceeded failed");

process.exit(0);
