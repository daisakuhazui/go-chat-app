package main

import (
	"go-chat-app/trace"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

type room struct {
	forward chan *message    // forward は他のクライアントに転送するためのメッセージを保持するチャネルです
	join    chan *client     // join はチャットルームに参加しようとしているクライアントのためのチャネルです
	leave   chan *client     // leave はチャットルームから退室しようとしているクライアントのためのチャネルです
	clients map[*client]bool // client には在室している全てのクライアントが保持されています
	tracer  trace.Tracer     // tracerはチャットルーム上で行われた操作のログを受け取ります
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// 参加
			r.clients[client] = true
			r.tracer.Trace("New client has joined: ", client.userData["name"])
		case client := <-r.leave:
			// 退室
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client leaved: ", client.userData["name"])
		case msg := <-r.forward:
			r.tracer.Trace("Got message: ", msg.Message)
			// 全てのクライアントにメッセージを転送
			for client := range r.clients {
				select {
				case client.send <- msg:
					// メッセージを送信
					r.tracer.Trace(" -- Has sent a message to client")
				default:
					// 送信に失敗
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace(" --  Fail to send a message. Cleaning up client.")
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP: ", err)
		return
	}
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatalln("Fail to get Cookie:", err)
		return
	}
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}

func newRoom() *room {
	// newRoom はすぐに利用できるチャットルームを生成して返します
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}
