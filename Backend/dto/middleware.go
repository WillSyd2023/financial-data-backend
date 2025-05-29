package dto

type Res struct {
	Success bool `json:"success"`
	Error   any  `json:"error"`
	Data    any  `json:"data"`
}

type ErrorType struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}
