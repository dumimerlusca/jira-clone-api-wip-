package app

import (
	"net/http"
)

func extractUserId(r *http.Request) string {
	userId := r.Context().Value(ContextKey("userId")).(string)
	return userId
}
