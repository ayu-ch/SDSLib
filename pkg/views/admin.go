package views

import (
	"html/template"
	"log"
	"net/http"

	"github.com/ayu-ch/SDSLib/pkg/types"
)

func Portal() *template.Template {
	tmpl, err := template.ParseFiles(
		"templates/base.tmpl",
		"templates/admin/adminPortal.tmpl",
	)
	if err != nil {
		panic(err)
	}
	return tmpl
}

func ManageBooks() *template.Template {
	tmpl, err := template.ParseFiles(
		"templates/base.tmpl",
		"templates/admin/manageBooks.tmpl",
	)
	if err != nil {
		panic(err)
	}
	return tmpl
}

func ChangePrivileges(w http.ResponseWriter, data []types.User) error {
	tmpl, err := template.ParseFiles("templates/base.tmpl", "templates/admin/changePrivileges.tmpl")
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
