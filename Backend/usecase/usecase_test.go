package usecase

import (
	"Backend/constant"
	"Backend/dto"
	mocks1 "Backend/mocks/repo"
	mocks2 "Backend/mocks/util"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
)

func TestUnitUsecaseParseOHLCV(t *testing.T) {
	testCases := []struct {
		name           string
		tsInput        func() *map[string]string
		expectedOutput func(*dto.DailyOHLCVRes)
		expectedError  error
	}{
		{
			name: "no open price",
			tsInput: func() *map[string]string {
				ts := make(map[string]string)
				return &ts
			},
			expectedOutput: func(output *dto.DailyOHLCVRes) {
				assert.Equal(t, nil, output)
			},
			expectedError: constant.ErrAlphaParseBody(
				"can't find open price as usual",
			),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			//given
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			uc := NewUsecase(new(mocks1.RepoItf), new(mocks2.HttpClientItf))

			//when
			output, err := uc.ParseOHLCV(c, tt.tsInput())

			//then
			tt.expectedOutput(output)
			assert.Equal(t, errors.Is(tt.expectedError, err), true)
		})
	}
}
