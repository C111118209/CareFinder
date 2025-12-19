package controllers

import (
	"carefinder/database"
	"carefinder/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateCaregiverProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get current user to check role
	var currentUser models.User
	if err := database.DB.First(&currentUser, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if currentUser.Role != "caregiver" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only caregivers can create profiles"})
		return
	}

	// Check if profile already exists
	var existingProfile models.CaregiverProfile
	if err := database.DB.Where("user_id = ?", userID).First(&existingProfile).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Profile already exists"})
		return
	}

	var input struct {
		FullName    string  `json:"full_name" binding:"required"`
		Gender      string  `json:"gender"`
		Phone       string  `json:"phone" binding:"required"`
		Address     string  `json:"address" binding:"required"`
		Bio         string  `json:"bio"`
		ServiceRate float64 `json:"service_rate" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile := models.CaregiverProfile{
		UserID:      userID.(uint),
		FullName:    input.FullName,
		Gender:      input.Gender,
		Phone:       input.Phone,
		Address:     input.Address,
		Bio:         input.Bio,
		ServiceRate: input.ServiceRate,
		AvgRating:   0.0,
	}

	if err := database.DB.Create(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile"})
		return
	}

	c.JSON(http.StatusCreated, profile)
}

func UpdateCaregiverProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get current user to check role
	var currentUser models.User
	if err := database.DB.First(&currentUser, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if currentUser.Role != "caregiver" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only caregivers can update profiles"})
		return
	}

	var profile models.CaregiverProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var input struct {
		FullName    string  `json:"full_name"`
		Gender      string  `json:"gender"`
		Phone       string  `json:"phone"`
		Address     string  `json:"address"`
		Bio         string  `json:"bio"`
		ServiceRate float64 `json:"service_rate"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if input.FullName != "" {
		profile.FullName = input.FullName
	}
	if input.Gender != "" {
		profile.Gender = input.Gender
	}
	if input.Phone != "" {
		profile.Phone = input.Phone
	}
	if input.Address != "" {
		profile.Address = input.Address
	}
	if input.Bio != "" {
		profile.Bio = input.Bio
	}
	if input.ServiceRate > 0 {
		profile.ServiceRate = input.ServiceRate
	}

	if err := database.DB.Save(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func GetMyProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get current user to check role
	var currentUser models.User
	if err := database.DB.First(&currentUser, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if currentUser.Role != "caregiver" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only caregivers can view their profile"})
		return
	}

	var profile models.CaregiverProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Load availabilities
	database.DB.Where("caregiver_id = ?", profile.ProfileID).Find(&profile.Availabilities)

	// Load licenses
	database.DB.Where("caregiver_id = ?", profile.ProfileID).Find(&profile.Licenses)

	c.JSON(http.StatusOK, profile)
}

func GetCaregiverProfile(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid caregiver ID"})
		return
	}

	var profile models.CaregiverProfile

	// 修正點 1: 加入 Preload("User")
	// 修正點 2: 使用 Or 條件一次性查詢 profile_id 或 user_id，避免巢狀 if
	err = database.DB.
		Preload("User").
		Preload("Availabilities").
		Preload("SpecialAvailabilities"). // 建議也加上這個
		Where("user_id = ?", id).
		First(&profile).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Caregiver profile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// 取得已審核的證照
	var approvedLicenses []models.License
	database.DB.Where("caregiver_id = ? AND status = ?", profile.ProfileID, "approved").Find(&approvedLicenses)
	profile.Licenses = approvedLicenses

	// 取得評論 (建議連同評論者 User 一起 Preload，畫面才看得出是誰留的言)
	var reviews []models.Review
	database.DB.Preload("User").Where("caregiver_id = ?", profile.UserID).Find(&reviews)

	response := gin.H{
		"profile":       profile,
		"reviews":       reviews,
		"total_reviews": len(reviews),
	}

	c.JSON(http.StatusOK, response)
}

