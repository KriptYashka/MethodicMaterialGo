package main

import (
    "embed"
    "encoding/json"
    "log"
    "net/http"
    "strconv"
    "sync"
    "time"
)

//go:embed static
var staticFiles embed.FS

type Task struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Done      bool      `json:"done"`
    CreatedAt time.Time `json:"created_at"`
}

type Store struct {
    mu     sync.RWMutex
    tasks  map[int]Task
    nextID int
}

func NewStore() *Store {
    return &Store{tasks: make(map[int]Task), nextID: 1}
}

func (s *Store) List() []Task {
    s.mu.RLock()
    defer s.mu.RUnlock()
    out := make([]Task, 0, len(s.tasks))
    for _, t := range s.tasks {
        out = append(out, t)
    }
    return out
}

func (s *Store) Create(title string) Task {
    s.mu.Lock()
    defer s.mu.Unlock()
    t := Task{ID: s.nextID, Title: title, CreatedAt: time.Now()}
    s.tasks[s.nextID] = t
    s.nextID++
    return t
}

func (s *Store) Update(id int, title string, done bool) (Task, bool) {
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

func (s *Store) Delete(id int) bool {
    s.mu.Lock()
    defer s.mu.Unlock()
    _, ok := s.tasks[id]
    if ok {
        delete(s.tasks, id)
    }
    return ok
}

func main() {
    store := NewStore()

    mux := http.NewServeMux()

    mux.HandleFunc("GET /api/tasks", func(w http.ResponseWriter, r *http.Request) {
        writeJSON(w, http.StatusOK, store.List())
    })

    mux.HandleFunc("POST /api/tasks", func(w http.ResponseWriter, r *http.Request) {
        var req struct{ Title string }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            writeError(w, http.StatusBadRequest, "invalid JSON")
            return
        }
        task := store.Create(req.Title)
        writeJSON(w, http.StatusCreated, task)
    })

    mux.HandleFunc("PUT /api/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
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
        task, ok := store.Update(id, req.Title, req.Done)
        if !ok {
            writeError(w, http.StatusNotFound, "not found")
            return
        }
        writeJSON(w, http.StatusOK, task)
    })

    mux.HandleFunc("DELETE /api/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
        id, err := strconv.Atoi(r.PathValue("id"))
        if err != nil {
            writeError(w, http.StatusBadRequest, "invalid id")
            return
        }
        if !store.Delete(id) {
            writeError(w, http.StatusNotFound, "not found")
            return
        }
        writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
    })

    // Serve static frontend
    fs := http.FileServer(http.FS(staticFiles))
    mux.Handle("GET /", fs)

    log.Println("Full app on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}

func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
    writeJSON(w, status, map[string]string{"error": msg})
}
