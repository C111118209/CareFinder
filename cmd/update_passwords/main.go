package main

import (
	"carefinder/database"
	"carefinder/models"
	"carefinder/utils"
	"fmt"
	"log"
	"strings"
)

func main() {
	// 連接資料庫
	database.Connect()

	// 取得所有使用者（排除所有 c111118209 開頭的 email）
	var users []models.User
	if err := database.DB.Where("email LIKE ?", "test%").Find(&users).Error; err != nil {
		log.Fatal("Failed to fetch users: ", err)
	}

	fmt.Printf("找到 %d 個使用者需要更新密碼\n", len(users))

	updatedCount := 0
	skippedCount := 0

	// 更新每個使用者的密碼
	for _, user := range users {
		// 檢查 email 是否包含 @
		if !strings.Contains(user.Email, "@") {
			fmt.Printf("跳過使用者 ID %d: email 格式不正確 (%s)\n", user.ID, user.Email)
			skippedCount++
			continue
		}

		// 提取 @ 前面的部分作為新密碼
		parts := strings.Split(user.Email, "@")
		if len(parts) == 0 || parts[0] == "" {
			fmt.Printf("跳過使用者 ID %d: 無法從 email 提取密碼 (%s)\n", user.ID, user.Email)
			skippedCount++
			continue
		}

		newPassword := parts[0]

		// 雜湊新密碼
		hashedPassword, err := utils.HashPassword(newPassword)
		if err != nil {
			fmt.Printf("錯誤：無法雜湊使用者 ID %d 的密碼: %v\n", user.ID, err)
			skippedCount++
			continue
		}

		// 更新密碼
		user.PasswordHash = hashedPassword
		if err := database.DB.Save(&user).Error; err != nil {
			fmt.Printf("錯誤：無法更新使用者 ID %d 的密碼: %v\n", user.ID, err)
			skippedCount++
			continue
		}

		fmt.Printf("✓ 已更新使用者 ID %d (%s) 的密碼為: %s\n", user.ID, user.Email, newPassword)
		updatedCount++
	}

	fmt.Printf("\n完成！\n")
	fmt.Printf("成功更新: %d 個使用者\n", updatedCount)
	fmt.Printf("跳過: %d 個使用者\n", skippedCount)
}
