package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"net/http"
)

// Echo the data received on the WebSocket.
func echoServer(ws *websocket.Conn) {
	//io.Copy(ws, ws)
	for {
		buffer := make([]byte, 512)
		readData, _ := ws.Read(buffer)
		str := string(buffer[:readData])
		fmt.Println(str)
		ws.Write(buffer[:readData])
	}
}

func connectionHandler() {
	http.Handle("/", websocket.Handler(echoServer))
	fmt.Println("Start listening after clients...")

	err := http.ListenAndServe(":12345", nil)
	// if err not equal null then panic (stop execution)
	if err != nil {
		panic("Error ListenAndServe: " + err.Error())
	}

}

func main() {
	fmt.Println("Initialize server...")
	connectionHandler()
}
