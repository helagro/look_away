package timer

import (
	"context"
	"fmt"
	"time"

	"github.com/helagro/look_away/internal/config"
	"github.com/helagro/look_away/internal/notifications"
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

	for {
		durations := []time.Duration{t.TimerDuration, t.BreakDuration}

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

					t.TimerDuration = time.Duration(cfg.Timer.DurationMinutes)
					t.BreakDuration = time.Duration(cfg.Timer.BreakSeconds)
					notifier := notifications.NewNotifier(cfg.Notifications)
					notifier.Notify(i)

					break innerloop
				case <-ticker.C:
					minutes := int(remaining.Minutes())
					seconds := int(remaining.Seconds()) % 60

					if minutes%10 == 0 && seconds == 0 {
						fmt.Printf("\r%02d:%02d remaining", minutes, seconds)
					}
				}
			}
			timer.Stop()
			ticker.Stop()
		}
	}
}
