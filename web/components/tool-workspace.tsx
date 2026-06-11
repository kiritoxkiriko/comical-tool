"use client";

import { AnimatePresence, motion } from "framer-motion";
import {
  CheckCircle2,
  Clipboard,
  Clock3,
  FileUp,
  ImageUp,
  Link2,
  LockKeyhole,
  Sparkles,
  Trash2,
  XCircle
} from "lucide-react";
import Image from "next/image";
import { useState } from "react";

import { ToolForm } from "@/components/tool-forms";
import type { ToastMessage, ToolMeta, ToolTab } from "@/components/tool-types";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

const modules: ToolMeta[] = [
  {
    value: "short",
    label: "短链接",
    title: "生成短链接",
    desc: "自定义 slug、过期时间和短域映射。",
    icon: Link2,
    image: "/huaji/original.png"
  },
  {
    value: "image",
    label: "图床",
    title: "上传图片",
    desc: "临时托管图片，并按需生成短链。",
    icon: ImageUp,
    image: "/huaji/eager.jpg"
  },
  {
    value: "clip",
    label: "剪贴板",
    title: "临时剪贴板",
    desc: "适合口令分享、少次读取和短时文本。",
    icon: Clipboard,
    image: "/huaji/watch.jpg"
  },
  {
    value: "file",
    label: "文件暂存",
    title: "暂存文件",
    desc: "文件默认短期留存，过期后可自动清理。",
    icon: FileUp,
    image: "/huaji/surprised.png"
  }
];

const details = [
  { icon: Clock3, label: "TTL", value: "所有模块支持过期时间" },
  { icon: LockKeyhole, label: "Guest", value: "已预留用户体系边界" },
  { icon: Trash2, label: "Cleanup", value: "支持删除、撤销和清理" }
] as const;

export function ToolWorkspace() {
  const [tab, setTab] = useState<ToolTab>("short");
  const [loading, setLoading] = useState(false);
  const [toasts, setToasts] = useState<ToastMessage[]>([]);
  const active = modules.find((item) => item.value === tab) ?? modules[0];

  function notify(message: Omit<ToastMessage, "id">) {
    const id = Date.now();
    setToasts((items) => [...items, { ...message, id }].slice(-3));
    window.setTimeout(() => removeToast(id), 4200);
  }

  function removeToast(id: number) {
    setToasts((items) => items.filter((item) => item.id !== id));
  }

  return (
    <main className="min-h-screen bg-[#fbfaf6] text-ink">
      <div className="mx-auto flex min-h-screen w-full max-w-6xl flex-col px-4 py-4 sm:px-6 lg:px-8">
        <Header />
        <section className="flex flex-1 py-6">
          <WorkspacePanel
            active={active}
            tab={tab}
            loading={loading}
            setTab={setTab}
            notify={notify}
            setLoading={setLoading}
          />
        </section>
      </div>
      <ToastStack items={toasts} onClose={removeToast} />
    </main>
  );
}

function Header() {
  return (
    <header className="flex items-center gap-3 border-b border-ink/10 py-4">
      <Image
        src="/logo.png"
        alt="comical-tool"
        width={40}
        height={40}
        className="h-10 w-10 rounded-xl object-cover"
        priority
      />
      <div className="min-w-0">
        <h1 className="text-xl font-black tracking-normal sm:text-2xl">comical-tool</h1>
        <p className="text-xs font-medium text-ink/50 sm:text-sm">短链、图床、剪贴板、文件暂存</p>
      </div>
      <div className="ml-auto hidden items-center gap-2 rounded-full border border-ink/10 bg-white px-3 py-2 text-sm font-semibold text-ink/65 sm:flex">
        <span className="h-2 w-2 rounded-full bg-[#27c46a]" />
        tool.sqlboy.me
      </div>
    </header>
  );
}

function WorkspacePanel(props: {
  active: ToolMeta;
  tab: ToolTab;
  loading: boolean;
  setTab: (value: ToolTab) => void;
  notify: (message: Omit<ToastMessage, "id">) => void;
  setLoading: (value: boolean) => void;
}) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.26, ease: "easeOut" }}
      className="flex w-full flex-col rounded-[28px] border border-ink/10 bg-white p-4 shadow-sm sm:p-6 lg:min-h-[620px]"
    >
      <WorkspaceHeader active={props.active} loading={props.loading} />
      <Tabs
        value={props.tab}
        onValueChange={(value) => props.setTab(value as ToolTab)}
        className="mt-6 flex flex-1 flex-col"
      >
        <ToolTabs />
        <ToolBody {...props} />
      </Tabs>
    </motion.div>
  );
}

