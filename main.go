package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Reminder: `defer` doesn't behave as expected in functions with log.Fatal, os.Exit, etc.
	rootCtx := context.Background()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})))

	if err := work(rootCtx, os.Args[1]); err != nil {
		panic(err)
	}
}

func work(ctx context.Context, configFile string) error {
	ctx, cancelRoot := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancelRoot()

	actions, err := parseConfigFile(configFile)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "Starting appearance watcher...")

	listener := &actionAppearanceListener{
		actions: actions,
	}

	return StartWatcher(ctx, listener)
}
