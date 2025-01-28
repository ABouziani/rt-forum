package ws

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"forum/server/models"
	"forum/server/utils"

	"github.com/gorilla/websocket"
)

type Message struct {
	Sender   *websocket.Conn
	Receiver *websocket.Conn
	Msg      string
}

type OnlineUsers struct {
	Online []string
}

var (
	Upgrader  = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	Clients   = make(map[string]*websocket.Conn)
	OnlineCh  = make(chan OnlineUsers)
	MessageCh = make(chan Message)
	Mu        sync.Mutex
)

func HandleWS(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	_, username, valid := models.ValidSession(r, db)

	if !valid {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	ws, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Error upgrading connection: %v\n", err)
		return
	}
	defer ws.Close()

	Mu.Lock()
	Clients[username] = ws
	Mu.Unlock()
	getOnlines()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			Mu.Lock()
			delete(Clients, username)
			Mu.Unlock()
			getOnlines()
			break
		}
		// var receivedMsg map[string]interface{}
		// err = json.Unmarshal(msg, &receivedMsg)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }

		fmt.Println(string(msg))

	}
}

func Broadcast() {
	for {
		// Get the next message from the broadcast channel
		msg := <-OnlineCh

		jsonData, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("Error serializing JSON:", err)
			return
		}

		// Send the message to all connected clients
		Mu.Lock()
		for id, client := range Clients {
			err := client.WriteMessage(websocket.TextMessage, jsonData)
			if err != nil {
				fmt.Printf("Error writing to client: %v\n", err)
				client.Close()      // Close the connection
				delete(Clients, id) // Remove the client from the map
			}
		}
		Mu.Unlock()
	}
}

func getOnlines() {
	var onlines []string
	for username := range Clients {
		onlines = append(onlines, username)
	}
	// connectedClients := strconv.Itoa(len(Clients))
	OnlineCh <- OnlineUsers{Online: onlines}
}

func SendMessage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	_, username, valid := models.ValidSession(r, db)
	if r.Method != http.MethodPost {
		utils.RenderError(db, w, r, http.StatusMethodNotAllowed, valid, username)
		return
	}
	if !valid {
		w.WriteHeader(401)
		return
	}
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(400)
		return
	}

	resp, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error reading requets body")
		return
	}

	fmt.Println("sender",username,string(resp))
}
