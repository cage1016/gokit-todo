package responses

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	httptransport "github.com/go-kit/kit/transport/http"
	"google.golang.org/grpc/status"

	"github.com/cage1016/gokit-todo/internal/pkg/errors"
)

const (
	contentType string = "application/json"
)

func JSONErrorDecoder(r *http.Response) error {
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return fmt.Errorf("expected JSON formatted error, got Content-Type %s", contentType)
	}
	var w ErrorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}

func EncodeJSONResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if headerer, ok := response.(httptransport.Headerer); ok {
		for k, values := range headerer.Headers() {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
	}
	code := http.StatusOK
	if sc, ok := response.(httptransport.StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)
	if code == http.StatusNoContent {
		return nil
	}

	if ar, ok := response.(Responser); ok {
		return json.NewEncoder(w).Encode(ar.Response())
	}

	return json.NewEncoder(w).Encode(response)
}

func ErrorEncodeJSONResponse(f func(errorVal errors.Error) (code int)) func(_ context.Context, err error, w http.ResponseWriter) {
	return func(_ context.Context, err error, w http.ResponseWriter) {
		code := http.StatusInternalServerError
		var message string
		var errs []errors.Errors
		w.Header().Set("Content-Type", contentType)
		if s, ok := status.FromError(err); !ok {
			// HTTP
			switch errorVal := err.(type) {
			case errors.Error:
				code = f(errorVal)

				if errorVal.Msg() != "" {
					message, errs = errorVal.Msg(), errorVal.Errors()
				}
			default:
				switch err {
				case io.ErrUnexpectedEOF, io.EOF:
					code = http.StatusBadRequest
				default:
					switch err.(type) {
					case *json.SyntaxError, *json.UnmarshalTypeError:
						code = http.StatusBadRequest
					}
				}

				errs = errors.FromError(err.Error())
				message = errs[0].Message
			}
		} else {
			// GRPC
			code = HTTPStatusFromCode(s.Code())
			errs = errors.FromError(s.Message())
			message = errs[0].Message
		}

		w.WriteHeader(code)
		json.NewEncoder(w).Encode(ErrorRes{ErrorResItem{code, message, errs}})
	}
}
