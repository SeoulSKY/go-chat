package main

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type Client struct {
	id      int64
	name    string
	conn    *websocket.Conn
	send    chan *Chat
	manager *ClientManager
}

type Chat struct {
	SenderName string    `json:"senderName"`
	Message    string    `json:"message"`
	Timestamp  time.Time `json:"timestamp"`
}

// NewClient New creates a new client
func NewClient(id int64, name string, conn *websocket.Conn, manager *ClientManager) *Client {
	return &Client{
		id,
		name,
		conn,
		make(chan *Chat),
		manager,
	}
}

// start the Client
func (c *Client) start() {
	c.manager.add <- c

	// read message from the client
	go func() {
		for {
			message := new(Chat)
			err := c.conn.ReadJSON(message)

			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println(err)
				}

				c.manager.remove <- c
				return
			}

			message.Timestamp = time.Now()
			log.Println(message.SenderName, "says:", message.Message)

			c.manager.broadcast <- message
		}
	}()

	// write message to the client
	go func() {
		for {
			select {
			case message := <-c.send:
				err := c.conn.WriteJSON(message)

				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						log.Println(err)
					}

					c.manager.remove <- c
					return
				}
			}
		}
	}()
}
