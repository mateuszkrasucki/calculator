package reversepolish

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mateuszkrasucki/calculator/pkg/calculator"
	"github.com/mateuszkrasucki/calculator/pkg/errors"
	"github.com/mateuszkrasucki/calculator/pkg/lexer"
)

// NumericItem is an interface covering lexer.Item with additional method to get numeric value as float64
type NumericItem interface {
	GetType() lexer.ItemType
	GetString() string
	GetValue() float64
}

type numericItem struct {
	stringValue string
	value       float64
}

type operatorsStack struct {
	stack []lexer.Item
}

type precedenceLevel int // higher the number higher the precedence

// ParseInfix provides parsing of infix mathematical operations for postif calculator
func ParseInfix(_ context.Context, input string) (calculator.OperationInterface, error) {
	l := lexer.Lex(input)

	items := []lexer.Item{}
	opStack := &operatorsStack{stack: []lexer.Item{}}

	for i := l.NextItem(); !isEmpty(i); i = l.NextItem() {
		switch {
		case isError(i):
			return nil, errors.NewParsingError(i.GetString())
		case isNumber(i):
			numItem, err := parseNumber(i)
			if err != nil {
				return nil, err
			}
			items = append(items, numItem)
		case isMathOperator(i):
			for {
				pop := true
				topItem := opStack.peek()

				switch {
				case isMathOperator(topItem) && (getPrecedenceLevel(topItem) > getPrecedenceLevel(i)):
					//do nothing; pop = true
				case isMathOperator(topItem) && (getPrecedenceLevel(topItem) == getPrecedenceLevel(i)) && isLeftAssociative(i):
					//do nothing; pop = true
				default:
					pop = false
				}

				if pop == true {
					opStack.pop()
					items = append(items, topItem)
				}

				break
			}
			opStack.push(i)
		case isLeftBracket(i):
			opStack.push(i)
		case isRightBracket(i):
			for poppedItem := opStack.pop(); !isLeftBracket(poppedItem); poppedItem = opStack.pop() {
				if isEmpty(poppedItem) {
					return nil, errors.NewParsingError("mismatched parantheses")
				}
				items = append(items, poppedItem)
			}
		default:
			return nil, errors.NewParsingError(fmt.Sprintf("invalid item returned from lexer: %s", i))
		}
	}

	for poppedItem := opStack.pop(); !isEmpty(poppedItem); poppedItem = opStack.pop() {
		if isBracket(poppedItem) {
			return nil, errors.NewParsingError("mismatched parantheses")
		}
		items = append(items, poppedItem)
	}

	return &rpnOperation{items}, nil
}

func parseNumber(item lexer.Item) (numericItem, error) {
	if !isNumber(item) {
		return numericItem{}, errors.NewParsingError(fmt.Sprintf("could not parse %s as a number", item.GetString()))
	}

	num, err := strconv.ParseFloat(item.GetString(), 64)
	if err != nil {
		return numericItem{}, errors.NewParsingErrorWrap(err, fmt.Sprintf("could not parse %s as a number", item.GetString()))
	}

	return numericItem{item.GetString(), num}, nil
}

func (i numericItem) GetType() lexer.ItemType {
	return lexer.Number
}

func (i numericItem) GetString() string {
	return i.stringValue
}

func (i numericItem) GetValue() float64 {
	return i.value
}

func isNumber(item lexer.Item) bool {
	if item.GetType() == lexer.Number {
		return true
	}
	return false
}

func getNumericValue(item lexer.Item) float64 {
	i, ok := item.(numericItem)
	if ok {
		return i.GetValue()
	}
	return 0.0
}

func isError(item lexer.Item) bool {
	if item.GetType() == lexer.Error {
		return true
	}

	return false
}

func isEmpty(item lexer.Item) bool {
	if item.GetType() == lexer.Empty {
		return true
	}

	return false
}

func isBracket(item lexer.Item) bool {
	if isLeftBracket(item) || isRightBracket(item) {
		return true
	}

	return false
}

func isLeftBracket(item lexer.Item) bool {
	switch typ := item.GetType(); {
	case typ == lexer.LeftParenthesis:
		return true
	default:
		return false
	}
}

func isRightBracket(item lexer.Item) bool {
	switch typ := item.GetType(); {
	case typ == lexer.RightParenthesis:
		return true
	default:
		return false
	}
}

func isOperator(item lexer.Item) bool {
	switch typ := item.GetType(); {
	case typ == lexer.Number:
		return false
	case typ == lexer.Error:
		return false
	case typ == lexer.Empty:
		return false
	default:
		return true
	}
}

func isMathOperator(item lexer.Item) bool {
	switch typ := item.GetType(); {
	case typ == lexer.Number:
		return false
	case typ == lexer.LeftParenthesis:
		return false
	case typ == lexer.RightParenthesis:
		return false
	case typ == lexer.Error:
		return false
	case typ == lexer.Empty:
		return false
	default:
		return true
	}
}

func getPrecedenceLevel(item lexer.Item) precedenceLevel {
	switch typ := item.GetType(); {
	case typ == lexer.Exponent:
		return 3
	case typ == lexer.Multiplication || typ == lexer.Division:
		return 2
	case typ == lexer.Addition || typ == lexer.Subtraction:
		return 1
	default:
		return 0
	}
}

func isLeftAssociative(item lexer.Item) bool {
	if !isOperator(item) {
		return false
	}

	switch typ := item.GetType(); {
	case typ == lexer.Exponent:
		return false
	default:
		return true
	}
}

func (s *operatorsStack) pop() lexer.Item {
	if len(s.stack) == 0 {
		return lexer.NewEmptyItem()
	}
	i := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]

	return i
}

func (s *operatorsStack) peek() lexer.Item {
	if len(s.stack) == 0 {
		return lexer.NewEmptyItem()
	}

	i := s.stack[len(s.stack)-1]

	return i
}

func (s *operatorsStack) push(i lexer.Item) error {
	if !isOperator(i) {
		return errors.NewCalcError("pushing invalid item to operators stack")
	}
	s.stack = append(s.stack, i)

	return nil
}
