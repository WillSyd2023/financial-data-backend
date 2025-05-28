package handler

import "Backend/usecase"

type HandlerItf interface{}

type Handler struct {
	uc usecase.UsecaseItf
}

func NewUsecase(uc usecase.UsecaseItf) *Handler {
	return &Handler{
		uc: uc,
	}
}
