# Занятие 5: Структуры и JSON + Практикум по Web

## Тема 9: Структуры + JSON

### Структуры в Go

Структура — составной тип, объединяющий поля разных типов.

```go
type User struct {
    ID        int
    Name      string
    Email     string
    CreatedAt time.Time
}
```

#### Создание и инициализация

```go
// Позиционная (не рекомендуется — хрупко)
u := User{1, "Alice", "alice@example.com", time.Now()}

// Именованная (рекомендуется)
u := User{
    ID:    1,
    Name:  "Alice",
    Email: "alice@example.com",
}

// Zero-значение
var u User
u.ID = 1
u.Name = "Alice"
```

#### Вложенные структуры

```go
type Address struct {
    City    string
    Street  string
    ZipCode string
}

type User struct {
    ID       int
    Name     string
    Address  Address  // вложенная структура
    Tags     []string // срез в структуре
}
```

#### Теги структур (struct tags)

```go
type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email,omitempty"`
    Password  string    `json:"-"`
    CreatedAt time.Time `json:"created_at"`
}
```

### JSON: маршалинг и анмаршалинг

```go
import "encoding/json"

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

// Marshal — структура → JSON (байты)
u := User{ID: 1, Name: "Alice"}
data, err := json.Marshal(u)
// data = {"id":1,"name":"Alice"}

// MarshalIndent — с отступами
data, _ := json.MarshalIndent(u, "", "  ")

// Unmarshal — JSON → структура
jsonStr := `{"id":1,"name":"Alice"}`
var u User
err := json.Unmarshal([]byte(jsonStr), &u)
```

#### Reader/Writer потоки

```go
// Decode из запроса
var u User
json.NewDecoder(r.Body).Decode(&u)

// Encode в ответ
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(u)
```

#### Кастомная маршализация

```go
type Color struct {
    R, G, B uint8
}

func (c Color) MarshalJSON() ([]byte, error) {
    return json.Marshal(fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B))
}

func (c *Color) UnmarshalJSON(data []byte) error {
    var s string
    if err := json.Unmarshal(data, &s); err != nil {
        return err
    }
    fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
    return nil
}
```

### Сериализация в Go

```go
// JSON — основная
// Gob — бинарный, Go-to-Go
// XML — для совместимости
// YAML — через сторонние библиотеки (gopkg.in/yaml.v3)
// Protocol Buffers — бинарный, через protobuf
```

---

## Тема 10: Практикум по Web (соревнование на 3 команды)

### Архитектура игры

```
┌─────────────────────────┐     HTTP/JSON     ┌──────────────────────┐
│  Pygame Field Viewer    │ ◄──────────────► │   Go Game Server     │
│  (графическое поле +    │                   │  (in-memory storage) │
│   команды + цены)       │                   │  + tick loop         │
└─────────────────────────┘                   │  + MACD/RSI market   │
                                              │  + quest system      │
┌─────────────────────────┐                   └──────────┬───────────┘
│  Pygame Exchange Viewer │ ◄──────────────►             │
│  (свечи + MACD + RSI)   │               ┌──────────────┴───────────┐
└─────────────────────────┘               │  Go CLI Client          │
                                           │  (команды для тестов)   │
                                           └──────────────────────────┘
```

### API эндпоинты игры

| Method | Path | Описание |
|--------|------|----------|
| GET | /api/field | Состояние поля |
| GET | /api/teams | Информация о командах |
| POST | /api/plant | Посадить фрукт |
| POST | /api/harvest/{cell_idx} | Собрать урожай |
| POST | /api/buy-cell | Купить ячейку |
| POST | /api/buy-tool | Купить инструмент |
| POST | /api/sell | Продать фрукты |
| GET | /api/market | Цены и свечи |
| GET | /api/quest/{id} | Получить вопрос |
| POST | /api/quest/{id} | Ответить на вопрос |
| GET | /api/house/{cell_idx} | Взаимодействие с домом |

### Game Design Summary

- **Поле**: 10×10, типы клеток: Земля, Вода, Горы, Легендарная, Дом
- **Фрукты**: Арбуз, Дыня, Малина — у каждой свои аффинити к типам клеток
- **Команды**: 2–3, стартовый капитал 100 монет
- **Тик**: каждую секунду, фрукты растут
- **Рост**: зависит от аффинити + можно ускорить удобрением
- **Сбор**: 30 секунд после созревания, иначе фрукт пропадает
- **Инструменты**: удобрение (×2 рост), забор (ячейка остаётся на 2 цикла)
- **Квесты**: GET вопрос → POST ответ → +деньги
- **Дом**: нельзя сажать, но можно получить случайную фразу

### Рынок и технические индикаторы

**Границы цен.** У каждого фрукта есть интервал:
- Арбуз: $5–$40
- Дыня: $3–$35
- Малина: $6–$45

**Дисперсия.** Определяет волатильность фрукта (малина самая волатильная — 0.7).

**MACD (Moving Average Convergence Divergence):**
- EMA(12) — быстрая скользящая средняя
- EMA(26) — медленная скользящая средняя
- MACD Line = EMA(12) − EMA(26)
- Signal Line = EMA(9) от MACD Line
- Гистограмма = MACD Line − Signal Line
- Положительная гистограмма → бычий сигнал (цена вероятнее пойдёт вверх)

**RSI (Relative Strength Index, 14 периодов):**
- RSI = 100 − 100/(1 + RS), где RS = средний прирост / средняя потеря
- RSI > 70 — перекупленность (oversold) → цена вероятнее пойдёт вниз
- RSI < 30 — перепроданность (overbought) → цена вероятнее пойдёт вверх

**Механизм pullback.** Если цена выходит за границы, следующие 2–3 шага цена принудительно возвращается в диапазон.

### Теги структур для JSON

Обратите внимание на используемые теги в коде сервера:

```go
type Fruit struct {
    Type      FruitType `json:"type"`
    TeamID    int       `json:"team_id"`
    PlantedAt time.Time `json:"planted_at"`
    Growth    float64   `json:"growth"`
    Ripe      bool      `json:"ripe"`
}
```

### Запуск игры

```bash
# 1. Сервер
cd lesson05/examples/server
go run .

# 2. Pygame Field Viewer (в другом терминале)
cd lesson05/examples/viewer
pip install -r requirements.txt
python viewer.py

# 3. Pygame Exchange Viewer (в третьем терминале)
cd lesson05/examples/exchange
pip install -r requirements.txt
python exchange.py

# 4. Go CLI-клиент (для тестовых запросов)
cd lesson05/examples/client
go run . field
go run . teams
go run . buy-cell 1 5
go run . plant 1 5 0
```
