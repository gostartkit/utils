package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

var jsonLexerPool = sync.Pool{
	New: func() any {
		return &JSONLexer{}
	},
}

func CreateJSONLexer(data []byte) *JSONLexer {
	lexer := jsonLexerPool.Get().(*JSONLexer)
	lexer.data = data
	lexer.pos = 0
	return lexer
}

func ReleaseJSONLexer(lexer *JSONLexer) {
	lexer.data = nil
	lexer.pos = 0
	jsonLexerPool.Put(lexer)
}

type JSONLexer struct {
	data []byte
	pos  int
}

func (l *JSONLexer) Peek() byte {
	if l.pos >= len(l.data) {
		return 0
	}
	return l.data[l.pos]
}

func (l *JSONLexer) Advance() {
	if l.pos < len(l.data) {
		l.pos++
	}
}

func (l *JSONLexer) SkipWhitespace() {
	for l.pos < len(l.data) && (l.data[l.pos] == ' ' || l.data[l.pos] == '\n' || l.data[l.pos] == '\t' || l.data[l.pos] == '\r') {
		l.pos++
	}
}

func (l *JSONLexer) Expect(c byte) error {
	l.SkipWhitespace()
	if l.Peek() != c {
		return fmt.Errorf("expected '%c', got '%c'", c, l.Peek())
	}
	l.Advance()
	return nil
}

func (l *JSONLexer) ReadString() (string, error) {
	if err := l.Expect('"'); err != nil {
		return "", err
	}

	start := l.pos
	for {
		if l.pos >= len(l.data) {
			return "", errors.New("unterminated string")
		}
		if l.data[l.pos] == '"' {
			break
		}
		if l.data[l.pos] == '\\' {
			l.pos++
		}
		l.pos++
	}

	result := string(l.data[start:l.pos])
	l.Advance()
	return result, nil
}

func (l *JSONLexer) ReadInt() (int, error) {
	start := l.pos
	for l.pos < len(l.data) && (l.data[l.pos] >= '0' && l.data[l.pos] <= '9') {
		l.pos++
	}
	val, err := strconv.Atoi(string(l.data[start:l.pos]))
	return val, err
}

func (l *JSONLexer) ReadBool() (bool, error) {
	if strings.HasPrefix(string(l.data[l.pos:]), "true") {
		l.pos += 4
		return true, nil
	}
	if strings.HasPrefix(string(l.data[l.pos:]), "false") {
		l.pos += 5
		return false, nil
	}
	return false, errors.New("invalid boolean value")
}

func (l *JSONLexer) ReadFloatArray() ([]float64, error) {
	if err := l.Expect('['); err != nil {
		return nil, err
	}
	var result []float64
	for {
		l.SkipWhitespace()
		if l.Peek() == ']' {
			l.Advance()
			break
		}
		start := l.pos
		for l.pos < len(l.data) && (l.data[l.pos] == '.' || (l.data[l.pos] >= '0' && l.data[l.pos] <= '9')) {
			l.pos++
		}
		val, err := strconv.ParseFloat(string(l.data[start:l.pos]), 64)
		if err != nil {
			return nil, err
		}
		result = append(result, val)

		l.SkipWhitespace()
		if l.Peek() == ',' {
			l.Advance()
		} else if l.Peek() == ']' {
			l.Advance()
			break
		} else {
			return nil, errors.New("unexpected character in array")
		}
	}
	return result, nil
}

func (l *JSONLexer) SkipValue() error {
	l.SkipWhitespace()
	switch l.Peek() {
	case '{':
		l.Advance()
		for l.Peek() != '}' {
			if err := l.SkipValue(); err != nil {
				return err
			}
			l.SkipWhitespace()
			if l.Peek() == ',' {
				l.Advance()
			}
		}
		l.Advance()
	case '[':
		l.Advance()
		for l.Peek() != ']' {
			if err := l.SkipValue(); err != nil {
				return err
			}
			l.SkipWhitespace()
			if l.Peek() == ',' {
				l.Advance()
			}
		}
		l.Advance()
	case '"':
		_, err := l.ReadString()
		return err
	default:
		for l.pos < len(l.data) && !unicode.IsSpace(rune(l.data[l.pos])) && l.data[l.pos] != ',' && l.data[l.pos] != '}' && l.data[l.pos] != ']' {
			l.pos++
		}
	}
	return nil
}
