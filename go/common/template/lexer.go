package template

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type item struct {
	typ itemType // The type of this item.
	val string   // The key value of this item.
}

// itemType identifies the type of lex items.
type itemType int

const (
	itemError itemType = iota
	itemEOF
	itemIdentifier
	itemText
	itemLeftDelim
	itemRightDelim

	leftDelim  = "{{"
	rightDelim = "}}"

	eof = -1
)

// lexStateFn represents the state of the scanner as a function that returns the next state.
type lexStateFn func(*lexer) lexStateFn

// lexer holds the state of the scanner.
type lexer struct {
	input      string // the string being scanned
	pos        int    // current position in the input
	start      int    // start position of this item
	width      int    // width of last rune read from input
	leftDelim  string
	rightDelim string
	items      chan item // channel of scanned items
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

// acceptUntil consumes a run of runes until one from the set is met.
func (l *lexer) acceptUntil(stopper string) {
	for !strings.ContainsRune(stopper, l.next()) {
	}
	l.backup()
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) ignoreSpace() {
	l.acceptRun(" ")
	l.ignore()
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) lexStateFn {
	l.items <- item{itemError, fmt.Sprintf(format, args...)}
	return nil
}

// nextItem returns the next item from the input.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) nextItem() item {
	return <-l.items
}

// lex creates a new scanner for the input string.
func lex(input, leftDelim, rightDelim string) *lexer {
	l := &lexer{
		input:      input,
		items:      make(chan item),
		leftDelim:  leftDelim,
		rightDelim: rightDelim,
	}
	go l.run()
	return l
}

// run the state machine for the lexer.
func (l *lexer) run() {
	defer close(l.items)
	for state := lexText; state != nil; {
		state = state(l)
	}
}

func lexText(l *lexer) lexStateFn {
	l.width = 0
	if x := strings.Index(l.input[l.pos:], l.leftDelim); x >= 0 {
		l.pos += x
		if l.pos > l.start {
			l.emit(itemText)
		}
		return lexLeftDelim
	}
	l.pos = len(l.input)
	// Correctly reached EOF.
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.emit(itemEOF)
	return nil
}

func lexLeftDelim(l *lexer) lexStateFn {
	l.pos += len(l.leftDelim)
	l.ignore()
	l.emit(itemLeftDelim)
	return lexIdentifier
}

// lexIdentifier scans an alphanumeric.
func lexIdentifier(l *lexer) lexStateFn {
	l.ignoreSpace()
	for {
		switch r := l.next(); {
		case isValidIdentifierRune(r):
			// absorb.
		default:
			l.backup()
			if l.pos == l.start {
				return l.errorf("expected identifier")
			}
			l.emit(itemIdentifier)
			l.ignoreSpace()
			x := len(rightDelim)
			if l.input[l.pos:l.pos+x] == rightDelim {
				return lexRightDelim
			}
			l.errorf("expected right delimiter")
			return nil
		}
	}
}

func lexRightDelim(l *lexer) lexStateFn {
	l.pos += len(l.rightDelim)
	l.ignore()
	l.emit(itemRightDelim)
	return lexText
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isValidIdentifierRune(r rune) bool {
	return strings.ContainsRune("_-[].", r) || unicode.IsLetter(r) || unicode.IsDigit(r)
}
