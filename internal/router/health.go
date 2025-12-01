package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yusufaniki/muslim_tech/internal/boostrap"
	"github.com/yusufaniki/muslim_tech/internal/service"
	"github.com/yusufaniki/muslim_tech/internal/handler"
)


func RegisterHealthRoutes(api  *gin.RouterGroup, app *boostrap.Application) {
	healthService  := service.NewHealthService(app.ConnPool)
	healthHandler  := handler.NewHealthHandler(app, healthService)
	
	api.GET("/health", healthHandler.Health)
	api.GET("/ping", healthHandler.Ping)
}