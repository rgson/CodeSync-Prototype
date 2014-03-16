(function() {
	
	// Make sure WebSocket is supported.
	if (!window["WebSocket"]) {
		alert("You browser does not support WebSocket, which is needed to perform document synchronization. Please update your browser.");
		return;
	}
	
	// ---------------------------------------------------

	/*
	* Constructor for the Document type.
	*/
	function Document(version, content) {
		this,version = version;
		this.content = content;
	}
	
	/*
	* Constructor for the Edit type.
	*/
	function Edit(version, diff) {
		this.v = version;
		this.diff = diff;
	}
	
	/*
	* Constructor for the Message type.
	*/
	function Message(type) {
		this.type = type;
		this.send = function() {
			conn.send(JSON.stringify(this));
		};
	}
	
	// ---------------------------------------------------
	
	// Set config values.
	var editor = document.getElementById("editor");		// The CodeMirror editor.
	var host = "ws://localhost:12345/";					// The endpoint containing the document server.
	var diffInterval = 1000;							// Time between diff calculations.

	// Establish connection.
	var conn = new WebSocket(wshost);
	conn.onclose = function(evt) {
		disconnect();
	}
	conn.onmessage = function(evt) {
		receive(evt.data);
	}
	
	// Set up documents.
	var clientText;				// The current state of the document.
	var clientShadow;			// The state of the document as of the last diff calculation.
	var backupShadow;			// The last state of the document confirmed received by the server.
	var edits = [];				// Queued edits not yet confirmed received by the server.
	
	// Request the full document.
	requestDocument();
	
	
	//TODO fix diff timing. Must wait for ACK or TimeOut after sending diffs before caluclating next diff.
	// Otherwise the clientShadow will be overwritten before it has been copied to backupShadow.
	
	// Set diff timer.
	//var diffIntervalId = setInterval(function(){ findDiff(); }, diffInterval);
	
	// ---------------------------------------------------
	
	/*
	* Function called when losing connection.
	*/
	function disconnect() {
		alert("You have lost connection to the synchronization server!");
	}
	
	/*
	* Function called when receiving a message.
	*/
	function receive(data) {
		try {
			var msg = JSON.parse(data);
			
			switch(msg.type) {
			case "diff":
				handleDiffMessage(msg);
				break;
			case "ack":
				handleAckMessage(msg);
				break;
			case "doc":
				handleDocMessage(msg);
				break;
			}
			
		} catch (ex) {
			// invalid json
		}
	}
	
	/*
	* Requests the full document from the document server.
	*/
	function requestDocument() {
		var message = new Message("req");
		message.send();
	}

	/*
	* Sends the queued edits to the server.
	*/
	function sendEdits() {
		var message = new Message("diff");
		message.received_v = clientShadow.remoteVersion;
		message.edits = edits;
		message.send();
	}
	
	/*
	* Sends an acknowledgement of a received edit to the server.
	*/
	function sendAck() {
		var message = new Message("ack");
		message.received_v = clientShadow.remoteVersion;
		message.send();
	}
	
	/*
	* Handles a received "diff" message.
	*/
	function handleDiffMessage(msg) {
		//TODO
	}
	
	/*
	* Handles a received "ack" message.
	*/
	function handleAckMessage(msg) {
		for(var i = 0; i < edits.length; i++) {
			if(edits[i].v <= msg.received_v)
				edits.shift();
			else
				break;
		}
	}
	
	/*
	* Handles a received "doc" message.
	*/
	function handleDocMessage(msg) {
		clientText = new Document(1, msg.content;
		
		clientShadow = clientText;
		clientShadow.remoteVersion = msg.v;
		
		backupShadow = clientShadow;
		
		edits = [];
	}	
	
	/*
	* Calculates the difference between the Client Text and Client Shadow.
	*/
	function findDiff() {
		// TODO
		// calculate diff
		// create edit
		// put edit in edits queue
	}
	
})();
