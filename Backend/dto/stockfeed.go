package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

// General use containers for stock data
type SymbolDataMeta struct {
	Symbol        string    `json:"symbol"`
	LastRefreshed time.Time `json:"last_refreshed"`
	Size          int       `json:"size"`
}

type DailyOHLCVRes struct {
	Day    time.Time                  `json:"day"`
	OHLC   map[string]decimal.Decimal `json:"ohlc"`
	Volume int                        `json:"volume"`
}

type DataPerSymbol struct {
	MetaData   *SymbolDataMeta
	TimeSeries []DailyOHLCVRes
}

type WeekRes struct {
	Monday    time.Time       `json:"monday"`
	Friday    time.Time       `json:"friday"`
	DailyData []DailyOHLCVRes `json:"daily_data"`
}

type StockDataRes struct {
	MetaData *SymbolDataMeta `json:"meta_data"`
	Weeks    []*WeekRes      `json:"weeks_covered"`
}

// GetSymbols
type GetSymbolsReq struct {
	Prefix string
}

type GetSymbolsSingle struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	Region string `json:"region"`
}

type GetSymbolsRes struct {
	BestMatches []GetSymbolsSingle `json:"best_matches"`
}

// CollectSymbol
type CollectSymbolReq struct {
	Symbol string
}

// DeleteSymbol
type DeleteSymbolReq struct {
	Symbol string
}
