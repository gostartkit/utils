package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func TryParse(val string, attributeType string) (interface{}, error) {
	switch attributeType {
	case "string":
		return val, nil
	case "*string":
		return &val, nil
	case "int":
		n, err := strconv.Atoi(val)
		if err != nil {
			return nil, err
		}
		return n, nil
	case "*int":
		n, err := strconv.Atoi(val)
		if err != nil {
			return nil, err
		}
		return &n, nil
	case "int8":
		n, err := strconv.ParseInt(val, 10, 8)
		if err != nil {
			return nil, err
		}
		return int8(n), nil
	case "*int8":
		n, err := strconv.ParseInt(val, 10, 8)
		if err != nil {
			return nil, err
		}
		n8 := int8(n)
		return &n8, nil
	case "int16":
		n, err := strconv.ParseInt(val, 10, 16)
		if err != nil {
			return nil, err
		}
		return int16(n), nil
	case "*int16":
		n, err := strconv.ParseInt(val, 10, 16)
		if err != nil {
			return nil, err
		}
		n16 := int16(n)
		return &n16, nil
	case "int32":
		n, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			return nil, err
		}
		return int32(n), nil
	case "*int32":
		n, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			return nil, err
		}
		n32 := int32(n)
		return &n32, nil
	case "int64":
		n, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, err
		}
		return n, nil
	case "*int64":
		n, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, err
		}
		return &n, nil
	case "uint":
		n, err := strconv.ParseUint(val, 10, 0)
		if err != nil {
			return nil, err
		}
		return uint(n), nil
	case "*uint":
		n, err := strconv.ParseUint(val, 10, 0)
		if err != nil {
			return nil, err
		}
		n0 := uint(n)
		return &n0, nil
	case "uint8":
		n, err := strconv.ParseUint(val, 10, 8)
		if err != nil {
			return nil, err
		}
		return uint8(n), nil
	case "*uint8":
		n, err := strconv.ParseUint(val, 10, 8)
		if err != nil {
			return nil, err
		}
		n8 := uint8(n)
		return &n8, nil
	case "uint16":
		n, err := strconv.ParseUint(val, 10, 16)
		if err != nil {
			return nil, err
		}
		return uint16(n), nil
	case "*uint16":
		n, err := strconv.ParseUint(val, 10, 16)
		if err != nil {
			return nil, err
		}
		n16 := uint16(n)
		return &n16, nil
	case "uint32":
		n, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			return nil, err
		}
		return uint32(n), nil
	case "*uint32":
		n, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			return nil, err
		}
		n32 := uint32(n)
		return &n32, nil
	case "uint64":
		n, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return nil, err
		}
		return n, nil
	case "*uint64":
		n, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return nil, err
		}
		return &n, nil
	case "float32":
		n, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return nil, err
		}
		return float32(n), nil
	case "*float32":
		n, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return nil, err
		}
		n32 := float32(n)
		return &n32, nil
	case "float64":
		n, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, err
		}
		return n, nil
	case "*float64":
		n, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, err
		}
		return &n, nil
	case "bool":
		n, err := strconv.ParseBool(val)
		if err != nil {
			return nil, err
		}
		return n, nil
	case "*bool":
		n, err := strconv.ParseBool(val)
		if err != nil {
			return nil, err
		}
		return &n, nil
	case "time.Time":
		n, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return nil, err
		}
		return n, nil
	case "*time.Time":
		n, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return nil, err
		}
		return &n, nil
	default:
		return nil, fmt.Errorf("TryParse: attributeType '%s' not supported", attributeType)
	}
}

// IsInt when goType is: int  int8  int16  int32  int64 return true
func IsInt(goType string) bool {
	// int  int8  int16  int32  int64
	switch goType {
	case "int", "int8", "int16", "int32", "int64":
		return true
	default:
		return false
	}
}

// Numeric if attribute type is Numeric return true, else false
func IsNumeric(attributeType string) bool {

	attributeType = strings.TrimPrefix(attributeType, "*")

	switch attributeType {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16",
		"uint32", "uint64", "float32", "float64", "complex64", "complex128":
		return true
	}

	return false
}

// IsBool if attribute type is Numeric return true, else false
func IsBool(attributeType string) bool {

	attributeType = strings.TrimPrefix(attributeType, "*")

	switch attributeType {
	case "bool":
		return true
	}

	return false
}
