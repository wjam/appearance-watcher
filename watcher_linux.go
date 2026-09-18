//go:build linux

package main

import (
	"context"
	"log/slog"
	"slices"

	"github.com/godbus/dbus/v5"
)

const (
	dbusSettingsInterface   = "org.freedesktop.portal.Settings"
	dbusSettingsMethod      = dbusSettingsInterface + ".Read"
	dbusSettingsMember      = "SettingChanged"
	dbusPortalObjectPath    = "/org/freedesktop/portal/desktop"
	dbusAppearanceNamespace = "org.freedesktop.appearance"
	dbusColorSchemeKey      = "color-scheme"
)

func handleColourScheme(ctx context.Context, listener AppearanceListener, colourScheme uint32) {
	if colourScheme == 1 {
		listener.AppearanceDark(ctx)
		return
	}

	listener.AppearanceLight(ctx)
}

func StartWatcher(ctx context.Context, listener AppearanceListener) error {
	conn, err := dbus.ConnectSessionBus(dbus.WithContext(ctx))
	if err != nil {
		return err
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			slog.ErrorContext(ctx, "Failed to close session bus", "error", err)
		}
	}()

	obj := conn.Object("org.freedesktop.portal.Desktop", dbusPortalObjectPath)

	var colourScheme uint32
	if err := obj.CallWithContext(ctx, dbusSettingsMethod, 0, dbusAppearanceNamespace, dbusColorSchemeKey).
		Store(&colourScheme); err != nil {
		return err
	}

	handleColourScheme(ctx, listener, colourScheme)

	if err := conn.AddMatchSignalContext(
		ctx,
		dbus.WithMatchInterface(dbusSettingsInterface),
		dbus.WithMatchMember(dbusSettingsMember),
		dbus.WithMatchObjectPath(dbusPortalObjectPath),
	); err != nil {
		return err
	}

	signals := make(chan *dbus.Signal, 10) //nolint:mnd // magic number copied from library example
	conn.Signal(signals)

	previousColourScheme := colourScheme

	doneCh := make(chan struct{}, 1)
	go func() {
		defer close(doneCh)
		defer func() { doneCh <- struct{}{} }()
		for {
			select {
			case <-ctx.Done():
				return
			case e := <-signals:
				if !slices.Contains(e.Body, dbusColorSchemeKey) || !slices.Contains(e.Body, dbusAppearanceNamespace) {
					continue
				}
				colourScheme := e.Body[len(e.Body)-1].(dbus.Variant).Value().(uint32)
				if colourScheme == previousColourScheme {
					// Changing appearance scheme in GNOME isn't an atomic write, meaning we tend to see _two_ signals
					// for each transition.
					continue
				}
				previousColourScheme = colourScheme

				handleColourScheme(ctx, listener, colourScheme)
			}
		}
	}()

	<-doneCh

	return nil
}
