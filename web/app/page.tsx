"use client";

import { FormEvent, useState } from "react";
import { AnimatePresence, motion } from "framer-motion";
import { Clipboard, FileUp, ImageUp, Link2, Loader2, Sparkles } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Textarea } from "@/components/ui/textarea";

const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

const modules = [
  { value: "short", label: "短链接", icon: Link2, image: "/huaji/original.png" },
  { value: "image", label: "图床", icon: ImageUp, image: "/huaji/eager.jpg" },
  { value: "clip", label: "剪贴板", icon: Clipboard, image: "/huaji/watch.jpg" },
  { value: "file", label: "文件暂存", icon: FileUp, image: "/huaji/surprised.png" }
] as const;

type Tab = (typeof modules)[number]["value"];

export default function Home() {
  const [tab, setTab] = useState<Tab>("short");
  const [result, setResult] = useState("等待操作结果");
  const [loading, setLoading] = useState(false);

  return (
    <main className="min-h-screen overflow-hidden bg-[#fff7df] text-ink">
      <div className="absolute inset-x-0 top-0 h-72 bg-[radial-gradient(circle_at_20%_20%,#ffd64d_0,#ffd64d_24%,transparent_25%),linear-gradient(120deg,#fff7df,#ffe9ec_45%,#eefbf6)]" />
      <div className="relative mx-auto max-w-6xl px-4 py-6">
        <header className="flex items-center gap-3">
          <img src="/logo.png" alt="comical-tool" className="h-11 w-11 rounded-xl border border-ink/10 bg-white object-cover shadow-sm" />
          <div>
            <h1 className="text-2xl font-black tracking-normal">comical-tool</h1>
            <p className="text-sm text-ink/60">短期资源工具台</p>
          </div>
          <div className="ml-auto hidden items-center gap-2 rounded-full border border-ink/10 bg-white/80 px-3 py-2 text-sm font-semibold text-ink/70 shadow-sm sm:flex">
            <Sparkles className="h-4 w-4 text-comicRed" />
            tool.sqlboy.me
          </div>
        </header>

        <section className="grid gap-6 py-8 lg:grid-cols-[1fr_380px]">
          <div className="rounded-2xl border border-ink/10 bg-white/85 p-4 shadow-xl shadow-comicRed/5 backdrop-blur md:p-6">
            <Tabs value={tab} onValueChange={(value) => setTab(value as Tab)}>
              <TabsList className="mb-6 flex w-full flex-wrap justify-start gap-1">
                {modules.map((item) => {
                  const Icon = item.icon;
                  return (
                    <TabsTrigger key={item.value} value={item.value} className="flex items-center gap-2">
                      <Icon className="h-4 w-4" />
                      {item.label}
                    </TabsTrigger>
                  );
                })}
              </TabsList>

              <AnimatePresence mode="wait">
                {modules.map((item) => (
                  <TabsContent key={item.value} value={item.value} forceMount hidden={tab !== item.value}>
                    {tab === item.value && (
                      <motion.div
                        initial={{ opacity: 0, y: 8 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: -8 }}
                        transition={{ duration: 0.18 }}
                        className="grid gap-6 md:grid-cols-[1fr_180px]"
                      >
                        <ToolForm tab={item.value} setResult={setResult} setLoading={setLoading} />
                        <aside className="hidden rounded-xl bg-comicYellow/25 p-4 md:block">
                          <img src={item.image} alt="" className="mx-auto h-32 w-32 rounded-xl object-contain" />
                          <p className="mt-4 text-sm font-medium text-ink/60">{caption(item.value)}</p>
                        </aside>
                      </motion.div>
                    )}
                  </TabsContent>
                ))}
              </AnimatePresence>
            </Tabs>
          </div>

          <ResultPanel loading={loading} result={result} />
        </section>
      </div>
    </main>
  );
}

function ToolForm(props: {
  tab: Tab;
  setResult: (value: string) => void;
  setLoading: (value: boolean) => void;
}) {
  if (props.tab === "short") return <ShortForm {...props} />;
  if (props.tab === "clip") return <ClipForm {...props} />;
  return <UploadForm kind={props.tab} {...props} />;
}

function ShortForm({ setResult, setLoading }: FormProps) {
  const [target, setTarget] = useState("");
  const [slug, setSlug] = useState("");
  const [ttl, setTTL] = useState("168h");
  return (
    <form onSubmit={(event) => submitJSON(event, "/api/short-links", { target_url: target, custom_slug: slug, ttl }, setResult, setLoading)} className="grid gap-4">
      <SectionTitle title="生成短链接" desc="支持自定义 slug、过期时间和独立短域映射。" />
      <Input value={target} onChange={(event) => setTarget(event.target.value)} placeholder="https://example.com" />
      <div className="grid gap-3 sm:grid-cols-2">
        <Input value={slug} onChange={(event) => setSlug(event.target.value)} placeholder="自定义 slug，可空" />
        <Input value={ttl} onChange={(event) => setTTL(event.target.value)} placeholder="168h" />
      </div>
      <Button variant="accent">创建短链</Button>
    </form>
  );
}

