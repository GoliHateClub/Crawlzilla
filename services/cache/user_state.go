package cache

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type UserState struct {
	redis *redis.Client
}

func CreateUserStateCache(ctx context.Context) *UserState {
	client := ctx.Value("redis").(*redis.Client)

	return &UserState{
		redis: client,
	}
}

func (s *UserState) SetUserState(ctx context.Context, chatID int64, state string) error {
	key := getStateKey(chatID)
	return s.redis.Set(ctx, key, state, time.Hour).Err()
}

func (s *UserState) GetUserState(ctx context.Context, chatID int64) (string, error) {
	key := getStateKey(chatID)
	state, err := s.redis.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return state, err
}

func (s *UserState) ClearUserState(ctx context.Context, chatID int64) error {
	key := getStateKey(chatID)
	return s.redis.Del(ctx, key).Err()
}

func getStateKey(chatID int64) string {
	return "conversation_state:" + strconv.FormatInt(chatID, 10)
}
