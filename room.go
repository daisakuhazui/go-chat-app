package main

type room struct {
	forward chan []byte      // forward は他のクライアントに転送するためのメッセージを保持するチャネルです
	join    chan *client     // join はチャットルームに参加しようとしているクライアントのためのチャネルです
	leave   chan *client     // leave はチャットルームから退室しようとしているクライアントのためのチャネルです
	client  map[*client]bool // client には在室している全てのクライアントが保持されています
}
