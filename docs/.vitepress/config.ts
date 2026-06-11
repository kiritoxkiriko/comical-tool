import { defineConfig } from "vitepress";

export default defineConfig({
  title: "comical-tool",
  description: "Utility platform for short-lived resources",
  themeConfig: {
    nav: [{ text: "Guide", link: "/guide/" }],
    sidebar: [
      { text: "Guide", link: "/guide/" },
      { text: "Development Plan", link: "/development-plan" }
    ]
  }
});
