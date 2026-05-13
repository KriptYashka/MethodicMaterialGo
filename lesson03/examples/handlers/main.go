package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "strings"
)

// ---- different handler forms ----

// 1. Handler interface
type healthHandler struct{}

func (h *healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// 2. HandlerFunc
func echoHandler(w http.ResponseWriter, r *http.Request) {
    body, _ := io.ReadAll(r.Body)
    defer r.Body.Close()

    w.Header().Set("Content-Type", "text/plain")
    fmt.Fprintf(w, "Method: %s\nPath: %s\nBody: %s", r.Method, r.URL.Path, string(body))
}

// 3. Closure returning HandlerFunc
func methodHandler(allowed string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != allowed {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
            return
        }
        fmt.Fprintf(w, "Handled with %s", allowed)
    }
}

// 4. Handler using http.HandlerFunc type conversion
func headersHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    for k, v := range r.Header {
        fmt.Fprintf(w, "%s: %s\n", k, strings.Join(v, ", "))
    }
}

func main() {
    mux := http.NewServeMux()

    mux.Handle("/health", &healthHandler{})
    mux.HandleFunc("/echo", echoHandler)
    mux.HandleFunc("/get-only", methodHandler(http.MethodGet))
    mux.HandleFunc("/headers", headersHandler)

    log.Println("Server on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
