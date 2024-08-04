package forum

import (
    "github.com/gorilla/websocket"
    "net/http"
    "database/sql"
    "log"
    "strings"
	"strconv"
    "fmt"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type Client struct {
    conn *websocket.Conn
    userID int
}

var clients = make(map[*Client]bool)

func HandleWebSocket(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }

    session, err := GetSession(r, db)
    if err != nil {
        log.Println(err)
        conn.Close()
        return
    }

    client := &Client{conn: conn, userID: session.UserID}
    clients[client] = true

    defer func() {
        delete(clients, client)
        conn.Close()
    }()

	for {
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }

        parts := strings.SplitN(string(p), ":", 2)
        if len(parts) != 2 {
            continue
        }

        receiverID, err := strconv.Atoi(parts[0])
        if err != nil {
            log.Println("Invalid receiver ID:", err)
            continue
        }
        content := parts[1]

        SaveMessage(db, client.userID, receiverID, content)

        for c := range clients {
            if c.userID == receiverID {
                err := c.conn.WriteMessage(messageType, []byte(fmt.Sprintf("%d:%s", client.userID, content)))
                if err != nil {
                    log.Println("Error sending message:", err)
                }
                break
            }
        }
    }
}

func SaveMessage(db *sql.DB, senderID int, receiverID int, content string) {
    _, err := db.Exec("INSERT INTO messages (sender_id, receiver_id, content) VALUES (?, ?, ?)",
        senderID, receiverID, content)
    if err != nil {
        log.Println(err)
    }
}