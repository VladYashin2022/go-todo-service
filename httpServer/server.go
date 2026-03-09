package httpServer

import (
	"cli_todo/model"
	"cli_todo/storage/postgres"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	repo       *postgres.TaskRepository
	httpServer *http.Server
}

func New(repo *postgres.TaskRepository) *Server {
	return &Server{
		repo: repo,
	}
}

func (s *Server) Run(addr string) error {

	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", s.handler)

	s.httpServer = &http.Server{
		Addr:    ":8080",
		Handler: Logging(mux),
	}

	return s.httpServer.ListenAndServe()

}

// server shutdown func
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// handler
func (s Server) handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetTasks(w, r)
	case http.MethodPost:
		s.handleCreateTask(w, r)
	case http.MethodDelete:
		s.handleDeleteTask(w, r)
	case http.MethodPut:
		s.handlePutTask(w, r)
	case http.MethodPatch:
		s.handlePatchTask(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

// GET
func (s *Server) handleGetTasks(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		//GET all
		tasks, err := s.repo.GetAllTasks()
		if err != nil {
			WriteError(w, "get tasks err", http.StatusInternalServerError)
			return
		}

		WriteJson(w, http.StatusOK, tasks)
		return
	} else {
		//GET by ID
		idTask, err := strconv.Atoi(idStr)
		if err != nil {
			WriteError(w, "URL query conv error", http.StatusBadRequest)
			return
		}

		jsonTask, err := s.repo.GetTaskByID(idTask)
		if errors.Is(err, postgres.ErrTaskNotFound) {
			WriteError(w, "find task error", http.StatusNotFound)
			return
		}
		if err != nil {
			WriteError(w, "get task error", http.StatusInternalServerError)
			return
		}

		WriteJson(w, http.StatusOK, jsonTask)
	}
}

// POST
func (s *Server) handleCreateTask(w http.ResponseWriter, r *http.Request) {
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

	dateTask, err := model.ReadDateTask(req.Date)
	if err != nil {
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println(dateTask)
	task, err := s.repo.CreateTask(req.Name, dateTask)
	if err != nil {
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJson(w, http.StatusCreated, task)
}

// DELETE
func (s *Server) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		WriteError(w, "no id", http.StatusBadRequest)
		return
	}

	idTask, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, "URL query conv error", http.StatusBadRequest)
		return
	}

	err = s.repo.DeleteTask(idTask)
	if errors.Is(err, postgres.ErrTaskNotFound) {
		WriteError(w, "task not found", http.StatusNotFound)
		return
	}
	if err != nil {
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PUT
func (s *Server) handlePutTask(w http.ResponseWriter, r *http.Request) {
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

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, "read body error", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Date == "" {
		WriteError(w, "Name or date in request body is empty", http.StatusBadRequest)
		return
	}

	//конвертируем строку в time.Time
	date, err := model.ReadDateTask(req.Date)
	if err != nil {
		WriteError(w, "invalid date in request", http.StatusBadRequest)
		return
	}

	task, err := s.repo.UpdateTask(idTask, req.Name, date)
	if errors.Is(err, postgres.ErrTaskNotFound) {
		WriteError(w, "task not found", http.StatusNotFound)
		return
	}
	if err != nil {
		WriteError(w, "task update error", http.StatusInternalServerError)
		return
	}
	WriteJson(w, http.StatusOK, task)
}

// PATCH
func (s *Server) handlePatchTask(w http.ResponseWriter, r *http.Request) {
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

	//validation
	if req.Name == nil && req.Date == nil {
		WriteError(w, "empty request", http.StatusBadRequest)
		return
	}
	if req.Name != nil && *req.Name == "" {
		WriteError(w, "empty field in request", http.StatusBadRequest)
		return
	}
	if req.Date != nil && *req.Date == "" {
		WriteError(w, "empty field in request", http.StatusBadRequest)
		return
	}

	//parse date
	var parsedDate *time.Time
	if req.Date != nil {
		t, err := model.ReadDateTask(*req.Date)

		if err != nil {
			WriteError(w, "wrong date format", http.StatusBadRequest)
			return
		}

		parsedDate = &t
	}

	task, err := s.repo.PatchTask(idTask, req.Name, parsedDate)

	if errors.Is(err, postgres.ErrTaskNotFound) {
		WriteError(w, "task not found", http.StatusNotFound)
		return
	}
	if err != nil {
		WriteError(w, "patch task error", http.StatusInternalServerError)
		return
	}

	WriteJson(w, http.StatusOK, task)
}

type requestTask struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

type patchTask struct {
	Name *string `json:"name"`
	Date *string `json:"date"`
}