package main

import "fmt"

func main() {
    // Creating slices
    var s1 []int
    s2 := []int{1, 2, 3}
    s3 := make([]int, 5)
    s4 := make([]int, 3, 10)

    fmt.Printf("s1: %v len=%d cap=%d\n", s1, len(s1), cap(s1))
    fmt.Printf("s2: %v len=%d cap=%d\n", s2, len(s2), cap(s2))
    fmt.Printf("s3: %v len=%d cap=%d\n", s3, len(s3), cap(s3))
    fmt.Printf("s4: %v len=%d cap=%d\n", s4, len(s4), cap(s4))

    // append
    fmt.Println("\n--- append ---")
    var nums []int
    for i := 1; i <= 5; i++ {
        nums = append(nums, i)
        fmt.Printf("after append %d: len=%d cap=%d\n", i, len(nums), cap(nums))
    }
    fmt.Println("nums:", nums)

    // append multiple
    nums = append(nums, 6, 7, 8)
    fmt.Println("after multi-append:", nums)

    // sub-slices
    fmt.Println("\n--- sub-slices ---")
    arr := [5]int{1, 2, 3, 4, 5}
    slice := arr[1:4]
    fmt.Printf("arr=%v slice=%v\n", arr, slice)

    slice[0] = 99 // changes original array!
    fmt.Printf("after modify: arr=%v slice=%v\n", arr, slice)

    // slice expressions
    s := []int{0, 1, 2, 3, 4, 5}
    fmt.Println("s[:2]:", s[:2])
    fmt.Println("s[2:]:", s[2:])
    fmt.Println("s[1:4]:", s[1:4])

    // full slice expression (limiting capacity)
    sub := s[1:3:4]
    fmt.Printf("sub=%v len=%d cap=%d\n", sub, len(sub), cap(sub))

    // copy
    fmt.Println("\n--- copy ---")
    src := []int{1, 2, 3}
    dst := make([]int, len(src))
    n := copy(dst, src)
    fmt.Printf("copied %d elements: dst=%v\n", n, dst)

    // slice internals (append may reallocate)
    fmt.Println("\n--- append internals ---")
    a := make([]int, 2, 4)
    a[0] = 1
    a[1] = 2
    b := append(a, 3) // still within cap, shares backing array
    a[0] = 99         // this will be visible in b too!
    fmt.Println("a:", a)
    fmt.Println("b:", b)
}
