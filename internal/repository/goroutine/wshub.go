package goroutine

import (
	"encoding/json"

	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
)

// WSHubRepo handling how to maintain client data
type WSHubRepo struct {
	// Registered clients.
	Clients map[*entity.Client]bool

	// Inbound messages from the clients.
	Broadcast chan []byte

	// Register requests from the clients.
	Register chan *entity.Client

	// Unregister requests from clients.
	Unregister chan *entity.Client
}

// NewHub is method to initiate wshub repo
func NewHub() repository.WSHub {
	return &WSHubRepo{
		Broadcast:  make(chan []byte),
		Register:   make(chan *entity.Client),
		Unregister: make(chan *entity.Client),
		Clients:    make(map[*entity.Client]bool),
	}
}

func (h *WSHubRepo) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				var jsonData map[string]interface{}
				err := json.Unmarshal(message, &jsonData)
				streamID := jsonData["stream_id"].(string)
				if streamID == client.StreamID && err == nil {
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(h.Clients, client)
					}
				}

			}
		}
	}
}
func (h *WSHubRepo) BroadcastMessage(message []byte) {
	h.Broadcast <- message
}

func (h *WSHubRepo) RegisterClient(c *entity.Client) {
	h.Register <- c
}

func (h *WSHubRepo) UnregisterClient(c *entity.Client) {
	h.Unregister <- c
}
