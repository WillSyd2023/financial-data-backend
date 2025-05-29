package dto

type SymbolRes struct {
	Symbol string `json:"1. symbol"`
	Name   string `json:"2. name"`
	Region string `json:"4. region"`
}

type SymbolsRes struct {
	BestMatches []SymbolRes `json:"bestMatches"`
}
