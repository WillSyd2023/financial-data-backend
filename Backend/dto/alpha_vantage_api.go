package dto

type AlphaSymbolRes struct {
	Symbol string `json:"1. symbol"`
	Name   string `json:"2. name"`
	Region string `json:"4. region"`
}

type AlphaSymbolsRes struct {
	BestMatches []AlphaSymbolRes `json:"bestMatches"`
}
