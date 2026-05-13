# Занятие 2: Функции, карты и обработка ошибок

## Тема 3: Карты и функции

### Функции

#### Объявление функций

```go
// Обычная функция
func add(a int, b int) int {
    return a + b
}

// Тип параметров можно указать один раз
func add(a, b int) int {
    return a + b
}
```

#### Возврат нескольких значений

```go
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

quotient, err := divide(10, 2)
```

#### Именованные возвращаемые значения

```go
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return // naked return
}
```

#### Variadic-функции (вариативные)

```go
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

sum(1, 2, 3)
sum(1, 2, 3, 4, 5)

// Распаковка среза
nums := []int{1, 2, 3}
sum(nums...)
```

#### Функции как значения

```go
fn := func(a, b int) int {
    return a + b
}
result := fn(3, 4) // 7
```

#### Замыкания (closures)

```go
func counter() func() int {
    i := 0
    return func() int {
        i++
        return i
    }
}

c := counter()
fmt.Println(c()) // 1
fmt.Println(c()) // 2
```

#### defer

Выполняется при выходе из функции (LIFO).

```go
func readFile(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close() // выполнится в конце
    // работа с файлом...
    return nil
}
```

### Карты (Maps)

Хеш-таблицы. Ссылочный тип.

#### Создание

```go
// Через make
ages := make(map[string]int)

// Литерал
ages := map[string]int{
    "Alice": 30,
    "Bob":   25,
}
```

#### Операции

```go
m := make(map[string]int)

// Добавление / обновление
m["key"] = 42

// Получение
val := m["key"]

// Comma-ok идиома
val, ok := m["key"]
if ok {
    fmt.Println("exists:", val)
}

// Удаление
delete(m, "key")

// Длина
len(m)

// Итерация (порядок не гарантируется!)
for k, v := range m {
    fmt.Printf("%s -> %d\n", k, v)
}
```

#### Map как множество

```go
set := make(map[string]bool)
set["apple"] = true

if set["apple"] {
    fmt.Println("apple is in set")
}
```

#### Важно

- Map — ссылочный тип. При присваивании копируется ссылка.
- Zero value для map — nil. Запись в nil-map вызывает panic.
- Map нельзя сравнивать через `==` (только с nil).
- Порядок итерации по map случайный.

---

## Тема 4: Обработка ошибок

### Интерфейс error

В Go ошибки — это значения, реализующие интерфейс:

```go
type error interface {
    Error() string
}
```

#### Создание ошибок

```go
// Простая ошибка
errors.New("something went wrong")

// С форматированием
fmt.Errorf("user %d not found", id)
```

### Sentinel-ошибки

```go
var ErrNotFound = errors.New("not found")

func GetUser(id int) (*User, error) {
    if id < 1 {
        return nil, ErrNotFound
    }
    return &User{ID: id}, nil
}

// Проверка
if errors.Is(err, ErrNotFound) {
    fmt.Println("user not found, creating new one")
}
```

### Ошибки с контекстом (wrapping)

```go
if err := doStep(); err != nil {
    return fmt.Errorf("do step failed: %w", err)
}

// Распаковка
if errors.Is(err, ErrNotFound) {
    // match
}

// Проверка типа
var netErr *net.DNSError
if errors.As(err, &netErr) {
    fmt.Println("DNS error:", netErr)
}
```

### Собственные типы ошибок

```go
type ValidationError struct {
    Field string
    Value any
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed on %s", e.Field)
}

func validate(v any) error {
    if v == nil {
        return &ValidationError{Field: "value", Value: v}
    }
    return nil
}
```

### panic и recover

Panic — аварийное завершение. Используется редко.

```go
func safeDiv(a, b int) (result int, ok bool) {
    defer func() {
        if r := recover(); r != nil {
            result = 0
            ok = false
        }
    }()
    return a / b, true
}
```

**Правила:**
- `panic` — только для truly exceptional ситуаций (не для обычных ошибок)
- `recover` — только внутри `defer`
- Используйте `errors.Is`/`errors.As` вместо проверки строк

### Best Practices

1. Всегда проверяйте ошибки
2. Не игнорируйте ошибки (`_` — только если осознанно)
3. Добавляйте контекст через `%w`
4. Используйте sentinel-ошибки для ожидаемых сценариев
5. Собственные типы — для дополнительных данных
6. Panic — только при инициализации программы
