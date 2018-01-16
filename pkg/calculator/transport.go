package calculator

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

func decodeFormParamRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	operation := r.FormValue("operation")

	return Request{
		Operation: operation,
	}, nil
}

func decodeJSONRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Body == nil {
		return nil, NewCalculatorError(InputError, "Body cannot be empty")
	}

	decoder := json.NewDecoder(r.Body)

	var request Request

	err := decoder.Decode(&request)
	if err != nil {
		return nil, NewWrappedCalculatorError(err, InputError, err.Error())
	}

	return request, nil
}

func encodePlainResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(Response)

	w.Header().Add("Content-type", "text/plain")
	result := strconv.FormatFloat(resp.Result, 'f', -1, 64)

	_, err := fmt.Fprint(w, result)
	if err != nil {
		return NewWrappedCalculatorError(err, EncodingError, err.Error())
	}
	return nil
}

func encodeJSONResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	r := response.(Response)
	jsonResp := Response{Result: r.Result}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	err := json.NewEncoder(w).Encode(jsonResp)
	if err != nil {
		return NewWrappedCalculatorError(err, EncodingError, err.Error())
	}
	return nil
}

func encodeHTMLResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	tmpl := `<form method="post">
        <input name="operation" required> <input type="submit" value="Calculate">
        </form>
        {{ if .Result }}<h1>{{ .Operation }} = {{ .Result }}</h1>{{ end }}
        {{ if .Error }}<h1>{{ .Error }}</h1>{{ end }}`

	resp := response.(Response)
	w.Header().Add("Content-type", "text/html")

	t, err := template.New("form").Parse(tmpl)
	if err != nil {
		return NewWrappedCalculatorError(err, EncodingError, err.Error())
	}

	t.Execute(w, resp)
	return nil
}

// NewHTTPHandler creates greeter handlers
func NewHTTPHandler(endpoint endpoint.Endpoint, logger log.Logger) http.Handler {
	m := http.NewServeMux()

	m.Handle("/api/calculate", httptransport.NewServer(
		EndpointLoggingMiddleware(log.With(logger, "endpoint", "/api/calculate"))(endpoint),
		decodeJSONRequest,
		encodeJSONResponse,
	))

	m.Handle("/calculator", httptransport.NewServer(
		EndpointLoggingMiddleware(log.With(logger, "endpoint", "/calculator"))(endpoint),
		decodeFormParamRequest,
		encodeHTMLResponse,
	))

	return m
}
