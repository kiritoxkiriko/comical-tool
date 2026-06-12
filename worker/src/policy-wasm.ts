import "./gen/wasm_exec.js";
import policyModule from "./gen/policy.wasm";

type PolicyWasm = NonNullable<typeof globalThis.comicalPolicy>;

let policyPromise: Promise<PolicyWasm> | undefined;

export function policyWasm(): Promise<PolicyWasm> {
  policyPromise ??= startPolicyWasm();
  return policyPromise;
}

async function startPolicyWasm(): Promise<PolicyWasm> {
  const go = new Go();
  const instance = new WebAssembly.Instance(policyModule, go.importObject);
  void go.run(instance);
  await Promise.resolve();
  if (!globalThis.comicalPolicy) {
    throw new Error("comical policy wasm did not initialize");
  }
  return globalThis.comicalPolicy;
}
