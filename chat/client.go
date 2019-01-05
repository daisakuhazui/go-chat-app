package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// client はチャットを行っている1人のユーザーを表します
type client struct {
	socket   *websocket.Conn        // socket はこのクライアントのための WebSocket です
	send     chan *message          // send はメッセージが送られるチャネルです
	room     *room                  // room はこのクライアントが参加しているチャットルームです
	userData map[string]interface{} // userData はユーザーに関する情報を保持しします
}

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			if avatarURL, ok := c.userData["avatar-url"]; ok {
				msg.AvatarURL = avatarURL.(string)
			}
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
