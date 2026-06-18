import type { Config } from "tailwindcss";

export default {
  content: ["./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        // Dark, app-like palette tuned for mobile (matches the Notes dark theme).
        bg: "#0b0b0f",
        surface: "#16161d",
        surface2: "#1f1f29",
        border: "#2a2a36",
        accent: "#f5a623", // the amber highlight from the Notes screenshots
        accent2: "#3b82f6",
        muted: "#8a8a99",
      },
    },
  },
  plugins: [],
} satisfies Config;
