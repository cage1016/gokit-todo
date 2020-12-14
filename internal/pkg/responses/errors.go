package responses

import (
	"github.com/cage1016/todo/internal/pkg/errors"
)

type ErrorWrapper struct {
	Error string `json:"error"`
}

type ErrorResItem struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Errors  []errors.Errors `json:"errors"`
}

type ErrorRes struct {
	Error ErrorResItem `json:"error"`
}
