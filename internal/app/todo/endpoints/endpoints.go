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
	AddEndpoint         endpoint.Endpoint `json:""`
	DeleteEndpoint      endpoint.Endpoint `json:""`
	UpdateEndpoint      endpoint.Endpoint `json:""`
	ListEndpoint        endpoint.Endpoint `json:""`
	CompleteEndpoint    endpoint.Endpoint `json:""`
	CompleteAllEndpoint endpoint.Endpoint `json:""`
	ClearEndpoint       endpoint.Endpoint `json:""`
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

	var completeEndpoint endpoint.Endpoint
	{
		method := "complete"
		completeEndpoint = MakeCompleteEndpoint(svc)
		completeEndpoint = opentracing.TraceServer(otTracer, method)(completeEndpoint)
		completeEndpoint = zipkin.TraceEndpoint(zipkinTracer, method)(completeEndpoint)
		completeEndpoint = LoggingMiddleware(log.With(logger, "method", method))(completeEndpoint)
		ep.CompleteEndpoint = completeEndpoint
	}

	var completeAllEndpoint endpoint.Endpoint
	{
		method := "completeAll"
		completeAllEndpoint = MakeCompleteAllEndpoint(svc)
		completeAllEndpoint = opentracing.TraceServer(otTracer, method)(completeAllEndpoint)
		completeAllEndpoint = zipkin.TraceEndpoint(zipkinTracer, method)(completeAllEndpoint)
		completeAllEndpoint = LoggingMiddleware(log.With(logger, "method", method))(completeAllEndpoint)
		ep.CompleteAllEndpoint = completeAllEndpoint
	}

	var clearEndpoint endpoint.Endpoint
	{
		method := "clear"
		clearEndpoint = MakeClearEndpoint(svc)
		clearEndpoint = opentracing.TraceServer(otTracer, method)(clearEndpoint)
		clearEndpoint = zipkin.TraceEndpoint(zipkinTracer, method)(clearEndpoint)
		clearEndpoint = LoggingMiddleware(log.With(logger, "method", method))(clearEndpoint)
		ep.ClearEndpoint = clearEndpoint
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
func (e Endpoints) Add(ctx context.Context, todo model.Todo) (res model.Todo, err error) {
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
func (e Endpoints) Update(ctx context.Context, id string, todo model.Todo) (res model.Todo, err error) {
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

// MakeCompleteAllEndpoint returns an endpoint that invokes CompleteAll on the service.
// Primarily useful in a server.
func MakeCompleteAllEndpoint(svc service.TodoService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		err := svc.CompleteAll(ctx)
		return CompleteAllResponse{}, err
	}
}

// CompleteAll implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) CompleteAll(ctx context.Context) (err error) {
	resp, err := e.CompleteAllEndpoint(ctx, CompleteAllRequest{})
	if err != nil {
		return
	}
	_ = resp.(CompleteAllResponse)
	return nil
}

// MakeClearEndpoint returns an endpoint that invokes Clear on the service.
// Primarily useful in a server.
func MakeClearEndpoint(svc service.TodoService) (ep endpoint.Endpoint) {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		err := svc.Clear(ctx)
		return ClearResponse{}, err
	}
}

// Clear implements the service interface, so Endpoints may be used as a service.
// This is primarily useful in the context of a client library.
func (e Endpoints) Clear(ctx context.Context) (err error) {
	resp, err := e.ClearEndpoint(ctx, ClearRequest{})
	if err != nil {
		return
	}
	_ = resp.(ClearResponse)
	return nil
}
