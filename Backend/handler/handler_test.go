package handler

import (
	"Backend/constant"
	"Backend/dto"
	"Backend/middleware"
	mocks "Backend/mocks/usecase"
	"Backend/usecase"
	"Backend/util"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	"github.com/stretchr/testify/mock"
)

func TestUnitHandlerGetSymbols(t *testing.T) {
	testCases := []struct {
		name           string
		link           string
		ucSetup        func(*gin.Context) usecase.UsecaseItf
		expectedStatus int
		expectedBody   string
		expectedError  func(*gin.Context)
	}{
		{
			name: "no query parameter provided",
			link: "/symbols",
			ucSetup: func(ctx *gin.Context) usecase.UsecaseItf {
				return new(mocks.UsecaseItf)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "",
			expectedError: func(ctx *gin.Context) {
				assert.Equal(t, len(ctx.Errors), 1)

				var ce constant.CustomError
				assert.Equal(t, errors.As(ctx.Errors[0], &ce), true)
				assert.Equal(t, errors.Is(ce, constant.ErrNoKeywords), true)
			},
		},
		{
			name: "usecase returns error",
			link: "/symbols?keywords=BA",
			ucSetup: func(ctx *gin.Context) usecase.UsecaseItf {
				mock := new(mocks.UsecaseItf)

				// input to usecase
				var req dto.GetSymbolsReq
				req.Prefix = "BA"

				// usecase mechanism
				mock.On("GetSymbols", ctx, &req).Return(nil, constant.ErrAPIExceed)

				return mock
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "",
			expectedError: func(ctx *gin.Context) {
				assert.Equal(t, len(ctx.Errors), 1)

				var ce constant.CustomError
				assert.Equal(t, errors.As(ctx.Errors[0], &ce), true)
				assert.Equal(t, errors.Is(ce, constant.ErrAPIExceed), true)
			},
		},
		{
			name: "handling successful usecase outcome",
			link: "/symbols?keywords=BA",
			ucSetup: func(ctx *gin.Context) usecase.UsecaseItf {
				mock := new(mocks.UsecaseItf)

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
				mock.On("GetSymbols", ctx, &req).Return(&symbols, nil)

				return mock
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{"data":{"best_matches":[` +
				`{"symbol":"BA","name":"Boeing Company","region":"United States"},` +
				`{"symbol":"BA.LON","name":"BAE Systems plc","region":"United Kingdom"}` +
				`]},"error":null,"message":null}`,
			expectedError: func(ctx *gin.Context) {
				assert.Equal(t, len(ctx.Errors), 0)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			//given
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			r := httptest.NewRequest("GET", tt.link, nil)
			c.Request = r

			hd := NewHandler(tt.ucSetup(c))

			//when
			hd.GetSymbols(c)

			//then
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
			tt.expectedError(c)
		})
	}
}
func TestIntegratedHandlerGetSymbols(t *testing.T) {
	// This is the unit test, except the middleware is also used
	// tests for output and status code
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

			middleware := middleware.NewMiddleware()
			hd := NewHandler(tt.ucSetup(c))

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
func TestUnitHandlerCollectSymbol(t *testing.T) {
	testCases := []struct {
		name           string
		link           string
		ucSetup        func(*gin.Context) usecase.UsecaseItf
		expectedStatus int
		expectedBody   string
		expectedError  func(*gin.Context)
	}{
		{
			name: "no path parameter provided",
			link: "/data",
			ucSetup: func(ctx *gin.Context) usecase.UsecaseItf {
				return new(mocks.UsecaseItf)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "",
			expectedError: func(ctx *gin.Context) {
				assert.Equal(t, len(ctx.Errors), 1)

				var ce constant.CustomError
				assert.Equal(t, errors.As(ctx.Errors[0], &ce), true)
				assert.Equal(t, errors.Is(ce, constant.ErrNoSymbol), true)
			},
		},
		{
			name: "usecase returns error",
			link: "/data/IBM",
			ucSetup: func(ctx *gin.Context) usecase.UsecaseItf {
				mock := new(mocks.UsecaseItf)

				// input to usecase
				var req dto.CollectSymbolReq
				req.Symbol = "IBM"

				// usecase mechanism
				mock.On("CollectSymbol", ctx, &req).Return(nil, constant.ErrAPIExceed)

				return mock
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "",
			expectedError: func(ctx *gin.Context) {
				assert.Equal(t, len(ctx.Errors), 1)

				var ce constant.CustomError
				assert.Equal(t, errors.As(ctx.Errors[0], &ce), true)
				log.Println(ce)
				assert.Equal(t, errors.Is(ce, constant.ErrAPIExceed), true)
			},
		},
		{
			name: "handling successful usecase outcome",
			link: "/data/AAPL",
			ucSetup: func(ctx *gin.Context) usecase.UsecaseItf {
				mock := new(mocks.UsecaseItf)

				// input to usecase
				var req dto.CollectSymbolReq
				req.Symbol = "AAPL"

				// output from usecase
				// - whole stock data structure
				var stockData dto.StockDataRes

				// - meta data (and time setup code)
				var meta dto.SymbolDataMeta
				meta.Symbol = "AAPL"
				meta.Size = 3

				dateGen := util.NewDateGenerator("2025-06-01")
				meta.LastRefreshed = dto.DateOnly(dateGen.Current())

				stockData.MetaData = &meta

				// - weeks
				ohlcvGen := util.NewOHLCVGenerator(dateGen, 100, 1)
				stockData.Weeks = []*dto.WeekRes{
					{
						Monday: dateGen.Next(),
						Friday: dateGen.Next(),
						DailyData: []dto.DailyOHLCVRes{
							ohlcvGen.Next(),
						},
					},
					{
						Monday: dateGen.Next(),
						Friday: dateGen.Next(),
						DailyData: []dto.DailyOHLCVRes{
							ohlcvGen.Next(),
							ohlcvGen.Next(),
						},
					},
				}

				// usecase mechanism
				mock.On("CollectSymbol", ctx, &req).Return(&stockData, nil)

				return mock
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   "",
			expectedError: func(ctx *gin.Context) {
				assert.Equal(t, len(ctx.Errors), 0)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			//given
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			r := httptest.NewRequest("POST", tt.link, nil)
			c.Request = r

			hd := NewHandler(tt.ucSetup(c))

			// - have to manually add params
			if len(tt.link) > len("/data")+1 {
				extraPath := tt.link[len("/data"):]

				var b byte = '/'
				assert.Equal(t, b, extraPath[0])

				params := gin.Params{
					{Key: "symbol", Value: extraPath[1:]},
				}
				c.Params = params
			}

			//when
			hd.CollectSymbol(c)

			//then
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
			tt.expectedError(c)
		})
	}
}
