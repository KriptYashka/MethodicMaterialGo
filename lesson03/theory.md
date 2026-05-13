# Занятие 3: API и HTTP-сервер на net/http

## Тема 5: API

### Что такое API?

API (Application Programming Interface) — контракт взаимодействия между компонентами системы.

#### Стили API

- **REST** — ресурсно-ориентированный, HTTP-методы, JSON
- **GraphQL** — единый endpoint, клиент выбирает поля
- **gRPC** — Protocol Buffers, HTTP/2, бинарный протокол
- **SOAP** — XML, строгие контракты (устаревает)

Фокус курса — REST API.

### RESTful API

**Принципы:**
- Ресурс — существительное (`/users`, `/orders`)
- HTTP-методы = CRUD-операции
- Stateless — каждый запрос содержит всю информацию
- Единый интерфейс (URI, методы, статусы)

#### HTTP-методы

| Метод    | CRUD   | Idempotent | Безопасный |
|----------|--------|------------|------------|
| GET      | Read   | Да         | Да         |
| POST     | Create | Нет        | Нет        |
| PUT      | Update (full) | Да  | Нет        |
| PATCH    | Update (partial) | Нет | Нет   |
| DELETE   | Delete | Да         | Нет        |

#### Коды ответа

| Код  | Описание |
|------|----------|
| 200  | OK             |
| 201  | Created        |
| 204  | No Content     |
| 400  | Bad Request    |
| 401  | Unauthorized   |
| 403  | Forbidden      |
| 404  | Not Found      |
| 409  | Conflict       |
| 422  | Unprocessable  |
| 500  | Internal Error |

#### Формат запроса/ответа

```json
// Request: POST /users
{
    "name": "Alice",
    "email": "alice@example.com"
}

// Response: 201 Created
{
    "id": 1,
    "name": "Alice",
    "email": "alice@example.com",
    "created_at": "2025-01-15T10:00:00Z"
}
```

---

## Тема 6: HTTP-сервер на net/http

### Структура HTTP-сервера

```
Client  ──►  Request  ──►  Handler  ──►  Response  ──►  Client
```

#### Основные типы

```go
type Handler interface {
    ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    f(w, r)
}
```

### Создание сервера

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, %s!", r.URL.Query().Get("name"))
    })

    http.ListenAndServe(":8080", nil)
}
```

### ResponseWriter

```go
// Запись тела ответа
w.Write([]byte("OK"))

// Установка статус-кода
w.WriteHeader(http.StatusNotFound)

// Установка заголовков
w.Header().Set("Content-Type", "application/json")

// Важно: WriteHeader должен быть после Header().Set() и до Write()
```

### Request

```go
// Query parameters
r.URL.Query().Get("name")
r.URL.Query()["tags"] // multiple values

// Path
r.URL.Path

// Method
r.Method

// Headers
r.Header.Get("Authorization")

// Body
body, _ := io.ReadAll(r.Body)
defer r.Body.Close()
```

### Чтение и запись JSON

```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

// Decode request
var u User
json.NewDecoder(r.Body).Decode(&u)

// Encode response
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(u)
```

### ServeMux — маршрутизатор по умолчанию

```go
mux := http.NewServeMux()
mux.HandleFunc("/users", listUsers)
mux.HandleFunc("/users/create", createUser)

http.ListenAndServe(":8080", mux)
```

### Middleware (обёртки)

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

mux := http.NewServeMux()
mux.HandleFunc("/", handler)

loggedMux := loggingMiddleware(mux)
http.ListenAndServe(":8080", loggedMux)
```

### Структура production-сервера

```go
srv := &http.Server{
    Addr:         ":8080",
    Handler:      mux,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  30 * time.Second,
}

if err := srv.ListenAndServe(); err != nil {
    log.Fatal(err)
}
```

### Graceful Shutdown

```go
srv := &http.Server{Addr: ":8080", Handler: mux}

go func() {
    srv.ListenAndServe()
}()

// wait for signal
sig := make(chan os.Signal, 1)
signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
<-sig

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
srv.Shutdown(ctx)
```
