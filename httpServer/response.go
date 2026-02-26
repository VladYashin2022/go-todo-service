package httpServer

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJson(w http.ResponseWriter, s int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s)

	if v == nil {
		return
	}

	err := json.NewEncoder(w).Encode(response{Data: v})
	if err != nil {
		return
	}
}

func WriteError(w http.ResponseWriter, msg string, s int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s)

	err := json.NewEncoder(w).Encode(response{
		Error: &apiError{Message: msg},
	})
	if err != nil {
		log.Panicln(err)
	}
}

type response struct {
	Data  any       `json:"data,omitempty"`
	Error *apiError `json:"error,omitempty"`
}

type apiError struct {
	Message string `json:"message"`
}
