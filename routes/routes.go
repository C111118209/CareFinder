package routes

import (
	"carefinder/controllers"
	"carefinder/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")

	// Serve HTML files
	r.StaticFile("/", "./static/index.html")
	r.StaticFile("/index.html", "./static/index.html")
	r.StaticFile("/login.html", "./static/login.html")
	r.StaticFile("/register.html", "./static/register.html")
	r.StaticFile("/dashboard.html", "./static/dashboard.html")
	r.StaticFile("/caregiver-search.html", "./static/caregiver-search.html")
	r.StaticFile("/caregiver-profile.html", "./static/caregiver-profile.html")
	r.StaticFile("/caregiver-profile-setup.html", "./static/caregiver-profile-setup.html")
	r.StaticFile("/caregiver-availability.html", "./static/caregiver-availability.html")
	r.StaticFile("/caregiver-licenses.html", "./static/caregiver-licenses.html")
	r.StaticFile("/contracts.html", "./static/contracts.html")
	r.StaticFile("/contract-details.html", "./static/contract-details.html")
	r.StaticFile("/create-contract.html", "./static/create-contract.html")
	r.StaticFile("/profile.html", "./static/profile.html")
	r.StaticFile("/admin-dashboard.html", "./static/admin-dashboard.html")

	// Note: Uploads are now served through protected routes only
	// r.Static("/uploads", "./uploads") - Removed for security

	// Group API routes
	api := r.Group("/api/v1")
	{
		// Public routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(middleware.JwtAuthMiddleware())
		{
			// User management
			users := protected.Group("/users")
			{
				users.GET("/me", controllers.GetCurrentUser)
				users.GET("/:id", controllers.GetUser)
				users.PUT("/:id", controllers.UpdateUser)
				users.PUT("/:id/password", controllers.UpdatePassword)
			}

			// Caregiver profile management
			caregivers := protected.Group("/caregivers")
			{
				caregivers.POST("/profile", controllers.CreateCaregiverProfile)
				caregivers.PUT("/profile", controllers.UpdateCaregiverProfile)
				caregivers.GET("/profile", controllers.GetMyProfile)
				caregivers.GET("/availability", controllers.GetAvailability)
				caregivers.PUT("/availability", controllers.UpdateAvailability)
				caregivers.GET("/special-availability", controllers.GetSpecialAvailability)
				caregivers.PUT("/special-availability", controllers.UpdateSpecialAvailability)
				caregivers.GET("/search", controllers.SearchCaregivers)
				caregivers.POST("/licenses", controllers.UploadLicense)
				caregivers.GET("/licenses", controllers.GetMyLicenses)
				caregivers.DELETE("/licenses/:id", controllers.DeleteLicense)
				// Protected route for license images
				caregivers.GET("/licenses/image/:filename", controllers.ServeLicenseImage)
			}

			// Contract management
			contracts := protected.Group("/contracts")
			{
				contracts.POST("", controllers.CreateContract)
				contracts.GET("", controllers.GetContracts)
				contracts.GET("/:id", controllers.GetContract)
				contracts.PUT("/:id/accept", controllers.AcceptContract)
				contracts.PUT("/:id/complete", controllers.CompleteContract)
				contracts.POST("/:id/renew", controllers.RenewContract)
			}

			// Review management
			reviews := protected.Group("/reviews")
			{
				reviews.POST("", controllers.CreateReview)
			}

			// Admin routes (require admin role)
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireAdmin())
			{
				// User management
				admin.GET("/users", controllers.GetAllUsers)
				admin.PUT("/users/:id/status", controllers.ToggleUserStatus)

				// Caregiver management
				admin.GET("/caregivers", controllers.GetAllCaregivers)
				admin.PUT("/caregivers/:id/status", controllers.UpdateCaregiverStatus)

				// Contract management
				admin.GET("/contracts", controllers.GetAllContracts)
				admin.PUT("/contracts/:id/status", controllers.UpdateContractStatus)

				// License review
				admin.GET("/licenses", controllers.GetAllLicenses)
				admin.GET("/licenses/pending", controllers.GetPendingLicenses)
				admin.PUT("/licenses/:id/review", controllers.ReviewLicense)

				// Statistics
				admin.GET("/statistics", controllers.GetStatistics)
			}
		}

		// Public caregiver profile and reviews (no auth required)
		api.GET("/caregivers/:id", controllers.GetCaregiverProfile)
		api.GET("/reviews/caregivers/:id", controllers.GetCaregiverReviews)
	}

	return r
}
