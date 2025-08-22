package mocks

import (
	"github.com/onunkwor/flypro-backend/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepo) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if u, ok := args.Get(0).(*models.User); ok {
		return u, args.Error(1)
	}

	return nil, args.Error(1)
}
func (m *MockUserRepo) GetUserByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if u, ok := args.Get(0).(*models.User); ok {
		return u, args.Error(1)
	}
	return nil, args.Error(1)
}
