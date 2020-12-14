package endpoints

import (
	"context"

	"github.com/cage1016/todo/internal/app/todo/model"
	"github.com/cage1016/todo/internal/app/todo/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
)

// Endpoints collects all of the endpoints that compose the todo service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	ListEndpoint          endpoint.Endpoint `json:""`
	GetEndpoint           endpoint.Endpoint `json:""`
	PostEndpoint          endpoint.Endpoint `json:""`
	CompleteEndpoint      endpoint.Endpoint `json:""`
	ClearCompleteEndpoint endpoint.Endpoint `json:""`
}

// New return a new instance of the endpoint that wraps the provided service.
func New(svc service.TodoService, logger log.Logger, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer) (ep Endpoints) {
	var listEndpoint endpoint.Endpoint
	{
		method := "list"
		listEndpoint = MakeListEndpoint(svc)
		listEndpoint = opentracing.TraceServer(otTracer, method)(listEndpoint)
		listEndpoint = zipkin.TraceEndpoint(zipkinTracer, method)(listEndpoint)
		listEndpoint = LoggingMiddleware(log.With(logger, "method", method))(listEndpoint)
		ep.ListEndpoint = listEndpoint
	}

	var getEndpoint endpoint.Endpoint
	{
		method := "get"
		getEndpoint = MakeGetEndpoint(svc)
		getEndpoint = opentracing.TraceServer(otTracer, method)(getEndpoint)
		getEndpoint = zipkin.TraceEndpoint(zipkinTracer, method)(getEndpoint)
		getEndpoint = LoggingMiddleware(log.With(logger, "method", method))(getEndpoint)
		ep.GetEndpoint = getEndpoint
	}

	var postEndpoint endpoint.Endpoint
	{
		method := "post"
		postEndpoint = MakePostEndpoint(svc)
		postEndpoint = opentracing.TraceServer(otTracer, method)(postEndpoint)
		postEndpoint = zipkin.TraceEndpoint(zipkinTracer, method)(postEndpoint)
		postEndpoint = LoggingMiddleware(log.With(logger, "method", method))(postEndpoint)
		ep.PostEndpoint = postEndpoint
	}

	var completeEndpoint endpoint.Endpoint
	{
		method := "complete"
		completeEndpoint = MakeCompleteEndpoint(svc)
		completeEndpoint = opentracing.TraceServer(otTracer, method)(completeEndpoint)
		completeEndpoint = zipkin.TraceEndpoint(zipkinTracer, method)(completeEndpoint)
		completeEndpoint = LoggingMiddleware(log.With(logger, "method", method))(completeEndpoint)
		ep.CompleteEndpoint = completeEndpoint
	}

	var clearCompleteEndpoint endpoint.Endpoint
	{
		method := "clearComplete"
		clearCompleteEndpoint = MakeClearCompleteEndpoint(svc)
		clearCompleteEndpoint = opentracing.TraceServer(otTracer, method)(clearCompleteEndpoint)
		clearCompleteEndpoint = zipkin.TraceEndpoint(zipkinTracer, method)(clearCompleteEndpoint)
		clearCompleteEndpoint = LoggingMiddleware(log.With(logger, "method", method))(clearCompleteEndpoint)
		ep.ClearCompleteEndpoint = clearCompleteEndpoint
	}

	return ep
}

// MakeListEndpoint returns an endpoint that invokes List on the service.
// Primarily useful in a server.
func MakeListEndpoint(svc service.TodoService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ListRequest)
		if err := req.validate(); err != nil {
			return GetResponse{}, err
		}
		res, err := svc.List(ctx, req.Filter)
		return ListResponse{Res: res}, err
	}
}

// List implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) List(ctx context.Context, filter string) (res []model.Todo, err error) {
	resp, err := e.ListEndpoint(ctx, ListRequest{Filter: filter})
	if err != nil {
		return
	}
	response := resp.(ListResponse)
	return response.Res, nil
}

// MakeGetEndpoint returns an endpoint that invokes Get on the service.
// Primarily useful in a server.
func MakeGetEndpoint(svc service.TodoService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetRequest)
		if err := req.validate(); err != nil {
			return GetResponse{}, err
		}
		res, err := svc.Get(ctx, req.Id)
		return GetResponse{Res: res}, err
	}
}

// Get implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) Get(ctx context.Context, id string) (res model.Todo, err error) {
	resp, err := e.GetEndpoint(ctx, GetRequest{Id: id})
	if err != nil {
		return
	}
	response := resp.(GetResponse)
	return response.Res, nil
}

// MakePostEndpoint returns an endpoint that invokes Post on the service.
// Primarily useful in a server.
func MakePostEndpoint(svc service.TodoService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PostRequest)
		if err := req.validate(); err != nil {
			return PostResponse{}, err
		}
		res, err := svc.Post(ctx, req.Todo)
		return PostResponse{Res: res}, err
	}
}

// Post implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) Post(ctx context.Context, todo model.Todo) (res model.Todo, err error) {
	resp, err := e.PostEndpoint(ctx, PostRequest{Todo: todo})
	if err != nil {
		return
	}
	response := resp.(PostResponse)
	return response.Res, nil
}

// MakeCompleteEndpoint returns an endpoint that invokes Complete on the service.
// Primarily useful in a server.
func MakeCompleteEndpoint(svc service.TodoService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CompleteRequest)
		if err := req.validate(); err != nil {
			return CompleteResponse{}, err
		}
		err := svc.Complete(ctx, req.Id)
		return CompleteResponse{}, err
	}
}

// Complete implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) Complete(ctx context.Context, id string) (err error) {
	resp, err := e.CompleteEndpoint(ctx, CompleteRequest{Id: id})
	if err != nil {
		return
	}
	_ = resp.(CompleteResponse)
	return nil
}

// MakeClearCompleteEndpoint returns an endpoint that invokes ClearComplete on the service.
// Primarily useful in a server.
func MakeClearCompleteEndpoint(svc service.TodoService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		err := svc.ClearComplete(ctx)
		return ClearCompleteResponse{}, err
	}
}

// ClearComplete implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) ClearComplete(ctx context.Context) (err error) {
	resp, err := e.ClearCompleteEndpoint(ctx, ClearCompleteRequest{})
	if err != nil {
		return
	}
	_ = resp.(ClearCompleteResponse)
	return nil
}
