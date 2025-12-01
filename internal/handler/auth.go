package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yusufaniki/muslim_tech/internal/cache"
	"github.com/yusufaniki/muslim_tech/internal/config"
	"github.com/yusufaniki/muslim_tech/internal/httpx"
	"github.com/yusufaniki/muslim_tech/internal/service"
	"github.com/yusufaniki/muslim_tech/internal/types"
)

type AuthHandler struct {
	authService  *service.AuthService
	cache        *cache.RedisCache
}


func NewAuthHandler(authService *service.AuthService, redisCache  *cache.RedisCache, cfg config.Config) *AuthHandler {
    return &AuthHandler{
		authService: authService,
		cache:  redisCache,
	}
}

func (a *AuthHandler) RegisterUser(c *gin.Context) {
	var input types.RegisterUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		httpx.BadRequestResponse(c, err)
		return
	}

	user, err := a.authService.Register(c.Request.Context(), input)
	if err != nil {
		if err == service.ErrEmailAlreadyExists {
			httpx.ConflictResponse(c, err)
		} else {
			httpx.InternalServerError(c, err)
		}
		return 
	}

	c.JSON(http.StatusCreated, user)
}
