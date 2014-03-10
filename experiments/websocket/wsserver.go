package main

import (
	"fmt"
    "net/http"

    "code.google.com/p/go.net/websocket"
)

// Echo the data received on the WebSocket.
func echoServer(ws *websocket.Conn) {
    //io.Copy(ws, ws)
    for {
    	in := make([]byte, 512)
    	n, _ := ws.Read(in)
    	s := string(in[:n])
    	fmt.Println(s)
    	ws.Write(in[:n])
    }
}

// This example demonstrates a trivial echo server.
func exampleHandler() {
    http.Handle("/", websocket.Handler(echoServer))
    err := http.ListenAndServe(":12345", nil)
    if err != nil {
        panic("ListenAndServe: " + err.Error())
    }
}

func main() {
	fmt.Println("Starting...")
	exampleHandler()
	fmt.Println("Done with ExampleHandler")
}
