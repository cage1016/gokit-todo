package transports

import (
	pb "github.com/cage1016/todo/pb/todo"
	"time"

	"github.com/cage1016/todo/internal/app/todo/model"
)

func ModelReqToPB(todo *model.TodoReq) *pb.ModelTodoReq {
	return &pb.ModelTodoReq{
		Text:      *todo.Text,
		Completed: *todo.Completed,
	}
}

func PBtoModelReq(todo *pb.ModelTodoReq) *model.TodoReq {
	return &model.TodoReq{
		Text:      &todo.Text,
		Completed: &todo.Completed,
	}
}

func ModelResToPB(todo *model.TodoRes) *pb.ModelTodoRes {
	return &pb.ModelTodoRes{
		Id:        todo.ID,
		CreatedAt: todo.CreatedAt.Format(time.RFC3339),
		UpdatedAt: todo.UpdatedAt.Format(time.RFC3339),
		Text:      todo.Text,
		Completed: todo.Completed,
	}
}

func PBtoModelRes(todo *pb.ModelTodoRes) *model.TodoRes {
	return &model.TodoRes{
		ID:        todo.Id,
		Completed: todo.Completed,
		Text:      todo.Text,
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
	}
}
