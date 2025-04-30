package ws

type Hub struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	SendTo     chan TargetedMessage
}

type TargetedMessage struct {
	RecipientIDs []string
	Message      []byte
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		SendTo:     make(chan TargetedMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.UserID] = client

		case client := <-h.Unregister:
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
			}

		case msg := <-h.SendTo:
			for _, id := range msg.RecipientIDs {
				if client, ok := h.Clients[id]; ok {
					client.Send <- msg.Message
				}
			}
		}
	}
}
