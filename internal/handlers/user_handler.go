package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/onunkwor/flypro-backend/internal/config"
	"github.com/onunkwor/flypro-backend/internal/dto"
	"github.com/onunkwor/flypro-backend/internal/models"
)

func CreateUser(c *gin.Context) {

	var request dto.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
		})
		return
	}

	user := models.User{
		Email: request.Email,
		Name:  request.Name,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to create user",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
		"user":    user,
	})
}
