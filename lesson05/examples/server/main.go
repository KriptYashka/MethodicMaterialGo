package main

import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"
    "time"

    "fruitserver/game"
)

type GameServer struct {
    g      *game.Game
    ticker *time.Ticker
}

func NewGameServer(g *game.Game) *GameServer {
    return &GameServer{g: g}
}

func (s *GameServer) Start() {
    s.ticker = time.NewTicker(game.TickInterval)
    marketTicker := time.NewTicker(5 * time.Second)
    go func() {
        for range s.ticker.C {
            s.g.TickOnce(time.Now())
        }
    }()
    go func() {
        for range marketTicker.C {
            s.g.UpdateMarket(time.Now())
        }
    }()
}

func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
    writeJSON(w, status, map[string]string{"error": msg})
}

// ---- middleware ----

func recovery(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("PANIC recovered: %v", err)
                http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}

func (s *GameServer) withReadLock(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        s.g.Mu.RLock()
        defer s.g.Mu.RUnlock()
        next(w, r)
    }
}

// ---- handlers ----

func (s *GameServer) handleField(w http.ResponseWriter, r *http.Request) {
    writeJSON(w, http.StatusOK, map[string]any{
        "width":  s.g.Width,
        "height": s.g.Height,
        "tick":   s.g.Tick,
        "cells":  s.g.Field,
    })
}

func (s *GameServer) handleTeams(w http.ResponseWriter, r *http.Request) {
    writeJSON(w, http.StatusOK, s.g.Teams)
}

func (s *GameServer) handlePlant(w http.ResponseWriter, r *http.Request) {
    var req game.PlantRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }
    teamID := r.URL.Query().Get("team_id")
    tid, _ := strconv.Atoi(teamID)
    if tid < 1 || tid > len(s.g.Teams) {
        writeError(w, http.StatusBadRequest, "invalid team_id")
        return
    }
    s.g.Mu.Lock()
    err := s.g.Plant(req.CellIndex, tid, game.FruitType(req.FruitType), time.Now())
    s.g.Mu.Unlock()
    if err != nil {
        writeError(w, http.StatusBadRequest, err.Error())
        return
    }
    writeJSON(w, http.StatusOK, map[string]string{"status": "planted"})
}

func (s *GameServer) handleHarvest(w http.ResponseWriter, r *http.Request) {
    idxStr := r.PathValue("cell_idx")
    cellIdx, _ := strconv.Atoi(idxStr)
    teamID := r.URL.Query().Get("team_id")
    tid, _ := strconv.Atoi(teamID)
    if tid < 1 || tid > len(s.g.Teams) {
        writeError(w, http.StatusBadRequest, "invalid team_id")
        return
    }
    s.g.Mu.Lock()
    fruit, err := s.g.Harvest(cellIdx, tid)
    s.g.Mu.Unlock()
    if err != nil {
        writeError(w, http.StatusBadRequest, err.Error())
        return
    }
    writeJSON(w, http.StatusOK, map[string]any{"status": "harvested", "fruit": fruit})
}

func (s *GameServer) handleBuyCell(w http.ResponseWriter, r *http.Request) {
    var req game.BuyCellRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }
    teamID := r.URL.Query().Get("team_id")
    tid, _ := strconv.Atoi(teamID)
    s.g.Mu.Lock()
    err := s.g.BuyCell(req.CellIndex, tid)
    s.g.Mu.Unlock()
    if err != nil {
        writeError(w, http.StatusBadRequest, err.Error())
        return
    }
    writeJSON(w, http.StatusOK, map[string]string{"status": "cell bought"})
}

func (s *GameServer) handleBuyTool(w http.ResponseWriter, r *http.Request) {
    var req game.BuyToolRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }
    teamID := r.URL.Query().Get("team_id")
    tid, _ := strconv.Atoi(teamID)
    s.g.Mu.Lock()
    err := s.g.BuyTool(req.CellIndex, tid, game.ToolType(req.ToolType))
    s.g.Mu.Unlock()
    if err != nil {
        writeError(w, http.StatusBadRequest, err.Error())
        return
    }
    writeJSON(w, http.StatusOK, map[string]string{"status": "tool bought"})
}

