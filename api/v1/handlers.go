package v1

import (
	"backend/internal/app/health"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handlers struct {
	DocsHandler   gin.HandlerFunc
	HealthHandler *health.Handler
}

func InitHandlers() *Handlers {
	// swagger initialization
	docsHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)

	healthHandler := health.NewHandler()
	return &Handlers{
		DocsHandler:   docsHandler,
		HealthHandler: healthHandler,
	}
}
