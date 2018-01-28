package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"

	calculator "github.com/mateuszkrasucki/calculator/pkg/calculator"
	rpn "github.com/mateuszkrasucki/calculator/pkg/reversepolish"
)

func getInput() string {
	input, err := readStdin()
	if err != nil {
		input, err = readFlag()
	}

	if err != nil {
		return ""
	}

	return input
}

func readStdin() (string, error) {
	fi, err := os.Stdin.Stat()

	if err != nil {
		return "", err
	}

	if fi.Mode()&os.ModeNamedPipe == 0 {
		return "", errors.New("StdIn not a named pipe")
	}

	b, _, err := bufio.NewReader(os.Stdin).ReadLine()
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func readFlag() (string, error) {
	input := flag.String("c", "", "Operation to calculate")
	flag.Parse()

	return *input, nil
}

func main() {
	var c calculator.Calculator
	{
		c = calculator.New(rpn.ParseInfix)
		c = calculator.ValidateMiddleware()(c)
	}

	result, err := c.Calculate(context.Background(), getInput())

	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
