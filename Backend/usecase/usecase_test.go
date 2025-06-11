package usecase

import (
	"Backend/dto"
	"Backend/repo"
	"Backend/util"
	"errors"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
)

func TestUnitUsecaseParseOHLCV(t *testing.T) {
	testCases := []struct {
		name           string
		rpSetup        func(*gin.Context) repo.RepoItf
		hcSetup        func(*gin.Context) util.HttpClientItf
		tsInput        func() *map[string]string
		expectedOutput dto.DailyOHLCVRes
		expectedError  error
	}{}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			//given
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			uc := NewUsecase(tt.rpSetup(c), tt.hcSetup(c))

			//when
			output, err := uc.ParseOHLCV(c, tt.tsInput())

			//then
			assert.Equal(t, reflect.DeepEqual(tt.expectedOutput, *output), true)
			assert.Equal(t, errors.Is(tt.expectedError, err), true)
		})
	}
}
