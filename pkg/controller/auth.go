package controller

import (
	"log"
	"net/http"
	"os"

	"github.com/ayu-ch/SDSLib/pkg/models"
	"github.com/ayu-ch/SDSLib/pkg/views"
	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID   int    `json:"userID"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := views.Login()
		err := tmpl.ExecuteTemplate(w, "base", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error executing login template: %v", err)
		}
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("name")
		password := r.FormValue("password")

		log.Printf("Attempting login for user: %s", username)

		user, err := models.Login(username, password)
		if err != nil {
			log.Printf("Login failed for user %s: %v", username, err)
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		claims := &Claims{
			UserID:   user.UserID,
			Username: user.Username,
			Role:     user.Role,
		}

		if user.Role != "Client" {
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
			return
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			log.Printf("Error generating token for user %s: %v", username, err)
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		log.Printf("Login successful for user: %s", username)

		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: tokenString,
			Path:  "/",
		})

		http.Redirect(w, r, "/home", http.StatusFound)
	}
}

func SignupPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := views.Signup()
		err := tmpl.ExecuteTemplate(w, "base", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("name")
		password := r.FormValue("password")

		err := models.Signup(username, password)
		if err != nil {
			http.Error(w, "Error signing up user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func AdminLoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := views.AdminLogin()
		err := tmpl.ExecuteTemplate(w, "base", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error executing login template: %v", err)
		}
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("name")
		password := r.FormValue("password")

		log.Printf("Attempting login for user: %s", username)

		user, err := models.Login(username, password)
		if err != nil {
			log.Printf("Login failed for user %s: %v", username, err)
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		claims := &Claims{
			UserID:   user.UserID,
			Username: user.Username,
			Role:     user.Role,
		}

		if user.Role != "Admin" {
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
			return
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			log.Printf("Error generating token for user %s: %v", username, err)
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		log.Printf("Login successful for user: %s", username)

		// Set token as cookie
		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: tokenString,
			Path:  "/", // Path to set cookie on
		})

		// Redirect user after successful login
		http.Redirect(w, r, "/admin", http.StatusFound)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		http.SetCookie(w, &http.Cookie{
			Name:  cookie.Name,
			Value: "",
			Path:  "/",
		})
	}

	http.Redirect(w, r, "/admin/login", http.StatusFound)
}
