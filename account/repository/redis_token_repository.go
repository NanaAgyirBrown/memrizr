package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/nanaagyirbrown/memrizr/handler/model"
	"github.com/nanaagyirbrown/memrizr/handler/model/apperrors"
	"log"
	"time"
)

type redisTokenRepository struct {
	Redis *redis.Client
}

func NewTokenRepository(redisClient *redis.Client) model.TokenRepository  {
	return &redisTokenRepository{
		Redis: redisClient,
	}
}

// SetRefreshToken stores a refresh token with an expiry time
func (r redisTokenRepository) SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error {
	log.Printf("SetRefreshToken called : %v\n", userID)
	// We'll store userID with token id so we can scan (non-blocking)
	// over the user's tokens and delete them in case of token leakage
	key := fmt.Sprintf("%s:%s", userID, tokenID)
	x := r.Redis.Set(ctx, key, 0, expiresIn)
	//err := r.Redis.Set(ctx, key, 0, expiresIn).Err();
	if x.Err() != nil {
		log.Printf("Could not SET refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, x.Err())
		return apperrors.NewInternal()
	}

	log.Printf(" %v", x)

	return nil
}

// DeleteRefreshToken used to delete old  refresh tokens
// Services may access this to revolve tokens
func (r redisTokenRepository) DeleteRefreshToken(ctx context.Context, userID string, prevTokenID string) error {
	key := fmt.Sprintf("%s:%s", userID, prevTokenID)
	if err := r.Redis.Del(ctx, key).Err(); err != nil {
		log.Printf("Could not delete refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, prevTokenID, err)
		return apperrors.NewInternal()
	}

	return nil
}

