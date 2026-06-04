# Skill: Разработка высоконагруженных сервисов на Go

Курс по Go для разработки высоконагруженных сервисов. Длительность: 2 недели (20 занятий).

## Структура

| Неделя | Занятия | Темы |
|--------|---------|------|
| 1 | lesson01–lesson05 | Основы Go, API, HTTP-сервер, REST, фронтенд, JSON, практикум |
| 2 | lesson06–lesson10 | SQL/SQLite, БД из Go, интерфейсы, DI, горутины, профилирование, итоговый проект, деплой |

## Reference

Теоретический материал по каждому занятию — файлы `theory.md`:

| Занятие | Темы | Файл |
|---------|------|------|
| lesson01 | Введение, переменные, типы, условия, циклы, срезы | [theory.md](lesson01/theory.md) |
| lesson02 | Карты, функции, обработка ошибок, panic/recover | [theory.md](lesson02/theory.md) |
| lesson03 | API (REST, методы, статусы), HTTP-сервер на net/http | [theory.md](lesson03/theory.md) |
| lesson04 | Маршрутизация (ServeMux), REST API, html/template, фронтенд | [theory.md](lesson04/theory.md) |
| lesson05 | Структуры, JSON, маршалинг, практикум по Web | [theory.md](lesson05/theory.md) |
| lesson06 | SQL, SQLite, database/sql, транзакции, подготовленные запросы | [theory.md](lesson06/theory.md) |
| lesson07 | Интерфейсы, DI, слоистая архитектура, тестирование с моками | [theory.md](lesson07/theory.md) |
| lesson08 | Горутины, каналы, select, sync, профилирование (pprof/trace) | [theory.md](lesson08/theory.md) |

## Examples

Рабочие примеры кода по каждому занятию:

| Занятие | Пример | Описание |
|---------|--------|----------|
| lesson01 | `examples/hello/` | Hello, World! |
| | `examples/variables/` | Объявление переменных, типы |
| | `examples/conditions/` | if/else, switch |
| | `examples/loops/` | Циклы for |
| | `examples/slices/` | Слайсы, append, copy |
| lesson02 | `examples/functions/` | Функции, замыкания, defer |
| | `examples/maps/` | Карты (map) |
| | `examples/errors/` | Обработка ошибок |
| | `examples/panic_recover/` | Panic и recover |
| lesson03 | `examples/simple_server/` | Простейший HTTP-сервер |
| | `examples/handlers/` | Обработчики запросов |
| | `examples/json_api/` | JSON API (decoder/encoder) |
| | `examples/middleware/` | Middleware-паттерны |
| lesson04 | `examples/mux_router/` | Маршрутизация через ServeMux |
| | `examples/rest_api/` | REST API структура |
| | `examples/frontend/` | HTML+JS+CSS фронтенд с Go-бэкендом |
| lesson05 | `examples/client/` | CLI-клиент для игры |
| | `examples/server/` | Go-сервер игры (game logic) |
| | `examples/viewer/` | Pygame-вьювер поля (Python) |
| | `examples/exchange/` | Pygame-вьювер биржи (Python) |
| lesson06 | `examples/main.go` | Полноценный Store: CRUD, транзакции, context |
| lesson07 | `examples/main.go` | REST + SQLite + DI: Handler -> Service -> Repository |
| | `examples/service_test.go` | Unit-тесты с mock-репозиторием |
| lesson08 | `examples/01_goroutines/` | Базовый запуск горутин |
| | `examples/02_channels/` | Каналы |
| | `examples/03_select/` | Select |
| | `examples/04_mutex/` | Mutex |
| | `examples/05_pprof/` | Профилирование pprof |
| | `examples/06_workerpool/` | Worker pool |
| | `examples/07_counter/` | Атомарный счётчик |
