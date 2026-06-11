"use client";

import { ArrowUpRight, UploadCloud } from "lucide-react";
import { type DragEvent, type FormEvent, type ReactNode, useState } from "react";

import type { ToastMessage, ToolMeta, ToolTab } from "@/components/tool-types";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";

const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export function ToolForm(props: {
  tab: ToolTab;
  meta: ToolMeta;
  notify: (message: Omit<ToastMessage, "id">) => void;
  setLoading: (value: boolean) => void;
}) {
  if (props.tab === "short") return <ShortForm {...props} />;
  if (props.tab === "clip") return <ClipForm {...props} />;
  return <UploadForm kind={props.tab} {...props} />;
}

function ShortForm({ meta, notify, setLoading }: FormProps) {
  const [target, setTarget] = useState("");
  const [slug, setSlug] = useState("");
  const [ttl, setTTL] = useState("168h");

  return (
    <form
      onSubmit={(event) =>
        submitJSON(event, "/api/short-links", { target_url: target, custom_slug: slug, ttl }, notify, setLoading)
      }
      className="grid gap-5"
    >
      <SectionTitle title={meta.title} desc={meta.desc} />
      <Field label="目标地址">
        <Input value={target} onChange={(event) => setTarget(event.target.value)} placeholder="https://example.com" />
      </Field>
      <div className="grid gap-4 sm:grid-cols-2">
        <Field label="自定义 slug">
          <Input value={slug} onChange={(event) => setSlug(event.target.value)} placeholder="可空" />
        </Field>
        <Field label="过期时间">
          <Input value={ttl} onChange={(event) => setTTL(event.target.value)} placeholder="168h" />
        </Field>
      </div>
      <SubmitButton>创建短链</SubmitButton>
    </form>
  );
}

function ClipForm({ meta, notify, setLoading }: FormProps) {
  const [content, setContent] = useState("");
  const [password, setPassword] = useState("");
  const [ttl, setTTL] = useState("1h");
  const [maxVisits, setMaxVisits] = useState("5");

  return (
    <form
      onSubmit={(event) =>
        submitJSON(
          event,
          "/api/clip",
          { content, password, ttl, max_visits: Number(maxVisits), link: true },
          notify,
          setLoading
        )
      }
      className="grid gap-5"
    >
      <SectionTitle title={meta.title} desc={meta.desc} />
      <Field label="剪贴内容">
        <Textarea value={content} onChange={(event) => setContent(event.target.value)} placeholder="粘贴临时文本" />
      </Field>
      <div className="grid gap-4 sm:grid-cols-3">
        <Field label="口令">
          <Input value={password} onChange={(event) => setPassword(event.target.value)} placeholder="可空" />
        </Field>
        <Field label="过期时间">
          <Input value={ttl} onChange={(event) => setTTL(event.target.value)} placeholder="1h" />
        </Field>
        <Field label="访问次数">
          <Input value={maxVisits} onChange={(event) => setMaxVisits(event.target.value)} placeholder="5" />
        </Field>
      </div>
      <SubmitButton>保存剪贴板</SubmitButton>
    </form>
  );
}

function UploadForm({ kind, meta, notify, setLoading }: FormProps & { kind: "image" | "file" }) {
  const [ttl, setTTL] = useState(kind === "image" ? "720h" : "168h");
  const [file, setFile] = useState<File | null>(null);

  return (
    <form onSubmit={(event) => upload(event, kind, ttl, file, notify, setLoading)} className="grid gap-5">
      <SectionTitle title={meta.title} desc={meta.desc} />
      <Dropzone kind={kind} file={file} onFile={setFile} />
      <Field label="过期时间">
        <Input value={ttl} onChange={(event) => setTTL(event.target.value)} placeholder="过期时间" />
      </Field>
      <SubmitButton>{kind === "image" ? "上传图片" : "上传文件"}</SubmitButton>
    </form>
  );
}

function Dropzone({ kind, file, onFile }: { kind: "image" | "file"; file: File | null; onFile: (file: File) => void }) {
  const [dragging, setDragging] = useState(false);
  const accept = kind === "image" ? "image/*" : undefined;

  function handleDrop(event: DragEvent<HTMLLabelElement>) {
    event.preventDefault();
    setDragging(false);
    const dropped = event.dataTransfer.files.item(0);
    if (dropped) onFile(dropped);
  }

  return (
    <label
      onDragOver={(event) => {
        event.preventDefault();
        setDragging(true);
      }}
      onDragLeave={() => setDragging(false)}
      onDrop={handleDrop}
      className={[
        "group grid min-h-44 cursor-pointer place-items-center rounded-3xl border border-dashed p-6 text-center transition",
        dragging
          ? "border-comicRed bg-comicRed/5"
          : "border-ink/20 bg-[#fbfaf6] hover:border-comicRed/60 hover:bg-comicYellow/10"
      ].join(" ")}
    >
      <input
        type="file"
        accept={accept}
        className="sr-only"
        onChange={(event) => {
          const selected = event.target.files?.item(0);
          if (selected) onFile(selected);
        }}
      />
      <span className="grid place-items-center gap-3">
        <span className="flex h-12 w-12 items-center justify-center rounded-2xl bg-white shadow-sm shadow-ink/5 transition group-hover:scale-105">
          <UploadCloud className="h-5 w-5 text-comicRed" />
        </span>
        <span>
          <span className="block text-base font-black">{file ? file.name : "拖拽文件到这里"}</span>
          <span className="mt-1 block text-sm text-ink/50">
            {file ? `${formatBytes(file.size)}，点击可重新选择` : "或点击选择文件，上传后自动生成短期资源"}
          </span>
        </span>
      </span>
    </label>
  );
}

