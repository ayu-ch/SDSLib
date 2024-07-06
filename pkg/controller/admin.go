package controller

import (
	"log"
	"net/http"

	"github.com/ayu-ch/SDSLib/pkg/models"
	"github.com/ayu-ch/SDSLib/pkg/views"
)

func AdminPortal(w http.ResponseWriter, r *http.Request) {
	tmpl := views.Portal()
	err := tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ManageBooksPage(w http.ResponseWriter, r *http.Request) {
	tmpl := views.ManageBooks()
	err := tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ChangePrivilegesPage(w http.ResponseWriter, r *http.Request) {
	users, err := models.FetchUsers()
	if err != nil {
		log.Printf("Error fetching users with pending requests: %s", err)
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}

	// Render the template with users data
	err = views.ChangePrivileges(w, users)
	if err != nil {
		log.Printf("Error rendering change privileges page: %s", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}

func AcceptUser(w http.ResponseWriter, r *http.Request) {
	UserID := r.FormValue("UserID")

	db, err := models.Connection()
	if err != nil {
		log.Printf("Error connecting to the database: %s", err)
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE User SET AdminRequest='Accepted' WHERE UserID=?", UserID)
	if err != nil {
		log.Printf("Error updating user admin request: %s", err)
		http.Error(w, "Error updating user admin request", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func DenyUser(w http.ResponseWriter, r *http.Request) {
	UserID := r.FormValue("UserID")

	db, err := models.Connection()
	if err != nil {
		log.Printf("Error connecting to the database: %s", err)
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE User SET AdminRequest='Denied' WHERE UserID=?", UserID)
	if err != nil {
		log.Printf("Error updating user admin request: %s", err)
		http.Error(w, "Error updating user admin request", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}
