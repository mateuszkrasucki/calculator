package calculator

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Request definition
type Request struct {
	Operation string `json:"operation"`
}

// Response definition
type Response struct {
	Operation string  `json:"operation,omitempty"`
	Result    float64 `json:"result"`
}

// MakeEndpoint creates endpoint for calculator
func MakeEndpoint(c Calculator) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(Request)

		result, err := c.Calculate(ctx, req.Operation)

		if err != nil {
			return nil, err
		}

		return Response{
			Operation: req.Operation,
			Result:    result,
		}, nil
	}
}
