package httpServer

import (
	"cli_todo/service"
	"cli_todo/storage"
	"encoding/json"
	"net/http"
	"strconv"
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
	case http.MethodDelete:
		handleDeleteTask(w, r)
	default:
		http.Error(w, "handler", http.StatusMethodNotAllowed)
	}

}

// GET
func handleGetTasks(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		//GET all
		data, err := storage.CreateJson(service.AllTasks)
		if err != nil {
			http.Error(w, "json create error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	} else {
		//GET by ID
		idTask, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "URL query conv error", http.StatusBadRequest)
			return
		}
		strTask, err := service.ReadTask(idTask, service.AllTasks)
		if err != nil {
			http.Error(w, "read task error", http.StatusNotFound)
			return
		}
		jsonTask, err := json.Marshal(strTask)
		if err != nil {
			http.Error(w, "json marshal error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonTask)
	}
}

// POST
func handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var req requestJson
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "json decode error", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Date == "" {
		http.Error(w, "empty parameter in request", http.StatusBadRequest)
		return
	}

	task, err := service.CreateTask(req.Name, req.Date)
	if err != nil {
		http.Error(w, "create task error", http.StatusInternalServerError)
		return
	}

	jsonTask, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "json marshal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte(jsonTask))

}

// DELETE
func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "no id", http.StatusNotFound)
		return
	} else {
		idTask, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "URL query conv error", http.StatusBadRequest)
			return
		}
		err = service.DeleteTask(idTask, &service.AllTasks)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

type requestJson struct {
	Name string `json:"name"`
	Date string `json:"date"`
}
