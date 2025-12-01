package middleware

import (
	"net/http"
	"strings"


    "github.com/yusufaniki/muslim_tech/pkg/auth"
	"github.com/yusufaniki/muslim_tech/pkg/logger"
	"github.com/gin-gonic/gin"
	
)


var log = logger.CreateZapLogger()

func AuthMiddleware(JWTManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		c.Abort()
		return
}

tokenString := strings.TrimPrefix(authHeader, "Bearer ")

claims, err := JWTManager.ValidateToken(tokenString)
if err != nil {
	log.Error("Error validating token", err)
	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	c.Abort()
	return
}

c.Set("user", claims)
c.Next()
	}
}

