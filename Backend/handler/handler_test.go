package handler

import (
	"Backend/usecase"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
)

func UnitTestHandlerGetSymbols(t *testing.T) {
	testCases := []struct {
		name          string
		ucSetup       func(*gin.Context) usecase.UsecaseItf
		expectedError func(*gin.Context)
		expectedBody  string
	}{}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			//given
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			r := httptest.NewRequest("GET", "/symbols", nil)
			c.Request = r

			hd := NewHandler(tt.ucSetup(c))

			//when
			hd.GetSymbols(c)

			//then
			assert.Equal(t, tt.expectedBody, w.Body.String())
			tt.expectedError(c)
		})
	}
}
