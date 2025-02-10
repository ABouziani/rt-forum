package ws

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"forum/server/models"
	"forum/server/utils"

	"github.com/gorilla/websocket"
)

type Message struct {
	Sender     string
	Receiver   string `json:"receiver"`
	Msg        string `json:"msg"`
	Created_at string `json:"created_at"`
}

type OnlineUsers struct {
	Online []string
}

type FetchStruct struct {
	Page     interface{}
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
	Broadcast(db)

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			Mu.Lock()
			delete(Clients, username)
			Mu.Unlock()
			Broadcast(db)
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
		receivedMsg.Created_at = time.Now().Format("02-01-06 15:04:05")
		SendMessage(receivedMsg.Sender, receivedMsg.Receiver, receivedMsg)
		Broadcast(db)

	}
}

func Broadcast(db *sql.DB) {
	// Get the next message from the broadcast channel

	// Send the message to all connected clients
	Mu.Lock()
	for uname, client := range Clients {
		relUsers := fetchRelated(db, uname)

		for i := 0; i < len(relUsers); i++ {
			if _, exist := Clients[relUsers[i].Uname]; exist {
				relUsers[i].Status = "online"
			} else {
				relUsers[i].Status = "offline"
			}
		}

		jsonData, err := json.Marshal(struct {
			Users []status
		}{Users: relUsers})
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

func FetchClinetNoConnect(db *sql.DB, conectClinet []string) ([]string, error) {
	placeholders := make([]string, len(conectClinet))
	newArray := make([]interface{}, len(conectClinet))
	for i := range conectClinet {
		placeholders[i] = "?"
		newArray[i] = conectClinet[i]
	}

	query := fmt.Sprintf("SELECT username FROM users WHERE username NOT IN (%s);", strings.Join(placeholders, ", "))
	rows, err := db.Query(query, newArray...)
	if err != nil {
		return nil, err
	}
	var notConetcClinet []string
	for rows.Next() {
		var user string
		err := rows.Scan(&user)
		if err != nil {
			return nil, err
		}
		notConetcClinet = append(notConetcClinet, user)
	}
	return notConetcClinet, nil
}

func SendMessage(sender, receiver string, data Message) {
	_, exist := Clients[sender]
	_, exist2 := Clients[receiver]

	if !exist {
		fmt.Println("error")
		return
	}
	if !exist2 {
		Clients[sender].WriteJSON(data)
		return
	}
	Clients[sender].WriteJSON(data)
	Clients[receiver].WriteJSON(data)
}

func StoreMsg(db *sql.DB, sender, receiver, msg string) error {
	query := `INSERT INTO messages (sender,receiver,msg,created_at) VALUES (?,?,?,?)`

	_, err := db.Exec(query, sender, receiver, msg, time.Now().Format("02-01-2006 15:04:05"))
	if err != nil {
		return err
	}

	return nil
}

func FetchMessages(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	_, sender, valid := models.ValidSession(r, db)

	if r.Method != http.MethodPost {
		utils.RenderError(db, w, r, http.StatusMethodNotAllowed, valid, sender)
		return
	}
	if !valid {
		w.WriteHeader(401)
		return
	}

	var rdata FetchStruct
	if err := json.NewDecoder(r.Body).Decode(&rdata); err != nil {
		return
	}

	msghistory, err := fetchdbMessages(db, sender, rdata.Receiver, rdata.Page.(float64))
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msghistory)
}

func fetchdbMessages(db *sql.DB, sender, receiver string, page float64) ([]Message, error) {
	// rows, _ := db.Query(`
	// SELECT sender, receiver, msg,

	// FROM messages
	// WHERE (sender = ? AND receiver = ?)
	// 		OR (receiver = ? AND sender = ?)
	// ORDER BY fcreated_at DESC
	// LIMIT 10 OFFSET ?;
	// `, sender, receiver, sender, receiver, page)

	rows, _ := db.Query("SELECT sender,receiver,msg,created_at FROM messages WHERE (sender = ? AND receiver = ?) OR (receiver = ? AND sender = ?) ORDER BY created_at DESC LIMIT 10 OFFSET ?;", sender, receiver, sender, receiver, page)
	var msgs []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.Sender, &msg.Receiver, &msg.Msg, &msg.Created_at)
		if err != nil {
			fmt.Println(err)
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

func RemoveUname(data []string, uname string) []string {
	index := slices.Index(data, uname)
	if index <= -1 {
		return data
	}
	return append(data[:index], data[index+1:]...)
}

type status struct {
	Status string
	Uname  string
}

func fetchRelated(db *sql.DB, username string) []status {
	rows, err := db.Query(`
	SELECT 
    conv.other_user, 
    m.created_at
FROM messages m
JOIN (
    SELECT 
        CASE 
            WHEN sender = ? THEN receiver 
            ELSE sender 
        END AS other_user,
        MAX(created_at) AS last_message_time
    FROM messages
    WHERE ? IN (sender, receiver)
    GROUP BY other_user
) AS conv 
ON (conv.other_user = m.sender OR conv.other_user = m.receiver) 
AND m.created_at = conv.last_message_time
ORDER BY m.created_at DESC;

	`, username, username)
	if err != nil {
		fmt.Println(err)
	}
	relUsers := []status{}

	for rows.Next() {
		var user status
		v := ""
		err := rows.Scan(&user.Uname, &v)
		if err != nil {
			fmt.Println(err)
		}
		relUsers = append(relUsers, user)
	}

	rows, err = db.Query(`
	SELECT u.username 
FROM users u
WHERE u.username != ?
AND u.username NOT IN (
    SELECT DISTINCT 
        CASE 
            WHEN m.sender = ? THEN m.receiver 
            ELSE m.sender 
        END 
    FROM messages m
    WHERE ? IN (m.sender, m.receiver)
);
	`, username, username, username)
	if err != nil {
		fmt.Println(err)
	}
	noRelUsers := []status{}
	for rows.Next() {
		var user status
		err := rows.Scan(&user.Uname)
		if err != nil {
			fmt.Println(err)
		}
		noRelUsers = append(noRelUsers, user)
	}

	sort.Slice(noRelUsers, func(i, j int) bool {
		return strings.Compare(noRelUsers[i].Uname, noRelUsers[j].Uname) == -1
	})

	relUsers = append(relUsers, noRelUsers...)

	return relUsers
}
