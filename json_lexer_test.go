package utils

import (
	"reflect"
	"testing"
)

// TestCreateAndReleaseJSONLexer tests the object pool functionality
func TestCreateAndReleaseJSONLexer(t *testing.T) {
	data := []byte(`{"key": "value"}`)
	lexer := CreateJSONLexer(data)
	if lexer.data == nil || len(lexer.data) != len(data) {
		t.Errorf("CreateJSONLexer failed to initialize data")
	}
	if lexer.pos != 0 || lexer.line != 1 || lexer.column != 1 {
		t.Errorf("CreateJSONLexer failed to initialize position: pos=%d, line=%d, column=%d", lexer.pos, lexer.line, lexer.column)
	}

	ReleaseJSONLexer(lexer)
	if lexer.data != nil || lexer.len != 0 || lexer.pos != 0 {
		t.Errorf("ReleaseJSONLexer failed to reset lexer")
	}
}

// TestPosition tests the Position method
func TestPosition(t *testing.T) {
	lexer := CreateJSONLexer([]byte("abc\ndef"))
	line, col := lexer.Position()
	if line != 1 || col != 1 {
		t.Errorf("Expected position (1, 1), got (%d, %d)", line, col)
	}
	lexer.Advance()
	lexer.Advance()
	lexer.Advance() // Move to newline
	lexer.Advance() // Move to 'd'
	line, col = lexer.Position()
	if line != 2 || col != 1 {
		t.Errorf("Expected position (2, 1) after newline, got (%d, %d)", line, col)
	}
	defer ReleaseJSONLexer(lexer)
}

// TestPeekAndAdvance tests the Peek and Advance methods
func TestPeekAndAdvance(t *testing.T) {
	lexer := CreateJSONLexer([]byte("ab"))
	if lexer.Peek() != 'a' {
		t.Errorf("Peek expected 'a', got '%c'", lexer.Peek())
	}
	lexer.Advance()
	if lexer.Peek() != 'b' {
		t.Errorf("Peek expected 'b', got '%c'", lexer.Peek())
	}
	lexer.Advance()
	if lexer.Peek() != 0 {
		t.Errorf("Peek expected EOF (0), got '%c'", lexer.Peek())
	}
	defer ReleaseJSONLexer(lexer)
}

// TestSkipWhitespace tests the SkipWhitespace method
func TestSkipWhitespace(t *testing.T) {
	lexer := CreateJSONLexer([]byte("  \t\n\r  a"))
	lexer.SkipWhitespace()
	if lexer.Peek() != 'a' {
		t.Errorf("SkipWhitespace failed, expected 'a', got '%c'", lexer.Peek())
	}
	defer ReleaseJSONLexer(lexer)
}

// TestExpect tests the Expect method
func TestExpect(t *testing.T) {
	lexer := CreateJSONLexer([]byte(`  ,`))
	if err := lexer.Expect(_CommaChar); err != nil {
		t.Errorf("Expect comma failed: %v", err)
	}
	if err := lexer.Expect(_QuoteChar); err == nil {
		t.Errorf("Expect quote should have failed")
	}
	defer ReleaseJSONLexer(lexer)
}

// TestReadString tests the ReadString method with various cases
func TestReadString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{`"hello"`, "hello", false},
		{`"hello \"world\""`, "hello \"world\"", false},
		{`"line\nbreak"`, "line\nbreak", false},
		{`"unicode \u0041"`, "unicode A", false},
		{`"incomplete`, "", true},
		{`invalid"`, "", true},
		{`"control \x01"`, "", true},
	}

	for _, tt := range tests {
		lexer := CreateJSONLexer([]byte(tt.input))
		result, err := lexer.ReadString()
		if (err != nil) != tt.wantErr {
			t.Errorf("ReadString(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && result != tt.expected {
			t.Errorf("ReadString(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
		ReleaseJSONLexer(lexer)
	}
}

// TestReadNumber tests the ReadNumber method with various cases
func TestReadNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		wantErr  bool
	}{
		{"123", 123, false},
		{"-456.78", -456.78, false},
		{"1.23e4", 12300, false},
		{"0.0", 0.0, false},
		{"-abc", 0, true},
		{"1.", 0, true},
		{".1", 0, true},
	}

	for _, tt := range tests {
		lexer := CreateJSONLexer([]byte(tt.input))
		result, err := lexer.ReadNumber()
		if (err != nil) != tt.wantErr {
			t.Errorf("ReadNumber(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && result != tt.expected {
			t.Errorf("ReadNumber(%q) = %f, expected %f", tt.input, result, tt.expected)
		}
		ReleaseJSONLexer(lexer)
	}
}

// TestReadInt tests the ReadInt method
func TestReadInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		wantErr  bool
	}{
		{"123", 123, false},
		{"-456", -456, false},
		{"1.23", 0, true},
		{"abc", 0, true},
	}

	for _, tt := range tests {
		lexer := CreateJSONLexer([]byte(tt.input))
		result, err := lexer.ReadInt()
		if (err != nil) != tt.wantErr {
			t.Errorf("ReadInt(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && result != tt.expected {
			t.Errorf("ReadInt(%q) = %d, expected %d", tt.input, result, tt.expected)
		}
		ReleaseJSONLexer(lexer)
	}
}

