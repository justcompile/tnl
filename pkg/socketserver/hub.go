package socketserver

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]*Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
}

func (h *Hub) GetClientForDomain(domain string) *Client {
	return h.clients[domain]
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.domain] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.domain]; ok {
				delete(h.clients, client.domain)
				close(client.send)
			}
		case message := <-h.broadcast:
			for _, client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client.domain)
				}
			}
		}
	}
}
