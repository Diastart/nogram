package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"
)

var DB *sql.DB

func InitDB() {
    var err error
    DB, err = sql.Open("sqlite3", "service/database/wasa.db")
    if err != nil {log.Fatal(err)}

    var tableExists bool
    err = DB.QueryRow("SELECT EXISTS (SELECT 1 FROM sqlite_master WHERE type = 'table' AND name = 'Users')").Scan(&tableExists)
    if err != nil {log.Fatal(err)}

    if !tableExists {
        schema, err := os.ReadFile("service/database/schema.sql")
        if err != nil {log.Fatal(err)}

        _, err = DB.Exec(string(schema))
        if err != nil {log.Fatal(err)}

        fmt.Println("Database schema initialized!")
    }
    fmt.Println("Successfully connected to SQLite database!")
}