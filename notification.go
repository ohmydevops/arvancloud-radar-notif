package main

import (
	"fmt"

	"github.com/gen2brain/beeep"
)

// Notifier is an interface for different notification backends.
// Only cares about getting the final title + message.
type Notifier interface {
	Notify(title string, message string) error
}

// NotifiersManager combines multiple notifiers into one.
type NotifiersManager struct {
	Notifiers []Notifier
}

func (m *NotifiersManager) Notify(title, message string) error {
	for _, n := range m.Notifiers {
		if err := n.Notify(title, message); err != nil {
			return err
		}
	}
	return nil
}

// DesktopNotifier sends desktop notifications using beeep.
type DesktopNotifier struct{}

func (d *DesktopNotifier) Notify(title, message string) error {
	if err := beeep.Notify(title, message, "./icon.png"); err != nil {
		return fmt.Errorf("desktop notification error: %v", err)
	}
	return nil
}

// Future backends:
//
// type EmailNotifier struct{}
// func (e *EmailNotifier) Notify(title, message string) error { ... }
//
// type TelegramNotifier struct{}
// func (s *SlackNotifier) Notify(title, message string) error { ... }
