package transports

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/rs/cors"

	"github.com/cage1016/todo/internal/app/todo/endpoints"
	"github.com/cage1016/todo/internal/app/todo/service"
	"github.com/cage1016/todo/internal/pkg/errors"
	"github.com/cage1016/todo/internal/pkg/responses"
)

const (
	contentType string = "application/json"
)

// ShowTodo godoc
// @Summary List
// @Description TODO
// @Tags TODO
// @Accept json
// @Produce json
// @Router /items [get]
func ListHandler(m *bone.Mux, endpoints endpoints.Endpoints, options []httptransport.ServerOption, otTracer stdopentracing.Tracer, logger log.Logger) {
	m.Get("/items", httptransport.NewServer(
		endpoints.ListEndpoint,
		decodeHTTPListRequest,
		responses.EncodeJSONResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "List", logger), kitjwt.HTTPToContext()))...,
	))
}

// ShowTodo godoc
// @Summary Get
// @Description TODO
// @Tags TODO
// @Accept json
// @Produce json
// @Router /items/{id} [get]
func GetHandler(m *bone.Mux, endpoints endpoints.Endpoints, options []httptransport.ServerOption, otTracer stdopentracing.Tracer, logger log.Logger) {
	m.Get("/items/:id", httptransport.NewServer(
		endpoints.GetEndpoint,
		decodeHTTPGetRequest,
		responses.EncodeJSONResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Get", logger), kitjwt.HTTPToContext()))...,
	))
}

// ShowTodo godoc
// @Summary Post
// @Description TODO
// @Tags TODO
// @Accept json
// @Produce json
// @Router /items [post]
func PostHandler(m *bone.Mux, endpoints endpoints.Endpoints, options []httptransport.ServerOption, otTracer stdopentracing.Tracer, logger log.Logger) {
	m.Post("/items", httptransport.NewServer(
		endpoints.PostEndpoint,
		decodeHTTPPostRequest,
		responses.EncodeJSONResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Post", logger), kitjwt.HTTPToContext()))...,
	))
}

// ShowTodo godoc
// @Summary Complete
// @Description TODO
// @Tags TODO
// @Accept json
// @Produce json
// @Router /items/{id} [put]
func CompleteHandler(m *bone.Mux, endpoints endpoints.Endpoints, options []httptransport.ServerOption, otTracer stdopentracing.Tracer, logger log.Logger) {
	m.Put("/items/:id", httptransport.NewServer(
		endpoints.CompleteEndpoint,
		decodeHTTPCompleteRequest,
		responses.EncodeJSONResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Complete", logger), kitjwt.HTTPToContext()))...,
	))
}

// ShowTodo godoc
// @Summary ClearComplete
// @Description TODO
// @Tags TODO
// @Accept json
// @Produce json
// @Router /items/clear  [delete]
func ClearCompleteHandler(m *bone.Mux, endpoints endpoints.Endpoints, options []httptransport.ServerOption, otTracer stdopentracing.Tracer, logger log.Logger) {
	m.Post("/items/clear", httptransport.NewServer(
		endpoints.ClearCompleteEndpoint,
		decodeHTTPClearCompleteRequest,
		responses.EncodeJSONResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "ClearComplete", logger), kitjwt.HTTPToContext()))...,
	))
}

// NewHTTPHandler returns a handler that makes a set of endpoints available on
// predefined paths.
func NewHTTPHandler(endpoints endpoints.Endpoints, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) http.Handler { // Zipkin HTTP Server Trace can either be instantiated per endpoint with a
	// provided operation name or a global tracing service can be instantiated
	// without an operation name and fed to each Go kit endpoint as ServerOption.
	// In the latter case, the operation name will be the endpoint's http method.
	// We demonstrate a global tracing service here.
	zipkinServer := zipkin.HTTPServerTrace(zipkinTracer)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(responses.ErrorEncodeJSONResponse(CustomErrorEncoder)),
		httptransport.ServerErrorLogger(logger),
		zipkinServer,
	}

	m := bone.New()
	ListHandler(m, endpoints, options, otTracer, logger)
	GetHandler(m, endpoints, options, otTracer, logger)
	PostHandler(m, endpoints, options, otTracer, logger)
	CompleteHandler(m, endpoints, options, otTracer, logger)
	ClearCompleteHandler(m, endpoints, options, otTracer, logger)
	return cors.AllowAll().Handler(m)
}

// decodeHTTPListRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func decodeHTTPListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.ListRequest
	return req, nil
}

// decodeHTTPGetRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func decodeHTTPGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.GetRequest
	req.Id = bone.GetValue(r, "id")
	return req, nil
}

// decodeHTTPPostRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func decodeHTTPPostRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.PostRequest
	err := json.NewDecoder(r.Body).Decode(&req.Todo)
	return req, err
}

// decodeHTTPCompleteRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func decodeHTTPCompleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.CompleteRequest
	req.Id = bone.GetValue(r, "id")
	return req, nil
}

// decodeHTTPClearCompleteRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func decodeHTTPClearCompleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.ClearCompleteRequest
	return req, nil
}

// NewHTTPClient returns an AddService backed by an HTTP server living at the
// remote instance. We expect instance to come from a service discovery system,
// so likely of the form "host:port". We bake-in certain middlewares,
// implementing the client library pattern.
func NewHTTPClient(instance string, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) (service.TodoService, error) { // Quickly sanitize the instance string.
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}

	// Zipkin HTTP Client Trace can either be instantiated per endpoint with a
	// provided operation name or a global tracing client can be instantiated
	// without an operation name and fed to each Go kit endpoint as ClientOption.
	// In the latter case, the operation name will be the endpoint's http method.
	zipkinClient := zipkin.HTTPClientTrace(zipkinTracer)

	// global client middlewares
	options := []httptransport.ClientOption{
		zipkinClient,
	}

	e := endpoints.Endpoints{}

	// Each individual endpoint is an http/transport.Client (which implements
	// endpoint.Endpoint) that gets wrapped with various middlewares. If you
	// made your own client library, you'd do this work there, so your server
	// could rely on a consistent set of client behavior.
	// The List endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var listEndpoint endpoint.Endpoint
	{
		listEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/list"),
			encodeHTTPListRequest,
			decodeHTTPListResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(otTracer, logger)))...,
		).Endpoint()
		listEndpoint = opentracing.TraceClient(otTracer, "List")(listEndpoint)
		listEndpoint = zipkin.TraceEndpoint(zipkinTracer, "List")(listEndpoint)
		e.ListEndpoint = listEndpoint
	}

	// The Get endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var getEndpoint endpoint.Endpoint
	{
		getEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/get"),
			encodeHTTPGetRequest,
			decodeHTTPGetResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(otTracer, logger)))...,
		).Endpoint()
		getEndpoint = opentracing.TraceClient(otTracer, "Get")(getEndpoint)
		getEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Get")(getEndpoint)
		e.GetEndpoint = getEndpoint
	}

	// The Post endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var postEndpoint endpoint.Endpoint
	{
		postEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/post"),
			encodeHTTPPostRequest,
			decodeHTTPPostResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(otTracer, logger)))...,
		).Endpoint()
		postEndpoint = opentracing.TraceClient(otTracer, "Post")(postEndpoint)
		postEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Post")(postEndpoint)
		e.PostEndpoint = postEndpoint
	}

	// The Complete endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var completeEndpoint endpoint.Endpoint
	{
		completeEndpoint = httptransport.NewClient(
			"PUT",
			copyURL(u, "/complete"),
			encodeHTTPCompleteRequest,
			decodeHTTPCompleteResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(otTracer, logger)))...,
		).Endpoint()
		completeEndpoint = opentracing.TraceClient(otTracer, "Complete")(completeEndpoint)
		completeEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Complete")(completeEndpoint)
		e.CompleteEndpoint = completeEndpoint
	}

	// The ClearComplete endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var clearCompleteEndpoint endpoint.Endpoint
	{
		clearCompleteEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/clearcomplete"),
			encodeHTTPClearCompleteRequest,
			decodeHTTPClearCompleteResponse,
			append(options, httptransport.ClientBefore(opentracing.ContextToHTTP(otTracer, logger)))...,
		).Endpoint()
		clearCompleteEndpoint = opentracing.TraceClient(otTracer, "ClearComplete")(clearCompleteEndpoint)
		clearCompleteEndpoint = zipkin.TraceEndpoint(zipkinTracer, "ClearComplete")(clearCompleteEndpoint)
		e.ClearCompleteEndpoint = clearCompleteEndpoint
	}

	// Returning the endpoint.Set as a service.Service relies on the
	// endpoint.Set implementing the Service methods. That's just a simple bit
	// of glue code.
	return e, nil
}

