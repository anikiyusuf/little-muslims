package httpx


import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
	"github.com/yusufaniki/muslim_tech/internal/boostrap" 
)


var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

type Response struct {
	*boostrap.Application

}

var ErrEmptyBody = errors.New("request body is empty")

func ReadJSON(c *gin.Context, data any) error {
	maxBytes := 1_048_578
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, int64(maxBytes))

	dec := json.NewDecoder(c.Request.Body)
	dec.DisallowUnknownFields()

	if c.Request.ContentLength == 0 {
		return ErrEmptyBody
	}

	if err := dec.Decode(data); err != nil {
		return err
	}
	return nil 
}


func ErrorResponse(c *gin.Context, status int, msg string, err interface{}) {
	Log.Errorw("Error occurred", "message", msg,  "status", status, "error", err)
	c.JSON(status, gin.H{"message": msg, "error":err})

}

func OkResponse(c *gin.Context, msg string, data any){
	Log.Infof("message: %s, response: %v", msg)
	c.JSON(http.StatusOK, gin.H{"message": msg, "data": data})
}

func ReadQueryAndValidate(c *gin.Context) (string, error) {
	 var token string 

	 c.Query("token")
	 return token, nil
}