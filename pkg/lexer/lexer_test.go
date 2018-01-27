package lexer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Item
	}{
		{
			"Success",
			"1+2*(3^2/0.5534)-   5.0",
			[]Item{
				item{Number, "1"},
				item{Addition, "+"},
				item{Number, "2"},
				item{Multiplication, "*"},
				item{LeftParenthesis, "("},
				item{Number, "3"},
				item{Exponent, "^"},
				item{Number, "2"},
				item{Division, "/"},
				item{Number, "0.5534"},
				item{RightParenthesis, ")"},
				item{Subtraction, "-"},
				item{Number, "5.0"},
			},
		},
		{
			"Error, cannot start with .",
			".5534-5.0",
			[]Item{
				item{Error, "invalid rune at: 0; could not lex: ."},
			},
		},
		{
			"Error, number cannot have two dots #1",
			"5.55.34-5.0",
			[]Item{
				item{Error, "invalid rune at: 4; could not lex: 5.55."},
			},
		},
		{
			"Error, number cannot have two dots #1",
			"5..5534-5.0",
			[]Item{
				item{Error, "invalid rune at: 2; could not lex: 5.."},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Lex(tt.input)

			result := []Item{}

			for i := l.NextItem(); i.GetType() != Empty; i = l.NextItem() {
				result = append(result, i)
			}

			if len(result) != len(tt.expected) {
				t.Errorf("expected: %v items, got: %v items", len(tt.expected), len(result))
			}

			if !cmp.Equal(tt.expected, result, cmp.AllowUnexported(item{})) {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}

}
