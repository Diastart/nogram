package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "strings"
    "os"
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

func enableCors(handler http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Max-Age", "1")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        handler(w, r)
    }
}

func initDB() {
    var err error
    db, err = sql.Open("sqlite3", "./wasa.db")
    if err != nil {log.Fatal(err)}

    var tableExists bool
    err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM sqlite_master WHERE type = 'table' AND name = 'Users')").Scan(&tableExists)
    if err != nil {log.Fatal(err)}

    if !tableExists {
        schema, err := os.ReadFile("schema.sql")
        if err != nil {log.Fatal(err)}

        _, err = db.Exec(string(schema))
        if err != nil {log.Fatal(err)}

        fmt.Println("Database schema initialized!")
    }
    fmt.Println("Successfully connected to SQLite database!")
}