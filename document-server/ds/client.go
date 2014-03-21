package ds

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
)

type Client struct {
	Connection *websocket.Conn
	Document *Document
}

func (c *Client) Handle() {
	
	// TODO Use ticker here to time the edits/sends
	//go calculateEdit()		// Run a timed goroutine to calculate edits
	//go sendEdits()			// Run a timed goroutine to send available edits
	
	for {
		buffer := make([]byte, 1024)
		length, err := c.Connection.Read(buffer)
		
		if err != nil {
			fmt.Println("Error listening:", err.Error)
			break; // Stops listening to this client!
		}
		
		msg, msgtype := ParseMessage(buffer[:length])
		c.handleMessage(msg, msgtype)
	}

}

func (c *Client) initialize() {

	c.Document = &Document{}
	c.Document.Initialize()
	c.sendDocument()

}

func (c *Client) handleMessage(msg interface{}, msgtype string) {

	switch msgtype {
	
	case MSGTYPE_EDIT:
		c.handleEditMessage(msg.(EditMessage))
		
	case MSGTYPE_ACK:
		c.handleAckMessage(msg.(AckMessage))
		
	case MSGTYPE_REQUEST:
		c.handleRequestMessage(msg.(RequestMessage))
		
	}
	
}

func (c *Client) handleEditMessage(msg EditMessage) {

	err := c.Document.ApplyEdits(msg.Version, msg.Edits)
	if err != nil {
		// Patch unsuccessful - reinitialize and accept loss.
		c.initialize()
	} else {
		// Patch successful - send ack.
		c.sendAck()
	}

}

func (c *Client) handleAckMessage(msg AckMessage) {
	c.Document.RemoveConfirmedEdits(msg.Version)
}

func (c *Client) handleRequestMessage(msg RequestMessage) {
	c.initialize()
}

func (c *Client) sendEdits() {
	msg := NewEditMessage(c.Document.GetEdits())
	c.Connection.Write(msg.ToJSON())
}

func (c *Client) sendAck() {
	msg := NewAckMessage(c.Document.Shadow.RemoteVersion)
	c.Connection.Write(msg.ToJSON())
}

func (c *Client) sendDocument() {
	msg := NewDocumentMessage(c.Document.Shadow.Content)
	c.Connection.Write(msg.ToJSON())
}
