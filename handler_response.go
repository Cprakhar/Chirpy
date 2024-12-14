package main

import (
	"encoding/json"
	"net/http"
)

func responseWithJSON(w http.ResponseWriter, code int, payload interface{}){
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}