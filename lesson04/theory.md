# Занятие 4: Маршрутизация, REST API и основы фронтенда

## Тема 7: Маршрутизация и REST API

### Маршрутизация в Go (Go 1.22+)

Начиная с Go 1.22, `http.ServeMux` поддерживает паттерны с методами и path-параметрами.

```go
mux := http.NewServeMux()

// Метод + путь
mux.HandleFunc("GET /users", listUsers)
mux.HandleFunc("POST /users", createUser)

// Path-параметры
mux.HandleFunc("GET /users/{id}", getUser)

// Строгий подпаттерн
mux.HandleFunc("GET /users/{id}/posts", getUserPosts)

// До Go 1.22 — ручной разбор
mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/users/")
})
```

### Извлечение path-параметров (Go 1.22+)

```go
func getUser(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    // ... 
}
```

### Проектирование REST API

#### Правила именования ресурсов

| Ресурс | GET | POST | PUT / {id} | DELETE / {id} |
|--------|-----|------|------------|---------------|
| /users | список | создать | обновить | удалить |
| /users/{id}/posts | посты пользователя | создать пост | — | — |

#### Структура REST-обработчика

```go
type UserHandler struct {
    store *UserStore
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        h.list(w, r)
    case http.MethodPost:
        h.create(w, r)
    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
}
```

---

## Тема 8: Основы фронтенда

### HTML-шаблоны в Go

Go имеет встроенный пакет `html/template` — безопасное экранирование HTML.

#### Синтаксис шаблонов

```html
<!-- Вывод значения -->
{{ .Name }}

<!-- Условие -->
{{ if .LoggedIn }}
    <p>Welcome, {{ .Name }}!</p>
{{ else }}
    <a href="/login">Login</a>
{{ end }}

<!-- Цикл -->
<ul>
{{ range .Items }}
    <li>{{ . }}</li>
{{ end }}
</ul>
```

#### Использование в Go

```go
tmpl := template.Must(template.ParseFiles("templates/index.html"))

data := struct {
    Name string
    Items []string
}{"Alice", []string{"A", "B", "C"}}

tmpl.Execute(w, data)
```

### Раздача статических файлов

```go
fs := http.FileServer(http.Dir("static"))
mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
```

### Embed — встраивание файлов в бинарник (Go 1.16+)

```go
//go:embed templates/*
var templateFS embed.FS

tmpl := template.Must(template.ParseFS(templateFS, "templates/*.html"))
```

### Пример: полноценный REST API + фронтенд

```
project/
├── main.go
├── static/
│   ├── index.html
│   ├── style.css
│   └── app.js
└── templates/
    └── page.html
```

#### Итоговая архитектура

```
Browser  ──►  /static/ (html, css, js) ──►  Client-side rendering
                │
                ▼  fetch /api/users
            JSON API (Go)  ──►  Database
```

### Структура ответа API

```json
// Успех
{
    "data": { ... },
    "meta": { "total": 42 }
}

// Ошибка
{
    "error": "invalid email",
    "code": "VALIDATION_ERROR"
}
```
