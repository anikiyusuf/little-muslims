package httpx

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestBadRequestRes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
		router.GET("/badrequest", func(c *gin.Context) {
			BadRequestResponse(c, errors.New("bad request error"))
})

w := httptest.NewRecorder()
req, _ := http.NewRequest(http.MethodGet, "/badrequest", nil)
router.ServeHTTP(w, req)

assert.Equal(t, http.StatusBadRequest, w.Code)
assert.Contains(t, w.Body.String(), "invalid request")
}

func TestNotFoundRes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
		router.GET("/notfound", func(c *gin.Context) {
			NotFoundResponse(c, errors.New("not found error"))
})

w := httptest.NewRecorder()
req, _ := http.NewRequest(http.MethodGet, "/notfound", nil)
router.ServeHTTP(w, req)

assert.Equal(t, http.StatusNotFound, w.Code)
assert.Contains(t, w.Body.String(), "not found error")
}

func TestConflictResponse(t *testing.T){
	gin.SetMode(gin.TestMode)
	router := gin.New()
		router.GET("/conflict", func(c *gin.Context) {
			ConflictResponse(c, errors.New("conflict error"))
})

w := httptest.NewRecorder()
req, _ := http.NewRequest(http.MethodGet, "/conflict", nil)
router.ServeHTTP(w, req)

assert.Equal(t, http.StatusConflict, w.Code)
assert.Contains(t, w.Body.String(), "conflict error")
}

func TestInternalServerError(t *testing.T){
	gin.SetMode(gin.TestMode)
	router := gin.New()
		router.GET("/internalservererror", func(c *gin.Context) {
			InternalServerError(c, errors.New("internal server error"))
})
w := httptest.NewRecorder()
req, _ := http.NewRequest(http.MethodGet, "/internalservererror", nil)
router.ServeHTTP(w, req)

assert.Equal(t, http.StatusUnauthorized, w.Code)
assert.Contains(t, w.Body.String(), "internal server error")
}

func TestForbiddenResponse(t *testing.T){
	gin.SetMode(gin.TestMode)
	router := gin.New()
		router.GET("/forbidden", func(c *gin.Context) {
			ForbiddenResponse(c, errors.New("forbidden error"))
})
w := httptest.NewRecorder()
req, _ := http.NewRequest(http.MethodGet, "/forbidden", nil)
router.ServeHTTP(w, req)

assert.Equal(t, http.StatusForbidden, w.Code)
assert.Contains(t, w.Body.String(), "forbidden error")
}

func TestUnauthorizedResponse(t *testing.T){
	gin.SetMode(gin.TestMode)
	router := gin.New()
		router.GET("/unauthorized", func(c *gin.Context) {
			UnauthorizedResponse(c, errors.New("unauthorized error"))
})
w := httptest.NewRecorder()
req, _ := http.NewRequest(http.MethodGet, "/unauthorized", nil)
router.ServeHTTP(w, req)

assert.Equal(t, http.StatusUnauthorized, w.Code)
assert.Contains(t, w.Body.String(), "unauthorized error")
}

func TestRateLimitExceededResponse(t *testing.T){
	gin.SetMode(gin.TestMode)
	router := gin.New()
		router.GET("/ratelimit", func(c *gin.Context) {
			RateLimitExceededResponse(c, "10")
})
w := httptest.NewRecorder()
req, _ := http.NewRequest(http.MethodGet, "/ratelimit", nil)
router.ServeHTTP(w, req)

assert.Equal(t, http.StatusTooManyRequests, w.Code)
assert.Contains(t, w.Body.String(), "rate limit exceeded, retry after10")
assert.Equal(t, "10", w.Header().Get("Retry-After"))

}
