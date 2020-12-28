package transports

import (
	"context"

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
	add         grpctransport.Handler `json:""`
	delete      grpctransport.Handler `json:""`
	update      grpctransport.Handler `json:""`
	list        grpctransport.Handler `json:""`
	complete    grpctransport.Handler `json:""`
	completeAll grpctransport.Handler `json:""`
	clear       grpctransport.Handler `json:""`
}

func (s *grpcServer) Add(ctx context.Context, req *pb.AddRequest) (rep *pb.AddResponse, err error) {
	_, rp, err := s.add.ServeGRPC(ctx, req)
	if err != nil {
		return nil, grpcEncodeError(errors.Cast(err))
	}
	rep = rp.(*pb.AddResponse)
	return rep, nil
}

func (s *grpcServer) Delete(ctx context.Context, req *pb.DeleteRequest) (rep *pb.DeleteResponse, err error) {
	_, rp, err := s.delete.ServeGRPC(ctx, req)
	if err != nil {
		return nil, grpcEncodeError(errors.Cast(err))
	}
	rep = rp.(*pb.DeleteResponse)
	return rep, nil
}

func (s *grpcServer) Update(ctx context.Context, req *pb.UpdateRequest) (rep *pb.UpdateResponse, err error) {
	_, rp, err := s.update.ServeGRPC(ctx, req)
	if err != nil {
		return nil, grpcEncodeError(errors.Cast(err))
	}
	rep = rp.(*pb.UpdateResponse)
	return rep, nil
}

func (s *grpcServer) List(ctx context.Context, req *pb.ListRequest) (rep *pb.ListResponse, err error) {
	_, rp, err := s.list.ServeGRPC(ctx, req)
	if err != nil {
		return nil, grpcEncodeError(errors.Cast(err))
	}
	rep = rp.(*pb.ListResponse)
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

func (s *grpcServer) CompleteAll(ctx context.Context, req *pb.CompleteAllRequest) (rep *pb.CompleteAllResponse, err error) {
	_, rp, err := s.completeAll.ServeGRPC(ctx, req)
	if err != nil {
		return nil, grpcEncodeError(errors.Cast(err))
	}
	rep = rp.(*pb.CompleteAllResponse)
	return rep, nil
}

func (s *grpcServer) Clear(ctx context.Context, req *pb.ClearRequest) (rep *pb.ClearResponse, err error) {
	_, rp, err := s.clear.ServeGRPC(ctx, req)
	if err != nil {
		return nil, grpcEncodeError(errors.Cast(err))
	}
	rep = rp.(*pb.ClearResponse)
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
		add: grpctransport.NewServer(
			endpoints.AddEndpoint,
			decodeGRPCAddRequest,
			encodeGRPCAddResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "Add", logger), kitjwt.GRPCToContext()))...,
		),

		delete: grpctransport.NewServer(
			endpoints.DeleteEndpoint,
			decodeGRPCDeleteRequest,
			encodeGRPCDeleteResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "Delete", logger), kitjwt.GRPCToContext()))...,
		),

		update: grpctransport.NewServer(
			endpoints.UpdateEndpoint,
			decodeGRPCUpdateRequest,
			encodeGRPCUpdateResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "Update", logger), kitjwt.GRPCToContext()))...,
		),

		list: grpctransport.NewServer(
			endpoints.ListEndpoint,
			decodeGRPCListRequest,
			encodeGRPCListResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "List", logger), kitjwt.GRPCToContext()))...,
		),

		complete: grpctransport.NewServer(
			endpoints.CompleteEndpoint,
			decodeGRPCCompleteRequest,
			encodeGRPCCompleteResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "Completed", logger), kitjwt.GRPCToContext()))...,
		),

		completeAll: grpctransport.NewServer(
			endpoints.CompleteAllEndpoint,
			decodeGRPCCompleteAllRequest,
			encodeGRPCCompleteAllResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "CompleteAll", logger), kitjwt.GRPCToContext()))...,
		),

		clear: grpctransport.NewServer(
			endpoints.ClearEndpoint,
			decodeGRPCClearRequest,
			encodeGRPCClearResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "Clear", logger), kitjwt.GRPCToContext()))...,
		),
	}
}

// decodeGRPCAddRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain request. Primarily useful in a server.
func decodeGRPCAddRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.AddRequest)
	return endpoints.AddRequest{Todo: PBtoModel(req.Todo)}, nil
}

// encodeGRPCAddResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain response to a gRPC reply. Primarily useful in a server.
func encodeGRPCAddResponse(_ context.Context, grpcReply interface{}) (res interface{}, err error) {
	reply := grpcReply.(endpoints.AddResponse)
	return &pb.AddResponse{Res: ModelToPB(reply.Res)}, grpcEncodeError(errors.Cast(reply.Err))
}

// decodeGRPCDeleteRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain request. Primarily useful in a server.
func decodeGRPCDeleteRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.DeleteRequest)
	return endpoints.DeleteRequest{Id: req.Id}, nil
}

// encodeGRPCDeleteResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain response to a gRPC reply. Primarily useful in a server.
func encodeGRPCDeleteResponse(_ context.Context, grpcReply interface{}) (res interface{}, err error) {
	reply := grpcReply.(endpoints.DeleteResponse)
	return &pb.DeleteResponse{}, grpcEncodeError(errors.Cast(reply.Err))
}

// decodeGRPCUpdateRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain request. Primarily useful in a server.
func decodeGRPCUpdateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UpdateRequest)
	return endpoints.UpdateRequest{Id: req.Id, Todo: PBtoModel(req.Todo)}, nil
}

// encodeGRPCUpdateResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain response to a gRPC reply. Primarily useful in a server.
func encodeGRPCUpdateResponse(_ context.Context, grpcReply interface{}) (res interface{}, err error) {
	reply := grpcReply.(endpoints.UpdateResponse)
	return &pb.UpdateResponse{Res: ModelToPB(reply.Res)}, grpcEncodeError(errors.Cast(reply.Err))
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
		todos = append(todos, ModelToPB(todo))
	}

	return &pb.ListResponse{Res: todos}, grpcEncodeError(errors.Cast(reply.Err))
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

// decodeGRPCCompleteAllRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain request. Primarily useful in a server.
func decodeGRPCCompleteAllRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	_ = grpcReq.(*pb.CompleteAllRequest)
	return endpoints.CompleteAllRequest{}, nil
}

// encodeGRPCCompleteAllResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain response to a gRPC reply. Primarily useful in a server.
func encodeGRPCCompleteAllResponse(_ context.Context, grpcReply interface{}) (res interface{}, err error) {
	reply := grpcReply.(endpoints.CompleteAllResponse)
	return &pb.CompleteAllResponse{}, grpcEncodeError(errors.Cast(reply.Err))
}

// decodeGRPCClearRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain request. Primarily useful in a server.
func decodeGRPCClearRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	_ = grpcReq.(*pb.ClearRequest)
	return endpoints.ClearRequest{}, nil
}

// encodeGRPCClearResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain response to a gRPC reply. Primarily useful in a server.
func encodeGRPCClearResponse(_ context.Context, grpcReply interface{}) (res interface{}, err error) {
	reply := grpcReply.(endpoints.ClearResponse)
	return &pb.ClearResponse{}, grpcEncodeError(errors.Cast(reply.Err))
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

	// The Add endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var addEndpoint endpoint.Endpoint
	{
		addEndpoint = grpctransport.NewClient(
			conn,
			"pb.Todo",
			"Add",
			encodeGRPCAddRequest,
			decodeGRPCAddResponse,
			pb.AddResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger), kitjwt.ContextToGRPC()))...,
		).Endpoint()
		addEndpoint = opentracing.TraceClient(otTracer, "Add")(addEndpoint)
	}

	// The Delete endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var deleteEndpoint endpoint.Endpoint
	{
		deleteEndpoint = grpctransport.NewClient(
			conn,
			"pb.Todo",
			"Delete",
			encodeGRPCDeleteRequest,
			decodeGRPCDeleteResponse,
			pb.DeleteResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger), kitjwt.ContextToGRPC()))...,
		).Endpoint()
		deleteEndpoint = opentracing.TraceClient(otTracer, "Delete")(deleteEndpoint)
	}

	// The Update endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var updateEndpoint endpoint.Endpoint
	{
		updateEndpoint = grpctransport.NewClient(
			conn,
			"pb.Todo",
			"Update",
			encodeGRPCUpdateRequest,
			decodeGRPCUpdateResponse,
			pb.UpdateResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger), kitjwt.ContextToGRPC()))...,
		).Endpoint()
		updateEndpoint = opentracing.TraceClient(otTracer, "Update")(updateEndpoint)
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

	// The Completed endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var completeEndpoint endpoint.Endpoint
	{
		completeEndpoint = grpctransport.NewClient(
			conn,
			"pb.Todo",
			"Completed",
			encodeGRPCCompleteRequest,
			decodeGRPCCompleteResponse,
			pb.CompleteResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger), kitjwt.ContextToGRPC()))...,
		).Endpoint()
		completeEndpoint = opentracing.TraceClient(otTracer, "Completed")(completeEndpoint)
	}

	// The CompleteAll endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var completeAllEndpoint endpoint.Endpoint
	{
		completeAllEndpoint = grpctransport.NewClient(
			conn,
			"pb.Todo",
			"CompleteAll",
			encodeGRPCCompleteAllRequest,
			decodeGRPCCompleteAllResponse,
			pb.CompleteAllResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger), kitjwt.ContextToGRPC()))...,
		).Endpoint()
		completeAllEndpoint = opentracing.TraceClient(otTracer, "CompleteAll")(completeAllEndpoint)
	}

	// The Clear endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var clearEndpoint endpoint.Endpoint
	{
		clearEndpoint = grpctransport.NewClient(
			conn,
			"pb.Todo",
			"Clear",
			encodeGRPCClearRequest,
			decodeGRPCClearResponse,
			pb.ClearResponse{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger), kitjwt.ContextToGRPC()))...,
		).Endpoint()
		clearEndpoint = opentracing.TraceClient(otTracer, "Clear")(clearEndpoint)
	}

	return endpoints.Endpoints{
		AddEndpoint:         addEndpoint,
		DeleteEndpoint:      deleteEndpoint,
		UpdateEndpoint:      updateEndpoint,
		ListEndpoint:        listEndpoint,
		CompleteEndpoint:    completeEndpoint,
		CompleteAllEndpoint: completeAllEndpoint,
		ClearEndpoint:       clearEndpoint,
	}
}

// encodeGRPCAddRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain Add request to a gRPC Add request. Primarily useful in a client.
func encodeGRPCAddRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(endpoints.AddRequest)
	return &pb.AddRequest{Todo: ModelToPB(req.Todo)}, nil
}

// decodeGRPCAddResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC Add reply to a user-domain Add response. Primarily useful in a client.
func decodeGRPCAddResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.AddResponse)
	return endpoints.AddResponse{Res: PBtoModel(reply.Res)}, nil
}

// encodeGRPCDeleteRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain Delete request to a gRPC Delete request. Primarily useful in a client.
func encodeGRPCDeleteRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(endpoints.DeleteRequest)
	return &pb.DeleteRequest{Id: req.Id}, nil
}

// decodeGRPCDeleteResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC Delete reply to a user-domain Delete response. Primarily useful in a client.
func decodeGRPCDeleteResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	_ = grpcReply.(*pb.DeleteResponse)
	return endpoints.DeleteResponse{}, nil
}

// encodeGRPCUpdateRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain Update request to a gRPC Update request. Primarily useful in a client.
func encodeGRPCUpdateRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(endpoints.UpdateRequest)
	return &pb.UpdateRequest{Id: req.Id, Todo: ModelToPB(req.Todo)}, nil
}

// decodeGRPCUpdateResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC Update reply to a user-domain Update response. Primarily useful in a client.
func decodeGRPCUpdateResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UpdateResponse)
	return endpoints.UpdateResponse{Res: PBtoModel(reply.Res)}, nil
}

// encodeGRPCListRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain List request to a gRPC List request. Primarily useful in a client.
func encodeGRPCListRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(endpoints.ListRequest)
	return &pb.ListRequest{Filter: req.Filter}, nil
}

// decodeGRPCListResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC List reply to a user-domain List response. Primarily useful in a client.
func decodeGRPCListResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.ListResponse)

	todos := []*model.Todo{}
	for _, todo := range reply.Res {
		todos = append(todos, PBtoModel(todo))
	}

	return endpoints.ListResponse{Res: todos}, nil
}

// encodeGRPCCompleteRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain Completed request to a gRPC Completed request. Primarily useful in a client.
func encodeGRPCCompleteRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(endpoints.CompleteRequest)
	return &pb.CompleteRequest{Id: req.Id}, nil
}

// decodeGRPCCompleteResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC Completed reply to a user-domain Completed response. Primarily useful in a client.
func decodeGRPCCompleteResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	_ = grpcReply.(*pb.CompleteResponse)
	return endpoints.CompleteResponse{}, nil
}

// encodeGRPCCompleteAllRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain CompleteAll request to a gRPC CompleteAll request. Primarily useful in a client.
func encodeGRPCCompleteAllRequest(_ context.Context, request interface{}) (interface{}, error) {
	_ = request.(endpoints.CompleteAllRequest)
	return &pb.CompleteAllRequest{}, nil
}

// decodeGRPCCompleteAllResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC CompleteAll reply to a user-domain CompleteAll response. Primarily useful in a client.
func decodeGRPCCompleteAllResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	_ = grpcReply.(*pb.CompleteAllResponse)
	return endpoints.CompleteAllResponse{}, nil
}

// encodeGRPCClearRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain Clear request to a gRPC Clear request. Primarily useful in a client.
func encodeGRPCClearRequest(_ context.Context, request interface{}) (interface{}, error) {
	_ = request.(endpoints.ClearRequest)
	return &pb.ClearRequest{}, nil
}

// decodeGRPCClearResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC Clear reply to a user-domain Clear response. Primarily useful in a client.
func decodeGRPCClearResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	_ = grpcReply.(*pb.ClearResponse)
	return endpoints.ClearResponse{}, nil
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
	// TODO write your own custom error check here
	case errors.Contains(err, kitjwt.ErrTokenContextMissing):
		return status.Error(codes.Unauthenticated, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
