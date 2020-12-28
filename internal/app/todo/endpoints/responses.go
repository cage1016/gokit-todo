package endpoints

import (
	"github.com/cage1016/todo/internal/app/todo/model"
	"net/http"
	
	httptransport "github.com/go-kit/kit/transport/http"
	
	"github.com/cage1016/todo/internal/pkg/responses"
	"github.com/cage1016/todo/internal/app/todo/service"
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

	_ httptransport.Headerer = (*CompleteResponse)(nil)

	_ httptransport.StatusCoder = (*CompleteResponse)(nil)

	_ httptransport.Headerer = (*CompleteAllResponse)(nil)

	_ httptransport.StatusCoder = (*CompleteAllResponse)(nil)

	_ httptransport.Headerer = (*ClearResponse)(nil)

	_ httptransport.StatusCoder = (*ClearResponse)(nil)
)

// AddResponse collects the response values for the Add method.
type AddResponse struct {
	Res *model.Todo `json:"res"`
	Err error      `json:"-"`
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
	Res *model.Todo `json:"res"`
	Err error      `json:"-"`
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
	Res []*model.Todo `json:"res"`
	Err error        `json:"-"`
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

// CompleteResponse collects the response values for the Complete method.
type CompleteResponse struct {
	Err error `json:"-"`
}

func (r CompleteResponse) StatusCode() int {
	return http.StatusNoContent // TBA
}

func (r CompleteResponse) Headers() http.Header {
	return http.Header{}
}

func (r CompleteResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version}
}

// CompleteAllResponse collects the response values for the CompleteAll method.
type CompleteAllResponse struct {
	Err error `json:"-"`
}

func (r CompleteAllResponse) StatusCode() int {
	return http.StatusNoContent // TBA
}

func (r CompleteAllResponse) Headers() http.Header {
	return http.Header{}
}

func (r CompleteAllResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version}
}

// ClearResponse collects the response values for the Clear method.
type ClearResponse struct {
	Err error `json:"-"`
}

func (r ClearResponse) StatusCode() int {
	return http.StatusNoContent // TBA
}

func (r ClearResponse) Headers() http.Header {
	return http.Header{}
}

func (r ClearResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version}
}

