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
		sendAlert(w, fmt.Sprintf("Error executing template: %s", err))
	}
}

func sendAlert(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<script type="text/javascript">
				alert("%s");
				window.history.back();
			</script>
		</head>
		<body></body>
		</html>
	`, message)
}

func Requests(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		sendAlert(w, "Cookie not found")
		return
	}
	tokenString := cookie.Value

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		sendAlert(w, "Unauthorized")
		return
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		userID := claims.UserID

		requests, err := models.GetRequestsByUserID(userID)
		if err != nil {
			sendAlert(w, fmt.Sprintf("Error fetching requests: %s", err))
			return
		}

		if len(requests) > 0 {
			var bookIDs []int
			for _, request := range requests {
				bookIDs = append(bookIDs, request.BookID)
			}

			books, err := models.GetBooksByIDsWithStatus(userID, bookIDs)
			if err != nil {
				sendAlert(w, fmt.Sprintf("Error fetching books with status: %s", err))
				return
			}

			tmpl := views.RequestsPage(books)
			err = tmpl.ExecuteTemplate(w, "base", books)
			if err != nil {
				sendAlert(w, fmt.Sprintf("Error executing template: %s", err))
				return
			}
		} else {
			fmt.Fprintf(w, "You do not have any requests yet")
		}
	} else {
		sendAlert(w, "Unauthorized")
		return
	}
}

func AcceptedBooks(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		sendAlert(w, "Cookie not found")
		return
	}
	tokenString := cookie.Value

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		sendAlert(w, "Unauthorized")
		return
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		userID := claims.UserID

		requests, err := models.GetAcceptedBooksByID(userID)
		if err != nil {
			sendAlert(w, fmt.Sprintf("Error fetching requests: %s", err))
			return
		}

		if len(requests) > 0 {
			var bookIDs []int
			for _, request := range requests {
				bookIDs = append(bookIDs, request.BookID)
			}

			books, err := models.GetBooksByIDsWithStatus(userID, bookIDs)
			if err != nil {
				sendAlert(w, fmt.Sprintf("Error fetching books with status: %s", err))
				return
			}

			tmpl := views.BooksPage(books)
			err = tmpl.ExecuteTemplate(w, "base", books)
			if err != nil {
				sendAlert(w, fmt.Sprintf("Error executing template: %s", err))
				return
			}
		} else {
			fmt.Fprintf(w, "You do not have any requests yet")
		}
	} else {
		sendAlert(w, "Unauthorized")
		return
	}
}

func RequestAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendAlert(w, "Invalid request method")
		return
	}

	db, err := models.Connection()
	if err != nil {
		sendAlert(w, fmt.Sprintf("Error connecting to the database: %s", err))
		return
	}
	defer db.Close()

	cookie, err := r.Cookie("token")
	if err != nil {
		sendAlert(w, "Cookie not found")
		return
	}
	tokenString := cookie.Value

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		sendAlert(w, "Unauthorized")
		return
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		userID := claims.UserID

		err = db.QueryRow("SELECT UserID FROM User WHERE Username = ?", claims.Username).Scan(&userID)
		if err != nil {
			sendAlert(w, "User not found")
			return
		}

		_, err = db.Exec("UPDATE User SET AdminRequest = 'Pending' WHERE UserID = ?", userID)
		if err != nil {
			sendAlert(w, "Failed to update user")
			return
		}else{
			sendAlert(w,"Request Sent!")
		}
	}
}

func AddBookRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		books, err := models.FetchBooks()
		if err != nil {
			sendAlert(w, fmt.Sprintf("Error fetching books: %s", err))
			return
		}

		views.RequestBooks(w, books)

	case http.MethodPost:
		cookie, err := r.Cookie("token")
		if err != nil {
			sendAlert(w, "Cookie not found")
			return
		}
		tokenString := cookie.Value

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			sendAlert(w, "Unauthorized")
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			userID := claims.UserID

			err := r.ParseForm()
			if err != nil {
				sendAlert(w, "Error parsing form")
				return
			}

			var bookIDs []int
			for _, v := range r.Form["selectedBooks"] {
				bookID, err := strconv.Atoi(v)
				if err != nil {
					sendAlert(w, "Invalid book ID")
					return
				}
				bookIDs = append(bookIDs, bookID)
			}

			err = models.AddBookRequests(userID, bookIDs)
			if err != nil {
				sendAlert(w, fmt.Sprintf("Error adding book requests: %s", err))
				return
			}

			http.Redirect(w, r, "/requests", http.StatusSeeOther)
		} else {
			sendAlert(w, "Unauthorized")
			return
		}

	default:
		sendAlert(w, "Invalid request method")
	}
}

func ReturnBooksPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		cookie, err := r.Cookie("token")
		if err != nil {
			sendAlert(w, "Cookie not found")
			return
		}
		tokenString := cookie.Value

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			sendAlert(w, "Unauthorized")
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			userID := claims.UserID
			requests, err := models.GetReturnableBooks(userID)
			if err != nil {
				sendAlert(w, fmt.Sprintf("Error fetching returnable books: %s", err))
				return
			}
			views.RenderReturnBooks(w, requests)
		} else {
			sendAlert(w, "Unauthorized")
		}

	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			sendAlert(w, "Error parsing form")
			return
		}

		var selectedRequests []int
		for _, v := range r.Form["selectedRequests"] {
			requestID, err := strconv.Atoi(v)
			if err != nil {
				sendAlert(w, "Invalid request ID")
				return
			}
			selectedRequests = append(selectedRequests, requestID)
		}

		err = models.ReturnBooks(selectedRequests)
		if err != nil {
			sendAlert(w, fmt.Sprintf("Error returning books: %s", err))
			return
		}

		http.Redirect(w, r, "/requests", http.StatusSeeOther)

	default:
		sendAlert(w, "Invalid request method")
	}
}

func BorrowHistoryHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		sendAlert(w, "Cookie not found")
		return
	}
	tokenString := cookie.Value

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		sendAlert(w, "Unauthorized")
		return
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		userID := claims.UserID
		history, err := models.GetBorrowingHistory(userID)
		if err != nil {
			sendAlert(w, "Error fetching borrowing history")
			return
		}

		views.RenderBorrowHistory(w, history)
	}
}
