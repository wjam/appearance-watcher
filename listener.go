package main

import (
	"context"
	"errors"
	"log/slog"
)

// AppearanceListener defines the contract for structs wanting theme updates.
type AppearanceListener interface {
	AppearanceDark(context.Context)
	AppearanceLight(context.Context)
}

type appearanceChangeAction interface {
	ToLight() error
	ToDark() error
}

type actionAppearanceListener struct {
	actions []appearanceChangeAction
}

func (l *actionAppearanceListener) AppearanceDark(ctx context.Context) {
	slog.InfoContext(ctx, "Appearance is now dark")
	var errs []error
	for _, a := range l.actions {
		errs = append(errs, a.ToDark())
	}

	if err := errors.Join(errs...); err != nil {
		slog.ErrorContext(ctx, "Error handling dark appearance", "err", err)
	}
}

func (l *actionAppearanceListener) AppearanceLight(ctx context.Context) {
	slog.InfoContext(ctx, "Appearance is now light")
	var errs []error
	for _, a := range l.actions {
		errs = append(errs, a.ToLight())
	}

	if err := errors.Join(errs...); err != nil {
		slog.ErrorContext(ctx, "Error handling light appearance", "err", err)
	}
}
