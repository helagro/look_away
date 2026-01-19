package notifications

import (
	"github.com/helagro/look_away/internal/config"
	"github.com/helagro/look_away/internal/utils"

	"github.com/gen2brain/beeep"
)

const (
	notificationHeader = "Blink timer"
	breakEndMessage    = "That's enough, go back to work!"
)

type NotificationType int

const (
	NotificationStart NotificationType = iota
	NotificationEnd
)

type Notifier struct {
	config config.NotificationConfig
}

func NewNotifier(cfg config.NotificationConfig) *Notifier {
	return &Notifier{config: cfg}
}

func (n *Notifier) Notify(i NotificationType) {
	message, err := getNotificationMessage(n.config.TextCommand)
	if err != nil {
		return
	}

	switch i {
	case NotificationStart:
		n.notifyStart(message)
	case NotificationEnd:
		n.notifyEnd()
	default:
		beeep.Alert(notificationHeader, "Invalid notification type", "")
	}
}

/* ================================= PRIVATE ================================ */

func (n *Notifier) notifyStart(message string) {
	if !n.config.ShowNotification {
		utils.Log("Notifications are disabled, skipping notification", false)
		return
	}

	if n.config.UseAlert {
		beeep.Alert(notificationHeader, message, "")
	} else {
		beeep.Notify(notificationHeader, message, "")
	}
}

func (n *Notifier) notifyEnd() {
	if n.config.ShowNotification {
		beeep.Notify(notificationHeader, breakEndMessage, "")
	}

	utils.Log(breakEndMessage, false)

	if n.config.UseAlert {
		for i := 0; i < 4; i++ {
			beeep.Beep(400, 100)
		}
	}
}
