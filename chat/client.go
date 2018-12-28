package main

import (
	"github.com/gorilla/websocket"
)

// client はチャットを行っている1人のユーザーを表します
type client struct {
	socket *websocket.Conn // socket はこのクライアントのための WebSocket です
	send   chan []byte     // send はメッセージが送られるチャネルです
	room   *room           // room はこのクライアントが参加しているチャットルームです
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
