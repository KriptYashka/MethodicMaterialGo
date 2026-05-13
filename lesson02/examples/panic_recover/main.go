package main

import (
    "errors"
    "fmt"
)

// safe division using recover
func safeDivide(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            result = 0
            err = errors.New("division by zero")
        }
    }()
    return a / b, nil
}

// panic on unexpected state
func mustInit() {
    panic("critical init failure") //nolint:forbidigo
}

// Teardown with defer
func work() (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("work panicked: %v", r)
        }
    }()
    // simulate work
    panic("unexpected state")
}

func main() {
    fmt.Println("--- safe divide ---")
    for _, pair := range [][2]int{{10, 2}, {5, 0}, {8, 4}} {
        result, err := safeDivide(pair[0], pair[1])
        if err != nil {
            fmt.Printf("%d/%d: error — %v\n", pair[0], pair[1], err)
        } else {
            fmt.Printf("%d/%d = %d\n", pair[0], pair[1], result)
        }
    }

    fmt.Println("\n--- panic in function ---")
    if err := work(); err != nil {
        fmt.Println("recovered:", err)
    }

    fmt.Println("\n--- program continues normally ---")
    fmt.Println("this line is reached because panics were recovered")
}
