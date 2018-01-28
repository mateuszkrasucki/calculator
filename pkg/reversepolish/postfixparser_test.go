package reversepolish

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/mateuszkrasucki/calculator/pkg/errors"
	"github.com/mateuszkrasucki/calculator/pkg/lexer"
)

func TestParseInfix(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedError error
		expected      []lexer.Item
	}{
		{
			"Success example from wiki",
			"3 + 4 * 2 / ( 1 - 5 ) ^  2 ^ 3",
			nil,
			[]lexer.Item{
				numericItem{"3", 3.0},
				numericItem{"4", 4.0},
				numericItem{"2", 2.0},
				lexer.NewItem(lexer.Multiplication, "*"),
				numericItem{"1", 1.0},
				numericItem{"5", 5.0},
				lexer.NewItem(lexer.Subtraction, "-"),
				numericItem{"2", 2.0},
				numericItem{"3", 3.0},
				lexer.NewItem(lexer.Exponent, "^"),
				lexer.NewItem(lexer.Exponent, "^"),
				lexer.NewItem(lexer.Division, "/"),
				lexer.NewItem(lexer.Addition, "+"),
			},
		},
		{
			"Error from lexer",
			"2+2..2",
			errors.NewParsingError("invalid rune at: 4; could not lex: 2.."),
			nil,
		},
		{
			"Mismatched parantheses #1",
			"2+(2*3*5",
			errors.NewParsingError("mismatched parantheses"),
			nil,
		},
		{
			"Mismatched parantheses #2",
			"2+(2*3*5))",
			errors.NewParsingError("mismatched parantheses"),
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := ParseInfix(context.Background(), tt.input)

			if (tt.expectedError != nil && err == nil) || (tt.expectedError == nil && err != nil) {
				t.Fatalf("expected error to be %v, got %v", tt.expectedError, err)
			}

			if tt.expectedError != nil && err != nil && !strings.Contains(err.Error(), tt.expectedError.Error()) {
				t.Fatalf("expected error to be %v, got %v", tt.expectedError, err)
			}

			if (tt.expected == nil && l != nil) || (l == nil && tt.expected != nil) {
				t.Fatalf("expected result to be %v, got %v", tt.expected, l)
			}

			if l != nil && tt.expected != nil {
				result, ok := l.(*rpnOperation)
				if !ok {
					t.Fatal("expected rpnOperations struct, casting failed")
				}

				if len(result.items) != len(tt.expected) {
					t.Fatalf("expected result to be %v, got %v", tt.expected, result.items)
				}

				for k := range result.items {
					if result.items[k].GetType() != tt.expected[k].GetType() || result.items[k].GetString() != tt.expected[k].GetString() || !cmp.Equal(getNumericValue(result.items[k]), getNumericValue(tt.expected[k])) {
						t.Fatalf("expected operation item to be %v at position %d, got %v", tt.expected[k], k, result.items[k])
					}
				}
			}
		})
	}

}
