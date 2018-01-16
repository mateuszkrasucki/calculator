package calculator

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// Error categories contants
const (
	InputError       = "InputError"
	ParsingError     = "ParsingError"
	CalculationError = "CalculationError"
	EncodingError    = "EncodingError"
	InternalError    = "InternalError"
)

var statusCodeDict = map[string]int{
	InputError:       http.StatusBadRequest,
	ParsingError:     http.StatusBadRequest,
	CalculationError: http.StatusInternalServerError,
	EncodingError:    http.StatusInternalServerError,
	InternalError:    http.StatusInternalServerError,
}

type calculatorError struct {
	error
	errorCategory    string
	errorDescription string
}

// NewCalculatorError returns new calculator error
func NewCalculatorError(errorCategory string, errorDescription string) error {
	return calculatorError{
		errors.New(errorCategory),
		errorCategory,
		errorDescription,
	}
}

// NewWrappedCalculatorError returns new calculator error with another error wrapped
func NewWrappedCalculatorError(err error, errorCategory string, errorDescription string) error {
	return calculatorError{
		errors.Wrap(err, errorCategory),
		errorCategory,
		errorDescription,
	}
}

// MarshallJSON returns error as a JSON string
func (e calculatorError) MarshalJSON() ([]byte, error) {
	errorRespStruct := struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description,omitempty"`
	}{Error: e.errorCategory, ErrorDescription: e.errorDescription}

	return json.Marshal(errorRespStruct)
}

// StatusCode returns HTTP status code appropriate for the error
func (e calculatorError) StatusCode() int {
	return statusCodeDict[e.errorCategory]
}
