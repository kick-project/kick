package parser

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Token identifies the token type
type Token int

// Pos represents the byte position in the original input text
type Pos int

// Known options
var options = map[string]any{
	"render": struct{}{},
	"ignore": struct{}{},
}

const (
	modeLine = "kick:"
	eof      = -1

	// Tokens
	ILLEGAL Token = iota
	END

	OPTION // "render" "ignore"
	RHS    // type=
	TYPE   // type value
)

// Item represents a token returned by the scanner
type Item struct {
	Type  Token
	Value string
}

func (i Item) String() string {
	return i.Value
}

type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner
type lexer struct {
	name  string    // used only for error reports.
	input string    // the string being scanned.
	start int       // start position of this item.
	pos   int       // current position in the input.
	width int       // width of last rune read.
	items chan Item // channel of scanned item.
	lines int       // maximum number of lines to scan.
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
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
func (l *lexer) emit(t Token) {
	l.items <- Item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	lineno := strings.Count(l.input[:l.pos], "\n") + 1
	l.items <- Item{ILLEGAL, fmt.Sprintf(l.name+":"+fmt.Sprint(lineno)+":"+format, args...)}
	return nil
}

// nextItem returns the next item from the input.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) nextItem() Item {
	item := <-l.items
	return item
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for state := lexTEXT; state != nil; {
		state = state(l)
	}
	close(l.items)
}

// lex the scanner for the input string
func lex(name, input string, lines int) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan Item),
		lines: lines,
	}
	go l.run()
	return l
}

// lexTEXT scan for a mode line
func lexTEXT(l *lexer) stateFn {
	var r rune
	for {
		r = l.next()
		switch {
		case r == '\n':
			l.ignore()
			if strings.Count(l.input[:l.start], "\n") >= l.lines {
				return nil
			}
		case !strings.HasPrefix(modeLine, l.input[l.start:l.pos]):
			l.ignore()
		case l.input[l.start:l.pos] == modeLine:
			l.ignore()
			return lexMLDATA
		case r == eof:
			return nil
		}
	}
}

// lexMLDATA lex mode line data
func lexMLDATA(l *lexer) stateFn {
	var r rune

LOOP:
	for {
		r = l.next()
		switch {
		case r == eof:
			l.ignore()
			l.emit(END)
			break LOOP
		case r == ' ':
			l.ignore()
		case isRHS(l.input[l.start:l.pos]) && isRHSEnd(l.peek()):
			return lexRHS
		case isOption(l.input[l.start:l.pos]) && isOptionEnd(l.peek()):
			return lexOPTION
		case isAlphaNumeric(r):
			continue LOOP
		default:
			fmt.Printf("%c\n", r)
			l.errorf("invalid character %c", r)
		}
	}
	return nil
}

func lexOPTION(l *lexer) stateFn {
	var id = l.input[l.start:l.pos]
	_, ok := options[id]
	if ok && isOptionEnd(l.peek()) {
		l.emit(OPTION)
		return lexMLDATA
	}
	return l.errorf("modeline error: unknown option: %s", id)
}

// TODO - Add RHS lexer
func lexRHS(l *lexer) stateFn {
	var id = l.input[l.start:l.pos]
	if id == "type" && l.peek() == '=' {
		l.emit(OPTION)
		l.next()
		l.ignore()
		return lexTYPE
	}
	return nil
}

func lexTYPE(l *lexer) stateFn {
	var r rune
	var p rune
LOOP:
	for {
		r = l.next()
		p = l.peek()
		switch {
		case isType(l.input[l.start:l.pos]):
			if p == ',' {
				l.emit(TYPE)
				continue LOOP
			} else if p == ' ' {
				l.emit(TYPE)
				return lexMLDATA
			} else if p == '\n' || p == eof {
				if len(l.input[l.start:l.pos]) > 0 {
					l.emit(TYPE)
				}
				return nil
			} else if isAlphaNumeric(p) {
				continue LOOP
			}
		case r == ',':
			l.ignore()
			return lexTYPE
		default:
			break LOOP
		}
	}
	return l.errorf("modeline error: illegal character %c", r)
}

// isAlphaNumeric returns true if r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

// isAlpha returns true if r is the start of an identifier
func isAlpha(r rune) bool {
	return unicode.IsLetter(r)
}

// isRHS returns true if rhs is a valid right hand side
func isRHS(rhs string) bool {
	startsWithAlpha := isAlpha(rune(rhs[0]))
	isAlphaNum := true
	for _, c := range rhs {
		if !isAlphaNumeric(c) {
			isAlphaNum = false
			break
		}
	}
	return startsWithAlpha && isAlphaNum
}

// isOption returns true if id is a valid identifier
func isOption(id string) bool {
	startsWithAlpha := isAlpha(rune(id[0]))
	isAlphaNum := true
	for _, c := range id {
		if !isAlphaNumeric(c) {
			isAlphaNum = false
			break
		}
	}
	return startsWithAlpha && isAlphaNum
}

// isType returns true if id is a valid type
func isType(typ string) bool {
	return isOption(typ)
}

// isRHSEnd returns true if r is an equals sign
func isRHSEnd(r rune) bool {
	return r == '='
}

// isOptionEnd returns true if r is an identifier boundary
func isOptionEnd(r rune) bool {
	return r == ' ' || r == eof
}

// Parse
func Parse(file, input string, lines int) (items []Item) {
	l := lex(file, input, lines)
	for {
		item := l.nextItem()
		if item.Type == 0 || item.Type == END || item.Type == ILLEGAL {
			break
		}
		items = append(items, item)
	}
	return
}
