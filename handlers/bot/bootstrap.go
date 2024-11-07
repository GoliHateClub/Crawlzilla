package bot

import (
	"github.com/GoliHateClub/Crawlzilla/config"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/net/proxy"
)

// UserState structure
type UserState struct {
	ChatID          int64
	CurrentCommand  string
	Data            map[string]interface{}
	LastInteraction time.Time
	Filter          filter
}
type filter struct {
	MinPrice     int    `json:"min_price,omitempty"`
	MaxPrice     int    `json:"max_price,omitempty"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	MinArea      int    `json:"min_area"`
	MaxArea      int    `json:"max_area"`
	Category     string `json:"category,omitempty"`
	MinAge       int    `json:"min_age,omitempty"`
	MaxAge       int    `json:"max_age,omitempty"`
	MinFloor     int    `json:"min_floor,omitempty"`
	MaxFloor     int    `json:"max_floor,omitempty"`
	HasElevator  bool   `json:"has_elevator,omitempty"`
	HasStorage   bool   `json:"has_storage,omitempty"`
	MinDate      string `json:"min_date,omitempty"`
	MaxDate      string `json:"max_date,omitempty"`
}
// BotServer struct with a mutex for thread-safe state access
type BotServer struct {
	bot            *tgbotapi.BotAPI
	webhookURL     string
	socksProxyAddr string
	listenAddr     string
	userStates     map[int64]*UserState
	client         *http.Client
	mu             sync.RWMutex
}

// SetMinPrice sets the minimum price filter for the user.
func (bs *BotServer) SetMinPrice(chatID int64, minPrice int) {
	userState := bs.GetUserState(chatID)
	userState.Filter.MinPrice = minPrice
}

// SetMaxPrice sets the maximum price filter for the user.
func (bs *BotServer) SetMaxPrice(chatID int64, maxPrice int) {
	userState := bs.GetUserState(chatID)
	userState.Filter.MaxPrice = maxPrice
}

// SetCity sets the city filter for the user.
func (bs *BotServer) SetCity(chatID int64, city string) {
	userState := bs.GetUserState(chatID)
	userState.Filter.City = city
}

// SetNeighborhood sets the neighborhood filter for the user.
func (bs *BotServer) SetNeighborhood(chatID int64, neighborhood string) {
	userState := bs.GetUserState(chatID)
	userState.Filter.Neighborhood = neighborhood
}

// SetMinArea sets the minimum area filter for the user.
func (bs *BotServer) SetMinArea(chatID int64, minArea int) {
	userState := bs.GetUserState(chatID)
	userState.Filter.MinArea = minArea
}

// SetMaxArea sets the maximum area filter for the user.
func (bs *BotServer) SetMaxArea(chatID int64, maxArea int) {
	userState := bs.GetUserState(chatID)
	userState.Filter.MaxArea = maxArea
}

// SetCategory sets the category filter for the user.
func (bs *BotServer) SetCategory(chatID int64, category string) {
	userState := bs.GetUserState(chatID)
	userState.Filter.Category = category
}

// SetMinAge sets the minimum age filter for the user.
func (bs *BotServer) SetMinAge(chatID int64, minAge int) {
	userState := bs.GetUserState(chatID)
	userState.Filter.MinAge = minAge
}

// SetMaxAge sets the maximum age filter for the user.
func (bs *BotServer) SetMaxAge(chatID int64, maxAge int) {
	userState := bs.GetUserState(chatID)
	userState.Filter.MaxAge = maxAge
}

// SetMinFloor sets the minimum floor filter for the user.
func (bs *BotServer) SetMinFloor(chatID int64, minFloor int) {
	userState := bs.GetUserState(chatID)
	userState.Filter.MinFloor = minFloor
}

// SetMaxFloor sets the maximum floor filter for the user.
func (bs *BotServer) SetMaxFloor(chatID int64, maxFloor int) {
	userState := bs.GetUserState(chatID)
	userState.Filter.MaxFloor = maxFloor
}

// SetHasElevator sets the has elevator filter for the user.
func (bs *BotServer) SetHasElevator(chatID int64, hasElevator bool) {
	userState := bs.GetUserState(chatID)
	userState.Filter.HasElevator = hasElevator
}

// SetHasStorage sets the has storage filter for the user.
func (bs *BotServer) SetHasStorage(chatID int64, hasStorage bool) {
	userState := bs.GetUserState(chatID)
	userState.Filter.HasStorage = hasStorage
}

// SetMinDate sets the minimum date filter for the user.
func (bs *BotServer) SetMinDate(chatID int64, minDate string) {
	userState := bs.GetUserState(chatID)
	userState.Filter.MinDate = minDate
}

// SetMaxDate sets the maximum date filter for the user.
func (bs *BotServer) SetMaxDate(chatID int64, maxDate string) {
	userState := bs.GetUserState(chatID)
	userState.Filter.MaxDate = maxDate
}

// GetMinPrice retrieves the minimum price filter for the user.
func (bs *BotServer) GetMinPrice(chatID int64) int {
	userState := bs.GetUserState(chatID)
	return userState.Filter.MinPrice
}

// GetMaxPrice retrieves the maximum price filter for the user.
func (bs *BotServer) GetMaxPrice(chatID int64) int {
	userState := bs.GetUserState(chatID)
	return userState.Filter.MaxPrice
}

// GetCity retrieves the city filter for the user.
func (bs *BotServer) GetCity(chatID int64) string {
	userState := bs.GetUserState(chatID)
	return userState.Filter.City
}

// GetNeighborhood retrieves the neighborhood filter for the user.
func (bs *BotServer) GetNeighborhood(chatID int64) string {
	userState := bs.GetUserState(chatID)
	return userState.Filter.Neighborhood
}

// GetMinArea retrieves the minimum area filter for the user.
func (bs *BotServer) GetMinArea(chatID int64) int {
	userState := bs.GetUserState(chatID)
	return userState.Filter.MinArea
}

// GetMaxArea retrieves the maximum area filter for the user.
func (bs *BotServer) GetMaxArea(chatID int64) int {
	userState := bs.GetUserState(chatID)
	return userState.Filter.MaxArea
}

// GetCategory retrieves the category filter for the user.
func (bs *BotServer) GetCategory(chatID int64) string {
	userState := bs.GetUserState(chatID)
	return userState.Filter.Category
}

// GetMinAge retrieves the minimum age filter for the user.
func (bs *BotServer) GetMinAge(chatID int64) int {
	userState := bs.GetUserState(chatID)
	return userState.Filter.MinAge
}

// GetMaxAge retrieves the maximum age filter for the user.
func (bs *BotServer) GetMaxAge(chatID int64) int {
	userState := bs.GetUserState(chatID)
	return userState.Filter.MaxAge
}

// GetMinFloor retrieves the minimum floor filter for the user.
func (bs *BotServer) GetMinFloor(chatID int64) int {
	userState := bs.GetUserState(chatID)
	return userState.Filter.MinFloor
}

// GetMaxFloor retrieves the maximum floor filter for the user.
func (bs *BotServer) GetMaxFloor(chatID int64) int {
	userState := bs.GetUserState(chatID)
	return userState.Filter.MaxFloor
}

// GetHasElevator retrieves the has elevator filter for the user.
func (bs *BotServer) GetHasElevator(chatID int64) bool {
	userState := bs.GetUserState(chatID)
	return userState.Filter.HasElevator
}

// GetHasStorage retrieves the has storage filter for the user.
func (bs *BotServer) GetHasStorage(chatID int64) bool {
	userState := bs.GetUserState(chatID)
	return userState.Filter.HasStorage
}

// GetMinDate retrieves the minimum date filter for the user.
func (bs *BotServer) GetMinDate(chatID int64) string {
	userState := bs.GetUserState(chatID)
	return userState.Filter.MinDate
}

// GetMaxDate retrieves the maximum date filter for the user.
func (bs *BotServer) GetMaxDate(chatID int64) string {
	userState := bs.GetUserState(chatID)
	return userState.Filter.MaxDate
}

func (bs *BotServer) GetUserState(chatID int64) *UserState {
	bs.mu.RLock()
	state, exists := bs.userStates[chatID]
	bs.mu.RUnlock()

	if !exists {
		// Initialize new state if not found
		state = &UserState{
			ChatID:          chatID,
			CurrentCommand:  "",
			Data:            make(map[string]interface{}),
			LastInteraction: time.Now(),
		}
		bs.mu.Lock()
		bs.userStates[chatID] = state
		bs.mu.Unlock()
	}
	return state
}

// SetUserState updates the user's current command.
func (bs *BotServer) SetUserState(chatID int64, command string) {
	bs.mu.Lock()
	if userState, exists := bs.userStates[chatID]; exists {
		userState.CurrentCommand = command
		userState.LastInteraction = time.Now()
	} else {
		bs.userStates[chatID] = &UserState{
			ChatID:          chatID,
			CurrentCommand:  command,
			Data:            make(map[string]interface{}),
			LastInteraction: time.Now(),
		}
	}
	bs.mu.Unlock()
}
func (bs *BotServer) ClearUserState(chatID int64) {
	bs.mu.Lock()
	delete(bs.userStates, chatID)
	bs.mu.Unlock()
}

// NewBotServer initializes a new BotServer instance.
func NewBotServer(botToken, webhookURL, socksProxyAddr, listenAddr string) (*BotServer, error) {
	// Create SOCKS5 dialer
	dialer, err := proxy.SOCKS5("tcp", socksProxyAddr, nil, proxy.Direct)
	if err != nil {
		return nil, err
	}

	// Configure HTTP client with the SOCKS5 dialer
	httpClient := &http.Client{
		Transport: &http.Transport{
			Dial: dialer.Dial,
		},
		Timeout: 10 * time.Second,
	}

	bot, err := tgbotapi.NewBotAPIWithClient(botToken, tgbotapi.APIEndpoint, httpClient)
	if err != nil {
		return nil, err
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &BotServer{
		bot:            bot,
		webhookURL:     webhookURL,
		socksProxyAddr: socksProxyAddr,
		listenAddr:     listenAddr,
		userStates:     make(map[int64]*UserState),
		client:         httpClient,
	}, nil
}

func (bs *BotServer) Start() error {
	wh, _ := tgbotapi.NewWebhook(bs.webhookURL + bs.bot.Token)
	_, err := bs.bot.Request(wh)
	if err != nil {
		return err
	}

	info, err := bs.bot.GetWebhookInfo()
	if err != nil {
		return err
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	err = bs.configBot()
	if err != nil {
		log.Printf("command config failed")
	}
	updates := bs.bot.ListenForWebhook("/" + bs.bot.Token)
	go func() {
		log.Printf("Starting server on %s", bs.listenAddr)
		if err := http.ListenAndServe(bs.listenAddr, nil); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	bs.handleUpdates(updates)
	return nil
}

func (bs *BotServer) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				bs.CommandHandler(update.Message)
			} else {
				bs.MessageHandler(update.Message)
			}
		} else if update.CallbackQuery != nil {
			bs.HandleCallbackQuery(update.CallbackQuery)
		}
	}
}
func (bs *BotServer) configBot() error {
	cmdCfg := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "start",
			Description: "start command",
		},
		tgbotapi.BotCommand{
			Command:     "echo",
			Description: "echo command",
		},
		tgbotapi.BotCommand{
			Command:     "help",
			Description: "help command",
		},
		tgbotapi.BotCommand{
			Command:     "admin",
			Description: "admin command",
		},
		tgbotapi.BotCommand{
			Command:     "set_ilter",
			Description: "setFilter command",
		},
		tgbotapi.BotCommand{
			Command:     "test",
			Description: "test command",
		},
		tgbotapi.BotCommand{
			Command:     "media",
			Description: "media command",
		},
		tgbotapi.BotCommand{
			Command:     "extension",
			Description: "extension command",
		},
	)
	_, err := bs.bot.Send(cmdCfg)
	return err
}

func RunBot() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	botServer, err := NewBotServer(cfg.BotToken, cfg.WebhookURL, cfg.SocksProxyAddr, cfg.ListenAddr)
	if err != nil {
		log.Fatalf("Failed to initialize bot server: %v", err)
	}

	if err := botServer.Start(); err != nil {
		log.Fatalf("Bot server error: %v", err)
	}
}

func (bs *BotServer) String(id int64) string {
	var sb strings.Builder
	if bs.GetMinPrice(id) > 0 {
		sb.WriteString(fmt.Sprintf("💰 حداقل قیمت: %d\n", bs.GetMinPrice(id)))
	}
	if bs.GetMaxPrice(id) > 0 {
		sb.WriteString(fmt.Sprintf("💰 حداکثر قیمت: %d\n", bs.GetMaxPrice(id)))
	}
	if bs.GetCity(id) != "" {
		sb.WriteString(fmt.Sprintf("🏙️ شهر: %s\n", bs.GetCity(id)))
	}
	if bs.GetNeighborhood(id) != "" {
		sb.WriteString(fmt.Sprintf("🏘️ محله: %s\n", bs.GetNeighborhood(id)))
	}
	if bs.GetMinArea(id) > 0 {
		sb.WriteString(fmt.Sprintf("📐 حداقل مساحت: %d متر مربع\n", bs.GetMinArea(id)))
	}
	if bs.GetMaxArea(id) > 0 {
		sb.WriteString(fmt.Sprintf("📐 حداکثر مساحت: %d متر مربع\n", bs.GetMaxArea(id)))
	}
	if bs.GetCategory(id) != "" {
		sb.WriteString(fmt.Sprintf("🏠 دسته‌بندی: %s\n", bs.GetCategory(id)))
	}
	if bs.GetMinAge(id) > 0 {
		sb.WriteString(fmt.Sprintf("⏳ حداقل سن: %d سال\n", bs.GetMinAge(id)))
	}
	if bs.GetMaxAge(id) > 0 {
		sb.WriteString(fmt.Sprintf("⏳ حداکثر سن: %d سال\n", bs.GetMaxAge(id)))
	}
	if bs.GetMinFloor(id) > 0 {
		sb.WriteString(fmt.Sprintf("🏢 حداقل طبقه: %d\n", bs.GetMinFloor(id)))
	}
	if bs.GetMaxFloor(id) > 0 {
		sb.WriteString(fmt.Sprintf("🏢 حداکثر طبقه: %d\n", bs.GetMaxFloor(id)))
	}
	if bs.GetHasElevator(id) {
		sb.WriteString("🛗 آسانسور: دارد\n")
	}
	if bs.GetHasStorage(id) {
		sb.WriteString("📦 انباری: دارد\n")
	}
	if bs.GetMinDate(id) != "" {
		sb.WriteString(fmt.Sprintf("📅 از تاریخ: %s\n", bs.GetMinDate(id)))
	}
	if bs.GetMaxDate(id) != "" {
		sb.WriteString(fmt.Sprintf("📅 تا تاریخ: %s\n", bs.GetMaxDate(id)))
	}

	return sb.String()
}
