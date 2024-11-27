package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"log"
	"strconv"
)

type UserCache struct {
	redis *redis.Client
}

type UserState struct {
	Conversation string
	Stage        string
	UserId       int64
	ChatId       int64
}

func CreateNewUserState(conversation string, data *tgbotapi.CallbackQuery) UserState {
	return UserState{
		UserId:       data.From.ID,
		ChatId:       data.Message.Chat.ID,
		Conversation: conversation,
		Stage:        "init",
	}
}

func CreateUserCache(ctx context.Context) *UserCache {
	client := ctx.Value("redis").(*redis.Client)

	return &UserCache{
		redis: client,
	}
}

func (s *UserCache) SetUserCache(ctx context.Context, chatID int64, state UserState) error {
	key := getUserStateKey(chatID)

	marshal, err := json.Marshal(state)

	if err != nil {
		return err
	}

	return s.redis.Set(ctx, key, marshal, 0).Err()
}

func (s *UserCache) UpdateUserCache(ctx context.Context, chatID int64, updates map[string]interface{}) error {
	// Get the current state from Redis
	currentState, err := s.GetUserCache(ctx, chatID)
	if err != nil {
		return err
	}

	// Apply updates to the current state
	for key, value := range updates {
		switch key {
		case "Conversation":
			if conv, ok := value.(string); ok {
				currentState.Conversation = conv
			}
		case "Stage":
			if stage, ok := value.(string); ok {
				currentState.Stage = stage
			}
		case "UserId":
			if userID, ok := value.(int64); ok {
				currentState.UserId = userID
			}
		case "ChatId":
			if chatID, ok := value.(int64); ok {
				currentState.ChatId = chatID
			}
		default:
			return errors.New("invalid field in updates")
		}
	}

	// Save the updated state back to Redis
	return s.SetUserCache(ctx, chatID, currentState)
}

func (s *UserCache) GetUserCache(ctx context.Context, chatID int64) (UserState, error) {
	key := getUserStateKey(chatID)
	stateStr, err := s.redis.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return UserState{}, nil
	}

	var state UserState
	err = json.Unmarshal([]byte(stateStr), &state)
	if err != nil {
		return UserState{}, err
	}

	return state, nil
}

func (s *UserCache) ClearUserCache(ctx context.Context, chatID int64) error {
	key := getUserStateKey(chatID)
	return s.redis.Del(ctx, key).Err()
}

func HandleUserStateError(logger *zap.Logger, state UserState, err error) {
	logger.Error(
		"Error updating action state",
		zap.Error(err),
		zap.String("user_id", strconv.Itoa(int(state.UserId))),
		zap.String("chat_id", strconv.Itoa(int(state.ChatId))),
		zap.String("action", "add_add_st1"),
	)
	log.Printf("Error updating user state: %v", err)
}

func getUserStateKey(chatID int64) string {
	return "conversation_state:" + strconv.FormatInt(chatID, 10)
}
