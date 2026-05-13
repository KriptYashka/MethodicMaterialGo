# Занятие 7: Интерфейсы и инъекции + Практикум с БД

## Тема 13: Интерфейсы и инъекции зависимостей

### Интерфейсы в Go

Интерфейс — набор методов. В Go интерфейсы реализуются неявно (duck typing).

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type ReadWriter interface {
    Reader
    Writer
}
```

#### Неявная реализация

```go
type FileStore struct{ path string }

func (f *FileStore) Read(p []byte) (int, error) {
    return os.Open(f.path).Read(p)
}

func (f *FileStore) Write(p []byte) (int, error) {
    return os.WriteFile(f.path, p, 0644)
}

// *FileStore реализует Reader и Writer автоматически
var r Reader = &FileStore{path: "data.txt"}
```

#### Зачем нужны интерфейсы

1. **Абстракция** — код работает с контрактом, а не с конкретной реализацией
2. **Тестируемость** — легко подменить реальную БД на mock
3. **Гибкость** — замена реализации без изменения вызывающего кода

#### Пустой интерфейс

```go
var x interface{} // любой тип
x = 42
x = "hello"
x = struct{}{}

val, ok := x.(int) // type assertion
switch v := x.(type) {
case int:    fmt.Println("int:", v)
case string: fmt.Println("string:", v)
default:     fmt.Println("unknown")
}
```

### Dependency Injection (DI)

DI — передача зависимостей через конструктор, а не создание внутри.

#### Проблема (без DI)

```go
type Handler struct{}

func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
    db, _ := sql.Open("sqlite3", "db.sqlite") // ❌ жёсткая связь
    rows, _ := db.Query(`SELECT ...`)
    // ...
}
```

#### Решение (с DI)

```go
type ProductRepository interface {
    List(ctx context.Context) ([]Product, error)
    GetByID(ctx context.Context, id int) (*Product, error)
}

type Handler struct {
    repo ProductRepository // ✅ интерфейс
}

func NewHandler(repo ProductRepository) *Handler {
    return &Handler{repo: repo}
}
```

При старте приложения создаём конкретные реализации и внедряем их (composition root):

```go
func main() {
    db, _ := sql.Open("sqlite3", "db.sqlite")
    repo := NewSQLiteProductRepo(db)   // конкретная реализация
    svc  := NewProductService(repo)    // бизнес-логика
    h    := NewHandler(svc)            // HTTP-слой
    http.HandleFunc("/products", h.GetProducts)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Принцип слоёв (Layered Architecture)

```
handlers  (HTTP, JSON)
    ↓
services  (бизнес-логика)
    ↓
repository (БД, внешние API)
```

Каждый слой общается с нижележащим через интерфейсы.

### Распространённые ошибки

- **Слишком большие интерфейсы** — лучше несколько маленьких (принцип ISP)
- **Интерфейсы на стороне производителя** — определяйте интерфейсы там, где они используются (на стороне потребителя)
- **Преждевременная абстракция** — не выделяйте интерфейс, пока нет второй реализации

---

## Тема 14: Практикум с БД

### Полный пример: HTTP + SQLite + DI

#### Repository layer

```go
type Product struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
    Qty   int     `json:"qty"`
}

type ProductRepository interface {
    List(ctx context.Context) ([]Product, error)
    GetByID(ctx context.Context, id int) (*Product, error)
    Create(ctx context.Context, p *Product) (int64, error)
    Update(ctx context.Context, p *Product) error
    Delete(ctx context.Context, id int) error
}
```

#### Service layer

```go
type ProductService struct {
    repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
    return &ProductService{repo: repo}
}

func (s *ProductService) ListProducts(ctx context.Context) ([]Product, error) {
    products, err := s.repo.List(ctx)
    if err != nil {
        return nil, fmt.Errorf("list products: %w", err)
    }
    if products == nil {
        products = []Product{} // вместо null
    }
    return products, nil
}
```

#### Handler layer

```go
type ProductHandler struct {
    svc *ProductService
}

func NewProductHandler(svc *ProductService) *ProductHandler {
    return &ProductHandler{svc: svc}
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
    products, err := h.svc.ListProducts(r.Context())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(products)
}
```

#### Тестирование с моками

```go
type mockRepo struct {
    products []Product
    err      error
}

func (m *mockRepo) List(ctx context.Context) ([]Product, error) {
    return m.products, m.err
}
// ... остальные методы

func TestListProducts(t *testing.T) {
    repo := &mockRepo{
        products: []Product{{ID: 1, Name: "Test", Price: 10}},
    }
    svc := NewProductService(repo)
    products, err := svc.ListProducts(context.Background())
    if err != nil {
        t.Fatal(err)
    }
    if len(products) != 1 {
        t.Fatalf("expected 1, got %d", len(products))
    }
}
```

### Миграции (продвинутый подход)

Использование `golang-migrate/migrate`:

```go
import (
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/sqlite3"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

m, err := migrate.New("file://migrations", "sqlite3://db.sqlite")
if err != nil {
    log.Fatal(err)
}
if err := m.Up(); err != nil && err != migrate.ErrNoChange {
    log.Fatal(err)
}
```

Структура папки `migrations/`:

```
migrations/
  000001_create_products.up.sql
  000001_create_products.down.sql
  000002_add_orders.up.sql
  000002_add_orders.down.sql
```
