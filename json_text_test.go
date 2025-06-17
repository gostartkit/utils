package utils

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestReadUntilNonWhitespace(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		bufSize   int
		pos       int
		n         int
		wantPos   int
		wantN     int
		wantErr   error
		wantChar  byte
		customBuf []byte
	}{
		{
			name:    "Only whitespaces",
			input:   " \t\n\r \t\t\r\n",
			bufSize: 8,
			pos:     0,
			n:       0,
			wantErr: io.EOF,
		},
		{
			name:     "Whitespace then 'a'",
			input:    " \n\t\r a",
			bufSize:  8,
			pos:      0,
			n:        0,
			wantChar: 'a',
		},
		{
			name:     "Immediate non-whitespace",
			input:    "a",
			bufSize:  4,
			pos:      0,
			n:        0,
			wantChar: 'a',
		},
		{
			name:     "Multiple buffers, non-whitespace after wrap",
			input:    " \t\r\n  \t  babcdefg",
			bufSize:  4,
			pos:      0,
			n:        0,
			wantChar: 'b',
		},
		{
			name:      "Whitespace after some initial data in buffer",
			input:     "   xyz",
			bufSize:   6,
			pos:       1,
			n:         4,
			wantChar:  'x',
			customBuf: []byte{' ', ' ', 'x', 'y', 'z', 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			buf := make([]byte, tt.bufSize)

			if tt.customBuf != nil {
				copy(buf, tt.customBuf)
			}

			pos, n, err := ReadUntilNonWhitespace(reader, buf, tt.pos, tt.n)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if pos >= n {
				t.Errorf("pos (%d) >= n (%d): invalid", pos, n)
			}

			if tt.wantChar != 0 && buf[pos] != tt.wantChar {
				t.Errorf("expected char %q at pos %d, got %q", tt.wantChar, pos, buf[pos])
			}
		})
	}
}
