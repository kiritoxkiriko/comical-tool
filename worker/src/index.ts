import { policyWasm } from "./policy-wasm";

const jsonHeaders = { "content-type": "application/json; charset=utf-8" };
const defaultShortTTLSeconds = 168 * 3600;
const defaultImageTTLSeconds = 720 * 3600;
const defaultClipTTLSeconds = 3600;
const defaultFileTTLSeconds = 168 * 3600;

export default {
  async fetch(request: Request, env: Env): Promise<Response> {
    const url = new URL(request.url);
    const meta = requestMeta(request);
    try {
      if (url.pathname === "/healthz" || url.pathname === "/api/health") return health(meta);
      if (url.pathname === "/api/short-links" && request.method === "POST") return createShort(request, env, meta);
      if (url.pathname.startsWith("/api/short-links/") && url.pathname.endsWith("/revoke"))
        return revokeShort(url, env, meta);
      if (url.pathname.startsWith("/short/") && request.method === "GET") return redirectShort(url, env, meta);
      if (url.pathname === "/api/clip" && request.method === "POST") return createClip(request, env, meta);
      if (url.pathname.startsWith("/api/clip/") && request.method === "GET") return getClip(url, env, meta);
      if (url.pathname.startsWith("/api/clip/") && request.method === "DELETE") return deleteClip(url, env, meta);
      if (url.pathname === "/api/images" && request.method === "POST") return uploadAsset(request, env, "image", meta);
      if (url.pathname === "/api/images" && request.method === "GET") return listAssets(env, "image", meta);
      if (url.pathname.startsWith("/api/images/") && request.method === "DELETE") return deleteAsset(url, env, meta);
      if (url.pathname === "/api/files" && request.method === "POST") return uploadAsset(request, env, "file", meta);
      if (url.pathname === "/api/files" && request.method === "GET") return listAssets(env, "file", meta);
      if (url.pathname.startsWith("/api/files/") && request.method === "DELETE") return deleteAsset(url, env, meta);
      if (url.pathname.startsWith("/api/assets/") && request.method === "GET") return getAsset(url, env, meta);
      if (url.pathname === "/api/admin/cleanup" && request.method === "POST") return adminCleanup(request, env, meta);
      if (request.method === "GET" && url.pathname.length > 1) return redirectShort(url, env, meta);
      return json({ error: "not_found", message: "route not found" }, 404, meta);
    } catch (error) {
      return json({ error: "internal", message: String(error) }, 500, meta);
    }
  },

  async scheduled(_controller: ScheduledController, env: Env, ctx: ExecutionContext): Promise<void> {
    ctx.waitUntil(cleanExpired(env));
  }
};

async function createShort(request: Request, env: Env, meta: RequestMeta): Promise<Response> {
  const body = (await request.json()) as ShortRequest;
  const targetURL = body.target_url?.trim() || "";
  if (!validTargetURL(targetURL)) return json({ error: "bad_request", message: "invalid target_url" }, 400, meta);
  const policy = await policyWasm();
  if (body.custom_slug && !policy.validateSlug(body.custom_slug))
    return json({ error: "bad_request", message: "invalid slug" }, 400, meta);
  const expiresAt = expiryUnix(policy, body.ttl, defaultShortTTLSeconds);
  if (expiresAt === invalidExpiry) return json({ error: "bad_request", message: "invalid ttl" }, 400, meta);
  const slug = body.custom_slug || policy.randomSlug();
  await env.DB.prepare("INSERT INTO short_links (id, slug, target_url, expires_at, created_at) VALUES (?, ?, ?, ?, ?)")
    .bind(crypto.randomUUID(), slug, targetURL, expiresAt, now())
    .run();
  return json(
    {
      slug,
      target_url: targetURL,
      short_url: publicURL(env, `/short/${slug}`),
      domain_urls: domainURLs(env, slug),
      mapped_urls: mappedURLs(env, slug),
      expires_at: expiresAt
    },
    200,
    meta
  );
}

