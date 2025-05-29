package dto

type GetSymbolsReq struct {
	Prefix string `json:"prefix" binding:"required,gte=1"`
}

type SymbolRes struct {
	Symbol string `json:"1. symbol"`
	Name   string `json:"2. name"`
	Region string `json:"4. region"`
}

type GetSymbolsRes struct {
	BestMatches []SymbolRes `json:"bestMatches"`
}
