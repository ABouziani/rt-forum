var ws; // Store the WebSocket instance
var pagee = 0;
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
                    if (uname != data.Active) {
                        let a = document.createElement('li')
                        a.style.cursor = "pointer"
                        a.innerText = uname
                        chatdiv.appendChild(a)
                        a.addEventListener('click', () => {
                            pagee = 0
                            getChatBox(uname, 0)
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
        <div onscroll="handleScroll('${receiver}','${pagee}')" class="chat-messages" id="chatMessages">
          <!-- Messages will appear here -->
        <br><br><br><br><br><br><br>
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
            Page: pagee,
        }),
    })
        .then(resp => resp.json())
        .then(data => {

            for (let i = data.length - 1; i >= 0; i--) {
                const chatMessages = document.getElementById("chatMessages");
                const messageElement = document.createElement("div");
                if (data[i].receiver == receiver) {
                    messageElement.className = "message sent";
                } else {
                    messageElement.className = "message received";
                }
                messageElement.textContent = data[i].msg;
                chatMessages.appendChild(messageElement);
                chatMessages.scrollTop = 100;
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


function handleScrollh(receiver) {
    const chatMessages = document.getElementById("chatMessages");
    if (chatMessages.scrollTop == 0) {
        console.log(pagee);
        
        getChatBox(receiver, pagee)
        pagee += 5
    }
}

let lastScrollTop = 0; // Track the last scroll position

function handleScroll(receiver, pagee) {
    const chatMessages = document.getElementById('chatMessages');
    
    // Check if the content overflows and scroll is possible
    if (chatMessages.scrollHeight > chatMessages.clientHeight) {
        const currentScrollTop = chatMessages.scrollTop;

        if (currentScrollTop < lastScrollTop) {
            // User is scrolling up
            triggerScrollUpFunction(receiver, pagee);
        }

        lastScrollTop = currentScrollTop; // Update the last scroll position
    }
}

function triggerScrollUpFunction(receiver) {
    // Your logic when scrolling up (e.g., load more messages)
    console.log('Scrolled up!', receiver, pagee);
    // Add your additional logic here
}