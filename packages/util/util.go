package util

import (
	"encoding/json"
	"io"
	"net/http"
)

func ReadAndUnmarshal(r io.Reader, pointerToV any) error {
	data, err := io.ReadAll(r)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, pointerToV)

	return err
}

func GetQueryParameter(name string, r *http.Request) string {
	return r.URL.Query().Get(name)
}
