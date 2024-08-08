package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"log"
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
			sendAlert(w, "Invalid quantity")
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
			sendAlert(w, "Failed to add book")
			return
		}

		http.Redirect(w, r, "/admin/books/list", http.StatusFound)
	}
}

func BooksPage(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("templates/base.tmpl", "templates/admin/listBooks.tmpl"))

	books, err := models.FetchBooks()
	if err != nil {
		sendAlert(w, err.Error())
		return
	}

	data := struct {
		Books []types.Books
	}{
		Books: books,
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		sendAlert(w, err.Error())
	}
}

func AcceptRequestsPage(w http.ResponseWriter, r *http.Request) {
	requests, err := models.GetPendingRequests()
	if err != nil {
		sendAlert(w, "Error fetching pending requests")
		return
	}

	err = views.AcceptRequests(w, requests)
	if err != nil {
		sendAlert(w, "Error rendering template")
		return
	}
}

func AcceptRequestHandler(w http.ResponseWriter, r *http.Request) {
	requestID := r.FormValue("requestID")

	db, err := models.Connection()
	if err != nil {
		sendAlert(w, "Error connecting to the database")
		return
	}
	defer db.Close()

	// Update request status to 'Accepted' and set AcceptDate
	_, err = db.Exec("UPDATE BookRequests SET Status='Accepted', AcceptDate = NOW() WHERE RequestID=?", requestID)
	if err != nil {
		sendAlert(w, "Error updating request status")
		return
	}

	// Get BookID from the request
	var bookID int
	err = db.QueryRow("SELECT BookID FROM BookRequests WHERE RequestID = ?", requestID).Scan(&bookID)
	if err != nil {
		sendAlert(w, "Error retrieving BookID")
		return
	}

	// Update Books quantity (decrease by 1)
	_, err = db.Exec("UPDATE Books SET Quantity = Quantity - 1 WHERE BookID=?", bookID)
	if err != nil {
		sendAlert(w, "Error updating book quantity")
		return
	}

	http.Redirect(w, r, "/admin/requests", http.StatusSeeOther)
}

func DenyRequestHandler(w http.ResponseWriter, r *http.Request) {
	requestID := r.FormValue("requestID")

	db, err := models.Connection()
	if err != nil {
		sendAlert(w, "Error connecting to the database")
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE BookRequests SET Status='Denied' WHERE RequestID=?", requestID)
	if err != nil {
		sendAlert(w, "Error updating request status")
		return
	}

	http.Redirect(w, r, "/admin/requests", http.StatusSeeOther)
}

func UpdateBooksPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		books, err := models.FetchBooks()
		if err != nil {
			sendAlert(w, "Failed to fetch books")
			return
		}

		err = views.UpdateBooks(w, books) 
		if err != nil {
			sendAlert(w, "Failed to render update page")
			return
		}
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		log.Printf("Form Values: %v", r.PostForm)

		bookIDs := r.Form["bookIDs[]"]

		quantityMap := make(map[string]string)

		for _, bookID := range bookIDs {
			quantityKey := fmt.Sprintf("quantities[%s]", bookID)

			quantity := r.FormValue(quantityKey)

			if quantity != "" {
				quantityMap[bookID] = quantity
			} else {
				sendAlert(w, fmt.Sprintf("Missing quantity for book ID %s", bookID))
				return
			}
		}

		if len(quantityMap) != len(bookIDs) {
			sendAlert(w, "Mismatch between book IDs and quantities")
			return
		}

		log.Printf("Final quantityMap: %v", quantityMap)

		err := models.UpdateBooksQuantity(bookIDs, quantityMap)
		if err != nil {
			sendAlert(w, fmt.Sprintf("Failed to update books: %s", err))
			return
		}

		http.Redirect(w, r, "/admin/books/list", http.StatusSeeOther)
	}
}
