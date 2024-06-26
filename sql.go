package utils

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// SqlFilter create sql for filter and args
func SqlFilter(filter string, w io.Writer, args *[]any, prefix string, fn func(key string, val string) (string, interface{}, error)) error {

	vals := filterParse(filter)

	l := len(vals)

	if l == 0 {
		return ErrFilter
	}

	var prev string
	var space bool
	var valid bool = false

	var spaceByte = []byte(" ")

	for i := 0; i < l; i++ {

		val := vals[i]

		switch val {
		case "eq", "ne", "gt", "ge", "lt", "le":
			if space {
				w.Write(spaceByte)
			} else {
				space = true
			}

			op := filterOp(val)

			i, val = filterNext(vals, i, l)

			n, v, err := fn(prev, val)

			if err != nil {
				return err
			}

			valid = true

			*args = append(*args, v)
			fmt.Fprintf(w, "%s`%s` %s ?", prefix, n, op)
		case "and", "or", "not":
			if space {
				w.Write(spaceByte)
			} else {
				space = true
			}
			w.Write([]byte(strings.ToUpper(val)))
		case "(", ")":
			w.Write([]byte(val))
		default:
			prev = val
		}
	}

	if !valid {
		return ErrFilter
	}

	return nil
}

// SqlOrderBy create sql for order by
func SqlOrderBy(orderBy string, w io.Writer, prefix string, fn func(variableName string) (string, string, string, error)) error {

	vals := orderByParse(orderBy)

	l := len(vals)

	if l == 0 {
		return ErrOrderBy
	}

	var valid bool = false

	var commaBytes = []byte(", ")
	var ascBytes = []byte(" ASC")
	var descBytes = []byte(" DESC")
	var backtickBytes = []byte("`")

	for i := 0; i < l; i++ {

		val := vals[i]

		switch val {
		case ",":
			w.Write(commaBytes)
		case "asc":
			w.Write(ascBytes)
		case "desc":
			w.Write(descBytes)
		default:
			n, _, _, err := fn(val)

			if err != nil {
				return err
			}

			valid = true

			w.Write([]byte(prefix))
			w.Write(backtickBytes)
			w.Write([]byte(n))
			w.Write(backtickBytes)
		}
	}

	if !valid {
		return ErrOrderBy
	}

	return nil
}

// filterParse parse filter to vals
func filterParse(filter string) []string {
	l := len(filter)

	prev := 0
	inSingle := false
	inQuotes := false

	vals := []string{}

	for pos := 0; pos < l; pos++ {
		r := filter[pos]

		switch r {
		case '(', ')':

			if inSingle || inQuotes {
				continue
			}

			if pos > prev {
				vals = append(vals, filter[prev:pos])
			}

			prev = pos + 1

			vals = append(vals, string(r))

		case ' ', '\t':

			if inSingle || inQuotes {
				continue
			}

			if pos > prev {
				vals = append(vals, filter[prev:pos])
			}

			prev = pos + 1
		case '\'':
			if inQuotes {
				continue
			}

			inSingle = !inSingle

			if pos > prev {
				vals = append(vals, filter[prev:pos])
			}

			prev = pos + 1
		case '"':
			if inSingle {
				continue
			}

			inQuotes = !inQuotes

			if pos > prev {
				vals = append(vals, filter[prev:pos])
			}

			prev = pos + 1
		}
	}

	if prev < l {
		vals = append(vals, filter[prev:])
	}

	return vals
}

// filterNext read next item
func filterNext(items []string, i int, l int) (int, string) {
	pos := i + 1

	if pos < l {
		return pos, items[pos]
	}

	return i, ""
}

// filterOp return op
func filterOp(key string) string {
	switch key {
	case "eq":
		return "="
	case "ne":
		return "<>"
	case "gt":
		return ">"
	case "ge":
		return ">="
	case "lt":
		return "<"
	case "le":
		return "<="
	default:
		return "="
	}
}

// orderByParse parse orderBy to vals
func orderByParse(orderBy string) []string {
	l := len(orderBy)

	prev := 0

	vals := []string{}

	for pos := 0; pos < l; pos++ {
		r := orderBy[pos]

		switch r {
		case ',':

			if pos > prev {
				vals = append(vals, orderBy[prev:pos])
			}

			prev = pos + 1

			vals = append(vals, string(r))

		case ' ', '\t':

			if pos > prev {
				vals = append(vals, orderBy[prev:pos])
			}

			prev = pos + 1
		}
	}

	if prev < l {
		vals = append(vals, orderBy[prev:])
	}

	return vals
}

// MySQLSplit split function for the Scanner
func MySQLSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {

	inDoubleQuotes := false
	inSingleQuotes := false

	for i := 0; i < len(data); i++ {

		switch data[i] {
		case '"':
			if !inSingleQuotes && i > 0 && data[i-1] != '\\' {
				inDoubleQuotes = !inDoubleQuotes
				continue
			}
		case '\'':
			if !inDoubleQuotes && i > 0 && data[i-1] != '\\' {
				inSingleQuotes = !inSingleQuotes
				continue
			}
		case ';':
			if inDoubleQuotes || inSingleQuotes {
				continue
			}

			return i + 1, data[:i], nil
		}
	}

	if !atEOF {
		return 0, nil, nil
	}

	return 0, data, bufio.ErrFinalToken
}

// Where parse where
func Where(where string, fn func(key string, op string, val string)) {

	vals := whereParse(where)

	l := len(vals)

	var prev string

	for i := 0; i < l; i++ {

		val := vals[i]

		switch strings.ToLower(val) {
		case "=", "<>", ">", ">=", "<", "<=":
			op := val
			i, val = whereNext(vals, i, l)
			fn(prev, op, val)
		case "and", "or", "not":
			fn(val, "", "")
		case "(", ")":
			fn(val, "", "")
		default:
			prev = val
		}
	}
}

// whereParse parse filter to vals
func whereParse(where string) []string {

	l := len(where)

	prev := 0
	inSingle := false
	inQuotes := false

	vals := []string{}

	for pos := 0; pos < l; pos++ {

		r := where[pos]

		switch r {
		case '(', ')':

			if inSingle || inQuotes {
				continue
			}

			if pos > prev {
				vals = append(vals, where[prev:pos])
			}

			prev = pos + 1

			vals = append(vals, string(r))

		case ' ', '\t':

			if inSingle || inQuotes {
				continue
			}

			if pos > prev {
				vals = append(vals, where[prev:pos])
			}

			prev = pos + 1
		case '\'':
			if inQuotes {
				continue
			}

			inSingle = !inSingle

			if pos > prev {
				vals = append(vals, where[prev:pos])
			}

			prev = pos + 1
		case '"':
			if inSingle {
				continue
			}

			inQuotes = !inQuotes

			if pos > prev {
				vals = append(vals, where[prev:pos])
			}

			prev = pos + 1
		}
	}

	if prev < l {
		vals = append(vals, where[prev:])
	}

	return vals
}

// whereNext read whereNext item
func whereNext(items []string, i int, l int) (int, string) {
	pos := i + 1

	if pos < l {
		return pos, items[pos]
	}

	return i, ""
}
