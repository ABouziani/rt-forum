package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"forum/server/models"
	"forum/server/utils"

	"github.com/gorilla/websocket"
)

func HandleWS(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	_, username, valid := models.ValidSession(r, db)
	if !valid {
		w.WriteHeader(401)
		return
	}
	var err error
	var ws *websocket.Conn
	ws, err = models.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	models.Mu.Lock()
	models.Clients[username] = ws
	models.Mu.Unlock()
	err = Broadcast(db)
	if err != nil {
		return
	}

	for {
		var receivedMsg models.Message
		err = ws.ReadJSON(&receivedMsg)
		if !checkUser(db, receivedMsg.Receiver) && receivedMsg.Receiver != "" {
			ws.WriteMessage(1, []byte("bad request!"))
			continue
		}
		if err != nil {
			handleDisconnect(username, db, valid, w, r)
			break
		}
		receivedMsg.Sender = username
		if receivedMsg.Type == "typpin" || receivedMsg.Type == "stoptypping" {
			err = SendMessage(receivedMsg.Sender, receivedMsg.Receiver, receivedMsg)
			if err != nil {
				handleDisconnect(username, db, valid, w, r)
				break
			}
			continue
		}
		if strings.TrimSpace(receivedMsg.Msg) == "" || len(strings.TrimSpace(receivedMsg.Msg)) > 100 {
			log.Println("Invalid message")
			continue
		}

		err = models.StoreMsg(db, receivedMsg.Sender, receivedMsg.Receiver, receivedMsg.Msg)
		if err != nil {
			handleDisconnect(username,db,valid,w,r)
			break
		}
		receivedMsg.Created_at = time.Now().Format("02-01-06 15:04:05")
		err = SendMessage(receivedMsg.Sender, receivedMsg.Receiver, receivedMsg)
		if err != nil {
			handleDisconnect(username,db,valid,w,r)
			break
		}
		err = Broadcast(db)
		if err != nil {
			break
		}
	}
}

func Broadcast(db *sql.DB) error {
	models.Mu.Lock()
	for uname, client := range models.Clients {
		err, relUsers := models.FetchRelated(db, uname)
		if err != nil {
			return err
		}
		for i := 0; i < len(relUsers); i++ {
			if _, exist := models.Clients[relUsers[i].Uname]; exist {
				relUsers[i].Status = "online"
			} else {
				relUsers[i].Status = "offline"
			}
		}

		jsonData := struct {
			Users []models.Status
		}{Users: relUsers}
		err = client.WriteJSON(jsonData)
		if err != nil {
			client.Close()
			delete(models.Clients, uname)
			return err
		}

	}
	models.Mu.Unlock()
	return nil
}

func SendMessage(sender, receiver string, data models.Message) error {
	_, exist := models.Clients[sender]
	_, exist2 := models.Clients[receiver]
	var err error
	if !exist {
		return fmt.Errorf("not exist")
	}
	if !exist2 {
		err = models.Clients[sender].WriteJSON(data)
		if err != nil {
			return err
		}
		return nil
	}
	if data.Type != "typpin" && data.Type != "stoptypping" {
		err = models.Clients[sender].WriteJSON(data)
		if err != nil {
			return err
		}
	}
	err = models.Clients[receiver].WriteJSON(data)
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

	var rdata models.FetchStruct
	if err := json.NewDecoder(r.Body).Decode(&rdata); err != nil {
		fmt.Println(err)
		return
	}

	msghistory, err := models.FetchdbMessages(db, sender, rdata.Receiver, rdata.Page.(float64))
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msghistory)
}

func checkUser(db *sql.DB, username string) bool {
	var exists bool
	query := "SELECT EXISTS(SELECT id FROM users WHERE username = ?)"
	db.QueryRow(query, username).Scan(&exists)
	return exists
}

func handleDisconnect(username string, db *sql.DB, valid bool, w http.ResponseWriter, r *http.Request) {
	models.Mu.Lock()
	delete(models.Clients, username)
	models.Mu.Unlock()
	Broadcast(db)
}
