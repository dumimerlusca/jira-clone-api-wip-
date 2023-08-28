package response

import (
	"encoding/json"
	"net/http"
)

func JSONWithHeaders(w http.ResponseWriter, status int, data any) error {
	d, err := json.Marshal(data)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	w.WriteHeader(status)

	w.Write(d)

	return nil
}

type ErrorResponse struct {
	Sucess bool   `json:"success"`
	Error  string `json:"error"`
}

func NewErrorResponse(w http.ResponseWriter, status int, msg string) {
	d := ErrorResponse{Sucess: false, Error: msg}

	err := JSONWithHeaders(w, status, d)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type SuccessResponse struct {
	Sucess bool `json:"success"`
	Data   any  `json:"data"`
}

func NewSuccessResponse(w http.ResponseWriter, status int, data any) {
	d := SuccessResponse{Sucess: true, Data: data}

	err := JSONWithHeaders(w, status, d)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
