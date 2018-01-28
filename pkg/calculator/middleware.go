package calculator

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/mateuszkrasucki/calculator/pkg/errors"
)

// Middleware type
type Middleware func(Calculator) Calculator

// ValidateMiddleware is a validator middleware for service
func ValidateMiddleware() Middleware {
	return func(next Calculator) Calculator {
		return validateMiddleware{next}
	}
}

type validateMiddleware struct {
	next Calculator
}

func (mw validateMiddleware) Calculate(ctx context.Context, input string) (float64, error) {
	matched, err := regexp.MatchString("^[ 0-9+\\(\\)\\^\\-*\\/\\.]*$", input)
	if err != nil {
		return 0, errors.NewCalcErrorWrap(err, "Validation regex failure")
	}

	if matched == false {
		return 0, errors.NewInputError("Invalid characters in input string")
	}

	if strings.TrimSpace(input) == "" {
		return 0, nil
	}

	return mw.next.Calculate(ctx, input)
}

// ServiceLoggingMiddleware is a logging middleware for service
func ServiceLoggingMiddleware(log log.Logger) Middleware {
	return func(next Calculator) Calculator {
		return loggingMiddleware{log, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   Calculator
}

func (mw loggingMiddleware) Calculate(ctx context.Context, input string) (result float64, err error) {
	mw.logger.Log("method", "Calculate", "operation", input)
	return mw.next.Calculate(ctx, input)
}

// EndpointLoggingMiddleware returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func EndpointLoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			defer func(begin time.Time) {
				logger.Log("transport_error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)
		}
	}
}
