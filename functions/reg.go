package forum

import (
	"database/sql"
	"net/http"
	"strings"
	"log"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HandleReg(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		log.Printf("Invalid method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nickname := strings.TrimSpace(r.FormValue("nickname"))
	age := strings.TrimSpace(r.FormValue("age"))
	gender := strings.TrimSpace(r.FormValue("gender"))
	firstname := strings.TrimSpace(r.FormValue("firstname"))
	lastname := strings.TrimSpace(r.FormValue("lastname"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	log.Printf("Received registration request: nickname=%s, age=%s, gender=%s, firstname=%s, lastname=%s, email=%s", 
		nickname, age, gender, firstname, lastname, email)

	if nickname == "" || age == "" || gender == "" || firstname == "" || lastname == "" || email == "" || password == "" {
		log.Printf("Missing required fields")
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		log.Printf("Invalid email format: %s", email)
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "An error occurred during registration", http.StatusInternalServerError)
		return
	}

	user := &User{
		NickName:  nickname,
		Age:       age,
		Gender:    gender,
		FirstName: firstname,
		LastName:  lastname,
		Email:     email,
		Password:  string(hashedPassword),
	}

	err = InsertUser(db, user)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
			http.Error(w, "This email is already taken!", http.StatusConflict)
		} else if strings.Contains(err.Error(), "UNIQUE constraint failed: users.nickname") {
			http.Error(w, "This username is already taken!", http.StatusConflict)
		} else {
			http.Error(w, "An error occurred during registration", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("User registered successfully: %s", nickname)

	// Create a session for the new user
	sessionID := CreateSession(w, user.UserID, db)
	if sessionID == "" {
		log.Printf("Error creating session for new user: %d", user.UserID)
		http.Error(w, "An error occurred during registration", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Registration successful")
}