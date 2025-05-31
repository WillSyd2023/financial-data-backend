package dto

// GetSymbols
type AlphaSymbolRes struct {
	Symbol string `json:"1. symbol"`
	Name   string `json:"2. name"`
	Region string `json:"4. region"`
}

type AlphaSymbolsRes struct {
	BestMatches []AlphaSymbolRes `json:"bestMatches"`
}

// CollectSymbol
type AlphaDailyOHLCVRes struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

type AlphaStockDataRes struct {
	MetaData   MetaData                      `json:"Meta Data"`
	TimeSeries map[string]AlphaDailyOHLCVRes `json:"Time Series (Daily)"`
}
