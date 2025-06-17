package utils

import (
	"fmt"
	"io"
	"strings"
)

// ReadUntilNonWhitespace reads from io.Reader until a non-whitespace character is found,
// using buf as a scratch buffer. It starts at position pos in buf, where n bytes are available.
// Whitespace is defined as space (' '), newline ('\n'), carriage return ('\r'), or tab ('\t') per JSON (RFC 8259).
// Returns the position of the first non-whitespace character, the number of bytes available in buf, and any error.
func ReadUntilNonWhitespace(r io.Reader, buf []byte, pos int, n int) (int, int, error) {

	l := len(buf)

	if n > l {
		return pos, n, fmt.Errorf("invalid buffer: n (%d) exceeds buffer length (%d)", n, l)
	}

	if pos > n {
		return pos, n, fmt.Errorf("invalid position: pos (%d) greater than n (%d)", pos, n)
	}

	var err error

	// Track start position for error reporting
	startPos := pos

	for {
		// Process buffer contents
		for pos < n {
			switch buf[pos] {
			case ' ', '\n', '\r', '\t':
				pos++
			default:
				return pos, n, nil // Found non-whitespace
			}
		}

		// If buffer is exhausted, read more data
		n, err = r.Read(buf)
		if err != nil && err != io.EOF {
			return pos, n, fmt.Errorf("read error at position %d: %w", startPos, err)
		}
		if n == 0 {
			return pos, n, io.EOF // No non-whitespace character found
		}
		pos = 0
	}
}

// ReadString reads a JSON string from the input reader `r`,
// using buffer `buf` with current offset `pos` and length `n`.
// It handles JSON escape sequences and Unicode surrogate pairs properly.
// Returns the parsed string, new position in buffer, or an error.
func ReadString(r io.Reader, buf []byte, pos, n int) (string, int, error) {
	// Refill buffer if position has already reached the end
	if pos >= n {
		var err error
		n, err = r.Read(buf)
		if err != nil {
			return "", pos, err
		}
		pos = 0
	}

	// JSON strings must start with a quote character
	if buf[pos] != '"' {
		return "", pos, fmt.Errorf("expecting '\"' at position %d", pos)
	}
	pos++

	var result strings.Builder
	result.Grow(64) // Preallocate capacity to avoid frequent realloc

	const (
		StateNormal  = iota // Normal character parsing
		StateEscape         // Parsing after encountering a backslash '\'
		StateUnicode        // Parsing 4-digit hex unicode escape \uXXXX
	)

	state := StateNormal
	var unicodeBuf [4]byte // Holds 4 hex digits for \uXXXX
	var surrogateHi rune   // Holds high surrogate, if any

	for {
		// If we reach the end of buffer, refill it
		if pos >= n {
			n = copy(buf, buf[pos:n]) // shift remaining data to front
			pos = 0
			n2, err := r.Read(buf[n:])
			n += n2
			if err != nil && err != io.EOF {
				return "", pos, err
			}
			if n == 0 {
				return "", pos, fmt.Errorf("unexpected EOF while parsing string")
			}
		}

		c := buf[pos]
		pos++

		switch state {
		case StateNormal:
			switch c {
			case '"':
				// End of string
				if surrogateHi != 0 {
					return "", pos, fmt.Errorf("incomplete unicode surrogate pair")
				}
				return result.String(), pos, nil
			case '\\':
				// Escape starts
				state = StateEscape
			default:
				// Normal character
				result.WriteByte(c)
			}

		case StateEscape:
			switch c {
			case '"', '\\', '/':
				result.WriteByte(c)
				state = StateNormal
			case 'b':
				result.WriteByte('\b')
				state = StateNormal
			case 'f':
				result.WriteByte('\f')
				state = StateNormal
			case 'n':
				result.WriteByte('\n')
				state = StateNormal
			case 'r':
				result.WriteByte('\r')
				state = StateNormal
			case 't':
				result.WriteByte('\t')
				state = StateNormal
			case 'u':
				// Begin parsing \uXXXX
				state = StateUnicode
			default:
				return "", pos, fmt.Errorf("invalid escape sequence \\%c", c)
			}

		case StateUnicode:
			// Read next 4 hex digits from the stream into unicodeBuf
			for i := 0; i < 4; i++ {
				if pos >= n {
					copy(buf, buf[pos:n])
					n = n - pos
					pos = 0
					n2, err := r.Read(buf[n:])
					n += n2
					if err != nil && err != io.EOF {
						return "", pos, err
					}
					if n < 4 {
						return "", pos, fmt.Errorf("incomplete \\uXXXX escape")
					}
				}
				unicodeBuf[i] = buf[pos]
				pos++
			}

			// Parse 4 hex digits into a Unicode code point
			code, err := parseHex4(unicodeBuf[:])
			if err != nil {
				return "", pos, fmt.Errorf("invalid \\u escape: %w", err)
			}
			r1 := rune(code)

			// Handle surrogate pairs: \uD800–\uDBFF followed by \uDC00–\uDFFF
			if isSurrogate(r1) {
				if surrogateHi == 0 && isHighSurrogate(r1) {
					// Save the high surrogate and expect a low surrogate next
					surrogateHi = r1
					state = StateEscape // Expecting next \uXXXX
					continue
				} else if surrogateHi != 0 && isLowSurrogate(r1) {
					// Combine high and low surrogates into a full codepoint
					full := decodeRune(surrogateHi, r1)
					result.WriteRune(full)
					surrogateHi = 0
				} else {
					return "", pos, fmt.Errorf("invalid unicode surrogate sequence")
				}
			} else {
				if surrogateHi != 0 {
					return "", pos, fmt.Errorf("unexpected low surrogate without high surrogate")
				}
				result.WriteRune(r1)
			}

			// Resume normal parsing after unicode
			state = StateNormal
		}
	}
}

