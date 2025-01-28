package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/bagusindars/simple-go-todo-list/api"
	"github.com/bagusindars/simple-go-todo-list/models"
)

var (
	todos  = []models.Todo{}
	nextId = 1
	todoMu sync.Mutex
)

type CreateTodoRequest struct {
	Title string `json:"title"`
}

type UpdateTodoStatusRequest struct {
	IsDone bool `json:"isDone"`
}

// GetTodos handle GET /todos
func GetTodos(w http.ResponseWriter, r *http.Request) {
	api.WriteResponse(w, "Todo loaded", 200, todos)
}

// CreateTodo handle POST /todos
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var data CreateTodoRequest

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		api.WriteResponse(w, "Invalid request body", http.StatusBadRequest, nil)
		return
	}

	if len(data.Title) == 0 {
		api.WriteResponse(w, "Title is required ", http.StatusBadRequest, nil)
		return
	}

	todoMu.Lock()

	var newTodo models.Todo = models.Todo{
		ID:     nextId,
		Title:  data.Title,
		IsDone: false,
	}

	todos = append(todos, newTodo)
	nextId++

	todoMu.Unlock()

	api.WriteResponse(w, "New todo successfully created", http.StatusCreated, newTodo)
}

// DeleteTodo handle DELETE /todos/{id}
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/todos/"):]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		api.WriteResponse(w, "Invalid ID", http.StatusBadRequest, nil)
		return
	}

	todoMu.Lock()
	defer todoMu.Unlock()

	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			api.WriteResponse(w, "Todo has been deleted", http.StatusOK, nil)
			return
		}
	}

	api.WriteResponse(w, "Todo not found", http.StatusNotFound, nil)
}

// UpdateTodoStatus handle PUT /todos/{id}
func UpdateTodoStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/todos/"):]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		api.WriteResponse(w, "Invalid ID", http.StatusBadRequest, nil)
		return
	}

	var data UpdateTodoStatusRequest

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		api.WriteResponse(w, "Invalid input", http.StatusBadRequest, nil)
		return
	}

	todoMu.Lock()
	defer todoMu.Unlock()

	for i, todo := range todos {
		if todo.ID == id {
			todos[i].IsDone = data.IsDone
			api.WriteResponse(w, "Todo has been updated", http.StatusOK, todos[i])
			return
		}
	}

	api.WriteResponse(w, "Todo not found", http.StatusNotFound, nil)
}
