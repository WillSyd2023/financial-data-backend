package middleware

import (
	"Backend/constant"
	"Backend/dto"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type MiddlewareItf interface {
	// error middleware
	Error() gin.HandlerFunc
	// authentication middleware
}

type Middleware struct {
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) Error() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there is no error
		if len(c.Errors) == 0 {
			return
		}

		// There is error; what error is it?
		err := c.Errors[0]

		// 1. Validation error from requests' JSON binding
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			validationErrors := make([]dto.ErrorType, 0)
			for _, fe := range ve {
				validationErrors = append(validationErrors, dto.ErrorType{
					Field:   fe.Field(),
					Message: fe.Error(),
				})
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.Res{
				Success: false,
				Error:   validationErrors,
			})
			return
		}

		// Custom error from `constant` repo
		var ce *constant.CustomError
		if errors.As(err, &ce) {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.Res{
				Success: false,
				Error:   ce.Error(),
			})
			return
		}

		// Unknown error, likely internal server error
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.Res{
			Success: false,
			Error:   err.Error(),
		})
	}
}
