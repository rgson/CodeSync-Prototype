package main

import (
	"fmt"
    "net/http"
	"encoding/json"
    "code.google.com/p/go.net/websocket"
)

type Document struct {
	Version int
	RemoteVersion int
	Content string
}

type Message struct {
	DocType string `json:"type"`
}

type DocumentMessage struct {
	DocType string `json:"type"`
	Version int `json:"v"`
	Content string `json:"content"`
}

type AckMessage struct {
	DocType string `json:"type"`
	Version int `json:"v"`
}

type EditMessage struct {
	DocType string `json:"type"`
	Version int `json:"v"`
	Edits []Edit `json:"edits"`
}

type Edit struct {
	Version int `json:"v"`
	Patch string `json:"patch"`
	MD5 string `json:"md5"`
}

var doc Document = Document{Version: 0, RemoteVersion: 0, Content: "The document's content goes here"}

func echoServer(ws *websocket.Conn) {
	i := 0
    for {
    	in := make([]byte, 512)
    	n, _ := ws.Read(in)
    	s := string(in[:n])
    	fmt.Print(i)
    	fmt.Println(": " + s)
    	i++
    	
    	msg := &Message{}
    	json.Unmarshal(in[:n], &msg)
    	switch msg.DocType {
    	case "req":
    		response := &DocumentMessage{DocType: "doc", Version: doc.Version, Content: doc.Content}
    		json, _ := json.Marshal(response)
    		ws.Write(json)
    	case "edit":
    		editMsg := &EditMessage{}
    		json.Unmarshal(in[:n], &editMsg)
    		doc.RemoteVersion = editMsg.Edits[len(editMsg.Edits)-1].Version
    		response := &AckMessage{DocType: "ack", Version: doc.RemoteVersion}
    		json, _ := json.Marshal(response)
    		ws.Write(json)
    	}
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
