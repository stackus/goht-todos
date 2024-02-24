package domain

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	ID          string
	Description string
	Completed   bool
	CreatedAt   time.Time
}

func NewTodo(description string) Todo {
	return Todo{
		ID:          NewTodoID(),
		Description: description,
		CreatedAt:   time.Now(),
	}
}

func NewTodoID() string {
	return uuid.New().String()
}
