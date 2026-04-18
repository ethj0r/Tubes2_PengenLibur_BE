package app

import (
	v1 "backend/api/v1"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	corsConfig := cors.DefaultConfig()
	origin, ok := os.LookupEnv("CORS_ALLOWED_ORIGIN")
	if !ok {
		origin = "http://localhost:3000"
	}

	corsConfig.AllowOrigins = []string{origin}
	corsConfig.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}

	handlers := v1.InitHandlers()

	v1.RegisterRoutes(r, handlers)

	r.Run(":" + port)
}
