package router

import (
	
	"github.com/yusufaniki/muslim_tech/internal/boostrap"
	"github.com/gin-gonic/gin"
	// "github.com/gin-contrib/cors"
)

const apiVersion = "/api/v1"

// SetupRoutes initializes and returns the Gin router
func SetupRoutes(app *boostrap.Application) *gin.Engine {
	r := gin.Default()


	// CORS middleware
	// setupMiddleware(r)


	public := r.Group(apiVersion)
	{
		RegisterHealthRoutes(public, app)
		RegisterPublicAuthRoutes(public, app)

	}

	return r
}
