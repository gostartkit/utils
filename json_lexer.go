package utils

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

const (
	_QuoteChar    = '"'
	_CommaChar    = ','
	_ColonChar    = ':'
	_BracketLeft  = '['
	_BracketRight = ']'
	_BraceLeft    = '{'
	_BraceRight   = '}'
)

var jsonLexerPool = sync.Pool{
	New: func() any {
		return &JSONLexer{}
	},
}

func CreateJSONLexer(data []byte) *JSONLexer {
	lexer := jsonLexerPool.Get().(*JSONLexer)
	lexer.data = data
	lexer.len = len(data)
	lexer.pos = 0
	lexer.line = 1
	lexer.column = 1
	return lexer
}

func ReleaseJSONLexer(lexer *JSONLexer) {
	lexer.data = nil
	lexer.len = 0
	lexer.pos = 0
	lexer.line = 1
	lexer.column = 1
	jsonLexerPool.Put(lexer)
}

type JSONLexer struct {
	data   []byte
	len    int
	pos    int
	line   int
	column int
}

// Position returns current line and column for error reporting
func (l *JSONLexer) Position() (int, int) {
	return l.line, l.column
}

// Peek returns the next byte in the input without advancing the lexer
// If the lexer is at the end of the input, it returns 0
// This is useful for lookahead when parsing JSON
func (l *JSONLexer) Peek() byte {
	if l.pos < l.len {
		return l.data[l.pos]
	}
	return 0
}

// Advance moves the lexer to the next byte in the input
// If the lexer is at the end of the input, it does nothing
// This is useful for advancing the lexer after reading a byte
func (l *JSONLexer) Advance() {
	if l.pos < l.len {
		if l.data[l.pos] == '\n' {
			l.line++
			l.column = 1
		} else {
			l.column++
		}
		l.pos++
	}
}

// SkipWhitespace advances the lexer until it reaches a non-whitespace character
// This is useful for skipping over whitespace between JSON tokens
func (l *JSONLexer) SkipWhitespace() {
	for l.pos < l.len {
		switch l.data[l.pos] {
		case ' ', '\r', '\n', '\t':
			l.Advance()
		default:
			return
		}
	}
}

// Expect checks if the next byte in the input is equal to the given byte
// If it is, the lexer advances to the next byte
// If it is not, an error is returned
// This is useful for checking if the next byte is a specific character
func (l *JSONLexer) Expect(c byte) error {
	l.SkipWhitespace()
	if l.pos >= l.len {
		return fmt.Errorf("line %d, column %d: expected '%c', got EOF", l.line, l.column, c)
	}
	if l.data[l.pos] != c {
		return fmt.Errorf("line %d, column %d: expected '%c', got '%c'", l.line, l.column, c, l.Peek())
	}
	l.Advance()
	return nil
}

// ReadString reads a JSON string from the input
// and returns it as a string
// The string may contain escaped characters
func (l *JSONLexer) ReadString() (string, error) {

	if err := l.Expect(_QuoteChar); err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.Grow(32)

	for l.pos < l.len {

		ch := l.data[l.pos]

		if ch == _QuoteChar {
			l.Advance()
			return sb.String(), nil
		}

		if ch == '\\' {
			l.Advance()
			if l.pos >= l.len {
				return "", fmt.Errorf("line %d, column %d: unterminated string escape", l.line, l.column)
			}

			switch l.data[l.pos] {
			case _QuoteChar:
				sb.WriteByte(_QuoteChar)
			case '\\':
				sb.WriteByte('\\')
			case '/':
				sb.WriteByte('/')
			case 'b':
				sb.WriteByte('\b')
			case 'f':
				sb.WriteByte('\f')
			case 'n':
				sb.WriteByte('\n')
			case 'r':
				sb.WriteByte('\r')
			case 't':
				sb.WriteByte('\t')
			case 'u':
				l.Advance()
				if l.pos+4 > l.len {
					return "", fmt.Errorf("line %d, column %d: incomplete unicode escape", l.line, l.column)
				}
				hex := string(l.data[l.pos : l.pos+4])
				val, err := strconv.ParseUint(hex, 16, 16)
				if err != nil {
					return "", fmt.Errorf("line %d, column %d: invalid unicode escape: %s", l.line, l.column, hex)
				}
				sb.WriteRune(rune(val))
				l.pos += 3
				l.column += 3
			default:
				return "", fmt.Errorf("line %d, column %d: invalid escape character: \\%c", l.line, l.column, l.data[l.pos])
			}
		} else if ch < 0x20 {
			return "", fmt.Errorf("line %d, column %d: unescaped control character", l.line, l.column)
		} else {
			sb.WriteByte(ch)
		}
		l.Advance()
	}

	return "", fmt.Errorf("line %d, column %d: unterminated string", l.line, l.column)
}

