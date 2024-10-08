package app

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

type contextKey string

func ContextKey(key string) contextKey {
	return contextKey(key)
}

func (app *application) authMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authH := strings.Split(r.Header.Get("Authorization"), " ")

		if len(authH) != 2 {
			app.unauthorizedRequest(w, "", nil)
			return
		}

		token, err := jwt.Parse(authH[1], func(token *jwt.Token) (interface{}, error) {

			secret := []byte(os.Getenv("JWT_SECRET"))

			return secret, nil
		})

		if err != nil {
			app.unauthorizedRequest(w, "", err)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			app.unauthorizedRequest(w, "", err)
			return
		}

		userId := claims["userId"].(string)

		ctx := context.WithValue(r.Context(), ContextKey("userId"), userId)

		r = r.WithContext(ctx)

		next(w, r)
	}
}

func (app *application) isProjectOwnerMW(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(ContextKey("userId"))
		projectId := mux.Vars(r)["projectId"]

		if _, ok := userId.(string); !ok {
			app.unauthorizedRequest(w, "", nil)
			return
		}

		var created_by_id string

		row := app.db.QueryRow(`SELECT created_by_id from projects WHERE id=$1 LIMIT 1`, projectId)

		err := row.Scan(&created_by_id)

		if err != nil {
			app.serverError(w, "", err)
			return
		}

		if userId != created_by_id {
			app.unauthorizedRequest(w, "current logged in user is not project leader", err)
			return
		}

		handler(w, r)
	}
}

// Works for routes containing {projectId} or {ticketId}
func (app *application) isProjectMemberMW(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(ContextKey("userId")).(string)
		projectId := mux.Vars(r)["projectId"]
		ticketId := mux.Vars(r)["ticketId"]

		if projectId == "" && ticketId != "" {
			ticket, err := app.queries.FindTicketById(ticketId)

			if err != nil {
				app.badRequest(w, err.Error(), nil)
				return
			}

			projectId = ticket.Project_id
		}

		if projectId == "" || userId == "" {
			app.badRequest(w, "projectId or userId or ticketId missing", nil)
			return
		}

		isMember, err := app.queries.IsProjectMember(userId, projectId)

		if err != nil {
			app.serverError(w, err.Error(), err)
			return
		}

		if !isMember {
			app.unauthorizedRequest(w, "you can't modify this project since you are not a member of it", nil)
			return
		}

		handler(w, r)
	}
}
