// +build !integration

package transports_test

import (
	"context"
	"net"
	"os"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/cage1016/gokit-todo/internal/app/todo/endpoints"
	"github.com/cage1016/gokit-todo/internal/app/todo/model"
	transports "github.com/cage1016/gokit-todo/internal/app/todo/transports/grpc"
	automocks "github.com/cage1016/gokit-todo/internal/mocks/app/todo/service"
	pb "github.com/cage1016/gokit-todo/pb/todo"
)

const (
	hostPort string = "localhost:8002"
)

func TestGrpcServer_Add(t *testing.T) {
	type fields struct {
		svc *automocks.MockTodoService
	}

	type args struct {
		todo *model.TodoReq
	}

	text := "aa"
	completed := true

	tests := []struct {
		name      string
		prepare   func(f *fields)
		args      args
		wantErr   bool
		checkFunc func(res *model.TodoRes, err error)
	}{
		{
			name: "grpc add todo",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.svc.EXPECT().Add(gomock.Any(), gomock.Any()).Return(&model.TodoRes{
						ID:        "iKe0KxpurIn0E_6vzUDAr",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						Text:      "aa",
						Completed: false,
					}, nil),
				)
			},
			args: args{todo: &model.TodoReq{
				Text:      &text,
				Completed: &completed,
			}},
			checkFunc: func(res *model.TodoRes, err error) {
				assert.Equal(t, "aa", res.Text)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				svc: automocks.NewMockTodoService(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			logger := log.NewLogfmtLogger(os.Stderr)
			zkt, _ := zipkin.NewTracer(nil, zipkin.WithNoopTracer(true))
			tracer := opentracing.GlobalTracer()

			// server
			server := grpc.NewServer()
			eps := endpoints.New(f.svc, logger, tracer, zkt)
			sc, err := net.Listen("tcp", hostPort)
			if err != nil {
				t.Fatalf("unable to listen: %+v", err)
			}
			defer server.GracefulStop()

			go func() {
				pb.RegisterTodoServer(server, transports.MakeGRPCServer(eps, tracer, zkt, logger))
				_ = server.Serve(sc)
			}()

			// client
			cc, err := grpc.Dial(hostPort, grpc.WithInsecure())
			if err != nil {
				t.Fatalf("unable to Dial: %+v", err)
			}
			svc := transports.NewGRPCClient(cc, tracer, zkt, logger)

			if res, err := svc.Add(context.Background(), tt.args.todo); (err != nil) != tt.wantErr {
				t.Errorf("svc.Post error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if tt.checkFunc != nil {
					tt.checkFunc(res, err)
				}
			}
		})
	}
}

func TestGrpcServer_Update(t *testing.T) {
	type fields struct {
		svc *automocks.MockTodoService
	}

	type args struct {
		id   string
		todo *model.TodoReq
	}

	text := "aa"
	completed := true

	tests := []struct {
		name      string
		prepare   func(f *fields)
		args      args
		wantErr   bool
		checkFunc func(res *model.TodoRes, err error)
	}{
		{
			name: "grpc update todo",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.svc.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(&model.TodoRes{
						ID:        "iKe0KxpurIn0E_6vzUDAr",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						Text:      "bb",
						Completed: false,
					}, nil),
				)
			},
			args: args{todo: &model.TodoReq{
				Text:      &text,
				Completed: &completed,
			}},
			checkFunc: func(res *model.TodoRes, err error) {
				assert.Equal(t, "bb", res.Text)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				svc: automocks.NewMockTodoService(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			logger := log.NewLogfmtLogger(os.Stderr)
			zkt, _ := zipkin.NewTracer(nil, zipkin.WithNoopTracer(true))
			tracer := opentracing.GlobalTracer()

			// server
			server := grpc.NewServer()
			eps := endpoints.New(f.svc, logger, tracer, zkt)
			sc, err := net.Listen("tcp", hostPort)
			if err != nil {
				t.Fatalf("unable to listen: %+v", err)
			}
			defer server.GracefulStop()

			go func() {
				pb.RegisterTodoServer(server, transports.MakeGRPCServer(eps, tracer, zkt, logger))
				_ = server.Serve(sc)
			}()

			// client
			cc, err := grpc.Dial(hostPort, grpc.WithInsecure())
			if err != nil {
				t.Fatalf("unable to Dial: %+v", err)
			}
			svc := transports.NewGRPCClient(cc, tracer, zkt, logger)

			if res, err := svc.Update(context.Background(), tt.args.id, tt.args.todo); (err != nil) != tt.wantErr {
				t.Errorf("svc.Update error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if tt.checkFunc != nil {
					tt.checkFunc(res, err)
				}
			}
		})
	}
}

