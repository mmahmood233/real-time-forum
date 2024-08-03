package forum		

import (
    "database/sql"
	"log"
)

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