package forum

import (
	"database/sql"
	"errors"
	"log"
)

func InsertUser(db *sql.DB, user *User) error {
	// Check if the user already exists
	existingUser, err := ValByEmailOrUsername(db, user.Email)
	if err != nil {
		log.Printf("Error checking existing user by email: %v", err)
		return err
	}
	if existingUser != nil {
		log.Printf("User with email %s already exists", user.Email)
		return errors.New("user with this email already exists")
	}

	existingUser, err = ValByEmailOrUsername(db, user.NickName)
	if err != nil {
		log.Printf("Error checking existing user by nickname: %v", err)
		return err
	}
	if existingUser != nil {
		log.Printf("User with nickname %s already exists", user.NickName)
		return errors.New("user with this username already exists")
	}

	insertUserSQL := `INSERT INTO users(nickname, age, gender, firstname, lastname, email, password1) VALUES (?, ?, ?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertUserSQL)
	if err != nil {
		log.Printf("Error preparing SQL statement: %v", err)
		return err
	}
	defer statement.Close()

	result, err := statement.Exec(user.NickName, user.Age, user.Gender, user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		log.Printf("Error executing SQL statement: %v", err)
		return err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		return err
	}

	user.UserID = int(userID)

	log.Printf("New user registered with ID: %d", user.UserID)

	return nil
}

func ValByEmailOrUsername(db *sql.DB, input string) (*User, error) {
	user := &User{}
	query := `SELECT user_id, nickname, age, gender, firstname, lastname, email, password1 FROM users WHERE email = ? OR nickname = ?`
	err := db.QueryRow(query, input, input).Scan(&user.UserID, &user.NickName, &user.Age, &user.Gender, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // no user found
		}
		log.Printf("Error querying user: %v", err)
		return nil, err
	}
	return user, nil
}
