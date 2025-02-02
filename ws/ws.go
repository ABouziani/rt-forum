package ws

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"sync"

	"forum/server/models"

	"github.com/gorilla/websocket"
)

type Message struct {
	Sender   string
	Receiver string `json:"receiver"`
	Msg      string `json:"msg"`
}

type OnlineUsers struct {
	Online []string
}

type FetchStruct struct {
	Page interface{}
	Receiver string
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
		w.WriteHeader(401)
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
	Broadcast(username)

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			Mu.Lock()
			delete(Clients, username)
			Mu.Unlock()
			Broadcast("")
			break
		}

		var receivedMsg Message
		err = json.Unmarshal(msg, &receivedMsg)
		if err != nil {
			fmt.Println(err)
			return
		}

		receivedMsg.Sender = username

		err = StoreMsg(db, receivedMsg.Sender, receivedMsg.Receiver, receivedMsg.Msg)
		if err != nil {
			fmt.Println(err)
			return
		}

		SendMessage(receivedMsg.Sender, receivedMsg.Receiver, receivedMsg)

	}
}

func Broadcast(username string) {
	// Get the next message from the broadcast channel
	var onlines []string
	for client := range Clients {
		onlines = append(onlines, client)
	}

	

	// Send the message to all connected clients
	Mu.Lock()
	
	for uname, client := range Clients {
		var temp = make([]string,len(onlines))
		copy(temp,onlines)
		jsonData, err := json.Marshal(struct {
			Online []string
		}{Online: RemoveUname(temp,uname)})
		if err != nil {
			fmt.Println("Error serializing JSON:", err)
			return
		}
		err = client.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			fmt.Printf("Error writing to client: %v\n", err)
			client.Close()         
			delete(Clients, uname) 
		}

	}
	Mu.Unlock()
}


func SendMessage(sender, receiver string, data Message) {
	_, exist := Clients[sender]
	_, exist2 := Clients[receiver]

	if !exist || !exist2 {
		fmt.Println("error")
		return
	}

	Clients[sender].WriteJSON(data)
	Clients[receiver].WriteJSON(data)
}

func StoreMsg(db *sql.DB, sender, receiver, msg string) error {
	query := `INSERT INTO messages (sender,receiver,msg) VALUES (?,?,?)`

	_, err := db.Exec(query, sender, receiver, msg)
	if err != nil {
		return err
	}

	return nil
}

func FetchMessages(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	_, sender, valid := models.ValidSession(r, db)

	if !valid {
		w.WriteHeader(401)
		return
	}

	var rdata FetchStruct
	if err := json.NewDecoder(r.Body).Decode(&rdata); err != nil {
		fmt.Println(err)
		return
	}

	msghistory, err := fetchdbMessages(db, sender, rdata.Receiver,rdata.Page.(float64))
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msghistory)
}

func fetchdbMessages(db *sql.DB, sender, receiver string, page float64) ([]Message, error) {
	rows, _ := db.Query("SELECT sender,receiver,msg FROM messages WHERE (sender = ? AND receiver = ?) OR (receiver = ? AND sender = ?) ORDER BY created_at DESC LIMIT 10 OFFSET ?;", sender, receiver, sender, receiver,page)
	var msgs []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.Sender, &msg.Receiver, &msg.Msg)
		if err != nil {
			fmt.Println(err)
		}
		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func RemoveUname(data []string, uname string) []string {
	index := slices.Index(data,uname)
	if index <=-1{
		return data
	}
	return append(data[:index],data[index+1:]...)
}
