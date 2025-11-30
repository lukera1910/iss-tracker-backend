package alert

import (
	"context"
	"fmt"
	"iss-tracker-backend/internal/notifier"
	"iss-tracker-backend/internal/scheduler"
	"iss-tracker-backend/internal/usecase"
	"time"
)

func RunAlertLoop(
	ctx context.Context,
	userLat float64,
	userLon float64,
	cache *scheduler.PositionCache,
	state *VisibilityState,
	notifier notifier.Notifier,
) {
	go func() {
		for {
			pos := cache.Get()

			if pos != nil {
				visibileNow := usecase.IsVisible(userLat, userLon, pos)

				// detecta transição: invisível -> visível
				if visibileNow && !state.LastVisible {
					notifier.Notify(ctx,
						fmt.Sprintf("A ISS está visível agora! Lat: %.2f Lon: %.2f",
						pos.Latitude, pos.Longitude),
					)	
				}

				state.LastVisible = visibileNow
			}

			time.Sleep(5 * time.Second)
		}
	}()
}