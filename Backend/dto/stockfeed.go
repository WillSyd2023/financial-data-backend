package dto

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

type DateOnly time.Time

func (d DateOnly) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d).Format("2006-01-02"))
}

func (d DateOnly) Weekday() time.Weekday {
	return time.Time(d).Weekday()
}

func (d DateOnly) AddDate(years, months, days int) DateOnly {
	return DateOnly(time.Time(d).AddDate(years, months, days))
}

func (d DateOnly) Before(e DateOnly) bool {
	return time.Time(d).Before(time.Time(e))
}

func (d DateOnly) After(e DateOnly) bool {
	return time.Time(d).After(time.Time(e))
}

// General use containers for stock data
type SymbolDataMeta struct {
	Symbol        string   `json:"symbol"`
	LastRefreshed DateOnly `json:"last_refreshed"`
	Size          int      `json:"size"`
}

type DailyOHLCVRes struct {
	Day    DateOnly                   `json:"day"`
	OHLC   map[string]decimal.Decimal `json:"ohlc"`
	Volume int                        `json:"volume"`
}

type DataPerSymbol struct {
	MetaData   *SymbolDataMeta
	TimeSeries []DailyOHLCVRes
}

type WeekRes struct {
	Monday    DateOnly        `json:"monday"`
	Friday    DateOnly        `json:"friday"`
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
