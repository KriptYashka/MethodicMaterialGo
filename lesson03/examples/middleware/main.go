package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

// ---- middleware ----

func logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
    })
}

func recovery(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("panic recovered: %v", err)
                http.Error(w, "internal server error", http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}

func withHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Server", "GoMiddlewareDemo")
        w.Header().Set("X-Duration-Options", "nope")
        next.ServeHTTP(w, r)
    })
}

// ---- handlers ----

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello from middleware demo!")
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
    panic("something went terribly wrong")
}

// ---- chaining ----

func chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
    for _, m := range middlewares {
        h = m(h)
    }
    return h
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", helloHandler)
    mux.HandleFunc("/panic", panicHandler)

    // apply middlewares
    wrapped := chain(mux, logging, recovery, withHeaders)

    log.Println("Server on :8080")
    log.Fatal(http.ListenAndServe(":8080", wrapped))
}
