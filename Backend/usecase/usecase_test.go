package usecase

import (
	"Backend/constant"
	"Backend/dto"
	mocks1 "Backend/mocks/repo"
	mocks2 "Backend/mocks/util"
	"Backend/repo"
	"Backend/util"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	"github.com/shopspring/decimal"
)

func TestUnitUsecaseParseOHLCV(t *testing.T) {
	testCases := []struct {
		name           string
		tsInput        func() *map[string]string
		expectedOutput func(*dto.DailyOHLCVRes)
		expectedError  func(error)
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
			expectedError: func(err error) {
				expected := constant.ErrAlphaParseBody(
					"can't find open price as usual",
				)
				assert.Equal(t, errors.Is(expected, err), true)
			},
		},
		{
			name: "unparseable open price",
			tsInput: func() *map[string]string {
				ts := make(map[string]string)
				ts["1. open"] = "one hundred"
				return &ts
			},
			expectedOutput: func(output *dto.DailyOHLCVRes) {
				assert.Equal(t, nil, output)
			},
			expectedError: func(err error) {
				var ce constant.CustomError
				assert.Equal(t, errors.As(err, &ce), true)
				assert.Equal(t, ce.StatusCode, http.StatusBadGateway)
				assert.Equal(
					t,
					strings.HasPrefix(
						ce.Message,
						"Alpha Vantage API response-body-parse error: ",
					),
					true,
				)
			},
		},
		{
			name: "no high price",
			tsInput: func() *map[string]string {
				ts := make(map[string]string)
				ts["1. open"] = "100"
				return &ts
			},
			expectedOutput: func(output *dto.DailyOHLCVRes) {
				assert.Equal(t, nil, output)
			},
			expectedError: func(err error) {
				expected := constant.ErrAlphaParseBody(
					"can't find high price as usual",
				)
				assert.Equal(t, errors.Is(expected, err), true)
			},
		},
		{
			name: "unparseable high price",
			tsInput: func() *map[string]string {
				ts := make(map[string]string)
				ts["1. open"] = "100"
				ts["2. high"] = "one hundred"
				return &ts
			},
			expectedOutput: func(output *dto.DailyOHLCVRes) {
				assert.Equal(t, nil, output)
			},
			expectedError: func(err error) {
				var ce constant.CustomError
				assert.Equal(t, errors.As(err, &ce), true)
				assert.Equal(t, ce.StatusCode, http.StatusBadGateway)
				assert.Equal(
					t,
					strings.HasPrefix(
						ce.Message,
						"Alpha Vantage API response-body-parse error: ",
					),
					true,
				)
			},
		},
		{
			name: "no low price",
			tsInput: func() *map[string]string {
				ts := make(map[string]string)
				ts["1. open"] = "100"
				ts["2. high"] = "100"
				return &ts
			},
			expectedOutput: func(output *dto.DailyOHLCVRes) {
				assert.Equal(t, nil, output)
			},
			expectedError: func(err error) {
				expected := constant.ErrAlphaParseBody(
					"can't find low price as usual",
				)
				assert.Equal(t, errors.Is(expected, err), true)
			},
		},
		{
			name: "unparseable low price",
			tsInput: func() *map[string]string {
				ts := make(map[string]string)
				ts["1. open"] = "100"
				ts["2. high"] = "100"
				ts["3. low"] = "one hundred"
				return &ts
			},
			expectedOutput: func(output *dto.DailyOHLCVRes) {
				assert.Equal(t, nil, output)
			},
			expectedError: func(err error) {
				var ce constant.CustomError
				assert.Equal(t, errors.As(err, &ce), true)
				assert.Equal(t, ce.StatusCode, http.StatusBadGateway)
				assert.Equal(
					t,
					strings.HasPrefix(
						ce.Message,
						"Alpha Vantage API response-body-parse error: ",
					),
					true,
				)
			},
		},
		{
			name: "no close price",
			tsInput: func() *map[string]string {
				ts := make(map[string]string)
				ts["1. open"] = "100"
				ts["2. high"] = "100"
				ts["3. low"] = "100"
				return &ts
			},
			expectedOutput: func(output *dto.DailyOHLCVRes) {
				assert.Equal(t, nil, output)
			},
			expectedError: func(err error) {
				expected := constant.ErrAlphaParseBody(
					"can't find close price as usual",
				)
				assert.Equal(t, errors.Is(expected, err), true)
			},
		},
		{
			name: "unparseable close price",
			tsInput: func() *map[string]string {
				ts := make(map[string]string)
				ts["1. open"] = "100"
				ts["2. high"] = "100"
				ts["3. low"] = "100"
				ts["4. close"] = "one hundred"
				return &ts
			},
			expectedOutput: func(output *dto.DailyOHLCVRes) {
				assert.Equal(t, nil, output)
			},
			expectedError: func(err error) {
				var ce constant.CustomError
				assert.Equal(t, errors.As(err, &ce), true)
				assert.Equal(t, ce.StatusCode, http.StatusBadGateway)
				assert.Equal(
					t,
					strings.HasPrefix(
						ce.Message,
						"Alpha Vantage API response-body-parse error: ",
					),
					true,
				)
			},
		},
		{
			name: "no volume",
			tsInput: func() *map[string]string {
				ts := make(map[string]string)
				ts["1. open"] = "100"
				ts["2. high"] = "100"
				ts["3. low"] = "100"
				ts["4. close"] = "100"
				return &ts
			},
			expectedOutput: func(output *dto.DailyOHLCVRes) {
				assert.Equal(t, nil, output)
			},
			expectedError: func(err error) {
				expected := constant.ErrAlphaParseBody(
					"can't find volume as usual",
				)
				assert.Equal(t, errors.Is(expected, err), true)
			},
		},
		{
			name: "unparseable volume",
			tsInput: func() *map[string]string {
				ts := make(map[string]string)
				ts["1. open"] = "100"
				ts["2. high"] = "100"
				ts["3. low"] = "100"
				ts["4. close"] = "100"
				ts["5. volume"] = "one hundred"
				return &ts
			},
			expectedOutput: func(output *dto.DailyOHLCVRes) {
				assert.Equal(t, nil, output)
			},
			expectedError: func(err error) {
				var ce constant.CustomError
				assert.Equal(t, errors.As(err, &ce), true)
				assert.Equal(t, ce.StatusCode, http.StatusBadGateway)
				assert.Equal(
					t,
					strings.HasPrefix(
						ce.Message,
						"Alpha Vantage API response-body-parse error: ",
					),
					true,
				)
			},
		},
		{
			name: "completely parseable",
			tsInput: func() *map[string]string {
				ts := make(map[string]string)
				ts["1. open"] = "100"
				ts["2. high"] = "100"
				ts["3. low"] = "100"
				ts["4. close"] = "100"
				ts["5. volume"] = "100"
				return &ts
			},
			expectedOutput: func(output *dto.DailyOHLCVRes) {
				var ohlcv dto.DailyOHLCVRes
				ohlcv.OHLC = make(map[string]decimal.Decimal)
				ohlcv.OHLC["open"] = decimal.NewFromInt(100)
				ohlcv.OHLC["high"] = decimal.NewFromInt(100)
				ohlcv.OHLC["low"] = decimal.NewFromInt(100)
				ohlcv.OHLC["close"] = decimal.NewFromInt(100)
				ohlcv.Volume = 100

				assert.Equal(t, reflect.DeepEqual(&ohlcv, output), true)
			},
			expectedError: func(err error) {
				assert.Equal(t, nil, err)
			},
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
			tt.expectedError(err)
		})
	}
}
func TestUnitUsecaseBuildStockData(t *testing.T) {
	timeDate := time.Date(2025, 6, 2, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		name           string
		dataInput      func() *dto.DataPerSymbol
		expectedOutput func() *dto.StockDataRes
	}{
		{
			name: "one week",
			dataInput: func() *dto.DataPerSymbol {
				data := new(dto.DataPerSymbol)

				dateGen := util.DateGenerator(timeDate.AddDate(0, 0, 2))
				ohlcvGen := util.NewOHLCVGenerator(
					&dateGen, 100, 100)

				data.TimeSeries = append(data.TimeSeries, ohlcvGen.Next())

				return data
			},
			expectedOutput: func() *dto.StockDataRes {
				output := new(dto.StockDataRes)

				week := new(dto.WeekRes)
				week.Monday = dto.DateOnly(timeDate)
				week.Friday = dto.DateOnly(timeDate).AddDate(0, 0, 4)

				dateGen := util.DateGenerator(timeDate.AddDate(0, 0, 2))
				ohlcvGen := util.NewOHLCVGenerator(
					&dateGen, 100, 100)

				week.DailyData = append(week.DailyData, ohlcvGen.Next())

				output.Weeks = append(output.Weeks, week)
				return output
			},
		},
		{
			name: "two weeks",
			dataInput: func() *dto.DataPerSymbol {
				data := new(dto.DataPerSymbol)

				dateGen := util.DateGenerator(timeDate.AddDate(0, 0, 2))
				ohlcvGen := util.NewOHLCVGenerator(
					&dateGen, 100, 100)

				data.TimeSeries = append(data.TimeSeries, ohlcvGen.Next())

				dateGen = util.DateGenerator(timeDate.AddDate(0, 0, 8))
				ohlcvGen.DateGen = &dateGen

				data.TimeSeries = append(data.TimeSeries, ohlcvGen.Next())
				data.TimeSeries = append(data.TimeSeries, ohlcvGen.Next())
				return data
			},
			expectedOutput: func() *dto.StockDataRes {
				output := new(dto.StockDataRes)

				week := new(dto.WeekRes)
				week.Monday = dto.DateOnly(timeDate)
				week.Friday = dto.DateOnly(timeDate).AddDate(0, 0, 4)

				dateGen := util.DateGenerator(timeDate.AddDate(0, 0, 2))
				ohlcvGen := util.NewOHLCVGenerator(
					&dateGen, 100, 100)

				week.DailyData = append(week.DailyData, ohlcvGen.Next())

				output.Weeks = append(output.Weeks, week)

				week = new(dto.WeekRes)
				week.Monday = dto.DateOnly(timeDate).AddDate(0, 0, 7)
				week.Friday = dto.DateOnly(timeDate).AddDate(0, 0, 11)

				dateGen = util.DateGenerator(timeDate.AddDate(0, 0, 8))
				ohlcvGen.DateGen = &dateGen

				week.DailyData = append(week.DailyData, ohlcvGen.Next())
				week.DailyData = append(week.DailyData, ohlcvGen.Next())

				output.Weeks = append(output.Weeks, week)
				return output
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			//given
			uc := NewUsecase(new(mocks1.RepoItf), new(mocks2.HttpClientItf))

			//when
			output := uc.BuildStockData(tt.dataInput())

			//then
			assert.Equal(t, reflect.DeepEqual(tt.expectedOutput(), output), true)
		})
	}
}
func TestUnitUsecaseCollectSymbol(t *testing.T) {
	testCases := []struct {
		name           string
		repoSetup      func() repo.RepoItf
		httpSetup      func() util.HttpClientItf
		inputReq       *dto.CollectSymbolReq
		expectedOutput func() *dto.StockDataRes
		expectedErr    func(error)
	}{}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			//given
			uc := NewUsecase(tt.repoSetup(), tt.httpSetup())
			c, _ := gin.CreateTestContext(httptest.NewRecorder())

			//when
			output, err := uc.CollectSymbol(c, tt.inputReq)

			//then
			assert.Equal(t, reflect.DeepEqual(tt.expectedOutput(), output), true)
			tt.expectedErr(err)
		})
	}
}
