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

func CreateContract(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, gin.H{"error": "Only users can create contracts"})
		return
	}

	var input struct {
		CaregiverID uint   `json:"caregiver_id" binding:"required"`
		StartDate   string `json:"start_date" binding:"required"`
		EndDate     string `json:"end_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format (YYYY-MM-DD)"})
		return
	}

	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format (YYYY-MM-DD)"})
		return
	}

	// Validate contract duration (max 6 months / 180 days)
	duration := endDate.Sub(startDate)
	if duration > 180*24*time.Hour {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Contract duration cannot exceed 180 days"})
		return
	}

	if duration <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End date must be after start date"})
		return
	}

	// Verify caregiver exists and is active
	var caregiver models.User
	if err := database.DB.First(&caregiver, input.CaregiverID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Caregiver not found"})
		return
	}

	if caregiver.Role != "caregiver" || !caregiver.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid caregiver"})
		return
	}

	contract := models.Contract{
		UserID:      userID.(uint),
		CaregiverID: input.CaregiverID,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      "pending",
		IsRenewal:   false,
	}

	if err := database.DB.Create(&contract).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create contract"})
		return
	}

	c.JSON(http.StatusCreated, contract)
}

func GetContracts(c *gin.Context) {
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

	var contracts []models.Contract
	switch currentUser.Role {
	case "user":
		database.DB.Where("user_id = ?", userID).Preload("Caregiver").Preload("User").Find(&contracts)
	case "caregiver":
		database.DB.Where("caregiver_id = ?", userID).Preload("User").Preload("Caregiver").Find(&contracts)
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid role"})
		return
	}

	c.JSON(http.StatusOK, contracts)
}

func GetContract(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	var contract models.Contract
	if err := database.DB.Preload("User").Preload("Caregiver").First(&contract, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check permission: user can only view own contracts
	var currentUser models.User
	database.DB.First(&currentUser, userID)
	if currentUser.Role != "admin" && contract.UserID != userID.(uint) && contract.CaregiverID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, contract)
}

func AcceptContract(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, gin.H{"error": "Only caregivers can accept contracts"})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
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

	// Verify contract belongs to this caregiver
	if contract.CaregiverID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Verify contract status is pending
	if contract.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Contract is not in pending status"})
		return
	}

	contract.Status = "active"
	if err := database.DB.Save(&contract).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update contract"})
		return
	}

	c.JSON(http.StatusOK, contract)
}

func CompleteContract(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
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

	// Verify user is part of this contract
	if contract.UserID != userID.(uint) && contract.CaregiverID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// For simplicity, we'll mark as completed when either party confirms
	// In a real system, you might want to track separate confirmation flags
	if contract.Status != "active" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Contract is not active"})
		return
	}

	contract.Status = "completed"
	if err := database.DB.Save(&contract).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update contract"})
		return
	}

	c.JSON(http.StatusOK, contract)
}

func RenewContract(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, gin.H{"error": "Only users can renew contracts"})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	var originalContract models.Contract
	if err := database.DB.First(&originalContract, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Verify contract belongs to this user
	if originalContract.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Verify contract can be renewed (within 30 days of end date)
	daysUntilEnd := int(time.Until(originalContract.EndDate).Hours() / 24)
	if daysUntilEnd > 30 || daysUntilEnd < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Contract can only be renewed within 30 days before end date"})
		return
	}

	var input struct {
		EndDate string `json:"end_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newEndDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format (YYYY-MM-DD)"})
		return
	}

	// Validate renewal duration (max 6 months from original start or new start)
	startDate := originalContract.EndDate.AddDate(0, 0, 1) // Start from day after original ends
	duration := newEndDate.Sub(startDate)
	if duration > 180*24*time.Hour {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Renewal duration cannot exceed 180 days"})
		return
	}

	if newEndDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End date must be after start date"})
		return
	}

	// Create renewal contract
	renewalContract := models.Contract{
		UserID:             userID.(uint),
		CaregiverID:        originalContract.CaregiverID,
		StartDate:          startDate,
		EndDate:            newEndDate,
		Status:             "pending",
		IsRenewal:          true,
		OriginalContractID: &originalContract.ContractID,
	}

	if err := database.DB.Create(&renewalContract).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create renewal contract"})
		return
	}

	c.JSON(http.StatusCreated, renewalContract)
}
