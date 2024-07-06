package controller

import (
	// "encoding/json"

	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/ayu-ch/SDSLib/pkg/models"
	"github.com/ayu-ch/SDSLib/pkg/types"
	"github.com/ayu-ch/SDSLib/pkg/views"
)

func AddBookPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		views.RenderAddBookPage(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		author := r.FormValue("author")
		genre := r.FormValue("genre")
		quantity := r.FormValue("quantity")

		// Convert quantity to int
		quantityInt, err := strconv.Atoi(quantity)
		if err != nil {
			http.Error(w, "Invalid quantity", http.StatusBadRequest)
			return
		}

		book := types.Books{
			Title:    title,
			Author:   author,
			Genre:    genre,
			Quantity: quantityInt,
		}

		err = models.AddBook(book)
		if err != nil {
			http.Error(w, "Failed to add book", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin/books/list", http.StatusFound)
	}
}

func BooksPage(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("templates/base.tmpl", "templates/admin/listBooks.tmpl"))

	books, err := models.FetchBooks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Books []types.Books
	}{
		Books: books,
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AcceptRequestsPage(w http.ResponseWriter, r *http.Request) {
	requests, err := models.GetPendingRequests()
	if err != nil {
		log.Printf("Error fetching pending requests: %s", err)
		http.Error(w, "Error fetching pending requests", http.StatusInternalServerError)
		return
	}

	err = views.AcceptRequests(w, requests)
	if err != nil {
		log.Printf("Error rendering template: %s", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func AcceptRequestHandler(w http.ResponseWriter, r *http.Request) {
	requestID := r.FormValue("requestID")

	db, err := models.Connection()
	if err != nil {
		log.Printf("Error connecting to the database: %s", err)
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Update request status to 'Accepted' and set AcceptDate
	_, err = db.Exec("UPDATE BookRequests SET Status='Accepted', AcceptDate = NOW() WHERE RequestID=?", requestID)
	if err != nil {
		log.Printf("Error updating request status: %s", err)
		http.Error(w, "Error updating request status", http.StatusInternalServerError)
		return
	}

	// Get BookID from the request
	var bookID int
	err = db.QueryRow("SELECT BookID FROM BookRequests WHERE RequestID = ?", requestID).Scan(&bookID)
	if err != nil {
		log.Printf("Error retrieving BookID: %s", err)
		http.Error(w, "Error retrieving BookID", http.StatusInternalServerError)
		return
	}

	// Update Books quantity (decrease by 1)
	_, err = db.Exec("UPDATE Books SET Quantity = Quantity - 1 WHERE BookID=?", bookID)
	if err != nil {
		log.Printf("Error updating book quantity: %s", err)
		http.Error(w, "Error updating book quantity", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/requests", http.StatusSeeOther)
}

func DenyRequestHandler(w http.ResponseWriter, r *http.Request) {
	requestID := r.FormValue("requestID")

	db, err := models.Connection()
	if err != nil {
		log.Printf("Error connecting to the database: %s", err)
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE BookRequests SET Status='Denied' WHERE RequestID=?", requestID)
	if err != nil {
		log.Printf("Error updating request status: %s", err)
		http.Error(w, "Error updating request status", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/requests", http.StatusSeeOther)
}
