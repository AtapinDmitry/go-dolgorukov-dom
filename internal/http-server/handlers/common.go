package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func DecodeJSONBody(r *http.Request, v interface{}) error {
	byteResult, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("could not read body: %w", err)
	}

	if err := json.Unmarshal(byteResult, v); err != nil {
		return fmt.Errorf("could not unmarshal body: %w", err)
	}

	return nil
}
