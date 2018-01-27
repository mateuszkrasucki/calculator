package lexer

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// ItemType type
type ItemType int

// Item interface
type Item interface {
	GetType() ItemType
	GetString() string
}

type item struct {
	typ   ItemType
	value string
}

// Lexer interface
type Lexer interface {
	NextItem() Item
}

type lexer struct {
	input    string
	start    int       // start position of current item
	pos      int       // current position of scanning in th input
	lastStep int       // length of last step
	items    chan Item // channel of scanned items
}

const eof = -1 // rune value used when reached end of string

// ItemType constants
const (
	Empty ItemType = iota
	Number
	LeftParenthesis
	RightParenthesis
	Addition
	Subtraction
	Multiplication
	Division
	Exponent
	Error
)

type stateFn func(*lexer) stateFn

// Lex returns lexer
func Lex(input string) Lexer {
	l, _ := lex(input)

	return l
}

// NextItem returns next lexed item
func (l *lexer) NextItem() Item {
	i, ok := <-l.items
	if ok {
		return i
	}

	return NewEmptyItem()
}

func lex(input string) (*lexer, chan Item) {
	l := &lexer{
		input: input,
		items: make(chan Item),
	}

	go l.run()

	return l, l.items
}

func (l *lexer) run() {
	for state := lexUnknown; state != nil; {
		state = state(l)
	}

	close(l.items)
}

func (l *lexer) emit(t ItemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
	l.lastStep = 0
}

func (l *lexer) emitError() {
	l.items <- item{
		Error,
		fmt.Sprintf("invalid rune at: %d; could not lex: %s", l.pos-1, l.input[l.start:l.pos]),
	}
	l.start = l.pos
	l.lastStep = 0
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.pos++
		l.lastStep = 1
		return eof
	}

	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += w
	l.lastStep = w

	return r
}

func (l *lexer) stepBack() {
	l.pos -= l.lastStep
	l.lastStep = 0
}

func (l *lexer) skip() {
	l.start = l.pos
	l.lastStep = 0
}

func lexUnknown(l *lexer) stateFn {
	switch r := l.next(); {
	case r == eof:
		return nil
	case unicode.IsDigit(r):
		return lexNumber
	case unicode.IsSpace(r):
		l.skip()
	case r == '(':
		l.emit(LeftParenthesis)
	case r == ')':
		l.emit(RightParenthesis)
	case r == '+':
		l.emit(Addition)
	case r == '-':
		l.emit(Subtraction)
	case r == '*':
		l.emit(Multiplication)
	case r == '/':
		l.emit(Division)
	case r == '^':
		l.emit(Exponent)
	default:
		l.emitError()
		return nil
	}

	return lexUnknown
}

func lexNumber(l *lexer) stateFn {
	dotCounter := 0

	for r := l.next(); r != eof && isPartOfNumber(r); r = l.next() {
		if r == '.' {
			dotCounter++
			if dotCounter > 1 {
				l.emitError()
				return nil
			}
		}
	}

	l.stepBack()
	l.emit(Number)

	return lexUnknown
}

func isPartOfNumber(r rune) bool {
	return r == '.' || unicode.IsDigit(r)
}

// NewItem returns lexer item
func NewItem(t ItemType, val string) Item {
	return item{t, val}
}

// NewEmptyItem returns empty lexer item
func NewEmptyItem() Item {
	return item{Empty, ""}
}

func (i item) GetType() ItemType {
	return i.typ
}

func (i item) GetString() string {
	return i.value
}
