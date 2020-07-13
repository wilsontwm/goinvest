package middleware

import (
	"context"
	"github.com/gorilla/mux"
	"goinvest/models"
	"goinvest/utils"
	"net/http"
	"strings"
)

// JwtAuthentication : Authenticate the authorization token in header
var JwtAuthentication = func() mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := make(map[string]interface{})
			// Check for authentication
			tokenHeader := r.Header.Get("Authorization")

			// If token is missing, then return error code 403 Unauthorized
			if tokenHeader == "" {
				utils.Fail(w, http.StatusUnauthorized, resp, "missing authorization token")
				return
			}

			// Check if the token format is correct, ie. Bearer {token}
			splitted := strings.Split(tokenHeader, " ")
			if len(splitted) != 2 {
				utils.Fail(w, http.StatusUnauthorized, resp, "invalid authorization token format")
				return
			}

			tokenPart := splitted[1] // Grab the second part
			token, err := models.UserAuthenticate(tokenPart)

			if err != nil {
				utils.Fail(w, http.StatusUnauthorized, resp, err.Error())
				return
			}

			// Set the user ID in the context
			ctx := context.Background()
			ctx = context.WithValue(ctx, models.ContextKeyUserID, token.UserID)
			r = r.WithContext(ctx)
			handler.ServeHTTP(w, r)
		})
	}
}
