package goblazer

import "bytes"

// CheckStringEmpty :
func CheckStringEmpty(s *string) bool {
	if s == nil || len(*s) <= 0 {
		return false
	}

	return true
}

// JoinStrings :
func JoinStrings(ss []string) string {
	var buff bytes.Buffer

	for i := 0; i < len(ss); i++ {
		buff.WriteString(ss[i])
	}

	return buff.String()
}
