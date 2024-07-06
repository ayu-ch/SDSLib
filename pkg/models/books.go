package models

import (
	"fmt"

	"github.com/ayu-ch/SDSLib/pkg/types"
)

func AddBook(book types.Books) error {
	db, err := Connection()
	if err != nil {
		return fmt.Errorf("error connecting to the database: %s", err)
	}
	defer db.Close()

	query := "INSERT INTO Books (Title, Author, Genre, Quantity) VALUES (?, ?, ?, ?)"
	_, err = db.Exec(query, book.Title, book.Author, book.Genre, book.Quantity)
	if err != nil {
		return fmt.Errorf("error inserting book: %s", err)
	}

	return nil
}

func FetchBooks() ([]types.Books, error) {
	db, err := Connection()
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %s", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT BookID, Title, Author, Genre, Quantity FROM Books")
	if err != nil {
		return nil, fmt.Errorf("error querying the database: %s", err)
	}
	defer rows.Close()

	var books []types.Books
	for rows.Next() {
		var book types.Books
		err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Genre, &book.Quantity)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %s", err)
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %s", err)
	}

	return books, nil
}

func GetPendingRequests() ([]types.BooksWithUser, error) {
	query := `
		SELECT 
			BookRequests.BookID as RequestedBookID, 
			BookRequests.UserID,
			BookRequests.RequestID,
			BookRequests.RequestDate,
			User.Username,
			Books.BookID as AvailableBookID,
			Books.Title,
			Books.Author,
			Books.Genre,
			Books.Quantity
		FROM BookRequests 
		JOIN User ON BookRequests.UserID = User.UserID 
		JOIN Books ON BookRequests.BookID = Books.BookID 
		WHERE BookRequests.Status = 'Pending'`

	db, err := Connection()
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %s", err)
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []types.BooksWithUser
	for rows.Next() {
		var req types.BooksWithUser
		if err := rows.Scan(
			&req.RequestedBookID,
			&req.UserID,
			&req.RequestID,
			&req.RequestDate,
			&req.Username,
			&req.AvailableBookID,
			&req.Title,
			&req.Author,
			&req.Genre,
			&req.Quantity,
		); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

func FetchUsers() ([]types.User, error) {
	query := `
		SELECT UserID, Username, Pass, Role, AdminRequest
		FROM User
		WHERE AdminRequest = 'Pending'`

	db, err := Connection()
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %s", err)
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []types.User
	for rows.Next() {
		var user types.User
		if err := rows.Scan(
			&user.UserID,
			&user.Username,
			&user.Pass,
			&user.Role,
			&user.AdminRequest,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
