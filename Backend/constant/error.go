package constant

import (
	"fmt"
	"net/http"
)

type CustomError struct {
	StatusCode int
	Message    string
}

func NewCError(StatusCode int, Message string) CustomError {
	return CustomError{StatusCode: StatusCode, Message: Message}
}

func (err CustomError) Error() string {
	return err.Message
}

var (
	// GetSymbols handler
	ErrNoKeywords = NewCError(http.StatusBadRequest, "please provide keywords")

	// CollectSymbol handler
	ErrNoSymbol     = NewCError(http.StatusBadRequest, "please provide symbol")
	ErrStockAlready = NewCError(http.StatusBadRequest, "The stock (symbol) is already tracked in the database and monitored regularly")
)

// Fetching data from e.g. Alpha Vantage API

func ErrAlphaGet(err error) error {
	return NewCError(
		http.StatusBadGateway,
		fmt.Sprintf(
			"Alpha Vantage API GET error: %s",
			err.Error(),
		),
	)
}

func ErrAlphaReadAll(err error) error {
	return NewCError(
		http.StatusBadGateway,
		fmt.Sprintf(
			"Alpha Vantage API body-io.ReadAll-parse error: %s",
			err.Error(),
		),
	)
}

func ErrAlphaUnmarshal(err error) error {
	return NewCError(
		http.StatusBadGateway,
		fmt.Sprintf(
			"Alpha Vantage API body-json.Unmarshal-parse error: %s",
			err.Error(),
		),
	)
}

func ErrAlphaParseDate(err error) error {
	return NewCError(
		http.StatusBadGateway,
		fmt.Sprintf(
			"Alpha Vantage API last-refreshed-date-parse error: %s",
			err.Error(),
		),
	)
}
