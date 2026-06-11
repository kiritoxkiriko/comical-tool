"use client";

import { Copy, ExternalLink, RefreshCcw, RotateCcw, Trash2 } from "lucide-react";
import { type FormEvent, useState } from "react";

import { apiBase, parseResponse, runToolAction } from "@/components/tool-api";
import type { ToastMessage, ToolTab } from "@/components/tool-types";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

type ActionProps = {
  tab: ToolTab;
  notify: (message: Omit<ToastMessage, "id">) => void;
  setLoading: (value: boolean) => void;
};

type Asset = {
  id: string;
  name: string;
  content_type: string;
  size: number;
  short_slug?: string;
  expires_at?: string;
};

type Clip = {
  content: string;
};

export function ToolActions(props: ActionProps) {
  if (props.tab === "short") return <ShortActions {...props} />;
  if (props.tab === "clip") return <ClipActions {...props} />;
  return <AssetActions kind={props.tab} {...props} />;
}

function ShortActions({ notify, setLoading }: ActionProps) {
  const [slug, setSlug] = useState("");
  return (
    <ActionShell title="短链管理" desc="输入短链路径后撤销访问。">
      <form
        onSubmit={(event) => revokeShort(event, slug, notify, setLoading)}
        className="flex flex-col gap-3 sm:flex-row"
      >
        <Input value={slug} onChange={(event) => setSlug(event.target.value)} placeholder="short-slug" />
        <Button variant="outline" className="shrink-0">
          <RotateCcw className="h-4 w-4" />
          撤销
        </Button>
      </form>
    </ActionShell>
  );
}

function ClipActions({ notify, setLoading }: ActionProps) {
  const [id, setID] = useState("");
  const [password, setPassword] = useState("");
  return (
    <ActionShell title="剪贴板管理" desc="读取内容会复制到系统剪贴板，删除只需要条目 ID。">
      <form onSubmit={(event) => readClip(event, id, password, notify, setLoading)} className="grid gap-3">
        <div className="grid gap-3 sm:grid-cols-[1fr_180px]">
          <Input value={id} onChange={(event) => setID(event.target.value)} placeholder="clipboard-id" />
          <Input value={password} onChange={(event) => setPassword(event.target.value)} placeholder="口令，可空" />
        </div>
        <div className="flex flex-wrap gap-2">
          <Button variant="outline">
            <Copy className="h-4 w-4" />
            读取并复制
          </Button>
          <Button type="button" variant="ghost" onClick={() => deleteClip(id, notify, setLoading)}>
            <Trash2 className="h-4 w-4" />
            删除
          </Button>
        </div>
      </form>
    </ActionShell>
  );
}

function AssetActions({ kind, notify, setLoading }: ActionProps & { kind: "image" | "file" }) {
  const [assets, setAssets] = useState<Asset[]>([]);

  return (
    <ActionShell title={kind === "image" ? "图片列表" : "文件列表"} desc="查看最近资源，打开、复制短链或删除。">
      <div className="flex justify-end">
        <Button type="button" variant="outline" onClick={() => loadAssets(kind, notify, setLoading, setAssets)}>
          <RefreshCcw className="h-4 w-4" />
          刷新
        </Button>
      </div>
      <div className="divide-y divide-ink/10 rounded-2xl border border-ink/10">
        {assets.length === 0 ? (
          <p className="p-4 text-sm text-ink/50">暂无资源。</p>
        ) : (
          assets.map((asset) => (
            <div key={asset.id} className="grid gap-3 p-4 sm:grid-cols-[minmax(0,1fr)_auto] sm:items-center">
              <div className="min-w-0">
                <p className="truncate text-sm font-bold">{asset.name}</p>
                <p className="mt-1 text-xs text-ink/45">
                  {formatBytes(asset.size)} · {asset.content_type || "application/octet-stream"}
                </p>
              </div>
              <div className="flex flex-wrap gap-2">
                <Button asChild variant="ghost" size="icon" title="打开资源">
                  <a
                    href={`${apiBase}/api/${kind === "image" ? "images" : "files"}/${asset.id}`}
                    target="_blank"
                    rel="noopener"
                  >
                    <ExternalLink className="h-4 w-4" />
                  </a>
                </Button>
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  title="复制短链"
                  disabled={!asset.short_slug}
                  onClick={() => copyText(`${apiBase}/short/${asset.short_slug}`, notify)}
                >
                  <Copy className="h-4 w-4" />
                </Button>
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  title="删除资源"
                  onClick={() => deleteAsset(kind, asset.id, notify, setLoading, setAssets)}
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              </div>
            </div>
          ))
        )}
      </div>
    </ActionShell>
  );
}

