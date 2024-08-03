package forum

import (
    "database/sql"
    "fmt"
    "net/http"
    "strings"
	"time"
)

func GetChatAreaHTML(db *sql.DB, r *http.Request) (string, error) {
    users, err := GetAllUsersStatus(db)
    if err != nil {
        return "", err
    }

    var chatHTML strings.Builder
    chatHTML.WriteString(`<div class="chat-area"><h3>Users</h3><ul>`)
    for _, user := range users {
        status := "offline"
        if HasActiveSession(db, user.UserID) {
            status = "online"
        }
        chatHTML.WriteString(fmt.Sprintf(`<li>%s - <span class="%s">%s</span></li>`, user.NickName, status, status))
    }
    chatHTML.WriteString(`</ul></div>`)

    return chatHTML.String(), nil
}

func HasActiveSession(db *sql.DB, userID int) bool {
    var count int
    err := db.QueryRow("SELECT COUNT(*) FROM sessions WHERE user_id = ? AND expires_at > ?", userID, time.Now()).Scan(&count)
    if err != nil {
        return false
    }
    return count > 0
}






func GetAllUsersStatus(db *sql.DB) ([]User, error) {
    query := `SELECT user_id, nickname FROM users`
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var u User
        if err := rows.Scan(&u.UserID, &u.NickName); err != nil {
            return nil, err
        }
        users = append(users, u)
    }

    return users, nil
}