// ReadNumber reads a JSON number from io.Reader, using buf as a scratch buffer.
// It starts reading at position pos in buf, where n bytes are available.
// Returns the parsed number as a string, the new position, and any error.
func ReadNumber(r io.Reader, buf []byte, pos, n int) (string, int, error) {
	// Validate buffer state
	if n > len(buf) {
		return "", pos, fmt.Errorf("invalid buffer: n (%d) exceeds buffer length (%d)", n, len(buf))
	}

	// Use strings.Builder for efficient number construction
	var result strings.Builder
	result.Grow(n - pos) // Pre-allocate capacity
	startPos := pos      // Track start position for error reporting

	// Ensure buffer has data
	if pos >= n {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return "", pos, fmt.Errorf("read error at position %d: %w", pos, err)
		}
		if n == 0 {
			return "", pos, fmt.Errorf("expected digit at position %d: %w", pos, io.EOF)
		}
		pos = 0
	}

	// Handle minus sign
	c := buf[pos]
	if c == '-' {
		result.WriteByte(c)
		pos++
		if pos >= n {
			n, err := r.Read(buf)
			if err != nil && err != io.EOF {
				return "", pos, fmt.Errorf("read error at position %d: %w", pos, err)
			}
			if n == 0 {
				return "", pos, fmt.Errorf("incomplete number: only minus sign at position %d", startPos)
			}
			pos = 0
		}
		c = buf[pos]
	}

	// Handle integer part
	if c == '0' {
		result.WriteByte(c)
		pos++
		// JSON forbids leading zeros (e.g., "01" or "-01")
		if pos < n && buf[pos] >= '0' && buf[pos] <= '9' {
			return "", pos, fmt.Errorf("invalid number: leading zero followed by digit at position %d", pos)
		}
	} else if c >= '1' && c <= '9' {
		result.WriteByte(c)
		pos++
		// Read subsequent digits
		for pos < n {
			c = buf[pos]
			if c >= '0' && c <= '9' {
				result.WriteByte(c)
				pos++
			} else {
				break
			}
		}
	} else {
		return "", pos, fmt.Errorf("expected digit at position %d, got %q", pos, c)
	}

	// Handle fraction part (optional)
	if pos < n && buf[pos] == '.' {
		result.WriteByte('.')
		pos++
		// Require at least one digit after decimal point
		if pos >= n {
			n, err := r.Read(buf)
			if err != nil && err != io.EOF {
				return "", pos, fmt.Errorf("read error at position %d: %w", pos, err)
			}
			if n == 0 {
				return "", pos, fmt.Errorf("incomplete number: no digits after decimal point at position %d", startPos)
			}
			pos = 0
		}
		c = buf[pos]
		if c < '0' || c > '9' {
			return "", pos, fmt.Errorf("expected digit after decimal point at position %d, got %q", pos, c)
		}
		result.WriteByte(c)
		pos++
		// Read additional digits
		for pos < n {
			c = buf[pos]
			if c >= '0' && c <= '9' {
				result.WriteByte(c)
				pos++
			} else {
				break
			}
		}
	}

	// Handle exponent part (optional)
	if pos < n && (buf[pos] == 'e' || buf[pos] == 'E') {
		result.WriteByte(buf[pos])
		pos++
		// Handle optional sign
		if pos < n && (buf[pos] == '+' || buf[pos] == '-') {
			result.WriteByte(buf[pos])
			pos++
		}
		// Require at least one digit after exponent
		if pos >= n {
			n, err := r.Read(buf)
			if err != nil && err != io.EOF {
				return "", pos, fmt.Errorf("read error at position %d: %w", pos, err)
			}
			if n == 0 {
				return "", pos, fmt.Errorf("incomplete number: no digits after exponent at position %d", startPos)
			}
			pos = 0
		}
		c = buf[pos]
		if c < '0' || c > '9' {
			return "", pos, fmt.Errorf("expected digit after exponent at position %d, got %q", pos, c)
		}
		result.WriteByte(c)
		pos++
		// Read additional digits
		for pos < n {
			c = buf[pos]
			if c >= '0' && c <= '9' {
				result.WriteByte(c)
				pos++
			} else {
				break
			}
		}
	}

	// Ensure a valid number was parsed
	if result.Len() == 0 || (result.Len() == 1 && result.String() == "-") {
		return "", pos, fmt.Errorf("invalid number at position %d", startPos)
	}

	return result.String(), pos, nil
}

