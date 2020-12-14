package service

import (
	"context"

	"github.com/cage1016/todo/internal/app/todo/model"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
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

func (lm loggingMiddleware) List(ctx context.Context, filter string) (res []model.Todo, err error) {
	defer func() {
		lm.logger.Log("method", "List", "filter", filter, "err", err)
	}()

	return lm.next.List(ctx, filter)
}

func (lm loggingMiddleware) Get(ctx context.Context, id string) (res model.Todo, err error) {
	defer func() {
		lm.logger.Log("method", "Get", "id", id, "err", err)
	}()

	return lm.next.Get(ctx, id)
}

func (lm loggingMiddleware) Post(ctx context.Context, todo model.Todo) (res model.Todo, err error) {
	defer func() {
		lm.logger.Log("method", "Post", "todo", todo, "err", err)
	}()

	return lm.next.Post(ctx, todo)
}

func (lm loggingMiddleware) Complete(ctx context.Context, id string) (err error) {
	defer func() {
		lm.logger.Log("method", "Complete", "id", id, "err", err)
	}()

	return lm.next.Complete(ctx, id)
}

func (lm loggingMiddleware) ClearComplete(ctx context.Context) (err error) {
	defer func() {
		lm.logger.Log("method", "ClearComplete", "err", err)
	}()

	return lm.next.ClearComplete(ctx)
}