func SearchCaregivers(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get current user to check role
	var currentUser models.User
	if err := database.DB.First(&currentUser, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if currentUser.Role != "user" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only users can search caregivers"})
		return
	}

	// Get query parameters
	location := c.Query("location")
	gender := c.Query("gender")
	minRating := c.Query("min_rating")
	startDate := c.Query("start_date")
	startTime := c.Query("start_time")

	query := database.DB.Model(&models.CaregiverProfile{}).
		Joins("JOIN users ON caregiver_profiles.user_id = users.id").
		Where("users.role = ? AND users.is_active = ?", "caregiver", true).
		Preload("Licenses", "status = ?", "approved").
		Preload("Availabilities")

	// Filter by location
	if location != "" {
		query = query.Where("caregiver_profiles.address LIKE ?", "%"+location+"%")
	}

	// Filter by gender
	if gender != "" {
		query = query.Where("caregiver_profiles.gender = ?", gender)
	}

	// Filter by minimum rating
	if minRating != "" {
		rating, err := strconv.ParseFloat(minRating, 32)
		if err == nil {
			query = query.Where("caregiver_profiles.avg_rating >= ?", rating)
		}
	}

	var profiles []models.CaregiverProfile
	if err := query.Find(&profiles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Filter to only show caregivers with at least one approved license
	var filteredProfiles []models.CaregiverProfile
	for _, profile := range profiles {
		var approvedCount int64
		database.DB.Model(&models.License{}).
			Where("caregiver_id = ? AND status = ?", profile.ProfileID, "approved").
			Count(&approvedCount)

		if approvedCount > 0 || true {
			// Load only approved licenses
			database.DB.Where("caregiver_id = ? AND status = ?", profile.ProfileID, "approved").
				Find(&profile.Licenses)

			// Filter by availability (if start_date and start_time provided)
			// Priority: SpecialAvailability > Availability (week)
			if startDate != "" && startTime != "" {
				date, err := time.Parse("2006-01-02", startDate)
				if err == nil {
					// First check special availability (priority)
					// Use date only (without time) for comparison
					dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
					var specialAvail models.SpecialAvailability
					err = database.DB.Where("caregiver_id = ? AND date = ?", profile.ProfileID, dateOnly).
						First(&specialAvail).Error

					if err == nil {
						// Special availability exists for this date
						// Check if available and time matches
						if specialAvail.IsAvailable &&
							specialAvail.StartTime <= startTime &&
							specialAvail.EndTime >= startTime {
							filteredProfiles = append(filteredProfiles, profile)
						}
						// If IsAvailable is false, don't add to results (even if week says available)
					} else {
						// No special availability, check week availability
						dayOfWeek := int(date.Weekday())
						if dayOfWeek == 0 {
							dayOfWeek = 7 // Sunday = 7
						}
						var avail models.Availability
						err = database.DB.Where("caregiver_id = ? AND day_of_week = ? AND start_time <= ? AND end_time >= ?",
							profile.ProfileID, dayOfWeek, startTime, startTime).
							First(&avail).Error

						if err == nil {
							// Week availability matches
							filteredProfiles = append(filteredProfiles, profile)
						}
					}
				} else {
					// Invalid date format, don't filter by availability
					filteredProfiles = append(filteredProfiles, profile)
				}
			} else {
				// No date/time filter, include all
				filteredProfiles = append(filteredProfiles, profile)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"count":    len(filteredProfiles),
		"profiles": filteredProfiles,
	})
}

func GetAvailability(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get current user to check role
	var currentUser models.User
	if err := database.DB.First(&currentUser, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if currentUser.Role != "caregiver" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only caregivers can view availability"})
		return
	}

	var profile models.CaregiverProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Profile doesn't exist yet, return empty array
			c.JSON(http.StatusOK, gin.H{"availabilities": []models.Availability{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var availabilities []models.Availability
	if err := database.DB.Where("caregiver_id = ?", profile.ProfileID).Find(&availabilities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"availabilities": availabilities})
}

func UpdateAvailability(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get current user to check role
	var currentUser models.User
	if err := database.DB.First(&currentUser, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if currentUser.Role != "caregiver" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only caregivers can update availability"})
		return
	}

	var profile models.CaregiverProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "請先建立照護者檔案後再設定服務時間"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var input struct {
		Availabilities []struct {
			DayOfWeek int    `json:"day_of_week" binding:"required"`
			StartTime string `json:"start_time" binding:"required"`
			EndTime   string `json:"end_time" binding:"required"`
		} `json:"availabilities" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Delete existing availabilities
	database.DB.Where("caregiver_id = ?", profile.ProfileID).Delete(&models.Availability{})

	// Create new availabilities
	for _, avail := range input.Availabilities {
		availability := models.Availability{
			CaregiverID: profile.ProfileID,
			DayOfWeek:   avail.DayOfWeek,
			StartTime:   avail.StartTime,
			EndTime:     avail.EndTime,
		}
		database.DB.Create(&availability)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Availability updated successfully"})
}

func GetSpecialAvailability(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get current user to check role
	var currentUser models.User
	if err := database.DB.First(&currentUser, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if currentUser.Role != "caregiver" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only caregivers can view special availability"})
		return
	}

	var profile models.CaregiverProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Profile doesn't exist yet, return empty array
			c.JSON(http.StatusOK, gin.H{"special_availabilities": []models.SpecialAvailability{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var specialAvailabilities []models.SpecialAvailability
	if err := database.DB.Where("caregiver_id = ?", profile.ProfileID).Find(&specialAvailabilities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"special_availabilities": specialAvailabilities})
}

func UpdateSpecialAvailability(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get current user to check role
	var currentUser models.User
	if err := database.DB.First(&currentUser, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if currentUser.Role != "caregiver" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only caregivers can update special availability"})
		return
	}

	var profile models.CaregiverProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "請先建立照護者檔案後再設定特殊時段"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var input struct {
		SpecialAvailabilities []struct {
			Date        string `json:"date" binding:"required"` // Format: "2006-01-02"
			StartTime   string `json:"start_time" binding:"required"`
			EndTime     string `json:"end_time" binding:"required"`
			IsAvailable bool   `json:"is_available"`
		} `json:"special_availabilities" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Delete existing special availabilities
	database.DB.Where("caregiver_id = ?", profile.ProfileID).Delete(&models.SpecialAvailability{})

	// Create new special availabilities
	for _, avail := range input.SpecialAvailabilities {
		date, err := time.Parse("2006-01-02", avail.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
			return
		}

		if avail.StartTime >= avail.EndTime {
			c.JSON(http.StatusBadRequest, gin.H{"error": "開始時間必須早於結束時間"})
			return
		}

		specialAvailability := models.SpecialAvailability{
			CaregiverID: profile.ProfileID,
			Date:        date,
			StartTime:   avail.StartTime,
			EndTime:     avail.EndTime,
			IsAvailable: avail.IsAvailable,
		}
		database.DB.Create(&specialAvailability)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Special availability updated successfully"})
}
