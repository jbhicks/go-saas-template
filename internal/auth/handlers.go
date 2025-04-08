package auth

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"fmt"

	"github.com/gorilla/mux"
	"github.com/pocketbase/pocketbase/core"
)

// Templates for auth pages
var templates = template.Must(template.ParseFiles(
	filepath.Join("internal", "templates", "login.html"),
	filepath.Join("internal", "templates", "register.html"),
	filepath.Join("internal", "templates", "forgot_password.html"),
	filepath.Join("internal", "templates", "reset_password.html"),
	filepath.Join("internal", "templates", "home.html"),
))

// LoginForm represents the login form data
type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Error    string `json:"error,omitempty"`
	Success  string `json:"success,omitempty"`
}

// RegisterForm represents the registration form data
type RegisterForm struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Error           string `json:"error,omitempty"`
	Success         string `json:"success,omitempty"`
}

// ForgotPasswordForm represents the forgot password form data
type ForgotPasswordForm struct {
	Email   string `json:"email"`
	Error   string `json:"error,omitempty"`
	Success string `json:"success,omitempty"`
}

// ResetPasswordForm represents the reset password form data
type ResetPasswordForm struct {
	Token           string `json:"token"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Error           string `json:"error,omitempty"`
	Success         string `json:"success,omitempty"`
}

// HomeData represents the data for the home page
type HomeData struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

// LoginHandler shows the login form
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Check for password reset success message
		resetSuccess := r.URL.Query().Get("reset_success")
		if resetSuccess == "true" {
			templates.ExecuteTemplate(w, "login.html", LoginForm{
				Success: "Your password has been reset successfully. You can now log in with your new password.",
			})
			return
		}
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
			Error: "Invalid email or password. If you forgot your password, use the 'Forgot Password' link below.",
		})
		return
	}

	// Validate password
	if !authRecord.ValidatePassword(password) {
		templates.ExecuteTemplate(w, "login.html", LoginForm{
			Email: email,
			Error: "Invalid email or password. If you forgot your password, use the 'Forgot Password' link below.",
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

	// Check if email already exists
	existingRecord, _ := PbClient.FindAuthRecordByEmail("users", email)
	if existingRecord != nil {
		templates.ExecuteTemplate(w, "register.html", RegisterForm{
			Email: email,
			Error: "An account with this email already exists. Please use the login page or reset your password.",
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

// ForgotPasswordHandler handles password reset requests
func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		templates.ExecuteTemplate(w, "forgot_password.html", nil)
		return
	}

	// Process form submission
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")

	// Validate input
	if email == "" {
		templates.ExecuteTemplate(w, "forgot_password.html", ForgotPasswordForm{
			Error: "Email is required",
		})
		return
	}

	// Make sure PocketBase client is initialized
	if PbClient == nil {
		templates.ExecuteTemplate(w, "forgot_password.html", ForgotPasswordForm{
			Email: email,
			Error: "Password reset system not available",
		})
		return
	}

	// Find user by email
	authRecord, err := PbClient.FindAuthRecordByEmail("users", email)
	if err != nil {
		// Don't reveal whether the email exists or not for security reasons
		templates.ExecuteTemplate(w, "forgot_password.html", ForgotPasswordForm{
			Success: "If an account with this email exists, password reset instructions have been sent.",
		})
		return
	}

	// Get user ID to include in the token for additional security
	userID := authRecord.Id

	// Generate a simple token (in production, use a proper token generation method)
	token := fmt.Sprintf("%s_%d", userID, time.Now().UnixNano())

	// Store the token in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "reset_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600, // 1 hour expiry
	})

	// Store the email in a cookie for the reset process
	http.SetCookie(w, &http.Cookie{
		Name:     "reset_email",
		Value:    email,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600, // 1 hour expiry
	})

	// Create a reset link
	resetLink := fmt.Sprintf("/auth/reset-password?token=%s", token)

	// Show success message with the reset link
	templates.ExecuteTemplate(w, "forgot_password.html", ForgotPasswordForm{
		Success: "Password reset instructions have been sent. For this demo, you can reset your password here: ",
		Email:   resetLink,
	})
}

// ResetPasswordHandler handles password reset form
func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Redirect(w, r, "/auth/forgot-password", http.StatusSeeOther)
			return
		}

		// Verify the token matches the cookie
		resetCookie, err := r.Cookie("reset_token")
		if err != nil || resetCookie.Value != token {
			templates.ExecuteTemplate(w, "reset_password.html", ResetPasswordForm{
				Error: "Invalid or expired reset token. Please request a new password reset.",
			})
			return
		}

		templates.ExecuteTemplate(w, "reset_password.html", ResetPasswordForm{
			Token: token,
		})
		return
	}

	// Process form submission
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	token := r.FormValue("token")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	// Validate inputs
	if token == "" || password == "" {
		templates.ExecuteTemplate(w, "reset_password.html", ResetPasswordForm{
			Token: token,
			Error: "Token and password are required",
		})
		return
	}

	// Validate passwords match
	if password != confirmPassword {
		templates.ExecuteTemplate(w, "reset_password.html", ResetPasswordForm{
			Token: token,
			Error: "Passwords do not match",
		})
		return
	}

	// Make sure PocketBase client is initialized
	if PbClient == nil {
		templates.ExecuteTemplate(w, "reset_password.html", ResetPasswordForm{
			Token: token,
			Error: "Password reset system not available",
		})
		return
	}

	// Verify the token from cookie
	resetCookie, err := r.Cookie("reset_token")
	if err != nil || resetCookie.Value != token {
		templates.ExecuteTemplate(w, "reset_password.html", ResetPasswordForm{
			Token: token,
			Error: "Invalid or expired reset token. Please request a new password reset.",
		})
		return
	}

	// Get the email from cookie
	emailCookie, err := r.Cookie("reset_email")
	if err != nil {
		templates.ExecuteTemplate(w, "reset_password.html", ResetPasswordForm{
			Token: token,
			Error: "Reset session expired. Please request a new password reset.",
		})
		return
	}

	// Find the user by email
	email := emailCookie.Value
	record, err := PbClient.FindAuthRecordByEmail("users", email)
	if err != nil {
		templates.ExecuteTemplate(w, "reset_password.html", ResetPasswordForm{
			Token: token,
			Error: "User not found. Please request a new password reset.",
		})
		return
	}

	// Update the password
	record.SetPassword(password)

	// Save the record
	if err := PbClient.Save(record); err != nil {
		templates.ExecuteTemplate(w, "reset_password.html", ResetPasswordForm{
			Token: token,
			Error: "Failed to update password: " + err.Error(),
		})
		return
	}

	// Clear the reset cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "reset_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "reset_email",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	// Redirect to login page with success message
	http.Redirect(w, r, "/auth/login?reset_success=true", http.StatusSeeOther)
}

// HomeRenderer renders the home page
func HomeRenderer(w http.ResponseWriter, r *http.Request) {
	// Get current authenticated user
	user := GetCurrentUser(r)
	if user == nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	// Get user email from the record
	email := user.Email()

	// Set proper content type and no-cache headers
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Render the home template with user data
	if err := templates.ExecuteTemplate(w, "home.html", HomeData{
		Email: email,
	}); err != nil {
		http.Error(w, "Error rendering home page: "+err.Error(), http.StatusInternalServerError)
	}
}
