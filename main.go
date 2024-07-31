package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "os"

    _ "github.com/mattn/go-sqlite3"

    forum "forum/functions"
)

var database *sql.DB

func main() {
    var err error
    database, err = sql.Open("sqlite3", "./forum.db")
    if err != nil {
        log.Fatal(err)
    }
    defer database.Close()

    // Execute the schema SQL file
    err = forum.ExecuteSQLFile(database, "./schema.sql")
    if err != nil {
        log.Fatalf("Error executing SQL file: %v", err)
    }

    // Open or create the log file in append mode
    logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal("Error opening log file:", err)
    }
    defer logFile.Close()
    log.SetOutput(logFile)

    // Set up route handlers
    http.HandleFunc("/", ServeMainPage)
    http.HandleFunc("/regToLog.js", ServeJavaScript)
    http.HandleFunc("/main.js", ServeJavaScript)
    http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
        forum.HandleReg(w, r, database)
    })
    http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        forum.HandleLogin(w, r, database)
    })
    http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
        forum.Logout(w, r, database)
    })
    http.HandleFunc("/create-post", func(w http.ResponseWriter, r *http.Request) {
        forum.CreatePost(w, r, database)
    })
    http.HandleFunc("/get-posts", func(w http.ResponseWriter, r *http.Request) {
        forum.GetPosts(w, r, database)
    })
    http.HandleFunc("/add-comment", func(w http.ResponseWriter, r *http.Request) {
        forum.CreateComment(w, r, database)
    })
    

    // Start the web server
    log.Println("Starting server on :8800")
    fmt.Println("Starting server on :8800")
    err = http.ListenAndServe(":8800", nil)
    if err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

func ServeMainPage(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    http.ServeFile(w, r, "temp/main.html")
}

func ServeJavaScript(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "temp/"+r.URL.Path[1:])
}
