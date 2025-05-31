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
	Open   decimal.Decimal `json:"open"`
	High   decimal.Decimal `json:"high"`
	Low    decimal.Decimal `json:"low"`
	Close  decimal.Decimal `json:"close"`
	Volume int             `json:"volume"`
}

type StockDataRes struct {
	MetaData   CollectSymbolMeta           `json:"meta_data"`
	TimeSeries map[time.Time]DailyOHLCVRes `json:"daily_time_series"`
}
