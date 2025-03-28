package timer

import (
	"context"
	"fmt"
	"jnsltk/look_away/internal/config"
	"jnsltk/look_away/internal/notifications"
	"os/exec"
	"time"
)

type Timer struct {
	TimerDuration time.Duration
	BreakDuration time.Duration
}

func NewTimer(duration time.Duration, breakDuration time.Duration) *Timer {
	return &Timer{
		duration,
		breakDuration,
	}
}

func (t *Timer) Start(ctx context.Context) {
	durations := []time.Duration{t.TimerDuration, t.BreakDuration}

	for {

		for i, duration := range durations {
			timer := time.NewTimer(duration)
			ticker := time.NewTicker(1 * time.Second)

		innerloop:
			for remaining := duration; remaining > 0; remaining -= 1 * time.Second {
				select {
				case <-ctx.Done():
					timer.Stop()
					ticker.Stop()
				case <-timer.C:
					cfg, err := config.LoadConfig()
					if err != nil {
						fmt.Println("Error loading config:", err)
						return
					}

					notifier := notifications.NewNotifier(cfg.Notifications)
					message := getNotificationMessage(i, cfg.Notifications.TextCommand)

					notifier.Notify(message)
					break innerloop
				case <-ticker.C:
					minutes := int(remaining.Minutes())
					seconds := int(remaining.Seconds()) % 60
					fmt.Printf("\r%02d:%02d remaining", minutes, seconds)
				}
			}
			timer.Stop()
			ticker.Stop()
		}
	}
}

func getNotificationMessage(i int, textCommand string) string {
	switch i {
	case 0:
		if textCommand != "" {
			cmd := exec.Command(textCommand)

			output, err := cmd.Output()
			if err == nil {
				return string(output)
			} else {

				return "Error executing command: " + err.Error()
			}
		} else {
			return "Time to rest your eyes! Look at least 20 ft (~6m) away for at least 20 seconds!"
		}
	case 1:
		return "That's enough, go back to work!"
	default:
		return ""
	}
}
