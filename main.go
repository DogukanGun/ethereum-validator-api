package main

import (
	_ "ethereum-validator-api/docs" // This is important - imports the swagger docs
	"ethereum-validator-api/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"os"
)

// @title           Ethereum Validator API
// @version         1.0
// @description     API that provides Ethereum validator information including sync committee duties and block rewards.

// @contact.name   API Support
// @contact.url    https://github.com/yourusername/ethereum-validator-api
// @contact.email  your-email@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:3001
// @BasePath  /

func main() {
	utils.InitializeENV(".env")
	router := gin.Default()

	// Set up CORS with proper configuration
	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		corsOrigin = "https://sf.dogukangun.de"
	}
	localCorsOrigin := "http://localhost:3003"
	
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{corsOrigin,localCorsOrigin},
		AllowMethods:     []string{"GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

	// Swagger documentation routes
	// Redirect /docs to /swagger/index.html for better UX
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	
	// Use the standard Swagger handler
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup the API endpoints
	err := utils.SetupEndpoints(router)
	if err != nil {
		log.Fatalf("Failed to setup endpoints: %v", err)
	}
	
	// Start the server
	log.Println("Server starting at http://localhost:3004")
	log.Println("Swagger UI available at http://localhost:3004/swagger/index.html")
	
	if err := router.Run(":3004"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
