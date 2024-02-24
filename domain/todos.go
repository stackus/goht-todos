package domain

import (
	"context"
)

type TodosStore interface {
	List(ctx context.Context) ([]Todo, error)
	Get(ctx context.Context, id string) (Todo, error)
	Create(ctx context.Context, todo Todo) (Todo, error)
	Update(ctx context.Context, id string, todo Todo) (Todo, error)
	Delete(ctx context.Context, id string) error
	Filter(ctx context.Context, filter string) ([]Todo, error)
	Reorder(ctx context.Context, todoIDs []string) error
}
