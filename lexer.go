package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

var lexerPool = sync.Pool{
	New: func() any {
		return &Lexer{}
	},
}

func CreateLexer(data []byte) *Lexer {
	lexer := lexerPool.Get().(*Lexer)
	lexer.data = data
	lexer.pos = 0
	return lexer
}

func ReleaseLexer(lexer *Lexer) {
	lexer.data = nil
	lexer.pos = 0
	lexerPool.Put(lexer)
}

type Lexer struct {
	data []byte
	pos  int
}

func (l *Lexer) Peek() byte {
	if l.pos >= len(l.data) {
		return 0
	}
	return l.data[l.pos]
}

func (l *Lexer) Advance() {
	if l.pos < len(l.data) {
		l.pos++
	}
}

func (l *Lexer) SkipWhitespace() {
	for l.pos < len(l.data) && (l.data[l.pos] == ' ' || l.data[l.pos] == '\n' || l.data[l.pos] == '\t' || l.data[l.pos] == '\r') {
		l.pos++
	}
}

func (l *Lexer) Expect(c byte) error {
	l.SkipWhitespace()
	if l.Peek() != c {
		return fmt.Errorf("expected '%c', got '%c'", c, l.Peek())
	}
	l.Advance()
	return nil
}

func (l *Lexer) ReadString() (string, error) {
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

func (l *Lexer) ReadInt() (int, error) {
	start := l.pos
	for l.pos < len(l.data) && (l.data[l.pos] >= '0' && l.data[l.pos] <= '9') {
		l.pos++
	}
	val, err := strconv.Atoi(string(l.data[start:l.pos]))
	return val, err
}

func (l *Lexer) ReadBool() (bool, error) {
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

func (l *Lexer) ReadFloatArray() ([]float64, error) {
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

func (l *Lexer) SkipValue() error {
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
