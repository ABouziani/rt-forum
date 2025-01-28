let ws; // Store the WebSocket instance

function getWebSocket() {
    if (!ws || ws.readyState === WebSocket.CLOSED) {
        ws = new WebSocket("ws://localhost:8080/ws");

        // Handle incoming messages
        ws.onmessage = (event) => {
            let data;
            try {
                // Parse the incoming JSON message
                data = JSON.parse(event.data);
            } catch (e) {
                console.error("Failed to parse message:", event.data);
                return; // Exit if JSON parsing fails
            }

            // Check if the message has a "Sender" property
            if (!data.Sender) {
                let chatdiv = document.getElementById('chat-section')
                chatdiv.innerText = ""
                for (const uname of data.Online) {
                    let a = document.createElement('li')
                    a.style.cursor = "pointer"
                    a.innerText = uname
                    chatdiv.appendChild(a)
                    a.addEventListener('click', () => {
                        getChatBox(uname)
                    })
                }
            }
        };

        // Handle errors
        ws.onerror = (error) => {
            console.error("WebSocket error:", error);
        };

        // Handle connection close and attempt reconnection
        ws.onclose = () => {
            console.log("WebSocket connection closed. Attempting to reconnect...");
            setTimeout(getWebSocket, 1000); // Reconnect after 1 second
        };

        // Log successful connection
        ws.onopen = () => {
            console.log("WebSocket connection established.");
        };
    }

    return ws; // Return the same WebSocket instance
}


window.addEventListener("load", getWebSocket);


function getChatBox(uname) {
    document.querySelector('.container').innerHTML = `
    <div class="chat-container">
    <div class="chat-messages" id="chatMessages">
      <!-- Messages will appear here -->
    <!-- <div class="message sent">Hello!</div>
      <div class="message received">Hi there!</div> -->
    </div>
    <div class="chat-input">
      <input type="text" id="chatInput" placeholder="Type your message...">
      <button onclick="sendMessage('${uname}')">Send</button>
    </div>
  </div>
    `
}


function sendMessage(uname){
    let message = document.querySelector('#chatInput').value
    fetch("/sendmessage", {
        method: "POST",
        header: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            receiver: uname,
            msg:message,
        }),
    })
}