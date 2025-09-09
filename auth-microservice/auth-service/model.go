package main

type User struct {
	ID        int
	Email     string
	Password  string
	Role      string
	CreatedAt string // หรือ time.Time ถ้า import "time"
}
