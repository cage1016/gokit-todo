package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"

	"github.com/cage1016/todo/internal/app/todo/model"
	"github.com/cage1016/todo/internal/app/todo/service"
)

// Endpoints collects all of the endpoints that compose the todo service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	AddEndpoint    endpoint.Endpoint `json:""`
	DeleteEndpoint endpoint.Endpoint `json:""`
	UpdateEndpoint endpoint.Endpoint `json:""`
	ListEndpoint   endpoint.Endpoint `json:""`
}

// New return a new instance of the endpoint that wraps the provided service.
func New(svc service.TodoService, logger log.Logger, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer) (ep Endpoints) {
	var addEndpoint endpoint.Endpoint
	{
		method := "add"
		addEndpoint = MakeAddEndpoint(svc)
		addEndpoint = opentracing.TraceServer(otTracer, method)(addEndpoint)
		addEndpoint = zipkin.TraceEndpoint(zipkinTracer, method)(addEndpoint)
		addEndpoint = LoggingMiddleware(log.With(logger, "method", method))(addEndpoint)
		ep.AddEndpoint = addEndpoint
	}

	var deleteEndpoint endpoint.Endpoint
	{
		method := "delete"
		deleteEndpoint = MakeDeleteEndpoint(svc)
		deleteEndpoint = opentracing.TraceServer(otTracer, method)(deleteEndpoint)
		deleteEndpoint = zipkin.TraceEndpoint(zipkinTracer, method)(deleteEndpoint)
		deleteEndpoint = LoggingMiddleware(log.With(logger, "method", method))(deleteEndpoint)
		ep.DeleteEndpoint = deleteEndpoint
	}

	var updateEndpoint endpoint.Endpoint
	{
		method := "update"
		updateEndpoint = MakeUpdateEndpoint(svc)
		updateEndpoint = opentracing.TraceServer(otTracer, method)(updateEndpoint)
		updateEndpoint = zipkin.TraceEndpoint(zipkinTracer, method)(updateEndpoint)
		updateEndpoint = LoggingMiddleware(log.With(logger, "method", method))(updateEndpoint)
		ep.UpdateEndpoint = updateEndpoint
	}

	var listEndpoint endpoint.Endpoint
	{
		method := "list"
		listEndpoint = MakeListEndpoint(svc)
		listEndpoint = opentracing.TraceServer(otTracer, method)(listEndpoint)
		listEndpoint = zipkin.TraceEndpoint(zipkinTracer, method)(listEndpoint)
		listEndpoint = LoggingMiddleware(log.With(logger, "method", method))(listEndpoint)
		ep.ListEndpoint = listEndpoint
	}

	return ep
}

// MakeAddEndpoint returns an endpoint that invokes Add on the service.
// Primarily useful in a server.
func MakeAddEndpoint(svc service.TodoService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddRequest)
		if err := req.validate(); err != nil {
			return AddResponse{}, err
		}
		res, err := svc.Add(ctx, req.Todo)
		return AddResponse{Res: res}, err
	}
}

// Add implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) Add(ctx context.Context, todo *model.TodoReq) (res *model.TodoRes, err error) {
	resp, err := e.AddEndpoint(ctx, AddRequest{Todo: todo})
	if err != nil {
		return
	}
	response := resp.(AddResponse)
	return response.Res, nil
}

// MakeDeleteEndpoint returns an endpoint that invokes Delete on the service.
// Primarily useful in a server.
func MakeDeleteEndpoint(svc service.TodoService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteRequest)
		if err := req.validate(); err != nil {
			return DeleteResponse{}, err
		}
		err := svc.Delete(ctx, req.Id)
		return DeleteResponse{}, err
	}
}

// Delete implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) Delete(ctx context.Context, id string) (err error) {
	resp, err := e.DeleteEndpoint(ctx, DeleteRequest{Id: id})
	if err != nil {
		return
	}
	_ = resp.(DeleteResponse)
	return nil
}

// MakeUpdateEndpoint returns an endpoint that invokes Update on the service.
// Primarily useful in a server.
func MakeUpdateEndpoint(svc service.TodoService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateRequest)
		if err := req.validate(); err != nil {
			return UpdateResponse{}, err
		}
		res, err := svc.Update(ctx, req.Id, req.Todo)
		return UpdateResponse{Res: res}, err
	}
}

// Update implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) Update(ctx context.Context, id string, todo *model.TodoReq) (res *model.TodoRes, err error) {
	resp, err := e.UpdateEndpoint(ctx, UpdateRequest{Id: id, Todo: todo})
	if err != nil {
		return
	}
	response := resp.(UpdateResponse)
	return response.Res, nil
}

// MakeListEndpoint returns an endpoint that invokes List on the service.
// Primarily useful in a server.
func MakeListEndpoint(svc service.TodoService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ListRequest)
		if err := req.validate(); err != nil {
			return ListResponse{}, err
		}
		res, err := svc.List(ctx)
		return ListResponse{Res: res}, err
	}
}

// List implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) List(ctx context.Context) (res []*model.TodoRes, err error) {
	resp, err := e.ListEndpoint(ctx, ListRequest{})
	if err != nil {
		return
	}
	response := resp.(ListResponse)
	return response.Res, nil
}
