package main

import (
	"code.google.com/p/go.net/websocket"
	"codesync/ds"
	"codesync/lb"
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

var clients int = 0
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

	go heartbeat()

	fmt.Println("Heartbeat running!")

	connectionHandler()

	fmt.Println("Stopped")
}

func heartbeat() {

	for {
		sendHeartbeat()
		time.Sleep(HEARTBEAT_TIMER)
	}

}

func sendHeartbeat() {

	conn, err := net.Dial("tcp", HEARTBEAT_HOST+":"+HEARTBEAT_PORT)
	if err != nil {
		fmt.Println("Error dialing:", err.Error())
		return
	}
	defer conn.Close()

	msg := lb.Message{Address: CONN_HOST+":"+CONN_PORT, Load: clients}
	_, err = conn.Write(msg.ToJSON())
	if err != nil {
		fmt.Println("Heartbeat failed:", err.Error())
	} else {
		fmt.Println("Sent heartbeat")
	}
	
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
