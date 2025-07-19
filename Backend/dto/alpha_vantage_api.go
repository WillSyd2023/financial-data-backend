package dto

type AlphaInfo struct {
	Info string `json:"Information"`
}

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
type AlphaCollectSymbolMeta struct {
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	OutputSize    string `json:"4. Output Size"`
}

type AlphaStockDataRes struct {
	MetaData   AlphaCollectSymbolMeta         `json:"Meta Data"`
	TimeSeries map[string](map[string]string) `json:"Time Series (Daily)"`
}
