# Занятие 1: Введение в Go и основы языка

## Тема 1: Введение в Go. Переменные и типы. Условия

### Что такое Go?

Go (Golang) — компилируемый, статически типизированный язык программирования, созданный в Google (2009). Авторы: Роберт Гриземер, Роб Пайк, Кен Томпсон.

**Ключевые особенности:**
- Простота и читаемость
- Быстрая компиляция
- Встроенная поддержка конкурентности (горутины)
- Сборка мусора
- Богатая стандартная библиотека
- Кроссплатформенность

### Установка и настройка

1. Скачать Go с [golang.org](https://golang.org/dl/)
2. Установить, проверить версию:
```bash
go version
```

### Первая программа

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

Запуск:
```bash
go run main.go
```

Компиляция:
```bash
go build -o myapp main.go
./myapp
```

### Основы Go-модулей

```bash
go mod init myproject
```

Файл `go.mod` содержит имя модуля и зависимости.

### Переменные

Go — статическая типизация. Тип указывается после имени или выводится автоматически.

```go
// Полное объявление
var name string = "Alice"

// С выводом типа
var age = 30

// Краткое объявление (только внутри функций)
city := "Moscow"

// Множественное объявление
var x, y int = 1, 2
var a, b = "hello", 42

// Групповое объявление
var (
    width  int    = 100
    height int    = 200
    title  string = "Image"
)

// Zero values (значения по умолчанию)
var i int       // 0
var f float64   // 0
var s string    // ""
var b bool      // false
```

### Базовые типы

```go
// Числовые
var a int8   = 127          // -128..127
var b int16  = 32767
var c int32  = 2147483647
var d int64  = 9223372036854775807
var e int    = 42           // платформозависимый (32 или 64 бита)

var u uint   = 100          // беззнаковый
var f float32 = 3.14
var g float64 = 2.71828

// Строки (immutable, UTF-8)
var s string = "Привет, Go!"

// Булевы
var ok bool = true

// Комплексные числа
var c complex128 = 1 + 2i
```

### Константы

```go
const Pi = 3.14159
const (
    StatusOK = 200
    StatusNotFound = 404
)

// iota — генератор последовательных значений
const (
    Red = iota // 0
    Green      // 1
    Blue       // 2
)
```

### Преобразование типов

В Go нет неявного приведения — только явное:

```go
var i int = 42
var f float64 = float64(i)
var s string = strconv.Itoa(i)
n, _ := strconv.Atoi("123")
```

### Условные операторы

#### if

```go
if x > 0 {
    fmt.Println("positive")
} else if x < 0 {
    fmt.Println("negative")
} else {
    fmt.Println("zero")
}
```

**Короткая запись с инициализацией:**

```go
if err := doSomething(); err != nil {
    fmt.Println("Error:", err)
}
```

#### switch

```go
// По значению
switch day {
case 1:
    fmt.Println("Monday")
case 2:
    fmt.Println("Tuesday")
default:
    fmt.Println("Unknown")
}

// Без выражения (как if-else)
switch {
case score >= 90:
    grade = "A"
case score >= 80:
    grade = "B"
default:
    grade = "F"
}
```

В Go в `case` не нужен `break` — выход происходит автоматически. Чтобы провалиться дальше, используйте `fallthrough`.

---

## Тема 2: Циклы и срезы

### Циклы

В Go есть только одна конструкция цикла — `for`. Нет `while` или `do-while`.

#### Классический for

```go
for i := 0; i < 10; i++ {
    fmt.Println(i)
}
```

#### for как while

```go
n := 1
for n < 100 {
    n *= 2
}
```

#### Бесконечный цикл

```go
for {
    // работает, пока не break
}
```

#### break и continue

```go
for i := 0; i < 10; i++ {
    if i%2 == 0 {
        continue
    }
    if i > 7 {
        break
    }
    fmt.Println(i)
}
```

### Массивы

Массив — фиксированного размера, передаётся по значению.

```go
var arr [5]int          // [0 0 0 0 0]
arr[0] = 1
arr2 := [3]int{1, 2, 3}
arr3 := [...]int{4, 5, 6} // размер определяется компилятором
```

### Срезы (Slices)

Срез — динамическая структура данных, основной инструмент для работы с последовательностями.

```go
// Создание
var s []int
s = append(s, 1, 2, 3)

// Литерал
s := []int{1, 2, 3}

// make — с указанием длины и ёмкости
s := make([]int, 5)      // len=5, cap=5
s := make([]int, 3, 10)  // len=3, cap=10

// Срез из массива
arr := [5]int{1, 2, 3, 4, 5}
slice := arr[1:3] // [2, 3]
```

#### Внутреннее устройство среза

Срез — это структура из трёх полей:
- **ptr** — указатель на первый элемент в массиве
- **len** — длина
- **cap** — ёмкость (до конца массива)

```
┌──────────┐
│ ptr ───────► [1] [2] [3] [4] [5]
│ len = 3  │
│ cap = 5  │
└──────────┘
```

#### Операции со срезами

```go
s := []int{1, 2, 3}

// append
s = append(s, 4, 5)     // [1 2 3 4 5]
s = append(s, 6, 7, 8)  // [1 2 3 4 5 6 7 8]

// copy
dst := make([]int, len(s))
copy(dst, s)

// Срез среза
sub := s[1:3]  // [2 3]
sub := s[:2]   // [1 2]
sub := s[2:]   // [3 4 5 6 7 8]

// Изменение ёмкости (full slice expression)
sub := s[1:3:4] // len=2, cap=3
```

#### Важно про append

Если len < cap — append записывает в существующий массив.
Если len == cap — создаётся новый массив в 2 раза больше.

```go
s := make([]int, 2, 4) // [0 0], cap=4
s2 := append(s, 1)     // разделяют общий массив
s[0] = 99              // s2[0] тоже изменится!
```

### range

```go
nums := []int{10, 20, 30}

for i, v := range nums {
    fmt.Printf("index=%d value=%d\n", i, v)
}

// Пропуск индекса
for _, v := range nums {
    fmt.Println(v)
}

// Только индекс
for i := range nums {
    fmt.Println(i)
}
```

### Строки как байтовые срезы

```go
s := "hello"
for i, b := range []byte(s) {
    fmt.Printf("%d: %d\n", i, b)
}

// Итерация по рунам (unicode)
for i, r := range "Привет" {
    fmt.Printf("%d: %c\n", i, r)
}
```
