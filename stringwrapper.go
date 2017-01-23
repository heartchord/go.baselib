package goblazer

import (
	"bytes"
	"strconv"
	"strings"
)

// TrueStrings is
var TrueStrings = []string{"true", "1", "yes"}

// FalseStrings is
var FalseStrings = []string{"false", "0", "no", ""}

// GetBoolString is
func GetBoolString(b bool) string {
	if b {
		return TrueStrings[0]
	}
	return FalseStrings[0]
}

// IsTrueString is
func IsTrueString(s string) bool {
	for _, v := range TrueStrings {
		if strings.EqualFold(v, s) {
			return true
		}
	}

	return false
}

// IsFalseString is
func IsFalseString(s string) bool {
	for _, v := range FalseStrings {
		if strings.EqualFold(v, s) {
			return true
		}
	}

	return false
}

// CheckStringEmpty :
func CheckStringEmpty(s string) bool {
	return len(s) > 0
}

// JoinStrings :
func JoinStrings(ss []string) string {
	var buff bytes.Buffer

	for i := 0; i < len(ss); i++ {
		buff.WriteString(ss[i])
	}

	return buff.String()
}

// StrSliceToIntSlice :
func StrSliceToIntSlice(ss []string) []int {
	ret := make([]int, len(ss))

	for i, v := range ss {
		if n, err := strconv.Atoi(v); err == nil {
			ret[i] = n
		}
	}

	return ret
}

// IntSliceToStrSlice is
func IntSliceToStrSlice(ns []int) []string {
	strs := make([]string, len(ns))
	for i, v := range ns {
		strs[i] = strconv.Itoa(v)
	}
	return strs
}

// StrSliceToInt32Slice :
func StrSliceToInt32Slice(ss []string) []int32 {
	ret := make([]int32, len(ss))
	for i, v := range ss {
		if n, err := strconv.ParseInt(v, 10, 32); err == nil {
			ret[i] = int32(n)
		}
	}

	return ret
}

// Int32SliceToStrSlice is
func Int32SliceToStrSlice(ns []int32) []string {
	strs := make([]string, len(ns))
	for i, v := range ns {
		strs[i] = strconv.FormatInt(int64(v), 10)
	}
	return strs
}

// StrSliceToInt64Slice :
func StrSliceToInt64Slice(ss []string) []int64 {
	ret := make([]int64, len(ss))
	for i, v := range ss {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			ret[i] = n
		}
	}

	return ret
}

// Int64SliceToStrSlice is
func Int64SliceToStrSlice(ns []int64) []string {
	strs := make([]string, len(ns))
	for i, v := range ns {
		strs[i] = strconv.FormatInt(v, 10)
	}
	return strs
}

// StrSliceToFloat32Slice :
func StrSliceToFloat32Slice(ss []string) []float32 {
	ret := make([]float32, len(ss))
	for i, v := range ss {
		if n, err := strconv.ParseFloat(v, 32); err == nil {
			ret[i] = float32(n)
		}
	}

	return ret
}

// Float32SliceToStrSlice is
func Float32SliceToStrSlice(fs []float32) []string {
	strs := make([]string, len(fs))
	for i, v := range fs {
		strs[i] = strconv.FormatFloat(float64(v), 'f', 30, 32)
	}
	return strs
}

// StrSliceToFloat64Slice :
func StrSliceToFloat64Slice(ss []string) []float64 {
	ret := make([]float64, len(ss))
	for i, v := range ss {
		if n, err := strconv.ParseFloat(v, 64); err == nil {
			ret[i] = n
		}
	}

	return ret
}

// Float64SliceToStrSlice is
func Float64SliceToStrSlice(fs []float64) []string {
	strs := make([]string, len(fs))
	for i, v := range fs {
		strs[i] = strconv.FormatFloat(v, 'f', 62, 64)
	}
	return strs
}

// StrSliceToBoolSlice :
func StrSliceToBoolSlice(ss []string) []bool {
	ret := make([]bool, len(ss))
	for i, v := range ss {
		ret[i] = IsTrueString(v)
	}
	return ret
}

// BoolSliceToStrSlice is
func BoolSliceToStrSlice(bs []bool) []string {
	strs := make([]string, len(bs))
	for i, v := range bs {
		strs[i] = GetBoolString(v)
	}
	return strs
}