// TestReadUint tests the ReadUint method
func TestReadUint(t *testing.T) {
	tests := []struct {
		input    string
		expected uint
		wantErr  bool
	}{
		{"123", 123, false},
		{"0", 0, false},
		{"-456", 0, true},
		{"1.23", 0, true},
	}

	for _, tt := range tests {
		lexer := CreateJSONLexer([]byte(tt.input))
		result, err := lexer.ReadUint()
		if (err != nil) != tt.wantErr {
			t.Errorf("ReadUint(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && result != tt.expected {
			t.Errorf("ReadUint(%q) = %d, expected %d", tt.input, result, tt.expected)
		}
		ReleaseJSONLexer(lexer)
	}
}

// TestReadBool tests the ReadBool method
func TestReadBool(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		wantErr  bool
	}{
		{"true", true, false},
		{"false", false, false},
		{"tru", false, true},
		{"fals", false, true},
	}

	for _, tt := range tests {
		lexer := CreateJSONLexer([]byte(tt.input))
		result, err := lexer.ReadBool()
		if (err != nil) != tt.wantErr {
			t.Errorf("ReadBool(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && result != tt.expected {
			t.Errorf("ReadBool(%q) = %v, expected %v", tt.input, result, tt.expected)
		}
		ReleaseJSONLexer(lexer)
	}
}

// TestReadNull tests the ReadNull method
func TestReadNull(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"null", false},
		{"nul", true},
		{"abc", true},
	}

	for _, tt := range tests {
		lexer := CreateJSONLexer([]byte(tt.input))
		err := lexer.ReadNull()
		if (err != nil) != tt.wantErr {
			t.Errorf("ReadNull(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
		ReleaseJSONLexer(lexer)
	}
}

// TestReadArrayFloat64 tests the ReadArrayFloat64 method
func TestReadArrayFloat64(t *testing.T) {
	tests := []struct {
		input    string
		expected []float64
		wantErr  bool
	}{
		{"[1, 2.3, -4]", []float64{1, 2.3, -4}, false},
		{"[]", []float64{}, false},
		{"[1, abc]", nil, true},
		{"[1 2]", nil, true}, // Missing comma
	}

	for _, tt := range tests {
		lexer := CreateJSONLexer([]byte(tt.input))
		result, err := lexer.ReadArrayFloat64()
		if (err != nil) != tt.wantErr {
			t.Errorf("ReadArrayFloat64(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("ReadArrayFloat64(%q) = %v, expected %v", tt.input, result, tt.expected)
		}
		ReleaseJSONLexer(lexer)
	}
}

// TestReadArrayString tests the ReadArrayString method
func TestReadArrayString(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
		wantErr  bool
	}{
		{`["a", "b c"]`, []string{"a", "b c"}, false},
		{`[]`, []string{}, false},
		{`["a", b]`, nil, true},
		{`["a" "b"]`, nil, true}, // Missing comma
	}

	for _, tt := range tests {
		lexer := CreateJSONLexer([]byte(tt.input))
		result, err := lexer.ReadArrayString()
		if (err != nil) != tt.wantErr {
			t.Errorf("ReadArrayString(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("ReadArrayString(%q) = %v, expected %v", tt.input, result, tt.expected)
		}
		ReleaseJSONLexer(lexer)
	}
}

// TestSkipValue tests the SkipValue method with various JSON structures
func TestSkipValue(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{`{"a": 1}`, false},
		{`[1, "b"]`, false},
		{`"string"`, false},
		{`true`, false},
		{`null`, false},
		{`123.45`, false},
		{`{`, true},          // Incomplete object
		{`["unclosed`, true}, // Incomplete array
	}

	for _, tt := range tests {
		lexer := CreateJSONLexer([]byte(tt.input))
		err := lexer.SkipValue()
		if (err != nil) != tt.wantErr {
			t.Errorf("SkipValue(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
		ReleaseJSONLexer(lexer)
	}
}

// BenchmarkReadString benchmarks the performance of ReadString
func BenchmarkReadString(b *testing.B) {
	data := []byte(`"this is a test string with \"escapes\" and unicode \u0041"`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer := CreateJSONLexer(data)
		_, _ = lexer.ReadString()
		ReleaseJSONLexer(lexer)
	}
}

// BenchmarkReadArrayFloat64 benchmarks the performance of ReadArrayFloat64
func BenchmarkReadArrayFloat64(b *testing.B) {
	data := []byte(`[1.23, 45.67, 89.01, -23.45]`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer := CreateJSONLexer(data)
		_, _ = lexer.ReadArrayFloat64()
		ReleaseJSONLexer(lexer)
	}
}
