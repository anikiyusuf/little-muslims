package httpx

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)


type TestStruct struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"gte=0,lte=130"`
}

func TestReadJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name 	      string
		body 	      string
		expectedError bool
	}{
		{
			name:      "Valid JSON",
			body: 	   `{"name":"John","age":30}`,
			expectedError: false,
		},
		{
			name:      "Invalid JSON",
			body: 	   `{"name":"John","age":}`,
			expectedError: false,
		},
		{
			name:      "invalid field type",
			body: 	   `{"name":"John","age":"thirty"}`,
			expectedError: true,
		},
		{
			name:      "exceeds max bytes",
			body: 	   `{"name":"` + string(bytes.Repeat([]byte("a"), 2_000_000)) + `","age":30}`,
			expectedError: true,
		},

	}

	for _, tt := range tests{
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.body))
            c.Request.Header.Set("Content-Type", "application/json")
			
			var ts TestStruct
			err := ReadJSON(c, &ts)

			if tt.expectedError{
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
	}