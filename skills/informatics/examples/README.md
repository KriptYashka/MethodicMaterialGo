# Examples / Примеры кода

Рабочие примеры Go и Python для каждого занятия.

## Основы Go

### lesson01
| Пример | Описание | Запуск |
|--------|----------|--------|
| [hello](../../lesson01/examples/hello/) | Первая программа | `go run .` |
| [variables](../../lesson01/examples/variables/) | Переменные и типы | `go run .` |
| [conditions](../../lesson01/examples/conditions/) | if/else/switch | `go run .` |
| [loops](../../lesson01/examples/loops/) | Циклы for | `go run .` |
| [slices](../../lesson01/examples/slices/) | Слайсы, append, copy | `go run .` |

### lesson02
| Пример | Описание | Запуск |
|--------|----------|--------|
| [functions](../../lesson02/examples/functions/) | Функции, defer, замыкания | `go run .` |
| [maps](../../lesson02/examples/maps/) | Карты, comma-ok | `go run .` |
| [errors](../../lesson02/examples/errors/) | Errors.Is, As, wrapping | `go run .` |
| [panic_recover](../../lesson02/examples/panic_recover/) | Panic/recover | `go run .` |

## Web

### lesson03
| Пример | Описание | Запуск |
|--------|----------|--------|
| [simple_server](../../lesson03/examples/simple_server/) | HTTP-сервер | `go run .` |
| [handlers](../../lesson03/examples/handlers/) | Кастомные Handler | `go run .` |
| [json_api](../../lesson03/examples/json_api/) | JSON API | `go run .` |
| [middleware](../../lesson03/examples/middleware/) | Middleware | `go run .` |

### lesson04
| Пример | Описание | Запуск |
|--------|----------|--------|
| [mux_router](../../lesson04/examples/mux_router/) | ServeMux, PathValue | `go run .` |
| [rest_api](../../lesson04/examples/rest_api/) | REST-эндпоинты | `go run .` |
| [frontend](../../lesson04/examples/frontend/) | Go + HTML/CSS/JS | `go run .` → http://localhost:8080 |

### lesson05
| Пример | Описание | Запуск |
|--------|----------|--------|
| [server](../../lesson05/examples/server/) | Go-сервер игры | `go run .` |
| [client](../../lesson05/examples/client/) | CLI-клиент | `go run .` |
| [viewer](../../lesson05/examples/viewer/) | Pygame поле (Python) | `python viewer.py` |
| [exchange](../../lesson05/examples/exchange/) | Pygame биржа (Python) | `python exchange.py` |

## Базы данных

### lesson06
| Пример | Описание | Запуск |
|--------|----------|--------|
| [main.go](../../lesson06/examples/main.go) | SQLite Store: CRUD, транзакции | `go run .` |

### lesson07
| Пример | Описание | Запуск |
|--------|----------|--------|
| [main.go](../../lesson07/examples/main.go) | REST + SQLite + DI (Handler/Service/Repo) | `go run .` |
| [service_test.go](../../lesson07/examples/service_test.go) | Unit-тесты с mock | `go test .` |

## Конкурентность

### lesson08
| Пример | Описание | Запуск |
|--------|----------|--------|
| [01_goroutines](../../lesson08/examples/01_goroutines/) | Базовые горутины | `go run .` |
| [02_channels](../../lesson08/examples/02_channels/) | Каналы | `go run .` |
| [03_select](../../lesson08/examples/03_select/) | Select | `go run .` |
| [04_mutex](../../lesson08/examples/04_mutex/) | Mutex | `go run .` |
| [05_pprof](../../lesson08/examples/05_pprof/) | Профилирование | `go run .` |
| [06_workerpool](../../lesson08/examples/06_workerpool/) | Worker pool | `go run .` |
| [07_counter](../../lesson08/examples/07_counter/) | Атомарный счётчик | `go run .` |
