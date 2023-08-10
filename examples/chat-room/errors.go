package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func ResponseErr(w http.ResponseWriter, httpCode int, err error) {
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
	log.Printf("ResponseErr: %v", err)
}
