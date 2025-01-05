package main

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

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
}

type Member struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
}

type CandidateRequest struct {
    ConversationId int `json:"conversationId"`
    UserId         int `json:"userId"`
}