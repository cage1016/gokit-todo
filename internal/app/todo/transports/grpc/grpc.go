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

	"github.com/cage1016/gokit-todo/internal/app/todo/endpoints"
	"github.com/cage1016/gokit-todo/internal/app/todo/model"
	"github.com/cage1016/gokit-todo/internal/app/todo/service"
	"github.com/cage1016/gokit-todo/internal/pkg/errors"
	pb "github.com/cage1016/gokit-todo/pb/todo"
)

type grpcServer struct {
	add    grpctransport.Handler `json:""`
	delete grpctransport.Handler `json:""`
	update grpctransport.Handler `json:""`
	list   grpctransport.Handler `json:""`
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
	}
}

// decodeGRPCAddRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain request. Primarily useful in a server.
func decodeGRPCAddRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.AddRequest)
	return endpoints.AddRequest{Todo: PBtoModelReq(req.Todo)}, nil
}

// encodeGRPCAddResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain response to a gRPC reply. Primarily useful in a server.
func encodeGRPCAddResponse(_ context.Context, grpcReply interface{}) (res interface{}, err error) {
	reply := grpcReply.(endpoints.AddResponse)
	return &pb.AddResponse{Res: ModelResToPB(reply.Res)}, grpcEncodeError(errors.Cast(reply.Err))
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
	return endpoints.UpdateRequest{Id: req.Id, Todo: PBtoModelReq(req.Todo)}, nil
}

// encodeGRPCUpdateResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain response to a gRPC reply. Primarily useful in a server.
func encodeGRPCUpdateResponse(_ context.Context, grpcReply interface{}) (res interface{}, err error) {
	reply := grpcReply.(endpoints.UpdateResponse)
	return &pb.UpdateResponse{Res: ModelResToPB(reply.Res)}, grpcEncodeError(errors.Cast(reply.Err))
}

// decodeGRPCListRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain request. Primarily useful in a server.
func decodeGRPCListRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	_ = grpcReq.(*pb.ListRequest)
	return endpoints.ListRequest{}, nil
}

// encodeGRPCListResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain response to a gRPC reply. Primarily useful in a server.
func encodeGRPCListResponse(_ context.Context, grpcReply interface{}) (res interface{}, err error) {
	reply := grpcReply.(endpoints.ListResponse)
	if reply.Err != nil {
		return &pb.ListResponse{}, grpcEncodeError(errors.Cast(reply.Err))
	}

	todos := []*pb.ModelTodoRes{}
	for _, todo := range reply.Res {
		todos = append(todos, ModelResToPB(todo))
	}

	return &pb.ListResponse{Res: todos}, grpcEncodeError(errors.Cast(reply.Err))
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

	return endpoints.Endpoints{
		AddEndpoint:    addEndpoint,
		DeleteEndpoint: deleteEndpoint,
		UpdateEndpoint: updateEndpoint,
		ListEndpoint:   listEndpoint,
	}
}

// encodeGRPCAddRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain Add request to a gRPC Add request. Primarily useful in a client.
func encodeGRPCAddRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(endpoints.AddRequest)
	return &pb.AddRequest{Todo: ModelReqToPB(req.Todo)}, nil
}

// decodeGRPCAddResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC Add reply to a user-domain Add response. Primarily useful in a client.
func decodeGRPCAddResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.AddResponse)
	return endpoints.AddResponse{Res: PBtoModelRes(reply.Res)}, nil
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
	return &pb.UpdateRequest{Id: req.Id, Todo: ModelReqToPB(req.Todo)}, nil
}

// decodeGRPCUpdateResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC Update reply to a user-domain Update response. Primarily useful in a client.
func decodeGRPCUpdateResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UpdateResponse)
	return endpoints.UpdateResponse{Res: PBtoModelRes(reply.Res)}, nil
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

	todos := []*model.TodoRes{}
	for _, todo := range reply.Res {
		todos = append(todos, PBtoModelRes(todo))
	}

	return endpoints.ListResponse{Res: todos}, nil
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
