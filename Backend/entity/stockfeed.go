package entity

type Symbol struct {
	Symbol string `json:"1. symbol"`
	Name   string `json:"2. name"`
	Region string `json:"4. region"`
}

type Symbols struct {
	BestMatches []Symbol `json:"bestMatches"`
}
