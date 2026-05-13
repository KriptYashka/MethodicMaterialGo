package game

import (
    "math"
    "math/rand"
    "sync"
    "time"
)

type Game struct {
    Mu     sync.RWMutex
    Field  []Cell     `json:"field"`
    Teams  []Team     `json:"teams"`
    Market *MarketState `json:"market"`
    Quests []Quest    `json:"quests"`
    Tick   int        `json:"tick"`
    Width  int        `json:"width"`
  Height int        `json:"height"`
}

func NewGame(width, height int, teamNames []string) *Game {
    g := &Game{
        Width:  width,
        Height: height,
        Field:  generateField(width, height),
        Market: NewMarketState(),
        Quests: generateQuests(),
        Tick:   0,
    }
    for i, name := range teamNames {
        g.Teams = append(g.Teams, Team{
            ID:        i + 1,
            Name:      name,
            Money:     100,
            Inventory: map[FruitType]int{},
        })
    }
    return g
}

func generateField(w, h int) []Cell {
    total := w * h
    cells := make([]Cell, total)
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    legendaryCount := total * 5 / 100
    legendaryPos := rng.Perm(total)[:legendaryCount]
    legendarySet := make(map[int]bool)
    for _, p := range legendaryPos {
        legendarySet[p] = true
    }
    houseCount := total * 5 / 100
    avail := make([]int, 0)
    for i := 0; i < total; i++ {
        if !legendarySet[i] {
            avail = append(avail, i)
        }
    }
    rng.Shuffle(len(avail), func(i, j int) { avail[i], avail[j] = avail[j], avail[i] })
    houseSet := make(map[int]bool)
    for _, p := range avail[:houseCount] {
        houseSet[p] = true
    }
    for i := 0; i < total; i++ {
        cells[i] = Cell{Index: i}
        switch {
        case legendarySet[i]:
            cells[i].Type = CellLegendary
        case houseSet[i]:
            cells[i].Type = CellHouse
        default:
            types := []CellType{CellWater, CellEarth, CellMountain}
            cells[i].Type = types[rng.Intn(len(types))]
        }
    }
    return cells
}

func generateQuests() []Quest {
    return []Quest{
        {ID: 1, Question: "What HTTP method is typically used to create a resource?", Answer: "POST", Reward: 30, Options: []string{"GET", "POST", "PUT", "DELETE"}},
        {ID: 2, Question: "Which Go keyword is used to launch a goroutine?", Answer: "go", Reward: 30, Options: []string{"go", "async", "spawn", "thread"}},
        {ID: 3, Question: "What is the zero value of an int in Go?", Answer: "0", Reward: 20, Options: []string{"0", "nil", "null", "undefined"}},
        {ID: 4, Question: "Which data structure is best for key-value storage in Go?", Answer: "map", Reward: 25, Options: []string{"map", "slice", "array", "struct"}},
        {ID: 5, Question: "What port does HTTP use by default?", Answer: "80", Reward: 20, Options: []string{"80", "443", "8080", "3000"}},
        {ID: 6, Question: "What package in Go handles JSON encoding/decoding?", Answer: "encoding/json", Reward: 35, Options: []string{"encoding/json", "json", "net/json", "fmt/json"}},
        {ID: 7, Question: "What is the HTTP status code for 'Not Found'?", Answer: "404", Reward: 20, Options: []string{"404", "400", "403", "500"}},
        {ID: 8, Question: "Which function creates a new slice with given length and capacity?", Answer: "make", Reward: 25, Options: []string{"make", "new", "create", "alloc"}},
        {ID: 9, Question: "What does 'defer' do in Go?", Answer: "delays execution until surrounding function returns", Reward: 40, Options: []string{"delays execution until surrounding function returns", "creates a new goroutine", "throws an error", "allocates memory"}},
        {ID: 10, Question: "What does HTML stand for?", Answer: "HyperText Markup Language", Reward: 20, Options: []string{"HyperText Markup Language", "High Tech Machine Learning", "Home Tool Markup Language", "HyperText Modern Language"}},
    }
}

func (g *Game) GetCell(idx int) *Cell {
    if idx < 0 || idx >= len(g.Field) {
        return nil
    }
    return &g.Field[idx]
}

func (g *Game) GetTeam(id int) *Team {
    for i := range g.Teams {
        if g.Teams[i].ID == id {
            return &g.Teams[i]
        }
    }
    return nil
}

