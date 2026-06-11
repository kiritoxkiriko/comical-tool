"use client";

import { AnimatePresence, motion } from "framer-motion";
import { CheckCircle2, XCircle } from "lucide-react";

import type { ToastMessage } from "@/components/tool-types";

export function ToastStack({ items, onClose }: { items: ToastMessage[]; onClose: (id: number) => void }) {
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
              initial={{ opacity: 0, y: -10, scale: 0.98 }}
              animate={{ opacity: 1, y: 0, scale: 1 }}
              exit={{ opacity: 0, y: -8, scale: 0.98 }}
              transition={{ duration: 0.18 }}
              onClick={() => onClose(item.id)}
              className="pointer-events-auto grid w-full grid-cols-[auto_1fr] gap-3 rounded-2xl border border-ink/10 bg-white p-4 text-left shadow-lg shadow-ink/10"
            >
              <Icon
                className={item.kind === "success" ? "mt-0.5 h-5 w-5 text-[#1f9d55]" : "mt-0.5 h-5 w-5 text-comicRed"}
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
