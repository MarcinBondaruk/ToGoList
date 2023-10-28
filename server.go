package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type ApiTodo struct {
	Title    string `json:"title"`
	Contents string `json:"contents"`
}

type Todo struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Contents string `json:"contents"`
}

func saveJson(t map[string]Todo) {
	todos, err := json.Marshal(t)
	if err != nil {
		fmt.Println("Failed to serialize Todos: ", err)
	}
	err = os.WriteFile("todos.json", todos, 0644)
	if err != nil {
		fmt.Println("Error writing to file: ", err)
	}
}

func handleNewTodo(t map[string]Todo, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var apiTodo ApiTodo
	if err := json.Unmarshal(body, &apiTodo); err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}

	// utworzyc todo i dodac do mapy
	newTodo := Todo{
		Id:       uuid.NewString(),
		Title:    apiTodo.Title,
		Contents: apiTodo.Contents,
	}

	t[newTodo.Id] = newTodo

	saveJson(t)

	fmt.Printf("> Got request data: %T %+v\n", t, t)
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

func retrieveTodos() map[string]Todo {
	todos := make(map[string]Todo)
	data, err := os.ReadFile("todos.json")
	if err != nil {
		fmt.Println("Failed to read todos.json")
		// return exception
	}

	if err := json.Unmarshal(data, &todos); err != nil {
		fmt.Println("Failed to unmarshal json data")
	}
	return todos
}

func main() {
	todos := retrieveTodos()
	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleNewTodo(todos, w, r)
		case http.MethodGet:
			getTodos(todos, w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Not Found"))
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
