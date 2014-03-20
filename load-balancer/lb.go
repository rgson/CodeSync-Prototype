package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	CONN_HOST = "localhost" // Lokal adress. Anpassas efter nätverket.
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
	DNS_TIMER = 300
)

var documentServer = make(map[string]int64)

func main() {
	go selectServer()
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	defer conn.Close()
	sec, _ := time.ParseDuration("10s")
	for {
		// Make a buffer to hold incoming data.
		buf := make([]byte, 4)

		conn.SetReadDeadline(time.Now().Add(sec))
		// Read the incoming connection into the buffer.
		_, err := conn.Read(buf)

		if err != nil {
			fmt.Println("Error reading:", err.Error())
			delete(documentServer, conn.RemoteAddr().String())
			break
		}

		buf2, _ := binary.Varint(buf)

		str := conn.RemoteAddr().String()
		documentServer[str] = buf2

		// Send the request to a new reciever.

	}

}

func selectServer() {
	for {
		var min string = ""

		for key, value := range documentServer {
			if min == "" {
				min = key
				continue
			}

			if documentServer[min] > value {
				min = key
			}
		}
		if min != "" {
			changeDNS(min)
		}
		time.Sleep(DNS_TIMER * time.Second)
	}
}

func changeDNS(ip string) {
	fmt.Println(ip) //TODO skaffa domännamn. Uppdatera DNS
}
