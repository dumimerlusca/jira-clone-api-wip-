package app

import (
	"encoding/json"
	"fmt"
	"io"
	"jira-clone/packages/db"
	"jira-clone/packages/response"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type registerPayload struct {
	Username string
	Password string
}

func (p *registerPayload) validate() error {
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
	body, err := io.ReadAll(r.Body)

	if err != nil {
		app.badRequest(w, err)
		return
	}

	var payload registerPayload

	err = json.Unmarshal(body, &payload)

	if err != nil {
		app.serverError(w, err)
		return
	}

	err = payload.validate()

	if err != nil {
		app.badRequest(w, err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 10)
	if err != nil {
		app.serverError(w, err)
		return
	}

	userId := uuid.NewString()

	signedToken, err := generateAuthToken(userId, payload.Username)

	if err != nil {
		app.serverError(w, err)
		return
	}

	createdUser, err := app.queries.CreateUser(userId, payload.Username, string(hashedPassword))

	if err != nil {
		app.serverError(w, err)
		return
	}

	type resPayload struct {
		Token string  `json:"token"`
		User  db.User `json:"user"`
	}

	response.JSONWithHeaders(w, http.StatusCreated, resPayload{Token: signedToken, User: *createdUser})
}

type loginPayload struct {
	Username string
	Password string
}

func (p *loginPayload) validate() error {
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
		app.badRequest(w, err)
		return
	}

	var payload loginPayload

	err = json.Unmarshal(body, &payload)

	if err != nil {
		app.badRequest(w, err)
		return
	}

	err = payload.validate()

	if err != nil {
		app.badRequest(w, err)
		return
	}

	user, err := app.queries.FindUserByUsername(payload.Username, true)

	if err != nil {
		app.badRequest(w, fmt.Errorf("bad credentials"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))

	if err != nil {
		app.badRequest(w, fmt.Errorf("bad credentials"))
		return
	}

	signedToken, err := generateAuthToken(user.Id, payload.Username)

	if err != nil {
		app.serverError(w, err)
		return
	}

	response.JSONWithHeaders(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{Token: signedToken})
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
