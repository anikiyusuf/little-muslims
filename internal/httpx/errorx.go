package httpx

import (
	"fmt"
	"net/http"
	"github.com/yusufaniki/muslim_tech/pkg/logger"
    "github.com/gin-gonic/gin"
)

var Log = logger.CreateZapLogger()

func BadRequestResponse(c *gin.Context, message interface{}){
	fmt.Println("message as interface", message)
	ErrorResponse(c, http.StatusBadRequest, "invalid request", message)
}

func NotFoundResponse(c *gin.Context, err error){
	c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
}


func ConflictResponse(c *gin.Context, err error){
	c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
}

func InternalServerError(c *gin.Context, err error){
	c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error})
}

func UnauthorizedResponse(c *gin.Context, err error){
	c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
}

func ForbiddenResponse(c *gin.Context, err error){
	c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
}

func RateLimitExceededResponse(c *gin.Context, retryAfter string){
	Log.Warnw("rate limit exceeded", "method", c.Request.Method, "path", c.Request.URL.Path)
     c.Header("Retry-After", retryAfter)
	 c.JSON(http.StatusTooManyRequests, gin.H{"error":"rate limit exceeded, retry after" + retryAfter})
}