async function revokeShort(url: URL, env: Env, meta: RequestMeta): Promise<Response> {
  const slug = url.pathname.split("/").at(-2);
  await env.DB.prepare("UPDATE short_links SET revoked_at = ? WHERE slug = ?").bind(now(), slug).run();
  return json({ revoked: true }, 200, meta);
}

async function redirectShort(url: URL, env: Env, meta: RequestMeta): Promise<Response> {
  const slug = url.pathname.split("/").filter(Boolean).pop() || "";
  const row = await env.DB.prepare("SELECT target_url, expires_at, revoked_at FROM short_links WHERE slug = ?")
    .bind(slug)
    .first<ShortRow>();
  if (!row) return json({ error: "not_found", message: "short link not found" }, 404, meta);
  const policy = await policyWasm();
  if (row.revoked_at || expiredUnix(policy, row.expires_at))
    return json({ error: "expired", message: "short link unavailable" }, 410, meta);
  const response = Response.redirect(row.target_url, 302);
  response.headers.set("x-request-id", meta.requestID);
  return response;
}

async function createClip(request: Request, env: Env, meta: RequestMeta): Promise<Response> {
  const body = (await request.json()) as ClipRequest;
  const id = crypto.randomUUID();
  const policy = await policyWasm();
  const expiresAt = expiryUnix(policy, body.ttl, defaultClipTTLSeconds);
  if (expiresAt === invalidExpiry) return json({ error: "bad_request", message: "invalid ttl" }, 400, meta);
  const item = {
    content: body.content,
    password_hash: await hashPassword(body.password || ""),
    max_visits: body.max_visits || 5,
    visit_count: 0,
    expires_at: expiresAt
  };
  await env.KV.put(`clip:${id}`, JSON.stringify(item), { expiration: expiresAt || undefined });
  return json({ id, expires_at: expiresAt }, 200, meta);
}

async function getClip(url: URL, env: Env, meta: RequestMeta): Promise<Response> {
  const id = url.pathname.split("/").pop() || "";
  const raw = await env.KV.get(`clip:${id}`);
  if (!raw) return json({ error: "not_found", message: "clipboard item not found" }, 404, meta);
  const item = JSON.parse(raw) as ClipItem;
  const policy = await policyWasm();
  if (expiredUnix(policy, item.expires_at))
    return json({ error: "expired", message: "clipboard item expired" }, 410, meta);
  if (!(await checkPassword(item.password_hash, url.searchParams.get("password") || "")))
    return json({ error: "forbidden", message: "invalid password" }, 403, meta);
  if (policy.visitLimitExceeded(item.max_visits, item.visit_count))
    return json({ error: "expired", message: "clipboard item exhausted" }, 410, meta);
  item.visit_count += 1;
  await env.KV.put(`clip:${id}`, JSON.stringify(item), { expiration: item.expires_at || undefined });
  return json({ id, content: item.content, visit_count: item.visit_count, expires_at: item.expires_at }, 200, meta);
}

async function deleteClip(url: URL, env: Env, meta: RequestMeta): Promise<Response> {
  const id = url.pathname.split("/").pop() || "";
  await env.KV.delete(`clip:${id}`);
  return json({ deleted: true }, 200, meta);
}

async function uploadAsset(request: Request, env: Env, kind: string, meta: RequestMeta): Promise<Response> {
  const form = await request.formData();
  const file = form.get("file");
  if (!(file instanceof File)) return json({ error: "bad_request", message: "file is required" }, 400, meta);
  const policy = await policyWasm();
  const expiresAt = expiryUnix(policy, String(form.get("ttl") || ""), defaultAssetTTLSeconds(kind));
  if (expiresAt === invalidExpiry) return json({ error: "bad_request", message: "invalid ttl" }, 400, meta);
  const maxVisits = kind === "file" ? parseMaxVisits(form.get("max_visits")) : 0;
  if (maxVisits < 0) return json({ error: "bad_request", message: "invalid max_visits" }, 400, meta);
  const passwordHash = kind === "file" ? await hashPassword(String(form.get("password") || "")) : "";
  const id = crypto.randomUUID();
  const key = `${kind}/${id}`;
  await env.BUCKET.put(key, file.stream(), { httpMetadata: { contentType: file.type } });
  await env.DB.prepare(
    "INSERT INTO assets (id, kind, name, content_type, size, object_key, password_hash, max_visits, visit_count, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
  )
    .bind(id, kind, file.name, file.type, file.size, key, passwordHash, maxVisits, 0, expiresAt, now())
    .run();
  return json(
    {
      id,
      kind,
      name: file.name,
      content_type: file.type,
      size: file.size,
      max_visits: maxVisits,
      expires_at: expiresAt
    },
    200,
    meta
  );
}

