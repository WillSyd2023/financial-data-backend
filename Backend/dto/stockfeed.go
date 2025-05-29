package dto

type GetSymbolsReq struct {
	Prefix string `json:"prefix" binding:"required,gte=1"`
}

type GetSymbolsSingle struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	Region string `json:"region"`
}

type GetSymbolsRes struct {
	BestMatches []GetSymbolsSingle `json:"best_matches"`
}
