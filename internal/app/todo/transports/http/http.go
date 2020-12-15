package transports

import (
	"context"
	"encoding/json"
	"net/http"

	kitjwt "github.com/go-kit/kit/auth/jwt"
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
	s := bone.GetQuery(r, "q")
	if len(s) > 1 {
		return nil, service.ErrInvalidQueryParams
	}

	if len(s) == 0 {
		return nil, service.ErrInvalidQueryParams
	}

	q := s[0]
	req.Filter = q
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

func CustomErrorEncoder(errorVal errors.Error) (code int) {
	switch {
	case errors.Contains(errorVal, service.ErrInvalidQueryParams),
		errors.Contains(errorVal, service.ErrMalformedEntity):
		code = http.StatusBadRequest
	case errors.Contains(errorVal, service.ErrNotFound):
		code = http.StatusNotFound
	}
	return
}
