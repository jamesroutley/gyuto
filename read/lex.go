package read

import (
	"bytes"
	"strings"
	"unicode/utf8"
)

type item struct {
	typ itemType
	val string
	row uint
	col uint
}

type itemType string

const eof rune = -1

const (
	itemLeftParen  itemType = "("
	itemRightParen          = ")"
	itemQuote               = "'"
	itemAtom                = "ATOM"
)

type stateFn func(*lexer) stateFn

type lexer struct {
	name     string
	input    string // The string being lexed
	start    int    // Start position of this item
	pos      int    // Current position in the input
	width    int
	items    chan item
	rowCount uint
	colCount uint
}

func lex(name, input string) (*lexer, chan item) {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
	}
	go l.run()
	return l, l.items
}

func (l *lexer) run() {
	defer close(l.items)
	for {
		switch r := l.next(); {
		case r == eof:
			return
		case r == '\n':
			l.colCount = 0
			l.rowCount++
			l.ignore()
		case isSpace(r):
			l.ignore()
		case r == '(':
			l.emit(itemLeftParen)
		case r == ')':
			l.emit(itemRightParen)
		case r == '\'':
			l.emit(itemQuote)
		case isAtom(r):
			lexAtom(l)
		default:
			panic("should not reach here")
		}
	}
}

func (l *lexer) emit(t itemType) {
	val := l.input[l.start:l.pos]
	col := l.colCount
	// If the item is multicharacter, we wish to display the column count
	// of the beginning, not the end
	if t == itemAtom {
		col = l.colCount - (uint(utf8.RuneCountInString(val)) - 1)
	}
	l.items <- item{
		typ: t,
		val: val,
		row: l.rowCount,
		col: col,
	}
	l.start = l.pos
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		// TODO: we really want to return an EOF here
		return eof
	}
	r, width := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = width
	l.pos += l.width
	l.colCount++
	return r
}

// backup steps back one rune
// It can only be called once per call of next
func (l *lexer) backup() {
	l.pos -= l.width
	l.colCount--
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func isSpace(r rune) bool {
	return strings.IndexRune(" \n\t\r", r) >= 0
}

// isAtom returns a bool indicating whether r is a valid atom rune.
// In lisp, most characters are valid for use in atoms. Things that aren't:
// - whitespace
// - parentheses
// - quotes
func isAtom(r rune) bool {
	if r == eof {
		return false
	}
	return strings.IndexRune(" \n\t\r()'", r) < 0
}

func lexAtom(l *lexer) {
	// We have already read the first rune in the atom when checking if it
	// is an atom or not. We move back a step, so we can read it again here
	l.backup()

	var buf bytes.Buffer
	r := l.next()
	for isAtom(r) {
		buf.WriteRune(r)
		r = l.next()
	}
	// We've read the first non-atom rune. Backup to 'un-read' it
	// TODO: I think a peek function might be useful here
	l.backup()

	l.emit(itemAtom)
}

// Consume the items in the channel and return a slice containing them.
// This breaks the lexer's concurrency and should only be used for debugging.
func itemChanToSlice(c chan item) (tokens []item) {
	for tok := range c {
		tokens = append(tokens, tok)
	}
	return
}
