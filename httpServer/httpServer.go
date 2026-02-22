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
	case http.MethodPut:
		handlePutTask(w, r)
	case http.MethodPatch:
		handlePatchTask(w, r)
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
	var req requestTask
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
		http.Error(w, "create task error", http.StatusBadRequest)
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
		http.Error(w, "no id", http.StatusBadRequest)
		return
	} else {
		idTask, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "URL query conv error", http.StatusBadRequest)
			return
		}
		err = service.DeleteTask(idTask, &service.AllTasks)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

// PUT
func handlePutTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "no ID in request", http.StatusBadRequest)
		return
	}

	idTask, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "URL query conv error", http.StatusBadRequest)
		return
	}
	//декодируем body в структуру
	var req requestTask
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "read body error", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Date == "" {
		http.Error(w, "Name or date in request body is empty", http.StatusBadRequest)
		return
	} else {
		err = service.UpdateAllTask(idTask, req.Name, req.Date, &service.AllTasks)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		err = storage.JsonUpdate(service.AllTasks) //обновляем json если нет ошибок
		if err != nil {
			http.Error(w, "update json error", http.StatusInternalServerError)
			return
		}

		updatedTaskJson, err := service.FindTaskJson(idTask, service.AllTasks)
		if err == service.ErrNotExist {
			http.Error(w, "Not exist", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "marshal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(updatedTaskJson)

	}
}

// PATCH
func handlePatchTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "no ID in request", http.StatusBadRequest)
		return
	}

	idTask, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "URL query conv error", http.StatusBadRequest)
		return
	}

	var req patchTask
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "read body error", http.StatusBadRequest)
		return
	}

	if req.Name == nil && req.Date == nil {
		http.Error(w, "empty request", http.StatusBadRequest)
		return
	}
	//validation
	if req.Name != nil && *req.Name == "" {
		http.Error(w, "empty field in request", http.StatusBadRequest)
		return
	}
	if req.Date != nil && *req.Date == "" {
		http.Error(w, "empty field in request", http.StatusBadRequest)
		return
	}

	//update name
	if req.Name != nil && *req.Name != "" {
		reqName := *req.Name
		err = service.UpdateName(idTask, reqName, &service.AllTasks)
		if err != nil {
			http.Error(w, "update task error", http.StatusBadRequest)
			return
		}
	}
	//update date
	if req.Date != nil && *req.Date != "" {
		reqDate := *req.Date
		err = service.UpdateDate(idTask, reqDate, &service.AllTasks)
		if err != nil {
			http.Error(w, "update task error", http.StatusBadRequest)
			return
		}
	}

	err = storage.JsonUpdate(service.AllTasks) //обновляем json если нет ошибок
	if err != nil {
		http.Error(w, "update json error", http.StatusInternalServerError)
		return
	}

	updatedTaskJson, err := service.FindTaskJson(idTask, service.AllTasks)
	if err == service.ErrNotExist {
		http.Error(w, "Not exist", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "marshal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(updatedTaskJson)
}

type requestTask struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

type patchTask struct {
	Name *string `json:"name"`
	Date *string `json:"date"`
}
