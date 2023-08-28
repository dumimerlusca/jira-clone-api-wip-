package util

import (
	"encoding/json"
	"io"
)

func ReadAndUnmarshal(r io.Reader, pointerToV any) error {
	data, err := io.ReadAll(r)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, pointerToV)

	return err
}
