package main

import "fmt"

func main() {
    // create
    ages := make(map[string]int)
    ages["Alice"] = 30
    ages["Bob"] = 25
    ages["Charlie"] = 35

    fmt.Println("ages:", ages)

    // literal
    scores := map[string]int{
        "math":    90,
        "physics": 85,
        "english": 78,
    }
    fmt.Println("scores:", scores)

    // get
    aliceAge := ages["Alice"]
    fmt.Println("Alice age:", aliceAge)

    // comma-ok idiom
    val, ok := ages["Nobody"]
    if !ok {
        fmt.Println("Nobody not found")
    } else {
        fmt.Println("Nobody age:", val)
    }

    // delete
    delete(ages, "Charlie")
    fmt.Println("after delete:", ages)

    // len
    fmt.Println("len:", len(ages))

    // iteration (order not guaranteed!)
    fmt.Println("\n--- iteration ---")
    for k, v := range scores {
        fmt.Printf("%s -> %d\n", k, v)
    }

    // map as set
    fmt.Println("\n--- set ---")
    set := make(map[string]bool)
    items := []string{"apple", "banana", "apple", "orange", "banana"}
    for _, item := range items {
        set[item] = true
    }
    fmt.Println("unique items:")
    for item := range set {
        fmt.Println(" -", item)
    }

    // map with struct keys
    fmt.Println("\n--- struct key ---")
    type Point struct {
        X, Y int
    }
    grid := make(map[Point]string)
    grid[Point{1, 2}] = "start"
    grid[Point{3, 4}] = "end"
    fmt.Println("point {1,2}:", grid[Point{1, 2}])

    // nested maps
    fmt.Println("\n--- nested maps ---")
    users := map[string]map[string]string{
        "alice": {
            "email": "alice@example.com",
            "role":  "admin",
        },
        "bob": {
            "email": "bob@example.com",
            "role":  "user",
        },
    }
    for name, info := range users {
        fmt.Printf("%s: email=%s role=%s\n", name, info["email"], info["role"])
    }
}
