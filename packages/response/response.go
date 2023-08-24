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
