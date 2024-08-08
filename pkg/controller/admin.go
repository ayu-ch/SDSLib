package controller

import (
	"log"
	"net/http"
    // "strconv"
	"github.com/ayu-ch/SDSLib/pkg/models"
	"github.com/ayu-ch/SDSLib/pkg/views"
)



func AdminPortal(w http.ResponseWriter, r *http.Request) {
	tmpl := views.Portal()
	err := tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		sendAlert(w, err.Error())
	}
}

func ManageBooksPage(w http.ResponseWriter, r *http.Request) {
	tmpl := views.ManageBooks()
	err := tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		sendAlert(w, err.Error())
	}
}

func ChangePrivilegesPage(w http.ResponseWriter, r *http.Request) {
	users, err := models.FetchUsers()
	if err != nil {
		log.Printf("Error fetching users with pending requests: %s", err)
		sendAlert(w, "Error fetching users")
		return
	}

	err = views.ChangePrivileges(w, users)
	if err != nil {
		log.Printf("Error rendering change privileges page: %s", err)
		sendAlert(w, "Error rendering page")
	}
}

func AcceptUser(w http.ResponseWriter, r *http.Request) {
	UserID := r.FormValue("UserID")
    log.Printf("UserID:%s",r.
	db, err := models.Connection()
	if err != nil {
		log.Printf("Error connecting to the database: %s", err)
		sendAlert(w, "Error connecting to the database")
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE User SET AdminRequest='Accepted' AND Role='Admin' WHERE UserID=?", UserID)
	if err != nil {
		log.Printf("Error updating user admin request: %s", err)
		sendAlert(w, "Error updating user admin request")
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func DenyUser(w http.ResponseWriter, r *http.Request) {
	UserID := r.FormValue("UserID")

	db, err := models.Connection()
	if err != nil {
		log.Printf("Error connecting to the database: %s", err)
		sendAlert(w, "Error connecting to the database")
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE User SET AdminRequest='Denied' WHERE UserID=?", UserID)
	if err != nil {
		log.Printf("Error updating user admin request: %s", err)
		sendAlert(w, "Error updating user admin request")
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}
