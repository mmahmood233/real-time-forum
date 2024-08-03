package forum

import (
    "net/http"
    "strconv"
    "strings"
    "time"
    "log"
    "database/sql"
    "fmt"
)

func CreateComment(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
    if r.Method != http.MethodPost {
        return fmt.Errorf("Invalid request method")
    }

    session, err := GetSession(r, db)
    if err != nil {
        return err
    }

    userID := session.UserID

    comContent := strings.TrimSpace(r.FormValue("commentCont"))
    if comContent == "" {
        return fmt.Errorf("Comment content cannot be empty")
    }

    postID := r.URL.Query().Get("postID")
    postIDInt, err := strconv.Atoi(postID)
    if err != nil {
        return fmt.Errorf("Invalid post ID: %v", err)
    }

    comment := &Comment{
        UserID:         userID,
        PostID:         postIDInt,
        CommentContent: comContent,
        CreatedAt:      time.Now(),
    }

    err = InsertComment(comment, db)
    if err != nil {
        return err
    }

    user, err := GetUserByID(db, userID)
    if err != nil {
        log.Printf("Error getting user info: %v", err)
    } else {
        log.Printf("New comment added - User: %s (ID: %d), Post ID: %d, Content: %s", user.NickName, userID, postIDInt, comContent)
    }

    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Comment added successfully"))
    return nil
}

// InsertComment inserts a new comment into the database
func InsertComment(comment *Comment, db *sql.DB) error {
    stmt, err := db.Prepare("INSERT INTO comments (user_id, post_id, comment_content, comment_created_at) VALUES (?, ?, ?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(comment.UserID, comment.PostID, comment.CommentContent, comment.CreatedAt)
    if err != nil {
        return err
    }
    return nil
}

func GetCommentsByPostID(postID int, db *sql.DB) ([]struct {
    Comment
    UserNickname string
}, error) {
    query := `
        SELECT c.comment_id, c.user_id, c.comment_content, c.comment_created_at, u.nickname,
               (SELECT COUNT(*) FROM comment_likes WHERE comment_id = c.comment_id) as like_count,
               (SELECT COUNT(*) FROM comment_dislikes WHERE comment_id = c.comment_id) as dislike_count
        FROM comments c
        JOIN users u ON c.user_id = u.user_id
        WHERE c.post_id = ?
    `
    rows, err := db.Query(query, postID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var comments []struct {
        Comment
        UserNickname string
    }

    for rows.Next() {
        var c Comment
        var userNickname string
        var createdAtStr string
        if err := rows.Scan(&c.CommentID, &c.UserID, &c.CommentContent, &createdAtStr, &userNickname, &c.LikeCount, &c.DislikeCount); err != nil {
            return nil, err
        }
        createdAt, err := time.Parse(time.RFC3339, createdAtStr)
        if err != nil {
            return nil, err
        }
        c.CreatedAt = createdAt

        comments = append(comments, struct {
            Comment
            UserNickname string
        }{c, userNickname})
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return comments, nil
}