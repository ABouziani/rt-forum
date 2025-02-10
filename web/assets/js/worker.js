var ws;
let ports = [];


function connectWebSocket() {
    if (!ws || ws.readyState === WebSocket.CLOSED) {
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
            console.log("WebSocket connection closed. Attempting to reconnect...");
            setTimeout(connectWebSocket, 1000);
        };
    }
}

onconnect = (event) => {
    const port = event.ports[0];
    ports.push(port)
    port.onmessage = (msgEvent) => {
        if (msgEvent.data == 'kill') {
            ws.close()
            return
        }
        if (ws && ws.readyState === WebSocket.OPEN) {
            ws.send(msgEvent.data);
        }
    };
    connectWebSocket();
};
