package scheduler

import (
	"log"
	"iss-tracker-backend/internal/adapter/opennotify"
	"iss-tracker-backend/internal/domain"
	"iss-tracker-backend/internal/usecase"
	"iss-tracker-backend/internal/websocket"
	"sync"
	"time"
)

type PositionCache struct {
	mu      sync.RWMutex
	Current *domain.ISSPosition
}

func (c *PositionCache) Get() *domain.ISSPosition {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Current
}

func (c *PositionCache) Set(p *domain.ISSPosition) {
	c.mu.Lock()
	c.Current = p
	c.mu.Unlock()
}

func StartISSPoller(cache *PositionCache) {
	client := opennotify.NewClient()

	go func() {
		for {
			pos, err := client.GetISSPosition()
			if err == nil {
				cache.Set(pos)

				// se websocket estiver conectado, envia atualização
				if websocket.ConnectionClient != nil {
					userLat := websocket.ConnectionClient.UserLat
					userLon := websocket.ConnectionClient.UserLon

					visible := usecase.IsVisible(userLat, userLon, pos)

					websocket.ConnectionClient.SendPosition(
						websocket.OutgoingMessage{
							Type:      "iss_position",
							Lat:       pos.Latitude,
							Lon:       pos.Longitude,
							Visible:   visible,
							Timestamp: pos.Timestamp,
						},
					)
				}
			} else {
				log.Println("Erro ao pegar posição da ISS:", err)
			}
			time.Sleep(5 * time.Second) // coleta a cada 5s
		}
	}()
}
