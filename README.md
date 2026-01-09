# CareFinder - 尋找照護者平台

一個使用 Golang (Gin) 後端和 Vue 3 (Vite) 前端打造的照護者媒合平台，採用前後端分離架構，具備完整的用戶、照護者及管理員功能。

## 功能特色

- **前後端分離架構**: 後端 Go-Gin API，前端 Vue 3 SPA。
- **用戶系統**:
  - 支援使用者、照護者、管理員三種角色。
  - JWT 身份驗證與註冊、登入功能。
  - 使用者個人資料管理。
- **照護者檔案**:
  - 建立與編輯詳細的個人檔案（含自我介紹、服務地點等）。
  - 設定可服務時間段。
  - 上傳與管理個人證照。
- **搜尋與配對**:
  - 依地點、性別、評分、可服務時間段等多維度篩選照護者。
  - 搜尋結果僅顯示持有「已審核」證照的照護者，確保品質。
- **合約管理**:
  - 使用者可向照護者發起服務合約邀請。
  - 完整的合約生命週期管理：從 `待處理`、`進行中`、`已完成` 到 `續約`。
  - 內建合約週期驗證（最長 6 個月）。
- **評價系統**:
  - 服務完成後，使用者可對照護者進行評分與留言。
  - 系統即時計算並更新照護者的平均評分。
- **管理員功能**:
  - 提供完整的管理後台介面。
  - 支援證照審核（通過/拒絕），確保平台上的照護者都經過驗證。
  - 提供用戶管理功能（查看列表、啟用/停用帳號）。
  - 支援服務合約管理（查看所有合約、修改狀態）。
- **響應式前端介面**:
  - 使用 Vue 3 + Pinia + Vue Router 打造的單頁應用。
  - 採用 Tailwind CSS，在不同裝置上皆有良好體驗。
  - 根據用戶角色顯示不同的導航與操作介面。

## 技術棧

### 後端
- **Golang** - 主要開發語言
- **Gin** - Web 框架
- **GORM** - ORM 框架
- **SQLite** - 資料庫（可切換為 PostgreSQL/MySQL）
- **JWT** - 身份驗證

### 前端
- **Vue 3** - 前端框架
- **Vue Router** - 路由管理
- **Pinia** - 狀態管理
- **Axios** - HTTP 客戶端
- **Tailwind CSS** - CSS 框架
- **Vite** - 構建工具

## 快速開始

### 前置需求

- Go 1.22 或更高版本
- Node.js 18 或更高版本
- npm 或 yarn
- 環境變數文件 `.env`

### 安裝步驟

1. **克隆專案**
```bash
git clone <repository-url>
cd CareFinder
```

2. **設置後端環境變數**

在項目根目錄創建 `.env` 文件：
```env
JWT_SECRET=your-secret-key-here
TOKEN_HOUR_LIFESPAN=24
```

3. **安裝後端依賴**
```bash
go mod download
```

4. **安裝前端依賴**
```bash
cd frontend
npm install
cd ..
```

5. **運行後端服務器**
```bash
go run main.go
```

後端服務器將在 `http://localhost:8080` 啟動

6. **運行前端開發伺服器**

在新的終端視窗中：
```bash
cd frontend
npm run dev
```

前端開發伺服器將在 `http://localhost:5173` 啟動

7. **訪問應用**

打開瀏覽器訪問：`http://localhost:5173`

## 專案結構

```
CareFinder/
├── controllers/          # 控制器（處理業務邏輯）
│   ├── auth.go          # 認證相關
│   ├── user.go          # 用戶管理
│   ├── caregiver.go     # 照護者管理
│   ├── contract.go      # 合約管理
│   └── review.go        # 評價管理
├── database/            # 資料庫連接
│   └── database.go
├── middleware/          # 中間件
│   └── auth.go          # JWT 認證中間件
├── models/              # 資料模型
│   └── models.go
├── routes/              # 路由配置
│   └── routes.go
├── frontend/            # Vue 3 前端應用
│   ├── src/
│   │   ├── components/  # Vue 組件
│   │   ├── views/       # 頁面組件
│   │   ├── router/      # 路由配置
│   │   ├── stores/      # Pinia 狀態管理
│   │   ├── services/    # API 服務層
│   │   └── assets/      # 靜態資源
│   ├── package.json
│   └── vite.config.js
├── static/              # 舊版前端文件（已棄用，保留作為參考）
├── utils/               # 工具函數
│   ├── password.go      # 密碼加密
│   └── token/           # JWT Token 處理
├── main.go              # 應用程式入口
├── go.mod               # Go 模組定義
└── README.md            # 本文件
```