//
func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}

// encodeHTTPListRequest is a transport/http.EncodeRequestFunc that
// JSON-encodes any request to the request body. Primarily useful in a client.
func encodeHTTPListRequest(_ context.Context, r *http.Request, request interface{}) (err error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// decodeHTTPListResponse is a transport/http.DecodeResponseFunc that decodes a
// JSON-encoded sum response from the HTTP response body. If the response has a
// non-200 status code, we will interpret that as an error and attempt to decode
// the specific error message from the response body. Primarily useful in a client.
func decodeHTTPListResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, responses.JSONErrorDecoder(r)
	}
	var resp endpoints.ListResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// encodeHTTPGetRequest is a transport/http.EncodeRequestFunc that
// JSON-encodes any request to the request body. Primarily useful in a client.
func encodeHTTPGetRequest(_ context.Context, r *http.Request, request interface{}) (err error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// decodeHTTPGetResponse is a transport/http.DecodeResponseFunc that decodes a
// JSON-encoded sum response from the HTTP response body. If the response has a
// non-200 status code, we will interpret that as an error and attempt to decode
// the specific error message from the response body. Primarily useful in a client.
func decodeHTTPGetResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, responses.JSONErrorDecoder(r)
	}
	var resp endpoints.GetResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// encodeHTTPPostRequest is a transport/http.EncodeRequestFunc that
// JSON-encodes any request to the request body. Primarily useful in a client.
func encodeHTTPPostRequest(_ context.Context, r *http.Request, request interface{}) (err error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// decodeHTTPPostResponse is a transport/http.DecodeResponseFunc that decodes a
// JSON-encoded sum response from the HTTP response body. If the response has a
// non-200 status code, we will interpret that as an error and attempt to decode
// the specific error message from the response body. Primarily useful in a client.
func decodeHTTPPostResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, responses.JSONErrorDecoder(r)
	}
	var resp endpoints.PostResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// encodeHTTPCompleteRequest is a transport/http.EncodeRequestFunc that
// JSON-encodes any request to the request body. Primarily useful in a client.
func encodeHTTPCompleteRequest(_ context.Context, r *http.Request, request interface{}) (err error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// decodeHTTPCompleteResponse is a transport/http.DecodeResponseFunc that decodes a
// JSON-encoded sum response from the HTTP response body. If the response has a
// non-200 status code, we will interpret that as an error and attempt to decode
// the specific error message from the response body. Primarily useful in a client.
func decodeHTTPCompleteResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, responses.JSONErrorDecoder(r)
	}
	var resp endpoints.CompleteResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// encodeHTTPClearCompleteRequest is a transport/http.EncodeRequestFunc that
// JSON-encodes any request to the request body. Primarily useful in a client.
func encodeHTTPClearCompleteRequest(_ context.Context, r *http.Request, request interface{}) (err error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// decodeHTTPClearCompleteResponse is a transport/http.DecodeResponseFunc that decodes a
// JSON-encoded sum response from the HTTP response body. If the response has a
// non-200 status code, we will interpret that as an error and attempt to decode
// the specific error message from the response body. Primarily useful in a client.
func decodeHTTPClearCompleteResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, responses.JSONErrorDecoder(r)
	}
	var resp endpoints.ClearCompleteResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func CustomErrorEncoder(errorVal errors.Error) (code int) {
	switch {
	// TODO write your own custom error check here
	}
	return
}