async function getAsset(url: URL, env: Env, meta: RequestMeta): Promise<Response> {
  const id = url.pathname.split("/").pop();
  const row = await env.DB.prepare(
    "SELECT object_key, content_type, expires_at, deleted_at, password_hash, max_visits, visit_count FROM assets WHERE id = ?"
  )
    .bind(id)
    .first<AssetRow>();
  if (!row) return json({ error: "not_found", message: "asset not found" }, 404, meta);
  const policy = await policyWasm();
  if (row.deleted_at || expiredUnix(policy, row.expires_at))
    return json({ error: "expired", message: "asset unavailable" }, 410, meta);
  if (policy.visitLimitExceeded(row.max_visits, row.visit_count))
    return json({ error: "expired", message: "asset exhausted" }, 410, meta);
  if (!(await checkPassword(row.password_hash, url.searchParams.get("password") || "")))
    return json({ error: "forbidden", message: "invalid password" }, 403, meta);
  const object = await env.BUCKET.get(row.object_key);
  if (!object) return json({ error: "not_found", message: "object not found" }, 404, meta);
  await env.DB.prepare("UPDATE assets SET visit_count = visit_count + 1 WHERE id = ?").bind(id).run();
  return new Response(object.body, { headers: { "content-type": row.content_type, "x-request-id": meta.requestID } });
}

async function listAssets(env: Env, kind: string, meta: RequestMeta): Promise<Response> {
  const result = await env.DB.prepare(
    "SELECT id, kind, name, content_type, size, short_slug, max_visits, visit_count, expires_at, created_at FROM assets WHERE kind = ? AND deleted_at IS NULL ORDER BY created_at DESC"
  )
    .bind(kind)
    .all();
  return json(result.results, 200, meta);
}

async function deleteAsset(url: URL, env: Env, meta: RequestMeta): Promise<Response> {
  const id = url.pathname.split("/").pop();
  const row = await env.DB.prepare("SELECT object_key FROM assets WHERE id = ?")
    .bind(id)
    .first<{ object_key: string }>();
  if (!row) return json({ error: "not_found", message: "asset not found" }, 404, meta);
  await env.BUCKET.delete(row.object_key);
  await env.DB.prepare("UPDATE assets SET deleted_at = ? WHERE id = ?").bind(now(), id).run();
  return json({ deleted: true }, 200, meta);
}

async function cleanExpired(env: Env): Promise<void> {
  const current = now();
  await env.DB.prepare("UPDATE assets SET deleted_at = ? WHERE expires_at IS NOT NULL AND expires_at < ?")
    .bind(current, current)
    .run();
}

async function adminCleanup(request: Request, env: Env, meta: RequestMeta): Promise<Response> {
  if (!adminAuthorized(request, env)) return json({ error: "unauthorized", message: "invalid admin token" }, 401, meta);
  await cleanExpired(env);
  return json({ cleanup: true }, 200, meta);
}

function json(value: unknown, status = 200, meta?: RequestMeta): Response {
  const payload = status >= 400 ? { error: errorPayload(value, meta) } : { data: value };
  const headers = new Headers(jsonHeaders);
  if (meta) headers.set("x-request-id", meta.requestID);
  return new Response(JSON.stringify(payload), { status, headers });
}

function health(meta: RequestMeta): Response {
  const headers = new Headers(jsonHeaders);
  headers.set("x-request-id", meta.requestID);
  return new Response(JSON.stringify({ ok: true }), { headers });
}

