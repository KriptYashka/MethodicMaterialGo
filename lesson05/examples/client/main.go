package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var baseURL = "http://localhost:3845"

type Team struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	Money     int            `json:"money"`
	Inventory map[string]int `json:"inventory"`
}

type Cell struct {
	Index  int    `json:"index"`
	Type   string `json:"type"`
	TeamID int    `json:"team_id"`
	Fruit  any    `json:"fruit,omitempty"`
}

type FieldResponse struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Cells  []Cell `json:"cells"`
}

type MarketState struct {
	Prices  map[string]float64 `json:"prices"`
	Candles map[string][]any   `json:"candles"`
}

func apiGet(path string) ([]byte, error) {
	resp, err := http.Get(baseURL + path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func apiPost(path string, body any) ([]byte, error) {
	data, _ := json.Marshal(body)
	resp, err := http.Post(baseURL+path, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func cmdField() {
	data, err := apiGet("/api/field")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var fr FieldResponse
	json.Unmarshal(data, &fr)
	fmt.Printf("Field: %dx%d\n", fr.Width, fr.Height)
	for i, c := range fr.Cells {
		if i > 0 && i%fr.Width == 0 {
			fmt.Println()
		}
		sym := "."
		switch c.Type {
		case "Water":
			sym = "~"
		case "Earth":
			sym = "#"
		case "Mountain":
			sym = "^"
		case "Legendary":
			sym = "*"
		case "House":
			sym = "H"
		}
		if c.TeamID > 0 {
			sym = strconv.Itoa(c.TeamID)
		}
		fmt.Printf("%2s ", sym)
	}
	fmt.Println()
}

func cmdTeams() {
	data, _ := apiGet("/api/teams")
	var teams []Team
	json.Unmarshal(data, &teams)
	for _, t := range teams {
		fmt.Printf("Team %d (%s): $%d\n", t.ID, t.Name, t.Money)
		for k, v := range t.Inventory {
			fmt.Printf("  %s: %d\n", k, v)
		}
	}
}

func cmdMarket() {
	data, _ := apiGet("/api/market")
	var ms MarketState
	json.Unmarshal(data, &ms)
	fmt.Println("Market prices:")
	for k, v := range ms.Prices {
		fmt.Printf("  %s: $%.2f\n", k, v)
	}
}

func cmdPlant(teamID, cellIdx, fruitType int) {
	body := map[string]int{
		"cell_index": cellIdx,
		"fruit_type": fruitType,
	}
	data, _ := apiPost("/api/plant?team_id="+strconv.Itoa(teamID), body)
	fmt.Println(string(data))
}

func cmdHarvest(teamID, cellIdx int) {
	data, _ := apiPost("/api/harvest/"+strconv.Itoa(cellIdx)+"?team_id="+strconv.Itoa(teamID), nil)
	fmt.Println(string(data))
}

func cmdBuyCell(teamID, cellIdx int) {
	body := map[string]int{"cell_index": cellIdx}
	data, _ := apiPost("/api/buy-cell?team_id="+strconv.Itoa(teamID), body)
	fmt.Println(string(data))
}

func cmdSell(teamID, fruitType, qty int) {
	body := map[string]int{"fruit_type": fruitType, "quantity": qty}
	data, _ := apiPost("/api/sell?team_id="+strconv.Itoa(teamID), body)
	fmt.Println(string(data))
}

func cmdQuest(questID int) {
	data, _ := apiGet("/api/quest/" + strconv.Itoa(questID))
	fmt.Println(string(data))
}

func cmdAnswer(questID int, answer string, teamID int) {
	body := map[string]any{"answer": answer, "team_id": teamID}
	data, _ := apiPost("/api/quest/"+strconv.Itoa(questID), body)
	fmt.Println(string(data))
}

func main() {
	flag.StringVar(&baseURL, "url", "http://localhost:3845", "server URL")
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Usage: fruitclient <command> [args...]")
		fmt.Println("Commands:")
		fmt.Println("  field")
		fmt.Println("  teams")
		fmt.Println("  market")
		fmt.Println("  plant <team_id> <cell_idx> <fruit_type(0=watermelon,1=melon,2=raspberry)>")
		fmt.Println("  harvest <team_id> <cell_idx>")
		fmt.Println("  buy-cell <team_id> <cell_idx>")
		fmt.Println("  sell <team_id> <fruit_type> <quantity>")
		fmt.Println("  quest <quest_id>")
		fmt.Println("  answer <quest_id> <answer> <team_id>")
		return
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "field":
		cmdField()
	case "teams":
		cmdTeams()
	case "market":
		cmdMarket()
	case "plant":
		tid, _ := strconv.Atoi(args[0])
		ci, _ := strconv.Atoi(args[1])
		ft, _ := strconv.Atoi(args[2])
		cmdPlant(tid, ci, ft)
	case "harvest":
		tid, _ := strconv.Atoi(args[0])
		ci, _ := strconv.Atoi(args[1])
		cmdHarvest(tid, ci)
	case "buy-cell":
		tid, _ := strconv.Atoi(args[0])
		ci, _ := strconv.Atoi(args[1])
		cmdBuyCell(tid, ci)
	case "sell":
		tid, _ := strconv.Atoi(args[0])
		ft, _ := strconv.Atoi(args[1])
		qty, _ := strconv.Atoi(args[2])
		cmdSell(tid, ft, qty)
	case "quest":
		qid, _ := strconv.Atoi(args[0])
		cmdQuest(qid)
	case "answer":
		qid, _ := strconv.Atoi(args[0])
		answer := strings.Join(args[1:len(args)-1], " ")
		tid, _ := strconv.Atoi(args[len(args)-1])
		cmdAnswer(qid, answer, tid)
	default:
		fmt.Println("Unknown command:", cmd)
	}
}
