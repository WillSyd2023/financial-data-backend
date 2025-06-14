package usecase

import (
	"Backend/constant"
	"Backend/dto"
	"Backend/repo"
	"Backend/util"
	"encoding/json"
	"fmt"
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
	ParseOHLCV(*gin.Context, *map[string]string) (*dto.DailyOHLCVRes, error)
	PrevWeekend(dto.DateOnly) dto.DateOnly
	NextWeek(dto.DateOnly) *dto.WeekRes
	BuildStockData(*dto.DataPerSymbol) *dto.StockDataRes

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
	err := json.Unmarshal(body, &info)
	if err != nil {
		return constant.ErrAlphaUnmarshal(err)
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

func (uc *Usecase) PrevWeekend(t dto.DateOnly) dto.DateOnly {
	for {
		weekday := t.Weekday()
		if weekday == time.Saturday || weekday == time.Sunday {
			return t
		}
		t = t.AddDate(0, 0, -1)
	}
}

func (uc *Usecase) NextWeek(t dto.DateOnly) *dto.WeekRes {
	week := &dto.WeekRes{}
	t1 := t
	for {
		weekday := t1.Weekday()
		if weekday == time.Monday {
			week.Monday = t1
			break
		}
		t1 = t1.AddDate(0, 0, 1)
	}
	t2 := week.Monday
	for {
		weekday := t2.Weekday()
		if weekday == time.Friday {
			week.Friday = t2
			break
		}
		t2 = t2.AddDate(0, 0, 1)
	}
	week.DailyData = make([]dto.DailyOHLCVRes, 0)
	return week
}

func (uc *Usecase) BuildStockData(data *dto.DataPerSymbol) *dto.StockDataRes {
	var stockData dto.StockDataRes
	stockData.MetaData = data.MetaData

	// Processing to divide time series to weeks for presentation
	var weekIndex int
	date := data.TimeSeries[0].Day
	date = uc.PrevWeekend(date)
	stockData.Weeks = append(stockData.Weeks, uc.NextWeek(date))
	thisWeek := stockData.Weeks[weekIndex]
	for _, day := range data.TimeSeries {
		if day.Day.After(thisWeek.Friday) {
			stockData.Weeks = append(stockData.Weeks, uc.NextWeek(thisWeek.Friday))
			weekIndex++
			thisWeek = stockData.Weeks[weekIndex]
		}

		thisWeek.DailyData = append(thisWeek.DailyData, day)
	}
	return &stockData
}

func (uc *Usecase) GetSymbols(ctx *gin.Context, req *dto.GetSymbolsReq) (*dto.AlphaSymbolsRes, error) {
	// Retrieve data from Alpha Vantage API
	url := fmt.Sprintf("https://www.alphavantage.co/"+
		"query?function=SYMBOL_SEARCH"+
		"&keywords=%s&apikey=%s",
		req.Prefix,
		os.Getenv("ALPHA_VANTAGE_API_KEY"),
	)

	response, err := uc.hc.Get(url)
	if err != nil {
		return nil, constant.ErrAlphaGet(err)
	}
	defer response.Body.Close()

	body, err := uc.hc.ReadAll(response.Body)
	if err != nil {
		return nil, constant.ErrAlphaReadAll(err)
	}

	// Check for e.g. API rate limit is exceeded
	err = uc.GetUnexpectedInfo(body)
	if err != nil {
		return nil, err
	}

	// Unmarshal body
	var symbols dto.AlphaSymbolsRes
	err = json.Unmarshal(body, &symbols)
	if err != nil {
		return nil, constant.ErrAlphaUnmarshal(err)
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
		os.Getenv("ALPHA_VANTAGE_API_KEY"),
	)

	response, err := uc.hc.Get(url)
	if err != nil {
		return nil, constant.ErrAlphaGet(err)
	}
	defer response.Body.Close()

	body, err := uc.hc.ReadAll(response.Body)
	if err != nil {
		return nil, constant.ErrAlphaReadAll(err)
	}

	// Check for e.g. API rate limit is exceeded
	err = uc.GetUnexpectedInfo(body)
	if err != nil {
		return nil, err
	}

	// Unmarshal body
	var alphaData dto.AlphaStockDataRes
	err = json.Unmarshal(body, &alphaData)
	if err != nil {
		return nil, constant.ErrAlphaUnmarshal(err)
	}

	alphaMeta := alphaData.MetaData

	// Process data from API:
	var metaData dto.SymbolDataMeta

	// 1. collect some metadata
	metaData.Symbol = alphaMeta.Symbol

	t, err := time.Parse(constant.LayoutISO, alphaMeta.LastRefreshed)
	if err != nil {
		return nil, constant.ErrAlphaParseBody(err.Error())
	}
	metaData.LastRefreshed = dto.DateOnly(t)

	// 2. collect first constant.DefaultStocksNum days of time series data
	date := metaData.LastRefreshed.AddDate(0, 0,
		-constant.DefaultStocksNum+1)
	date = uc.PrevWeekend(date)
	dateTime := time.Time(date)
	timeSeries := make([]dto.DailyOHLCVRes, 0)
	for key, value := range alphaData.TimeSeries {
		keyDate, err := time.Parse(constant.LayoutISO, key)
		if err != nil {
			return nil, constant.ErrAlphaParseBody(err.Error())
		}

		if !keyDate.Before(dateTime) {

			ohlcv, err := uc.ParseOHLCV(ctx, &value)

			if err != nil {
				return nil, constant.ErrAlphaParseBody(err.Error())
			}
			ohlcv.Day = dto.DateOnly(keyDate)
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

	dataForSym := &dto.DataPerSymbol{
		MetaData: &metaData, TimeSeries: timeSeries}

	// Insert data
	err = uc.rp.InsertNewSymbolData(ctx, dataForSym)
	if err != nil {
		return nil, err
	}

	// Processing to divide time series to weeks for presentation
	// just before returning
	return uc.BuildStockData(dataForSym), nil
}

func (uc *Usecase) DeleteSymbol(ctx *gin.Context, req *dto.DeleteSymbolReq) error {
	// repo
	return uc.rp.DeleteSymbol(ctx, req)
}

func (uc *Usecase) StoredData(ctx *gin.Context) ([]*dto.StockDataRes, error) {
	// repo
	dataPerSymbol, err := uc.rp.StoredData(ctx)
	if err != nil {
		return nil, err
	}

	// assemble data for presentation
	stockData := make([]*dto.StockDataRes, 0)
	for _, datum := range dataPerSymbol {
		stockData = append(stockData, uc.BuildStockData(&datum))
	}

	return stockData, nil
}
