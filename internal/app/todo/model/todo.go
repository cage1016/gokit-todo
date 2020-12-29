package model

import (
	"context"
	"encoding/json"
	"time"
)

type Todo struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Text      string    `json:"text"`
	Completed bool      `json:"completed"`
}

func (p Todo) MarshalJSON() ([]byte, error) {
	type Alias Todo
	return json.Marshal(&struct {
		Alias
		UpdatedAt string `json:"updatedAt"`
		CreatedAt string `json:"createdAt"`
	}{
		Alias:     (Alias)(p),
		UpdatedAt: p.UpdatedAt.Format(time.RFC3339),
		CreatedAt: p.CreatedAt.Format(time.RFC3339),
	})
}

func (t *Todo) UnmarshalJSON(data []byte) error {
	type Alias Todo

	pr := &struct {
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &pr); err != nil {
		return nil
	}

	t.UpdatedAt = time.Time{}
	t.CreatedAt = time.Time{}

	return nil
}

//go:generate mockgen -destination ../../../../internal/mocks/app/todo/model/todo.go -package=automocks . TodoRepository
type TodoRepository interface {
	Add(context.Context, *Todo) error
	Delete(context.Context, string) error
	Update(context.Context, *Todo) error
	List(context.Context) (res []*Todo, err error)
	Get(context.Context, string) (res *Todo, err error)
}

type TodoReq struct {
	Text      *string `json:"text"`
	Completed *bool   `json:"completed"`
}

type TodoRes Todo