// SkipValue skips a single JSON value from io.Reader, using buf as a scratch buffer.
// It starts at position pos in buf, where n bytes are available.
// The value can be a string, number, object, array, true, false, or null (per RFC 8259).
// Returns the new position, number of bytes available in buf, and any error.
func SkipValue(r io.Reader, buf []byte, pos, n int) (newPos, newN int, err error) {
	// Validate buffer state
	if n > len(buf) {
		return pos, n, fmt.Errorf("invalid buffer: n (%d) exceeds buffer length (%d)", n, len(buf))
	}

	// Skip leading whitespace
	pos, n, err = ReadUntilNonWhitespace(r, buf, pos, n)
	if err != nil {
		return pos, n, err
	}

	// Check if buffer is exhausted
	if pos >= n {
		n, err = r.Read(buf)
		if err != nil && err != io.EOF {
			return pos, n, fmt.Errorf("read error at position %d: %w", pos, err)
		}
		if n == 0 {
			return pos, n, fmt.Errorf("expected JSON value at position %d: %w", pos, io.EOF)
		}
		pos = 0
	}

	// Handle value based on first character
	c := buf[pos]
	switch c {
	case '"':
		// Handle string
		_, pos, err = ReadString(r, buf, pos, n)
		if err != nil {
			return pos, n, fmt.Errorf("invalid string at position %d: %w", pos, err)
		}
		return pos, n, nil

	case '{', '[':
		// Handle object or array
		opening := c
		depth := 1
		pos++

		for depth > 0 {
			// Skip whitespace
			pos, n, err = ReadUntilNonWhitespace(r, buf, pos, n)
			if err != nil {
				return pos, n, fmt.Errorf("error in nested structure at position %d: %w", pos, err)
			}

			// Check if buffer is exhausted
			if pos >= n {
				n, err = r.Read(buf)
				if err != nil && err != io.EOF {
					return pos, n, fmt.Errorf("read error at position %d: %w", pos, err)
				}
				if n == 0 {
					return pos, n, fmt.Errorf("unclosed %c at position %d: %w", opening, pos, io.EOF)
				}
				pos = 0
			}

			c = buf[pos]
			if c == '"' {
				// Skip string
				_, pos, err = ReadString(r, buf, pos, n)
				if err != nil {
					return pos, n, fmt.Errorf("invalid string in %c at position %d: %w", opening, pos, err)
				}
			} else if c == '{' || c == '[' {
				// Increase depth for nested object or array
				depth++
				pos++
			} else if c == '}' || c == ']' {
				// Decrease depth and check matching
				if (c == '}' && opening != '{') || (c == ']' && opening != '[') {
					return pos, n, fmt.Errorf("mismatched closing %c for %c at position %d", c, opening, pos)
				}
				depth--
				pos++
			} else {
				// Skip other values (number, true, false, null)
				_, pos, err = SkipValue(r, buf, pos, n)
				if err != nil {
					return pos, n, fmt.Errorf("invalid value in %c at position %d: %w", opening, pos, err)
				}
			}
		}
		return pos, n, nil

	case 't':
		// Handle true
		if pos+3 >= n {
			// Move remaining bytes to start and read more
			if pos < n {
				copy(buf[0:], buf[pos:n])
				n -= pos
				pos = 0
			} else {
				n = 0
				pos = 0
			}
			n2, err := r.Read(buf[n:])
			if err != nil && err != io.EOF {
				return pos, n, fmt.Errorf("read error at position %d: %w", pos, err)
			}
			n += n2
			if pos+3 >= n {
				return pos, n, fmt.Errorf("incomplete true at position %d: %w", pos, io.EOF)
			}
		}
		if string(buf[pos:pos+4]) != "true" {
			return pos, n, fmt.Errorf("expected true at position %d, got %q", pos, string(buf[pos:pos+4]))
		}
		pos += 4
		return pos, n, nil

	case 'f':
		// Handle false
		if pos+4 >= n {
			if pos < n {
				copy(buf[0:], buf[pos:n])
				n -= pos
				pos = 0
			} else {
				n = 0
				pos = 0
			}
			n2, err := r.Read(buf[n:])
			if err != nil && err != io.EOF {
				return pos, n, fmt.Errorf("read error at position %d: %w", pos, err)
			}
			n += n2
			if pos+4 >= n {
				return pos, n, fmt.Errorf("incomplete false at position %d: %w", pos, io.EOF)
			}
		}
		if string(buf[pos:pos+5]) != "false" {
			return pos, n, fmt.Errorf("expected false at position %d, got %q", pos, string(buf[pos:pos+5]))
		}
		pos += 5
		return pos, n, nil

	case 'n':
		// Handle null
		if pos+3 >= n {
			if pos < n {
				copy(buf[0:], buf[pos:n])
				n -= pos
				pos = 0
			} else {
				n = 0
				pos = 0
			}
			n2, err := r.Read(buf[n:])
			if err != nil && err != io.EOF {
				return pos, n, fmt.Errorf("read error at position %d: %w", pos, err)
			}
			n += n2
			if pos+3 >= n {
				return pos, n, fmt.Errorf("incomplete null at position %d: %w", pos, io.EOF)
			}
		}
		if string(buf[pos:pos+4]) != "null" {
			return pos, n, fmt.Errorf("expected null at position %d, got %q", pos, string(buf[pos:pos+4]))
		}
		pos += 4
		return pos, n, nil

	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		// Handle number
		_, pos, err = ReadNumber(r, buf, pos, n)
		if err != nil {
			return pos, n, fmt.Errorf("invalid number at position %d: %w", pos, err)
		}
		return pos, n, nil

	default:
		return pos, n, fmt.Errorf("invalid JSON value starting with %q at position %d", c, pos)
	}
}

// parseUint64
func parseUint64(s string) (uint64, error) {
	var n uint64
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, io.ErrUnexpectedEOF
		}
		n = n*10 + uint64(c-'0')
	}
	return n, nil
}

func parseHex4(buf []byte) (uint16, error) {
	var value uint16
	for i := 0; i < 4; i++ {
		c := buf[i]
		var digit uint16
		switch {
		case '0' <= c && c <= '9':
			digit = uint16(c - '0')
		case 'a' <= c && c <= 'f':
			digit = uint16(c - 'a' + 10)
		case 'A' <= c && c <= 'F':
			digit = uint16(c - 'A' + 10)
		default:
			return 0, fmt.Errorf("invalid hexadecimal character %c", c)
		}
		value = value<<4 | digit
	}
	return value, nil
}

func isSurrogate(r rune) bool {
	return r >= 0xD800 && r <= 0xDFFF
}

func isHighSurrogate(r rune) bool {
	return r >= 0xD800 && r <= 0xDBFF
}

func isLowSurrogate(r rune) bool {
	return r >= 0xDC00 && r <= 0xDFFF
}

func decodeRune(high, low rune) rune {
	return ((high-0xD800)<<10 | (low - 0xDC00)) + 0x10000
}
