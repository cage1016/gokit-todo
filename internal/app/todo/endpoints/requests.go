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
	Todo *model.Todo `json:"todo"`
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
	Id   string      `json:"id"`
	Todo *model.Todo `json:"todo"`
}

func (r UpdateRequest) validate() error {
	return nil
}

// ListRequest collects the request parameters for the List method.
type ListRequest struct {
	Filter string `json:"filter"`
}

func (r ListRequest) validate() error {
	if r.Filter != service.ALL && r.Filter != service.ACTIVE && r.Filter != service.COMPLETE {
		return service.ErrInvalidQueryParams
	}
	return nil
}

// CompleteRequest collects the request parameters for the Completed method.
type CompleteRequest struct {
	Id string `json:"id"`
}

func (r CompleteRequest) validate() error {
	if r.Id == "" {
		return service.ErrMalformedEntity
	}
	return nil
}

// CompleteAllRequest collects the request parameters for the CompleteAll method.
type CompleteAllRequest struct {
}

func (r CompleteAllRequest) validate() error {
	return nil // TBA
}

// ClearRequest collects the request parameters for the Clear method.
type ClearRequest struct {
}

func (r ClearRequest) validate() error {
	return nil // TBA
}
