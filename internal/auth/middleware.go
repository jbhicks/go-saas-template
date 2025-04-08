package auth

import (
	"context"
	"net/http"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// PocketBase client instance will be set from main
var PbClient *pocketbase.PocketBase

// User context key
type contextKey string

const userContextKey contextKey = "user"

// AuthMiddleware checks if user is authenticated
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get auth token from cookie
		cookie, err := r.Cookie("pb_auth")

		// If no cookie or error, redirect to login
		if err != nil || cookie == nil || cookie.Value == "" {
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}

		// Verify token using PocketBase
		token := cookie.Value

		// According to go-records documentation, use FindAuthRecordByToken
		authRecord, err := PbClient.FindAuthRecordByToken(token, core.TokenTypeAuth)
		if err != nil || authRecord == nil {
			// Invalid token, redirect to login
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}

		// Store user in request context
		ctx := context.WithValue(r.Context(), userContextKey, authRecord)

		// Continue to next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetCurrentUser returns the current authenticated user or nil
func GetCurrentUser(r *http.Request) *core.Record {
	// Get user from request context
	user, ok := r.Context().Value(userContextKey).(*core.Record)
	if ok {
		return user
	}

	// No user in context, try to authenticate with cookie
	cookie, err := r.Cookie("pb_auth")
	if err != nil || cookie == nil || cookie.Value == "" {
		return nil
	}

	// Verify token using PocketBase
	token := cookie.Value

	// According to go-records documentation, use FindAuthRecordByToken
	authRecord, err := PbClient.FindAuthRecordByToken(token, core.TokenTypeAuth)
	if err != nil || authRecord == nil {
		return nil
	}

	return authRecord
}
