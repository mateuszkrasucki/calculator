package main

import (
	"context"
	"fmt"

	"github.com/mateuszkrasucki/calculator/pkg/calculator"
)

func main() {
	var c calculator.Calculator
	{
		c = calculator.New(calculator.SimpleParse)
		c = calculator.ValidateMiddleware()(c)
	}

	result, err := c.Calculate(context.Background(), calculator.GetInput())

	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