function ActionShell({ title, desc, children }: { title: string; desc: string; children: React.ReactNode }) {
  return (
    <section className="grid gap-4 border-t border-ink/10 pt-6">
      <div>
        <h3 className="text-lg font-black tracking-normal">{title}</h3>
        <p className="mt-1 text-sm leading-6 text-ink/50">{desc}</p>
      </div>
      {children}
    </section>
  );
}

async function revokeShort(
  event: FormEvent,
  slug: string,
  notify: ActionProps["notify"],
  setLoading: ActionProps["setLoading"]
) {
  event.preventDefault();
  if (!slug) return notify({ kind: "error", title: "缺少路径", description: "请输入要撤销的短链路径。" });
  await runToolAction(notify, setLoading, async () => {
    const res = await fetch(`${apiBase}/api/short-links/${encodeURIComponent(slug)}/revoke`, { method: "POST" });
    return parseResponse(res);
  });
}

async function readClip(
  event: FormEvent,
  id: string,
  password: string,
  notify: ActionProps["notify"],
  setLoading: ActionProps["setLoading"]
) {
  event.preventDefault();
  if (!id) return notify({ kind: "error", title: "缺少 ID", description: "请输入剪贴板条目 ID。" });
  const query = new URLSearchParams({ password });
  const payload = await runToolAction(notify, setLoading, async () => {
    const res = await fetch(`${apiBase}/api/clip/${encodeURIComponent(id)}?${query.toString()}`);
    return parseResponse<Clip>(res);
  });
  if (payload && typeof payload === "object" && "content" in payload) {
    await copyText(String((payload as Clip).content), notify);
  }
}

async function deleteClip(id: string, notify: ActionProps["notify"], setLoading: ActionProps["setLoading"]) {
  if (!id) return notify({ kind: "error", title: "缺少 ID", description: "请输入要删除的剪贴板条目 ID。" });
  await runToolAction(notify, setLoading, async () => {
    const res = await fetch(`${apiBase}/api/clip/${encodeURIComponent(id)}`, { method: "DELETE" });
    return parseResponse(res);
  });
}

async function loadAssets(
  kind: "image" | "file",
  notify: ActionProps["notify"],
  setLoading: ActionProps["setLoading"],
  setAssets: (items: Asset[]) => void
) {
  const path = kind === "image" ? "/api/images" : "/api/files";
  const payload = await runToolAction(notify, setLoading, async () =>
    parseResponse<Asset[]>(await fetch(apiBase + path))
  );
  if (Array.isArray(payload)) setAssets(payload);
}

async function deleteAsset(
  kind: "image" | "file",
  id: string,
  notify: ActionProps["notify"],
  setLoading: ActionProps["setLoading"],
  setAssets: (updater: (items: Asset[]) => Asset[]) => void
) {
  const path = kind === "image" ? "/api/images" : "/api/files";
  const payload = await runToolAction(notify, setLoading, async () => {
    const res = await fetch(`${apiBase}${path}/${encodeURIComponent(id)}`, { method: "DELETE" });
    return parseResponse(res);
  });
  if (payload) setAssets((items) => items.filter((asset) => asset.id !== id));
}

async function copyText(value: string, notify: ActionProps["notify"]) {
  try {
    await navigator.clipboard.writeText(value);
    notify({ kind: "success", title: "已复制", description: value });
  } catch {
    notify({ kind: "error", title: "复制失败", description: "浏览器未允许访问剪贴板。" });
  }
}

function formatBytes(size: number): string {
  if (size < 1024) return `${size} B`;
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`;
  return `${(size / 1024 / 1024).toFixed(1)} MB`;
}
