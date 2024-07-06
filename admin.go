package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUsername, dbPassword, dbHost, dbName))
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter admin username: ")
	adminUsername, _ := reader.ReadString('\n')
	adminUsername = strings.TrimSpace(adminUsername)

	fmt.Print("Enter admin password: ")
	adminPassword, _ := reader.ReadString('\n')
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
