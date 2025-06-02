package handler

import (
	"Backend/constant"
	"Backend/dto"
	"Backend/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerItf interface {
	GetSymbols(*gin.Context)
	CollectSymbol(*gin.Context)
	DeleteSymbol(*gin.Context)
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
	symbols, err := hd.uc.GetSymbols(ctx, &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	// process response before returning
	var GetSymbolsRes dto.GetSymbolsRes
	GetSymbolsRes.BestMatches = make([]dto.GetSymbolsSingle, 0)
	for _, symbol := range symbols.BestMatches {
		var singleRes dto.GetSymbolsSingle
		singleRes.Symbol = symbol.Symbol
		singleRes.Name = symbol.Name
		singleRes.Region = symbol.Region
		GetSymbolsRes.BestMatches = append(
			GetSymbolsRes.BestMatches,
			singleRes,
		)
	}

	// return response
	ctx.JSON(http.StatusOK,
		gin.H{
			"message": nil,
			"error":   nil,
			"data":    GetSymbolsRes,
		})
}

func (hd *Handler) CollectSymbol(ctx *gin.Context) {
	// request validation
	symbol := ctx.Param("symbol")
	if symbol == "" {
		ctx.Error(constant.ErrNoSymbol)
		return
	}
	var req dto.CollectSymbolReq
	req.Symbol = symbol

	// usecase
	stockData, err := hd.uc.CollectSymbol(ctx, &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated,
		gin.H{
			"message": nil,
			"error":   nil,
			"data":    stockData,
		})
}

func (hd *Handler) DeleteSymbol(ctx *gin.Context) {
	// request validation
	symbol := ctx.Param("symbol")
	if symbol == "" {
		ctx.Error(constant.ErrNoSymbol)
		return
	}
}
