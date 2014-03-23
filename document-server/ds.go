package main

import (
	"code.google.com/p/go.net/websocket"
	"ds"
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	CONN_HOST       = "localhost"     // Listening network.
	CONN_PORT       = "4343"          // Listening port.
	HEARTBEAT_HOST  = "localhost"     // Target host for heartbeat monitor.
	HEARTBEAT_PORT  = "3434"          // Target port for heartbeat monitor.
	HEARTBEAT_TIMER = 8 * time.Second // Heartbeat interval.
)

var clients int64 = 0
var nextId int = 0

// Echo the data received on the WebSocket.
func serveClient(ws *websocket.Conn) {
	id := addClient()
	client := ds.Client{ID: id, Connection: ws}
	client.Handle()
	dropClient()
}

func connectionHandler() {
	http.Handle("/", websocket.Handler(serveClient))
	err := http.ListenAndServe(CONN_HOST+":"+CONN_PORT, nil)
	// if err not equal null then panic (stop execution)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
}

func main() {
	fmt.Println("Starting document server...")

	// Dial heartbeat monitor.
	conn, err := net.Dial("tcp", HEARTBEAT_HOST+":"+HEARTBEAT_PORT)
	if err != nil {
		fmt.Println("Error dialing:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	// Run a timed goroutine to send heartbeats.
	heartbeatTicker := time.NewTicker(HEARTBEAT_TIMER)
	go func() {
		for _ = range heartbeatTicker.C {
			buf := make([]byte, 4)
			binary.PutVarint(buf, clients)
			_, err := conn.Write(buf)
			if err != nil {
				fmt.Println("Heartbeat failed:", err.Error())
			} else {
				fmt.Println("Sent heartbeat")
			}
		}
	}()
	defer heartbeatTicker.Stop()

	fmt.Println("Heartbeat running!")

	connectionHandler()

	fmt.Println("Stopped")
}

func addClient() int {
	clients++
	nextId++
	id := nextId
	fmt.Printf("Client connected. %d clients connected.\n", clients)
	return id
}

func dropClient() {
	clients--
	fmt.Printf("Client disconnected. %d clients connected.\n", clients)
}
