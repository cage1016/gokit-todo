package transports

import (
	"context"
	"time"

	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cage1016/todo/internal/app/todo/endpoints"
	"github.com/cage1016/todo/internal/app/todo/model"
	"github.com/cage1016/todo/internal/app/todo/service"
	"github.com/cage1016/todo/internal/pkg/errors"
	pb "github.com/cage1016/todo/pb/todo"
)

type grpcServer struct {
	list          grpctransport.Handler `json:""`
	get           grpctransport.Handler `json:""`
	post          grpctransport.Handler `json:""`
	complete      grpctransport.Handler `json:""`
	clearComplete grpctransport.Handler `json:""`
}

func (s *grpcServer) List(ctx context.Context, req *pb.ListRequest) (rep *pb.ListResponse, err error) {
	_, rp, err := s.list.ServeGRPC(ctx, req)
	if err != nil {
		return nil, grpcEncodeError(errors.Cast(err))
	}
	rep = rp.(*pb.ListResponse)
	return rep, nil
}

func (s *grpcServer) Get(ctx context.Context, req *pb.GetRequest) (rep *pb.GetResponse, err error) {
	_, rp, err := s.get.ServeGRPC(ctx, req)
	if err != nil {
		return nil, grpcEncodeError(errors.Cast(err))
	}
	rep = rp.(*pb.GetResponse)
	return rep, nil
}

func (s *grpcServer) Post(ctx context.Context, req *pb.PostRequest) (rep *pb.PostResponse, err error) {
	_, rp, err := s.post.ServeGRPC(ctx, req)
	if err != nil {
		return nil, grpcEncodeError(errors.Cast(err))
	}
	rep = rp.(*pb.PostResponse)
	return rep, nil
}

func (s *grpcServer) Complete(ctx context.Context, req *pb.CompleteRequest) (rep *pb.CompleteResponse, err error) {
	_, rp, err := s.complete.ServeGRPC(ctx, req)
	if err != nil {
		return nil, grpcEncodeError(errors.Cast(err))
	}
	rep = rp.(*pb.CompleteResponse)
	return rep, nil
}

func (s *grpcServer) ClearComplete(ctx context.Context, req *pb.ClearCompleteRequest) (rep *pb.ClearCompleteResponse, err error) {
	_, rp, err := s.clearComplete.ServeGRPC(ctx, req)
	if err != nil {
		return nil, grpcEncodeError(errors.Cast(err))
	}
	rep = rp.(*pb.ClearCompleteResponse)
	return rep, nil
}

// MakeGRPCServer makes a set of endpoints available as a gRPC server.
func MakeGRPCServer(endpoints endpoints.Endpoints, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) (req pb.TodoServer) { // Zipkin GRPC Server Trace can either be instantiated per gRPC method with a
	// provided operation name or a global tracing service can be instantiated
	// without an operation name and fed to each Go kit gRPC server as a
	// ServerOption.
	// In the latter case, the operation name will be the endpoint's grpc method
	// path if used in combination with the Go kit gRPC Interceptor.
	//
	// In this example, we demonstrate a global Zipkin tracing service with
	// Go kit gRPC Interceptor.
	zipkinServer := zipkin.GRPCServerTrace(zipkinTracer)

	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
		zipkinServer,
	}

	return &grpcServer{
		list: grpctransport.NewServer(
			endpoints.ListEndpoint,
			decodeGRPCListRequest,
			encodeGRPCListResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "List", logger), kitjwt.GRPCToContext()))...,
		),

		get: grpctransport.NewServer(
			endpoints.GetEndpoint,
			decodeGRPCGetRequest,
			encodeGRPCGetResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "Get", logger), kitjwt.GRPCToContext()))...,
		),

		post: grpctransport.NewServer(
			endpoints.PostEndpoint,
			decodeGRPCPostRequest,
			encodeGRPCPostResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "Post", logger), kitjwt.GRPCToContext()))...,
		),

		complete: grpctransport.NewServer(
			endpoints.CompleteEndpoint,
			decodeGRPCCompleteRequest,
			encodeGRPCCompleteResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "Complete", logger), kitjwt.GRPCToContext()))...,
		),

		clearComplete: grpctransport.NewServer(
			endpoints.ClearCompleteEndpoint,
			decodeGRPCClearCompleteRequest,
			encodeGRPCClearCompleteResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "ClearComplete", logger), kitjwt.GRPCToContext()))...,
		),
	}
}

// decodeGRPCListRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain request. Primarily useful in a server.
func decodeGRPCListRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ListRequest)
	return endpoints.ListRequest{Filter: req.Filter}, nil
}

// encodeGRPCListResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain response to a gRPC reply. Primarily useful in a server.
func encodeGRPCListResponse(_ context.Context, grpcReply interface{}) (res interface{}, err error) {
	reply := grpcReply.(endpoints.ListResponse)
	if reply.Err != nil {
		return &pb.ListResponse{}, grpcEncodeError(errors.Cast(reply.Err))
	}

	todos := []*pb.ModelTodo{}
	for _, todo := range reply.Res {
		todos = append(todos, ModelToPBTodo(todo))
	}

	return &pb.ListResponse{Res: todos}, grpcEncodeError(errors.Cast(reply.Err))
}

// decodeGRPCGetRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain request. Primarily useful in a server.
func decodeGRPCGetRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetRequest)
	return endpoints.GetRequest{Id: req.Id}, nil
}

// encodeGRPCGetResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain response to a gRPC reply. Primarily useful in a server.
func encodeGRPCGetResponse(_ context.Context, grpcReply interface{}) (res interface{}, err error) {
	reply := grpcReply.(endpoints.GetResponse)
	return &pb.GetResponse{Res: ModelToPBTodo(reply.Res)}, grpcEncodeError(errors.Cast(reply.Err))
}

// decodeGRPCPostRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain request. Primarily useful in a server.
func decodeGRPCPostRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.PostRequest)
	return endpoints.PostRequest{Todo: PBtoModelTodo(req.Todo)}, nil
}

// encodeGRPCPostResponse is a transport/ModelTodo(todo *pb.Modeltogrpc.EncodeResponseFunc that converts a
// user-domain response to a gRPC reply. Primarily useful in a server.
func encodeGRPCPostResponse(_ context.Context, grpcReply interface{}) (res interface{}, err error) {
	reply := grpcReply.(endpoints.PostResponse)
	return &pb.PostResponse{Res: ModelToPBTodo(reply.Res)}, grpcEncodeError(errors.Cast(reply.Err))
}

// decodeGRPCCompleteRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain request. Primarily useful in a server.
func decodeGRPCCompleteRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CompleteRequest)
	return endpoints.CompleteRequest{Id: req.Id}, nil
}

// encodeGRPCCompleteResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain response to a gRPC reply. Primarily useful in a server.
func encodeGRPCCompleteResponse(_ context.Context, grpcReply interface{}) (res interface{}, err error) {
	reply := grpcReply.(endpoints.CompleteResponse)
	return &pb.CompleteResponse{}, grpcEncodeError(errors.Cast(reply.Err))
}

// decodeGRPCClearCompleteRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain request. Primarily useful in a server.
func decodeGRPCClearCompleteRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	_ = grpcReq.(*pb.ClearCompleteRequest)
	return endpoints.ClearCompleteRequest{}, nil
}

// encodeGRPCClearCompleteResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain response to a gRPC reply. Primarily useful in a server.
func encodeGRPCClearCompleteResponse(_ context.Context, grpcReply interface{}) (res interface{}, err error) {
	reply := grpcReply.(endpoints.ClearCompleteResponse)
	return &pb.ClearCompleteResponse{}, grpcEncodeError(errors.Cast(reply.Err))
}

// NewGRPCClient returns an AddService backed by a gRPC server at the other end
// of the conn. The caller is responsible for constructing the conn, and
// eventually closing the underlying transport. We bake-in certain middlewares,
// implementing the client library pattern.
func NewGRPCClient(conn *grpc.ClientConn, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) service.TodoService { // Zipkin GRPC Client Trace can either be instantiated per gRPC method with a
	// provided operation name or a global tracing client can be instantiated
	// without an operation name and fed to each Go kit client as ClientOption.
	// In the latter case, the operation name will be the endpoint's grpc method
	// path.
	//
	// In this example, we demonstrace a global tracing client.
	zipkinClient := zipkin.GRPCClientTrace(zipkinTracer)

	// global client middlewares
	options := []grpctransport.ClientOption{
		zipkinClient,
	}

	// The List endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var listEndpoint endpoint.Endpoint
	{
		listEndpoint = grpctransport.NewClient(
			conn,
			"pb.Todo",
			"List",
			encodeGRPCListRequest,
			decodeGRPCListResponse,
			pb.ListResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger), kitjwt.ContextToGRPC()))...,
		).Endpoint()
		listEndpoint = opentracing.TraceClient(otTracer, "List")(listEndpoint)
	}

	// The Get endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var getEndpoint endpoint.Endpoint
	{
		getEndpoint = grpctransport.NewClient(
			conn,
			"pb.Todo",
			"Get",
			encodeGRPCGetRequest,
			decodeGRPCGetResponse,
			pb.GetResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger), kitjwt.ContextToGRPC()))...,
		).Endpoint()
		getEndpoint = opentracing.TraceClient(otTracer, "Get")(getEndpoint)
	}

	// The Post endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var postEndpoint endpoint.Endpoint
	{
		postEndpoint = grpctransport.NewClient(
			conn,
			"pb.Todo",
			"Post",
			encodeGRPCPostRequest,
			decodeGRPCPostResponse,
			pb.PostResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger), kitjwt.ContextToGRPC()))...,
		).Endpoint()
		postEndpoint = opentracing.TraceClient(otTracer, "Post")(postEndpoint)
	}

	// The Complete endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var completeEndpoint endpoint.Endpoint
	{
		completeEndpoint = grpctransport.NewClient(
			conn,
			"pb.Todo",
			"Complete",
			encodeGRPCCompleteRequest,
			decodeGRPCCompleteResponse,
			pb.CompleteResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger), kitjwt.ContextToGRPC()))...,
		).Endpoint()
		completeEndpoint = opentracing.TraceClient(otTracer, "Complete")(completeEndpoint)
	}

	// The ClearComplete endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var clearCompleteEndpoint endpoint.Endpoint
	{
		clearCompleteEndpoint = grpctransport.NewClient(
			conn,
			"pb.Todo",
			"ClearComplete",
			encodeGRPCClearCompleteRequest,
			decodeGRPCClearCompleteResponse,
			pb.ClearCompleteResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger), kitjwt.ContextToGRPC()))...,
		).Endpoint()
		clearCompleteEndpoint = opentracing.TraceClient(otTracer, "ClearComplete")(clearCompleteEndpoint)
	}

	return endpoints.Endpoints{
		ListEndpoint:          listEndpoint,
		GetEndpoint:           getEndpoint,
		PostEndpoint:          postEndpoint,
		CompleteEndpoint:      completeEndpoint,
		ClearCompleteEndpoint: clearCompleteEndpoint,
	}
}

// encodeGRPCListRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain List request to a gRPC List request. Primarily useful in a client.
func encodeGRPCListRequest(_ context.Context, request interface{}) (interface{}, error) {
	_ = request.(endpoints.ListRequest)
	return &pb.ListRequest{}, nil
}

// decodeGRPCListResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC List reply to a user-domain List response. Primarily useful in a client.
func decodeGRPCListResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.ListResponse)

	todos := []model.Todo{}
	for _, todo := range reply.Res {
		todos = append(todos, PBtoModelTodo(todo))
	}

	return endpoints.ListResponse{Res: todos}, nil
}

// encodeGRPCGetRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain Get request to a gRPC Get request. Primarily useful in a client.
func encodeGRPCGetRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(endpoints.GetRequest)
	return &pb.GetRequest{Id: req.Id}, nil
}

// decodeGRPCGetResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC Get reply to a user-domain Get response. Primarily useful in a client.
func decodeGRPCGetResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.GetResponse)
	return endpoints.GetResponse{Res: PBtoModelTodo(reply.Res)}, nil
}

// encodeGRPCPostRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain Post request to a gRPC Post request. Primarily useful in a client.
func encodeGRPCPostRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(endpoints.PostRequest)
	return &pb.PostRequest{Todo: ModelToPBTodo(req.Todo)}, nil
}

// decodeGRPCPostResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC Post reply to a user-domain Post response. Primarily useful in a client.
func decodeGRPCPostResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.PostResponse)
	return endpoints.PostResponse{Res: PBtoModelTodo(reply.Res)}, nil
}

// encodeGRPCCompleteRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain Complete request to a gRPC Complete request. Primarily useful in a client.
func encodeGRPCCompleteRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(endpoints.CompleteRequest)
	return &pb.CompleteRequest{Id: req.Id}, nil
}

// decodeGRPCCompleteResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC Complete reply to a user-domain Complete response. Primarily useful in a client.
func decodeGRPCCompleteResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	_ = grpcReply.(*pb.CompleteResponse)
	return endpoints.CompleteResponse{}, nil
}

// encodeGRPCClearCompleteRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain ClearComplete request to a gRPC ClearComplete request. Primarily useful in a client.
func encodeGRPCClearCompleteRequest(_ context.Context, request interface{}) (interface{}, error) {
	_ = request.(endpoints.ClearCompleteRequest)
	return &pb.ClearCompleteRequest{}, nil
}

// decodeGRPCClearCompleteResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC ClearComplete reply to a user-domain ClearComplete response. Primarily useful in a client.
func decodeGRPCClearCompleteResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	_ = grpcReply.(*pb.ClearCompleteResponse)
	return endpoints.ClearCompleteResponse{}, nil
}

func grpcEncodeError(err errors.Error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if ok {
		return status.Error(st.Code(), st.Message())
	}

	switch {
	case errors.Contains(err, service.ErrInvalidQueryParams),
		errors.Contains(err, service.ErrMalformedEntity):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Contains(err, service.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Contains(err, kitjwt.ErrTokenContextMissing):
		return status.Error(codes.Unauthenticated, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}

func ModelToPBTodo(todo model.Todo) *pb.ModelTodo {
	return &pb.ModelTodo{
		Id:        todo.ID,
		CreatedAt: todo.CreatedAt.Format(time.RFC3339),
		UpdatedAt: todo.UpdatedAt.Format(time.RFC3339),
		Text:      todo.Text,
		Complete:  todo.Complete,
	}
}

func PBtoModelTodo(todo *pb.ModelTodo) model.Todo {
	return model.Todo{
		ID: todo.Id,
		CreatedAt: func() time.Time {
			t, err := time.Parse(time.RFC3339, todo.CreatedAt)
			if err != nil {
				return time.Time{}
			}
			return t
		}(),
		UpdatedAt: func() time.Time {
			t, err := time.Parse(time.RFC3339, todo.UpdatedAt)
			if err != nil {
				return time.Time{}
			}
			return t
		}(),
		Text:     todo.Text,
		Complete: todo.Complete,
	}
}
