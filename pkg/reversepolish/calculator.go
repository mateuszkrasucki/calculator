package reversepolish

import (
	"context"
	"fmt"

	"github.com/mateuszkrasucki/calculator/pkg/errors"
	"github.com/mateuszkrasucki/calculator/pkg/lexer"
	"github.com/mateuszkrasucki/calculator/pkg/simplecalculator"
)

type rpnOperation struct {
	items []lexer.Item
}

type numericStack struct {
	stack []float64
}

func (o rpnOperation) Calculate(ctx context.Context) (float64, error) {
	stack := numericStack{[]float64{}}

	for _, i := range o.items {
		switch {
		case isMathOperator(i):
			if stack.length() < 2 {
				return 0.0, errors.NewCalculationError("not enough operands on stack")
			}

			operand2 := stack.pop()
			operand1 := stack.pop()
			simpleOp := simplecalculator.NewOperation(i.GetString(), operand1, operand2)
			r, err := simpleOp.Calculate(ctx)
			if err != nil {
				return 0.0, errors.NewCalculationErrorWrap(err, fmt.Sprintf("failed calculating simple operation %f %s %f", operand1, i.GetString(), operand2))
			}

			stack.push(r)
		case isNumber(i):
			r := i.(numericItem)
			stack.push(r.GetValue())
		default:
			return 0.0, errors.NewCalculationError(fmt.Sprintf("invalid item in the RPN operation: %s", i.GetString()))
		}
	}
	if stack.length() != 1 {
		return 0.0, errors.NewCalculationError("too many operands on the stack at the end of calculation")
	}

	return stack.pop(), nil
}

func (s *numericStack) length() int {
	return len(s.stack)
}

func (s *numericStack) pop() float64 {
	if len(s.stack) == 0 {
		return 0.0
	}

	i := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]

	return i
}

func (s *numericStack) push(i float64) {
	s.stack = append(s.stack, i)
}
