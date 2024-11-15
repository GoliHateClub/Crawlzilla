package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

type ActionCache struct {
	redis *redis.Client
}

type ActionState struct {
	Conversation string
	Action       string
	UserId       int64
	ChatId       int64
	ActionData   map[string]interface{}
}

func CreateNewActionState(conversation string, data *tgbotapi.CallbackQuery, name string, action map[string]interface{}) ActionState {
	return ActionState{
		UserId:       int64(data.From.ID),
		ChatId:       data.Message.Chat.ID,
		Conversation: conversation,
		Action:       name,
		ActionData:   action,
	}
}

func CreateActionCache(ctx context.Context) *ActionCache {
	client := ctx.Value("redis").(*redis.Client)

	return &ActionCache{
		redis: client,
	}
}

func (s *ActionCache) SetUserState(ctx context.Context, chatID int64, state ActionState) error {
	key := getActionStateKey(chatID)

	marshal, err := json.Marshal(state)

	if err != nil {
		return err
	}

	return s.redis.Set(ctx, key, marshal, 0).Err()
}

func (s *ActionCache) GetActionState(ctx context.Context, chatID int64) (ActionState, error) {
	key := getActionStateKey(chatID)
	stateStr, err := s.redis.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return ActionState{}, nil
	}

	var state ActionState
	err = json.Unmarshal([]byte(stateStr), &state)
	if err != nil {
		return ActionState{}, err
	}

	return state, nil
}

func (s *ActionCache) ClearActionState(ctx context.Context, chatID int64) error {
	key := getActionStateKey(chatID)
	return s.redis.Del(ctx, key).Err()
}

func getActionStateKey(chatID int64) string {
	return "action_state:" + strconv.FormatInt(chatID, 10)
}
