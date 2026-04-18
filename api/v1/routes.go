package v1

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, handlers *Handlers) {
	r.GET("/swagger/*any", handlers.DocsHandler)
	r.GET("/health", handlers.HealthHandler.HealthCheck)
}
