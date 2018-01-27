package calculator

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/mateuszkrasucki/calculator/pkg/errors"
)

func TestValidationMiddleware(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	calcServiceMock := NewMockCalculator(mockCtrl)
	c := ValidateMiddleware()(calcServiceMock)

	tests := []struct {
		name          string
		input         string
		serviceCalls  int
		serviceOutput float64
		serviceError  error
		output        float64
		err           error
	}{
		{
			"Successful validation",
			"2+2",
			1,
			4.0,
			nil,
			4.0,
			nil,
		},
		{
			"Successful validation, error passed from service",
			"2+2",
			1,
			0.0,
			errors.NewCalculationError(""),
			0.0,
			errors.NewCalculationError(""),
		},
		{
			"Successful validation, instant result from empty input",
			"",
			0,
			0.0,
			nil,
			0.0,
			nil,
		},
		{
			"Failed validation, invalid character",
			"2a+2",
			0,
			0.0,
			nil,
			0.0,
			errors.NewInputError("Invalid characters in input string"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calcServiceMock.EXPECT().
				Calculate(gomock.Any(), tt.input).
				Return(tt.serviceOutput, tt.serviceError).
				Times(tt.serviceCalls)

			res, err := c.Calculate(context.Background(), tt.input)

			if (tt.err != nil && err == nil) || (tt.err == nil && err != nil) {
				t.Fatalf("expected error to be %v, got %v", tt.err, err)
			}

			if tt.err != nil && err != nil && err.Error() != tt.err.Error() {
				t.Fatalf("expected error to be %v, got %v", tt.err, err)
			}

			if tt.output != res {
				t.Errorf("expected result to be %v, got %v", tt.output, res)
			}
		})
	}
}
