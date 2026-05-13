"""
Fruit Exchange Viewer (Pygame + asyncio)
"""
import asyncio
import sys

import httpx
import pygame

from config import (
    SERVER_HOST, SERVER_PORT,
    FRUIT_WATERMELON, FRUIT_MELON, FRUIT_RASPBERRY,
    FRUIT_NAMES, FRUIT_COLORS,
    WINDOW_WIDTH, WINDOW_HEIGHT,
    BG_COLOR, GREEN, RED, WHITE, GRAY, YELLOW, ORANGE,
    PRICE_CHART_HEIGHT, MACD_CHART_HEIGHT, RSI_CHART_HEIGHT,
    CHART_MARGIN, CHART_LEFT, CANDLE_WIDTH, CANDLE_GAP,
    MAX_CANDLES,
)

BASE_URL = f"http://{SERVER_HOST}:{SERVER_PORT}"
FRUIT_KEYS = [FRUIT_WATERMELON, FRUIT_MELON, FRUIT_RASPBERRY]
FONT_NAME = "Impact"

BTN_WIDTH = 180
BTN_HEIGHT = 40
BTN_GAP = 10


async def fetch_json(client, path):
    try:
        resp = await client.get(f"{BASE_URL}{path}", timeout=5)
        return resp.json()
    except Exception as e:
        return {"error": str(e), "prices": {}, "candles": {}, "indicators": {}}


def compute_visible_range(data_len, chart_width):
    total_w = CANDLE_WIDTH + CANDLE_GAP
    max_fit = chart_width // total_w
    if data_len <= max_fit:
        return 0, data_len
    return data_len - max_fit, data_len


