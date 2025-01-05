package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "strings"
)

var db *sql.DB

func getUserId(request *http.Request) (int, error) {
    token := strings.TrimPrefix(request.Header.Get("Authorization"), "Bearer ")
    if token == "" {
        return 0, fmt.Errorf("no authorization token")
    }

    var userId int
    err := db.QueryRow("SELECT id FROM Users WHERE token = ?", token).Scan(&userId)
    if err != nil {
        if err == sql.ErrNoRows {
            return 0, fmt.Errorf("invalid token")
        }
        return 0, fmt.Errorf("database error: %v", err)
    }

    return userId, nil
}

func initDB(){
	var err error
	db, err = sql.Open("mysql", "root:Abcdefg_1234@tcp(localhost:3306)/WASADB")
	if err != nil {log.Fatal(err)}
	err = db.Ping()
    if err != nil {log.Fatal(err)}
    fmt.Println("Successfully connected to Wasa database!")
}