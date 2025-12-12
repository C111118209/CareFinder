package controllers

import (
	"carefinder/database"
	"carefinder/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ========== 用戶管理 ==========

// GetAllUsers 獲取所有用戶（管理員）
func GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Remove password hashes
	for i := range users {
		users[i].PasswordHash = ""
	}

	c.JSON(http.StatusOK, users)
}

// ToggleUserStatus 啟用/停用用戶帳號
func ToggleUserStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// 不能停用管理員帳號
	if user.Role == "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot deactivate admin account"})
		return
	}

	user.IsActive = !user.IsActive
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User status updated successfully",
		"user":    user,
	})
}

// ========== 照護者管理 ==========

// GetAllCaregivers 獲取所有照護者檔案（管理員）
func GetAllCaregivers(c *gin.Context) {
	var profiles []models.CaregiverProfile
	if err := database.DB.Preload("User").Find(&profiles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, profiles)
}

// UpdateCaregiverStatus 更新照護者狀態（管理員）
func UpdateCaregiverStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile ID"})
		return
	}

	var input struct {
		IsActive bool `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var profile models.CaregiverProfile
	if err := database.DB.First(&profile, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Caregiver profile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Update user status
	var user models.User
	if err := database.DB.First(&user, profile.UserID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	user.IsActive = input.IsActive
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update caregiver status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Caregiver status updated successfully",
		"user":    user,
	})
}

// ========== 合約管理 ==========

// GetAllContracts 獲取所有合約（管理員）
func GetAllContracts(c *gin.Context) {
	var contracts []models.Contract
	if err := database.DB.Preload("User").Preload("Caregiver").Find(&contracts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, contracts)
}

// UpdateContractStatus 更新合約狀態（管理員）
func UpdateContractStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	var input struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status
	validStatuses := map[string]bool{
		"pending":   true,
		"active":    true,
		"completed": true,
		"canceled":  true,
	}
	if !validStatuses[input.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	var contract models.Contract
	if err := database.DB.First(&contract, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	contract.Status = input.Status
	if err := database.DB.Save(&contract).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update contract status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Contract status updated successfully",
		"contract": contract,
	})
}

// ========== 證照審核 ==========

// GetPendingLicenses 獲取待審核證照列表（管理員）
func GetPendingLicenses(c *gin.Context) {
	var licenses []models.License
	if err := database.DB.Where("status = ?", "pending").
		Order("created_at ASC").
		Find(&licenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Load caregiver profiles for each license
	type LicenseWithProfile struct {
		models.License
		CaregiverProfile *models.CaregiverProfile `json:"caregiver_profile,omitempty"`
	}

	result := make([]LicenseWithProfile, len(licenses))
	for i, license := range licenses {
		result[i].License = license
		var profile models.CaregiverProfile
		if err := database.DB.Preload("User").First(&profile, license.CaregiverID).Error; err == nil {
			result[i].CaregiverProfile = &profile
		}
	}

	c.JSON(http.StatusOK, result)
}

// GetAllLicenses 獲取所有證照（管理員）
func GetAllLicenses(c *gin.Context) {
	var licenses []models.License
	if err := database.DB.Order("created_at DESC").Find(&licenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Load caregiver profiles
	type LicenseWithProfile struct {
		models.License
		CaregiverProfile *models.CaregiverProfile `json:"caregiver_profile,omitempty"`
	}

	result := make([]LicenseWithProfile, len(licenses))
	for i, license := range licenses {
		result[i].License = license
		var profile models.CaregiverProfile
		if err := database.DB.Preload("User").First(&profile, license.CaregiverID).Error; err == nil {
			result[i].CaregiverProfile = &profile
		}
	}

	c.JSON(http.StatusOK, result)
}

// ReviewLicense 審核證照（通過/拒絕）
func ReviewLicense(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid license ID"})
		return
	}

	var input struct {
		Status string `json:"status" binding:"required"` // "approved" or "rejected"
		Note   string `json:"note"`                      // Optional note
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status
	if input.Status != "approved" && input.Status != "rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status must be 'approved' or 'rejected'"})
		return
	}

	var license models.License
	if err := database.DB.First(&license, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "License not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if license.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "License is not pending review"})
		return
	}

	license.Status = input.Status
	if err := database.DB.Save(&license).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update license status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "License reviewed successfully",
		"license": license,
	})
}

// ========== 統計報告 ==========

// GetStatistics 獲取系統營運數據與報告
func GetStatistics(c *gin.Context) {
	var stats struct {
		Users struct {
			Total      int64 `json:"total"`
			Active     int64 `json:"active"`
			Inactive   int64 `json:"inactive"`
			ByRole     map[string]int64 `json:"by_role"`
		} `json:"users"`
		Caregivers struct {
			Total  int64 `json:"total"`
			Active int64 `json:"active"`
		} `json:"caregivers"`
		Contracts struct {
			Total     int64 `json:"total"`
			Pending   int64 `json:"pending"`
			Active    int64 `json:"active"`
			Completed int64 `json:"completed"`
			Canceled  int64 `json:"canceled"`
		} `json:"contracts"`
		Licenses struct {
			Total    int64 `json:"total"`
			Pending  int64 `json:"pending"`
			Approved int64 `json:"approved"`
			Rejected int64 `json:"rejected"`
		} `json:"licenses"`
		Reviews struct {
			Total      int64   `json:"total"`
			AvgRating  float64 `json:"avg_rating"`
		} `json:"reviews"`
	}

	// User statistics
	database.DB.Model(&models.User{}).Count(&stats.Users.Total)
	database.DB.Model(&models.User{}).Where("is_active = ?", true).Count(&stats.Users.Active)
	database.DB.Model(&models.User{}).Where("is_active = ?", false).Count(&stats.Users.Inactive)
	
	stats.Users.ByRole = make(map[string]int64)
	var userCount, caregiverCount, adminCount int64
	database.DB.Model(&models.User{}).Where("role = ?", "user").Count(&userCount)
	database.DB.Model(&models.User{}).Where("role = ?", "caregiver").Count(&caregiverCount)
	database.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount)
	stats.Users.ByRole["user"] = userCount
	stats.Users.ByRole["caregiver"] = caregiverCount
	stats.Users.ByRole["admin"] = adminCount

	// Caregiver statistics
	database.DB.Model(&models.CaregiverProfile{}).Count(&stats.Caregivers.Total)
	database.DB.Model(&models.CaregiverProfile{}).
		Joins("JOIN users ON caregiver_profiles.user_id = users.id").
		Where("users.is_active = ?", true).
		Count(&stats.Caregivers.Active)

	// Contract statistics
	database.DB.Model(&models.Contract{}).Count(&stats.Contracts.Total)
	database.DB.Model(&models.Contract{}).Where("status = ?", "pending").Count(&stats.Contracts.Pending)
	database.DB.Model(&models.Contract{}).Where("status = ?", "active").Count(&stats.Contracts.Active)
	database.DB.Model(&models.Contract{}).Where("status = ?", "completed").Count(&stats.Contracts.Completed)
	database.DB.Model(&models.Contract{}).Where("status = ?", "canceled").Count(&stats.Contracts.Canceled)

	// License statistics
	database.DB.Model(&models.License{}).Count(&stats.Licenses.Total)
	database.DB.Model(&models.License{}).Where("status = ?", "pending").Count(&stats.Licenses.Pending)
	database.DB.Model(&models.License{}).Where("status = ?", "approved").Count(&stats.Licenses.Approved)
	database.DB.Model(&models.License{}).Where("status = ?", "rejected").Count(&stats.Licenses.Rejected)

	// Review statistics
	database.DB.Model(&models.Review{}).Count(&stats.Reviews.Total)
	var avgRatingResult struct {
		AvgRating float64
	}
	database.DB.Model(&models.Review{}).
		Select("AVG(rating) as avg_rating").
		Scan(&avgRatingResult)
	stats.Reviews.AvgRating = avgRatingResult.AvgRating

	c.JSON(http.StatusOK, stats)
}

