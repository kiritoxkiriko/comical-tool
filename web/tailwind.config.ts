import type { Config } from "tailwindcss";

const config: Config = {
  content: ["./app/**/*.{ts,tsx}", "./components/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        comicYellow: "#ffd64d",
        comicRed: "#e83f37",
        ink: "#211815"
      }
    }
  },
  plugins: []
};

export default config;
