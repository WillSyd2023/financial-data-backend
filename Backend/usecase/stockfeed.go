package usecase

import (
	"Backend/constant"
	"Backend/dto"
	"Backend/repo"
	"Backend/util"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type UsecaseItf interface {
	// Helper methods
	GetUnexpectedInfo([]byte) error
	PrevWeekend(time.Time) time.Time
	NextMonday(time.Time) time.Time
	NextFriday(time.Time) time.Time
	NextWeek(time.Time) *dto.WeekRes
	ParseOHLCV(*gin.Context, *map[string]string) (*dto.DailyOHLCVRes, error)

	// Main methods
	GetSymbols(*gin.Context, *dto.GetSymbolsReq) (*dto.AlphaSymbolsRes, error)
	CollectSymbol(*gin.Context, *dto.CollectSymbolReq) (*dto.StockDataRes, error)
	DeleteSymbol(*gin.Context, *dto.DeleteSymbolReq) error
	StoredData(*gin.Context) ([]*dto.StockDataRes, error)
}

type Usecase struct {
	rp repo.RepoItf
	hc util.HttpClientItf
}

func NewUsecase(rp repo.RepoItf, hc util.HttpClientItf) *Usecase {
	return &Usecase{
		rp: rp,
		hc: hc,
	}
}

func (uc *Usecase) GetUnexpectedInfo(body []byte) error {
	var info dto.AlphaInfo
	readErr := json.Unmarshal(body, &info)
	if readErr != nil {
		return constant.ErrAlphaUnmarshal(readErr)
	}

	// Indicate if this is not an information-JSON body
	if info.Info == "" {
		return nil
	}

	// Erase any trace of my API key
	info.Info = strings.ReplaceAll(info.Info,
		os.Getenv("ALPHA_VANTAGE_API_KEY"), "[REDACTED]")

	// Simplify exceed-API-limit message
	if info.Info == constant.APIExceedLimit {
		return constant.ErrAPIExceed
	}

	// For any unexpected error I have never seen before
	return constant.NewCError(http.StatusBadGateway, info.Info)
}

func (uc *Usecase) PrevWeekend(t time.Time) time.Time {
	for {
		weekday := t.Weekday()
		if weekday == time.Saturday || weekday == time.Sunday {
			return t
		}
		t = t.AddDate(0, 0, -1)
	}
}

func (uc *Usecase) NextMonday(t time.Time) time.Time {
	for {
		weekday := t.Weekday()
		if weekday == time.Monday {
			return t
		}
		t = t.AddDate(0, 0, 1)
	}
}

func (uc *Usecase) NextFriday(t time.Time) time.Time {
	for {
		weekday := t.Weekday()
		if weekday == time.Friday {
			return t
		}
		t = t.AddDate(0, 0, 1)
	}
}

func (uc *Usecase) NextWeek(t time.Time) *dto.WeekRes {
	week := &dto.WeekRes{}
	week.Monday = uc.NextMonday(t)
	week.Friday = uc.NextFriday(week.Monday)
	week.DailyData = make([]dto.DailyOHLCVRes, 0)
	return week
}

func (uc *Usecase) ParseOHLCV(ctx *gin.Context, timeSeries *map[string]string) (*dto.DailyOHLCVRes, error) {
	TimeSeries := *timeSeries
	var ohlcv dto.DailyOHLCVRes
	ohlcv.OHLC = make(map[string]decimal.Decimal)

	// - OHLC
	for _, value := range []string{"1. open", "2. high",
		"3. low", "4. close"} {

		parts := strings.Split(value, " ")
		text, ok := TimeSeries[value]

		if !ok {
			return nil, constant.ErrAlphaParseBody(
				fmt.Sprintf("can't find %s price as usual", parts[1]),
			)
		}
		dec, err := decimal.NewFromString(text)

		if err != nil {
			return nil, constant.ErrAlphaParseBody(err.Error())
		}
		ohlcv.OHLC[parts[1]] = dec
	}

	// - Volume
	text, ok := TimeSeries["5. volume"]
	if !ok {
		return nil, constant.ErrAlphaParseBody(
			"can't find volume as usual")
	}
	vol, err := strconv.Atoi(text)
	if err != nil {
		return nil, constant.ErrAlphaParseBody(err.Error())
	}
	ohlcv.Volume = vol

	return &ohlcv, nil
}

func (uc *Usecase) GetSymbols(ctx *gin.Context, req *dto.GetSymbolsReq) (*dto.AlphaSymbolsRes, error) {
	// Retrieve data from Alpha Vantage API
	url := fmt.Sprintf("https://www.alphavantage.co/"+
		"query?function=SYMBOL_SEARCH"+
		"&keywords=%s&apikey=%s",
		req.Prefix,
		"demo",
	)

	response, err := uc.hc.Get(url)
	if err != nil {
		return nil, constant.ErrAlphaGet(err)
	}
	defer response.Body.Close()

	body, readErr := uc.hc.ReadAll(response.Body)
	if readErr != nil {
		return nil, constant.ErrAlphaReadAll(err)
	}

	// Check for e.g. API rate limit is exceeded
	err = uc.GetUnexpectedInfo(body)
	if err != nil {
		return nil, err
	}

	// Unmarshal body
	var symbols dto.AlphaSymbolsRes
	readErr = json.Unmarshal(body, &symbols)
	if readErr != nil {
		return nil, constant.ErrAlphaUnmarshal(readErr)
	}

	return &symbols, nil
}

func (uc *Usecase) CollectSymbol(ctx *gin.Context, req *dto.CollectSymbolReq) (*dto.StockDataRes, error) {
	// Check if symbol is in database already
	exists, err := uc.rp.CheckSymbolExists(ctx, req)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, constant.ErrStockAlready
	}

	// Retrieve data from Alpha Vantage API
	url := fmt.Sprintf("https://www.alphavantage.co/"+
		"query?function=TIME_SERIES_DAILY"+
		"&symbol=%s&apikey=%s",
		req.Symbol,
		"demo",
	)

	response, err := uc.hc.Get(url)

	if err != nil {
		return nil, constant.ErrAlphaGet(err)
	}
	defer response.Body.Close()

	body, readErr := uc.hc.ReadAll(response.Body)
	if readErr != nil {
		return nil, constant.ErrAlphaReadAll(err)
	}

	// Check for e.g. API rate limit is exceeded
	err = uc.GetUnexpectedInfo(body)
	log.Println(err)
	if err != nil {
		return nil, err
	}

	// Unmarshal body
	var alphaData dto.AlphaStockDataRes
	readErr = json.Unmarshal(body, &alphaData)
	if readErr != nil {
		return nil, constant.ErrAlphaUnmarshal(readErr)
	}

	alphaMeta := alphaData.MetaData

	// Process data from API:
	var stockData dto.StockDataRes
	var metaData dto.CollectSymbolMeta

	// 1. collect metadata
	metaData.Symbol = alphaMeta.Symbol

	t, err := time.Parse(constant.LayoutISO, alphaMeta.LastRefreshed)
	if err != nil {
		return nil, constant.ErrAlphaParseBody(err.Error())
	}
	metaData.LastRefreshed = t

	// (there is a default size of stocks to be recorded per symbol)
	stockData.MetaData = &metaData

	// 2. collect first constant.DefaultStocksNum days of time series data
	date := metaData.LastRefreshed.AddDate(0, 0,
		-constant.DefaultStocksNum+1)
	date = uc.PrevWeekend(date)
	timeSeries := make([]dto.DailyOHLCVRes, 0)
	for key, value := range alphaData.TimeSeries {
		keyDate, err := time.Parse(constant.LayoutISO, key)
		if err != nil {
			return nil, constant.ErrAlphaParseBody(err.Error())
		}

		if !keyDate.Before(date) {

			ohlcv, err := uc.ParseOHLCV(ctx, &value)

			if err != nil {
				return nil, constant.ErrAlphaParseBody(err.Error())
			}
			ohlcv.Day = keyDate
			timeSeries = append(timeSeries, *ohlcv)
		}
	}

	// - figure out number of time series data kept
	metaData.Size = len(timeSeries)

	// 3. sort the kept time series data
	sort.SliceStable(timeSeries, func(i, j int) bool {
		return timeSeries[i].Day.Before(
			timeSeries[j].Day,
		)
	})

	// Insert data
	err = uc.rp.InsertNewSymbolData(ctx, &stockData, timeSeries)
	if err != nil {
		return nil, err
	}

	// Processing to divide time series to weeks for presentation
	var weekIndex int
	stockData.Weeks = append(stockData.Weeks, uc.NextWeek(date))
	thisWeek := stockData.Weeks[weekIndex]
	for _, day := range timeSeries {
		if day.Day.After(thisWeek.Friday) {
			stockData.Weeks = append(stockData.Weeks, uc.NextWeek(thisWeek.Friday))
			weekIndex++
			thisWeek = stockData.Weeks[weekIndex]
		}

		thisWeek.DailyData = append(thisWeek.DailyData, day)
	}

	return &stockData, nil
}

func (uc *Usecase) DeleteSymbol(ctx *gin.Context, req *dto.DeleteSymbolReq) error {
	// repo
	return uc.rp.DeleteSymbol(ctx, req)
}

func (uc *Usecase) StoredData(ctx *gin.Context) ([]*dto.StockDataRes, error) {
	// repo
	return uc.rp.StoredData(ctx)
}
