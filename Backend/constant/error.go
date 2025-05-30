package constant

import "net/http"

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
