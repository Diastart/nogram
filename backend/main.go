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
    "time"
    "strconv"
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

type Message struct {
    Id         int       `json:"id"`
    Content    string    `json:"content"`
    Time       string    `json:"time"`
    SenderId   int       `json:"senderId"`
    SenderName string    `json:"senderName"`
}

type MessageRequest struct {
    ConversationId int    `json:"conversationId"`
    Content        string `json:"content"`
}

type GroupRequest struct {
    GroupName    string `json:"groupName"`
    GroupMembers []int  `json:"groupMembers"`
}

type GroupResponse struct {
    ID        int    `json:"id"`
    Groupname string `json:"groupname"`
    Photo     string `json:"photo"`
}

type RedirectRequest struct {
    MessageId      int `json:"messageId"`
    ConversationId int `json:"conversationId"`
    SenderName string  `json:"senderName"`
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

func handleMessages(response http.ResponseWriter, request *http.Request) {
    switch request.Method {
    case http.MethodGet:
        handlegetMessages(response, request)
    case http.MethodPost:
        handlepostMessages(response, request)
    case http.MethodPut:
        handleputMessages(response, request)
    case http.MethodDelete:
        handledeleteMessages(response, request)
    default:
        http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func handlegetMessages(response http.ResponseWriter, request *http.Request) {
    conversationIdStr := request.URL.Query().Get("conversationId")
    conversationId, err := strconv.Atoi(conversationIdStr)
    if err != nil {
        http.Error(response, "invalid conversation id", http.StatusBadRequest)
        return
    }

    rows, err := db.Query(`
        SELECT id, sender_id, content, time
        FROM Messages 
        WHERE conversation_id = ?
        ORDER BY time ASC`,
        conversationId)
    if err != nil {
        log.Printf("Query error: %v", err)
        http.Error(response, "database error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    messages := []Message{}
    for rows.Next() {
        var msg Message
        var timeStr string
        if err := rows.Scan(&msg.Id, &msg.SenderId, &msg.Content, &timeStr); err != nil {
            log.Printf("Scan error: %v", err)
            http.Error(response, "error scanning messages", http.StatusInternalServerError)
            return
        }

        t, err := time.Parse("2006-01-02 15:04:05", timeStr)
        if err != nil {
            log.Printf("Time parse error: %v", err)
            http.Error(response, "error parsing time", http.StatusInternalServerError)
            return
        }
        msg.Time = t.Format(time.RFC3339)

        err = db.QueryRow("SELECT username FROM Users WHERE id = ?", msg.SenderId).Scan(&msg.SenderName)
        if err != nil {
            log.Printf("Username query error: %v", err)
            http.Error(response, "error getting sender name", http.StatusInternalServerError)
            return
        } 
        messages = append(messages, msg)
    }
    response.Header().Set("Content-Type", "application/json")
    json.NewEncoder(response).Encode(messages)
}

func handlepostMessages(response http.ResponseWriter, request *http.Request) {
    senderId, err := getUserId(request)
    if err != nil {
        http.Error(response, err.Error(), http.StatusUnauthorized)
        return
    }

    var messageRequest MessageRequest
    if err := json.NewDecoder(request.Body).Decode(&messageRequest); err != nil {
        http.Error(response, err.Error(), http.StatusBadRequest)
        return
    }

    _, err = db.Exec(`
        INSERT INTO Messages (sender_id, conversation_id, content) 
        VALUES (?, ?, ?)`,
        senderId, messageRequest.ConversationId, messageRequest.Content)
    if err != nil {
        http.Error(response, "error creating message", http.StatusInternalServerError)
        return
    }

    response.WriteHeader(http.StatusCreated)
}

func handleputMessages(response http.ResponseWriter, request *http.Request) {
    myId, err := getUserId(request)
    if err != nil {
        http.Error(response, err.Error(), http.StatusUnauthorized)
        return
    }

    var messageRequest RedirectRequest
    if err := json.NewDecoder(request.Body).Decode(&messageRequest); err != nil {
        http.Error(response, err.Error(), http.StatusBadRequest)
        return
    }

    var originalMessage Message 
    err = db.QueryRow("SELECT content FROM Messages WHERE id = ?", 
        messageRequest.MessageId).Scan(&originalMessage.Content)
    if err != nil {
        http.Error(response, "error finding message", http.StatusInternalServerError)
        return
    }

    _, err = db.Exec("INSERT INTO Messages (sender_id, conversation_id, content) VALUES (?, ?, ?)",
        myId, messageRequest.ConversationId, "*Redirected " + messageRequest.SenderName + "* " + originalMessage.Content)
    if err != nil {
        http.Error(response, "error creating message", http.StatusInternalServerError)
        return
    }

    response.WriteHeader(http.StatusCreated)
}

func handledeleteMessages(response http.ResponseWriter, request *http.Request) {
    messageId := request.URL.Query().Get("messageId")
   
    _, err := db.Exec("DELETE FROM Messages WHERE id = ?", messageId)
    if err != nil {
        http.Error(response, "error deleting message", http.StatusInternalServerError)
        return
    }

    response.WriteHeader(http.StatusOK)
}

func handleGroups(response http.ResponseWriter, request *http.Request) {
    switch request.Method {
    case http.MethodPost:
        handlepostGroups(response, request)
    case http.MethodGet:
        handlegetGroups(response, request)
    default:
        http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func handlepostGroups(response http.ResponseWriter, request *http.Request) {
    userId, err := getUserId(request)
    if err != nil {
        http.Error(response, err.Error(), http.StatusUnauthorized)
        return
    }

    var groupRequest GroupRequest
    if err := json.NewDecoder(request.Body).Decode(&groupRequest); err != nil {
        http.Error(response, err.Error(), http.StatusBadRequest)
        return
    }

    photoURL := fmt.Sprintf("https://api.dicebear.com/7.x/identicon/svg?seed=%s", groupRequest.GroupName)

    tx, err := db.Begin()
    if err != nil {
        http.Error(response, "database error", http.StatusInternalServerError)
        return
    }

    result, err := tx.Exec("INSERT INTO `Groups` (groupname, photo) VALUES (?, ?)", 
        groupRequest.GroupName, photoURL)
    if err != nil {
        tx.Rollback()
        http.Error(response, "error creating group", http.StatusInternalServerError)
        return
    }

    groupId, err := result.LastInsertId()
    if err != nil {
        tx.Rollback()
        http.Error(response, "error getting group id", http.StatusInternalServerError)
        return
    }

    _, err = tx.Exec("INSERT INTO GroupsMembers (group_id, user_id) VALUES (?, ?)", groupId, userId)
    if err != nil {
        tx.Rollback()
        http.Error(response, "error adding current member", http.StatusInternalServerError)
        return
    }

    for _, memberId := range groupRequest.GroupMembers {
        _, err = tx.Exec("INSERT INTO GroupsMembers (group_id, user_id) VALUES (?, ?)",
            groupId, memberId)
        if err != nil {
            tx.Rollback()
            http.Error(response, "error adding group members", http.StatusInternalServerError)
            return
        }
    }

    _, err = tx.Exec("INSERT INTO Conversations (group_id) VALUES (?)", groupId)
    if err != nil {
        tx.Rollback()
        http.Error(response, "error creating conversation", http.StatusInternalServerError)
        return
    }

    if err = tx.Commit(); err != nil {
        http.Error(response, "error committing transaction", http.StatusInternalServerError)
        return
    }

    response.WriteHeader(http.StatusCreated)
}

func handlegetGroups(response http.ResponseWriter, request *http.Request) {
    userId, err := getUserId(request)
    if err != nil {
        http.Error(response, err.Error(), http.StatusUnauthorized)
        return
    }

    rows, err := db.Query("SELECT group_id FROM GroupsMembers WHERE user_id = ?", userId)
    if err != nil {
        http.Error(response, "database error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    groupIds := []int{}
    for rows.Next() {
        var groupId int
        if err := rows.Scan(&groupId); err != nil {
            http.Error(response, "error scanning group ids", http.StatusInternalServerError)
            return
        }
        groupIds = append(groupIds, groupId)
    }

    groups := []GroupResponse{}
    for _, id := range groupIds {
        var group GroupResponse
        err := db.QueryRow("SELECT id, groupname, photo FROM `Groups` WHERE id = ?", id).
            Scan(&group.ID, &group.Groupname, &group.Photo)
        if err != nil {
            http.Error(response, "error getting group info", http.StatusInternalServerError)
            return
        }
        groups = append(groups, group)
    }

    response.Header().Set("Content-Type", "application/json")
    json.NewEncoder(response).Encode(groups)
}

func initDB(){
	var err error
	db, err = sql.Open("mysql", "root:Abcdefg_1234@tcp(localhost:3306)/WASADB")
	if err != nil {log.Fatal(err)}
	err = db.Ping()
    if err != nil {log.Fatal(err)}
    fmt.Println("Successfully connected to Wasa database!")
}

func handleConversationsGroups(response http.ResponseWriter, request *http.Request) {
    groupId := request.URL.Query().Get("groupId")

    var conversationPhoto string
    err := db.QueryRow("SELECT photo FROM `Groups` WHERE id = ?", groupId).Scan(&conversationPhoto)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(response, "group not found", http.StatusNotFound)
            return
        }
        http.Error(response, "database error", http.StatusInternalServerError)
        return
    }

    var conversationId int
    err = db.QueryRow("SELECT id FROM Conversations WHERE group_id = ?", groupId).Scan(&conversationId)
    if err != nil {
        http.Error(response, "error getting conversation", http.StatusInternalServerError)
        return
    }

    response.Header().Set("Content-Type", "application/json")
    json.NewEncoder(response).Encode(map[string]interface{}{
        "conversationId": conversationId,
        "conversationPhoto": conversationPhoto,
    })
}

func main() {
	initDB()
    defer db.Close()
	http.HandleFunc("/api/session", handleSession)
    http.HandleFunc("/api/dialogs", handleDialogs)
    http.HandleFunc("/api/companions", handleCompanions)
    http.HandleFunc("/api/conversations", handleConversations)
    http.HandleFunc("/api/messages", handleMessages)
    http.HandleFunc("/api/groups", handleGroups)
    http.HandleFunc("/api/conversations/groups", handleConversationsGroups)
	http.ListenAndServe(":8080", nil)
}