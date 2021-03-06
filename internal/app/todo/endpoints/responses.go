package endpoints

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/cage1016/gokit-todo/internal/app/todo/model"
	"github.com/cage1016/gokit-todo/internal/app/todo/service"
	"github.com/cage1016/gokit-todo/internal/pkg/responses"
)

var (
	_ httptransport.Headerer = (*AddResponse)(nil)

	_ httptransport.StatusCoder = (*AddResponse)(nil)

	_ httptransport.Headerer = (*DeleteResponse)(nil)

	_ httptransport.StatusCoder = (*DeleteResponse)(nil)

	_ httptransport.Headerer = (*UpdateResponse)(nil)

	_ httptransport.StatusCoder = (*UpdateResponse)(nil)

	_ httptransport.Headerer = (*ListResponse)(nil)

	_ httptransport.StatusCoder = (*ListResponse)(nil)
)

// AddResponse collects the response values for the Add method.
type AddResponse struct {
	Res *model.TodoRes `json:"res"`
	Err error          `json:"-"`
}

func (r AddResponse) StatusCode() int {
	return http.StatusCreated // TBA
}

func (r AddResponse) Headers() http.Header {
	return http.Header{}
}

func (r AddResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version, Data: r.Res}
}

// DeleteResponse collects the response values for the Delete method.
type DeleteResponse struct {
	Err error `json:"-"`
}

func (r DeleteResponse) StatusCode() int {
	return http.StatusNoContent // TBA
}

func (r DeleteResponse) Headers() http.Header {
	return http.Header{}
}

func (r DeleteResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version}
}

// UpdateResponse collects the response values for the Update method.
type UpdateResponse struct {
	Res *model.TodoRes `json:"res"`
	Err error          `json:"-"`
}

func (r UpdateResponse) StatusCode() int {
	return http.StatusOK // TBA
}

func (r UpdateResponse) Headers() http.Header {
	return http.Header{}
}

func (r UpdateResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version, Data: r.Res}
}

// ListResponse collects the response values for the List method.
type ListResponse struct {
	Res []*model.TodoRes `json:"res"`
	Err error            `json:"-"`
}

func (r ListResponse) StatusCode() int {
	return http.StatusOK // TBA
}

func (r ListResponse) Headers() http.Header {
	return http.Header{}
}

func (r ListResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version, Data: r.Res}
}

// CompleteAllResponse collects the response values for the CompleteAll method.
type CompleteAllResponse struct {
	Err error `json:"-"`
}
