package forum

import (
    "net/http"
    "strings"
    "time"
    "log"
    "database/sql"
    // "encoding/json"
    "fmt"
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

func GetPosts(db *sql.DB) (string, error) {
    query := `
        SELECT p.post_id, p.user_id, p.post_content, p.post_created_at, u.nickname,
               (SELECT COUNT(*) FROM post_likes WHERE post_id = p.post_id) as like_count,
               (SELECT COUNT(*) FROM post_dislikes WHERE post_id = p.post_id) as dislike_count
        FROM posts p
        JOIN users u ON p.user_id = u.user_id
        ORDER BY p.post_created_at DESC
    `

    rows, err := db.Query(query)
    if err != nil {
        return "", err
    }
    defer rows.Close()

    var postsHTML strings.Builder

    for rows.Next() {
        var p Post
        var u User
        var createdAtStr string
        if err := rows.Scan(&p.PostID, &p.UserID, &p.PostContent, &createdAtStr, &u.NickName, &p.LikeCount, &p.DislikeCount); err != nil {
            return "", err
        }
        createdAt, err := time.Parse(time.RFC3339, createdAtStr)
        if err != nil {
            return "", err
        }
        p.CreatedAt = createdAt

        categories, err := GetCategoriesByPostID(p.PostID, db)
        if err != nil {
            return "", err
        }
        if len(categories) == 0 {
            categories = append(categories, Category{CatName: "None"})
        }

        comments, err := GetCommentsByPostID(p.PostID, db)
        if err != nil {
            return "", err
        }

        postsHTML.WriteString(fmt.Sprintf(`
            <div class="post" data-id="%d">
                <h3>%s</h3>
                <p>%s</p>
                <small>Category: %s</small>
                <small>Created at: %s</small>
                <div class="comments">
        `, p.PostID, u.NickName, p.PostContent, categories[0].CatName, p.CreatedAt.Format("2006-01-02 15:04:05")))

        for _, comment := range comments {
            postsHTML.WriteString(fmt.Sprintf(`
                <div class="comment">
                    <p>%s</p>
                    <small>By %s on %s</small>
                </div>
            `, comment.CommentContent, comment.UserNickname, comment.CreatedAt.Format("2006-01-02 15:04:05")))
        }

        postsHTML.WriteString(`
                </div>
                <form class="comment-form">
                    <input type="text" name="comment" placeholder="Add a comment" required>
                    <button type="submit">Comment</button>
                </form>
            </div>
        `)
    }

    if err := rows.Err(); err != nil {
        return "", err
    }

    return postsHTML.String(), nil
}

func GetCategoriesByPostID(postID int, db *sql.DB) ([]Category, error) {
	query := `
        SELECT c.cat_id, c.cat_name
        FROM categories c
        JOIN post_categories pc ON c.cat_id = pc.category_id
        WHERE pc.post_id = ?
    `
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.CatID, &c.CatName); err != nil {
			return nil, err
		}
		c.PostID = postID
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
