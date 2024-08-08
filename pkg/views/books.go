package views

import (
	"html/template"
	"log"
	"net/http"

	"github.com/ayu-ch/SDSLib/pkg/types"
)
type BooksData struct {
    Books []types.Books
}
func ListPage() *template.Template {
	temp := template.Must(template.ParseFiles("templates/admin/list.html"))
	return temp
}

func RenderAddBookPage(w http.ResponseWriter, data interface{}) {
	var (
		tmpl = template.Must(template.ParseFiles("templates/base.tmpl", "templates/admin/addBook.tmpl"))
	)
	err := tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Error rendering template: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func AcceptRequests(w http.ResponseWriter, data []types.BooksWithUser) error {
	tmpl, err := template.ParseFiles("templates/base.tmpl", "templates/admin/acceptRequests.tmpl")
	if err != nil {
		log.Printf("Error parsing template: %s", err)
		return err
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Error executing template: %s", err)
		return err
	}

	return nil
}

func UpdateBooks(w http.ResponseWriter, books []types.Books) error {
    tmpl, err := template.ParseFiles("templates/base.tmpl", "templates/admin/updateBooks.tmpl")
    if err != nil {
        log.Printf("Error parsing template: %s", err)
        return err
    }

    // Wrap the books in a struct
    data := struct {
        Books []types.Books
    }{
        Books: books,
    }

    err = tmpl.ExecuteTemplate(w, "base", data)
    if err != nil {
        log.Printf("Error executing template: %s", err)
        return err
    }

    return nil
}

