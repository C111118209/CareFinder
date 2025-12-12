# CareFinder 待辦事項 (TODO)

本文檔列出根據 `spec.md` 規格，目前架構中尚未實踐的功能。

## ✅ 已實現功能

- [x] 資料庫模型定義（User, CaregiverProfile, License, Availability, Contract, Review）
- [x] 資料庫連接與自動遷移
- [x] 用戶註冊功能 (`POST /api/v1/auth/register`)
- [x] 用戶登入功能 (`POST /api/v1/auth/login`)
- [x] 密碼加密（bcrypt）
- [x] JWT Token 生成
- [x] JWT 認證中間件（驗證 token 有效性）
- [x] Token 驗證工具函數（解析 token 獲取 user_id）
- [x] `GET /api/v1/users/{id}` - 取得使用者基本資料
- [x] `PUT /api/v1/users/{id}` - 更新使用者基本資料
- [x] `POST /api/v1/caregivers/profile` - 建立照護者檔案
- [x] `PUT /api/v1/caregivers/profile` - 更新照護者檔案
- [x] `GET /api/v1/caregivers/{id}` - 取得單一照護者檔案（包含評價、證照、可服務時間段）
- [x] `GET /api/v1/caregivers/search` - 搜尋/篩選照護者（支援地點、性別、評分、可服務時間段篩選）
- [x] `PUT /api/v1/caregivers/availability` - 更新可服務時間段（批次更新）
- [x] `POST /api/v1/contracts` - 發起新合約邀請（含合約週期驗證）
- [x] `GET /api/v1/contracts` - 取得我的所有合約
- [x] `GET /api/v1/contracts/{id}` - 取得單一合約詳情
- [x] `PUT /api/v1/contracts/{id}/accept` - 照護者接受合約
- [x] `PUT /api/v1/contracts/{id}/complete` - 服務完成確認
- [x] `POST /api/v1/contracts/{id}/renew` - 發起續約（含續約邏輯驗證）
- [x] `POST /api/v1/reviews` - 提交評價與留言（含評分期限驗證、自動更新平均評分）
- [x] `GET /api/v1/reviews/caregivers/{id}` - 取得某照護者的所有評價
- [x] 合約週期驗證（最長 6 個月/180 天）
- [x] 續約邏輯（結束前 30 天內、最長 6 個月限制）
- [x] 評分期限驗證（服務完成後 7 天內）
- [x] 證照審核狀態影響搜尋（只顯示有已審核證照的照護者）
- [x] 平均評分自動計算
- [x] 照護者檔案管理前端頁面（建立/編輯）
- [x] 服務時間設定前端頁面
- [x] 證照上傳與管理前端頁面
- [x] 合約發起功能改進（表單式介面）
- [x] 角色區分的導航系統（照護者與使用者不同介面）
- [x] 登入API返回用戶資訊（包含角色）

---

## 🔨 待實現功能

### 1. 認證與授權中間件

- [ ] 角色權限檢查中間件（User/Caregiver/Admin）- 目前在各控制器中手動檢查

### 2. 用戶管理 API

- [x] ~~`GET /api/v1/users/{id}` - 取得使用者基本資料~~ ✅ 已完成
- [x] ~~`PUT /api/v1/users/{id}` - 更新使用者基本資料~~ ✅ 已完成

### 3. 照護者檔案管理

- [x] ~~`POST /api/v1/caregivers/profile` - 建立照護者檔案~~ ✅ 已完成
- [x] ~~`PUT /api/v1/caregivers/profile` - 更新照護者檔案~~ ✅ 已完成
- [x] ~~`GET /api/v1/caregivers/{id}` - 取得單一照護者檔案~~ ✅ 已完成
- [x] ~~`GET /api/v1/caregivers/search` - 搜尋/篩選照護者~~ ✅ 已完成
- [ ] 搜尋功能中支援 `service_type` 參數：
      1. 日常生活協助
      （包含吃飯、洗澡、如廁、移位、整理儀容等基本 ADL）
      2. 安全看顧與陪伴
      （防走失、陪伴聊天、安撫躁動、確保環境安全）
      3. 日常活動與認知刺激
      （簡易運動、認知活動、回憶治療、維持社交互動）
      4. 醫療與健康管理
      （量血壓、提醒吃藥、健康觀察、簡易護理）
      5. 外出與就醫陪同
      （陪同看診、外出散步、協助外出交通）
      6. 家務支援

### 4. 證照管理

- [x] ~~`POST /api/v1/caregivers/licenses` - 上傳新證照~~ ✅ 已完成
  - 需要 Caregiver 角色
  - 檔案上傳功能（目前使用本地儲存，待整合 AWS S3 或 Google Cloud Storage）
  - 證照狀態預設為 'pending'