function SubmitButton({ children }: { children: ReactNode }) {
  return (
    <Button variant="accent" className="mt-2 h-11 w-full sm:w-fit">
      {children}
      <ArrowUpRight className="h-4 w-4" />
    </Button>
  );
}

function Field({ label, children }: { label: string; children: ReactNode }) {
  return (
    <div className="grid gap-2">
      <span className="text-sm font-bold text-ink/62">{label}</span>
      {children}
    </div>
  );
}

function SectionTitle({ title, desc }: { title: string; desc: string }) {
  return (
    <div className="border-b border-ink/10 pb-5">
      <h2 className="text-2xl font-black tracking-normal">{title}</h2>
      <p className="mt-2 text-sm leading-6 text-ink/55">{desc}</p>
    </div>
  );
}

async function submitJSON(
  event: FormEvent,
  path: string,
  body: unknown,
  notify: (message: Omit<ToastMessage, "id">) => void,
  setLoading: (value: boolean) => void
) {
  event.preventDefault();
  await run(notify, setLoading, async () => {
    const res = await fetch(apiBase + path, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body)
    });
    return parseResponse(res);
  });
}

async function upload(
  event: FormEvent<HTMLFormElement>,
  kind: "image" | "file",
  ttl: string,
  file: File | null,
  notify: (message: Omit<ToastMessage, "id">) => void,
  setLoading: (value: boolean) => void
) {
  event.preventDefault();
  if (!file) {
    notify({ kind: "error", title: "请选择文件", description: "拖拽一个文件到上传区域，或点击上传区域选择文件。" });
    return;
  }
  await run(notify, setLoading, async () => {
    const data = new FormData();
    data.set("file", file);
    data.set("ttl", ttl);
    data.set("link", "true");
    const path = kind === "image" ? "/api/images" : "/api/files";
    const res = await fetch(apiBase + path, { method: "POST", body: data });
    return parseResponse(res);
  });
}

async function run(
  notify: (message: Omit<ToastMessage, "id">) => void,
  setLoading: (value: boolean) => void,
  action: () => Promise<unknown>
) {
  setLoading(true);
  try {
    const payload = await action();
    notify({ kind: "success", title: "操作成功", description: primaryMessage(payload) });
  } catch (error) {
    notify({ kind: "error", title: "操作失败", description: error instanceof Error ? error.message : String(error) });
  } finally {
    setLoading(false);
  }
}

async function parseResponse(res: Response): Promise<unknown> {
  const text = await res.text();
  const payload = parsePayload(text);
  if (!res.ok) {
    throw new Error(errorMessage(payload));
  }
  return payload;
}

function parsePayload(text: string): unknown {
  if (!text) return {};
  try {
    return JSON.parse(text);
  } catch {
    return text;
  }
}

function primaryMessage(payload: unknown): string {
  if (typeof payload === "string") return payload || "请求已完成。";
  if (!payload || typeof payload !== "object") return "请求已完成。";
  const value = payload as Record<string, unknown>;
  const primary = value.short_url ?? value.id ?? value.slug ?? value.deleted ?? value.revoked;
  if (typeof primary === "string") return primary;
  if (typeof primary === "boolean") return primary ? "已完成。" : "请求已完成。";
  return "请求已完成。";
}

function errorMessage(payload: unknown): string {
  if (typeof payload === "string") return payload || "请求失败。";
  if (!payload || typeof payload !== "object") return "请求失败。";
  const value = payload as Record<string, unknown>;
  if (typeof value.message === "string") return value.message;
  if (typeof value.error === "string") return value.error;
  return "请求失败。";
}

function formatBytes(size: number): string {
  if (size < 1024) return `${size} B`;
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`;
  return `${(size / 1024 / 1024).toFixed(1)} MB`;
}

type FormProps = {
  meta: ToolMeta;
  notify: (message: Omit<ToastMessage, "id">) => void;
  setLoading: (value: boolean) => void;
};
