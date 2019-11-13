package gateway

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	//broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		//broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Send(deviceID string, message string) {
	// TODO: find the client by the id
	// TODO: add message to client.send channel
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:

			// TODO: manage this by ID
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			//case message := <-h.broadcast:
			//for client := range h.clients {
			//select {
			//case client.send <- message:
			//default:
			//close(client.send)
			//delete(h.clients, client)
			//}
			//}
		}
	}
}
