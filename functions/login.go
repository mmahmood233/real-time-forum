package forum

import (
	"database/sql"
	"net/http"
	"log"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HandleLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := ValByEmailOrUsername(db, username)
	if err != nil {
		log.Printf("Error validating user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Create a new session
	sessionID := CreateSession(w, user.UserID, db)
	if sessionID == "" {
		log.Printf("Error creating session for user: %d", user.UserID)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Login successful")
}

func Logout(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	session, err := GetSession(r, db)
	if err != nil {
		log.Printf("Error getting session: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err = RemoveSession(db, session.SessionID)
	if err != nil {
		log.Printf("Error removing session: %v", err)
	}

	// Clear the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func IsLoggedIn(r *http.Request, db *sql.DB) bool {
	session, err := GetSession(r, db)
	if err != nil {
		return false
	}

	return session != nil && session.ExpiresAt.After(time.Now())
}

func GetUserByID(db *sql.DB, userID int) (*User, error) {
	user := &User{}
	query := `SELECT user_id, email, nickname FROM users WHERE user_id = ?`
	err := db.QueryRow(query, userID).Scan(&user.UserID, &user.Email, &user.NickName)
	if err != nil {
		return nil, err
	}
	return user, nil
}


