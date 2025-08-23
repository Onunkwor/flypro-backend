package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/onunkwor/flypro-backend/internal/dto"
	"github.com/onunkwor/flypro-backend/internal/models"
	"github.com/onunkwor/flypro-backend/internal/repository"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type UserService struct {
	repo  repository.UserRepository
	redis *redis.Client
}

func NewUserService(r *redis.Client, repo repository.UserRepository) *UserService {
	return &UserService{redis: r, repo: repo}
}

func (s *UserService) CreateUser(user *models.User) (dto.UserResponse, error) {
	existing, err := s.repo.FindByEmail(user.Email)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return dto.UserResponse{}, err
	}
	if existing != nil {
		return dto.UserResponse{}, ErrEmailAlreadyExists
	}

	if err := s.repo.CreateUser(user); err != nil {
		return dto.UserResponse{}, err
	}

	return dto.NewUserResponse(user), nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (dto.UserResponse, error) {
	key := fmt.Sprintf("user:%d", id)
	if s.redis != nil {
		val, err := s.redis.Get(ctx, key).Result()
		if err == nil {
			var user dto.UserResponse
			if unmarshalErr := json.Unmarshal([]byte(val), &user); unmarshalErr == nil {
				return user, nil
			}
		} else if err != redis.Nil {
			return dto.UserResponse{}, err
		}
	}

	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return dto.UserResponse{}, err
	}

	response := dto.NewUserResponse(user)
	bytes, _ := json.Marshal(response)
	if s.redis != nil {
		s.redis.Set(ctx, key, bytes, time.Hour)
	}
	return response, nil

}
