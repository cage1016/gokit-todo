package endpoints

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/cage1016/todo/internal/app/todo/model"
	"github.com/cage1016/todo/internal/app/todo/service"
	"github.com/cage1016/todo/internal/pkg/responses"
)

var (
	_ httptransport.Headerer = (*ListResponse)(nil)

	_ httptransport.StatusCoder = (*ListResponse)(nil)

	_ httptransport.Headerer = (*GetResponse)(nil)

	_ httptransport.StatusCoder = (*GetResponse)(nil)

	_ httptransport.Headerer = (*PostResponse)(nil)

	_ httptransport.StatusCoder = (*PostResponse)(nil)

	_ httptransport.Headerer = (*CompleteResponse)(nil)

	_ httptransport.StatusCoder = (*CompleteResponse)(nil)

	_ httptransport.Headerer = (*ClearCompleteResponse)(nil)

	_ httptransport.StatusCoder = (*ClearCompleteResponse)(nil)
)

// ListResponse collects the response values for the List method.
type ListResponse struct {
	Res []model.Todo `json:"res"`
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

// GetResponse collects the response values for the Get method.
type GetResponse struct {
	Res model.Todo `json:"res"`
	Err error      `json:"-"`
}

func (r GetResponse) StatusCode() int {
	return http.StatusOK // TBA
}

func (r GetResponse) Headers() http.Header {
	return http.Header{}
}

func (r GetResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version, Data: r.Res}
}

// PostResponse collects the response values for the Post method.
type PostResponse struct {
	Res model.Todo `json:"res"`
	Err error      `json:"-"`
}

func (r PostResponse) StatusCode() int {
	return http.StatusOK // TBA
}

func (r PostResponse) Headers() http.Header {
	return http.Header{}
}

func (r PostResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version, Data: r.Res}
}

// CompleteResponse collects the response values for the Complete method.
type CompleteResponse struct {
	Err error `json:"-"`
}

func (r CompleteResponse) StatusCode() int {
	return http.StatusOK // TBA
}

func (r CompleteResponse) Headers() http.Header {
	return http.Header{}
}

func (r CompleteResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version}
}

// ClearCompleteResponse collects the response values for the ClearComplete method.
type ClearCompleteResponse struct {
	Err error `json:"-"`
}

func (r ClearCompleteResponse) StatusCode() int {
	return http.StatusNoContent // TBA
}

func (r ClearCompleteResponse) Headers() http.Header {
	return http.Header{}
}

func (r ClearCompleteResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version}
}
