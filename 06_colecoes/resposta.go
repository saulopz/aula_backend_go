package main

import (
	"encoding/json"
	"net/http"
)

func responderJSON(w http.ResponseWriter, status int, dado any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(dado)
}
