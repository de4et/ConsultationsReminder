package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

var cfg *Config

// TODO: add structure levels?
func main() {
	// reading config
	cfg = ReadConfig("config.yml")

	// setup global logger (os.Stdout for local; log.txt for prod)
	setupLogger(cfg.Env)

	slog.Info("started ConsultationsReminder")
	slog.Debug("->", slog.String("config", fmt.Sprintf("%+v", cfg)))

	// get full schedule from file
	groups := ParseFile(cfg.SchedulePath)

	// setup bot
	bot := CreateBot(cfg.BotKey)

	// setting up reminders for specific group
	notifier := CreateNotifier(withTGBot(bot))
	notifier.SetupLessonsReminders(groups[cfg.GroupNumber])

	// chillin
	slog.Info("Waiting...")
	for {
		time.Sleep(100 * time.Millisecond)
	}
}

func createLogFile(fileName string) *os.File {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	return file
}

func setupLogger(env string) {
	var log *slog.Logger
	const logFilePath = "log.txt"

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewTextHandler(createLogFile(logFilePath), &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	slog.SetDefault(log)
}
