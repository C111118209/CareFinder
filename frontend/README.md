# CareFinder Frontend

Vue 3 前端應用程式，使用 Vite 作為構建工具。

## 技術棧

- Vue 3
- Vue Router
- Pinia (狀態管理)
- Axios (HTTP 客戶端)
- Tailwind CSS
- Vite

## 安裝依賴

```bash
npm install
```

## 開發

啟動開發伺服器（會在 http://localhost:5173 運行）：

```bash
npm run dev
```

## 構建生產版本

```bash
npm run build
```

構建後的檔案會在 `dist` 目錄中。

## 預覽生產構建

```bash
npm run preview
```

## 環境變數

複製 `.env.example` 為 `.env` 並根據需要修改：

```
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

## 項目結構

```
frontend/
├── src/
│   ├── assets/          # 靜態資源（CSS 等）
│   ├── components/      # 可重用的 Vue 組件
│   ├── router/          # Vue Router 配置
│   ├── services/        # API 服務層
│   ├── stores/          # Pinia 狀態管理
│   ├── views/           # 頁面組件
│   ├── App.vue          # 根組件
│   └── main.js          # 應用程式入口
├── index.html
├── package.json
├── vite.config.js
└── tailwind.config.js
```

## 後端 API

確保後端 API 服務正在運行（預設為 http://localhost:8080）。

開發模式下，Vite 會自動代理 `/api` 請求到後端伺服器（見 `vite.config.js`）。