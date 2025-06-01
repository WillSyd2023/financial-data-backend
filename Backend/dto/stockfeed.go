package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

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

type CollectSymbolMeta struct {
	Symbol        string    `json:"symbol"`
	LastRefreshed time.Time `json:"last_refreshed"`
	Size          int       `json:"size"`
}

type DailyOHLCVRes struct {
	Day    time.Time                  `json:"day"`
	OHLC   map[string]decimal.Decimal `json:"ohlc"`
	Volume int                        `json:"volume"`
}

type StockDataRes struct {
	MetaData   CollectSymbolMeta `json:"meta_data"`
	TimeSeries []DailyOHLCVRes   `json:"daily_time_series"`
}
