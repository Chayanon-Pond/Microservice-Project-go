package main

// Todo represents a simple todo item owned by an email.
type Todo struct {
    ID    int    `json:"id"`
    Owner string `json:"owner"`
    Text  string `json:"text"`
}
