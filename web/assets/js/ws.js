


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
            console.log(data);

            // Check if the message has a "Sender" property
            if (data.Online) {

                let chatdiv = document.getElementById('chat-section')
                chatdiv.innerText = ""
                for (const uname of data.Online) {
                    if (uname != data.Active) {
                        let a = document.createElement('li')
                        a.className = 'user'
                        a.style.cursor = "pointer"
                        a.innerHTML = `<span class="fa-regular fa-user"></span> <span style="margin-top:5px;" class="status-dot online"></span>${uname}`
                        chatdiv.appendChild(a)
                        a.addEventListener('click', () => {
                            pagee = 0
                            getChatBox(uname)
                        })
                    }
                }
                for (const uname of data.NotOnline) {
                    if (uname != data.Active) {
                        let a = document.createElement('li')
                        a.className = 'user'
                        a.style.cursor = "pointer"
                        a.innerHTML = `<span class="fa-regular fa-user"></span> <span style="margin-top:5px;" class="status-dot offline"></span>${uname}`
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


function getChatBox(receiver) {
    if (pagee < 10) {
        document.querySelector('.container').innerHTML = `
        <p style="margin:auto;" class="currentPage">Conversation</p>
            <div class="chat-container">
            <div style="padding-left:20px;padding-bottom:20px;"><span class="fa-regular fa-user"></span><span id="receiver" style="margin-left:10px;">${receiver}</span></div>
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
            if (data) {
                for (let i = 0; i < data.length; i++) {
                    const chatMessages = document.getElementById("chatMessages");
                    const messageElement = document.createElement("div");
                    if (data[i].receiver == receiver) {
                        messageElement.className = "message sent";
                    } else {
                        messageElement.className = "message received";
                    }
                    messageElement.innerHTML = `
        <div class="header">
        <span style="color:black;font-weight: 700;" class="username">${data[i].Sender}</span>
        <span style="margin-left:50px;" class="timestamp">${data[i].created_at}</span>
    </div>
    <div style="margin-top:15px; text-align: left;" class="content">
        <span>${data[i].msg}</span>
    </div>
        `
                    chatMessages.prepend(messageElement);
                    chatMessages.scrollTop = 100;
                }
            }

        })
}


function addMsg(data) {
    if (document.querySelector('#receiver')) {

        let receiver = document.querySelector('#receiver').innerText
        console.log(receiver);

        const chatInput = document.getElementById("chatInput");
        const chatMessages = document.getElementById("chatMessages");
        const messageElement = document.createElement("div");
        if (data.receiver == receiver) {
            messageElement.className = "message sent";
        } else {
            messageElement.className = "message received";
        }
        messageElement.innerHTML = `
        <div class="header">
        <span style="color:black;font-weight: 700;" class="username">${data.Sender}</span>
        <span style="margin-left:auto;" class="timestamp">${data.created_at}</span>
    </div>
    <div style="margin-top:15px; text-align: left;" class="content">
        <span>${data.msg}</span>
    </div>
        `
        chatMessages.appendChild(messageElement);

        chatInput.value = "";
        chatMessages.scrollTop = chatMessages.scrollHeight;
    }
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

const trchatbox = debounce(getChatBox, 2000)

ws = getWebSocket()