## API 端點

### 認證
- `POST /api/v1/auth/register` - 註冊
- `POST /api/v1/auth/login` - 登入

### 用戶管理
- `GET /api/v1/users/{id}` - 取得用戶資料
- `PUT /api/v1/users/{id}` - 更新用戶資料

### 照護者管理
- `POST /api/v1/caregivers/profile` - 建立照護者檔案
- `PUT /api/v1/caregivers/profile` - 更新照護者檔案
- `GET /api/v1/caregivers/{id}` - 取得照護者檔案
- `GET /api/v1/caregivers/search` - 搜尋照護者
- `PUT /api/v1/caregivers/availability` - 更新可服務時間

### 證照管理
- `POST /api/v1/caregivers/licenses` - 上傳新證照
- `GET /api/v1/caregivers/licenses` - 取得我的證照
- `DELETE /api/v1/caregivers/licenses/:id` - 刪除證照

### 合約管理
- `POST /api/v1/contracts` - 發起合約
- `GET /api/v1/contracts` - 取得所有合約
- `GET /api/v1/contracts/{id}` - 取得單一合約
- `PUT /api/v1/contracts/{id}/accept` - 接受合約
- `PUT /api/v1/contracts/{id}/complete` - 完成合約
- `POST /api/v1/contracts/{id}/renew` - 續約

### 評價管理
- `POST /api/v1/reviews` - 提交評價
- `GET /api/v1/reviews/caregivers/{id}` - 取得照護者評價

### 管理員 (Admin)
- `GET /api/v1/admin/licenses` - 取得所有（或待審核）證照
- `PUT /api/v1/admin/licenses/{id}/status` - 更新證照狀態（審核）
- `GET /api/v1/admin/users` - 取得所有用戶
- `PUT /api/v1/admin/users/{id}/status` - 更新用戶狀態（啟用/停用）
- `GET /api/v1/admin/contracts` - 取得所有合約
- `PUT /api/v1/admin/contracts/{id}/status` - 更新合約狀態

## 使用說明

### 註冊帳號

1. 訪問註冊頁面
2. 填寫電子郵件和密碼
3. 選擇身份（使用者/照護者）
4. 提交註冊

### 照護者註冊後

1. 登入後進入儀表板
2. 建立照護者檔案（填寫姓名、電話、地址、收費標準等）
3. 設定可服務時間段
4. 上傳證照（等待管理員審核）。

### 使用者搜尋照護者

1. 登入後進入「搜尋照護者」頁面
2. 使用篩選條件搜尋（地點、性別、評分、時間等）
3. 查看照護者詳細資料
4. 發起合約邀請

### 合約管理

- **使用者**：可以發起合約、查看合約、完成服務、續約、評價
- **照護者**：可以接受合約、查看合約、完成服務

## 開發說明

### 前端開發

前端使用 Vue 3 + Vite，開發時會自動熱重載。詳細說明請參考 `frontend/README.md`。

開發模式：
```bash
cd frontend
npm run dev
```

構建生產版本：
```bash
cd frontend
npm run build
```

### 資料庫

目前使用 SQLite，資料庫文件為 `carefinder.db`。如需切換到 PostgreSQL 或 MySQL，請修改 `database/database.go`。

## 待實現功能

目前專案已完成核心功能，未來規劃實現以下項目（詳見 `todo.md`）：

- **進階授權**: 實作角色權限檢查中間件，取代控制器中的手動檢查。
- **功能完善**:
  - 搜尋功能支援 `service_type` 篩選。
  - 驗證合約時間是否在照護者可服務時間段內。
  - 服務完成需由雙方確認。
- **管理員後台**:
  - 擴充管理員功能，如：編輯用戶詳細資料、查看系統日誌等。
  - 建立系統營運數據儀表板（用戶統計、合約統計等）。
- **通知系統**:
  - 透過 Email 或站內信通知合約狀態變更、新邀請、證照審核結果等。
- **技術優化**:
  - 整合雲端儲存（如 AWS S3）取代本地檔案上傳。
  - 建立統一的錯誤處理與日誌記錄機制。
  - 撰寫單元測試與整合測試。
  - 產生 API 文件（Swagger/OpenAPI）。