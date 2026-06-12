import { mkdtemp, rm } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join, resolve } from "node:path";
import { pathToFileURL } from "node:url";
import * as esbuild from "esbuild";

const tmp = await mkdtemp(join(tmpdir(), "comical-worker-api-"));

try {
  const outfile = join(tmp, "worker.mjs");
  await esbuild.build({
    entryPoints: [resolve("src/index.ts")],
    outfile,
    bundle: true,
    format: "esm",
    platform: "neutral",
    target: "es2022",
    plugins: [
      {
        name: "policy-wasm-stub",
        setup(build) {
          build.onResolve({ filter: /^\.\/policy-wasm$/ }, () => ({
            path: "policy-wasm-stub",
            namespace: "stub"
          }));
          build.onLoad({ filter: /.*/, namespace: "stub" }, () => ({
            loader: "js",
            contents: `
              export async function policyWasm() {
                return {
                  validateSlug: (slug) => /^[a-zA-Z0-9_-]{3,64}$/.test(slug),
                  randomSlug: () => "random-slug",
                  expiryUnix: () => 0,
                  expiredUnix: () => false,
                  visitLimitExceeded: () => false
                };
              }
            `
          }));
        }
      }
    ]
  });

  const worker = (await import(pathToFileURL(outfile).href)).default;
  const response = await worker.fetch(
    new Request("https://tool.sqlboy.me/api/short-links", {
      method: "POST",
      headers: { "content-type": "application/json", "x-request-id": "req-test" },
      body: JSON.stringify({
        target_url: "https://example.com",
        custom_slug: "dup",
        ttl: "1h"
      })
    }),
    duplicateSlugEnv()
  );

  if (response.status !== 409) {
    throw new Error(`expected duplicate slug status 409, got ${response.status}: ${await response.text()}`);
  }
  const payload = await response.json();
  if (payload.error?.code !== "conflict") {
    throw new Error(`expected conflict error code, got ${JSON.stringify(payload)}`);
  }
} finally {
  await rm(tmp, { recursive: true, force: true });
}

function duplicateSlugEnv() {
  return {
    PUBLIC_BASE_URL: "https://tool.sqlboy.me",
    SHORT_DOMAIN_MAPPINGS: "{}",
    DB: {
      prepare() {
        return {
          bind() {
            return {
              async run() {
                throw new Error("D1_ERROR: UNIQUE constraint failed: short_links.slug");
              }
            };
          }
        };
      }
    },
    KV: {},
    BUCKET: {}
  };
}
