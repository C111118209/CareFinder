# 尋找照護者 (CareFinder) - 平台規格

## 📋 平台功能規格 (Functional Specifications)

### 1. 核心使用者角色

- **平台管理者 (Admin):**
  - 管理所有用戶、照護者、服務合約、證照審核。
  - 查看系統營運數據與報告。
  - 管理網站內容（例如：公告、使用條款）。

- **使用者/家屬 (User):**
  - 註冊/登入/修改個人資料。
  - 瀏覽、搜尋、篩選照護者。
  - 發送服務需求或合約邀請。
  - 管理合約（簽訂、續約、取消）。
  - 服務完成後，對照護者進行留言與評分。

- **照護者 (Caregiver):**
  - 註冊/登入/修改個人資料。
  - 填寫/更新個人檔案 (包含：服務項目、經驗、證照、收費標準、可服務時間段)。
  - 上傳證照供平台審核。
  - 接收、接受或拒絕服務邀請。
  - 管理合約（簽訂、續約、完成）。

### 2. 照護者管理 (Caregiver Management)

- **照護者檔案:** 包含基本資料、自我介紹、服務項目、收費標準。
- **證照審核機制:** 照護者上傳證照，平台管理者審核（通過/拒絕）。未通過審核的照護者無法被使用者搜尋。
- **可服務時間段:** 照護者可設定每日或每週的特定時間區塊 (例如：每週一 09:00-12:00, 14:00-17:00)。
- **搜尋與篩選:**
  - 使用者可依據地點、服務項目、證照、性別、可服務時間段、評分等條件進行搜尋與篩選。

### 3. 合約與服務管理 (Contract & Service Management)

- **合約週期:** 最短無限制，最長六個月 (180天)。
- **合約續約:** 可以在合約結束前 (例如：提前 30 天內) 由使用者發起續約，照護者同意後生效。續約期仍需遵守最長六個月的限制。
- **服務完成確認:** 服務結束時，雙方（使用者與照護者）需在平台確認服務完成。
- **評分與留言:** 服務完成確認後，使用者必須在期限內（例如：7 天內）對照護者進行星級評分 (1-5 星) 和文字留言。

### 4. 系統功能 (System Features)

- **通知系統:** 合約狀態變更、新服務邀請、證照審核結果、續約提醒等。
- **資料安全:** 密碼加密儲存、資料傳輸使用 HTTPS。

---

## 💾 資料庫設計 (Database Schema Design)
*以下設計基於 Golang 開發中常用的 SQL (例如 PostgreSQL, MySQL) 關係型資料庫。*

### 1. `User` (使用者/家屬 & 照護者 - 統一用戶表)
| 欄位名稱 (Field Name) | 資料型態 (Data Type) | 說明 (Description) | 備註 (Notes) |
|---|---|---|---|
| `user_id` | UUID/INT | 主鍵 (Primary Key) | 平台唯一 ID |
| `email` | VARCHAR(255) | 登入信箱 | Unique |
| `password_hash` | VARCHAR(255) | 密碼雜湊值 | |
| `role` | VARCHAR(50) | 角色 | 'user', 'caregiver', 'admin' |
| `is_active` | BOOLEAN | 帳號是否啟用 | |
| `created_at` | TIMESTAMP | 建立時間 | |
| `updated_at` | TIMESTAMP | 更新時間 | |

### 2. `CaregiverProfile` (照護者檔案)
| 欄位名稱 (Field Name) | 資料型態 (Data Type) | 說明 (Description) | 備註 (Notes) |
|---|---|---|---|
| `profile_id` | UUID/INT | 主鍵 | |
| `user_id` | UUID/INT | 外鍵 (Foreign Key) to `User` | |
| `full_name` | VARCHAR(100) | 姓名 | |
| `gender` | VARCHAR(10) | 性別 | |
| `phone` | VARCHAR(20) | 聯絡電話 | |
| `address` | TEXT | 服務地點/居住地 | 用於搜尋 |
| `bio` | TEXT | 自我介紹/經驗 | |
| `avg_rating` | NUMERIC(2,1) | 平均分數 | 系統計算，範圍 1.0 - 5.0 |
| `service_rate`| DECIMAL | 基本時薪/日薪 | |

### 3. `License` (照護者證照)
| 欄位名稱 (Field Name) | 資料型態 (Data Type) | 說明 (Description) | 備註 (Notes) |
|---|---|---|---|
| `license_id` | UUID/INT | 主鍵 | |
| `caregiver_id` | UUID/INT | 外鍵 to `CaregiverProfile` | |
| `name` | VARCHAR(100) | 證照名稱 | |
| `issue_date` | DATE | 發證日期 | |
| `expiry_date` | DATE | 有效期限 | |
| `status` | VARCHAR(50) | 審核狀態 | 'pending', 'approved', 'rejected' |
| `proof_url` | VARCHAR(255) | 證明文件儲存路徑 | |

### 4. `Availability` (照護者可服務時間段)
| 欄位名稱 (Field Name) | 資料型態 (Data Type) | 說明 (Description) | 備註 (Notes) |
|---|---|---|---|
| `availability_id`| UUID/INT | 主鍵 | |
| `caregiver_id` | UUID/INT | 外鍵 to `CaregiverProfile` | |
| `day_of_week` | INT | 星期幾 | 1=Mon, 7=Sun |
| `start_time` | TIME | 可服務開始時間 | |
| `end_time` | TIME | 可服務結束時間 | |

