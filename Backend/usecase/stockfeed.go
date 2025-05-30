package usecase

import (
	"Backend/constant"
	"Backend/dto"
	"Backend/repo"
	"Backend/util"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

type UsecaseItf interface {
	GetSymbols(*gin.Context, *dto.GetSymbolsReq) (*dto.AlphaSymbolsRes, error)
	CollectSymbol(*gin.Context, *dto.CollectSymbolReq) error
}

type Usecase struct {
	rp repo.RepoItf
	hc util.HttpClientItf
}

func NewUsecase(rp repo.RepoItf, hc util.HttpClientItf) *Usecase {
	return &Usecase{
		rp: rp,
		hc: hc,
	}
}

func (uc *Usecase) GetSymbols(ctx *gin.Context, req *dto.GetSymbolsReq) (*dto.AlphaSymbolsRes, error) {
	// Retrieve data from Alpha Vantage API
	url := fmt.Sprintf("https://www.alphavantage.co/"+
		"query?function=SYMBOL_SEARCH"+
		"&keywords=%s&apikey=%s",
		req.Prefix,
		os.Getenv("ALPHA_VANTAGE_API_KEY"),
	)

	response, err := uc.hc.Get(url)
	if err != nil {
		return nil, constant.ErrAlphaGet(err)
	}
	defer response.Body.Close()

	body, readErr := uc.hc.ReadAll(response.Body)
	if readErr != nil {
		return nil, constant.ErrAlphaReadAll(err)
	}

	// Unmarshal body
	var symbols dto.AlphaSymbolsRes
	readErr = json.Unmarshal(body, &symbols)
	if readErr != nil {
		return nil, constant.ErrAlphaUnmarshal(err)
	}

	return &symbols, nil
}

func (uc *Usecase) CollectSymbol(ctx *gin.Context, req *dto.CollectSymbolReq) error {
	// Check if symbol is in database already
	exists, err := uc.rp.CheckSymbolExists(ctx, req)
	if err != nil {
		return err
	}
	if exists {
		return constant.ErrStockAlready
	}

	// Retrieve data from Alpha Vantage API
	url := fmt.Sprintf("https://www.alphavantage.co/"+
		"query?function=TIME_SERIES_DAILY"+
		"&symbol=%s&apikey=%s",
		req.Symbol,
		os.Getenv("ALPHA_VANTAGE_API_KEY"),
	)

	response, err := uc.hc.Get(url)
	if err != nil {
		return constant.ErrAlphaGet(err)
	}
	defer response.Body.Close()

	body, readErr := uc.hc.ReadAll(response.Body)
	if readErr != nil {
		return constant.ErrAlphaReadAll(err)
	}

	return nil
}
