package postgres

import "github.com/cage1016/todo/internal/app/todo/model"

type Todo = model.Todo

func ModelToDB(todo model.Todo) Todo {
	return Todo(todo)
}

func DBtoModel(todo Todo) model.Todo {
	return model.Todo(todo)
}
