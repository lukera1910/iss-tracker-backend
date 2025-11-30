package httpadapter

import (
	"encoding/json"
	"net/http"
	"strconv"

	"iss-tracker-backend/internal/scheduler"
	"iss-tracker-backend/internal/usecase"
)

type Handler struct {
	Cache *scheduler.PositionCache
}

func (h *Handler) GetISSNow(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(h.Cache.Current)
}

func (h *Handler) GetVisibility(w http.ResponseWriter, r *http.Request) {
	lat, _ := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	lon, _ := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)

	visible := usecase.IsVisible(lat, lon, h.Cache.Current)

	json.NewEncoder(w).Encode(map[string]any{
		"visible": visible,
		"position": h.Cache.Current,
	})
}