func (g *Game) TickOnce(now time.Time) {
    g.Mu.Lock()
    defer g.Mu.Unlock()
    g.Tick++

    for i := range g.Field {
        cell := &g.Field[i]
        if cell.Fruit == nil || cell.Fruit.Ripe {
            continue
        }
        aff := Affinity(cell.Fruit.Type, cell.Type)
        if cell.Tool == ToolFertilizer {
            aff *= 2.0
        }
        cell.Fruit.Growth += 0.02 * aff
        if cell.Fruit.Growth >= GrowthThreshold {
            cell.Fruit.Ripe = true
            cell.ExpireAt = now.Add(HarvestWindow)
        }
    }

    for i := range g.Field {
        cell := &g.Field[i]
        if cell.Fruit != nil && cell.Fruit.Ripe && !now.Before(cell.ExpireAt) {
            cell.Fruit = nil
            if cell.Tool != ToolFence {
                cell.TeamID = 0
            }
            cell.Tool = 0
        }
    }

    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    if rng.Float64() < 0.02 {
        idx := rng.Intn(len(g.Field))
        cell := &g.Field[idx]
        if cell.TeamID == 0 && cell.Type != CellHouse && cell.Type != CellLegendary {
            cell.Type = CellLegendary
        }
    }
}

func (g *Game) Plant(cellIdx, teamID int, fruitType FruitType, now time.Time) error {
    cell := g.GetCell(cellIdx)
    if cell == nil {
        return ErrInvalidCell
    }
    if cell.TeamID != teamID {
        return ErrNotYourCell
    }
    if cell.Type == CellHouse {
        return ErrHouseCell
    }
    if cell.Fruit != nil {
        return ErrCellOccupied
    }
    cell.Fruit = &Fruit{
        Type:   fruitType,
        TeamID: teamID,
        PlantedAt: now,
    }
    return nil
}

func (g *Game) Harvest(cellIdx, teamID int) (*Fruit, error) {
    cell := g.GetCell(cellIdx)
    if cell == nil {
        return nil, ErrInvalidCell
    }
    if cell.Fruit == nil {
        return nil, ErrNoFruit
    }
    if cell.Fruit.TeamID != teamID {
        return nil, ErrNotYourFruit
    }
    if !cell.Fruit.Ripe {
        return nil, ErrNotRipe
    }
    fruit := cell.Fruit
    cell.Fruit = nil
    if cell.Tool != ToolFence {
        cell.TeamID = 0
    }
    cell.Tool = 0
    team := g.GetTeam(teamID)
    team.Inventory[fruit.Type]++
    return fruit, nil
}

func (g *Game) BuyCell(cellIdx, teamID int) error {
    cell := g.GetCell(cellIdx)
    if cell == nil {
        return ErrInvalidCell
    }
    if cell.TeamID != 0 {
        return ErrCellTaken
    }
    team := g.GetTeam(teamID)
    price := CellPrice(cell.Type)
    if team.Money < price {
        return ErrNotEnoughMoney
    }
    team.Money -= price
    cell.TeamID = teamID
    return nil
}

func (g *Game) BuyTool(cellIdx, teamID int, tool ToolType) error {
    cell := g.GetCell(cellIdx)
    if cell == nil {
        return ErrInvalidCell
    }
    if cell.TeamID != teamID {
        return ErrNotYourCell
    }
    team := g.GetTeam(teamID)
    price := ToolPrices[tool]
    if team.Money < price {
        return ErrNotEnoughMoney
    }
    team.Money -= price
    cell.Tool = tool
    return nil
}

func (g *Game) SellFruits(teamID int, fruitType FruitType, qty int) (int, error) {
    team := g.GetTeam(teamID)
    if team.Inventory[fruitType] < qty {
        return 0, ErrNotEnoughFruit
    }
    price := g.Market.Prices[fruitType]
    revenue := int(price * float64(qty))
    team.Inventory[fruitType] -= qty
    team.Money += revenue
    return revenue, nil
}

// ---- EMA ----
func ema(values []float64, period int) []float64 {
    if len(values) == 0 {
        return nil
    }
    k := 2.0 / float64(period+1)
    result := make([]float64, len(values))
    result[0] = values[0]
    for i := 1; i < len(values); i++ {
        result[i] = values[i]*k + result[i-1]*(1-k)
    }
    return result
}

func (ms *MarketState) calcMACD(prices []float64) (macdLine, signalLine, histogram float64) {
    if len(prices) < 26 {
        return 0, 0, 0
    }
    e12 := ema(prices, 12)
    e26 := ema(prices, 26)

    // macd line = ema12 - ema26
    macdVals := make([]float64, len(prices))
    for i := range macdVals {
        macdVals[i] = e12[i] - e26[i]
    }

    // signal = ema(9) of macd line
    signalVals := ema(macdVals, 9)

    macdLine = macdVals[len(macdVals)-1]
    signalLine = signalVals[len(signalVals)-1]
    histogram = macdLine - signalLine
    return
}

