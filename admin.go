package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	dbUsername := "root"
	dbPassword := "Physics@1504"
	dbHost := "localhost"
	dbName := "Library"

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUsername, dbPassword, dbHost, dbName))
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	adminUsername := os.Args[1]
	adminPassword := os.Args[2]

	adminUsername = strings.TrimSpace(adminUsername)
	adminPassword = strings.TrimSpace(adminPassword)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
	}

	insertQuery := `INSERT INTO User (Username, Pass, Role) VALUES (?, ?, ?)`
	_, err = db.Exec(insertQuery, adminUsername, hashedPassword, "Admin")
	if err != nil {
		log.Fatalf("Error inserting admin user: %v", err)
	}

	fmt.Println("Admin user added successfully.")
}