package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	IsDone bool   `json:"isDone"`
}

type CreateTodoRequest struct {
	Title string `json:"title"`
}

type UpdateTodoStatusRequest struct {
	IsDone bool `json:"isDone"`
}

var (
	todos  = []Todo{}
	nextId = 1
	todoMu sync.Mutex
)

func main() {
	port := ":8090"

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetTodos(w, r)
		case http.MethodPost:
			CreateTodo(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			UpdateTodoStatus(w, r)
		case http.MethodDelete:
			DeleteTodo(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server started on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// GetTodos handle GET /todos
func GetTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// CreateTodo handle POST /todos
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var data CreateTodoRequest

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(data.Title) == 0 {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	todoMu.Lock()

	var newTodo Todo = Todo{
		ID:     nextId,
		Title:  data.Title,
		IsDone: false,
	}

	todos = append(todos, newTodo)
	nextId++

	todoMu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}

// DeleteTodo handle DELETE /todos/{id}
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/todos/"):]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	todoMu.Lock()
	defer todoMu.Unlock()

	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Todo not found", http.StatusNotFound)
}

// UpdateTodoStatus handle PUT /todos/{id}
func UpdateTodoStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/todos/"):]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var data UpdateTodoStatusRequest

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	todoMu.Lock()
	defer todoMu.Unlock()

	for i, todo := range todos {
		if todo.ID == id {
			todos[i].IsDone = data.IsDone
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(todos[i])
			return
		}
	}

	http.Error(w, "Todo not found", http.StatusNotFound)
}
