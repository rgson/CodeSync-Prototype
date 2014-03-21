package main

import (
	"code.google.com/p/go.net/websocket"
	"ds"
	"fmt"
	"net/http"
	"os"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "12345"
)

var nextId int = 0

// Echo the data received on the WebSocket.
func serveClient(ws *websocket.Conn) {
	id := nextId
	nextId++
	fmt.Println("Found client:", id)
	client := ds.Client{ID: id, Connection: ws}
	client.Handle()
	fmt.Println("Lost client:", id)
}

func connectionHandler() {
	http.Handle("/", websocket.Handler(serveClient))
	err := http.ListenAndServe(CONN_HOST+":"+CONN_PORT, nil)
	// if err not equal null then panic (stop execution)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
}

func main() {
	fmt.Println("Starting document server...")
	connectionHandler()
	fmt.Println("Stopped")
}