func (s *GameServer) handleSell(w http.ResponseWriter, r *http.Request) {
    var req game.SellRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }
    teamID := r.URL.Query().Get("team_id")
    tid, _ := strconv.Atoi(teamID)
    s.g.Mu.Lock()
    revenue, err := s.g.SellFruits(tid, game.FruitType(req.FruitType), req.Quantity)
    s.g.Mu.Unlock()
    if err != nil {
        writeError(w, http.StatusBadRequest, err.Error())
        return
    }
    writeJSON(w, http.StatusOK, map[string]any{"status": "sold", "revenue": revenue})
}

func (s *GameServer) handleMarket(w http.ResponseWriter, r *http.Request) {
    writeJSON(w, http.StatusOK, s.g.Market)
}

func (s *GameServer) handleQuestGet(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")
    qid, _ := strconv.Atoi(idStr)
    for _, q := range s.g.Quests {
        if q.ID == qid {
            writeJSON(w, http.StatusOK, map[string]any{
                "id":       q.ID,
                "question": q.Question,
                "options":  q.Options,
                "reward":   q.Reward,
            })
            return
        }
    }
    writeError(w, http.StatusNotFound, "quest not found")
}

func (s *GameServer) handleQuestAnswer(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")
    qid, _ := strconv.Atoi(idStr)
    var req struct {
        Answer string `json:"answer"`
        TeamID int    `json:"team_id"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }
    correct, reward, err := s.g.AnswerQuest(qid, req.Answer)
    if err != nil {
        writeError(w, http.StatusNotFound, err.Error())
        return
    }
    if correct {
        s.g.Mu.Lock()
        team := s.g.GetTeam(req.TeamID)
        if team != nil {
            team.Money += reward
        }
        s.g.Mu.Unlock()
    }
    writeJSON(w, http.StatusOK, map[string]any{"correct": correct, "reward": reward})
}

func (s *GameServer) handleHouse(w http.ResponseWriter, r *http.Request) {
    idxStr := r.PathValue("cell_idx")
    cellIdx, _ := strconv.Atoi(idxStr)
    s.g.Mu.RLock()
    cell := s.g.GetCell(cellIdx)
    s.g.Mu.RUnlock()
    if cell == nil || cell.Type != game.CellHouse {
        writeError(w, http.StatusBadRequest, "not a house cell")
        return
    }
    phrase := s.g.HouseInteract(cellIdx)
    writeJSON(w, http.StatusOK, map[string]string{"phrase": phrase})
}

func main() {
    g := game.NewGame(10, 10, []string{"Team Alpha", "Team Beta", "Team Gamma"})
    s := NewGameServer(g)
    s.Start()

    mux := http.NewServeMux()
    mux.HandleFunc("GET /api/field", s.withReadLock(s.handleField))
    mux.HandleFunc("GET /api/teams", s.withReadLock(s.handleTeams))
    mux.HandleFunc("POST /api/plant", s.handlePlant)
    mux.HandleFunc("POST /api/harvest/{cell_idx}", s.handleHarvest)
    mux.HandleFunc("POST /api/buy-cell", s.handleBuyCell)
    mux.HandleFunc("POST /api/buy-tool", s.handleBuyTool)
    mux.HandleFunc("POST /api/sell", s.handleSell)
    mux.HandleFunc("GET /api/market", s.withReadLock(s.handleMarket))
    mux.HandleFunc("GET /api/quest/{id}", s.withReadLock(s.handleQuestGet))
    mux.HandleFunc("POST /api/quest/{id}", s.handleQuestAnswer)
    mux.HandleFunc("GET /api/house/{cell_idx}", s.handleHouse)

    addr := ":8080"
    log.Printf("Fruit Game Server starting on %s", addr)
    log.Printf("Teams: %d players on %dx%d field", len(g.Teams), g.Width, g.Height)
    log.Fatal(http.ListenAndServe(addr, recovery(cors(mux))))
}

func cors(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        next.ServeHTTP(w, r)
    })
}
