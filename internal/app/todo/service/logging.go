package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/cage1016/todo/internal/app/todo/model"
)

type loggingMiddleware struct {
	logger log.Logger  `json:""`
	next   TodoService `json:""`
}

// LoggingMiddleware takes a logger as a dependency
// and returns a ServiceMiddleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next TodoService) TodoService {
		return loggingMiddleware{level.Info(logger), next}
	}
}

func (lm loggingMiddleware) Add(ctx context.Context, todo *model.Todo) (res *model.Todo, err error) {
	defer func() {
		lm.logger.Log("method", "Add", "todo", todo, "err", err)
	}()

	return lm.next.Add(ctx, todo)
}

func (lm loggingMiddleware) Delete(ctx context.Context, id string) (err error) {
	defer func() {
		lm.logger.Log("method", "Delete", "id", id, "err", err)
	}()

	return lm.next.Delete(ctx, id)
}

func (lm loggingMiddleware) Update(ctx context.Context, id string, todo *model.Todo) (res *model.Todo, err error) {
	defer func() {
		lm.logger.Log("method", "Update", "id", id, "todo", todo, "err", err)
	}()

	return lm.next.Update(ctx, id, todo)
}

func (lm loggingMiddleware) List(ctx context.Context, filter string) (res []*model.Todo, err error) {
	defer func() {
		lm.logger.Log("method", "List", "filter", filter, "err", err)
	}()

	return lm.next.List(ctx, filter)
}

func (lm loggingMiddleware) Complete(ctx context.Context, id string) (err error) {
	defer func() {
		lm.logger.Log("method", "Complete", "id", id, "err", err)
	}()

	return lm.next.Complete(ctx, id)
}

func (lm loggingMiddleware) CompleteAll(ctx context.Context) (err error) {
	defer func() {
		lm.logger.Log("method", "CompleteAll", "err", err)
	}()

	return lm.next.CompleteAll(ctx)
}

func (lm loggingMiddleware) Clear(ctx context.Context) (err error) {
	defer func() {
		lm.logger.Log("method", "Clear", "err", err)
	}()

	return lm.next.Clear(ctx)
}
