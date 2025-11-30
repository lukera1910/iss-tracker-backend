package opennotify

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"iss-tracker-backend/internal/domain"
)

type Client struct {
	Http *http.Client
}

func NewClient() *Client {
	return &Client {
		Http: &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *Client) GetISSPosition() (*domain.ISSPosition, error) {
	resp, err := c.Http.Get("http://api.open-notify.org/iss-now.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var payload struct {
		Message string `json:"message"`
		Timestamp int64 `json:"timestamp"`
		ISS struct {
			Latitude string `json:"latitude"`
			Longitude string `json:"longitude"`
		} `json:"iss_position"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	lat, _ := strconv.ParseFloat(payload.ISS.Latitude, 64)
	lon, _ := strconv.ParseFloat(payload.ISS.Longitude, 64)

	return &domain.ISSPosition {
		Latitude: lat,
		Longitude: lon,
		Timestamp: payload.Timestamp,
	}, nil
}