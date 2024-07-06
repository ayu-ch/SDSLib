package views

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/ayu-ch/SDSLib/pkg/types"
)

func Home(data interface{}) *template.Template {
	tmpl := template.New("home")
	tmpl, err := tmpl.ParseFiles(
		"templates/base.tmpl",
		"templates/client/home.tmpl",
	)
	if err != nil {
		panic(err)
	}
	return tmpl
}

func RequestsPage(data []types.BooksWithStatus) *template.Template {
	tmpl, err := template.New("requests").
		ParseFiles(
			filepath.Join("templates", "base.tmpl"),
			filepath.Join("templates", "client/requestsClient.tmpl"),
		)
	if err != nil {
		log.Fatalf("error parsing template files: %v", err)
	}
	return tmpl
}

func BooksPage(data []types.BooksWithStatus) *template.Template {
	tmpl, err := template.New("requests").
		ParseFiles(
			filepath.Join("templates", "base.tmpl"),
			filepath.Join("templates", "client/acceptedBooks.tmpl"),
		)
	if err != nil {
		log.Fatalf("error parsing template files: %v", err)
	}
	return tmpl
}

func RequestBooks(w http.ResponseWriter, books []types.Books) {
	tmpl, err := template.ParseFiles("templates/base.tmpl", "templates/client/requestBooks.tmpl")
	if err != nil {
		log.Printf("Error parsing template: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Books []types.Books
	}{
		Books: books,
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Error executing template: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func RenderReturnBooks(w http.ResponseWriter, data []types.BookRequestsWithTitle) {
	tmpl, err := template.New("requests").
		ParseFiles(
			filepath.Join("templates", "base.tmpl"),
			filepath.Join("templates", "client/returnBooks.tmpl"),
		)
	if err != nil {
		log.Fatalf("error parsing template files: %v", err)
	}

	tmpl.ExecuteTemplate(w, "base", struct {
		Requests []types.BookRequestsWithTitle
	}{Requests: data})
}

func RenderBorrowHistory(w http.ResponseWriter, history []types.BookRequestsWithTitle) {
	tmpl, err := template.New("base").ParseFiles(
		filepath.Join("templates", "base.tmpl"),
		filepath.Join("templates", "client/borrowHistory.tmpl"),
	)
	if err != nil {
		log.Fatalf("error parsing template files: %v", err)
	}

	data := struct {
		History []types.BookRequestsWithTitle
	}{
		History: history,
	}

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
