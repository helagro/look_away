package notifications

import (
	"os/exec"

	"github.com/helagro/look_away/internal/config"
	"github.com/helagro/look_away/internal/utils"

	"github.com/gen2brain/beeep"
)

const (
	notificationHeader = "Look away"
	breakEndMessage    = "That's enough, go back to work!"
)

type Notifier struct {
	config config.NotificationConfig
}

func NewNotifier(cfg config.NotificationConfig) *Notifier {
	return &Notifier{config: cfg}
}

func (n *Notifier) Notify(i int) {
	switch i {
	case 0:
		n.NotifyStart()
	case 1:
		n.NotifyEnd()
	default:
		beeep.Alert(notificationHeader, "Invalid notification type", "")
	}

}

func (n *Notifier) NotifyStart() {
	message := getCommandMessage(n.config.TextCommand)
	utils.Log(message)

	if n.config.UseAlert {
		beeep.Alert(notificationHeader, message, "")
	} else {
		beeep.Notify(notificationHeader, message, "")
	}
}

func getCommandMessage(textCommand string) string {
	if textCommand != "" {
		cmd := exec.Command(textCommand)

		output, err := cmd.Output()
		if err == nil {
			return string(output)
		} else {
			return "Error executing command: " + err.Error()
		}
	} else {
		return "Time to rest your eyes! Look at least 6m away for at least 20 seconds!"
	}
}

func (n *Notifier) NotifyEnd() {
	beeep.Notify(notificationHeader, breakEndMessage, "")
	utils.Log(breakEndMessage)

	if n.config.UseAlert {
		for i := 0; i < 4; i++ {
			beeep.Beep(400, 100)
		}
	}
}
