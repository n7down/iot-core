package gateway

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]*Client

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client, 10),
		unregister: make(chan *Client, 10),
		clients:    make(map[string]*Client),
	}
}

func (h *Hub) Send(deviceID string, message string) {
	if _, ok := h.clients[deviceID]; ok {
		h.clients[deviceID].Send <- []byte(message)
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			//if _, ok := h.clients[client.ID]; ok {
			//delete(h.clients, client.ID)
			//close(client.Send)
			//log.Info(fmt.Sprintf("Client disconnected: %v", client))
			//}

			h.clients[client.ID] = client
			log.Info(fmt.Sprintf("Device connected: %v", client))
		case client := <-h.unregister:
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.Send)
				// TODO: call detach?
				log.Info(fmt.Sprintf("Device disconnected: %v", client))
			}
		}
	}
}
