package main

import (
    "fmt"
    "strconv"
)

func main() {
    // var declaration
    var name string = "Alice"
    var age = 30
    city := "Moscow"

    fmt.Println(name, age, city)

    // Multiple variables
    var x, y int = 1, 2
    a, b := "hello", 42
    fmt.Println(x, y, a, b)

    // Group declaration
    var (
        width  int    = 100
        height int    = 200
        title  string = "Image"
    )
    fmt.Println(width, height, title)

    // Zero values
    var (
        zeroInt    int
        zeroFloat  float64
        zeroString string
        zeroBool   bool
    )
    fmt.Printf("zero values: %d, %f, %q, %v\n",
        zeroInt, zeroFloat, zeroString, zeroBool)

    // Type conversion
    var i int = 42
    var f float64 = float64(i)
    var s string = strconv.Itoa(i)
    fmt.Printf("int=%d float=%f string=%s\n", i, f, s)

    n, _ := strconv.Atoi("123")
    fmt.Printf("parsed: %d\n", n+1)

    // Constants
    const Pi = 3.14159
    const (
        StatusOK       = 200
        StatusNotFound = 404
    )
    fmt.Println(Pi, StatusOK, StatusNotFound)

    // iota
    const (
        Red  = iota
        Green
        Blue
    )
    fmt.Println(Red, Green, Blue)
}