function WorkspaceHeader({ active, loading }: { active: ToolMeta; loading: boolean }) {
  return (
    <div className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_220px]">
      <div>
        <div className="mb-3 inline-flex items-center gap-2 rounded-full bg-comicYellow/25 px-3 py-1 text-xs font-bold text-ink/70">
          <Sparkles className="h-3.5 w-3.5 text-comicRed" />
          {loading ? "processing" : "guest workspace"}
        </div>
        <h2 className="max-w-2xl text-4xl font-black tracking-normal text-ink sm:text-5xl">临时资源处理台</h2>
        <p className="mt-3 max-w-xl text-sm leading-6 text-ink/58">
          用最少步骤生成短链、上传临时资源、保存一次性文本。所有操作保留过期、删除和后续用户体系扩展边界。
        </p>
      </div>
      <motion.div
        key={active.value}
        initial={{ opacity: 0, scale: 0.96, rotate: -2 }}
        animate={{ opacity: 1, scale: 1, rotate: 0 }}
        transition={{ duration: 0.22 }}
        className="hidden items-end justify-end lg:flex"
      >
        <Image src={active.image} alt="" width={160} height={160} className="h-36 w-36 rounded-3xl object-contain" />
      </motion.div>
    </div>
  );
}

function ToolTabs() {
  return (
    <TabsList className="grid w-full grid-cols-2 gap-1 rounded-2xl bg-[#f3f0e8] p-1 shadow-none sm:grid-cols-4">
      {modules.map((item) => {
        const Icon = item.icon;
        return (
          <TabsTrigger
            key={item.value}
            value={item.value}
            className="h-11 rounded-xl px-2 text-xs font-bold text-ink/55 data-[state=active]:bg-white data-[state=active]:text-ink data-[state=active]:shadow-sm sm:text-sm"
          >
            <Icon className="h-4 w-4" />
            {item.label}
          </TabsTrigger>
        );
      })}
    </TabsList>
  );
}

function ToolBody(props: {
  active: ToolMeta;
  tab: ToolTab;
  notify: (message: Omit<ToastMessage, "id">) => void;
  setLoading: (value: boolean) => void;
}) {
  return (
    <div className="relative mt-7 flex-1">
      <AnimatePresence mode="wait">
        {modules.map((item) => (
          <TabsContent key={item.value} value={item.value} forceMount hidden={props.tab !== item.value}>
            {props.tab === item.value && (
              <motion.div
                initial={{ opacity: 0, y: 12 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -8 }}
                transition={{ duration: 0.2, ease: "easeOut" }}
                className="grid gap-8 lg:grid-cols-[minmax(0,1fr)_210px]"
              >
                <ToolForm tab={item.value} meta={item} notify={props.notify} setLoading={props.setLoading} />
                <ModuleAside active={item} />
              </motion.div>
            )}
          </TabsContent>
        ))}
      </AnimatePresence>
    </div>
  );
}

function ModuleAside({ active }: { active: ToolMeta }) {
  return (
    <aside className="hidden border-l border-ink/10 pl-6 lg:block">
      <p className="text-xs font-bold uppercase text-ink/35">Current module</p>
      <h3 className="mt-2 text-2xl font-black tracking-normal">{active.label}</h3>
      <p className="mt-2 text-sm leading-6 text-ink/55">{active.desc}</p>
      <div className="mt-8 space-y-5">
        {details.map((item) => {
          const Icon = item.icon;
          return (
            <div key={item.label} className="flex gap-3">
              <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-[#f3f0e8]">
                <Icon className="h-4 w-4 text-comicRed" />
              </div>
              <div>
                <p className="text-sm font-bold">{item.label}</p>
                <p className="text-sm text-ink/50">{item.value}</p>
              </div>
            </div>
          );
        })}
      </div>
    </aside>
  );
}

function ToastStack({ items, onClose }: { items: ToastMessage[]; onClose: (id: number) => void }) {
  return (
    <div
      aria-live="polite"
      className="pointer-events-none fixed inset-x-4 top-4 z-50 flex flex-col items-end gap-3 sm:inset-x-auto sm:right-5 sm:w-[360px]"
    >
      <AnimatePresence>
        {items.map((item) => {
          const Icon = item.kind === "success" ? CheckCircle2 : XCircle;
          return (
            <motion.button
              type="button"
              key={item.id}
              initial={{ opacity: 0, y: -12, scale: 0.98 }}
              animate={{ opacity: 1, y: 0, scale: 1 }}
              exit={{ opacity: 0, y: -8, scale: 0.98 }}
              transition={{ duration: 0.18 }}
              onClick={() => onClose(item.id)}
              className="pointer-events-auto grid w-full grid-cols-[auto_1fr] gap-3 rounded-2xl border border-ink/10 bg-white p-4 text-left shadow-lg shadow-ink/10"
            >
              <Icon
                className={item.kind === "success" ? "mt-0.5 h-5 w-5 text-[#20a65a]" : "mt-0.5 h-5 w-5 text-comicRed"}
              />
              <span>
                <span className="block text-sm font-black">{item.title}</span>
                <span className="mt-1 block break-words text-sm leading-5 text-ink/55">{item.description}</span>
              </span>
            </motion.button>
          );
        })}
      </AnimatePresence>
    </div>
  );
}