// ReadNumber reads a JSON number from the input
// and returns it as a float64
// The number may be an integer or a floating point number
func (l *JSONLexer) ReadNumber() (float64, error) {
	l.SkipWhitespace()
	start := l.pos

	if l.Peek() == '-' {
		l.Advance()
	}

	if l.pos >= l.len || !isDigit(l.data[l.pos]) {
		return 0, fmt.Errorf("line %d, column %d: invalid number", l.line, l.column)
	}

	for l.pos < l.len && isDigit(l.data[l.pos]) {
		l.Advance()
	}

	if l.pos < l.len && l.data[l.pos] == '.' {
		l.Advance()
		if l.pos >= l.len || !isDigit(l.data[l.pos]) {
			return 0, fmt.Errorf("line %d, column %d: invalid decimal number", l.line, l.column)
		}
		for l.pos < l.len && isDigit(l.data[l.pos]) {
			l.Advance()
		}
	}

	if l.pos < l.len && (l.data[l.pos] == 'e' || l.data[l.pos] == 'E') {
		l.Advance()
		if l.pos < l.len && (l.data[l.pos] == '+' || l.data[l.pos] == '-') {
			l.Advance()
		}
		if l.pos >= l.len || !isDigit(l.data[l.pos]) {
			return 0, fmt.Errorf("line %d, column %d: invalid exponent", l.line, l.column)
		}
		for l.pos < l.len && isDigit(l.data[l.pos]) {
			l.Advance()
		}
	}

	numStr := string(l.data[start:l.pos])
	val, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("line %d, column %d: %v", l.line, l.column, err)
	}
	return val, nil
}

// ReadInt reads a JSON number from the input
// and returns it as an integer
// The number may be an integer or a floating point number
// If the number is a floating point number, an error is returned
func (l *JSONLexer) ReadInt() (int, error) {
	val, err := l.ReadNumber()
	if err != nil {
		return 0, err
	}
	if float64(int(val)) != val {
		return 0, fmt.Errorf("line %d, column %d: expected integer, got float", l.line, l.column)
	}
	return int(val), nil
}

// ReadUint reads a JSON number as an unsigned integer, rejecting negatives and floats
func (l *JSONLexer) ReadUint() (uint, error) {
	val, err := l.ReadNumber()
	if err != nil {
		return 0, err
	}
	if val < 0 {
		return 0, fmt.Errorf("line %d, column %d: expected unsigned integer, got negative number", l.line, l.column)
	}
	if float64(uint(val)) != val {
		return 0, fmt.Errorf("line %d, column %d: expected unsigned integer, got float", l.line, l.column)
	}
	return uint(val), nil
}

// ReadUint reads a JSON number from the input
// and returns it as an unsigned integer
// The number may be an integer or a floating point number
// If the number is a floating point number or a negative integer, an error is returned
func (l *JSONLexer) ReadBool() (bool, error) {
	l.SkipWhitespace()
	if l.pos >= l.len {
		return false, fmt.Errorf("line %d, column %d: invalid boolean value", l.line, l.column)
	}
	switch l.data[l.pos] {
	case 't':
		if l.pos+4 <= l.len && l.data[l.pos+1] == 'r' && l.data[l.pos+2] == 'u' && l.data[l.pos+3] == 'e' {
			l.pos += 4
			l.column += 4
			return true, nil
		}
	case 'f':
		if l.pos+5 <= l.len && l.data[l.pos+1] == 'a' && l.data[l.pos+2] == 'l' && l.data[l.pos+3] == 's' && l.data[l.pos+4] == 'e' {
			l.pos += 5
			l.column += 5
			return false, nil
		}
	}
	return false, fmt.Errorf("line %d, column %d: invalid boolean value", l.line, l.column)
}

