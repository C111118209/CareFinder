package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"carefinder/database"
	"carefinder/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

//go:embed all:frontend/dist
var frontendFS embed.FS

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	// Initialize Database
	database.Connect()

	// Setup Router from routes package
	r := routes.SetupRouter()

	// Get the sub-filesystem for the 'dist' directory
	distFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Fatal("failed to create sub-filesystem: ", err)
	}

	// Serve the frontend using a Gin middleware
	r.Use(func(c *gin.Context) {
		path := c.Request.URL.Path

		// If the request is for an API endpoint, pass to the next handler
		if strings.HasPrefix(path, "/api") {
			c.Next()
			return
		}

		// Attempt to serve the static file
		// http.FileServer will automatically handle serving index.html for directories
		fileServer := http.FileServer(http.FS(distFS))

		// Check if the file exists in the embedded filesystem
		// We need to trim the leading slash for the Open call
		if _, err := distFS.Open(strings.TrimPrefix(path, "/")); err != nil {
			// If the file doesn't exist, it's likely a client-side route.
			// Rewrite the request to serve the root index.html,
			// so the Vue app can handle routing.
			c.Request.URL.Path = "/"
		}

		// Let the file server handle the request
		fileServer.ServeHTTP(c.Writer, c.Request)

		// Abort the Gin chain to prevent other handlers from running
		c.Abort()
	})

	// Start Server
	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}