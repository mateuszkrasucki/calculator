package calculator

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/mateuszkrasucki/calculator/pkg/calculator"
	"github.com/mateuszkrasucki/calculator/pkg/errors"
)

type simpleOperation struct {
	arg1     float64
	arg2     float64
	operator string
}

// Parse provides parsing for simple two argument, one operator mathematical operations
func Parse(_ context.Context, input string) (calculator.OperationInterface, error) {
	matched, err := regexp.MatchString(".*[\\+\\-\\/\\*]+.*[\\+\\-\\/\\*]+.*", input)

	if err != nil {
		return nil, errors.NewCalcErrorWrap(err, "Validation regex failure")
	}

	if matched {
		return nil, errors.NewParsingError("Operation contains more than one operation sign")
	}

	var operator string
	var args []string
	var arg1 float64
	var arg2 float64

	switch {
	case strings.Contains(input, "+"):
		operator = "+"
		args = strings.Split(input, "+")
	case strings.Contains(input, "*"):
		operator = "*"
		args = strings.Split(input, "*")
	case strings.Contains(input, "-"):
		operator = "-"
		args = strings.Split(input, "-")
	case strings.Contains(input, "/"):
		operator = "/"
		args = strings.Split(input, "/")
	default:
		return nil, errors.NewParsingError("Operation does not contain operation sign")
	}

	if arg1, err = strconv.ParseFloat(args[0], 64); err != nil {
		return nil, errors.NewParsingErrorWrap(err, "First operation argument could not be parsed to number")
	}

	if arg2, err = strconv.ParseFloat(args[1], 64); err != nil {
		return nil, errors.NewParsingErrorWrap(err, "Second operation argument could not be parsed to number")
	}

	return &simpleOperation{arg1: arg1, arg2: arg2, operator: operator}, nil
}

func (operation *simpleOperation) Calculate(_ context.Context) (result float64, err error) {
	switch operation.operator {
	case "+":
		return operation.arg1 + operation.arg2, nil
	case "*":
		return operation.arg1 * operation.arg2, nil
	case "-":
		return operation.arg1 - operation.arg2, nil
	case "/":
		return operation.arg1 / operation.arg2, nil
	default:
		return 0, errors.NewCalculationError("Calculation error")
	}
}
