package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidationErrorResponse(c *gin.Context, err map[string]string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error":   "validation error",
		"details": err,
	})
}

func InternalServerErrorResponse(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error":   "internal server error",
		"details": err.Error(),
	})
}

func FormatValidationError(err error) map[string]string {
	out := map[string]string{}
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range errs {
			out[fe.Field()] = fe.Tag()
		}
	} else {
		out["error"] = err.Error()
	}
	return out
}
