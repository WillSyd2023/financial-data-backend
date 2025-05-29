package handler

import (
	"Backend/constant"
	"Backend/dto"
	"Backend/usecase"

	"github.com/gin-gonic/gin"
)

type HandlerItf interface {
	GetSymbols(*gin.Context)
}

type Handler struct {
	uc usecase.UsecaseItf
}

func NewHandler(uc usecase.UsecaseItf) *Handler {
	return &Handler{
		uc: uc,
	}
}

func (hd *Handler) GetSymbols(ctx *gin.Context) {
	// request validation
	keywords := ctx.Query("keywords")
	if keywords == "" {
		ctx.Error(constant.ErrNoKeywords)
		return
	}
	var req dto.GetSymbolsReq
	req.Prefix = keywords

	// usecase
	hd.uc.GetSymbols(ctx, &req)
}
