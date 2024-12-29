package main

import (
	"fmt"
    "database/sql"
    "log"
    "net/http"
    _ "github.com/go-sql-driver/mysql"
	"encoding/json"
    "github.com/google/uuid"
    "strings"
)

var db *sql.DB

type SessionRequest struct {
    Username string `json:"username"`
}

type SessionResponse struct {
    Token string `json:"token"`
}

type Companion struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
}

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

func handleSession(response http.ResponseWriter, request *http.Request) {
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
    err := db.QueryRow("SELECT token FROM Users WHERE username = ?", req.Username).Scan(&existingToken)

    if err != nil && err != sql.ErrNoRows {
        http.Error(response, "Database error", http.StatusInternalServerError)
        return
    }

    if err == nil {
        response.Header().Set("Content-Type", "application/json")
        json.NewEncoder(response).Encode(SessionResponse{
            Token: existingToken,
        })
        return
    }

    token := uuid.New().String()
    photoURL := fmt.Sprintf("https://api.dicebear.com/7.x/avataaars/svg?seed=%s", req.Username)

    _, err = db.Exec("INSERT INTO Users (username, token, photo) VALUES (?, ?, ?)", req.Username, token, photoURL)
    if err != nil {
        http.Error(response, "Error creating user", http.StatusInternalServerError)
        return
    }

    response.Header().Set("Content-Type", "application/json")
    json.NewEncoder(response).Encode(SessionResponse{
        Token: token,
    })
}

func handleDialogs(response http.ResponseWriter, request *http.Request) {
    companion1Id, err := getUserId(request)
    if err != nil {
        http.Error(response, err.Error(), http.StatusUnauthorized)
        return
    }
 
    companionUsername := request.URL.Query().Get("companion")
    var companion2Id int
    err = db.QueryRow("SELECT id FROM Users WHERE username = ?", companionUsername).Scan(&companion2Id)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(response, "companion not found", http.StatusNotFound)
            return
        }
        http.Error(response, "database error", http.StatusInternalServerError)
        return
    }

    var exists bool
    err = db.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM Dialogs 
            WHERE (companion1Id = ? AND companion2Id = ?) 
            OR (companion1Id = ? AND companion2Id = ?)
        )`,
        companion1Id, companion2Id,
        companion2Id, companion1Id,
    ).Scan(&exists)
    if err != nil {
        http.Error(response, "database error", http.StatusInternalServerError)
        return
    }

   if exists {
       http.Error(response, "dialog already exists", http.StatusConflict)
       return
   }

    result, err := db.Exec("INSERT INTO Dialogs (companion1Id, companion2Id) VALUES (?, ?)", companion1Id, companion2Id)
    if err != nil {
        http.Error(response, "error creating dialog", http.StatusInternalServerError)
        return
    }

    dialogId, err := result.LastInsertId()
    if err != nil {
        http.Error(response, "error getting dialog id", http.StatusInternalServerError)
        return
    }

    _, err = db.Exec("INSERT INTO Conversations (dialog_id) VALUES (?)", dialogId)
    if err != nil {
        http.Error(response, "error creating conversation", http.StatusInternalServerError)
        return
    }

    response.WriteHeader(http.StatusCreated)
}

func handleCompanions(response http.ResponseWriter, request *http.Request) {
    userId, err := getUserId(request)
    if err != nil {
        http.Error(response, err.Error(), http.StatusUnauthorized)
        return
    }

    rows1, err := db.Query("SELECT companion2Id FROM Dialogs WHERE companion1Id = ?", userId)
    if err != nil {
        http.Error(response, "database error", http.StatusInternalServerError)
        return
    }
    defer rows1.Close()

    rows2, err := db.Query("SELECT companion1Id FROM Dialogs WHERE companion2Id = ?", userId)
    if err != nil {
        http.Error(response, "database error", http.StatusInternalServerError)
        return
    }
    defer rows2.Close()

    companionIds := []int{}

    for rows1.Next() {
        var companionId int
        if err := rows1.Scan(&companionId); err != nil {
            http.Error(response, "error scanning companions", http.StatusInternalServerError)
            return
        }
        companionIds = append(companionIds, companionId)
    }

    for rows2.Next() {
        var companionId int
        if err := rows2.Scan(&companionId); err != nil {
            http.Error(response, "error scanning companions", http.StatusInternalServerError)
            return
        }
        companionIds = append(companionIds, companionId)
    }

    companions := []Companion{}

    for _, id := range companionIds {
        var companion Companion
        err := db.QueryRow("SELECT id, username FROM Users WHERE id = ?", id).Scan(&companion.ID, &companion.Username)
        if err != nil {
            http.Error(response, "error getting companion info", http.StatusInternalServerError)
            return
        }
        companions = append(companions, companion)
    }

    response.Header().Set("Content-Type", "application/json")
    json.NewEncoder(response).Encode(companions)
}

func handleConversations(response http.ResponseWriter, request *http.Request) {
    myId, err := getUserId(request)
    if err != nil {
        http.Error(response, err.Error(), http.StatusUnauthorized)
        return
    }

    companionId := request.URL.Query().Get("companionId")

    var dialogId int
    err = db.QueryRow(`
        SELECT id FROM Dialogs 
        WHERE (companion1Id = ? AND companion2Id = ?) 
        OR (companion1Id = ? AND companion2Id = ?)`,
        myId, companionId,
        companionId, myId,
    ).Scan(&dialogId)

    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(response, "dialog not found", http.StatusNotFound)
            return
        }
        http.Error(response, "database error", http.StatusInternalServerError)
        return
    }

    var conversationId int
    err = db.QueryRow("SELECT id FROM Conversations WHERE dialog_id = ?", dialogId).Scan(&conversationId)
    if err != nil {
        http.Error(response, "error getting conversation", http.StatusInternalServerError)
        return
    }

    var conversationPhoto string
    err = db.QueryRow("SELECT photo FROM Users WHERE id = ?", companionId).Scan(&conversationPhoto)
    if err != nil {
        http.Error(response, "error getting photo", http.StatusInternalServerError)
        return
    }

    response.Header().Set("Content-Type", "application/json")
    json.NewEncoder(response).Encode(map[string]interface{}{
        "conversationId": conversationId,
        "conversationPhoto": conversationPhoto,
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
    http.HandleFunc("/api/dialogs", handleDialogs)
    http.HandleFunc("/api/companions", handleCompanions)
    http.HandleFunc("/api/conversations", handleConversations)
	http.ListenAndServe(":8080", nil)
}