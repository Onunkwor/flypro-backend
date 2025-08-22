package services_test

import (
	"context"
	"testing"

	"github.com/onunkwor/flypro-backend/internal/config"
	"github.com/onunkwor/flypro-backend/internal/models"
	"github.com/onunkwor/flypro-backend/internal/repository"
	"github.com/onunkwor/flypro-backend/internal/repository/mocks"
	"github.com/onunkwor/flypro-backend/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	svc := services.NewUserService(config.Redis, mockRepo)

	existingUser := &models.User{Email: "test@example.com"}

	// Repo should return an existing user when searching
	mockRepo.On("FindByEmail", existingUser.Email).Return(existingUser, nil)

	err := svc.CreateUser(existingUser)

	assert.Equal(t, services.ErrEmailAlreadyExists, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	svc := services.NewUserService(config.Redis, mockRepo)

	newUser := &models.User{Email: "cake@gmail.com"}

	// Repo should return nil user when searching
	mockRepo.On("FindByEmail", newUser.Email).Return(nil, nil)
	// Repo should successfully create the user
	mockRepo.On("CreateUser", newUser).Return(nil)

	err := svc.CreateUser(newUser)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	svc := services.NewUserService(config.Redis, mockRepo)
	userID := uint(2)

	mockRepo.On("GetUserByID", userID).Return(nil, repository.ErrUserNotFound)

	user, err := svc.GetUserByID(context.Background(), userID)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, repository.ErrUserNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)
	svc := services.NewUserService(config.Redis, mockRepo)
	userID := uint(3)
	expectedUser := &models.User{ID: userID, Email: "cake@gmail.com"}
	mockRepo.On("GetUserByID", userID).Return(expectedUser, nil)

	user, err := svc.GetUserByID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}
