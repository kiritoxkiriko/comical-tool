"use client";

import { FormEvent, useState } from "react";

const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";
const tabs = ["short", "image", "clip", "file"] as const;
type Tab = (typeof tabs)[number];

export default function Home() {
  const [tab, setTab] = useState<Tab>("short");
  const [result, setResult] = useState("");

  return (
    <main className="min-h-screen bg-[#fff8e1] text-ink">
      <header className="border-b-4 border-ink bg-comicYellow">
        <div className="mx-auto flex max-w-5xl items-center gap-4 px-4 py-4">
          <img
            src="https://i.loli.net/2021/02/11/JLHnIjOvFl7PC4o.png"
            alt="comical-tool"
            className="h-12 w-12 rounded border-2 border-ink bg-white"
          />
          <h1 className="text-2xl font-black">comical-tool</h1>
          <nav className="ml-auto flex gap-2">
            {tabs.map((item) => (
              <button
                key={item}
                onClick={() => setTab(item)}
                className={`border-2 border-ink px-3 py-2 font-bold ${
                  tab === item ? "bg-comicRed text-white" : "bg-white"
                }`}
              >
                {label(item)}
              </button>
            ))}
          </nav>
        </div>
      </header>

      <section className="mx-auto grid max-w-5xl gap-6 px-4 py-8 md:grid-cols-[1fr_0.9fr]">
        <div className="border-4 border-ink bg-white p-5 shadow-[8px_8px_0_#211815]">
          {tab === "short" && <ShortForm onResult={setResult} />}
          {tab === "image" && <UploadForm kind="image" onResult={setResult} />}
          {tab === "clip" && <ClipForm onResult={setResult} />}
          {tab === "file" && <UploadForm kind="file" onResult={setResult} />}
        </div>
        <pre className="min-h-80 overflow-auto border-4 border-ink bg-ink p-4 text-sm text-comicYellow">
          {result || "API result"}
        </pre>
      </section>
    </main>
  );
}

function ShortForm({ onResult }: { onResult: (value: string) => void }) {
  const [target, setTarget] = useState("");
  const [slug, setSlug] = useState("");
  const [ttl, setTTL] = useState("168h");
  return (
    <form onSubmit={(event) => postJSON(event, "/api/short-links", { target_url: target, custom_slug: slug, ttl }, onResult)} className="grid gap-4">
      <h2 className="text-xl font-black">短链接</h2>
      <Text value={target} onChange={setTarget} placeholder="https://example.com" />
      <Text value={slug} onChange={setSlug} placeholder="自定义短链，可空" />
      <Text value={ttl} onChange={setTTL} placeholder="过期时间，如 168h" />
      <Submit>创建短链</Submit>
    </form>
  );
}

function ClipForm({ onResult }: { onResult: (value: string) => void }) {
  const [content, setContent] = useState("");
  const [password, setPassword] = useState("");
  const [ttl, setTTL] = useState("1h");
  return (
    <form onSubmit={(event) => postJSON(event, "/api/clip", { content, password, ttl, link: true }, onResult)} className="grid gap-4">
      <h2 className="text-xl font-black">临时剪贴板</h2>
      <textarea value={content} onChange={(event) => setContent(event.target.value)} className="min-h-40 border-2 border-ink p-3" />
      <Text value={password} onChange={setPassword} placeholder="口令，可空" />
      <Text value={ttl} onChange={setTTL} placeholder="过期时间，如 1h" />
      <Submit>保存剪贴板</Submit>
    </form>
  );
}

function UploadForm({ kind, onResult }: { kind: "image" | "file"; onResult: (value: string) => void }) {
  const [ttl, setTTL] = useState(kind === "image" ? "720h" : "168h");
  return (
    <form onSubmit={(event) => upload(event, kind, ttl, onResult)} className="grid gap-4">
      <h2 className="text-xl font-black">{kind === "image" ? "图床" : "文件暂存"}</h2>
      <input name="file" type="file" className="border-2 border-ink p-3" />
      <Text value={ttl} onChange={setTTL} placeholder="过期时间" />
      <Submit>上传</Submit>
    </form>
  );
}

function Text(props: { value: string; onChange: (value: string) => void; placeholder: string }) {
  return <input value={props.value} onChange={(event) => props.onChange(event.target.value)} placeholder={props.placeholder} className="border-2 border-ink p-3" />;
}

function Submit({ children }: { children: React.ReactNode }) {
  return <button className="border-2 border-ink bg-comicRed px-4 py-3 font-black text-white">{children}</button>;
}

async function postJSON(event: FormEvent, path: string, body: unknown, onResult: (value: string) => void) {
  event.preventDefault();
  const res = await fetch(apiBase + path, { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(body) });
  onResult(JSON.stringify(await res.json(), null, 2));
}

async function upload(event: FormEvent<HTMLFormElement>, kind: "image" | "file", ttl: string, onResult: (value: string) => void) {
  event.preventDefault();
  const data = new FormData(event.currentTarget);
  data.set("ttl", ttl);
  data.set("link", "true");
  const path = kind === "image" ? "/api/images" : "/api/files";
  const res = await fetch(apiBase + path, { method: "POST", body: data });
  onResult(JSON.stringify(await res.json(), null, 2));
}

function label(tab: Tab) {
  return { short: "短链接", image: "图床", clip: "剪贴板", file: "文件" }[tab];
}
