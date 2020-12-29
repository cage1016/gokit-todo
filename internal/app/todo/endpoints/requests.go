package endpoints

import (
	"github.com/cage1016/todo/internal/app/todo/model"
	"github.com/cage1016/todo/internal/app/todo/service"
)

type Request interface {
	validate() error
}

// AddRequest collects the request parameters for the Add method.
type AddRequest struct {
	Todo *model.TodoReq `json:"todo"`
}

func (r AddRequest) validate() error {
	return nil // TBA
}

// DeleteRequest collects the request parameters for the Delete method.
type DeleteRequest struct {
	Id string `json:"id"`
}

func (r DeleteRequest) validate() error {
	if r.Id == "" {
		return service.ErrMalformedEntity
	}
	return nil
}

// UpdateRequest collects the request parameters for the Update method.
type UpdateRequest struct {
	Id   string         `json:"id"`
	Todo *model.TodoReq `json:"todo"`
}

func (r UpdateRequest) validate() error {
	return nil
}

// ListRequest collects the request parameters for the List method.
type ListRequest struct {
}

func (r ListRequest) validate() error {
	return nil
}
