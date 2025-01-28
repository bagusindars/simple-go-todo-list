package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bagusindars/simple-go-todo-list/api"
	"github.com/bagusindars/simple-go-todo-list/handlers"
)

func main() {
	port := ":8090"

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetTodos(w, r)
		case http.MethodPost:
			handlers.CreateTodo(w, r)
		default:
			api.WriteResponse(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
		}
	})

	http.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			handlers.UpdateTodoStatus(w, r)
		case http.MethodDelete:
			handlers.DeleteTodo(w, r)
		default:
			api.WriteResponse(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
		}
	})

	fmt.Println("Server started on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
