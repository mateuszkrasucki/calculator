package errors

import (
	"encoding/json"
	"net/http"

	pkgerrors "github.com/pkg/errors"
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

type calcError struct {
	error
	category    string
	description string
}

// NewCalcError returns new calculator error
func NewCalcError(description string) error {
	return newCalcErrorCategorized(InternalError, description)
}

func newCalcErrorCategorized(category string, description string) error {
	return calcError{
		pkgerrors.New(category),
		category,
		description,
	}
}

// NewCalcErrorWrap returns new calculator error with another error wrapped
func NewCalcErrorWrap(err error, description string) error {
	return newCalcErrorWrapCategorized(err, InternalError, description)
}

func newCalcErrorWrapCategorized(err error, category string, description string) error {
	return calcError{
		pkgerrors.Wrap(err, category),
		category,
		description,
	}
}

// NewInputError returns new calculator error
func NewInputError(description string) error {
	return newCalcErrorCategorized(InputError, description)
}

// NewInputErrorWrap returns new calculator error with another error wrapped
func NewInputErrorWrap(err error, description string) error {
	return newCalcErrorWrapCategorized(err, InputError, description)
}

// NewParsingError returns new calculator error
func NewParsingError(description string) error {
	return newCalcErrorCategorized(ParsingError, description)
}

// NewParsingErrorWrap returns new calculator error with another error wrapped
func NewParsingErrorWrap(err error, description string) error {
	return newCalcErrorWrapCategorized(err, ParsingError, description)
}

// NewCalculationError returns new calculator error
func NewCalculationError(description string) error {
	return newCalcErrorCategorized(CalculationError, description)
}

// NewCalculationErrorWrap returns new calculator error with another error wrapped
func NewCalculationErrorWrap(err error, description string) error {
	return newCalcErrorWrapCategorized(err, CalculationError, description)
}

// NewEncodingError returns new calculator error
func NewEncodingError(description string) error {
	return newCalcErrorCategorized(EncodingError, description)
}

// NewEncodingErrorWrap returns new calculator error with another error wrapped
func NewEncodingErrorWrap(err error, description string) error {
	return newCalcErrorWrapCategorized(err, EncodingError, description)
}

// MarshallJSON returns error as a JSON string
func (e calcError) MarshalJSON() ([]byte, error) {
	errorRespStruct := struct {
		Error       string `json:"error"`
		Description string `json:"error_description,omitempty"`
	}{Error: e.category, Description: e.description}

	return json.Marshal(errorRespStruct)
}

// StatusCode returns HTTP status code appropriate for the error
func (e calcError) StatusCode() int {
	return statusCodeDict[e.category]
}
