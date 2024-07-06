package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ayu-ch/SDSLib/pkg/models"
	"github.com/ayu-ch/SDSLib/pkg/views"
	"github.com/dgrijalva/jwt-go"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl := views.Home(nil)
	err := tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Requests(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Cookie not found", http.StatusUnauthorized)
		return
	}
	tokenString := cookie.Value

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		userID := claims.UserID

		requests, err := models.GetRequestsByUserID(userID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching requests: %s", err), http.StatusInternalServerError)
			return
		}

		if len(requests) > 0 {
			var bookIDs []int
			for _, request := range requests {
				bookIDs = append(bookIDs, request.BookID)
			}

			books, err := models.GetBooksByIDsWithStatus(userID, bookIDs)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error fetching books with status: %s", err), http.StatusInternalServerError)
				return
			}

			tmpl := views.RequestsPage(books)
			err = tmpl.ExecuteTemplate(w, "base", books)
			if err != nil {
				http.Error(w, fmt.Sprintf("error executing template: %s", err), http.StatusInternalServerError)
				return
			}
		} else {
			fmt.Fprintf(w, "You do not have any requests yet")
		}
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
}

func AcceptedBooks(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Cookie not found", http.StatusUnauthorized)
		return
	}
	tokenString := cookie.Value

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		userID := claims.UserID

		requests, err := models.GetAcceptedBooksByID(userID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching requests: %s", err), http.StatusInternalServerError)
			return
		}

		if len(requests) > 0 {
			var bookIDs []int
			for _, request := range requests {
				bookIDs = append(bookIDs, request.BookID)
			}

			books, err := models.GetBooksByIDsWithStatus(userID, bookIDs)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error fetching books with status: %s", err), http.StatusInternalServerError)
				return
			}

			tmpl := views.BooksPage(books)
			err = tmpl.ExecuteTemplate(w, "base", books)
			if err != nil {
				http.Error(w, fmt.Sprintf("error executing template: %s", err), http.StatusInternalServerError)
				return
			}
		} else {
			fmt.Fprintf(w, "You do not have any requests yet")
		}
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
}

func RequestAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	db, err := models.Connection()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error connecting to the database: %s", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Cookie not found", http.StatusUnauthorized)
		return
	}
	tokenString := cookie.Value

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		userID := claims.UserID

		err = db.QueryRow("SELECT UserID FROM User WHERE Username = ?", claims.Username).Scan(&userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("UPDATE User SET AdminRequest = 'Pending' WHERE UserID = ?", userID)
		if err != nil {
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}
	}
}

func AddBookRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		books, err := models.FetchBooks()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching books: %s", err), http.StatusInternalServerError)
			return
		}

		views.RequestBooks(w, books)

	case http.MethodPost:
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Cookie not found", http.StatusUnauthorized)
			return
		}
		tokenString := cookie.Value

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			userID := claims.UserID

			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Error parsing form", http.StatusBadRequest)
				return
			}

			var bookIDs []int
			for _, v := range r.Form["selectedBooks"] {
				bookID, err := strconv.Atoi(v)
				if err != nil {
					http.Error(w, "Invalid book ID", http.StatusBadRequest)
					return
				}
				bookIDs = append(bookIDs, bookID)
			}

			err = models.AddBookRequests(userID, bookIDs)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error adding book requests: %s", err), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/home/requests", http.StatusSeeOther)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func ReturnBooksPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Cookie not found", http.StatusUnauthorized)
			return
		}
		tokenString := cookie.Value

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			userID := claims.UserID
			requests, err := models.GetReturnableBooks(userID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error fetching returnable books: %s", err), http.StatusInternalServerError)
				return
			}
			views.RenderReturnBooks(w, requests)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}

	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		var selectedRequests []int
		for _, v := range r.Form["selectedRequests"] {
			requestID, err := strconv.Atoi(v)
			if err != nil {
				http.Error(w, "Invalid request ID", http.StatusBadRequest)
				return
			}
			selectedRequests = append(selectedRequests, requestID)
		}

		err = models.ReturnBooks(selectedRequests)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error returning books: %s", err), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/home/requests", http.StatusSeeOther)

	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func BorrowHistoryHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Cookie not found", http.StatusUnauthorized)
		return
	}
	tokenString := cookie.Value

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		userID := claims.UserID
		history, err := models.GetBorrowingHistory(userID)
		if err != nil {
			http.Error(w, "Error fetching borrowing history", http.StatusInternalServerError)
			return
		}

		views.RenderBorrowHistory(w, history)
	}
}
