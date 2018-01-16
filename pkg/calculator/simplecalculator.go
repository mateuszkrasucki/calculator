package calculator

import (
	"context"
	"regexp"
	"strconv"
	"strings"
)

type simpleOperation struct {
	arg1    float64
	arg2    float64
	operand string
}

// SimpleParse provides parsing for simple two argument, one operand mathematical operations
func SimpleParse(_ context.Context, input string) (operation OperationInterface, err error) {
	matched, err := regexp.MatchString(".*[\\+\\-\\/\\*]+.*[\\+\\-\\/\\*]+.*", input)

	if err != nil {
		return nil, NewWrappedCalculatorError(err, InternalError, "Validation regex failure")
	}

	if matched {
		return nil, NewCalculatorError(ParsingError, "Operation contains more than one operation sign")
	}

	var operand string
	var args []string
	var arg1 float64
	var arg2 float64

	switch {
	case strings.Contains(input, "+"):
		operand = "+"
		args = strings.Split(input, "+")
	case strings.Contains(input, "*"):
		operand = "*"
		args = strings.Split(input, "*")
	case strings.Contains(input, "-"):
		operand = "-"
		args = strings.Split(input, "-")
	case strings.Contains(input, "/"):
		operand = "/"
		args = strings.Split(input, "/")
	default:
		return nil, NewCalculatorError(ParsingError, "Operation does not contain operation sign")
	}

	if arg1, err = strconv.ParseFloat(args[0], 64); err != nil {
		return nil, NewWrappedCalculatorError(err, ParsingError, "First operation argument could not be parsed to number")
	}

	if arg2, err = strconv.ParseFloat(args[1], 64); err != nil {
		return nil, NewWrappedCalculatorError(err, ParsingError, "Second operation argument could not be parsed to number")
	}

	return &simpleOperation{arg1: arg1, arg2: arg2, operand: operand}, nil
}

func (operation simpleOperation) calculate(_ context.Context) (result float64, err error) {
	switch operation.operand {
	case "+":
		return operation.arg1 + operation.arg2, nil
	case "*":
		return operation.arg1 * operation.arg2, nil
	case "-":
		return operation.arg1 - operation.arg2, nil
	case "/":
		return operation.arg1 / operation.arg2, nil
	default:
		return 0, NewCalculatorError(CalculationError, "Calculation error")
	}
}
