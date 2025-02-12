
var pagee = 0;
var worker = new SharedWorker('/assets/js/worker.js');
worker.port.start();

worker.port.onmessage = (event) => {
    
    if (event.data == 'logout' || event.data == 'login') {
        refetchLogin('/login')
        return
    }

    let data;
    try {
        data = JSON.parse(event.data);
    } catch (e) {
        console.error("Failed to parse message:", event.data);
        return;
    }

    let chatdiv = document.getElementById('chat-section');
    let chatdivmobile = document.getElementById('chat-mobile');
    if (chatdiv && !data.msg && chatdivmobile) {
        chatdiv.innerText = "";
        chatdivmobile.innerText = "";
        if (data.Users) {
            for (const user of data.Users) {
                let a = document.createElement('li');
                a.className = 'user';
                a.style.cursor = "pointer";
                a.innerHTML = `<span class="fa-regular fa-user"></span> <span style="margin-top:5px;" class="status-dot ${user.Status}"></span>${user.Uname}`;
                let b = a.cloneNode(true);
                chatdiv.appendChild(b);
                chatdivmobile.appendChild(a);
                a.addEventListener('click', () => {
                    pagee = 0;
                    getChatBox(user.Uname);
                });
                b.addEventListener('click', () => {
                    pagee = 0;
                    getChatBox(user.Uname);
                });
            }
        }
    } else if (data.msg) {
        addMsg(data);
    }
};

worker.port.postMessage('helooo')

function getChatBox(receiver, s) {
    if (pagee < 10) {
        document.querySelector('.container').innerHTML = `
        <button class="nav-button" onclick="displayMobileNav()">
                <i class="fa-solid fa-bars"></i>
            </button>
        <p style="margin:auto;" class="currentPage">Conversation</p>
            <div class="chat-container">
            <div style="padding-left:20px;padding-bottom:20px;"><span class="fa-regular fa-user"></span><span id="receiver" style="margin-left:10px;">${receiver}</span></div>
            <div onscroll="handleScroll('${receiver}','${pagee}')" class="chat-messages" id="chatMessages">
            
            </div>
            <div class="chat-input">
              <input type="text" id="chatInput" placeholder="Type your message...">
              <button onclick="trsendMessage('${receiver}')">Send</button>
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
        <span style="margin-right:20px;color:black;font-weight: 700;" class="username">${data[i].Sender}</span>
        <span class="timestamp">${data[i].created_at}</span>
    </div>
    <div style="margin-top:15px; text-align: left;" class="content">
        <span>${data[i].msg}</span>
    </div>
        `
                    chatMessages.prepend(messageElement);
                    chatMessages.scrollTop = chatMessages.scrollHeight;
                }
                if (s) {
                    chatMessages.scrollTop = s
                } else {
                    chatMessages.scrollTop = chatMessages.scrollHeight;
                }

            }

        })
}


function addMsg(data) {
    if (document.querySelector('#receiver')) {

        let receiver = document.querySelector('#receiver').innerText

        const chatInput = document.getElementById("chatInput");
        const chatMessages = document.getElementById("chatMessages");
        const messageElement = document.createElement("div");
        if (data.receiver == receiver) {
            pagee++
            messageElement.className = "message sent";
        } else if (data.Sender == receiver) {
            pagee++
            messageElement.className = "message received";
        } else {
            alert(`${data.Sender} sent you a message`)
            return
        }
        messageElement.innerHTML = `
        <div class="header">
        <span style="margin-right:20px;color:black;font-weight: 700;" class="username">${data.Sender}</span>
        <span style="margin-left:auto;" class="timestamp">${data.created_at}</span>
    </div>
    <div style="margin-top:15px; text-align: left;" class="content">
        <span>${data.msg}</span>
    </div>
        `
        chatMessages.appendChild(messageElement);

        chatInput.value = "";
        chatMessages.scrollTop = chatMessages.scrollHeight;
    }else if (data.Sender){
        alert(`${data.Sender} sent you a message`)
    }
}

function sendMessage(uname) {
    let message = document.querySelector('#chatInput').value

    if (message == "" || message.length > 100){
        alert('Check your message and retry!')
        return
    }
    
    worker.port.postMessage(JSON.stringify({
        Receiver: uname,
        Msg: message,
    }))
}


function handleScroll(receiver) {
    const chatMessages = document.getElementById("chatMessages");
    if (chatMessages.scrollTop == 0) {
        pagee += 10
        trchatbox(receiver, 100)
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


function throttle(fn, delay) {
    let last = 0;
    return function () {
        const now = +new Date();
        if (now - last > delay) {
            fn.apply(this, arguments);
            last = now;
        }
    };
}

const trchatbox = debounce(getChatBox, 2000)
const trsendMessage = throttle(sendMessage,1000)

async function refetchLogin(request) {
    fetch(request,{
        headers : {
            'request':'refetch',
        },
    }).then(resp => resp.text())
        .then(html => {
            document.documentElement.innerHTML = html
        })
}



function logout() {
    worker.port.postMessage('kill')
    fetch('/logout', {
        method: 'POST',
    })
        .then(async response => {
            if (response.status === 200) {
                refetchLogin('/login')
            }
        })
        .catch(() => {
            writeError(logerror, "red", 'Network error, please try again later!', 1500);
        });
}
