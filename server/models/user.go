package models

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GetUserInfo(db *sql.DB, username string) (int, string, error) {
	var user_id int
	var hashedPassword string
	err := db.QueryRow("SELECT id,password FROM users WHERE username = ? OR email= ?", username,username).Scan(&user_id, &hashedPassword)
	if err != nil {
		return 0, "", err
	}
	return user_id, hashedPassword, nil
}

func StoreUser(db *sql.DB, email, username, password ,firstname, lastname, gender string , age int) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return -1, err
	}

	query := `INSERT INTO users (email,username,password,firstname, lastname, gender, age) VALUES (?,?,?,?,?,?,?)`
	result, err := db.Exec(query, email, username, hashedPassword,firstname, lastname, gender,age)
	if err != nil {
		return -1, fmt.Errorf("%v", err)
	}

	userID, _ := result.LastInsertId()

	return userID, nil
}
