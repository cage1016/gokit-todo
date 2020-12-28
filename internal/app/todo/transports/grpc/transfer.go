package transports

import (
	"time"

	pb "github.com/cage1016/todo/pb/todo"

	"github.com/cage1016/todo/internal/app/todo/model"
)

func ModelToPB(todo *model.Todo) *pb.ModelTodo {
	return &pb.ModelTodo{
		Id:        todo.ID,
		CreatedAt: todo.CreatedAt.Format(time.RFC3339),
		UpdatedAt: todo.UpdatedAt.Format(time.RFC3339),
		Text:      todo.Text,
		Complete:  todo.Complete,
	}
}

func PBtoModel(todo *pb.ModelTodo) *model.Todo {
	return &model.Todo{
		ID: todo.Id,
		CreatedAt: func() time.Time {
			t, err := time.Parse(time.RFC3339, todo.CreatedAt)
			if err != nil {
				return time.Time{}
			}
			return t
		}(),
		UpdatedAt: func() time.Time {
			t, err := time.Parse(time.RFC3339, todo.UpdatedAt)
			if err != nil {
				return time.Time{}
			}
			return t
		}(),
		Text:     todo.Text,
		Complete: todo.Complete,
	}
}
