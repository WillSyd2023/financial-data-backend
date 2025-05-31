package usecase

import (
	"Backend/constant"
	"Backend/dto"
	"Backend/repo"
	"Backend/util"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type UsecaseItf interface {
	parseOHLCV(*gin.Context, *map[string]string) (*dto.DailyOHLCVRes, error)
	GetSymbols(*gin.Context, *dto.GetSymbolsReq) (*dto.AlphaSymbolsRes, error)
	CollectSymbol(*gin.Context, *dto.CollectSymbolReq) (*dto.StockDataRes, error)
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

func (uc *Usecase) parseOHLCV(ctx *gin.Context, timeSeries *map[string]string) (*dto.DailyOHLCVRes, error) {
	TimeSeries := *timeSeries
	var ohlcv dto.DailyOHLCVRes

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
		os.Getenv("ALPHA_VANTAGE_API_KEY"),
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

	// Unmarshal body
	var symbols dto.AlphaSymbolsRes
	readErr = json.Unmarshal(body, &symbols)
	if readErr != nil {
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

	body, readErr := uc.hc.ReadAll(response.Body)
	if readErr != nil {
		return nil, constant.ErrAlphaReadAll(err)
	}

	// Unmarshal body
	var alphaData dto.AlphaStockDataRes
	readErr = json.Unmarshal(body, &alphaData)
	if readErr != nil {
		return nil, constant.ErrAlphaUnmarshal(err)
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
	metaData.Size = constant.DefaultStocksNum

	stockData.MetaData = metaData

	// 2. collect first constant.DefaultStocksNum days of time series data
	earliestDate := metaData.LastRefreshed.AddDate(0, 0, -metaData.Size+1)
	for key, value := range alphaData.TimeSeries {
		keyDate, err := time.Parse(constant.LayoutISO, key)
		if err != nil {
			return nil, constant.ErrAlphaParseBody(err.Error())
		}

		if !keyDate.Before(earliestDate) {
			ohlcv, err := uc.parseOHLCV(ctx, &value)
			if err != nil {
				return nil, constant.ErrAlphaParseBody(err.Error())
			}
			ohlcv.Day = keyDate
			stockData.TimeSeries = append(stockData.TimeSeries, ohlcv)
		}
	}

	return &stockData, nil
}
