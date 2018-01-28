package reversepolish

import (
	"context"
	"strings"
	"testing"

	//"github.com/google/go-cmp/cmp"

	"github.com/mateuszkrasucki/calculator/pkg/errors"
	"github.com/mateuszkrasucki/calculator/pkg/lexer"
)

func TestReversePolishCalculate(t *testing.T) {
	tests := []struct {
		name           string
		operation      rpnOperation
		expectedResult float64
		expectedError  error
	}{
		{
			"Success +",
			rpnOperation{
				[]lexer.Item{
					numericItem{"15", 15.0},
					numericItem{"7", 7.0},
					numericItem{"1", 1.0},
					numericItem{"1", 1.0},
					lexer.NewItem(lexer.Addition, "+"),
					lexer.NewItem(lexer.Subtraction, "-"),
					lexer.NewItem(lexer.Division, "/"),
					numericItem{"3", 3.0},
					lexer.NewItem(lexer.Multiplication, "*"),
					numericItem{"2", 2.0},
					numericItem{"1", 1.0},
					numericItem{"1", 1.0},
					lexer.NewItem(lexer.Addition, "+"),
					lexer.NewItem(lexer.Addition, "+"),
					lexer.NewItem(lexer.Subtraction, "-"),
				},
			},
			5.0,
			nil,
		},
		{
			"Error no operands on stack",
			rpnOperation{
				[]lexer.Item{
					numericItem{"1", 1.0},
					lexer.NewItem(lexer.Addition, "+"),
					numericItem{"1", 1.0},
					lexer.NewItem(lexer.Addition, "+"),
				},
			},
			0.0,
			errors.NewCalculationError("not enough operands on stack"),
		},
		{
			"Error invalid item",
			rpnOperation{
				[]lexer.Item{
					numericItem{"1", 1.0},
					numericItem{"1", 1.0},
					lexer.NewItem(lexer.LeftParenthesis, "("),
				},
			},
			0.0,
			errors.NewCalculationError("invalid item in the RPN operation: ("),
		},
		{
			"Error too many operands, too litle operations",
			rpnOperation{
				[]lexer.Item{
					numericItem{"1", 1.0},
					numericItem{"1", 1.0},
					numericItem{"1", 1.0},
					lexer.NewItem(lexer.Addition, "+"),
				},
			},
			0.0,
			errors.NewCalculationError("too many operands on the stack at the end of calculation"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.operation.Calculate(context.Background())

			if (tt.expectedError != nil && err == nil) || (tt.expectedError == nil && err != nil) {
				t.Fatalf("expected error to be %v, got %v", tt.expectedError, err)
			}

			if tt.expectedError != nil && err != nil && !strings.Contains(err.Error(), tt.expectedError.Error()) {
				t.Fatalf("expected error to be %v, got %v", tt.expectedError, err)
			}

			if tt.expectedResult != result {
				t.Errorf("expected result to be %v, got %v", tt.expectedResult, result)
			}
		})
	}
}
