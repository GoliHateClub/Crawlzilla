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

func (s *ActionCache) UpdateActionCache(ctx context.Context, chatID int64, updates map[string]interface{}) error {
	// Get the current state from Redis
	currentState, err := s.GetActionState(ctx, chatID)
	if err != nil {
		return err
	}

	// Apply updates to the current state
	for key, value := range updates {
		switch key {
		case "Conversation":
			if conversation, ok := value.(string); ok {
				currentState.Conversation = conversation
			}
		case "Action":
			if action, ok := value.(string); ok {
				currentState.Action = action
			}
		case "UserId":
			if userID, ok := value.(int64); ok {
				currentState.UserId = userID
			}
		case "ChatId":
			if chatID, ok := value.(int64); ok {
				currentState.ChatId = chatID
			}
		case "ActionData":
			if actionData, ok := value.(map[string]interface{}); ok {
				currentState.ActionData = actionData
			}
		default:
			return errors.New("invalid field in updates")
		}
	}

	// Save the updated state back to Redis
	return s.SetUserState(ctx, chatID, currentState)
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

func HandleActionStateError(logger *zap.Logger, state UserState, err error) {
	logger.Error(
		"Error updating user state",
		zap.Error(err),
		zap.String("user_id", strconv.Itoa(int(state.UserId))),
		zap.String("chat_id", strconv.Itoa(int(state.ChatId))),
	)
	log.Printf("Error updating user state: %v", err)
}

func getActionStateKey(chatID int64) string {
	return "action_state:" + strconv.FormatInt(chatID, 10)
}
