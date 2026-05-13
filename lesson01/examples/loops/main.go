package main

import "fmt"

func main() {
    // Classic for loop
    fmt.Println("--- classic for ---")
    for i := 0; i < 5; i++ {
        fmt.Printf("i = %d\n", i)
    }

    // For as while
    fmt.Println("\n--- for as while ---")
    n := 1
    for n < 100 {
        n *= 2
    }
    fmt.Printf("n = %d\n", n)

    // Infinite loop with break
    fmt.Println("\n--- infinite with break ---")
    sum := 0
    for {
        sum++
        if sum > 5 {
            break
        }
    }
    fmt.Printf("sum = %d\n", sum)

    // Continue
    fmt.Println("\n--- continue (odd numbers) ---")
    for i := 0; i < 10; i++ {
        if i%2 == 0 {
            continue
        }
        fmt.Printf("%d ", i)
    }
    fmt.Println()

    // range with slice
    fmt.Println("\n--- range over slice ---")
    nums := []int{10, 20, 30, 40, 50}
    for i, v := range nums {
        fmt.Printf("nums[%d] = %d\n", i, v)
    }

    // range with index only
    fmt.Println("\n--- range with index only ---")
    for i := range nums {
        nums[i] *= 2
    }
    fmt.Println(nums)

    // range with values ignored
    fmt.Println("\n--- range with values only ---")
    for _, v := range nums {
        fmt.Printf("value: %d\n", v)
    }

    // range over string (runes)
    fmt.Println("\n--- range over string ---")
    for i, r := range "Привет, Go!" {
        fmt.Printf("%d: %c (%U)\n", i, r, r)
    }
}
