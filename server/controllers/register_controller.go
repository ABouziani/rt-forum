package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"

	"forum/server/models"
	"forum/server/utils"
)

func GetRegisterPage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Header.Get("request") != "refetch" {
		utils.RenderError(db, w, r, 404, false, "")
		return
	}
	if r.Method != http.MethodGet {
		utils.RenderError(db, w, r, http.StatusMethodNotAllowed, false, "")
		return
	}

	var valid bool
	if _, _, valid = models.ValidSession(r, db); valid {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	err := utils.RenderTemplate(db, w, r, "register", http.StatusOK, nil, false, "")
	if err != nil {
		log.Println(err)
		utils.RenderError(db, w, r, http.StatusInternalServerError, false, "")
		return
	}
}

func Signup(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		utils.RenderError(db, w, r, http.StatusMethodNotAllowed, false, "")
		return
	}

	var valid bool
	if _, _, valid = models.ValidSession(r, db); valid {
		w.WriteHeader(302)
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(400)
		return
	}

	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	passwordConfirmation := r.FormValue("password-confirmation")
	firstname := r.FormValue("firstname")
	lastname := r.FormValue("lastname")
	gender := r.FormValue("gender")
	email = strings.ToLower(strings.TrimSpace(email))
	age, er := strconv.Atoi(r.FormValue("age"))
	if er != nil || !utils.IsValidEmail(email) || !CheckData(lastname, firstname, gender, password, passwordConfirmation, username, age) {
		w.WriteHeader(400)
		return
	}

	_, err := models.StoreUser(db, email, username, password, firstname, lastname, gender, age)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.username" || err.Error() == "UNIQUE constraint failed: users.email" {
			w.WriteHeader(304)
			return
		}

		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
}

func CheckData(lastname, firstname, gender, password, passwordConfirmation, username string, age int) bool {
	if strings.Contains(username, " ") || len(strings.TrimSpace(username)) < 4 || len(strings.TrimSpace(lastname)) < 4 ||
		len(strings.TrimSpace(firstname)) < 4 || (gender != "male" && gender != "female") || len(password) < 6 {
		return false
	}

	if password != passwordConfirmation {
		return false
	}

	if age < 18 {
		return false
	}

	return true
}
