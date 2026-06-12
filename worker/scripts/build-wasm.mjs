import { execFileSync } from "node:child_process";
import { chmod, copyFile, mkdir, rm, writeFile } from "node:fs/promises";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const root = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const outDir = resolve(root, "src/gen");
const goEnv = { ...process.env, GOTOOLCHAIN: "go1.26.0" };

await mkdir(outDir, { recursive: true });
const goroot = execFileSync("go", ["env", "GOROOT"], { env: goEnv, encoding: "utf8" }).trim();
const wasmExec = resolve(outDir, "wasm_exec.js");
await rm(wasmExec, { force: true });
await copyFile(resolve(goroot, "lib/wasm/wasm_exec.js"), wasmExec);
await chmod(wasmExec, 0o644);
await writeFile(resolve(outDir, "wasm_exec.js.d.ts"), "export {};\n");
execFileSync("go", ["build", "-trimpath", "-o", resolve(outDir, "policy.wasm"), "."], {
  cwd: resolve(root, "wasm/policy"),
  env: { ...goEnv, GOWORK: "off", GOOS: "js", GOARCH: "wasm" },
  stdio: "inherit"
});
await writeFile(
  resolve(outDir, "policy.wasm.d.ts"),
  "declare const module: WebAssembly.Module;\nexport default module;\n"
);
