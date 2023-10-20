package main

import (
	"encoding/json"
	"net/http"
)

func parseBody[T any](w http.ResponseWriter, r *http.Request) T {
	var data T
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return data
	}

	defer r.Body.Close()

	return data
}
