package service

import (
	"context"
	"time"

	gonanoid "github.com/matoous/go-nanoid"

	"github.com/cage1016/gokit-todo/internal/app/todo/model"
	"github.com/cage1016/gokit-todo/internal/pkg/errors"
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
	Add(ctx context.Context, todo *model.TodoReq) (res *model.TodoRes, err error)
	// [method=delete,expose=true,router=items/:id]
	Delete(ctx context.Context, id string) (err error)
	// [method=put,expose=true,router=items/:id]
	Update(ctx context.Context, id string, todo *model.TodoReq) (res *model.TodoRes, err error)
	// [method=get,expose=true,router=items]
	List(ctx context.Context) (res []*model.TodoRes, err error)
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
func (to *stubTodoService) Add(ctx context.Context, todo *model.TodoReq) (res *model.TodoRes, err error) {
	id, _ := gonanoid.ID(21)

	t := new(model.Todo)
	t.ID = id
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	if todo.Completed != nil {
		t.Completed = *todo.Completed
	}
	if todo.Text != nil {
		t.Text = *todo.Text
	}
	if err := to.repo.Add(ctx, t); err != nil {
		return res, err
	}
	x := model.TodoRes(*t)
	return &x, nil
}

// Implement the business logic of Delete
func (to *stubTodoService) Delete(ctx context.Context, id string) (err error) {
	return to.repo.Delete(ctx, id)
}

// Implement the business logic of Update
func (to *stubTodoService) Update(ctx context.Context, id string, todo *model.TodoReq) (res *model.TodoRes, err error) {
	dt, err := to.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	dt.UpdatedAt = time.Now()
	if todo.Completed != nil {
		dt.Completed = *todo.Completed
	}
	if todo.Text != nil {
		dt.Text = *todo.Text
	}

	if err := to.repo.Update(ctx, dt); err != nil {
		return nil, err
	}
	x := model.TodoRes(*dt)
	return &x, nil
}

// Implement the business logic of List
func (to *stubTodoService) List(ctx context.Context) (res []*model.TodoRes, err error) {
	res = make([]*model.TodoRes, 0)

	rr, err := to.repo.List(ctx)
	if err != nil {
		return
	}

	for _, r := range rr {
		item := model.TodoRes(*r)
		res = append(res, &item)
	}
	return
}
