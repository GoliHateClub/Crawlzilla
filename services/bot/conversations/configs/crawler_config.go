package configs

import (
	"Crawlzilla/services/cache"
	crawlerConfigService "Crawlzilla/services/crawler"
	"context"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ConfigCrawlerConversation(ctx context.Context, state cache.UserState, update tgbotapi.Update) {
	bot := ctx.Value("bot").(*tgbotapi.BotAPI)
	userStates := ctx.Value("user_state").(*cache.UserCache)
	actionStates := ctx.Value("action_state").(*cache.ActionCache)

	switch state.Stage {
	case "init":
		// Ask the user for the first parameter
		msg := tgbotapi.NewMessage(state.ChatId, "لطفاً زمان کرال (crawl time) را بر حسب دقیقه وارد کنید:")
		bot.Send(msg)

		// Update user state
		userStates.SetUserCache(ctx, state.ChatId, cache.UserState{
			ChatId:       state.ChatId,
			UserId:       state.UserId,
			Conversation: "config_crawler",
			Stage:        "get_crawl_time",
		})

	case "get_crawl_time":
		// Parse crawl time
		crawlTime, err := strconv.Atoi(strings.TrimSpace(update.Message.Text))
		if err != nil || crawlTime <= 0 {
			bot.Send(tgbotapi.NewMessage(state.ChatId, "زمان وارد شده نامعتبر است. لطفاً عددی معتبر وارد کنید."))
			return
		}

		// Update action state with crawl time
		currentAction, _ := actionStates.GetActionState(ctx, state.ChatId)
		actionData := currentAction.ActionData
		if actionData == nil {
			actionData = make(map[string]interface{})
		}
		actionData["crawlTime"] = crawlTime

		actionStates.SetUserState(ctx, state.ChatId, cache.ActionState{
			ChatId:       state.ChatId,
			UserId:       state.UserId,
			Conversation: "config_crawler",
			Action:       "set_crawler_config",
			ActionData:   actionData,
		})

		// Ask for the next parameter
		bot.Send(tgbotapi.NewMessage(state.ChatId, "لطفاً زمان اسکرپ صفحات (page scrap time) را بر حسب ثانیه وارد کنید:"))

		// Update user state
		userStates.UpdateUserCache(ctx, state.ChatId, map[string]interface{}{
			"Stage": "get_page_scrap_time",
		})

	case "get_page_scrap_time":
		// Parse page scrap time
		pageScrapTime, err := strconv.Atoi(strings.TrimSpace(update.Message.Text))
		if err != nil || pageScrapTime <= 0 {
			bot.Send(tgbotapi.NewMessage(state.ChatId, "زمان وارد شده نامعتبر است. لطفاً عددی معتبر وارد کنید."))
			return
		}

		// Update action state with page scrap time
		currentAction, _ := actionStates.GetActionState(ctx, state.ChatId)
		actionData := currentAction.ActionData
		actionData["pageScrapTime"] = pageScrapTime

		actionStates.SetUserState(ctx, state.ChatId, cache.ActionState{
			ChatId:       state.ChatId,
			UserId:       state.UserId,
			Conversation: "config_crawler",
			Action:       "set_crawler_config",
			ActionData:   actionData,
		})

		// Ask for the next parameter
		bot.Send(tgbotapi.NewMessage(state.ChatId, "لطفاً حداکثر تعداد آگهی‌ها (ad count) را وارد کنید:"))

		// Update user state
		userStates.UpdateUserCache(ctx, state.ChatId, map[string]interface{}{
			"Stage": "get_ad_count",
		})

	case "get_ad_count":
		// Parse ad count
		adCount, err := strconv.Atoi(strings.TrimSpace(update.Message.Text))
		if err != nil || adCount <= 0 {
			bot.Send(tgbotapi.NewMessage(state.ChatId, "عدد وارد شده نامعتبر است. لطفاً عددی معتبر وارد کنید."))
			return
		}

		// Update action state with ad count
		currentAction, _ := actionStates.GetActionState(ctx, state.ChatId)
		actionData := currentAction.ActionData
		actionData["adCount"] = adCount

		actionStates.SetUserState(ctx, state.ChatId, cache.ActionState{
			ChatId:       state.ChatId,
			UserId:       state.UserId,
			Conversation: "config_crawler",
			Action:       "set_crawler_config",
			ActionData:   actionData,
		})

		// Ask for the next parameter
		bot.Send(tgbotapi.NewMessage(state.ChatId, "لطفاً حداکثر تعداد اسکرول‌ها (max scroll) را وارد کنید:"))

		// Update user state
		userStates.UpdateUserCache(ctx, state.ChatId, map[string]interface{}{
			"Stage": "get_max_scroll",
		})

	case "get_max_scroll":
		// Parse max scroll
		maxScroll, err := strconv.Atoi(strings.TrimSpace(update.Message.Text))
		if err != nil || maxScroll <= 0 {
			bot.Send(tgbotapi.NewMessage(state.ChatId, "عدد وارد شده نامعتبر است. لطفاً عددی معتبر وارد کنید."))
			return
		}

		// Retrieve all action data
		currentAction, _ := actionStates.GetActionState(ctx, state.ChatId)
		actionData := currentAction.ActionData

		// Safely convert values from actionData
		crawlTime, ok := actionData["crawlTime"].(float64)
		if !ok {
			bot.Send(tgbotapi.NewMessage(state.ChatId, "خطا در بازیابی زمان کرال!"))
			return
		}
		pageScrapTime, ok := actionData["pageScrapTime"].(float64)
		if !ok {
			bot.Send(tgbotapi.NewMessage(state.ChatId, "خطا در بازیابی زمان اسکرپ!"))
			return
		}
		adCount, ok := actionData["adCount"].(float64)
		if !ok {
			bot.Send(tgbotapi.NewMessage(state.ChatId, "خطا در بازیابی تعداد آگهی‌ها!"))
			return
		}

		// Call the crawler config service
		crawlerConfigService.SetCrawlerConfig(
			int(crawlTime),     // Convert float64 to int
			int(pageScrapTime), // Convert float64 to int
			int(adCount),       // Convert float64 to int
			maxScroll,          // Already an int
		)

		// Clear cache
		userStates.ClearUserCache(ctx, state.ChatId)
		actionStates.ClearActionState(ctx, state.ChatId)

		// Inform user of success
		bot.Send(tgbotapi.NewMessage(state.ChatId, "تنظیمات کرالر با موفقیت اعمال شد!"))
	}
}
