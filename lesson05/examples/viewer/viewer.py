"""
Fruit Game Field Viewer (Pygame + asyncio)
Поле слева. Справа: Teams слева, Market справа, Legend внизу.
"""
import asyncio
import sys

import httpx
import pygame

from config import (
    SERVER_HOST, SERVER_PORT,
    CELL_COLORS, CELL_SYMBOLS, CELL_NAMES,
    FRUIT_COLORS, FRUIT_NAMES,
    TEAM_COLORS,
    WINDOW_WIDTH, WINDOW_HEIGHT, CELL_SIZE,
)

BASE_URL = f"http://{SERVER_HOST}:{SERVER_PORT}"

BLACK = (0, 0, 0)
WHITE = (255, 255, 255)
GRAY = (200, 200, 200)
DARK_GRAY = (50, 50, 50)
GREEN = (100, 255, 100)
RED = (255, 100, 100)

LAYOUT_LEFT = 20
LAYOUT_TOP = 55
FIELD_COLS = 10
FIELD_ROWS = 10
FONT_NAME = "Impact"


async def fetch_json(client, path):
    try:
        resp = await client.get(f"{BASE_URL}{path}", timeout=5)
        return resp.json()
    except Exception as e:
        return {"error": str(e)}


def draw_field(screen, font, cells, ox, oy):
    for idx, cell in enumerate(cells):
        x = ox + (idx % FIELD_COLS) * CELL_SIZE
        y = oy + (idx // FIELD_COLS) * CELL_SIZE
        rect = pygame.Rect(x, y, CELL_SIZE, CELL_SIZE)

        ctype = cell.get("type", 1)
        color = CELL_COLORS.get(ctype, GRAY)
        pygame.draw.rect(screen, color, rect)
        pygame.draw.rect(screen, BLACK, rect, 1)

        team_id = cell.get("team_id", 0)
        fruit = cell.get("fruit")

        if fruit:
            ft = fruit.get("type", 0)
            ripe = fruit.get("ripe", False)
            fcolor = FRUIT_COLORS.get(ft, WHITE)
            cx = x + CELL_SIZE // 2
            cy = y + CELL_SIZE // 2
            r = CELL_SIZE // 4
            pygame.draw.circle(screen, fcolor, (cx, cy), r)
            if ripe:
                pygame.draw.circle(screen, WHITE, (cx, cy), r, 2)
        elif team_id > 0:
            tcolor = TEAM_COLORS.get(team_id, WHITE)
            text = font.render(str(team_id), True, tcolor)
            screen.blit(text, (x + CELL_SIZE // 2 - text.get_width() // 2,
                               y + CELL_SIZE // 2 - text.get_height() // 2))
        else:
            sym = CELL_SYMBOLS.get(ctype, "?")
            text = font.render(sym, True, DARK_GRAY)
            screen.blit(text, (x + CELL_SIZE // 2 - text.get_width() // 2,
                               y + CELL_SIZE // 2 - text.get_height() // 2))

        if fruit and not fruit.get("ripe", False):
            growth = fruit.get("growth", 0)
            bar_w = CELL_SIZE - 4
            bar_h = 4
            bar_x = x + 2
            bar_y = y + CELL_SIZE - bar_h - 2
            pygame.draw.rect(screen, DARK_GRAY, (bar_x, bar_y, bar_w, bar_h))
            pygame.draw.rect(screen, (0, 255, 0), (bar_x, bar_y, int(bar_w * growth), bar_h))


def draw_teams(screen, font, teams, x, y, max_w):
    text = font.render("TEAMS", True, WHITE)
    screen.blit(text, (x, y))
    y += 30
    for team in teams:
        tcolor = TEAM_COLORS.get(team["id"], WHITE)
        line = f"#{team['id']} {team['name']}"
        text = font.render(line, True, tcolor)
        screen.blit(text, (x, y)); y += 22
        money = font.render(f"${team['money']}", True, WHITE)
        screen.blit(money, (x + 10, y)); y += 22
        inv = team.get("inventory", {})
        parts = []
        for k, v in inv.items():
            if v > 0:
                fname = FRUIT_NAMES.get(int(k), k)
                parts.append(f"{fname}:{v}")
        if parts:
            inv_text = font.render("  ".join(parts), True, GRAY)
            screen.blit(inv_text, (x + 10, y))
        y += 28


def draw_market(screen, font, market, x, y, max_w):
    text = font.render("MARKET", True, WHITE)
    screen.blit(text, (x, y))
    y += 30

    prices = market.get("prices", {})
    indicators = market.get("indicators", {})

    for raw_key, price in prices.items():
        try:
            ft = int(raw_key)
        except (ValueError, TypeError):
            ft = 0
        fname = FRUIT_NAMES.get(ft, raw_key)
        ind = indicators.get(raw_key, {})
        rsi = ind.get("rsi", 50)
        hist = ind.get("histogram", 0)

        dir_sym = chr(8594)
        if hist > 0.5:
            dir_sym = chr(8593)
        elif hist < -0.5:
            dir_sym = chr(8595)

        rsi_color = GREEN if rsi < 40 else RED if rsi > 60 else WHITE
        line = font.render(f"{fname} ${price:.2f}{dir_sym}", True, WHITE)
        screen.blit(line, (x, y)); y += 22
        rsi_text = font.render(f"RSI{rsi:.1f}", True, rsi_color)
        screen.blit(rsi_text, (x + 10, y))
        y += 28


def draw_legend(screen, font, x, y):
    text = font.render("LEGEND", True, WHITE)
    screen.blit(text, (x, y))
    y += 24
    items = [
        ("~", CELL_NAMES.get(0), CELL_COLORS.get(0)),
        ("#", CELL_NAMES.get(1), CELL_COLORS.get(1)),
        ("^", CELL_NAMES.get(2), CELL_COLORS.get(2)),
        ("*", CELL_NAMES.get(3), CELL_COLORS.get(3)),
        ("H", CELL_NAMES.get(4), CELL_COLORS.get(4)),
    ]
    for sym, name, color in items:
        line = font.render(f"  {sym}={name}", True, color)
        screen.blit(line, (x, y)); y += 18
    y += 6
    for line in ["W/M/R=fruit", "O-ring=ripe", "bar=growth", "num=owner"]:
        t = font.render(line, True, GRAY)
        screen.blit(t, (x, y)); y += 16


async def main():
    pygame.init()
    screen = pygame.display.set_mode((WINDOW_WIDTH, WINDOW_HEIGHT))
    pygame.display.set_caption("Fruit Game - Field Viewer")
    title_font = pygame.font.SysFont(FONT_NAME, 26)
    font = pygame.font.SysFont(FONT_NAME, 16)
    small = pygame.font.SysFont(FONT_NAME, 14)

    running = True
    field_data = {"cells": [], "width": 10, "height": 10}
    teams_data = []
    market_data = {"prices": {}, "indicators": {}}

    field_w = FIELD_COLS * CELL_SIZE
    sep_x = LAYOUT_LEFT + field_w + 20
    panel_x = sep_x + 14
    panel_w = WINDOW_WIDTH - panel_x - 14
    col_w = panel_w // 2

    async with httpx.AsyncClient() as client:
        while running:
            for event in pygame.event.get():
                if event.type == pygame.QUIT:
                    running = False

            field_data, teams_data, market_data = await asyncio.gather(
                fetch_json(client, "/api/field"),
                fetch_json(client, "/api/teams"),
                fetch_json(client, "/api/market"),
            )

            if isinstance(field_data, dict) and "error" not in field_data:
                cells = field_data.get("cells", [])
            else:
                cells = []
            if not isinstance(teams_data, list):
                teams_data = []
            if not isinstance(market_data, dict) or "error" in market_data:
                market_data = {"prices": {}, "indicators": {}}

            screen.fill(BLACK)

            tick = field_data.get("tick", 0) if isinstance(field_data, dict) else 0
            top = title_font.render(f"Fruit Game  Tick {tick}", True, WHITE)
            screen.blit(top, (LAYOUT_LEFT, 12))

            pygame.draw.line(screen, DARK_GRAY, (sep_x, 0), (sep_x, WINDOW_HEIGHT), 2)

            draw_field(screen, font, cells, LAYOUT_LEFT, LAYOUT_TOP)

            # right side: Teams | Market
            py = LAYOUT_TOP
            draw_teams(screen, font, teams_data, panel_x, py, col_w)
            draw_market(screen, font, market_data, panel_x + col_w + 8, py, col_w)

            # legend below both
            ly = py + len(teams_data) * 75 + 20
            draw_legend(screen, small, panel_x, ly)

            status = small.render(f"Connected to {BASE_URL}", True, GREEN)
            screen.blit(status, (LAYOUT_LEFT, WINDOW_HEIGHT - 24))

            pygame.display.flip()
            await asyncio.sleep(0.5)

    pygame.quit()
    sys.exit()


if __name__ == "__main__":
    asyncio.run(main())
