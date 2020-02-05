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

	command chan string
}

func NewHub(c chan string) *Hub {
	return &Hub{
		register:   make(chan *Client, 10),
		unregister: make(chan *Client, 10),
		clients:    make(map[string]*Client),
		command:    c,
	}
}

func (h *Hub) Send(deviceID string, message string) {
	if _, ok := h.clients[deviceID]; ok {
		h.clients[deviceID].Send <- []byte(message)
	} else {
		log.Error(fmt.Sprintf("Device %s is not connected", deviceID))
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

				// TODO: call detach
				h.command <- fmt.Sprintf("%s detach", client.ID)

				// TODO: call unsubscribe
				//h.command <- fmt.Sprintf("%s unsubscribe", client.ID)

				log.Info(fmt.Sprintf("Device disconnected: %v", client))
			}
		}
	}
}
