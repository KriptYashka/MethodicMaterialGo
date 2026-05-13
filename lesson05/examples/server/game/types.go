package game

import "time"

type CellType int

const (
    CellWater     CellType = iota
    CellEarth
    CellMountain
    CellLegendary
    CellHouse
)

func (ct CellType) String() string {
    switch ct {
    case CellWater:
        return "Water"
    case CellEarth:
        return "Earth"
    case CellMountain:
        return "Mountain"
    case CellLegendary:
        return "Legendary"
    case CellHouse:
        return "House"
    default:
        return "Unknown"
    }
}

func (ct CellType) Color() string {
    switch ct {
    case CellWater:
        return "blue"
    case CellEarth:
        return "brown"
    case CellMountain:
        return "gray"
    case CellLegendary:
        return "gold"
    case CellHouse:
        return "green"
    default:
        return "white"
    }
}

type FruitType int

const (
    FruitWatermelon FruitType = iota
    FruitMelon
    FruitRaspberry
)

var FruitNames = map[FruitType]string{
    FruitWatermelon: "watermelon",
    FruitMelon:      "melon",
    FruitRaspberry:  "raspberry",
}

type ToolType int

const (
    ToolFertilizer ToolType = iota
    ToolFence
)

var ToolNames = map[ToolType]string{
    ToolFertilizer: "fertilizer",
    ToolFence:      "fence",
}

var ToolPrices = map[ToolType]int{
    ToolFertilizer: 50,
    ToolFence:      80,
}

// fruit market config
type FruitConfig struct {
    MinPrice   float64
    MaxPrice   float64
    Dispersion float64
    Pullback   int // remaining pullback ticks (runtime)
}

var FruitMarketConfig = map[FruitType]*FruitConfig{
    FruitWatermelon: {MinPrice: 5, MaxPrice: 40, Dispersion: 0.6},
    FruitMelon:      {MinPrice: 3, MaxPrice: 35, Dispersion: 0.5},
    FruitRaspberry:  {MinPrice: 6, MaxPrice: 45, Dispersion: 0.7},
}

type Fruit struct {
    Type      FruitType `json:"type"`
    TeamID    int       `json:"team_id"`
    PlantedAt time.Time `json:"planted_at"`
    Growth    float64   `json:"growth"`
    Ripe      bool      `json:"ripe"`
}

type Cell struct {
    Index    int      `json:"index"`
    Type     CellType `json:"type"`
    TeamID   int      `json:"team_id"`
    Fruit    *Fruit   `json:"fruit,omitempty"`
    Tool     ToolType `json:"tool,omitempty"`
    ExpireAt time.Time `json:"expire_at,omitempty"`
}

type Team struct {
    ID        int                `json:"id"`
    Name      string             `json:"name"`
    Money     int                `json:"money"`
    Inventory map[FruitType]int  `json:"inventory"`
}

type PlantRequest struct {
    CellIndex int `json:"cell_index"`
    FruitType int `json:"fruit_type"`
}

type SellRequest struct {
    FruitType int `json:"fruit_type"`
    Quantity  int `json:"quantity"`
}

type BuyCellRequest struct {
    CellIndex int `json:"cell_index"`
}

type BuyToolRequest struct {
    CellIndex int `json:"cell_index"`
    ToolType  int `json:"tool_type"`
}

type Quest struct {
    ID       int      `json:"id"`
    Question string   `json:"question"`
    Answer   string   `json:"-"`
    Reward   int      `json:"reward"`
    Options  []string `json:"options,omitempty"`
}

type Candle struct {
    Time  time.Time `json:"time"`
    Open  float64   `json:"open"`
    High  float64   `json:"high"`
    Low   float64   `json:"low"`
    Close float64   `json:"close"`
}

type MarketIndicators struct {
    MACDLine   float64 `json:"macd_line"`
    SignalLine float64 `json:"signal_line"`
    Histogram  float64 `json:"histogram"`
    RSI        float64 `json:"rsi"`
}

type MarketState struct {
    Prices      map[FruitType]float64          `json:"prices"`
    Candles     map[FruitType][]Candle          `json:"candles"`
    Indicators  map[FruitType]MarketIndicators  `json:"indicators"`
    priceHistory map[FruitType][]float64
    pullbackTicks map[FruitType]int
    pullbackDir   map[FruitType]int  // +1 = force up, -1 = force down
}

func NewMarketState() *MarketState {
    ms := &MarketState{
        Prices: map[FruitType]float64{
            FruitWatermelon: 15.0,
            FruitMelon:      12.0,
            FruitRaspberry:  18.0,
        },
        Candles:    map[FruitType][]Candle{},
        Indicators: map[FruitType]MarketIndicators{},
        priceHistory: map[FruitType][]float64{},
        pullbackTicks: map[FruitType]int{},
        pullbackDir:   map[FruitType]int{},
    }
    // seed initial history for indicators
    for _, ft := range []FruitType{FruitWatermelon, FruitMelon, FruitRaspberry} {
        h := make([]float64, 30)
        for i := range h {
            h[i] = ms.Prices[ft]
        }
        ms.priceHistory[ft] = h
    }
    return ms
}

func CellPrice(cellType CellType) int {
    switch cellType {
    case CellWater:
        return 30
    case CellEarth:
        return 25
    case CellMountain:
        return 35
    case CellLegendary:
        return 100
    default:
        return 99999
    }
}

func Affinity(fruit FruitType, cell CellType) float64 {
    switch {
    case cell == CellLegendary:
        return 4.0
    case cell == CellHouse:
        return 0.0
    case fruit == FruitWatermelon && cell == CellWater:
        return 2.0
    case fruit == FruitWatermelon && cell == CellEarth:
        return 2.0
    case fruit == FruitWatermelon && cell == CellMountain:
        return 0.5
    case fruit == FruitMelon && cell == CellWater:
        return 2.0
    case fruit == FruitMelon && cell == CellEarth:
        return 0.5
    case fruit == FruitMelon && cell == CellMountain:
        return 2.0
    case fruit == FruitRaspberry && cell == CellWater:
        return 0.5
    case fruit == FruitRaspberry && cell == CellEarth:
        return 2.0
    case fruit == FruitRaspberry && cell == CellMountain:
        return 2.0
    default:
        return 1.0
    }
}

const GrowthThreshold = 1.0
const HarvestWindow = 30 * time.Second
const TickInterval = 1 * time.Second

var HousePhrases = []string{
    "This is my house, get out!",
    "Who's there?",
    "Go away, I'm growing my own fruits!",
    "This land is protected by ancient spirits!",
    "You shall not pass!",
    "My precious... fruits!",
    "No trespassing! Violators will be watered!",
    "I'm watching you...",
}
