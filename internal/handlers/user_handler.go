package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onunkwor/flypro-backend/internal/dto"
	"github.com/onunkwor/flypro-backend/internal/models"
	"github.com/onunkwor/flypro-backend/internal/repository"
	"github.com/onunkwor/flypro-backend/internal/services"
	"github.com/onunkwor/flypro-backend/internal/utils"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}
func (h *UserHandler) CreateUser(c *gin.Context) {

	var request dto.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		formatted := utils.FormatValidationError(err)
		utils.ValidationErrorResponse(c, formatted)
		return
	}

	user := models.User{
		Email: request.Email,
		Name:  request.Name,
	}
	if err := h.service.CreateUser(&user); err != nil {
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}
		utils.InternalServerErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
		"user":    user,
	})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		utils.InternalServerErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}
