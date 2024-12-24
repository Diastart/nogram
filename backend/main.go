package main

import (
	"fmt"
    "database/sql"
    "log"
    "net/http"
    _ "github.com/go-sql-driver/mysql"
	"encoding/json"
    "github.com/google/uuid"
)

var db *sql.DB

type SessionRequest struct {
    Username string `json:"username"`
}

type SessionResponse struct {
    Token string `json:"token"`
}

func handleSession(response http.ResponseWriter, request *http.Request){
	if request.Method != http.MethodPost {
        http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

	var req SessionRequest
    if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
        http.Error(response, err.Error(), http.StatusBadRequest)
        return
    }

    var existingToken string
    err := db.QueryRow(`
        SELECT Tokens.token 
        FROM Users 
        JOIN Tokens ON Users.id = Tokens.user_id 
        WHERE Users.username = ?`, 
        req.Username).Scan(&existingToken)

    if err == nil {
        response.Header().Set("Content-Type", "application/json")
        json.NewEncoder(response).Encode(SessionResponse{
            Token: existingToken,
        })
        return
    }

	tx, err := db.Begin()
    if err != nil {
        http.Error(response, "Database error", http.StatusInternalServerError)
        return
    }

	result, err := tx.Exec("INSERT INTO Users (username) VALUES (?)", req.Username)
	if err != nil {
        tx.Rollback()
        http.Error(response, "Error creating user", http.StatusInternalServerError)
        return
    }

	user_id, err := result.LastInsertId()
    if err != nil {
        tx.Rollback()
        http.Error(response, "Error getting user id", http.StatusInternalServerError)
        return
    }

	token := uuid.New().String()
	_, err = tx.Exec("INSERT INTO Tokens (token, user_id) VALUES (?, ?)", token, user_id)
	if err != nil {
        tx.Rollback()
        http.Error(response, "Error creating token", http.StatusInternalServerError)
        return
    }

	if err = tx.Commit(); err != nil {
        http.Error(response, "Error committing transaction", http.StatusInternalServerError)
        return
    }

	response.Header().Set("Content-Type", "application/json")
    json.NewEncoder(response).Encode(SessionResponse{
        Token: token,
    })
}

func initDB(){
	var err error
	db, err = sql.Open("mysql", "root:Abcdefg_1234@tcp(localhost:3306)/WASADB")
	if err != nil {log.Fatal(err)}
	err = db.Ping()
    if err != nil {log.Fatal(err)}
    fmt.Println("Successfully connected to Wasa database!")
}

func main() {
	initDB()
    defer db.Close()
	http.HandleFunc("/api/session", handleSession)
	http.ListenAndServe(":8080", nil)
}