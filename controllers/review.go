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

func CreateReview(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, gin.H{"error": "Only users can create reviews"})
		return
	}

	var input struct {
		ContractID uint   `json:"contract_id" binding:"required"`
		Rating     int    `json:"rating" binding:"required"`
		Comment    string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate rating range
	if input.Rating < 1 || input.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rating must be between 1 and 5"})
		return
	}

	// Get contract
	var contract models.Contract
	if err := database.DB.First(&contract, input.ContractID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Verify contract belongs to this user
	if contract.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Verify contract is completed
	if contract.Status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Contract must be completed before review"})
		return
	}

	// Check if review already exists
	var existingReview models.Review
	if err := database.DB.Where("contract_id = ?", input.ContractID).First(&existingReview).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Review already exists for this contract"})
		return
	}

	// Verify review is within 7 days of contract completion
	// For simplicity, we'll check if contract end_date is within 7 days
	daysSinceEnd := int(time.Since(contract.EndDate).Hours() / 24)
	if daysSinceEnd > 7 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Review must be submitted within 7 days of contract completion"})
		return
	}

	// Create review
	review := models.Review{
		ContractID:  input.ContractID,
		UserID:      userID.(uint),
		CaregiverID: contract.CaregiverID,
		Rating:      input.Rating,
		Comment:     input.Comment,
		CreatedAt:   time.Now(),
	}

	if err := database.DB.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	// Update caregiver's average rating
	updateCaregiverRating(contract.CaregiverID)

	c.JSON(http.StatusCreated, review)
}

func GetCaregiverReviews(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid caregiver ID"})
		return
	}

	// Verify caregiver exists
	var caregiver models.User
	if err := database.DB.First(&caregiver, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Caregiver not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if caregiver.Role != "caregiver" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is not a caregiver"})
		return
	}

	var reviews []models.Review
	if err := database.DB.Where("caregiver_id = ?", id).
		Preload("User").
		Order("created_at DESC").
		Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"caregiver_id": id,
		"count":        len(reviews),
		"reviews":      reviews,
	})
}

// Helper function to update caregiver's average rating
func updateCaregiverRating(caregiverID uint) {
	var reviews []models.Review
	database.DB.Where("caregiver_id = ?", caregiverID).Find(&reviews)

	if len(reviews) == 0 {
		return
	}

	var totalRating int
	for _, review := range reviews {
		totalRating += review.Rating
	}

	avgRating := float32(totalRating) / float32(len(reviews))

	// Get caregiver profile
	var profile models.CaregiverProfile
	if err := database.DB.Where("user_id = ?", caregiverID).First(&profile).Error; err == nil {
		profile.AvgRating = avgRating
		database.DB.Save(&profile)
	}
}