// ReadNull reads a JSON null value from the input
// If the next 4 bytes are "null", the lexer advances past them
// If the next 4 bytes are not "null", an error is returned
func (l *JSONLexer) ReadNull() error {
	l.SkipWhitespace()
	if l.pos+4 <= l.len && l.data[l.pos] == 'n' && l.data[l.pos+1] == 'u' && l.data[l.pos+2] == 'l' && l.data[l.pos+3] == 'l' {
		l.pos += 4
		l.column += 4
		return nil
	}
	return fmt.Errorf("line %d, column %d: expected null", l.line, l.column)
}

// ReadArrayString reads a JSON array of strings
func (l *JSONLexer) ReadArrayString() ([]string, error) {
	if err := l.Expect(_BracketLeft); err != nil {
		return nil, err
	}

	result := make([]string, 0)

	l.SkipWhitespace()
	if l.Peek() == _BracketRight {
		l.Advance()
		return result, nil
	}

	for {
		val, err := l.ReadString()
		if err != nil {
			return nil, err
		}
		result = append(result, val)

		l.SkipWhitespace()
		if l.Peek() == _BracketRight {
			l.Advance()
			break
		}
		if err := l.Expect(_CommaChar); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// ReadArrayFloat64 reads a JSON array of floating point numbers
func (l *JSONLexer) ReadArrayFloat64() ([]float64, error) {
	if err := l.Expect(_BracketLeft); err != nil {
		return nil, err
	}

	result := make([]float64, 0)

	l.SkipWhitespace()
	if l.Peek() == _BracketRight {
		l.Advance()
		return result, nil
	}

	for {
		val, err := l.ReadNumber()
		if err != nil {
			return nil, err
		}
		result = append(result, val)

		l.SkipWhitespace()
		if l.Peek() == _BracketRight {
			l.Advance()
			break
		}
		if err := l.Expect(_CommaChar); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// SkipValue skips over a JSON value in the input
// This is useful for skipping over JSON values when parsing JSON
// It can be used to skip over JSON objects, arrays, strings, numbers, booleans, and null values
func (l *JSONLexer) SkipValue() error {
	l.SkipWhitespace()
	switch l.Peek() {
	case _BraceLeft:
		l.Advance()
		if l.Peek() == _BraceRight {
			l.Advance()
			return nil
		}
		for {
			if _, err := l.ReadString(); err != nil {
				return err
			}
			if err := l.Expect(_ColonChar); err != nil {
				return err
			}
			if err := l.SkipValue(); err != nil {
				return err
			}
			l.SkipWhitespace()
			if l.Peek() == _BraceRight {
				l.Advance()
				break
			}
			if err := l.Expect(_CommaChar); err != nil {
				return err
			}
		}
	case _BracketLeft:
		l.Advance()
		if l.Peek() == _BracketRight {
			l.Advance()
			return nil
		}
		for {
			if err := l.SkipValue(); err != nil {
				return err
			}
			l.SkipWhitespace()
			if l.Peek() == _BracketRight {
				l.Advance()
				break
			}
			if err := l.Expect(_CommaChar); err != nil {
				return err
			}
		}
	case _QuoteChar:
		_, err := l.ReadString()
		return err
	case 't', 'f':
		_, err := l.ReadBool()
		return err
	case 'n':
		return l.ReadNull()
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		_, err := l.ReadNumber()
		return err
	default:
		return fmt.Errorf("line %d, column %d: unexpected character: %c", l.line, l.column, l.Peek())
	}
	return nil
}

// isDigit returns true if the given byte is a digit character
func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
