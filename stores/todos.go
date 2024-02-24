package stores

import (
	"context"
	"math"
	"slices"
	"strings"
	"sync"

	"github.com/stackus/errors"

	"github.com/stackus/goht-todos/domain"
)

type todosStore struct {
	mu    sync.RWMutex
	todos []domain.Todo
}

func NewTodosStore() domain.TodosStore {
	return &todosStore{
		todos: []domain.Todo{},
	}
}

func (s *todosStore) List(_ context.Context) ([]domain.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.todos, nil
}

func (s *todosStore) Get(_ context.Context, id string) (domain.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, todo := range s.todos {
		if todo.ID == id {
			return todo, nil
		}
	}
	return domain.Todo{}, errors.ErrNotFound.Msgf("todo with id %s not found", id)
}

func (s *todosStore) Create(_ context.Context, todo domain.Todo) (domain.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.todos = append(s.todos, todo)
	return todo, nil
}

func (s *todosStore) Update(_ context.Context, id string, todo domain.Todo) (domain.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, t := range s.todos {
		if t.ID == id {
			s.todos[i].Completed = todo.Completed
			s.todos[i].Description = todo.Description
			return s.todos[i], nil
		}
	}
	return domain.Todo{}, errors.ErrNotFound.Msgf("todo with id %s not found", id)
}

func (s *todosStore) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, todo := range s.todos {
		if todo.ID == id {
			s.todos = append(s.todos[:i], s.todos[i+1:]...)
			return nil
		}
	}
	return errors.ErrNotFound.Msgf("todo with id %s not found", id)
}

func (s *todosStore) Filter(_ context.Context, filter string) ([]domain.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var list []domain.Todo
	for _, todo := range s.todos {
		if strings.Contains(todo.Description, filter) {
			list = append(list, todo)
		}
	}
	return list, nil
}

func (s *todosStore) Reorder(_ context.Context, todoIDs []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// get original and new positions; from the _possibly_ filtered list of todos
	origOrder := make([]int, len(todoIDs))
	posTrans := make([]int, len(todoIDs))
	index := 0
	for i, todo := range s.todos {
		if newPos := slices.Index(todoIDs, todo.ID); newPos != -1 {
			origOrder[newPos] = i
			posTrans[newPos] = newPos - index
			index++
		}
	}

	// the todo with the greatest distance to move was the one the user wanted to move
	moveIndex := 0
	for i := 1; i < len(posTrans); i++ {
		if math.Abs(float64(posTrans[i])) > math.Abs(float64(posTrans[moveIndex])) {
			moveIndex = i
		}
	}

	oldIndex := origOrder[moveIndex]
	newIndex := 0
	switch posTrans[moveIndex] > 0 {
	case true: // move todo to the right
		newIndex = origOrder[moveIndex-1]
	case false: // move todo to the left
		newIndex = origOrder[moveIndex+1]
	}

	// move the todo to the new position from the old position
	if newIndex < 0 {
		newIndex = 0
	}
	if newIndex >= len(s.todos) {
		newIndex = len(s.todos) - 1
	}
	todo := s.todos[oldIndex]
	s.todos = append(s.todos[:oldIndex], s.todos[oldIndex+1:]...)
	s.todos = append(s.todos[:newIndex], append([]domain.Todo{todo}, s.todos[newIndex:]...)...)

	return nil
}
