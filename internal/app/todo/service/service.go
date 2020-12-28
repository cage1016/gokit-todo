package service

import (
	"context"
	"time"

	gonanoid "github.com/matoous/go-nanoid"

	"github.com/cage1016/todo/internal/app/todo/model"
	"github.com/cage1016/todo/internal/pkg/errors"
	"github.com/go-kit/kit/log"
)

var (
	ErrMalformedEntity = errors.New("malformed entity specification")

	ErrNotFound = errors.New("non-existent entity")

	ErrInvalidQueryParams = errors.New("invalid query params")
)

const (
	ALL      = "all"
	ACTIVE   = "active"
	COMPLETE = "complete"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(TodoService) TodoService

// Service describes a service that adds things together
// Implement yor service methods methods.
// e.x: Foo(ctx context.Context, s string)(rs string, err error)
//go:generate mockgen -destination ../../../../internal/mocks/app/todo/service/todoservice.go -package=automocks . TodoService
type TodoService interface {
	// [method=post,expose=true,router=items]
	Add(ctx context.Context, todo *model.Todo) (res *model.Todo, err error)
	// [method=delete,expose=true,router=items/:id]
	Delete(ctx context.Context, id string) (err error)
	// [method=put,expose=true,router=items/:id]
	Update(ctx context.Context, id string, todo *model.Todo) (res *model.Todo, err error)
	// [method=get,expose=true,router=items]
	List(ctx context.Context, filter string) (res []*model.Todo, err error)
	// [method=post,expose=true,router=items/completes/:id]
	Complete(ctx context.Context, id string) (err error)
	// [method=post,expose=true,router=items/completes]
	CompleteAll(ctx context.Context) (err error)
	// [method=post,expose=true,router=items/clear]
	Clear(ctx context.Context) (err error)
}

// the concrete implementation of service interface
type stubTodoService struct {
	repo   model.TodoRepository
	logger log.Logger
}

// New return a new instance of the service.
// If you want to add service middleware this is the place to put them.
func New(repo model.TodoRepository, logger log.Logger) (s TodoService) {
	var svc TodoService
	{
		svc = &stubTodoService{repo: repo, logger: logger}
		svc = LoggingMiddleware(logger)(svc)
	}
	return svc
}

// Implement the business logic of Add
func (to *stubTodoService) Add(ctx context.Context, todo *model.Todo) (res *model.Todo, err error) {
	id, _ := gonanoid.ID(21)
	todo.ID = id
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()

	if err := to.repo.Add(ctx, todo); err != nil {
		return res, err
	}
	return todo, nil
}

// Implement the business logic of Delete
func (to *stubTodoService) Delete(ctx context.Context, id string) (err error) {
	return to.repo.Delete(ctx, id)
}

// Implement the business logic of Update
func (to *stubTodoService) Update(ctx context.Context, id string, todo *model.Todo) (res* model.Todo, err error) {
	todo.UpdatedAt = time.Now()
	if err := to.repo.Update(ctx, todo); err != nil {
		return todo, err
	}
	return todo, nil
}

// Implement the business logic of List
func (to *stubTodoService) List(ctx context.Context, filter string) (res []*model.Todo, err error) {
	return to.repo.List(ctx, filter)
}

// Implement the business logic of Complete
func (to *stubTodoService) Complete(ctx context.Context, id string) (err error) {
	return to.repo.Complete(ctx, id)
}

// Implement the business logic of CompleteAll
func (to *stubTodoService) CompleteAll(ctx context.Context) (err error) {
	return to.repo.CompleteAll(ctx)
}

// Implement the business logic of Clear
func (to *stubTodoService) Clear(ctx context.Context) (err error) {
	return to.repo.Clear(ctx)
}
