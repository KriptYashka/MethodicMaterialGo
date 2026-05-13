# Занятие 6: SQL и SQLite + Работа с БД из Go

## Тема 11: Введение в SQL и SQLite

### Что такое SQL?

SQL (Structured Query Language) — язык структурированных запросов для работы с реляционными базами данных.

Основные операции (CRUD):
- `CREATE` — создание таблиц
- `INSERT` — вставка данных
- `SELECT` — чтение
- `UPDATE` — обновление
- `DELETE` — удаление

### SQLite

SQLite — встраиваемая реляционная БД без отдельного сервера. Хранит всю БД в одном файле. Идеальна для обучения, прототипирования и микросервисов.

```sql
-- Создание таблицы
CREATE TABLE IF NOT EXISTS products (
    id    INTEGER PRIMARY KEY AUTOINCREMENT,
    name  TEXT    NOT NULL,
    price REAL    NOT NULL DEFAULT 0,
    qty   INTEGER NOT NULL DEFAULT 0
);

-- Вставка
INSERT INTO products (name, price, qty) VALUES ('Widget', 9.99, 100);

-- Выборка
SELECT id, name, price FROM products WHERE price > 5 ORDER BY price DESC;

-- Обновление
UPDATE products SET qty = qty - 1 WHERE id = 1;

-- Удаление
DELETE FROM products WHERE qty = 0;
```

### Типы данных в SQLite

| Тип    | Описание                  |
|--------|---------------------------|
| INTEGER | Целое число (int64)      |
| REAL    | Число с плавающей точкой |
| TEXT    | Строка (UTF-8)           |
| BLOB    | Бинарные данные          |
| NULL    | Отсутствие значения      |

### Агрегатные функции

```sql
SELECT COUNT(*) FROM products;
SELECT AVG(price) FROM products;
SELECT SUM(qty) FROM products;
SELECT category, MAX(price) FROM products GROUP BY category;
```

### JOIN

```sql
SELECT o.id, p.name, o.qty
FROM orders o
JOIN products p ON o.product_id = p.id
WHERE o.user_id = 1;
```

---

## Тема 12: Работа с БД из Go

### Пакет database/sql

Go предоставляет интерфейс `database/sql` для работы с SQL-БД. Драйвер подключается отдельно.

```go
import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3" // драйвер SQLite
)
```

### Подключение

```go
db, err := sql.Open("sqlite3", "./app.db")
if err != nil {
    log.Fatal(err)
}
defer db.Close()

// Проверка соединения
if err := db.Ping(); err != nil {
    log.Fatal(err)
}
```

> `sql.Open` не создаёт подключение — только инициализирует пул. Реальное соединение открывается при первом запросе.

### Exec — выполнение без возврата строк

```go
_, err := db.Exec(`CREATE TABLE IF NOT EXISTS products (
    id    INTEGER PRIMARY KEY AUTOINCREMENT,
    name  TEXT    NOT NULL,
    price REAL    NOT NULL
)`)
```

```go
result, err := db.Exec(`INSERT INTO products (name, price) VALUES (?, ?)`, "Widget", 9.99)
id, _ := result.LastInsertId()
n, _  := result.RowsAffected()
```

> Используйте `?` placeholders — это защищает от SQL-инъекций.

### QueryRow — одна строка

```go
var name string
var price float64
err := db.QueryRow(`SELECT name, price FROM products WHERE id = ?`, 1).Scan(&name, &price)
if err == sql.ErrNoRows {
    log.Println("not found")
}
```

### Query — много строк

```go
rows, err := db.Query(`SELECT id, name, price FROM products WHERE price > ?`, minPrice)
if err != nil {
    log.Fatal(err)
}
defer rows.Close()

for rows.Next() {
    var p Product
    if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
        log.Fatal(err)
    }
    fmt.Println(p)
}
if err := rows.Err(); err != nil {
    log.Fatal(err)
}
```

### Транзакции

```go
tx, err := db.Begin()
if err != nil {
    return err
}
defer tx.Rollback() // откат, если не вызван Commit

_, err = tx.Exec(`UPDATE products SET qty = qty - 1 WHERE id = ?`, productID)
if err != nil {
    return err
}
_, err = tx.Exec(`INSERT INTO orders (user_id, product_id) VALUES (?, ?)`, userID, productID)
if err != nil {
    return err
}

return tx.Commit()
```

### Подготовленные запросы (Prepared Statements)

```go
stmt, err := db.Prepare(`INSERT INTO products (name, price) VALUES (?, ?)`)
if err != nil {
    log.Fatal(err)
}
defer stmt.Close()

stmt.Exec("A", 10)
stmt.Exec("B", 20) // переиспользуем план запроса
```

### Конфигурация пула соединений

```go
db.SetMaxOpenConns(25)       // макс. открытых соединений
db.SetMaxIdleConns(5)        // макс. в простое
db.SetConnMaxLifetime(5 * time.Minute) // время жизни соединения
```

### RawBytes и NULL

```go
var price sql.NullFloat64
err := db.QueryRow(`SELECT price FROM products WHERE id = ?`, id).Scan(&price)
if price.Valid {
    fmt.Println(price.Float64)
} else {
    fmt.Println("NULL")
}
```

### Миграции (ручные)

Простая схема — хранить SQL-файлы в папке `migrations/`:

```
migrations/
  001_create_products.sql
  002_add_category.sql
```

Применять по очереди, отслеживая версию в таблице `schema_migrations`.

```go
// Упрощённый пример применения миграций
migrations := []string{
    `CREATE TABLE IF NOT EXISTS products (... )`,
    `ALTER TABLE products ADD COLUMN category TEXT`,
}
for i, m := range migrations {
    if _, err := db.Exec(m); err != nil {
        log.Fatalf("migration %d failed: %v", i+1, err)
    }
}
```

### Готовый пример структуры

```go
type Product struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
    Qty   int     `json:"qty"`
}

type Store struct {
    db *sql.DB
}

func NewStore(db *sql.DB) *Store {
    return &Store{db: db}
}

func (s *Store) List(ctx context.Context) ([]Product, error) {
    rows, err := s.db.QueryContext(ctx, `SELECT id, name, price, qty FROM products`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []Product
    for rows.Next() {
        var p Product
        if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Qty); err != nil {
            return nil, err
        }
        products = append(products, p)
    }
    return products, rows.Err()
}
```

### Context в запросах

Все методы `database/sql` поддерживают `context.Context`:
- `ExecContext`
- `QueryContext`
- `QueryRowContext`
- `PrepareContext`

Это позволяет контролировать таймауты:

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

var name string
err := db.QueryRowContext(ctx, `SELECT name FROM products WHERE id = ?`, id).Scan(&name)
```
