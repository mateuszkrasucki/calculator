package calculator

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/google/go-cmp/cmp"

	"github.com/mateuszkrasucki/calculator/pkg/errors"
)

func TestApi(t *testing.T) {
	endpoint := func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(Request)

		switch req.Operation {
		case errors.InputError:
			return nil, errors.NewInputError(errors.InputError)
		case errors.ParsingError:
			return nil, errors.NewParsingError(errors.ParsingError)
		case errors.CalculationError:
			return nil, errors.NewCalculationError(errors.CalculationError)
		case errors.EncodingError:
			return nil, errors.NewEncodingError(errors.EncodingError)
		}

		return Response{
			Operation: req.Operation,
			Result:    0,
		}, nil
	}

	handler := NewHTTPHandler(endpoint, log.NewNopLogger())

	type respBodyStruct struct {
		Result           float64 `json:"result"`
		Error            string  `json:"error"`
		ErrorDescription string  `json:"error_description"`
	}

	tests := []struct {
		name           string
		reqBody        string
		wantStatusCode int
		wantBody       respBodyStruct
	}{
		{
			"API success",
			"{\"operation\": \"2+2\"}",
			http.StatusOK,
			respBodyStruct{Result: 0},
		},
		{
			"API InputError",
			"{\"operation\": \"InputError\"}",
			http.StatusBadRequest,
			respBodyStruct{Error: errors.InputError, ErrorDescription: errors.InputError},
		},
		{
			"API ParsingError",
			"{\"operation\": \"ParsingError\"}",
			http.StatusBadRequest,
			respBodyStruct{Error: errors.ParsingError, ErrorDescription: errors.ParsingError},
		},
		{
			"API CalculationError",
			"{\"operation\": \"CalculationError\"}",
			http.StatusInternalServerError,
			respBodyStruct{Error: errors.CalculationError, ErrorDescription: errors.CalculationError},
		},
		{
			"API EncodingError",
			"{\"operation\": \"EncodingError\"}",
			http.StatusInternalServerError,
			respBodyStruct{Error: errors.EncodingError, ErrorDescription: errors.EncodingError},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodPost, "/api/calculate", strings.NewReader(tt.reqBody))

			if err != nil {
				t.Fatalf("expected error to be nil, got '%v'", err)
			}

			// handle request
			handler.ServeHTTP(rw, req)

			if expect, got := tt.wantStatusCode, rw.Code; expect != got {
				t.Fatalf("expected '%v', got '%v'", expect, got)
			}

			var bodyDecoded respBodyStruct

			err = json.NewDecoder(rw.Body).Decode(&bodyDecoded)
			if err != nil {
				t.Fatalf("expected nil, got '%v'", err)
			}

			if !cmp.Equal(tt.wantBody, bodyDecoded) {
				t.Errorf("expected '%v', got '%v'", tt.wantBody, bodyDecoded)
			}
		})
	}
}
