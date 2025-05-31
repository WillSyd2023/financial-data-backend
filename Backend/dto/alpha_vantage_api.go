package dto

import (
	"github.com/shopspring/decimal"
)

type AlphaSymbolRes struct {
	Symbol string `json:"1. symbol"`
	Name   string `json:"2. name"`
	Region string `json:"4. region"`
}

type AlphaSymbolsRes struct {
	BestMatches []AlphaSymbolRes `json:"bestMatches"`
}

type DailyOHLCV struct {
	Open   decimal.Decimal `json:"1. open"`
	High   decimal.Decimal `json:"2. high"`
	Low    decimal.Decimal `json:"3. low"`
	Close  decimal.Decimal `json:"4. close"`
	Volume int             `json:"5. volume"`
}

type StockData struct {
	MetaData struct {
		Symbol        string `json:"2. Symbol"`
		LastRefreshed string `json:"3. Last Refreshed"`
		OutputSize    string `json:"4. Output Size"`
	} `json:"Meta Data"`
	TimeSeries map[string]DailyOHLCV `json:"Time Series (Daily)"`
}
