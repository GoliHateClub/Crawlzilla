package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

type UserState struct {
	redis *redis.Client
}

type State struct {
	Conversation string
	Stage        string
	UserId       int64
	ChatId       int64
}

func CreateNewState(conversation string, data *tgbotapi.CallbackQuery) State {
	return State{
		UserId:       int64(data.From.ID),
		ChatId:       data.Message.Chat.ID,
		Conversation: conversation,
		Stage:        "init",
	}
}

func CreateUserStateCache(ctx context.Context) *UserState {
	client := ctx.Value("redis").(*redis.Client)

	return &UserState{
		redis: client,
	}
}

func (s *UserState) SetUserState(ctx context.Context, chatID int64, state State) error {
	key := getStateKey(chatID)

	marshal, err := json.Marshal(state)

	if err != nil {
		return err
	}

	return s.redis.Set(ctx, key, marshal, 0).Err()
}

func (s *UserState) GetUserState(ctx context.Context, chatID int64) (State, error) {
	key := getStateKey(chatID)
	stateStr, err := s.redis.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return State{}, nil
	}

	var state State
	err = json.Unmarshal([]byte(stateStr), &state)
	if err != nil {
		return State{}, err
	}

	return state, nil
}

func (s *UserState) ClearUserState(ctx context.Context, chatID int64) error {
	key := getStateKey(chatID)
	return s.redis.Del(ctx, key).Err()
}

func getStateKey(chatID int64) string {
	return "conversation_state:" + strconv.FormatInt(chatID, 10)
}
