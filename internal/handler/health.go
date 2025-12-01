package handler


import (
	"github.com/gin-gonic/gin"
	"github.com/yusufaniki/muslim_tech/internal/httpx"
	"github.com/yusufaniki/muslim_tech/internal/service"
	"github.com/yusufaniki/muslim_tech/internal/boostrap"
)

type HealthHandler struct {
	HealthService *service.HealthService
	app           *boostrap.Application
}

func NewHealthHandler(app *boostrap.Application, healthService *service.HealthService) *HealthHandler{
	return &HealthHandler{HealthService: healthService, app: app}
}

type HealthResponse struct {
	Data      string `json:"data"`
	Message    string `json:"message"`  
}

func (h *HealthHandler) Health(c *gin.Context){
	stats := h.HealthService.Health()

	if stats["status"] == "down"{
		httpx.ErrorResponse(c, 500, stats["error"], nil)
		return
	}

	httpx.OkResponse(c, "service is healthy", stats)
}


func (h *HealthHandler) Ping(c *gin.Context) {
	httpx.OkResponse(c, "pong", map[string]string{"status":"alive"})
}