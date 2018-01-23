package calculator

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/mateuszkrasucki/calculator/pkg/errors"
)

func TestSimpleParse(t *testing.T) {
	tests := []struct {
		name              string
		operation         string
		expectedOperation *simpleOperation
		expectedError     error
	}{
		{
			"Success +",
			"2+3.1",
			&simpleOperation{
				arg1:    2,
				arg2:    3.1,
				operand: "+",
			},
			nil,
		},
		{
			"Success -",
			"3-2.1",
			&simpleOperation{
				arg1:    3,
				arg2:    2.1,
				operand: "-",
			},
			nil,
		},
		{
			"Success *",
			"3.2*2.1",
			&simpleOperation{
				arg1:    3.2,
				arg2:    2.1,
				operand: "*",
			},
			nil,
		},
		{
			"Success /",
			"3/2",
			&simpleOperation{
				arg1:    3,
				arg2:    2,
				operand: "/",
			},
			nil,
		},
		{
			"Error ++",
			"3++2",
			nil,
			errors.NewParsingError(""),
		},
		{
			"Error + - *",
			"3+2-3*4",
			nil,
			errors.NewParsingError(""),
		},
		{
			"Error no operand",
			"3",
			nil,
			errors.NewParsingError(""),
		},
		{
			"Cannot parse first number #1",
			"aa*2",
			nil,
			errors.NewParsingError(""),
		},
		{
			"Cannot parse first number #2",
			"aa+bb",
			nil,
			errors.NewParsingError(""),
		},
		{
			"Cannot parse second number",
			"2+bb",
			nil,
			errors.NewParsingError(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operation, err := Parse(context.Background(), tt.operation)

			if (tt.expectedError != nil && err == nil) || (tt.expectedError == nil && err != nil) {
				t.Errorf("expected error to be %v, got %v", tt.expectedError, err)
			}

			if tt.expectedError != nil && err != nil && !strings.Contains(err.Error(), tt.expectedError.Error()) {
				t.Errorf("expected error to be %v, got %v", tt.expectedError, err)
			}

			if tt.expectedOperation == nil {
				if operation != nil {
					t.Errorf("expected operation to be %v, got %v", tt.expectedOperation, operation)
				}
			} else if !cmp.Equal(tt.expectedOperation, operation, cmp.AllowUnexported(simpleOperation{})) {
				t.Errorf("expected operation to be %v, got %v", tt.expectedOperation, operation)
			}
		})
	}
}

func TestSimpleOperationCalculate(t *testing.T) {
	tests := []struct {
		name           string
		operation      simpleOperation
		expectedResult float64
		expectedError  error
	}{
		{
			"Success +",
			simpleOperation{
				arg1:    2,
				arg2:    3.1,
				operand: "+",
			},
			2.0 + 3.1,
			nil,
		},
		{
			"Success -",
			simpleOperation{
				arg1:    3,
				arg2:    2.5,
				operand: "-",
			},
			3.0 - 2.5,
			nil,
		},
		{
			"Success *",
			simpleOperation{
				arg1:    3,
				arg2:    2,
				operand: "*",
			},
			3 * 2,
			nil,
		},
		{
			"Success /",
			simpleOperation{
				arg1:    5,
				arg2:    2,
				operand: "/",
			},
			5.0 / 2.0,
			nil,
		},
		{
			"Error ++",
			simpleOperation{
				arg1:    5,
				arg2:    2,
				operand: "++",
			},
			0,
			errors.NewCalculationError(""),
		},
		{
			"Error no operand",
			simpleOperation{
				arg1:    5,
				arg2:    2,
				operand: "",
			},
			0,
			errors.NewCalculationError(""),
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
