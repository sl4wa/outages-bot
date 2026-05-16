package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/sl4wa/outages-bot/internal/schedule/loe"
	"github.com/sl4wa/outages-bot/internal/schedule/notifier"
	"github.com/sl4wa/outages-bot/internal/schedule/persistence"
	"github.com/sl4wa/outages-bot/internal/schedule/schedule"
	"github.com/sl4wa/outages-bot/internal/schedule/telegram"
	"github.com/sl4wa/outages-bot/internal/shared/subscribers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type appConfig struct {
	StatePath        string
	HTTPCachePath    string
	APIURL           string
	TelegramBotToken string
	TelegramUsersDir string
	Zone             *time.Location
}

func main() {
	var interval time.Duration
	flag.DurationVar(&interval, "interval", 0, "Run repeatedly with this interval between runs (e.g. 60s). If zero, run once and exit.")
	flag.Parse()

	log.SetOutput(os.Stdout)
	_ = godotenv.Load()

	config, err := configFromEnv()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bot, err := tgbotapi.NewBotAPI(config.TelegramBotToken)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create Telegram bot: %v\n", err)
		os.Exit(1)
	}

	cache := loe.NewHTTPCache(config.HTTPCachePath)
	provider := loe.Provider{
		LoadPayload: func(ctx context.Context) (string, error) {
			return cache.Fetch(ctx, config.APIURL)
		},
		OnPayloadAccepted: cache.Commit,
	}
	runner := notifier.Runner{
		Provider: provider,
		Store:    persistence.NewCSVStateStore(config.StatePath),
		Notifier: telegram.UserNotifier{
			Sender:      telegram.BotSender{Bot: bot},
			Subscribers: subscribers.NewFileStore(config.TelegramUsersDir),
			Logger:      log.Default(),
		},
		Zone: config.Zone,
	}

	runFn := func(ctx context.Context) error {
		return runScheduleOnce(ctx, runner, config)
	}

	if interval <= 0 {
		if err := runFn(context.Background()); err != nil {
			fmt.Fprintf(os.Stderr, "schedule notification failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := runScheduleLoop(ctx, runFn, interval, log.Default()); err != nil {
		fmt.Fprintf(os.Stderr, "schedule notification failed: %v\n", err)
		os.Exit(1)
	}
}

func runScheduleOnce(ctx context.Context, runner notifier.Runner, config appConfig) error {
	now := time.Now().In(config.Zone)
	today := schedule.NormalizeDate(now)
	result, err := runner.Execute(ctx, []time.Time{today, today.AddDate(0, 0, 1)})
	if err != nil {
		return err
	}
	if result == nil {
		log.Printf("No schedule found for %s and %s. State preserved.", schedule.FormatStateDate(today), schedule.FormatStateDate(today.AddDate(0, 0, 1)))
		return nil
	}
	if result.Changed {
		log.Printf("Updated schedule state for %d changed date(s) at %s", len(result.ChangedDates), config.StatePath)
	} else {
		log.Printf("Schedule unchanged. Skipping state write.")
	}
	if result.Notified {
		log.Printf("Telegram user broadcast sent.")
	} else if len(result.ChangedDates) > 0 {
		log.Printf("Schedule changed, but no Telegram users received the message.")
	}
	return nil
}

func runScheduleLoop(ctx context.Context, runFn func(context.Context) error, interval time.Duration, logger *log.Logger) error {
	if interval <= 0 {
		return fmt.Errorf("interval must be positive, got %v", interval)
	}

	for {
		if ctx.Err() != nil {
			return nil
		}
		if err := runFn(ctx); err != nil {
			logger.Printf("Schedule notification error: %v", err)
		}
		timer := time.NewTimer(interval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil
		case <-timer.C:
		}
	}
}

func configFromEnv() (appConfig, error) {
	zone, err := time.LoadLocation("Europe/Kyiv")
	if err != nil {
		return appConfig{}, err
	}
	dataDirEnv := os.Getenv("DATA_DIR")
	if dataDirEnv == "" {
		return appConfig{}, fmt.Errorf("DATA_DIR must be set")
	}
	dataDir := absPath(dataDirEnv)
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		return appConfig{}, fmt.Errorf("TELEGRAM_BOT_TOKEN must be set")
	}
	apiURL := os.Getenv("SCHEDULE_API_URL")
	if apiURL == "" {
		return appConfig{}, fmt.Errorf("SCHEDULE_API_URL must be set")
	}
	return appConfig{
		StatePath:        filepath.Join(dataDir, persistence.StateFileName),
		HTTPCachePath:    filepath.Join(dataDir, loe.DefaultCacheFileName),
		APIURL:           apiURL,
		TelegramBotToken: token,
		TelegramUsersDir: filepath.Join(dataDir, "users"),
		Zone:             zone,
	}, nil
}

func absPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return filepath.Clean(abs)
}
