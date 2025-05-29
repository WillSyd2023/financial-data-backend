package constant

type CustomError struct {
	Message string
}

func NewCError(Message string) CustomError {
	return CustomError{Message: Message}
}

func (err CustomError) Error() string {
	return err.Message
}

var (
	// GetSymbols handler
	ErrNoKeywords = NewCError("please provide keywords")
)
