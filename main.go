package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var usersDB []UserResponse

type UserRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
type UserResponse struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Status string `json:"status"`
}

func handleEverything(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%s] %s\n", r.Method, r.URL.Path)
	switch r.URL.Path {
	case "/ping":
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			response := PingResponse{Status: "ok"}
			json.NewEncoder(w).Encode(response)
		} else {
			http.Error(w, "Метод не поддерживается!", http.StatusMethodNotAllowed)
		}
	case "/echo":
		if r.Method == http.MethodPost {
			var req UserRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if req.Name == "" {
				http.Error(w, "Имя не должно быть пустым", http.StatusBadRequest)
				return
			}
			if req.Age < 18 {
				http.Error(w, "Доступ только для совершеннолетних", http.StatusBadRequest)
				return
			}
			response := UserResponse{req.Name, req.Age, "approved"}
			usersDB = append(usersDB, response)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			http.Error(w, "Метод не поддерживается!", http.StatusMethodNotAllowed)
		}
	case "/users":
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(usersDB)
		}
	default:
		w.Write([]byte("Я такого не знаю?!"))
	}
}

type PingResponse struct {
	Status string `json:"status"`
}

func main() {
	http.HandleFunc("/", handleEverything)

	http.ListenAndServe(":8080", nil)

}