function errorPayload(value: unknown, meta?: RequestMeta): { code: string; message: string; request_id?: string } {
  if (!value || typeof value !== "object")
    return { code: "internal", message: String(value || "request failed"), request_id: meta?.requestID };
  const record = value as Record<string, unknown>;
  return {
    code: typeof record.error === "string" ? record.error : "internal",
    message: typeof record.message === "string" ? record.message : "request failed",
    request_id: meta?.requestID
  };
}

function requestMeta(request: Request): RequestMeta {
  return { requestID: request.headers.get("x-request-id")?.trim() || crypto.randomUUID() };
}

const invalidExpiry = -1;

function validTargetURL(value: string): boolean {
  try {
    const parsed = new URL(value);
    return parsed.protocol === "http:" || parsed.protocol === "https:";
  } catch {
    return false;
  }
}

function expiryUnix(
  policy: PolicyWasm,
  ttl: string | undefined,
  fallbackSeconds: number
): number | null | typeof invalidExpiry {
  const value = policy.expiryUnix(ttl || "", fallbackSeconds);
  if (value === invalidExpiry) return invalidExpiry;
  if (value <= 0) return null;
  return value;
}

function expiredUnix(policy: PolicyWasm, timestamp: number | null): boolean {
  return timestamp !== null && policy.expiredUnix(timestamp);
}

function defaultAssetTTLSeconds(kind: string): number {
  return kind === "image" ? defaultImageTTLSeconds : defaultFileTTLSeconds;
}

function parseMaxVisits(value: unknown): number {
  if (value === null || value === "") return 0;
  const parsed = Number.parseInt(String(value), 10);
  return Number.isNaN(parsed) ? -1 : parsed;
}

function now(): number {
  return Math.floor(Date.now() / 1000);
}

function publicURL(env: Env, path: string): string {
  return env.PUBLIC_BASE_URL.replace(/\/$/, "") + path;
}

function domainURLs(env: Env, slug: string): Record<string, string> {
  const mappings = shortDomainMappings(env);
  return Object.fromEntries(Object.keys(mappings).map((host) => [host, `https://${host}/${slug}`]));
}

function mappedURLs(env: Env, slug: string): Record<string, string> {
  const mappings = shortDomainMappings(env);
  return Object.fromEntries(
    Object.entries(mappings).map(([host, base]) => [host, `${base.replace(/\/$/, "")}/${slug}`])
  );
}

function shortDomainMappings(env: Env): Record<string, string> {
  if (!env.SHORT_DOMAIN_MAPPINGS) return {};
  return JSON.parse(env.SHORT_DOMAIN_MAPPINGS) as Record<string, string>;
}

function adminAuthorized(request: Request, env: Env): boolean {
  const token = (env as Env & { ADMIN_TOKEN?: string }).ADMIN_TOKEN?.trim();
  if (!token) return false;
  return request.headers.get("authorization")?.trim() === `Bearer ${token}`;
}

async function hashPassword(password: string): Promise<string> {
  if (!password) return "";
  const data = new TextEncoder().encode(password);
  const digest = await crypto.subtle.digest("SHA-256", data);
  return btoa(String.fromCharCode(...new Uint8Array(digest)));
}

async function checkPassword(hash: string, password: string): Promise<boolean> {
  return hash === "" || hash === (await hashPassword(password));
}

type ShortRequest = { target_url: string; custom_slug?: string; ttl?: string };
type ClipRequest = { content: string; password?: string; max_visits?: number; ttl?: string };
type ShortRow = { target_url: string; expires_at: number | null; revoked_at: number | null };
type AssetRow = {
  object_key: string;
  content_type: string;
  expires_at: number | null;
  deleted_at: number | null;
  password_hash: string;
  max_visits: number;
  visit_count: number;
};
type RequestMeta = { requestID: string };
type PolicyWasm = Awaited<ReturnType<typeof policyWasm>>;
type ClipItem = {
  content: string;
  password_hash: string;
  max_visits: number;
  visit_count: number;
  expires_at: number | null;
};
