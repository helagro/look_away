package timer

import (
	"context"
	"fmt"
	"time"

	"github.com/helagro/look_away/internal/config"
	"github.com/helagro/look_away/internal/notifications"
)

type Timer struct {
}

func NewTimer(duration time.Duration, breakDuration time.Duration) *Timer {
	return &Timer{}
}

func (t *Timer) Start(ctx context.Context) {

	for {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		durations := []time.Duration{cfg.GetTimerDuration(), cfg.GetBreakSeconds()}

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
					notifier := notifications.NewNotifier(cfg.Notifications)

					if i == 0 {
						notifier.Notify(notifications.NotificationStart)
					} else {
						notifier.Notify(notifications.NotificationEnd)
					}

					break innerloop
				case <-ticker.C:
					minutes := int(remaining.Minutes())

					if minutes%10 == 0 {
						seconds := int(remaining.Seconds()) % 60

						if seconds == 0 {
							fmt.Printf("\r%02d:%02d remaining\n", minutes, seconds)
						}
					}
				}
			}
			timer.Stop()
			ticker.Stop()
		}
	}
}
