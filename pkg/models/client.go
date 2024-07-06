package models

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ayu-ch/SDSLib/pkg/types"
)

func GetRequestsByUserID(userID int) ([]types.BookRequests, error) {
	db, err := Connection()
	if err != nil {
		log.Printf("Error connecting to the database: %s", err)
		return nil, fmt.Errorf("error connecting to the database: %s", err)
	}
	defer db.Close()

	query := "SELECT RequestID, UserID, BookID, Status FROM BookRequests WHERE UserID = ? AND Status IN ('Pending', 'Denied')"
	rows, err := db.Query(query, userID)
	if err != nil {
		log.Printf("Error executing query: %s", err)
		return nil, fmt.Errorf("error executing query: %s", err)
	}
	defer rows.Close()

	var requests []types.BookRequests
	for rows.Next() {
		var request types.BookRequests
		if err := rows.Scan(&request.RequestID, &request.UserID, &request.BookID, &request.Status); err != nil {
			log.Printf("Error scanning row: %s", err)
			return nil, fmt.Errorf("error scanning row: %s", err)
		}
		requests = append(requests, request)

		log.Printf("Fetched request: %+v", request)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %s", err)
		return nil, fmt.Errorf("error iterating rows: %s", err)
	}

	log.Printf("Fetched %d requests successfully for UserID %d", len(requests), userID)

	return requests, nil
}

func GetBooksByIDsWithStatus(userID int, bookIDs []int) ([]types.BooksWithStatus, error) {
	db, err := Connection()
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %s", err)
	}
	defer db.Close()

	placeholders := make([]string, len(bookIDs))
	args := make([]interface{}, len(bookIDs)+1)
	for i, id := range bookIDs {
		placeholders[i] = "?"
		args[i] = id
	}
	args[len(bookIDs)] = userID

	query := fmt.Sprintf(`
			SELECT b.BookID, b.Title, b.Author, b.Genre, b.Quantity, br.Status
			FROM Books b
			JOIN BookRequests br ON b.BookID = br.BookID
			WHERE b.BookID IN (%s) AND br.UserID = ?
	`, strings.Join(placeholders, ","))

	log.Printf("Executing query: %s with args: %v", query, args)

	var books []types.BooksWithStatus
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var book types.BooksWithStatus
		if err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Genre, &book.Quantity, &book.Status); err != nil {
			return nil, fmt.Errorf("error scanning row: %s", err)
		}
		books = append(books, book)

		log.Printf("Fetched book: %+v", book)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %s", err)
	}

	log.Printf("Fetched %d books successfully for UserID %d", len(books), userID)

	return books, nil
}

func GetAcceptedBooksByID(userID int) ([]types.BookRequests, error) {
	db, err := Connection()
	if err != nil {
		log.Printf("Error connecting to the database: %s", err)
		return nil, fmt.Errorf("error connecting to the database: %s", err)
	}
	defer db.Close()

	query := "SELECT RequestID, UserID, BookID, Status FROM BookRequests WHERE UserID = ? AND Status = 'Accepted'"
	rows, err := db.Query(query, userID)
	if err != nil {
		log.Printf("Error executing query: %s", err)
		return nil, fmt.Errorf("error executing query: %s", err)
	}
	defer rows.Close()

	var requests []types.BookRequests
	for rows.Next() {
		var request types.BookRequests
		if err := rows.Scan(&request.RequestID, &request.UserID, &request.BookID, &request.Status); err != nil {
			log.Printf("Error scanning row: %s", err)
			return nil, fmt.Errorf("error scanning row: %s", err)
		}
		requests = append(requests, request)

		log.Printf("Fetched request: %+v", request)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %s", err)
		return nil, fmt.Errorf("error iterating rows: %s", err)
	}

	log.Printf("Fetched %d requests successfully for UserID %d", len(requests), userID)

	return requests, nil
}

