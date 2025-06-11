package middleware

import (
	"Backend/constant"
	"Backend/dto"
	"Backend/handler"
	mocks "Backend/mocks/usecase"
	"Backend/usecase"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/mock"
)

func TestMiddlewareError(t *testing.T) {
	testCases := []struct {
		name           string
		handle         func(c *gin.Context)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "no error",
			handle: func(c *gin.Context) {
			},
			expectedStatus: http.StatusOK,
			expectedBody:   ``,
		},
		{
			name: "validation errors - empty",
			handle: func(c *gin.Context) {
				c.Error(validator.ValidationErrors{})
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"success":false,"error":[],"data":null}`,
		},
		{
			name: "validation errors - non-empty 1",
			handle: func(c *gin.Context) {
				type request struct {
					Field string `json:"field" binding:"required"`
				}

				var r request
				errorArg := c.ShouldBindQuery(&r)

				if errorArg != nil {
					c.Error(errorArg)
				}
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: `{"success":false,` +
				`"error":[{"field":"Field",` +
				`"message":"` +
				`Key: 'request.Field' Error:Field validation for 'Field' failed on the 'required' tag` +
				`"}],` +
				`"data":null}`,
		},
		{
			name: "validation errors - non-empty 2",
			handle: func(c *gin.Context) {
				type request struct {
					Field string `json:"field" binding:"required"`
					Var   string `json:"variable" binding:"required"`
				}

				var r request
				errorArg := c.ShouldBindQuery(&r)

				if errorArg != nil {
					c.Error(errorArg)
				}
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: `{"success":false,` +
				`"error":[{"field":"Field",` +
				`"message":"` +
				`Key: 'request.Field' Error:Field validation for 'Field' failed on the 'required' tag` +
				`"},{"field":"Var",` +
				`"message":"` +
				`Key: 'request.Var' Error:Field validation for 'Var' failed on the 'required' tag` +
				`"}],` +
				`"data":null}`,
		},
		{
			name: "custom error 1",
			handle: func(c *gin.Context) {
				c.Error(constant.CustomError{
					StatusCode: http.StatusBadRequest,
					Message:    "custom error message",
				})
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: `{"success":false,` +
				`"error":"custom error message","data":null}`,
		},
		{
			name: "custom error 2",
			handle: func(c *gin.Context) {
				c.Error(constant.CustomError{
					StatusCode: http.StatusBadGateway,
					Message:    "custom error message",
				})
			},
			expectedStatus: http.StatusBadGateway,
			expectedBody: `{"success":false,` +
				`"error":"custom error message","data":null}`,
		},
		{
			name: "interval server error",
			handle: func(c *gin.Context) {
				c.Error(errors.New("unknown error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: `{"success":false,` +
				`"error":"unknown error","data":null}`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			//given
			recorder := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(recorder)

			middleware := NewMiddleware()

			engine.GET("/", middleware.Error(), tt.handle)
			r := httptest.NewRequest("", "/", nil)

			//when
			engine.ServeHTTP(recorder, r)

			//then
			assert.Equal(t, tt.expectedStatus, recorder.Code)
			assert.Equal(t, tt.expectedBody, recorder.Body.String())
		})
	}
}
func TestIntegratedHandlerGetSymbols(t *testing.T) {
	// This is the unit test, except the middleware is also used
	// tests for output and status code
	// This can be used to partly test that middleware works in general,
	// regardless of handler
	testCases := []struct {
		name           string
		link           string
		ucSetup        func(*gin.Context) usecase.UsecaseItf
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "no query parameter provided",
			link: "/symbols",
			ucSetup: func(ctx *gin.Context) usecase.UsecaseItf {
				return new(mocks.UsecaseItf)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: `{"success":false,` +
				`"error":"please provide keywords",` +
				`"data":null}`},
		{
			name: "usecase returns error",
			link: "/symbols?keywords=BA",
			ucSetup: func(ctx *gin.Context) usecase.UsecaseItf {
				mocked := new(mocks.UsecaseItf)

				// input to usecase
				var req dto.GetSymbolsReq
				req.Prefix = "BA"

				// usecase mechanism
				contextMatcher := mock.MatchedBy(func(c *gin.Context) bool {
					// Verify the query parameter was properly extracted
					prefix, exists := c.GetQuery("keywords")
					return exists && prefix == "BA"
				})
				mocked.On("GetSymbols", contextMatcher, &req).Return(nil, constant.ErrAPIExceed)

				return mocked
			},
			expectedStatus: http.StatusBadGateway,
			expectedBody: `{"success":false,` +
				`"error":"exceeded API-use limit today",` +
				`"data":null}`,
		},
		{
			name: "handling successful usecase outcome",
			link: "/symbols?keywords=BA",
			ucSetup: func(ctx *gin.Context) usecase.UsecaseItf {
				mocked := new(mocks.UsecaseItf)

				// input to usecase
				var req dto.GetSymbolsReq
				req.Prefix = "BA"

				// output from usecase
				var symbols dto.AlphaSymbolsRes
				bestMatches := []dto.AlphaSymbolRes{
					{
						Symbol: "BA",
						Name:   "Boeing Company",
						Region: "United States",
					},
					{
						Symbol: "BA.LON",
						Name:   "BAE Systems plc",
						Region: "United Kingdom",
					},
				}
				symbols.BestMatches = bestMatches

				// usecase mechanism
				contextMatcher := mock.MatchedBy(func(c *gin.Context) bool {
					// Verify the query parameter was properly extracted
					prefix, exists := c.GetQuery("keywords")
					return exists && prefix == "BA"
				})
				mocked.On("GetSymbols", contextMatcher, &req).Return(&symbols, nil)

				return mocked
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{"data":{"best_matches":[` +
				`{"symbol":"BA","name":"Boeing Company","region":"United States"},` +
				`{"symbol":"BA.LON","name":"BAE Systems plc","region":"United Kingdom"}` +
				`]},"error":null,"message":null}`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			//given
			recorder := httptest.NewRecorder()
			c, engine := gin.CreateTestContext(recorder)

			middleware := NewMiddleware()
			hd := handler.NewHandler(tt.ucSetup(c))

			engine.GET("/symbols", middleware.Error(), hd.GetSymbols)

			r := httptest.NewRequest("GET", tt.link, nil)

			//when
			engine.ServeHTTP(recorder, r)

			//then
			assert.Equal(t, tt.expectedStatus, recorder.Code)
			assert.Equal(t, tt.expectedBody, recorder.Body.String())
		})
	}
}
