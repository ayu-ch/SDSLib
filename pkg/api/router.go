package api

import (
	"fmt"
	"net/http"

	"github.com/ayu-ch/SDSLib/pkg/controller"
	"github.com/ayu-ch/SDSLib/pkg/middleware"
)

func Start() {

	http.Handle("/", middleware.Auths(http.HandlerFunc(controller.SignupPage)))
	http.Handle("/login", middleware.Auths(http.HandlerFunc(controller.LoginPage)))
	http.Handle("/admin/login", middleware.Auths(http.HandlerFunc(controller.AdminLoginPage)))
	http.HandleFunc("/logout", controller.Logout)

	http.Handle("/home", middleware.IsLoggedIn(http.HandlerFunc(controller.HomePage)))
	http.Handle("/request", middleware.IsLoggedIn(http.HandlerFunc(controller.AddBookRequest)))
	http.Handle("/requests", middleware.IsLoggedIn(http.HandlerFunc(controller.Requests)))
	http.Handle("/requestAdmin", middleware.IsLoggedIn(http.HandlerFunc(controller.RequestAdmin)))
	http.Handle("/books", middleware.IsLoggedIn(http.HandlerFunc(controller.AcceptedBooks)))
	http.Handle("/return", middleware.IsLoggedIn(http.HandlerFunc(controller.ReturnBooksPage)))
	http.Handle("/borrowHistory", middleware.IsLoggedIn(http.HandlerFunc(controller.BorrowHistoryHandler)))

	http.Handle("/admin", middleware.IsAdmin(http.HandlerFunc(controller.AdminPortal)))

	http.Handle("/admin/users", middleware.IsAdmin(http.HandlerFunc(controller.ChangePrivilegesPage)))
	http.Handle("/admin/users/accept", middleware.IsAdmin(http.HandlerFunc(controller.AcceptUser)))
	http.Handle("/admin/users/deny", middleware.IsAdmin(http.HandlerFunc(controller.DenyUser)))

	http.Handle("/admin/requests", middleware.IsAdmin(http.HandlerFunc(controller.AcceptRequestsPage)))
	http.Handle("/admin/requests/accept", middleware.IsAdmin(http.HandlerFunc(controller.AcceptRequestHandler)))
	http.Handle("/admin/requests/deny", middleware.IsAdmin(http.HandlerFunc(controller.DenyRequestHandler)))

	http.Handle("/admin/books", middleware.IsAdmin(http.HandlerFunc(controller.ManageBooksPage)))
	http.Handle("/admin/books/list", middleware.IsAdmin(http.HandlerFunc(controller.BooksPage)))
	http.Handle("/admin/books/add", middleware.IsAdmin(http.HandlerFunc(controller.AddBookPage)))
	http.Handle("/admin/books/update", middleware.IsAdmin(http.HandlerFunc(controller.UpdateBooksPage)))

	fmt.Println("Starting server on :8000")
	http.ListenAndServe(":8000", nil)
}
