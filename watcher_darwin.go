//go:build darwin

package main

/*
#cgo LDFLAGS: -framework Foundation
const char* GetCurrentAppearance();
*/
import "C"
import (
	"context"
	"time"
)

func handleTheme(ctx context.Context, listener AppearanceListener, theme string) {
	if theme == "Dark" {
		listener.AppearanceDark(ctx)
		return
	}
	listener.AppearanceLight(ctx)
}

func StartWatcher(ctx context.Context, listener AppearanceListener) error {
	previousTheme := C.GoString(C.GetCurrentAppearance())
	handleTheme(ctx, listener, previousTheme)

	tick := time.NewTicker(1 * time.Second)

	// While it is possible to watch for changes to the desktop appearance, limitations with the Apple AppKit means this
	// way is simpler and testable. Apple AppKit requires all work to be done on the main thread while Go happily swaps
	// threads on a whim and Go tests are never done on the main thread.

	for {
		select {
		case <-tick.C:
			theme := C.GoString(C.GetCurrentAppearance())
			if previousTheme == theme {
				continue
			}
			handleTheme(ctx, listener, theme)
			previousTheme = theme
		case <-ctx.Done():
			return nil
		}
	}
}
