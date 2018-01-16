package calculator

import (
	"context"
	"errors"
	"strings"
	"testing"
)

type mockOperation struct {
	Operation string
}

func (o mockOperation) calculate(_ context.Context) (result float64, err error) {
	if o.Operation == CalculationError {
		return 0, errors.New(CalculationError)
	}

	return float64(len(o.Operation)), nil
}

func mockParser(_ context.Context, operation string) (OperationInterface, error) {
	return &mockOperation{Operation: operation}, nil
}

func mockParserError(_ context.Context, operation string) (OperationInterface, error) {
	return nil, errors.New(operation)
}

func TestCalculate(t *testing.T) {
	c := New(mockParser)

	tests := []struct {
		name           string
		operation      string
		expectedResult float64
		expectedError  error
	}{
		{
			"Success",
			"2+2",
			3,
			nil,
		},
		{
			"Error",
			CalculationError,
			0,
			NewCalculatorError(CalculationError, ""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := c.Calculate(context.Background(), tt.operation)

			if (tt.expectedError != nil && err == nil) || (tt.expectedError == nil && err != nil) {
				t.Fatalf("expected error to be %v, got %v", tt.expectedError, err)
			}

			if tt.expectedError != nil && err != nil && !strings.Contains(err.Error(), tt.expectedError.Error()) {
				t.Fatalf("expected error to be %v, got %v", tt.expectedError, err)
			}

			if result != tt.expectedResult {
				t.Errorf("expected result to be %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestParsingError(t *testing.T) {
	expectedError := errors.New("operation")
	expectedResult := 0.0

	c := New(mockParserError)

	result, err := c.Calculate(context.Background(), "operation")

	if (expectedError != nil && err == nil) || (expectedError == nil && err != nil) {
		t.Fatalf("expected error to be %v, got %v", expectedError, err)
	}

	if expectedError != nil && err != nil && !strings.Contains(err.Error(), expectedError.Error()) {
		t.Fatalf("expected error to be %v, got %v", expectedError, err)
	}

	if result != expectedResult {
		t.Errorf("expected result to be %v, got %v", expectedResult, result)
	}
}