func AddBookRequests(userID int, bookIDs []int) error {
	db, err := Connection()
	if err != nil {
		return err
	}
	defer db.Close()

	checkStmt, err := db.Prepare("SELECT COUNT(*) FROM BookRequests WHERE UserID = ? AND BookID = ? AND Status IN ('Accepted', 'Pending')")
	if err != nil {
		log.Printf("Error preparing check statement: %s", err)
		return err
	}
	defer checkStmt.Close()

	insertStmt, err := db.Prepare("INSERT INTO BookRequests (UserID, BookID, RequestDate, Status) VALUES (?, ?, ?, 'Pending')")
	if err != nil {
		log.Printf("Error preparing insert statement: %s", err)
		return err
	}
	defer insertStmt.Close()

	for _, bookID := range bookIDs {
		var count int
		err := checkStmt.QueryRow(userID, bookID).Scan(&count)
		if err != nil {
			log.Printf("Error checking existing book request: %s", err)
			return err
		}

		if count > 0 {
			return fmt.Errorf("you already have or requested this book: %d", bookID)
		}

		_, err = insertStmt.Exec(userID, bookID, time.Now())
		if err != nil {
			log.Printf("Error inserting book request: %s", err)
			return err
		}
	}

	return nil
}

func GetReturnableBooks(userID int) ([]types.BookRequestsWithTitle, error) {
	db, err := Connection()
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %s", err)
	}
	defer db.Close()

	query := `
		SELECT br.RequestID, br.UserID, br.BookID, br.RequestDate, br.Status, b.Title
		FROM BookRequests br
		JOIN Books b ON br.BookID = b.BookID
		WHERE br.Status = 'Accepted' AND br.UserID = ?
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying the database: %s", err)
	}
	defer rows.Close()

	var requests []types.BookRequestsWithTitle
	for rows.Next() {
		var req types.BookRequestsWithTitle
		var requestDate string
		if err := rows.Scan(&req.RequestID, &req.UserID, &req.BookID, &requestDate, &req.Status, &req.Title); err != nil {
			return nil, fmt.Errorf("error scanning row: %s", err)
		}

		req.RequestDate, err = time.Parse("2006-01-02 15:04:05", requestDate)
		if err != nil {
			return nil, fmt.Errorf("error parsing date: %s", err)
		}

		requests = append(requests, req)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %s", err)
	}

	return requests, nil
}

func ReturnBooks(selectedRequests []int) error {
	db, err := Connection()
	if err != nil {
		return fmt.Errorf("error connecting to the database: %s", err)
	}
	defer db.Close()

	for _, requestID := range selectedRequests {
		_, err := db.Exec(`UPDATE BookRequests SET Status = 'Returned' WHERE RequestID = ?`, requestID)
		if err != nil {
			return fmt.Errorf("error updating BookRequests table for RequestID %d: %s", requestID, err)
		}

		var bookID int
		err = db.QueryRow(`SELECT BookID FROM BookRequests WHERE RequestID = ?`, requestID).Scan(&bookID)
		if err != nil {
			return fmt.Errorf("error retrieving BookID for RequestID %d: %s", requestID, err)
		}

		_, err = db.Exec(`UPDATE Books SET Quantity = Quantity + 1 WHERE BookID = ?`, bookID)
		if err != nil {
			return fmt.Errorf("error updating Books table for BookID %d: %s", bookID, err)
		}
	}

	return nil
}

func GetBorrowingHistory(userID int) ([]types.BookRequestsWithTitle, error) {
	db, err := Connection()
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}
	defer db.Close()

	query := `
		SELECT br.AcceptDate, br.RequestID, br.BookID, b.Title 
		FROM BookRequests br 
		INNER JOIN Books b ON br.BookID = b.BookID 
		WHERE br.UserID = ? AND (br.Status = 'Accepted' OR br.Status = 'Returned')
	`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var history []types.BookRequestsWithTitle
	for rows.Next() {
		var req types.BookRequestsWithTitle
		var acceptDate string
		if err := rows.Scan(&acceptDate, &req.RequestID, &req.BookID, &req.Title); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		req.RequestDate, err = time.Parse("2006-01-02 15:04:05", acceptDate)
		if err != nil {
			return nil, fmt.Errorf("error parsing AcceptDate: %v", err)
		}
		history = append(history, req)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error with rows: %v", err)
	}

	return history, nil
}
