const jsonHeaders = { "content-type": "application/json; charset=utf-8" };

export default {
  async fetch(request: Request, env: Env): Promise<Response> {
    const url = new URL(request.url);
    try {
      if (url.pathname === "/healthz") return json({ ok: true });
      if (url.pathname === "/api/short-links" && request.method === "POST") return createShort(request, env);
      if (url.pathname.startsWith("/api/short-links/") && url.pathname.endsWith("/revoke")) return revokeShort(url, env);
      if (url.pathname.startsWith("/short/") && request.method === "GET") return redirectShort(url, env);
      if (url.pathname === "/api/clip" && request.method === "POST") return createClip(request, env);
      if (url.pathname.startsWith("/api/clip/") && request.method === "GET") return getClip(url, env);
      if (url.pathname === "/api/images" && request.method === "POST") return uploadAsset(request, env, "image");
      if (url.pathname === "/api/files" && request.method === "POST") return uploadAsset(request, env, "file");
      if (url.pathname.startsWith("/api/assets/") && request.method === "GET") return getAsset(url, env);
      if (request.method === "GET" && url.pathname.length > 1) return redirectShort(url, env);
      return json({ error: "not_found", message: "route not found" }, 404);
    } catch (error) {
      return json({ error: "internal", message: String(error) }, 500);
    }
  },

  async scheduled(_controller: ScheduledController, env: Env, ctx: ExecutionContext): Promise<void> {
    ctx.waitUntil(cleanExpired(env));
  }
};

async function createShort(request: Request, env: Env): Promise<Response> {
  const body = await request.json() as ShortRequest;
  const slug = body.custom_slug || randomSlug();
  const expiresAt = parseExpiry(body.ttl);
  await env.DB.prepare(
    "INSERT INTO short_links (id, slug, target_url, expires_at, created_at) VALUES (?, ?, ?, ?, ?)"
  ).bind(crypto.randomUUID(), slug, body.target_url, expiresAt, now()).run();
  return json({
    slug,
    target_url: body.target_url,
    short_url: publicURL(env, "/short/" + slug),
    domain_urls: domainURLs(env, slug),
    mapped_urls: mappedURLs(env, slug),
    expires_at: expiresAt
  });
}

async function revokeShort(url: URL, env: Env): Promise<Response> {
  const slug = url.pathname.split("/").at(-2);
  await env.DB.prepare("UPDATE short_links SET revoked_at = ? WHERE slug = ?").bind(now(), slug).run();
  return json({ revoked: true });
}

async function redirectShort(url: URL, env: Env): Promise<Response> {
  const slug = url.pathname.split("/").filter(Boolean).pop() || "";
  const row = await env.DB.prepare("SELECT target_url, expires_at, revoked_at FROM short_links WHERE slug = ?")
    .bind(slug).first<ShortRow>();
  if (!row) return json({ error: "not_found", message: "short link not found" }, 404);
  if (row.revoked_at || expired(row.expires_at)) return json({ error: "expired", message: "short link unavailable" }, 410);
  return Response.redirect(row.target_url, 302);
}

async function createClip(request: Request, env: Env): Promise<Response> {
  const body = await request.json() as ClipRequest;
  const id = crypto.randomUUID();
  const expiresAt = parseExpiry(body.ttl || "1h");
  const item = { content: body.content, password_hash: await hashPassword(body.password || ""), max_visits: body.max_visits || 5, visit_count: 0, expires_at: expiresAt };
  await env.KV.put("clip:" + id, JSON.stringify(item), { expiration: expiresAt || undefined });
  return json({ id, expires_at: expiresAt });
}

async function getClip(url: URL, env: Env): Promise<Response> {
  const id = url.pathname.split("/").pop() || "";
  const raw = await env.KV.get("clip:" + id);
  if (!raw) return json({ error: "not_found", message: "clipboard item not found" }, 404);
  const item = JSON.parse(raw) as ClipItem;
  if (expired(item.expires_at)) return json({ error: "expired", message: "clipboard item expired" }, 410);
  if (!(await checkPassword(item.password_hash, url.searchParams.get("password") || ""))) return json({ error: "forbidden", message: "invalid password" }, 403);
  if (item.max_visits > 0 && item.visit_count >= item.max_visits) return json({ error: "expired", message: "clipboard item exhausted" }, 410);
  item.visit_count += 1;
  await env.KV.put("clip:" + id, JSON.stringify(item), { expiration: item.expires_at || undefined });
  return json({ id, content: item.content, visit_count: item.visit_count, expires_at: item.expires_at });
}

