// +build !integration

package service_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/cage1016/gokit-todo/internal/app/todo/model"
	"github.com/cage1016/gokit-todo/internal/app/todo/service"
	automocks "github.com/cage1016/gokit-todo/internal/mocks/app/todo/model"
)

func TestStubTodoService_Add(t *testing.T) {
	type fields struct {
		repo *automocks.MockTodoRepository
	}

	type args struct {
		todo *model.TodoReq
	}

	text := "aa"
	completed := false

	tests := []struct {
		name      string
		prepare   func(f *fields)
		args      args
		wantErr   bool
		checkFunc func(res *model.TodoRes, err error)
	}{
		{
			name: "Add todo",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().Add(context.Background(), gomock.Any()).Return(nil),
				)
			},
			args: args{todo: &model.TodoReq{
				Text:      &text,
				Completed: &completed,
			}},
			wantErr: false,
			checkFunc: func(res *model.TodoRes, err error) {
				assert.Nil(t, err, fmt.Sprintf("should return nil: expected nil got %v", err))
				assert.Equal(t, "aa", res.Text, fmt.Sprintf("text: expected aa got %v", res.Text))
			},
		},
		{
			name: "Add todo fail",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().Add(context.Background(), gomock.Any()).Return(sql.ErrNoRows),
				)
			},
			args: args{todo: &model.TodoReq{
				Text:      &text,
				Completed: &completed,
			}},
			wantErr: true,
			checkFunc: func(res *model.TodoRes, err error) {
				assert.Equal(t, err, sql.ErrNoRows, fmt.Sprintf("err: expected sql.ErrNoRows got %v", err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				repo: automocks.NewMockTodoRepository(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			svc := service.New(f.repo, log.NewLogfmtLogger(os.Stderr))
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

func TestLoggingMiddleware_List(t *testing.T) {
	type fields struct {
		repo *automocks.MockTodoRepository
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
			name: "list todo",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().List(context.Background()).Return([]*model.Todo{
						{
							ID:        "b5z2zC5c9O6~Ns_qLVmn~",
							Completed: false,
							Text:      "aa",
						},
					}, nil),
				)
			},
			wantErr: false,
			checkFunc: func(res []*model.TodoRes, err error) {
				assert.Nil(t, err, fmt.Sprintf("should return nil: expected nil got %v", err))
				assert.Equal(t, len(res), 1, fmt.Sprintf("count res: expected 1 got %v", len(res)))
			},
		},
		{
			name: "list todo fial",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().List(context.Background()).Return([]*model.Todo{}, sql.ErrNoRows),
				)
			},
			wantErr: true,
			checkFunc: func(res []*model.TodoRes, err error) {
				assert.Equal(t, err, sql.ErrNoRows, fmt.Sprintf("err: expected sql.ErrNoRows got %v", err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				repo: automocks.NewMockTodoRepository(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			svc := service.New(f.repo, log.NewLogfmtLogger(os.Stderr))
			if res, err := svc.List(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("svc.List error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if tt.checkFunc != nil {
					tt.checkFunc(res, err)
				}
			}
		})
	}
}

func TestLoggingMiddleware_Delete(t *testing.T) {
	type fields struct {
		repo *automocks.MockTodoRepository
	}

	type args struct {
		id   string
		todo *model.TodoReq
	}

	text := "aa"
	completed := false

	tests := []struct {
		name      string
		prepare   func(f *fields)
		args      args
		wantErr   bool
		checkFunc func(err error)
	}{
		{
			name: "delete todo",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().Delete(context.Background(), gomock.Any()).Return(nil),
				)
			},
			args: args{
				"b5z2zC5c9O6~Ns_qLVmn~",
				&model.TodoReq{
					Text:      &text,
					Completed: &completed,
				}},
			wantErr: false,
		},
		{
			name: "Add todo fail",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().Delete(context.Background(), gomock.Any()).Return(sql.ErrNoRows),
				)
			},
			args: args{todo: &model.TodoReq{
				Text:      &text,
				Completed: &completed,
			}},
			wantErr: true,
			checkFunc: func(err error) {
				assert.Equal(t, err, sql.ErrNoRows, fmt.Sprintf("err: expected sql.ErrNoRows got %v", err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				repo: automocks.NewMockTodoRepository(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			svc := service.New(f.repo, log.NewLogfmtLogger(os.Stderr))
			if err := svc.Delete(context.Background(), tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("svc.Delete error = %v, wantErr %v", err, tt.wantErr)
			} else {
				if tt.checkFunc != nil {
					tt.checkFunc(err)
				}
			}
		})
	}
}

func TestLoggingMiddleware_Update(t *testing.T) {
	type fields struct {
		repo *automocks.MockTodoRepository
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
			name: "Update todo",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().Get(context.Background(), gomock.Any()).Return(&model.Todo{
						ID:        "b5z2zC5c9O6~Ns_qLVmn~",
						Text:      "aa",
						Completed: true,
					}, nil),
					f.repo.EXPECT().Update(context.Background(), gomock.Any()).Return(nil),
				)
			},
			wantErr: false,
			args: args{
				"b5z2zC5c9O6~Ns_qLVmn~",
				&model.TodoReq{
					Text:      &text,
					Completed: &completed,
				},
			},
			checkFunc: func(res *model.TodoRes, err error) {
				assert.Equal(t, err, nil, fmt.Sprintf("err: expected sql.ErrNoRows got %v", err))
			},
		},
		{
			name: "Update todo fail",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().Get(context.Background(), gomock.Any()).Return(nil, sql.ErrNoRows),
					//f.repo.EXPECT().Update(context.Background(), gomock.Any()).Return(nil),
				)
			},
			wantErr: true,
			args: args{
				"b5z2zC5c9O6~Ns_qLVmn~",
				&model.TodoReq{
					Text:      &text,
					Completed: &completed,
				},
			},
			checkFunc: func(res *model.TodoRes, err error) {
				assert.Equal(t, err, sql.ErrNoRows, fmt.Sprintf("err: expected sql.ErrNoRows got %v", err))
			},
		},
		{
			name: "Update todo fail 2",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.repo.EXPECT().Get(context.Background(), gomock.Any()).Return(&model.Todo{
						ID:        "b5z2zC5c9O6~Ns_qLVmn~",
						Text:      "aa",
						Completed: true,
					}, nil),
					f.repo.EXPECT().Update(context.Background(), gomock.Any()).Return(sql.ErrNoRows),
				)
			},
			wantErr: true,
			args: args{
				"b5z2zC5c9O6~Ns_qLVmn~",
				&model.TodoReq{
					Text:      &text,
					Completed: &completed,
				},
			},
			checkFunc: func(res *model.TodoRes, err error) {
				assert.Equal(t, err, sql.ErrNoRows, fmt.Sprintf("err: expected sql.ErrNoRows got %v", err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				repo: automocks.NewMockTodoRepository(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			svc := service.New(f.repo, log.NewLogfmtLogger(os.Stderr))
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
