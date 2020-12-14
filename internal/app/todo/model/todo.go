package model

import (
	"encoding/json"
	"time"
)

type Todo struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Text      string    `json:"text"`
	Complete  bool      `json:"complete"`
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
