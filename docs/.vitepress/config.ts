import { defineConfig } from "vitepress";

export default defineConfig({
  title: "comical-tool",
  description: "Utility platform for short-lived resources",
  themeConfig: {
    nav: [
      { text: "Guide", link: "/guide/" },
      { text: "API", link: "/api" },
      { text: "Cloudflare", link: "/cloudflare" }
    ],
    sidebar: [
      { text: "Guide", link: "/guide/" },
      { text: "Quick Start", link: "/quick-start" },
      { text: "Project Structure", link: "/structure" },
      { text: "Configuration", link: "/configuration" },
      { text: "API", link: "/api" },
      { text: "CLI", link: "/cli" },
      { text: "Local Development", link: "/local-development" },
      { text: "Docker", link: "/docker" },
      { text: "Cloudflare", link: "/cloudflare" },
      { text: "Migrations", link: "/migrations" },
      { text: "Storage Backends", link: "/storage" },
      { text: "Troubleshooting", link: "/troubleshooting" },
      { text: "Development Plan", link: "/development-plan" }
    ]
  }
});
