package main

type room struct {
	forward chan []byte // forward は他のクライアントに転送するためのメッセージを保持するチャネルです
}
