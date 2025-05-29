package usecase

import (
	"Backend/constant"
	"Backend/dto"
	"Backend/repo"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type UsecaseItf interface {
	GetSymbols(*gin.Context, *dto.GetSymbolsReq) error
}

type Usecase struct {
	rp repo.RepoItf
}

func NewUsecase(rp repo.RepoItf) *Usecase {
	return &Usecase{
		rp: rp,
	}
}

func (uc *Usecase) GetSymbols(ctx *gin.Context, req *dto.GetSymbolsReq) error {
	// Retrieve data from Alpha Vantage API
	url := fmt.Sprintf("https://www.alphavantage.co/"+
		"query?function=SYMBOL_SEARCH"+
		"&keywords=%s&apikey=%s",
		req.Prefix,
		os.Getenv("ALPHA_VANTAGE_API_KEY"),
	)

	response, err := http.Get(url)
	if err != nil {
		return constant.NewCError(
			http.StatusBadGateway,
			fmt.Sprintf(
				"Alpha Vantage API GET error: %s",
				err.Error(),
			),
		)
	}
	defer response.Body.Close()

	body, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		return constant.NewCError(
			http.StatusBadGateway,
			fmt.Sprintf(
				"Alpha Vantage API response-body-parse error: %s",
				readErr.Error(),
			),
		)
	}
	log.Print(string(body))

	return nil
}
