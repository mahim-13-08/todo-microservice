package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/mahim-13-08/todo-microservice/internal/model"
	"github.com/mahim-13-08/todo-microservice/internal/repository"
	"github.com/mahim-13-08/todo-microservice/internal/service"
	"github.com/mahim-13-08/todo-microservice/pkg/response"
	"github.com/mahim-13-08/todo-microservice/pkg/validator"
)

// TodoHandler handles HTTP requests for the Todo resource.
type TodoHandler struct {
	svc service.TodoService
}

// NewTodoHandler creates a new TodoHandler.
func NewTodoHandler(svc service.TodoService) *TodoHandler {
	return &TodoHandler{svc: svc}
}

// RegisterRoutes wires all Todo routes onto the given mux.
// Routes:
//
//	POST   /api/v1/todos        – create
//	GET    /api/v1/todos        – list all
//	GET    /api/v1/todos/{id}   – get by id
//	PUT    /api/v1/todos/{id}   – full update
//	PATCH  /api/v1/todos/{id}   – partial update
//	DELETE /api/v1/todos/{id}   – delete
func (h *TodoHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/todos", h.handleCollection)
	mux.HandleFunc("/api/v1/todos/", h.handleSingle)
}

// handleCollection routes GET / POST on /api/v1/todos.
func (h *TodoHandler) handleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listTodos(w, r)
	case http.MethodPost:
		h.createTodo(w, r)
	default:
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// handleSingle routes GET / PUT / PATCH / DELETE on /api/v1/todos/{id}.
func (h *TodoHandler) handleSingle(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid todo id")
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.getTodo(w, r, id)
	case http.MethodPut, http.MethodPatch:
		h.updateTodo(w, r, id)
	case http.MethodDelete:
		h.deleteTodo(w, r, id)
	default:
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// createTodo handles POST /api/v1/todos.
func (h *TodoHandler) createTodo(w http.ResponseWriter, r *http.Request) {
	var req model.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		response.Error(w, http.StatusUnprocessableEntity, validator.ValidationErrors(err))
		return
	}
	todo, err := h.svc.CreateTodo(r.Context(), &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create todo")
		return
	}
	response.Success(w, http.StatusCreated, "todo created successfully", todo)
}

// listTodos handles GET /api/v1/todos.
func (h *TodoHandler) listTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.svc.GetAllTodos(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to retrieve todos")
		return
	}
	if todos == nil {
		todos = []*model.Todo{} // return empty array, not null
	}
	response.Success(w, http.StatusOK, "todos retrieved successfully", todos)
}

// getTodo handles GET /api/v1/todos/{id}.
func (h *TodoHandler) getTodo(w http.ResponseWriter, r *http.Request, id int64) {
	todo, err := h.svc.GetTodoByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "todo not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to retrieve todo")
		return
	}
	response.Success(w, http.StatusOK, "todo retrieved successfully", todo)
}

// updateTodo handles PUT/PATCH /api/v1/todos/{id}.
func (h *TodoHandler) updateTodo(w http.ResponseWriter, r *http.Request, id int64) {
	var req model.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		response.Error(w, http.StatusUnprocessableEntity, validator.ValidationErrors(err))
		return
	}
	todo, err := h.svc.UpdateTodo(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "todo not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to update todo")
		return
	}
	response.Success(w, http.StatusOK, "todo updated successfully", todo)
}

// deleteTodo handles DELETE /api/v1/todos/{id}.
func (h *TodoHandler) deleteTodo(w http.ResponseWriter, r *http.Request, id int64) {
	if err := h.svc.DeleteTodo(r.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "todo not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to delete todo")
		return
	}
	response.Success(w, http.StatusOK, "todo deleted successfully", nil)
}

// extractID parses the trailing segment of a URL path as an int64.
func extractID(path string) (int64, error) {
	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
	raw := parts[len(parts)-1]
	return strconv.ParseInt(raw, 10, 64)
}
