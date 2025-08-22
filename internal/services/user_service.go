package services

import (
	"errors"

	"github.com/onunkwor/flypro-backend/internal/models"
	"github.com/onunkwor/flypro-backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user *models.User) error {
	existing, err := s.repo.FindByEmail(user.Email)
	if err == nil && existing != nil {
		return ErrEmailAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return s.repo.CreateUser(user)
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil

}
