package forum

import (
    "encoding/json"
    "net/http"
    "time"
    "log"
    "database/sql"
)

func AddComment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var input struct {
        PostID      int    `json:"postID"`
        CommentCont string `json:"commentCont"`
    }

    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    sessionObj, err := GetSession(r, db)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    comment := &Comment{
        PostID:         input.PostID,
        UserID:         sessionObj.UserID,
        CommentContent: input.CommentCont,
        CreatedAt:      time.Now(),
    }

    lastInsertID, err := InsertComment(comment, db)
    if err != nil {
        log.Printf("Error inserting comment: %v", err)
        http.Error(w, "Error creating comment", http.StatusInternalServerError)
        return
    }

    // Respond with the created comment
    comment.CommentID = int(lastInsertID)
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(comment); err != nil {
        log.Printf("Error encoding comment response: %v", err)
        http.Error(w, "Error encoding response", http.StatusInternalServerError)
    }
}

func InsertComment(comment *Comment, db *sql.DB) (int64, error) {
    stmt, err := db.Prepare("INSERT INTO comments (post_id, user_id, comment_content, comment_created_at) VALUES (?, ?, ?, ?)")
    if err != nil {
        return 0, err
    }
    defer stmt.Close()

    result, err := stmt.Exec(comment.PostID, comment.UserID, comment.CommentContent, comment.CreatedAt)
    if err != nil {
        return 0, err
    }

    return result.LastInsertId()
}
