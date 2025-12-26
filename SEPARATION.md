# 前後端分離說明

本專案已完成前後端分離，前端使用 Vue 3，後端使用 Go + Gin。

## 架構變更

### 後端變更

1. **移除靜態文件服務**
   - 移除了 `r.Static()` 和 `r.StaticFile()` 等靜態文件服務
   - 後端現在只提供 API 端點

2. **添加 CORS 支持**
   - 使用 `github.com/gin-contrib/cors` 中間件
   - 允許來自 `http://localhost:5173` 和 `http://localhost:3000` 的請求
   - 支持必要的 HTTP 方法和頭部

3. **API 端點保持不變**
   - 所有 API 端點仍在 `/api/v1` 路徑下
   - API 接口行為保持不變

### 前端變更

1. **從 HTML 遷移到 Vue 3**
   - 使用 Vite 作為構建工具
   - 使用 Vue Router 進行路由管理
   - 使用 Pinia 進行狀態管理

2. **項目結構**
   ```
   frontend/
   ├── src/
   │   ├── components/    # 可重用組件（Navbar 等）
   │   ├── views/         # 頁面組件
   │   ├── router/        # 路由配置
   │   ├── stores/        # Pinia stores（認證狀態等）
   │   ├── services/      # API 服務層
   │   └── assets/        # 靜態資源（CSS）
   ```

3. **API 調用**
   - 使用 Axios 進行 HTTP 請求
   - 開發模式下使用 Vite 代理轉發到後端
   - 自動添加 Authorization header

## 運行方式

### 開發環境

1. **啟動後端**
   ```bash
   go run main.go
   ```
   後端運行在 `http://localhost:8080`

2. **啟動前端**
   ```bash
   cd frontend
   npm install  # 首次運行需要
   npm run dev
   ```
   前端運行在 `http://localhost:5173`

3. **訪問應用**
   瀏覽器訪問 `http://localhost:5173`

### 生產環境

1. **構建前端**
   ```bash
   cd frontend
   npm run build
   ```
   構建結果在 `frontend/dist` 目錄

2. **部署選項**
   - 方案 1：使用 Nginx 等 Web 伺服器提供前端靜態文件，代理 API 請求到後端
   - 方案 2：後端提供前端靜態文件（需要重新配置路由）

## 配置

### 環境變數

前端可在 `.env` 文件中配置：
```
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

如果未設置，開發模式下會使用 Vite 代理（相對路徑 `/api/v1`）。

### CORS 配置

後端 CORS 配置在 `routes/routes.go` 中：
```go
config := cors.Config{
    AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
    AllowCredentials: true,
}
```

生產環境需要根據實際前端域名更新 `AllowOrigins`。

## 遷移注意事項

1. **舊的 HTML 文件**
   - `static/` 目錄中的 HTML 文件已不再使用
   - 可以保留作為參考，或刪除

2. **API 兼容性**
   - 所有 API 端點保持不變
   - 前端調用方式與原來相似，只是使用 Vue 組件和 Axios

3. **認證機制**
   - JWT Token 仍存儲在 localStorage
   - 認證流程保持不變

## 下一步

- [ ] 完善 Vue 組件功能
- [ ] 添加錯誤處理和用戶提示
- [ ] 優化 UI/UX
- [ ] 添加單元測試
- [ ] 配置生產環境部署