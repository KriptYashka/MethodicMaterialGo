# Занятие 8: Легковесные потоки (горутины) + Профилирование

## Тема 15: Легковесные потоки (горутины)

### Горутины

Горутина — легковесный поток выполнения, запускаемый с ключевого слова `go`.

```go
go func() {
    fmt.Println("Hello from goroutine")
}()
```

Горутины:
- Занимают ~4 КБ стека (против 1-8 МБ у OS-потоков)
- Стоимость создания ~1 мкс
- Планируются Go runtime (M:N scheduling), а не OS
- Могут быть миллионы в одном процессе

### Каналы

Канал — типизированный "трубопровод" для передачи данных между горутинами.

```go
ch := make(chan int)   // небуферизированный
ch := make(chan int, 10) // буферизированный

ch <- 42  // отправка
val := <-ch // получение
close(ch)   // закрытие
```

#### Небуферизированные каналы

Отправка блокируется до получения, получение блокируется до отправки — синхронизация "рукопожатие".

```go
func worker(ch chan string) {
    msg := <-ch
    fmt.Println("received:", msg)
}

ch := make(chan string)
go worker(ch)
ch <- "hello" // блокируется, пока worker не получит
```

#### Буферизированные каналы

Отправка блокируется только когда буфер заполнен, получение — когда пуст.

```go
ch := make(chan int, 3)
ch <- 1
ch <- 2
ch <- 3
// ch <- 4 // блокировка — буфер полон
```

### Паттерн: генератор

```go
func fib(n int) <-chan int {
    ch := make(chan int)
    go func() {
        a, b := 0, 1
        for i := 0; i < n; i++ {
            ch <- a
            a, b = b, a+b
        }
        close(ch)
    }()
    return ch
}

for x := range fib(10) {
    fmt.Println(x)
}
```

### Выбор (select)

`select` ждёт один из нескольких каналов:

```go
select {
case msg := <-ch1:
    fmt.Println("from ch1:", msg)
case msg := <-ch2:
    fmt.Println("from ch2:", msg)
case <-time.After(1 * time.Second):
    fmt.Println("timeout")
default:
    fmt.Println("no one ready")
}
```

### Паттерн: timeout

```go
ch := doWork()
select {
case result := <-ch:
    fmt.Println(result)
case <-time.After(3 * time.Second):
    fmt.Println("timeout")
}
```

### Паттерн: fan-out (воркеры)

```go
jobs := make(chan int, 100)
results := make(chan int, 100)

// 3 воркера
for w := 0; w < 3; w++ {
    go func(id int) {
        for job := range jobs {
            results <- job * 2
        }
    }(w)
}

// отправляем работу
for j := 0; j < 10; j++ {
    jobs <- j
}
close(jobs)

// собираем результаты
for r := 0; r < 10; r++ {
    <-results
}
```

### Паттерн: fan-in (мультиплексирование)

```go
func merge(cs ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    for _, c := range cs {
        wg.Add(1)
        go func(ch <-chan int) {
            defer wg.Done()
            for v := range ch {
                out <- v
            }
        }(c)
    }
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}
```

### waitgroup

```go
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(n int) {
        defer wg.Done()
        fmt.Println(n)
    }(i)
}
wg.Wait()
```

### Mutex

```go
type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *Counter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}
```

### Once

```go
var once sync.Once
for i := 0; i < 10; i++ {
    go func() {
        once.Do(func() { fmt.Println("only once") })
    }()
}
```

### Atomic

```go
var counter atomic.Int64
counter.Add(1)
val := counter.Load()
```

### GOMAXPROCS

```go
runtime.GOMAXPROCS(4) // ограничить число OS-потоков
```

По умолчанию — число ядер CPU.

### Data Race

Гонка данных — одновременный доступ к памяти без синхронизации.

```bash
go run -race main.go
```

Обнаружение: `go run -race`, `go build -race`, `go test -race`.

---

## Тема 16: Профилирование

### pprof

Go включает встроенный профайлер `runtime/pprof` и `net/http/pprof`.

#### HTTP-профилирование

```go
import _ "net/http/pprof"

// В main()
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

Доступные профили:
- `/debug/pprof/` — список
- `/debug/pprof/profile` — CPU (30 сек)
- `/debug/pprof/heap` — память
- `/debug/pprof/goroutine` — стек горутин
- `/debug/pprof/block` — блокировки
- `/debug/pprof/mutex` — contention

#### Сбор профиля

```bash
# CPU
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Память
go tool pprof http://localhost:6060/debug/pprof/heap

# Горутины
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

#### Анализ в pprof

```bash
top       # top по времени/памяти
list func # показать строки функции
web       # граф в браузере
peek      # кто вызывает
traces    # стек-трейсы
```

В вебе (через `-http`):

```bash
go tool pprof -http :8081 ~/pprof/pprof.samples.cpu.001.pb.gz
```

### Бенчмарки

```go
// file_test.go
func BenchmarkSum(b *testing.B) {
    nums := make([]int, 1000)
    for i := range nums {
        nums[i] = i
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        sum(nums)
    }
}
```

```bash
go test -bench=. -benchmem
```

### trace

```go
import "runtime/trace"

f, _ := os.Create("trace.out")
trace.Start(f)
defer trace.Stop()
```

```bash
go tool trace trace.out
```

### GODEBUG и GOTRACEBACK

```bash
GODEBUG=gctrace=1 ./app      # логи GC
GOTRACEBACK=all ./app         # стек всех горутин при панике
GOTRACEBUG=schedtrace=1000 ./app  # лог шедулера каждые 1000 мс
```
