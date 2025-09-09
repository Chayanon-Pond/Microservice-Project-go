package main

import "errors"

// simple in-memory repository
type UserRepo struct {
    users map[string]User // key by email
}

func NewUserRepo() *UserRepo {
    r := &UserRepo{users: make(map[string]User)}
    // seed with an example user
    r.users["alice@example.com"] = User{ID: 1, Email: "alice@example.com", Name: "Alice", Role: "user"}
    r.users["admin@example.com"] = User{ID: 2, Email: "admin@example.com", Name: "Admin", Role: "admin"}
    return r
}

func (r *UserRepo) GetByEmail(email string) (User, error) {
    if u, ok := r.users[email]; ok {
        return u, nil
    }
    return User{}, errors.New("user not found")
}
