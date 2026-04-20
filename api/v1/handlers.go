package v1

import (
	"backend/internal/app/health"
	"backend/internal/app/lca"
	"backend/internal/app/traverse"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handlers struct {
	DocsHandler     gin.HandlerFunc
	HealthHandler   *health.Handler
	TraverseHandler *traverse.Handler
	LCAHandler      *lca.Handler
}

func InitHandlers() *Handlers {
	docsHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)
	return &Handlers{
		DocsHandler:     docsHandler,
		HealthHandler:   health.NewHandler(),
		TraverseHandler: traverse.NewHandler(),
		LCAHandler:      lca.NewHandler(),
	}
}
