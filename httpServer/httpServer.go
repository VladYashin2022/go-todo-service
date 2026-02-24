package httpServer

import (
	"cli_todo/service"
	"cli_todo/storage"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func Run(addr string) error {

	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", handler)

	err := http.ListenAndServe(addr, Logging(mux))
	return err

}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s,%s,%s", r.Method, r.URL.Query(), r.URL.RawQuery)
		next.ServeHTTP(w, r)
	})
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
		WriteJson(w, http.StatusOK, service.AllTasks)
		return
	} else {
		//GET by ID
		idTask, err := strconv.Atoi(idStr)
		if err != nil {
			WriteError(w, "URL query conv error", http.StatusBadRequest)
			return
		}
		jsonTask, err := service.FindTask(idTask, service.AllTasks)
		if err != nil {
			WriteError(w, "find task error", http.StatusNotFound)
			return
		}

		WriteJson(w, http.StatusOK, jsonTask)
	}
}

// POST
func handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var req requestTask
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		WriteError(w, "json decode error", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Date == "" {
		WriteError(w, "empty parameter in request", http.StatusBadRequest)
		return
	}

	task, err := service.CreateTask(req.Name, req.Date)
	if err != nil {
		WriteError(w, "create task error", http.StatusBadRequest)
		return
	}

	WriteJson(w, http.StatusCreated, task)
}

// DELETE
func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		WriteError(w, "no id", http.StatusBadRequest)
		return
	} else {
		idTask, err := strconv.Atoi(idStr)
		if err != nil {
			WriteError(w, "URL query conv error", http.StatusBadRequest)
			return
		}
		err = service.DeleteTask(idTask, &service.AllTasks)
		if err != nil {
			WriteError(w, err.Error(), http.StatusNotFound)
			return
		}

		WriteJson(w, http.StatusNoContent, nil)
	}
}

// PUT
func handlePutTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		WriteError(w, "no ID in request", http.StatusBadRequest)
		return
	}

	idTask, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, "URL query conv error", http.StatusBadRequest)
		return
	}
	//декодируем body в структуру
	var req requestTask
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		WriteError(w, "read body error", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Date == "" {
		WriteError(w, "Name or date in request body is empty", http.StatusBadRequest)
		return
	} else {
		err = service.UpdateAllTask(idTask, req.Name, req.Date, &service.AllTasks)
		if err != nil {
			WriteError(w, err.Error(), http.StatusNotFound)
			return
		}

		err = storage.JsonUpdate(service.AllTasks) //обновляем json если нет ошибок
		if err != nil {
			WriteError(w, "update json error", http.StatusInternalServerError)
			return
		}

		updatedTaskJson, err := service.FindTask(idTask, service.AllTasks)
		if err != nil {
			WriteError(w, "Not exist", http.StatusNotFound)
			return
		}

		WriteJson(w, http.StatusOK, updatedTaskJson)
	}
}

// PATCH
func handlePatchTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		WriteError(w, "no ID in request", http.StatusBadRequest)
		return
	}

	idTask, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, "URL query conv error", http.StatusBadRequest)
		return
	}

	var req patchTask
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		WriteError(w, "read body error", http.StatusBadRequest)
		return
	}

	if req.Name == nil && req.Date == nil {
		WriteError(w, "empty request", http.StatusBadRequest)
		return
	}
	//validation
	if req.Name != nil && *req.Name == "" {
		WriteError(w, "empty field in request", http.StatusBadRequest)
		return
	}
	if req.Date != nil && *req.Date == "" {
		WriteError(w, "empty field in request", http.StatusBadRequest)
		return
	}

	//update name
	if req.Name != nil && *req.Name != "" {
		reqName := *req.Name
		err = service.UpdateName(idTask, reqName, &service.AllTasks)
		if err != nil {
			WriteError(w, "update task error", http.StatusBadRequest)
			return
		}
	}
	//update date
	if req.Date != nil && *req.Date != "" {
		reqDate := *req.Date
		err = service.UpdateDate(idTask, reqDate, &service.AllTasks)
		if err != nil {
			WriteError(w, "update task error", http.StatusBadRequest)
			return
		}
	}

	err = storage.JsonUpdate(service.AllTasks) //обновляем json если нет ошибок
	if err != nil {
		WriteError(w, "update json error", http.StatusInternalServerError)
		return
	}

	updatedTaskJson, err := service.FindTask(idTask, service.AllTasks)
	if err != nil {
		WriteError(w, "Not exist", http.StatusInternalServerError)
		return
	}

	WriteJson(w, http.StatusOK, updatedTaskJson)
}

type requestTask struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

type patchTask struct {
	Name *string `json:"name"`
	Date *string `json:"date"`
}
