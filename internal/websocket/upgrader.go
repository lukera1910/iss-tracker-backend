package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,

	// permite conexão de qualquer origem (para testes/local)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}