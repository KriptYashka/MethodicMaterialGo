package main

import (
    "errors"
    "fmt"
)

// sentinel errors
var (
    ErrNotFound = errors.New("not found")
    ErrInvalid  = errors.New("invalid input")
)

// custom error type
type ValidationError struct {
    Field string
    Value any
    Msg   string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error: %s=%v — %s", e.Field, e.Value, e.Msg)
}

// ---- business logic ----

type User struct {
    ID   int
    Name string
    Age  int
}

var users = map[int]User{
    1: {ID: 1, Name: "Alice", Age: 30},
    2: {ID: 2, Name: "Bob", Age: 25},
}

func GetUser(id int) (*User, error) {
    if id < 1 {
        return nil, ErrInvalid
    }
    u, ok := users[id]
    if !ok {
        return nil, fmt.Errorf("get user: %w", ErrNotFound)
    }
    return &u, nil
}

func ValidateUser(u User) error {
    if u.Name == "" {
        return &ValidationError{Field: "Name", Value: u.Name, Msg: "must not be empty"}
    }
    if u.Age < 0 || u.Age > 150 {
        return &ValidationError{Field: "Age", Value: u.Age, Msg: "out of range"}
    }
    return nil
}

func SaveUser(u User) error {
    if err := ValidateUser(u); err != nil {
        return fmt.Errorf("save user: %w", err)
    }
    users[u.ID] = u
    return nil
}

func main() {
    fmt.Println("--- sentinel errors ---")
    _, err := GetUser(99)
    if errors.Is(err, ErrNotFound) {
        fmt.Println("handled: user not found, create new")
    }
    fmt.Println("full error:", err)

    fmt.Println("\n--- wrapped error ---")
    err = SaveUser(User{ID: 3, Name: "", Age: 25})
    if err != nil {
        fmt.Println("save error:", err)
        var ve *ValidationError
        if errors.As(err, &ve) {
            fmt.Printf("validation field=%s value=%v msg=%s\n", ve.Field, ve.Value, ve.Msg)
        }
    }

    fmt.Println("\n--- multiple error checks ---")
    for _, id := range []int{0, 1, 99} {
        u, err := GetUser(id)
        switch {
        case errors.Is(err, ErrInvalid):
            fmt.Printf("id=%d: invalid id\n", id)
        case errors.Is(err, ErrNotFound):
            fmt.Printf("id=%d: not found\n", id)
        case err != nil:
            fmt.Printf("id=%d: unexpected error: %v\n", id, err)
        default:
            fmt.Printf("id=%d: %+v\n", id, *u)
        }
    }
}
