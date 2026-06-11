import "./globals.css";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "comical-tool",
  description: "Short links, image hosting, clipboard, and temporary files"
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="zh-CN">
      <body>{children}</body>
    </html>
  );
}