def draw_buttons(screen, font, selected_key):
    total_w = len(FRUIT_KEYS) * BTN_WIDTH + (len(FRUIT_KEYS) - 1) * BTN_GAP
    start_x = (WINDOW_WIDTH - total_w) // 2

    buttons = []
    for i, fk in enumerate(FRUIT_KEYS):
        x = start_x + i * (BTN_WIDTH + BTN_GAP)
        y = 10
        rect = pygame.Rect(x, y, BTN_WIDTH, BTN_HEIGHT)
        is_active = (fk == selected_key)
        color = FRUIT_COLORS.get(fk, WHITE)
        bg = color if is_active else (40, 40, 60)
        border = color

        pygame.draw.rect(screen, bg, rect)
        pygame.draw.rect(screen, border, rect, 3 if is_active else 1)

        label = FRUIT_NAMES.get(fk, str(fk)).upper()
        text_color = (0, 0, 0) if is_active else color
        txt = font.render(label, True, text_color)
        screen.blit(txt, (rect.centerx - txt.get_width() // 2,
                          rect.centery - txt.get_height() // 2))
        buttons.append((rect, fk))
    return buttons


def draw_candle_chart(screen, font, fruit_key, candles, indicators, rect):
    x0, y0, w, h = rect
    pygame.draw.rect(screen, (30, 30, 45), rect)
    pygame.draw.rect(screen, GRAY, rect, 1)

    fname = FRUIT_NAMES.get(fruit_key, str(fruit_key))
    color = FRUIT_COLORS.get(fruit_key, WHITE)
    price = candles[-1]["close"] if candles else 0
    ind = indicators.get(str(fruit_key), {})
    rsi = ind.get("rsi", 50)
    macd_line = ind.get("macd_line", 0)

    title_text = f"{fname.upper()}  ${price:.2f}  RSI:{rsi:.1f}  MACD:{macd_line:.2f}"
    tsurf = font.render(title_text, True, color)
    screen.blit(tsurf, (x0 + 5, y0 + 3))

    if len(candles) < 2:
        no_text = font.render("Not enough data yet...", True, GRAY)
        screen.blit(no_text, (x0 + 50, y0 + h // 2))
        return

    chart_area_y = y0 + 25
    chart_area_h = h - 35
    chart_area_x = x0 + CHART_LEFT
    chart_area_w = w - CHART_LEFT - 20

    low = min(c["low"] for c in candles) * 0.98
    high = max(c["high"] for c in candles) * 1.02
    if high - low < 0.01:
        high = low + 1

    start, end = compute_visible_range(len(candles), chart_area_w)
    visible = candles[start:end]
    step = chart_area_w / max(len(visible), 1)

    grid_lines = 5
    for i in range(grid_lines + 1):
        gy = chart_area_y + int(chart_area_h * i / grid_lines)
        pygame.draw.line(screen, (50, 50, 70), (chart_area_x, gy),
                         (chart_area_x + chart_area_w, gy), 1)
        price_val = high - (high - low) * i / grid_lines
        plabel = font.render(f"${price_val:.1f}", True, GRAY)
        screen.blit(plabel, (chart_area_x - 55, gy - 8))

    for i, c in enumerate(visible):
        cx = int(chart_area_x + i * step)
        ccolor = GREEN if c["close"] >= c["open"] else RED

        oy = chart_area_y + int(chart_area_h * (high - c["open"]) / (high - low))
        cy = chart_area_y + int(chart_area_h * (high - c["close"]) / (high - low))
        ly = chart_area_y + int(chart_area_h * (high - c["low"]) / (high - low))
        hy = chart_area_y + int(chart_area_h * (high - c["high"]) / (high - low))

        pygame.draw.line(screen, ccolor, (cx, ly), (cx, hy), 1)
        body_top = min(oy, cy)
        body_h = max(abs(cy - oy), 2)
        body_w = max(int(step * 0.7), 2)
        pygame.draw.rect(screen, ccolor, (cx - body_w // 2, body_top, body_w, body_h))

    if visible:
        last = visible[-1]
        last_y = chart_area_y + int(chart_area_h * (high - last["close"]) / (high - low))
        plab = font.render(f"${last['close']:.2f}", True, WHITE)
        screen.blit(plab, (chart_area_x + chart_area_w - plab.get_width() - 5, last_y - 10))


def draw_macd(screen, font, fruit_key, candles, indicators, rect):
    x0, y0, w, h = rect
    pygame.draw.rect(screen, (30, 30, 45), rect)
    pygame.draw.rect(screen, GRAY, rect, 1)

    fname = FRUIT_NAMES.get(fruit_key, str(fruit_key))
    color = FRUIT_COLORS.get(fruit_key, WHITE)
    label = font.render(f"{fname.upper()} MACD", True, color)
    screen.blit(label, (x0 + 5, y0 + 3))

    closes = [c["close"] for c in candles]
    if len(closes) < 26:
        return

    macd_vals = _ema(closes, 12)
    macd_vals2 = _ema(closes, 26)
    macd_line_series = [macd_vals[i] - macd_vals2[i] for i in range(len(macd_vals))]
    sig_series = _ema(macd_line_series, 9)

    chart_area_x = x0 + CHART_LEFT
    chart_area_w = w - CHART_LEFT - 20
    chart_area_y = y0 + 25
    chart_area_h = h - 30

    start, end = compute_visible_range(len(macd_line_series), chart_area_w)
    step = chart_area_w / max(end - start, 1)

    all_vals = macd_line_series[start:end] + sig_series[start:end]
    if not all_vals:
        return
    lo = min(all_vals) - 0.5
    hi = max(all_vals) + 0.5
    if hi - lo < 0.01:
        hi = lo + 1

    zy = chart_area_y + int(chart_area_h * hi / (hi - lo))
    pygame.draw.line(screen, GRAY, (chart_area_x, zy),
                     (chart_area_x + chart_area_w, zy), 1)

    for i in range(start, end):
        cx = int(chart_area_x + (i - start) * step)
        idx = i - start
        if idx < len(macd_line_series) - start - 1:
            hist_val = macd_line_series[i] - sig_series[i]
            bar_h = int(chart_area_h * abs(hist_val) / (hi - lo))
            bar_y = chart_area_y + int(chart_area_h * (
                hi - max(macd_line_series[i], sig_series[i])) / (hi - lo))
            bar_color = GREEN if hist_val > 0 else RED
            bw = max(int(step * 0.8), 2)
            pygame.draw.rect(screen, bar_color, (cx - bw // 2, bar_y, bw, max(bar_h, 1)))

            if idx > 0:
                px = int(chart_area_x + (i - 1 - start) * step)
                py = chart_area_y + int(chart_area_h * (
                    hi - macd_line_series[i - 1]) / (hi - lo))
                nx = cx
                ny = chart_area_y + int(chart_area_h * (
                    hi - macd_line_series[i]) / (hi - lo))
                pygame.draw.line(screen, (100, 200, 255), (px, py), (nx, ny), 2)

            if idx > 0 and i - start < len(sig_series) - 1:
                px = int(chart_area_x + (i - 1 - start) * step)
                py = chart_area_y + int(chart_area_h * (
                    hi - sig_series[i - 1]) / (hi - lo))
                nx = cx
                ny = chart_area_y + int(chart_area_h * (
                    hi - sig_series[i]) / (hi - lo))
                pygame.draw.line(screen, ORANGE, (px, py), (nx, ny), 2)

    last_m = macd_line_series[-1] if macd_line_series else 0
    last_s = sig_series[-1] if sig_series else 0
    last_h = last_m - last_s
    info = font.render(f"MACD:{last_m:.2f}  Signal:{last_s:.2f}  Hist:{last_h:.2f}",
                       True, WHITE)
    screen.blit(info, (x0 + 5, y0 + h - 20))


def _ema(values, period):
    k = 2.0 / (period + 1)
    result = [values[0]]
    for i in range(1, len(values)):
        result.append(values[i] * k + result[-1] * (1 - k))
    return result


def draw_rsi(screen, font, fruit_key, candles, indicators, rect):
    x0, y0, w, h = rect
    pygame.draw.rect(screen, (30, 30, 45), rect)
    pygame.draw.rect(screen, GRAY, rect, 1)

    fname = FRUIT_NAMES.get(fruit_key, str(fruit_key))
    color = FRUIT_COLORS.get(fruit_key, WHITE)
    label = font.render(f"{fname.upper()} RSI", True, color)
    screen.blit(label, (x0 + 5, y0 + 3))

    closes = [c["close"] for c in candles]
    if len(closes) < 15:
        return

    rsi_vals = []
    for i in range(14, len(closes)):
        gains = 0.0
        losses = 0.0
        for j in range(i - 13, i + 1):
            diff = closes[j] - closes[j - 1]
            if diff > 0:
                gains += diff
            else:
                losses -= diff
        avg_gain = gains / 14
        avg_loss = losses / 14
        if avg_loss == 0:
            rsi_vals.append(100.0)
        else:
            rs = avg_gain / avg_loss
            rsi_vals.append(100 - 100 / (1 + rs))

    if not rsi_vals:
        return

    chart_area_x = x0 + CHART_LEFT
    chart_area_w = w - CHART_LEFT - 20
    chart_area_y = y0 + 25
    chart_area_h = h - 30

    start, end = compute_visible_range(len(rsi_vals), chart_area_w)
    step = chart_area_w / max(end - start, 1)

    for level, lcolor in [(70, RED), (30, GREEN), (50, GRAY)]:
        ly = chart_area_y + int(chart_area_h * (100 - level) / 100)
        lw = 1 if level == 50 else 2
        pygame.draw.line(screen, lcolor, (chart_area_x, ly),
                         (chart_area_x + chart_area_w, ly), lw)
        ltxt = font.render(str(level), True, lcolor)
        screen.blit(ltxt, (chart_area_x + chart_area_w + 3, ly - 8))

    for i in range(start, end):
        cx = int(chart_area_x + (i - start) * step)
        idx = i - start
        val = rsi_vals[i]
        vy = chart_area_y + int(chart_area_h * (100 - val) / 100)
        dot_color = GREEN if val < 40 else RED if val > 60 else WHITE
        pygame.draw.circle(screen, dot_color, (cx, vy), 2)
        if idx > 0:
            px = int(chart_area_x + (i - 1 - start) * step)
            py = chart_area_y + int(chart_area_h * (100 - rsi_vals[i - 1]) / 100)
            pygame.draw.line(screen, dot_color, (px, py), (cx, vy), 1)

    cur = font.render(f"RSI: {rsi_vals[-1]:.1f}", True, WHITE)
    screen.blit(cur, (x0 + 5, y0 + h - 20))


async def main():
    pygame.init()
    screen = pygame.display.set_mode((WINDOW_WIDTH, WINDOW_HEIGHT))
    pygame.display.set_caption("Fruit Exchange - Price Charts")
    font = pygame.font.SysFont(FONT_NAME, 16)
    big_font = pygame.font.SysFont(FONT_NAME, 24)
    small_font = pygame.font.SysFont(FONT_NAME, 12)

    selected_fruit = FRUIT_WATERMELON
    running = True

    async with httpx.AsyncClient() as client:
        while running:
            for event in pygame.event.get():
                if event.type == pygame.QUIT:
                    running = False
                elif event.type == pygame.KEYDOWN:
                    if event.key == pygame.K_ESCAPE:
                        running = False
                    elif event.key == pygame.K_1:
                        selected_fruit = FRUIT_WATERMELON
                    elif event.key == pygame.K_2:
                        selected_fruit = FRUIT_MELON
                    elif event.key == pygame.K_3:
                        selected_fruit = FRUIT_RASPBERRY
                elif event.type == pygame.MOUSEBUTTONDOWN:
                    for rect, fk in buttons:
                        if rect.collidepoint(event.pos):
                            selected_fruit = fk
                            break

            market_data = await fetch_json(client, "/api/market")
            prices = market_data.get("prices", {})
            candles_raw = market_data.get("candles", {})
            indicators = market_data.get("indicators", {})

            screen.fill(BG_COLOR)

            buttons = draw_buttons(screen, font, selected_fruit)

            sk = str(selected_fruit)
            price = prices.get(sk, 0)
            fname = FRUIT_NAMES.get(selected_fruit, sk).upper()
            color = FRUIT_COLORS.get(selected_fruit, WHITE)
            price_text = big_font.render(f"{fname}  ${price:.2f}", True, color)
            screen.blit(price_text, ((WINDOW_WIDTH - price_text.get_width()) // 2, 60))

            hint = small_font.render("Click or press 1/2/3 to switch", True, GRAY)
            screen.blit(hint, ((WINDOW_WIDTH - hint.get_width()) // 2, 88))

            margin = 20
            chart_top = 115
            chart_w = WINDOW_WIDTH - 2 * margin

            candles = candles_raw.get(sk, [])
            ind = indicators.get(sk, {})

            price_rect = (margin, chart_top, chart_w, PRICE_CHART_HEIGHT)
            draw_candle_chart(screen, small_font, selected_fruit, candles, ind, price_rect)

            macd_y = chart_top + PRICE_CHART_HEIGHT + 8
            macd_w = chart_w // 2 - margin // 2
            macd_rect = (margin, macd_y, macd_w, MACD_CHART_HEIGHT)
            draw_macd(screen, small_font, selected_fruit, candles, ind, macd_rect)

            rsi_x = margin + macd_w + margin
            rsi_rect = (rsi_x, macd_y, chart_w - macd_w - margin, RSI_CHART_HEIGHT)
            draw_rsi(screen, small_font, selected_fruit, candles, ind, rsi_rect)

            help_y = macd_y + MACD_CHART_HEIGHT + 5
            rsi_val = ind.get("rsi", 50) if ind else 50
            hist_val = ind.get("histogram", 0) if ind else 0
            rsi_status = "oversold" if rsi_val < 30 else "overbought" if rsi_val > 70 else "neutral"
            macd_status = "bullish" if hist_val > 0 else "bearish"
            info = small_font.render(
                f"Signal: {macd_status}  |  RSI: {rsi_status}  |  Histogram: {hist_val:+.2f}",
                True, YELLOW,
            )
            screen.blit(info, (margin, help_y))

            footer = small_font.render("ESC to exit", True, GRAY)
            screen.blit(footer, (20, WINDOW_HEIGHT - 20))

            pygame.display.flip()
            await asyncio.sleep(2)

    pygame.quit()
    sys.exit()


if __name__ == "__main__":
    asyncio.run(main())
