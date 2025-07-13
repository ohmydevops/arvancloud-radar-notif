package notification

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

func NewNotofiersManager(notifiers []Notifier) *NotifiersManager {
	return &NotifiersManager{
		Notifiers: notifiers,
	}
}

func (m *NotifiersManager) Notify(title, message string) error {
	for _, n := range m.Notifiers {
		if err := n.Notify(title, message); err != nil {
			return err
		}
	}
	return nil
}

// ConsoleNotifier prints logs on console.
type ConsoleNotifier struct {
}

func NewConsoleNotifier() *ConsoleNotifier {
	return &ConsoleNotifier{}
}
func (c *ConsoleNotifier) Notify(title, message string) error {
	_, err := fmt.Printf("[%s] %s", title, message)
	return err
}

// DesktopNotifier sends desktop notifications using beeep.
type DesktopNotifier struct {
	IconPath          string
	NotificationTitle string
}

func NewDesktopNotofier(NotificationTitle, iconPath string) *DesktopNotifier {
	return &DesktopNotifier{
		IconPath:          iconPath,
		NotificationTitle: NotificationTitle,
	}
}

func (d *DesktopNotifier) Notify(title, message string) error {
	if beeep.AppName != d.NotificationTitle {
		beeep.AppName = d.NotificationTitle
	}
	if err := beeep.Notify(title, message, d.IconPath); err != nil {
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
