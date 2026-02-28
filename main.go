package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"outages-bot/internal/application/notification"
	"outages-bot/internal/application/subscription"
	"outages-bot/internal/client/outageapi"
	tgclient "outages-bot/internal/client/telegram"
	"outages-bot/internal/cmd"
	"outages-bot/internal/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "outages-bot",
		Short: "Telegram bot for power outage notifications in Lviv, Ukraine",
	}

	rootCmd.AddCommand(botCmd())
	rootCmd.AddCommand(notifierCmd())
	rootCmd.AddCommand(outagesCmd())
	rootCmd.AddCommand(usersCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustBotAPI() *tgbotapi.BotAPI {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "TELEGRAM_BOT_TOKEN environment variable is required")
		os.Exit(1)
	}
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create Telegram bot: %v\n", err)
		os.Exit(1)
	}
	return api
}

func dataDir() string {
	return getEnv("DATA_DIR", "data")
}

func botCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bot",
		Short: "Run the Telegram bot (long-running)",
		Run: func(cmd *cobra.Command, args []string) {
			api := mustBotAPI()
			dir := dataDir()

			userRepo, err := repository.NewFileUserRepository(filepath.Join(dir, "users"))
			if err != nil {
				log.Fatalf("Failed to create user repository: %v", err)
			}

			streetRepo, err := repository.NewFileStreetRepository(filepath.Join(dir, "streets.csv"))
			if err != nil {
				log.Fatalf("Failed to create street repository: %v", err)
			}

			runner := cli.NewBotRunner(cli.BotRunnerConfig{
				Bot:                     api,
				SearchStreetService:     subscription.NewSearchStreetService(streetRepo),
				ShowSubscriptionService: subscription.NewShowSubscriptionService(userRepo),
				SaveSubscriptionService: subscription.NewSaveSubscriptionService(userRepo),
				UserRepo:                userRepo,
			})
			defer runner.Close()

			log.Printf("Bot started as @%s", api.Self.UserName)
			runner.Run()
		},
	}
}

func notifierCmd() *cobra.Command {
	var interval time.Duration

	cmd := &cobra.Command{
		Use:   "notifier",
		Short: "Fetch outages and send notifications",
		RunE: func(cmd *cobra.Command, args []string) error {
			api := mustBotAPI()
			dir := dataDir()

			userRepo, err := repository.NewFileUserRepository(filepath.Join(dir, "users"))
			if err != nil {
				return fmt.Errorf("failed to create user repository: %w", err)
			}

			outageProvider := outageapi.NewProvider("", nil, nil)
			sender := tgclient.NewNotificationSender(api)
			fetchService := notification.NewOutageFetchService(outageProvider)
			notificationService := notification.NewService(sender, userRepo, log.Default())

			runFn := func(ctx context.Context) error {
				return cli.RunNotifierCommand(ctx, fetchService, notificationService, log.Default())
			}

			if interval <= 0 {
				return runFn(context.Background())
			}

			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			return runNotifierLoop(ctx, runFn, interval, log.Default())
		},
	}

	cmd.Flags().DurationVar(&interval, "interval", 0, "Run repeatedly with this interval between runs (e.g. 60s). If zero, run once and exit.")

	return cmd
}

func runNotifierLoop(ctx context.Context, runFn func(context.Context) error, interval time.Duration, logger *log.Logger) error {
	if interval <= 0 {
		return fmt.Errorf("interval must be positive, got %v", interval)
	}

	for {
		if ctx.Err() != nil {
			return nil
		}
		if err := runFn(ctx); err != nil {
			logger.Printf("Notifier error: %v", err)
		}
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(interval):
		}
	}
}

func outagesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "outages",
		Short: "Print a table of current outages",
		RunE: func(cmd *cobra.Command, args []string) error {
			outageProvider := outageapi.NewProvider("", nil, nil)
			return cli.RunOutagesCommand(context.Background(), outageProvider, os.Stdout)
		},
	}
}

func usersCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "users",
		Short: "List all subscribed users",
		RunE: func(cmd *cobra.Command, args []string) error {
			api := mustBotAPI()
			dir := dataDir()

			userRepo, err := repository.NewFileUserRepository(filepath.Join(dir, "users"))
			if err != nil {
				return fmt.Errorf("failed to create user repository: %w", err)
			}

			userInfoProvider := tgclient.NewUserInfoProvider(api)
			cli.RunUsersCommand(userRepo, userInfoProvider, os.Stdout, log.Default())
			return nil
		},
	}
}
