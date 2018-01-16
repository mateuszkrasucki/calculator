package calculator

import (
	"bufio"
	"errors"
	"flag"
	"os"
)

// GetInput reads input from Stdin and if not available tries to read from execution flag
func GetInput() string {
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
