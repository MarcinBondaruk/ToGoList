package main

import (
	"encoding/json"
	"errors"
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

func saveToJsonFile(t map[string]Todo) error {
	todos, err := json.Marshal(t)

	if err != nil {
		return errors.New("error serializing todos: " + err.Error())
	}

	err = os.WriteFile("todos.json", todos, 0644)

	if err != nil {
		return errors.New("error writing to file: " + err.Error())
	}

	return nil
}

func createTodo(t map[string]Todo, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	var apiTodo ApiTodo
	if err := json.Unmarshal(body, &apiTodo); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}

	newTodo := Todo{
		Id:       uuid.NewString(),
		Title:    apiTodo.Title,
		Contents: apiTodo.Contents,
	}

	t[newTodo.Id] = newTodo

	if err = saveToJsonFile(t); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Internal Server Error"))
	}

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

func deleteTodos(t map[string]Todo, w http.ResponseWriter, r *http.Request) {
	if err := os.Truncate("todos.json", 0); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Not Found"))
	}

	for k := range t {
		delete(t, k)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

func reconstituteTodos() (map[string]Todo, error) {
	todos := make(map[string]Todo)
	data, err := os.ReadFile("todos.json")

	if err != nil {
		return nil, errors.New("failed to read todos.json")
	}

	if len(data) != 0 {
		if err := json.Unmarshal(data, &todos); err != nil {
			return nil, errors.New("failed to unmarshal json data")
		}
	}

	return todos, nil
}

func main() {
	todos, _ := reconstituteTodos()

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createTodo(todos, w, r)
		case http.MethodGet:
			getTodos(todos, w, r)
		case http.MethodDelete:
			deleteTodos(todos, w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Not Found"))
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
