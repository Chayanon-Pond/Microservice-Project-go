package main

// User represents a simple user profile stored in the user-service.
type User struct {
    ID    int    `json:"id"`
    Email string `json:"email"`
    Name  string `json:"name"`
    Role  string `json:"role"`
}
