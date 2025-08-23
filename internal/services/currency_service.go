package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

type CurrencyService struct {
	redis  *redis.Client
	apiKey string
}

type CurrencyConverter interface {
	Convert(ctx context.Context, amount float64, from, to string) (float64, error)
}

func NewCurrencyService(r *redis.Client, apiKey string) *CurrencyService {
	return &CurrencyService{redis: r, apiKey: apiKey}
}

func (s *CurrencyService) getRate(ctx context.Context, from, to string) (float64, error) {
	key := fmt.Sprintf("rate:%s:%s", from, to)

	val, err := s.redis.Get(ctx, key).Result()
	if err == nil {
		var rate float64
		json.Unmarshal([]byte(val), &rate)
		return rate, nil
	}
	url := fmt.Sprintf("https://open.er-api.com/v6/latest/%s", from)
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch exchange rate: %w", err)
	}
	defer resp.Body.Close()
	var data struct {
		Rates map[string]float64 `json:"rates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, fmt.Errorf("failed to parse exchange rate: %w", err)
	}
	rate, ok := data.Rates[to]
	if !ok {
		return 0, fmt.Errorf("currency %s not supported", to)
	}

	// 3. Cache in Redis (6 hours)
	rateJSON, _ := json.Marshal(rate)
	s.redis.Set(ctx, key, rateJSON, 6*time.Hour)

	return rate, nil
}

func (s *CurrencyService) Convert(ctx context.Context, amount float64, from, to string) (float64, error) {
	rate, err := s.getRate(ctx, from, to)
	if err != nil {
		return 0, err
	}
	return amount * rate, nil
}
