package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type Todo struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func newTodo(t map[string]Todo, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	newId := uuid.NewString()

	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var todo Todo
	if err := json.Unmarshal(body, &todo); err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}

	t[newId] = todo

	fmt.Printf("> Got request data: %T %+v\n", todo, todo)
	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

func getTodos(t map[string]Todo, w http.ResponseWriter, r *http.Request) {
	jsonResponse, err := json.Marshal(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func main() {
	todos := make(map[string]Todo)
	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			newTodo(todos, w, r)
		case http.MethodGet:
			getTodos(todos, w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Not Found"))
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
