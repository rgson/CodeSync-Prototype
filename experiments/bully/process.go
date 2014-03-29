package main

import (
	"os"
	"fmt"
	"net"
	"flag"
	"time"
	"encoding/json"
)

const (
	CONFIG_FILE = "processes.json"
	TIMEOUT = 5 * time.Second
	PING_INTERVAL = 10 * time.Second
	MSG_PING = "ping"
	MSG_VOTE = "vote"
	MSG_ALIVE = "alive"
	MSG_LEADER = "leader"
)

var id int
var processes = make(map[int]string)
var leader int

func main() {

	// Read ID from flag
	flag.IntVar(&id, "id", -1, "The process' ID")
	flag.Parse()
	
	if id == -1 {
		fmt.Println("A process ID must be provided with the '-id' flag.")
		os.Exit(1)
	}
	
	// Read processes from config
	readConfig()
	
	if _, exists := processes[id]; !exists {
		fmt.Println("Could not find the provided ID in the config file.")
		os.Exit(1)
	}
	
	// Listen for other processes
	listenAddr := processes[id]
	delete(processes, id)
	go listen(listenAddr)
	
	// Guess leader (highest ID)
	guessLeader()
	
	// Ping leader at interval
	for {
		if leader != id {
			fmt.Printf("%d: Pinging leader.\n", id)
			alive := pingLeader()
			if !alive {
			fmt.Printf("%d: No response from leader!\n", id)
				announceVote()
			}
		}
		time.Sleep(PING_INTERVAL)
	}

}

func readConfig() {

	file, err := os.Open(CONFIG_FILE)
	if err != nil {
		fmt.Println("Failed to open configuration file:", CONFIG_FILE)
		os.Exit(1)
	}
	defer file.Close()
	
	decoder := json.NewDecoder(file)
	procs := &Processes{}
	fmt.Printf("%v\n", decoder.Decode(&procs))
	
	for _, proc := range procs.Processes {
		processes[proc.ID] = fmt.Sprintf("%s:%d", proc.Host, proc.Port)
	}

}

func listen(addr string) {

	l, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	
	fmt.Printf("%d: Listening...\n", id)
	
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleConn(conn)
	}

}

func handleConn(conn net.Conn) {
	defer conn.Close()
	
	fmt.Printf("%d: Got connection.\n", id)
	
	// Read from connection
	buf := make([]byte, 512)
	conn.SetReadDeadline(time.Now().Add(TIMEOUT))
	n, err := conn.Read(buf)
	if err != nil {
		return
	}
	
	// Unmarshal message
	msg := Message{}
	err = json.Unmarshal(buf[:n], &msg)
	
	// Handle message
	switch msg.Type {
	case MSG_VOTE:
		fmt.Printf("%d: Got vote message.\n", id)
		go announceVote()
		fallthrough
	case MSG_PING:
		fmt.Printf("%d: Sending 'alive' message.\n", id)
		response := Message{Type: MSG_ALIVE, Sender: id}
		conn.Write(response.ToJSON())
	case MSG_LEADER:
		if msg.Sender > id {
			leader = msg.Sender
			fmt.Printf("%d: New leader %d.\n", id, leader)
		} else {
			fmt.Printf("%d: Leader with lower ID. Bullying!\n", id)
			go announceVote()
		}
	}

}

func guessLeader() {

	leader = id
	for k, _ := range processes {
		if leader < k {
			leader = k
		}
	}
	fmt.Printf("%d: New leader.\n", id)
	
	if leader == id {
		announceVictory()
	}
	
}

func announceVictory() {

	fmt.Printf("%d: I'm the new leader!\n", id)
	
	leader = id
	for pid, _ := range processes {
		msg := Message{Type: MSG_LEADER, Sender: id}
		send(pid, msg)
	}

}

func announceVote() {

	fmt.Printf("%d: The leader is down!\n", id)

	c := make(chan bool)
	
	for pid, _ := range processes {
		if id < pid {
			go ask(pid, c)
		}
	}
	
	select {
	case <-c:
		return
	case <-time.After(TIMEOUT):
		fmt.Printf("%d: No anwser from higher processes.\n", id)
		announceVictory()
	}			

}

func pingLeader() bool {

	// Dial leader
	conn, err := dial(leader)
	if err != nil {
		return false
	}
	defer conn.Close()
	
	// Send ping
	msg := Message{Type: MSG_PING, Sender: id}
	_, err = conn.Write(msg.ToJSON())
	if err != nil {
		return false
	}
	
	// Wait for answer
	buf := make([]byte, 512)
	conn.SetReadDeadline(time.Now().Add(TIMEOUT))
	_, err = conn.Read(buf)
	if err != nil {
		return false
	}

	return true

}

func dial(pid int) (conn net.Conn, err error) {

	conn, err = net.Dial("tcp", processes[pid])
	if err != nil {
		err = fmt.Errorf("Error dialing process %d", pid)
		fmt.Printf("%d: %s\n", id, err.Error())
	}
	
	return

}

func send(pid int, msg Message) error {

	conn, err := dial(pid)
	if err == nil {
		_, err = conn.Write(msg.ToJSON())
	}
	
	return err

}

func ask(pid int, c chan bool) {

	msg := Message{Type: MSG_VOTE, Sender: id}

	// Dial process
	conn, err := dial(pid)
	if err != nil {
		return
	}
	defer conn.Close()

	// Send vote announcement
	_, err = conn.Write(msg.ToJSON())
	if err != nil {
		return
	}

	// Wait for answer
	buf := make([]byte, 512)
	conn.SetReadDeadline(time.Now().Add(TIMEOUT))
	_, err = conn.Read(buf)
	if err != nil {
		return
	}

	fmt.Printf("%d: Answer from %d.\n", id, pid)
	c <- true

}

type Processes struct {
	Processes []Process `json:"processes"`
}

type Process struct {
	ID int `json:"id"`
	Host string `json:"host"`
	Port int `json:"port"`
}

type Message struct {
	Type string `json:"type"`
	Sender int `json:"sender"`
}

func (m Message) ToJSON() []byte {
	msg, _ := json.Marshal(m)
	return msg
}
