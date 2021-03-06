package service

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/cage1016/gokit-todo/internal/app/todo/model"
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

func (lm loggingMiddleware) Add(ctx context.Context, todo *model.TodoReq) (res *model.TodoRes, err error) {
	defer func() {
		lm.logger.Log("method", "Add", "todo", fmt.Sprintf("%v", todo), "err", err)
	}()

	return lm.next.Add(ctx, todo)
}

func (lm loggingMiddleware) Delete(ctx context.Context, id string) (err error) {
	defer func() {
		lm.logger.Log("method", "Delete", "id", id, "err", err)
	}()

	return lm.next.Delete(ctx, id)
}

func (lm loggingMiddleware) Update(ctx context.Context, id string, todo *model.TodoReq) (res *model.TodoRes, err error) {
	defer func() {
		lm.logger.Log("method", "Update", "id", id, "todo", fmt.Sprintf("%v", todo), "err", err)
	}()

	return lm.next.Update(ctx, id, todo)
}

func (lm loggingMiddleware) List(ctx context.Context) (res []*model.TodoRes, err error) {
	defer func() {
		lm.logger.Log("method", "List", "err", err)
	}()

	return lm.next.List(ctx)
}
