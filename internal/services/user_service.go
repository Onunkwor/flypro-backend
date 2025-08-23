package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
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

func (s *UserService) CreateUser(user *models.User) error {
	existing, err := s.repo.FindByEmail(user.Email)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return err
	}
	if existing != nil {
		return ErrEmailAlreadyExists
	}
	return s.repo.CreateUser(user)
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	key := fmt.Sprintf("user:%d", id)
	if s.redis != nil {
		val, err := s.redis.Get(ctx, key).Result()
		if err == nil {
			var user models.User
			if unmarshalErr := json.Unmarshal([]byte(val), &user); unmarshalErr == nil {
				return &user, nil
			}
		} else if err != redis.Nil {
			return nil, err
		}
	}
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	bytes, _ := json.Marshal(user)
	if s.redis != nil {
		s.redis.Set(ctx, key, bytes, time.Hour)
	}
	return user, nil

}
