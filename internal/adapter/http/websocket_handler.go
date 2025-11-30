package httpadapter

import (
	"log"
	"net/http"

	wss "iss-tracker-backend/internal/websocket"
)

func WebSockerHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wss.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WS upgrade error:", err)
		return
	}

	client := wss.NewClient(conn)
	// armazena o cliente exportado no pacote websocket (evita dependência circular)
	wss.ConnectionClient = client

	go func() {
		client.ReadLoop()
		client.Close()
		wss.ConnectionClient = nil
	}()

	log.Println("Cliente conectado!")

	go client.ReadLoop()

	// aqui **não bloqueamos**, o write é feito pelo poller de iss
}
