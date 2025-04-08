package auth

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/pocketbase/pocketbase/core"
)

// Templates for auth pages
var templates = template.Must(template.ParseFiles(
	filepath.Join("internal", "templates", "login.html"),
	filepath.Join("internal", "templates", "register.html"),
))

// LoginForm represents the login form data
type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Error    string `json:"error,omitempty"`
}

// RegisterForm represents the registration form data
type RegisterForm struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Error           string `json:"error,omitempty"`
}

// LoginHandler shows the login form
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		templates.ExecuteTemplate(w, "login.html", nil)
		return
	}

	// Process form submission
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	// Validate input
	if email == "" || password == "" {
		templates.ExecuteTemplate(w, "login.html", LoginForm{
			Error: "Email and password are required",
		})
		return
	}

	// Make sure PocketBase client is initialized
	if PbClient == nil {
		templates.ExecuteTemplate(w, "login.html", LoginForm{
			Email: email,
			Error: "Authentication system not available",
		})
		return
	}

	// Find user by email
	authRecord, err := PbClient.FindAuthRecordByEmail("users", email)
	if err != nil {
		templates.ExecuteTemplate(w, "login.html", LoginForm{
			Email: email,
			Error: "Invalid email or password",
		})
		return
	}

	// Validate password
	if !authRecord.ValidatePassword(password) {
		templates.ExecuteTemplate(w, "login.html", LoginForm{
			Email: email,
			Error: "Invalid email or password",
		})
		return
	}

	// Generate auth token
	token, err := authRecord.NewAuthToken()
	if err != nil {
		templates.ExecuteTemplate(w, "login.html", LoginForm{
			Email: email,
			Error: "Failed to create authentication token",
		})
		return
	}

	// Set cookie with the auth token
	http.SetCookie(w, &http.Cookie{
		Name:     "pb_auth",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// RegisterHandler shows the registration form
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		templates.ExecuteTemplate(w, "register.html", nil)
		return
	}

	// Process form submission
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	// Validate inputs
	if email == "" || password == "" {
		templates.ExecuteTemplate(w, "register.html", RegisterForm{
			Error: "Email and password are required",
		})
		return
	}

	// Validate passwords match
	if password != confirmPassword {
		templates.ExecuteTemplate(w, "register.html", RegisterForm{
			Email: email,
			Error: "Passwords do not match",
		})
		return
	}

	// Make sure PocketBase client is initialized
	if PbClient == nil {
		templates.ExecuteTemplate(w, "register.html", RegisterForm{
			Email: email,
			Error: "Registration system not available",
		})
		return
	}

	// Find the users collection
	collection, err := PbClient.FindCollectionByNameOrId("users")
	if err != nil {
		templates.ExecuteTemplate(w, "register.html", RegisterForm{
			Email: email,
			Error: "User system not configured correctly",
		})
		return
	}

	// Create a new user record
	record := core.NewRecord(collection)
	record.SetEmail(email)
	record.SetPassword(password)

	// Save the record
	if err := PbClient.Save(record); err != nil {
		templates.ExecuteTemplate(w, "register.html", RegisterForm{
			Email: email,
			Error: "Registration failed: " + err.Error(),
		})
		return
	}

	// Generate auth token for the new user
	token, err := record.NewAuthToken()
	if err != nil {
		// Registration succeeded but token generation failed - redirect to login
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	// Set cookie with the auth token
	http.SetCookie(w, &http.Cookie{
		Name:     "pb_auth",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// LogoutHandler logs the user out
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the auth cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "pb_auth",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	// Redirect to login page
	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}

// PocketBaseAuthHandler forwards authentication requests to PocketBase's API
func PocketBaseAuthHandler(w http.ResponseWriter, r *http.Request) {
	// Make sure PocketBase client is initialized
	if PbClient == nil {
		http.Error(w, "Authentication system not available", http.StatusInternalServerError)
		return
	}

	// Extract the action from the URL
	vars := mux.Vars(r)
	action := vars["action"]

	// Parse form data from the request
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	var result interface{}

	// Convert form data to appropriate format
	formData := make(map[string]any)
	for key, values := range r.Form {
		if len(values) > 0 {
			formData[key] = values[0]
		}
	}

	// Process the auth request based on action
	switch action {
	case "login":
		email, _ := formData["email"].(string)
		password, _ := formData["password"].(string)

		// Find the user by email
		record, err := PbClient.FindAuthRecordByEmail("users", email)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusBadRequest)
			return
		}

		// Validate password
		if !record.ValidatePassword(password) {
			http.Error(w, "Invalid email or password", http.StatusBadRequest)
			return
		}

		// Generate auth token
		token, err := record.NewAuthToken()
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// Set cookie with auth token
		http.SetCookie(w, &http.Cookie{
			Name:     "pb_auth",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(24 * time.Hour),
		})

		result = map[string]any{
			"token": token,
			"user":  record.PublicExport(),
		}

	case "register":
		// Create a new user record
		collection, err := PbClient.FindCollectionByNameOrId("users")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		record := core.NewRecord(collection)

		// Set user fields from form data
		email, _ := formData["email"].(string)
		password, _ := formData["password"].(string)

		record.SetEmail(email)
		record.SetPassword(password)

		// Save the record
		if err := PbClient.Save(record); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Generate auth token
		token, err := record.NewAuthToken()
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// Set cookie with auth token
		http.SetCookie(w, &http.Cookie{
			Name:     "pb_auth",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(24 * time.Hour),
		})

		result = map[string]any{
			"token": token,
			"user":  record.PublicExport(),
		}

	case "refresh":
		// Get token from request
		token, ok := formData["token"].(string)
		if !ok {
			// Try to get from auth cookie
			cookie, err := r.Cookie("pb_auth")
			if err != nil || cookie.Value == "" {
				http.Error(w, "Missing token", http.StatusBadRequest)
				return
			}
			token = cookie.Value
		}

		// Find the auth record by token
		record, err := PbClient.FindAuthRecordByToken(token, core.TokenTypeAuth)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Generate new token
		newToken, err := record.NewAuthToken()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the new token cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "pb_auth",
			Value:    newToken,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(24 * time.Hour),
		})

		result = map[string]any{
			"token": newToken,
			"user":  record.PublicExport(),
		}

	default:
		http.Error(w, "Unsupported auth action", http.StatusBadRequest)
		return
	}

	// Return the result as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// AuthRefresh refreshes an authentication token
func AuthRefresh(w http.ResponseWriter, r *http.Request) {
	// Get auth token from cookie
	cookie, err := r.Cookie("pb_auth")
	if err != nil {
		http.Error(w, "No authentication token", http.StatusBadRequest)
		return
	}

	// Make sure PocketBase client is initialized
	if PbClient == nil {
		http.Error(w, "Authentication system not available", http.StatusInternalServerError)
		return
	}

	// Find the auth record by token
	record, err := PbClient.FindAuthRecordByToken(cookie.Value, core.TokenTypeAuth)
	if err != nil {
		// If token is invalid, clear the cookie and redirect to login
		http.SetCookie(w, &http.Cookie{
			Name:     "pb_auth",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		})
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	// Generate a new token
	newToken, err := record.NewAuthToken()
	if err != nil {
		http.Error(w, "Failed to refresh token", http.StatusInternalServerError)
		return
	}

	// Set cookie with the new token
	http.SetCookie(w, &http.Cookie{
		Name:     "pb_auth",
		Value:    newToken,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token":   newToken,
		"user":    record.PublicExport(),
	})
}
