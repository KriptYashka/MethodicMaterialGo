package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
)

type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}

// in-memory store
var users = []User{
    {ID: 1, Name: "Alice", Email: "alice@example.com", CreatedAt: time.Now()},
    {ID: 2, Name: "Bob", Email: "bob@example.com", CreatedAt: time.Now()},
}
var nextID = 3

func writeJSON(w http.ResponseWriter, status int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
    writeJSON(w, status, ErrorResponse{Error: msg})
}

func listUsers(w http.ResponseWriter, r *http.Request) {
    writeJSON(w, http.StatusOK, users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid JSON body")
        return
    }
    if req.Name == "" || req.Email == "" {
        writeError(w, http.StatusBadRequest, "name and email are required")
        return
    }

    user := User{
        ID:        nextID,
        Name:      req.Name,
        Email:     req.Email,
        CreatedAt: time.Now(),
    }
    nextID++
    users = append(users, user)

    writeJSON(w, http.StatusCreated, user)
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            listUsers(w, r)
        case http.MethodPost:
            createUser(w, r)
        default:
            writeError(w, http.StatusMethodNotAllowed, "method not allowed")
        }
    })

    addr := ":8080"
    fmt.Printf("JSON API server on %s\n", addr)
    log.Fatal(http.ListenAndServe(addr, mux))
}

// Test:
//   curl http://localhost:8080/users
//   curl -X POST -d '{"name":"Charlie","email":"charlie@example.com"}' http://localhost:8080/users