- [x] ~~`GET /api/v1/caregivers/licenses` - 取得我的證照~~ ✅ 已完成
- [x] ~~`DELETE /api/v1/caregivers/licenses/:id` - 刪除證照~~ ✅ 已完成
- [x] ~~證照上傳前端頁面~~ ✅ 已完成
- [ ] 證照審核機制
  - 管理員審核介面（通過/拒絕）
  - 證照審核結果通知

### 5. 可服務時間段管理

- [x] ~~`PUT /api/v1/caregivers/availability` - 更新可服務時間段（批次更新）~~ ✅ 已完成

### 6. 合約管理

- [x] ~~`POST /api/v1/contracts` - 發起新合約邀請~~ ✅ 已完成
- [x] ~~`GET /api/v1/contracts` - 取得我的所有合約~~ ✅ 已完成
- [x] ~~`GET /api/v1/contracts/{id}` - 取得單一合約詳情~~ ✅ 已完成
- [x] ~~`PUT /api/v1/contracts/{id}/accept` - 照護者接受合約~~ ✅ 已完成
- [x] ~~`PUT /api/v1/contracts/{id}/complete` - 服務完成確認~~ ✅ 已完成
- [x] ~~`POST /api/v1/contracts/{id}/renew` - 發起續約~~ ✅ 已完成
- [ ] 驗證照護者可服務時間段（合約時間需在照護者可服務時間內）

### 7. 評價與留言系統

- [x] ~~`POST /api/v1/reviews` - 提交評價與留言~~ ✅ 已完成
- [x] ~~`GET /api/v1/reviews/caregivers/{id}` - 取得某照護者的所有評價~~ ✅ 已完成

### 8. 通知系統

- [ ] 合約狀態變更通知
  - 合約被接受時通知使用者
  - 合約完成時通知雙方
- [ ] 新服務邀請通知
  - 使用者發起合約時通知照護者
- [ ] 證照審核結果通知
  - 管理員審核證照後通知照護者
- [ ] 續約提醒通知
  - 合約結束前 30 天提醒使用者可續約

### 9. 管理員功能

- [ ] 管理所有用戶
  - 查看、啟用/停用用戶帳號
- [ ] 管理照護者
  - 查看所有照護者檔案
  - 管理照護者狀態
- [ ] 管理服務合約
  - 查看所有合約
  - 管理合約狀態
- [ ] 證照審核
  - 查看待審核證照列表
  - 審核證照（通過/拒絕）
- [ ] 查看系統營運數據與報告
  - 用戶統計
  - 合約統計
  - 評價統計
- [ ] 管理網站內容
  - 公告管理
  - 使用條款管理

### 10. 業務邏輯驗證

- [x] ~~合約週期驗證（最長 6 個月/180 天）~~ ✅ 已完成
- [x] ~~續約邏輯（結束前 30 天內、最長 6 個月限制）~~ ✅ 已完成
- [ ] 服務完成確認邏輯（雙方都需確認才能完成）- 目前簡化為單方確認即可
- [x] ~~評分期限驗證（服務完成後 7 天內）~~ ✅ 已完成
- [x] ~~證照審核狀態影響搜尋（只顯示有已審核證照的照護者）~~ ✅ 已完成
- [x] ~~平均評分自動計算~~ ✅ 已完成

### 11. 技術實作

- [x] ~~圖片檔案上傳功能（基本實現）~~ ✅ 已完成
  - 目前使用本地儲存（`./uploads/licenses/`）
  - 支援 multipart/form-data 上傳
  - [ ] 整合 AWS S3 或 Google Cloud Storage（待實現）
  - [x] ~~證照檔案儲存~~ ✅ 已完成
- [ ] 環境變數管理
  - 使用 `spf13/viper` 管理配置檔
  - JWT Secret、資料庫連線等配置
- [ ] 錯誤處理與日誌
  - 統一的錯誤回應格式
  - 日誌記錄機制
- [ ] API 文檔
  - Swagger/OpenAPI 文檔生成
- [ ] 單元測試與整合測試
  - 控制器測試
  - 資料庫操作測試
- [ ] HTTPS 配置
  - 資料傳輸使用 HTTPS（生產環境）

---

## 📝 備註

- 目前使用 SQLite 作為資料庫，規格建議使用 PostgreSQL 或 MySQL
- 需要考慮資料庫遷移策略（從 SQLite 遷移到 PostgreSQL/MySQL）
- 所有 API 端點都需要適當的錯誤處理和驗證
- 建議實作 API 版本控制機制
- 考慮實作 API 速率限制（Rate Limiting）

