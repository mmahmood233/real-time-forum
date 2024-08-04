package forum

import (
    "database/sql"
    "net/http"
    "strconv"
    "fmt"
	"strings"
	"log"
)

func GetMessages(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    userID, err := strconv.Atoi(r.URL.Query().Get("userId"))
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    session, err := GetSession(r, db)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    rows, err := db.Query(`
        SELECT sender_id, receiver_id, content, created_at
        FROM messages
        WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
        ORDER BY created_at ASC
    `, session.UserID, userID, userID, session.UserID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var messagesHTML strings.Builder
    for rows.Next() {
        var senderID, receiverID int
        var content, createdAt string
        if err := rows.Scan(&senderID, &receiverID, &content, &createdAt); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        nickname := GetUserNickname(db, senderID)
        messageClass := "received"
        if senderID == session.UserID {
            messageClass = "sent"
        }
        messagesHTML.WriteString(fmt.Sprintf("<div class='message %s'><span class='sender'>%s</span>: %s <span class='timestamp'>(%s)</span></div>", messageClass, nickname, content, createdAt))
    }

    w.Header().Set("Content-Type", "text/html")
    w.Write([]byte(messagesHTML.String()))
}

func GetUserNickname(db *sql.DB, userID int) string {
    var nickname string
    err := db.QueryRow("SELECT nickname FROM users WHERE user_id = ?", userID).Scan(&nickname)
    if err != nil {
        return "Unknown"
    }
    return nickname
}

func GetChatArea(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    users, err := GetAllUsersStatus(db)
    if err != nil {
        log.Printf("Error getting users: %v", err)
        http.Error(w, "Failed to get users", http.StatusInternalServerError)
        return
    }

    var chatHTML strings.Builder
    chatHTML.WriteString("<h3>Users</h3><ul>")
    for _, user := range users {
        status := "offline"
        if HasActiveSession(db, user.UserID) {
            status = "online"
        }
        chatHTML.WriteString(fmt.Sprintf(`<li data-user-id="%d">%s - <span class="%s">%s</span></li>`, 
            user.UserID, user.NickName, status, status))
    }
    chatHTML.WriteString("</ul>")

    w.Header().Set("Content-Type", "text/html")
    w.Write([]byte(chatHTML.String()))
}