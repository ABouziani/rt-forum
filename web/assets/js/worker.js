var ws;
let ports = [];


function connectWebSocket() {
    ws = new WebSocket("ws://localhost:8080/ws");

    ws.onmessage = (event) => {
        ports.forEach(port => {
            port.postMessage(event.data);
        })

    };

    ws.onerror = (error) => {
        console.error("WebSocket error:", error);
    };

    ws.onclose = () => {
        setTimeout(connectWebSocket, 1000);
    };
}

onconnect = (event) => {
    const port = event.ports[0];
    ports.push(port)
    port.onmessage = (msgEvent) => {
        if (msgEvent.data == 'login') {
            ports.forEach(port => {
                port.postMessage('login');
            })
            // connectWebSocket();
            // return
        }
        if (msgEvent.data == 'kill') {
            ws.close()
            ports.forEach(port => {
                port.postMessage('logout');
            })
            return
        }


        if (ws && ws.readyState === WebSocket.OPEN) {
            ws.send(msgEvent.data);
        }


    };

    port.onclose = () => {
        let ind = ports.indexOf(port)
        if (ind != -1) {
            ports.splice(ind, 1)
        }
    }
    if ((!ws || ws.readyState === WebSocket.CLOSED)) {
        connectWebSocket();
    }
};


