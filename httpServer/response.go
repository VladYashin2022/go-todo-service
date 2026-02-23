package httpServer

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, s int, v []byte) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s)

	if v == nil {
		return nil
	}

	_, err := w.Write(v)
	if err != nil {
		return err
	}

	return nil
}

func WriteError(w http.ResponseWriter, msg string, s int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s)

	errStruct := responseError{Error: msg, Status: s}
	err := json.NewEncoder(w).Encode(errStruct)
	if err != nil {
		w.Write([]byte(`{"error":"Encoder error"}`))
	}
}

type responseError struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}
