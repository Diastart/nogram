package main

import (
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	initDB()
	defer db.Close()
	http.HandleFunc("/api/health", enableCors(handleHealth))
	http.HandleFunc("/api/profile/photo", enableCors(handleProfilePhoto))
	http.HandleFunc("/api/profile/username", enableCors(handleProfileUsername))
	http.HandleFunc("/api/profile", enableCors(handleProfile))
	http.HandleFunc("/api/session", enableCors(handleSession))
	http.HandleFunc("/api/users", enableCors(handleUsers))
	http.HandleFunc("/api/dialogs", enableCors(handleDialogs))
	http.HandleFunc("/api/companions", enableCors(handleCompanions))
	http.HandleFunc("/api/conversations", enableCors(handleConversations))
	http.HandleFunc("/api/messages", enableCors(handleMessages))
	http.HandleFunc("/api/groups", enableCors(handleGroups))
	http.HandleFunc("/api/members", enableCors(handleMembers))
	http.HandleFunc("/api/candidates", enableCors(handleCandidates))
	http.HandleFunc("/api/conversations/groups", enableCors(handleConversationsGroups))
	http.HandleFunc("/api/reactions", enableCors(handleReactions))
	http.HandleFunc("/api/latest/messages", enableCors(handleLatestMessages))
	http.HandleFunc("/api/latest/messages/groups", enableCors(handleLatestGroupMessages))
	http.HandleFunc("/api/messages/photo", enableCors(handleMessagesPhoto))
	http.HandleFunc("/api/groups/photo", enableCors(handleGroupsPhoto))
	http.HandleFunc("/api/groups/groupname", enableCors(handleGroupsGroupName))

	http.ListenAndServe(":8080", nil)
}
