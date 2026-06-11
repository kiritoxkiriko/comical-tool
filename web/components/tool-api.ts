import type { ToastMessage } from "@/components/tool-types";

export const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export async function runToolAction(
  notify: (message: Omit<ToastMessage, "id">) => void,
  setLoading: (value: boolean) => void,
  action: () => Promise<unknown>
) {
  setLoading(true);
  try {
    const payload = await action();
    notify({ kind: "success", title: "操作成功", description: primaryMessage(payload) });
    return payload;
  } catch (error) {
    notify({ kind: "error", title: "操作失败", description: error instanceof Error ? error.message : String(error) });
    return undefined;
  } finally {
    setLoading(false);
  }
}

export async function parseResponse<T = unknown>(res: Response): Promise<T> {
  const text = await res.text();
  const payload = parsePayload(text);
  if (!res.ok) {
    throw new Error(errorMessage(payload));
  }
  return payload as T;
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
