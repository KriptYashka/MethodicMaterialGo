package main

import (
    "errors"
    "fmt"
)

func main() {
    // if-else
    x := 10
    if x > 0 {
        fmt.Println("positive")
    } else if x < 0 {
        fmt.Println("negative")
    } else {
        fmt.Println("zero")
    }

    // if with short statement
    if err := doSomething(); err != nil {
        fmt.Println("Error:", err)
    }

    // switch by value
    day := 3
    switch day {
    case 1:
        fmt.Println("Monday")
    case 2:
        fmt.Println("Tuesday")
    case 3:
        fmt.Println("Wednesday")
    default:
        fmt.Println("Unknown day")
    }

    // switch without expression (if-else chain)
    score := 85
    var grade string
    switch {
    case score >= 90:
        grade = "A"
    case score >= 80:
        grade = "B"
    case score >= 70:
        grade = "C"
    default:
        grade = "F"
    }
    fmt.Printf("Score %d -> Grade %s\n", score, grade)

    // fallthrough example
    i := 1
    switch i {
    case 1:
        fmt.Println("one")
        fallthrough
    case 2:
        fmt.Println("two")
    }
}

func doSomething() error {
    return errors.New("something went wrong")
}
