// +build !integration

package service_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/cage1016/todo/internal/app/todo/model"
	"github.com/cage1016/todo/internal/app/todo/service"
	automocks "github.com/cage1016/todo/internal/mocks/app/todo/model"
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