### 5. `Contract` (服務合約)
| 欄位名稱 (Field Name) | 資料型態 (Data Type) | 說明 (Description) | 備註 (Notes) |
|---|---|---|---|
| `contract_id` | UUID/INT | 主鍵 | |
| `user_id` | UUID/INT | 外鍵 to `User` (發起者/使用者) | |
| `caregiver_id`| UUID/INT | 外鍵 to `User` (照護者) | |
| `start_date` | DATE | 合約開始日期 | |
| `end_date` | DATE | 合約結束日期 | 最長 6 個月 |
| `status` | VARCHAR(50) | 合約狀態 | 'pending', 'active', 'completed', 'canceled' |
| `is_renewal` | BOOLEAN | 是否為續約合約 | |
| `original_contract_id` | UUID/INT | 自參考外鍵 (Self-referencing Foreign Key) | 連結到被續約的合約 ID |

### 6. `Review` (評價與留言)
| 欄位名稱 (Field Name) | 資料型態 (Data Type) | 說明 (Description) | 備註 (Notes) |
|---|---|---|---|
| `review_id` | UUID/INT | 主鍵 | |
| `contract_id` | UUID/INT | 外鍵 to `Contract` | Unique (一個合約只能有一份評價) |
| `user_id` | UUID/INT | 外鍵 to `User` (評價者) | |
| `caregiver_id`| UUID/INT | 外鍵 to `User` (被評價者) | |
| `rating` | INT | 星級評分 | 1 - 5 |
| `comment` | TEXT | 留言內容 | |
| `created_at` | TIMESTAMP | 評價時間 | |

---

## 💻 API 接口設計 (API Endpoints Design)
*基於 Golang 常用的 RESTful API 風格，並假設使用 JSON 格式進行資料交換。*

### 1. 用戶與認證 (Auth & User)
| 動作 (Method) | 路徑 (Path) | 說明 (Description) | 身份驗證 (Auth) |
|---|---|---|---|
| POST | `/api/v1/auth/register` | 註冊新用戶/照護者 | Public |
| POST | `/api/v1/auth/login` | 登入 | Public |
| GET | `/api/v1/users/{id}` | 取得使用者基本資料 | User/Caregiver/Admin |
| PUT | `/api/v1/users/{id}` | 更新使用者基本資料 | User/Caregiver |

### 2. 照護者檔案與搜尋 (Caregiver Profile & Search)
| 動作 (Method) | 路徑 (Path) | 說明 (Description) | 身份驗證 (Auth) |
|---|---|---|---|
| POST | `/api/v1/caregivers/profile` | 建立照護者檔案 | Caregiver |
| PUT | `/api/v1/caregivers/profile` | 更新照護者檔案 | Caregiver |
| GET | `/api/v1/caregivers/search` | 搜尋/篩選照護者 | User |
| GET | `/api/v1/caregivers/{id}` | 取得單一照護者檔案 (含評價/證照) | User |
| POST | `/api/v1/caregivers/licenses` | 上傳新證照 | Caregiver |
| PUT | `/api/v1/caregivers/availability` | 更新可服務時間段 (批次更新) | Caregiver |

> **GET `/api/v1/caregivers/search` 範例查詢參數 (Query Params):**
> - `location`: 地點 (e.g., Taipei)
> - `service_type`: 服務類型 (e.g., Home Care)
> - `start_date`: YYYY-MM-DD
> - `start_time`: HH:MM

### 3. 合約管理 (Contract Management)
| 動作 (Method) | 路徑 (Path) | 說明 (Description) | 身份驗證 (Auth) |
|---|---|---|---|
| POST | `/api/v1/contracts` | 發起新合約邀請 | User |
| GET | `/api/v1/contracts` | 取得我的所有合約 (使用者/照護者) | User/Caregiver |
| GET | `/api/v1/contracts/{id}` | 取得單一合約詳情 | User/Caregiver |
| PUT | `/api/v1/contracts/{id}/accept` | 照護者接受合約 | Caregiver |
| PUT | `/api/v1/contracts/{id}/complete` | 服務完成確認 | User/Caregiver |
| POST | `/api/v1/contracts/{id}/renew` | 發起續約 | User |

### 4. 評價與留言 (Review & Rating)
| 動作 (Method) | 路徑 (Path) | 說明 (Description) | 身份驗證 (Auth) |
|---|---|---|---|
| POST | `/api/v1/reviews` | 提交評價與留言 | User |
| GET | `/api/v1/reviews/caregivers/{id}` | 取得某照護者的所有評價 | Public |

---

## ⚙️ 技術架構建議 (Technology Stack Recommendation)

### Golang 後端架構
| 類別 (Category) | 建議技術 (Recommended Tech) | 說明 (Notes) |
|---|---|---|
| 主語言/框架 | Golang | 高性能、高併發、快速啟動 |
| Web 框架 | Gin 或 Echo | 輕量、高效能的 Web 框架，適合 RESTful API |
| 資料庫驅動 | `lib/pq` (for PostgreSQL) 或 `go-sql-driver/mysql` | 根據選擇的資料庫 |
| ORM/Query Builder | Gorm 或 SQLC | 方便操作資料庫，減少手寫 SQL 錯誤 |
| 認證/JWT | `dgrijalva/jwt-go` | 處理用戶登入狀態 |
| 環境變數 | `spf13/viper` | 管理配置檔和環境變數 |
| 圖片儲存 | AWS S3 或 Google Cloud Storage | 儲存證照等圖片檔案，不建議直接存於伺服器 |
