import "./globals.css";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "comical-tool",
  description: "Short links, image hosting, clipboard, and temporary files",
  icons: {
    icon: "/logo.png",
    apple: "/logo.png"
  }
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="zh-CN">
      <body>{children}</body>
    </html>
  );
}
