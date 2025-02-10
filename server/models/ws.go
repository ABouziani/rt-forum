package models

import (
	"database/sql"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

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

type Status struct {
	Status string
	Uname  string
}

var (
	Upgrader  = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	Clients   = make(map[string]*websocket.Conn)
	OnlineCh  = make(chan OnlineUsers)
	MessageCh = make(chan Message)
	Mu        sync.Mutex
)

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


func StoreMsg(db *sql.DB, sender, receiver, msg string) error {
	query := `INSERT INTO messages (sender,receiver,msg,created_at) VALUES (?,?,?,?)`

	_, err := db.Exec(query, sender, receiver, msg, time.Now().Format("02-01-2006 15:04:05"))
	if err != nil {
		return err
	}

	return nil
}

func FetchdbMessages(db *sql.DB, sender, receiver string, page float64) ([]Message, error) {
	rows, er := db.Query("SELECT sender,receiver,msg,created_at FROM messages WHERE (sender = ? AND receiver = ?) OR (receiver = ? AND sender = ?) ORDER BY created_at DESC LIMIT 10 OFFSET ?;", sender, receiver, sender, receiver, page)
	if er != nil {
		return nil,er
	}
	var msgs []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.Sender, &msg.Receiver, &msg.Msg, &msg.Created_at)
		if err != nil {
			return nil,err
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

func FetchRelated(db *sql.DB, username string) (error,[]Status) {
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
		return err,nil
	}
	relUsers := []Status{}

	for rows.Next() {
		var user Status
		v := ""
		err := rows.Scan(&user.Uname, &v)
		if err != nil {
			return err,nil
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
		return err ,nil
	}
	noRelUsers := []Status{}
	for rows.Next() {
		var user Status
		err := rows.Scan(&user.Uname)
		if err != nil {
			return err ,nil
		}
		noRelUsers = append(noRelUsers, user)
	}

	sort.Slice(noRelUsers, func(i, j int) bool {
		return strings.Compare(noRelUsers[i].Uname, noRelUsers[j].Uname) == -1
	})

	relUsers = append(relUsers, noRelUsers...)

	return nil,relUsers
}
