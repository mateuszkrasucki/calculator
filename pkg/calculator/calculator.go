package calculator

import (
	"context"
)

// OperationInterface represents parsed mathematical operation that can be calculated
type OperationInterface interface {
	calculate(context.Context) (float64, error)
}

type parser func(context.Context, string) (OperationInterface, error)

// Calculator interface, accepts context and string represeting mathematical operation to be calculated
type Calculator interface {
	Calculate(context.Context, string) (float64, error)
}

type calculator struct {
	parse parser
}

// New returns new Calculator with provided parsing function
func New(parsingFunc parser) Calculator {
	return calculator{parse: parsingFunc}
}

// Calculate result of mathemtical operation passed as string
func (c calculator) Calculate(ctx context.Context, input string) (result float64, err error) {
	operation, err := c.parse(ctx, input)
	if err != nil {
		return 0, err
	}

	res, err := operation.calculate(ctx)
	return res, err
}
