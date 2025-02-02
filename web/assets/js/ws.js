


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

    return ws; 
}


window.addEventListener("load", getWebSocket);


function getChatBox(receiver) {
    if (pagee < 10) {
        document.querySelector('.container').innerHTML = `
            <div class="chat-container">
            <div id="receiver">${receiver}</div>
            <div onscroll="handleScroll('${receiver}','${pagee}')" class="chat-messages" id="chatMessages">
            
            </div>
            <div class="chat-input">
              <input type="text" id="chatInput" placeholder="Type your message...">
              <button onclick="sendMessage('${receiver}')">Send</button>
            </div>
          </div>
            `
    }
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
            if (data){
                for (let i = 0; i < data.length; i++) {
                    const chatMessages = document.getElementById("chatMessages");
                    const messageElement = document.createElement("div");
                    if (data[i].receiver == receiver) {
                        messageElement.className = "message sent";
                    } else {
                        messageElement.className = "message received";
                    }
                    messageElement.textContent = data[i].msg;
                    chatMessages.prepend(messageElement);
                    chatMessages.scrollTop = 100;
                }
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

    ws.send(JSON.stringify({
        Receiver: uname,
        Msg: message,
    }))
}


function handleScroll(receiver) {
    const chatMessages = document.getElementById("chatMessages");
    if (chatMessages.scrollTop == 0) {
        pagee += 10
        trchatbox(receiver, pagee)
    }
}

function debounce(fn, delay) {
    let timer = null;
    return function () {
        let context = this;
        let args = arguments;
        clearTimeout(timer);
        timer = setTimeout(function () {
            fn.apply(context, args);
        }, delay);
    };
}

const trchatbox = debounce(getChatBox,2000)