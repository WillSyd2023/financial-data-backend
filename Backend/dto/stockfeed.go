package dto

type GetSymbolsReq struct {
	Prefix string `json:"prefix" binding:"required,gte=1"`
}
