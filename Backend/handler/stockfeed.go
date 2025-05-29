package handler

import (
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

func NewUsecase(uc usecase.UsecaseItf) *Handler {
	return &Handler{
		uc: uc,
	}
}

func (hd *Handler) GetSymbols(ctx *gin.Context) {
	// request validation
	var req dto.GetSymbolsReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}
}
