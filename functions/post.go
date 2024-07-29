package forum

import (
    "net/http"
    "strings"
    "time"
    "log"
    "database/sql"
    "encoding/json"
)

func CreatePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    sessionObj, err := GetSession(r, db)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    postContent := strings.TrimSpace(r.FormValue("postCont"))
    category := r.FormValue("catCont")

    if postContent == "" {
        http.Error(w, "No post content", http.StatusBadRequest)
        return
    }

    // Create a new Post struct
    post := &Post{
        UserID:      sessionObj.UserID,
        PostContent: postContent,
        CreatedAt:   time.Now(),
    }

    // Insert the post into the database
    lastInsertID, err := InsertPost(post, db)
    if err != nil {
        log.Printf("Error inserting post: %v", err)
        http.Error(w, "Error creating post", http.StatusInternalServerError)
        return
    }

    log.Printf("New post created - ID: %d, Content: %s, Category: %s", lastInsertID, postContent, category)

    // Insert category for the post
    postCategory := &PostCategory{
        PostID:     int(lastInsertID),
        CategoryID: getCategoryID(category, db),
    }
    err = InsertPostCategory(postCategory, db)
    if err != nil {
        log.Printf("Error inserting post category: %v", err)
        // Note: We're not returning an error here as the post has been created successfully
    }

    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Post created successfully"))
}

func InsertPost(post *Post, db *sql.DB) (int64, error) {
    stmt, err := db.Prepare("INSERT INTO posts (user_id, post_content, post_created_at) VALUES (?, ?, ?)")
    if err != nil {
        return 0, err
    }
    defer stmt.Close()

    result, err := stmt.Exec(post.UserID, post.PostContent, post.CreatedAt)
    if err != nil {
        return 0, err
    }

    return result.LastInsertId()
}

func InsertPostCategory(postCategory *PostCategory, db *sql.DB) error {
    stmt, err := db.Prepare("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(postCategory.PostID, postCategory.CategoryID)
    return err
}

func getCategoryID(categoryName string, db *sql.DB) int {
    var categoryID int
    err := db.QueryRow("SELECT cat_id FROM categories WHERE cat_name = ?", categoryName).Scan(&categoryID)
    if err != nil {
        if err == sql.ErrNoRows {
            // Category doesn't exist, create it
            stmt, err := db.Prepare("INSERT INTO categories (cat_name) VALUES (?)")
            if err != nil {
                log.Printf("Error preparing category insert statement: %v", err)
                return 0
            }
            defer stmt.Close()

            result, err := stmt.Exec(categoryName)
            if err != nil {
                log.Printf("Error inserting new category: %v", err)
                return 0
            }

            lastInsertID, err := result.LastInsertId()
            if err != nil {
                log.Printf("Error getting last insert ID for category: %v", err)
                return 0
            }

            return int(lastInsertID)
        }
        log.Printf("Error querying category: %v", err)
        return 0
    }
    return categoryID
}

func GetPosts(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Fetch posts from the database
    rows, err := db.Query(`
        SELECT p.post_id, p.user_id, p.post_content, p.post_created_at, u.nickname, c.cat_name
        FROM posts p
        JOIN users u ON p.user_id = u.user_id
        LEFT JOIN post_categories pc ON p.post_id = pc.post_id
        LEFT JOIN categories c ON pc.category_id = c.cat_id
        ORDER BY p.post_created_at DESC
    `)
    if err != nil {
        log.Printf("Error fetching posts: %v", err)
        http.Error(w, "Error fetching posts", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var posts []map[string]interface{}
    for rows.Next() {
        var post Post
        var nickname, category string
        err := rows.Scan(&post.PostID, &post.UserID, &post.PostContent, &post.CreatedAt, &nickname, &category)
        if err != nil {
            log.Printf("Error scanning post row: %v", err)
            continue
        }
        posts = append(posts, map[string]interface{}{
            "id":        post.PostID,
            "userId":    post.UserID,
            "content":   post.PostContent,
            "createdAt": post.CreatedAt,
            "category":  category,
            "author":    nickname,
        })
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}