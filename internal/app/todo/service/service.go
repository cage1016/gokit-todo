package service

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	gonanoid "github.com/matoous/go-nanoid"

	"github.com/cage1016/todo/internal/app/todo/model"
	"github.com/cage1016/todo/internal/pkg/errors"
)

var (
	ErrMalformedEntity = errors.New("malformed entity specification")

	ErrNotFound = errors.New("non-existent entity")

	ErrInvalidQueryParams = errors.New("invalid query params")
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(TodoService) TodoService

// Service describes a service that adds things together
// Implement yor service methods methods.
// e.x: Foo(ctx context.Context, s string)(rs string, err error)
//go:generate mockgen -destination ../../../../internal/mocks/app/todo/service/todoservice.go -package=automocks . TodoService
type TodoService interface {
	// [method=get,expose=true,router=items]
	List(ctx context.Context, filter string) (res []model.Todo, err error)
	// [method=get,expose=true,router=items/[a-zA-Z0-9_-~]{21}]
	Get(ctx context.Context, id string) (res model.Todo, err error)
	// [method=post,expose=true,router=/items /]
	Post(ctx context.Context, todo model.Todo) (res model.Todo, err error)
	// [method=put,expose=true,router=/items/[a-zA-Z0-9_-~]{21}]
	Complete(ctx context.Context, id string) (err error)
	// [method=post,expose=true,router=/items/clear]
	ClearComplete(ctx context.Context) (err error)
}

// the concrete implementation of service interface
type stubTodoService struct {
	logger log.Logger
	repo   model.TodoRepository
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

// Implement the business logic of List
func (to *stubTodoService) List(ctx context.Context, filter string) (res []model.Todo, err error) {
	return to.repo.List(ctx, filter)
}

// Implement the business logic of Get
func (to *stubTodoService) Get(ctx context.Context, todoID string) (res model.Todo, err error) {
	return to.repo.Get(ctx, todoID)
}

// Implement the business logic of Post
func (to *stubTodoService) Post(ctx context.Context, todo model.Todo) (res model.Todo, err error) {
	id, _ := gonanoid.ID(21)

	todo.ID = id
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()

	if err := to.repo.Save(ctx, todo); err != nil {
		return res, err
	}
	return todo, nil
}

// Implement the business logic of Complete
func (to *stubTodoService) Complete(ctx context.Context, todoID string) (err error) {
	return to.repo.Complete(ctx, todoID)
}

// Implement the business logic of ClearComplete
func (to *stubTodoService) ClearComplete(ctx context.Context) (err error) {
	return to.repo.Clear(ctx)
}
