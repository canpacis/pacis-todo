import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";

const __dirname = dirname(fileURLToPath(import.meta.url));

export default defineConfig({
  plugins: [
    tailwindcss(),
  ],
  server: {
    hmr: {
      port: 5173,
    },
  },
  build: {
    manifest: true,
    lib: {
      entry: {
        script: resolve(__dirname, "src/web/main.ts"),
        stream: resolve(__dirname, "src/web/stream.ts"),
        style: resolve(__dirname, "src/web/style.css"),
      },
      name: "App",
      fileName: "app",
    },
    outDir: "build/static",
    cssCodeSplit: true,
    rollupOptions: {
      output: {
        entryFileNames: "[name].[hash].js",
        chunkFileNames: "[name].[hash].js",
        assetFileNames: "[name].[hash].[ext]",
      },
    },
  },
});
