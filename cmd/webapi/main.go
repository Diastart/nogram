package main

import (
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"nogram/cmd/healthcheck"
	"nogram/service/database"
	"nogram/service/api"
	"nogram/cmd/auth"
)

func main() {
	database.InitDB()
	defer database.DB.Close()
	http.HandleFunc("/api/health", auth.EnableCors(healthcheck.HandleHealth))
	http.HandleFunc("/api/profile/photo", auth.EnableCors(api.HandleProfilePhoto))
	http.HandleFunc("/api/profile/username", auth.EnableCors(api.HandleProfileUsername))
	http.HandleFunc("/api/profile", auth.EnableCors(api.HandleProfile))
	http.HandleFunc("/api/session", auth.EnableCors(api.HandleSession))
	http.HandleFunc("/api/users", auth.EnableCors(api.HandleUsers))
	http.HandleFunc("/api/dialogs", auth.EnableCors(api.HandleDialogs))
	http.HandleFunc("/api/companions", auth.EnableCors(api.HandleCompanions))
	http.HandleFunc("/api/conversations", auth.EnableCors(api.HandleConversations))
	http.HandleFunc("/api/messages", auth.EnableCors(api.HandleMessages))
	http.HandleFunc("/api/groups", auth.EnableCors(api.HandleGroups))
	http.HandleFunc("/api/members", auth.EnableCors(api.HandleMembers))
	http.HandleFunc("/api/candidates", auth.EnableCors(api.HandleCandidates))
	http.HandleFunc("/api/conversations/groups", auth.EnableCors(api.HandleConversationsGroups))
	http.HandleFunc("/api/reactions", auth.EnableCors(api.HandleReactions))
	http.HandleFunc("/api/latest/messages", auth.EnableCors(api.HandleLatestMessages))
	http.HandleFunc("/api/latest/messages/groups", auth.EnableCors(api.HandleLatestGroupMessages))
	http.HandleFunc("/api/messages/photo", auth.EnableCors(api.HandleMessagesPhoto))
	http.HandleFunc("/api/groups/photo", auth.EnableCors(api.HandleGroupsPhoto))
	http.HandleFunc("/api/groups/groupname", auth.EnableCors(api.HandleGroupsGroupName))

	http.ListenAndServe(":8080", nil)
}
