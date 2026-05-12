import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import tailwindcss from "@tailwindcss/vite";
import path from "path/win32";

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), tailwindcss()],
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
      // Если тебе нужно проксировать Swagger, добавь и его
      "/swagger": {
        target: "http://localhost:8080",
        changeOrigin: true,
        configure: (proxy) => {
          proxy.on("proxyReq", (proxyReq, req, res) => {
            // Если зашли совсем без слэша, делаем браузерный редирект
            if (req.url === "/swagger") {
              res.writeHead(301, { Location: "/swagger/" });
              res.end();
            }
          });
        },
        // Когда запрос уже со слэшем (/swagger/),
        // Vite незаметно для браузера спросит у Go именно index.html
        rewrite: (path) => path.replace(/\/swagger\/$/, "/swagger/index.html"),
      },
    },
  },
});
