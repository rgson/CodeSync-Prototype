package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	CONN_HOST = "localhost" 		// Listening network.
	CONN_PORT = "3434"				// Listening port.
	DS_TIMEOUT = 10 * time.Second	// Timeout period for a document server's heartbeat signal.
	DNS_TIMER = 300 * time.Second	// Interval för DNS updates.
)

var documentServer = make(map[string]int64)

func main() {
	
	// Listen for incoming connections.
	l, err := net.Listen("tcp", CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	
	// Run a timed goroutine to select the best server.
	dnsTicker := time.NewTicker(DNS_TIMER)
	go func() {
		for _ = range dnsTicker.C {
			selectServer()
		}
	}()
	defer dnsTicker.Stop()
	
	// Actively listen for connections.
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
	
	remoteAddr := conn.RemoteAddr().String()
		
	fmt.Println("Got connection:", remoteAddr)
	
	for {
		// Make a buffer to hold incoming data.
		buf := make([]byte, 4)

		conn.SetReadDeadline(time.Now().Add(DS_TIMEOUT))
		// Read the incoming data into the buffer.
		_, err := conn.Read(buf)

		if err != nil {
			fmt.Println("Error reading:", err.Error())
			delete(documentServer, remoteAddr)
			printServers()
			break
		}
		
		clients, _ := binary.Varint(buf)
		documentServer[remoteAddr] = clients
		fmt.Printf("Heartbeat:\t%s\t%d\n", remoteAddr, clients)
		
		printServers()
	}
}

func selectServer() {
	var min string = ""

	// Find server with least connections.
	for ip, clients := range documentServer {
		if min == "" {
			min = ip
			continue
		}

		if documentServer[min] > clients {
			min = ip
		}
	}
	
	// Make sure a server was found (i.e. at least one server exists).
	if min != "" {
		// Change the active server.
		changeDNS(min)
	}
}

func changeDNS(ip string) {
	//TODO skaffa domännamn. Uppdatera DNS
	fmt.Printf("DNS: Active server is %s with %d clients.\n", ip, documentServer[ip])
}

func printServers() {
	fmt.Printf("LST: ------\n")
	for key, value := range documentServer {
		fmt.Printf("LST: %s\t%d\n", key, value)
	}
	fmt.Printf("LST: ------\n")
}
