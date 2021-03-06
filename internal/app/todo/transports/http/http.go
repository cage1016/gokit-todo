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

	"github.com/cage1016/gokit-todo/internal/app/todo/endpoints"
	"github.com/cage1016/gokit-todo/internal/app/todo/service"
	"github.com/cage1016/gokit-todo/internal/pkg/errors"
	"github.com/cage1016/gokit-todo/internal/pkg/responses"
)

const (
	contentType string = "application/json"
)

// ShowTodo godoc
// @Summary Add
// @Description TODO
// @Tags TODO
// @Accept json
// @Produce json
// @Router /items [post]
func AddHandler(m *bone.Mux, endpoints endpoints.Endpoints, options []httptransport.ServerOption, otTracer stdopentracing.Tracer, logger log.Logger) {
	m.Post("/items", httptransport.NewServer(
		endpoints.AddEndpoint,
		decodeHTTPAddRequest,
		responses.EncodeJSONResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Add", logger), kitjwt.HTTPToContext()))...,
	))
}

// ShowTodo godoc
// @Summary Delete
// @Description TODO
// @Tags TODO
// @Accept json
// @Produce json
// @Router /items/:id [delete]
func DeleteHandler(m *bone.Mux, endpoints endpoints.Endpoints, options []httptransport.ServerOption, otTracer stdopentracing.Tracer, logger log.Logger) {
	m.Delete("/items/:id", httptransport.NewServer(
		endpoints.DeleteEndpoint,
		decodeHTTPDeleteRequest,
		responses.EncodeJSONResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Delete", logger), kitjwt.HTTPToContext()))...,
	))
}

// ShowTodo godoc
// @Summary Update
// @Description TODO
// @Tags TODO
// @Accept json
// @Produce json
// @Router /items/:id [put]
func UpdateHandler(m *bone.Mux, endpoints endpoints.Endpoints, options []httptransport.ServerOption, otTracer stdopentracing.Tracer, logger log.Logger) {
	m.Patch("/items/:id", httptransport.NewServer(
		endpoints.UpdateEndpoint,
		decodeHTTPUpdateRequest,
		responses.EncodeJSONResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Update", logger), kitjwt.HTTPToContext()))...,
	))
}

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
	AddHandler(m, endpoints, options, otTracer, logger)
	DeleteHandler(m, endpoints, options, otTracer, logger)
	UpdateHandler(m, endpoints, options, otTracer, logger)
	ListHandler(m, endpoints, options, otTracer, logger)
	return cors.AllowAll().Handler(m)
}

// decodeHTTPAddRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func decodeHTTPAddRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.AddRequest
	err := json.NewDecoder(r.Body).Decode(&req.Todo)
	return req, err
}

// decodeHTTPDeleteRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func decodeHTTPDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.DeleteRequest
	req.Id = bone.GetValue(r, "id")
	return req, nil
}

// decodeHTTPUpdateRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func decodeHTTPUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.UpdateRequest
	req.Id = bone.GetValue(r, "id")
	err := json.NewDecoder(r.Body).Decode(&req.Todo)
	return req, err
}

// decodeHTTPListRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func decodeHTTPListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.ListRequest
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
