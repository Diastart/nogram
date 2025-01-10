package main

import (
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	initDB()
	defer db.Close()
	http.HandleFunc("/api/profile/photo", handleProfilePhoto)
	http.HandleFunc("/api/profile/username", handleProfileUsername)
	http.HandleFunc("/api/profile", handleProfile)
	http.HandleFunc("/api/session", handleSession)
	http.HandleFunc("/api/users", handleUsers)
	http.HandleFunc("/api/dialogs", handleDialogs)
	http.HandleFunc("/api/companions", handleCompanions)
	http.HandleFunc("/api/conversations", handleConversations)
	http.HandleFunc("/api/messages", handleMessages)
	http.HandleFunc("/api/groups", handleGroups)
	http.HandleFunc("/api/members", handleMembers)
	http.HandleFunc("/api/candidates", handleCandidates)
	http.HandleFunc("/api/conversations/groups", handleConversationsGroups)
	http.HandleFunc("/api/reactions", handleReactions)
	http.HandleFunc("/api/latest/messages", handleLatestMessages)
	http.HandleFunc("/api/latest/messages/groups", handleLatestGroupMessages)
	http.HandleFunc("/api/messages/photo", handleMessagesPhoto)
	http.HandleFunc("/api/groups/photo", handleGroupsPhoto)
	http.HandleFunc("/api/groups/groupname", handleGroupsGroupName)

	http.ListenAndServe(":8080", nil)
}
