package dto

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
