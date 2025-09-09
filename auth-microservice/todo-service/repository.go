package main

import "sync"

type TodoRepo struct {
    mu    sync.Mutex
    items map[int]Todo
    next  int
}

func NewTodoRepo() *TodoRepo {
    return &TodoRepo{items: make(map[int]Todo), next: 1}
}

func (r *TodoRepo) ListByOwner(owner string) []Todo {
    r.mu.Lock()
    defer r.mu.Unlock()
    var out []Todo
    for _, t := range r.items {
        if t.Owner == owner {
            out = append(out, t)
        }
    }
    return out
}

func (r *TodoRepo) Add(owner, text string) Todo {
    r.mu.Lock()
    defer r.mu.Unlock()
    id := r.next
    r.next++
    t := Todo{ID: id, Owner: owner, Text: text}
    r.items[id] = t
    return t
}

func (r *TodoRepo) Delete(owner string, id int) bool {
    r.mu.Lock()
    defer r.mu.Unlock()
    t, ok := r.items[id]
    if !ok || t.Owner != owner {
        return false
    }
    delete(r.items, id)
    return true
}
