package app

import (
	"encoding/json"
	"fmt"
	"io"
	"jira-clone/packages/models"
	"jira-clone/packages/response"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type registerRequestPayload struct {
	Username string
	Password string
}

type registerResponsePayload struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func (p *registerRequestPayload) validate() error {
	if p.Username == "" {
		return fmt.Errorf("username required")
	}

	if p.Password == "" {
		return fmt.Errorf("password required")
	}

	if len(p.Password) <= 6 {
		return fmt.Errorf("password should have more than 6 characters")
	}

	return nil
}

func (app *application) registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		app.badRequest(w, "body must not be empty", nil)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		app.badRequest(w, "error decoding request body", err)
		return
	}

	var payload registerRequestPayload

	err = json.Unmarshal(body, &payload)

	if err != nil {
		app.badRequest(w, "error decoding request body", err)
		return
	}

	err = payload.validate()

	if err != nil {
		app.badRequest(w, err.Error(), err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 10)
	if err != nil {
		app.serverError(w, "", err)
		return
	}

	userId := uuid.NewString()

	signedToken, err := generateAuthToken(userId, payload.Username)

	if err != nil {
		app.serverError(w, "", err)
		return
	}

	createdUser, err := app.queries.CreateUser(userId, payload.Username, string(hashedPassword))

	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == pgerrcode.UniqueViolation {
				app.badRequest(w, "username already taken", err)
				return
			}
		}

		app.serverError(w, "", err)
		return
	}

	response.JSONWithHeaders(w, http.StatusCreated, registerResponsePayload{Token: signedToken, User: *createdUser})
}

type loginRequestPayload struct {
	Username string
	Password string
}

type loginResponsePayload struct {
	Token string `json:"token"`
}

func (p *loginRequestPayload) validate() error {
	if p.Password == "" {
		return fmt.Errorf("password required")
	}
	if p.Username == "" {
		return fmt.Errorf("username required")
	}
	return nil
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		app.badRequest(w, "error reading req body", err)
		return
	}

	var payload loginRequestPayload

	err = json.Unmarshal(body, &payload)

	if err != nil {
		app.badRequest(w, "error decoding req body", err)
		return
	}

	err = payload.validate()

	if err != nil {
		app.badRequest(w, err.Error(), err)
		return
	}

	user, err := app.queries.FindUserByUsername(payload.Username, true)

	if err != nil {
		app.unauthorizedRequest(w, "bad credentials", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))

	if err != nil {
		app.unauthorizedRequest(w, "bad credentials", err)
		return
	}

	signedToken, err := generateAuthToken(user.Id, payload.Username)

	if err != nil {
		app.serverError(w, "", err)
		return
	}

	response.JSONWithHeaders(w, http.StatusOK, loginResponsePayload{Token: signedToken})
}

func generateAuthToken(userId string, username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()
	claims["userId"] = userId
	claims["username"] = username

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return signedToken, err
}
