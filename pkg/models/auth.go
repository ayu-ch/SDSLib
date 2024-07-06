package models

import (
	"database/sql"
	"fmt"

	"github.com/ayu-ch/SDSLib/pkg/types"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func Login(username, password string) (types.User, error) {
	var user types.User
	var hashedPassword string
	db, err := Connection()
	if err != nil {
		return user, fmt.Errorf("error connecting to the database: %s", err)
	}
	defer db.Close()

	query := "SELECT UserID, Username, Pass, Role FROM User WHERE Username = ?;"
	err = db.QueryRow(query, username).Scan(&user.UserID, &user.Username, &hashedPassword, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("invalid credentials")
		}
		return user, fmt.Errorf("error querying the database: %s", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return user, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

func Signup(username, password string) error {
	db, err := Connection()
	if err != nil {
		return fmt.Errorf("error connecting to the database: %s", err)
	}
	defer db.Close()

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return fmt.Errorf("error hashing password: %s", err)
	}

	query := "INSERT INTO User (Username, Pass, Role) VALUES (?, ?, ?)"
	_, err = db.Exec(query, username, hashedPassword, "Client")
	if err != nil {
		return fmt.Errorf("error inserting user: %s", err)
	}

	return nil
}