func (ms *MarketState) calcRSI(prices []float64) float64 {
    period := 14
    if len(prices) < period+1 {
        return 50
    }
    gains := 0.0
    losses := 0.0
    for i := len(prices) - period; i < len(prices); i++ {
        diff := prices[i] - prices[i-1]
        if diff > 0 {
            gains += diff
        } else {
            losses -= diff
        }
    }
    avgGain := gains / float64(period)
    avgLoss := losses / float64(period)
    if avgLoss == 0 {
        return 100
    }
    rs := avgGain / avgLoss
    return 100 - (100 / (1 + rs))
}

func (g *Game) UpdateMarket(now time.Time) {
    g.Mu.Lock()
    defer g.Mu.Unlock()

    rng := rand.New(rand.NewSource(time.Now().UnixNano()))

    for _, ft := range []FruitType{FruitWatermelon, FruitMelon, FruitRaspberry} {
        cfg := FruitMarketConfig[ft]
        oldPrice := g.Market.Prices[ft]

        // MACD/RSI-based direction probability
        dirProb := 0.5 // base: 50% up
        hist := g.Market.priceHistory[ft]
        if len(hist) >= 26 {
            _, _, macdHist := g.Market.calcMACD(hist)
            rsi := g.Market.calcRSI(hist)

            // RSI influence
            if rsi > 70 {
                dirProb -= 0.15 // overbought → more likely to go down
            } else if rsi < 30 {
                dirProb += 0.15 // oversold → more likely to go up
            }

            // MACD histogram influence
            if macdHist > 0 {
                dirProb += 0.10 // positive momentum → more likely up
            } else {
                dirProb -= 0.10
            }
        }

        // pullback mechanism
        pb := g.Market.pullbackTicks[ft]
        if pb > 0 {
            dir := g.Market.pullbackDir[ft]
            if dir > 0 {
                dirProb = 0.85
            } else {
                dirProb = 0.15
            }
            g.Market.pullbackTicks[ft] = pb - 1
        }

        // generate price change
        delta := (rng.Float64() - 0.5) * 2.0 * cfg.Dispersion * 4.0
        if rng.Float64() > dirProb {
            delta = -math.Abs(delta)
        } else {
            delta = math.Abs(delta)
        }
        newPrice := oldPrice + delta

        // check bounds — trigger pullback if exceeded
        if newPrice < cfg.MinPrice {
            newPrice = cfg.MinPrice
            g.Market.pullbackTicks[ft] = 2 + rng.Intn(2) // 2-3 ticks
            g.Market.pullbackDir[ft] = 1 // push up
        } else if newPrice > cfg.MaxPrice {
            newPrice = cfg.MaxPrice
            g.Market.pullbackTicks[ft] = 2 + rng.Intn(2)
            g.Market.pullbackDir[ft] = -1 // push down
        }

        g.Market.Prices[ft] = newPrice

        // record candle
        g.Market.Candles[ft] = append(g.Market.Candles[ft], Candle{
            Time:  now,
            Open:  oldPrice,
            Close: newPrice,
            High:  math.Max(oldPrice, newPrice),
            Low:   math.Min(oldPrice, newPrice),
        })
        if len(g.Market.Candles[ft]) > 200 {
            g.Market.Candles[ft] = g.Market.Candles[ft][len(g.Market.Candles[ft])-200:]
        }

        // update price history for indicators
        g.Market.priceHistory[ft] = append(g.Market.priceHistory[ft], newPrice)
        if len(g.Market.priceHistory[ft]) > 60 {
            g.Market.priceHistory[ft] = g.Market.priceHistory[ft][len(g.Market.priceHistory[ft])-60:]
        }
        hist = g.Market.priceHistory[ft]

        // calculate indicators
        macdLine, signalLine, histogram := g.Market.calcMACD(hist)
        rsi := g.Market.calcRSI(hist)
        g.Market.Indicators[ft] = MarketIndicators{
            MACDLine:   macdLine,
            SignalLine: signalLine,
            Histogram:  histogram,
            RSI:        rsi,
        }
    }
}

func (g *Game) AnswerQuest(questID int, answer string) (bool, int, error) {
    for _, q := range g.Quests {
        if q.ID == questID {
            if q.Answer == answer {
                return true, q.Reward, nil
            }
            return false, 0, nil
        }
    }
    return false, 0, ErrQuestNotFound
}

func (g *Game) HouseInteract(cellIdx int) string {
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    return HousePhrases[rng.Intn(len(HousePhrases))]
}
