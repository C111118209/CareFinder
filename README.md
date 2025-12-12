# CareFinder - 尋找照護者平台

一個使用 Golang 後端和 Tailwind CSS 前端的照護者配對平台。

## 功能特色

- ✅ 用戶註冊與登入（支援使用者/照護者/管理員角色）
- ✅ 照護者檔案管理
- ✅ 智能搜尋與篩選照護者
- ✅ 合約管理（發起、接受、完成、續約）
- ✅ 評價與評分系統
- ✅ 響應式設計，使用 Tailwind CSS

## 技術棧

### 後端
- **Golang** - 主要開發語言
- **Gin** - Web 框架
- **GORM** - ORM 框架
- **SQLite** - 資料庫（可切換為 PostgreSQL/MySQL）
- **JWT** - 身份驗證

### 前端
- **HTML5** - 標記語言
- **Tailwind CSS** - 使用 CDN 版本
- **Vanilla JavaScript** - 無框架 JavaScript

## 快速開始

### 前置需求

- Go 1.22 或更高版本
- 環境變數文件 `.env`

### 安裝步驟

1. **克隆專案**
```bash
git clone <repository-url>
cd CareFinder
```

2. **設置環境變數**

創建 `.env` 文件：
```env
JWT_SECRET=your-secret-key-here
TOKEN_HOUR_LIFESPAN=24
```

3. **安裝 Go 依賴**
```bash
go mod download
```

4. **運行後端服務器**
```bash
go run main.go
```

服務器將在 `http://localhost:8080` 啟動

5. **訪問前端**

打開瀏覽器訪問：
- 首頁：`http://localhost:8080/`
- 登入：`http://localhost:8080/login.html`
- 註冊：`http://localhost:8080/register.html`

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
├── routes/               # 路由配置
│   └── routes.go
├── static/              # 靜態文件（前端）
│   ├── css/             # CSS 文件
│   ├── js/              # JavaScript 文件
│   └── *.html           # HTML 頁面
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
4. 上傳證照（待實現）

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

前端使用 Tailwind CSS CDN，無需額外構建步驟。如需使用本地 Tailwind：

```bash
npm install
npm run build-css
```

### 資料庫

目前使用 SQLite，資料庫文件為 `carefinder.db`。如需切換到 PostgreSQL 或 MySQL，請修改 `database/database.go`。

## 待實現功能

詳見 `todo.md` 文件，包括：
- 證照上傳與審核
- 管理員功能
- 通知系統
- 更多進階功能

## 授權

MIT License

## 貢獻

歡迎提交 Issue 和 Pull Request！

