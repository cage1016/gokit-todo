package endpoints

import (
	"github.com/cage1016/todo/internal/app/todo/model"
	"github.com/cage1016/todo/internal/app/todo/service"
)

type Request interface {
	validate() error
}

// ListRequest collects the request parameters for the List method.
type ListRequest struct {
}

func (r ListRequest) validate() error {
	return nil // TBA
}

// GetRequest collects the request parameters for the Get method.
type GetRequest struct {
	Id string `json:"id"`
}

func (r GetRequest) validate() error {
	if r.Id == "" {
		return service.ErrMalformedEntity
	}
	return nil
}

// PostRequest collects the request parameters for the Post method.
type PostRequest struct {
	Todo model.Todo `json:"todo"`
}

func (r PostRequest) validate() error {
	return nil // TBA
}

// CompleteRequest collects the request parameters for the Complete method.
type CompleteRequest struct {
	Id string `json:"id"`
}

func (r CompleteRequest) validate() error {
	if r.Id == "" {
		return service.ErrMalformedEntity
	}
	return nil
}

// ClearCompleteRequest collects the request parameters for the ClearComplete method.
type ClearCompleteRequest struct {
}

func (r ClearCompleteRequest) validate() error {
	return nil // TBA
}
