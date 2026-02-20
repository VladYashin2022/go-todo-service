package httpServer

import (
	"cli_todo/service"
	"cli_todo/storage"
	"encoding/json"
	"net/http"
)

func Run(addr string) error {

	//создаем router
	mux := http.NewServeMux()

	//создаем handlers
	mux.HandleFunc("/tasks", handler)

	//слушаем порт
	err := http.ListenAndServe(addr, mux)

	return err

}

// handler
func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetTasks(w, r)
	case http.MethodPost:
		handleCreateTask(w, r)
	default:
		http.Error(w, "handler", http.StatusMethodNotAllowed)
	}

}

// GET
func handleGetTasks(w http.ResponseWriter, r *http.Request) {
	data, err := storage.CreateJson(service.AllTasks)
	if err != nil {
		http.Error(w, "json create error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// POST
func handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var req requestJson
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "json decode error", http.StatusBadRequest)
		return
	}
	if req.Name == "" && req.Date == "" {
		http.Error(w, "empty parameter in request", http.StatusBadRequest)
		return
	}

	_, err = service.CreateTask(req.Name, req.Date)
	if err != nil {
		http.Error(w, "create task error", http.StatusInternalServerError)
	}
	w.Write([]byte("The task was created successfully."))

}

type requestJson struct {
	Name string `json:"name"`
	Date string `json:"date"`
}
