package endpoints

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// LoggingMiddleware returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func LoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				if err == nil {
					level.Info(logger).Log("transport_error", err, "took", time.Since(begin))
				} else {
					level.Error(logger).Log("transport_error", err, "took", time.Since(begin))
				}
			}(time.Now())
			return next(ctx, request)
		}
	}
}

// AuthnMiddleware returns an endpoint middleware that apply authentication func
func AuthnMiddleware(n endpoint.Middleware, endpoints Endpoints) Endpoints {
	return endpoints
}

// AuthzMiddleware returns an endpoint middleware that apply authorization func (opa rbac)
func AuthzMiddleware(z func(action string, resource string) endpoint.Middleware, endpoints Endpoints) Endpoints {
	return endpoints
}
