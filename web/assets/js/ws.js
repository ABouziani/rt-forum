var ws; // Store the WebSocket instance

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
            if (data.Online) {
                console.log(data.Active);
                
                let chatdiv = document.getElementById('chat-section')
                chatdiv.innerText = ""
                for (const uname of data.Online) {
                    if (uname != data.Active){
                        let a = document.createElement('li')
                        a.style.cursor = "pointer"
                        a.innerText = uname
                        chatdiv.appendChild(a)
                        a.addEventListener('click', () => {
                            getChatBox(uname)
                        })
                    }
                }
            } else if (data.msg) {
                addMsg(data)
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


function getChatBox(receiver) {
    document.querySelector('.container').innerHTML = `
        <div class="chat-container">
        <div id="receiver">${receiver}</div>
        <div class="chat-messages" id="chatMessages">
          <!-- Messages will appear here -->
        <!-- <div class="message sent">Hello!</div>
          <div class="message received">Hi there!</div> -->
        </div>
        <div class="chat-input">
          <input type="text" id="chatInput" placeholder="Type your message...">
          <button onclick="sendMessage('${receiver}')">Send</button>
        </div>
      </div>
        `

    fetch("/fetchmessages", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            Receiver: receiver,
        }),
    })
        .then(resp => resp.json())
        .then(data => {
            console.log(data);
            
            for (let x of data) {
                console.log(x);
                
                const chatMessages = document.getElementById("chatMessages");
                const messageElement = document.createElement("div");
                if (x.receiver == receiver) {
                    messageElement.className = "message sent";
                } else {
                    messageElement.className = "message received";
                }
                messageElement.textContent = x.msg;
                chatMessages.appendChild(messageElement);
            }
        })
}


function addMsg(data) {
    let receiver = document.querySelector('#receiver').innerText
    let message = data.msg
    const chatInput = document.getElementById("chatInput");
    const chatMessages = document.getElementById("chatMessages");
    const messageElement = document.createElement("div");
    if (data.receiver == receiver) {
        messageElement.className = "message sent";
    } else {
        messageElement.className = "message received";
    }
    messageElement.textContent = message;
    chatMessages.appendChild(messageElement);
    chatInput.value = "";
    chatMessages.scrollTop = chatMessages.scrollHeight;
}

function sendMessage(uname) {
    let message = document.querySelector('#chatInput').value
    // const chatInput = document.getElementById("chatInput");
    // const chatMessages = document.getElementById("chatMessages");
    ws.send(JSON.stringify({
        Receiver: uname,
        Msg: message,
    }))
    // fetch("/sendmessage", {
    //     method: "POST",
    //     header: {
    //         "Content-Type": "application/json"
    //     },
    //     body: ,
    // }).then(data => {
    //     if (data.ok) {
    //         const messageElement = document.createElement("div");
    //         messageElement.className = "message sent";
    //         messageElement.textContent = message;
    //         chatMessages.appendChild(messageElement);
    //         chatInput.value = "";
    //         chatMessages.scrollTop = chatMessages.scrollHeight;
    //     }
    // })
}