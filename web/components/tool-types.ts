import type { LucideIcon } from "lucide-react";

export type ToolTab = "short" | "image" | "clip" | "file";

export type ToolMeta = {
  value: ToolTab;
  label: string;
  title: string;
  desc: string;
  icon: LucideIcon;
  image: string;
};

export type ToastKind = "success" | "error";

export type ToastMessage = {
  id: number;
  kind: ToastKind;
  title: string;
  description: string;
};
