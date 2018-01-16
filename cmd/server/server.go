package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"

	"github.com/mateuszkrasucki/calculator/pkg/calculator"
)

func main() {
	addr := flag.String("addr", ":8080", "Interface and port to listen on")
	flag.Parse()

	// Create a single logger, which we'll use and give to other components.
	var logger log.Logger
	{
		logger = log.NewJSONLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// Create calculator service
	var c calculator.Calculator
	{
		c = calculator.New(calculator.SimpleParse)
		c = calculator.ServiceLoggingMiddleware(logger)(c)
		c = calculator.ValidateMiddleware()(c)
	}

	endpoint := calculator.MakeEndpoint(c)
	handler := calculator.NewHTTPHandler(endpoint, logger)

	logger.Log("transport", "http", "listen", *addr)
	err := http.ListenAndServe(*addr, handler)
	if err != nil {
		logger.Log("transport", "http", "during", "listen", "err", err)
	}
}
