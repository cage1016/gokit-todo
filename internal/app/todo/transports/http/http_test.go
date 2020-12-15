package transports_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"github.com/stretchr/testify/assert"

	"github.com/cage1016/todo/internal/app/todo/endpoints"
	"github.com/cage1016/todo/internal/app/todo/model"
	transports "github.com/cage1016/todo/internal/app/todo/transports/http"
	automocks "github.com/cage1016/todo/internal/mocks/app/todo/service"
	test "github.com/cage1016/todo/test/util"
)

func TestAddHandler(t *testing.T) {
	type fields struct {
		svc *automocks.MockTodoService
	}
	type args struct {
		method, url string
		body        string
	}

	tests := []struct {
		name      string
		prepare   func(f *fields)
		args      args
		wantErr   bool
		checkFunc func(res *http.Response, err error, body []byte)
	}{
		{
			name: "add todo",
			prepare: func(f *fields) {
				gomock.InOrder(
					f.svc.EXPECT().Add(gomock.Any(), gomock.Any()).Return(model.Todo{
						ID:        "iKe0KxpurIn0E_6vzUDAr",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						Text:      "aa",
						Complete:  false,
					}, nil),
				)
			},
			wantErr: false,
			args: args{
				method: http.MethodPost,
				url:    "/items",
				body:   `{"text":"aa"}`,
			},
			checkFunc: func(res *http.Response, err error, body []byte) {
				assert.Nil(t, err, fmt.Sprintf("unexpected error %s", err))
				assert.Equal(t, http.StatusOK, res.StatusCode, fmt.Sprintf("status should be 204: got %d", res.StatusCode))
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

			eps := endpoints.New(f.svc, logger, tracer, zkt)
			ts := httptest.NewServer(transports.NewHTTPHandler(eps, tracer, zkt, logger))
			defer ts.Close()

			req := test.TestRequest{
				Client:      ts.Client(),
				Method:      tt.args.method,
				URL:         fmt.Sprintf("%s%s", ts.URL, tt.args.url),
				ContentType: "application/json",
				Body:        strings.NewReader(tt.args.body),
			}

			if res, err := req.Make(); (err != nil) != tt.wantErr {
				t.Errorf("%s: unexpected error %s", tt.name, err)
			} else {
				body, _ := ioutil.ReadAll(res.Body)
				if tt.checkFunc != nil {
					tt.checkFunc(res, err, body)
				}
			}
		})
	}
}
