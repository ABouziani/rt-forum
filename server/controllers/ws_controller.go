package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"forum/server/models"
	"forum/server/utils"
)

func HandleWS(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	_, username, valid := models.ValidSession(r, db)

	if !valid {
		w.WriteHeader(401)
		return
	}

	ws, err := models.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	models.Mu.Lock()
	models.Clients[username] = ws
	models.Mu.Unlock()
	Broadcast(db)

	for {
		var receivedMsg models.Message
		err := ws.ReadJSON(&receivedMsg)
		if receivedMsg.Msg=="" || len(receivedMsg.Msg)>100{
			w.WriteHeader(400)
			continue
		}
		if err != nil {
			models.Mu.Lock()
			delete(models.Clients, username)
			models.Mu.Unlock()
			Broadcast(db)
			break
		}

		receivedMsg.Sender = username

		err = models.StoreMsg(db, receivedMsg.Sender, receivedMsg.Receiver, receivedMsg.Msg)
		if err != nil {
			return
		}
		receivedMsg.Created_at = time.Now().Format("02-01-06 15:04:05")
		SendMessage(receivedMsg.Sender, receivedMsg.Receiver, receivedMsg)
		Broadcast(db)

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

func SendMessage(sender, receiver string, data models.Message) {
	_, exist := models.Clients[sender]
	_, exist2 := models.Clients[receiver]

	if !exist {
		return
	}
	if !exist2 {
		models.Clients[sender].WriteJSON(data)
		return
	}
	models.Clients[sender].WriteJSON(data)
	models.Clients[receiver].WriteJSON(data)
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
