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
)

func TestApi(t *testing.T) {
	endpoint := func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(Request)

		switch req.Operation {
		case InputError:
			return nil, NewCalculatorError(InputError, InputError)
		case ParsingError:
			return nil, NewCalculatorError(ParsingError, ParsingError)
		case CalculationError:
			return nil, NewCalculatorError(CalculationError, CalculationError)
		case EncodingError:
			return nil, NewCalculatorError(EncodingError, EncodingError)
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
			respBodyStruct{Error: InputError, ErrorDescription: InputError},
		},
		{
			"API ParsingError",
			"{\"operation\": \"ParsingError\"}",
			http.StatusBadRequest,
			respBodyStruct{Error: ParsingError, ErrorDescription: ParsingError},
		},
		{
			"API CalculationError",
			"{\"operation\": \"CalculationError\"}",
			http.StatusInternalServerError,
			respBodyStruct{Error: CalculationError, ErrorDescription: CalculationError},
		},
		{
			"API EncodingError",
			"{\"operation\": \"EncodingError\"}",
			http.StatusInternalServerError,
			respBodyStruct{Error: EncodingError, ErrorDescription: EncodingError},
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
