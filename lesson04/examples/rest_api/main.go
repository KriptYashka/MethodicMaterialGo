package main

import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"
    "sync"
    "time"
)

type Task struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Done      bool      `json:"done"`
    CreatedAt time.Time `json:"created_at"`
}

type TaskStore struct {
    mu     sync.RWMutex
    tasks  map[int]Task
    nextID int
}

func NewTaskStore() *TaskStore {
    return &TaskStore{
        tasks:  make(map[int]Task),
        nextID: 1,
    }
}

func (s *TaskStore) List() []Task {
    s.mu.RLock()
    defer s.mu.RUnlock()
    result := make([]Task, 0, len(s.tasks))
    for _, t := range s.tasks {
        result = append(result, t)
    }
    return result
}

func (s *TaskStore) Get(id int) (Task, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    t, ok := s.tasks[id]
    return t, ok
}

func (s *TaskStore) Create(title string) Task {
    s.mu.Lock()
    defer s.mu.Unlock()
    t := Task{
        ID:        s.nextID,
        Title:     title,
        Done:      false,
        CreatedAt: time.Now(),
    }
    s.tasks[s.nextID] = t
    s.nextID++
    return t
}

func (s *TaskStore) Update(id int, title string, done bool) (Task, bool) {
    s.mu.Lock()
    defer s.mu.Unlock()
    t, ok := s.tasks[id]
    if !ok {
        return Task{}, false
    }
    if title != "" {
        t.Title = title
    }
    t.Done = done
    s.tasks[id] = t
    return t, true
}

func (s *TaskStore) Delete(id int) bool {
    s.mu.Lock()
    defer s.mu.Unlock()
    _, ok := s.tasks[id]
    if ok {
        delete(s.tasks, id)
    }
    return ok
}

// ---- HTTP handlers ----

type TaskHandler struct {
    store *TaskStore
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
    tasks := h.store.List()
    writeJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Title string `json:"title"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid JSON")
        return
    }
    if req.Title == "" {
        writeError(w, http.StatusBadRequest, "title is required")
        return
    }
    task := h.store.Create(req.Title)
    writeJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid id")
        return
    }
    task, ok := h.store.Get(id)
    if !ok {
        writeError(w, http.StatusNotFound, "task not found")
        return
    }
    writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid id")
        return
    }
    var req struct {
        Title string `json:"title"`
        Done  bool   `json:"done"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid JSON")
        return
    }
    task, ok := h.store.Update(id, req.Title, req.Done)
    if !ok {
        writeError(w, http.StatusNotFound, "task not found")
        return
    }
    writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid id")
        return
    }
    if !h.store.Delete(id) {
        writeError(w, http.StatusNotFound, "task not found")
        return
    }
    writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// ---- helpers ----

func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
    writeJSON(w, status, map[string]string{"error": msg})
}

func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func main() {
    store := NewTaskStore()
    handler := &TaskHandler{store: store}

    mux := http.NewServeMux()
    mux.HandleFunc("GET /api/tasks", handler.List)
    mux.HandleFunc("POST /api/tasks", handler.Create)
    mux.HandleFunc("GET /api/tasks/{id}", handler.Get)
    mux.HandleFunc("PUT /api/tasks/{id}", handler.Update)
    mux.HandleFunc("DELETE /api/tasks/{id}", handler.Delete)

    wrapped := corsMiddleware(mux)

    log.Println("REST API on :8080")
    log.Fatal(http.ListenAndServe(":8080", wrapped))
}

// Test:
//   curl http://localhost:8080/api/tasks
//   curl -X POST -d '{"title":"Learn Go"}' http://localhost:8080/api/tasks
//   curl -X PUT -d '{"title":"Learn Go","done":true}' http://localhost:8080/api/tasks/1
//   curl -X DELETE http://localhost:8080/api/tasks/1
