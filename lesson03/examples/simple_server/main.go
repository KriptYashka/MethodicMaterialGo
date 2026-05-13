package main

import (
    "fmt"
    "log"
    "net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    if name == "" {
        name = "World"
    }
    fmt.Fprintf(w, "Hello, %s!", name)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"status":"ok"}`))
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/hello", helloHandler)
    mux.HandleFunc("/health", healthHandler)

    srv := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    log.Println("Server starting on :8080")
    if err := srv.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}