function ClipForm({ setResult, setLoading }: FormProps) {
  const [content, setContent] = useState("");
  const [password, setPassword] = useState("");
  const [ttl, setTTL] = useState("1h");
  const [maxVisits, setMaxVisits] = useState("5");
  return (
    <form onSubmit={(event) => submitJSON(event, "/api/clip", { content, password, ttl, max_visits: Number(maxVisits), link: true }, setResult, setLoading)} className="grid gap-4">
      <SectionTitle title="临时剪贴板" desc="适合口令分享、少次读取和短时效内容。" />
      <Textarea value={content} onChange={(event) => setContent(event.target.value)} placeholder="粘贴临时文本" />
      <div className="grid gap-3 sm:grid-cols-3">
        <Input value={password} onChange={(event) => setPassword(event.target.value)} placeholder="口令，可空" />
        <Input value={ttl} onChange={(event) => setTTL(event.target.value)} placeholder="1h" />
        <Input value={maxVisits} onChange={(event) => setMaxVisits(event.target.value)} placeholder="访问次数" />
      </div>
      <Button variant="accent">保存剪贴板</Button>
    </form>
  );
}

function UploadForm({ kind, setResult, setLoading }: FormProps & { kind: "image" | "file" }) {
  const [ttl, setTTL] = useState(kind === "image" ? "720h" : "168h");
  return (
    <form onSubmit={(event) => upload(event, kind, ttl, setResult, setLoading)} className="grid gap-4">
      <SectionTitle title={kind === "image" ? "上传图片" : "暂存文件"} desc="上传后自动写入元数据，可同时生成短链。" />
      <Input name="file" type="file" />
      <Input value={ttl} onChange={(event) => setTTL(event.target.value)} placeholder="过期时间" />
      <Button variant="accent">{kind === "image" ? "上传图片" : "上传文件"}</Button>
    </form>
  );
}

function ResultPanel({ loading, result }: { loading: boolean; result: string }) {
  return (
    <aside className="rounded-2xl border border-ink/10 bg-ink p-4 text-white shadow-xl shadow-ink/10">
      <div className="mb-3 flex items-center justify-between">
        <h2 className="text-sm font-bold text-comicYellow">响应结果</h2>
        {loading && <Loader2 className="h-4 w-4 animate-spin text-comicYellow" />}
      </div>
      <pre className="min-h-80 overflow-auto whitespace-pre-wrap rounded-xl bg-black/25 p-4 text-sm leading-6 text-[#fff2b8]">
        {result}
      </pre>
      <img src="/huaji/ghost.png" alt="" className="ml-auto mt-4 h-20 w-20 opacity-80" />
    </aside>
  );
}

function SectionTitle({ title, desc }: { title: string; desc: string }) {
  return (
    <div>
      <h2 className="text-xl font-black tracking-normal">{title}</h2>
      <p className="mt-1 text-sm text-ink/60">{desc}</p>
    </div>
  );
}

async function submitJSON(event: FormEvent, path: string, body: unknown, setResult: (value: string) => void, setLoading: (value: boolean) => void) {
  event.preventDefault();
  await run(setResult, setLoading, async () => {
    const res = await fetch(apiBase + path, { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(body) });
    return res.json();
  });
}

async function upload(event: FormEvent<HTMLFormElement>, kind: "image" | "file", ttl: string, setResult: (value: string) => void, setLoading: (value: boolean) => void) {
  event.preventDefault();
  await run(setResult, setLoading, async () => {
    const data = new FormData(event.currentTarget);
    data.set("ttl", ttl);
    data.set("link", "true");
    const path = kind === "image" ? "/api/images" : "/api/files";
    const res = await fetch(apiBase + path, { method: "POST", body: data });
    return res.json();
  });
}

async function run(setResult: (value: string) => void, setLoading: (value: boolean) => void, action: () => Promise<unknown>) {
  setLoading(true);
  try {
    setResult(JSON.stringify(await action(), null, 2));
  } catch (error) {
    setResult(JSON.stringify({ error: String(error) }, null, 2));
  } finally {
    setLoading(false);
  }
}

function caption(tab: Tab) {
  return {
    short: "把长链接压成好转发的短路径。",
    image: "临时图片托管，短期分享更顺手。",
    clip: "一次性文本和口令内容放这里。",
    file: "文件暂存默认 7 天，按需过期。"
  }[tab];
}

type FormProps = {
  setResult: (value: string) => void;
  setLoading: (value: boolean) => void;
};
