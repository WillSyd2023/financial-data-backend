package usecase

import (
	"Backend/dto"
	"Backend/repo"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type UsecaseItf interface {
	GetSymbols(*gin.Context, *dto.GetSymbolsReq)
}

type Usecase struct {
	rp repo.RepoItf
}

func NewUsecase(rp repo.RepoItf) *Usecase {
	return &Usecase{
		rp: rp,
	}
}

func (uc *Usecase) GetSymbols(ctx *gin.Context, req *dto.GetSymbolsReq) {
	// Retrieve data from Alpha Vantage API
	url := fmt.Sprintf("https://www.alphavantage.co/"+
		"query?function=SYMBOL_SEARCH"+
		"&keywords=%s&apikey=%s",
		req.Prefix,
		os.Getenv("ALPHA_VANTAGE_API_KEY"),
	)

	response, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer response.Body.Close()

	body, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		fmt.Print(readErr.Error())
	}
	fmt.Println(string(body))
}
