# Fruit Game — конфигурация клиента

SERVER_HOST = "127.0.0.1"
SERVER_PORT = 8080

# CellType enum (должен совпадать с сервером)
CELL_WATER = 0
CELL_EARTH = 1
CELL_MOUNTAIN = 2
CELL_LEGENDARY = 3
CELL_HOUSE = 4

# Цвета клеток (ключи — CellType int)
CELL_COLORS = {
    CELL_WATER: (64, 128, 255),
    CELL_EARTH: (139, 90, 43),
    CELL_MOUNTAIN: (128, 128, 128),
    CELL_LEGENDARY: (255, 215, 0),
    CELL_HOUSE: (34, 139, 34),
}

# Символы клеток
CELL_SYMBOLS = {
    CELL_WATER: "~",
    CELL_EARTH: "#",
    CELL_MOUNTAIN: "^",
    CELL_LEGENDARY: "*",
    CELL_HOUSE: "H",
}

# Названия клеток
CELL_NAMES = {
    CELL_WATER: "Water",
    CELL_EARTH: "Earth",
    CELL_MOUNTAIN: "Mountain",
    CELL_LEGENDARY: "Legendary",
    CELL_HOUSE: "House",
}

# FruitType enum
FRUIT_WATERMELON = 0
FRUIT_MELON = 1
FRUIT_RASPBERRY = 2

# Цвета фруктов
FRUIT_COLORS = {
    FRUIT_WATERMELON: (50, 180, 50),
    FRUIT_MELON: (255, 200, 50),
    FRUIT_RASPBERRY: (200, 50, 50),
}

# Названия фруктов
FRUIT_NAMES = {
    FRUIT_WATERMELON: "watermelon",
    FRUIT_MELON: "melon",
    FRUIT_RASPBERRY: "raspberry",
}

# Цвета команд
TEAM_COLORS = {
    1: (0, 120, 255),
    2: (255, 80, 80),
    3: (255, 200, 0),
}

# Размеры окна
WINDOW_WIDTH = 1024
WINDOW_HEIGHT = 768
CELL_SIZE = 48