async function uploadAsset(request: Request, env: Env, kind: string): Promise<Response> {
  const form = await request.formData();
  const file = form.get("file");
  if (!(file instanceof File)) return json({ error: "bad_request", message: "file is required" }, 400);
  const id = crypto.randomUUID();
  const key = `${kind}/${id}`;
  await env.BUCKET.put(key, file.stream(), { httpMetadata: { contentType: file.type } });
  const expiresAt = parseExpiry(String(form.get("ttl") || ""));
  await env.DB.prepare(
    "INSERT INTO assets (id, kind, name, content_type, size, object_key, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
  ).bind(id, kind, file.name, file.type, file.size, key, expiresAt, now()).run();
  return json({ id, kind, name: file.name, content_type: file.type, size: file.size, expires_at: expiresAt });
}

async function getAsset(url: URL, env: Env): Promise<Response> {
  const id = url.pathname.split("/").pop();
  const row = await env.DB.prepare("SELECT object_key, content_type, expires_at, deleted_at FROM assets WHERE id = ?")
    .bind(id).first<AssetRow>();
  if (!row) return json({ error: "not_found", message: "asset not found" }, 404);
  if (row.deleted_at || expired(row.expires_at)) return json({ error: "expired", message: "asset unavailable" }, 410);
  const object = await env.BUCKET.get(row.object_key);
  if (!object) return json({ error: "not_found", message: "object not found" }, 404);
  return new Response(object.body, { headers: { "content-type": row.content_type } });
}

async function cleanExpired(env: Env): Promise<void> {
  const current = now();
  await env.DB.prepare("UPDATE assets SET deleted_at = ? WHERE expires_at IS NOT NULL AND expires_at < ?").bind(current, current).run();
}

function json(value: unknown, status = 200): Response {
  return new Response(JSON.stringify(value), { status, headers: jsonHeaders });
}

function randomSlug(): string {
  const bytes = new Uint8Array(4);
  crypto.getRandomValues(bytes);
  return [...bytes].map((value) => value.toString(16).padStart(2, "0")).join("");
}

function parseExpiry(ttl?: string): number | null {
  if (!ttl) return null;
  const match = /^(\d+)(m|h|d)$/.exec(ttl);
  if (!match) return null;
  const units = { m: 60, h: 3600, d: 86400 } as const;
  return now() + Number(match[1]) * units[match[2] as keyof typeof units];
}

function expired(timestamp: number | null): boolean {
  return timestamp !== null && timestamp < now();
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
  return Object.fromEntries(Object.entries(mappings).map(([host, base]) => [host, base.replace(/\/$/, "") + "/" + slug]));
}

function shortDomainMappings(env: Env): Record<string, string> {
  if (!env.SHORT_DOMAIN_MAPPINGS) return {};
  return JSON.parse(env.SHORT_DOMAIN_MAPPINGS) as Record<string, string>;
}

async function hashPassword(password: string): Promise<string> {
  if (!password) return "";
  const data = new TextEncoder().encode(password);
  const digest = await crypto.subtle.digest("SHA-256", data);
  return btoa(String.fromCharCode(...new Uint8Array(digest)));
}

async function checkPassword(hash: string, password: string): Promise<boolean> {
  return hash === "" || hash === await hashPassword(password);
}

type ShortRequest = { target_url: string; custom_slug?: string; ttl?: string };
type ClipRequest = { content: string; password?: string; max_visits?: number; ttl?: string };
type ShortRow = { target_url: string; expires_at: number | null; revoked_at: number | null };
type AssetRow = { object_key: string; content_type: string; expires_at: number | null; deleted_at: number | null };
type ClipItem = { content: string; password_hash: string; max_visits: number; visit_count: number; expires_at: number | null };
