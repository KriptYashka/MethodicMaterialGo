package main

import "fmt"

// ---- basic functions ----

func add(a, b int) int {
    return a + b
}

// multiple return values
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

// named returns
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return
}

// variadic
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

// function as value
func operate(op string) func(int, int) int {
    switch op {
    case "add":
        return func(a, b int) int { return a + b }
    case "sub":
        return func(a, b int) int { return a - b }
    default:
        return func(a, b int) int { return 0 }
    }
}

// closure — counter
func counter() func() int {
    i := 0
    return func() int {
        i++
        return i
    }
}

// closure — fibonacci
func fibonacci() func() int {
    a, b := 0, 1
    return func() int {
        result := a
        a, b = b, a+b
        return result
    }
}

// defer example
func deferDemo() {
    defer fmt.Println("3. defer: cleanup")
    defer fmt.Println("2. defer: second")
    fmt.Println("1. normal execution")
}

func main() {
    fmt.Println("--- basic ---")
    fmt.Println("add(3,5):", add(3, 5))

    q, err := divide(10, 3)
    if err != nil {
        fmt.Println("error:", err)
    } else {
        fmt.Println("divide(10,3):", q)
    }

    fmt.Println("split(17):", split(17))

    fmt.Println("\n--- variadic ---")
    fmt.Println("sum(1,2,3):", sum(1, 2, 3))
    fmt.Println("sum(1,2,3,4,5):", sum(1, 2, 3, 4, 5))

    nums := []int{10, 20, 30}
    fmt.Println("sum(nums...):", sum(nums...))

    fmt.Println("\n--- function as value ---")
    op := operate("add")
    fmt.Println("op(5,3):", op(5, 3))

    fmt.Println("\n--- closures ---")
    c := counter()
    fmt.Println(c()) // 1
    fmt.Println(c()) // 2
    fmt.Println(c()) // 3

    fmt.Println("\n--- fibonacci ---")
    fib := fibonacci()
    for i := 0; i < 10; i++ {
        fmt.Printf("%d ", fib())
    }
    fmt.Println()

    fmt.Println("\n--- defer ---")
    deferDemo()
}
