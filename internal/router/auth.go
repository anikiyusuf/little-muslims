package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yusufaniki/muslim_tech/internal/boostrap"
	"github.com/yusufaniki/muslim_tech/internal/service"
	"github.com/yusufaniki/muslim_tech/internal/handler"
)



const (
	authPrefix = "/auth"
)


func RegisterPublicAuthRoutes(api *gin.RouterGroup, app *boostrap.Application) {
	authService := service.NewAuthService(app.ConnPool, *app.Repository, *app.Queue, *app.Cache, *app.JWTManager)
	authHandler := handler.NewAuthHandler(authService, app.Cache, *app.Config)

	api.POST(authPrefix+"/register", authHandler.RegisterUser)
    api.POST(authPrefix+"/verify-token", authHandler.VerifyUser)
	api.POST(authPrefix+"/login", authHandler.Login)
	api.POST(authPrefix+"/forgot-password", authHandler.ForgetPassword)
	api.POST(authPrefix+"/reset-password", authHandler.ResetPassword)
	// api.POST(authPrefix+"/refresh-token", authHandler.RefreshToken)
	api.POST(authPrefix+"/resend-token", authHandler.ResendVerificationCode)

}