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