func TestGrpcServer_List(t *testing.T) {
	type fields struct {
		svc *automocks.MockTodoService
	}

	type args struct {
	}

	tests := []struct {
		name      string
		prepare   func(f *fields)
		args      args
		wantErr   bool
		checkFunc func(res []*model.TodoRes, err error)
	}{
		{
			name: "grpc list todo",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.svc.EXPECT().List(gomock.Any()).Return([]*model.TodoRes{{
						ID:        "iKe0KxpurIn0E_6vzUDAr",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						Text:      "aa",
						Completed: false,
					}}, nil),
				)
			},
			checkFunc: func(res []*model.TodoRes, err error) {
				assert.Equal(t, len(res), 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				svc: automocks.NewMockTodoService(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			logger := log.NewLogfmtLogger(os.Stderr)
			zkt, _ := zipkin.NewTracer(nil, zipkin.WithNoopTracer(true))
			tracer := opentracing.GlobalTracer()

			// server
			server := grpc.NewServer()
			eps := endpoints.New(f.svc, logger, tracer, zkt)
			sc, err := net.Listen("tcp", hostPort)
			if err != nil {
				t.Fatalf("unable to listen: %+v", err)
			}
			defer server.GracefulStop()

			go func() {
				pb.RegisterTodoServer(server, transports.MakeGRPCServer(eps, tracer, zkt, logger))
				_ = server.Serve(sc)
			}()

			// client
			cc, err := grpc.Dial(hostPort, grpc.WithInsecure())
			if err != nil {
				t.Fatalf("unable to Dial: %+v", err)
			}
			svc := transports.NewGRPCClient(cc, tracer, zkt, logger)

			if res, err := svc.List(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("svc.Update error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if tt.checkFunc != nil {
					tt.checkFunc(res, err)
				}
			}
		})
	}
}

func TestGrpcServer_Delete(t *testing.T) {
	type fields struct {
		svc *automocks.MockTodoService
	}

	type args struct {
		id string
	}

	tests := []struct {
		name      string
		prepare   func(f *fields)
		args      args
		wantErr   bool
		checkFunc func(err error)
	}{
		{
			name: "grpc delete todo",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.svc.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil),
				)
			},
			args: args{id: "iKe0KxpurIn0E_6vzUDAr"},
			checkFunc: func(err error) {
				assert.Nil(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				svc: automocks.NewMockTodoService(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			logger := log.NewLogfmtLogger(os.Stderr)
			zkt, _ := zipkin.NewTracer(nil, zipkin.WithNoopTracer(true))
			tracer := opentracing.GlobalTracer()

			// server
			server := grpc.NewServer()
			eps := endpoints.New(f.svc, logger, tracer, zkt)
			sc, err := net.Listen("tcp", hostPort)
			if err != nil {
				t.Fatalf("unable to listen: %+v", err)
			}
			defer server.GracefulStop()

			go func() {
				pb.RegisterTodoServer(server, transports.MakeGRPCServer(eps, tracer, zkt, logger))
				_ = server.Serve(sc)
			}()

			// client
			cc, err := grpc.Dial(hostPort, grpc.WithInsecure())
			if err != nil {
				t.Fatalf("unable to Dial: %+v", err)
			}
			svc := transports.NewGRPCClient(cc, tracer, zkt, logger)

			if err := svc.Delete(context.Background(), tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("svc.Update error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if tt.checkFunc != nil {
					tt.checkFunc(err)
				}
			}
		})
	}
}
