package controllers

import (
	"carefinder/database"
	"carefinder/models"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UploadLicense(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, gin.H{"error": "Only caregivers can upload licenses"})
		return
	}

	// Get caregiver profile
	var profile models.CaregiverProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Caregiver profile not found. Please create your profile first."})
		return
	}

	// Parse form data (support both JSON and multipart form)
	var input struct {
		Name      string `form:"name" json:"name" binding:"required"`
		IssueDate string `form:"issue_date" json:"issue_date" binding:"required"`
		ExpiryDate string `form:"expiry_date" json:"expiry_date" binding:"required"`
	}

	// Try to bind as multipart form first, then JSON
	if err := c.ShouldBind(&input); err != nil {
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// Parse dates
	issueDate, err := time.Parse("2006-01-02", input.IssueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue_date format (YYYY-MM-DD)"})
		return
	}

	expiryDate, err := time.Parse("2006-01-02", input.ExpiryDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expiry_date format (YYYY-MM-DD)"})
		return
	}

	// Handle file upload
	file, err := c.FormFile("proof")
	var proofURL string

	if err == nil && file != nil {
		// For now, save to local storage
		// In production, upload to S3 or similar
		uploadDir := "./uploads/licenses"
		
		// Create directory if not exists
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}
		
		filename := filepath.Join(uploadDir, strconv.FormatUint(uint64(userID.(uint)), 10)+"_"+strconv.FormatInt(time.Now().Unix(), 10)+filepath.Ext(file.Filename))
		
		// Note: In production, use proper file storage service
		proofURL = "/uploads/licenses/" + filepath.Base(filename)
		
		// Save file
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file: " + err.Error()})
			return
		}
	} else {
		// Allow creating license without file for now (can be added later)
		proofURL = ""
	}

	// Create license
	license := models.License{
		CaregiverID: profile.ProfileID,
		Name:        input.Name,
		IssueDate:   issueDate,
		ExpiryDate:  expiryDate,
		Status:      "pending",
		ProofURL:    proofURL,
	}

	if err := database.DB.Create(&license).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create license"})
		return
	}

	c.JSON(http.StatusCreated, license)
}

func GetMyLicenses(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, gin.H{"error": "Only caregivers can view licenses"})
		return
	}

	// Get caregiver profile
	var profile models.CaregiverProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Caregiver profile not found"})
		return
	}

	// Get licenses
	var licenses []models.License
	if err := database.DB.Where("caregiver_id = ?", profile.ProfileID).Order("created_at DESC").Find(&licenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, licenses)
}

func DeleteLicense(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, gin.H{"error": "Only caregivers can delete licenses"})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid license ID"})
		return
	}

	// Get caregiver profile
	var profile models.CaregiverProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Caregiver profile not found"})
		return
	}

	// Verify license belongs to this caregiver
	var license models.License
	if err := database.DB.Where("license_id = ? AND caregiver_id = ?", id, profile.ProfileID).First(&license).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "License not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Delete license
	if err := database.DB.Delete(&license).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete license"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "License deleted successfully"})
}

