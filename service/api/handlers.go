package api

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	"regexp"

	"github.com/google/uuid"
	"nogram/service/database"
	"nogram/cmd/auth"
)

func HandleSession(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SessionRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	if matched, _ := regexp.MatchString(`^[A-Za-z0-9_]{3,16}$`, req.Username); !matched {
        http.Error(response, "username does not satisfy constraints", http.StatusBadRequest)
        return
    }

	var existingToken string
	err := database.DB.QueryRow("SELECT token FROM Users WHERE username = ?", req.Username).Scan(&existingToken)

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

	_, err = database.DB.Exec("INSERT INTO Users (username, token, photo) VALUES (?, ?, ?)", req.Username, token, photoURL)
	if err != nil {
		http.Error(response, "Error creating user", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(SessionResponse{
		Token: token,
	})
}

func HandleUsers(response http.ResponseWriter, request *http.Request) {
	rows, err := database.DB.Query("SELECT id, username FROM Users")
	if err != nil {
		http.Error(response, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			http.Error(response, "error scanning users", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(users)
}

func HandleDialogs(response http.ResponseWriter, request *http.Request) {
	companion1Id, err := auth.GetUserId(request)
	if err != nil {
		http.Error(response, err.Error(), http.StatusUnauthorized)
		return
	}

	companionUsername := request.URL.Query().Get("companion")
	var companion2Id int
	err = database.DB.QueryRow("SELECT id FROM Users WHERE username = ?", companionUsername).Scan(&companion2Id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(response, "companion not found", http.StatusNotFound)
			return
		}
		http.Error(response, "database error", http.StatusInternalServerError)
		return
	}

	if companion1Id == companion2Id {
		http.Error(response, "you can not add yourself", http.StatusConflict)
		return
	}

	var exists bool
	err = database.DB.QueryRow(`
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

	result, err := database.DB.Exec("INSERT INTO Dialogs (companion1Id, companion2Id) VALUES (?, ?)", companion1Id, companion2Id)
	if err != nil {
		http.Error(response, "error creating dialog", http.StatusInternalServerError)
		return
	}

	dialogId, err := result.LastInsertId()
	if err != nil {
		http.Error(response, "error getting dialog id", http.StatusInternalServerError)
		return
	}

	_, err = database.DB.Exec("INSERT INTO Conversations (dialog_id) VALUES (?)", dialogId)
	if err != nil {
		http.Error(response, "error creating conversation", http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusCreated)
}

func HandleCompanions(response http.ResponseWriter, request *http.Request) {
	userId, err := auth.GetUserId(request)
	if err != nil {
		http.Error(response, err.Error(), http.StatusUnauthorized)
		return
	}

	rows1, err := database.DB.Query("SELECT companion2Id FROM Dialogs WHERE companion1Id = ?", userId)
	if err != nil {
		http.Error(response, "database error", http.StatusInternalServerError)
		return
	}
	defer rows1.Close()

	rows2, err := database.DB.Query("SELECT companion1Id FROM Dialogs WHERE companion2Id = ?", userId)
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
		err := database.DB.QueryRow("SELECT id, username FROM Users WHERE id = ?", id).Scan(&companion.ID, &companion.Username)
		if err != nil {
			http.Error(response, "error getting companion info", http.StatusInternalServerError)
			return
		}
		companions = append(companions, companion)
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(companions)
}

func HandleConversations(response http.ResponseWriter, request *http.Request) {
	myId, err := auth.GetUserId(request)
	if err != nil {
		http.Error(response, err.Error(), http.StatusUnauthorized)
		return
	}

	companionId := request.URL.Query().Get("companionId")

	var dialogId int
	err = database.DB.QueryRow(`
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
	err = database.DB.QueryRow("SELECT id FROM Conversations WHERE dialog_id = ?", dialogId).Scan(&conversationId)
	if err != nil {
		http.Error(response, "error getting conversation", http.StatusInternalServerError)
		return
	}

	var conversationPhoto string
	err = database.DB.QueryRow("SELECT photo FROM Users WHERE id = ?", companionId).Scan(&conversationPhoto)
	if err != nil {
		http.Error(response, "error getting photo", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(map[string]interface{}{
		"conversationId":    conversationId,
		"conversationPhoto": conversationPhoto,
	})
}

func HandleMessages(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		HandlegetMessages(response, request)
	case http.MethodPost:
		HandlepostMessages(response, request)
	case http.MethodPut:
		HandleputMessages(response, request)
	case http.MethodDelete:
		HandledeleteMessages(response, request)
	default:
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandlegetMessages(response http.ResponseWriter, request *http.Request) {
	currentUserId, err := auth.GetUserId(request)
	if err != nil {
		http.Error(response, err.Error(), http.StatusUnauthorized)
		return
	}

	conversationIdStr := request.URL.Query().Get("conversationId")
	conversationId, err := strconv.Atoi(conversationIdStr)
	if err != nil {
		http.Error(response, "invalid conversation id", http.StatusBadRequest)
		return
	}

	rows, err := database.DB.Query(`
        SELECT id, sender_id, content, time, reaction, checkmark, photo
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

	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(response, "error starting transaction", http.StatusInternalServerError)
		return
	}

	messages := []Message{}
	for rows.Next() {
		var msg Message
		var timeStr string
		var reaction sql.NullString
		var content sql.NullString
		var photo sql.NullString

		if err := rows.Scan(&msg.Id, &msg.SenderId, &content, &timeStr, &reaction, &msg.Checkmark, &photo); err != nil {
			tx.Rollback()
			log.Printf("Scan error: %v", err)
			http.Error(response, "error scanning messages", http.StatusInternalServerError)
			return
		}

		if content.Valid {
			msg.Content = content.String
		} else {
			msg.Content = ""
		}

		if photo.Valid {
			msg.Photo = photo.String
		} else {
			msg.Photo = ""
		}

		if reaction.Valid {
			msg.Reaction = reaction.String
		} else {
			msg.Reaction = ""
		}

		if msg.SenderId != currentUserId && msg.Checkmark == "✓" {
			_, err = tx.Exec("UPDATE Messages SET checkmark = '✓✓' WHERE id = ?", msg.Id)
			if err != nil {
				tx.Rollback()
				log.Printf("Error updating checkmark: %v", err)
				http.Error(response, "error when trying to change checkmark field", http.StatusInternalServerError)
				return
			}
			msg.Checkmark = "✓✓"
		}

		t, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			log.Printf("Time parse error: %v", err)
			http.Error(response, "error parsing time", http.StatusInternalServerError)
			return
		}
		msg.Time = t.Format(time.RFC3339)

		err = tx.QueryRow("SELECT username FROM Users WHERE id = ?", msg.SenderId).Scan(&msg.SenderName)
		if err != nil {
			log.Printf("Username query error: %v", err)
			http.Error(response, "error getting sender name", http.StatusInternalServerError)
			return
		}
		messages = append(messages, msg)
	}

	if err = tx.Commit(); err != nil {
		http.Error(response, "error committing transaction", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(messages)
}

func HandlepostMessages(response http.ResponseWriter, request *http.Request) {
	senderId, err := auth.GetUserId(request)
	if err != nil {
		http.Error(response, err.Error(), http.StatusUnauthorized)
		return
	}

	var messageRequest MessageRequest
	if err := json.NewDecoder(request.Body).Decode(&messageRequest); err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec(`
        INSERT INTO Messages (sender_id, conversation_id, content, photo) 
        VALUES (?, ?, ?, ?)`,
		senderId, messageRequest.ConversationId, messageRequest.Content, messageRequest.Photo)
	if err != nil {
		http.Error(response, "error creating message", http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusCreated)
}

func HandleputMessages(response http.ResponseWriter, request *http.Request) {
	myId, err := auth.GetUserId(request)
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
	err = database.DB.QueryRow("SELECT content, photo FROM Messages WHERE id = ?",
		messageRequest.MessageId).Scan(&originalMessage.Content, &originalMessage.Photo)
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(response, "error finding message", http.StatusInternalServerError)
		return
	}

	_, err = database.DB.Exec("INSERT INTO Messages (sender_id, conversation_id, content, photo) VALUES (?, ?, ?, ?)",
		myId, messageRequest.ConversationId, "🪁Redirected🪁 "+messageRequest.SenderName+": "+originalMessage.Content, messageRequest.Photo)
	if err != nil {
		http.Error(response, "error creating message", http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusCreated)
}

func HandledeleteMessages(response http.ResponseWriter, request *http.Request) {
	messageId := request.URL.Query().Get("messageId")

	_, err := database.DB.Exec("DELETE FROM Messages WHERE id = ?", messageId)
	if err != nil {
		http.Error(response, "error deleting message", http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
}

func HandleGroups(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		HandlepostGroups(response, request)
	case http.MethodGet:
		HandlegetGroups(response, request)
	default:
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandlepostGroups(response http.ResponseWriter, request *http.Request) {
	userId, err := auth.GetUserId(request)
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

	tx, err := database.DB.Begin()
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

func HandlegetGroups(response http.ResponseWriter, request *http.Request) {
	userId, err := auth.GetUserId(request)
	if err != nil {
		http.Error(response, err.Error(), http.StatusUnauthorized)
		return
	}

	rows, err := database.DB.Query("SELECT group_id FROM GroupsMembers WHERE user_id = ?", userId)
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
		err := database.DB.QueryRow("SELECT id, groupname, photo FROM `Groups` WHERE id = ?", id).
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

func HandleGroupsPhoto(response http.ResponseWriter, request *http.Request) {
	err := request.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(response, "error parsing form", http.StatusBadRequest)
		return
	}

	conversationId := request.FormValue("conversationId")

	file, fileHeader, err := request.FormFile("photo")
	if err != nil {
		http.Error(response, "error getting photo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes := make([]byte, fileHeader.Size)
	if _, err = file.Read(fileBytes); err != nil {
		http.Error(response, "error reading file", http.StatusInternalServerError)
		return
	}

	photoBase64 := base64.StdEncoding.EncodeToString(fileBytes)
	photoURL := fmt.Sprintf("data:%s;base64,%s", fileHeader.Header.Get("Content-Type"), photoBase64)

	var groupId int
	err = database.DB.QueryRow("SELECT group_id FROM Conversations WHERE id = ?", conversationId).Scan(&groupId)
	if err != nil {
		http.Error(response, "error getting group id", http.StatusInternalServerError)
		return
	}

	_, err = database.DB.Exec("UPDATE `Groups` SET photo = ? WHERE id = ?", photoURL, groupId)
	if err != nil {
		http.Error(response, "error updating group photo", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
    json.NewEncoder(response).Encode(map[string]string{
        "photo": photoURL,
    })
}

func HandleCandidates(response http.ResponseWriter, request *http.Request) {
	var candidateRequest CandidateRequest
	if err := json.NewDecoder(request.Body).Decode(&candidateRequest); err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	var groupId int
	err := database.DB.QueryRow("SELECT group_id FROM Conversations WHERE id = ?",
		candidateRequest.ConversationId).Scan(&groupId)
	if err != nil {
		http.Error(response, "error getting group id", http.StatusInternalServerError)
		return
	}

	_, err = database.DB.Exec("INSERT INTO GroupsMembers (group_id, user_id) VALUES (?, ?)",
		groupId, candidateRequest.UserId)
	if err != nil {
		http.Error(response, "error adding member", http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(map[string]int{
		"groupId": groupId,
	})
}

func HandleMembers(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		HandlegetMembers(response, request)
	case http.MethodDelete:
		HandledeleteMembers(response, request)
	default:
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandlegetMembers(response http.ResponseWriter, request *http.Request) {
	groupId := request.URL.Query().Get("groupId")

	rows, err := database.DB.Query("SELECT user_id FROM GroupsMembers WHERE group_id = ?", groupId)
	if err != nil {
		http.Error(response, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var members []Member
	for rows.Next() {
		var userId int
		if err := rows.Scan(&userId); err != nil {
			http.Error(response, "error scanning user ids", http.StatusInternalServerError)
			return
		}

		var member Member
		err := database.DB.QueryRow("SELECT id, username FROM Users WHERE id = ?", userId).
			Scan(&member.ID, &member.Username)
		if err != nil {
			http.Error(response, "error getting member info", http.StatusInternalServerError)
			return
		}
		members = append(members, member)
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(members)
}

func HandledeleteMembers(response http.ResponseWriter, request *http.Request) {
	userId, err := auth.GetUserId(request)
	if err != nil {
		http.Error(response, err.Error(), http.StatusUnauthorized)
		return
	}

	conversationId := request.URL.Query().Get("conversationId")

	var groupId int
	err = database.DB.QueryRow("SELECT group_id FROM Conversations WHERE id = ?",
		conversationId).Scan(&groupId)
	if err != nil {
		http.Error(response, "error getting group id", http.StatusInternalServerError)
		return
	}

	_, err = database.DB.Exec("DELETE FROM GroupsMembers WHERE group_id = ? AND user_id = ?", groupId, userId)
	if err != nil {
		http.Error(response, "error deleting members", http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusOK)
}

func HandleConversationsGroups(response http.ResponseWriter, request *http.Request) {
	groupId := request.URL.Query().Get("groupId")

	var conversationPhoto string
	err := database.DB.QueryRow("SELECT photo FROM `Groups` WHERE id = ?", groupId).Scan(&conversationPhoto)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(response, "group not found", http.StatusNotFound)
			return
		}
		http.Error(response, "database error", http.StatusInternalServerError)
		return
	}

	var conversationId int
	err = database.DB.QueryRow("SELECT id FROM Conversations WHERE group_id = ?", groupId).Scan(&conversationId)
	if err != nil {
		http.Error(response, "error getting conversation", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(map[string]interface{}{
		"conversationId":    conversationId,
		"conversationPhoto": conversationPhoto,
	})
}

func HandleReactions(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPut:
		HandleputReactions(response, request)
	case http.MethodDelete:
		HandledeleteReactions(response, request)
	default:
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleputReactions(response http.ResponseWriter, request *http.Request) {
	var reactionRequest ReactionRequest
	if err := json.NewDecoder(request.Body).Decode(&reactionRequest); err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec("UPDATE Messages SET reaction = ? WHERE id = ?",
		reactionRequest.Reaction, reactionRequest.MessageId)
	if err != nil {
		http.Error(response, "error updating reaction", http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
}

func HandledeleteReactions(response http.ResponseWriter, request *http.Request) {
	messageId := request.URL.Query().Get("messageId")

	_, err := database.DB.Exec("UPDATE Messages SET reaction = NULL WHERE id = ?",
		messageId)
	if err != nil {
		http.Error(response, "error updating reaction", http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
}

func HandleLatestMessages(response http.ResponseWriter, request *http.Request) {
	userId, err := auth.GetUserId(request)
	if err != nil {
		http.Error(response, err.Error(), http.StatusUnauthorized)
		return
	}

	var companions []Companion
	if err := json.NewDecoder(request.Body).Decode(&companions); err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	latestMessages := []LatestMessage{}

	for _, companion := range companions {
		var conversationId int
		err := database.DB.QueryRow(`
            SELECT conversation.id 
            FROM Conversations conversation
            JOIN Dialogs dialog ON conversation.dialog_id = dialog.id
            WHERE (dialog.companion1Id = ? AND dialog.companion2Id = ?)
               OR (dialog.companion1Id = ? AND dialog.companion2Id = ?)`,
			userId, companion.ID,
			companion.ID, userId,
		).Scan(&conversationId)

		if err != nil {
			http.Error(response, "conversation not found", http.StatusNotFound)
			return
		}

		var latestMsg LatestMessage
		var timeStr string
		err = database.DB.QueryRow(`
            SELECT content, time
            FROM Messages
            WHERE conversation_id = ?
            ORDER BY time DESC
            LIMIT 1`,
			conversationId,
		).Scan(&latestMsg.Content, &timeStr)

		if err != nil {
			continue
		}

		t, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			http.Error(response, "error parsing time", http.StatusInternalServerError)
			return
		}
		latestMsg.Time = t.Format(time.RFC3339)
		latestMsg.CompanionId = companion.ID

		latestMessages = append(latestMessages, latestMsg)
	}

	sort.Slice(latestMessages, func(i, j int) bool {
		return latestMessages[i].Time > latestMessages[j].Time
	})

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(latestMessages)
}

func HandleLatestGroupMessages(response http.ResponseWriter, request *http.Request) {
	var groups []GroupResponse
	if err := json.NewDecoder(request.Body).Decode(&groups); err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	latestMessages := []LatestGroupMessage{}

	for _, group := range groups {
		var conversationId int
		err := database.DB.QueryRow("SELECT id FROM Conversations WHERE group_id = ?",
			group.ID).Scan(&conversationId)
		if err != nil {
			continue
		}

		var latestMsg LatestGroupMessage
		var timeStr string
		err = database.DB.QueryRow(`
            SELECT content, time
            FROM Messages
            WHERE conversation_id = ?
            ORDER BY time DESC
            LIMIT 1`,
			conversationId,
		).Scan(&latestMsg.Content, &timeStr)

		if err != nil {
			continue
		}

		t, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			continue
		}
		latestMsg.Time = t.Format(time.RFC3339)
		latestMsg.GroupId = group.ID

		latestMessages = append(latestMessages, latestMsg)
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(latestMessages)
}

func HandleProfileUsername(response http.ResponseWriter, request *http.Request) {
	userId, err := auth.GetUserId(request)

	if err != nil {
		http.Error(response, err.Error(), http.StatusUnauthorized)
		return
	}

	var usernameRequest UsernameRequest
	if err := json.NewDecoder(request.Body).Decode(&usernameRequest); err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("UPDATE Users SET username = ? WHERE id = ?", usernameRequest.Username, userId)
	if err != nil {
		http.Error(response, "error updating username", http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
}

func HandleProfilePhoto(response http.ResponseWriter, request *http.Request) {
	// Ensures only PUT requests are accepted for photo updates
	if request.Method != http.MethodPut {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from the authentication token
	userId, err := auth.GetUserId(request)
	if err != nil {
		http.Error(response, err.Error(), http.StatusUnauthorized)
		return
	}

	// Parse the multipart form data, limiting file size to 10MB
	err = request.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(response, "error parsing multipart form", http.StatusBadRequest)
		return
	}

	// Get the file from the form field named "photo"
	file, fileHeader, err := request.FormFile("photo")
	if err != nil {
		http.Error(response, "error getting photo from form: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Ensure file is closed after function completes
	defer file.Close()

	// Read the entire file into memory as bytes
	fileBytes := make([]byte, fileHeader.Size)
	if _, err = file.Read(fileBytes); err != nil {
		http.Error(response, "error reading file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify that the uploaded file is an image
	contentType := fileHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		http.Error(response, "file type not supported", http.StatusBadRequest)
		return
	}

	// Convert the image to a base64 string for storage
	photoBase64 := base64.StdEncoding.EncodeToString(fileBytes)
	// Create a data URL that can be used directly in HTML <img> tags
	photoURL := fmt.Sprintf("data:%s;base64,%s", contentType, photoBase64)

	// Update the user's photo in the database
	result, err := database.DB.Exec("UPDATE Users SET photo = ? WHERE id = ?", photoURL, userId)
	if err != nil {
		http.Error(response, "error updating photo in database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify that the update was successful and affected exactly one row
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(response, "no user found or no update made", http.StatusInternalServerError)
		return
	}

	// Return success status
	response.WriteHeader(http.StatusOK)
}

func HandleProfile(response http.ResponseWriter, request *http.Request) {
	userId, err := auth.GetUserId(request)
	if err != nil {
		http.Error(response, err.Error(), http.StatusUnauthorized)
		return
	}

	var profile ProfileResponse
	err = database.DB.QueryRow("SELECT username, photo FROM Users WHERE id = ?", userId).
		Scan(&profile.Username, &profile.Photo)
	if err != nil {
		http.Error(response, "error getting profile", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(profile)
}

func HandleMessagesPhoto(response http.ResponseWriter, request *http.Request) {
	// Get user ID from token
	userId, err := auth.GetUserId(request)
	if err != nil {
		http.Error(response, err.Error(), http.StatusUnauthorized)
		return
	}

	// Parse multipart form
	err = request.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(response, "error parsing form", http.StatusBadRequest)
		return
	}

	// Get conversation ID
	conversationId := request.FormValue("conversationId")
	content := request.FormValue("content")

	// Get the file
	file, fileHeader, err := request.FormFile("photo")
	if err != nil {
		http.Error(response, "error getting photo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file bytes
	fileBytes := make([]byte, fileHeader.Size)
	if _, err = file.Read(fileBytes); err != nil {
		http.Error(response, "error reading file", http.StatusInternalServerError)
		return
	}

	// Convert to base64
	photoBase64 := base64.StdEncoding.EncodeToString(fileBytes)
	photoURL := fmt.Sprintf("data:%s;base64,%s", fileHeader.Header.Get("Content-Type"), photoBase64)

	// Insert into database
	_, err = database.DB.Exec(
		"INSERT INTO Messages (sender_id, conversation_id, photo, content) VALUES (?, ?, ?, ?)",
		userId, conversationId, photoURL, content,
	)
	if err != nil {
		http.Error(response, "error saving message", http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusCreated)
}

func HandleGroupsGroupName(response http.ResponseWriter, request *http.Request) {
	var groupNameRequest GroupNameRequest
	if err := json.NewDecoder(request.Body).Decode(&groupNameRequest); err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
 
	var groupId int
	err := database.DB.QueryRow("SELECT group_id FROM Conversations WHERE id = ?", 
		groupNameRequest.ConversationId).Scan(&groupId)
	if err != nil {
		http.Error(response, "error getting group id", http.StatusInternalServerError)
		return
	}
 
	_, err = database.DB.Exec("UPDATE `Groups` SET groupname = ? WHERE id = ?", 
		groupNameRequest.Groupname, groupId)
	if err != nil {
		http.Error(response, "error updating group name", http.StatusInternalServerError)
		return
	}
 
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(map[string]string{
		"groupname": groupNameRequest.Groupname,
	})
}