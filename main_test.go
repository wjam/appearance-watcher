package main

import (
	"context"
	"log/slog"
	"sync"
	"testing"
	"time"
)

type mockAction struct {
	mu     sync.Mutex
	Themes []string
}

func (m *mockAction) ToDark() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Themes = append(m.Themes, "dark")
	return nil
}

func (m *mockAction) ToLight() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Themes = append(m.Themes, "light")
	return nil
}

func TestStartWatcher(t *testing.T) {
	mock := &mockAction{}

	ctx, cancel := context.WithTimeout(t.Context(), 500*time.Millisecond)
	defer cancel()

	defaultSlog := slog.Default()
	t.Cleanup(func() {
		slog.SetDefault(defaultSlog)
	})
	slog.SetDefault(slog.New(slog.NewTextHandler(t.Output(), &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})))

	err := StartWatcher(ctx, &actionAppearanceListener{actions: []appearanceChangeAction{mock}})
	if err != nil {
		t.Fatal(err)
	}

	mock.mu.Lock()
	defer mock.mu.Unlock()

	if len(mock.Themes) == 0 {
		t.Error("Expected at least an initial theme capture, but got nothing")
	}
}
