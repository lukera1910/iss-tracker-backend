package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn

	// localização enviada pelo usuário
	UserLat float64
	UserLon float64

	closed  bool

	mu sync.RWMutex
}

func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
	c.conn.Close()
}

// ConnectionClient é a referência global (um cliente) usada por quem precisa
// enviar atualizações para o cliente WebSocket. Exportamos aqui para evitar
// ciclos entre pacotes (por exemplo: scheduler -> adapter/http -> scheduler).
var ConnectionClient *Client

// mensagens recebidas do cliente
type IncomingMessage struct {
	Type string  `json:"type"` // ex: "update_location"
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

// mensagens enviadas para o cliente
type OutgoingMessage struct {
	Type      string  `json:"type"` // "iss_position" ou "visibility"
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Visible   bool    `json:"visible"`
	Timestamp int64   `json:"timestamp"`
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

// loop que recebe mensagens do usuário
func (c *Client) ReadLoop() {
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("WS read error:", err)
			return
		}

		var incoming IncomingMessage
		if err := json.Unmarshal(msg, &incoming); err != nil {
			log.Println("invalid ws message:", err)
			continue
		}

		if incoming.Type == "update_location" {
			c.mu.Lock()
			c.UserLat = incoming.Lat
			c.UserLon = incoming.Lon
			c.mu.Unlock()

			log.Printf("Nova localização recebida: lat=%.4f lon=%.4f\n",
				incoming.Lat, incoming.Lon)
		}
	}
}

// envia dados da iss para o cliente
func (c *Client) SendPosition(msg OutgoingMessage) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	return c.conn.WriteJSON(msg)
}
