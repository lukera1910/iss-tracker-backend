package main

import (
	"context"
	"log"
	"net/http"

	httpadapter "iss-tracker-backend/internal/adapter/http"
	"iss-tracker-backend/internal/alert"
	"iss-tracker-backend/internal/notifier"
	"iss-tracker-backend/internal/scheduler"
)

func main() {
	cache := &scheduler.PositionCache{}
	scheduler.StartISSPoller(cache)

	ctx := context.Background()

	userLat := -23.55
	userLon := -46.63

	alertState := &alert.VisibilityState{}
	notifier := notifier.LogNotifier{}

	alert.RunAlertLoop(ctx, userLat, userLon, cache, alertState, notifier)

	handler := &httpadapter.Handler{
		Cache: cache,
	}

	http.HandleFunc("/iss/now", handler.GetISSNow)
	http.HandleFunc("/iss/visible", handler.GetVisibility)
	http.HandleFunc("/ws", httpadapter.WebSockerHandler)

	log.Println("ISS Tracker rodando em http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}