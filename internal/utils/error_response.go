package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidationErrorResponse(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error":   "validation error",
		"details": err.Error(),
	})
}

func InternalServerErrorResponse(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error":   "internal server error",
		"details": err.Error(),
	})
}
